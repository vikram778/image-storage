package producer

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

// KafkaConfig for getting pusher feedbacks
type KafkaConfig struct {
	Brokers string
	Topics  []string
	Config  *sarama.Config
	Group   string
	Logger  *logrus.Logger
}

func NewKafkaProducer(brokers string, topics []string, config *sarama.Config, group string, logger *logrus.Logger) *KafkaConfig {
	return &KafkaConfig{
		Brokers: brokers,
		Topics:  topics,
		Config:  config,
		Group:   group,
		Logger:  logger,
	}
}

// FormatConfiguration ...
func FormatConfiguration() (q *KafkaConfig, err error) {

	q = &KafkaConfig{}

	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)

	log.Println("Formatting sarama configuration")

	q.Logger = logrus.New()
	q.Logger.SetFormatter(&logrus.JSONFormatter{})
	q.Logger.SetOutput(os.Stdout)

	l := q.Logger.WithFields(logrus.Fields{
		"method": "kafka.FormatConfiguration",
	})

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */

	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.ClientID = os.Getenv("KAFKA_CONSUMER_GROUP_ID")
	config.Metadata.RefreshFrequency = 2 * time.Minute
	config.Metadata.Retry.Max = 3
	config.Net.KeepAlive = 1 * time.Minute
	config.Net.DialTimeout = 30 * time.Second
	config.Net.ReadTimeout = 180 * time.Second
	config.Net.WriteTimeout = 30 * time.Second

	config.Producer.Retry.Max = 3
	config.Producer.Retry.Backoff = 30 * time.Second

	assignor := "range"
	if os.Getenv("ASSIGNOR") != "" {
		assignor = os.Getenv("ASSIGNOR")
	}
	switch assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		l.WithError(err).Error("[Error] Unrecognized consumer group partition assignor")
		return q, err
	}

	q.Config = config

	// get broker
	clusterEnv := os.Getenv("CLUSTER_ENV")
	if clusterEnv == "LOCAL" {
		q.Brokers = os.Getenv("LOCAL_CLUSTER")
	}

	q.Topics = strings.Split(os.Getenv("KAFKA_TOPICS"), ",")
	q.Group = os.Getenv("KAFKA_CONSUMER_GROUP_ID")
	return q, nil
}

var Jobs = make(chan map[string]interface{})

func KafkaDispatcher(k *KafkaConfig) {

	wg := &sync.WaitGroup{}

	wg.Add(10)
	go func(bool) {
		defer wg.Done()
		for {
			select {
			case j := <-Jobs:
				k.ProduceLoop(j)

			}
		}
	}(false)
	wg.Wait()
}

// SubscribeLoop ...
func (q *KafkaConfig) ProduceLoop(msg map[string]interface{}) error {
	l := q.Logger.WithFields(logrus.Fields{
		"method":  "kafka.ProduceLoop",
		"brokers": q.Brokers,
		"topics":  q.Topics,
		"group":   q.Group,
		"message": msg,
	})

	// Start with a client
	client, err := sarama.NewClient(strings.Split(q.Brokers, ","), q.Config)
	if err != nil {
		l.WithError(err).Error("[Error] Starting client")
		return err
	}
	defer func() {
		if err = client.Close(); err != nil {
			l.WithError(err).Error("[Error] Closing client")
		}
	}()

	// Start a new consumer group
	producer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		l.WithError(err).Error("[Error] Starting producer")
		return err
	}
	defer func() {
		if err = producer.Close(); err != nil {
			l.WithError(err).Error("[Error] Closing producer")
		}
	}()

	// Track errors
	go func() {
		for err := range producer.Errors() {
			l.WithError(err).Error("Errors")
		}
	}()

	msgbyt, _ := json.Marshal(msg)
	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	var enqueued, errors int
	pmsg := &sarama.ProducerMessage{
		Topic: msg["topic"].(string),
		Key:   nil,
		Value: sarama.StringEncoder(string(msgbyt)),
	}
	producer.Input() <- pmsg

	l.Infof("Enqueued: %d; errors: %d\n", enqueued, errors)
	return err
}

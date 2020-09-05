package logs

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"image-storage/filesystem"

	"github.com/willf/pad"
)

const (
	// NewInstanceMsg sets the message to indicate the start of the log
	NewInstanceMsg = "START"
	// EndInstanceMsg sets the message to indicate the end of the log
	EndInstanceMsg = "END"
	// LogLevelDebug defines a normal debug log
	LogLevelDebug = "DEBUG"
	// LogLevelPanic defines a panic log
	LogLevelPanic = "PANIC"
	// LogLevelFatal defines a fatal log
	LogLevelFatal = "FATAL"
	// DateFormat defines the log date format
	DateFormat = time.RFC3339
)

var (
	osExit = os.Exit
)

// Log represents information about a rest server log.
type Log struct {
	entries   []Entry
	folder    string
	indentier string
	showCount bool
	startTime int64
}

// Entry represents information about a rest server log entry.
type Entry struct {
	Level   string
	Message string
	Time    time.Time
}

func (l Log) getDate(t time.Time) string {
	return t.Format(DateFormat)
}

// New creates new instance of Log
func New() *Log {
	var log Log

	log.folder = "/tmp/log"

	if os.Getenv("LOG_FOLDER") != "" {
		log.folder = os.Getenv("LOG_FOLDER")
	}
	if exist, _ := filesystem.Exist(log.folder); !exist {
		filesystem.Mkdir(log.folder)
	}

	log.entries = make([]Entry, 1)
	log.entries[0] = Entry{
		Message: NewInstanceMsg,
		Time:    time.Now(),
	}

	log.startTime = log.TimeMs()

	return &log
}

// ShowCount log
func (l *Log) ShowCount(show bool) {
	l.showCount = show
}

// GetIdentify log
func (l *Log) GetIdentify() string {
	return l.indentier
}

// GetCount log
func (l *Log) GetCount() bool {
	return l.showCount
}

// SetIdentify log
func (l *Log) SetIdentify(tag string) {
	l.indentier = pad.Right(tag, 15, " ")
}

func (l *Log) addEntry(level string, v ...interface{}) {
	l.entries = append(
		l.entries,
		Entry{
			Level:   level,
			Message: fmt.Sprint(v...),
			Time:    time.Now(),
		},
	)
}

// Entries returns all the entries
func (l *Log) Entries() []Entry {
	return l.entries
}

// Print a regular log
func (l *Log) Print(v ...interface{}) {
	l.addEntry(LogLevelDebug, v...)
}

// Panic then throws a panic with the same message afterwards
func (l *Log) Panic(v ...interface{}) {
	l.addEntry(LogLevelPanic, v...)
	panic(fmt.Sprint(v...))
}

// ThrowFatalTest allows Fatal to be testable
func (l *Log) ThrowFatalTest(msg string) {
	defer func() { osExit = os.Exit }()
	osExit = func(int) {}
	l.Fatal(msg)
}

// Fatal is equivalent to Print() and followed by a call to os.Exit(1)
func (l *Log) Fatal(v ...interface{}) {
	l.addEntry(LogLevelFatal, v...)
	l.Dump()
	osExit(1)
}

// LastEntry returns the last inserted log
func (l *Log) LastEntry() Entry {
	return l.entries[len(l.entries)-1]
}

// Count returns number of inserted logs
func (l *Log) Count() int {
	return len(l.entries)
}

// TimeMs ...
func (l *Log) TimeMs() int64 {
	return time.Now().UnixNano() / 1000000
}

// Dump will print all the messages to the io.
func (l *Log) Dump() {
	var (
		line, format    string
		filename, lines string
		params          []interface{}
	)

	l.Print("Elapse ", l.TimeMs()-l.startTime, " ms")
	l.addEntry("", EndInstanceMsg+"\n")

	len := len(l.entries)
	for i := 0; i < len; i++ {
		format = "%s\t%s"
		params = []interface{}{
			l.getDate(l.entries[i].Time),
			l.entries[i].Level,
		}

		if l.indentier != "" {
			params = append(params, l.indentier)
			format = format + "\t%s"
		}

		if l.showCount {
			params = append(params, pad.Left(strconv.Itoa(i), 3, "0"))
			format = format + "  %s"
		}

		params = append(params, l.entries[i].Message)

		format = format + "  %s\n"
		line = fmt.Sprintf(format, params...)
		lines = lines + line
		fmt.Print(line)
	}

	// create folder for program if not exist
	p := l.folder + string(os.PathSeparator)
	if exist, _ := filesystem.Exist(p); !exist {
		filesystem.Mkdir(p)
	}

	filename = fmt.Sprintf("%s/%d.%s.log", p, time.Now().UnixNano())
	file, _ := os.Create(filename)
	defer file.Close()
	l.entries = make([]Entry, 0)
	file.WriteString(lines)

}

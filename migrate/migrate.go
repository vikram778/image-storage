package migrate

import (
	"errors"
	"strconv"
	"strings"
	"image-storage/logs"

	"bitbucket.org/liamstask/goose/lib/goose"
)

// Migration ...
type Migration struct {
	ID  int64
	SQL string
}

// Migrations ...
type Migrations []Migration

// Process ...
func Process(conf *goose.DBConf, migration Migrations) (err error) {
	var (
		log = logs.New()
		ms  = Migrations{}
	)

	defer func() {
		log.Dump()
	}()

	log.Print("Goose migration")

	log.Print("Connect to DB")
	db, err := goose.OpenDBFromDBConf(conf)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Print("Get current version")
	current, err := goose.EnsureDBVersion(conf, db)
	if err != nil {
		return err
	}

	log.Print("Get target version")
	target, _ := strconv.ParseInt("-1", 10, 64)
	for _, m := range migration {
		if m.ID > target {
			target = m.ID
		}
	}

	if target == -1 {
		err = errors.New("No valid version found")
		return
	}

	if target == current {
		log.Print("Nothing to migrate!")
		err = nil
		return
	}

	log.Print("Current Version ID : ", current)
	log.Print("Target Version ID : ", target)

	for _, m := range migration {
		if versionFilter(m.ID, current, target) {
			ms = append(ms, Migration{ID: m.ID, SQL: m.SQL})
		}
	}

	for _, m := range ms {
		txn, err := db.Begin()
		if err != nil {
			return err
		}

		for _, query := range strings.SplitAfter(m.SQL, ";") {
			if query == "" {
				continue
			}

			if _, err = txn.Exec(query); err != nil {
				txn.Rollback()

				log.Print("ERROR : Version ID ", m.ID)
				log.Print("QUERY : ", query)

				return err
			}
		}

		if err = goose.FinalizeMigration(conf, txn, true, m.ID); err != nil {
			return err
		}

		log.Print("OK : Version ID ", m.ID)
	}

	return

}

// versionFilter ...
func versionFilter(v, current, target int64) bool {

	if target > current {
		return v > current && v <= target
	}

	if target < current {
		return v <= current && v > target
	}

	return false
}

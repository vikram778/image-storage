package db

import "github.com/jmoiron/sqlx"

// DB represents database config
type DB struct {
	Connection *sqlx.DB
	Status     string
	Driver     string
	Open       string
}

// Databases maps all databases used by this application
type Databases map[string]*DB

// New creates a new (DB)Database object
func New(driver string, open string) (*DB, error) {
	db := &DB{
		Driver: driver,
		Open:   open,
	}

	if _, err := db.Connect(); err != nil {
		return nil, err
	}

	return db, nil
}

// Connect sets up a connection using the current credentials
func (d *DB) Connect() (*sqlx.DB, error) {
	db, err := sqlx.Connect(d.Driver, d.Open)
	d.Connection = db
	d.Status = "open"
	return db, err
}

// Close sets up a connection using the current credentials
func (d *DB) Close() {
	if d.Status == "open" {
		d.Connection.Close()
		d.Status = "close"
	}
}

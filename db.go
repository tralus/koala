package koala

import (
	"github.com/jmoiron/sqlx"
)

type DB struct {
	DBConfig DBConfig
}

// It creates a new DB instance
func NewDB(config DBConfig) *DB {
	return &DB{config}
}

// It connects to database
func (db *DB) Connect() (*sqlx.DB, error) {
	config := db.DBConfig
	
	sqlxdb, err := sqlx.Connect(config.Driver, config.Datasource)
	
	db.Configure(sqlxdb)
	
	return sqlxdb, err
}

// It sets default values for the connection
func (db *DB) Configure(sqlxDB *sqlx.DB) {
	sqlxDB.SetMaxOpenConns(db.DBConfig.MaxOpenConns)
	sqlxDB.SetMaxIdleConns(0)
}
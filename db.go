package koala

import (
	"github.com/jmoiron/sqlx"
)

type DBAdapter interface {
	Connect() (*sqlx.DB, error)
}

type DB struct {
	DBConfig DBConfig
}

func NewDB(config DBConfig) DB {
	return DB{config}
}

func (db *DB) Connect(f func (DBConfig) (*sqlx.DB, error)) (*sqlx.DB, error) {
	sqlxdb, err := f(db.DBConfig)
	
	db.Configure(sqlxdb)
	
	return sqlxdb, err
}

func (db *DB) Configure(sqlxDB *sqlx.DB) {
	sqlxDB.SetMaxOpenConns(db.DBConfig.MaxOpenConns)
	sqlxDB.SetMaxIdleConns(0)
}

func SqlxDBConnect(config DBConfig) (*sqlx.DB, error) {
	return sqlx.Connect(config.Driver, config.Datasource)
}
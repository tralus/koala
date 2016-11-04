package koala

import (
	"github.com/jmoiron/sqlx"
)

// DBConfig represents the database config
type DBConfig struct {
	Driver       string
	Datasource   string
	MaxOpenConns int
	MaxIdleConns int
}

// NewDBConfig creates an instance of DatabaseConfig
func NewDBConfig(dr string, dt string) DBConfig {
	return DBConfig{Driver: dr, Datasource: dt}
}

// ConnectDB connects for the database
func ConnectDB(c DBConfig) (*sqlx.DB, error) {
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 20
	}

	sqlxdb, err := sqlx.Connect(c.Driver, c.Datasource)

	sqlxdb.SetMaxOpenConns(c.MaxOpenConns)
	sqlxdb.SetMaxIdleConns(c.MaxIdleConns)

	return sqlxdb, err
}

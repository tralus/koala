package sqlxdb

import (
	"github.com/jmoiron/sqlx"
	"github.com/tralus/koala/errors"
)

// Config represents the database config
type Config struct {
	Driver string
	DSN    string

	MaxOpenConns int
	MaxIdleConns int
}

// NewConfig creates an instance of Config
func NewConfig(dr string, dt string) Config {
	return Config{Driver: dr, DSN: dt}
}

// Connect connects for the database
func Connect(c Config) (*sqlx.DB, error) {
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 20
	}

	db, err := sqlx.Connect(c.Driver, c.DSN)

	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetMaxIdleConns(c.MaxIdleConns)

	return db, err
}

// DatabaseError error type for database error
type DatabaseError struct {
	errors.BaseError
}

// NewDatabaseError creates a new DataBaseError instance
func NewDatabaseError(err error) error {
	return DatabaseError{errors.NewBaseError(err)}
}

// IsDatabaseError verifies if error is an DataBaseError
func IsDatabaseError(err error) bool {
	_, ok := errors.Cause(err).(DatabaseError)
	return ok
}

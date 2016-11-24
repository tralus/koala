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

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetMaxIdleConns(c.MaxIdleConns)

	return db, err
}

// DatabaseError error type for database error
type DatabaseError struct {
	Msg   string
	Stack string
}

// GetStack gets stack trace error
func (e DatabaseError) GetStack() string {
	return e.Stack
}

// Built-in interface
func (e DatabaseError) Error() string {
	return e.Msg
}

// NewDatabaseError creates a new DataBaseError instance
func NewDatabaseError(m string) error {
	s, _ := errors.StackTrace()
	return DatabaseError{m, s}
}

// IsDatabaseError verifies if error is an DataBaseError
func IsDatabaseError(e error) bool {
	_, ok := e.(DatabaseError)
	return ok
}

package sqlxtpl

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/tralus/koala/db"
	"github.com/tralus/koala/errors"

	"gopkg.in/guregu/null.v3"
)

// EmptyResultDataError error type for an empty result database
type EmptyResultDataError struct {
	Msg   string
	Stack string
}

// GetStack gets stack trace error
func (e EmptyResultDataError) GetStack() string {
	return e.Stack
}

// Built-in interface
func (e EmptyResultDataError) Error() string {
	return e.Msg
}

// NewEmptyResultDataError creates an EmptyResultDataError instance
func NewEmptyResultDataError(m string) error {
	s, _ := errors.StackTrace()
	return EmptyResultDataError{m, s}
}

// IsEmptyResultDataError verifies if error is an EmptyResultDataError
func IsEmptyResultDataError(e error) bool {
	_, ok := e.(EmptyResultDataError)
	return ok
}

// TxSQLSetter interface for a struct that supports transaction
type TxSQLSetter interface {
	SetTx(tx *sqlx.Tx)
}

// TransactedSQL should be embedded on sql repositories
type TransactedSQL struct {
	tx *sqlx.Tx
}

// SetTx sets the sql transaction
func (t *TransactedSQL) SetTx(tx *sqlx.Tx) {
	t.tx = tx
}

// Tx gets the sql transaction
func (t *TransactedSQL) Tx() *sqlx.Tx {
	return t.tx
}

// SqlxTpl is a template for database queries
type SqlxTpl struct {
	TransactedSQL
	DB *sqlx.DB
}

// UnsafeSelect executes unsafe select on the database connection
func (s SqlxTpl) UnsafeSelect(dest interface{}, query string, args ...interface{}) error {
	err := s.DB.Select(dest, query, args)

	if err != nil {
		if err == sql.ErrNoRows {
			return NewEmptyResultDataError("Empty result data.")
		}
		return db.NewDatabaseError("Database Error - UnsafeSelect Tx: " + err.Error())
	}

	return nil
}

// NewSqlxTpl creates a SqlxTpl instance
func NewSqlxTpl(db *sqlx.DB) *SqlxTpl {
	return &SqlxTpl{TransactedSQL{}, db}
}

// TxDo executes a callback function with a shared transaction
func (s *SqlxTpl) TxDo(do func(tx *sqlx.Tx) error) error {
	tx, err := Begin(s.DB)

	if err != nil {
		return db.NewDatabaseError(
			"Database Error - Begin Tx: " + err.Error())
	}

	err = do(tx)

	if err != nil {
		errback := Rollback(tx)

		if errback != nil {
			return db.NewDatabaseError(
				"Database Error - Roolback: " + errback.Error())
		}

		return db.NewDatabaseError(
			"Database Error - TxDo Callback: " + err.Error())
	}

	err = Commit(tx)

	if err != nil {
		return err
	}

	return nil
}

// Begin creates a sqlx transaction
func Begin(d *sqlx.DB) (*sqlx.Tx, error) {
	tx, err := d.Beginx()

	if err != nil {
		return nil, db.NewDatabaseError(
			"Database Error - Begin Tx: " + err.Error())
	}

	return tx, nil
}

// Rollback undoes queries of the transaction
func Rollback(tx *sqlx.Tx) error {
	if tx == nil {
		return nil
	}

	err := tx.Rollback()

	if err != nil {
		return db.NewDatabaseError(
			"Database Error - Can`t Rollback: " + err.Error())
	}

	return nil
}

// Commit applies queries of the transaction
func Commit(tx *sqlx.Tx) error {
	if tx == nil {
		return nil
	}

	err := tx.Commit()

	if err != nil {
		return db.NewDatabaseError(
			"Database Error - Can`t Commit: " + err.Error())
	}

	return nil
}

// NamedExec executes a query
// If a transaction is setted, the query runs over it
func (s *SqlxTpl) NamedExec(query string, arg interface{}) (sql.Result, error) {
	var sqlResult sql.Result
	var err error

	t := s.Tx()

	if t != nil {
		sqlResult, err = t.NamedExec(query, arg)
	} else {
		sqlResult, err = s.DB.NamedExec(query, arg)
	}

	if err != nil {
		return nil, db.NewDatabaseError(
			"Database Error - NamedExec: " + err.Error())
	}

	return sqlResult, nil
}

// TxNamedExec executes the query with a transaction
func (s *SqlxTpl) TxNamedExec(tx *sqlx.Tx, query string, arg interface{}) (sql.Result, error) {
	if tx == nil {
		return nil, db.NewDatabaseError(
			"Database Error - Tx is not a valid instance.")
	}

	sqlResult, err := tx.NamedExec(query, arg)

	if err != nil {
		return nil, db.NewDatabaseError(
			"Database Error - NamedExec: " + err.Error())
	}

	return sqlResult, nil
}

// NullInt returns a invalid null.Int when i is zero
func NullInt(i int) null.Int {
	return null.NewInt(int64(i), i != 0)
}

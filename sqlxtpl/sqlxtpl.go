package sqlxtpl

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/tralus/koala/errors"
	"github.com/tralus/koala/sqlxdb"

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

// UnsafeQueryx executes unsafe query on the database connection
func (s SqlxTpl) UnsafeQueryx(query string) (*sqlx.Rows, error) {
	rows, err := s.DB.Unsafe().Queryx(query)

	if err != nil {
		if err == sql.ErrNoRows {
			return rows, NewEmptyResultDataError("Empty result data.")
		}
		return rows, sqlxdb.NewDatabaseError(err.Error())
	}

	return rows, nil
}

// UnsafeSelect executes unsafe select on the database connection
func (s SqlxTpl) UnsafeSelect(dest interface{}, query string, args ...interface{}) error {
	err := s.DB.Unsafe().Select(dest, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return NewEmptyResultDataError("Empty result data.")
		}
		return sqlxdb.NewDatabaseError(err.Error())
	}

	return nil
}

// Select executes safe select on the database connection
func (s SqlxTpl) Select(dest interface{}, query string, args ...interface{}) error {
	err := s.DB.Select(dest, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return NewEmptyResultDataError("Empty result data.")
		}
		return sqlxdb.NewDatabaseError(err.Error())
	}

	return nil
}

// UnsafeGet executes unsafe get on the database connection
func (s SqlxTpl) UnsafeGet(dest interface{}, query string, args ...interface{}) error {
	err := s.DB.Unsafe().Get(dest, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return NewEmptyResultDataError("Empty result data.")
		}
		return sqlxdb.NewDatabaseError(err.Error())
	}

	return nil
}

// Get executes safe get on the database connection
func (s SqlxTpl) Get(dest interface{}, query string, args ...interface{}) error {
	err := s.DB.Get(dest, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return NewEmptyResultDataError("Empty result data.")
		}
		return sqlxdb.NewDatabaseError(err.Error())
	}

	return nil
}

// NewSqlxTpl creates a SqlxTpl instance
func NewSqlxTpl(db *sqlx.DB) SqlxTpl {
	return SqlxTpl{TransactedSQL{}, db}
}

// TxDo executes a callback function with a shared transaction
func (s *SqlxTpl) TxDo(do func(tx *sqlx.Tx) error) error {
	tx, err := Begin(s.DB)

	if err != nil {
		return sqlxdb.NewDatabaseError(err.Error())
	}

	err = do(tx)

	if err != nil {
		errback := Rollback(tx)

		if errback != nil {
			return sqlxdb.NewDatabaseError(errback.Error())
		}

		return sqlxdb.NewDatabaseError(err.Error())
	}

	err = Commit(tx)

	if err != nil {
		return err
	}

	return nil
}

// Begin creates a sqlx transaction
func Begin(db *sqlx.DB) (*sqlx.Tx, error) {
	tx, err := db.Beginx()

	if err != nil {
		return nil, sqlxdb.NewDatabaseError(err.Error())
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
		return sqlxdb.NewDatabaseError(err.Error())
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
		return sqlxdb.NewDatabaseError(err.Error())
	}

	return nil
}

// NamedExec executes a query
// If a transaction is setted, the query runs over it
func (s SqlxTpl) NamedExec(query string, arg interface{}) (sql.Result, error) {
	var sqlResult sql.Result
	var err error

	t := s.Tx()

	if t != nil {
		sqlResult, err = t.NamedExec(query, arg)
	} else {
		sqlResult, err = s.DB.NamedExec(query, arg)
	}

	if err != nil {
		return nil, sqlxdb.NewDatabaseError(err.Error())
	}

	return sqlResult, nil
}

// TxNamedExec executes the query with a transaction
func (s SqlxTpl) TxNamedExec(tx *sqlx.Tx, query string, arg interface{}) (sql.Result, error) {
	if tx == nil {
		return nil, sqlxdb.NewDatabaseError("Tx is not a valid instance.")
	}

	sqlResult, err := tx.NamedExec(query, arg)

	if err != nil {
		return nil, sqlxdb.NewDatabaseError(err.Error())
	}

	return sqlResult, nil
}

// NullInt returns a invalid null.Int when i is zero
func NullInt(i int) null.Int {
	return null.NewInt(int64(i), i != 0)
}

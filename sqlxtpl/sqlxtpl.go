package sqlxtpl

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/tralus/koala/errors"
	"github.com/tralus/koala/sqlxdb"

	"gopkg.in/guregu/null.v3"
)

const emptyResultData = "Empty Result Data"

const dbError = "Database Error"

// EmptyResultDataError error type for an empty result database
type EmptyResultDataError struct {
	errors.BaseError
}

// NewEmptyResultDataError creates an EmptyResultDataError instance
func NewEmptyResultDataError(err error) error {
	return EmptyResultDataError{errors.NewBaseError(err)}
}

// IsEmptyResultDataError verifies if error is an EmptyResultDataError
func IsEmptyResultDataError(err error) bool {
	_, ok := errors.Cause(err).(EmptyResultDataError)
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

// ParseRows is used as a callback function to parse query result
type ParseRows func(r *sqlx.Rows) error

// This function is used like an unique point to process rows.
// The parse function passes the next row for the caller to work.
func processRows(rows *sqlx.Rows, err error, parse ParseRows) error {
	if err != nil {
		if err == sql.ErrNoRows {
			return NewEmptyResultDataError(
				errors.Wrap(err, "Empty result data."),
			)
		}

		return sqlxdb.NewDatabaseError(errors.Wrap(err, dbError))
	}

	for rows.Next() {
		if err = parse(rows); err != nil {
			return err
		}
	}

	return nil
}

// NamedQuery executes a safe named query
func (s SqlxTpl) NamedQuery(query string, arg interface{}, parse ParseRows) error {
	rows, err := s.DB.NamedQuery(query, arg)
	return processRows(rows, err, parse)
}

// UnsafeNamedQuery executes an unsafe named query
func (s SqlxTpl) UnsafeNamedQuery(query string, arg interface{}, parse ParseRows) error {
	rows, err := s.DB.Unsafe().NamedQuery(query, arg)
	return processRows(rows, err, parse)
}

// Queryx executes a safe query
func (s SqlxTpl) Queryx(query string, args []interface{}, parse ParseRows) error {
	rows, err := s.DB.Queryx(query, args...)
	return processRows(rows, err, parse)
}

// UnsafeQueryx executes an unsafe query
func (s SqlxTpl) UnsafeQueryx(query string, parse ParseRows, args ...interface{}) error {
	rows, err := s.DB.Unsafe().Queryx(query, args...)
	return processRows(rows, err, parse)
}

// UnsafeSelect executes an unsafe select
func (s SqlxTpl) UnsafeSelect(dest interface{}, query string, args ...interface{}) error {
	err := s.DB.Unsafe().Select(dest, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return NewEmptyResultDataError(
				errors.Wrap(err, emptyResultData))
		}
		return sqlxdb.NewDatabaseError(errors.Wrap(err, dbError))
	}

	return nil
}

// Select executes a safe select
func (s SqlxTpl) Select(dest interface{}, query string, args ...interface{}) error {
	err := s.DB.Select(dest, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return NewEmptyResultDataError(
				errors.Wrap(err, emptyResultData))
		}
		return sqlxdb.NewDatabaseError(errors.Wrap(err, dbError))
	}

	return nil
}

// UnsafeGet executes unsafe get on the database connection
func (s SqlxTpl) UnsafeGet(dest interface{}, query string, args ...interface{}) error {
	err := s.DB.Unsafe().Get(dest, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return NewEmptyResultDataError(
				errors.Wrap(err, emptyResultData))
		}
		return sqlxdb.NewDatabaseError(errors.Wrap(err, dbError))
	}

	return nil
}

// Get executes safe get on the database connection
func (s SqlxTpl) Get(dest interface{}, query string, args ...interface{}) error {
	err := s.DB.Get(dest, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return NewEmptyResultDataError(
				errors.Wrap(err, emptyResultData))
		}
		return sqlxdb.NewDatabaseError(errors.Wrap(err, dbError))
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
		return sqlxdb.NewDatabaseError(
			errors.Wrap(err, dbError),
		)
	}

	err = do(tx)

	if err != nil {
		errback := Rollback(tx)

		if errback != nil {
			return sqlxdb.NewDatabaseError(
				errors.Wrap(errback, dbError),
			)
		}

		return sqlxdb.NewDatabaseError(
			errors.Wrap(err, dbError),
		)
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
		return nil, sqlxdb.NewDatabaseError(errors.Wrap(err, dbError))
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
		return sqlxdb.NewDatabaseError(errors.Wrap(err, dbError))
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
		return sqlxdb.NewDatabaseError(errors.Wrap(err, dbError))
	}

	return nil
}

// NamedExec executes a namedexec query
// If a transaction is setted, the query runs over it
func (s SqlxTpl) NamedExec(query string, arg interface{}) (result sql.Result, err error) {
	tx := s.Tx()

	if tx != nil {
		result, err = tx.NamedExec(query, arg)
	} else {
		result, err = s.DB.NamedExec(query, arg)
	}

	if err != nil {
		return nil, sqlxdb.NewDatabaseError(errors.Wrap(err, dbError))
	}

	return
}

// TxNamedExec executes the query with a transaction
func (s SqlxTpl) TxNamedExec(tx *sqlx.Tx, query string, arg interface{}) (sql.Result, error) {
	if tx == nil {
		return nil, sqlxdb.NewDatabaseError(
			errors.New("Tx is not a valid instance."))
	}

	sqlResult, err := tx.NamedExec(query, arg)

	if err != nil {
		return nil, sqlxdb.NewDatabaseError(
			errors.Wrap(err, dbError))
	}

	return sqlResult, nil
}

// NullInt returns a invalid null.Int when i is zero
func NullInt(i int) null.Int {
	return null.NewInt(int64(i), i != 0)
}

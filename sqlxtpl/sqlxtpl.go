package sqlxtpl

import (
	"database/sql"
	
	"github.com/jmoiron/sqlx"
	"github.com/tralus/koala/errors"
)

// Generic error type for database aspects
type DatabaseError struct {
	Msg string
	Stack string
}

// RuntimeError interface
func (e DatabaseError) GetStack() string {
	return e.Stack
}

// RuntimeError interface
func (e DatabaseError) Error() string {
	return e.Msg
}

// It creates a new DataBaseError instance
func NewDatabaseError(m string) error {
	s, _ := errors.StackTrace()
	return DatabaseError{m, s} 
}

// It verifies if error is an DataBaseError type
func IsDatabaseError(e error) bool {
	_, ok := e.(DatabaseError)
	return ok
}

// Interface for a struct that supports transaction
type TxSqlSetter interface {
	SetTx(tx *sqlx.Tx)
}

// TransactedSql should be embedded on sql repositories
type TransactedSql struct {
	tx *sqlx.Tx
}

// Set the sql transaction
func (t *TransactedSql) SetTx(tx *sqlx.Tx) {
	t.tx = tx
}

// Get the sql transaction
func (t *TransactedSql) Tx() *sqlx.Tx {
	return t.tx
}

// SqlxTpl is a template for database queries
type SqlxTpl struct {
	TransactedSql
	
	DB *sqlx.DB
}

// It creates a SqlxTpl instance
func NewSqlxTpl(db *sqlx.DB) *SqlxTpl {
	return &SqlxTpl{TransactedSql{}, db}
}

// It executes a callback function with a shared transaction
func (sqlxTpl *SqlxTpl) TxDo(do func(tx *sqlx.Tx) error) error {
	tx, err := Begin(sqlxTpl.DB)
	
	if (err != nil) {
		return NewDatabaseError(
			"Database Error - Begin Tx: " + err.Error())
	}
	
	err = do(tx)
	
	if (err != nil) {
		errback := Rollback(tx)
		
		if (errback != nil) {
			return NewDatabaseError(
				"Database Error - Roolback: " + errback.Error())
		}
		
		return NewDatabaseError(
			"Database Error - TxDo Callback: " + err.Error()) 
	}
	
	err = Commit(tx)
	
	if (err != nil) {
		return err
	}
	
	return nil
}

// It creates a sqlx transaction
func Begin(db *sqlx.DB) (*sqlx.Tx, error) {
	tx, err := db.Beginx()
	
	if (err != nil) {
		return nil, NewDatabaseError(
			"Database Error - Begin Tx: " + err.Error())
	}
	
	return tx, nil
}

// It undoes queries of the transaction
func Rollback(tx *sqlx.Tx) error {
	if (tx == nil) {
		return nil
	}
	
	err := tx.Rollback()
    	
	if (err != nil) {
		return NewDatabaseError(
			"Database Error - Can`t Rollback: " + err.Error())
	}
	
	return nil
}

// It applies queries of the transaction
func Commit(tx *sqlx.Tx) error {
	if (tx == nil) {
		return nil
	}
	
	err := tx.Commit()
	
	if (err != nil) {
		return NewDatabaseError(
			"Database Error - Can`t Commit: " + err.Error())
	}
	
	return nil
}

// It executes the query without a transaction
// If a transaction is setted, the query runs over it
func (sqlxTpl *SqlxTpl) NamedExec(query string, arg interface{}) (sql.Result, error) {
	var sqlResult sql.Result
	var err error
	
	tx := sqlxTpl.Tx()
	
	if tx != nil {
		sqlResult, err = tx.NamedExec(query, arg)
	} else {
		sqlResult, err = sqlxTpl.DB.NamedExec(query, arg)
	}
	
   	if (err != nil) {
		return nil, NewDatabaseError(
    		"Database Error - NamedExec: " + err.Error())
	}
   	
   	return sqlResult, nil
}

// It executes the query with a transaction
func (sqlxTpl *SqlxTpl) TxNamedExec(tx *sqlx.Tx, query string, arg interface{}) (sql.Result, error) {
	if (tx == nil) {
		return nil, NewDatabaseError(
    		"Database Error - Tx is not a valid instance.")
	}
	
	sqlResult, err := tx.NamedExec(query, arg)
	
	if (err != nil) {
		return nil, NewDatabaseError(
    		"Database Error - NamedExec: " + err.Error())
	}
	
	return sqlResult, nil
}
package sqlxtpl

import (
	"database/sql"
	
	"github.com/jmoiron/sqlx"
)

// Generic error type for database aspects
type DataBaseError struct {
	Msg string
	Stack string
}

// RuntimeError interface
func (e DataBaseError) GetStack() string {
	return e.Stack
}

// RuntimeError interface
func (e DataBaseError) Error() string {
	return e.Msg
}

// It creates a new DataBaseError instance
func NewDataBaseError(m string) error {
	s, _ := StackTrace()
	return DataBaseError{m, s} 
}

// It verifies if error is an DataBaseError type
func IsDataBaseError(e error) bool {
	_, ok := e.(DataBaseError)
	return ok
}

// SqlxTpl is a template for database queries
type SqlxTpl struct {
	DB *sqlx.DB
}

// It creates a SqlxTpl instance
func NewSqlxTpl(db *sqlx.DB) *SqlxTpl {
	return &SqlxTpl{db}
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
		errback := Roolback(tx)
		
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
func Roolback(tx *sqlx.Tx) error {
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
func (sqlxTpl *SqlxTpl) NamedExec(query string, arg interface{}) (sql.Result, error) {
	sqlResult, err := sqlxTpl.DB.NamedExec(query, arg)
	
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
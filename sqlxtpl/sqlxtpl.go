package sqlxtpl

import (
	"reflect"
	"database/sql"
	
	"github.com/jmoiron/sqlx"
)

type DataBaseError struct {
	Msg string
}

func (err DataBaseError) Error() string {
	return err.Msg 
}

func NewDatabaseError(msg string) error {
	return DataBaseError{msg}
}

func IsDataBaseError(err error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(DataBaseError{})
}

type SqlxTpl struct {
	DB *sqlx.DB
}

func NewSqlxTpl(db *sqlx.DB) *SqlxTpl {
	return &SqlxTpl{db}
}

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

func Begin(db *sqlx.DB) (*sqlx.Tx, error) {
	tx, err := db.Beginx()
	
	if (err != nil) {
		return nil, NewDatabaseError(
			"Database Error - Begin Tx: " + err.Error())
	}
	
	return tx, nil
}

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

func (sqlxTpl *SqlxTpl) NamedExec(query string, arg interface{}) (sql.Result, error) {
	sqlResult, err := sqlxTpl.DB.NamedExec(query, arg)
	
   	if (err != nil) {
		return nil, NewDatabaseError(
    		"Database Error - NamedExec: " + err.Error())
	}
   	
   	return sqlResult, nil
}

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
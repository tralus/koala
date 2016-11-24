package sqlxtest

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/jmoiron/sqlx"
	"github.com/tralus/koala/migration"
	"github.com/tralus/koala/sqlxdb"
)

const migrationsDir = "migrations"

var errDatabaseNameEmpty = errors.New("The database name is empty in the dsn.")
var errDriverNotSupported = errors.New("Driver not supported.")

// DB connection for tests
var DB *sqlx.DB

// DatabaseSetter defines an interface for database configurators
type DatabaseSetter interface {
	ConnectDatabase() (*sqlx.DB, error)

	CreateDatabase() error

	DropDatabase() error
}

// TestDB represents a test database
type TestDB struct {
	oldDB    *sqlx.DB
	dbConfig sqlxdb.Config
}

// NewTestDB creates a TestDB instance
func NewTestDB(c sqlxdb.Config) (TestDB, error) {
	db, err := sqlxdb.Connect(c)

	if err != nil {
		return TestDB{}, err
	}

	return TestDB{db, c}, nil
}

// Postgres represents a test database for postgres
type Postgres struct {
	TestDB
}

// NewPostgres creates a Postgres instance
func NewPostgres(c sqlxdb.Config) (Postgres, error) {
	db, err := NewTestDB(c)

	if err != nil {
		return Postgres{}, err
	}

	return Postgres{db}, nil
}

// getDatabaseName gets the database name for postgres test database
func (p Postgres) getDatabaseName() (string, error) {
	dsn, err := p.getDSN()

	if err != nil {
		return "", err
	}

	dbName, err := GetDatabaseNameFromDSN(dsn)

	if err != nil {
		return "", err
	}

	return dbName, nil
}

// getDSN gets the dsn for postgres test database
func (p Postgres) getDSN() (string, error) {
	dsn := p.dbConfig.DSN

	u, err := url.Parse(dsn)

	if err != nil {
		return "", err
	}

	n, err := GetDatabaseNameFromDSN(dsn)

	if err != nil {
		return "", err
	}

	u.Path = u.Path[:1] + n + "_test"

	return u.String(), nil
}

// CreateDatabase - Implements DatabaseSetter
func (p Postgres) CreateDatabase() error {
	dbName, err := p.getDatabaseName()

	if err != nil {
		return err
	}

	sql := "CREATE DATABASE " + dbName

	if _, err := p.oldDB.Exec(sql); err != nil {
		return err
	}

	return nil
}

// ConnectDatabase - Implements DatabaseSetter
func (p Postgres) ConnectDatabase() (*sqlx.DB, error) {
	dsn, err := p.getDSN()

	if err != nil {
		return nil, err
	}

	c := sqlxdb.NewConfig(p.dbConfig.Driver, dsn)

	return sqlxdb.Connect(c)
}

// DropDatabase - Implements DatabaseSetter
func (p Postgres) DropDatabase() error {
	n, err := p.getDatabaseName()

	if err != nil {
		return err
	}

	sql := "DROP DATABASE IF EXISTS " + n

	if _, err := p.oldDB.Exec(sql); err != nil {
		return err
	}

	return nil
}

// GetDatabaseSetter gets a database setter for the driver on c (sqlxdb.Config)
func GetDatabaseSetter(c sqlxdb.Config) (DatabaseSetter, error) {
	if c.Driver == "postgres" || c.Driver == "postgresql" {
		return NewPostgres(c)
	}
	return nil, errDriverNotSupported
}

// GetDatabaseNameFromDSN gests the database name from the dsn
func GetDatabaseNameFromDSN(dsn string) (string, error) {
	u, err := url.Parse(dsn)

	if err != nil {
		return "", err
	}

	n := u.Path[1:]

	if n == "" {
		return "", errDatabaseNameEmpty
	}

	return n, nil
}

// Setup represents the test setup
type Setup struct {
	dbConfig sqlxdb.Config
	dbSetter DatabaseSetter
}

// NewSetup creates a Setup instance
func NewSetup(dbConfig sqlxdb.Config) (Setup, error) {
	st, err := GetDatabaseSetter(dbConfig)

	if err != nil {
		return Setup{}, err
	}

	return Setup{dbConfig, st}, nil
}

// Run runs the test setup
func (s Setup) Run() error {
	if err := s.Destroy(); err != nil {
		return err
	}

	fmt.Println("Creating the test database...")

	if err := s.dbSetter.CreateDatabase(); err != nil {
		return err
	}

	DB, err := s.dbSetter.ConnectDatabase()

	if err != nil {
		return err
	}

	source := migration.NewFileMigration(migrationsDir)

	migration := migration.NewMigration(
		DB.DB, // sql.DB from sqlx.DB
		s.dbConfig.Driver,
		source,
	)

	fmt.Println("Applying migrations...")

	i, err := migration.Up()

	if err != nil {
		return err
	}

	fmt.Printf("Number of migrations applied: %d.\n", i)

	return nil
}

// Destroy clears all after the test setup
func (s Setup) Destroy() error {
	fmt.Println("Destroying test database...")

	if err := s.dbSetter.DropDatabase(); err != nil {
		return err
	}

	return nil
}

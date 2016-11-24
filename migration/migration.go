package migration

import (
	"database/sql"

	"github.com/rubenv/sql-migrate"
)

// Migration represents a migration runner
type Migration struct {
	schema string
	table  string

	db         *sql.DB
	dialect    string
	migrations migrate.MigrationSource
}

// NewFileMigration creates a migrate.MigrationSource instance
func NewFileMigration(path string) migrate.MigrationSource {
	return migrate.FileMigrationSource{Dir: path}
}

// NewMigration creates an Migration instance
func NewMigration(db *sql.DB, dialect string, m migrate.MigrationSource) Migration {
	return Migration{db: db, dialect: dialect, migrations: m}
}

// Up runs the migrations to up
func (m Migration) Up() (int, error) {
	m.configureSQLMigrate()
	return migrate.Exec(m.db, m.dialect, m.migrations, migrate.Up)
}

// Down runs the migration to down
func (m Migration) Down() (int, error) {
	m.configureSQLMigrate()
	return migrate.Exec(m.db, m.dialect, m.migrations, migrate.Down)
}

func (m Migration) configureSQLMigrate() {
	schema := "public"
	table := "sql_migrations"

	if m.schema != "" {
		schema = m.schema
	}

	if m.table != "" {
		table = m.table
	}

	migrate.SetSchema(schema)
	migrate.SetTable(table)
}

// SetSchema sets the schema for migrations table
func (m Migration) SetSchema(s string) {
	m.schema = s
}

// SetTable sets the migrations table name
func (m Migration) SetTable(t string) {
	m.table = t
}

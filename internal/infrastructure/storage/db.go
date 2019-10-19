package storage

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	migration "github.com/rubenv/sql-migrate"
)

//DBConfig contains a database configuration.
type DBConfig struct {
	DriverName      string `validate:"required"`
	DSN             string `validate:"required"`
	MaxOpenConns    int    `validate:"required"`
	MaxIdleConns    int    `validate:"required"`
	ConnMaxLifeTime time.Duration
	Migrations      *MigrationsConfig `validate:"required"`
}

// MigrationsConfig contains SQL database migration configurations.
type MigrationsConfig struct {
	Dialect string `yaml:"Dialect" validate:"required"`
	Dir     string `yaml:"Dir" validate:"required"`
}

// NewDB creates a new SQL database.
func NewDB(config *DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open(config.DriverName, config.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifeTime)

	if err := migrate(config.Migrations, db.DB); err != nil {
		return nil, err
	}

	return db, nil
}

func migrate(config *MigrationsConfig, db *sql.DB) error {
	migrations := &migration.FileMigrationSource{Dir: config.Dir}
	_, err := migration.Exec(db, config.Dialect, migrations, migration.Up)
	if err != nil {
		return err
	}
	return nil
}

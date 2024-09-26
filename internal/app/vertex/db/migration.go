package db

import (
	"database/sql"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func MigrateDatabase(path string, databaseName string) *sql.DB {
	zap.L().Info("starting database migrations")

	db, err := sql.Open("sqlite3", path)

	if err != nil {
		zap.L().Fatal("failed to open sqlite database", zap.Error(err))
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})

	if err != nil {
		zap.L().Fatal("failed to create sqlite driver", zap.Error(err))
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Migration files path
		databaseName,        // Database name
		driver,
	)
	if err != nil {
		zap.L().Fatal("failed to initialize migration", zap.Error(err))
	}

	// Apply all migrations
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		zap.L().Fatal("migration failed", zap.Error(err))
	}

	zap.L().Info("database migration applied successfully")

	return db
}

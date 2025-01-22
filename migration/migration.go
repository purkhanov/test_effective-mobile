package migration

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"music/internal/config"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const migrationsPath = "./migration/schemas"

func StartMigrate(cfg config.Config, dbUrl string, logger *slog.Logger) error {
	if err := createDBIfNotExists(cfg, logger); err != nil {
		return err
	}

	m, err := migrate.New("file://"+migrationsPath, dbUrl)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Info("No migrations to apply", slog.String("info", err.Error()))
			return nil
		}

		return err
	}

	logger.Debug("Apply migrations")

	if err := closeMigrationResources(m); err != nil {
		return err
	}

	return nil
}

func createDBIfNotExists(cfg config.Config, logger *slog.Logger) error {
	dbUrl := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/postgres?sslmode=disable",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort,
	)
	logger.Debug("DB url", slog.String("DBurl", dbUrl))

	time.Sleep(10 * time.Second)
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return err
	}
	defer db.Close()

	var exists bool

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1);"

	err = db.QueryRowContext(ctx, query, cfg.DBName).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if !exists {
		query = fmt.Sprintf("CREATE DATABASE %s;", cfg.DBName)

		if _, err := db.ExecContext(ctx, query); err != nil {
			return err
		}

		logger.Debug("Created database")
	} else {
		logger.Debug("Database already exists")
	}

	return nil
}

func closeMigrationResources(m *migrate.Migrate) error {
	sourceErr, dbErr := m.Close()

	if dbErr != nil {
		return dbErr
	}

	if sourceErr != nil {
		return sourceErr
	}

	return nil
}

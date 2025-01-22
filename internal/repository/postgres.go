package repository

import (
	"database/sql"
	"fmt"
	"log"
	"music/internal/config"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgresDB(cfg config.Config) DB {
	db, err := sql.Open("postgres", GetDBUrl(cfg))
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	return &Postgres{db: db}
}

func (p *Postgres) GetDB() *sql.DB {
	return p.db
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

func GetDBUrl(cfg config.Config) string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUsername, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode,
	)
}

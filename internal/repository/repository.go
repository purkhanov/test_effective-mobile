package repository

import (
	"database/sql"
)

type DB interface {
	GetDB() *sql.DB
	Close() error
}

type Repository struct {
	Music
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Music: newMusicPostgres(db),
	}
}

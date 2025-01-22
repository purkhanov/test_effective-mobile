package service

import (
	"log/slog"
	"music/internal/repository"
)

type Service struct {
	Music
}

func NewService(repos *repository.Repository, logger *slog.Logger) *Service {
	return &Service{
		Music: newMusicService(repos.Music, logger),
	}
}

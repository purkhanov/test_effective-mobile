package controller

import (
	"log/slog"
	"music/internal/service"
)

type Controller struct {
	Music
}

func NewController(services *service.Service, logger *slog.Logger) *Controller {
	return &Controller{
		Music: newMusicController(services.Music, logger),
	}
}

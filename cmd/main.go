package main

import (
	"fmt"
	"log/slog"
	"music/internal/app"
	"music/internal/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.LoadConfig()

	logger := config.SetupLogger(cfg.Mode)

	logger.Debug("Loaded configuration")

	application := app.NewHTTPServer(cfg, logger)

	logger.Info(fmt.Sprintf("Starting application on port: %s", cfg.Port))

	go application.Run()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	logger.Info("Stopping application", slog.String("signal", sign.String()))

	if err := application.Shutdown(); err != nil {
		logger.Error("error on shutting down server", slog.String("err", err.Error()))
	} else {
		logger.Error("Application stopped")
	}
}

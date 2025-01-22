package config

import (
	"log"
	"log/slog"
	"os"
)

const logPath = "./logs/out.log"

func SetupLogger(mode string) *slog.Logger {
	var logger *slog.Logger
	var logFile *os.File
	var err error

	if mode == ModeProd {
		logFile, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_RDONLY, 0600)
		if err != nil {
			log.Printf("failed to open log file: %v", err)
			logFile = os.Stdout
		}
	}

	switch mode {
	case ModeLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case modeDev:
		logger = slog.New(
			slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case ModeProd:
		logger = slog.New(
			slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelWarn}),
		)
	}

	return logger
}

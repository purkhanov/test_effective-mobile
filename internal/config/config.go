package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Url string

const (
	ModeLocal = "local"
	modeDev   = "dev"
	ModeProd  = "prod"
)

type Config struct {
	Mode string
	Port string

	DBHost     string
	DBPort     string
	DBPassword string
	DBUsername string
	DBName     string
	DBSSLMode  string
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Can not load env file, loading from environment or using default values")
	}

	var cfg Config

	Url = getEnv("URL", "http://localhost")

	cfg.Mode = getEnv("MODE", ModeLocal)
	cfg.Port = getEnv("PORT", "8000")

	cfg.DBHost = getEnv("DB_HOST", "localhost")
	cfg.DBPort = getEnv("DB_PORT", "5432")
	cfg.DBPassword = getEnv("POSTGRES_PASSWORD", "postgres")
	cfg.DBUsername = getEnv("POSTGRES_USER", "postgres")
	cfg.DBName = getEnv("POSTGRES_DB", "postgres")
	cfg.DBSSLMode = getEnv("DB_SSLMode", "disable")

	return cfg
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return value
}

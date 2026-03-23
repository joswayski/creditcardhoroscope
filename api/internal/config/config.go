package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("No .env file found", "error", err)
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		slog.Warn("Error loading port from .env, using default 8080")
		portString = "8080"
	}

	return Config{
		Port: portString,
	}
}

package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string // Defaults to 8080

	// Required
	AIBaseURL string
	AIAPIKey  string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("No .env file found", "error", err)
	}

	requiredEnvErrors := []string{}

	portString := os.Getenv("PORT")
	if portString == "" {
		slog.Warn("Error loading port from .env, using default 8080")
		portString = "8080"
	}

	aiBaseURL := os.Getenv("AI_BASE_URL")
	if aiBaseURL == "" {
		requiredEnvErrors = append(requiredEnvErrors, "AI_BASE_URL")
	}

	aiAPIKey := os.Getenv("AI_API_KEY")
	if aiAPIKey == "" {
		requiredEnvErrors = append(requiredEnvErrors, "AI_API_KEY")
	}

	if len(requiredEnvErrors) > 0 {
		slog.Error("Missing required environment variables, cannot start!")
		for i, v := range requiredEnvErrors {
			slog.Error(fmt.Sprintf("%d) %s", i+1, v))
		}

		os.Exit(1)
	}

	return Config{
		Port:      portString,
		AIBaseURL: aiBaseURL,
		AIAPIKey:  aiAPIKey,
	}
}

// internal/config/config.go

package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	OpenWeatherMapAPIKey     string
	PostgresConnectionString string
	ServerPort               string
	ResendApiKey             string
}

func LoadConfig() (*Config, error) {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	cfg := &Config{
		PostgresConnectionString: os.Getenv("POSTGRES_CONNECTION_STRING"),
		OpenWeatherMapAPIKey:     os.Getenv("OPENWEATHERMAP_API_KEY"),
		ServerPort:               os.Getenv("PORT"),
		ResendApiKey:             os.Getenv("RESEND_API_KEY"),
	}

	// Add basic validation
	if cfg.PostgresConnectionString == "" {
		return nil, fmt.Errorf("DATABASE_URL not set")
	}

	if cfg.ServerPort == "" {
		cfg.ServerPort = "8080" // Default port if not set
	}

	return cfg, nil
}

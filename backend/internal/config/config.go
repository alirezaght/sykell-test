package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	Environment string
	TemporalHostPort string
	Namespace		string
	LogLevel    string
	LogFormat   string
}

// DefaultTimeout is the default timeout for db operations
const DefaultTimeout = 5 * time.Second

// Load reads configuration from environment variables and returns a Config struct
func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	cfg := &Config{
		Port:        getEnv("PORT", "7070"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:password@localhost/sykell_db?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		Environment: getEnv("ENVIRONMENT", "development"),
		TemporalHostPort: getEnv("TEMPORAL_HOST_PORT", "localhost:7233"),
		Namespace:   getEnv("TEMPORAL_NAMESPACE", "default"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		LogFormat:   getEnv("LOG_FORMAT", "json"),
	}

	return cfg, nil
}

// getEnv retrieves the value of the environment variable named by the key.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
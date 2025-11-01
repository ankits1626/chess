// Package config manages application settings.
package config

import (
	"os"
	"time"
)

// Config holds app settings.
type Config struct {
	// Server settings
	Port        string
	Environment string
	GinMode     string

	// Database settings
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	ShutdownTimeout time.Duration
}

// Load reads config from env with defaults.
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		GinMode:     getEnv("GIN_MODE", "debug"),

		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "chess_coach"),
		DBPassword:      getEnv("DB_PASSWORD", "chess_coach_dev"),
		DBName:          getEnv("DB_NAME", "chess_coach_dev"),
		DBSSLMode:       getEnv("DB_SSLMODE", "disable"),
		ShutdownTimeout: 5 * time.Second,
	}
}

// DSN returns PostgreSQL connection string.
func (c *Config) DSN() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=" + c.DBSSLMode
}

// getEnv retrieves env var or returns default.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

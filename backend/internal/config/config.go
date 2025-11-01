// Package config manages application settings.
package config

import "os"

// Config holds app settings.
type Config struct {
	// Port is the HTTP server port
	Port string
	// Environment is the app environment
	Environment string
	// GinMode is the Gin framework mode
	GinMode string
}

// Load reads config from env with defaults.
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		GinMode:     getEnv("GIN_MODE", "debug"),
	}
}

// getEnv retrieves env var or returns default.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

package config

import (
    "fmt"
    "log/slog"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Port        string
    Environment string
    DatabaseURL string
    RedisURL    string
}

func Load() (*Config, error) {
    // Load .env file (optional - won't fail if missing)
    if err := godotenv.Load(); err != nil {
        slog.Warn("no .env file found, using environment variables")
    }

    cfg := &Config{
        Port:        getEnv("PORT", "8080"),
        Environment: getEnv("APP_ENV", "development"),
        DatabaseURL: getEnv("DATABASE_URL", ""),
        RedisURL:    getEnv("REDIS_URL", ""),
    }

    if err := cfg.Validate(); err != nil {
        return nil, err
    }

    return cfg, nil
}

func (c *Config) Validate() error {
    // Add validation rules as needed
    if c.Port == "" {
        return fmt.Errorf("PORT is required")
    }
    return nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func (c *Config) IsDevelopment() bool {
    return c.Environment == "development"
}

func (c *Config) IsProduction() bool {
    return c.Environment == "production"
}

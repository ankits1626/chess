package server

import (
	"testing"

	"github.com/ankits1626/chess-coach-backend/internal/config"
)

func TestNew_CreatesServer(t *testing.T) {
	t.Skip("Skipping test that requires database connection. Run integration tests with test database.")

	cfg := &config.Config{
		Port:        "8080",
		Environment: "test",
		GinMode:     "test",
	}

	// In a real test environment, you would set up a test database here
	// For example: db := setupTestDB(t)
	// var db *database.DB = nil
	// srv := New(cfg, db)

	_ = cfg // Suppress unused variable warning
}

func TestServer_StartAndShutdown(t *testing.T) {
	t.Skip("Skipping test that requires database connection. Run integration tests with test database.")
}

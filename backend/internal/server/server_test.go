package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ankits1626/chess-coach-backend/internal/config"
)

func TestNew_CreatesServer(t *testing.T) {
	cfg := &config.Config{
		Port:        "8080",
		Environment: "test",
		GinMode:     "test",
	}

	srv := New(cfg)

	if srv == nil {
		t.Error("Expected server to be created, got nil")
	}
}

func TestServer_StartAndShutdown(t *testing.T) {
	cfg := &config.Config{
		Port:        "8081",
		Environment: "test",
		GinMode:     "test",
	}

	srv := New(cfg)

	// Start server in goroutine
	errChan := make(chan error, 1)
	go func() {
		if err := srv.Start(); err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Test server is responding
	resp, err := http.Get("http://localhost:8081/")
	if err != nil {
		t.Fatalf("Server not responding: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		t.Errorf("Shutdown failed: %v", err)
	}

	// Check for startup errors
	select {
	case err := <-errChan:
		t.Errorf("Server error: %v", err)
	default:
		// No error, good
	}
}

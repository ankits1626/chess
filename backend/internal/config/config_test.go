package config

import (
	"os"
	"testing"
)

func TestLoad_WithDefaults(t *testing.T) {
	// Clear environment variables
	os.Unsetenv("PORT")
	os.Unsetenv("ENVIRONMENT")
	os.Unsetenv("GIN_MODE")

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("Expected Port to be '8080', got '%s'", cfg.Port)
	}
	if cfg.Environment != "development" {
		t.Errorf("Expected Environment to be 'development', got '%s'", cfg.Environment)
	}
	if cfg.GinMode != "debug" {
		t.Errorf("Expected GinMode to be 'debug', got '%s'", cfg.GinMode)
	}
}

func TestLoad_WithEnvironmentVariables(t *testing.T) {
	os.Setenv("PORT", "3000")
	os.Setenv("ENVIRONMENT", "production")
	os.Setenv("GIN_MODE", "release")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("ENVIRONMENT")
		os.Unsetenv("GIN_MODE")
	}()

	cfg := Load()

	if cfg.Port != "3000" {
		t.Errorf("Expected Port to be '3000', got '%s'", cfg.Port)
	}
	if cfg.Environment != "production" {
		t.Errorf("Expected Environment to be 'production', got '%s'", cfg.Environment)
	}
	if cfg.GinMode != "release" {
		t.Errorf("Expected GinMode to be 'release', got '%s'", cfg.GinMode)
	}
}

func TestGetEnv_WithValue(t *testing.T) {
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	result := getEnv("TEST_KEY", "default")

	if result != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", result)
	}
}

func TestGetEnv_WithDefault(t *testing.T) {
	os.Unsetenv("NON_EXISTENT_KEY")

	result := getEnv("NON_EXISTENT_KEY", "default")

	if result != "default" {
		t.Errorf("Expected 'default', got '%s'", result)
	}
}

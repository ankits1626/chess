package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestHealthHandler_ReturnsHealthy(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Call handler
	HealthHandler(c)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Parse response
	var response HealthResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check response values
	if response.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response.Status)
	}
	if response.Version != "1.0" {
		t.Errorf("Expected version '1.0', got '%s'", response.Version)
	}
}

func TestHealthHandler_ReturnsJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	HealthHandler(c)

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json; charset=utf-8" {
		t.Errorf("Expected Content-Type 'application/json; charset=utf-8', got '%s'", contentType)
	}
}

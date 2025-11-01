// Package v1 provides API v1 endpoints.
package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthResponse represents health check response.
type HealthResponse struct {
	// Status indicates if service is healthy
	Status string `json:"status"`
	// Version is the API version
	Version string `json:"version"`
}

// HealthHandler handles health check endpoint.
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:  "healthy",
		Version: "1.0",
	})
}

// Package v1 provides API v1 routes registration.
package v1

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all v1 API routes.
func RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", HealthHandler)

		// Future endpoints
		// v1.POST("/analyze", AnalyzeHandler)
		// v1.GET("/games", ListGamesHandler)
	}
}

// Package router configures HTTP routes.
package router

import (
	"net/http"

	"github.com/ankits1626/chess-coach-backend/internal/router/v1"
	"github.com/gin-gonic/gin"
)

// Setup creates configured Gin router.
func Setup() *gin.Engine {
	r := gin.Default()

	// Register v1 API routes
	v1.RegisterRoutes(r)

	// Root endpoint
	r.GET("/", rootHandler)

	return r
}

// rootHandler handles the root endpoint.
func rootHandler(c *gin.Context) {
	c.String(http.StatusOK, "Chess Coach API - Server Running!")
}

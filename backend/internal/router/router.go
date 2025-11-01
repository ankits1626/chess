// Package router configures HTTP routes.
package router

import (
	"net/http"

	v1 "github.com/ankits1626/chess-coach-backend/internal/router/v1"
	"github.com/gin-gonic/gin"

	_ "github.com/ankits1626/chess-coach-backend/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Setup creates configured Gin router.
func Setup() *gin.Engine {
	r := gin.Default()

	// Register v1 API routes
	v1.RegisterRoutes(r)

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Root endpoint
	r.GET("/", rootHandler)

	return r
}

// rootHandler handles the root endpoint.
func rootHandler(c *gin.Context) {
	c.String(http.StatusOK, "Chess Coach API - Server Running!")
}

// Package router provides HTTP routing setup.
package router

import (
	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/ankits1626/chess-coach-backend/docs"
)

// Setup creates and configures the Gin router.
func Setup(db *database.DB) *gin.Engine {
	r := gin.Default()

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Register API v1 routes
	RegisterV1Routes(r, db)

	return r
}

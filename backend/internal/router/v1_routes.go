// Package router provides HTTP routing setup.
package router

import (
	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/handler/v1/game"
	"github.com/ankits1626/chess-coach-backend/internal/handler/v1/health"
	"github.com/ankits1626/chess-coach-backend/internal/handler/v1/move" // ADD
	"github.com/ankits1626/chess-coach-backend/internal/handler/v1/user"
	"github.com/ankits1626/chess-coach-backend/internal/repository"
	"github.com/gin-gonic/gin"
)

// RegisterV1Routes registers API v1 routes.
func RegisterV1Routes(r *gin.Engine, db *database.DB) {
	v1 := r.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", health.Check)

		// User routes
		userRepo := repository.NewUserRepository(db)
		userHandler := user.NewHandler(userRepo)

		// Game routes
		gameRepo := repository.NewGameRepository(db)
		gameHandler := game.NewHandler(gameRepo)

		// Move routes - ADD ENTIRE BLOCK
		moveRepo := repository.NewMoveRepository(db)
		moveHandler := move.NewHandler(moveRepo)

		users := v1.Group("/users")
		{
			users.GET("", userHandler.List)
			users.GET("/:id", userHandler.Get)
			users.POST("", userHandler.Create)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)

			// User's games sub-resource
			users.GET("/:id/games", gameHandler.ListByUser)
		}

		games := v1.Group("/games")
		{
			games.GET("", gameHandler.List)
			games.GET("/:id", gameHandler.Get)
			games.POST("", gameHandler.Create)

			// Game's moves sub-resource - ADD
			games.GET("/:id/moves", moveHandler.ListByGame)
			games.POST("/:id/moves", moveHandler.Create)
		}

		// Move routes - ADD
		moves := v1.Group("/moves")
		{
			moves.GET("/:id", moveHandler.Get)
		}
	}
}

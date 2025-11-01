// @title Chess Coach API
// @version 1.0
// @description API for chess game analysis and coaching
// @host localhost:8080
// @BasePath /api/v1

// Package main bootstraps the Chess Coach API.
package main

import (
	"context"
	"log"

	"github.com/ankits1626/chess-coach-backend/internal/app"
	"github.com/ankits1626/chess-coach-backend/internal/config"
	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	appLogger := logger.NewStdLogger()

	// Connect to database
	ctx := context.Background()
	db, err := database.NewDB(ctx, cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create and run application
	application := app.New(cfg, db, appLogger)
	defer application.Close()

	appLogger.Info("Database connected successfully")

	// Run application (blocks until shutdown)
	if err := application.Run(ctx); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}

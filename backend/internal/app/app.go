// Package app manages application lifecycle.
package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ankits1626/chess-coach-backend/internal/config"
	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/logger"
	"github.com/ankits1626/chess-coach-backend/internal/server"
	"github.com/ankits1626/chess-coach-backend/internal/websocket"
	"github.com/ankits1626/chess-coach-backend/internal/websocket/handlers"
	"github.com/ankits1626/chess-coach-backend/internal/websocket/handlers/player"
)

// App manages application lifecycle.
type App struct {
	config    *config.Config
	db        *database.DB
	server    server.Server
	logger    logger.Logger
	wsHub     *websocket.Hub
	aiService player.AIService // AI service for computer players
}

// New creates new application.
func New(cfg *config.Config, db *database.DB, log logger.Logger) *App {
	// Create chess service
	chessService := handlers.NewChessService()

	// Create AI service for computer players
	aiService, err := player.NewStockfishService()
	if err != nil {
		log.Fatalf("Failed to initialize AI service: %v", err)
	}

	// Create game manager with AI service
	gameManager := handlers.NewGameManager(db, chessService, aiService)

	// Create handler router
	handler := handlers.NewHandlerRouter(gameManager)
	hub := websocket.NewHub(handler)
	return &App{
		config:    cfg,
		db:        db,
		server:    server.New(cfg, db, hub),
		logger:    log,
		wsHub:     hub,
		aiService: aiService,
	}
}

// Run starts application and handles graceful shutdown.
func (a *App) Run(ctx context.Context) error {
	// Start WebSocket hub
	go a.wsHub.Run()
	a.logger.Info("WebSocket Hub started")

	// Start server in goroutine
	go func() {
		if err := a.server.Start(); err != nil {
			a.logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	a.logger.Infof("Server started on port %s", a.config.Port)
	a.logger.Infof("Swagger UI: http://localhost:%s/swagger/index.html", a.config.Port)
	a.logger.Infof("WebSocket endpoint: ws://localhost:%s/ws?user_id=test", a.config.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("Shutting down server...")

	// Shutdown WebSocket hub first
	a.wsHub.Shutdown()
	a.logger.Info("WebSocket Hub stopped")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(ctx, a.config.ShutdownTimeout)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.logger.Errorf("Server forced to shutdown: %v", err)
		return err
	}

	a.logger.Info("Server exited gracefully")
	return nil
}

// Close closes application resources.
func (a *App) Close() {
	if a.aiService != nil {
		a.aiService.Close()
	}
	if a.db != nil {
		a.db.Close()
	}
}

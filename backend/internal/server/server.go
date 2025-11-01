// Package server manages HTTP server lifecycle.
package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ankits1626/chess-coach-backend/internal/config"
	"github.com/ankits1626/chess-coach-backend/internal/router"
	"github.com/gin-gonic/gin"
)

// Server handles HTTP server start/stop.
type Server interface {
	// Start begins listening for HTTP requests.
	Start() error
	// Shutdown gracefully shuts down the server.
	Shutdown(ctx context.Context) error
}

// server is the concrete implementation of Server.
type server struct {
	httpServer *http.Server
	config     *config.Config
}

// New creates a new Server with given config.
func New(cfg *config.Config) Server {
	gin.SetMode(cfg.GinMode)
	r := router.Setup()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	return &server{
		httpServer: httpServer,
		config:     cfg,
	}
}

// Start implements Server.Start.
func (s *server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown implements Server.Shutdown.
func (s *server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

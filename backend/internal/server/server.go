package server

import (
	"log/slog"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/gorilla/websocket"
	"github.com/google/uuid"

	"chess-coach/backend/internal/db"
	"chess-coach/backend/internal/middleware"
	ws "chess-coach/backend/internal/websocket"
)

type Server struct {
	app     *fiber.App
	logger  *slog.Logger
	hub     *ws.Hub
	queries *db.Queries
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")

		// Allow file:// protocol (null or empty origin) for local testing
		if origin == "" || origin == "null" {
			return true
		}

		// Check against allowed origins
		allowedOrigins := middleware.AllowedOrigins()
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				return true
			}
		}
		return false
	},
}

func New(logger *slog.Logger, queries *db.Queries) *Server {
	hub := ws.NewHub()

	app := fiber.New(fiber.Config{
		AppName:      "Chess Coach API",
		ServerHeader: "Chess Coach",
	})

	srv := &Server{
		app:     app,
		logger:  logger,
		hub:     hub,
		queries: queries,
	}

	// Start hub
	go hub.Run()

	return srv
}

func (s *Server) SetupRoutes() {
	// Middleware
	s.app.Use(middleware.CORS())

	// HTTP Routes
	s.app.Get("/health", s.handleHealth)
	s.app.Get("/", s.handleRoot)

	// WebSocket Route
	s.app.Get("/ws", adaptor.HTTPHandlerFunc(s.handleWebSocket))
}

func (s *Server) handleHealth(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "chess-coach-backend",
	})
}

func (s *Server) handleRoot(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Chess Coach Backend API",
		"version": "0.2.0",
	})
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("websocket upgrade failed", "error", err)
		return
	}

	// Create client
	clientID := uuid.New().String()
	s.logger.Info("new websocket connection", "client_id", clientID)

	client := ws.NewClient(s.hub, conn, clientID)
	s.hub.Register <- client

	// Start client goroutines
	go client.WritePump()
	go client.ReadPump()
}

func (s *Server) Start(port string) error {
	s.logger.Info("starting server", "port", port)
	return s.app.Listen(":" + port)
}

func (s *Server) Shutdown() error {
	s.logger.Info("shutting down server")
	return s.app.Shutdown()
}

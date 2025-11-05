package websocket

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		// TODO Phase 5: Implement proper CORS checking
		// For now, allow all origins during development
		return true
	},
}

// Handler handles WebSocket connections.
type Handler struct {
	hub *Hub
}

// NewHandler creates a new WebSocket handler.
func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

// ServeWS handles WebSocket connection upgrade.
func (h *Handler) ServeWS(c *gin.Context) {
	// TODO Phase 5: Get user ID from JWT token
	// For now, use query parameter
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id query parameter required"})
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	// Create new client
	client := NewClient(conn, h.hub, userID)

	// Register client with hub
	h.hub.register <- client

	log.Printf("WebSocket connection established for user: %s (client: %s)", userID, client.ID)

	// Start goroutines for this client
	go client.writePump()
	go client.readPump(context.Background(), h.handleMessage)
}

// handleMessage routes incoming WebSocket messages.
// This is a placeholder - full implementation in Phase 4.
func (h *Handler) handleMessage(ctx context.Context, client *Client, data []byte) {
	// Parse message
	msg, err := FromJSON(data)
	if err != nil {
		log.Printf("Failed to parse message from client %s: %v", client.ID, err)
		errResp := NewErrorResponse("", "Invalid JSON message")
		client.SendMessage(errResp)
		return
	}

	log.Printf("Received message from client %s: type=%s, action=%s, event=%s",
		client.ID, msg.Type, msg.Action, msg.Event)

	// TODO Phase 4: Implement actual message handlers
	// For now, echo back a simple response
	switch msg.Type {
	case TypeRequest:
		// Echo response for testing
		resp := NewResponse(msg.ID, map[string]interface{}{
			"echo":    "Request received",
			"action":  msg.Action,
			"message": "Handlers will be implemented in Phase 4",
		})
		client.SendMessage(resp)

	case TypeEvent:
		log.Printf("Event received: %s (events don't require responses)", msg.Event)

	default:
		errResp := NewErrorResponse(msg.ID, "Unknown message type")
		client.SendMessage(errResp)
	}

	// If hub has a handler (Phase 4), call it
	if h.hub.handler != nil {
		h.hub.handler.HandleMessage(ctx, client, msg)
	}
}

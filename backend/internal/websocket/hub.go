package websocket

import (
	"context"
	"log"
	"sync"
)

// MessageHandler defines the interface for handling WebSocket messages.
// This will be implemented in handler.go (Phase 4).
type MessageHandler interface {
	HandleMessage(ctx context.Context, client *Client, msg *Message)
}

// Hub maintains the set of active clients and broadcasts messages to clients.
type Hub struct {
	// Registered clients (clientID -> client)
	clients map[*Client]bool

	// Game rooms (gameID -> room)
	rooms map[string]*Room

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Shutdown signal
	shutdown chan struct{}

	// Mutex for thread-safe operations
	mu sync.RWMutex

	// Message handler (injected in Phase 4)
	handler MessageHandler
}

// NewHub creates a new Hub.
func NewHub(handler MessageHandler) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]*Room),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		shutdown:   make(chan struct{}),
		handler:    handler,
	}
}

// Run starts the hub's main event loop.
// This should be called in a goroutine.
func (h *Hub) Run() {
	log.Println("WebSocket Hub started")

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client registered: %s (user: %s). Total clients: %d", client.ID, client.UserID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				// Leave room if in one
				if client.room != nil {
					client.LeaveRoom()
				}

				log.Printf("Client unregistered: %s (user: %s). Total clients: %d", client.ID, client.UserID, len(h.clients))
			}
			h.mu.Unlock()

			// Clean up empty rooms after client disconnect
			h.CleanupEmptyRooms()

		case <-h.shutdown:
			log.Println("WebSocket Hub shutting down...")
			h.mu.Lock()
			// Close all client connections
			for client := range h.clients {
				close(client.send)
			}
			h.clients = make(map[*Client]bool)
			h.rooms = make(map[string]*Room)
			h.mu.Unlock()
			log.Println("WebSocket Hub stopped")
			return
		}
	}
}

// GetOrCreateRoom gets an existing room or creates a new one.
func (h *Hub) GetOrCreateRoom(gameID string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.rooms[gameID]
	if !exists {
		room = NewRoom(gameID)
		h.rooms[gameID] = room
		log.Printf("Room created: %s", gameID)
	}
	return room
}

// GetRoom gets a room by ID (returns nil if not found).
func (h *Hub) GetRoom(gameID string) *Room {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.rooms[gameID]
}

// RemoveRoom removes a room from the hub.
func (h *Hub) RemoveRoom(gameID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.rooms, gameID)
	log.Printf("Room removed: %s", gameID)
}

// CleanupEmptyRooms removes all rooms with no clients.
func (h *Hub) CleanupEmptyRooms() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for gameID, room := range h.rooms {
		if room.ClientCount() == 0 {
			delete(h.rooms, gameID)
			log.Printf("Empty room cleaned up: %s", gameID)
		}
	}
}

// Shutdown gracefully shuts down the hub.
func (h *Hub) Shutdown() {
	log.Println("Initiating WebSocket Hub shutdown...")
	close(h.shutdown)
}

// ClientCount returns the number of connected clients.
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// RoomCount returns the number of active rooms.
func (h *Hub) RoomCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.rooms)
}

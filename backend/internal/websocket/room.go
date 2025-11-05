package websocket

import (
	"log"
	"sync"
)

// Room represents a game room with multiple clients.
type Room struct {
	// Room ID (same as game ID)
	ID string

	// Clients in this room
	clients map[*Client]bool

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// NewRoom creates a new room.
func NewRoom(id string) *Room {
	return &Room{
		ID:      id,
		clients: make(map[*Client]bool),
	}
}

// AddClient adds a client to the room.
func (r *Room) AddClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[client] = true
	log.Printf("Client %s added to room %s. Room size: %d", client.ID, r.ID, len(r.clients))
}

// RemoveClient removes a client from the room.
func (r *Room) RemoveClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, client)
	log.Printf("Client %s removed from room %s. Room size: %d", client.ID, r.ID, len(r.clients))
}

// Broadcast sends a message to all clients in the room.
func (r *Room) Broadcast(msg *Message) {
	data, err := msg.ToJSON()
	if err != nil {
		log.Printf("Error encoding message for broadcast in room %s: %v", r.ID, err)
		return
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	for client := range r.clients {
		select {
		case client.send <- data:
			// Message queued successfully
		default:
			// Client's send buffer is full, skip this client
			log.Printf("Client %s send buffer full, skipping broadcast message", client.ID)
		}
	}

	log.Printf("Broadcast message to %d clients in room %s", len(r.clients), r.ID)
}

// BroadcastExcept sends a message to all clients except the specified one.
func (r *Room) BroadcastExcept(msg *Message, except *Client) {
	data, err := msg.ToJSON()
	if err != nil {
		log.Printf("Error encoding message for broadcast in room %s: %v", r.ID, err)
		return
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	sentCount := 0
	for client := range r.clients {
		if client == except {
			continue
		}

		select {
		case client.send <- data:
			sentCount++
		default:
			log.Printf("Client %s send buffer full, skipping broadcast message", client.ID)
		}
	}

	log.Printf("Broadcast message to %d clients in room %s (excluding sender)", sentCount, r.ID)
}

// ClientCount returns the number of clients in the room.
func (r *Room) ClientCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clients)
}

// GetClients returns a slice of all clients in the room.
func (r *Room) GetClients() []*Client {
	r.mu.RLock()
	defer r.mu.RUnlock()

	clients := make([]*Client, 0, len(r.clients))
	for client := range r.clients {
		clients = append(clients, client)
	}
	return clients
}

// HasClient checks if a client is in this room.
func (r *Room) HasClient(client *Client) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.clients[client]
	return exists
}

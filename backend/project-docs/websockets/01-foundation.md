# Phase 1: Foundation - Core WebSocket Components

**Goal**: Build core WebSocket infrastructure with all fixes applied

**Duration**: 4-5 hours

**Status**: 📋 Ready to implement

---

## 📋 What You'll Build

By the end of this phase, you'll have:
- ✅ Complete message protocol with proper imports
- ✅ Client wrapper with readPump/writePump
- ✅ Hub with connection management
- ✅ Room-based game organization
- ✅ All code compiling successfully
- ✅ Context propagation throughout
- ✅ Graceful shutdown support

---

## 🎯 Prerequisites

Before starting:
- [x] Read [00-overview-and-fixes.md](./00-overview-and-fixes.md)
- [ ] Go 1.25.3 installed
- [ ] Backend server currently running
- [ ] Terminal ready for commands

---

## 📦 Step 1: Install Dependencies (5 minutes)

```bash
cd /Users/ankit/code/learn/chess-coach/backend

# Install WebSocket library
go get github.com/gorilla/websocket

# Install UUID library (likely already installed)
go get github.com/google/uuid

# Verify installation
go mod tidy
```

**Verify**:
```bash
grep "github.com/gorilla/websocket" go.mod
grep "github.com/google/uuid" go.mod
```

---

## 📁 Step 2: Create Package Structure (2 minutes)

```bash
# Create websocket package directory
mkdir -p internal/websocket

# Verify
ls -la internal/websocket
```

---

## 📝 Step 3: Create Message Protocol (30 minutes)

### File: `internal/websocket/message.go`

**Fixes Applied**:
- ✅ Issue #1: Added `generateID()` function
- ✅ All imports present
- ✅ Complete implementation (no TODOs)

**Create the file**:

```go
// Package websocket provides real-time communication for chess games.
package websocket

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// MessageType defines the type of WebSocket message.
type MessageType string

const (
	// TypeRequest is a client request requiring a response
	TypeRequest MessageType = "request"
	// TypeResponse is a server response to a request
	TypeResponse MessageType = "response"
	// TypeEvent is a server-initiated event (no response expected)
	TypeEvent MessageType = "event"
)

// Message is the base WebSocket message structure.
type Message struct {
	ID        string                 `json:"id,omitempty"`        // For request/response matching
	Type      MessageType            `json:"type"`                // request|response|event
	Action    string                 `json:"action,omitempty"`    // Specific action (for requests)
	Event     string                 `json:"event,omitempty"`     // Event name (for events)
	Data      map[string]interface{} `json:"data,omitempty"`      // Payload
	Error     string                 `json:"error,omitempty"`     // Error message (for responses)
	Success   bool                   `json:"success,omitempty"`   // Success flag (for responses)
	Timestamp int64                  `json:"timestamp,omitempty"` // Unix timestamp
}

// NewRequest creates a request message.
func NewRequest(action string, data map[string]interface{}) *Message {
	return &Message{
		ID:        generateID(),
		Type:      TypeRequest,
		Action:    action,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// NewResponse creates a success response message.
func NewResponse(requestID string, data map[string]interface{}) *Message {
	return &Message{
		ID:        requestID,
		Type:      TypeResponse,
		Success:   true,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// NewErrorResponse creates an error response message.
func NewErrorResponse(requestID string, err string) *Message {
	return &Message{
		ID:        requestID,
		Type:      TypeResponse,
		Success:   false,
		Error:     err,
		Timestamp: time.Now().Unix(),
	}
}

// NewEvent creates an event message.
func NewEvent(event string, data map[string]interface{}) *Message {
	return &Message{
		Type:      TypeEvent,
		Event:     event,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// ToJSON converts message to JSON bytes.
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON parses JSON into message.
func FromJSON(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}

// generateID generates a unique message ID.
func generateID() string {
	return uuid.New().String()
}
```

**Test**:
```bash
cd internal/websocket
go build .
# Should compile successfully
```

---

## 📝 Step 4: Create Client Wrapper (45 minutes)

### File: `internal/websocket/client.go`

**Fixes Applied**:
- ✅ Issue #2: Added `fmt` import
- ✅ Issue #4: Added context propagation
- ✅ Issue #13: JoinRoom returns error
- ✅ Proper error handling throughout

**Create the file**:

```go
package websocket

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 1 * 1024 * 1024 // 1 MB
)

// Client represents a WebSocket client connection.
type Client struct {
	// Unique client ID
	ID string

	// User ID from authentication
	UserID string

	// Current game ID (if in a game)
	GameID string

	// WebSocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte

	// Reference to hub
	hub *Hub

	// Current room (if joined)
	room *Room

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// NewClient creates a new WebSocket client.
func NewClient(conn *websocket.Conn, hub *Hub, userID string) *Client {
	return &Client{
		ID:     generateID(),
		UserID: userID,
		conn:   conn,
		send:   make(chan []byte, 256),
		hub:    hub,
	}
}

// readPump pumps messages from the WebSocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump(ctx context.Context, messageHandler func(context.Context, *Client, []byte)) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for client %s: %v", c.ID, err)
			}
			break
		}

		// Handle message with context
		messageHandler(ctx, c, message)
	}
}

// writePump pumps messages from the hub to the WebSocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("Write error for client %s: %v", c.ID, err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendMessage sends a message to the client.
func (c *Client) SendMessage(msg *Message) error {
	data, err := msg.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	select {
	case c.send <- data:
		return nil
	default:
		return fmt.Errorf("client send buffer full")
	}
}

// JoinRoom adds client to a room.
// Returns error if client is already in a room or room is full.
func (c *Client) JoinRoom(room *Room) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if already in a room
	if c.room != nil {
		return fmt.Errorf("client already in room %s", c.room.ID)
	}

	// Check room capacity (chess is max 2 players)
	if room.ClientCount() >= 2 {
		return fmt.Errorf("room is full (max 2 players)")
	}

	c.room = room
	c.GameID = room.ID
	room.AddClient(c)

	log.Printf("Client %s (user %s) joined room %s", c.ID, c.UserID, room.ID)
	return nil
}

// LeaveRoom removes client from current room.
func (c *Client) LeaveRoom() {
	c.mu.Lock()
	room := c.room
	c.room = nil
	c.GameID = ""
	c.mu.Unlock()

	if room != nil {
		room.RemoveClient(c)
		log.Printf("Client %s (user %s) left room %s", c.ID, c.UserID, room.ID)
	}
}

// GetRoom returns the current room (thread-safe).
func (c *Client) GetRoom() *Room {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.room
}
```

**Test**:
```bash
cd internal/websocket
go build .
# Should compile successfully
```

---

## 📝 Step 5: Create Hub (60 minutes)

### File: `internal/websocket/hub.go`

**Fixes Applied**:
- ✅ Issue #4: Context support throughout
- ✅ Issue #14: Room cleanup method
- ✅ Issue #15: Graceful shutdown
- ✅ Message handler dependency injection

**Create the file**:

```go
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
```

**Test**:
```bash
cd internal/websocket
go build .
# Should compile successfully
```

---

## 📝 Step 6: Create Room Management (30 minutes)

### File: `internal/websocket/room.go`

**Fixes Applied**:
- ✅ Issue #10: Added `log` import
- ✅ Thread-safe operations
- ✅ Proper error handling

**Create the file**:

```go
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
```

**Test**:
```bash
cd internal/websocket
go build .
# Should compile successfully
```

---

## ✅ Step 7: Final Verification (10 minutes)

### Compile All Files

```bash
cd /Users/ankit/code/learn/chess-coach/backend

# Build the entire project
go build ./...

# Should see no errors
```

### Verify Package Structure

```bash
ls -la internal/websocket/

# Should see:
# client.go
# hub.go
# message.go
# room.go
```

### Run Go Vet

```bash
go vet ./internal/websocket/

# Should report no issues
```

### Check Imports

```bash
# Verify all imports are used
go mod tidy

# Check for any unused dependencies
```

---

## 📋 Phase 1 Checklist

Go to [implementation-checklist.md](./implementation-checklist.md) and check off:

**Phase 1: Foundation**
- [x] Install `gorilla/websocket`
- [x] Install `google/uuid`
- [x] Create directory `internal/websocket`
- [x] Create `message.go` with all functions
- [x] Create `client.go` with pumps and context
- [x] Create `hub.go` with shutdown support
- [x] Create `room.go` with thread safety
- [x] All files compile without errors
- [x] Run `go vet ./internal/websocket`
- [x] Run `go build ./...`

---

## 🎯 What You've Accomplished

### ✅ Issues Fixed in This Phase

| Issue # | Description | Status |
|---------|-------------|--------|
| 1 | Missing `generateID()` | ✅ Fixed in message.go |
| 2 | Missing `fmt` import | ✅ Fixed in client.go |
| 4 | Missing context | ✅ Added context propagation |
| 10 | Missing `log` import | ✅ Fixed in room.go |
| 13 | Poor room error handling | ✅ JoinRoom returns error |
| 14 | No room cleanup | ✅ CleanupEmptyRooms() added |
| 15 | No graceful shutdown | ✅ Shutdown() method added |

### ✅ Components Built

1. **Message Protocol** - Type-safe message handling
2. **Client Wrapper** - Concurrent read/write with pumps
3. **Hub** - Connection registry and lifecycle
4. **Room** - Game-based client grouping

### ✅ Production Features

- Thread-safe operations with mutexes
- Context propagation for cancellation
- Graceful shutdown support
- Comprehensive logging
- Error handling throughout
- Room capacity limits (2 players for chess)

---

## 🚀 Next Steps

**Phase 1 is complete!** Your WebSocket foundation is ready.

**Next**: [02-integration.md](./02-integration.md) - Connect to existing app

In Phase 2, you'll:
- Wire hub into app lifecycle
- Create WebSocket HTTP upgrade handler
- Register `/ws` endpoint in router
- Connect to your existing repositories

---

## 📚 Testing Your Work (Optional)

Want to test Phase 1 before moving on? Create a simple test:

### File: `internal/websocket/message_test.go`

```go
package websocket

import (
	"testing"
)

func TestNewRequest(t *testing.T) {
	req := NewRequest("testAction", map[string]interface{}{
		"key": "value",
	})

	if req.Type != TypeRequest {
		t.Errorf("Expected type %s, got %s", TypeRequest, req.Type)
	}

	if req.Action != "testAction" {
		t.Errorf("Expected action testAction, got %s", req.Action)
	}

	if req.ID == "" {
		t.Error("Expected non-empty ID")
	}
}

func TestNewResponse(t *testing.T) {
	resp := NewResponse("req-123", map[string]interface{}{
		"result": "success",
	})

	if resp.Type != TypeResponse {
		t.Errorf("Expected type %s, got %s", TypeResponse, resp.Type)
	}

	if !resp.Success {
		t.Error("Expected success to be true")
	}

	if resp.ID != "req-123" {
		t.Errorf("Expected ID req-123, got %s", resp.ID)
	}
}

func TestMessageSerialization(t *testing.T) {
	msg := NewEvent("testEvent", map[string]interface{}{
		"data": "test",
	})

	// Serialize
	data, err := msg.ToJSON()
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}

	// Deserialize
	parsed, err := FromJSON(data)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	if parsed.Type != TypeEvent {
		t.Errorf("Expected type %s, got %s", TypeEvent, parsed.Type)
	}

	if parsed.Event != "testEvent" {
		t.Errorf("Expected event testEvent, got %s", parsed.Event)
	}
}
```

**Run tests**:
```bash
cd internal/websocket
go test -v
```

---

## 💡 Key Concepts Review

Before moving to Phase 2, make sure you understand:

### 1. **The Pump Pattern**
- `readPump()` - Dedicated goroutine for reading
- `writePump()` - Dedicated goroutine for writing
- Communication via `send` channel

### 2. **The Hub Pattern**
- Central coordinator for all connections
- Event loop processes register/unregister/broadcast
- Thread-safe with channels and mutexes

### 3. **Message Protocol**
- Request/Response for actions needing confirmation
- Events for broadcasts (no response expected)
- ID matching for request/response correlation

### 4. **Room Management**
- Each game has its own room
- Broadcast to all players in a game
- Automatic cleanup when empty

---

**Ready for Phase 2?** → [02-integration.md](./02-integration.md)

---

**Last Updated**: 2025-11-05
**Status**: ✅ Complete and ready to implement
**Estimated Time**: 4-5 hours

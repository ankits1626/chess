# Step 25: WebSocket Implementation Plan for Chess-Coach

**Date**: November 3, 2025
**Status**: Planning
**Dependencies**: Lessons 1-3 completed

---

## Executive Summary

Integrate production-ready WebSocket communication into chess-coach backend to enable:
- Real-time game moves
- Live game state synchronization
- Player-to-player communication
- Game room management
- AI move streaming

**Timeline**: 2-3 days of focused development
**Complexity**: Medium (leveraging lessons learned)

---

## Phase 1: Foundation (Day 1)

### 1.1 Project Structure

```
backend/
├── internal/
│   ├── websocket/
│   │   ├── client.go          # WebSocket client wrapper
│   │   ├── hub.go             # Connection registry + rooms
│   │   ├── message.go         # Message protocol definitions
│   │   ├── handler.go         # Message routing and validation
│   │   └── room.go            # Game room management
│   ├── game/
│   │   ├── state.go           # Game state management
│   │   └── validator.go       # Move validation
│   └── api/
│       └── websocket.go       # HTTP upgrade handler
└── go.mod
```

### 1.2 Core Components

#### A. Message Protocol (`internal/websocket/message.go`)

```go
package websocket

import (
	"encoding/json"
	"time"
)

// MessageType defines the type of WebSocket message
type MessageType string

const (
	TypeRequest  MessageType = "request"
	TypeResponse MessageType = "response"
	TypeEvent    MessageType = "event"
)

// Message is the base WebSocket message structure
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

// Request creates a request message
func NewRequest(action string, data map[string]interface{}) *Message {
	return &Message{
		ID:        generateID(),
		Type:      TypeRequest,
		Action:    action,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// Response creates a response message
func NewResponse(requestID string, data map[string]interface{}) *Message {
	return &Message{
		ID:        requestID,
		Type:      TypeResponse,
		Success:   true,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// ErrorResponse creates an error response
func NewErrorResponse(requestID string, err string) *Message {
	return &Message{
		ID:        requestID,
		Type:      TypeResponse,
		Success:   false,
		Error:     err,
		Timestamp: time.Now().Unix(),
	}
}

// Event creates an event message
func NewEvent(event string, data map[string]interface{}) *Message {
	return &Message{
		Type:      TypeEvent,
		Event:     event,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

// ToJSON converts message to JSON bytes
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON parses JSON into message
func FromJSON(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}
```

#### B. Client Structure (`internal/websocket/client.go`)

```go
package websocket

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1 * 1024 * 1024 // 1 MB
)

// Client represents a WebSocket client connection
type Client struct {
	ID       string
	UserID   string // From JWT authentication
	conn     *websocket.Conn
	send     chan []byte
	hub      *Hub
	room     *Room
	mu       sync.RWMutex
}

// NewClient creates a new client
func NewClient(conn *websocket.Conn, hub *Hub, userID string) *Client {
	return &Client{
		ID:     generateID(),
		UserID: userID,
		conn:   conn,
		send:   make(chan []byte, 256),
		hub:    hub,
	}
}

// readPump handles incoming messages from client
func (c *Client) readPump() {
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
		_, messageData, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error: %v", err)
			}
			break
		}

		// Parse message
		msg, err := FromJSON(messageData)
		if err != nil {
			log.Printf("Invalid JSON from client %s: %v", c.ID, err)
			errResp := NewErrorResponse("", "Invalid JSON")
			c.SendMessage(errResp)
			continue
		}

		// Route message to handler
		c.hub.handleMessage(c, msg)
	}
}

// writePump handles outgoing messages to client
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
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
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

// SendMessage sends a message to the client
func (c *Client) SendMessage(msg *Message) error {
	data, err := msg.ToJSON()
	if err != nil {
		return err
	}

	select {
	case c.send <- data:
		return nil
	default:
		return fmt.Errorf("client send buffer full")
	}
}

// JoinRoom adds client to a room
func (c *Client) JoinRoom(room *Room) {
	c.mu.Lock()
	c.room = room
	c.mu.Unlock()
	room.AddClient(c)
}

// LeaveRoom removes client from current room
func (c *Client) LeaveRoom() {
	c.mu.Lock()
	room := c.room
	c.room = nil
	c.mu.Unlock()

	if room != nil {
		room.RemoveClient(c)
	}
}
```

#### C. Hub (Connection Registry) (`internal/websocket/hub.go`)

```go
package websocket

import (
	"log"
	"sync"
)

// Hub maintains active clients and rooms
type Hub struct {
	clients    map[*Client]bool
	rooms      map[string]*Room
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	handler    *MessageHandler
}

// NewHub creates a new Hub
func NewHub(handler *MessageHandler) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]*Room),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		handler:    handler,
	}
}

// Run starts the hub's event loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client %s registered. Total: %d", client.ID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				// Leave room if in one
				if client.room != nil {
					client.LeaveRoom()
				}
			}
			h.mu.Unlock()
			log.Printf("Client %s unregistered. Total: %d", client.ID, len(h.clients))
		}
	}
}

// GetOrCreateRoom gets existing room or creates new one
func (h *Hub) GetOrCreateRoom(gameID string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.rooms[gameID]
	if !exists {
		room = NewRoom(gameID)
		h.rooms[gameID] = room
		log.Printf("Room %s created", gameID)
	}
	return room
}

// RemoveRoom removes a room
func (h *Hub) RemoveRoom(gameID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.rooms, gameID)
	log.Printf("Room %s removed", gameID)
}

// handleMessage routes incoming messages to appropriate handlers
func (h *Hub) handleMessage(client *Client, msg *Message) {
	switch msg.Type {
	case TypeRequest:
		h.handler.HandleRequest(client, msg)
	case TypeEvent:
		h.handler.HandleEvent(client, msg)
	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}
```

#### D. Room Management (`internal/websocket/room.go`)

```go
package websocket

import (
	"sync"
)

// Room represents a game room with multiple clients
type Room struct {
	ID      string
	clients map[*Client]bool
	mu      sync.RWMutex
}

// NewRoom creates a new room
func NewRoom(id string) *Room {
	return &Room{
		ID:      id,
		clients: make(map[*Client]bool),
	}
}

// AddClient adds a client to the room
func (r *Room) AddClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[client] = true
}

// RemoveClient removes a client from the room
func (r *Room) RemoveClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, client)
}

// Broadcast sends a message to all clients in the room
func (r *Room) Broadcast(msg *Message) {
	data, err := msg.ToJSON()
	if err != nil {
		log.Printf("Error encoding message: %v", err)
		return
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	for client := range r.clients {
		select {
		case client.send <- data:
		default:
			log.Printf("Client %s send buffer full, skipping", client.ID)
		}
	}
}

// BroadcastExcept sends message to all except the sender
func (r *Room) BroadcastExcept(msg *Message, except *Client) {
	data, err := msg.ToJSON()
	if err != nil {
		return
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	for client := range r.clients {
		if client == except {
			continue
		}
		select {
		case client.send <- data:
		default:
			// Skip if buffer full
		}
	}
}

// ClientCount returns number of clients in room
func (r *Room) ClientCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clients)
}
```

---

## Phase 2: Message Handlers (Day 1-2)

### 2.1 Handler Structure (`internal/websocket/handler.go`)

```go
package websocket

import (
	"log"
)

// MessageHandler handles incoming WebSocket messages
type MessageHandler struct {
	hub *Hub
	// Add dependencies (DB, game service, etc.)
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(hub *Hub) *MessageHandler {
	return &MessageHandler{
		hub: hub,
	}
}

// HandleRequest handles request-type messages
func (h *MessageHandler) HandleRequest(client *Client, msg *Message) {
	switch msg.Action {
	case "createGame":
		h.handleCreateGame(client, msg)
	case "joinGame":
		h.handleJoinGame(client, msg)
	case "makeMove":
		h.handleMakeMove(client, msg)
	case "getGameState":
		h.handleGetGameState(client, msg)
	case "leaveGame":
		h.handleLeaveGame(client, msg)
	default:
		resp := NewErrorResponse(msg.ID, "Unknown action: "+msg.Action)
		client.SendMessage(resp)
	}
}

// HandleEvent handles event-type messages
func (h *MessageHandler) HandleEvent(client *Client, msg *Message) {
	switch msg.Event {
	case "chat":
		h.handleChat(client, msg)
	default:
		log.Printf("Unknown event: %s", msg.Event)
	}
}
```

### 2.2 Game Actions

#### Create Game
```go
func (h *MessageHandler) handleCreateGame(client *Client, msg *Message) {
	// 1. Validate input
	timeControl, ok := msg.Data["timeControl"].(string)
	if !ok {
		client.SendMessage(NewErrorResponse(msg.ID, "timeControl required"))
		return
	}

	// 2. Create game in database
	gameID := generateGameID()
	// ... create game logic ...

	// 3. Create room
	room := h.hub.GetOrCreateRoom(gameID)
	client.JoinRoom(room)

	// 4. Send response
	resp := NewResponse(msg.ID, map[string]interface{}{
		"gameId": gameID,
		"url":    "/game/" + gameID,
		"color":  "white", // Creator is white
	})
	client.SendMessage(resp)

	// 5. Broadcast event to room
	event := NewEvent("gameCreated", map[string]interface{}{
		"gameId": gameID,
		"players": map[string]interface{}{
			"white": client.UserID,
		},
	})
	room.Broadcast(event)
}
```

#### Join Game
```go
func (h *MessageHandler) handleJoinGame(client *Client, msg *Message) {
	gameID, ok := msg.Data["gameId"].(string)
	if !ok {
		client.SendMessage(NewErrorResponse(msg.ID, "gameId required"))
		return
	}

	// Get room
	room := h.hub.GetOrCreateRoom(gameID)

	// Check if room is full
	if room.ClientCount() >= 2 {
		client.SendMessage(NewErrorResponse(msg.ID, "Game is full"))
		return
	}

	// Join room
	client.JoinRoom(room)

	// Send response
	resp := NewResponse(msg.ID, map[string]interface{}{
		"gameId": gameID,
		"color":  "black", // Second player is black
	})
	client.SendMessage(resp)

	// Broadcast player joined
	event := NewEvent("playerJoined", map[string]interface{}{
		"userId": client.UserID,
		"color":  "black",
	})
	room.BroadcastExcept(event, client)

	// Start game if both players present
	if room.ClientCount() == 2 {
		startEvent := NewEvent("gameStarted", map[string]interface{}{
			"gameId": gameID,
		})
		room.Broadcast(startEvent)
	}
}
```

#### Make Move
```go
func (h *MessageHandler) handleMakeMove(client *Client, msg *Message) {
	from, _ := msg.Data["from"].(string)
	to, _ := msg.Data["to"].(string)

	if from == "" || to == "" {
		client.SendMessage(NewErrorResponse(msg.ID, "from and to required"))
		return
	}

	// Validate client is in a room
	if client.room == nil {
		client.SendMessage(NewErrorResponse(msg.ID, "Not in a game"))
		return
	}

	// Validate move (check turn, legal move, etc.)
	// ... validation logic ...

	// Update game state in database
	// ... update logic ...

	// Send success response to player
	resp := NewResponse(msg.ID, map[string]interface{}{
		"status": "success",
	})
	client.SendMessage(resp)

	// Broadcast move to opponent
	moveEvent := NewEvent("opponentMove", map[string]interface{}{
		"from":      from,
		"to":        to,
		"timestamp": time.Now().Unix(),
	})
	client.room.BroadcastExcept(moveEvent, client)
}
```

---

## Phase 3: HTTP Integration (Day 2)

### 3.1 WebSocket Endpoint (`internal/api/websocket.go`)

```go
package api

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	ws "chess-coach/internal/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implement proper origin checking
		return true
	},
}

// HandleWebSocket upgrades HTTP to WebSocket
func HandleWebSocket(hub *ws.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Authenticate user from JWT token
		userID, err := authenticateRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 2. Upgrade connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Upgrade error: %v", err)
			return
		}

		// 3. Create client
		client := ws.NewClient(conn, hub, userID)

		// 4. Register client
		hub.register <- client

		// 5. Start goroutines
		go client.writePump()
		go client.readPump()
	}
}

// authenticateRequest extracts and validates JWT token
func authenticateRequest(r *http.Request) (string, error) {
	// Get token from query param or header
	token := r.URL.Query().Get("token")
	if token == "" {
		token = r.Header.Get("Authorization")
	}

	// Validate token and extract user ID
	// ... JWT validation logic ...

	return "user123", nil // Placeholder
}
```

### 3.2 Main Server Setup

```go
package main

import (
	"chess-coach/internal/api"
	"chess-coach/internal/websocket"
	"log"
	"net/http"
)

func main() {
	// Create hub and handler
	handler := websocket.NewMessageHandler(nil)
	hub := websocket.NewHub(handler)
	handler.hub = hub

	// Start hub
	go hub.Run()

	// Setup routes
	http.HandleFunc("/ws", api.HandleWebSocket(hub))

	// Existing REST routes
	// ...

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## Phase 4: Testing (Day 2-3)

### 4.1 Unit Tests

```go
// internal/websocket/message_test.go
func TestMessageCreation(t *testing.T) {
	msg := NewRequest("createGame", map[string]interface{}{
		"timeControl": "5+0",
	})

	assert.Equal(t, TypeRequest, msg.Type)
	assert.Equal(t, "createGame", msg.Action)
	assert.NotEmpty(t, msg.ID)
}

// internal/websocket/room_test.go
func TestRoomBroadcast(t *testing.T) {
	room := NewRoom("game123")
	// ... test broadcast logic ...
}
```

### 4.2 Integration Tests

Create test client in Go:

```go
// test/websocket_test.go
func TestGameFlow(t *testing.T) {
	// 1. Connect two clients
	client1 := connectTestClient("user1")
	client2 := connectTestClient("user2")

	// 2. Client1 creates game
	gameID := client1.CreateGame("5+0")

	// 3. Client2 joins game
	client2.JoinGame(gameID)

	// 4. Client1 makes move
	client1.MakeMove("e2", "e4")

	// 5. Client2 receives move
	move := client2.WaitForEvent("opponentMove")
	assert.Equal(t, "e2", move.Data["from"])
}
```

### 4.3 Load Testing

```bash
# Using wscat or custom tool
for i in {1..100}; do
    wscat -c ws://localhost:8080/ws?token=test &
done
```

---

## Phase 5: Frontend Integration (Day 3)

### 5.1 React WebSocket Hook

```typescript
// frontend/src/hooks/useWebSocket.ts
import { useEffect, useRef, useState } from 'react';

interface Message {
  id?: string;
  type: 'request' | 'response' | 'event';
  action?: string;
  event?: string;
  data?: any;
  error?: string;
  success?: boolean;
}

export function useWebSocket(url: string) {
  const ws = useRef<WebSocket | null>(null);
  const [connected, setConnected] = useState(false);
  const [lastMessage, setLastMessage] = useState<Message | null>(null);
  const pendingRequests = useRef<Map<string, (msg: Message) => void>>(new Map());

  useEffect(() => {
    ws.current = new WebSocket(url);

    ws.current.onopen = () => {
      console.log('WebSocket connected');
      setConnected(true);
    };

    ws.current.onmessage = (event) => {
      const msg: Message = JSON.parse(event.data);

      // Handle responses
      if (msg.type === 'response' && msg.id) {
        const resolver = pendingRequests.current.get(msg.id);
        if (resolver) {
          resolver(msg);
          pendingRequests.current.delete(msg.id);
          return;
        }
      }

      // Handle events
      setLastMessage(msg);
    };

    ws.current.onclose = () => {
      console.log('WebSocket disconnected');
      setConnected(false);
    };

    return () => {
      ws.current?.close();
    };
  }, [url]);

  const sendRequest = (action: string, data: any): Promise<Message> => {
    return new Promise((resolve) => {
      const id = `req-${Date.now()}`;
      const msg: Message = {
        id,
        type: 'request',
        action,
        data,
      };

      pendingRequests.current.set(id, resolve);
      ws.current?.send(JSON.stringify(msg));
    });
  };

  const sendEvent = (event: string, data: any) => {
    const msg: Message = {
      type: 'event',
      event,
      data,
    };
    ws.current?.send(JSON.stringify(msg));
  };

  return {
    connected,
    lastMessage,
    sendRequest,
    sendEvent,
  };
}
```

### 5.2 Game Component

```typescript
// frontend/src/components/Game.tsx
import { useWebSocket } from '../hooks/useWebSocket';

export function Game({ gameId }: { gameId: string }) {
  const { connected, lastMessage, sendRequest, sendEvent } =
    useWebSocket(`ws://localhost:8080/ws?token=${getToken()}`);

  useEffect(() => {
    if (connected && gameId) {
      // Join game
      sendRequest('joinGame', { gameId });
    }
  }, [connected, gameId]);

  useEffect(() => {
    if (lastMessage?.event === 'opponentMove') {
      // Update board with opponent's move
      const { from, to } = lastMessage.data;
      updateBoard(from, to);
    }
  }, [lastMessage]);

  const makeMove = async (from: string, to: string) => {
    const response = await sendRequest('makeMove', { from, to });
    if (!response.success) {
      alert(response.error);
    }
  };

  return (
    <div>
      <div>Status: {connected ? 'Connected' : 'Disconnected'}</div>
      {/* Render chess board */}
    </div>
  );
}
```

---

## Phase 6: Production Readiness

### 6.1 Security

- [ ] JWT token validation
- [ ] Rate limiting per client
- [ ] Input sanitization
- [ ] CORS configuration
- [ ] TLS/WSS in production

### 6.2 Monitoring

```go
// Add metrics
var (
	activeConnections = prometheus.NewGauge(...)
	messagesReceived  = prometheus.NewCounter(...)
	messagesSent      = prometheus.NewCounter(...)
)
```

### 6.3 Configuration

```yaml
# config.yaml
websocket:
  ping_period: 54s
  pong_wait: 60s
  write_wait: 10s
  max_message_size: 1048576
  max_clients_per_room: 2
  read_buffer_size: 4096
  write_buffer_size: 4096
```

---

## Testing Checklist

- [ ] Unit tests for message protocol
- [ ] Unit tests for room management
- [ ] Integration test: Create and join game
- [ ] Integration test: Make and receive moves
- [ ] Integration test: Client disconnect/reconnect
- [ ] Integration test: Room cleanup
- [ ] Load test: 100 concurrent connections
- [ ] Load test: 50 concurrent games
- [ ] Manual test: Network disconnect
- [ ] Manual test: Server restart
- [ ] Manual test: Invalid messages

---

## Deployment

### Docker
```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o chess-coach
EXPOSE 8080
CMD ["./chess-coach"]
```

### Environment Variables
```
JWT_SECRET=your-secret
DATABASE_URL=postgres://...
WS_ORIGIN=https://chess-coach.com
```

---

## Next Steps

1. **Review this plan** with the team
2. **Start with Phase 1** - Foundation
3. **Test incrementally** after each phase
4. **Document as you go** - Update this file with actual implementation details

---

## Success Criteria

✅ Two players can create and join a game
✅ Players can make moves in real-time
✅ Moves are validated before broadcasting
✅ Disconnected players can reconnect
✅ Server handles 100+ concurrent connections
✅ All tests pass
✅ Frontend integrates smoothly

---

**Let's build this! 🚀**

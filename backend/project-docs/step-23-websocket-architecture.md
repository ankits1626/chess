# WebSocket Architecture for Chess Coach

**Future-proof architecture for computer play, multiplayer, and AI coaching**

---

## Table of Contents
1. [Overview](#overview)
2. [Architecture Design](#architecture-design)
3. [Phase 1: WebSocket Infrastructure](#phase-1-websocket-infrastructure)
4. [Phase 2: Computer Play via WebSocket](#phase-2-computer-play-via-websocket)
5. [Phase 3: Multiplayer Support](#phase-3-multiplayer-support)
6. [Phase 4: AI Coach Integration](#phase-4-ai-coach-integration)
7. [Testing & Deployment](#testing--deployment)

---

## Overview

### Why WebSocket First?

**Your future requirements:**
1. ✅ Human vs. Computer (current)
2. ✅ Human vs. Human online (future)
3. ✅ AI Coach chat during games (future)

**All of these benefit from real-time bidirectional communication!**

### Architecture Benefits

```
┌──────────────────────────────────────────┐
│  Single WebSocket Connection Per User    │
├──────────────────────────────────────────┤
│  Handles:                                │
│  • Computer move generation              │
│  • Opponent move notifications           │
│  • AI coach streaming responses          │
│  • Position analysis                     │
│  • Game state synchronization            │
└──────────────────────────────────────────┘
```

**One infrastructure, multiple features!**

---

## Architecture Design

### Message Types

```go
// All messages flow through WebSocket
type WSMessage struct {
    Type    string      `json:"type"`
    GameID  string      `json:"game_id,omitempty"`
    Payload interface{} `json:"payload"`
}

// Message Types:
// - "game.move"         : Player makes move
// - "game.move.computer": Computer move response
// - "game.move.opponent": Opponent move notification
// - "coach.question"    : Ask AI coach
// - "coach.response"    : AI coach answer (streaming)
// - "game.join"         : Join game room
// - "game.leave"        : Leave game room
// - "game.sync"         : Sync game state
```

### Connection Flow

```
Client connects
    ↓
Authenticate (JWT)
    ↓
Join game room
    ↓
Subscribe to events
    ↓
Send/receive messages
    ↓
Leave room on disconnect
```

---

## Phase 1: WebSocket Infrastructure

### Step 1.1: Install Dependencies

```bash
cd backend
go get github.com/gorilla/websocket
go get github.com/google/uuid
```

---

### Step 1.2: Create WebSocket Package Structure

```bash
mkdir -p internal/websocket
mkdir -p internal/websocket/message
```

**Directory structure:**
```
internal/websocket/
├── hub.go           # Connection pool manager
├── client.go        # Individual client connection
├── handler.go       # HTTP → WebSocket upgrade
├── room.go          # Game room management
└── message/
    ├── types.go     # Message type definitions
    └── handlers.go  # Message routing
```

---

### Step 1.3: Create Hub (Connection Manager)

**File:** `internal/websocket/hub.go`

```go
// Package websocket provides real-time communication for chess games.
package websocket

import (
    "sync"

    "github.com/google/uuid"
)

// Hub maintains active client connections and broadcasts messages.
type Hub struct {
    // Registered clients (userID -> client)
    clients map[string]*Client

    // Game rooms (gameID -> set of clients)
    rooms map[string]map[*Client]bool

    // Register client
    register chan *Client

    // Unregister client
    unregister chan *Client

    // Broadcast message to room
    broadcast chan *BroadcastMessage

    // Mutex for thread-safe operations
    mu sync.RWMutex
}

// BroadcastMessage contains message and target room.
type BroadcastMessage struct {
    RoomID  string
    Message []byte
    Exclude *Client // Optional: exclude sender
}

// NewHub creates a new WebSocket hub.
func NewHub() *Hub {
    return &Hub{
        clients:    make(map[string]*Client),
        rooms:      make(map[string]map[*Client]bool),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        broadcast:  make(chan *BroadcastMessage),
    }
}

// Run starts the hub's main loop.
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client.UserID] = client
            h.mu.Unlock()

        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client.UserID]; ok {
                delete(h.clients, client.UserID)
                close(client.send)

                // Remove from all rooms
                for roomID, room := range h.rooms {
                    if _, inRoom := room[client]; inRoom {
                        delete(room, client)
                        if len(room) == 0 {
                            delete(h.rooms, roomID)
                        }
                    }
                }
            }
            h.mu.Unlock()

        case message := <-h.broadcast:
            h.mu.RLock()
            room := h.rooms[message.RoomID]
            for client := range room {
                // Skip excluded client (usually sender)
                if message.Exclude != nil && client == message.Exclude {
                    continue
                }

                select {
                case client.send <- message.Message:
                default:
                    // Client buffer full, close connection
                    close(client.send)
                    delete(h.clients, client.UserID)
                    delete(room, client)
                }
            }
            h.mu.RUnlock()
        }
    }
}

// JoinRoom adds client to a game room.
func (h *Hub) JoinRoom(roomID string, client *Client) {
    h.mu.Lock()
    defer h.mu.Unlock()

    if h.rooms[roomID] == nil {
        h.rooms[roomID] = make(map[*Client]bool)
    }
    h.rooms[roomID][client] = true
}

// LeaveRoom removes client from a game room.
func (h *Hub) LeaveRoom(roomID string, client *Client) {
    h.mu.Lock()
    defer h.mu.Unlock()

    if room, ok := h.rooms[roomID]; ok {
        delete(room, client)
        if len(room) == 0 {
            delete(h.rooms, roomID)
        }
    }
}

// SendToRoom broadcasts message to all clients in room.
func (h *Hub) SendToRoom(roomID string, message []byte, exclude *Client) {
    h.broadcast <- &BroadcastMessage{
        RoomID:  roomID,
        Message: message,
        Exclude: exclude,
    }
}

// SendToUser sends message to specific user.
func (h *Hub) SendToUser(userID string, message []byte) {
    h.mu.RLock()
    defer h.mu.RUnlock()

    if client, ok := h.clients[userID]; ok {
        select {
        case client.send <- message:
        default:
            // Buffer full
        }
    }
}
```

---

### Step 1.4: Create Client Connection

**File:** `internal/websocket/client.go`

```go
package websocket

import (
    "encoding/json"
    "log"
    "time"

    "github.com/gorilla/websocket"
)

const (
    // Time allowed to write message to peer
    writeWait = 10 * time.Second

    // Time allowed to read next pong message from peer
    pongWait = 60 * time.Second

    // Send pings to peer with this period (must be less than pongWait)
    pingPeriod = (pongWait * 9) / 10

    // Maximum message size allowed from peer
    maxMessageSize = 512 * 1024 // 512 KB
)

// Client represents a WebSocket client connection.
type Client struct {
    UserID string
    GameID string
    hub    *Hub
    conn   *websocket.Conn
    send   chan []byte
}

// NewClient creates a new WebSocket client.
func NewClient(userID string, hub *Hub, conn *websocket.Conn) *Client {
    return &Client{
        UserID: userID,
        hub:    hub,
        conn:   conn,
        send:   make(chan []byte, 256),
    }
}

// ReadPump pumps messages from WebSocket connection to hub.
func (c *Client) ReadPump(messageHandler func(*Client, []byte)) {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()

    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    c.conn.SetReadLimit(maxMessageSize)

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("WebSocket error: %v", err)
            }
            break
        }

        // Handle message
        messageHandler(c, message)
    }
}

// WritePump pumps messages from hub to WebSocket connection.
func (c *Client) WritePump() {
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
                // Hub closed the channel
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            w, err := c.conn.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            w.Write(message)

            // Add queued messages to current websocket message
            n := len(c.send)
            for i := 0; i < n; i++ {
                w.Write([]byte{'\n'})
                w.Write(<-c.send)
            }

            if err := w.Close(); err != nil {
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

// SendJSON sends JSON message to client.
func (c *Client) SendJSON(v interface{}) error {
    data, err := json.Marshal(v)
    if err != nil {
        return err
    }

    select {
    case c.send <- data:
    default:
        // Buffer full
    }

    return nil
}
```

---

### Step 1.5: Create Message Types

**File:** `internal/websocket/message/types.go`

```go
// Package message defines WebSocket message types.
package message

// Type represents message type.
type Type string

const (
    // Game messages
    TypeGameJoin          Type = "game.join"
    TypeGameLeave         Type = "game.leave"
    TypeGameMove          Type = "game.move"
    TypeGameMoveComputer  Type = "game.move.computer"
    TypeGameMoveOpponent  Type = "game.move.opponent"
    TypeGameSync          Type = "game.sync"
    TypeGameEnd           Type = "game.end"

    // Coach messages
    TypeCoachQuestion  Type = "coach.question"
    TypeCoachResponse  Type = "coach.response"
    TypeCoachAnalysis  Type = "coach.analysis"

    // System messages
    TypeError   Type = "error"
    TypePing    Type = "ping"
    TypePong    Type = "pong"
)

// Message is the base WebSocket message structure.
type Message struct {
    Type    Type        `json:"type"`
    GameID  string      `json:"game_id,omitempty"`
    Payload interface{} `json:"payload"`
}

// GameJoinPayload for joining a game.
type GameJoinPayload struct {
    GameID   string `json:"game_id"`
    PlayerID string `json:"player_id"`
}

// GameMovePayload for making a move.
type GameMovePayload struct {
    Move       string `json:"move"`        // UCI format
    MoveSAN    string `json:"move_san"`    // SAN format
    FEN        string `json:"fen"`         // Position after move
    MoveNumber int    `json:"move_number"`
    Side       string `json:"side"`        // "white" or "black"
}

// ComputerMovePayload for computer move response.
type ComputerMovePayload struct {
    Move       string `json:"move"`
    MoveSAN    string `json:"move_san"`
    Evaluation int    `json:"evaluation"`
    FEN        string `json:"fen"`
    Depth      int    `json:"depth"`
}

// CoachQuestionPayload for AI coach question.
type CoachQuestionPayload struct {
    Question string `json:"question"`
    FEN      string `json:"fen,omitempty"`
    Context  string `json:"context,omitempty"`
}

// CoachResponsePayload for AI coach answer.
type CoachResponsePayload struct {
    Answer    string `json:"answer"`
    Streaming bool   `json:"streaming"`
    Token     string `json:"token,omitempty"` // For streaming responses
    Done      bool   `json:"done"`
}

// ErrorPayload for error messages.
type ErrorPayload struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

---

## Phase 2: Computer Play via WebSocket

### Step 2.1: Create Game Message Handler

**File:** `internal/websocket/message/game_handler.go`

```go
package message

import (
    "context"
    "encoding/json"
    "fmt"
    "log"

    "github.com/ankits1626/chess-coach-backend/internal/service"
)

// GameHandler handles game-related messages.
type GameHandler struct {
    engineService *service.EngineService
}

// NewGameHandler creates a new game message handler.
func NewGameHandler(engineService *service.EngineService) *GameHandler {
    return &GameHandler{
        engineService: engineService,
    }
}

// HandleMove handles user move and generates computer response.
func (h *GameHandler) HandleMove(ctx context.Context, payload GameMovePayload) (*ComputerMovePayload, error) {
    // Request computer move from engine
    req := service.SuggestMoveRequest{
        FEN:        payload.FEN,
        Difficulty: 10, // TODO: Get from game settings
    }

    response, err := h.engineService.SuggestMove(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("failed to get computer move: %w", err)
    }

    return &ComputerMovePayload{
        Move:       response.Move,
        MoveSAN:    response.MoveSAN,
        Evaluation: response.Evaluation,
        Depth:      10,
    }, nil
}
```

---

### Step 2.2: Create WebSocket Handler

**File:** `internal/websocket/handler.go`

```go
package websocket

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/ankits1626/chess-coach-backend/internal/service"
    wsmsg "github.com/ankits1626/chess-coach-backend/internal/websocket/message"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // TODO: Implement proper CORS check
        return true
    },
}

// Handler handles WebSocket connections.
type Handler struct {
    hub           *Hub
    gameHandler   *wsmsg.GameHandler
    engineService *service.EngineService
}

// NewHandler creates a new WebSocket handler.
func NewHandler(hub *Hub, engineService *service.EngineService) *Handler {
    return &Handler{
        hub:           hub,
        gameHandler:   wsmsg.NewGameHandler(engineService),
        engineService: engineService,
    }
}

// ServeWS handles WebSocket connection upgrade.
func (h *Handler) ServeWS(c *gin.Context) {
    // TODO: Get user ID from JWT token
    userID := c.Query("user_id")
    if userID == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id required"})
        return
    }

    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Printf("Failed to upgrade connection: %v", err)
        return
    }

    client := NewClient(userID, h.hub, conn)
    h.hub.register <- client

    // Start read/write pumps
    go client.WritePump()
    go client.ReadPump(h.handleMessage)
}

// handleMessage routes messages to appropriate handlers.
func (h *Handler) handleMessage(client *Client, data []byte) {
    var msg wsmsg.Message
    if err := json.Unmarshal(data, &msg); err != nil {
        log.Printf("Failed to parse message: %v", err)
        return
    }

    switch msg.Type {
    case wsmsg.TypeGameJoin:
        h.handleGameJoin(client, msg)

    case wsmsg.TypeGameLeave:
        h.handleGameLeave(client, msg)

    case wsmsg.TypeGameMove:
        h.handleGameMove(client, msg)

    case wsmsg.TypeCoachQuestion:
        h.handleCoachQuestion(client, msg)

    case wsmsg.TypePing:
        h.handlePing(client)

    default:
        log.Printf("Unknown message type: %s", msg.Type)
    }
}

// handleGameJoin adds client to game room.
func (h *Handler) handleGameJoin(client *Client, msg wsmsg.Message) {
    var payload wsmsg.GameJoinPayload
    data, _ := json.Marshal(msg.Payload)
    if err := json.Unmarshal(data, &payload); err != nil {
        log.Printf("Invalid game join payload: %v", err)
        return
    }

    client.GameID = payload.GameID
    h.hub.JoinRoom(payload.GameID, client)

    log.Printf("User %s joined game %s", client.UserID, payload.GameID)
}

// handleGameLeave removes client from game room.
func (h *Handler) handleGameLeave(client *Client, msg wsmsg.Message) {
    if client.GameID != "" {
        h.hub.LeaveRoom(client.GameID, client)
        log.Printf("User %s left game %s", client.UserID, client.GameID)
    }
}

// handleGameMove processes player move and generates computer response.
func (h *Handler) handleGameMove(client *Client, msg wsmsg.Message) {
    var payload wsmsg.GameMovePayload
    data, _ := json.Marshal(msg.Payload)
    if err := json.Unmarshal(data, &payload); err != nil {
        log.Printf("Invalid move payload: %v", err)
        return
    }

    // Get computer move
    computerMove, err := h.gameHandler.HandleMove(client.conn.UnderlyingConn().(*net.TCPConn).LocalAddr().Network(), payload)
    if err != nil {
        client.SendJSON(wsmsg.Message{
            Type: wsmsg.TypeError,
            Payload: wsmsg.ErrorPayload{
                Code:    "COMPUTER_MOVE_FAILED",
                Message: err.Error(),
            },
        })
        return
    }

    // Send computer move back to client
    client.SendJSON(wsmsg.Message{
        Type:    wsmsg.TypeGameMoveComputer,
        GameID:  msg.GameID,
        Payload: computerMove,
    })
}

// handleCoachQuestion handles AI coach questions.
func (h *Handler) handleCoachQuestion(client *Client, msg wsmsg.Message) {
    // TODO: Implement AI coach integration
    client.SendJSON(wsmsg.Message{
        Type: wsmsg.TypeCoachResponse,
        Payload: wsmsg.CoachResponsePayload{
            Answer: "AI Coach coming soon!",
            Done:   true,
        },
    })
}

// handlePing responds to ping.
func (h *Handler) handlePing(client *Client) {
    client.SendJSON(wsmsg.Message{Type: wsmsg.TypePong})
}
```

---

## Summary

### What This Architecture Gives You

✅ **Computer Play** - WebSocket-based move generation
✅ **Multiplayer Ready** - Already has room management
✅ **AI Coach Ready** - Message types defined, streaming support
✅ **Scalable** - Hub pattern handles many connections
✅ **Reusable** - One infrastructure for all features

### Next Steps

1. **Integrate with existing REST API** (save moves to DB)
2. **Add authentication** (JWT validation)
3. **Build frontend WebSocket client**
4. **Test computer gameplay**
5. **Add multiplayer support**
6. **Integrate AI coach**

---

**Estimated Time:**
- Phase 1: 4-5 hours (infrastructure)
- Phase 2: 2-3 hours (computer play)
- Phase 3: 2-3 hours (multiplayer)
- Phase 4: 3-4 hours (AI coach)

**Total: ~12-15 hours for complete system**

---

**Last Updated**: 2025-11-01

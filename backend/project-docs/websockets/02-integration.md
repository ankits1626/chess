# Phase 2: Integration - Connect WebSocket to Existing App

**Goal**: Wire WebSocket infrastructure into your existing chess-coach application

**Duration**: 2-3 hours

**Status**: 📋 Ready to implement

---

## 📋 What You'll Do

By the end of this phase:
- ✅ Hub starts when app starts
- ✅ Hub shuts down gracefully when app stops
- ✅ WebSocket endpoint accessible at `/ws`
- ✅ Placeholder authentication (real auth in Phase 5)
- ✅ Message routing connected (handlers in Phase 4)

---

## 🎯 Prerequisites

- [x] Phase 1 complete (WebSocket package built)
- [x] Backend server currently running
- [x] Understanding of your app structure

---

## 📝 Step 1: Create Upgrade Handler (45 minutes)

This file handles HTTP → WebSocket upgrade and message routing.

### File: `internal/websocket/upgrade.go`

**Create new file**:

```go
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
```

**Test compilation**:
```bash
cd internal/websocket
go build .
# Should compile successfully
```

---

## 📝 Step 2: Modify App Lifecycle (30 minutes)

Wire the hub into your application's lifecycle.

### File: `internal/app/app.go` (Modify)

**Read current file first**:
```bash
cat internal/app/app.go
```

**Apply these changes**:

```go
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
	"github.com/ankits1626/chess-coach-backend/internal/websocket"  // ← ADD THIS IMPORT
)

// App manages application lifecycle.
type App struct {
	config *config.Config
	db     *database.DB
	server server.Server
	logger logger.Logger
	wsHub  *websocket.Hub  // ← ADD THIS FIELD
}

// New creates new application.
func New(cfg *config.Config, db *database.DB, log logger.Logger) *App {
	// Create WebSocket hub (nil handler for now, Phase 4 will add real handler)
	hub := websocket.NewHub(nil)  // ← ADD THIS

	return &App{
		config: cfg,
		db:     db,
		server: server.New(cfg, db, hub),  // ← PASS HUB TO SERVER
		logger: log,
		wsHub:  hub,  // ← STORE HUB
	}
}

// Run starts application and handles graceful shutdown.
func (a *App) Run(ctx context.Context) error {
	// Start WebSocket hub
	go a.wsHub.Run()  // ← ADD THIS
	a.logger.Info("WebSocket Hub started")

	// Start server in goroutine
	go func() {
		if err := a.server.Start(); err != nil {
			a.logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	a.logger.Infof("Server started on port %s", a.config.Port)
	a.logger.Infof("Swagger UI: http://localhost:%s/swagger/index.html", a.config.Port)
	a.logger.Infof("WebSocket endpoint: ws://localhost:%s/ws?user_id=test", a.config.Port)  // ← ADD THIS

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("Shutting down server...")

	// Shutdown WebSocket hub first
	a.wsHub.Shutdown()  // ← ADD THIS
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
	if a.db != nil {
		a.db.Close()
	}
}
```

**Key changes**:
1. ✅ Import websocket package
2. ✅ Add `wsHub` field to App struct
3. ✅ Create hub in `New()`
4. ✅ Pass hub to server
5. ✅ Start hub in `Run()`
6. ✅ Shutdown hub before server
7. ✅ Log WebSocket endpoint

---

## 📝 Step 3: Modify Server (20 minutes)

Update server to accept and use the hub.

### File: `internal/server/server.go` (Modify)

**Read current file first**:
```bash
cat internal/server/server.go
```

**Apply these changes**:

```go
// Package server manages HTTP server lifecycle.
package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ankits1626/chess-coach-backend/internal/config"
	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/router"
	"github.com/ankits1626/chess-coach-backend/internal/websocket"  // ← ADD THIS IMPORT
	"github.com/gin-gonic/gin"
)

// Server handles HTTP server start/stop.
type Server interface {
	Start() error
	Shutdown(ctx context.Context) error
}

// server is the concrete implementation of Server.
type server struct {
	httpServer *http.Server
	config     *config.Config
	db         *database.DB
}

// New creates a new Server with given config.
func New(cfg *config.Config, db *database.DB, hub *websocket.Hub) Server {  // ← ADD HUB PARAMETER
	gin.SetMode(cfg.GinMode)
	r := router.Setup(db, hub)  // ← PASS HUB TO ROUTER

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	return &server{
		httpServer: httpServer,
		config:     cfg,
		db:         db,
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
```

**Key changes**:
1. ✅ Import websocket package
2. ✅ Add `hub` parameter to `New()`
3. ✅ Pass hub to router

---

## 📝 Step 4: Modify Router (30 minutes)

Register the WebSocket endpoint in your router.

### File: `internal/router/router.go` (Modify)

**Read current file first**:
```bash
cat internal/router/router.go
```

**Apply these changes**:

```go
// Package router provides HTTP routing setup.
package router

import (
	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/websocket"  // ← ADD THIS IMPORT
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/ankits1626/chess-coach-backend/docs"
)

// Setup creates and configures the Gin router.
func Setup(db *database.DB, hub *websocket.Hub) *gin.Engine {  // ← ADD HUB PARAMETER
	r := gin.Default()

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// WebSocket endpoint (before API routes)
	wsHandler := websocket.NewHandler(hub)  // ← ADD THIS
	r.GET("/ws", wsHandler.ServeWS)         // ← ADD THIS

	// Register API v1 routes
	RegisterV1Routes(r, db)

	return r
}
```

**Key changes**:
1. ✅ Import websocket package
2. ✅ Add `hub` parameter to `Setup()`
3. ✅ Create WebSocket handler
4. ✅ Register `/ws` endpoint

---

## 📝 Step 5: Verify Everything Compiles (10 minutes)

### Build the entire project

```bash
cd /Users/ankit/code/learn/chess-coach/backend

# Clean and build
go mod tidy
go build ./...
```

**Expected output**: No errors

### Run go vet

```bash
go vet ./...
```

**Expected output**: No issues

### Check imports

```bash
# Verify websocket package is importable
go list -m github.com/gorilla/websocket
go list -m github.com/google/uuid
```

---

## 📝 Step 6: Start the Server (5 minutes)

### Start your server

```bash
# If using Air (hot reload)
air

# Or directly
go run cmd/server/main.go
```

### Verify logs

You should see:
```
WebSocket Hub started
Server started on port 8080
Swagger UI: http://localhost:8080/swagger/index.html
WebSocket endpoint: ws://localhost:8080/ws?user_id=test
```

---

## 📝 Step 7: Test WebSocket Connection (20 minutes)

### Option A: Test with wscat (Recommended)

**Install wscat**:
```bash
npm install -g wscat
```

**Connect to WebSocket**:
```bash
wscat -c 'ws://localhost:8080/ws?user_id=alice'
```

**You should see**:
```
Connected (press CTRL+C to quit)
>
```

**In server logs, you should see**:
```
WebSocket connection established for user: alice (client: <uuid>)
Client registered: <uuid> (user: alice). Total clients: 1
```

**Send a test message**:
```json
{"type":"request","action":"test","data":{"hello":"world"}}
```

**You should receive**:
```json
{
  "id": "<uuid>",
  "type": "response",
  "success": true,
  "data": {
    "echo": "Request received",
    "action": "test",
    "message": "Handlers will be implemented in Phase 4"
  },
  "timestamp": 1699000000
}
```

**Disconnect** (Ctrl+C):
```
Disconnected
```

**Server logs should show**:
```
Client unregistered: <uuid> (user: alice). Total clients: 0
```

---

### Option B: Test with Browser Console

**Open browser** to `http://localhost:8080/swagger/index.html`

**Open Developer Console** (F12)

**Run this JavaScript**:
```javascript
// Connect
const ws = new WebSocket('ws://localhost:8080/ws?user_id=bob');

// Handle connection open
ws.onopen = () => {
    console.log('Connected!');

    // Send test message
    ws.send(JSON.stringify({
        type: 'request',
        action: 'test',
        data: { hello: 'from browser' }
    }));
};

// Handle messages
ws.onmessage = (event) => {
    console.log('Received:', JSON.parse(event.data));
};

// Handle close
ws.onclose = () => {
    console.log('Disconnected');
};

// Handle errors
ws.onerror = (error) => {
    console.error('WebSocket error:', error);
};
```

**You should see**:
```
Connected!
Received: {
  id: "...",
  type: "response",
  success: true,
  data: { ... }
}
```

**To disconnect**:
```javascript
ws.close();
```

---

## 📝 Step 8: Test Multiple Connections (15 minutes)

### Test concurrent connections

**Terminal 1**:
```bash
wscat -c 'ws://localhost:8080/ws?user_id=alice'
```

**Terminal 2**:
```bash
wscat -c 'ws://localhost:8080/ws?user_id=bob'
```

**Server logs should show**:
```
Client registered: <uuid1> (user: alice). Total clients: 1
Client registered: <uuid2> (user: bob). Total clients: 2
```

**Disconnect both** (Ctrl+C in each terminal)

**Server logs should show**:
```
Client unregistered: <uuid1> (user: alice). Total clients: 1
Client unregistered: <uuid2> (user: bob). Total clients: 0
```

---

## 📝 Step 9: Test Graceful Shutdown (10 minutes)

### Start server with active connection

1. **Start server**: `air` or `go run cmd/server/main.go`
2. **Connect client**: `wscat -c 'ws://localhost:8080/ws?user_id=test'`
3. **Verify connection** in logs

### Shutdown server

**Press Ctrl+C in server terminal**

**Server logs should show**:
```
^C
Shutting down server...
Initiating WebSocket Hub shutdown...
WebSocket Hub shutting down...
WebSocket Hub stopped
Server exited gracefully
```

**Client should disconnect gracefully**

---

## ✅ Phase 2 Checklist

Go to [implementation-checklist.md](./implementation-checklist.md) and check off:

**Phase 2: Integration**
- [x] Create `internal/websocket/upgrade.go`
- [x] Import websocket package in `app.go`
- [x] Add `wsHub` field to App struct
- [x] Create hub in `App.New()`
- [x] Pass hub to server
- [x] Start hub with `go hub.Run()`
- [x] Add hub shutdown before server shutdown
- [x] Import websocket package in `server.go`
- [x] Add hub parameter to `server.New()`
- [x] Pass hub to router
- [x] Import websocket package in `router.go`
- [x] Add hub parameter to `router.Setup()`
- [x] Create WebSocket handler
- [x] Register `/ws` endpoint
- [x] Server compiles and starts
- [x] Navigate to `/ws` (connection test)
- [x] Check logs for graceful shutdown
- [x] Test with wscat or browser

---

## 🎯 What You've Accomplished

### ✅ Issues Fixed in This Phase

| Issue # | Description | Status |
|---------|-------------|--------|
| 6 | No repository DI | 🟡 Prepared (Phase 4) |
| 7 | No app integration | ✅ Fixed - Hub in app lifecycle |
| 8 | No router integration | ✅ Fixed - /ws endpoint registered |

### ✅ Integration Points

1. **App Lifecycle** - Hub starts/stops with app
2. **Server** - Accepts and uses hub
3. **Router** - WebSocket endpoint accessible
4. **Graceful Shutdown** - Hub closes before server
5. **Logging** - All events logged

### ✅ Testing Capabilities

- Connect via wscat
- Connect via browser
- Multiple concurrent connections
- Message echo (placeholder)
- Graceful disconnection
- Server shutdown

---

## 🚀 Next Steps

**Phase 2 is complete!** Your WebSocket is integrated into your app.

**Next**: [03-utilities.md](./03-utilities.md) - Helper functions

In Phase 3, you'll create:
- UUID conversion utilities (`pgtype.UUID` ↔ `string`)
- Context helpers
- Validation utilities
- Error response helpers

---

## 🐛 Troubleshooting

### Server won't compile

**Check**:
```bash
go mod tidy
go build ./...
```

**Common issues**:
- Missing imports
- Typos in function signatures
- Hub parameter not passed through

---

### Can't connect to WebSocket

**Check**:
1. Server is running
2. URL is correct: `ws://localhost:8080/ws?user_id=test`
3. Port 8080 is not blocked
4. Check server logs for errors

**Test with curl** (should fail to upgrade):
```bash
curl -i http://localhost:8080/ws?user_id=test
# Should see "Upgrade required"
```

---

### Messages not received

**Check**:
1. Message is valid JSON
2. Client is connected (check logs)
3. Server logs show message received
4. No errors in writePump

**Enable verbose logging**:
```go
// In upgrade.go handleMessage, add:
log.Printf("Raw message: %s", string(data))
```

---

### Multiple clients disconnect

**Check**:
1. Hub.Run() is running (check logs)
2. No panic in hub event loop
3. Channels not closed prematurely

**Add debug logging**:
```go
// In hub.go Run(), add at start:
log.Println("Hub event loop iteration")
```

---

## 📊 Verification Checklist

Before moving to Phase 3, verify:

- [ ] Server starts without errors
- [ ] Hub starts (check logs)
- [ ] Can connect to `/ws?user_id=test`
- [ ] Can send message and receive response
- [ ] Can disconnect gracefully
- [ ] Multiple clients can connect
- [ ] Server shutdown closes all connections
- [ ] All logs show expected messages

**All checked?** You're ready for Phase 3! ✅

---

## 💡 Key Concepts Review

### Integration Points

1. **App** → Creates and owns hub
2. **Server** → Receives hub, passes to router
3. **Router** → Creates handler, registers endpoint
4. **Handler** → Upgrades connections, manages clients

### Data Flow

```
HTTP Request → Router → Handler
                          ↓
                    Upgrade to WebSocket
                          ↓
                    Create Client
                          ↓
                    Register with Hub
                          ↓
              Start readPump + writePump
```

### Lifecycle

```
App Start
  ├─ Hub.Run() starts
  ├─ Server starts
  └─ Ready for connections

App Shutdown
  ├─ Hub.Shutdown() called
  ├─ All clients disconnected
  ├─ Hub.Run() exits
  └─ Server.Shutdown() called
```

---

**Ready for Phase 3?** → [03-utilities.md](./03-utilities.md)

---

**Last Updated**: 2025-11-05
**Status**: ✅ Complete and ready to implement
**Estimated Time**: 2-3 hours
**Dependencies**: Phase 1 complete

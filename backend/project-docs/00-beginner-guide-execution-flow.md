# Chess Coach Backend - Complete Execution Flow Guide for Beginners

> **Last Updated:** 2025-11-05
> **Purpose:** Understand how the Go backend server runs from start to finish, including HTTP REST API, WebSocket real-time communication, and database interactions.

---

## Table of Contents

1. [Quick Overview](#quick-overview)
2. [Directory Structure](#directory-structure)
3. [The Big Picture: How Everything Works Together](#the-big-picture)
4. [Phase 1: Application Startup](#phase-1-application-startup)
5. [Phase 2: Server Initialization](#phase-2-server-initialization)
6. [Phase 3: HTTP REST API Flow](#phase-3-http-rest-api-flow)
7. [Phase 4: WebSocket Real-Time Communication](#phase-4-websocket-real-time-communication)
8. [Phase 5: Database Layer](#phase-5-database-layer)
9. [Complete Request Flow Examples](#complete-request-flow-examples)
10. [Concurrency & Goroutines](#concurrency--goroutines)
11. [Key Design Patterns](#key-design-patterns)
12. [Troubleshooting & Common Issues](#troubleshooting--common-issues)

---

## Quick Overview

The Chess Coach backend is a Go application that provides:

- **REST API** for CRUD operations (users, games, moves)
- **WebSocket** for real-time chess game updates
- **PostgreSQL** database for persistent storage
- **Graceful shutdown** for clean server termination

**Tech Stack:**
- **Language:** Go 1.23+
- **Web Framework:** Gin (HTTP router)
- **WebSocket:** gorilla/websocket
- **Database:** PostgreSQL with pgx driver
- **SQL Type Safety:** sqlc (generates Go code from SQL)

---

## Directory Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go                    # 🚀 ENTRY POINT - Program starts here
│
├── internal/
│   ├── app/
│   │   └── app.go                     # 🎯 Application lifecycle manager
│   │
│   ├── config/
│   │   └── config.go                  # ⚙️ Environment configuration
│   │
│   ├── server/
│   │   └── server.go                  # 🌐 HTTP server wrapper
│   │
│   ├── router/
│   │   ├── router.go                  # 🛣️ Main router setup
│   │   └── v1_routes.go               # 📍 API v1 route registration
│   │
│   ├── handler/v1/
│   │   ├── user/handler.go            # 👤 User HTTP endpoints
│   │   ├── game/handler.go            # ♟️ Game HTTP endpoints
│   │   ├── move/handler.go            # 🎯 Move HTTP endpoints
│   │   └── health/handler.go          # ❤️ Health check
│   │
│   ├── repository/
│   │   ├── user_repository.go         # 💾 User data access
│   │   ├── game_repository.go         # 💾 Game data access
│   │   └── move_repository.go         # 💾 Move data access
│   │
│   ├── database/
│   │   ├── connection.go              # 🔌 Database connection pool
│   │   ├── db.go                      # 🗄️ sqlc generated base
│   │   └── *.sql.go                   # 📝 sqlc generated queries
│   │
│   ├── websocket/
│   │   ├── hub.go                     # 🎪 WebSocket hub (manager)
│   │   ├── client.go                  # 🔌 WebSocket client connection
│   │   ├── room.go                    # 🏠 Game room management
│   │   ├── message.go                 # 📨 Message types
│   │   └── upgrade.go                 # 🔄 HTTP → WebSocket upgrade
│   │
│   └── logger/
│       └── logger.go                  # 📋 Logging utility
│
└── sql/
    ├── schema/                        # 🗃️ Database schema
    └── queries/                       # 📝 SQL queries (for sqlc)
```

---

## The Big Picture

### System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLIENT LAYER                            │
│  - Web Browser (HTTP REST API)                                  │
│  - WebSocket Client (Real-time updates)                         │
└────────────────┬────────────────────────────────────┬───────────┘
                 │                                    │
                 ▼                                    ▼
┌────────────────────────────┐    ┌──────────────────────────────┐
│    HTTP REST API           │    │    WEBSOCKET ENDPOINT        │
│    (Gin Router)            │    │    GET /ws                   │
│                            │    │                              │
│  GET  /api/v1/users        │    │  • Upgrade HTTP → WS         │
│  POST /api/v1/users        │    │  • Register client           │
│  GET  /api/v1/games        │    │  • Start read/write pumps    │
│  POST /api/v1/games        │    │                              │
│  ...                       │    │                              │
└────────────┬───────────────┘    └──────────┬───────────────────┘
             │                               │
             ▼                               ▼
┌────────────────────────┐    ┌─────────────────────────────────┐
│   HANDLER LAYER        │    │      WEBSOCKET HUB              │
│   (HTTP Handlers)      │    │   (Connection Manager)          │
│                        │    │                                 │
│  • Parse requests      │    │  • Manages all clients          │
│  • Validate input      │    │  • Manages game rooms           │
│  • Call repository     │    │  • Routes messages              │
│  • Return JSON         │    │  • Broadcasts updates           │
└────────────┬───────────┘    └─────────────┬───────────────────┘
             │                               │
             ▼                               ▼
┌──────────────────────────────────────────────────────────────────┐
│                     REPOSITORY LAYER                             │
│   (Business Logic & Data Access)                                 │
│                                                                  │
│  • UserRepository                                                │
│  • GameRepository                                                │
│  • MoveRepository                                                │
└─────────────────────────────┬────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│                     DATABASE LAYER                               │
│   (sqlc Generated Queries + pgxpool)                             │
│                                                                  │
│  • Type-safe SQL queries                                         │
│  • Connection pooling                                            │
│  • Transaction support                                           │
└─────────────────────────────┬────────────────────────────────────┘
                              │
                              ▼
┌──────────────────────────────────────────────────────────────────┐
│                   POSTGRESQL DATABASE                            │
│                                                                  │
│  Tables: users, games, moves                                     │
└──────────────────────────────────────────────────────────────────┘
```

---

## Phase 1: Application Startup

### Entry Point: `cmd/server/main.go`

When you run `go run cmd/server/main.go`, this is what happens step by step:

```
┌─────────────────────────────────────────────────────────────────┐
│                    APPLICATION BOOTSTRAP                        │
└─────────────────────────────────────────────────────────────────┘

1️⃣  Load Configuration
    ├─ Read environment variables
    ├─ Set defaults (PORT=8080, DB settings, etc.)
    └─ File: internal/config/config.go
        ↓

2️⃣  Initialize Logger
    ├─ Create standard logger instance
    └─ File: internal/logger/logger.go
        ↓

3️⃣  Connect to Database
    ├─ Create PostgreSQL connection pool (pgxpool)
    ├─ Ping database to verify connection
    ├─ File: internal/database/connection.go
    └─ Connection string: postgres://user:pass@host:port/dbname
        ↓

4️⃣  Create Application Instance
    ├─ Initialize WebSocket Hub
    ├─ Initialize HTTP Server (with router)
    ├─ File: internal/app/app.go
    └─ Wire all components together
        ↓

5️⃣  Run Application
    ├─ Start WebSocket Hub (goroutine)
    ├─ Start HTTP Server (goroutine)
    ├─ Log server information
    └─ Wait for shutdown signal (SIGINT/SIGTERM)
```

### Code Flow in `main.go`

```go
// File: cmd/server/main.go

func main() {
    ctx := context.Background()

    // Step 1: Load configuration
    cfg := config.Load()

    // Step 2: Initialize logger
    appLogger := logger.NewStdLogger()

    // Step 3: Connect to database
    db, err := database.NewDB(ctx, cfg.DSN())
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    // Step 4: Create application
    application := app.New(cfg, db, appLogger)

    // Step 5: Run application (blocks until shutdown)
    if err := application.Run(ctx); err != nil {
        log.Fatal("Application error:", err)
    }
}
```

---

## Phase 2: Server Initialization

### Application Lifecycle: `internal/app/app.go`

The `app.App` struct manages the entire application lifecycle:

```go
type App struct {
    config *config.Config
    db     *database.DB
    server *server.Server
    logger logger.Logger
    wsHub  *websocket.Hub
}
```

### Initialization Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                   app.New() - Create Application                │
└─────────────────────────────────────────────────────────────────┘

1️⃣  Create WebSocket Hub
    └─ websocket.NewHub(nil)
        • Hub manages all WebSocket connections
        • Runs in its own goroutine
        ↓

2️⃣  Create HTTP Server
    └─ server.New(cfg, db, hub)
        • Sets up Gin router
        • Registers all routes
        • Configures HTTP server
        ↓

3️⃣  Return App Instance
    └─ Contains all initialized components
```

### Runtime Flow: `app.Run()`

```
┌─────────────────────────────────────────────────────────────────┐
│                    app.Run() - Start Server                     │
└─────────────────────────────────────────────────────────────────┘

1️⃣  Start WebSocket Hub (Background Goroutine)
    go wsHub.Run()
    • Listens for client registrations
    • Manages rooms and message routing
    • Runs until shutdown signal
        ↓

2️⃣  Start HTTP Server (Background Goroutine)
    go server.Start()
    • Listens on configured port (default: 8080)
    • Handles HTTP and WebSocket requests
    • Runs until shutdown signal
        ↓

3️⃣  Log Server Information
    • Server running on: http://localhost:8080
    • Swagger UI: http://localhost:8080/swagger/index.html
    • WebSocket: ws://localhost:8080/ws
        ↓

4️⃣  Wait for Shutdown Signal
    • Blocks and waits for SIGINT (Ctrl+C) or SIGTERM
        ↓

5️⃣  Graceful Shutdown (When signal received)
    a) Shutdown WebSocket Hub
       • Stop accepting new connections
       • Close all active connections
       • Clean up rooms

    b) Shutdown HTTP Server (with 5s timeout)
       • Stop accepting new requests
       • Wait for active requests to complete
       • Force close after timeout

    c) Database connections auto-close (defer in main)
```

### Router Setup: `internal/router/router.go`

```
┌─────────────────────────────────────────────────────────────────┐
│              router.Setup() - Configure Routes                  │
└─────────────────────────────────────────────────────────────────┘

1️⃣  Create Gin Engine
    gin.Default()
    • Logger middleware (logs requests)
    • Recovery middleware (panic recovery)
        ↓

2️⃣  Register Swagger Documentation
    GET /swagger/*any
    • API documentation UI
    • Auto-generated from code comments
        ↓

3️⃣  Register WebSocket Endpoint
    GET /ws
    • Upgrades HTTP → WebSocket connection
    • Handled by websocket.Handler
        ↓

4️⃣  Register API v1 Routes
    /api/v1/...
    • User routes
    • Game routes
    • Move routes
    • Health check
```

---

## Phase 3: HTTP REST API Flow

### Route Registration: `internal/router/v1_routes.go`

```
┌─────────────────────────────────────────────────────────────────┐
│                    API v1 ROUTES                                │
└─────────────────────────────────────────────────────────────────┘

/api/v1/health
    └─ GET → Health check (always returns OK)

/api/v1/users
    ├─ GET    → List all users
    ├─ POST   → Create new user
    ├─ GET    /:id → Get user by ID
    ├─ PUT    /:id → Update user
    ├─ DELETE /:id → Delete user
    └─ GET    /:id/games → Get user's games

/api/v1/games
    ├─ GET    → List all games
    ├─ POST   → Create new game
    ├─ GET    /:id → Get game by ID
    └─ GET    /:id/moves → Get game's moves

/api/v1/moves
    └─ GET /:id → Get move by ID
```

### Request Flow Example: Creating a User

Let's trace what happens when you send `POST /api/v1/users`:

```
┌─────────────────────────────────────────────────────────────────┐
│            REQUEST: POST /api/v1/users                          │
│            BODY: {"username": "alice", "email": "a@b.com"}      │
└─────────────────────────────────────────────────────────────────┘

1️⃣  Gin Router (router/v1_routes.go:37)
    • Matches route: POST /api/v1/users
    • Calls: userHandler.Create
        ↓

2️⃣  Handler (handler/v1/user/handler.go)
    func (h *UserHandler) Create(c *gin.Context) {

        a) Bind JSON Request
           • Parse request body
           • Validate required fields (username, email)
           • Map to CreateRequest struct

        b) Set Default Rating
           • If not provided, default to 1200

        c) Call Repository
           user, err := h.repo.Create(ctx, username, email, rating)

        d) Handle Response
           • If error → return 500
           • If success → return 201 with user data
    }
        ↓

3️⃣  Repository (repository/user_repository.go)
    func (r *UserRepository) Create(...) {

        a) Build Database Parameters
           params := database.CreateUserParams{
               Username: username,
               Email:    email,
               Rating:   rating,
           }

        b) Call Database Layer
           user, err := r.db.CreateUser(ctx, params)

        c) Return Result
           • User struct or error
    }
        ↓

4️⃣  Database Layer (database/users.sql.go - generated by sqlc)
    func (q *Queries) CreateUser(ctx, params) {

        a) Execute SQL Query
           INSERT INTO users (username, email, rating)
           VALUES ($1, $2, $3)
           RETURNING id, username, email, rating, created_at

        b) Scan Result
           • Map database row to User struct

        c) Return User or Error
    }
        ↓

5️⃣  PostgreSQL Database
    • Executes INSERT statement
    • Returns new row with generated ID
        ↓

6️⃣  Response Flow (Reverse Direction)
    Database → Repository → Handler → Client

    Final Response:
    {
        "id": "550e8400-...",
        "username": "alice",
        "email": "a@b.com",
        "rating": 1200,
        "created_at": "2025-11-05T12:00:00Z"
    }
```

### Layer Responsibilities

| Layer | Responsibility | Example File |
|-------|---------------|-------------|
| **Handler** | HTTP request/response, validation | `handler/v1/user/handler.go` |
| **Repository** | Business logic, data transformation | `repository/user_repository.go` |
| **Database** | SQL queries, type conversion | `database/users.sql.go` |

---

## Phase 4: WebSocket Real-Time Communication

WebSocket enables real-time bidirectional communication between clients and server. Perfect for chess games where moves need to be instantly shared!

### WebSocket Components

```
┌─────────────────────────────────────────────────────────────────┐
│                    WEBSOCKET ARCHITECTURE                       │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────┐
│      HUB        │  Central manager (1 per server)
│  (websocket/    │  • Registers/unregisters clients
│   hub.go)       │  • Manages game rooms
│                 │  • Routes messages
└────────┬────────┘
         │
         ├─── Manages ───┐
         │               │
         ▼               ▼
┌──────────────┐  ┌──────────────┐
│   CLIENT 1   │  │   CLIENT 2   │  One per WebSocket connection
│ (websocket/  │  │ (websocket/  │  • Read pump (inbound)
│  client.go)  │  │  client.go)  │  • Write pump (outbound)
└──────┬───────┘  └───────┬──────┘
       │                  │
       └────── In Room ───┘
              │
              ▼
       ┌──────────────┐
       │     ROOM     │  One per game
       │ (websocket/  │  • Contains 2 clients (chess players)
       │  room.go)    │  • Broadcasts messages to both
       └──────────────┘
```

### Connection Flow

#### Step 1: Establish WebSocket Connection

```
┌─────────────────────────────────────────────────────────────────┐
│    CLIENT → SERVER: ws://localhost:8080/ws?user_id=alice        │
└─────────────────────────────────────────────────────────────────┘

1️⃣  HTTP Request arrives at Gin Router
    GET /ws?user_id=alice
        ↓

2️⃣  WebSocket Handler (websocket/upgrade.go)
    func (h *Handler) ServeWS(c *gin.Context) {

        a) Extract user_id from query parameter
           userID := c.Query("user_id")
           // TODO Phase 5: Extract from JWT token instead

        b) Upgrade HTTP → WebSocket
           conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
           • Uses gorilla/websocket library
           • Switches protocol from HTTP to WebSocket

        c) Create Client
           client := websocket.NewClient(conn, hub, userID)
           • Unique client ID generated
           • References hub for communication

        d) Register Client with Hub
           hub.register <- client
           • Sends client to hub's register channel

        e) Start Client Goroutines
           go client.writePump()  // Handles outbound messages
           go client.readPump()   // Handles inbound messages
    }
```

#### Step 2: Hub Manages Clients

```
┌─────────────────────────────────────────────────────────────────┐
│              HUB EVENT LOOP (websocket/hub.go)                  │
│              Running continuously in goroutine                  │
└─────────────────────────────────────────────────────────────────┘

Hub Structure:
    type Hub struct {
        clients    map[*Client]bool      // All connected clients
        rooms      map[string]*Room      // Game rooms (gameID → room)
        register   chan *Client          // New client channel
        unregister chan *Client          // Remove client channel
        shutdown   chan struct{}         // Shutdown signal
    }

Event Loop:
    for {
        select {

        case client := <-h.register:
            // New client connected
            ├─ Add to clients map
            └─ Log connection

        case client := <-h.unregister:
            // Client disconnected
            ├─ Remove from clients map
            ├─ Close client send channel
            ├─ Remove from room (if in one)
            ├─ Delete empty rooms
            └─ Log disconnection

        case <-h.shutdown:
            // Server shutting down
            ├─ Close all client connections
            ├─ Clear clients and rooms
            └─ Exit goroutine
        }
    }
```

#### Step 3: Client Read/Write Pumps

Each WebSocket client has TWO goroutines:

##### Read Pump (Receives messages from client)

```
┌─────────────────────────────────────────────────────────────────┐
│           CLIENT.READPUMP() (websocket/client.go)               │
│           Runs continuously in goroutine                        │
└─────────────────────────────────────────────────────────────────┘

Setup:
    • Set read deadline (60 seconds)
    • Set pong handler (resets deadline on pong)

Read Loop:
    for {
        1️⃣  Read message from WebSocket
            _, message, err := c.conn.ReadMessage()

        2️⃣  If error (connection closed, timeout)
            → Break loop and unregister

        3️⃣  Process message
            messageHandler(ctx, c, message)
            • Parse JSON
            • Route by message type
            • Send response if needed
    }

On Exit:
    • Unregister from hub
    • Close connection
```

##### Write Pump (Sends messages to client)

```
┌─────────────────────────────────────────────────────────────────┐
│          CLIENT.WRITEPUMP() (websocket/client.go)               │
│          Runs continuously in goroutine                         │
└─────────────────────────────────────────────────────────────────┘

Setup:
    • Create ping ticker (54 seconds)

Write Loop:
    for {
        select {

        case message := <-c.send:
            // Message to send to client
            1️⃣  Set write deadline (10 seconds)
            2️⃣  Write message to WebSocket
            3️⃣  If error → exit loop

        case <-ticker.C:
            // Ping timeout
            1️⃣  Send ping message
            2️⃣  Keep connection alive
        }
    }

On Exit:
    • Close connection
```

### Message Flow

#### Message Structure (`websocket/message.go`)

```json
{
    "id": "unique-message-id",
    "type": "request|response|event",
    "action": "join_game",
    "event": "player_joined",
    "data": {
        "game_id": "550e8400-...",
        "player": "alice"
    },
    "error": "",
    "success": true,
    "timestamp": 1699200000
}
```

#### Message Types

| Type | Direction | Purpose | Response Required? |
|------|-----------|---------|-------------------|
| **request** | Client → Server | Request action | Yes |
| **response** | Server → Client | Reply to request | No |
| **event** | Server → Client | Notify event | No |

#### Example: Join Game Flow

```
┌─────────────────────────────────────────────────────────────────┐
│              EXAMPLE: PLAYER JOINS GAME ROOM                    │
└─────────────────────────────────────────────────────────────────┘

1️⃣  Client sends request:
    {
        "id": "msg-123",
        "type": "request",
        "action": "join_game",
        "data": {
            "game_id": "game-456"
        }
    }
        ↓

2️⃣  Server (readPump) receives message
        ↓

3️⃣  handleMessage() processes request
    a) Parse JSON
    b) Validate message type and action
    c) Get or create room for game_id
    d) Add client to room
        ↓

4️⃣  Server sends response to client:
    {
        "id": "msg-123",
        "type": "response",
        "success": true,
        "data": {
            "message": "Joined game-456"
        }
    }
        ↓

5️⃣  Server broadcasts event to ALL players in room:
    {
        "type": "event",
        "event": "player_joined",
        "data": {
            "player": "alice",
            "game_id": "game-456"
        }
    }
```

### Room Broadcasting

```
┌─────────────────────────────────────────────────────────────────┐
│                ROOM BROADCASTING (websocket/room.go)            │
└─────────────────────────────────────────────────────────────────┘

Room Structure:
    type Room struct {
        ID      string
        clients map[*Client]bool
    }

Methods:

1️⃣  Broadcast(message []byte)
    • Sends message to ALL clients in room

    for client := range r.clients {
        select {
        case client.send <- message:
            // Message queued
        default:
            // Channel full, close client
        }
    }

2️⃣  BroadcastExcept(message []byte, exceptClient *Client)
    • Sends to all EXCEPT one client
    • Used for: "opponent made move" notifications

    for client := range r.clients {
        if client == exceptClient {
            continue
        }
        client.send <- message
    }
```

---

## Phase 5: Database Layer

### Database Connection (`internal/database/connection.go`)

```
┌─────────────────────────────────────────────────────────────────┐
│              DATABASE CONNECTION POOL                           │
└─────────────────────────────────────────────────────────────────┘

Function: database.NewDB(ctx, dsn)

1️⃣  Create Connection Pool
    pool, err := pgxpool.New(ctx, dsn)
    • DSN: postgres://user:pass@host:port/dbname?sslmode=disable
    • Uses pgx library (PostgreSQL driver)
    • Connection pooling for performance
        ↓

2️⃣  Verify Connection
    err = pool.Ping(ctx)
    • Tests database connectivity
    • Fails fast if database unavailable
        ↓

3️⃣  Create DB Wrapper
    return &DB{
        Queries: New(pool),
        pool:    pool,
    }
```

### Database Structure

```go
type DB struct {
    *Queries
    pool *pgxpool.Pool
}
```

### sqlc: Type-Safe SQL

The project uses **sqlc** to generate Go code from SQL queries.

#### How sqlc Works

```
┌─────────────────────────────────────────────────────────────────┐
│                     sqlc WORKFLOW                               │
└─────────────────────────────────────────────────────────────────┘

1️⃣  Write SQL Queries
    File: sql/queries/users.sql

    -- name: GetUser :one
    SELECT * FROM users WHERE id = $1;

    -- name: CreateUser :one
    INSERT INTO users (username, email, rating)
    VALUES ($1, $2, $3)
    RETURNING *;
        ↓

2️⃣  Run sqlc Generate
    $ sqlc generate
        ↓

3️⃣  sqlc Generates Go Code
    File: internal/database/users.sql.go

    func (q *Queries) GetUser(ctx, id) (User, error)
    func (q *Queries) CreateUser(ctx, params) (User, error)
        ↓

4️⃣  Use in Repository
    user, err := db.GetUser(ctx, userID)
```

---

## Complete Request Flow Examples

### Example 1: REST API - Get User by ID

```
GET /api/v1/users/550e8400-e29b-41d4-a716-446655440000

1️⃣  Gin Router → userHandler.GetByID
2️⃣  Handler parses UUID, validates
3️⃣  Calls repository.GetByID
4️⃣  Repository calls database.GetUser
5️⃣  Database executes SQL query
6️⃣  Response flows back: DB → Repo → Handler → Client
```

### Example 2: WebSocket - Real-Time Move

```
Alice makes move e2→e4:

1️⃣  Alice sends WebSocket message (type: request)
2️⃣  readPump receives and processes
3️⃣  handleMessage validates and stores in DB
4️⃣  Server sends response to Alice (type: response)
5️⃣  Server broadcasts event to Bob (type: event)
6️⃣  Both players see updated board in real-time
```

---

## Concurrency & Goroutines

### Goroutine Architecture

```
MAIN GOROUTINE
├─ WebSocket Hub (1 goroutine)
├─ HTTP Server (1 goroutine)
└─ Per Client (2 goroutines each)
   ├─ readPump
   └─ writePump
```

### Channel Communication

```
Hub Channels:
├─ register (unbuffered)
├─ unregister (unbuffered)
└─ shutdown (unbuffered)

Client Channels:
└─ send (buffered: 256)
```

---

## Key Design Patterns

### 1. Layered Architecture

```
Handler → Repository → Database
```

### 2. Dependency Injection

```
main → app → server → router → handlers
```

### 3. Repository Pattern

```
Abstraction over data access
```

### 4. Hub-and-Spoke (WebSocket)

```
Central hub manages all clients
```

### 5. Read/Write Pumps (WebSocket)

```
Separate goroutines for read/write
```

---

## Troubleshooting & Common Issues

### Issue 1: Server Won't Start

```bash
# Port already in use
lsof -i :8080
kill -9 <PID>
```

### Issue 2: Database Connection Failed

```bash
# Check PostgreSQL
pg_isready -h localhost -p 5432
brew services start postgresql
```

### Issue 3: WebSocket Connection Fails

```javascript
// Use ws:// not http://
const ws = new WebSocket('ws://localhost:8080/ws?user_id=test');
```

### Issue 4: SQL Query Not Found

```bash
# Regenerate sqlc code
sqlc generate
```

---

## Summary

```
1️⃣  Bootstrap → Load config, connect DB
2️⃣  Initialize → Create hub, server, router
3️⃣  Runtime → Start goroutines
4️⃣  HTTP → Router → Handler → Repo → DB
5️⃣  WebSocket → Upgrade → Hub → Rooms → Clients
6️⃣  Shutdown → Graceful cleanup
```

---

**Happy coding! 🚀**

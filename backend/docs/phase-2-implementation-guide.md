# Phase 2 Implementation Guide

**Target Audience:** Backend Engineer
**Prerequisites:** Phase 1 complete (minimal HTTP server running)
**Estimated Time:** 2-4 weeks
**Difficulty:** Intermediate

---

## Table of Contents

1. [Overview](#overview)
2. [Architecture Goals](#architecture-goals)
3. [Implementation Roadmap](#implementation-roadmap)
4. [Detailed Instructions](#detailed-instructions)
5. [Testing Strategy](#testing-strategy)
6. [Troubleshooting](#troubleshooting)

---

## Overview

Phase 2 transforms the minimal HTTP server into a production-ready backend with:

- **Fiber v3** web framework
- **WebSocket** real-time communication
- **PostgreSQL** database with type-safe queries
- **Structured logging** and error handling
- **Testing infrastructure**

---

## Architecture Goals

### Current State (Phase 1)
```
[Client] → HTTP → [net/http server] → /health endpoint
```

### Target State (Phase 2)
```
[Frontend]
    ↓
[Fiber v3 + CORS Middleware]
    ↓
[WebSocket Hub] ←→ [PostgreSQL]
    ↓              ↓
[Game Manager] → [PGN Parser]
    ↓
[AI Coach (Future Phase 3)]
```

---

## Implementation Roadmap

### Week 1: Framework & Logging
- [ ] Day 1-2: Fiber v3 migration
- [ ] Day 3: Structured logging setup
- [ ] Day 4: CORS middleware
- [ ] Day 5: Environment configuration

### Week 2: WebSocket Infrastructure
- [ ] Day 1-2: Hub/Client pattern implementation
- [ ] Day 3: Connection management & heartbeat
- [ ] Day 4: Message routing
- [ ] Day 5: WebSocket testing

### Week 3: Database Layer
- [ ] Day 1: PostgreSQL Docker setup
- [ ] Day 2: sqlc installation & configuration
- [ ] Day 3: Schema design & migrations
- [ ] Day 4: Generated queries integration
- [ ] Day 5: Database testing

### Week 4: Game Logic
- [ ] Day 1-2: PGN parser implementation
- [ ] Day 3: Game state manager
- [ ] Day 4: WebSocket + Game integration
- [ ] Day 5: End-to-end testing

---

## Detailed Instructions

## Task 1: Fiber v3 Migration

### Objective
Replace `net/http` with Fiber v3 for better performance, middleware support, and Express.js-like API.

### Steps

#### 1.1 Add Fiber Dependency
```bash
cd backend
go get github.com/gofiber/fiber/v3
go mod tidy
```

#### 1.2 Create Server Package
```bash
mkdir -p internal/server
touch internal/server/server.go
```

**File: `internal/server/server.go`**
```go
package server

import (
    "log/slog"
    "github.com/gofiber/fiber/v3"
)

type Server struct {
    app    *fiber.App
    logger *slog.Logger
}

func New(logger *slog.Logger) *Server {
    app := fiber.New(fiber.Config{
        AppName:           "Chess Coach API",
        EnablePrintRoutes: true,
        ServerHeader:      "Chess Coach",
    })

    return &Server{
        app:    app,
        logger: logger,
    }
}

func (s *Server) SetupRoutes() {
    // Health check
    s.app.Get("/health", s.handleHealth)

    // API info
    s.app.Get("/", s.handleRoot)
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

func (s *Server) Start(port string) error {
    s.logger.Info("starting server", "port", port)
    return s.app.Listen(":" + port)
}

func (s *Server) Shutdown() error {
    s.logger.Info("shutting down server")
    return s.app.Shutdown()
}
```

#### 1.3 Update main.go
**File: `cmd/server/main.go`**
```go
package main

import (
    "log/slog"
    "os"
    "os/signal"
    "syscall"

    "chess-coach/backend/internal/server"
)

func main() {
    // Initialize structured logger
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
    slog.SetDefault(logger)

    // Get port from environment
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // Create server
    srv := server.New(logger)
    srv.SetupRoutes()

    // Graceful shutdown
    go func() {
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
        <-sigChan

        logger.Info("received shutdown signal")
        if err := srv.Shutdown(); err != nil {
            logger.Error("shutdown error", "error", err)
        }
    }()

    // Start server
    if err := srv.Start(port); err != nil {
        logger.Error("server error", "error", err)
        os.Exit(1)
    }
}
```

#### 1.4 Test the Migration
```bash
# Rebuild and restart
make restart

# Test endpoints
curl http://localhost:8080/health
curl http://localhost:8080/

# Check logs for Fiber startup message
make backend-logs
```

**Expected Output:**
```json
{"status":"ok","service":"chess-coach-backend"}
```

---

## Task 2: CORS Middleware

### Objective
Enable frontend (React) to communicate with backend from different origin.

### Steps

#### 2.1 Add CORS Package
```bash
go get github.com/gofiber/fiber/v3/middleware/cors
go mod tidy
```

#### 2.2 Create Middleware Package
```bash
mkdir -p internal/middleware
touch internal/middleware/cors.go
```

**File: `internal/middleware/cors.go`**
```go
package middleware

import (
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/cors"
)

func CORS() fiber.Handler {
    return cors.New(cors.Config{
        AllowOrigins: []string{
            "http://localhost:5173",  // Vite dev server
            "http://localhost:3000",  // Alternative dev port
        },
        AllowMethods: []string{
            fiber.MethodGet,
            fiber.MethodPost,
            fiber.MethodPut,
            fiber.MethodDelete,
            fiber.MethodOptions,
        },
        AllowHeaders: []string{
            "Origin",
            "Content-Type",
            "Accept",
            "Authorization",
        },
        AllowCredentials: true,
        MaxAge:           300, // 5 minutes
    })
}
```

#### 2.3 Apply Middleware
**Update: `internal/server/server.go`**
```go
import (
    "chess-coach/backend/internal/middleware"
    // ... other imports
)

func (s *Server) SetupRoutes() {
    // Middleware
    s.app.Use(middleware.CORS())

    // Routes
    s.app.Get("/health", s.handleHealth)
    s.app.Get("/", s.handleRoot)
}
```

#### 2.4 Test CORS
```bash
# Test with curl (simulate browser preflight)
curl -X OPTIONS http://localhost:8080/health \
  -H "Origin: http://localhost:5173" \
  -H "Access-Control-Request-Method: GET" \
  -v

# Check for CORS headers in response:
# Access-Control-Allow-Origin: http://localhost:5173
# Access-Control-Allow-Credentials: true
```

---

## Task 3: Environment Configuration

### Objective
Load environment variables from `.env` file using `godotenv`.

### Steps

#### 3.1 Add Dependency
```bash
go get github.com/joho/godotenv
go mod tidy
```

#### 3.2 Create Config Package
```bash
mkdir -p internal/config
touch internal/config/config.go
```

**File: `internal/config/config.go`**
```go
package config

import (
    "fmt"
    "log/slog"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Port        string
    Environment string
    DatabaseURL string
    RedisURL    string
}

func Load() (*Config, error) {
    // Load .env file (optional - won't fail if missing)
    if err := godotenv.Load(); err != nil {
        slog.Warn("no .env file found, using environment variables")
    }

    cfg := &Config{
        Port:        getEnv("PORT", "8080"),
        Environment: getEnv("APP_ENV", "development"),
        DatabaseURL: getEnv("DATABASE_URL", ""),
        RedisURL:    getEnv("REDIS_URL", ""),
    }

    if err := cfg.Validate(); err != nil {
        return nil, err
    }

    return cfg, nil
}

func (c *Config) Validate() error {
    // Add validation rules as needed
    if c.Port == "" {
        return fmt.Errorf("PORT is required")
    }
    return nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func (c *Config) IsDevelopment() bool {
    return c.Environment == "development"
}

func (c *Config) IsProduction() bool {
    return c.Environment == "production"
}
```

#### 3.3 Update main.go
**File: `cmd/server/main.go`**
```go
import (
    "chess-coach/backend/internal/config"
    "chess-coach/backend/internal/server"
    // ... other imports
)

func main() {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        slog.Error("failed to load config", "error", err)
        os.Exit(1)
    }

    // Initialize logger
    logLevel := slog.LevelInfo
    if cfg.IsDevelopment() {
        logLevel = slog.LevelDebug
    }

    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: logLevel,
    }))
    slog.SetDefault(logger)

    logger.Info("loaded configuration",
        "environment", cfg.Environment,
        "port", cfg.Port,
    )

    // Create and start server
    srv := server.New(logger)
    srv.SetupRoutes()

    // ... rest of main.go
}
```

#### 3.4 Create .env File
```bash
cd backend
cp .env.example .env
```

**Edit `.env`:**
```bash
PORT=8080
APP_ENV=development
DATABASE_URL=postgresql://postgres:password@postgres:5432/chess_coach?sslmode=disable
REDIS_URL=redis://redis:6379
```

---

## Task 4: WebSocket Hub/Client Implementation

### Objective
Implement WebSocket server with connection management, broadcasting, and heartbeat.

### Steps

#### 4.1 Add WebSocket Dependency
```bash
go get github.com/gorilla/websocket
go mod tidy
```

#### 4.2 Create WebSocket Package Structure
```bash
mkdir -p internal/websocket
touch internal/websocket/hub.go
touch internal/websocket/client.go
touch internal/websocket/message.go
```

#### 4.3 Define Message Types
**File: `internal/websocket/message.go`**
```go
package websocket

import "encoding/json"

type MessageType string

const (
    MessageTypePing     MessageType = "ping"
    MessageTypePong     MessageType = "pong"
    MessageTypeJoin     MessageType = "join"
    MessageTypeLeave    MessageType = "leave"
    MessageTypeMove     MessageType = "move"
    MessageTypeAnalysis MessageType = "analysis"
    MessageTypeError    MessageType = "error"
)

type Message struct {
    Type    MessageType     `json:"type"`
    Payload json.RawMessage `json:"payload,omitempty"`
}

type ErrorPayload struct {
    Error string `json:"error"`
}

type MovePayload struct {
    GameID string `json:"game_id"`
    Move   string `json:"move"` // UCI or SAN notation
}
```

#### 4.4 Implement Client
**File: `internal/websocket/client.go`**
```go
package websocket

import (
    "log/slog"
    "time"

    "github.com/gorilla/websocket"
)

const (
    writeWait      = 10 * time.Second
    pongWait       = 60 * time.Second
    pingPeriod     = (pongWait * 9) / 10
    maxMessageSize = 512 * 1024 // 512KB
)

type Client struct {
    hub  *Hub
    conn *websocket.Conn
    send chan []byte
    id   string
}

func NewClient(hub *Hub, conn *websocket.Conn, id string) *Client {
    return &Client{
        hub:  hub,
        conn: conn,
        send: make(chan []byte, 256),
        id:   id,
    }
}

// ReadPump reads messages from WebSocket connection
func (c *Client) ReadPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()

    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                slog.Error("websocket error", "client_id", c.id, "error", err)
            }
            break
        }

        // Broadcast to hub
        c.hub.broadcast <- message
    }
}

// WritePump writes messages to WebSocket connection
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
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            w, err := c.conn.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            w.Write(message)

            // Add queued messages
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
```

#### 4.5 Implement Hub
**File: `internal/websocket/hub.go`**
```go
package websocket

import (
    "log/slog"
)

type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
            slog.Info("client connected",
                "client_id", client.id,
                "total_clients", len(h.clients),
            )

        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
                slog.Info("client disconnected",
                    "client_id", client.id,
                    "total_clients", len(h.clients),
                )
            }

        case message := <-h.broadcast:
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
        }
    }
}

func (h *Hub) ClientCount() int {
    return len(h.clients)
}
```

#### 4.6 Add WebSocket Handler to Server
**Update: `internal/server/server.go`**
```go
import (
    "chess-coach/backend/internal/websocket"
    "github.com/gofiber/fiber/v3"
    gorillaWs "github.com/gorilla/websocket"
    "github.com/google/uuid"
)

type Server struct {
    app    *fiber.App
    logger *slog.Logger
    hub    *websocket.Hub
}

var upgrader = gorillaWs.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // Allow all origins in development
        // TODO: Restrict in production
        return true
    },
}

func New(logger *slog.Logger) *Server {
    hub := websocket.NewHub()

    app := fiber.New(fiber.Config{
        AppName:           "Chess Coach API",
        EnablePrintRoutes: true,
    })

    srv := &Server{
        app:    app,
        logger: logger,
        hub:    hub,
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
    s.app.Get("/ws", s.handleWebSocket)
}

func (s *Server) handleWebSocket(c fiber.Ctx) error {
    // Upgrade to WebSocket
    conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
    if err != nil {
        s.logger.Error("websocket upgrade failed", "error", err)
        return err
    }

    // Create client
    clientID := uuid.New().String()
    client := websocket.NewClient(s.hub, conn, clientID)

    s.hub.register <- client

    // Start client goroutines
    go client.WritePump()
    go client.ReadPump()

    return nil
}
```

#### 4.7 Test WebSocket
Create test file: `backend/test-websocket.html`
```html
<!DOCTYPE html>
<html>
<head>
    <title>WebSocket Test</title>
</head>
<body>
    <h1>Chess Coach WebSocket Test</h1>
    <div>
        <button onclick="connect()">Connect</button>
        <button onclick="disconnect()">Disconnect</button>
        <button onclick="sendPing()">Send Ping</button>
    </div>
    <div>
        <h3>Messages:</h3>
        <pre id="messages"></pre>
    </div>

    <script>
        let ws = null;
        const messages = document.getElementById('messages');

        function connect() {
            ws = new WebSocket('ws://localhost:8080/ws');

            ws.onopen = () => {
                log('Connected!');
            };

            ws.onmessage = (event) => {
                log('Received: ' + event.data);
            };

            ws.onclose = () => {
                log('Disconnected');
            };

            ws.onerror = (error) => {
                log('Error: ' + error);
            };
        }

        function disconnect() {
            if (ws) {
                ws.close();
            }
        }

        function sendPing() {
            if (ws && ws.readyState === WebSocket.OPEN) {
                const msg = JSON.stringify({ type: 'ping', payload: {} });
                ws.send(msg);
                log('Sent: ' + msg);
            }
        }

        function log(msg) {
            messages.textContent += new Date().toISOString() + ' - ' + msg + '\n';
        }
    </script>
</body>
</html>
```

Open in browser and test connection.

---

## Task 5: PostgreSQL Setup

### Objective
Add PostgreSQL database with migrations and type-safe queries using sqlc.

### Steps

#### 5.1 Update docker-compose.yml
**File: `docker-compose.yml` (at root)**
```yaml
services:
  backend:
    image: golang:1.23-alpine
    container_name: chess-coach-backend
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - APP_ENV=development
      - DATABASE_URL=postgresql://postgres:password@postgres:5432/chess_coach?sslmode=disable
    volumes:
      - ./backend:/app
      - /app/tmp
    command: |
      sh -c "
        go install github.com/air-verse/air@v1.52.3 &&
        air -c .air.toml
      "
    working_dir: /app
    restart: unless-stopped
    networks:
      - chess-coach-network
    depends_on:
      postgres:
        condition: service_healthy

  postgres:
    image: postgres:16-alpine
    container_name: chess-coach-postgres
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=password
      - POSTGRES_DB=chess_coach
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-PATCH", "pg_isready", "-U", "postgres"]
      interval: 5s
      timeout: 5s
      retries: 5
    networks:
      - chess-coach-network

networks:
  chess-coach-network:
    driver: bridge

volumes:
  postgres_data:
```

#### 5.2 Install sqlc
```bash
# On macOS
brew install sqlc

# Or download binary from https://github.com/sqlc-dev/sqlc
```

#### 5.3 Create Database Schema
```bash
mkdir -p backend/internal/db/migrations
touch backend/internal/db/migrations/001_initial_schema.sql
```

**File: `backend/internal/db/migrations/001_initial_schema.sql`**
```sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Games table
CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pgn TEXT NOT NULL,
    title VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, completed, abandoned
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Game moves table (for analysis)
CREATE TABLE moves (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    move_number INTEGER NOT NULL,
    move_san VARCHAR(10) NOT NULL, -- Standard Algebraic Notation
    move_uci VARCHAR(10) NOT NULL, -- Universal Chess Interface
    fen VARCHAR(255) NOT NULL, -- Position after move
    analysis JSONB, -- AI analysis result
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_games_user_id ON games(user_id);
CREATE INDEX idx_games_status ON games(status);
CREATE INDEX idx_moves_game_id ON moves(game_id);
CREATE INDEX idx_moves_game_move ON moves(game_id, move_number);
```

#### 5.4 Configure sqlc
**File: `backend/sqlc.yaml`**
```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "internal/db/queries"
    schema: "internal/db/migrations"
    gen:
      go:
        package: "db"
        out: "internal/db"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_interface: true
        emit_empty_slices: true
```

#### 5.5 Create Queries
```bash
mkdir -p backend/internal/db/queries
touch backend/internal/db/queries/users.sql
touch backend/internal/db/queries/games.sql
```

**File: `backend/internal/db/queries/users.sql`**
```sql
-- name: CreateUser :one
INSERT INTO users (email, username)
VALUES ($1, $2)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC;
```

**File: `backend/internal/db/queries/games.sql`**
```sql
-- name: CreateGame :one
INSERT INTO games (user_id, pgn, title)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetGame :one
SELECT * FROM games
WHERE id = $1 LIMIT 1;

-- name: ListUserGames :many
SELECT * FROM games
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateGamePGN :exec
UPDATE games
SET pgn = $2, updated_at = NOW()
WHERE id = $1;
```

#### 5.6 Generate Code
```bash
cd backend
sqlc generate
```

This creates type-safe Go code in `internal/db/`.

#### 5.7 Add Database Connection
```bash
go get github.com/jackc/pgx/v5
go mod tidy
```

**Create: `internal/db/db.go`**
```go
package db

import (
    "context"
    "fmt"

    "github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
    pool, err := pgxpool.New(ctx, databaseURL)
    if err != nil {
        return nil, fmt.Errorf("unable to create connection pool: %w", err)
    }

    // Test connection
    if err := pool.Ping(ctx); err != nil {
        return nil, fmt.Errorf("unable to ping database: %w", err)
    }

    return pool, nil
}
```

#### 5.8 Update main.go
```go
import (
    "context"
    "chess-coach/backend/internal/db"
)

func main() {
    // ... existing config and logger setup

    // Connect to database
    ctx := context.Background()
    pool, err := db.Connect(ctx, cfg.DatabaseURL)
    if err != nil {
        logger.Error("database connection failed", "error", err)
        os.Exit(1)
    }
    defer pool.Close()

    logger.Info("connected to database")

    // Create queries instance
    queries := db.New(pool)

    // Pass to server
    srv := server.New(logger, queries)

    // ... rest of main
}
```

#### 5.9 Run Migrations
For now, manually run the migration:
```bash
# Connect to postgres container
docker exec -it chess-coach-postgres psql -U postgres -d chess_coach

# Copy/paste the SQL from 001_initial_schema.sql
# Later: Use golang-migrate for automated migrations
```

---

## Task 6: Game Logic Foundation

### Objective
Implement PGN parser and basic game state management.

### Steps

#### 6.1 Add Chess Library
```bash
go get github.com/notnil/chess
go mod tidy
```

#### 6.2 Create Game Package
```bash
mkdir -p internal/game
touch internal/game/manager.go
touch internal/game/pgn.go
```

#### 6.3 Implement PGN Parser
**File: `internal/game/pgn.go`**
```go
package game

import (
    "fmt"
    "strings"

    "github.com/notnil/chess"
)

type PGNParser struct{}

func NewPGNParser() *PGNParser {
    return &PGNParser{}
}

func (p *PGNParser) Parse(pgnString string) (*chess.Game, error) {
    reader := strings.NewReader(pgnString)
    pgn, err := chess.PGN(reader)
    if err != nil {
        return nil, fmt.Errorf("invalid PGN: %w", err)
    }

    game := chess.NewGame(pgn)
    return game, nil
}

func (p *PGNParser) Validate(pgnString string) error {
    _, err := p.Parse(pgnString)
    return err
}

func (p *PGNParser) GetMoves(game *chess.Game) []string {
    moves := []string{}
    for _, move := range game.Moves() {
        moves = append(moves, move.String())
    }
    return moves
}

func (p *PGNParser) GetFEN(game *chess.Game) string {
    return game.Position().String()
}
```

#### 6.4 Implement Game Manager
**File: `internal/game/manager.go`**
```go
package game

import (
    "context"
    "fmt"

    "chess-coach/backend/internal/db"
    "github.com/google/uuid"
    "github.com/notnil/chess"
)

type Manager struct {
    queries *db.Queries
    parser  *PGNParser
}

func NewManager(queries *db.Queries) *Manager {
    return &Manager{
        queries: queries,
        parser:  NewPGNParser(),
    }
}

func (m *Manager) CreateGame(ctx context.Context, userID uuid.UUID, pgnString string) (*db.Game, error) {
    // Validate PGN
    if err := m.parser.Validate(pgnString); err != nil {
        return nil, fmt.Errorf("invalid PGN: %w", err)
    }

    // Create game in database
    game, err := m.queries.CreateGame(ctx, db.CreateGameParams{
        UserID: userID,
        Pgn:    pgnString,
        Title:  "New Game",
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create game: %w", err)
    }

    return &game, nil
}

func (m *Manager) GetGame(ctx context.Context, gameID uuid.UUID) (*chess.Game, error) {
    // Fetch from database
    dbGame, err := m.queries.GetGame(ctx, gameID)
    if err != nil {
        return nil, fmt.Errorf("game not found: %w", err)
    }

    // Parse PGN
    game, err := m.parser.Parse(dbGame.Pgn)
    if err != nil {
        return nil, fmt.Errorf("failed to parse game: %w", err)
    }

    return game, nil
}

func (m *Manager) UpdateGame(ctx context.Context, gameID uuid.UUID, pgnString string) error {
    // Validate PGN
    if err := m.parser.Validate(pgnString); err != nil {
        return fmt.Errorf("invalid PGN: %w", err)
    }

    // Update database
    err := m.queries.UpdateGamePGN(ctx, db.UpdateGamePGNParams{
        ID:  gameID,
        Pgn: pgnString,
    })
    if err != nil {
        return fmt.Errorf("failed to update game: %w", err)
    }

    return nil
}
```

---

## Testing Strategy

### Unit Tests

#### Test Structure
```
backend/
├── internal/
│   ├── game/
│   │   ├── manager.go
│   │   ├── manager_test.go
│   │   ├── pgn.go
│   │   └── pgn_test.go
```

#### Install Testing Dependencies
```bash
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/require
go mod tidy
```

#### Example: PGN Parser Test
**File: `internal/game/pgn_test.go`**
```go
package game

import (
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestPGNParser_Parse(t *testing.T) {
    parser := NewPGNParser()

    t.Run("valid PGN", func(t *testing.T) {
        pgn := `1. e4 e5 2. Nf3 Nc6`
        game, err := parser.Parse(pgn)

        require.NoError(t, err)
        assert.NotNil(t, game)

        moves := parser.GetMoves(game)
        assert.Equal(t, 4, len(moves))
    })

    t.Run("invalid PGN", func(t *testing.T) {
        pgn := `1. e4 e5 2. Zz9`
        _, err := parser.Parse(pgn)

        assert.Error(t, err)
    })
}
```

#### Run Tests
```bash
cd backend
go test ./...

# With coverage
go test -cover ./...

# Verbose
go test -v ./...
```

### Integration Tests

#### Example: WebSocket Test
**File: `internal/websocket/hub_test.go`**
```go
package websocket

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestHub_ClientRegistration(t *testing.T) {
    hub := NewHub()
    go hub.Run()

    // Create mock client
    client := &Client{
        hub:  hub,
        send: make(chan []byte),
        id:   "test-client",
    }

    // Register
    hub.register <- client
    time.Sleep(100 * time.Millisecond)

    assert.Equal(t, 1, hub.ClientCount())

    // Unregister
    hub.unregister <- client
    time.Sleep(100 * time.Millisecond)

    assert.Equal(t, 0, hub.ClientCount())
}
```

---

## Troubleshooting

### Database Connection Issues

**Problem:** `unable to connect to database`

**Solutions:**
```bash
# Check postgres is running
docker ps | grep postgres

# Check postgres logs
docker logs chess-coach-postgres

# Test connection manually
docker exec -it chess-coach-postgres psql -U postgres -d chess_coach

# Verify DATABASE_URL in .env
cat backend/.env | grep DATABASE_URL
```

### WebSocket Connection Refused

**Problem:** WebSocket upgrade fails

**Solutions:**
```bash
# Check CORS headers
curl -X OPTIONS http://localhost:8080/ws \
  -H "Origin: http://localhost:5173" -v

# Check upgrader CheckOrigin function
# Should return true in development

# Verify WebSocket endpoint
curl http://localhost:8080/ws
# Should return: Bad Request (expected for HTTP, not WS)
```

### sqlc Generation Errors

**Problem:** `sqlc generate` fails

**Solutions:**
```bash
# Verify sqlc.yaml syntax
cat backend/sqlc.yaml

# Check SQL syntax in queries/
# Common issue: Missing semicolons

# Ensure schema files exist
ls backend/internal/db/migrations/

# Update sqlc
brew upgrade sqlc
```

### Hot Reload Not Working

**Problem:** Changes not reflected

**Solutions:**
```bash
# Check Air is running
make backend-logs | grep air

# Force rebuild
make restart

# Check .air.toml config
cat backend/.air.toml

# Ensure volume mount is correct
docker inspect chess-coach-backend | grep Mounts
```

---

## Validation Checklist

After completing Phase 2, verify:

- [ ] Fiber server starts without errors
- [ ] `/health` endpoint returns 200
- [ ] CORS headers present in responses
- [ ] Environment variables loaded from `.env`
- [ ] WebSocket connection successful
- [ ] WebSocket echo working
- [ ] PostgreSQL container healthy
- [ ] Database migrations applied
- [ ] sqlc code generated
- [ ] Database queries execute
- [ ] PGN parser validates games
- [ ] Unit tests pass
- [ ] Integration tests pass

### Test Commands
```bash
# Health check
curl http://localhost:8080/health

# CORS check
curl -I http://localhost:8080/health

# Database check
docker exec -it chess-coach-postgres psql -U postgres -d chess_coach -c "SELECT COUNT(*) FROM users;"

# WebSocket check (open test-websocket.html in browser)

# Run tests
cd backend && go test ./...
```

---

## Next Steps (Phase 3 Preview)

After Phase 2 completion:

1. **AI Coach Integration**
   - Claude API integration
   - Move analysis pipeline
   - Response streaming

2. **Authentication**
   - JWT implementation
   - User registration/login
   - Protected routes

3. **Redis Caching**
   - Position caching
   - Rate limiting
   - Session storage

4. **Production Deployment**
   - AWS ECS setup
   - RDS PostgreSQL
   - Application Load Balancer
   - CloudWatch monitoring

---

**Document Version:** 1.0
**Last Updated:** November 1, 2025
**Author:** Chess Coach Team
**Status:** Ready for Implementation

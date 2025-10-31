# Chess Coach Backend: Go Tech Stack (2025)

**Language:** Go 1.23+
**Team Size:** 2-3 developers
**Deployment:** AWS ECS Fargate
**Priority:** Production-ready, type-safe, high-performance

---

## Executive Summary

This is our **finalized Go tech stack** for Chess Coach backend, based on 2025 ecosystem research, benchmarks, and production best practices.

### Core Stack

| Component | Library | Why This Choice |
|-----------|---------|-----------------|
| **Web Framework** | Fiber v3 | Express-like API, 36K req/s, fasthttp-based |
| **WebSocket** | gorilla/websocket | Battle-tested, 1M+ connections proven |
| **Database** | sqlc | Type-safe SQL, zero runtime overhead |
| **Cache** | go-redis v9 | Official Redis client, pipeline support |
| **Message Queue** | nats.go | Official NATS client, JetStream support |
| **Logging** | slog (stdlib) | Zero deps, structured, AWS CloudWatch friendly |
| **Validation** | go-playground/validator v10 | Most popular, supports custom validators |
| **Config** | viper | Environment + file config, AWS Secrets Manager |
| **Testing** | testify + httptest | Standard testing with assertions |
| **AWS SDK** | aws-sdk-go-v2 | Official v2 SDK, modular imports |

---

## 1. Web Framework: Fiber v3

### Why Fiber?

**Performance (2025 Benchmarks):**
- **36K requests/second** (vs Gin 34K, Echo 34K)
- **2.8ms median latency** (vs Gin/Echo 3ms)
- Built on `fasthttp` (fastest HTTP engine for Go)

**Developer Experience:**
- Express.js-like API (familiar for JS developers)
- Excellent documentation
- Active maintenance (v3 released 2025)

**Real-World Fit:**
- Perfect for WebSocket upgrades
- Low memory footprint
- Middleware ecosystem

### Installation

```bash
go get -u github.com/gofiber/fiber/v3
go get -u github.com/gofiber/websocket/v2
```

### Sample Code

```go
package main

import (
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/websocket/v2"
)

func main() {
    app := fiber.New(fiber.Config{
        AppName:      "Chess Coach API",
        ServerHeader: "Fiber",
        // Enable body limit for uploads
        BodyLimit: 10 * 1024 * 1024, // 10MB
    })

    // Middleware
    app.Use(logger.New())
    app.Use(cors.New())
    app.Use(recover.New())

    // Health check (for ALB)
    app.Get("/health", func(c fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "healthy",
            "time":   time.Now().Unix(),
        })
    })

    // WebSocket upgrade
    app.Use("/ws", func(c fiber.Ctx) error {
        if websocket.IsWebSocketUpgrade(c) {
            return c.Next()
        }
        return fiber.ErrUpgradeRequired
    })

    // WebSocket endpoint
    app.Get("/ws", websocket.New(handleWebSocket))

    // REST API
    api := app.Group("/api/v1")
    api.Post("/games", createGame)
    api.Get("/games/:id", getGame)

    app.Listen(":3000")
}
```

### Why Not Gin or Echo?

**Gin:**
- Slightly slower (34K req/s vs 36K)
- Uses `net/http` (slower than fasthttp)
- ✅ Still excellent, but Fiber edges ahead

**Echo:**
- Also excellent (34K req/s)
- More opinionated structure
- ✅ Good alternative if you prefer stronger conventions

**Verdict:** Fiber for raw performance + Express-like ergonomics

---

## 2. WebSocket: gorilla/websocket

### Why gorilla/websocket?

**Production Proven:**
- 6+ years battle-tested
- Used by Discord, Twitch, many others
- Handles 1M+ concurrent connections

**vs nhooyr.io/websocket (coder/websocket):**
- nhooyr is "more idiomatic" but less proven
- gorilla has richer API, better docs
- Performance difference negligible in real-world

**Consensus (2025):**
> "Gorilla WebSocket is best for most production use-cases, rich API, extensive documentation. A very mature library."
> — Go Forum 2025

### Installation

```bash
go get -u github.com/gorilla/websocket
```

### Sample WebSocket Hub Pattern

```go
// internal/websocket/hub.go
package websocket

import (
    "sync"
    "github.com/gorilla/websocket"
)

type Hub struct {
    // Registered clients (userId -> connection)
    clients    map[string]*websocket.Conn
    clientsMux sync.RWMutex

    // Broadcast channel
    broadcast chan Message

    // Register/Unregister
    register   chan *Client
    unregister chan *Client
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[string]*websocket.Conn),
        broadcast:  make(chan Message, 256),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clientsMux.Lock()
            h.clients[client.UserID] = client.Conn
            h.clientsMux.Unlock()

        case client := <-h.unregister:
            h.clientsMux.Lock()
            if _, ok := h.clients[client.UserID]; ok {
                delete(h.clients, client.UserID)
                client.Conn.Close()
            }
            h.clientsMux.Unlock()

        case message := <-h.broadcast:
            h.clientsMux.RLock()
            for _, conn := range h.clients {
                go conn.WriteJSON(message)
            }
            h.clientsMux.RUnlock()
        }
    }
}

func (h *Hub) BroadcastToGame(gameID string, message Message) {
    // Broadcast to all clients in a specific game
    h.clientsMux.RLock()
    defer h.clientsMux.RUnlock()

    for userID, conn := range h.clients {
        // Check if user is in this game (via Redis or in-memory map)
        if h.isUserInGame(userID, gameID) {
            go conn.WriteJSON(message)
        }
    }
}
```

### Connection Handling

```go
// internal/websocket/client.go
type Client struct {
    UserID string
    GameID string
    Conn   *websocket.Conn
    Hub    *Hub
}

func (c *Client) ReadPump() {
    defer func() {
        c.Hub.unregister <- c
        c.Conn.Close()
    }()

    // Set read deadline
    c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
    c.Conn.SetPongHandler(func(string) error {
        c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        return nil
    })

    for {
        var msg Message
        err := c.Conn.ReadJSON(&msg)
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("error: %v", err)
            }
            break
        }

        // Handle message (publish to NATS, etc.)
        c.HandleMessage(msg)
    }
}

func (c *Client) WritePump() {
    ticker := time.NewTicker(54 * time.Second)
    defer func() {
        ticker.Stop()
        c.Conn.Close()
    }()

    for {
        select {
        case <-ticker.C:
            c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
```

---

## 3. Database: sqlc (Type-Safe SQL)

### Why sqlc over GORM/Ent?

**Performance (2025 Benchmarks):**
- GORM: 59.3ms for 15K rows (2x slower, uses reflection)
- Ent: Fast, but complex for simple queries
- **sqlc: 31.7ms** (essentially raw SQL performance)

**Type Safety:**
- **Compile-time safety** (errors at build, not runtime)
- Generates Go structs from SQL schemas
- No runtime reflection overhead

**Philosophy:**
> "You write SQL, sqlc generates type-safe Go code"

### Installation

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

### Configuration

```yaml
# sqlc.yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "internal/db/queries/"
    schema: "internal/db/schema/"
    gen:
      go:
        package: "db"
        out: "internal/db"
        emit_json_tags: true
        emit_db_tags: true
        emit_interface: true
        emit_empty_slices: true
```

### Schema Definition

```sql
-- internal/db/schema/001_games.sql
CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pgn TEXT NOT NULL,
    white_player VARCHAR(255) NOT NULL,
    black_player VARCHAR(255) NOT NULL,
    result VARCHAR(10),
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_games_created_at ON games(created_at DESC);
CREATE INDEX idx_games_players ON games(white_player, black_player);
```

### Queries Definition

```sql
-- internal/db/queries/games.sql

-- name: CreateGame :one
INSERT INTO games (
    pgn, white_player, black_player, result, metadata
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetGame :one
SELECT * FROM games
WHERE id = $1 LIMIT 1;

-- name: ListRecentGames :many
SELECT * FROM games
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetUserGames :many
SELECT * FROM games
WHERE white_player = $1 OR black_player = $1
ORDER BY created_at DESC;
```

### Generated Code Usage

```go
package main

import (
    "context"
    "database/sql"
    "github.com/yourteam/chess-coach/internal/db"
    _ "github.com/lib/pq"
)

func main() {
    conn, _ := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    queries := db.New(conn)

    // Type-safe queries!
    game, err := queries.CreateGame(context.Background(), db.CreateGameParams{
        Pgn:         "1. e4 e5 2. Nf3",
        WhitePlayer: "Alice",
        BlackPlayer: "Bob",
        Result:      sql.NullString{String: "1-0", Valid: true},
        Metadata:    json.RawMessage(`{"time_control": "10+0"}`),
    })

    // All fields are type-safe
    fmt.Printf("Created game: %s\n", game.ID)
}
```

### Why Not GORM?

**GORM Issues:**
- Runtime reflection (slow, not type-safe)
- 2x slower in benchmarks
- Hidden query complexity
- Magic behavior (auto-migrations can be dangerous)

**When to use GORM:**
- Rapid prototyping
- Small apps (<1K users)
- Team unfamiliar with SQL

**For Chess Coach: sqlc is better** (performance + type safety)

---

## 4. Cache: go-redis v9

### Why go-redis?

**Official Client:**
- Maintained by Redis team
- Supports Redis 7+ features
- Excellent performance

**Features for Chess Coach:**
- Pub/Sub (for game room broadcasts)
- Pipelining (batch operations)
- Lua scripting
- Cluster support

### Installation

```bash
go get -u github.com/redis/go-redis/v9
```

### Usage

```go
package cache

import (
    "context"
    "github.com/redis/go-redis/v9"
)

type Client struct {
    rdb *redis.Client
}

func NewClient(addr string) *Client {
    rdb := redis.NewClient(&redis.Options{
        Addr:         addr,
        Password:     "", // from env
        DB:           0,
        PoolSize:     100,
        MinIdleConns: 10,
    })

    return &Client{rdb: rdb}
}

// Store active game state
func (c *Client) SetGameState(ctx context.Context, gameID string, state GameState) error {
    data, _ := json.Marshal(state)
    return c.rdb.Set(ctx, "game:"+gameID, data, 1*time.Hour).Err()
}

// Pub/Sub for game room
func (c *Client) PublishMove(ctx context.Context, gameID string, move Move) error {
    data, _ := json.Marshal(move)
    return c.rdb.Publish(ctx, "game:"+gameID, data).Err()
}

func (c *Client) SubscribeToGame(ctx context.Context, gameID string) *redis.PubSub {
    return c.rdb.Subscribe(ctx, "game:"+gameID)
}

// Connection mapping (userId -> ECS task ID)
func (c *Client) MapUserToTask(ctx context.Context, userID, taskID string) error {
    return c.rdb.HSet(ctx, "ws:connections", userID, taskID).Err()
}

func (c *Client) GetUserTask(ctx context.Context, userID string) (string, error) {
    return c.rdb.HGet(ctx, "ws:connections", userID).Result()
}
```

---

## 5. Message Queue: nats.go

### Why NATS?

**For Chess Coach:**
- Sub-10ms latency
- Built-in JetStream (persistence)
- Lightweight (single binary)
- Perfect for event-driven architecture

### Installation

```bash
go get -u github.com/nats-io/nats.go
```

### Usage

```go
package messaging

import (
    "github.com/nats-io/nats.go"
)

type NATSClient struct {
    nc *nats.Conn
    js nats.JetStreamContext
}

func NewNATSClient(url string) (*NATSClient, error) {
    nc, err := nats.Connect(url,
        nats.RetryOnFailedConnect(true),
        nats.MaxReconnects(10),
    )
    if err != nil {
        return nil, err
    }

    js, err := nc.JetStream()
    if err != nil {
        return nil, err
    }

    return &NATSClient{nc: nc, js: js}, nil
}

// Publish move to subject
func (n *NATSClient) PublishMove(gameID string, move Move) error {
    data, _ := json.Marshal(move)
    _, err := n.js.Publish("chess.moves."+gameID, data)
    return err
}

// Subscribe to moves
func (n *NATSClient) SubscribeMoves(gameID string, handler func(Move)) error {
    _, err := n.js.Subscribe("chess.moves."+gameID, func(msg *nats.Msg) {
        var move Move
        json.Unmarshal(msg.Data, &move)
        handler(move)
        msg.Ack()
    }, nats.Durable("moves-consumer"))
    return err
}

// AI request/response pattern
func (n *NATSClient) RequestAI(ctx context.Context, prompt string) (string, error) {
    data, _ := json.Marshal(map[string]string{"prompt": prompt})

    msg, err := n.nc.RequestWithContext(ctx, "chess.ai.query", data)
    if err != nil {
        return "", err
    }

    var response map[string]string
    json.Unmarshal(msg.Data, &response)
    return response["text"], nil
}
```

---

## 6. Logging: slog (Standard Library)

### Why slog over zerolog/zap?

**2025 Consensus:**
- **slog is now in stdlib** (Go 1.21+)
- Zero external dependencies
- Efficient memory usage (40 B/op)
- AWS CloudWatch friendly (structured JSON)

**Performance:**
- Zerolog: Fastest (but only slightly)
- **slog: Best memory efficiency**
- zap: Fastest but most memory (168 B/op)

**For small teams:** stdlib > external deps

### Usage

```go
package main

import (
    "log/slog"
    "os"
)

func main() {
    // JSON handler for production (CloudWatch)
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
        AddSource: true, // Add file:line
    }))

    slog.SetDefault(logger)

    // Structured logging
    slog.Info("Server started",
        "port", 3000,
        "env", os.Getenv("ENV"),
    )

    slog.Error("Failed to connect to database",
        "error", err,
        "host", dbHost,
        "attempt", 3,
    )

    // Context-aware logging
    ctx := context.WithValue(context.Background(), "requestId", "abc123")
    logger.InfoContext(ctx, "Processing request",
        "userId", user.ID,
        "endpoint", "/api/games",
    )
}
```

**Output (CloudWatch-ready JSON):**
```json
{
  "time": "2025-10-31T17:30:00Z",
  "level": "INFO",
  "msg": "Server started",
  "port": 3000,
  "env": "production",
  "source": "main.go:25"
}
```

---

## 7. Validation: go-playground/validator

### Why validator?

**Most popular:**
- 14K+ stars on GitHub
- Used by Gin, Echo internally
- Comprehensive tag-based validation

### Installation

```bash
go get -u github.com/go-playground/validator/v10
```

### Usage

```go
package main

import (
    "github.com/go-playground/validator/v10"
)

type MoveRequest struct {
    From   string `json:"from" validate:"required,len=2,alphanum"`
    To     string `json:"to" validate:"required,len=2,alphanum"`
    GameID string `json:"gameId" validate:"required,uuid"`
}

var validate = validator.New()

func handleMove(c fiber.Ctx) error {
    var req MoveRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
    }

    // Validate
    if err := validate.Struct(req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "Validation failed",
            "details": err.Error(),
        })
    }

    // Process move...
    return c.JSON(fiber.Map{"success": true})
}

// Custom validators
func init() {
    validate.RegisterValidation("chess_square", func(fl validator.FieldLevel) bool {
        val := fl.Field().String()
        return regexp.MustCompile(`^[a-h][1-8]$`).MatchString(val)
    })
}

type Move struct {
    From string `validate:"required,chess_square"`
    To   string `validate:"required,chess_square"`
}
```

---

## 8. Configuration: viper

### Why viper?

**Features:**
- Reads from env, files, AWS Secrets Manager
- Hot-reloading (useful for dev)
- Type-safe getters

### Installation

```bash
go get -u github.com/spf13/viper
```

### Usage

```go
package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    Port        int
    DatabaseURL string
    RedisURL    string
    NATSURL     string
    AnthropicKey string
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("./config")

    // Environment variables take precedence
    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        // Config file not required if all env vars set
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, err
        }
    }

    return &Config{
        Port:         viper.GetInt("PORT"),
        DatabaseURL:  viper.GetString("DATABASE_URL"),
        RedisURL:     viper.GetString("REDIS_URL"),
        NATSURL:      viper.GetString("NATS_URL"),
        AnthropicKey: viper.GetString("ANTHROPIC_API_KEY"),
    }, nil
}
```

**config.yaml:**
```yaml
# Development defaults
port: 3000
database_url: "postgresql://dev:dev@localhost:5432/chess_coach"
redis_url: "redis://localhost:6379"
nats_url: "nats://localhost:4222"
```

**Production (env vars only):**
```bash
export PORT=3000
export DATABASE_URL="postgresql://..."
export ANTHROPIC_API_KEY="sk-ant-..."
```

---

## 9. Testing: testify + httptest

### Why testify?

**Assertions:**
- Cleaner than raw if/else
- Better error messages

### Installation

```bash
go get -u github.com/stretchr/testify
```

### Usage

```go
package game

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestCreateGame(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    defer db.Close()

    // Execute
    game, err := CreateGame(db, "Alice", "Bob")

    // Assert
    require.NoError(t, err)
    assert.NotEmpty(t, game.ID)
    assert.Equal(t, "Alice", game.WhitePlayer)
    assert.Equal(t, "Bob", game.BlackPlayer)
}

// HTTP endpoint testing
func TestCreateGameEndpoint(t *testing.T) {
    app := setupTestApp()

    req := httptest.NewRequest("POST", "/api/games", strings.NewReader(`{
        "whitePlayer": "Alice",
        "blackPlayer": "Bob"
    }`))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)
    require.NoError(t, err)
    assert.Equal(t, 201, resp.StatusCode)

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    assert.NotEmpty(t, result["id"])
}
```

---

## 10. AWS SDK: aws-sdk-go-v2

### Why v2?

**Modern:**
- Modular imports (only import what you need)
- Better performance
- Context-aware

### Installation

```bash
go get -u github.com/aws/aws-sdk-go-v2/config
go get -u github.com/aws/aws-sdk-go-v2/service/s3
go get -u github.com/aws/aws-sdk-go-v2/service/secretsmanager
```

### Usage

```go
package aws

import (
    "context"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type AWSClient struct {
    s3      *s3.Client
    secrets *secretsmanager.Client
}

func NewAWSClient(ctx context.Context) (*AWSClient, error) {
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        return nil, err
    }

    return &AWSClient{
        s3:      s3.NewFromConfig(cfg),
        secrets: secretsmanager.NewFromConfig(cfg),
    }, nil
}

// Upload PGN to S3
func (a *AWSClient) UploadPGN(ctx context.Context, gameID, pgn string) error {
    _, err := a.s3.PutObject(ctx, &s3.PutObjectInput{
        Bucket: aws.String("chess-coach-pgn-archives"),
        Key:    aws.String(fmt.Sprintf("games/%s.pgn", gameID)),
        Body:   strings.NewReader(pgn),
    })
    return err
}

// Get secret from Secrets Manager
func (a *AWSClient) GetSecret(ctx context.Context, name string) (string, error) {
    result, err := a.secrets.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
        SecretId: aws.String(name),
    })
    if err != nil {
        return "", err
    }
    return *result.SecretString, nil
}
```

---

## Complete Dependency List

### go.mod

```go
module github.com/yourteam/chess-coach

go 1.23

require (
    // Web framework
    github.com/gofiber/fiber/v3 v3.0.0
    github.com/gofiber/websocket/v2 v2.2.1

    // WebSocket
    github.com/gorilla/websocket v1.5.1

    // Database
    github.com/lib/pq v1.10.9
    github.com/sqlc-dev/sqlc v1.26.0

    // Cache
    github.com/redis/go-redis/v9 v9.5.0

    // Message queue
    github.com/nats-io/nats.go v1.34.0

    // Validation
    github.com/go-playground/validator/v10 v10.19.0

    // Config
    github.com/spf13/viper v1.18.2

    // Testing
    github.com/stretchr/testify v1.9.0

    // AWS
    github.com/aws/aws-sdk-go-v2 v1.26.0
    github.com/aws/aws-sdk-go-v2/config v1.27.0
    github.com/aws/aws-sdk-go-v2/service/s3 v1.51.0
    github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.28.0

    // Utilities
    github.com/google/uuid v1.6.0
    github.com/joho/godotenv v1.5.1
)
```

---

## Project Structure

```
chess-coach-backend/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point
├── internal/
│   ├── websocket/
│   │   ├── hub.go                  # Connection manager
│   │   ├── client.go               # Per-client handler
│   │   └── handler.go              # Message handlers
│   ├── game/
│   │   ├── manager.go              # Game state management
│   │   ├── moves.go                # Move validation
│   │   └── stockfish.go            # Engine integration
│   ├── ai/
│   │   ├── claude.go               # Claude API client
│   │   └── streaming.go            # Stream handling
│   ├── db/
│   │   ├── queries/                # SQL queries (sqlc input)
│   │   ├── schema/                 # SQL schemas
│   │   ├── db.go                   # Generated code (sqlc output)
│   │   └── migrations/             # DB migrations
│   ├── cache/
│   │   └── redis.go                # Redis client
│   ├── messaging/
│   │   └── nats.go                 # NATS client
│   ├── api/
│   │   ├── handlers/               # HTTP handlers
│   │   └── middleware/             # Custom middleware
│   └── config/
│       └── config.go               # Configuration
├── pkg/
│   └── types/
│       └── messages.go             # Shared types
├── tests/
│   ├── integration/
│   └── unit/
├── scripts/
│   ├── migrate.sh                  # Run migrations
│   └── seed.sh                     # Seed data
├── config/
│   ├── config.yaml                 # Local config
│   └── config.production.yaml      # Production config
├── docker-compose.yml
├── Dockerfile
├── .air.toml                       # Hot reload config
├── sqlc.yaml                       # sqlc config
├── go.mod
├── go.sum
└── README.md
```

---

## Development Workflow

### Initial Setup

```bash
# 1. Install Go 1.23+
brew install go

# 2. Clone and init
git clone <repo>
cd chess-coach-backend
go mod download

# 3. Install tools
go install github.com/cosmtrek/air@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# 4. Generate code
sqlc generate

# 5. Start infrastructure
docker-compose up -d

# 6. Run migrations
./scripts/migrate.sh up

# 7. Start server (hot reload)
air
```

### Daily Development

```bash
# Terminal 1: Infrastructure
docker-compose up

# Terminal 2: Backend (hot reload)
air

# Backend running at http://localhost:3000
# WebSocket at ws://localhost:3000/ws
```

**Hot reload time:** ~100ms (instant feedback)

---

## Performance Targets

### Benchmarks (Expected)

| Metric | Target | How We Achieve |
|--------|--------|----------------|
| **HTTP Req/Sec** | 30K+ | Fiber + fasthttp |
| **WebSocket Connections** | 100K+ | gorilla/websocket + goroutines |
| **Message Latency** | <10ms | NATS JetStream |
| **Database Query** | <5ms | sqlc + connection pooling |
| **Memory (10K users)** | <500MB | Go's efficiency |
| **Cold Start** | <100ms | Single binary |

---

## AWS Deployment

### Dockerfile

```dockerfile
# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 3000
CMD ["./server"]
```

**Image size:** ~15MB

### GitHub Actions (CI/CD)

```yaml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.23'

      - name: Run tests
        run: go test ./...

      - name: Build and push
        run: |
          docker build -t chess-coach .
          docker push $ECR_REGISTRY/chess-coach:latest

      - name: Deploy to ECS
        run: aws ecs update-service --force-new-deployment
```

---

## Summary: Final Stack

```
┌─────────────────────────────────────────────┐
│  Language: Go 1.23+                         │
│  Framework: Fiber v3                        │
│  WebSocket: gorilla/websocket               │
│  Database: PostgreSQL + sqlc                │
│  Cache: Redis (go-redis v9)                 │
│  Queue: NATS (nats.go)                      │
│  Logging: slog (stdlib)                     │
│  Validation: validator v10                  │
│  Config: viper                              │
│  Testing: testify                           │
│  AWS: aws-sdk-go-v2                         │
└─────────────────────────────────────────────┘
```

**Next Step:** Create initial project structure with working examples?

---

**Document Version:** 1.0
**Last Updated:** October 31, 2025
**Status:** Finalized and ready for implementation

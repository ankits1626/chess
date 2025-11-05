# Step 15: Docker + Database Setup (Go Native Way)

## Goal
Dockerize Chess Coach backend with PostgreSQL using Go best practices (pgx + sqlc).

---

## Tech Stack Decision (2025 Best Practices)

✅ **PostgreSQL** - Database
✅ **pgx/v5** - PostgreSQL driver (fastest)
✅ **sqlc** - SQL → Go code generation (type-safe)
✅ **golang-migrate** - Database migrations
✅ **ozzo-validation** - Request validation (code-based)
✅ **Docker Compose** - Container orchestration

**No GORM, following industry standard!**

---

## Architecture Overview

```
chess-coach/
├── backend/
│   ├── cmd/server/main.go
│   ├── internal/
│   │   ├── config/          # Configuration
│   │   ├── database/        # DB connection, sqlc generated code
│   │   ├── repository/      # Data access layer
│   │   └── router/
│   ├── db/
│   │   ├── migrations/      # SQL migrations
│   │   ├── queries/         # SQL queries (sqlc input)
│   │   └── schema.sql       # Database schema
│   ├── sqlc.yaml           # sqlc config
│   ├── Dockerfile
│   ├── Dockerfile.dev
│   ├── docker compose.yml
│   └── Makefile
└── database/               # Shared (optional)
    └── init/
        └── 01-init.sql
```

---

## Implementation Steps

### Phase 1: Database Schema & Folder Structure

#### Step 1.1: Create Database Folders
**Time:** 2 minutes

```bash
cd /Users/ankit/code/learn/chess-coach/backend

# Create database directories
mkdir -p db/migrations
mkdir -p db/queries
mkdir -p internal/database
mkdir -p internal/repository

# Create root database folder (shared)
cd ..
mkdir -p database/init
```

---

#### Step 1.2: Define Database Schema
**Time:** 10 minutes

**File:** `backend/db/schema.sql`

```sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    rating INTEGER DEFAULT 1200,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Games table
CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    white_player_id UUID REFERENCES users(id),
    black_player_id UUID REFERENCES users(id),
    pgn TEXT NOT NULL,
    result VARCHAR(10), -- 1-0, 0-1, 1/2-1/2, *
    time_control INTEGER, -- in seconds
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT valid_result CHECK (result IN ('1-0', '0-1', '1/2-1/2', '*'))
);

-- Moves table (for detailed analysis)
CREATE TABLE moves (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    move_number INTEGER NOT NULL,
    side VARCHAR(5) NOT NULL, -- 'white' or 'black'
    move_san VARCHAR(20) NOT NULL, -- Standard Algebraic Notation (e.g., 'e4', 'Nf3')
    move_uci VARCHAR(10) NOT NULL, -- UCI format (e.g., 'e2e4')
    fen TEXT NOT NULL, -- Board position after move
    time_taken INTEGER, -- milliseconds

    CONSTRAINT valid_side CHECK (side IN ('white', 'black'))
);

-- Indexes for performance
CREATE INDEX idx_games_white_player ON games(white_player_id);
CREATE INDEX idx_games_black_player ON games(black_player_id);
CREATE INDEX idx_moves_game_id ON moves(game_id);
CREATE INDEX idx_moves_game_move ON moves(game_id, move_number);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);

-- Updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply triggers
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_games_updated_at BEFORE UPDATE ON games
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

---

#### Step 1.3: Create Initial Migration
**Time:** 5 minutes

**File:** `backend/db/migrations/000001_init_schema.up.sql`

```sql
-- Copy content from schema.sql
-- This is the "up" migration (apply changes)
```

**File:** `backend/db/migrations/000001_init_schema.down.sql`

```sql
-- Rollback migration (undo changes)
DROP TRIGGER IF EXISTS update_games_updated_at ON games;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();

DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_moves_game_move;
DROP INDEX IF EXISTS idx_moves_game_id;
DROP INDEX IF EXISTS idx_games_black_player;
DROP INDEX IF EXISTS idx_games_white_player;

DROP TABLE IF EXISTS moves;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "uuid-ossp";
```

---

### Phase 2: sqlc Configuration & Queries

#### Step 2.1: Create sqlc Config
**Time:** 3 minutes

**File:** `backend/sqlc.yaml`

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "db/queries"
    schema: "db/schema.sql"
    gen:
      go:
        package: "database"
        out: "internal/database"
        sql_package: "pgx/v5"
        emit_json_tags: true
        emit_prepared_queries: false
        emit_interface: true
        emit_exact_table_names: false
        emit_empty_slices: true
```

---

#### Step 2.2: Write SQL Queries
**Time:** 15 minutes

**File:** `backend/db/queries/users.sql`

```sql
-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CreateUser :one
INSERT INTO users (username, email, rating)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET username = $2, email = $3, rating = $4
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
```

**File:** `backend/db/queries/games.sql`

```sql
-- name: GetGame :one
SELECT * FROM games
WHERE id = $1 LIMIT 1;

-- name: ListGames :many
SELECT * FROM games
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListUserGames :many
SELECT * FROM games
WHERE white_player_id = $1 OR black_player_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateGame :one
INSERT INTO games (white_player_id, black_player_id, pgn, result, time_control)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateGame :one
UPDATE games
SET pgn = $2, result = $3
WHERE id = $1
RETURNING *;

-- name: DeleteGame :exec
DELETE FROM games
WHERE id = $1;
```

**File:** `backend/db/queries/moves.sql`

```sql
-- name: GetMove :one
SELECT * FROM moves
WHERE id = $1 LIMIT 1;

-- name: ListGameMoves :many
SELECT * FROM moves
WHERE game_id = $1
ORDER BY move_number ASC;

-- name: CreateMove :one
INSERT INTO moves (game_id, move_number, side, move_san, move_uci, fen, time_taken)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: DeleteGameMoves :exec
DELETE FROM moves
WHERE game_id = $1;
```

---

### Phase 3: Docker Setup

#### Step 3.1: Create Dockerfile (Production)
**Time:** 5 minutes

**File:** `backend/Dockerfile`

```dockerfile
# Build stage
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Install sqlc
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Install golang-migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate sqlc code
RUN sqlc generate

# Generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init -g cmd/server/main.go

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /root/

# Copy binary
COPY --from=builder /app/main .

# Copy migrations
COPY --from=builder /app/db/migrations ./db/migrations

# Copy Swagger docs
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./main"]
```

---

#### Step 3.2: Create Dockerfile.dev (Development)
**Time:** 3 minutes

**File:** `backend/Dockerfile.dev`

```dockerfile
FROM golang:1.25.3-alpine

WORKDIR /app

# Install development tools
RUN apk add --no-cache git make postgresql-client

# Install Air (hot reload)
RUN go install github.com/air-verse/air@latest

# Install sqlc
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Install swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Install golang-migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Generate sqlc code
RUN sqlc generate

# Generate Swagger
RUN swag init -g cmd/server/main.go

EXPOSE 8080

CMD ["air"]
```

---

#### Step 3.3: Create docker compose.yml
**Time:** 10 minutes

**File:** `backend/docker compose.yml`

```yaml
version: '3.8'

services:
  # PostgreSQL Database
  postgres:
    image: postgres:16-alpine
    container_name: chess-coach-db
    restart: unless-stopped
    environment:
      POSTGRES_USER: chess_coach
      POSTGRES_PASSWORD: chess_coach_dev
      POSTGRES_DB: chess_coach_dev
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./db/schema.sql:/docker-entrypoint-initdb.d/01-schema.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U chess_coach -d chess_coach_dev"]
      interval: 5s
      timeout: 5s
      retries: 5

  # Go Backend API
  api:
    build:
      context: .
      dockerfile: Dockerfile.dev
    container_name: chess-coach-api
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - GIN_MODE=debug
      - ENVIRONMENT=development
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=chess_coach
      - DB_PASSWORD=chess_coach_dev
      - DB_NAME=chess_coach_dev
      - DB_SSLMODE=disable
    depends_on:
      postgres:
        condition: service_healthy
    volumes:
      - .:/app
      - /app/tmp
    command: sh -c "sqlc generate && swag init -g cmd/server/main.go && air"

  # pgAdmin (Database UI)
  pgadmin:
    image: dpage/pgadmin4:latest
    container_name: chess-coach-pgadmin
    restart: unless-stopped
    environment:
      PGADMIN_DEFAULT_EMAIL: admin@chesscoach.com
      PGADMIN_DEFAULT_PASSWORD: admin
      PGADMIN_CONFIG_SERVER_MODE: 'False'
    ports:
      - "5050:80"
    depends_on:
      - postgres
    volumes:
      - pgadmin_data:/var/lib/pgadmin

volumes:
  postgres_data:
  pgadmin_data:
```

---

#### Step 3.4: Create .dockerignore
**Time:** 2 minutes

**File:** `backend/.dockerignore`

```
# Git
.git
.gitignore

# Documentation
project-docs/
*.md
!README.md

# Tests
*_test.go

# Temporary files
tmp/
*.tmp
*.log

# IDE
.vscode/
.idea/

# Dependencies (will download in container)
vendor/

# OS
.DS_Store

# Generated files (will regenerate)
internal/database/*.go
docs/

# Database
*.db
*.sqlite
```

---

### Phase 4: Go Code Integration

#### Step 4.1: Update Config
**Time:** 5 minutes

**File:** `internal/config/config.go` (update)

```go
// Package config manages application settings.
package config

import "os"

// Config holds app settings.
type Config struct {
	// Server settings
	Port        string
	Environment string
	GinMode     string

	// Database settings
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// Load reads config from env with defaults.
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		GinMode:     getEnv("GIN_MODE", "debug"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "chess_coach"),
		DBPassword: getEnv("DB_PASSWORD", "chess_coach_dev"),
		DBName:     getEnv("DB_NAME", "chess_coach_dev"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// DSN returns PostgreSQL connection string.
func (c *Config) DSN() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=" + c.DBSSLMode
}

// getEnv retrieves env var or returns default.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
```

---

#### Step 4.2: Create Database Connection
**Time:** 10 minutes

**File:** `internal/database/db.go`

```go
// Package database provides database connection and queries.
package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB wraps pgxpool for dependency injection.
type DB struct {
	*Queries
	pool *pgxpool.Pool
}

// New creates database connection pool.
func New(dsn string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	log.Println("Database connected successfully")

	return &DB{
		Queries: New(pool),
		pool:    pool,
	}, nil
}

// Close closes database connection pool.
func (db *DB) Close() {
	db.pool.Close()
}

// Health checks database connection.
func (db *DB) Health(ctx context.Context) error {
	return db.pool.Ping(ctx)
}
```

---

#### Step 4.3: Create Repository Layer (SOLID)
**Time:** 10 minutes

**File:** `internal/repository/user.go`

```go
// Package repository provides data access layer.
package repository

import (
	"context"
	"fmt"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/google/uuid"
)

// UserRepository handles user data access.
type UserRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*database.User, error)
	GetByUsername(ctx context.Context, username string) (*database.User, error)
	Create(ctx context.Context, params database.CreateUserParams) (*database.User, error)
	List(ctx context.Context, limit, offset int32) ([]database.User, error)
}

// userRepo implements UserRepository.
type userRepo struct {
	queries *database.Queries
}

// NewUserRepository creates user repository.
func NewUserRepository(queries *database.Queries) UserRepository {
	return &userRepo{queries: queries}
}

// GetByID retrieves user by ID.
func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*database.User, error) {
	user, err := r.queries.GetUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return &user, nil
}

// GetByUsername retrieves user by username.
func (r *userRepo) GetByUsername(ctx context.Context, username string) (*database.User, error) {
	user, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get user by username: %w", err)
	}
	return &user, nil
}

// Create creates new user.
func (r *userRepo) Create(ctx context.Context, params database.CreateUserParams) (*database.User, error) {
	user, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &user, nil
}

// List retrieves paginated users.
func (r *userRepo) List(ctx context.Context, limit, offset int32) ([]database.User, error) {
	users, err := r.queries.ListUsers(ctx, database.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return users, nil
}
```

---

#### Step 4.4: Update Server to Include Database
**Time:** 5 minutes

**File:** `internal/server/server.go` (update)

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

// New creates a new Server with given config and database.
func New(cfg *config.Config, db *database.DB) Server {
	gin.SetMode(cfg.GinMode)
	r := router.Setup(db)

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
	// Close database connection
	if s.db != nil {
		s.db.Close()
	}
	return s.httpServer.Shutdown(ctx)
}
```

---

#### Step 4.5: Update Main
**Time:** 5 minutes

**File:** `cmd/server/main.go` (update)

```go
// @title Chess Coach API
// @version 1.0
// @description API for chess game analysis and coaching
// @host localhost:8080
// @BasePath /api/v1

// Package main bootstraps the Chess Coach API.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ankits1626/chess-coach-backend/internal/config"
	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/server"
)

// main starts server with graceful shutdown.
func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := database.New(cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create server
	srv := server.New(cfg, db)

	// Start server in goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.Port)
	log.Printf("Swagger UI: http://localhost:%s/swagger/index.html", cfg.Port)

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
```

---

### Phase 5: Makefile & Commands

#### Step 5.1: Create Makefile
**Time:** 10 minutes

**File:** `backend/Makefile`

```makefile
.PHONY: help install-tools sqlc-generate migrate-up migrate-down docker-up docker-down logs test

help:
	@echo "Chess Coach Backend - Commands"
	@echo ""
	@echo "Development:"
	@echo "  make install-tools   - Install sqlc, migrate, etc."
	@echo "  make sqlc-generate   - Generate sqlc code"
	@echo "  make migrate-up      - Run migrations"
	@echo "  make migrate-down    - Rollback migrations"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-up       - Start all containers"
	@echo "  make docker-down     - Stop all containers"
	@echo "  make docker-rebuild  - Rebuild and restart"
	@echo "  make logs            - View logs"
	@echo ""
	@echo "Testing:"
	@echo "  make test            - Run tests"
	@echo "  make test-coverage   - Run tests with coverage"

install-tools:
	@echo "Installing Go tools..."
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Tools installed!"

sqlc-generate:
	sqlc generate

migrate-up:
	migrate -path db/migrations -database "postgres://chess_coach:chess_coach_dev@localhost:5432/chess_coach_dev?sslmode=disable" up

migrate-down:
	migrate -path db/migrations -database "postgres://chess_coach:chess_coach_dev@localhost:5432/chess_coach_dev?sslmode=disable" down

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-rebuild:
	docker compose up -d --build

logs:
	docker compose logs -f api

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
```

---

### Phase 6: Install Dependencies

#### Step 6.1: Add Go Dependencies
**Time:** 3 minutes

```bash
cd backend

# Database
go get github.com/jackc/pgx/v5
go get github.com/jackc/pgx/v5/pgxpool

# UUID
go get github.com/google/uuid

# Validation
go get github.com/go-ozzo/ozzo-validation/v4

# Update go.mod
go mod tidy
```

---

## Execution Order

### Initial Setup (One Time)

```bash
# 1. Install tools
make install-tools

# 2. Generate sqlc code
make sqlc-generate

# 3. Start Docker
make docker-up

# 4. Run migrations (inside container)
docker compose exec api migrate -path db/migrations -database "postgres://chess_coach:chess_coach_dev@postgres:5432/chess_coach_dev?sslmode=disable" up

# 5. View logs
make logs
```

### Daily Development

```bash
# Start
make docker-up

# View logs
make logs

# Stop
make docker-down
```

---

## Verification

### Step 1: Check Containers Running
```bash
docker compose ps
```

**Expected:** 3 containers (postgres, api, pgadmin)

### Step 2: Test Database Connection
```bash
docker compose exec postgres psql -U chess_coach -d chess_coach_dev -c "\dt"
```

**Expected:** List of tables (users, games, moves)

### Step 3: Test API
```bash
curl http://localhost:8080/api/v1/health
```

### Step 4: Test Swagger
Open: http://localhost:8080/swagger/index.html

### Step 5: Test pgAdmin
Open: http://localhost:5050
Login: admin@chesscoach.com / admin

---

## Next Steps (After Setup)

1. **Create user endpoints** (Step 16)
2. **Add validation** with ozzo-validation
3. **Create game endpoints**
4. **Add tests** for repositories
5. **Add chess library** integration

---

## Summary

**What We Set Up:**
- ✅ PostgreSQL in Docker
- ✅ pgx driver (fastest)
- ✅ sqlc code generation
- ✅ Database migrations
- ✅ Repository pattern (SOLID)
- ✅ Docker Compose orchestration
- ✅ Hot reload in development

**Commands:**
```bash
make docker-up      # Start everything
make sqlc-generate  # Generate code
make migrate-up     # Run migrations
make logs           # View logs
```

---

**Total Time:** ~2 hours for complete setup

**Ready to implement? Say "start docker database"!**

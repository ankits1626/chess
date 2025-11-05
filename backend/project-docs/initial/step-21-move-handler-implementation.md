# Move Handler Implementation Guide

**Complete guide for implementing the Move handler in the Chess Coach backend**

---

## Table of Contents
1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Phase 1: Analyze Existing Code](#phase-1-analyze-existing-code)
4. [Phase 2: Create Handler Structure](#phase-2-create-handler-structure)
5. [Phase 3: Register Routes](#phase-3-register-routes)
6. [Phase 4: Build and Test](#phase-4-build-and-test)
7. [Phase 5: Manual Testing](#phase-5-manual-testing)
8. [Phase 6: Advanced Testing](#phase-6-advanced-testing)
9. [Troubleshooting](#troubleshooting)

---

## Overview

### What We're Building

The Move handler provides endpoints for managing individual chess moves within games:

**Endpoints:**
- `GET /api/v1/games/:id/moves` - List all moves for a game
- `GET /api/v1/moves/:id` - Get single move details
- `POST /api/v1/games/:id/moves` - Add move to a game

**Key Features:**
- Ordered move list (by move number)
- Chess notation support (SAN and UCI)
- Position tracking (FEN)
- Time tracking per move
- Cascading delete (moves deleted when game is deleted)

---

## Prerequisites

### What Already Exists

✅ Database schema (moves table)
✅ sqlc queries generated (`internal/database/moves.sql.go`)
✅ Database model (`database.Move`)
✅ Feature-based handler structure (`handler/v1/`)
✅ Repository pattern established
✅ Route versioning setup

### Time Estimate

**Total: 30-45 minutes**
- Phase 1: 5 minutes (analysis)
- Phase 2: 20-25 minutes (implementation)
- Phase 3: 5 minutes (routes)
- Phase 4-6: 10-15 minutes (testing)

---

## Phase 1: Analyze Existing Code

### Step 1.1: Review Database Schema

**File**: `db/schema.sql` (lines 28-40)

```sql
CREATE TABLE moves (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    move_number INTEGER NOT NULL,
    side VARCHAR(5) NOT NULL, -- 'white' or 'black'
    move_san VARCHAR(20) NOT NULL, -- e.g., 'e4', 'Nf3'
    move_uci VARCHAR(10) NOT NULL, -- e.g., 'e2e4'
    fen TEXT NOT NULL, -- Board position after move
    time_taken INTEGER, -- milliseconds (nullable)

    CONSTRAINT valid_side CHECK (side IN ('white', 'black'))
);
```

**Key observations:**
- `game_id` is foreign key with CASCADE delete
- `move_number` orders moves sequentially
- `side` is constrained to 'white' or 'black'
- `time_taken` is optional (nullable)
- No `created_at`/`updated_at` (moves are immutable)

---

### Step 1.2: Review Generated Database Code

**File**: `internal/database/moves.sql.go`

**Available methods:**
```go
// GetMove retrieves single move by ID
func (q *Queries) GetMove(ctx context.Context, id pgtype.UUID) (Move, error)

// ListGameMoves retrieves all moves for a game, ordered by move_number
func (q *Queries) ListGameMoves(ctx context.Context, gameID pgtype.UUID) ([]Move, error)

// CreateMove adds new move to a game
func (q *Queries) CreateMove(ctx context.Context, arg CreateMoveParams) (Move, error)

// DeleteGameMoves deletes all moves for a game
func (q *Queries) DeleteGameMoves(ctx context.Context, gameID pgtype.UUID) error
```

**CreateMoveParams:**
```go
type CreateMoveParams struct {
    GameID     pgtype.UUID `json:"game_id"`
    MoveNumber int32       `json:"move_number"`
    Side       string      `json:"side"`        // "white" or "black"
    MoveSan    string      `json:"move_san"`    // "e4", "Nf3"
    MoveUci    string      `json:"move_uci"`    // "e2e4"
    Fen        string      `json:"fen"`         // Position after move
    TimeTaken  pgtype.Int4 `json:"time_taken"`  // Milliseconds (nullable)
}
```

**Move model:**
```go
type Move struct {
    ID         pgtype.UUID `json:"id"`
    GameID     pgtype.UUID `json:"game_id"`
    MoveNumber int32       `json:"move_number"`
    Side       string      `json:"side"`
    MoveSan    string      `json:"move_san"`
    MoveUci    string      `json:"move_uci"`
    Fen        string      `json:"fen"`
    TimeTaken  pgtype.Int4 `json:"time_taken"` // Nullable
}
```

---

### Step 1.3: Understand Nullable Fields

**Only one nullable field:**
- `TimeTaken` - `pgtype.Int4` (nullable integer)

**Type conversions needed:**
- `pgtype.UUID` ↔ `uuid.UUID`
- `pgtype.Int4` → `*int32` (for JSON response)

---

## Phase 2: Create Handler Structure

### Step 2.1: Create Directory

```bash
cd backend/internal
mkdir -p handler/v1/move
```

---

### Step 2.2: Create DTO File

**File**: `internal/handler/v1/move/dto.go`

```go
// Package move handles move-related HTTP requests for API v1.
package move

import (
	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/google/uuid"
)

// CreateRequest represents move creation request.
type CreateRequest struct {
	MoveNumber int32  `json:"move_number" binding:"required,min=1"`
	Side       string `json:"side" binding:"required,oneof=white black"`
	MoveSan    string `json:"move_san" binding:"required"`
	MoveUci    string `json:"move_uci" binding:"required"`
	Fen        string `json:"fen" binding:"required"`
	TimeTaken  *int32 `json:"time_taken,omitempty"` // Optional: milliseconds
}

// Response represents move in API responses.
type Response struct {
	ID         uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	GameID     uuid.UUID `json:"game_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	MoveNumber int32     `json:"move_number" example:"1"`
	Side       string    `json:"side" example:"white"`
	MoveSan    string    `json:"move_san" example:"e4"`
	MoveUci    string    `json:"move_uci" example:"e2e4"`
	Fen        string    `json:"fen" example:"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"`
	TimeTaken  *int32    `json:"time_taken,omitempty" example:"1500"`
}

// ToResponse converts database.Move to Response.
func ToResponse(m database.Move) Response {
	resp := Response{
		ID:         uuid.UUID(m.ID.Bytes),
		GameID:     uuid.UUID(m.GameID.Bytes),
		MoveNumber: m.MoveNumber,
		Side:       m.Side,
		MoveSan:    m.MoveSan,
		MoveUci:    m.MoveUci,
		Fen:        m.Fen,
	}

	// Handle nullable TimeTaken
	if m.TimeTaken.Valid {
		timeTaken := m.TimeTaken.Int32
		resp.TimeTaken = &timeTaken
	}

	return resp
}

// ToResponses converts slice of database.Move to Response slice.
func ToResponses(moves []database.Move) []Response {
	responses := make([]Response, len(moves))
	for i, m := range moves {
		responses[i] = ToResponse(m)
	}
	return responses
}
```

**Key points:**
- `Side` validated with `oneof=white black`
- `TimeTaken` is pointer for optional field
- `MoveNumber` must be >= 1
- FEN string stores position after move

---

### Step 2.3: Create Handler File

**File**: `internal/handler/v1/move/handler.go`

```go
// Package move handles move-related HTTP requests for API v1.
package move

import (
	"net/http"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Handler handles move endpoints.
type Handler struct {
	repo *repository.MoveRepository
}

// NewHandler creates move handler.
func NewHandler(repo *repository.MoveRepository) *Handler {
	return &Handler{repo: repo}
}

// ListByGame retrieves all moves for a game.
// @Summary List game moves
// @Description Get all moves for a specific game, ordered by move number
// @Tags moves
// @Accept json
// @Produce json
// @Param game_id path string true "Game ID"
// @Success 200 {array} Response
// @Router /games/{game_id}/moves [get]
func (h *Handler) ListByGame(c *gin.Context) {
	gameID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	moves, err := h.repo.ListByGame(c.Request.Context(), gameID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch moves"})
		return
	}

	c.JSON(http.StatusOK, ToResponses(moves))
}

// Get retrieves move by ID.
// @Summary Get move
// @Description Get move by ID
// @Tags moves
// @Accept json
// @Produce json
// @Param id path string true "Move ID"
// @Success 200 {object} Response
// @Router /moves/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid move ID"})
		return
	}

	move, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "move not found"})
		return
	}

	c.JSON(http.StatusOK, ToResponse(move))
}

// Create creates new move for a game.
// @Summary Create move
// @Description Add new move to a game
// @Tags moves
// @Accept json
// @Produce json
// @Param game_id path string true "Game ID"
// @Param move body CreateRequest true "Move data"
// @Success 201 {object} Response
// @Router /games/{game_id}/moves [post]
func (h *Handler) Create(c *gin.Context) {
	gameID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build create params with type conversions
	params := database.CreateMoveParams{
		GameID:     pgtype.UUID{Bytes: gameID, Valid: true},
		MoveNumber: req.MoveNumber,
		Side:       req.Side,
		MoveSan:    req.MoveSan,
		MoveUci:    req.MoveUci,
		Fen:        req.Fen,
	}

	// Handle optional TimeTaken
	if req.TimeTaken != nil && *req.TimeTaken > 0 {
		params.TimeTaken = pgtype.Int4{Int32: *req.TimeTaken, Valid: true}
	}

	move, err := h.repo.Create(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create move"})
		return
	}

	c.JSON(http.StatusCreated, ToResponse(move))
}
```

**Key implementation details:**
- Uses `:id` parameter consistently (for game_id in nested routes)
- Validates side with Gin binding (`oneof=white black`)
- Handles nullable `TimeTaken` properly
- Returns 201 Created for POST
- Returns 404 for missing moves

---

### Step 2.4: Create Repository File

**File**: `internal/repository/move_repository.go`

```go
package repository

import (
	"context"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// MoveRepository handles move business logic.
type MoveRepository struct {
	db *database.DB
}

// NewMoveRepository creates move repository.
func NewMoveRepository(db *database.DB) *MoveRepository {
	return &MoveRepository{db: db}
}

// GetByID retrieves move by ID.
func (r *MoveRepository) GetByID(ctx context.Context, id uuid.UUID) (database.Move, error) {
	return r.db.GetMove(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

// ListByGame retrieves all moves for a game.
func (r *MoveRepository) ListByGame(ctx context.Context, gameID uuid.UUID) ([]database.Move, error) {
	return r.db.ListGameMoves(ctx, pgtype.UUID{Bytes: gameID, Valid: true})
}

// Create creates new move.
func (r *MoveRepository) Create(ctx context.Context, params database.CreateMoveParams) (database.Move, error) {
	return r.db.CreateMove(ctx, params)
}

// DeleteByGame deletes all moves for a game.
func (r *MoveRepository) DeleteByGame(ctx context.Context, gameID uuid.UUID) error {
	return r.db.DeleteGameMoves(ctx, pgtype.UUID{Bytes: gameID, Valid: true})
}
```

**Repository responsibilities:**
- Type conversion (`uuid.UUID` ↔ `pgtype.UUID`)
- Business logic isolation
- Shared across API versions

---

## Phase 3: Register Routes

### Step 3.1: Update Router File

**File**: `internal/router/v1_routes.go`

**Add move handler initialization and routes:**

```go
// Package router provides HTTP routing setup.
package router

import (
	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/handler/v1/game"
	"github.com/ankits1626/chess-coach-backend/internal/handler/v1/health"
	"github.com/ankits1626/chess-coach-backend/internal/handler/v1/move"    // ADD
	"github.com/ankits1626/chess-coach-backend/internal/handler/v1/user"
	"github.com/ankits1626/chess-coach-backend/internal/repository"
	"github.com/gin-gonic/gin"
)

// RegisterV1Routes registers API v1 routes.
func RegisterV1Routes(r *gin.Engine, db *database.DB) {
	v1 := r.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", health.Check)

		// User routes
		userRepo := repository.NewUserRepository(db)
		userHandler := user.NewHandler(userRepo)

		// Game routes
		gameRepo := repository.NewGameRepository(db)
		gameHandler := game.NewHandler(gameRepo)

		// Move routes - ADD ENTIRE BLOCK
		moveRepo := repository.NewMoveRepository(db)
		moveHandler := move.NewHandler(moveRepo)

		users := v1.Group("/users")
		{
			users.GET("", userHandler.List)
			users.GET("/:id", userHandler.Get)
			users.POST("", userHandler.Create)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)

			// User's games sub-resource
			users.GET("/:id/games", gameHandler.ListByUser)
		}

		games := v1.Group("/games")
		{
			games.GET("", gameHandler.List)
			games.GET("/:id", gameHandler.Get)
			games.POST("", gameHandler.Create)

			// Game's moves sub-resource - ADD
			games.GET("/:id/moves", moveHandler.ListByGame)
			games.POST("/:id/moves", moveHandler.Create)
		}

		// Move routes - ADD
		moves := v1.Group("/moves")
		{
			moves.GET("/:id", moveHandler.Get)
		}
	}
}
```

**Route structure:**
- `GET /api/v1/games/:id/moves` - List moves (nested under game)
- `POST /api/v1/games/:id/moves` - Create move (nested under game)
- `GET /api/v1/moves/:id` - Get single move (top-level)

**Why this structure?**
- Moves belong to games (nested routes make sense)
- Getting a specific move by ID doesn't require game context
- Follows RESTful conventions

---

## Phase 4: Build and Test

### Step 4.1: Verify Directory Structure

```bash
tree internal/handler/v1/move
```

**Expected output:**
```
internal/handler/v1/move/
├── dto.go
└── handler.go
```

---

### Step 4.2: Run go mod tidy

```bash
go mod tidy
```

---

### Step 4.3: Build Application

```bash
go build -o bin/server cmd/server/main.go
```

**Expected:** No errors

**If errors occur:**
- Check import paths
- Verify all files created
- Run `go mod tidy` again

---

### Step 4.4: Run Tests

```bash
go test ./...
```

**Expected:** All tests pass (some skipped)

---

### Step 4.5: Run go vet

```bash
go vet ./...
```

**Expected:** No issues

---

## Phase 5: Manual Testing

### Step 5.1: Start Server

```bash
# Terminal 1: Start database (if not running)
docker compose up -d postgres

# Terminal 2: Start server
./bin/server
```

**Expected output:**
```
Database connected successfully
Server starting on :8080
```

---

### Step 5.2: Prepare Test Data

First, ensure you have a game to work with:

```bash
# Get a game ID
curl -s http://localhost:8080/api/v1/games | jq -r '.[0].id'
# Example output: 71c24ce3-2a64-4f20-bfe4-8983f3d4d732
```

---

### Step 5.3: Test Create Move (POST)

```bash
# Replace GAME_ID with actual game ID
curl -X POST http://localhost:8080/api/v1/games/GAME_ID/moves \
  -H "Content-Type: application/json" \
  -d '{
    "move_number": 1,
    "side": "white",
    "move_san": "e4",
    "move_uci": "e2e4",
    "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
    "time_taken": 1500
  }' | jq .
```

**Expected response (201 Created):**
```json
{
  "id": "a1b2c3d4-...",
  "game_id": "71c24ce3-...",
  "move_number": 1,
  "side": "white",
  "move_san": "e4",
  "move_uci": "e2e4",
  "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
  "time_taken": 1500
}
```

---

### Step 5.4: Test Create Move Without TimeTaken

```bash
# Test optional field
curl -X POST http://localhost:8080/api/v1/games/GAME_ID/moves \
  -H "Content-Type: application/json" \
  -d '{
    "move_number": 2,
    "side": "black",
    "move_san": "e5",
    "move_uci": "e7e5",
    "fen": "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2"
  }' | jq .
```

**Expected:** Success with `time_taken: null`

---

### Step 5.5: Test List Game Moves (GET)

```bash
curl -s http://localhost:8080/api/v1/games/GAME_ID/moves | jq .
```

**Expected response (200 OK):**
```json
[
  {
    "id": "a1b2c3d4-...",
    "game_id": "71c24ce3-...",
    "move_number": 1,
    "side": "white",
    "move_san": "e4",
    "move_uci": "e2e4",
    "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
    "time_taken": 1500
  },
  {
    "id": "b2c3d4e5-...",
    "game_id": "71c24ce3-...",
    "move_number": 2,
    "side": "black",
    "move_san": "e5",
    "move_uci": "e7e5",
    "fen": "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2"
  }
]
```

**Verify:** Moves are ordered by `move_number` ASC

---

### Step 5.6: Test Get Single Move (GET)

```bash
# Get first move ID from previous response
MOVE_ID="a1b2c3d4-..."

curl -s http://localhost:8080/api/v1/moves/$MOVE_ID | jq .
```

**Expected response (200 OK):**
```json
{
  "id": "a1b2c3d4-...",
  "game_id": "71c24ce3-...",
  "move_number": 1,
  "side": "white",
  "move_san": "e4",
  "move_uci": "e2e4",
  "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
  "time_taken": 1500
}
```

---

### Step 5.7: Test Error Cases

**Invalid game ID:**
```bash
curl -s http://localhost:8080/api/v1/games/invalid-uuid/moves | jq .
```
**Expected:** 400 Bad Request

**Invalid move ID:**
```bash
curl -s http://localhost:8080/api/v1/moves/invalid-uuid | jq .
```
**Expected:** 400 Bad Request

**Non-existent move:**
```bash
curl -s http://localhost:8080/api/v1/moves/00000000-0000-0000-0000-000000000000 | jq .
```
**Expected:** 404 Not Found

**Invalid side:**
```bash
curl -X POST http://localhost:8080/api/v1/games/GAME_ID/moves \
  -H "Content-Type: application/json" \
  -d '{
    "move_number": 3,
    "side": "invalid",
    "move_san": "Nf3",
    "move_uci": "g1f3",
    "fen": "..."
  }' | jq .
```
**Expected:** 400 Bad Request (validation error)

**Missing required field:**
```bash
curl -X POST http://localhost:8080/api/v1/games/GAME_ID/moves \
  -H "Content-Type: application/json" \
  -d '{
    "move_number": 3,
    "side": "white",
    "move_san": "Nf3"
  }' | jq .
```
**Expected:** 400 Bad Request (missing move_uci, fen)

---

## Phase 6: Advanced Testing

### Step 6.1: Test Complete Game

Create a full opening sequence:

```bash
# 1. e4
curl -X POST http://localhost:8080/api/v1/games/GAME_ID/moves \
  -H "Content-Type: application/json" \
  -d '{
    "move_number": 1,
    "side": "white",
    "move_san": "e4",
    "move_uci": "e2e4",
    "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
  }'

# 1... e5
curl -X POST http://localhost:8080/api/v1/games/GAME_ID/moves \
  -H "Content-Type: application/json" \
  -d '{
    "move_number": 2,
    "side": "black",
    "move_san": "e5",
    "move_uci": "e7e5",
    "fen": "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2"
  }'

# 2. Nf3
curl -X POST http://localhost:8080/api/v1/games/GAME_ID/moves \
  -H "Content-Type: application/json" \
  -d '{
    "move_number": 3,
    "side": "white",
    "move_san": "Nf3",
    "move_uci": "g1f3",
    "fen": "rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2"
  }'

# Verify all moves
curl -s http://localhost:8080/api/v1/games/GAME_ID/moves | jq '.[] | {move_number, side, move_san}'
```

**Expected output:**
```json
{
  "move_number": 1,
  "side": "white",
  "move_san": "e4"
}
{
  "move_number": 2,
  "side": "black",
  "move_san": "e5"
}
{
  "move_number": 3,
  "side": "white",
  "move_san": "Nf3"
}
```

---

### Step 6.2: Test Cascading Delete

Verify moves are deleted when game is deleted:

```bash
# Count moves before deletion
curl -s http://localhost:8080/api/v1/games/GAME_ID/moves | jq 'length'
# Output: 3

# Delete the game
curl -X DELETE http://localhost:8080/api/v1/games/GAME_ID

# Try to fetch moves (should be gone)
curl -s http://localhost:8080/api/v1/games/GAME_ID/moves | jq .
# Expected: [] (empty array) or 404
```

---

### Step 6.3: Verify Database Directly

```bash
docker compose exec postgres psql -U postgres -d chess_coach -c \
  "SELECT move_number, side, move_san FROM moves WHERE game_id = 'GAME_ID' ORDER BY move_number;"
```

**Expected:** Matches API responses

---

## Troubleshooting

### Issue: Type mismatch errors

**Problem:** `cannot use uuid.UUID as pgtype.UUID`

**Solution:** Ensure proper conversion:
```go
pgtype.UUID{Bytes: id, Valid: true}
```

---

### Issue: Route conflict with game routes

**Problem:**
```
panic: ':id' in new path '/api/v1/games/:id/moves'
```

**Solution:** This shouldn't happen as `/games/:id/moves` is distinct from `/games/:id`. If it does:
1. Ensure move routes are registered AFTER single game GET route
2. Verify route group structure

---

### Issue: Nullable field handling

**Problem:** Panic when accessing `TimeTaken`

**Solution:** Always check `Valid` flag:
```go
if m.TimeTaken.Valid {
    timeTaken := m.TimeTaken.Int32
    resp.TimeTaken = &timeTaken
}
```

---

### Issue: Validation error - invalid side

**Problem:** 400 error when creating move

**Solution:** Ensure side is exactly `"white"` or `"black"` (lowercase)

---

### Issue: Foreign key constraint violation

**Problem:**
```
ERROR: insert or update on table "moves" violates foreign key constraint
```

**Root Cause:** Game ID doesn't exist

**Solution:**
1. Verify game exists: `curl http://localhost:8080/api/v1/games/GAME_ID`
2. Use valid game ID

---

### Issue: Move ordering incorrect

**Problem:** Moves returned in wrong order

**Solution:** Check SQL query has `ORDER BY move_number ASC`. Already present in `db/queries/moves.sql`:
```sql
-- name: ListGameMoves :many
SELECT * FROM moves
WHERE game_id = $1
ORDER BY move_number ASC;  -- ✓ Correct
```

---

## Summary

### What We Built

✅ **Move Handler** (`internal/handler/v1/move/`)
- `handler.go` - 3 endpoint handlers
- `dto.go` - Request/Response types

✅ **Move Repository** (`internal/repository/move_repository.go`)
- Data access layer
- Type conversions
- Shared across versions

✅ **Routes** (nested and top-level)
- `GET /api/v1/games/:id/moves` - List game moves
- `POST /api/v1/games/:id/moves` - Create move
- `GET /api/v1/moves/:id` - Get single move

✅ **Features**
- Ordered move sequences
- Chess notation (SAN + UCI)
- Position tracking (FEN)
- Time tracking (optional)
- Proper validation
- Cascading delete

---

### Testing Checklist

- [ ] Build succeeds (`go build`)
- [ ] Tests pass (`go test ./...`)
- [ ] No vet issues (`go vet ./...`)
- [ ] Create move with time_taken
- [ ] Create move without time_taken
- [ ] List game moves (ordered)
- [ ] Get single move
- [ ] Invalid game ID (400)
- [ ] Invalid move ID (400)
- [ ] Non-existent move (404)
- [ ] Invalid side validation (400)
- [ ] Missing required fields (400)
- [ ] Complete game sequence
- [ ] Cascading delete works

---

### Next Steps

After completing Move handler:

1. **Update Swagger Documentation**
   ```bash
   swag init -g cmd/server/main.go -o docs
   ```

2. **Add Unit Tests**
   - `internal/handler/v1/move/handler_test.go`
   - `internal/repository/move_repository_test.go`

3. **Consider Additional Features**
   - Move annotations (?, !, !!, ??)
   - Position evaluation
   - Best move suggestions
   - Opening name detection

4. **Authentication/Authorization**
   - JWT-based auth
   - User ownership validation
   - Protected endpoints

---

**Last Updated**: 2025-11-01

**Related Docs**:
- [step-20-game-handler-implementation.md](./step-20-game-handler-implementation.md)
- [step-19-package-reorganization.md](./step-19-package-reorganization.md)
- [api-versioning-strategy.md](./api-versioning-strategy.md)
- [request-response-flow.md](./request-response-flow.md)

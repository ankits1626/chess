# Step 20: Game Handler Implementation

**Objective**: Implement HTTP handlers for game operations to expose chess games via REST API.

**Time**: 3-4 hours
**Difficulty**: Medium

---

## Overview

You already have:
- ✅ Database schema for games
- ✅ SQL queries for games (sqlc-generated)
- ✅ GameRepository with data access methods

What you'll build:
- 🎯 Game HTTP handlers (CRUD operations)
- 🎯 Game DTOs (request/response models)
- 🎯 Route registration
- 🎯 Swagger documentation
- 🎯 End-to-end testing

---

## Prerequisites

Before starting, verify these files exist:

```bash
# Check repository exists
ls backend/internal/repository/game_repository.go

# Check sqlc-generated code
ls backend/internal/database/games.sql.go

# Check game queries
ls backend/db/queries/games.sql
```

---

## Architecture Overview

```
Client Request (POST /api/v1/games)
    ↓
Router (v1_routes.go) - Matches route
    ↓
Handler (game/handler.go) - Parse JSON, validate
    ↓
DTO (game/dto.go) - Convert to internal types
    ↓
Repository (game_repository.go) - Business logic
    ↓
Database (games.sql.go - sqlc) - Execute SQL
    ↓
PostgreSQL - Store/retrieve data
    ↓
Response (JSON) - Return game data
```

---

## Phase 1: Analyze Existing Code

**Time**: 15 minutes

### Step 1.1: Review Game Repository

**File**: `internal/repository/game_repository.go`

Open and review the available methods:

```bash
cat backend/internal/repository/game_repository.go
```

**Expected methods:**
- `GetByID(ctx, id)` - Get single game
- `List(ctx, limit, offset)` - List all games (paginated)
- `ListByUser(ctx, userID, limit, offset)` - Get user's games
- `Create(ctx, params)` - Create new game

### Step 1.2: Review Database Schema

**File**: `internal/database/models.go` (sqlc-generated)

```go
type Game struct {
    ID            pgtype.UUID      `json:"id"`
    WhitePlayerID pgtype.UUID      `json:"white_player_id"`
    BlackPlayerID pgtype.UUID      `json:"black_player_id"`
    Pgn           string           `json:"pgn"`
    Result        pgtype.Text      `json:"result"`
    TimeControl   pgtype.Int4      `json:"time_control"`
    CreatedAt     pgtype.Timestamp `json:"created_at"`
    UpdatedAt     pgtype.Timestamp `json:"updated_at"`
}
```

**Note the types:**
- `pgtype.UUID` - Database UUIDs (need conversion to `uuid.UUID`)
- `pgtype.Text` - Nullable text (need to handle `Valid` flag)
- `pgtype.Int4` - Nullable int32 (need to handle `Valid` flag)
- `pgtype.Timestamp` - Database timestamp (convert to `time.Time`)

### Step 1.3: Review SQL Queries

**File**: `db/queries/games.sql`

```bash
cat backend/db/queries/games.sql
```

Understand what operations are available.

---

## Phase 2: Create Game Handler Structure

**Time**: 20 minutes

### Step 2.1: Create Directory

```bash
cd backend/internal/handler/v1
mkdir -p game
```

### Step 2.2: Create DTO File

**File**: `internal/handler/v1/game/dto.go`

```go
// Package game handles game-related HTTP requests for API v1.
package game

import (
	"time"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/google/uuid"
)

// CreateRequest represents game creation request.
type CreateRequest struct {
	WhitePlayerID uuid.UUID `json:"white_player_id" binding:"required"`
	BlackPlayerID uuid.UUID `json:"black_player_id" binding:"required"`
	Pgn           string    `json:"pgn" binding:"required"`
	Result        string    `json:"result,omitempty"`        // Optional: "1-0", "0-1", "1/2-1/2"
	TimeControl   int32     `json:"time_control,omitempty"` // Optional: seconds
}

// UpdateRequest represents game update request.
type UpdateRequest struct {
	Pgn         string `json:"pgn" binding:"required"`
	Result      string `json:"result,omitempty"`
	TimeControl int32  `json:"time_control,omitempty"`
}

// Response represents game in API responses.
type Response struct {
	ID            uuid.UUID  `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	WhitePlayerID uuid.UUID  `json:"white_player_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	BlackPlayerID uuid.UUID  `json:"black_player_id" example:"550e8400-e29b-41d4-a716-446655440002"`
	Pgn           string     `json:"pgn" example:"1. e4 e5 2. Nf3 Nc6"`
	Result        *string    `json:"result,omitempty" example:"1-0"`
	TimeControl   *int32     `json:"time_control,omitempty" example:"600"`
	CreatedAt     time.Time  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt     time.Time  `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// ToResponse converts database.Game to Response.
func ToResponse(g database.Game) Response {
	resp := Response{
		ID:            uuid.UUID(g.ID.Bytes),
		WhitePlayerID: uuid.UUID(g.WhitePlayerID.Bytes),
		BlackPlayerID: uuid.UUID(g.BlackPlayerID.Bytes),
		Pgn:           g.Pgn,
		CreatedAt:     g.CreatedAt.Time,
		UpdatedAt:     g.UpdatedAt.Time,
	}

	// Handle nullable Result
	if g.Result.Valid {
		result := g.Result.String
		resp.Result = &result
	}

	// Handle nullable TimeControl
	if g.TimeControl.Valid {
		timeControl := g.TimeControl.Int32
		resp.TimeControl = &timeControl
	}

	return resp
}

// ToResponses converts slice of database.Game to Response slice.
func ToResponses(games []database.Game) []Response {
	responses := make([]Response, len(games))
	for i, g := range games {
		responses[i] = ToResponse(g)
	}
	return responses
}
```

**Key Points:**
- `CreateRequest` - Required fields for creating a game
- `UpdateRequest` - Fields that can be updated
- `Response` - API response format with proper JSON tags
- Nullable fields use pointers (`*string`, `*int32`)
- Conversion functions handle `pgtype` → Go types

---

### Step 2.3: Create Handler File

**File**: `internal/handler/v1/game/handler.go`

```go
// Package game handles game-related HTTP requests for API v1.
package game

import (
	"net/http"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Handler handles game endpoints.
type Handler struct {
	repo *repository.GameRepository
}

// NewHandler creates game handler.
func NewHandler(repo *repository.GameRepository) *Handler {
	return &Handler{repo: repo}
}

// List lists games with pagination.
// @Summary List games
// @Description Get paginated list of games
// @Tags games
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} Response
// @Router /games [get]
func (h *Handler) List(c *gin.Context) {
	var params struct {
		Limit  int32 `form:"limit" binding:"omitempty,min=1,max=100"`
		Offset int32 `form:"offset" binding:"omitempty,min=0"`
	}

	// Default limit
	if params.Limit == 0 {
		params.Limit = 10
	}

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	games, err := h.repo.List(c.Request.Context(), params.Limit, params.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch games"})
		return
	}

	c.JSON(http.StatusOK, ToResponses(games))
}

// Get retrieves game by ID.
// @Summary Get game
// @Description Get game by ID
// @Tags games
// @Accept json
// @Produce json
// @Param id path string true "Game ID"
// @Success 200 {object} Response
// @Router /games/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game ID"})
		return
	}

	game, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}

	c.JSON(http.StatusOK, ToResponse(game))
}

// ListByUser retrieves user's games.
// @Summary List user's games
// @Description Get paginated list of games for a specific user
// @Tags games
// @Accept json
// @Produce json
// @Param user_id path string true "User ID"
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} Response
// @Router /users/{user_id}/games [get]
func (h *Handler) ListByUser(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var params struct {
		Limit  int32 `form:"limit" binding:"omitempty,min=1,max=100"`
		Offset int32 `form:"offset" binding:"omitempty,min=0"`
	}

	if params.Limit == 0 {
		params.Limit = 10
	}

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	games, err := h.repo.ListByUser(c.Request.Context(), userID, params.Limit, params.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user games"})
		return
	}

	c.JSON(http.StatusOK, ToResponses(games))
}

// Create creates new game.
// @Summary Create game
// @Description Create new chess game
// @Tags games
// @Accept json
// @Produce json
// @Param game body CreateRequest true "Game data"
// @Success 201 {object} Response
// @Router /games [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Build create params with type conversions
	params := database.CreateGameParams{
		WhitePlayerID: pgtype.UUID{Bytes: req.WhitePlayerID, Valid: true},
		BlackPlayerID: pgtype.UUID{Bytes: req.BlackPlayerID, Valid: true},
		Pgn:           req.Pgn,
	}

	// Handle optional Result
	if req.Result != "" {
		params.Result = pgtype.Text{String: req.Result, Valid: true}
	}

	// Handle optional TimeControl
	if req.TimeControl > 0 {
		params.TimeControl = pgtype.Int4{Int32: req.TimeControl, Valid: true}
	}

	game, err := h.repo.Create(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create game"})
		return
	}

	c.JSON(http.StatusCreated, ToResponse(game))
}
```

**Key Points:**
- `List()` - Paginated list of all games
- `Get()` - Get single game by ID
- `ListByUser()` - Get games for specific user (white or black player)
- `Create()` - Create new game with validation
- Type conversions: `uuid.UUID` → `pgtype.UUID`, `int32` → `pgtype.Int4`
- Nullable fields handled with `Valid` flag
- Swagger annotations for API documentation

---

## Phase 3: Register Routes

**Time**: 10 minutes

### Step 3.1: Update v1_routes.go

**File**: `internal/router/v1_routes.go`

Add game routes after user routes:

```go
// Package router provides HTTP routing setup.
package router

import (
	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/handler/v1/game"    // ADD
	"github.com/ankits1626/chess-coach-backend/internal/handler/v1/health"
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

		// Game routes - ADD ENTIRE BLOCK BEFORE USER ROUTES
		gameRepo := repository.NewGameRepository(db)
		gameHandler := game.NewHandler(gameRepo)

		users := v1.Group("/users")
		{
			users.GET("", userHandler.List)
			users.GET("/:id", userHandler.Get)
			users.POST("", userHandler.Create)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)

			// User's games sub-resource - Use :id not :user_id to avoid conflict
			users.GET("/:id/games", gameHandler.ListByUser)  // ADD
		}

		games := v1.Group("/games")
		{
			games.GET("", gameHandler.List)       // List all games
			games.GET("/:id", gameHandler.Get)    // Get single game
			games.POST("", gameHandler.Create)    // Create game
		}
	}
}
```

**Route Structure:**
```
GET    /api/v1/games              - List all games
GET    /api/v1/games/:id          - Get game by ID
POST   /api/v1/games              - Create new game
GET    /api/v1/users/:user_id/games - Get user's games
```

---

## Phase 4: Build and Test

**Time**: 30 minutes

### Step 4.1: Tidy Dependencies

```bash
cd backend
go mod tidy
```

### Step 4.2: Regenerate Swagger Docs

```bash
swag init -g cmd/server/main.go
```

**Expected output:**
```
Generate swagger docs....
Generating game.Response
Generating game.CreateRequest
Generating game.UpdateRequest
...
create docs.go at docs/docs.go
```

### Step 4.3: Build Application

```bash
go build -o bin/server cmd/server/main.go
```

**Should compile without errors.**

### Step 4.4: Start Application

```bash
# Ensure database is running
docker compose up -d postgres

# Run server
./bin/server
```

**Expected output:**
```
INFO: Database connected successfully
INFO: Server started on port 8080
INFO: Swagger UI: http://localhost:8080/swagger/index.html
```

---

## Phase 5: Manual Testing

**Time**: 30 minutes

### Step 5.1: Test Swagger UI

1. Open browser: `http://localhost:8080/swagger/index.html`
2. Verify you see new game endpoints:
   - ✅ `GET /api/v1/games`
   - ✅ `GET /api/v1/games/{id}`
   - ✅ `POST /api/v1/games`
   - ✅ `GET /api/v1/users/{user_id}/games`

### Step 5.2: Create Test Users

First, create two users for testing:

```bash
# Create white player
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "magnus",
    "email": "magnus@chess.com",
    "rating": 2800
  }'

# Save the ID from response
WHITE_PLAYER_ID="<id-from-response>"

# Create black player
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "hikaru",
    "email": "hikaru@chess.com",
    "rating": 2750
  }'

# Save the ID from response
BLACK_PLAYER_ID="<id-from-response>"
```

### Step 5.3: Create Game

```bash
curl -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -d '{
    "white_player_id": "'$WHITE_PLAYER_ID'",
    "black_player_id": "'$BLACK_PLAYER_ID'",
    "pgn": "1. e4 e5 2. Nf3 Nc6 3. Bb5",
    "result": "1-0",
    "time_control": 600
  }'
```

**Expected Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "white_player_id": "...",
  "black_player_id": "...",
  "pgn": "1. e4 e5 2. Nf3 Nc6 3. Bb5",
  "result": "1-0",
  "time_control": 600,
  "created_at": "2025-11-01T10:30:00Z",
  "updated_at": "2025-11-01T10:30:00Z"
}
```

### Step 5.4: List All Games

```bash
curl http://localhost:8080/api/v1/games
```

**Expected Response (200 OK):**
```json
[
  {
    "id": "...",
    "white_player_id": "...",
    "black_player_id": "...",
    "pgn": "1. e4 e5 2. Nf3 Nc6 3. Bb5",
    "result": "1-0",
    "time_control": 600,
    "created_at": "2025-11-01T10:30:00Z",
    "updated_at": "2025-11-01T10:30:00Z"
  }
]
```

### Step 5.5: Get Single Game

```bash
GAME_ID="<id-from-previous-response>"
curl http://localhost:8080/api/v1/games/$GAME_ID
```

### Step 5.6: Get User's Games

```bash
curl http://localhost:8080/api/v1/users/$WHITE_PLAYER_ID/games
```

**Should return games where user is white or black player.**

### Step 5.7: Test Pagination

```bash
# Create more games (run create command multiple times)

# Test pagination
curl "http://localhost:8080/api/v1/games?limit=5&offset=0"
curl "http://localhost:8080/api/v1/games?limit=5&offset=5"
```

### Step 5.8: Test Error Cases

```bash
# Invalid game ID
curl http://localhost:8080/api/v1/games/invalid-id
# Expected: 400 Bad Request

# Non-existent game
curl http://localhost:8080/api/v1/games/00000000-0000-0000-0000-000000000000
# Expected: 404 Not Found

# Missing required fields
curl -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -d '{"pgn": "1. e4"}'
# Expected: 400 Bad Request with validation errors
```

---

## Phase 6: Advanced Testing (Optional)

**Time**: 30 minutes

### Step 6.1: Test with Invalid PGN

```bash
curl -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -d '{
    "white_player_id": "'$WHITE_PLAYER_ID'",
    "black_player_id": "'$BLACK_PLAYER_ID'",
    "pgn": "invalid chess moves"
  }'
```

**Note:** Currently no PGN validation. This will be added in future steps.

### Step 6.2: Test with Non-Existent Players

```bash
curl -X POST http://localhost:8080/api/v1/games \
  -H "Content-Type: application/json" \
  -d '{
    "white_player_id": "00000000-0000-0000-0000-000000000000",
    "black_player_id": "00000000-0000-0000-0000-000000000001",
    "pgn": "1. e4 e5"
  }'
```

**Expected:** Database foreign key constraint violation (500 error).

### Step 6.3: Performance Test

Create multiple games and test pagination:

```bash
# Create 20 games
for i in {1..20}; do
  curl -X POST http://localhost:8080/api/v1/games \
    -H "Content-Type: application/json" \
    -d '{
      "white_player_id": "'$WHITE_PLAYER_ID'",
      "black_player_id": "'$BLACK_PLAYER_ID'",
      "pgn": "1. e4 e5 '${i}'. Nf3"
    }' &
done
wait

# Test pagination performance
time curl "http://localhost:8080/api/v1/games?limit=10&offset=0"
```

---

## Phase 7: Verify Database

**Time**: 10 minutes

### Step 7.1: Check Database via pgAdmin

1. Open pgAdmin: `http://localhost:5050`
2. Login: `admin@chesscoach.com` / `admin`
3. Navigate to: `chess_coach` → `Schemas` → `public` → `Tables` → `games`
4. Right-click → `View/Edit Data` → `All Rows`

**Verify:**
- Games are stored correctly
- UUIDs are valid
- Foreign keys reference existing users
- Timestamps are set
- Nullable fields handled correctly

### Step 7.2: SQL Queries

```sql
-- Count games
SELECT COUNT(*) FROM games;

-- List recent games with player info
SELECT
    g.id,
    g.pgn,
    g.result,
    u1.username as white_player,
    u2.username as black_player,
    g.created_at
FROM games g
JOIN users u1 ON g.white_player_id = u1.id
JOIN users u2 ON g.black_player_id = u2.id
ORDER BY g.created_at DESC
LIMIT 10;

-- Games per user
SELECT
    u.username,
    COUNT(*) as game_count
FROM users u
LEFT JOIN games g ON u.id = g.white_player_id OR u.id = g.black_player_id
GROUP BY u.username
ORDER BY game_count DESC;
```

---

## Troubleshooting

### Issue: Import errors

**Solution:**
```bash
cd backend
go mod tidy
go clean -cache
```

### Issue: Swagger not showing game endpoints

**Solution:**
```bash
swag init -g cmd/server/main.go --parseDependency --parseInternal
```

### Issue: Type conversion errors

**Problem:** `cannot use uuid.UUID as pgtype.UUID`

**Solution:** Ensure you're using proper conversion:
```go
pgtype.UUID{Bytes: id, Valid: true}
```

### Issue: Route conflict - ':user_id' conflicts with ':id'

**Problem:**
```
panic: ':user_id' in new path '/api/v1/users/:user_id/games' conflicts with
existing wildcard ':id' in existing prefix '/api/v1/users/:id'
```

**Root Cause:** Gin's router cannot distinguish between:
- `/users/:id` (get user by ID)
- `/users/:user_id/games` (get user's games)

Both start with `/users/:` and Gin sees this as conflicting wildcards.

**Solution:** Use the same parameter name `:id` for both routes:

```go
// WRONG - Causes conflict
users.GET("/:id", userHandler.Get)
users.GET("/:user_id/games", gameHandler.ListByUser)

// CORRECT - Use same parameter name
users.GET("/:id", userHandler.Get)
users.GET("/:id/games", gameHandler.ListByUser)
```

Then in the handler, extract it as `c.Param("id")`:

```go
func (h *Handler) ListByUser(c *gin.Context) {
    userID, err := uuid.Parse(c.Param("id"))  // Use "id" not "user_id"
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }
    // ... rest of handler
}
```

---

### Issue: Nullable field handling

**Problem:** Panic when accessing nullable fields

**Solution:** Always check `Valid` flag:
```go
if g.Result.Valid {
    result := g.Result.String
    resp.Result = &result
}
```

### Issue: Foreign key constraint violation

**Problem:** Creating game with non-existent user IDs

**Solution:** Verify users exist before creating game:
```bash
curl http://localhost:8080/api/v1/users/$WHITE_PLAYER_ID
```

---

## Final Project Structure

```
backend/
├── internal/
│   ├── handler/
│   │   └── v1/
│   │       ├── game/              # ✨ NEW
│   │       │   ├── dto.go
│   │       │   └── handler.go
│   │       ├── health/
│   │       │   └── handler.go
│   │       └── user/
│   │           ├── dto.go
│   │           └── handler.go
│   ├── repository/
│   │   ├── game_repository.go     # ✅ Already exists
│   │   └── user_repository.go
│   └── router/
│       ├── router.go
│       └── v1_routes.go           # ✨ Updated
└── docs/
    ├── docs.go                    # ✨ Updated (Swagger)
    ├── swagger.json
    └── swagger.yaml
```

---

## API Endpoints Summary

### Games

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/games` | List all games (paginated) |
| GET | `/api/v1/games/:id` | Get game by ID |
| POST | `/api/v1/games` | Create new game |
| GET | `/api/v1/users/:user_id/games` | Get user's games |

### Request Examples

**Create Game:**
```json
POST /api/v1/games
{
  "white_player_id": "uuid",
  "black_player_id": "uuid",
  "pgn": "1. e4 e5 2. Nf3",
  "result": "1-0",
  "time_control": 600
}
```

**List Games:**
```
GET /api/v1/games?limit=10&offset=0
```

---

## What's Next?

After completing this step, you can:

1. **Add Update/Delete operations** - Modify existing games
2. **Add Move handlers** - Store individual moves
3. **Add PGN validation** - Ensure valid chess notation
4. **Add game search** - Filter by player, date, result
5. **Add statistics** - Win/loss records, ELO changes
6. **Add authentication** - Protect game creation

---

## Completion Checklist

- [ ] Created `internal/handler/v1/game/dto.go`
- [ ] Created `internal/handler/v1/game/handler.go`
- [ ] Updated `internal/router/v1_routes.go`
- [ ] Ran `go mod tidy`
- [ ] Regenerated Swagger docs
- [ ] Application builds successfully
- [ ] Server starts without errors
- [ ] Swagger UI shows game endpoints
- [ ] Can create users
- [ ] Can create games
- [ ] Can list games
- [ ] Can get single game
- [ ] Can get user's games
- [ ] Pagination works
- [ ] Error handling works
- [ ] Database stores games correctly

---

## Key Learnings

### 1. DTO Pattern
- Separate API models from database models
- Use pointers for nullable fields
- Provide conversion functions

### 2. Type Conversions
- `uuid.UUID` → `pgtype.UUID{Bytes: id, Valid: true}`
- `int32` → `pgtype.Int4{Int32: value, Valid: true}`
- `string` → `pgtype.Text{String: value, Valid: true}`

### 3. Nullable Fields
- Database: `pgtype.Text`, `pgtype.Int4`
- API Response: `*string`, `*int32`
- Always check `Valid` flag before accessing

### 4. REST API Design
- Use nested routes for relationships: `/users/:user_id/games`
- Consistent error responses
- Proper HTTP status codes
- Pagination for lists

### 5. Repository Pattern
- Handler calls repository
- Repository handles business logic
- Database layer executes SQL
- Clear separation of concerns

---

**Last Updated**: 2025-11-01

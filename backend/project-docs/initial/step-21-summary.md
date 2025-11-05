# Move Handler Implementation - Quick Reference

**Quick guide for implementing Chess Move tracking**

---

## What You're Building

**Endpoints:**
- `GET /api/v1/games/:id/moves` - List all moves in a game (ordered)
- `POST /api/v1/games/:id/moves` - Add move to game
- `GET /api/v1/moves/:id` - Get single move details

**Time:** 30-45 minutes

---

## Quick Start

### 1. Create Files (5 minutes)

```bash
cd backend/internal
mkdir -p handler/v1/move
```

Create 3 files:
1. `handler/v1/move/dto.go` - Request/Response types
2. `handler/v1/move/handler.go` - HTTP handlers
3. `repository/move_repository.go` - Data access

### 2. Update Routes (2 minutes)

Edit `router/v1_routes.go`:
- Add move handler import
- Initialize move repository
- Add routes under `/games/:id/moves`
- Add route for `/moves/:id`

### 3. Test (10 minutes)

```bash
go build -o bin/server cmd/server/main.go
./bin/server

# Create move
curl -X POST http://localhost:8080/api/v1/games/GAME_ID/moves \
  -H "Content-Type: application/json" \
  -d '{
    "move_number": 1,
    "side": "white",
    "move_san": "e4",
    "move_uci": "e2e4",
    "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
  }'

# List moves
curl http://localhost:8080/api/v1/games/GAME_ID/moves
```

---

## Key Implementation Points

### DTO (dto.go)

**CreateRequest:**
```go
type CreateRequest struct {
    MoveNumber int32  `json:"move_number" binding:"required,min=1"`
    Side       string `json:"side" binding:"required,oneof=white black"`
    MoveSan    string `json:"move_san" binding:"required"`
    MoveUci    string `json:"move_uci" binding:"required"`
    Fen        string `json:"fen" binding:"required"`
    TimeTaken  *int32 `json:"time_taken,omitempty"` // Optional
}
```

**Response:**
```go
type Response struct {
    ID         uuid.UUID `json:"id"`
    GameID     uuid.UUID `json:"game_id"`
    MoveNumber int32     `json:"move_number"`
    Side       string    `json:"side"`
    MoveSan    string    `json:"move_san"`
    MoveUci    string    `json:"move_uci"`
    Fen        string    `json:"fen"`
    TimeTaken  *int32    `json:"time_taken,omitempty"`
}
```

### Handler (handler.go)

**3 methods:**
1. `ListByGame(c *gin.Context)` - GET /games/:id/moves
2. `Get(c *gin.Context)` - GET /moves/:id
3. `Create(c *gin.Context)` - POST /games/:id/moves

**Pattern:** Same as Game handler
- Parse UUID from `c.Param("id")`
- Call repository
- Convert to DTO with `ToResponse()`
- Return JSON

### Repository (move_repository.go)

**4 methods:**
```go
func (r *MoveRepository) GetByID(ctx, id) (database.Move, error)
func (r *MoveRepository) ListByGame(ctx, gameID) ([]database.Move, error)
func (r *MoveRepository) Create(ctx, params) (database.Move, error)
func (r *MoveRepository) DeleteByGame(ctx, gameID) error
```

---

## Database Schema

```sql
CREATE TABLE moves (
    id UUID PRIMARY KEY,
    game_id UUID REFERENCES games(id) ON DELETE CASCADE,
    move_number INTEGER NOT NULL,
    side VARCHAR(5) CHECK (side IN ('white', 'black')),
    move_san VARCHAR(20) NOT NULL,  -- "e4", "Nf3"
    move_uci VARCHAR(10) NOT NULL,  -- "e2e4"
    fen TEXT NOT NULL,              -- Position after move
    time_taken INTEGER              -- Milliseconds (nullable)
);
```

**Key points:**
- Cascading delete (moves deleted with game)
- Ordered by `move_number`
- Side constrained to 'white' or 'black'
- `time_taken` is optional

---

## Common Issues

### 1. Nullable Field Handling

**Problem:** `TimeTaken` is nullable

**Solution in DTO:**
```go
// Request
TimeTaken *int32 `json:"time_taken,omitempty"`

// Response
if m.TimeTaken.Valid {
    timeTaken := m.TimeTaken.Int32
    resp.TimeTaken = &timeTaken
}
```

**Solution in Handler:**
```go
if req.TimeTaken != nil && *req.TimeTaken > 0 {
    params.TimeTaken = pgtype.Int4{Int32: *req.TimeTaken, Valid: true}
}
```

### 2. Type Conversions

**Always convert:**
- `uuid.UUID` → `pgtype.UUID{Bytes: id, Valid: true}`
- `pgtype.UUID` → `uuid.UUID(m.ID.Bytes)`
- `pgtype.Int4` → `*int32` (with Valid check)

### 3. Route Parameter

**Use `:id` consistently:**
```go
// Route: /games/:id/moves
gameID := c.Param("id")  // Not "game_id"!
```

### 4. Side Validation

**Must be lowercase:**
- ✅ `"white"`, `"black"`
- ❌ `"White"`, `"BLACK"`

---

## Testing Sequence

1. ✅ Create move with time_taken
2. ✅ Create move without time_taken
3. ✅ List game moves (verify ordering)
4. ✅ Get single move
5. ✅ Invalid game ID (400)
6. ✅ Invalid move ID (400)
7. ✅ Invalid side (400)
8. ✅ Missing required fields (400)

---

## Example: Complete Opening

```bash
# Get game ID
GAME_ID=$(curl -s http://localhost:8080/api/v1/games | jq -r '.[0].id')

# 1. e4
curl -X POST http://localhost:8080/api/v1/games/$GAME_ID/moves \
  -H "Content-Type: application/json" \
  -d '{
    "move_number": 1,
    "side": "white",
    "move_san": "e4",
    "move_uci": "e2e4",
    "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
  }'

# 1... e5
curl -X POST http://localhost:8080/api/v1/games/$GAME_ID/moves \
  -H "Content-Type: application/json" \
  -d '{
    "move_number": 2,
    "side": "black",
    "move_san": "e5",
    "move_uci": "e7e5",
    "fen": "rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq e6 0 2"
  }'

# List all moves
curl -s http://localhost:8080/api/v1/games/$GAME_ID/moves | jq '.[] | {move_number, side, move_san}'
```

---

## Route Registration

```go
// RegisterV1Routes in internal/router/v1_routes.go
func RegisterV1Routes(r *gin.Engine, db *database.DB) {
    v1 := r.Group("/api/v1")

    // ... existing code ...

    // Move routes - ADD THIS
    moveRepo := repository.NewMoveRepository(db)
    moveHandler := move.NewHandler(moveRepo)

    games := v1.Group("/games")
    {
        games.GET("", gameHandler.List)
        games.GET("/:id", gameHandler.Get)
        games.POST("", gameHandler.Create)

        // Game's moves sub-resource - ADD THIS
        games.GET("/:id/moves", moveHandler.ListByGame)
        games.POST("/:id/moves", moveHandler.Create)
    }

    // Top-level move routes - ADD THIS
    moves := v1.Group("/moves")
    {
        moves.GET("/:id", moveHandler.Get)
    }
}
```

---

## Benefits

After completing this:
- ✅ Full game replay capability
- ✅ Move-by-move analysis ready
- ✅ Complete Chess Coach API (Users → Games → Moves)
- ✅ Foundation for AI analysis
- ✅ Opening detection ready

---

## Full Implementation Guide

See [step-21-move-handler-implementation.md](./step-21-move-handler-implementation.md) for:
- Complete code listings
- Detailed explanations
- Advanced testing scenarios
- Database verification
- Troubleshooting guide

---

**Next After Completion:**
1. Update Swagger docs (`swag init`)
2. Add unit tests
3. Implement authentication
4. Add move analysis features

---

**Last Updated**: 2025-11-01

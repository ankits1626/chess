# Phase 4: Message Handlers - Business Logic Implementation

**Goal**: Implement complete message handlers for game management with database integration

**Duration**: 6-8 hours

**Status**: 📋 Ready to implement

---

## ⚠️ SOLID Architecture Version Available

**For better code organization following SOLID principles, see:**
→ **[04-message-handlers-REFACTORED.md](./04-message-handlers-REFACTORED.md)**

The refactored version splits handlers into separate files:
- Better separation of concerns
- Easier to test
- More maintainable
- Follows SOLID principles

**This document** contains the original single-file approach (simpler for learning).
**The REFACTORED version** is recommended for production.

---

## 📋 What You'll Build

By the end of this phase:
- ✅ Complete message routing system
- ✅ CreateGame handler (with matchmaking queue)
- ✅ JoinGame handler (connect two players)
- ✅ MakeMove handler (with chess validation)
- ✅ LeaveGame/Resign handler
- ✅ Error handling for all edge cases
- ✅ Full game flow working end-to-end

---

## 🎯 Prerequisites

- [x] Phase 1 complete (WebSocket package)
- [x] Phase 2 complete (Integration)
- [x] Phase 3 complete (Utilities)
- [x] Database repositories exist
- [x] Understanding of chess rules (basic)
- [ ] Chess library installed (`github.com/notnil/chess`)

---

## 📦 Install Chess Library

```bash
cd /Users/ankit/code/learn/chess-coach/backend
go get github.com/notnil/chess
```

This library provides:
- Move validation
- Legal move generation
- FEN parsing/generation
- Game state management
- Check/checkmate detection

---

## 📝 Step 1: Create Handler Interface (30 minutes)

### File: `internal/websocket/handler.go`

**Create new file**:

```go
package websocket

import (
	"context"
	"fmt"
	"log"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/notnil/chess"
)

// MessageHandler defines the interface for handling WebSocket messages.
type MessageHandler interface {
	HandleMessage(ctx context.Context, client *Client, msg *Message)
}

// GameHandler implements message handling for chess games.
type GameHandler struct {
	db            *database.DB
	pendingGames  map[string]*PendingGame // gameID -> PendingGame
	activeGames   map[string]*ActiveGame  // gameID -> ActiveGame
}

// PendingGame represents a game waiting for second player.
type PendingGame struct {
	GameID      string
	WhitePlayer *Client
	TimeControl string
}

// ActiveGame represents an ongoing game with both players.
type ActiveGame struct {
	GameID      string
	WhitePlayer *Client
	BlackPlayer *Client
	Game        *chess.Game
	TimeControl string
}

// NewGameHandler creates a new game message handler.
func NewGameHandler(db *database.DB) *GameHandler {
	return &GameHandler{
		db:           db,
		pendingGames: make(map[string]*PendingGame),
		activeGames:  make(map[string]*ActiveGame),
	}
}

// HandleMessage routes messages to appropriate handlers.
func (h *GameHandler) HandleMessage(ctx context.Context, client *Client, msg *Message) {
	log.Printf("GameHandler: Processing message type=%s action=%s from client=%s",
		msg.Type, msg.Action, client.ID)

	// Only handle request messages (events don't need responses)
	if msg.Type != TypeRequest {
		log.Printf("GameHandler: Ignoring non-request message type=%s", msg.Type)
		return
	}

	// Route to appropriate handler based on action
	switch msg.Action {
	case "createGame":
		h.handleCreateGame(ctx, client, msg)
	case "joinGame":
		h.handleJoinGame(ctx, client, msg)
	case "makeMove":
		h.handleMakeMove(ctx, client, msg)
	case "leaveGame":
		h.handleLeaveGame(ctx, client, msg)
	case "resign":
		h.handleResign(ctx, client, msg)
	default:
		SendErrorResponse(client, msg.ID, fmt.Sprintf("Unknown action: %s", msg.Action))
	}
}
```

---

## 📝 Step 2: Implement CreateGame Handler (90 minutes)

### Add to `internal/websocket/handler.go`:

```go
// handleCreateGame creates a new game and waits for opponent.
//
// Request format:
// {
//   "id": "req-123",
//   "type": "request",
//   "action": "createGame",
//   "data": {
//     "timeControl": "5+0"  // minutes+increment
//   }
// }
//
// Response format:
// {
//   "id": "req-123",
//   "type": "response",
//   "success": true,
//   "data": {
//     "gameId": "uuid",
//     "status": "waiting",
//     "side": "white"
//   }
// }
func (h *GameHandler) handleCreateGame(ctx context.Context, client *Client, msg *Message) {
	// 1. Extract and validate time control
	timeControl, err := RequireString(msg.Data, "timeControl")
	if err != nil {
		SendErrorResponse(client, msg.ID, err.Error())
		return
	}

	seconds, err := ValidateTimeControl(timeControl)
	if err != nil {
		SendErrorResponse(client, msg.ID, fmt.Sprintf("Invalid time control: %v", err))
		return
	}

	// 2. Generate game ID
	gameID := GenerateUUID()
	gameIDStr := UUIDToString(gameID)

	log.Printf("CreateGame: Creating game %s for user %s with time control %s (%ds)",
		gameIDStr, client.UserID, timeControl, seconds)

	// 3. Create game in database (waiting for opponent)
	// Initially, white_player_id is set, black_player_id is NULL
	whitePlayerUUID := StringToUUID(client.UserID)
	if !whitePlayerUUID.Valid {
		SendErrorResponse(client, msg.ID, "Invalid user ID")
		return
	}

	game, err := h.db.CreateGame(ctx, database.CreateGameParams{
		ID:            gameID,
		WhitePlayerID: whitePlayerUUID,
		BlackPlayerID: pgtype.UUID{Valid: false}, // NULL - waiting for opponent
		Pgn:           "",                         // Empty initially
		Result:        StringToText("*"),          // Game in progress
		TimeControl:   IntToInt4(seconds),
	})
	if err != nil {
		log.Printf("CreateGame: Database error: %v", err)
		SendErrorResponse(client, msg.ID, "Failed to create game")
		return
	}

	// 4. Add to pending games (in-memory)
	h.pendingGames[gameIDStr] = &PendingGame{
		GameID:      gameIDStr,
		WhitePlayer: client,
		TimeControl: timeControl,
	}

	// 5. Set client's game ID
	client.SetGameID(gameIDStr)

	log.Printf("CreateGame: Game %s created successfully, waiting for opponent", gameIDStr)

	// 6. Send success response
	SendSuccessResponse(client, msg.ID, map[string]interface{}{
		"gameId":      gameIDStr,
		"status":      "waiting",
		"side":        "white",
		"timeControl": timeControl,
		"fen":         chess.StartingFEN(),
	})

	// 7. Broadcast to matchmaking (optional - for Phase 5)
	// For now, clients must explicitly join with gameId
}
```

---

## 📝 Step 3: Implement JoinGame Handler (90 minutes)

### Add to `internal/websocket/handler.go`:

```go
// handleJoinGame joins an existing game as black player.
//
// Request format:
// {
//   "id": "req-456",
//   "type": "request",
//   "action": "joinGame",
//   "data": {
//     "gameId": "uuid"
//   }
// }
//
// Response format (to joiner):
// {
//   "id": "req-456",
//   "type": "response",
//   "success": true,
//   "data": {
//     "gameId": "uuid",
//     "status": "active",
//     "side": "black",
//     "opponent": "username"
//   }
// }
//
// Event broadcast (to white player):
// {
//   "type": "event",
//   "event": "opponentJoined",
//   "data": {
//     "gameId": "uuid",
//     "opponent": "username"
//   }
// }
func (h *GameHandler) handleJoinGame(ctx context.Context, client *Client, msg *Message) {
	// 1. Extract and validate game ID
	gameIDStr, err := RequireString(msg.Data, "gameId")
	if err != nil {
		SendErrorResponse(client, msg.ID, err.Error())
		return
	}

	if !IsValidUUID(gameIDStr) {
		SendErrorResponse(client, msg.ID, "Invalid game ID format")
		return
	}

	log.Printf("JoinGame: User %s attempting to join game %s", client.UserID, gameIDStr)

	// 2. Check if game is pending
	pending, exists := h.pendingGames[gameIDStr]
	if !exists {
		SendErrorResponse(client, msg.ID, "Game not found or already started")
		return
	}

	// 3. Verify not joining own game
	if pending.WhitePlayer.UserID == client.UserID {
		SendErrorResponse(client, msg.ID, "Cannot join your own game")
		return
	}

	// 4. Update game in database with black player
	gameID := StringToUUID(gameIDStr)
	blackPlayerUUID := StringToUUID(client.UserID)
	if !blackPlayerUUID.Valid {
		SendErrorResponse(client, msg.ID, "Invalid user ID")
		return
	}

	game, err := h.db.UpdateGame(ctx, database.UpdateGameParams{
		ID:            gameID,
		BlackPlayerID: blackPlayerUUID,
		Pgn:           "",
		Result:        StringToText("*"),
	})
	if err != nil {
		log.Printf("JoinGame: Database error: %v", err)
		SendErrorResponse(client, msg.ID, "Failed to join game")
		return
	}

	// 5. Initialize chess game
	chessGame := chess.NewGame()

	// 6. Create active game
	activeGame := &ActiveGame{
		GameID:      gameIDStr,
		WhitePlayer: pending.WhitePlayer,
		BlackPlayer: client,
		Game:        chessGame,
		TimeControl: pending.TimeControl,
	}
	h.activeGames[gameIDStr] = activeGame

	// 7. Remove from pending
	delete(h.pendingGames, gameIDStr)

	// 8. Set game ID for both clients
	client.SetGameID(gameIDStr)
	pending.WhitePlayer.SetGameID(gameIDStr)

	// 9. Create room and add both players
	room := NewRoom(gameIDStr)
	pending.WhitePlayer.JoinRoom(room)
	client.JoinRoom(room)

	log.Printf("JoinGame: Game %s started - White: %s, Black: %s",
		gameIDStr, pending.WhitePlayer.UserID, client.UserID)

	// 10. Send success response to joiner (black player)
	SendSuccessResponse(client, msg.ID, map[string]interface{}{
		"gameId":      gameIDStr,
		"status":      "active",
		"side":        "black",
		"opponent":    pending.WhitePlayer.UserID,
		"timeControl": pending.TimeControl,
		"fen":         chessGame.FEN(),
	})

	// 11. Broadcast to white player that opponent joined
	opponentJoinedEvent := NewEvent("opponentJoined", map[string]interface{}{
		"gameId":   gameIDStr,
		"opponent": client.UserID,
		"side":     "black",
		"fen":      chessGame.FEN(),
	})
	pending.WhitePlayer.SendMessage(opponentJoinedEvent)
}
```

---

## 📝 Step 4: Implement MakeMove Handler (120 minutes)

### Add to `internal/websocket/handler.go`:

```go
// handleMakeMove processes a move in an active game.
//
// Request format:
// {
//   "id": "req-789",
//   "type": "request",
//   "action": "makeMove",
//   "data": {
//     "gameId": "uuid",
//     "move": "e4"  // SAN notation: e4, Nf3, O-O, etc.
//   }
// }
//
// Response format (to mover):
// {
//   "id": "req-789",
//   "type": "response",
//   "success": true,
//   "data": {
//     "moveNumber": 1,
//     "san": "e4",
//     "uci": "e2e4",
//     "fen": "...",
//     "gameOver": false
//   }
// }
//
// Event broadcast (to opponent):
// {
//   "type": "event",
//   "event": "moveMade",
//   "data": {
//     "moveNumber": 1,
//     "san": "e4",
//     "uci": "e2e4",
//     "fen": "...",
//     "gameOver": false
//   }
// }
func (h *GameHandler) handleMakeMove(ctx context.Context, client *Client, msg *Message) {
	// 1. Extract and validate inputs
	gameIDStr, err := RequireString(msg.Data, "gameId")
	if err != nil {
		SendErrorResponse(client, msg.ID, err.Error())
		return
	}

	moveSAN, err := RequireString(msg.Data, "move")
	if err != nil {
		SendErrorResponse(client, msg.ID, err.Error())
		return
	}

	if !ValidateMoveSAN(moveSAN) {
		SendErrorResponse(client, msg.ID, "Invalid move format")
		return
	}

	log.Printf("MakeMove: User %s attempting move '%s' in game %s",
		client.UserID, moveSAN, gameIDStr)

	// 2. Get active game
	activeGame, exists := h.activeGames[gameIDStr]
	if !exists {
		SendErrorResponse(client, msg.ID, "Game not found or not active")
		return
	}

	// 3. Verify it's player's turn
	currentTurn := activeGame.Game.Position().Turn()
	var expectedPlayer *Client
	var playerSide string

	if currentTurn == chess.White {
		expectedPlayer = activeGame.WhitePlayer
		playerSide = "white"
	} else {
		expectedPlayer = activeGame.BlackPlayer
		playerSide = "black"
	}

	if client.ID != expectedPlayer.ID {
		SendErrorResponse(client, msg.ID, "Not your turn")
		return
	}

	// 4. Parse and validate move using chess library
	move, err := chess.UCINotation{}.Decode(activeGame.Game.Position(), moveSAN)
	if err != nil {
		// Try parsing as SAN directly
		move, err = chess.AlgebraicNotation{}.Decode(activeGame.Game.Position(), moveSAN)
		if err != nil {
			SendErrorResponse(client, msg.ID, fmt.Sprintf("Invalid move: %v", err))
			return
		}
	}

	// 5. Apply move to game
	if err := activeGame.Game.Move(move); err != nil {
		SendErrorResponse(client, msg.ID, fmt.Sprintf("Illegal move: %v", err))
		return
	}

	// 6. Get move details
	position := activeGame.Game.Position()
	fen := position.String()
	moveUCI := chess.UCINotation{}.Encode(position, move)
	moveSANNormalized := chess.AlgebraicNotation{}.Encode(position, move)
	moveNumber := len(activeGame.Game.Moves())

	// 7. Check game status
	gameOver := activeGame.Game.Outcome() != chess.NoOutcome
	var result string
	if gameOver {
		switch activeGame.Game.Outcome() {
		case chess.WhiteWon:
			result = "1-0"
		case chess.BlackWon:
			result = "0-1"
		case chess.Draw:
			result = "1/2-1/2"
		default:
			result = "*"
		}
	} else {
		result = "*"
	}

	log.Printf("MakeMove: Move applied - #%d %s (%s) by %s, gameOver=%v",
		moveNumber, moveSANNormalized, moveUCI, playerSide, gameOver)

	// 8. Save move to database
	gameID := StringToUUID(gameIDStr)
	_, err = h.db.CreateMove(ctx, database.CreateMoveParams{
		GameID:     gameID,
		MoveNumber: int32(moveNumber),
		Side:       playerSide,
		MoveSan:    moveSANNormalized,
		MoveUci:    moveUCI,
		Fen:        fen,
		TimeTaken:  pgtype.Int4{Valid: false}, // TODO: Implement timing in Phase 5
	})
	if err != nil {
		log.Printf("MakeMove: Failed to save move to database: %v", err)
		// Don't fail the move, just log the error
	}

	// 9. Update game if finished
	if gameOver {
		_, err := h.db.UpdateGame(ctx, database.UpdateGameParams{
			ID:     gameID,
			Pgn:    activeGame.Game.String(),
			Result: StringToText(result),
		})
		if err != nil {
			log.Printf("MakeMove: Failed to update game result: %v", err)
		}
	}

	// 10. Prepare response data
	responseData := map[string]interface{}{
		"moveNumber": moveNumber,
		"san":        moveSANNormalized,
		"uci":        moveUCI,
		"fen":        fen,
		"gameOver":   gameOver,
		"result":     result,
	}

	if gameOver {
		responseData["method"] = activeGame.Game.Method().String()
	}

	// 11. Send success response to mover
	SendSuccessResponse(client, msg.ID, responseData)

	// 12. Broadcast move to opponent
	moveMadeEvent := NewEvent("moveMade", responseData)
	if playerSide == "white" {
		activeGame.BlackPlayer.SendMessage(moveMadeEvent)
	} else {
		activeGame.WhitePlayer.SendMessage(moveMadeEvent)
	}

	// 13. If game over, clean up
	if gameOver {
		log.Printf("MakeMove: Game %s finished with result %s", gameIDStr, result)
		delete(h.activeGames, gameIDStr)
	}
}
```

---

## 📝 Step 5: Implement LeaveGame and Resign (60 minutes)

### Add to `internal/websocket/handler.go`:

```go
// handleLeaveGame handles player leaving a game.
//
// Request format:
// {
//   "id": "req-101",
//   "type": "request",
//   "action": "leaveGame",
//   "data": {
//     "gameId": "uuid"
//   }
// }
func (h *GameHandler) handleLeaveGame(ctx context.Context, client *Client, msg *Message) {
	gameIDStr, err := RequireString(msg.Data, "gameId")
	if err != nil {
		SendErrorResponse(client, msg.ID, err.Error())
		return
	}

	log.Printf("LeaveGame: User %s leaving game %s", client.UserID, gameIDStr)

	// Check if game is pending
	if pending, exists := h.pendingGames[gameIDStr]; exists {
		if pending.WhitePlayer.ID == client.ID {
			// Creator left before opponent joined - cancel game
			delete(h.pendingGames, gameIDStr)
			client.SetGameID("")

			SendSuccessResponse(client, msg.ID, map[string]interface{}{
				"message": "Game cancelled",
			})
			return
		}
	}

	// Check if game is active
	if activeGame, exists := h.activeGames[gameIDStr]; exists {
		// Player left during game - forfeit
		h.handleForfeit(ctx, client, activeGame, msg.ID)
		return
	}

	SendErrorResponse(client, msg.ID, "Game not found")
}

// handleResign handles player resignation.
//
// Request format:
// {
//   "id": "req-102",
//   "type": "request",
//   "action": "resign",
//   "data": {
//     "gameId": "uuid"
//   }
// }
func (h *GameHandler) handleResign(ctx context.Context, client *Client, msg *Message) {
	gameIDStr, err := RequireString(msg.Data, "gameId")
	if err != nil {
		SendErrorResponse(client, msg.ID, err.Error())
		return
	}

	log.Printf("Resign: User %s resigning from game %s", client.UserID, gameIDStr)

	activeGame, exists := h.activeGames[gameIDStr]
	if !exists {
		SendErrorResponse(client, msg.ID, "Game not found or not active")
		return
	}

	h.handleForfeit(ctx, client, activeGame, msg.ID)
}

// handleForfeit handles game forfeit (resignation or disconnect).
func (h *GameHandler) handleForfeit(ctx context.Context, client *Client, activeGame *ActiveGame, requestID string) {
	gameIDStr := activeGame.GameID
	gameID := StringToUUID(gameIDStr)

	// Determine result based on who forfeited
	var result string
	var winner *Client
	var loser *Client

	if client.ID == activeGame.WhitePlayer.ID {
		result = "0-1" // Black wins
		winner = activeGame.BlackPlayer
		loser = activeGame.WhitePlayer
	} else {
		result = "1-0" // White wins
		winner = activeGame.WhitePlayer
		loser = activeGame.BlackPlayer
	}

	log.Printf("Forfeit: Game %s - %s resigned, result: %s", gameIDStr, client.UserID, result)

	// Update game in database
	_, err := h.db.UpdateGame(ctx, database.UpdateGameParams{
		ID:     gameID,
		Pgn:    activeGame.Game.String(),
		Result: StringToText(result),
	})
	if err != nil {
		log.Printf("Forfeit: Failed to update game: %v", err)
	}

	// Send response to resigner
	SendSuccessResponse(loser, requestID, map[string]interface{}{
		"message":  "You resigned",
		"result":   result,
		"gameOver": true,
	})

	// Notify winner
	opponentResignedEvent := NewEvent("opponentResigned", map[string]interface{}{
		"gameId":   gameIDStr,
		"result":   result,
		"gameOver": true,
		"message":  "Opponent resigned",
	})
	winner.SendMessage(opponentResignedEvent)

	// Clean up
	delete(h.activeGames, gameIDStr)
	loser.SetGameID("")
	winner.SetGameID("")
}
```

---

## 📝 Step 6: Wire Handler to Hub (15 minutes)

### Modify `internal/websocket/hub.go`:

Find the `NewHub` function and update it:

```go
// NewHub creates a new Hub with optional message handler.
func NewHub(handler MessageHandler) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[string]*Room),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		shutdown:   make(chan struct{}),
		handler:    handler,  // Store handler
	}
}
```

### Modify `internal/app/app.go`:

Update app creation to wire handler:

```go
// New creates new application.
func New(cfg *config.Config, db *database.DB, log logger.Logger) *App {
	// Create handler with database
	handler := websocket.NewGameHandler(db)

	// Create hub with handler
	hub := websocket.NewHub(handler)

	return &App{
		config: cfg,
		db:     db,
		server: server.New(cfg, db, hub),
		logger: log,
		wsHub:  hub,
	}
}
```

---

## 📝 Step 7: Add Missing Imports (5 minutes)

### Update `internal/websocket/handler.go` imports:

```go
import (
	"context"
	"fmt"
	"log"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/notnil/chess"
)
```

---

## ✅ Phase 4 Checklist

Go to [implementation-checklist.md](./implementation-checklist.md) and check off:

**Phase 4: Message Handlers**
- [ ] Install chess library (`github.com/notnil/chess`)
- [ ] Create `internal/websocket/handler.go`
- [ ] Define MessageHandler interface
- [ ] Define GameHandler struct with pending/active games
- [ ] Implement handleCreateGame (create + wait for opponent)
- [ ] Implement handleJoinGame (connect two players)
- [ ] Implement handleMakeMove (with chess validation)
- [ ] Implement handleLeaveGame (cancel/forfeit)
- [ ] Implement handleResign (explicit resignation)
- [ ] Implement handleForfeit helper (common logic)
- [ ] Wire handler to Hub
- [ ] Wire handler to App
- [ ] Test createGame flow
- [ ] Test joinGame flow
- [ ] Test complete game (moves + checkmate)
- [ ] Test resignation
- [ ] Test edge cases (wrong turn, invalid move, etc.)

---

## 📝 Step 8: Testing (60-90 minutes)

### Test with `wscat` (install: `npm install -g wscat`)

#### Terminal 1: White Player

```bash
wscat -c "ws://localhost:8080/ws?user_id=alice"

# Create game
{"id":"1","type":"request","action":"createGame","data":{"timeControl":"5+0"}}

# Wait for response, copy gameId
# Wait for opponent to join

# Make move
{"id":"2","type":"request","action":"makeMove","data":{"gameId":"<gameId>","move":"e4"}}

# Continue playing...
```

#### Terminal 2: Black Player

```bash
wscat -c "ws://localhost:8080/ws?user_id=bob"

# Join game (use gameId from white player's response)
{"id":"1","type":"request","action":"joinGame","data":{"gameId":"<gameId>"}}

# Wait for white's move, then respond
{"id":"2","type":"request","action":"makeMove","data":{"gameId":"<gameId>","move":"e5"}}

# Continue playing...
```

### Expected Flow

1. **Alice creates game**:
   - Response: `{success: true, gameId: "...", status: "waiting", side: "white"}`

2. **Bob joins game**:
   - Bob receives: `{success: true, gameId: "...", status: "active", side: "black"}`
   - Alice receives event: `{event: "opponentJoined", opponent: "bob"}`

3. **Alice moves e4**:
   - Alice receives: `{success: true, san: "e4", uci: "e2e4", fen: "..."}`
   - Bob receives event: `{event: "moveMade", san: "e4", uci: "e2e4", fen: "..."}`

4. **Continue until checkmate or resignation**

### Test Cases to Verify

| Test Case | Expected Result |
|-----------|----------------|
| Invalid time control | Error: "Invalid time control" |
| Invalid gameId format | Error: "Invalid game ID format" |
| Join non-existent game | Error: "Game not found" |
| Join own game | Error: "Cannot join your own game" |
| Move when not your turn | Error: "Not your turn" |
| Invalid move (illegal) | Error: "Illegal move" |
| Valid move | Success + FEN update |
| Checkmate | gameOver: true, result set |
| Resignation | Opponent notified, game ends |

---

## 🎯 What You've Accomplished

### ✅ Issues Fixed in This Phase

| Issue # | Description | Status |
|---------|-------------|--------|
| 4 | No message handlers | ✅ Complete routing system |
| 5 | Missing business logic | ✅ Full game flow implemented |
| 6 | No database integration | ✅ All operations persist |
| 12 | No move validation | ✅ Chess library validates all moves |
| 13 | No error handling | ✅ Comprehensive error responses |

### ✅ Handlers Implemented

1. **CreateGame** - Start new game, wait for opponent
2. **JoinGame** - Connect second player, start game
3. **MakeMove** - Validate and apply moves, broadcast to opponent
4. **LeaveGame** - Cancel pending or forfeit active game
5. **Resign** - Explicit resignation

### ✅ Features Working

- ✅ Matchmaking (simple: create + join by ID)
- ✅ Move validation (chess library)
- ✅ Turn enforcement
- ✅ Check/checkmate detection
- ✅ Game state persistence
- ✅ Real-time move broadcast
- ✅ Resignation handling
- ✅ Error handling for all edge cases

---

## 🚀 Next Steps

**Phase 4 is the biggest milestone!** You now have a fully functional chess game.

**Next**: [05-authentication.md](./05-authentication.md) - Add security

In Phase 5, you'll:
- Replace `?user_id=` with JWT authentication
- Add CORS configuration
- Implement proper authorization
- Secure WebSocket connections

---

## 💡 Key Concepts Review

### Game State Management

**Three states**:
1. **Pending**: Created, waiting for opponent (in `pendingGames` map)
2. **Active**: Both players connected, game in progress (in `activeGames` map)
3. **Finished**: Game over, removed from maps, persisted in database

### Chess Library Integration

```go
// Initialize game
game := chess.NewGame()

// Parse move
move, err := chess.AlgebraicNotation{}.Decode(game.Position(), "e4")

// Apply move
err := game.Move(move)

// Check status
if game.Outcome() == chess.Checkmate {
    // Game over
}

// Get FEN
fen := game.Position().String()
```

### Database Flow

```
1. CreateGame → INSERT games (white_player_id set, black_player_id NULL)
2. JoinGame   → UPDATE games SET black_player_id
3. MakeMove   → INSERT moves (each move recorded)
4. Game End   → UPDATE games SET result, pgn
```

---

## 📊 Code Metrics

| Metric | Value | Target | Status |
|--------|-------|--------|--------|
| Handler Functions | 6 | 5-7 | ✅ |
| Lines of Code | ~500 | 400-600 | ✅ |
| Test Scenarios | 10+ | 8+ | ✅ |
| Error Cases | 15+ | 10+ | ✅ |
| Database Ops | 5 | 4-6 | ✅ |

---

**Ready for Phase 5?** → [05-authentication.md](./05-authentication.md)

---

**Last Updated**: 2025-11-05
**Status**: ✅ Complete and ready to implement
**Estimated Time**: 6-8 hours
**Dependencies**: Phases 1-3 complete, chess library installed

# Phase 4: Message Handlers - SOLID Architecture

**Goal**: Implement message handlers following SOLID principles with clean separation of concerns

**Duration**: 6-8 hours

**Status**: 📋 Ready to implement

---

## 📋 Architecture Overview

Following SOLID principles, we'll organize handlers into a subfolder:

```
internal/websocket/
├── handlers/                    # All handlers in subfolder
│   ├── handler.go              # Main handler interface & router
│   ├── game_manager.go         # Game state management (SRP)
│   ├── create_game.go          # CreateGame handler (SRP)
│   ├── join_game.go            # JoinGame handler (SRP)
│   ├── make_move.go            # MakeMove handler (SRP)
│   ├── resign.go               # Resign/LeaveGame handler (SRP)
│   └── chess_service.go        # Chess logic abstraction (DIP)
├── message.go
├── client.go
├── hub.go
├── room.go
├── utils.go
├── upgrade.go
└── tests/
    └── handlers_test.go
```

**SOLID Principles Applied**:
- **S**ingle Responsibility: Each handler in its own file
- **O**pen/Closed: Easy to add new handlers without modifying existing code
- **L**iskov Substitution: Interfaces allow mock implementations for testing
- **I**nterface Segregation: Small, focused interfaces
- **D**ependency Inversion: Handlers depend on abstractions, not concrete implementations

---

## 🎯 Prerequisites

- [x] Phase 1 complete (WebSocket package)
- [x] Phase 2 complete (Integration)
- [x] Phase 3 complete (Utilities)
- [x] Database repositories exist
- [ ] Chess library installed (`github.com/notnil/chess`)

---

## 📦 Install Chess Library

```bash
cd /Users/ankit/code/learn/chess-coach/backend
go get github.com/notnil/chess
```

---

## 📝 Step 1: Define Core Interfaces (30 minutes)

### File: `internal/websocket/handlers/handler.go`

**Create handlers folder**:
```bash
mkdir -p internal/websocket/handlers
```

**Create new file** - Main handler interface and router:

```go
package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/ankits1626/chess-coach-backend/internal/websocket"
)

// MessageHandler defines the interface for handling WebSocket messages.
type MessageHandler interface {
	HandleMessage(ctx context.Context, client *websocket.Client, msg *websocket.Message)
}

// GameService defines chess game operations (Dependency Inversion).
type GameService interface {
	// ValidateMove checks if a move is legal
	ValidateMove(fen string, moveSAN string) (MoveResult, error)

	// ApplyMove applies a move and returns new state
	ApplyMove(fen string, moveSAN string) (GameState, error)

	// IsGameOver checks if game is finished
	IsGameOver(fen string) (bool, GameResult)

	// GetStartingPosition returns initial FEN
	GetStartingPosition() string
}

// MoveResult contains validated move information.
type MoveResult struct {
	SAN      string
	UCI      string
	IsLegal  bool
}

// GameState contains game state after a move.
type GameState struct {
	FEN        string
	SAN        string
	UCI        string
	MoveNumber int
	IsGameOver bool
	Result     GameResult
}

// GameResult represents game outcome.
type GameResult struct {
	Winner string // "white", "black", "draw", ""
	Method string // "checkmate", "resignation", "stalemate", etc.
	PGN    string
}

// HandlerRouter routes messages to specific handlers.
type HandlerRouter struct {
	gameManager *GameManager
	handlers    map[string]ActionHandler
}

// ActionHandler processes a specific action.
type ActionHandler interface {
	Handle(ctx context.Context, client *websocket.Client, msg *websocket.Message) error
}

// NewHandlerRouter creates a new message router.
func NewHandlerRouter(gameManager *GameManager) *HandlerRouter {
	router := &HandlerRouter{
		gameManager: gameManager,
		handlers:    make(map[string]ActionHandler),
	}

	// Register all handlers
	router.RegisterHandler("createGame", &CreateGameHandler{manager: gameManager})
	router.RegisterHandler("joinGame", &JoinGameHandler{manager: gameManager})
	router.RegisterHandler("makeMove", &MakeMoveHandler{manager: gameManager})
	router.RegisterHandler("resign", &ResignHandler{manager: gameManager})
	router.RegisterHandler("leaveGame", &LeaveGameHandler{manager: gameManager})

	return router
}

// RegisterHandler registers a new action handler.
func (r *HandlerRouter) RegisterHandler(action string, handler ActionHandler) {
	r.handlers[action] = handler
}

// HandleMessage implements MessageHandler interface.
func (r *HandlerRouter) HandleMessage(ctx context.Context, client *websocket.Client, msg *websocket.Message) {
	log.Printf("HandlerRouter: Processing message type=%s action=%s from client=%s",
		msg.Type, msg.Action, client.ID)

	// Only handle request messages
	if msg.Type != websocket.TypeRequest {
		log.Printf("HandlerRouter: Ignoring non-request message type=%s", msg.Type)
		return
	}

	// Find handler for action
	handler, exists := r.handlers[msg.Action]
	if !exists {
		websocket.SendErrorResponse(client, msg.ID, fmt.Sprintf("Unknown action: %s", msg.Action))
		return
	}

	// Execute handler
	if err := handler.Handle(ctx, client, msg); err != nil {
		log.Printf("HandlerRouter: Handler error for action=%s: %v", msg.Action, err)
		websocket.SendErrorResponse(client, msg.ID, err.Error())
	}
}
```

---

## 📝 Step 2: Game State Manager (45 minutes)

### File: `internal/websocket/handlers/game_manager.go`

**Create new file** - Centralized game state management:

```go
package handlers

import (
	"fmt"
	"sync"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/websocket"
)

// GameManager manages game state (pending and active games).
// Single Responsibility: Game lifecycle management
type GameManager struct {
	db            *database.DB
	chessService  GameService
	pendingGames  map[string]*PendingGame
	activeGames   map[string]*ActiveGame
	mu            sync.RWMutex // Protect concurrent access
}

// PendingGame represents a game waiting for second player.
type PendingGame struct {
	GameID      string
	WhitePlayer *websocket.Client
	TimeControl string
}

// ActiveGame represents an ongoing game with both players.
type ActiveGame struct {
	GameID      string
	WhitePlayer *websocket.Client
	BlackPlayer *websocket.Client
	CurrentFEN  string
	MoveCount   int
	TimeControl string
}

// NewGameManager creates a new game manager.
func NewGameManager(db *database.DB, chessService GameService) *GameManager {
	return &GameManager{
		db:           db,
		chessService: chessService,
		pendingGames: make(map[string]*PendingGame),
		activeGames:  make(map[string]*ActiveGame),
	}
}

// CreatePendingGame adds a new pending game.
func (m *GameManager) CreatePendingGame(gameID string, creator *websocket.Client, timeControl string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.pendingGames[gameID] = &PendingGame{
		GameID:      gameID,
		WhitePlayer: creator,
		TimeControl: timeControl,
	}
}

// GetPendingGame retrieves a pending game.
func (m *GameManager) GetPendingGame(gameID string) (*PendingGame, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	game, exists := m.pendingGames[gameID]
	return game, exists
}

// RemovePendingGame removes a pending game.
func (m *GameManager) RemovePendingGame(gameID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.pendingGames, gameID)
}

// ActivateGame moves a game from pending to active.
func (m *GameManager) ActivateGame(gameID string, blackPlayer *websocket.Client) (*ActiveGame, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	pending, exists := m.pendingGames[gameID]
	if !exists {
		return nil, fmt.Errorf("pending game not found")
	}

	// Create active game
	activeGame := &ActiveGame{
		GameID:      gameID,
		WhitePlayer: pending.WhitePlayer,
		BlackPlayer: blackPlayer,
		CurrentFEN:  m.chessService.GetStartingPosition(),
		MoveCount:   0,
		TimeControl: pending.TimeControl,
	}

	m.activeGames[gameID] = activeGame
	delete(m.pendingGames, gameID)

	return activeGame, nil
}

// GetActiveGame retrieves an active game.
func (m *GameManager) GetActiveGame(gameID string) (*ActiveGame, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	game, exists := m.activeGames[gameID]
	return game, exists
}

// UpdateGameState updates the FEN and move count.
func (m *GameManager) UpdateGameState(gameID string, newFEN string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	game, exists := m.activeGames[gameID]
	if !exists {
		return fmt.Errorf("active game not found")
	}

	game.CurrentFEN = newFEN
	game.MoveCount++

	return nil
}

// RemoveActiveGame removes an active game.
func (m *GameManager) RemoveActiveGame(gameID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.activeGames, gameID)
}

// GetPlayerSide returns the side (white/black) for a client in a game.
func (m *GameManager) GetPlayerSide(gameID string, client *websocket.Client) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	game, exists := m.activeGames[gameID]
	if !exists {
		return "", fmt.Errorf("game not found")
	}

	if game.WhitePlayer.ID == client.ID {
		return "white", nil
	}
	if game.BlackPlayer.ID == client.ID {
		return "black", nil
	}

	return "", fmt.Errorf("client not in game")
}

// GetOpponent returns the opponent client.
func (m *GameManager) GetOpponent(gameID string, client *websocket.Client) (*websocket.Client, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	game, exists := m.activeGames[gameID]
	if !exists {
		return nil, fmt.Errorf("game not found")
	}

	if game.WhitePlayer.ID == client.ID {
		return game.BlackPlayer, nil
	}
	if game.BlackPlayer.ID == client.ID {
		return game.WhitePlayer, nil
	}

	return nil, fmt.Errorf("client not in game")
}

// GetDB returns the database connection.
func (m *GameManager) GetDB() *database.DB {
	return m.db
}

// GetChessService returns the chess service.
func (m *GameManager) GetChessService() GameService {
	return m.chessService
}
```

---

## 📝 Step 3: Chess Service Implementation (30 minutes)

### File: `internal/websocket/handlers/chess_service.go`

**Create new file** - Chess logic abstraction:

```go
package handlers

import (
	"fmt"

	"github.com/notnil/chess"
)

// ChessService implements GameService using notnil/chess library.
type ChessService struct{}

// NewChessService creates a new chess service.
func NewChessService() *ChessService {
	return &ChessService{}
}

// ValidateMove implements GameService.ValidateMove.
func (s *ChessService) ValidateMove(fen string, moveSAN string) (MoveResult, error) {
	// Parse FEN
	fenFunc, err := chess.FEN(fen)
	if err != nil {
		return MoveResult{}, fmt.Errorf("invalid FEN: %w", err)
	}

	game := chess.NewGame(fenFunc)

	// Try to parse move as UCI first
	move, err := chess.UCINotation{}.Decode(game.Position(), moveSAN)
	if err != nil {
		// Try as SAN
		move, err = chess.AlgebraicNotation{}.Decode(game.Position(), moveSAN)
		if err != nil {
			return MoveResult{IsLegal: false}, fmt.Errorf("invalid move notation: %w", err)
		}
	}

	// Check if move is legal
	validMoves := game.ValidMoves()
	isLegal := false
	for _, validMove := range validMoves {
		if validMove == move {
			isLegal = true
			break
		}
	}

	if !isLegal {
		return MoveResult{IsLegal: false}, fmt.Errorf("illegal move")
	}

	// Get UCI and SAN notation
	uci := chess.UCINotation{}.Encode(game.Position(), move)
	san := chess.AlgebraicNotation{}.Encode(game.Position(), move)

	return MoveResult{
		SAN:     san,
		UCI:     uci,
		IsLegal: true,
	}, nil
}

// ApplyMove implements GameService.ApplyMove.
func (s *ChessService) ApplyMove(fen string, moveSAN string) (GameState, error) {
	// Parse FEN
	fenFunc, err := chess.FEN(fen)
	if err != nil {
		return GameState{}, fmt.Errorf("invalid FEN: %w", err)
	}

	game := chess.NewGame(fenFunc)

	// Parse and apply move
	move, err := chess.UCINotation{}.Decode(game.Position(), moveSAN)
	if err != nil {
		move, err = chess.AlgebraicNotation{}.Decode(game.Position(), moveSAN)
		if err != nil {
			return GameState{}, fmt.Errorf("invalid move: %w", err)
		}
	}

	if err := game.Move(move); err != nil {
		return GameState{}, fmt.Errorf("failed to apply move: %w", err)
	}

	// Get new state
	position := game.Position()
	newFEN := position.String()
	moveNumber := len(game.Moves())

	// Get move notation
	uci := chess.UCINotation{}.Encode(position, move)
	san := chess.AlgebraicNotation{}.Encode(position, move)

	// Check if game is over
	isGameOver, result := s.IsGameOver(newFEN)

	return GameState{
		FEN:        newFEN,
		SAN:        san,
		UCI:        uci,
		MoveNumber: moveNumber,
		IsGameOver: isGameOver,
		Result:     result,
	}, nil
}

// IsGameOver implements GameService.IsGameOver.
func (s *ChessService) IsGameOver(fen string) (bool, GameResult) {
	fenFunc, err := chess.FEN(fen)
	if err != nil {
		return false, GameResult{}
	}

	game := chess.NewGame(fenFunc)
	outcome := game.Outcome()

	if outcome == chess.NoOutcome {
		return false, GameResult{}
	}

	result := GameResult{
		Method: game.Method().String(),
		PGN:    game.String(),
	}

	switch outcome {
	case chess.WhiteWon:
		result.Winner = "white"
	case chess.BlackWon:
		result.Winner = "black"
	case chess.Draw:
		result.Winner = "draw"
	}

	return true, result
}

// GetStartingPosition implements GameService.GetStartingPosition.
func (s *ChessService) GetStartingPosition() string {
	return chess.StartingFEN()
}
```

---

## 📝 Step 4: CreateGame Handler (30 minutes)

### File: `internal/websocket/handlers/create_game.go`

**Create new file** - Single responsibility: Create game:

```go
package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/websocket"
	"github.com/jackc/pgx/v5/pgtype"
)

// CreateGameHandler handles game creation.
type CreateGameHandler struct {
	manager *GameManager
}

// Handle implements ActionHandler.Handle.
//
// Request format:
// {
//   "id": "req-123",
//   "type": "request",
//   "action": "createGame",
//   "data": {
//     "timeControl": "5+0"
//   }
// }
func (h *CreateGameHandler) Handle(ctx context.Context, client *websocket.Client, msg *websocket.Message) error {
	// 1. Extract and validate time control
	timeControl, err := websocket.RequireString(msg.Data, "timeControl")
	if err != nil {
		return err
	}

	seconds, err := websocket.ValidateTimeControl(timeControl)
	if err != nil {
		return fmt.Errorf("invalid time control: %w", err)
	}

	// 2. Generate game ID
	gameID := websocket.GenerateUUID()
	gameIDStr := websocket.UUIDToString(gameID)

	log.Printf("CreateGame: Creating game %s for user %s with time control %s (%ds)",
		gameIDStr, client.UserID, timeControl, seconds)

	// 3. Validate white player ID
	whitePlayerUUID := websocket.StringToUUID(client.UserID)
	if !whitePlayerUUID.Valid {
		return fmt.Errorf("invalid user ID")
	}

	// 4. Create game in database
	_, err = h.manager.GetDB().CreateGame(ctx, database.CreateGameParams{
		ID:            gameID,
		WhitePlayerID: whitePlayerUUID,
		BlackPlayerID: pgtype.UUID{Valid: false}, // NULL - waiting
		Pgn:           "",
		Result:        websocket.StringToText("*"),
		TimeControl:   websocket.IntToInt4(seconds),
	})
	if err != nil {
		log.Printf("CreateGame: Database error: %v", err)
		return fmt.Errorf("failed to create game")
	}

	// 5. Add to pending games
	h.manager.CreatePendingGame(gameIDStr, client, timeControl)

	// 6. Set client's game ID
	client.SetGameID(gameIDStr)

	log.Printf("CreateGame: Game %s created successfully, waiting for opponent", gameIDStr)

	// 7. Send success response
	websocket.SendSuccessResponse(client, msg.ID, map[string]interface{}{
		"gameId":      gameIDStr,
		"status":      "waiting",
		"side":        "white",
		"timeControl": timeControl,
		"fen":         h.manager.GetChessService().GetStartingPosition(),
	})

	return nil
}
```

---

## 📝 Step 5: JoinGame Handler (45 minutes)

### File: `internal/websocket/handlers/join_game.go`

**Create new file** - Single responsibility: Join game:

```go
package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/websocket"
)

// JoinGameHandler handles joining existing games.
type JoinGameHandler struct {
	manager *GameManager
}

// Handle implements ActionHandler.Handle.
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
func (h *JoinGameHandler) Handle(ctx context.Context, client *websocket.Client, msg *websocket.Message) error {
	// 1. Extract and validate game ID
	gameIDStr, err := websocket.RequireString(msg.Data, "gameId")
	if err != nil {
		return err
	}

	if !websocket.IsValidUUID(gameIDStr) {
		return fmt.Errorf("invalid game ID format")
	}

	log.Printf("JoinGame: User %s attempting to join game %s", client.UserID, gameIDStr)

	// 2. Get pending game
	pending, exists := h.manager.GetPendingGame(gameIDStr)
	if !exists {
		return fmt.Errorf("game not found or already started")
	}

	// 3. Verify not joining own game
	if pending.WhitePlayer.UserID == client.UserID {
		return fmt.Errorf("cannot join your own game")
	}

	// 4. Update game in database
	gameID := websocket.StringToUUID(gameIDStr)
	blackPlayerUUID := websocket.StringToUUID(client.UserID)
	if !blackPlayerUUID.Valid {
		return fmt.Errorf("invalid user ID")
	}

	_, err = h.manager.GetDB().UpdateGame(ctx, database.UpdateGameParams{
		ID:     gameID,
		Pgn:    "",
		Result: websocket.StringToText("*"),
	})
	if err != nil {
		log.Printf("JoinGame: Database error: %v", err)
		return fmt.Errorf("failed to join game")
	}

	// 5. Activate game
	activeGame, err := h.manager.ActivateGame(gameIDStr, client)
	if err != nil {
		return fmt.Errorf("failed to activate game: %w", err)
	}

	// 6. Set game IDs
	client.SetGameID(gameIDStr)
	pending.WhitePlayer.SetGameID(gameIDStr)

	// 7. Create room and add both players
	room := websocket.NewRoom(gameIDStr)
	pending.WhitePlayer.JoinRoom(room)
	client.JoinRoom(room)

	log.Printf("JoinGame: Game %s started - White: %s, Black: %s",
		gameIDStr, pending.WhitePlayer.UserID, client.UserID)

	// 8. Send success response to joiner
	websocket.SendSuccessResponse(client, msg.ID, map[string]interface{}{
		"gameId":      gameIDStr,
		"status":      "active",
		"side":        "black",
		"opponent":    pending.WhitePlayer.UserID,
		"timeControl": pending.TimeControl,
		"fen":         activeGame.CurrentFEN,
	})

	// 9. Notify white player
	opponentJoinedEvent := websocket.NewEvent("opponentJoined", map[string]interface{}{
		"gameId":   gameIDStr,
		"opponent": client.UserID,
		"side":     "black",
		"fen":      activeGame.CurrentFEN,
	})
	pending.WhitePlayer.SendMessage(opponentJoinedEvent)

	return nil
}
```

---

## 📝 Step 6: MakeMove Handler (60 minutes)

### File: `internal/websocket/handlers/make_move.go`

**Create new file** - Single responsibility: Process moves:

```go
package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/websocket"
	"github.com/jackc/pgx/v5/pgtype"
)

// MakeMoveHandler handles move processing.
type MakeMoveHandler struct {
	manager *GameManager
}

// Handle implements ActionHandler.Handle.
//
// Request format:
// {
//   "id": "req-789",
//   "type": "request",
//   "action": "makeMove",
//   "data": {
//     "gameId": "uuid",
//     "move": "e4"
//   }
// }
func (h *MakeMoveHandler) Handle(ctx context.Context, client *websocket.Client, msg *websocket.Message) error {
	// 1. Extract and validate inputs
	gameIDStr, err := websocket.RequireString(msg.Data, "gameId")
	if err != nil {
		return err
	}

	moveSAN, err := websocket.RequireString(msg.Data, "move")
	if err != nil {
		return err
	}

	if !websocket.ValidateMoveSAN(moveSAN) {
		return fmt.Errorf("invalid move format")
	}

	log.Printf("MakeMove: User %s attempting move '%s' in game %s",
		client.UserID, moveSAN, gameIDStr)

	// 2. Get active game
	activeGame, exists := h.manager.GetActiveGame(gameIDStr)
	if !exists {
		return fmt.Errorf("game not found or not active")
	}

	// 3. Verify it's player's turn
	playerSide, err := h.manager.GetPlayerSide(gameIDStr, client)
	if err != nil {
		return err
	}

	currentTurn := h.getCurrentTurn(activeGame.CurrentFEN)
	if playerSide != currentTurn {
		return fmt.Errorf("not your turn")
	}

	// 4. Apply move using chess service
	gameState, err := h.manager.GetChessService().ApplyMove(activeGame.CurrentFEN, moveSAN)
	if err != nil {
		return fmt.Errorf("invalid move: %w", err)
	}

	// 5. Update game state
	if err := h.manager.UpdateGameState(gameIDStr, gameState.FEN); err != nil {
		return fmt.Errorf("failed to update game state: %w", err)
	}

	log.Printf("MakeMove: Move applied - #%d %s (%s) by %s, gameOver=%v",
		gameState.MoveNumber, gameState.SAN, gameState.UCI, playerSide, gameState.IsGameOver)

	// 6. Save move to database
	gameID := websocket.StringToUUID(gameIDStr)
	_, err = h.manager.GetDB().CreateMove(ctx, database.CreateMoveParams{
		GameID:     gameID,
		MoveNumber: int32(gameState.MoveNumber),
		Side:       playerSide,
		MoveSan:    gameState.SAN,
		MoveUci:    gameState.UCI,
		Fen:        gameState.FEN,
		TimeTaken:  pgtype.Int4{Valid: false}, // TODO: Phase 5
	})
	if err != nil {
		log.Printf("MakeMove: Failed to save move: %v", err)
	}

	// 7. Update game if finished
	if gameState.IsGameOver {
		result := h.getResult(gameState.Result)
		_, err := h.manager.GetDB().UpdateGame(ctx, database.UpdateGameParams{
			ID:     gameID,
			Pgn:    gameState.Result.PGN,
			Result: websocket.StringToText(result),
		})
		if err != nil {
			log.Printf("MakeMove: Failed to update game result: %v", err)
		}
	}

	// 8. Prepare response data
	responseData := map[string]interface{}{
		"moveNumber": gameState.MoveNumber,
		"san":        gameState.SAN,
		"uci":        gameState.UCI,
		"fen":        gameState.FEN,
		"gameOver":   gameState.IsGameOver,
	}

	if gameState.IsGameOver {
		responseData["result"] = h.getResult(gameState.Result)
		responseData["method"] = gameState.Result.Method
	}

	// 9. Send success response
	websocket.SendSuccessResponse(client, msg.ID, responseData)

	// 10. Broadcast to opponent
	opponent, err := h.manager.GetOpponent(gameIDStr, client)
	if err == nil {
		moveMadeEvent := websocket.NewEvent("moveMade", responseData)
		opponent.SendMessage(moveMadeEvent)
	}

	// 11. Clean up if game over
	if gameState.IsGameOver {
		h.manager.RemoveActiveGame(gameIDStr)
	}

	return nil
}

// getCurrentTurn determines whose turn it is from FEN.
func (h *MakeMoveHandler) getCurrentTurn(fen string) string {
	// FEN format: "pieces active_color ..."
	// Second field is 'w' or 'b'
	parts := strings.Fields(fen)
	if len(parts) >= 2 {
		if parts[1] == "w" {
			return "white"
		}
		return "black"
	}
	return "white" // Default
}

// getResult converts GameResult to result string.
func (h *MakeMoveHandler) getResult(result GameResult) string {
	switch result.Winner {
	case "white":
		return "1-0"
	case "black":
		return "0-1"
	case "draw":
		return "1/2-1/2"
	default:
		return "*"
	}
}
```

---

## 📝 Step 7: Resign Handler (30 minutes)

### File: `internal/websocket/handlers/resign.go`

**Create new file** - Single responsibility: Handle resignations:

```go
package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/websocket"
)

// ResignHandler handles game resignations.
type ResignHandler struct {
	manager *GameManager
}

// Handle implements ActionHandler.Handle for resignation.
func (h *ResignHandler) Handle(ctx context.Context, client *websocket.Client, msg *websocket.Message) error {
	gameIDStr, err := websocket.RequireString(msg.Data, "gameId")
	if err != nil {
		return err
	}

	log.Printf("Resign: User %s resigning from game %s", client.UserID, gameIDStr)

	activeGame, exists := h.manager.GetActiveGame(gameIDStr)
	if !exists {
		return fmt.Errorf("game not found or not active")
	}

	return h.handleForfeit(ctx, client, activeGame, msg.ID)
}

// handleForfeit processes game forfeit.
func (h *ResignHandler) handleForfeit(ctx context.Context, client *websocket.Client, activeGame *ActiveGame, requestID string) error {
	gameIDStr := activeGame.GameID
	gameID := websocket.StringToUUID(gameIDStr)

	// Determine result
	playerSide, err := h.manager.GetPlayerSide(gameIDStr, client)
	if err != nil {
		return err
	}

	var result string
	var winner *websocket.Client
	var loser *websocket.Client

	if playerSide == "white" {
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
	_, err = h.manager.GetDB().UpdateGame(ctx, database.UpdateGameParams{
		ID:     gameID,
		Pgn:    "", // TODO: Get PGN from game state
		Result: websocket.StringToText(result),
	})
	if err != nil {
		log.Printf("Forfeit: Failed to update game: %v", err)
	}

	// Send response to resigner
	websocket.SendSuccessResponse(loser, requestID, map[string]interface{}{
		"message":  "You resigned",
		"result":   result,
		"gameOver": true,
	})

	// Notify winner
	opponentResignedEvent := websocket.NewEvent("opponentResigned", map[string]interface{}{
		"gameId":   gameIDStr,
		"result":   result,
		"gameOver": true,
		"message":  "Opponent resigned",
	})
	winner.SendMessage(opponentResignedEvent)

	// Clean up
	h.manager.RemoveActiveGame(gameIDStr)
	loser.SetGameID("")
	winner.SetGameID("")

	return nil
}

// LeaveGameHandler handles leaving games.
type LeaveGameHandler struct {
	manager *GameManager
}

// Handle implements ActionHandler.Handle for leaving.
func (h *LeaveGameHandler) Handle(ctx context.Context, client *websocket.Client, msg *websocket.Message) error {
	gameIDStr, err := websocket.RequireString(msg.Data, "gameId")
	if err != nil {
		return err
	}

	log.Printf("LeaveGame: User %s leaving game %s", client.UserID, gameIDStr)

	// Check pending game
	if pending, exists := h.manager.GetPendingGame(gameIDStr); exists {
		if pending.WhitePlayer.ID == client.ID {
			h.manager.RemovePendingGame(gameIDStr)
			client.SetGameID("")

			websocket.SendSuccessResponse(client, msg.ID, map[string]interface{}{
				"message": "Game cancelled",
			})
			return nil
		}
	}

	// Check active game - treat as resignation
	if activeGame, exists := h.manager.GetActiveGame(gameIDStr); exists {
		resignHandler := &ResignHandler{manager: h.manager}
		return resignHandler.handleForfeit(ctx, client, activeGame, msg.ID)
	}

	return fmt.Errorf("game not found")
}
```

---

## 📝 Step 8: Wire Everything Together (15 minutes)

### Update `internal/app/app.go`:

```go
import (
	// ... other imports
	"github.com/ankits1626/chess-coach-backend/internal/websocket"
	"github.com/ankits1626/chess-coach-backend/internal/websocket/handlers"
)

// New creates new application.
func New(cfg *config.Config, db *database.DB, log logger.Logger) *App {
	// Create chess service
	chessService := handlers.NewChessService()

	// Create game manager
	gameManager := handlers.NewGameManager(db, chessService)

	// Create handler router
	handler := handlers.NewHandlerRouter(gameManager)

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

## ✅ Phase 4 Checklist

**Architecture Setup**
- [ ] Chess library already installed ✅ (in go.mod)
- [ ] Create `handlers/` subfolder
- [ ] Create `handlers/handler.go` (interfaces + router)
- [ ] Create `handlers/game_manager.go` (state management)
- [ ] Create `handlers/chess_service.go` (chess abstraction)

**Handler Implementation**
- [ ] Create `handlers/create_game.go`
- [ ] Create `handlers/join_game.go`
- [ ] Create `handlers/make_move.go`
- [ ] Create `handlers/resign.go`

**Integration**
- [ ] Wire all components in `app.go`
- [ ] Test each handler individually
- [ ] Test complete game flow
- [ ] Verify error handling

---

## 📊 File Organization

| File | Lines | Responsibility | Principle |
|------|-------|----------------|-----------|
| `handlers/handler.go` | ~100 | Routing & interfaces | ISP, DIP |
| `handlers/game_manager.go` | ~150 | State management | SRP |
| `handlers/chess_service.go` | ~100 | Chess logic | DIP |
| `handlers/create_game.go` | ~80 | Create game | SRP |
| `handlers/join_game.go` | ~100 | Join game | SRP |
| `handlers/make_move.go` | ~150 | Move processing | SRP |
| `handlers/resign.go` | ~100 | Resignations | SRP |
| **Total** | **~780** | - | SOLID ✅ |

**All handlers organized in `internal/websocket/handlers/` subfolder** ✅

---

## 🎯 Benefits of This Architecture

### 1. **Single Responsibility (SRP)**
- Each handler has ONE reason to change
- Easy to understand and maintain
- Clear separation of concerns

### 2. **Open/Closed (OCP)**
- Add new handlers without modifying existing code
- Just implement `ActionHandler` interface

### 3. **Liskov Substitution (LSP)**
- Can swap `ChessService` with mock for testing
- All handlers depend on interfaces

### 4. **Interface Segregation (ISP)**
- Small, focused interfaces
- Handlers only depend on what they need

### 5. **Dependency Inversion (DIP)**
- Handlers depend on `GameManager` interface
- Chess logic behind `GameService` interface
- Easy to test and mock

---

## 🧪 Testing Benefits

With this architecture, you can easily:

```go
// Mock chess service for testing
type MockChessService struct{}

func (m *MockChessService) ValidateMove(fen, move string) (MoveResult, error) {
	return MoveResult{SAN: "e4", UCI: "e2e4", IsLegal: true}, nil
}

// Test handler in isolation
func TestCreateGameHandler(t *testing.T) {
	mockDB := &MockDB{}
	mockChess := &MockChessService{}
	manager := NewGameManager(mockDB, mockChess)
	handler := &CreateGameHandler{manager: manager}

	// Test...
}
```

---

## 🚀 Next Steps

After implementing all handlers:

1. Test each handler individually
2. Test complete game flow
3. Move to Phase 5 (Authentication)

---

**Ready for Phase 5?** → [05-authentication.md](./05-authentication.md)

---

**Last Updated**: 2025-11-05
**Status**: ✅ SOLID architecture ready
**Estimated Time**: 6-8 hours
**Dependencies**: Phases 1-3 complete

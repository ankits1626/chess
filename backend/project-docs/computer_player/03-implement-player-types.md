# Step 3: Implement Player Types

**Duration**: 2-3 hours

**Status**: ⬜ Not Started

**Goal**: Implement HumanPlayer and ComputerPlayer with clean, testable code

---

## 🎯 What You'll Build

```
internal/websocket/handlers/player/
├── types.go              # ✅ Already exists
├── game_mode.go          # ✅ Already exists (needs fix)
├── player.go             # ✅ Already exists
├── factory.go            # ✅ Already exists (needs refactor)
├── config.go             # 🆕 NEW - Separate config types
├── human_player.go       # 🆕 NEW - Human player implementation
├── computer_player.go    # 🆕 NEW - Computer player implementation
└── ai_service.go         # 🆕 NEW - Stockfish UCI wrapper
```

---

## 📋 Overview

### Part 1: Refactor Configs (30 min)
Fix the Interface Segregation Principle violation we identified in [IMPROVEMENTS.md](./IMPROVEMENTS.md#issue-2-playerconfig-violates-isp).

### Part 2: Implement AIService (45 min)
Create Stockfish UCI wrapper using `github.com/notnil/chess/uci`.

### Part 3: Implement HumanPlayer (30 min)
Wrap existing WebSocket client behavior.

### Part 4: Implement ComputerPlayer (45 min)
Integrate with AIService to get moves from Stockfish.

### Part 5: Update Factory (15 min)
Use new config types and player implementations.

---

## 🏗️ Part 1: Refactor Player Configs

### Problem Recap

Current `PlayerConfig` violates ISP:
```go
// ❌ BAD - All players see all config
type PlayerConfig struct {
    Client     *websocket.Client  // Only Human needs
    Difficulty string             // Only Computer needs
    APIKey     string             // Only LLM needs
    Model      string             // Only LLM needs
    RemoteURL  string             // Only Network needs
}
```

### Solution: Separate Configs

**File**: `internal/websocket/handlers/player/config.go`

```go
package player

import (
	"fmt"

	"github.com/ankits1626/chess-coach-backend/internal/websocket"
)

// PlayerConfig is the base interface for all player configurations
type PlayerConfig interface {
	Validate() error
	GetPlayerType() PlayerType
}

// HumanPlayerConfig configures a human player
type HumanPlayerConfig struct {
	Client *websocket.Client
}

func (c HumanPlayerConfig) Validate() error {
	if c.Client == nil {
		return fmt.Errorf("client required for human player")
	}
	return nil
}

func (c HumanPlayerConfig) GetPlayerType() PlayerType {
	return PlayerTypeHuman
}

// ComputerPlayerConfig configures a computer player
type ComputerPlayerConfig struct {
	Difficulty string // "easy", "medium", "hard"
}

func (c ComputerPlayerConfig) Validate() error {
	if c.Difficulty == "" {
		c.Difficulty = "medium" // default
	}

	valid := c.Difficulty == "easy" ||
	         c.Difficulty == "medium" ||
	         c.Difficulty == "hard"

	if !valid {
		return fmt.Errorf("invalid difficulty: %s (must be easy, medium, or hard)", c.Difficulty)
	}

	return nil
}

func (c ComputerPlayerConfig) GetPlayerType() PlayerType {
	return PlayerTypeComputer
}

// Future configs can be added here...

// LLMPlayerConfig configures an LLM-powered player (future)
type LLMPlayerConfig struct {
	APIKey string
	Model  string
}

func (c LLMPlayerConfig) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("API key required for LLM player")
	}
	if c.Model == "" {
		c.Model = "gpt-4" // default
	}
	return nil
}

func (c LLMPlayerConfig) GetPlayerType() PlayerType {
	return PlayerTypeLLM
}
```

**Why this is better**:
- ✅ Each config has only what it needs (ISP)
- ✅ Type-safe - can't pass wrong config to wrong player
- ✅ Validation at config level
- ✅ Easy to add new player types

---

## 🏗️ Part 2: Implement AIService

### Goal
Create a wrapper around Stockfish using UCI protocol.

### Dependencies

First, add the UCI package:
```bash
go get github.com/notnil/chess/uci
```

### File: `internal/websocket/handlers/player/ai_service.go`

```go
package player

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
)

// AIService provides AI move calculation using chess engines
type AIService interface {
	GetBestMove(ctx context.Context, game *chess.Game, difficulty string) (*chess.Move, error)
	Close() error
}

// StockfishService implements AIService using Stockfish chess engine
type StockfishService struct {
	engine     *uci.Engine
	enginePath string
}

// NewStockfishService creates a new Stockfish AI service
func NewStockfishService() (*StockfishService, error) {
	// Get Stockfish path from environment
	stockfishPath := os.Getenv("STOCKFISH_PATH")

	// Fallback to common paths if not set
	if stockfishPath == "" {
		paths := []string{
			"/usr/local/bin/stockfish",    // Our Docker build location
			"/usr/games/stockfish",        // Alpine package location
			"/opt/homebrew/bin/stockfish", // macOS Apple Silicon
			"/usr/local/bin/stockfish",    // macOS Intel
			"/usr/bin/stockfish",          // Linux
		}

		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				stockfishPath = path
				break
			}
		}

		if stockfishPath == "" {
			return nil, fmt.Errorf("stockfish not found, set STOCKFISH_PATH environment variable")
		}
	}

	// Create UCI engine
	eng, err := uci.New(stockfishPath)
	if err != nil {
		return nil, fmt.Errorf("failed to start stockfish at %s: %w", stockfishPath, err)
	}

	// Initialize engine
	if err := eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
		return nil, fmt.Errorf("failed to initialize stockfish: %w", err)
	}

	return &StockfishService{
		engine:     eng,
		enginePath: stockfishPath,
	}, nil
}

// GetBestMove calculates the best move for the current position
func (s *StockfishService) GetBestMove(ctx context.Context, game *chess.Game, difficulty string) (*chess.Move, error) {
	// Set position
	cmdPos := uci.CmdPosition{Position: game.Position()}
	cmdGo := uci.CmdGo{MoveTime: s.getMoveTime(difficulty)}

	// Run calculation
	if err := s.engine.Run(cmdPos, cmdGo); err != nil {
		return nil, fmt.Errorf("stockfish calculation failed: %w", err)
	}

	// Get best move
	searchResults := s.engine.SearchResults()
	if searchResults == nil || searchResults.BestMove == uci.NoMove {
		return nil, fmt.Errorf("stockfish returned no move")
	}

	// Convert UCI move to chess.Move
	move := searchResults.BestMove

	// Find the move in legal moves
	for _, legalMove := range game.ValidMoves() {
		if legalMove.S1() == move.From && legalMove.S2() == move.To {
			// Handle promotion
			if move.Promotion != chess.NoPieceType {
				if legalMove.Promo() == move.Promotion {
					return legalMove, nil
				}
			} else {
				return legalMove, nil
			}
		}
	}

	return nil, fmt.Errorf("stockfish move not found in legal moves: %s", move.String())
}

// getMoveTime returns the thinking time based on difficulty
func (s *StockfishService) getMoveTime(difficulty string) time.Duration {
	switch difficulty {
	case "easy":
		return 100 * time.Millisecond  // Quick, weaker moves
	case "hard":
		return 2 * time.Second          // Stronger, more calculated
	default: // "medium"
		return 500 * time.Millisecond   // Balanced
	}
}

// Close stops the Stockfish engine
func (s *StockfishService) Close() error {
	if s.engine != nil {
		return s.engine.Close()
	}
	return nil
}
```

**Key Points**:
- ✅ Environment variable for Stockfish path
- ✅ Fallback to common locations
- ✅ Difficulty-based thinking time
- ✅ Proper error handling
- ✅ UCI protocol communication

---

## 🏗️ Part 3: Implement HumanPlayer

### Goal
Wrap the existing WebSocket client to implement the Player interface.

### File: `internal/websocket/handlers/player/human_player.go`

```go
package player

import (
	"context"
	"fmt"
	"time"

	"github.com/ankits1626/chess-coach-backend/internal/websocket"
)

// HumanPlayer represents a human player connected via WebSocket
type HumanPlayer struct {
	client *websocket.Client
}

// NewHumanPlayer creates a new human player
func NewHumanPlayer(client *websocket.Client) *HumanPlayer {
	return &HumanPlayer{
		client: client,
	}
}

// GetID returns the player's unique identifier
func (h *HumanPlayer) GetID() string {
	return h.client.ID
}

// GetType returns the player type
func (h *HumanPlayer) GetType() PlayerType {
	return PlayerTypeHuman
}

// RequestMove waits for the human player to make a move via WebSocket
func (h *HumanPlayer) RequestMove(ctx context.Context, fen string) (string, error) {
	// Create a channel to receive the move
	moveCh := make(chan string, 1)
	errCh := make(chan error, 1)

	// Set up a temporary handler for move responses
	// This would integrate with your existing WebSocket message handling
	// For now, this is a placeholder that shows the pattern

	// In practice, you'd:
	// 1. Register a pending move request for this client
	// 2. When the client sends a move message, fulfill this request
	// 3. Use the context for timeout handling

	select {
	case move := <-moveCh:
		return move, nil
	case err := <-errCh:
		return "", err
	case <-ctx.Done():
		return "", ErrMoveTimeout
	case <-time.After(5 * time.Minute): // Generous timeout for human
		return "", ErrMoveTimeout
	}
}

// SendMove sends a move notification to the human player
func (h *HumanPlayer) SendMove(move MoveNotification) error {
	// Send move to client via WebSocket
	return h.client.Send(websocket.Message{
		Type: "move",
		Data: map[string]interface{}{
			"gameId":     move.GameID,
			"moveSAN":    move.MoveSAN,
			"moveUCI":    move.MoveUCI,
			"fen":        move.FEN,
			"moveNumber": move.MoveNumber,
			"isGameOver": move.IsGameOver,
			"result":     move.Result,
		},
	})
}

// NotifyGameStart sends game start notification
func (h *HumanPlayer) NotifyGameStart(info GameStartInfo) error {
	return h.client.Send(websocket.Message{
		Type: "game_start",
		Data: map[string]interface{}{
			"gameId":      info.GameID,
			"mode":        info.Mode,
			"opponent":    info.Opponent,
			"yourSide":    info.YourSide,
			"timeControl": info.TimeControl,
			"startingFEN": info.StartingFEN,
		},
	})
}

// NotifyGameEnd sends game end notification
func (h *HumanPlayer) NotifyGameEnd(result GameResult) error {
	return h.client.Send(websocket.Message{
		Type: "game_end",
		Data: map[string]interface{}{
			"winner": result.Winner,
			"method": result.Method,
			"pgn":    result.PGN,
		},
	})
}

// IsConnected checks if the player is still connected
func (h *HumanPlayer) IsConnected() bool {
	return h.client != nil && h.client.IsConnected()
}

// Cleanup performs any necessary cleanup
func (h *HumanPlayer) Cleanup() error {
	// Human player cleanup is handled by WebSocket disconnect
	return nil
}
```

**Key Points**:
- ✅ Wraps existing WebSocket client
- ✅ Implements all Player interface methods
- ✅ Uses context for timeout handling
- ✅ Sends structured messages to client

**Note**: The `RequestMove` method needs integration with your existing WebSocket message handling. You'll need to:
1. Store pending move requests by client ID
2. When a move message arrives, fulfill the pending request
3. Handle timeouts gracefully

---

## 🏗️ Part 4: Implement ComputerPlayer

### Goal
Create an AI player that uses Stockfish to calculate moves.

### File: `internal/websocket/handlers/player/computer_player.go`

```go
package player

import (
	"context"
	"fmt"
	"log"

	"github.com/notnil/chess"
)

// ComputerPlayer represents an AI player powered by a chess engine
type ComputerPlayer struct {
	id         string
	aiService  AIService
	difficulty string
	currentGame *chess.Game
}

// NewComputerPlayer creates a new computer player
func NewComputerPlayer(aiService AIService, difficulty string) *ComputerPlayer {
	return &ComputerPlayer{
		id:         fmt.Sprintf("computer-%s", difficulty),
		aiService:  aiService,
		difficulty: difficulty,
	}
}

// GetID returns the player's unique identifier
func (c *ComputerPlayer) GetID() string {
	return c.id
}

// GetType returns the player type
func (c *ComputerPlayer) GetType() PlayerType {
	return PlayerTypeComputer
}

// RequestMove calculates the best move using the AI engine
func (c *ComputerPlayer) RequestMove(ctx context.Context, fen string) (string, error) {
	// Parse FEN to create game state
	fenFunc, err := chess.FEN(fen)
	if err != nil {
		return "", fmt.Errorf("invalid FEN: %w", err)
	}

	game := chess.NewGame(fenFunc)
	c.currentGame = game

	// Check if game is over
	if game.Outcome() != chess.NoOutcome {
		return "", fmt.Errorf("game is already over")
	}

	// Get best move from AI
	move, err := c.aiService.GetBestMove(ctx, game, c.difficulty)
	if err != nil {
		return "", fmt.Errorf("AI move calculation failed: %w", err)
	}

	// Return move in UCI format (e.g., "e2e4")
	return move.String(), nil
}

// SendMove receives move notifications (computer doesn't need to see moves)
func (c *ComputerPlayer) SendMove(move MoveNotification) error {
	// Computer player doesn't need move notifications
	// But we log for debugging
	log.Printf("[ComputerPlayer %s] Received move: %s", c.id, move.MoveSAN)
	return nil
}

// NotifyGameStart receives game start notification
func (c *ComputerPlayer) NotifyGameStart(info GameStartInfo) error {
	log.Printf("[ComputerPlayer %s] Game started: %s vs %s", c.id, info.YourSide, info.Opponent)
	return nil
}

// NotifyGameEnd receives game end notification
func (c *ComputerPlayer) NotifyGameEnd(result GameResult) error {
	log.Printf("[ComputerPlayer %s] Game ended: %s by %s", c.id, result.Winner, result.Method)
	return nil
}

// IsConnected always returns true for computer players
func (c *ComputerPlayer) IsConnected() bool {
	return true
}

// Cleanup performs any necessary cleanup
func (c *ComputerPlayer) Cleanup() error {
	c.currentGame = nil
	return nil
}
```

**Key Points**:
- ✅ Uses AIService to calculate moves
- ✅ Returns moves in UCI format
- ✅ Always "connected" (doesn't disconnect)
- ✅ Difficulty-based play strength
- ✅ Logging for debugging

---

## 🏗️ Part 5: Update Factory

### File: `internal/websocket/handlers/player/factory.go`

Replace the existing content with:

```go
package player

import (
	"fmt"
)

// PlayerFactory creates player instances based on type
type PlayerFactory struct {
	aiService AIService
}

// NewPlayerFactory creates a new player factory
func NewPlayerFactory(aiService AIService) *PlayerFactory {
	return &PlayerFactory{
		aiService: aiService,
	}
}

// CreatePlayer creates a player of the specified type with the given config
func (f *PlayerFactory) CreatePlayer(config PlayerConfig) (Player, error) {
	// Validate config
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Create player based on type
	switch config.GetPlayerType() {
	case PlayerTypeHuman:
		humanConfig, ok := config.(HumanPlayerConfig)
		if !ok {
			return nil, fmt.Errorf("expected HumanPlayerConfig, got %T", config)
		}
		return NewHumanPlayer(humanConfig.Client), nil

	case PlayerTypeComputer:
		computerConfig, ok := config.(ComputerPlayerConfig)
		if !ok {
			return nil, fmt.Errorf("expected ComputerPlayerConfig, got %T", config)
		}
		if f.aiService == nil {
			return nil, fmt.Errorf("AI service not configured")
		}
		return NewComputerPlayer(f.aiService, computerConfig.Difficulty), nil

	case PlayerTypeLLM:
		return nil, fmt.Errorf("LLM player not yet implemented")

	case PlayerTypeNetwork:
		return nil, fmt.Errorf("network player not yet implemented")

	default:
		return nil, fmt.Errorf("unsupported player type: %s", config.GetPlayerType())
	}
}
```

**What Changed**:
- ✅ Now accepts `PlayerConfig` interface (not concrete struct)
- ✅ Uses type assertion to get specific config
- ✅ Validates config before creating player
- ✅ Type-safe player creation

---

## 🏗️ Part 6: Fix GameMode Default

Quick fix from [IMPROVEMENTS.md](./IMPROVEMENTS.md#issue-1-gamemode-default-case).

### File: `internal/websocket/handlers/player/game_mode.go`

Change line 40:
```go
default:
    return PlayerTypeHuman, PlayerTypeComputer  // ✅ Better default
```

---

## ✅ Verification

### Test 1: Config Validation

```go
// Test in your test file or main
config := player.ComputerPlayerConfig{
    Difficulty: "invalid",
}

err := config.Validate()
// Should error: "invalid difficulty: invalid"

config.Difficulty = "hard"
err = config.Validate()
// Should succeed: err == nil
```

### Test 2: AIService

```go
// Create AIService
aiService, err := player.NewStockfishService()
if err != nil {
    log.Fatal(err)
}
defer aiService.Close()

// Test move calculation
game := chess.NewGame()
ctx := context.Background()

move, err := aiService.GetBestMove(ctx, game, "medium")
if err != nil {
    log.Fatal(err)
}

log.Printf("Stockfish suggests: %s", move)
```

### Test 3: Factory Usage

```go
// Create factory
aiService, _ := player.NewStockfishService()
factory := player.NewPlayerFactory(aiService)

// Create computer player
config := player.ComputerPlayerConfig{Difficulty: "hard"}
computerPlayer, err := factory.CreatePlayer(config)
if err != nil {
    log.Fatal(err)
}

log.Printf("Created %s player: %s",
    computerPlayer.GetType(),
    computerPlayer.GetID())
```

---

## 📊 File Structure After Step 3

```
internal/websocket/handlers/player/
├── types.go              # ✅ PlayerType enum
├── game_mode.go          # ✅ GameMode enum (fixed default)
├── player.go             # ✅ Player interface
├── config.go             # ✅ NEW - Separate config types
├── factory.go            # ✅ Updated - Uses new configs
├── ai_service.go         # ✅ NEW - Stockfish wrapper
├── human_player.go       # ✅ NEW - Human implementation
└── computer_player.go    # ✅ NEW - Computer implementation
```

**Total lines added**: ~400-500 lines of clean, tested code

---

## 🎯 Next Step

Once Step 3 is complete:

→ **[Step 4: Refactor GameManager](./04-refactor-game-manager.md)**

You'll integrate the player system into the existing game management infrastructure.

---

## 💡 Tips

1. **Start with tests**: Write unit tests for each component
2. **Test AIService first**: Make sure Stockfish integration works
3. **Use logs**: Add logging to track move calculations
4. **Handle errors gracefully**: AI can fail, handle it properly
5. **Keep it simple**: Don't over-engineer on first pass

---

**Estimated Time**: 2-3 hours
**Actual Time**: ________

**Last Updated**: 2025-11-05

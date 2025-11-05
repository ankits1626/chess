# Step 1: Design Player Interface

**Duration**: 1 hour

**Status**: ⬜ Not Started

**Goal**: Design the core Player interface following SOLID principles

---

## 🎯 What We're Building

A **Player interface** that can represent ANY chess player:
- Human (via WebSocket)
- Computer (via Stockfish)
- LLM (via API - future)
- Network Player (via API - future)

This is the **foundation** of our extensible architecture!

---

## 📋 Player Interface Design

### File: `internal/websocket/handlers/player/player.go`

```go
package player

import (
	"context"
	"errors"
)

// Player represents any chess player in the system.
// This interface allows for different player implementations:
// - HumanPlayer (WebSocket client)
// - ComputerPlayer (AI engine)
// - LLMPlayer (GPT-4, Claude, etc.)
type Player interface {
	// Core identification
	GetID() string
	GetType() PlayerType

	// Move handling
	// RequestMove asks the player for their next move given current position
	// - For Human: waits for WebSocket message
	// - For Computer: runs AI calculation
	// - For LLM: calls API
	RequestMove(ctx context.Context, fen string) (string, error)

	// SendMove notifies the player about a move (their own or opponent's)
	SendMove(move MoveNotification) error

	// Game events
	NotifyGameStart(info GameStartInfo) error
	NotifyGameEnd(result GameResult) error

	// Lifecycle
	IsConnected() bool
	Cleanup() error
}

// MoveNotification contains information about a move
type MoveNotification struct {
	GameID     string
	MoveSAN    string
	MoveUCI    string
	FEN        string
	MoveNumber int
	IsGameOver bool
	Result     *GameResult // nil if game not over
}

// GameStartInfo contains game start information
type GameStartInfo struct {
	GameID      string
	Mode        GameMode
	Opponent    string // Opponent ID or "Computer" or "GPT-4"
	YourSide    string // "white" or "black"
	TimeControl string
	StartingFEN string
}

// GameResult contains game result information
type GameResult struct {
	Winner string // "white", "black", "draw"
	Method string // "checkmate", "stalemate", "resignation", etc.
	PGN    string
}

// ErrPlayerDisconnected is returned when a player is not connected
var ErrPlayerDisconnected = errors.New("player disconnected")

// ErrMoveTimeout is returned when RequestMove times out
var ErrMoveTimeout = errors.New("move request timeout")
```

---

## 📋 Player Types

### File: `internal/websocket/handlers/player/types.go`

```go
package player

// PlayerType represents the type of player
type PlayerType string

const (
	PlayerTypeHuman    PlayerType = "human"
	PlayerTypeComputer PlayerType = "computer"
	PlayerTypeLLM      PlayerType = "llm"
	PlayerTypeNetwork  PlayerType = "network" // Future: remote player via API
)

// String returns string representation
func (pt PlayerType) String() string {
	return string(pt)
}

// IsValid checks if player type is valid
func (pt PlayerType) IsValid() bool {
	switch pt {
	case PlayerTypeHuman, PlayerTypeComputer, PlayerTypeLLM, PlayerTypeNetwork:
		return true
	}
	return false
}
```

---

## 📋 Game Modes

### File: `internal/websocket/handlers/player/game_mode.go`

```go
package player

// GameMode defines who plays against whom
type GameMode string

const (
	GameModeHumanVsHuman       GameMode = "human_vs_human"
	GameModeHumanVsComputer    GameMode = "human_vs_computer"
	GameModeComputerVsComputer GameMode = "computer_vs_computer"
	GameModeHumanVsLLM         GameMode = "human_vs_llm"
)

// String returns string representation
func (gm GameMode) String() string {
	return string(gm)
}

// IsValid checks if game mode is valid
func (gm GameMode) IsValid() bool {
	switch gm {
	case GameModeHumanVsHuman, GameModeHumanVsComputer,
	     GameModeComputerVsComputer, GameModeHumanVsLLM:
		return true
	}
	return false
}

// GetPlayerTypes returns the two player types for this mode
func (gm GameMode) GetPlayerTypes() (PlayerType, PlayerType) {
	switch gm {
	case GameModeHumanVsHuman:
		return PlayerTypeHuman, PlayerTypeHuman
	case GameModeHumanVsComputer:
		return PlayerTypeHuman, PlayerTypeComputer
	case GameModeComputerVsComputer:
		return PlayerTypeComputer, PlayerTypeComputer
	case GameModeHumanVsLLM:
		return PlayerTypeHuman, PlayerTypeLLM
	default:
		return PlayerTypeHuman, PlayerTypeHuman
	}
}
```

---

## 📋 Player Factory

### File: `internal/websocket/handlers/player/factory.go`

```go
package player

import (
	"fmt"

	"github.com/ankits1626/chess-coach-backend/internal/websocket"
)

// PlayerFactory creates player instances based on type
type PlayerFactory struct {
	aiService AIService // Will be defined in Step 3
}

// NewPlayerFactory creates a new player factory
func NewPlayerFactory(aiService AIService) *PlayerFactory {
	return &PlayerFactory{
		aiService: aiService,
	}
}

// PlayerConfig holds configuration for creating a player
type PlayerConfig struct {
	// For Human players
	Client *websocket.Client

	// For Computer players
	Difficulty string // "easy", "medium", "hard"

	// For LLM players (future)
	APIKey string
	Model  string

	// For Network players (future)
	RemoteURL string
}

// CreatePlayer creates a player of the specified type
func (f *PlayerFactory) CreatePlayer(playerType PlayerType, config PlayerConfig) (Player, error) {
	if !playerType.IsValid() {
		return nil, fmt.Errorf("invalid player type: %s", playerType)
	}

	switch playerType {
	case PlayerTypeHuman:
		if config.Client == nil {
			return nil, fmt.Errorf("client required for human player")
		}
		return NewHumanPlayer(config.Client), nil

	case PlayerTypeComputer:
		if f.aiService == nil {
			return nil, fmt.Errorf("AI service not configured")
		}
		if config.Difficulty == "" {
			config.Difficulty = "medium"
		}
		return NewComputerPlayer(f.aiService, config.Difficulty), nil

	case PlayerTypeLLM:
		return nil, fmt.Errorf("LLM player not yet implemented")

	case PlayerTypeNetwork:
		return nil, fmt.Errorf("network player not yet implemented")

	default:
		return nil, fmt.Errorf("unsupported player type: %s", playerType)
	}
}

// AIService interface (will be properly defined in Step 3)
type AIService interface {
	GetBestMove(ctx context.Context, fen string, difficulty string) (string, error)
	Close() error
}
```

---

## 🎨 SOLID Principles Verification

### ✅ Single Responsibility Principle (SRP)
- `Player` interface: Represents a player's capabilities
- `PlayerFactory`: Creates player instances
- `PlayerType`: Defines player types
- `GameMode`: Defines game modes

Each type has ONE reason to change.

### ✅ Open/Closed Principle (OCP)
```go
// Adding new player type (LLM):
// 1. Create LLMPlayer struct
// 2. Implement Player interface
// 3. Add case to factory

// NO changes to:
// - Player interface ✅
// - Existing implementations ✅
// - GameManager ✅
```

### ✅ Liskov Substitution Principle (LSP)
```go
func handleMove(player Player, fen string) {
    move, _ := player.RequestMove(ctx, fen)
    // Works for ANY Player implementation
}
```

### ✅ Interface Segregation Principle (ISP)
- Player interface has only essential methods
- No method forced on implementations that don't need it

### ✅ Dependency Inversion Principle (DIP)
```go
// GameManager depends on Player (abstraction)
type GameManager struct {
    whitePlayer Player  // NOT HumanPlayer or ComputerPlayer
    blackPlayer Player
}
```

---

## 📁 Implementation Steps

### 1. Create Player Package Directory

```bash
cd /Users/ankit/code/learn/chess-coach/backend

mkdir -p internal/websocket/handlers/player
```

### 2. Create Files

Create these files in order:

1. **types.go** - Player types enum
2. **game_mode.go** - Game mode enum
3. **player.go** - Player interface
4. **factory.go** - Player factory (with stub AIService)

### 3. Verify Compilation

```bash
cd /Users/ankit/code/learn/chess-coach/backend

# Should compile (even with stub implementations)
go build ./internal/websocket/handlers/player/...
```

---

## ✅ Verification Checklist

- [ ] Created `internal/websocket/handlers/player/` directory
- [ ] Created `types.go` with PlayerType enum
- [ ] Created `game_mode.go` with GameMode enum
- [ ] Created `player.go` with Player interface
- [ ] Created `factory.go` with PlayerFactory
- [ ] Code compiles without errors
- [ ] All SOLID principles verified

---

## 📊 What This Enables

With this interface in place, we can now:

1. **Step 2-3**: Implement `HumanPlayer` and `ComputerPlayer`
2. **Step 4**: Update GameManager to use `Player` interface
3. **Step 5**: Create games with any mode
4. **Future**: Add `LLMPlayer` by just implementing the interface

---

## 🎯 Design Benefits

### Before (Tightly Coupled)
```go
type Game struct {
    whiteClient *websocket.Client  // Must be human
    blackClient *websocket.Client  // Must be human
}
```

### After (Loosely Coupled)
```go
type Game struct {
    whitePlayer Player  // Can be ANY player type
    blackPlayer Player  // Can be ANY player type
}
```

---

## 🚀 Usage Example

```go
// Create players via factory
factory := NewPlayerFactory(aiService)

// Human vs Computer game
humanPlayer, _ := factory.CreatePlayer(PlayerTypeHuman, PlayerConfig{
    Client: websocketClient,
})

computerPlayer, _ := factory.CreatePlayer(PlayerTypeComputer, PlayerConfig{
    Difficulty: "medium",
})

// Create game with mode
game := NewGame(GameModeHumanVsComputer, humanPlayer, computerPlayer)

// Request move works for BOTH
humanMove, _ := humanPlayer.RequestMove(ctx, currentFEN)
computerMove, _ := computerPlayer.RequestMove(ctx, currentFEN)
```

---

## 📝 Complete File Structure After This Step

```
internal/websocket/handlers/
├── player/                     # NEW
│   ├── player.go              # Player interface + types
│   ├── types.go               # PlayerType enum
│   ├── game_mode.go           # GameMode enum
│   └── factory.go             # PlayerFactory
```

---

## 🎯 Next Step

Once this step is complete and verified:

→ **[Step 2: Install & Test Stockfish](./02-install-stockfish.md)**

We'll install the AI engine that will power ComputerPlayer.

---

## 💡 Key Takeaway

This step creates the **abstraction layer** that makes everything else possible. Get this right, and the rest becomes straightforward implementation!

**Core Principle**: Program to an interface, not an implementation.

---

**Estimated Time**: 1 hour
**Actual Time**: ________

---

**Last Updated**: 2025-11-05

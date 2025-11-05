# Step 4: Refactor GameManager

**Goal**: Integrate the Player system into GameManager to support different game modes (Human vs Computer, Human vs Human, etc.)

**Estimated Time**: 2 hours

---

## 📋 Overview

Currently, `GameManager` assumes all games are Human vs Human:
- `PendingGame` stores a single `WhitePlayer *websocket.Client`
- `ActiveGame` stores `WhitePlayer` and `BlackPlayer` as `*websocket.Client`
- No concept of game modes or player types

**What We'll Do**:
1. Add game mode support to `PendingGame` and `ActiveGame`
2. Store players as `player.Player` interface instead of `*websocket.Client`
3. Add AIService dependency for computer players
4. Update game creation to accept game mode
5. Update game activation to create appropriate players
6. Add helper methods for player management

---

## 🎯 Current vs Desired State

### Current State (Human vs Human only):

```go
type PendingGame struct {
    GameID      string
    WhitePlayer *websocket.Client  // Always human
    TimeControl string
}

type ActiveGame struct {
    GameID      string
    WhitePlayer *websocket.Client  // Always human
    BlackPlayer *websocket.Client  // Always human
    CurrentFEN  string
    MoveCount   int
    TimeControl string
}
```

### Desired State (Any game mode):

```go
type PendingGame struct {
    GameID      string
    Mode        player.GameMode     // NEW: Game mode
    WhitePlayer player.Player       // NEW: Player interface
    Difficulty  string              // NEW: For computer difficulty
    TimeControl string
}

type ActiveGame struct {
    GameID      string
    Mode        player.GameMode     // NEW: Game mode
    WhitePlayer player.Player       // NEW: Player interface
    BlackPlayer player.Player       // NEW: Player interface
    CurrentFEN  string
    MoveCount   int
    TimeControl string
}
```

---

## 📝 Implementation Plan

### Part 1: Update GameManager Struct (10 min)

**File**: `internal/websocket/handlers/game_manager.go`

**Add**:
1. Import player package
2. Add `playerFactory *player.PlayerFactory` field
3. Add `aiService player.AIService` field
4. Update `NewGameManager` constructor

**Before**:
```go
type GameManager struct {
    db           *database.DB
    chessService GameService
    pendingGames map[string]*PendingGame
    activeGames  map[string]*ActiveGame
    mu           sync.RWMutex
}

func NewGameManager(db *database.DB, chessService GameService) *GameManager {
    return &GameManager{
        db:           db,
        chessService: chessService,
        pendingGames: make(map[string]*PendingGame),
        activeGames:  make(map[string]*ActiveGame),
    }
}
```

**After**:
```go
type GameManager struct {
    db            *database.DB
    chessService  GameService
    aiService     player.AIService      // NEW
    playerFactory *player.PlayerFactory // NEW
    pendingGames  map[string]*PendingGame
    activeGames   map[string]*ActiveGame
    mu            sync.RWMutex
}

func NewGameManager(db *database.DB, chessService GameService, aiService player.AIService) *GameManager {
    return &GameManager{
        db:            db,
        chessService:  chessService,
        aiService:     aiService,
        playerFactory: player.NewPlayerFactory(aiService), // NEW
        pendingGames:  make(map[string]*PendingGame),
        activeGames:   make(map[string]*ActiveGame),
    }
}
```

---

### Part 2: Update PendingGame and ActiveGame (15 min)

**File**: `internal/websocket/handlers/game_manager.go`

**Update structs**:

```go
type PendingGame struct {
    GameID      string
    Mode        player.GameMode // NEW
    WhitePlayer player.Player   // CHANGED from *websocket.Client
    Difficulty  string          // NEW (for computer difficulty)
    TimeControl string
}

type ActiveGame struct {
    GameID      string
    Mode        player.GameMode // NEW
    WhitePlayer player.Player   // CHANGED from *websocket.Client
    BlackPlayer player.Player   // CHANGED from *websocket.Client
    CurrentFEN  string
    MoveCount   int
    TimeControl string
}
```

---

### Part 3: Update CreatePendingGame Method (20 min)

**File**: `internal/websocket/handlers/game_manager.go`

**Current**:
```go
func (m *GameManager) CreatePendingGame(gameID string, creator *websocket.Client, timeControl string) {
    m.mu.Lock()
    defer m.mu.Unlock()

    m.pendingGames[gameID] = &PendingGame{
        GameID:      gameID,
        WhitePlayer: creator,
        TimeControl: timeControl,
    }
}
```

**New**:
```go
func (m *GameManager) CreatePendingGame(
    gameID string,
    mode player.GameMode,
    creator *websocket.Client,
    difficulty string,
    timeControl string,
) error {
    m.mu.Lock()
    defer m.mu.Unlock()

    // Determine white and black player types based on mode
    whiteType, blackType := mode.GetPlayerTypes()

    // Create white player
    var whitePlayer player.Player
    var err error

    if whiteType == player.PlayerTypeHuman {
        config := player.HumanPlayerConfig{Client: creator}
        whitePlayer, err = m.playerFactory.CreatePlayer(config)
        if err != nil {
            return fmt.Errorf("failed to create white player: %w", err)
        }
    } else if whiteType == player.PlayerTypeComputer {
        config := player.ComputerPlayerConfig{Difficulty: difficulty}
        whitePlayer, err = m.playerFactory.CreatePlayer(config)
        if err != nil {
            return fmt.Errorf("failed to create white player: %w", err)
        }
    }

    m.pendingGames[gameID] = &PendingGame{
        GameID:      gameID,
        Mode:        mode,
        WhitePlayer: whitePlayer,
        Difficulty:  difficulty,
        TimeControl: timeControl,
    }

    return nil
}
```

---

### Part 4: Update ActivateGame Method (25 min)

**File**: `internal/websocket/handlers/game_manager.go`

**Current**:
```go
func (m *GameManager) ActivateGame(gameID string, blackPlayer *websocket.Client) (*ActiveGame, error) {
    m.mu.Lock()
    defer m.mu.Unlock()

    pending, exists := m.pendingGames[gameID]
    if !exists {
        return nil, fmt.Errorf("pending game not found")
    }

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
```

**New**:
```go
func (m *GameManager) ActivateGame(gameID string, blackClient *websocket.Client) (*ActiveGame, error) {
    m.mu.Lock()
    defer m.mu.Unlock()

    pending, exists := m.pendingGames[gameID]
    if !exists {
        return nil, fmt.Errorf("pending game not found")
    }

    // Determine black player type from mode
    _, blackType := pending.Mode.GetPlayerTypes()

    // Create black player based on type
    var blackPlayer player.Player
    var err error

    if blackType == player.PlayerTypeHuman {
        if blackClient == nil {
            return nil, fmt.Errorf("black client required for human player")
        }
        config := player.HumanPlayerConfig{Client: blackClient}
        blackPlayer, err = m.playerFactory.CreatePlayer(config)
        if err != nil {
            return nil, fmt.Errorf("failed to create black player: %w", err)
        }
    } else if blackType == player.PlayerTypeComputer {
        config := player.ComputerPlayerConfig{Difficulty: pending.Difficulty}
        blackPlayer, err = m.playerFactory.CreatePlayer(config)
        if err != nil {
            return nil, fmt.Errorf("failed to create black player: %w", err)
        }
    }

    activeGame := &ActiveGame{
        GameID:      gameID,
        Mode:        pending.Mode,
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
```

---

### Part 5: Update Helper Methods (20 min)

**File**: `internal/websocket/handlers/game_manager.go`

#### 5.1 Update GetPlayerSide

**Current**:
```go
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
```

**New**:
```go
func (m *GameManager) GetPlayerSide(gameID string, client *websocket.Client) (string, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    game, exists := m.activeGames[gameID]
    if !exists {
        return "", fmt.Errorf("game not found")
    }

    if game.WhitePlayer.GetID() == client.ID {
        return "white", nil
    }
    if game.BlackPlayer.GetID() == client.ID {
        return "black", nil
    }

    return "", fmt.Errorf("client not in game")
}
```

#### 5.2 Update GetOpponent

**Current**:
```go
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
```

**New**:
```go
func (m *GameManager) GetOpponent(gameID string, client *websocket.Client) (player.Player, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    game, exists := m.activeGames[gameID]
    if !exists {
        return nil, fmt.Errorf("game not found")
    }

    if game.WhitePlayer.GetID() == client.ID {
        return game.BlackPlayer, nil
    }
    if game.BlackPlayer.GetID() == client.ID {
        return game.WhitePlayer, nil
    }

    return nil, fmt.Errorf("client not in game")
}
```

#### 5.3 Add New Helper Methods

**Add these new methods**:

```go
// GetPlayer returns the player for a given side
func (m *GameManager) GetPlayer(gameID string, side string) (player.Player, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    game, exists := m.activeGames[gameID]
    if !exists {
        return nil, fmt.Errorf("game not found")
    }

    if side == "white" {
        return game.WhitePlayer, nil
    } else if side == "black" {
        return game.BlackPlayer, nil
    }

    return nil, fmt.Errorf("invalid side: %s", side)
}

// GetGameMode returns the mode for a game
func (m *GameManager) GetGameMode(gameID string) (player.GameMode, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    game, exists := m.activeGames[gameID]
    if !exists {
        return "", fmt.Errorf("game not found")
    }

    return game.Mode, nil
}

// IsComputerTurn checks if it's a computer player's turn
func (m *GameManager) IsComputerTurn(gameID string, fen string) bool {
    m.mu.RLock()
    defer m.mu.RUnlock()

    game, exists := m.activeGames[gameID]
    if !exists {
        return false
    }

    // Parse FEN to determine whose turn it is
    // FEN format: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
    // The field after the position is the active color: 'w' or 'b'
    fields := strings.Fields(fen)
    if len(fields) < 2 {
        return false
    }

    activeColor := fields[1]
    if activeColor == "w" && game.WhitePlayer.GetType() == player.PlayerTypeComputer {
        return true
    }
    if activeColor == "b" && game.BlackPlayer.GetType() == player.PlayerTypeComputer {
        return true
    }

    return false
}
```

---

### Part 6: Add Computer Move Handler (30 min)

**File**: `internal/websocket/handlers/game_manager.go`

**Add this new method**:

```go
// HandleComputerMove triggers the computer to make a move if it's their turn
func (m *GameManager) HandleComputerMove(ctx context.Context, gameID string) error {
    m.mu.RLock()
    game, exists := m.activeGames[gameID]
    m.mu.RUnlock()

    if !exists {
        return fmt.Errorf("game not found")
    }

    // Check if it's computer's turn
    if !m.IsComputerTurn(gameID, game.CurrentFEN) {
        return nil // Not computer's turn, nothing to do
    }

    // Determine which player is the computer
    var computerPlayer player.Player
    var side string

    if game.WhitePlayer.GetType() == player.PlayerTypeComputer {
        computerPlayer = game.WhitePlayer
        side = "white"
    } else {
        computerPlayer = game.BlackPlayer
        side = "black"
    }

    // Request move from computer
    moveUCI, err := computerPlayer.RequestMove(ctx, game.CurrentFEN)
    if err != nil {
        log.Printf("Computer move failed for game %s: %v", gameID, err)
        return fmt.Errorf("computer move failed: %w", err)
    }

    log.Printf("Computer (%s) in game %s plays: %s", side, gameID, moveUCI)

    // Validate and apply the move
    newFEN, moveSAN, err := m.chessService.MakeMove(game.CurrentFEN, moveUCI)
    if err != nil {
        log.Printf("Invalid computer move %s for game %s: %v", moveUCI, gameID, err)
        return fmt.Errorf("invalid computer move: %w", err)
    }

    // Update game state
    if err := m.UpdateGameState(gameID, newFEN); err != nil {
        return fmt.Errorf("failed to update game state: %w", err)
    }

    // Check if game is over
    isGameOver := m.chessService.IsGameOver(newFEN)
    var result string
    if isGameOver {
        result = m.chessService.GetGameResult(newFEN)
        log.Printf("Game %s ended: %s", gameID, result)
    }

    // Notify both players
    moveNotification := player.MoveNotification{
        GameID:     gameID,
        MoveSAN:    moveSAN,
        MoveUCI:    moveUCI,
        FEN:        newFEN,
        MoveNumber: game.MoveCount,
        IsGameOver: isGameOver,
        Result:     result,
    }

    // Send to white player
    if err := game.WhitePlayer.SendMove(moveNotification); err != nil {
        log.Printf("Failed to send move to white player: %v", err)
    }

    // Send to black player
    if err := game.BlackPlayer.SendMove(moveNotification); err != nil {
        log.Printf("Failed to send move to black player: %v", err)
    }

    // If game is over, clean up
    if isGameOver {
        m.RemoveActiveGame(gameID)
    }

    return nil
}
```

---

## 🔧 Compile Fixes Needed

After these changes, you'll need to update files that use GameManager:

### Files to Update:

1. **create_game.go** - Update CreateGameHandler
2. **join_game.go** - Update JoinGameHandler
3. **make_move.go** - Update to trigger computer moves
4. **main.go** or server initialization - Pass AIService to NewGameManager

We'll handle these in **Step 5** (Create Mode Handler) and **Step 6** (Update MakeMove).

For now, the GameManager refactoring is complete!

---

## ✅ Verification Checklist

After implementation:

- [ ] GameManager has `aiService` and `playerFactory` fields
- [ ] NewGameManager accepts `aiService` parameter
- [ ] PendingGame has `Mode`, `Difficulty`, and `Player` interface
- [ ] ActiveGame has `Mode` and `Player` interfaces
- [ ] CreatePendingGame accepts mode and difficulty
- [ ] ActivateGame creates appropriate player types
- [ ] GetPlayerSide uses `Player.GetID()`
- [ ] GetOpponent returns `player.Player`
- [ ] New helper methods added (GetPlayer, GetGameMode, IsComputerTurn)
- [ ] HandleComputerMove implemented
- [ ] Code compiles (may have errors in other files - expected!)

---

## 📊 Summary

**Files Modified**: 1 file (`game_manager.go`)

**Lines Added**: ~150 lines

**New Methods**: 4 new methods

**Updated Methods**: 5 methods updated

**Breaking Changes**: Yes (other handlers will need updates in Steps 5 & 6)

---

## 🎯 Next Step

Once GameManager is refactored, move to:

**→ [Step 5: Create Mode Handler](./05-create-mode-handler.md)**

This will update the handlers (CreateGame, JoinGame) to use the new GameManager API with game modes.

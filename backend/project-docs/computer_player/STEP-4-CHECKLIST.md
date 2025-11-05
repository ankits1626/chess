# Step 4: Refactor GameManager - Checklist

**Goal**: Integrate Player system into GameManager for game mode support

**Estimated Time**: 2 hours

---

## 📋 Quick Overview

| Part | Task | Time | Status |
|------|------|------|--------|
| 1. Update GameManager Struct | Add AIService & PlayerFactory | 10 min | ⬜ |
| 2. Update Game Structs | Add Mode & Player interfaces | 15 min | ⬜ |
| 3. Update CreatePendingGame | Accept mode & create players | 20 min | ⬜ |
| 4. Update ActivateGame | Create black player by mode | 25 min | ⬜ |
| 5. Update Helper Methods | Use Player interface | 20 min | ⬜ |
| 6. Add Computer Move Handler | HandleComputerMove method | 30 min | ⬜ |

---

## 📝 Part 1: Update GameManager Struct (10 min)

**File**: `internal/websocket/handlers/game_manager.go`

- [ ] Add import for player package:
  ```go
  "github.com/ankits1626/chess-coach-backend/internal/websocket/handlers/player"
  ```
- [ ] Add `aiService player.AIService` field to GameManager
- [ ] Add `playerFactory *player.PlayerFactory` field to GameManager
- [ ] Update NewGameManager to accept `aiService player.AIService` parameter
- [ ] Initialize playerFactory: `playerFactory: player.NewPlayerFactory(aiService)`

**Verification**:
```bash
# Should have compilation errors in other files - that's expected!
go build ./internal/websocket/handlers/...
```

---

## 📝 Part 2: Update Game Structs (15 min)

**File**: `internal/websocket/handlers/game_manager.go`

### PendingGame Updates:
- [ ] Add `Mode player.GameMode` field
- [ ] Change `WhitePlayer *websocket.Client` to `WhitePlayer player.Player`
- [ ] Add `Difficulty string` field

### ActiveGame Updates:
- [ ] Add `Mode player.GameMode` field
- [ ] Change `WhitePlayer *websocket.Client` to `WhitePlayer player.Player`
- [ ] Change `BlackPlayer *websocket.Client` to `BlackPlayer player.Player`

**Verification**:
```bash
go build ./internal/websocket/handlers/game_manager.go
# Will have errors in methods - fix in next parts
```

---

## 📝 Part 3: Update CreatePendingGame (20 min)

**File**: `internal/websocket/handlers/game_manager.go`

- [ ] Change method signature:
  ```go
  func (m *GameManager) CreatePendingGame(
      gameID string,
      mode player.GameMode,
      creator *websocket.Client,
      difficulty string,
      timeControl string,
  ) error
  ```
- [ ] Get player types from mode: `whiteType, blackType := mode.GetPlayerTypes()`
- [ ] Create white player based on `whiteType`:
  - [ ] If `PlayerTypeHuman`: use `HumanPlayerConfig{Client: creator}`
  - [ ] If `PlayerTypeComputer`: use `ComputerPlayerConfig{Difficulty: difficulty}`
- [ ] Call `m.playerFactory.CreatePlayer(config)` to create player
- [ ] Update PendingGame creation with Mode, Player, Difficulty
- [ ] Return error instead of void

**Verification**:
```bash
go build ./internal/websocket/handlers/game_manager.go
```

---

## 📝 Part 4: Update ActivateGame (25 min)

**File**: `internal/websocket/handlers/game_manager.go`

- [ ] Get player types: `_, blackType := pending.Mode.GetPlayerTypes()`
- [ ] Create black player based on `blackType`:
  - [ ] If `PlayerTypeHuman`:
    - [ ] Check `blackClient != nil`
    - [ ] Use `HumanPlayerConfig{Client: blackClient}`
  - [ ] If `PlayerTypeComputer`:
    - [ ] Use `ComputerPlayerConfig{Difficulty: pending.Difficulty}`
- [ ] Call `m.playerFactory.CreatePlayer(config)` to create player
- [ ] Add `Mode: pending.Mode` to ActiveGame creation
- [ ] Handle errors properly

**Verification**:
```bash
go build ./internal/websocket/handlers/game_manager.go
```

---

## 📝 Part 5: Update Helper Methods (20 min)

**File**: `internal/websocket/handlers/game_manager.go`

### 5.1 Update GetPlayerSide:
- [ ] Change `game.WhitePlayer.ID` to `game.WhitePlayer.GetID()`
- [ ] Change `game.BlackPlayer.ID` to `game.BlackPlayer.GetID()`

### 5.2 Update GetOpponent:
- [ ] Change return type from `*websocket.Client` to `player.Player`
- [ ] Change `game.WhitePlayer.ID` to `game.WhitePlayer.GetID()`
- [ ] Change `game.BlackPlayer.ID` to `game.BlackPlayer.GetID()`

### 5.3 Add New Helper Methods:

- [ ] Add `GetPlayer(gameID string, side string) (player.Player, error)`
  - [ ] Return white or black player based on side
  - [ ] Return error for invalid side

- [ ] Add `GetGameMode(gameID string) (player.GameMode, error)`
  - [ ] Return game mode
  - [ ] Return error if game not found

- [ ] Add `IsComputerTurn(gameID string, fen string) bool`
  - [ ] Add import: `"strings"`
  - [ ] Parse FEN to get active color (field 2)
  - [ ] Check if active player is computer
  - [ ] Return true if computer's turn

**Verification**:
```bash
go build ./internal/websocket/handlers/game_manager.go
```

---

## 📝 Part 6: Add Computer Move Handler (30 min)

**File**: `internal/websocket/handlers/game_manager.go`

- [ ] Add import: `"context"` (if not already present)
- [ ] Add import: `"log"` (if not already present)
- [ ] Create `HandleComputerMove(ctx context.Context, gameID string) error` method:

  - [ ] Get active game (with RLock)
  - [ ] Check if it's computer's turn using `IsComputerTurn()`
  - [ ] If not computer's turn, return nil
  - [ ] Determine which player is computer (white or black)
  - [ ] Call `computerPlayer.RequestMove(ctx, game.CurrentFEN)`
  - [ ] Validate move using `m.chessService.MakeMove(currentFEN, moveUCI)`
  - [ ] Update game state with new FEN
  - [ ] Check if game is over
  - [ ] Create `MoveNotification` struct
  - [ ] Send move to white player
  - [ ] Send move to black player
  - [ ] If game over, remove active game
  - [ ] Log computer moves for debugging

**Verification**:
```bash
go build ./internal/websocket/handlers/game_manager.go
```

---

## ✅ Final Verification Checklist

### GameManager Structure:
- [ ] `aiService player.AIService` field added
- [ ] `playerFactory *player.PlayerFactory` field added
- [ ] NewGameManager accepts aiService parameter
- [ ] PlayerFactory initialized in constructor

### PendingGame:
- [ ] `Mode player.GameMode` field added
- [ ] `WhitePlayer player.Player` (not *websocket.Client)
- [ ] `Difficulty string` field added

### ActiveGame:
- [ ] `Mode player.GameMode` field added
- [ ] `WhitePlayer player.Player` (not *websocket.Client)
- [ ] `BlackPlayer player.Player` (not *websocket.Client)

### Methods Updated:
- [ ] CreatePendingGame accepts mode & difficulty
- [ ] CreatePendingGame creates players via factory
- [ ] ActivateGame creates black player by mode
- [ ] GetPlayerSide uses `Player.GetID()`
- [ ] GetOpponent returns `player.Player`

### New Methods Added:
- [ ] GetPlayer(gameID, side) method
- [ ] GetGameMode(gameID) method
- [ ] IsComputerTurn(gameID, fen) method
- [ ] HandleComputerMove(ctx, gameID) method

### Compilation:
- [ ] game_manager.go compiles successfully
```bash
go build ./internal/websocket/handlers/game_manager.go
```

### Expected Status:
- [ ] GameManager file compiles ✅
- [ ] Other handler files have errors ⚠️ (expected - will fix in Steps 5 & 6)

---

## 🚨 Expected Compilation Errors

After completing Step 4, these files will have errors (to be fixed in Steps 5 & 6):

1. **create_game.go**:
   - `CreatePendingGame` signature mismatch (needs mode & difficulty)

2. **join_game.go**:
   - May need updates for new ActivateGame behavior

3. **make_move.go**:
   - Needs to call `HandleComputerMove()` after human moves

4. **Server initialization** (main.go or similar):
   - NewGameManager needs AIService parameter

**Don't fix these yet!** We'll address them systematically in Steps 5 & 6.

---

## 📊 Progress Tracking

Track your progress:

```
Part 1: Update GameManager Struct     [ ] → [x]  ✅
Part 2: Update Game Structs           [ ] → [x]  ✅
Part 3: Update CreatePendingGame      [ ] → [x]  ✅
Part 4: Update ActivateGame           [ ] → [x]  ✅
Part 5: Update Helper Methods         [ ] → [x]  ✅
Part 6: Add Computer Move Handler     [ ] → [x]  ✅
```

**When all parts complete**: ✅ Step 4 Done! → Move to Step 5

---

## 🎯 Next Step

Once all checkboxes are complete:

→ **[Step 5: Create Mode Handler](./05-create-mode-handler.md)**

**What you'll do**:
- Update CreateGameHandler to accept game mode
- Update JoinGameHandler to auto-activate computer games
- Fix compilation errors from Step 4 changes

---

**Questions?** Check the full guide: [04-refactor-game-manager.md](./04-refactor-game-manager.md)

**Ready?** Start with Part 1 (Update GameManager Struct)!

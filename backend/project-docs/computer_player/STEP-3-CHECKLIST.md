# Step 3: Implement Player Types - Checklist

**Goal**: Implement HumanPlayer and ComputerPlayer with AIService

**Estimated Time**: 2-3 hours

---

## 📋 Quick Overview

| Part | File | Time | Status |
|------|------|------|--------|
| 1. Config Refactor | `config.go` | 30 min | ⬜ |
| 2. AIService | `ai_service.go` | 45 min | ⬜ |
| 3. HumanPlayer | `human_player.go` | 30 min | ⬜ |
| 4. ComputerPlayer | `computer_player.go` | 45 min | ⬜ |
| 5. Update Factory | `factory.go` | 15 min | ⬜ |
| 6. Fix GameMode | `game_mode.go` | 5 min | ⬜ |

---

## 📝 Part 1: Config Refactor (30 min)

### Create `config.go`

- [ ] Create file `internal/websocket/handlers/player/config.go`
- [ ] Define `PlayerConfig` interface with `Validate()` and `GetPlayerType()`
- [ ] Implement `HumanPlayerConfig` struct
- [ ] Implement `ComputerPlayerConfig` struct
- [ ] Add validation logic for each config
- [ ] (Optional) Add `LLMPlayerConfig` stub for future

**Verification**:
```bash
go build ./internal/websocket/handlers/player/...
```

---

## 📝 Part 2: AIService Implementation (45 min)

### Install Dependency

```bash
cd /Users/ankit/code/learn/chess-coach/backend
go get github.com/notnil/chess/uci
```

### Create `ai_service.go`

- [ ] Create file `internal/websocket/handlers/player/ai_service.go`
- [ ] Define `AIService` interface
- [ ] Implement `StockfishService` struct
- [ ] Add `NewStockfishService()` constructor
  - [ ] Get STOCKFISH_PATH from environment
  - [ ] Fallback to common paths
  - [ ] Initialize UCI engine
- [ ] Implement `GetBestMove(ctx, game, difficulty)` method
  - [ ] Set position
  - [ ] Calculate with time based on difficulty
  - [ ] Return best move
- [ ] Implement `getMoveTime(difficulty)` helper
  - [ ] easy: 100ms
  - [ ] medium: 500ms
  - [ ] hard: 2s
- [ ] Implement `Close()` method

**Verification**:
```bash
# Test Stockfish path
echo $STOCKFISH_PATH

# Build
go build ./internal/websocket/handlers/player/...
```

---

## 📝 Part 3: HumanPlayer Implementation (30 min)

### Create `human_player.go`

- [ ] Create file `internal/websocket/handlers/player/human_player.go`
- [ ] Define `HumanPlayer` struct
- [ ] Implement `NewHumanPlayer(client)` constructor
- [ ] Implement `GetID()` method
- [ ] Implement `GetType()` method → return `PlayerTypeHuman`
- [ ] Implement `RequestMove(ctx, fen)` method
  - [ ] Use channel for async move handling
  - [ ] Handle context timeout
  - [ ] Return move from WebSocket
- [ ] Implement `SendMove(move)` method
- [ ] Implement `NotifyGameStart(info)` method
- [ ] Implement `NotifyGameEnd(result)` method
- [ ] Implement `IsConnected()` method
- [ ] Implement `Cleanup()` method

**Note**: `RequestMove` needs integration with existing WebSocket handlers.

**Verification**:
```bash
go build ./internal/websocket/handlers/player/...
```

---

## 📝 Part 4: ComputerPlayer Implementation (45 min)

### Create `computer_player.go`

- [ ] Create file `internal/websocket/handlers/player/computer_player.go`
- [ ] Define `ComputerPlayer` struct
  - [ ] id string
  - [ ] aiService AIService
  - [ ] difficulty string
  - [ ] currentGame *chess.Game
- [ ] Implement `NewComputerPlayer(aiService, difficulty)` constructor
- [ ] Implement `GetID()` method → return `"computer-{difficulty}"`
- [ ] Implement `GetType()` method → return `PlayerTypeComputer`
- [ ] Implement `RequestMove(ctx, fen)` method
  - [ ] Parse FEN to chess.Game
  - [ ] Check if game is over
  - [ ] Call aiService.GetBestMove()
  - [ ] Return move in UCI format
- [ ] Implement `SendMove(move)` method (just log)
- [ ] Implement `NotifyGameStart(info)` method (just log)
- [ ] Implement `NotifyGameEnd(result)` method (just log)
- [ ] Implement `IsConnected()` method → return `true`
- [ ] Implement `Cleanup()` method

**Verification**:
```bash
go build ./internal/websocket/handlers/player/...
```

---

## 📝 Part 5: Update Factory (15 min)

### Update `factory.go`

- [ ] Change `CreatePlayer` signature to accept `PlayerConfig` interface
- [ ] Remove old `PlayerConfig` struct
- [ ] Update `CreatePlayer` implementation:
  - [ ] Call `config.Validate()`
  - [ ] Switch on `config.GetPlayerType()`
  - [ ] Type assert to specific config type
  - [ ] Create appropriate player
- [ ] Keep `NewPlayerFactory` as-is

**Verification**:
```bash
go build ./internal/websocket/handlers/player/...
```

---

## 📝 Part 6: Fix GameMode Default (5 min)

### Update `game_mode.go`

- [ ] Find line ~40 (default case)
- [ ] Change from: `return PlayerTypeHuman, PlayerTypeHuman`
- [ ] To: `return PlayerTypeHuman, PlayerTypeComputer`

**Verification**:
```bash
go build ./internal/websocket/handlers/player/...
```

---

## ✅ Final Verification Checklist

### Files Created
- [ ] `config.go` exists with all config types
- [ ] `ai_service.go` exists with Stockfish wrapper
- [ ] `human_player.go` exists with full implementation
- [ ] `computer_player.go` exists with full implementation

### Files Updated
- [ ] `factory.go` uses new config interface
- [ ] `game_mode.go` has correct default

### Dependencies
- [ ] `github.com/notnil/chess/uci` installed
- [ ] Appears in `go.mod`

### Compilation
- [ ] All files compile without errors
```bash
go build ./internal/websocket/handlers/player/...
```

### Vet Check
- [ ] No issues from go vet
```bash
go vet ./internal/websocket/handlers/player/...
```

---

## 🧪 Manual Testing

### Test 1: Create AIService

```bash
# In Go playground or test file
aiService, err := player.NewStockfishService()
if err != nil {
    log.Fatal(err)
}
defer aiService.Close()

game := chess.NewGame()
move, err := aiService.GetBestMove(context.Background(), game, "medium")
if err != nil {
    log.Fatal(err)
}

log.Printf("Stockfish suggests: %s", move)
```

Expected: Should print a valid opening move like `e2e4` or `d2d4`

### Test 2: Create Computer Player

```go
factory := player.NewPlayerFactory(aiService)

config := player.ComputerPlayerConfig{Difficulty: "hard"}
computerPlayer, err := factory.CreatePlayer(config)
if err != nil {
    log.Fatal(err)
}

log.Printf("Player ID: %s", computerPlayer.GetID())
log.Printf("Player Type: %s", computerPlayer.GetType())
```

Expected:
```
Player ID: computer-hard
Player Type: computer
```

### Test 3: Request Move from Computer

```go
fen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
move, err := computerPlayer.RequestMove(context.Background(), fen)
if err != nil {
    log.Fatal(err)
}

log.Printf("Computer plays: %s", move)
```

Expected: Should return a valid UCI move like `e2e4`

---

## 🚨 Common Issues

### Issue: "stockfish not found"
**Solution**:
```bash
# Check STOCKFISH_PATH
echo $STOCKFISH_PATH

# Should be: /usr/local/bin/stockfish (in Docker)
# Or: /opt/homebrew/bin/stockfish (on Mac)

# Test manually
stockfish
# Type 'quit' to exit
```

### Issue: "package chess/uci not found"
**Solution**:
```bash
go get github.com/notnil/chess/uci
go mod tidy
```

### Issue: Build errors in factory.go
**Solution**: Make sure you removed the old `PlayerConfig` struct and updated the method signature.

### Issue: Type assertion errors
**Solution**: Check that you're using the correct config type for each player.

---

## 📊 Progress Tracking

| Component | Lines | Time | Status |
|-----------|-------|------|--------|
| config.go | ~80 | 30 min | ⬜ |
| ai_service.go | ~120 | 45 min | ⬜ |
| human_player.go | ~120 | 30 min | ⬜ |
| computer_player.go | ~100 | 45 min | ⬜ |
| factory.go (update) | ~50 | 15 min | ⬜ |
| game_mode.go (fix) | ~2 | 5 min | ⬜ |
| **Total** | **~470** | **2h 50min** | ⬜ |

---

## 🎯 Next Step

Once all checkboxes are complete:

→ **[Step 4: Refactor GameManager](./04-refactor-game-manager.md)**

**What you'll do**:
- Integrate Player system into game flow
- Update game creation to support modes
- Add computer move handling

---

**Questions?** Check the full guide: [03-implement-player-types.md](./03-implement-player-types.md)

**Ready?** Start with Part 1 (Config Refactor)!

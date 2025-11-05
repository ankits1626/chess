# Step 5: Create Mode Handler

**Goal**: Update CreateGameHandler and JoinGameHandler to support game modes

**Estimated Time**: 1.5 hours

---

## 📋 Overview

After Step 4, we have:
- ✅ GameManager with Player system
- ✅ Support for different game modes
- ⚠️ Handlers that still assume Human vs Human

**What We'll Do**:
1. Update CreateGameHandler to accept `mode` and `difficulty` parameters
2. Update JoinGameHandler to handle computer games (auto-activate)
3. Update server initialization to pass AIService to GameManager
4. Fix compilation errors from Step 4

---

## 🎯 Tasks

### Part 1: Update CreateGameHandler (30 min)

**File**: `internal/websocket/handlers/create_game.go`

**Changes**:
1. Extract `mode` from request data (default: "humanVsComputer")
2. Extract `difficulty` from request data (default: "medium")
3. Validate game mode
4. Update `CreatePendingGame` call to include mode and difficulty
5. If mode is HumanVsComputer, auto-activate immediately (no waiting for second player)

**Request format**:
```json
{
  "id": "req-123",
  "type": "request",
  "action": "createGame",
  "data": {
    "mode": "humanVsComputer",
    "difficulty": "medium",
    "timeControl": "5+0"
  }
}
```

---

### Part 2: Update JoinGameHandler (20 min)

**File**: `internal/websocket/handlers/join_game.go`

**Changes**:
1. Check if pending game mode requires second human player
2. If mode is ComputerVsComputer or similar, reject join attempt
3. Update ActivateGame call (should work with Step 4 changes)

---

### Part 3: Initialize AIService (30 min)

**File**: Server initialization (main.go or cmd/server/main.go)

**Changes**:
1. Create AIService instance:
   ```go
   aiService, err := player.NewStockfishService()
   if err != nil {
       log.Fatal(err)
   }
   defer aiService.Close()
   ```
2. Pass aiService to NewGameManager:
   ```go
   gameManager := handlers.NewGameManager(db, chessService, aiService)
   ```

---

### Part 4: Auto-Activate Computer Games (10 min)

**File**: `internal/websocket/handlers/create_game.go`

**Logic**:
After creating pending game, if mode doesn't need second human:
```go
if mode == player.GameModeHumanVsComputer {
    // Auto-activate with nil blackClient (computer will be created)
    activeGame, err := h.manager.ActivateGame(gameIDStr, nil)
    if err != nil {
        return err
    }

    // Trigger computer move if computer is white
    go h.manager.HandleComputerMove(ctx, gameIDStr)
}
```

---

## ✅ Verification

After completion:
- [ ] Can create game with mode="humanVsComputer"
- [ ] Computer games auto-activate
- [ ] Computer makes first move if playing white
- [ ] Human vs Human still works (mode="humanVsHuman")
- [ ] All handlers compile

---

**Detailed implementation guide coming when you reach this step!**

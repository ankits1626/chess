# Step 6: Update MakeMove Handler

**Goal**: Update MakeMove handler to trigger computer moves after human moves

**Estimated Time**: 1 hour

---

## 📋 Overview

After Step 5, we have:
- ✅ Games can be created with different modes
- ✅ Computer games auto-activate
- ⚠️ Computer doesn't respond to human moves

**What We'll Do**:
1. Update MakeMoveHandler to call HandleComputerMove after human moves
2. Add proper error handling for computer move failures
3. Test the complete flow

---

## 🎯 Tasks

### Part 1: Update MakeMoveHandler (30 min)

**File**: `internal/websocket/handlers/make_move.go`

**Changes**:
1. After successful human move, check game mode
2. If game includes computer player, trigger computer move:
   ```go
   // After updating game state and notifying players
   mode, err := h.manager.GetGameMode(gameID)
   if err == nil && mode != player.GameModeHumanVsHuman {
       // Trigger computer move in background
       go func() {
           ctx := context.Background()
           if err := h.manager.HandleComputerMove(ctx, gameID); err != nil {
               log.Printf("Computer move failed: %v", err)
           }
       }()
   }
   ```

---

### Part 2: Test Full Flow (30 min)

**Using wscat**:

1. **Create human vs computer game**:
   ```json
   {
     "id": "1",
     "type": "request",
     "action": "createGame",
     "data": {
       "mode": "humanVsComputer",
       "difficulty": "easy",
       "timeControl": "5+0"
     }
   }
   ```

2. **Make move as human**:
   ```json
   {
     "id": "2",
     "type": "request",
     "action": "makeMove",
     "data": {
       "gameId": "...",
       "move": "e2e4"
     }
   }
   ```

3. **Expect**:
   - Your move succeeds
   - Computer responds with its move
   - You receive computer's move notification

---

## ✅ Verification

After completion:
- [ ] Human can make moves
- [ ] Computer responds to human moves
- [ ] Computer moves are valid
- [ ] Game continues until checkmate/draw
- [ ] Different difficulty levels work
- [ ] Can play full game against computer

---

**Detailed implementation guide coming when you reach this step!**

# Chess Engine Integration - Quick Summary

**Enable human vs. computer gameplay with Stockfish**

---

## 🎯 Goal

Add chess engine to your backend so users can play against the computer.

---

## 📋 Implementation Phases

### **Phase 1: Install Stockfish** (30 min)
```bash
# macOS
brew install stockfish

# Verify
stockfish
# Type: uci (should show: uciok)
# Type: quit
```

### **Phase 2: Create UCI Client** (2 hours)

**Add dependency:**
```bash
go get github.com/notnil/chess
```

**Create files:**
1. `internal/engine/uci.go` - UCI protocol client
2. `internal/engine/config.go` - Engine configuration
3. `internal/service/engine_service.go` - Business logic

### **Phase 3: Add API Endpoints** (1 hour)

**New endpoints:**
- `POST /api/v1/engine/suggest` - Get computer move
- `POST /api/v1/engine/evaluate` - Evaluate position

**Create:**
- `internal/handler/v1/engine/handler.go`

**Update:**
- `internal/router/v1_routes.go`

### **Phase 4: Database Updates** (1 hour)

**Add columns:**
```sql
ALTER TABLE games ADD COLUMN game_type VARCHAR(20);
ALTER TABLE games ADD COLUMN difficulty INTEGER;
ALTER TABLE moves ADD COLUMN evaluation INTEGER;
```

**Create computer user:**
```sql
INSERT INTO users (username, email, rating)
VALUES ('Computer', 'computer@chesscoach.ai', 3200);
```

### **Phase 5: Frontend Integration** (3 hours)

**Create:**
- `frontend/src/services/engineService.ts` - API client
- `frontend/src/components/ComputerGameControls.tsx` - UI

**Update:**
- Game store to handle computer moves
- Board component to trigger computer after user move

---

## 🧪 Quick Test

```bash
# Start server
./bin/server

# Test move suggestion
curl -X POST http://localhost:8080/api/v1/engine/suggest \
  -H "Content-Type: application/json" \
  -d '{
    "fen": "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
    "difficulty": 10
  }' | jq .

# Expected: {"move":"e2e4","move_san":"e4","evaluation":25}
```

---

## 🎮 How It Works

```
User makes move → Frontend
    ↓
POST /api/v1/games/:id/moves (save move)
    ↓
POST /api/v1/engine/suggest (get computer move)
    ↓
Backend → Stockfish → Best move
    ↓
POST /api/v1/games/:id/moves (save computer move)
    ↓
Frontend updates board
```

---

## 🔑 Key Files

**Backend:**
- `internal/engine/uci.go` - Stockfish communication
- `internal/service/engine_service.go` - Move suggestion logic
- `internal/handler/v1/engine/handler.go` - API endpoints

**Frontend:**
- `src/services/engineService.ts` - API calls
- `src/components/ComputerGameControls.tsx` - UI controls
- `src/store/gameStore.ts` - State management

---

## 💡 Difficulty Levels

| Level | Depth | Strength | Description |
|-------|-------|----------|-------------|
| 1-5   | 1-5   | Beginner | Makes mistakes |
| 6-10  | 6-10  | Intermediate | Decent play |
| 11-15 | 11-15 | Advanced | Strong player |
| 16-20 | 16-20 | Master | Near perfect |

---

## 🚀 After Implementation

You'll have:
- ✅ Computer opponent with adjustable difficulty
- ✅ Position evaluation (centipawns)
- ✅ Move hints/suggestions
- ✅ Full game recording
- ✅ Foundation for advanced features

---

## 📚 Full Guide

See [step-22-chess-engine-integration.md](./step-22-chess-engine-integration.md) for:
- Complete code listings
- Detailed explanations
- Error handling
- Testing procedures
- Advanced features

---

**Total Time: 4-6 hours**

**Next Steps:**
1. Opening book integration
2. Game analysis (blunder detection)
3. Multi-line analysis
4. Tournament mode

---

**Last Updated**: 2025-11-01

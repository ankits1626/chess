# Phase 4: Architecture Comparison

**Which implementation should you choose?**

---

## 📋 Two Approaches Available

### 1. **Single-File Approach** ([04-message-handlers.md](./04-message-handlers.md))
All handlers in `internal/websocket/handler.go`

### 2. **SOLID Architecture** ([04-message-handlers-REFACTORED.md](./04-message-handlers-REFACTORED.md)) ⭐ **Recommended**
Handlers split across multiple files following SOLID principles

---

## 🔍 Detailed Comparison

| Aspect | Single-File | SOLID Architecture |
|--------|-------------|-------------------|
| **Total Files** | 1 file | 7 files |
| **Total Lines** | ~500 lines | ~780 lines |
| **Learning Curve** | Easier (all in one place) | Moderate (need to navigate) |
| **Maintainability** | Lower (large file) | Higher (focused files) |
| **Testability** | Harder to mock | Easy to mock interfaces |
| **Team Collaboration** | Merge conflicts likely | Each dev owns a handler |
| **Adding Features** | Modify large file | Add new file |
| **SOLID Principles** | ❌ Violates SRP | ✅ Follows all SOLID |
| **Production Ready** | Acceptable | Excellent |

---

## 📁 File Structure Comparison

### Single-File Approach

```
internal/websocket/
├── handler.go              # ~500 lines - EVERYTHING
├── message.go
├── client.go
├── hub.go
├── room.go
├── utils.go
└── upgrade.go
```

**Pros**:
- Simple to navigate (one file)
- Good for learning
- Quick to implement

**Cons**:
- Large file gets unwieldy
- Hard to test individual handlers
- Merge conflicts in teams
- Violates Single Responsibility

---

### SOLID Architecture ⭐

```
internal/websocket/
├── handler.go              # ~100 lines - Router & interfaces
├── game_manager.go         # ~150 lines - State management
├── chess_service.go        # ~100 lines - Chess abstraction
├── handler_create_game.go  # ~80 lines  - CreateGame only
├── handler_join_game.go    # ~100 lines - JoinGame only
├── handler_make_move.go    # ~150 lines - MakeMove only
├── handler_resign.go       # ~100 lines - Resign/Leave only
├── message.go
├── client.go
├── hub.go
├── room.go
├── utils.go
└── upgrade.go
```

**Pros**:
- Each file has ONE responsibility
- Easy to find and modify specific handler
- Easy to test (interfaces for mocking)
- Team-friendly (no merge conflicts)
- Professional architecture
- Follows all SOLID principles

**Cons**:
- More files to navigate
- Slightly more setup code

---

## 🎯 SOLID Principles Explained

### Single Responsibility Principle (SRP)
**Problem**: Single file has multiple responsibilities
**Solution**: Each handler in its own file

```go
// ❌ Single-File: One file does everything
handler.go
  - Routes messages
  - Creates games
  - Joins games
  - Makes moves
  - Handles resignations

// ✅ SOLID: Each file has one job
handler.go               → Routes messages only
handler_create_game.go   → Creates games only
handler_join_game.go     → Joins games only
handler_make_move.go     → Makes moves only
handler_resign.go        → Handles resignations only
```

### Open/Closed Principle (OCP)
**Problem**: Adding new handler requires editing large file
**Solution**: Just add a new handler file

```go
// ❌ Single-File: Must edit handler.go switch statement
func (h *GameHandler) HandleMessage(...) {
    switch msg.Action {
    case "createGame": ...
    case "joinGame": ...
    case "newFeature": ...  // EDIT LARGE FILE
    }
}

// ✅ SOLID: Just create new file, register handler
handler_new_feature.go  // NEW FILE - no edits to existing code
router.RegisterHandler("newFeature", &NewFeatureHandler{})
```

### Liskov Substitution Principle (LSP)
**Problem**: Hard to swap implementations
**Solution**: Depend on interfaces

```go
// ❌ Single-File: Directly uses chess library
game := chess.NewGame()  // Tightly coupled

// ✅ SOLID: Use interface
type GameService interface {
    ValidateMove(...) (MoveResult, error)
    ApplyMove(...) (GameState, error)
}

// Can swap with MockChessService for testing
```

### Interface Segregation Principle (ISP)
**Problem**: Large interfaces force unnecessary dependencies
**Solution**: Small, focused interfaces

```go
// ❌ Single-File: Handlers access everything
type GameHandler struct {
    db            *database.DB
    pendingGames  map[string]*PendingGame
    activeGames   map[string]*ActiveGame
    // All handlers share this state
}

// ✅ SOLID: Handlers only see what they need
type ActionHandler interface {
    Handle(ctx, client, msg) error  // Just one method
}
```

### Dependency Inversion Principle (DIP)
**Problem**: Handlers depend on concrete implementations
**Solution**: Handlers depend on abstractions

```go
// ❌ Single-File: Depend on concrete chess library
func handleMove(...) {
    game := chess.NewGame()  // Concrete dependency
}

// ✅ SOLID: Depend on interface
type GameService interface {
    ApplyMove(...) (GameState, error)
}

func (h *MakeMoveHandler) Handle(...) {
    state, err := h.manager.GetChessService().ApplyMove(...)
}
```

---

## 🧪 Testing Comparison

### Single-File Testing

```go
// Hard to test individual handlers
func TestHandleMessage(t *testing.T) {
    handler := &GameHandler{
        db: realDB,  // Need real DB
        pendingGames: ...,
        activeGames: ...,
    }
    // Tests entire handler at once
}
```

### SOLID Testing

```go
// Easy to test individual handlers
func TestCreateGameHandler(t *testing.T) {
    mockDB := &MockDB{}
    mockChess := &MockChessService{
        ApplyMoveFunc: func(...) (GameState, error) {
            return GameState{...}, nil  // Controlled response
        },
    }

    manager := NewGameManager(mockDB, mockChess)
    handler := &CreateGameHandler{manager: manager}

    // Test just CreateGame in isolation
    err := handler.Handle(ctx, client, msg)
    assert.NoError(t, err)
}
```

---

## 💼 Real-World Scenarios

### Scenario 1: Bug in Move Validation

**Single-File**:
1. Open 500-line handler.go
2. Find handleMakeMove function (scroll, search)
3. Fix bug
4. Risk: Might break other handlers

**SOLID**:
1. Open handler_make_move.go (150 lines, focused)
2. Bug is obvious (only move logic here)
3. Fix bug
4. No risk to other handlers

---

### Scenario 2: Add "Draw Offer" Feature

**Single-File**:
1. Edit handler.go (already 500 lines)
2. Add case to switch statement
3. Add handleDrawOffer function
4. File now 600 lines
5. Merge conflict if teammate editing same file

**SOLID**:
1. Create handler_draw_offer.go (new file)
2. Implement DrawOfferHandler
3. Register in one line
4. No merge conflicts
5. Easy to review (small PR)

---

### Scenario 3: Team of 3 Developers

**Single-File**:
- Dev A: Working on createGame in handler.go
- Dev B: Working on makeMove in handler.go
- Dev C: Working on resign in handler.go
- **Result**: 🔴 Constant merge conflicts

**SOLID**:
- Dev A: Working on handler_create_game.go
- Dev B: Working on handler_make_move.go
- Dev C: Working on handler_resign.go
- **Result**: ✅ No conflicts, parallel development

---

## 📊 Code Metrics

| Metric | Single-File | SOLID |
|--------|-------------|-------|
| Lines per file (avg) | 500 | 100 |
| Cyclomatic complexity | High | Low |
| Test coverage | Hard to achieve | Easy to achieve |
| Bug isolation | Difficult | Easy |
| Onboarding time | Medium | Fast (clear structure) |
| Refactoring risk | High | Low |

---

## 🎓 Learning Recommendation

### If You're New to Go/WebSockets
→ Start with **Single-File** to understand flow
→ Then refactor to **SOLID** for production

### If You're Experienced
→ Go directly with **SOLID Architecture** ⭐

---

## 🚀 Our Recommendation

**Use SOLID Architecture** ([04-message-handlers-REFACTORED.md](./04-message-handlers-REFACTORED.md))

**Reasons**:
1. ✅ Chess library already installed (line 47 in go.mod)
2. ✅ You're asking about SOLID - shows maturity
3. ✅ Better for your portfolio
4. ✅ Industry standard approach
5. ✅ Easier to explain in interviews

The extra 30 minutes to set up SOLID will save you hours later.

---

## 📝 Migration Path

If you already implemented Single-File:

1. **Keep it working** (don't break existing code)
2. **Create new files** alongside handler.go
3. **Move handlers one-by-one**:
   - Create handler_create_game.go
   - Copy handleCreateGame code
   - Update to use GameManager
   - Test
   - Repeat for other handlers
4. **Remove old handler.go** when done

---

## ✅ Final Decision Matrix

Choose **Single-File** if:
- [ ] You're learning WebSockets for the first time
- [ ] You want to prototype quickly
- [ ] You're working solo
- [ ] You won't add more handlers

Choose **SOLID** if:
- [x] You want production-quality code
- [x] You're building a portfolio project
- [x] You might work in a team
- [x] You plan to add more features
- [x] You care about testability
- [x] You want to follow best practices

---

**Decision**: We recommend **SOLID Architecture** ⭐

**Next**: Implement [04-message-handlers-REFACTORED.md](./04-message-handlers-REFACTORED.md)

---

**Last Updated**: 2025-11-05
**Recommendation**: SOLID Architecture for production readiness

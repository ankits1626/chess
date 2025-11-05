# Game Modes - Master Implementation Plan

**Goal**: Build extensible chess game system supporting multiple opponent types

**Design Principles**: SOLID, Plugin Architecture, Extensible

**Total Duration**: 10-14 hours

**Status**: 📋 Planning Complete

---

## 🎯 Architecture Vision

```
┌─────────────────────────────────────────────────────────────┐
│                    Game Mode System                         │
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐       │
│  │   Human      │  │   Computer   │  │     LLM      │       │
│  │   Player     │  │   Player     │  │   Player     │       │
│  │  (WebSocket) │  │ (Stockfish)  │  │  (GPT-4)     │       │
│  └──────────────┘  └──────────────┘  └──────────────┘       │
│         │                 │                 │               │
│         └─────────────────┴─────────────────┘               │
│                          │                                  │
│                   ┌──────▼──────┐                           │
│                   │   Player    │                           │
│                   │  Interface  │  (Dependency Inversion)   │
│                   └──────┬──────┘                           │
│                          │                                  │
│                   ┌──────▼──────┐                           │
│                   │    Game     │                           │
│                   │   Manager   │                           │
│                   └──────┬──────┘                           │
│                          │                                  │
│         ┌────────────────┼────────────────┐                 │
│         │                │                │                 │
│    ┌────▼────┐      ┌────▼────┐      ┌────▼────┐            │
│    │  Human  │      │Computer │      │  Human  │            │
│    │   vs    │      │   vs    │      │   vs    │            │
│    │ Human   │      │Computer │      │   LLM   │            │
│    └─────────┘      └─────────┘      └─────────┘            │
└─────────────────────────────────────────────────────────────┘
```

---

## 🏗️ Game Modes to Support

### Phase 1: Foundation (Steps 1-6)
1. **Human vs Computer** ⭐ (First Milestone)
   - Human plays white/black
   - Stockfish AI as opponent
   - WebSocket for human, AI for computer

### Phase 2: Future Modes (Extensible)
2. **Human vs Human**
   - Already implemented (Phase 4)
   - Two WebSocket clients

3. **Computer vs Computer**
   - Two AI engines battle
   - For testing/demo purposes

4. **Human vs LLM**
   - Human via WebSocket
   - LLM (GPT-4, Claude) via API
   - Interesting for chess coaching

---

## 📚 Implementation Steps

### **Milestone 1: Backend System (Steps 1-6)**

#### **Step 1: Design Player Interface** (1 hour)
→ [01-design-player-interface.md](./01-design-player-interface.md)
- Create `Player` interface (SOLID - Interface Segregation)
- Define `PlayerType` enum (Human, Computer, LLM)
- Create `PlayerFactory` for instantiation
- Design game mode configuration

#### **Step 2: Install & Test Stockfish** (30 minutes)
→ [02-install-stockfish.md](./02-install-stockfish.md)
- Install Stockfish engine
- Verify installation
- Test UCI protocol
- Install Go UCI wrapper

#### **Step 3: Implement Player Types** (2 hours)
→ [03-implement-player-types.md](./03-implement-player-types.md)
- Create `HumanPlayer` (wraps WebSocket client)
- Create `ComputerPlayer` (wraps Stockfish)
- Implement `Player` interface for both
- Add difficulty levels for ComputerPlayer

#### **Step 4: Refactor GameManager** (2 hours)
→ [04-refactor-game-manager.md](./04-refactor-game-manager.md)
- Add `GameMode` enum
- Update game structures to use `Player` interface
- Add mode-specific game maps
- Implement game creation with modes

#### **Step 5: Create Game Mode Handler** (2 hours)
→ [05-create-game-mode-handler.md](./05-create-game-mode-handler.md)
- Create `CreateGameWithModeHandler`
- Support Human vs Computer mode
- Handle player instantiation via factory
- Manage game flow based on mode

#### **Step 6: Update MakeMove for Modes** (1.5 hours)
→ [06-update-make-move-modes.md](./06-update-make-move-modes.md)
- Detect game mode
- Route move to appropriate handler
- Trigger computer/LLM response if needed
- Handle async move generation

---

### **Milestone 2: Backend Testing (Steps 7-8)**

#### **Step 7: Wire Backend Components** (1 hour)
→ [07-wire-backend.md](./07-wire-backend.md)
- Initialize PlayerFactory
- Wire AI services
- Update app.go
- Test compilation

#### **Step 8: Test with wscat** (1.5 hours)
→ [08-test-with-wscat.md](./08-test-with-wscat.md)
- Install wscat
- Create game via WebSocket
- Make moves
- Verify computer responses
- Test edge cases

---

### **Milestone 3: Frontend Integration (Steps 9-12)**

#### **Step 9: Frontend - WebSocket Service** (1 hour)
→ [09-frontend-websocket.md](./09-frontend-websocket.md)
- Create WebSocket service
- Handle game mode messages
- Handle computer move events

#### **Step 10: Frontend - Game Mode Selector** (1.5 hours)
→ [10-frontend-game-mode-selector.md](./10-frontend-game-mode-selector.md)
- Create mode selector UI
- Support Human vs Computer
- Difficulty and side selection
- Extensible for future modes

#### **Step 11: Frontend - Update Game Store** (1 hour)
→ [11-frontend-game-store.md](./11-frontend-game-store.md)
- Add game mode state
- Handle mode-specific logic
- Update move handling

#### **Step 12: Frontend - Update App** (1 hour)
→ [12-frontend-update-app.md](./12-frontend-update-app.md)
- Integrate mode selector
- Show game mode in UI
- Handle cleanup

---

### **Milestone 4: Polish (Step 13)**

#### **Step 13: Testing & Polish** (2 hours)
→ [13-testing-polish.md](./13-testing-polish.md)
- End-to-end testing
- Error handling
- UI polish
- Documentation

---

## 🎨 SOLID Principles Applied

### Single Responsibility Principle (SRP)
```
✅ Player interface - ONE responsibility: represent a player
✅ HumanPlayer - handle human interactions
✅ ComputerPlayer - handle AI interactions
✅ GameManager - manage game lifecycle
✅ Each handler - handle ONE action
```

### Open/Closed Principle (OCP)
```
✅ New player types (LLM) can be added without modifying existing code
✅ Just implement Player interface and register in factory
```

### Liskov Substitution Principle (LSP)
```
✅ Any Player implementation can replace another
✅ GameManager works with Player interface, not concrete types
```

### Interface Segregation Principle (ISP)
```
✅ Player interface has only essential methods
✅ No client forced to depend on methods it doesn't use
```

### Dependency Inversion Principle (DIP)
```
✅ GameManager depends on Player interface (abstraction)
✅ Not on HumanPlayer or ComputerPlayer (concrete)
```

---

## 📋 Player Interface Design Preview

```go
// Player represents any chess player (human, computer, LLM)
type Player interface {
    // Core identification
    GetID() string
    GetType() PlayerType

    // Move handling
    RequestMove(ctx context.Context, fen string) (string, error)
    SendMove(move MoveNotification) error

    // Game events
    NotifyGameStart(info GameStartInfo) error
    NotifyGameEnd(result GameResult) error

    // Lifecycle
    IsConnected() bool
    Cleanup() error
}

// PlayerType enum
type PlayerType string

const (
    PlayerTypeHuman    PlayerType = "human"
    PlayerTypeComputer PlayerType = "computer"
    PlayerTypeLLM      PlayerType = "llm"
)

// GameMode defines who plays against whom
type GameMode string

const (
    GameModeHumanVsHuman    GameMode = "human_vs_human"
    GameModeHumanVsComputer GameMode = "human_vs_computer"
    GameModeComputerVsComputer GameMode = "computer_vs_computer"
    GameModeHumanVsLLM      GameMode = "human_vs_llm"
)
```

---

## 📊 Progress Tracking

| Step | Milestone | Status | Duration | Notes |
|------|-----------|--------|----------|-------|
| 1. Design Player Interface | Backend | ⬜ Not Started | 1 hour | Core architecture |
| 2. Install Stockfish | Backend | ⬜ Not Started | 30 min | One-time setup |
| 3. Implement Player Types | Backend | ⬜ Not Started | 2 hours | HumanPlayer, ComputerPlayer |
| 4. Refactor GameManager | Backend | ⬜ Not Started | 2 hours | Use Player interface |
| 5. Create Mode Handler | Backend | ⬜ Not Started | 2 hours | Handle game creation |
| 6. Update MakeMove | Backend | ⬜ Not Started | 1.5 hours | Mode-aware moves |
| 7. Wire Backend | Backend | ⬜ Not Started | 1 hour | Integration |
| 8. Test with wscat | Testing | ⬜ Not Started | 1.5 hours | Validate backend |
| 9. WebSocket Service | Frontend | ⬜ Not Started | 1 hour | |
| 10. Mode Selector | Frontend | ⬜ Not Started | 1.5 hours | |
| 11. Update Store | Frontend | ⬜ Not Started | 1 hour | |
| 12. Update App | Frontend | ⬜ Not Started | 1 hour | |
| 13. Testing & Polish | Polish | ⬜ Not Started | 2 hours | |

**Legend**: ⬜ Not Started | 🟡 In Progress | ✅ Complete | ❌ Blocked

---

## 🎯 Milestones

### ✅ Milestone 1: Backend Complete (Steps 1-6)
- Player interface designed
- Computer player working
- Game modes supported
- Compiles and runs

### ✅ Milestone 2: Backend Validated (Steps 7-8)
- Tested with wscat
- All game flows work
- Error handling verified

### ✅ Milestone 3: Frontend Complete (Steps 9-12)
- UI for mode selection
- Human vs Computer playable
- Smooth user experience

### ✅ Milestone 4: Production Ready (Step 13)
- Fully tested
- Polished UI
- Documented

---

## 🚀 Future Extensibility

### Adding New Player Type (e.g., LLM Player)

**Step 1**: Implement Player interface
```go
type LLMPlayer struct {
    id       string
    apiKey   string
    model    string
}

func (p *LLMPlayer) RequestMove(ctx context.Context, fen string) (string, error) {
    // Call OpenAI/Claude API
    // Get move from LLM
    return move, nil
}

// Implement other Player methods...
```

**Step 2**: Register in PlayerFactory
```go
func (f *PlayerFactory) CreatePlayer(playerType PlayerType, config PlayerConfig) (Player, error) {
    switch playerType {
    case PlayerTypeHuman:
        return NewHumanPlayer(config.Client), nil
    case PlayerTypeComputer:
        return NewComputerPlayer(config.Difficulty), nil
    case PlayerTypeLLM:  // NEW
        return NewLLMPlayer(config.APIKey, config.Model), nil
    }
}
```

**Step 3**: Add game mode
```go
const GameModeHumanVsLLM GameMode = "human_vs_llm"
```

**Done!** No changes to GameManager or existing handlers needed. ✅

---

## 📁 File Structure

```
internal/websocket/handlers/
├── handler.go              # Router & interfaces
├── game_manager.go         # Game lifecycle (updated)
├── chess_service.go        # Chess logic
├── create_game.go          # Existing multiplayer
├── join_game.go            # Existing multiplayer
├── make_move.go            # Updated for modes
├── resign.go               # Existing
│
├── player/                 # NEW - Player abstraction
│   ├── player.go           # Player interface
│   ├── types.go            # PlayerType, GameMode enums
│   ├── factory.go          # PlayerFactory
│   ├── human_player.go     # HumanPlayer implementation
│   ├── computer_player.go  # ComputerPlayer implementation
│   └── llm_player.go       # (Future) LLM implementation
│
├── ai/                     # NEW - AI engines
│   ├── ai_service.go       # AIService interface
│   ├── stockfish.go        # Stockfish implementation
│   └── difficulty.go       # Difficulty levels
│
└── game_mode_handler.go    # NEW - Mode-based game creation
```

---

## 🎯 Next Step

Start with **[Step 1: Design Player Interface](./01-design-player-interface.md)** →

This is the foundation - get this right and everything else becomes plug-and-play!

---

## 📖 Key Design Decisions

1. **Player Interface**: All player types implement same interface
2. **Factory Pattern**: PlayerFactory creates appropriate player types
3. **Game Mode Enum**: Explicit modes for clarity and type safety
4. **Separation of Concerns**: AI logic separate from player abstraction
5. **Backward Compatibility**: Existing Human vs Human mode unchanged

---

**Last Updated**: 2025-11-05
**Architecture**: SOLID, Plugin-based, Extensible

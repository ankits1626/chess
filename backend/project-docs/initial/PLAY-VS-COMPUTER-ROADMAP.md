# Play vs Computer - Implementation Roadmap

**Goal**: Allow users to play chess against a computer AI without authentication

**Duration**: 8-12 hours

**Status**: 📋 Planning phase

---

## 🎯 Overview

Users open the React app → Click "Play vs Computer" → Immediately start playing against AI

**No authentication required** - Sessions identified by temporary client IDs

---

## 📊 Current State Analysis

### Frontend (React + Vite + TypeScript + Zustand)
- ✅ Chess.js library installed (v1.4.0)
- ✅ Zustand state management
- ✅ GameBoard component exists
- ✅ Live play mode already implemented
- ✅ Move validation working
- ✅ Promotion dialog working

### Backend (Go + WebSocket + PostgreSQL)
- ✅ WebSocket infrastructure (Phases 1-3 complete)
- ✅ Phase 4 handlers implemented (CreateGame, JoinGame, MakeMove, etc.)
- ✅ Chess validation service using notnil/chess
- ✅ Database schema with games and moves tables
- ❌ **Missing**: Computer AI opponent
- ❌ **Missing**: "Play vs Computer" game type

---

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                        User's Browser                       │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  React Frontend (No Auth)                             │  │
│  │  - Click "Play vs Computer"                           │  │
│  │  - Display chess board                                │  │
│  │  - Make moves                                         │  │
│  │  - Receive AI moves                                   │  │
│  └───────────────────────────────────────────────────────┘  │
│                            │                                 │
│                     WebSocket Connection                     │
│                            │                                 │
└────────────────────────────┼─────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                      Backend Server                         │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  WebSocket Handler                                   │   │
│  │  - createComputerGame action                         │   │
│  │  - makeMove action (user's move)                     │   │
│  │  - computerMove event (AI's response)                │   │
│  └──────────────────────────────────────────────────────┘   │
│                            │                                │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Chess AI Engine                                     │   │
│  │  - Analyze position                                  │   │
│  │  - Generate best move                                │   │
│  │  - Difficulty levels (easy, medium, hard)            │   │
│  └──────────────────────────────────────────────────────┘   │
│                            │                                │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Database (PostgreSQL)                               │   │
│  │  - Store computer games                              │   │
│  │  - black_player_id = NULL or special "computer" UUID │   │
│  └──────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

---

## 🔧 Implementation Plan

### **Phase 1: Backend - Chess AI Engine (3-4 hours)**

#### Step 1.1: Choose AI Approach (30 minutes)

**Option A: Simple Random Moves** (Quick prototype)
- ✅ Fast to implement (30 min)
- ✅ No external dependencies
- ❌ Not challenging for players

**Option B: Stockfish Integration** ⭐ **RECOMMENDED**
- ✅ World-class chess AI
- ✅ Adjustable difficulty
- ✅ Already available as Go package
- ⏱️ Takes 2-3 hours to integrate

**Option C: Minimax Algorithm** (Custom implementation)
- ✅ Educational value
- ✅ Full control
- ❌ Complex to implement well (10+ hours)
- ❌ Slower than Stockfish

**Decision**: Start with **Option B (Stockfish)** for production quality

#### Step 1.2: Install Stockfish (15 minutes)

```bash
# Install Stockfish on macOS
brew install stockfish

# Install Go wrapper for Stockfish
cd backend
go get github.com/notnil/chess/uci
```

#### Step 1.3: Create AI Service (2 hours)

**File**: `internal/websocket/handlers/ai_service.go`

```go
package handlers

import (
    "context"
    "fmt"
    "time"

    "github.com/notnil/chess"
    "github.com/notnil/chess/uci"
)

// AIService provides computer chess opponent functionality.
type AIService interface {
    GetBestMove(ctx context.Context, fen string, difficulty Difficulty) (string, error)
    Close() error
}

// Difficulty levels for AI
type Difficulty int

const (
    DifficultyEasy   Difficulty = 1  // ~800-1200 ELO
    DifficultyMedium Difficulty = 2  // ~1200-1600 ELO
    DifficultyHard   Difficulty = 3  // ~1600-2000 ELO
)

// StockfishAI implements AIService using Stockfish engine.
type StockfishAI struct {
    engine *uci.Engine
}

// NewStockfishAI creates a new Stockfish AI service.
func NewStockfishAI(enginePath string) (*StockfishAI, error) {
    // enginePath: "/opt/homebrew/bin/stockfish" on macOS
    eng, err := uci.New(enginePath)
    if err != nil {
        return nil, fmt.Errorf("failed to start engine: %w", err)
    }

    // Initialize engine
    if err := eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
        return nil, fmt.Errorf("failed to initialize engine: %w", err)
    }

    return &StockfishAI{engine: eng}, nil
}

// GetBestMove returns the best move for the given position.
func (ai *StockfishAI) GetBestMove(ctx context.Context, fen string, difficulty Difficulty) (string, error) {
    // Set position
    fenFunc, err := chess.FEN(fen)
    if err != nil {
        return "", fmt.Errorf("invalid FEN: %w", err)
    }

    game := chess.NewGame(fenFunc)

    // Configure difficulty (limit strength)
    var searchTime time.Duration
    var depth int

    switch difficulty {
    case DifficultyEasy:
        searchTime = 100 * time.Millisecond
        depth = 5
    case DifficultyMedium:
        searchTime = 500 * time.Millisecond
        depth = 10
    case DifficultyHard:
        searchTime = 2 * time.Second
        depth = 15
    default:
        searchTime = 500 * time.Millisecond
        depth = 10
    }

    // Set position in engine
    if err := ai.engine.Run(uci.CmdPosition{Position: game.Position()}, uci.CmdIsReady); err != nil {
        return "", fmt.Errorf("failed to set position: %w", err)
    }

    // Search for best move
    cmdGo := uci.CmdGo{
        Depth:    depth,
        MoveTime: searchTime,
    }

    if err := ai.engine.Run(cmdGo); err != nil {
        return "", fmt.Errorf("search failed: %w", err)
    }

    // Get result
    searchResults := ai.engine.SearchResults()
    if searchResults == nil || searchResults.BestMove == nil {
        return "", fmt.Errorf("no move found")
    }

    // Convert to UCI notation
    move := searchResults.BestMove
    uci := chess.UCINotation{}.Encode(game.Position(), move)

    return uci, nil
}

// Close shuts down the engine.
func (ai *StockfishAI) Close() error {
    if ai.engine != nil {
        return ai.engine.Close()
    }
    return nil
}
```

#### Step 1.4: Create Computer Game Handler (1 hour)

**File**: `internal/websocket/handlers/create_computer_game.go`

```go
package handlers

import (
    "context"
    "fmt"
    "log"

    "github.com/ankits1626/chess-coach-backend/internal/database"
    "github.com/ankits1626/chess-coach-backend/internal/websocket"
    "github.com/jackc/pgx/v5/pgtype"
)

// CreateComputerGameHandler handles computer game creation.
type CreateComputerGameHandler struct {
    manager *GameManager
}

// Handle implements ActionHandler.Handle.
// Expected message:
// {
//   "id": "req-123",
//   "type": "request",
//   "action": "createComputerGame",
//   "data": {
//     "difficulty": "medium",  // "easy", "medium", "hard"
//     "playerSide": "white"    // "white" or "black"
//   }
// }
func (h *CreateComputerGameHandler) Handle(ctx context.Context, client *websocket.Client, msg *websocket.Message) error {
    // 1. Extract and validate inputs
    difficulty, err := websocket.RequireString(msg.Data, "difficulty")
    if err != nil {
        return err
    }

    playerSide, err := websocket.RequireString(msg.Data, "playerSide")
    if err != nil {
        return err
    }

    if difficulty != "easy" && difficulty != "medium" && difficulty != "hard" {
        return fmt.Errorf("invalid difficulty: must be easy, medium, or hard")
    }

    if playerSide != "white" && playerSide != "black" {
        return fmt.Errorf("invalid playerSide: must be white or black")
    }

    log.Printf("CreateComputerGame: User %s starting game as %s with difficulty %s",
        client.UserID, playerSide, difficulty)

    // 2. Validate player ID
    playerUUID := websocket.StringToUUID(client.UserID)
    if !playerUUID.Valid {
        return fmt.Errorf("invalid user ID")
    }

    // 3. Create game in database
    var whitePlayerID, blackPlayerID pgtype.UUID
    if playerSide == "white" {
        whitePlayerID = playerUUID
        blackPlayerID = pgtype.UUID{Valid: false} // Computer is black
    } else {
        whitePlayerID = pgtype.UUID{Valid: false} // Computer is white
        blackPlayerID = playerUUID
    }

    game, err := h.manager.GetDB().CreateGame(ctx, database.CreateGameParams{
        WhitePlayerID: whitePlayerID,
        BlackPlayerID: blackPlayerID,
        Pgn:           "",
        Result:        websocket.StringToText("*"),
        TimeControl:   websocket.IntToInt4(0), // Untimed
    })
    if err != nil {
        log.Printf("CreateComputerGame: Database error: %v", err)
        return fmt.Errorf("failed to create game")
    }

    gameIDStr := websocket.UUIDToString(game.ID)

    // 4. Create active computer game
    h.manager.CreateComputerGame(gameIDStr, client, playerSide, difficulty)

    // 5. Set client's game ID
    client.SetGameID(gameIDStr)

    log.Printf("CreateComputerGame: Game %s created, player is %s", gameIDStr, playerSide)

    // 6. Get starting position
    startingFEN := h.manager.GetChessService().GetStartingPosition()

    // 7. If computer plays white, make first move
    var firstMove *ComputerMoveResult
    if playerSide == "black" {
        firstMove, err = h.manager.MakeComputerMove(ctx, gameIDStr)
        if err != nil {
            log.Printf("CreateComputerGame: Failed to make computer move: %v", err)
            return fmt.Errorf("failed to make computer move")
        }
    }

    // 8. Send success response
    responseData := map[string]interface{}{
        "gameId":     gameIDStr,
        "status":     "active",
        "side":       playerSide,
        "difficulty": difficulty,
        "fen":        startingFEN,
        "vsComputer": true,
    }

    if firstMove != nil {
        responseData["fen"] = firstMove.NewFEN
        responseData["computerMove"] = map[string]interface{}{
            "san": firstMove.SAN,
            "uci": firstMove.UCI,
            "fen": firstMove.NewFEN,
        }
    }

    websocket.SendSuccessResponse(client, msg.ID, responseData)

    return nil
}

// ComputerMoveResult holds the result of a computer move.
type ComputerMoveResult struct {
    SAN        string
    UCI        string
    NewFEN     string
    IsGameOver bool
    Result     GameResult
}
```

#### Step 1.5: Update GameManager (1 hour)

**File**: `internal/websocket/handlers/game_manager.go`

Add computer game support:

```go
// Add to GameManager struct
type GameManager struct {
    db            *database.DB
    chessService  GameService
    aiService     AIService  // NEW
    pendingGames  map[string]*PendingGame
    activeGames   map[string]*ActiveGame
    computerGames map[string]*ComputerGame  // NEW
    mu            sync.RWMutex
}

// ComputerGame represents a game vs computer.
type ComputerGame struct {
    GameID      string
    Player      *websocket.Client
    PlayerSide  string // "white" or "black"
    CurrentFEN  string
    MoveCount   int
    Difficulty  string // "easy", "medium", "hard"
}

// NewGameManager - update constructor
func NewGameManager(db *database.DB, chessService GameService, aiService AIService) *GameManager {
    return &GameManager{
        db:            db,
        chessService:  chessService,
        aiService:     aiService,
        pendingGames:  make(map[string]*PendingGame),
        activeGames:   make(map[string]*ActiveGame),
        computerGames: make(map[string]*ComputerGame),
    }
}

// CreateComputerGame creates a new computer game.
func (m *GameManager) CreateComputerGame(gameID string, player *websocket.Client, playerSide, difficulty string) {
    m.mu.Lock()
    defer m.mu.Unlock()

    m.computerGames[gameID] = &ComputerGame{
        GameID:     gameID,
        Player:     player,
        PlayerSide: playerSide,
        CurrentFEN: m.chessService.GetStartingPosition(),
        MoveCount:  0,
        Difficulty: difficulty,
    }
}

// GetComputerGame retrieves a computer game.
func (m *GameManager) GetComputerGame(gameID string) (*ComputerGame, bool) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    game, exists := m.computerGames[gameID]
    return game, exists
}

// MakeComputerMove generates and applies a computer move.
func (m *GameManager) MakeComputerMove(ctx context.Context, gameID string) (*ComputerMoveResult, error) {
    m.mu.Lock()
    computerGame, exists := m.computerGames[gameID]
    m.mu.Unlock()

    if !exists {
        return nil, fmt.Errorf("computer game not found")
    }

    // Convert difficulty string to enum
    var difficulty Difficulty
    switch computerGame.Difficulty {
    case "easy":
        difficulty = DifficultyEasy
    case "medium":
        difficulty = DifficultyMedium
    case "hard":
        difficulty = DifficultyHard
    default:
        difficulty = DifficultyMedium
    }

    // Get best move from AI
    moveUCI, err := m.aiService.GetBestMove(ctx, computerGame.CurrentFEN, difficulty)
    if err != nil {
        return nil, fmt.Errorf("AI failed to find move: %w", err)
    }

    // Apply move
    gameState, err := m.chessService.ApplyMove(computerGame.CurrentFEN, moveUCI)
    if err != nil {
        return nil, fmt.Errorf("failed to apply move: %w", err)
    }

    // Update game state
    m.mu.Lock()
    computerGame.CurrentFEN = gameState.FEN
    computerGame.MoveCount++
    m.mu.Unlock()

    return &ComputerMoveResult{
        SAN:        gameState.SAN,
        UCI:        gameState.UCI,
        NewFEN:     gameState.FEN,
        IsGameOver: gameState.IsGameOver,
        Result:     gameState.Result,
    }, nil
}

// RemoveComputerGame removes a computer game.
func (m *GameManager) RemoveComputerGame(gameID string) {
    m.mu.Lock()
    defer m.mu.Unlock()

    delete(m.computerGames, gameID)
}
```

#### Step 1.6: Update MakeMove Handler (30 minutes)

**File**: `internal/websocket/handlers/make_move.go`

Add computer response after player move:

```go
func (h *MakeMoveHandler) Handle(ctx context.Context, client *websocket.Client, msg *websocket.Message) error {
    // ... existing validation code ...

    // Check if this is a computer game
    computerGame, isComputerGame := h.manager.GetComputerGame(gameIDStr)

    if isComputerGame {
        return h.handleComputerGameMove(ctx, client, msg, computerGame, moveSAN)
    }

    // ... existing multiplayer game code ...
}

func (h *MakeMoveHandler) handleComputerGameMove(
    ctx context.Context,
    client *websocket.Client,
    msg *websocket.Message,
    computerGame *ComputerGame,
    moveSAN string,
) error {
    gameIDStr := computerGame.GameID

    // 1. Apply player's move
    gameState, err := h.manager.GetChessService().ApplyMove(computerGame.CurrentFEN, moveSAN)
    if err != nil {
        return fmt.Errorf("invalid move: %w", err)
    }

    // 2. Update game state
    if err := h.manager.UpdateComputerGameState(gameIDStr, gameState.FEN); err != nil {
        return fmt.Errorf("failed to update game state")
    }

    // 3. Save player move to database
    // ... save move ...

    // 4. Check if game over after player move
    if gameState.IsGameOver {
        return h.handleGameOver(ctx, client, msg, gameIDStr, gameState)
    }

    // 5. Send response to player
    websocket.SendSuccessResponse(client, msg.ID, map[string]interface{}{
        "gameId":     gameIDStr,
        "move":       gameState.SAN,
        "fen":        gameState.FEN,
        "gameOver":   false,
        "vsComputer": true,
    })

    // 6. Make computer move
    computerMoveResult, err := h.manager.MakeComputerMove(ctx, gameIDStr)
    if err != nil {
        log.Printf("Computer move failed: %v", err)
        return nil // Don't fail player's move
    }

    // 7. Save computer move to database
    // ... save move ...

    // 8. Send computer move to player
    computerMoveEvent := websocket.NewEvent("computerMove", map[string]interface{}{
        "gameId":   gameIDStr,
        "move":     computerMoveResult.SAN,
        "fen":      computerMoveResult.NewFEN,
        "gameOver": computerMoveResult.IsGameOver,
    })

    if computerMoveResult.IsGameOver {
        computerMoveEvent.Data["result"] = h.getResult(computerMoveResult.Result)
        computerMoveEvent.Data["method"] = computerMoveResult.Result.Method
    }

    client.SendMessage(computerMoveEvent)

    // 9. Clean up if game over
    if computerMoveResult.IsGameOver {
        h.manager.RemoveComputerGame(gameIDStr)
    }

    return nil
}
```

#### Step 1.7: Wire Everything (30 minutes)

**File**: `internal/app/app.go`

```go
func New(cfg *config.Config, db *database.DB, log logger.Logger) *App {
    // Create chess service
    chessService := handlers.NewChessService()

    // Create AI service
    aiService, err := handlers.NewStockfishAI("/opt/homebrew/bin/stockfish")
    if err != nil {
        log.Fatal("Failed to initialize AI:", err)
    }

    // Create game manager (now with AI)
    gameManager := handlers.NewGameManager(db, chessService, aiService)

    // Create handler router
    router := handlers.NewHandlerRouter(gameManager)

    // Register computer game handler
    router.RegisterHandler("createComputerGame", &handlers.CreateComputerGameHandler{manager: gameManager})

    // Create hub
    hub := websocket.NewHub(router)

    return &App{
        config: cfg,
        db:     db,
        server: server.New(cfg, db, hub),
        logger: log,
        wsHub:  hub,
        aiService: aiService, // Store for cleanup
    }
}

// Cleanup on shutdown
func (a *App) Shutdown() {
    if a.aiService != nil {
        a.aiService.Close()
    }
}
```

---

### **Phase 2: Frontend - UI Updates (2-3 hours)**

#### Step 2.1: Add Game Mode Selection (1 hour)

**File**: `frontend/app/src/components/game/GameModeSelector.tsx`

```tsx
import { useState } from 'react';

interface GameModeSelectorProps {
  onStartComputerGame: (difficulty: string, side: string) => void;
  onStartMultiplayer: () => void;
}

export default function GameModeSelector({ onStartComputerGame, onStartMultiplayer }: GameModeSelectorProps) {
  const [showComputerOptions, setShowComputerOptions] = useState(false);
  const [difficulty, setDifficulty] = useState('medium');
  const [side, setSide] = useState('white');

  if (showComputerOptions) {
    return (
      <div className="bg-gray-700 rounded-lg p-6 space-y-4">
        <h2 className="text-xl font-bold">Play vs Computer</h2>

        <div>
          <label className="block text-sm font-medium mb-2">Difficulty</label>
          <div className="grid grid-cols-3 gap-2">
            {['easy', 'medium', 'hard'].map((d) => (
              <button
                key={d}
                onClick={() => setDifficulty(d)}
                className={`px-4 py-2 rounded ${
                  difficulty === d
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-600 hover:bg-gray-500'
                }`}
              >
                {d.charAt(0).toUpperCase() + d.slice(1)}
              </button>
            ))}
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium mb-2">Play as</label>
          <div className="grid grid-cols-2 gap-2">
            {['white', 'black'].map((s) => (
              <button
                key={s}
                onClick={() => setSide(s)}
                className={`px-4 py-2 rounded ${
                  side === s
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-600 hover:bg-gray-500'
                }`}
              >
                {s.charAt(0).toUpperCase() + s.slice(1)}
              </button>
            ))}
          </div>
        </div>

        <div className="flex gap-2">
          <button
            onClick={() => onStartComputerGame(difficulty, side)}
            className="flex-1 bg-green-600 hover:bg-green-700 text-white px-6 py-3 rounded-lg font-medium"
          >
            Start Game
          </button>
          <button
            onClick={() => setShowComputerOptions(false)}
            className="px-4 py-3 bg-gray-600 hover:bg-gray-500 rounded-lg"
          >
            Back
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-gray-700 rounded-lg p-6 space-y-3">
      <h2 className="text-xl font-bold mb-4">Start a Game</h2>

      <button
        onClick={() => setShowComputerOptions(true)}
        className="w-full bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-medium"
      >
        🤖 Play vs Computer
      </button>

      <button
        onClick={onStartMultiplayer}
        className="w-full bg-purple-600 hover:bg-purple-700 text-white px-6 py-3 rounded-lg font-medium"
      >
        👥 Play vs Human
      </button>
    </div>
  );
}
```

#### Step 2.2: Create WebSocket Service (1 hour)

**File**: `frontend/app/src/services/websocketService.ts`

```typescript
export interface WebSocketMessage {
  id: string;
  type: 'request' | 'response' | 'event';
  action?: string;
  data?: any;
  error?: string;
}

class WebSocketService {
  private ws: WebSocket | null = null;
  private messageHandlers: Map<string, (msg: WebSocketMessage) => void> = new Map();
  private eventHandlers: Map<string, (data: any) => void> = new Map();

  connect(url: string): Promise<void> {
    return new Promise((resolve, reject) => {
      this.ws = new WebSocket(url);

      this.ws.onopen = () => {
        console.log('WebSocket connected');
        resolve();
      };

      this.ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        reject(error);
      };

      this.ws.onmessage = (event) => {
        const msg: WebSocketMessage = JSON.parse(event.data);

        if (msg.type === 'response') {
          const handler = this.messageHandlers.get(msg.id);
          if (handler) {
            handler(msg);
            this.messageHandlers.delete(msg.id);
          }
        } else if (msg.type === 'event' && msg.action) {
          const handler = this.eventHandlers.get(msg.action);
          if (handler) {
            handler(msg.data);
          }
        }
      };

      this.ws.onclose = () => {
        console.log('WebSocket disconnected');
      };
    });
  }

  send(action: string, data: any): Promise<WebSocketMessage> {
    return new Promise((resolve, reject) => {
      if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
        reject(new Error('WebSocket not connected'));
        return;
      }

      const id = `req-${Date.now()}-${Math.random()}`;
      const message: WebSocketMessage = {
        id,
        type: 'request',
        action,
        data,
      };

      this.messageHandlers.set(id, (response) => {
        if (response.error) {
          reject(new Error(response.error));
        } else {
          resolve(response);
        }
      });

      this.ws.send(JSON.stringify(message));

      // Timeout after 10 seconds
      setTimeout(() => {
        if (this.messageHandlers.has(id)) {
          this.messageHandlers.delete(id);
          reject(new Error('Request timeout'));
        }
      }, 10000);
    });
  }

  on(eventName: string, handler: (data: any) => void) {
    this.eventHandlers.set(eventName, handler);
  }

  off(eventName: string) {
    this.eventHandlers.delete(eventName);
  }

  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
    this.messageHandlers.clear();
    this.eventHandlers.clear();
  }
}

export const wsService = new WebSocketService();
```

#### Step 2.3: Update Game Store (1 hour)

**File**: `frontend/app/src/store/useGameStore.ts`

Add computer game state and actions:

```typescript
interface GameState {
  // ... existing state ...

  // Computer game state
  isComputerGame: boolean;
  computerDifficulty: 'easy' | 'medium' | 'hard' | null;
  playerSide: 'white' | 'black' | null;
  waitingForComputer: boolean;

  // Actions
  startComputerGame: (difficulty: string, side: string) => Promise<void>;
  handleComputerMove: (san: string, fen: string) => void;
}

export const useGameStore = create<GameState>((set, get) => ({
  // ... existing state ...

  isComputerGame: false,
  computerDifficulty: null,
  playerSide: null,
  waitingForComputer: false,

  startComputerGame: async (difficulty: string, side: string) => {
    try {
      // Connect to WebSocket
      await wsService.connect('ws://localhost:8080/ws');

      // Listen for computer moves
      wsService.on('computerMove', (data) => {
        get().handleComputerMove(data.move, data.fen);

        if (data.gameOver) {
          console.log('Game over:', data.result);
        }
      });

      // Create computer game
      const response = await wsService.send('createComputerGame', {
        difficulty,
        playerSide: side,
      });

      const newGame = new Chess();

      // If computer made first move, apply it
      if (response.data.computerMove) {
        newGame.move(response.data.computerMove.san);
      }

      set({
        game: newGame,
        mode: 'live',
        isComputerGame: true,
        computerDifficulty: difficulty as any,
        playerSide: side as any,
        waitingForComputer: false,
        selectedSquare: null,
        validMoves: [],
        lastMove: null,
      });
    } catch (error) {
      console.error('Failed to start computer game:', error);
      throw error;
    }
  },

  handleComputerMove: (san: string, fen: string) => {
    const { game } = get();

    try {
      game.move(san);
      set({
        game: Object.assign(Object.create(Object.getPrototypeOf(game)), game),
        waitingForComputer: false,
        lastMove: null, // Could extract from/to from SAN if needed
      });
    } catch (error) {
      console.error('Failed to apply computer move:', error);
    }
  },

  // Update selectSquare to handle computer games
  selectSquare: (square: Square) => {
    const { game, isComputerGame, playerSide, waitingForComputer, mode } = get();

    if (mode === 'replay') return;
    if (waitingForComputer) return;

    // For computer games, only allow moves on player's turn
    if (isComputerGame) {
      const isWhiteTurn = game.turn() === 'w';
      const isPlayerTurn = (playerSide === 'white' && isWhiteTurn) ||
                          (playerSide === 'black' && !isWhiteTurn);

      if (!isPlayerTurn) return;
    }

    // ... existing move logic ...

    // After successful move in computer game
    if (isComputerGame && move) {
      set({ waitingForComputer: true });

      // Send move to backend
      wsService.send('makeMove', {
        gameId: 'current-game-id', // Store game ID in state
        move: move.san,
      }).catch(error => {
        console.error('Failed to send move:', error);
        set({ waitingForComputer: false });
      });
    }
  },
}));
```

#### Step 2.4: Update App Component (30 minutes)

**File**: `frontend/app/src/App.tsx`

```tsx
import GameModeSelector from '@/components/game/GameModeSelector';

function App() {
  const [showModeSelector, setShowModeSelector] = useState(true);
  const isComputerGame = useGameStore(state => state.isComputerGame);
  const startComputerGame = useGameStore(state => state.startComputerGame);

  const handleStartComputerGame = async (difficulty: string, side: string) => {
    try {
      await startComputerGame(difficulty, side);
      setShowModeSelector(false);
    } catch (error) {
      alert('Failed to start game: ' + error);
    }
  };

  return (
    <div className="min-h-screen bg-gray-800 text-white p-4">
      {showModeSelector ? (
        <div className="max-w-md mx-auto mt-20">
          <GameModeSelector
            onStartComputerGame={handleStartComputerGame}
            onStartMultiplayer={() => {
              // TODO: Implement multiplayer
              alert('Multiplayer coming soon!');
            }}
          />
        </div>
      ) : (
        <div className="max-w-[1600px] mx-auto">
          {/* Existing game board UI */}
          <GameBoard />
          <GameInfo />

          {isComputerGame && (
            <button
              onClick={() => {
                setShowModeSelector(true);
                // Cleanup WebSocket
                wsService.disconnect();
              }}
              className="mt-4 bg-red-600 hover:bg-red-700 px-4 py-2 rounded"
            >
              End Game
            </button>
          )}
        </div>
      )}
    </div>
  );
}
```

---

### **Phase 3: Testing & Polish (2-3 hours)**

#### Step 3.1: Manual Testing (1 hour)

Test scenarios:
- [ ] Start computer game as white, easy difficulty
- [ ] Start computer game as black, easy difficulty
- [ ] Start computer game as white, medium difficulty
- [ ] Start computer game as white, hard difficulty
- [ ] Make valid moves
- [ ] Try illegal moves
- [ ] Complete a full game (checkmate)
- [ ] Complete a full game (stalemate)
- [ ] Resign from computer game
- [ ] Disconnect and reconnect

#### Step 3.2: Error Handling (1 hour)

Add proper error handling:
- WebSocket connection failures
- Stockfish engine crashes
- Invalid moves
- Game not found errors
- Timeout handling

#### Step 3.3: UI Polish (1 hour)

- Add loading spinner while waiting for computer move
- Add "Computer is thinking..." indicator
- Add move highlight for computer moves
- Add sound effects (optional)
- Add game over modal with result

---

## 📝 File Checklist

### Backend Files to Create/Modify:
- [ ] `internal/websocket/handlers/ai_service.go` (NEW)
- [ ] `internal/websocket/handlers/create_computer_game.go` (NEW)
- [ ] `internal/websocket/handlers/game_manager.go` (MODIFY)
- [ ] `internal/websocket/handlers/make_move.go` (MODIFY)
- [ ] `internal/websocket/handlers/handler.go` (MODIFY - register handler)
- [ ] `internal/app/app.go` (MODIFY - wire AI service)

### Frontend Files to Create/Modify:
- [ ] `frontend/app/src/components/game/GameModeSelector.tsx` (NEW)
- [ ] `frontend/app/src/services/websocketService.ts` (NEW)
- [ ] `frontend/app/src/store/useGameStore.ts` (MODIFY)
- [ ] `frontend/app/src/App.tsx` (MODIFY)
- [ ] `frontend/app/src/components/game/GameInfo.tsx` (MODIFY - show difficulty)

---

## 🎯 Success Criteria

- ✅ User can click "Play vs Computer" without login
- ✅ User can choose difficulty (easy, medium, hard)
- ✅ User can choose side (white or black)
- ✅ Game starts immediately
- ✅ Computer responds to moves within 1-3 seconds
- ✅ Computer plays legal moves
- ✅ Game ends properly (checkmate, stalemate, resignation)
- ✅ No authentication required

---

## 🚀 Quick Start Guide (After Implementation)

### Backend:
```bash
# Install Stockfish
brew install stockfish

# Install dependencies
cd backend
go get github.com/notnil/chess/uci

# Run server
go run cmd/server/main.go
```

### Frontend:
```bash
cd frontend/app
npm install
npm run dev
```

### Play:
1. Open http://localhost:5173
2. Click "Play vs Computer"
3. Choose difficulty and side
4. Start playing!

---

## 🔄 Future Enhancements (Optional)

**Phase 4: Advanced Features** (4-6 hours)
- [ ] Move hints (show best move)
- [ ] Position analysis
- [ ] Opening book integration
- [ ] Endgame tablebase
- [ ] Save/load computer games
- [ ] Game review with computer analysis
- [ ] Multiple AI personalities

**Phase 5: Multiplayer** (8-10 hours)
- [ ] Create game and share link
- [ ] Join game by ID
- [ ] Matchmaking system
- [ ] Rating system

---

## 📊 Estimated Timeline

| Phase | Task | Time |
|-------|------|------|
| 1 | Choose AI approach | 30 min |
| 1 | Install Stockfish | 15 min |
| 1 | Create AI Service | 2 hours |
| 1 | Create Computer Game Handler | 1 hour |
| 1 | Update GameManager | 1 hour |
| 1 | Update MakeMove Handler | 30 min |
| 1 | Wire everything | 30 min |
| 2 | Add Game Mode Selector UI | 1 hour |
| 2 | Create WebSocket Service | 1 hour |
| 2 | Update Game Store | 1 hour |
| 2 | Update App Component | 30 min |
| 3 | Manual Testing | 1 hour |
| 3 | Error Handling | 1 hour |
| 3 | UI Polish | 1 hour |
| **Total** | | **12 hours** |

---

## 🎓 Learning Resources

- Stockfish UCI Protocol: https://www.shredderchess.com/download/div/uci.zip
- notnil/chess documentation: https://pkg.go.dev/github.com/notnil/chess
- WebSocket best practices: https://developer.mozilla.org/en-US/docs/Web/API/WebSockets_API

---

**Ready to implement?** Start with Phase 1, Step 1.1! 🚀

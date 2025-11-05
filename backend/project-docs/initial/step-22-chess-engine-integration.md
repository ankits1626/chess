# Chess Engine Integration Guide

**Complete guide for integrating Stockfish chess engine into the Chess Coach backend**

---

## Table of Contents
1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Phase 1: Install Stockfish](#phase-1-install-stockfish)
4. [Phase 2: UCI Protocol Client](#phase-2-uci-protocol-client)
5. [Phase 3: Engine Service Layer](#phase-3-engine-service-layer)
6. [Phase 4: API Endpoints](#phase-4-api-endpoints)
7. [Phase 5: Testing](#phase-5-testing)
8. [Database Schema Updates](#database-schema-updates)
9. [Frontend Integration](#frontend-integration)

---

## Overview

### What We're Building

Add Stockfish chess engine integration to enable:
- 🤖 Computer move generation
- 📊 Position evaluation (centipawns)
- 💡 Move hints/suggestions
- 🎯 Adjustable difficulty (depth 1-20)
- 🧠 Multi-line analysis

### Architecture

```
Frontend (React)
    ↓
Backend API (Gin)
    ↓
Engine Service
    ↓
UCI Protocol Client
    ↓
Stockfish Binary
```

---

## Prerequisites

### What You Already Have

✅ Complete REST API (users, games, moves)
✅ Move storage (SAN, UCI, FEN)
✅ PostgreSQL database
✅ Go backend with Gin
✅ React frontend with chess.js

### Time Estimate

**Total: 4-6 hours**
- Phase 1: 30 minutes (install Stockfish)
- Phase 2: 2 hours (UCI client)
- Phase 3: 1 hour (service layer)
- Phase 4: 1 hour (API endpoints)
- Phase 5: 30-60 minutes (testing)

---

## Phase 1: Install Stockfish

### Step 1.1: Install Stockfish Binary

**Option A: macOS (Homebrew)**
```bash
brew install stockfish
which stockfish
# Output: /opt/homebrew/bin/stockfish
```

**Option B: Linux (apt)**
```bash
sudo apt-get update
sudo apt-get install stockfish
which stockfish
# Output: /usr/games/stockfish
```

**Option C: Manual Download**
```bash
# Download from https://stockfishchess.org/download/
cd ~/Downloads
wget https://stockfishchess.org/files/stockfish-16-linux-x64.zip
unzip stockfish-16-linux-x64.zip
sudo mv stockfish-16-linux-x64/stockfish /usr/local/bin/
chmod +x /usr/local/bin/stockfish
```

**Verify installation:**
```bash
stockfish
# Should enter Stockfish command prompt
# Type: uci
# Should show: uciok
# Type: quit
```

---

### Step 1.2: Test UCI Protocol

```bash
echo -e "uci\nquit" | stockfish
```

**Expected output:**
```
Stockfish 16 by the Stockfish developers
id name Stockfish 16
id author the Stockfish developers
option name Threads type spin default 1 min 1 max 1024
option name Hash type spin default 16 min 1 max 33554432
...
uciok
```

---

## Phase 2: UCI Protocol Client

### Step 2.1: Add Dependencies

```bash
cd backend
go get github.com/notnil/chess
```

**Update `go.mod`:**
```go
module github.com/ankits1626/chess-coach-backend

go 1.23

require (
    github.com/notnil/chess v1.9.0  // ADD THIS
    // ... existing dependencies
)
```

---

### Step 2.2: Create Engine Package

**File:** `internal/engine/uci.go`

```go
// Package engine provides chess engine integration via UCI protocol.
package engine

import (
    "bufio"
    "fmt"
    "io"
    "os/exec"
    "strings"
    "time"
)

// UCIEngine represents a UCI-compliant chess engine.
type UCIEngine struct {
    cmd    *exec.Cmd
    stdin  io.WriteCloser
    stdout *bufio.Scanner
    ready  bool
}

// NewUCIEngine creates and initializes a UCI engine.
func NewUCIEngine(enginePath string) (*UCIEngine, error) {
    cmd := exec.Command(enginePath)

    stdin, err := cmd.StdinPipe()
    if err != nil {
        return nil, fmt.Errorf("failed to get stdin pipe: %w", err)
    }

    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return nil, fmt.Errorf("failed to get stdout pipe: %w", err)
    }

    if err := cmd.Start(); err != nil {
        return nil, fmt.Errorf("failed to start engine: %w", err)
    }

    engine := &UCIEngine{
        cmd:    cmd,
        stdin:  stdin,
        stdout: bufio.NewScanner(stdout),
    }

    // Initialize UCI protocol
    if err := engine.init(); err != nil {
        engine.Close()
        return nil, err
    }

    return engine, nil
}

// init initializes the UCI protocol.
func (e *UCIEngine) init() error {
    // Send UCI command
    if err := e.send("uci"); err != nil {
        return err
    }

    // Wait for uciok
    for e.stdout.Scan() {
        line := e.stdout.Text()
        if line == "uciok" {
            break
        }
    }

    // Send isready
    if err := e.send("isready"); err != nil {
        return err
    }

    // Wait for readyok
    for e.stdout.Scan() {
        line := e.stdout.Text()
        if line == "readyok" {
            e.ready = true
            break
        }
    }

    return nil
}

// send sends a command to the engine.
func (e *UCIEngine) send(command string) error {
    _, err := fmt.Fprintln(e.stdin, command)
    return err
}

// SetPosition sets the board position using FEN.
func (e *UCIEngine) SetPosition(fen string) error {
    return e.send(fmt.Sprintf("position fen %s", fen))
}

// SetPositionWithMoves sets position with move history.
func (e *UCIEngine) SetPositionWithMoves(moves []string) error {
    moveList := strings.Join(moves, " ")
    return e.send(fmt.Sprintf("position startpos moves %s", moveList))
}

// BestMove gets the best move for current position.
// depth: search depth (1-20, higher = stronger but slower)
// timeout: maximum time to think (milliseconds)
func (e *UCIEngine) BestMove(depth int, timeout time.Duration) (string, error) {
    // Send go command with depth
    if err := e.send(fmt.Sprintf("go depth %d", depth)); err != nil {
        return "", err
    }

    // Set timeout
    deadline := time.Now().Add(timeout)

    // Wait for bestmove
    var bestMove string
    for e.stdout.Scan() {
        if time.Now().After(deadline) {
            return "", fmt.Errorf("engine timeout")
        }

        line := e.stdout.Text()
        if strings.HasPrefix(line, "bestmove") {
            parts := strings.Fields(line)
            if len(parts) >= 2 {
                bestMove = parts[1]
                break
            }
        }
    }

    if bestMove == "" {
        return "", fmt.Errorf("no best move found")
    }

    return bestMove, nil
}

// Evaluate gets position evaluation in centipawns.
func (e *UCIEngine) Evaluate(depth int) (int, error) {
    if err := e.send(fmt.Sprintf("go depth %d", depth)); err != nil {
        return 0, err
    }

    var evaluation int
    for e.stdout.Scan() {
        line := e.stdout.Text()

        // Parse info lines for score
        if strings.HasPrefix(line, "info") && strings.Contains(line, "score cp") {
            parts := strings.Fields(line)
            for i, part := range parts {
                if part == "cp" && i+1 < len(parts) {
                    fmt.Sscanf(parts[i+1], "%d", &evaluation)
                }
            }
        }

        // Stop at bestmove
        if strings.HasPrefix(line, "bestmove") {
            break
        }
    }

    return evaluation, nil
}

// Close stops the engine.
func (e *UCIEngine) Close() error {
    if err := e.send("quit"); err != nil {
        return err
    }

    e.stdin.Close()

    // Give engine time to quit gracefully
    done := make(chan error, 1)
    go func() {
        done <- e.cmd.Wait()
    }()

    select {
    case <-time.After(5 * time.Second):
        // Force kill if not exited
        return e.cmd.Process.Kill()
    case err := <-done:
        return err
    }
}
```

---

### Step 2.3: Create Engine Configuration

**File:** `internal/engine/config.go`

```go
package engine

import (
    "os/exec"
    "runtime"
)

// Config holds engine configuration.
type Config struct {
    EnginePath string
    DefaultDepth int
    TimeoutMS int
}

// DefaultConfig returns default engine configuration.
func DefaultConfig() *Config {
    return &Config{
        EnginePath:   findStockfish(),
        DefaultDepth: 10,
        TimeoutMS:    5000, // 5 seconds
    }
}

// findStockfish attempts to find Stockfish binary.
func findStockfish() string {
    // Try common paths
    paths := []string{
        "/opt/homebrew/bin/stockfish",  // macOS Homebrew (M1/M2)
        "/usr/local/bin/stockfish",     // macOS Homebrew (Intel) / Linux manual
        "/usr/bin/stockfish",            // Linux apt
        "/usr/games/stockfish",          // Ubuntu/Debian
        "stockfish",                     // PATH
    }

    for _, path := range paths {
        if _, err := exec.LookPath(path); err == nil {
            return path
        }
    }

    // Windows
    if runtime.GOOS == "windows" {
        return "stockfish.exe"
    }

    // Default fallback
    return "stockfish"
}
```

---

## Phase 3: Engine Service Layer

### Step 3.1: Create Engine Service

**File:** `internal/service/engine_service.go`

```go
// Package service provides business logic layer.
package service

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/ankits1626/chess-coach-backend/internal/engine"
    "github.com/notnil/chess"
)

// EngineService manages chess engine operations.
type EngineService struct {
    engine *engine.UCIEngine
    config *engine.Config
    mu     sync.Mutex
}

// NewEngineService creates a new engine service.
func NewEngineService(cfg *engine.Config) (*EngineService, error) {
    if cfg == nil {
        cfg = engine.DefaultConfig()
    }

    eng, err := engine.NewUCIEngine(cfg.EnginePath)
    if err != nil {
        return nil, fmt.Errorf("failed to initialize engine: %w", err)
    }

    return &EngineService{
        engine: eng,
        config: cfg,
    }, nil
}

// SuggestMoveRequest contains parameters for move suggestion.
type SuggestMoveRequest struct {
    FEN        string   // Current position
    Moves      []string // Move history in UCI format (optional)
    Difficulty int      // 1-20 (depth)
}

// SuggestMoveResponse contains the suggested move.
type SuggestMoveResponse struct {
    Move       string `json:"move"`        // UCI format (e.g., "e2e4")
    MoveSAN    string `json:"move_san"`    // SAN format (e.g., "e4")
    Evaluation int    `json:"evaluation"`  // Centipawns
}

// SuggestMove gets the best move for a position.
func (s *EngineService) SuggestMove(ctx context.Context, req SuggestMoveRequest) (*SuggestMoveResponse, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Set difficulty (depth)
    depth := req.Difficulty
    if depth <= 0 {
        depth = s.config.DefaultDepth
    }
    if depth > 20 {
        depth = 20
    }

    // Set position
    if req.FEN != "" {
        if err := s.engine.SetPosition(req.FEN); err != nil {
            return nil, fmt.Errorf("failed to set position: %w", err)
        }
    } else if len(req.Moves) > 0 {
        if err := s.engine.SetPositionWithMoves(req.Moves); err != nil {
            return nil, fmt.Errorf("failed to set moves: %w", err)
        }
    } else {
        return nil, fmt.Errorf("either FEN or moves must be provided")
    }

    // Get best move
    timeout := time.Duration(s.config.TimeoutMS) * time.Millisecond
    moveUCI, err := s.engine.BestMove(depth, timeout)
    if err != nil {
        return nil, fmt.Errorf("failed to get best move: %w", err)
    }

    // Convert UCI to SAN using chess.notnil/chess
    game, err := s.gameFromRequest(req)
    if err != nil {
        return nil, fmt.Errorf("failed to parse position: %w", err)
    }

    // Parse UCI move
    move, err := chess.UCINotation{}.Decode(game.Position(), moveUCI)
    if err != nil {
        return nil, fmt.Errorf("failed to parse UCI move: %w", err)
    }

    // Get SAN notation
    moveSAN := chess.AlgebraicNotation{}.Encode(game.Position(), move)

    // Get evaluation
    evaluation, _ := s.engine.Evaluate(depth)

    return &SuggestMoveResponse{
        Move:       moveUCI,
        MoveSAN:    moveSAN,
        Evaluation: evaluation,
    }, nil
}

// gameFromRequest creates a chess.Game from request.
func (s *EngineService) gameFromRequest(req SuggestMoveRequest) (*chess.Game, error) {
    if req.FEN != "" {
        fen, err := chess.FEN(req.FEN)
        if err != nil {
            return nil, err
        }
        return chess.NewGame(fen), nil
    }

    game := chess.NewGame()
    for _, moveUCI := range req.Moves {
        move, err := chess.UCINotation{}.Decode(game.Position(), moveUCI)
        if err != nil {
            return nil, fmt.Errorf("invalid move %s: %w", moveUCI, err)
        }
        if err := game.Move(move); err != nil {
            return nil, fmt.Errorf("illegal move %s: %w", moveUCI, err)
        }
    }
    return game, nil
}

// EvaluatePosition evaluates a chess position.
func (s *EngineService) EvaluatePosition(ctx context.Context, fen string, depth int) (int, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if depth <= 0 {
        depth = s.config.DefaultDepth
    }

    if err := s.engine.SetPosition(fen); err != nil {
        return 0, err
    }

    return s.engine.Evaluate(depth)
}

// Close shuts down the engine.
func (s *EngineService) Close() error {
    return s.engine.Close()
}
```

---

## Phase 4: API Endpoints

### Step 4.1: Create Engine Handler

**File:** `internal/handler/v1/engine/handler.go`

```go
// Package engine handles chess engine-related HTTP requests for API v1.
package engine

import (
    "net/http"

    "github.com/ankits1626/chess-coach-backend/internal/service"
    "github.com/gin-gonic/gin"
)

// Handler handles engine endpoints.
type Handler struct {
    engineService *service.EngineService
}

// NewHandler creates engine handler.
func NewHandler(engineService *service.EngineService) *Handler {
    return &Handler{engineService: engineService}
}

// SuggestMoveRequest represents move suggestion request.
type SuggestMoveRequest struct {
    FEN        string   `json:"fen,omitempty"`
    Moves      []string `json:"moves,omitempty"`
    Difficulty int      `json:"difficulty" binding:"omitempty,min=1,max=20"`
}

// SuggestMove suggests best move for position.
// @Summary Suggest move
// @Description Get engine's best move for a position
// @Tags engine
// @Accept json
// @Produce json
// @Param request body SuggestMoveRequest true "Position and difficulty"
// @Success 200 {object} service.SuggestMoveResponse
// @Router /engine/suggest [post]
func (h *Handler) SuggestMove(c *gin.Context) {
    var req SuggestMoveRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Validate: need either FEN or moves
    if req.FEN == "" && len(req.Moves) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "either fen or moves must be provided"})
        return
    }

    serviceReq := service.SuggestMoveRequest{
        FEN:        req.FEN,
        Moves:      req.Moves,
        Difficulty: req.Difficulty,
    }

    response, err := h.engineService.SuggestMove(c.Request.Context(), serviceReq)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, response)
}

// EvaluateRequest represents position evaluation request.
type EvaluateRequest struct {
    FEN   string `json:"fen" binding:"required"`
    Depth int    `json:"depth" binding:"omitempty,min=1,max=20"`
}

// EvaluatePosition evaluates a chess position.
// @Summary Evaluate position
// @Description Get centipawn evaluation for a position
// @Tags engine
// @Accept json
// @Produce json
// @Param request body EvaluateRequest true "FEN and depth"
// @Success 200 {object} map[string]int
// @Router /engine/evaluate [post]
func (h *Handler) EvaluatePosition(c *gin.Context) {
    var req EvaluateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    evaluation, err := h.engineService.EvaluatePosition(c.Request.Context(), req.FEN, req.Depth)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "evaluation": evaluation,
        "fen":        req.FEN,
    })
}
```

---

### Step 4.2: Register Engine Routes

**Update:** `internal/router/v1_routes.go`

```go
import (
    // ... existing imports
    "github.com/ankits1626/chess-coach-backend/internal/handler/v1/engine"  // ADD
    "github.com/ankits1626/chess-coach-backend/internal/service"            // ADD
)

func RegisterV1Routes(r *gin.Engine, db *database.DB) {
    v1 := r.Group("/api/v1")
    {
        // ... existing routes ...

        // Engine routes - ADD THIS BLOCK
        engineService, err := service.NewEngineService(nil)
        if err != nil {
            panic(fmt.Sprintf("Failed to initialize engine: %v", err))
        }
        engineHandler := engine.NewHandler(engineService)

        engineGroup := v1.Group("/engine")
        {
            engineGroup.POST("/suggest", engineHandler.SuggestMove)
            engineGroup.POST("/evaluate", engineHandler.EvaluatePosition)
        }
    }
}
```

---

## Phase 5: Testing

### Step 5.1: Build and Run

```bash
go mod tidy
go build -o bin/server cmd/server/main.go
./bin/server
```

**Expected output:**
```
Database connected successfully
Server starting on :8080
```

---

### Step 5.2: Test Move Suggestion

```bash
# Suggest move from starting position
curl -X POST http://localhost:8080/api/v1/engine/suggest \
  -H "Content-Type: application/json" \
  -d '{
    "fen": "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
    "difficulty": 10
  }' | jq .
```

**Expected response:**
```json
{
  "move": "e2e4",
  "move_san": "e4",
  "evaluation": 25
}
```

---

### Step 5.3: Test with Move History

```bash
# After 1.e4 e5, get move suggestion
curl -X POST http://localhost:8080/api/v1/engine/suggest \
  -H "Content-Type: application/json" \
  -d '{
    "moves": ["e2e4", "e7e5"],
    "difficulty": 10
  }' | jq .
```

**Expected:** Stockfish suggests `g1f3` (Nf3)

---

### Step 5.4: Test Position Evaluation

```bash
curl -X POST http://localhost:8080/api/v1/engine/evaluate \
  -H "Content-Type: application/json" \
  -d '{
    "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1",
    "depth": 15
  }' | jq .
```

**Expected:**
```json
{
  "evaluation": 30,
  "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"
}
```

**Evaluation meanings:**
- Positive = White is better
- Negative = Black is better
- ~100 centipawns = 1 pawn advantage

---

## Database Schema Updates

### Add Game Type and Difficulty

**File:** `db/migrations/000002_add_engine_support.up.sql`

```sql
-- Add game type
ALTER TABLE games ADD COLUMN game_type VARCHAR(20) DEFAULT 'human_vs_human';
ALTER TABLE games ADD CONSTRAINT valid_game_type
  CHECK (game_type IN ('human_vs_human', 'human_vs_computer', 'computer_vs_computer'));

-- Add difficulty for computer games
ALTER TABLE games ADD COLUMN difficulty INTEGER;
ALTER TABLE games ADD CONSTRAINT valid_difficulty
  CHECK (difficulty IS NULL OR (difficulty >= 1 AND difficulty <= 20));

-- Add evaluation to moves
ALTER TABLE moves ADD COLUMN evaluation INTEGER;

-- Create computer player
INSERT INTO users (id, username, email, rating)
VALUES (
  '00000000-0000-0000-0000-000000000001',
  'Computer',
  'computer@chesscoach.ai',
  3200
) ON CONFLICT DO NOTHING;
```

**Down migration:** `db/migrations/000002_add_engine_support.down.sql`

```sql
ALTER TABLE games DROP COLUMN IF EXISTS game_type;
ALTER TABLE games DROP COLUMN IF EXISTS difficulty;
ALTER TABLE moves DROP COLUMN IF EXISTS evaluation;
DELETE FROM users WHERE username = 'Computer';
```

---

## Frontend Integration

### Step 6.1: Create Engine Service

**File:** `frontend/src/services/engineService.ts`

```typescript
import axios from 'axios';

const API_BASE = 'http://localhost:8080/api/v1';

export interface SuggestMoveRequest {
  fen?: string;
  moves?: string[];
  difficulty?: number;
}

export interface SuggestMoveResponse {
  move: string;        // UCI format
  move_san: string;    // SAN format
  evaluation: number;  // Centipawns
}

export interface EvaluateRequest {
  fen: string;
  depth?: number;
}

export interface EvaluateResponse {
  evaluation: number;
  fen: string;
}

export const engineService = {
  async suggestMove(request: SuggestMoveRequest): Promise<SuggestMoveResponse> {
    const response = await axios.post(`${API_BASE}/engine/suggest`, request);
    return response.data;
  },

  async evaluatePosition(request: EvaluateRequest): Promise<EvaluateResponse> {
    const response = await axios.post(`${API_BASE}/engine/evaluate`, request);
    return response.data;
  },
};
```

---

### Step 6.2: Update Game Store

**File:** `frontend/src/store/gameStore.ts`

```typescript
import { create } from 'zustand';
import { Chess } from 'chess.js';
import { engineService } from '../services/engineService';

interface GameState {
  game: Chess;
  isComputerGame: boolean;
  isComputerTurn: boolean;
  computerDifficulty: number;
  thinking: boolean;

  // Actions
  makeComputerMove: () => Promise<void>;
  setComputerGame: (enabled: boolean, difficulty: number) => void;
}

export const useGameStore = create<GameState>((set, get) => ({
  game: new Chess(),
  isComputerGame: false,
  isComputerTurn: false,
  computerDifficulty: 10,
  thinking: false,

  makeComputerMove: async () => {
    const { game, computerDifficulty } = get();

    set({ thinking: true });

    try {
      const response = await engineService.suggestMove({
        fen: game.fen(),
        difficulty: computerDifficulty,
      });

      // Make the move
      game.move(response.move_san);
      set({
        game: new Chess(game.fen()),
        isComputerTurn: false,
        thinking: false
      });
    } catch (error) {
      console.error('Computer move failed:', error);
      set({ thinking: false });
    }
  },

  setComputerGame: (enabled, difficulty) => {
    set({
      isComputerGame: enabled,
      computerDifficulty: difficulty
    });
  },
}));
```

---

### Step 6.3: Add Computer Game UI

**File:** `frontend/src/components/ComputerGameControls.tsx`

```typescript
import { useState } from 'react';
import { useGameStore } from '../store/gameStore';

export const ComputerGameControls = () => {
  const [difficulty, setDifficulty] = useState(10);
  const { isComputerGame, thinking, setComputerGame } = useGameStore();

  const startComputerGame = () => {
    setComputerGame(true, difficulty);
  };

  const stopComputerGame = () => {
    setComputerGame(false, 10);
  };

  return (
    <div className="p-4 border rounded">
      <h3 className="font-bold mb-2">Play vs Computer</h3>

      {!isComputerGame ? (
        <>
          <label className="block mb-2">
            Difficulty: {difficulty}
            <input
              type="range"
              min="1"
              max="20"
              value={difficulty}
              onChange={(e) => setDifficulty(parseInt(e.target.value))}
              className="w-full"
            />
          </label>
          <button
            onClick={startComputerGame}
            className="bg-blue-500 text-white px-4 py-2 rounded"
          >
            Start Game
          </button>
        </>
      ) : (
        <div>
          <p className="mb-2">Difficulty: {difficulty}</p>
          {thinking && <p className="text-blue-500">Computer is thinking...</p>}
          <button
            onClick={stopComputerGame}
            className="bg-red-500 text-white px-4 py-2 rounded"
          >
            End Game
          </button>
        </div>
      )}
    </div>
  );
};
```

---

## Summary

### What We Built

✅ **UCI Protocol Client** - Communicates with Stockfish
✅ **Engine Service** - Business logic for move suggestions
✅ **API Endpoints** - `/engine/suggest` and `/engine/evaluate`
✅ **Database Schema** - Game type, difficulty, evaluation
✅ **Frontend Integration** - Engine service and UI controls

### Testing Checklist

- [ ] Stockfish installed and working
- [ ] UCI client can communicate with engine
- [ ] Move suggestion endpoint works
- [ ] Position evaluation endpoint works
- [ ] Difficulty levels work (1-20)
- [ ] Frontend can request computer moves
- [ ] Computer thinking indicator shows
- [ ] Game is recorded in database

---

## Next Steps

1. **Add Game Flow Integration**
   - Automatically call computer move after user move
   - Handle game end detection
   - Save computer games to database

2. **Enhanced Features**
   - Opening book integration
   - Multi-line analysis (show alternative moves)
   - Blunder detection
   - Hint system

3. **Performance Optimization**
   - Engine pooling for concurrent games
   - Caching common positions
   - Background position analysis

---

**Last Updated**: 2025-11-01

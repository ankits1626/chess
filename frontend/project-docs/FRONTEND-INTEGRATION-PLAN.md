# Frontend Integration Plan: Computer Player

**Date**: 2025-11-06

**Status**: Ready to implement

---

## 🎯 Goal

Integrate the backend computer player functionality into the React frontend, allowing users to play chess against the computer with selectable difficulty levels.

---

## 📋 Current Frontend State

### Technology Stack:
- **React** 19.2.0 with TypeScript
- **State Management**: Zustand 5.0.8
- **Chess Logic**: chess.js 1.4.0
- **Build Tool**: Vite 5.0.0
- **Styling**: Tailwind CSS 4.1.16

### Current Game Modes:
1. **'live'**: Local play (human vs human on same device)
2. **'replay'**: View imported games from Chess.com

### Key Files:
- State: `src/store/useGameStore.ts` (310 lines)
- Services: `src/services/chesscomApi.ts` (110 lines)
- UI: `src/components/game/GameInfo.tsx` (55 lines)
- Board: `src/components/board/GameBoard.tsx` (110 lines)

---

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        User Interface                           │
├─────────────────────────────────────────────────────────────────┤
│  GameSetup Modal                                                │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │ • Select Opponent: [Human | Computer]                    │   │
│  │ • Select Color: [White | Black]                          │   │
│  │ • Select Difficulty: [Easy | Medium | Hard]              │   │
│  │ • Time Control: [5+0 | 10+0 | 15+10]                     │   │
│  │                                         [Start Game]     │   │
│  └──────────────────────────────────────────────────────────┘   │
│                              ↓                                  │
│  GameBoard + GameInfo                                           │
│  ┌──────────────────────────────────────────────────────────-┐  │
│  │ Board                  │  Game Info                       │  │
│  │ [Chess Board UI]       │  Turn: White                     │  │
│  │                        │  Status: Computer thinking...    │  │
│  │                        │  Opponent: Computer (Medium)     │  │
│  └───────────────────────────────────────────────────────-───┘  │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                      State Management (Zustand)                 │
├─────────────────────────────────────────────────────────────────┤
│  useGameStore                                                   │
│  • mode: 'live' | 'replay' | 'computer'                         │
│  • opponentType: 'human' | 'computer'                           │
│  • computerDifficulty: 'easy' | 'medium' | 'hard'               │
│  • isComputerThinking: boolean                                  │
│  • playerColor: 'white' | 'black'                               │
│  • gameId: string                                               │
│                                                                 │
│  Actions:                                                       │
│  • startComputerGame(color, difficulty)                         │
│  • handleComputerMove(move)                                     │
│  • selectSquare(square) [modified]                              │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                      Services Layer                             │
├─────────────────────────────────────────────────────────────────┤
│  GameService (WebSocket)                                        │
│  • connect(userId)                                              │
│  • createGame(mode, difficulty, timeControl)                    │
│  • makeMove(gameId, move)                                       │
│  • onMoveEvent(callback)                                        │
│  • onGameEnd(callback)                                          │
│  • disconnect()                                                 │
└─────────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────────┐
│                      Backend (WebSocket)                        │
│  ws://localhost:8080/ws?user_id=<UUID>                          │
└─────────────────────────────────────────────────────────────────┘
```

---

## 📝 Implementation Steps

### Step 1: Create WebSocket Game Service (2 hours)

**File**: `src/services/gameService.ts`

**Features**:
- WebSocket connection management
- Message sending/receiving with request/response pattern
- Event handlers for move events
- Reconnection logic
- Error handling

**Key Methods**:
```typescript
class GameService {
  connect(userId: string): Promise<void>
  createGame(mode: GameMode, difficulty: Difficulty, timeControl: string): Promise<GameInfo>
  makeMove(gameId: string, move: string): Promise<MoveResult>
  onMoveEvent(callback: (move: MoveEvent) => void): void
  onGameEnd(callback: (result: GameResult) => void): void
  disconnect(): void
}
```

---

### Step 2: Update GameStore (2 hours)

**File**: `src/store/useGameStore.ts`

**New State Fields**:
```typescript
interface GameState {
  // Existing fields...

  // NEW: Computer game fields
  opponentType: 'human' | 'computer'
  computerDifficulty: 'easy' | 'medium' | 'hard'
  isComputerThinking: boolean
  playerColor: 'white' | 'black'
  gameId: string | null
  connectionStatus: 'disconnected' | 'connecting' | 'connected' | 'error'
}
```

**New Actions**:
```typescript
interface GameActions {
  // NEW: Computer game actions
  startComputerGame: (color: 'white' | 'black', difficulty: Difficulty) => Promise<void>
  handleComputerMove: (moveUCI: string) => void
  setComputerThinking: (thinking: boolean) => void
  endGame: (result: GameResult) => void
}
```

**Modified Actions**:
- `selectSquare()`: After human move, check if computer's turn and wait for computer response
- `reset()`: Also reset computer game state

---

### Step 3: Create GameSetup Component (1.5 hours)

**File**: `src/components/game/GameSetup.tsx`

**UI Elements**:
```tsx
<GameSetup onStart={handleStart}>
  <OpponentSelector
    value={opponent}
    onChange={setOpponent}
    options={['human', 'computer']}
  />

  {opponent === 'computer' && (
    <>
      <ColorSelector
        value={color}
        onChange={setColor}
        options={['white', 'black']}
      />
      <DifficultySelector
        value={difficulty}
        onChange={setDifficulty}
        options={['easy', 'medium', 'hard']}
      />
    </>
  )}

  <TimeControlSelector
    value={timeControl}
    onChange={setTimeControl}
    options={['5+0', '10+0', '15+10']}
  />

  <StartButton onClick={handleStartGame} />
</GameSetup>
```

---

### Step 4: Update GameInfo Component (1 hour)

**File**: `src/components/game/GameInfo.tsx`

**New Display Elements**:
- Computer thinking indicator with loading animation
- Opponent type and difficulty display
- Connection status indicator
- Error messages from WebSocket

**Example**:
```tsx
{opponentType === 'computer' && (
  <div className="opponent-info">
    <span>Opponent: Computer ({difficulty})</span>
    {isComputerThinking && (
      <div className="thinking-indicator">
        <Spinner />
        <span>Computer is thinking...</span>
      </div>
    )}
  </div>
)}
```

---

### Step 5: Integrate WebSocket Events (1.5 hours)

**File**: `src/App.tsx` or create `src/hooks/useGameConnection.ts`

**Event Handlers**:
```typescript
useEffect(() => {
  if (opponentType === 'computer' && gameId) {
    // Listen for computer moves
    gameService.onMoveEvent((moveEvent) => {
      handleComputerMove(moveEvent.moveUCI)
      setComputerThinking(false)
    })

    // Listen for game end
    gameService.onGameEnd((result) => {
      endGame(result)
    })

    // Cleanup on unmount
    return () => gameService.disconnect()
  }
}, [opponentType, gameId])
```

---

### Step 6: Update Move Flow (1 hour)

**File**: `src/store/useGameStore.ts`

**Modified `selectSquare()` Logic**:
```typescript
selectSquare: (square: Square) => {
  // Existing move logic...

  // After successful human move
  if (get().opponentType === 'computer') {
    const playerColor = get().playerColor
    const currentTurn = get().game.turn() // 'w' or 'b'
    const isComputerTurn = (playerColor === 'white' && currentTurn === 'b') ||
                           (playerColor === 'black' && currentTurn === 'w')

    if (isComputerTurn) {
      set({ isComputerThinking: true })
      // Computer move will come via WebSocket event
    }
  }
}
```

---

### Step 7: Error Handling & Edge Cases (1 hour)

**Scenarios to Handle**:
1. WebSocket connection failure
2. Computer move timeout
3. Invalid move from computer (shouldn't happen but handle gracefully)
4. User disconnection during game
5. Reconnection with game state recovery

**Implementation**:
- Add error state to GameStore
- Show error notifications to user
- Implement reconnection with exponential backoff
- Store gameId in localStorage for recovery

---

## 🎮 User Flow

### Starting a Computer Game:

1. User clicks "New Game" button
2. GameSetup modal appears
3. User selects:
   - Opponent: Computer
   - Color: White (or Black)
   - Difficulty: Medium
   - Time Control: 5+0
4. User clicks "Start Game"
5. Frontend:
   - Connects to WebSocket
   - Sends `createGame` request with mode="human_vs_computer"
   - Receives gameId and initial state
6. If user chose Black:
   - Computer makes first move automatically
   - Frontend receives move event
   - Board updates

### Playing Moves:

1. User clicks source square → highlights valid moves
2. User clicks destination square → move executes
3. Frontend sends `makeMove` to backend
4. Backend validates and applies move
5. Backend triggers computer move calculation
6. Frontend shows "Computer is thinking..." indicator
7. Computer move event arrives via WebSocket
8. Frontend applies computer move
9. Board updates, thinking indicator disappears
10. Repeat until game over

### Game End:

1. Checkmate/stalemate detected
2. `gameEnd` event received from backend
3. Frontend shows result modal
4. Options: New Game, Review Game, Exit

---

## 🗂️ File Structure

```
frontend/app/src/
├── services/
│   ├── chesscomApi.ts         (existing)
│   └── gameService.ts         (NEW - WebSocket service)
│
├── store/
│   └── useGameStore.ts        (MODIFY - add computer state/actions)
│
├── components/
│   ├── game/
│   │   ├── GameInfo.tsx       (MODIFY - show computer status)
│   │   ├── GameSetup.tsx      (NEW - game creation modal)
│   │   └── GameResult.tsx     (NEW - end game modal)
│   │
│   ├── board/
│   │   └── GameBoard.tsx      (minor mods for computer mode)
│   │
│   └── ui/
│       ├── Spinner.tsx        (NEW - loading indicator)
│       └── Select.tsx         (NEW - reusable selector)
│
├── hooks/
│   └── useGameConnection.ts   (NEW - WebSocket event handling)
│
└── types/
    └── game.ts                (MODIFY - add computer types)
```

---

## 🔧 TypeScript Types

```typescript
// src/types/game.ts

export type GameMode = 'live' | 'replay' | 'computer'

export type OpponentType = 'human' | 'computer'

export type Difficulty = 'easy' | 'medium' | 'hard'

export type PlayerColor = 'white' | 'black'

export type ConnectionStatus = 'disconnected' | 'connecting' | 'connected' | 'error'

export interface GameInfo {
  gameId: string
  status: 'waiting' | 'active' | 'finished'
  side: PlayerColor
  mode: string
  difficulty: Difficulty
  timeControl: string
  fen: string
}

export interface MoveResult {
  moveNumber: number
  san: string
  uci: string
  fen: string
  gameOver: boolean
  result?: GameResult
}

export interface MoveEvent {
  gameId: string
  moveSAN: string
  moveUCI: string
  fen: string
  moveNumber: number
  isGameOver: boolean
  result: GameResult | null
}

export interface GameResult {
  winner: 'white' | 'black' | 'draw'
  method: 'checkmate' | 'resignation' | 'stalemate' | 'timeout' | 'draw'
  pgn: string
}

export interface WebSocketMessage {
  id: string
  type: 'request' | 'response' | 'event'
  action?: string
  event?: string
  data?: any
  error?: string
  success?: boolean
}
```

---

## ✅ Testing Checklist

### Unit Tests:
- [ ] GameService connection/disconnection
- [ ] GameService message sending/receiving
- [ ] GameStore computer game state transitions
- [ ] GameStore move validation with computer opponent

### Integration Tests:
- [ ] Complete game flow: start → moves → end
- [ ] Computer move response handling
- [ ] Error scenarios (connection loss, timeout)
- [ ] Reconnection logic

### Manual Testing:
- [ ] Create game as white vs computer (easy)
- [ ] Create game as black vs computer (medium)
- [ ] Play complete game to checkmate
- [ ] Test resignation
- [ ] Test connection loss during game
- [ ] Test all difficulty levels
- [ ] Test different time controls
- [ ] Test switching between game modes

---

## 📊 Estimated Timeline

| Step | Description | Time | Dependencies |
|------|-------------|------|--------------|
| 1 | WebSocket Service | 2h | None |
| 2 | GameStore Updates | 2h | Step 1 |
| 3 | GameSetup Component | 1.5h | Step 2 |
| 4 | GameInfo Updates | 1h | Step 2 |
| 5 | WebSocket Integration | 1.5h | Steps 1-2 |
| 6 | Move Flow Updates | 1h | Steps 1-2 |
| 7 | Error Handling | 1h | All previous |
| 8 | Testing & Polish | 2h | All previous |
| **Total** | | **12h** | |

---

## 🎯 Success Criteria

✅ User can create a game against the computer
✅ User can select difficulty level
✅ User can play as white or black
✅ Computer responds automatically after human moves
✅ "Computer thinking" indicator shows during calculation
✅ Game ends properly with correct result
✅ Error states handled gracefully
✅ UI is responsive and intuitive
✅ No console errors or warnings

---

## 🚀 Next Steps

1. Review this plan with the team
2. Start with Step 1 (WebSocket Service)
3. Iterate through steps 2-7
4. Comprehensive testing
5. Deploy to staging
6. User acceptance testing
7. Production deployment

---

**Ready to start implementation!** 🎮♟️

# Chess Coach Frontend - Complete Execution Flow Guide for Beginners

> **Last Updated:** 2025-11-07
> **Status:** ✅ Fully Implemented & Working
> **Purpose:** Understand how the React frontend runs from start to finish, including state management with Zustand, WebSocket real-time communication, chess.js integration, and UI rendering with Tailwind CSS.

---

## Table of Contents

1. [Quick Overview](#quick-overview)
2. [Directory Structure](#directory-structure)
3. [The Big Picture: How Everything Works Together](#the-big-picture)
4. [Phase 1: Application Startup](#phase-1-application-startup)
5. [Phase 2: State Management with Zustand](#phase-2-state-management-with-zustand)
6. [Phase 3: Component Hierarchy & Rendering](#phase-3-component-hierarchy--rendering)
7. [Phase 4: User Interactions & Game Logic](#phase-4-user-interactions--game-logic)
8. [Phase 5: WebSocket Communication](#phase-5-websocket-communication)
9. [Phase 6: Computer Game Flow](#phase-6-computer-game-flow)
10. [Phase 7: Replay Mode & PGN Import](#phase-7-replay-mode--pgn-import)
11. [Complete User Flow Examples](#complete-user-flow-examples)
12. [Key Design Patterns](#key-design-patterns)
13. [Troubleshooting & Common Issues](#troubleshooting--common-issues)

---

## Quick Overview

The Chess Coach frontend is a React application that provides:

- **Interactive Chess Board** with drag-and-drop piece movement
- **Live Chess Games** against human or computer opponents
- **PGN Replay Mode** for reviewing games from Chess.com
- **WebSocket Integration** for real-time game updates
- **State Management** using Zustand for global state
- **Responsive UI** built with Tailwind CSS

**Tech Stack:**
- **Framework:** React 19.2.0
- **Build Tool:** Vite 5.x
- **State Management:** Zustand 5.x
- **Chess Logic:** chess.js 1.4.0
- **Styling:** Tailwind CSS 4.x
- **Language:** TypeScript 5.9.3
- **Real-time Communication:** WebSocket API

---

## Directory Structure

```
frontend/app/
├── src/
│   ├── main.tsx                           # 🚀 ENTRY POINT - React app starts here
│   ├── App.tsx                            # 📱 Root application component
│   │
│   ├── store/
│   │   └── useGameStore.ts                # 🗃️ Zustand global state store
│   │
│   ├── services/
│   │   ├── gameService.ts                 # 🌐 WebSocket game service
│   │   └── chesscomApi.ts                 # 🔌 Chess.com API integration
│   │
│   ├── components/
│   │   ├── board/
│   │   │   ├── GameBoard.tsx              # ♟️ Main chessboard component
│   │   │   ├── Square.tsx                 # ⬜ Individual square component
│   │   │   └── Piece.tsx                  # 🎭 Chess piece renderer
│   │   │
│   │   ├── game/
│   │   │   ├── GameController.tsx         # 🎮 Game controls & logic
│   │   │   ├── GameInfo.tsx               # ℹ️ Game status display
│   │   │   ├── GameSetup.tsx              # ⚙️ Computer game setup dialog
│   │   │   ├── MoveHistory.tsx            # 📜 Move list display
│   │   │   └── PromotionDialog.tsx        # 👑 Pawn promotion UI
│   │   │
│   │   ├── replay/
│   │   │   ├── ReplayControls.tsx         # ⏯️ Replay navigation controls
│   │   │   ├── MoveList.tsx               # 📋 Clickable move list
│   │   │   └── PlayerDisplay.tsx          # 👤 Player name display
│   │   │
│   │   ├── importer/
│   │   │   ├── GameImporter.tsx           # 📥 Import games from Chess.com
│   │   │   ├── GameList.tsx               # 📚 List of imported games
│   │   │   ├── GameListItem.tsx           # 📄 Single game item
│   │   │   └── ArchivePaginator.tsx       # 📖 Pagination for archives
│   │   │
│   │   └── debug/
│   │       ├── DebugPanel.tsx             # 🐛 Development debug tools
│   │       ├── FenLoader.tsx              # 🔧 Load position from FEN
│   │       └── ScenarioPicker.tsx         # 🎯 Test scenarios
│   │
│   ├── hooks/
│   │   ├── useGameController.ts           # 🎯 Game control hook
│   │   └── useDebugPanel.ts               # 🐛 Debug panel hook
│   │
│   ├── types/
│   │   ├── chess.ts                       # ♟️ Chess-specific types
│   │   └── game.ts                        # 🎮 Game mode & WebSocket types
│   │
│   └── constants/
│       └── debugScenarios.ts              # 🧪 Test scenarios
│
├── public/                                # 📁 Static assets
├── index.html                             # 🌐 HTML entry point
├── vite.config.ts                         # ⚙️ Vite configuration
├── tsconfig.json                          # 📝 TypeScript config
├── tailwind.config.js                     # 🎨 Tailwind config
└── package.json                           # 📦 Dependencies
```

---

## The Big Picture

### System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                        BROWSER LAYER                            │
│   - User interacts with chess board                            │
│   - Clicks squares, selects pieces, makes moves                │
└────────────────┬────────────────────────────────────┬───────────┘
                 │                                    │
                 ▼                                    ▼
┌────────────────────────────┐    ┌──────────────────────────────┐
│    REACT COMPONENTS        │    │    WEBSOCKET SERVICE         │
│    (UI Layer)              │    │    (gameService.ts)          │
│                            │    │                              │
│  • GameBoard               │◄───┤  • Connect/Disconnect        │
│  • Square                  │    │  • Send requests             │
│  • Piece                   │    │  • Receive events            │
│  • GameInfo                │    │  • Auto-reconnect            │
│  • GameSetup               │    │                              │
│  • ReplayControls          │    │                              │
└────────────┬───────────────┘    └──────────┬───────────────────┘
             │                               │
             ▼                               ▼
┌──────────────────────────────────────────────────────────────────┐
│                    ZUSTAND STORE                                 │
│                  (useGameStore.ts)                               │
│                                                                  │
│  • Global State Management                                      │
│  • Game state (Chess instance)                                  │
│  • Selected square & valid moves                                │
│  • Game mode (live/replay/computer)                             │
│  • Player colors & connection status                            │
│  • Actions: selectSquare, makeMove, handlePromotion             │
└─────────────────────────────┬────────────────────────────────────┘
                              │
             ┌────────────────┴────────────────┐
             ▼                                 ▼
┌──────────────────────────┐    ┌──────────────────────────────┐
│    CHESS.JS LIBRARY      │    │  CHESS.COM API               │
│    (Game Logic)          │    │  (chesscomApi.ts)            │
│                          │    │                              │
│  • Chess rules engine    │    │  • Fetch user archives       │
│  • Move validation       │    │  • Fetch monthly games       │
│  • Game state (FEN)      │    │  • Parse PGN data            │
│  • Checkmate detection   │    │                              │
│  • PGN parsing           │    │                              │
└──────────────────────────┘    └──────────────────────────────┘
             │                                 │
             ▼                                 ▼
┌──────────────────────────────────────────────────────────────────┐
│                   BACKEND SERVER (Go)                            │
│                                                                  │
│  • WebSocket hub                                                 │
│  • Game management                                               │
│  • Stockfish AI integration                                     │
│  • Database persistence                                          │
└──────────────────────────────────────────────────────────────────┘
```

---

## Phase 1: Application Startup

### Entry Point: `src/main.tsx`

When you run `npm run dev`, Vite starts a development server and this is what happens:

```
┌─────────────────────────────────────────────────────────────────┐
│                    APPLICATION BOOTSTRAP                        │
└─────────────────────────────────────────────────────────────────┘

1️⃣  Vite Dev Server Starts
    ├─ Loads index.html
    ├─ Processes TypeScript/JSX
    ├─ Applies Tailwind CSS
    └─ Hot Module Replacement (HMR) enabled
        ↓

2️⃣  Load main.tsx (Entry Point)
    ├─ Import React & ReactDOM
    ├─ Import root App component
    ├─ Import global CSS (Tailwind)
    └─ File: src/main.tsx
        ↓

3️⃣  Initialize React Root
    ├─ Find root element in DOM
    ├─ Create React root (createRoot)
    └─ Enable StrictMode (dev checks)
        ↓

4️⃣  Render App Component
    ├─ Mount <App /> to DOM
    ├─ Initialize Zustand store
    ├─ Render component tree
    └─ File: src/App.tsx
        ↓

5️⃣  Application Ready
    ├─ UI visible in browser
    ├─ WebSocket service ready (not connected)
    └─ User can interact with board
```

### Code Flow in `main.tsx`

```tsx
// File: src/main.tsx

import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'         // Tailwind CSS
import App from './App.tsx'

// Step 1: Find root element
createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />  // Step 2: Render App component
  </StrictMode>,
)
```

---

## Phase 2: State Management with Zustand

### Zustand Store: `src/store/useGameStore.ts`

Zustand provides a lightweight, hook-based state management solution:

```
┌─────────────────────────────────────────────────────────────────┐
│                    ZUSTAND STORE STRUCTURE                      │
└─────────────────────────────────────────────────────────────────┘

STATE:
├─ game: Chess                    # chess.js instance
├─ selectedSquare: Square | null  # Currently selected square
├─ validMoves: Square[]           # Legal moves for selected piece
├─ pendingMove: PendingPromotion  # Move awaiting promotion choice
├─ lastMove: LastMove             # Previous move (from, to)
├─ mode: GameMode                 # 'live' | 'replay' | 'computer'
│
├─ REPLAY STATE:
│  ├─ replayMoves: Move[]         # All moves in replay
│  ├─ replayIndex: number         # Current move index
│  ├─ isAutoplaying: boolean      # Auto-advance moves
│  ├─ whitePlayer: string         # White player name
│  └─ blackPlayer: string         # Black player name
│
├─ COMPUTER GAME STATE:
│  ├─ opponentType: OpponentType  # 'human' | 'computer'
│  ├─ computerDifficulty          # 'easy' | 'medium' | 'hard'
│  ├─ isComputerThinking: boolean # AI calculating move
│  ├─ playerColor: PlayerColor    # 'white' | 'black'
│  ├─ gameId: string              # WebSocket game ID
│  ├─ connectionStatus            # WebSocket status
│  └─ gameError: string           # Error messages
│
└─ IMPORTER STATE:
   ├─ gameArchives: string[]      # Chess.com archive URLs
   ├─ importedGames: Game[]       # Fetched games
   └─ isGameListLoading: boolean  # Loading indicator

ACTIONS:
├─ selectSquare(square)           # Handle square click
├─ handlePromotion(piece)         # Choose promotion piece
├─ resetGame()                    # Clear all state
├─ loadPgn(pgn)                   # Load PGN for replay
│
├─ REPLAY ACTIONS:
│  ├─ goToMove(index)             # Jump to specific move
│  ├─ nextMove()                  # Advance one move
│  ├─ prevMove()                  # Go back one move
│  ├─ toggleAutoplay()            # Start/stop auto-advance
│  └─ goToFirstMove()             # Go to starting position
│
├─ COMPUTER GAME ACTIONS:
│  ├─ startComputerGame()         # Connect & create game
│  ├─ handleComputerMove(uci)     # Apply AI move to board
│  ├─ endGame(result)             # Handle game end
│  └─ disconnectGame()            # Close WebSocket
│
└─ IMPORTER ACTIONS:
   ├─ fetchGameArchives(user)     # Get Chess.com archives
   ├─ fetchGamesForArchive(url)   # Get games from archive
   └─ loadPgnFromGame(game)       # Load game for replay
```

### How Zustand Works

```typescript
// Creating a Zustand store
export const useGameStore = create<GameState>((set, get) => ({
  // Initial state
  game: new Chess(),
  selectedSquare: null,

  // Actions that modify state
  selectSquare: (square) => {
    const { game } = get();  // Get current state
    // ... logic ...
    set({ selectedSquare: square });  // Update state
  },
}));

// Using in components
function MyComponent() {
  // Subscribe to specific state
  const game = useGameStore(state => state.game);
  const selectSquare = useGameStore(state => state.selectSquare);

  // Component re-renders when subscribed state changes
}
```

**Key Benefits:**
- 🚀 No Provider wrapper needed
- 🎯 Subscribe to specific state slices
- ⚡ Minimal re-renders
- 🔧 Simple API (set, get)

---

## Phase 3: Component Hierarchy & Rendering

### Component Tree

```
App.tsx (Root)
├─ GameBoard                      # Chess board with 64 squares
│  └─ Square (×64)                # Individual squares
│     └─ Piece                    # Chess piece image
│
├─ GameInfo                       # Game status, turn indicator
│  ├─ Turn display                # "White's turn" / "Black's turn"
│  ├─ Game status                 # Active / Checkmate / Stalemate
│  ├─ Board flip button           # Rotate board
│  └─ Reset button                # New game
│
├─ GameImporter                   # Import from Chess.com
│  ├─ Username input              # Chess.com username
│  ├─ ArchivePaginator            # Navigate archives
│  └─ GameList                    # List of games
│     └─ GameListItem (×N)        # Individual game
│
├─ ReplayControls                 # Replay navigation (if mode=replay)
│  ├─ First move button           # ⏮
│  ├─ Previous button             # ◀
│  ├─ Play/Pause button           # ⏯
│  ├─ Next button                 # ▶
│  └─ Last move button            # ⏭
│
├─ MoveList                       # Clickable move list (if mode=replay)
│  └─ Move items                  # 1. e4 e5 2. Nf3 Nc6 ...
│
├─ PlayerDisplay (×2)             # Player names (if mode=replay)
│  ├─ Top player                  # Depends on board orientation
│  └─ Bottom player               # Depends on board orientation
│
├─ PromotionDialog                # Pawn promotion UI (conditional)
│  └─ Piece selection             # Q, R, B, N
│
├─ GameSetup                      # Computer game setup (conditional)
│  ├─ Color selection             # White / Black
│  ├─ Difficulty selection        # Easy / Medium / Hard
│  └─ Start button                # Begin game
│
└─ DebugPanel (dev only)          # Development tools
   ├─ FEN loader                  # Load custom position
   └─ ScenarioPicker              # Test scenarios
```

### Rendering Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                    REACT RENDERING CYCLE                        │
└─────────────────────────────────────────────────────────────────┘

1️⃣  Initial Render
    ├─ App component mounts
    ├─ useGameStore() subscribes to state
    ├─ Components read initial state
    └─ Render full component tree
        ↓

2️⃣  User Interaction (e.g., click square)
    ├─ onClick handler fires
    ├─ Calls store action (selectSquare)
    └─ Action updates Zustand state
        ↓

3️⃣  Zustand Notifies Subscribers
    ├─ Components subscribed to changed state
    ├─ React schedules re-render
    └─ Only affected components re-render
        ↓

4️⃣  Component Re-renders
    ├─ Read new state from store
    ├─ Calculate derived values
    ├─ Update virtual DOM
    └─ React commits changes to real DOM
        ↓

5️⃣  UI Updates
    └─ User sees updated board state
```

**Example: Square Selection Flow**

```tsx
// User clicks square "e2"
<Square onClick={() => selectSquare("e2")} />

// selectSquare action in store
selectSquare: (square) => {
  const { game, selectedSquare } = get();

  if (!selectedSquare) {
    // First click: select piece
    const piece = game.get(square);
    if (piece && piece.color === game.turn()) {
      const moves = game.moves({ square, verbose: true });
      set({
        selectedSquare: square,
        validMoves: moves.map(m => m.to)
      });
    }
  } else {
    // Second click: make move
    const move = game.move({ from: selectedSquare, to: square });
    if (move) {
      set({
        game: {...game},
        selectedSquare: null,
        validMoves: [],
        lastMove: { from: selectedSquare, to: square }
      });
    }
  }
}

// Components re-render with new state
// - Selected square gets highlighted
// - Valid move indicators appear
```

---

## Phase 4: User Interactions & Game Logic

### Move Execution Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                    MOVE EXECUTION FLOW                          │
└─────────────────────────────────────────────────────────────────┘

1️⃣  User Clicks Source Square (e.g., "e2")
    ├─ GameBoard renders 64 Square components
    ├─ User clicks square with white pawn
    ├─ onClick={onSquareClick("e2")} fires
    └─ Calls store.selectSquare("e2")
        ↓

2️⃣  Select Piece & Show Valid Moves
    ├─ Check if piece exists on e2
    ├─ Check if piece color matches current turn
    ├─ Get legal moves using chess.js:
    │  const moves = game.moves({ square: "e2", verbose: true })
    │  // Returns: [{ from: "e2", to: "e3" }, { from: "e2", to: "e4" }]
    ├─ Update state:
    │  set({ selectedSquare: "e2", validMoves: ["e3", "e4"] })
    └─ Components re-render:
       • e2 square highlighted (isSelected=true)
       • e3, e4 squares show move indicators (isValidMove=true)
        ↓

3️⃣  User Clicks Destination Square (e.g., "e4")
    ├─ User clicks valid destination
    ├─ onClick={onSquareClick("e4")} fires
    └─ Calls store.selectSquare("e4")
        ↓

4️⃣  Execute Move
    ├─ Check if move is to promote pawn
    │  if (piece.type === 'p' && (square.endsWith('1') || square.endsWith('8')))
    │     └─ Show PromotionDialog (pending state)
    │
    ├─ Otherwise, execute move:
    │  const move = game.move({ from: "e2", to: "e4" })
    │  // chess.js validates and executes move
    │  // Returns: { san: "e4", from: "e2", to: "e4", ... }
    │
    ├─ Update state:
    │  set({
    │    game: {...game},  // Trigger re-render
    │    selectedSquare: null,
    │    validMoves: [],
    │    lastMove: { from: "e2", to: "e4" }
    │  })
    │
    └─ If computer game, send move to backend:
       gameService.makeMove(gameId, move.san)
        ↓

5️⃣  Board Updates
    ├─ GameBoard re-renders with new position
    ├─ Piece moved from e2 to e4
    ├─ Last move squares highlighted
    ├─ Turn indicator switches to black
    └─ Check for game over (checkmate/stalemate)
```

### Promotion Handling

```
┌─────────────────────────────────────────────────────────────────┐
│                    PAWN PROMOTION FLOW                          │
└─────────────────────────────────────────────────────────────────┘

1️⃣  User Moves Pawn to 8th Rank
    ├─ Detect promotion move:
    │  piece.type === 'p' && square.endsWith('8')
    ├─ Store pending move:
    │  set({ pendingMove: { from: "e7", to: "e8", color: "w" } })
    └─ Clear selection:
       set({ selectedSquare: null, validMoves: [] })
        ↓

2️⃣  Show Promotion Dialog
    ├─ PromotionDialog renders (conditional)
    ├─ Display 4 piece choices: Q, R, B, N
    ├─ User clicks Queen (Q)
    └─ Calls handlePromotion("q")
        ↓

3️⃣  Execute Promotion Move
    ├─ Get pending move from state
    ├─ Execute with promotion:
    │  const move = game.move({
    │    from: "e7",
    │    to: "e8",
    │    promotion: "q"
    │  })
    ├─ Update state:
    │  set({
    │    game: {...game},
    │    pendingMove: null,
    │    lastMove: { from: "e7", to: "e8" }
    │  })
    └─ If computer game, send to backend:
       gameService.makeMove(gameId, move.san)  // "e8=Q"
        ↓

4️⃣  Dialog Closes & Board Updates
    └─ Queen appears on e8
```

---

## Phase 5: WebSocket Communication

### WebSocket Service: `src/services/gameService.ts`

```
┌─────────────────────────────────────────────────────────────────┐
│                    WEBSOCKET ARCHITECTURE                       │
└─────────────────────────────────────────────────────────────────┘

GameService Class:
├─ ws: WebSocket                  # WebSocket connection
├─ pendingRequests: Map           # Request/response tracking
├─ eventHandlers: Map             # Event listeners
├─ reconnectAttempts: number      # Auto-reconnect counter
└─ Methods:
   ├─ connect(userId)             # Establish connection
   ├─ disconnect()                # Close connection
   ├─ sendRequest(action, data)   # Send request, await response
   ├─ handleMessage(data)         # Process incoming message
   ├─ on(event, handler)          # Register event listener
   └─ off(event, handler)         # Unregister listener
```

### WebSocket Message Protocol

```json
// REQUEST (from frontend to backend)
{
  "id": "req-1",
  "type": "request",
  "action": "createGame",
  "data": {
    "mode": "human_vs_computer",
    "difficulty": "medium",
    "timeControl": "5+0",
    "playerColor": "white"
  }
}

// RESPONSE (from backend to frontend)
{
  "id": "req-1",
  "type": "response",
  "success": true,
  "data": {
    "gameId": "abc-123",
    "status": "active",
    "side": "white",
    "fen": "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
  }
}

// EVENT (from backend to frontend)
{
  "id": "evt-1",
  "type": "event",
  "event": "move",
  "data": {
    "gameId": "abc-123",
    "moveSAN": "e5",
    "moveUCI": "e7e5",
    "fen": "...",
    "isGameOver": false
  }
}
```

### Connection Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                    WEBSOCKET CONNECTION FLOW                    │
└─────────────────────────────────────────────────────────────────┘

1️⃣  User Starts Computer Game
    ├─ User clicks "Play vs Computer"
    ├─ GameSetup dialog appears
    ├─ User selects white, medium difficulty
    └─ Clicks "Start Game"
        ↓

2️⃣  Connect to WebSocket
    ├─ Call gameService.connect(userId)
    ├─ Create WebSocket:
    │  ws = new WebSocket("ws://localhost:8080/ws?user_id=...")
    ├─ Set up event handlers:
    │  • onopen → resolve promise
    │  • onmessage → handleMessage()
    │  • onerror → reject promise
    │  • onclose → handleDisconnect()
    └─ Wait for connection
        ↓

3️⃣  Connection Established
    ├─ onopen fires
    ├─ Update state:
    │  set({ connectionStatus: 'connected' })
    └─ Ready to send requests
        ↓

4️⃣  Create Game Request
    ├─ Call gameService.createGame(...)
    ├─ Generate unique request ID: "req-1"
    ├─ Send message over WebSocket
    ├─ Store pending request with timeout
    └─ Wait for response
        ↓

5️⃣  Receive Response
    ├─ onmessage fires with response
    ├─ Parse JSON: { id: "req-1", type: "response", ... }
    ├─ Find pending request by ID
    ├─ Clear timeout
    ├─ Resolve promise with data
    └─ Remove from pending requests
        ↓

6️⃣  Update Game State
    ├─ Store returns game info
    ├─ Update Zustand state:
    │  set({
    │    mode: 'computer',
    │    gameId: gameInfo.gameId,
    │    playerColor: 'white'
    │  })
    └─ Register event listeners:
       • onMove → handle computer moves
       • onGameEnd → handle game completion
        ↓

7️⃣  Game Active
    └─ WebSocket connection maintained
       • Send move requests
       • Receive move events
       • Auto-reconnect on disconnect
```

### Request/Response Pattern

```typescript
// gameService.ts

async createGame(mode, difficulty, timeControl, playerColor) {
  // Wrap request in promise
  return this.sendRequest('createGame', {
    mode, difficulty, timeControl, playerColor
  });
}

private sendRequest<T>(action: string, data: any): Promise<T> {
  return new Promise((resolve, reject) => {
    const id = `req-${++this.requestId}`;
    const message = { id, type: 'request', action, data };

    // 10 second timeout
    const timeout = setTimeout(() => {
      this.pendingRequests.delete(id);
      reject(new Error('Request timeout'));
    }, 10000);

    // Store promise callbacks
    this.pendingRequests.set(id, { resolve, reject, timeout });

    // Send message
    this.ws.send(JSON.stringify(message));
  });
}

private handleMessage(data: string) {
  const message = JSON.parse(data);

  if (message.type === 'response') {
    // Match response to request
    const pending = this.pendingRequests.get(message.id);
    if (pending) {
      clearTimeout(pending.timeout);
      this.pendingRequests.delete(message.id);

      if (message.success) {
        pending.resolve(message.data);  // Resolve promise
      } else {
        pending.reject(new Error(message.error));
      }
    }
  }
}
```

### Event Handling

```typescript
// Register event listener
gameService.onMove((moveEvent: MoveEvent) => {
  console.log('Computer moved:', moveEvent.moveSAN);
  get().handleComputerMove(moveEvent.moveUCI);
});

// In gameService.ts
onMove(callback: (move: MoveEvent) => void) {
  this.on('move', callback);  // Register
  return () => this.off('move', callback);  // Return unsubscribe
}

private on(event: string, handler: EventHandler) {
  if (!this.eventHandlers.has(event)) {
    this.eventHandlers.set(event, []);
  }
  this.eventHandlers.get(event)!.push(handler);
}

// When event arrives
private handleMessage(data: string) {
  const message = JSON.parse(data);

  if (message.type === 'event') {
    const handlers = this.eventHandlers.get(message.event);
    if (handlers) {
      handlers.forEach(handler => handler(message.data));
    }
  }
}
```

---

## Phase 6: Computer Game Flow

### Complete Computer Game Sequence

```
┌─────────────────────────────────────────────────────────────────┐
│           COMPUTER GAME: HUMAN (WHITE) VS AI (BLACK)            │
└─────────────────────────────────────────────────────────────────┘

1️⃣  User Initiates Game
    ├─ Click "Play vs Computer" button
    ├─ GameSetup dialog opens
    ├─ Select: White pieces, Medium difficulty
    └─ Click "Start Game"
        ↓

2️⃣  Frontend: Start Computer Game
    ├─ Call store.startComputerGame('white', 'medium', userId)
    ├─ Set connectionStatus: 'connecting'
    └─ Clear any errors
        ↓

3️⃣  Connect to WebSocket
    ├─ gameService.connect(userId)
    ├─ ws = new WebSocket("ws://localhost:8080/ws?user_id=...")
    ├─ Wait for onopen event
    └─ Set connectionStatus: 'connected'
        ↓

4️⃣  Create Game Request
    ├─ gameService.createGame(
    │    'human_vs_computer',
    │    'medium',
    │    '5+0',
    │    'white'
    │  )
    ├─ Send WebSocket message:
    │  {
    │    "id": "req-1",
    │    "type": "request",
    │    "action": "createGame",
    │    "data": { mode, difficulty, timeControl, playerColor }
    │  }
    └─ Wait for response
        ↓

5️⃣  Backend Creates Game
    ├─ CreateGameHandler processes request
    ├─ Creates Game instance
    ├─ Creates HumanPlayer (white)
    ├─ Creates ComputerPlayer (black, Stockfish medium)
    ├─ Stores in GameManager
    └─ Sends response:
       {
         "id": "req-1",
         "type": "response",
         "success": true,
         "data": {
           "gameId": "abc-123",
           "status": "active",
           "side": "white",
           "fen": "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
         }
       }
        ↓

6️⃣  Frontend Updates State
    ├─ Receive response
    ├─ Update Zustand store:
    │  set({
    │    mode: 'computer',
    │    opponentType: 'computer',
    │    computerDifficulty: 'medium',
    │    playerColor: 'white',
    │    gameId: 'abc-123',
    │    game: new Chess(),
    │    whitePlayer: 'You',
    │    blackPlayer: 'Computer'
    │  })
    ├─ Register event listeners:
    │  • gameService.onMove(handleComputerMove)
    │  • gameService.onGameEnd(endGame)
    └─ Close GameSetup dialog
        ↓

7️⃣  Human's First Move (e4)
    ├─ User clicks e2 pawn
    ├─ selectSquare('e2') → highlights square, shows valid moves
    ├─ User clicks e4
    ├─ selectSquare('e4') → executes move
    ├─ chess.js validates and applies: game.move({ from: 'e2', to: 'e4' })
    ├─ Update board display
    ├─ Send move to backend:
    │  gameService.makeMove('abc-123', 'e4')
    │  {
    │    "id": "req-2",
    │    "type": "request",
    │    "action": "makeMove",
    │    "data": { "gameId": "abc-123", "move": "e4" }
    │  }
    └─ Set isComputerThinking: true (show loading indicator)
        ↓

8️⃣  Backend Processes Move
    ├─ MakeMoveHandler receives move
    ├─ Validates it's human's turn
    ├─ Applies move to backend game state
    ├─ Broadcasts move event to room
    ├─ Triggers ComputerPlayer.PlayTurn()
    └─ Computer player starts thinking
        ↓

9️⃣  Computer Calculates Move
    ├─ ComputerPlayer sends FEN to Stockfish
    ├─ Stockfish calculates (500ms for medium)
    ├─ Returns best move: "e7e5"
    ├─ Backend applies move
    ├─ Backend broadcasts move event:
    │  {
    │    "type": "event",
    │    "event": "move",
    │    "data": {
    │      "gameId": "abc-123",
    │      "moveSAN": "e5",
    │      "moveUCI": "e7e5",
    │      "fen": "...",
    │      "isGameOver": false
    │    }
    │  }
    └─ Sends response to makeMove request:
       {
         "id": "req-2",
         "type": "response",
         "success": true,
         "data": { "moveNumber": 1, "san": "e5", ... }
       }
        ↓

🔟  Frontend Receives Computer Move
    ├─ onMove event handler fires
    ├─ Call handleComputerMove('e7e5')
    ├─ Parse UCI move:
    │  const move = game.move('e7e5')  // chess.js accepts UCI
    ├─ Update state:
    │  set({
    │    game: {...game},
    │    lastMove: { from: 'e7', to: 'e5' },
    │    isComputerThinking: false
    │  })
    ├─ Board updates with computer's move
    └─ Turn indicator: "White's turn"
        ↓

1️⃣1️⃣  Game Continues
    ├─ Human makes next move
    ├─ Computer responds
    ├─ Repeat until game over
    └─ If checkmate/stalemate:
       Backend sends gameEnd event
       Frontend shows result
```

### State Transitions

```
DISCONNECTED
    ↓ (User clicks "Play vs Computer")
CONNECTING
    ↓ (WebSocket opens)
CONNECTED
    ↓ (Create game request)
ACTIVE GAME
    ├─ Human's turn → Computer thinking → Human's turn
    ├─ (Repeat)
    └─ Game over → Show result
    ↓ (User resets or disconnects)
DISCONNECTED
```

---

## Phase 7: Replay Mode & PGN Import

### PGN Import Flow

```
┌─────────────────────────────────────────────────────────────────┐
│                    PGN IMPORT & REPLAY                          │
└─────────────────────────────────────────────────────────────────┘

1️⃣  User Enters Chess.com Username
    ├─ Type username in GameImporter
    ├─ Click "Load Games"
    └─ Call store.fetchGameArchives('hikaru')
        ↓

2️⃣  Fetch User Archives
    ├─ Call Chess.com API:
    │  GET https://api.chess.com/pub/player/hikaru/games/archives
    ├─ Response: array of monthly archive URLs
    │  [
    │    "https://api.chess.com/pub/player/hikaru/games/2024/11",
    │    "https://api.chess.com/pub/player/hikaru/games/2024/10",
    │    ...
    │  ]
    ├─ Store archives in state
    └─ Auto-fetch latest month's games
        ↓

3️⃣  Fetch Monthly Games
    ├─ Call Chess.com API:
    │  GET https://api.chess.com/pub/player/hikaru/games/2024/11
    ├─ Response: { games: [...] }
    │  Each game includes:
    │  • white: { username, rating }
    │  • black: { username, rating }
    │  • pgn: "full PGN string"
    │  • time_control: "600"
    │  • end_time: 1699123456
    ├─ Reverse array (newest first)
    └─ Store in state: importedGames
        ↓

4️⃣  Display Game List
    ├─ GameList renders
    ├─ Show each game:
    │  • White vs Black
    │  • Result (1-0, 0-1, ½-½)
    │  • Date & time
    │  • Time control
    └─ User clicks game
        ↓

5️⃣  Load PGN for Replay
    ├─ Call store.loadPgnFromGame(game)
    ├─ Extract PGN string from game object
    └─ Call store.loadPgn(pgn)
        ↓

6️⃣  Parse PGN
    ├─ Extract player names:
    │  [White "Hikaru"]
    │  [Black "MagnusCarlsen"]
    ├─ Extract moves (SAN format):
    │  "1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 ..."
    ├─ Clean PGN:
    │  • Remove comments { ... }
    │  • Remove variations ( ... )
    │  • Remove move numbers
    ├─ Parse moves using regex:
    │  const sanMoves = moveText.match(/([NBRQK]?[a-h]?[1-8]?x?[a-h][1-8](?:=[NBRQ])?[+#]?|O-O(?:-O)?)/g)
    │  // Returns: ["e4", "e5", "Nf3", "Nc6", ...]
    └─ Validate moves:
       const tempGame = new Chess()
       moves = sanMoves.map(san => tempGame.move(san))
        ↓

7️⃣  Enter Replay Mode
    ├─ Create fresh Chess instance
    ├─ Update state:
    │  set({
    │    mode: 'replay',
    │    game: new Chess(),
    │    replayMoves: moves,  // Array of Move objects
    │    replayIndex: -1,     // Starting position
    │    whitePlayer: 'Hikaru',
    │    blackPlayer: 'MagnusCarlsen',
    │    selectedSquare: null,
    │    validMoves: [],
    │  })
    ├─ UI updates:
    │  • Show ReplayControls
    │  • Show MoveList
    │  • Show PlayerDisplay (top & bottom)
    │  • Hide "Play vs Computer" button
    └─ Board shows starting position
        ↓

8️⃣  Navigate Through Moves
    ├─ User clicks "Next" button
    ├─ Call store.nextMove()
    ├─ Increment replayIndex: -1 → 0
    ├─ Replay moves from start:
    │  const tempGame = new Chess()
    │  for (let i = 0; i <= 0; i++) {
    │    tempGame.move(replayMoves[i].san)  // Apply "e4"
    │  }
    ├─ Update state:
    │  set({ game: tempGame, replayIndex: 0 })
    └─ Board shows position after 1. e4
        ↓

9️⃣  Autoplay
    ├─ User clicks "Play" button
    ├─ Call store.toggleAutoplay()
    ├─ Start interval:
    │  setInterval(() => {
    │    if (replayIndex < replayMoves.length - 1) {
    │      nextMove()
    │    } else {
    │      stopAutoplay()
    │    }
    │  }, 1000)  // 1 move per second
    ├─ Update state:
    │  set({ isAutoplaying: true })
    └─ Moves advance automatically
        ↓

🔟  Jump to Specific Move
    ├─ User clicks move in MoveList (e.g., move 5)
    ├─ Call store.goToMove(4)  // 0-indexed
    ├─ Stop autoplay
    ├─ Replay from start to move 4:
    │  const tempGame = new Chess()
    │  for (let i = 0; i <= 4; i++) {
    │    tempGame.move(replayMoves[i].san)
    │  }
    ├─ Update state:
    │  set({ game: tempGame, replayIndex: 4 })
    └─ Board shows position after move 5
```

### Replay Controls

```
⏮  First: goToMove(-1)           # Starting position
◀  Previous: prevMove()           # Go back one move
⏯  Play/Pause: toggleAutoplay()   # Auto-advance moves
▶  Next: nextMove()               # Advance one move
⏭  Last: goToMove(moves.length-1) # Final position
```

### Keyboard Shortcuts (Replay Mode)

```
Arrow Left  → Previous move
Arrow Right → Next move
Spacebar    → Play/Pause autoplay
Home        → First move (starting position)
End         → Last move (final position)
```

---

## Complete User Flow Examples

### Example 1: Quick Game vs Computer

```
1. Open app → Board shows starting position
2. Click "Play vs Computer" button
3. Select White pieces, Medium difficulty
4. Click "Start Game"
   → WebSocket connects
   → Game created on backend
   → Dialog closes
5. Click e2 pawn → Square highlights, valid moves appear (e3, e4)
6. Click e4 → Pawn moves, "Computer thinking..." appears
7. Wait 500ms → Computer plays e5, board updates
8. Continue playing until checkmate
9. Backend sends gameEnd event
10. UI shows "Checkmate! Black wins"
```

### Example 2: Import & Replay Game

```
1. Type "hikaru" in username field
2. Click "Load Games"
   → API fetches archives
   → Latest month's games load
3. Browse game list
4. Click game: "Hikaru vs MagnusCarlsen"
   → PGN loads and parses
   → Mode switches to 'replay'
   → Board shows starting position
5. Click "Play" button
   → Autoplay starts
   → Moves advance every 1 second
6. Click move 15 in move list
   → Autoplay stops
   → Jump to move 15
7. Use arrow keys to navigate
8. Press Home to go back to start
```

### Example 3: Pawn Promotion

```
1. Play game until pawn reaches 7th rank
2. Click pawn on e7
3. Click e8 (8th rank)
   → PromotionDialog appears
   → 4 piece options: Q, R, B, N
4. Click Queen
   → Move executes with promotion
   → Queen appears on e8
   → If computer game, move sent to backend
5. Continue playing
```

---

## Key Design Patterns

### 1. **Centralized State Management (Zustand)**
```
Single source of truth for all application state
Components subscribe to specific slices
Actions encapsulate all state mutations
```

### 2. **Service Layer Pattern**
```
GameService handles all WebSocket communication
ChessComAPI handles external API calls
Business logic separated from UI components
```

### 3. **Component Composition**
```
Small, focused components (Square, Piece)
Composed into larger components (GameBoard)
Props flow down, events bubble up
```

### 4. **Request/Response Pattern (WebSocket)**
```
Every request gets unique ID
Promises wrap async requests
Timeout handling for failed requests
```

### 5. **Event-Driven Architecture**
```
WebSocket events trigger state updates
State updates trigger UI re-renders
User interactions dispatch actions
```

### 6. **Optimistic UI Updates**
```
Move executes immediately in chess.js
Board updates instantly (no network delay)
If backend rejects, rollback state
```

### 7. **Separation of Concerns**
```
Types: TypeScript interfaces
State: Zustand store
UI: React components
Logic: chess.js + services
Styling: Tailwind CSS classes
```

---

## Troubleshooting & Common Issues

### Issue 1: WebSocket Won't Connect

```
Error: "WebSocket connection failed"

Solutions:
1. Check backend is running:
   cd backend
   go run cmd/server/main.go

2. Verify WebSocket URL in gameService.ts:
   wsUrl: 'ws://localhost:8080/ws'

3. Check browser console for errors

4. Test connection manually:
   const ws = new WebSocket('ws://localhost:8080/ws?user_id=test');
   ws.onopen = () => console.log('Connected!');
```

### Issue 2: Move Not Updating Board

```
Problem: Click square, nothing happens

Debug:
1. Check console for errors
2. Verify mode is 'live' or 'computer' (not 'replay')
3. Check if it's player's turn (computer game)
4. Verify chess.js validates move:
   const game = new Chess();
   console.log(game.move({ from: 'e2', to: 'e4' }));
```

### Issue 3: Computer Not Responding

```
Problem: Computer doesn't move after human move

Debug:
1. Check WebSocket connection status
2. Verify backend received move (check backend logs)
3. Check for WebSocket event in browser console:
   gameService.onMove(move => console.log('Move event:', move))
4. Verify gameId is set in state
5. Check isComputerThinking flag
```

### Issue 4: PGN Import Fails

```
Error: "Failed to load PGN"

Solutions:
1. Verify username exists on Chess.com
2. Check network tab for API responses
3. Try different month if current month has no games
4. Check PGN format:
   • Must have move text after headers
   • Moves must be valid SAN notation
```

### Issue 5: Board Not Flipping

```
Problem: Board doesn't flip when playing as black

Debug:
1. Check playerColor in state:
   useGameStore.getState().playerColor
2. Verify auto-flip logic in GameBoard.tsx:
   const autoFlip = playerColor === 'black';
   const isFlipped = playerColor ? (isBoardFlipped ? !autoFlip : autoFlip) : isBoardFlipped;
3. Try manual flip button
```

### Issue 6: Hot Reload Issues (Vite)

```
Problem: Changes not reflecting in browser

Solutions:
1. Hard refresh: Cmd+Shift+R (Mac) or Ctrl+Shift+R (Windows)
2. Clear Vite cache:
   rm -rf node_modules/.vite
3. Restart dev server:
   npm run dev
4. Check for TypeScript errors in terminal
```

### Issue 7: State Not Updating

```
Problem: Zustand state changes but UI doesn't update

Debug:
1. Verify you're creating new object (not mutating):
   ❌ game.move(...); set({ game })
   ✅ const newGame = {...game}; newGame.move(...); set({ game: newGame })

2. Check subscription in component:
   const game = useGameStore(state => state.game)  ✅
   const { game } = useGameStore()  ❌ (over-subscription)

3. Use React DevTools to inspect state
```

---

## Summary

```
1️⃣  Bootstrap → Vite loads React app, mounts to DOM
2️⃣  State → Zustand initializes global state (chess.js instance)
3️⃣  Render → Component tree renders based on state
4️⃣  Interaction → User clicks square → selectSquare action
5️⃣  Move → chess.js validates → state updates → UI re-renders
6️⃣  WebSocket → Connect → Create game → Send moves → Receive events
7️⃣  Computer → User moves → Backend AI calculates → Event → Apply move
8️⃣  Replay → Fetch PGN → Parse moves → Navigate with controls
```

---

## Key Features Implemented

✅ Interactive chessboard with move validation
✅ Live games (human vs human locally)
✅ Computer opponent via WebSocket (easy/medium/hard)
✅ PGN import from Chess.com
✅ Replay mode with autoplay
✅ Keyboard shortcuts (arrow keys, spacebar)
✅ Pawn promotion dialog
✅ Board flipping (auto for black, manual toggle)
✅ Responsive design (Tailwind CSS)
✅ WebSocket auto-reconnect
✅ TypeScript type safety
✅ Development debug panel (FEN loader)

---

## Technology Deep Dive

### Vite Build Tool

**Why Vite?**
- ⚡ Instant server start (no bundling in dev)
- 🔥 Hot Module Replacement (HMR)
- 📦 Optimized production builds
- 🎯 Native ESM support

**How it works:**
```
Development:
1. Start dev server (instant)
2. Serve source files as ES modules
3. Transform TypeScript/JSX on-the-fly
4. Hot reload on file changes

Production:
1. Bundle with Rollup
2. Minify & tree-shake
3. Generate optimized static files
4. Output to dist/
```

### Zustand State Management

**Why Zustand?**
- 🎯 Minimal boilerplate (no Provider, no reducers)
- ⚡ Excellent performance (selective subscriptions)
- 🔧 Simple API (set, get)
- 📦 Small bundle size (~1KB)

**Comparison to Redux:**
```
Redux:
• Actions, reducers, store, Provider
• Lots of boilerplate
• ~10KB bundle size

Zustand:
• Single create() call
• Direct state mutations via set()
• ~1KB bundle size
```

### chess.js Library

**What it provides:**
- ♟️ Complete chess rules engine
- ✅ Move validation (legal moves only)
- 📊 Game state (FEN, PGN)
- 🏁 Game over detection (checkmate, stalemate, draw)
- 📝 Move generation
- 🔄 Undo/redo moves

**Example usage:**
```typescript
import { Chess } from 'chess.js';

const game = new Chess();

// Get legal moves for a piece
const moves = game.moves({ square: 'e2', verbose: true });
// [{ from: 'e2', to: 'e3', ... }, { from: 'e2', to: 'e4', ... }]

// Make a move
const move = game.move({ from: 'e2', to: 'e4' });
// { san: 'e4', from: 'e2', to: 'e4', piece: 'p', ... }

// Check game status
console.log(game.isCheckmate());  // false
console.log(game.turn());         // 'b' (black to move)
console.log(game.fen());          // Current position

// Load PGN
game.loadPgn('1. e4 e5 2. Nf3 Nc6');
```

### Tailwind CSS

**Why Tailwind?**
- 🎨 Utility-first CSS
- 📱 Responsive design out-of-the-box
- 🎯 No context switching (styles in JSX)
- 📦 Purge unused CSS in production

**Example:**
```tsx
<button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
  Play vs Computer
</button>

// Responsive:
<div className="w-full md:w-1/2 lg:w-1/3">
  {/* Full width on mobile, half on tablet, third on desktop */}
</div>
```

---

**Happy coding! ♟️⚡**

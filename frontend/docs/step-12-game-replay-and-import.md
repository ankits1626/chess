# Step 12: Game Viewer & Replay Mode

**Goal:** To fetch a user's game from the public chess.com API, load it into the application, and provide UI controls to replay the game move by move.

**This step serves as a crucial bridge to AI integration by forcing us to handle external game data and build a robust replay system.**

---

## 1. High-Level Plan

1.  **UI for Game Import:**
    *   Add a new component, `GameImporter.tsx`, which will contain a text input for a chess.com username and a button to fetch their most recent game.

2.  **API Service for Chess.com:**
    *   Enhance `src/services/api.ts` to include functions for calling the public chess.com API. Specifically, we'll need to fetch a user's monthly archives and then the PGN data from a specific game.
    *   *Reference:* `https://api.chess.com/pub/player/{username}/games/archives`

3.  **State Management Refactor (Zustand):**
    *   As this feature introduces more complex state (is replaying, current move number, full PGN), now is the perfect time to refactor from render props to **Zustand**, as previously discussed.
    *   The store will manage the core `chess.js` instance and add state for the full move list and the current move index.

4.  **Game Replay Controls:**
    *   Create a new component, `ReplayControls.tsx`, with buttons for:
        *   `<<` (First Move)
        *   `<` (Previous Move)
        *   `Play/Pause` (Automatically step through moves)
        *   `>` (Next Move)
        *   `>>` (Last Move)
    *   Integrate these controls, likely below the `GameBoard`.

5.  **Update Game Logic:**
    *   The core logic will need to be updated to handle two modes: "Live Play" and "Replay".
    *   In "Replay" mode, the `onSquareClick` interaction will be disabled, and the board will be updated based on the replay controls, not user moves.

---

## 2. Architectural & Implementation Details

### Zustand Store (`src/store/useGameStore.ts`)

This will replace `GameController.tsx`.

```typescript
import create from 'zustand';
import { Chess } from 'chess.js';

interface GameState {
  game: Chess;
  pgn: string | null;
  moves: any[]; // from game.history({ verbose: true })
  currentMoveIndex: number;
  isPlaying: boolean;

  // Actions
  loadPgn: (pgn: string) => void;
  goToMove: (index: number) => void;
  togglePlay: () => void;
  // ... other actions like selectSquare, makeMove
}

export const useGameStore = create<GameState>((set, get) => ({
  game: new Chess(),
  pgn: null,
  moves: [],
  currentMoveIndex: -1,
  isPlaying: false,

  loadPgn: (pgn) => {
    const newGame = new Chess();
    newGame.loadPgn(pgn);
    set({ 
      game: newGame, 
      pgn, 
      moves: newGame.history({ verbose: true }),
      currentMoveIndex: newGame.history().length - 1
    });
  },

  goToMove: (index) => {
    const { moves } = get();
    if (index < -1 || index >= moves.length) return;

    const tempGame = new Chess();
    for (let i = 0; i <= index; i++) {
      tempGame.move(moves[i]);
    }
    set({ game: tempGame, currentMoveIndex: index });
  },

  togglePlay: () => set(state => ({ isPlaying: !state.isPlaying })),
  
  // ... other actions to be implemented
}));
```

### `GameImporter.tsx` Component

A simple form to fetch a user's game.

```typescript
// src/components/GameImporter.tsx
import { useState } from 'react';
import { useGameStore } from '../store/useGameStore';
import { fetchUserGames, fetchGamePgn } from '../services/chesscomApi';

const GameImporter = () => {
  const [username, setUsername] = useState('hikaru'); // Default for easy testing
  const loadPgn = useGameStore(state => state.loadPgn);

  const handleFetch = async () => {
    try {
      const archives = await fetchUserGames(username);
      const lastGameUrl = archives.pop(); // Get the most recent month's games
      if (lastGameUrl) {
        const pgn = await fetchGamePgn(lastGameUrl);
        loadPgn(pgn);
      }
    } catch (error) {
      console.error("Failed to fetch game:", error);
    }
  };

  return (
    <div className="flex gap-2 p-4 bg-gray-900 rounded-lg">
      <input 
        type="text" 
        value={username} 
        onChange={e => setUsername(e.target.value)}
        className="bg-gray-700 p-2 rounded text-white"
      />
      <button onClick={handleFetch} className="bg-blue-600 p-2 rounded text-white">Fetch Latest Game</button>
    </div>
  );
};
```

### `ReplayControls.tsx` Component

UI for game navigation.

```typescript
// src/components/ReplayControls.tsx
import { useGameStore } from '../store/useGameStore';

const ReplayControls = () => {
  const { goToMove, currentMoveIndex, moves, isPlaying, togglePlay } = useGameStore();

  const handleNext = () => goToMove(currentMoveIndex + 1);
  const handlePrev = () => goToMove(currentMoveIndex - 1);
  const handleFirst = () => goToMove(-1); // -1 represents the initial board state
  const handleLast = () => goToMove(moves.length - 1);

  // useEffect for play/pause functionality would go here

  return (
    <div className="flex justify-center items-center gap-4 p-2 bg-gray-800 rounded-lg">
      <button onClick={handleFirst}>{"«"}</button>
      <button onClick={handlePrev}>{"<"}</button>
      <button onClick={togglePlay}>{isPlaying ? "❚❚" : "▶"}</button>
      <button onClick={handleNext}>{">"}</button>
      <button onClick={handleLast}>{"»"}</button>
    </div>
  );
};
```

### `App.tsx` Refactor

`App.tsx` will no longer use `GameController`. It will be simplified to orchestrate the new components.

```typescript
// App.tsx (Simplified Example)
import { useGameStore } from './store/useGameStore';
import GameBoard from './components/GameBoard';
import GameInfo from './components/GameInfo';
import GameImporter from './components/GameImporter';
import ReplayControls from './components/ReplayControls';

function App() {
  const game = useGameStore(state => state.game);
  const pgn = useGameStore(state => state.pgn);

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-800 text-white p-4">
      <GameImporter />
      <div className="flex flex-col lg:flex-row gap-4 lg:gap-8 items-start mt-4">
        <div>
          <GameBoard game={game} /* ... other props ... */ />
          {pgn && <ReplayControls />}
        </div>
        <GameInfo game={game} /* ... other props ... */ />
      </div>
    </div>
  );
}
```

---

## 3. Phased Implementation Plan

1.  **Phase 1: State Refactor**
    *   Install Zustand (`pnpm add zustand`).
    *   Create `src/store/useGameStore.ts`.
    *   Migrate all state and logic from `GameController.tsx` into the new Zustand store.
    *   Refactor `App.tsx`, `GameBoard.tsx`, and `GameInfo.tsx` to use the `useGameStore` hook instead of receiving props from `GameController`.
    *   Delete `GameController.tsx`.

2.  **Phase 2: API Integration**
    *   Create `src/services/chesscomApi.ts`.
    *   Implement `fetchUserGames` to get the list of monthly archives.
    *   Implement `fetchGamePgn` to fetch the raw PGN data from a game archive URL.
    *   Create and integrate the `GameImporter.tsx` component into `App.tsx`.

3.  **Phase 3: Replay Functionality**
    *   Implement the `goToMove` action in the Zustand store.
    *   Create the `ReplayControls.tsx` component.
    *   Add the play/pause logic using a `useEffect` hook that respects the `isPlaying` state.
    *   Conditionally render `ReplayControls` in `App.tsx` only when a game has been imported (i.e., `pgn` is not null).

---

This plan provides a robust foundation for viewing and replaying games, which will be essential for the AI Coach to analyze and comment on specific moves later.

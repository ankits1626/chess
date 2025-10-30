# Step 12: Game Viewer & Replay Mode - Critical Feedback

## Overall Assessment

**Status: ⚠️ GOOD FOUNDATION - Critical issues need addressing**

The plan provides a solid foundation for game replay and Zustand migration, but there are **critical architectural issues**, **API design problems**, and **missing features** that need to be fixed before implementation.

---

## 🔴 Critical Issues

### 1. Zustand Import - DEPRECATED (Line 44)

**Issue:**
```typescript
import create from 'zustand';
```

**Problem:** `create` import is **deprecated** as of Zustand v4 (2022).

**Correct Import (Zustand v4+):**
```typescript
import { create } from 'zustand';
```

**Why:** Zustand moved to named exports for better tree-shaking and TypeScript support.

---

### 2. Chess.js API Misuse - CRITICAL BUG (Lines 83-86)

**Issue:**
```typescript
goToMove: (index) => {
  const { moves } = get();
  const tempGame = new Chess();
  for (let i = 0; i <= index; i++) {
    tempGame.move(moves[i]);  // ❌ WRONG - moves[i] is a Move object, not a string
  }
  set({ game: tempGame, currentMoveIndex: index });
}
```

**Problem:** `game.move()` expects a move string (SAN) or move object with `{ from, to }`, not the verbose move object directly.

**Correct Implementation:**
```typescript
goToMove: (index) => {
  const { moves } = get();
  if (index < -1 || index >= moves.length) return;

  const tempGame = new Chess();
  if (index >= 0) {
    // Load from PGN or replay moves using SAN notation
    for (let i = 0; i <= index; i++) {
      tempGame.move(moves[i].san);  // ✅ Use .san property
    }
  }
  set({ game: tempGame, currentMoveIndex: index });
}
```

---

### 3. Chess.com API - WRONG APPROACH (Lines 112-116)

**Issue:**
```typescript
const archives = await fetchUserGames(username);
const lastGameUrl = archives.pop(); // Get the most recent month's games
if (lastGameUrl) {
  const pgn = await fetchGamePgn(lastGameUrl);
  loadPgn(pgn);
}
```

**Problems:**
1. **Wrong data structure:** `archives` is a list of URLs to monthly archives, not individual games
2. **Wrong URL:** You're passing a monthly archive URL to `fetchGamePgn`, not a game URL
3. **No game selection:** Gets all games from a month, not the latest single game

**Correct API Flow:**
```typescript
// Step 1: Get archives list
// https://api.chess.com/pub/player/{username}/games/archives
// Returns: { "archives": ["https://api.chess.com/pub/player/{username}/games/2025/10", ...] }

// Step 2: Fetch specific month's games
// https://api.chess.com/pub/player/{username}/games/2025/10
// Returns: { "games": [{ "pgn": "...", "url": "...", ... }, ...] }

// Step 3: Extract PGN from the first/last game
const archives = await fetchUserGames(username);
const lastArchiveUrl = archives.archives[archives.archives.length - 1];
const monthGames = await fetchMonthGames(lastArchiveUrl);
const latestGame = monthGames.games[monthGames.games.length - 1];
const pgn = latestGame.pgn;
loadPgn(pgn);
```

---

### 4. Store Type Definition - INCOMPLETE (Lines 47-59)

**Issue:**
```typescript
interface GameState {
  game: Chess;
  pgn: string | null;
  moves: any[]; // ❌ any[] is a code smell
  currentMoveIndex: number;
  isPlaying: boolean;

  // Missing critical state from current GameController:
  // - selectedSquare
  // - validMoves
  // - pendingMove
  // - lastMove
  // - debugActions
}
```

**Problem:** Plan doesn't account for migrating existing live play state.

**Complete State Interface:**
```typescript
import { Move } from 'chess.js';

interface GameState {
  // Chess game state
  game: Chess;

  // Live play state (from GameController)
  selectedSquare: Square | null;
  validMoves: Square[];
  pendingMove: PendingPromotion | null;
  lastMove: LastMove;

  // Replay state (new)
  mode: 'live' | 'replay';
  pgn: string | null;
  moves: Move[];  // ✅ Properly typed
  currentMoveIndex: number;
  isPlaying: boolean;

  // Actions
  // Live play actions
  selectSquare: (square: Square) => void;
  handlePromotion: (piece: PieceType) => void;
  resetGame: () => void;

  // Replay actions
  loadPgn: (pgn: string) => void;
  goToMove: (index: number) => void;
  togglePlay: () => void;
  nextMove: () => void;
  prevMove: () => void;

  // Debug actions (dev-only)
  debugActions?: DebugActions;
}
```

---

### 5. Mode Switching - NOT ADDRESSED (Line 32-33)

**Issue:**
> In "Replay" mode, the `onSquareClick` interaction will be disabled

**Problem:** Plan doesn't explain:
1. How to switch between modes?
2. Can user exit replay mode?
3. What happens to live game when loading replay?
4. Should there be a "New Game" vs "Import Game" distinction?

**Better Design:**
```typescript
// Add mode state
mode: 'live' | 'replay';

// Add mode switching
switchToLiveMode: () => {
  set({
    mode: 'live',
    pgn: null,
    moves: [],
    currentMoveIndex: -1,
    game: new Chess()
  });
},

switchToReplayMode: (pgn: string) => {
  const game = new Chess();
  game.loadPgn(pgn);
  set({
    mode: 'replay',
    pgn,
    moves: game.history({ verbose: true }),
    currentMoveIndex: game.history().length - 1,
    game
  });
},

// In GameBoard:
const mode = useGameStore(state => state.mode);
const onSquareClick = mode === 'live' ? selectSquare : undefined;
```

---

### 6. Auto-Play Logic - MISSING (Line 153)

**Issue:**
```typescript
// useEffect for play/pause functionality would go here
```

**Problem:** Plan doesn't provide the auto-play implementation.

**Implementation:**
```typescript
// In ReplayControls.tsx
import { useEffect } from 'react';

const ReplayControls = () => {
  const { isPlaying, currentMoveIndex, moves, goToMove } = useGameStore();

  useEffect(() => {
    if (!isPlaying) return;

    // Auto-advance every 1 second
    const interval = setInterval(() => {
      if (currentMoveIndex < moves.length - 1) {
        goToMove(currentMoveIndex + 1);
      } else {
        // Reached end, stop playing
        useGameStore.getState().togglePlay();
      }
    }, 1000);

    return () => clearInterval(interval);
  }, [isPlaying, currentMoveIndex, moves.length, goToMove]);

  // ... rest of component
};
```

---

## ⚠️ Major Design Issues

### 7. PGN Loading - INCOMPLETE (Lines 68-76)

**Issue:**
```typescript
loadPgn: (pgn) => {
  const newGame = new Chess();
  newGame.loadPgn(pgn);
  set({
    game: newGame,
    pgn,
    moves: newGame.history({ verbose: true }),
    currentMoveIndex: newGame.history().length - 1  // Shows final position
  });
}
```

**Problems:**
1. Always shows final position (not useful for replay)
2. Doesn't clear live game state
3. Doesn't handle PGN parsing errors
4. Doesn't extract game metadata (players, result, date)

**Better Implementation:**
```typescript
loadPgn: (pgn: string) => {
  try {
    const newGame = new Chess();

    // Parse PGN (may throw if invalid)
    const success = newGame.loadPgn(pgn);
    if (!success) {
      throw new Error('Invalid PGN');
    }

    // Extract metadata
    const header = newGame.header();

    // Get move history
    const moves = newGame.history({ verbose: true });

    // Reset to starting position for replay
    const replayGame = new Chess();

    set({
      mode: 'replay',
      game: replayGame,  // ✅ Start at beginning
      pgn,
      moves,
      currentMoveIndex: -1,  // ✅ Before first move
      isPlaying: false,

      // Clear live game state
      selectedSquare: null,
      validMoves: [],
      pendingMove: null,
      lastMove: null,

      // Store metadata
      metadata: {
        white: header.White || 'Unknown',
        black: header.Black || 'Unknown',
        result: header.Result || '*',
        date: header.Date || '',
      }
    });
  } catch (error) {
    console.error('Failed to load PGN:', error);
    // Could set an error state here
  }
}
```

---

### 8. Chess.com API Service - INCOMPLETE (Lines 100-135)

**Issue:** Plan shows component but not the actual API service implementation.

**Complete Implementation Needed:**

```typescript
// src/services/chesscomApi.ts

interface ChessComArchivesResponse {
  archives: string[];  // List of archive URLs
}

interface ChessComGamesResponse {
  games: Array<{
    url: string;
    pgn: string;
    time_control: string;
    end_time: number;
    rated: boolean;
    // ... other fields
  }>;
}

export const fetchUserGames = async (username: string): Promise<ChessComArchivesResponse> => {
  const response = await fetch(
    `https://api.chess.com/pub/player/${username}/games/archives`
  );

  if (!response.ok) {
    throw new Error(`Failed to fetch archives: ${response.statusText}`);
  }

  return response.json();
};

export const fetchMonthGames = async (archiveUrl: string): Promise<ChessComGamesResponse> => {
  const response = await fetch(archiveUrl);

  if (!response.ok) {
    throw new Error(`Failed to fetch games: ${response.statusText}`);
  }

  return response.json();
};

export const fetchLatestGame = async (username: string): Promise<string> => {
  // Get archives list
  const archives = await fetchUserGames(username);

  if (archives.archives.length === 0) {
    throw new Error('No games found for user');
  }

  // Get most recent month
  const latestArchiveUrl = archives.archives[archives.archives.length - 1];
  const monthGames = await fetchMonthGames(latestArchiveUrl);

  if (monthGames.games.length === 0) {
    throw new Error('No games in latest archive');
  }

  // Get most recent game
  const latestGame = monthGames.games[monthGames.games.length - 1];
  return latestGame.pgn;
};
```

---

### 9. GameImporter UX - POOR (Lines 106-134)

**Issues:**
1. No loading state (API calls take time)
2. No error handling UI
3. No success feedback
4. No game preview before loading

**Better UX:**
```typescript
const GameImporter = () => {
  const [username, setUsername] = useState('hikaru');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [gamePreview, setGamePreview] = useState<GamePreview | null>(null);

  const loadPgn = useGameStore(state => state.loadPgn);

  const handleFetch = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const pgn = await fetchLatestGame(username);

      // Parse for preview
      const game = new Chess();
      game.loadPgn(pgn);
      const header = game.header();

      setGamePreview({
        white: header.White,
        black: header.Black,
        result: header.Result,
        date: header.Date,
        pgn
      });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch game');
    } finally {
      setIsLoading(false);
    }
  };

  const handleLoad = () => {
    if (gamePreview) {
      loadPgn(gamePreview.pgn);
      setGamePreview(null);
    }
  };

  return (
    <div className="p-4 bg-gray-900 rounded-lg space-y-4">
      <div className="flex gap-2">
        <input
          type="text"
          value={username}
          onChange={e => setUsername(e.target.value)}
          placeholder="Chess.com username"
          className="bg-gray-700 p-2 rounded text-white flex-grow"
          disabled={isLoading}
        />
        <button
          onClick={handleFetch}
          disabled={isLoading}
          className="bg-blue-600 hover:bg-blue-500 disabled:bg-gray-600 p-2 rounded text-white"
        >
          {isLoading ? 'Fetching...' : 'Fetch Latest Game'}
        </button>
      </div>

      {error && (
        <div className="bg-red-900/50 border border-red-500 p-3 rounded text-red-200">
          {error}
        </div>
      )}

      {gamePreview && (
        <div className="bg-gray-800 p-4 rounded space-y-2">
          <h3 className="font-bold">Game Preview</h3>
          <p><span className="text-gray-400">White:</span> {gamePreview.white}</p>
          <p><span className="text-gray-400">Black:</span> {gamePreview.black}</p>
          <p><span className="text-gray-400">Result:</span> {gamePreview.result}</p>
          <p><span className="text-gray-400">Date:</span> {gamePreview.date}</p>
          <button
            onClick={handleLoad}
            className="bg-green-600 hover:bg-green-500 p-2 rounded text-white w-full"
          >
            Load Game
          </button>
        </div>
      )}
    </div>
  );
};
```

---

### 10. ReplayControls - MISSING FEATURES (Lines 140-164)

**Missing:**
1. Move number display
2. Keyboard shortcuts (arrow keys)
3. Speed control (1x, 2x, 0.5x)
4. Progress bar/slider
5. Move annotations display

**Enhanced Implementation:**
```typescript
const ReplayControls = () => {
  const {
    isPlaying,
    currentMoveIndex,
    moves,
    goToMove,
    togglePlay
  } = useGameStore();

  // Keyboard shortcuts
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'ArrowLeft') handlePrev();
      if (e.key === 'ArrowRight') handleNext();
      if (e.key === ' ') { e.preventDefault(); togglePlay(); }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [currentMoveIndex, moves.length]);

  const handleNext = () => goToMove(Math.min(currentMoveIndex + 1, moves.length - 1));
  const handlePrev = () => goToMove(Math.max(currentMoveIndex - 1, -1));
  const handleFirst = () => goToMove(-1);
  const handleLast = () => goToMove(moves.length - 1);

  // Auto-play logic
  useEffect(() => {
    if (!isPlaying) return;

    const interval = setInterval(() => {
      if (currentMoveIndex < moves.length - 1) {
        goToMove(currentMoveIndex + 1);
      } else {
        togglePlay(); // Stop at end
      }
    }, 1000);

    return () => clearInterval(interval);
  }, [isPlaying, currentMoveIndex, moves.length]);

  return (
    <div className="bg-gray-800 rounded-lg p-4 space-y-3">
      {/* Move counter */}
      <div className="text-center text-gray-300">
        Move {currentMoveIndex + 1} of {moves.length}
        {currentMoveIndex >= 0 && (
          <span className="ml-2 text-white font-mono">
            {moves[currentMoveIndex].san}
          </span>
        )}
      </div>

      {/* Progress slider */}
      <input
        type="range"
        min="-1"
        max={moves.length - 1}
        value={currentMoveIndex}
        onChange={(e) => goToMove(parseInt(e.target.value))}
        className="w-full"
      />

      {/* Control buttons */}
      <div className="flex justify-center items-center gap-2">
        <button
          onClick={handleFirst}
          disabled={currentMoveIndex === -1}
          className="px-3 py-2 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-800 rounded"
        >
          {"⏮"}
        </button>
        <button
          onClick={handlePrev}
          disabled={currentMoveIndex === -1}
          className="px-3 py-2 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-800 rounded"
        >
          {"◀"}
        </button>
        <button
          onClick={togglePlay}
          className="px-4 py-2 bg-blue-600 hover:bg-blue-500 rounded"
        >
          {isPlaying ? "❚❚" : "▶"}
        </button>
        <button
          onClick={handleNext}
          disabled={currentMoveIndex === moves.length - 1}
          className="px-3 py-2 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-800 rounded"
        >
          {"▶"}
        </button>
        <button
          onClick={handleLast}
          disabled={currentMoveIndex === moves.length - 1}
          className="px-3 py-2 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-800 rounded"
        >
          {"⏭"}
        </button>
      </div>

      {/* Keyboard hint */}
      <div className="text-center text-xs text-gray-400">
        Arrow keys: ←/→  |  Space: Play/Pause
      </div>
    </div>
  );
};
```

---

## Missing Features

### 11. Game Metadata Display

**Should show:**
- Player names (White vs Black)
- Game result (1-0, 0-1, 1/2-1/2)
- Date played
- Time control
- Opening name (if available)

### 12. Move List Navigation

**Clickable move list:**
- Click any move in history to jump to that position
- Current move highlighted
- Move annotations displayed

### 13. Export Functionality

**Should support:**
- Copy PGN to clipboard
- Download PGN file
- Share game URL

### 14. Multiple Game Selection

**Current plan only fetches one game:**
- Should show list of recent games
- User selects which to view
- Search/filter by opponent, date, result

---

## Implementation Order (Corrected)

### Phase 1: Zustand Migration (Week 1)

1. ✅ Install Zustand: `pnpm add zustand`
2. ✅ Create complete store interface
3. ✅ Migrate live game state (selectedSquare, validMoves, etc.)
4. ✅ Migrate live game actions (selectSquare, makeMove, etc.)
5. ✅ Keep debugActions separation
6. ✅ Update all components to use store
7. ✅ Test thoroughly - no regressions
8. ✅ Delete GameController.tsx

**Critical:** Don't add replay features yet. Just migrate existing functionality.

### Phase 2: API Integration (Week 2)

1. ✅ Create `src/services/chesscomApi.ts`
2. ✅ Implement `fetchUserGames`, `fetchMonthGames`, `fetchLatestGame`
3. ✅ Add error handling and types
4. ✅ Test API calls manually
5. ✅ Create GameImporter component with loading/error states
6. ✅ Add game preview before loading

### Phase 3: Replay Core (Week 3)

1. ✅ Add replay state to store (mode, pgn, moves, currentMoveIndex)
2. ✅ Implement `loadPgn` action with proper error handling
3. ✅ Implement `goToMove` action (fixed version)
4. ✅ Add mode switching logic
5. ✅ Disable board interaction in replay mode
6. ✅ Test PGN loading and navigation

### Phase 4: Replay UI (Week 4)

1. ✅ Create ReplayControls component
2. ✅ Add auto-play logic
3. ✅ Add keyboard shortcuts
4. ✅ Add progress slider
5. ✅ Add move counter
6. ✅ Style and polish

### Phase 5: Enhancements (Week 5)

1. ✅ Add game metadata display
2. ✅ Add clickable move list
3. ✅ Add export functionality
4. ✅ Add multiple game selection
5. ✅ Add speed control

---

## Testing Strategy

### Unit Tests (Recommended)

```typescript
// __tests__/store/useGameStore.test.ts
describe('useGameStore', () => {
  test('loadPgn should parse valid PGN', () => {
    const { loadPgn } = useGameStore.getState();
    loadPgn('1. e4 e5 2. Nf3');

    const { moves, mode } = useGameStore.getState();
    expect(moves).toHaveLength(3);
    expect(mode).toBe('replay');
  });

  test('goToMove should navigate correctly', () => {
    // Load a game first
    const { loadPgn, goToMove } = useGameStore.getState();
    loadPgn('1. e4 e5 2. Nf3 Nc6');

    // Go to move 2 (e5)
    goToMove(1);

    const { currentMoveIndex, game } = useGameStore.getState();
    expect(currentMoveIndex).toBe(1);
    expect(game.fen()).toContain('e5'); // Verify board state
  });
});
```

### Integration Tests

1. Fetch real game from chess.com API
2. Load into application
3. Navigate through moves
4. Verify board updates correctly

---

## Summary of Critical Fixes Needed

| Issue | Severity | Fix Required |
|-------|----------|-------------|
| Zustand import | 🔴 Critical | Change to `import { create }` |
| goToMove API misuse | 🔴 Critical | Use `moves[i].san` not `moves[i]` |
| Chess.com API flow | 🔴 Critical | Fix archive URL handling |
| Incomplete state | 🔴 Critical | Add all GameController state |
| Mode switching | ⚠️ High | Add mode management logic |
| Auto-play logic | ⚠️ High | Implement useEffect |
| PGN loading | ⚠️ High | Start at beginning, handle errors |
| API service | ⚠️ High | Complete implementation |
| UX feedback | ⚠️ Medium | Add loading/error states |
| Enhanced controls | ⚠️ Medium | Add slider, keyboard, etc. |

---

## Recommendation

**Do NOT start implementation until these issues are fixed:**

1. ✅ Update Zustand import syntax
2. ✅ Fix `goToMove` to use `.san` property
3. ✅ Correct Chess.com API flow
4. ✅ Complete state interface with all live game state
5. ✅ Add mode switching logic
6. ✅ Implement auto-play properly
7. ✅ Add comprehensive error handling
8. ✅ Improve GameImporter UX

**Estimated Time:** 3-4 weeks for complete, polished implementation.

**Alternative:** Skip Step 12 and go straight to AI coach with manual PGN input. Add game import later.

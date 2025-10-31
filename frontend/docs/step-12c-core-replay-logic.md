# Step 12c: Core Replay Logic

**Goal:** To implement the foundational logic within the Zustand store that allows the application to load a PGN and programmatically navigate through its moves, setting the board to the state after each move.

**This step corresponds to Phase 3 of the feedback on the original Step 12 plan.**

---

## 1. Implementation Plan

### 1. Update Zustand Store (`src/store/useGameStore.ts`)

We need to add actions to the `useGameStore` that allow us to navigate through the `replayMoves` array and update the `game` object accordingly.

**Key additions to `GameState` interface:**
*   `goToMove: (index: number) => void;`
*   `nextMove: () => void;`
*   `prevMove: () => void;`
*   `firstMove: () => void;`
*   `lastMove: () => void;`

**Example Implementation:**

```typescript
// src/store/useGameStore.ts

// ... existing imports and types ...

interface GameState {
  // ... existing state properties ...

  // Replay state
  mode: 'live' | 'replay';
  replayMoves: Move[]; // `Move` type from chess.js
  replayIndex: number; // -1 for initial position, 0 for first move, etc.

  // ... existing actions ...

  // New Replay Actions
  goToMove: (index: number) => void;
  nextMove: () => void;
  prevMove: () => void;
  firstMove: () => void;
  lastMove: () => void;
}

export const useGameStore = create<GameState>((set, get) => ({
  // ... existing initial state ...

  // New Replay State
  mode: 'live',
  replayMoves: [],
  replayIndex: -1, // -1 means initial board state

  // ... existing actions ...

  // Replay Actions
  goToMove: (index: number) => {
    const { replayMoves } = get();

    // Ensure index is within bounds
    if (index < -1 || index >= replayMoves.length) {
      console.warn(`Attempted to go to invalid move index: ${index}`);
      return;
    }

    const tempGame = new Chess();
    // Replay moves up to the specified index
    for (let i = 0; i <= index; i++) {
      tempGame.move(replayMoves[i].san); // Use .san property for moves
    }

    set({
      game: Object.assign(Object.create(Object.getPrototypeOf(tempGame)), tempGame),
      replayIndex: index,
      // Clear live game state when navigating in replay mode
      selectedSquare: null,
      validMoves: [],
      pendingMove: null,
      lastMove: null,
    });
  },

  nextMove: () => {
    const { replayIndex, replayMoves, goToMove } = get();
    if (replayIndex < replayMoves.length - 1) {
      goToMove(replayIndex + 1);
    }
  },

  prevMove: () => {
    const { replayIndex, goToMove } = get();
    if (replayIndex > -1) {
      goToMove(replayIndex - 1);
    }
  },

  firstMove: () => {
    get().goToMove(-1); // Go to initial board state
  },

  lastMove: () => {
    const { replayMoves, goToMove } = get();
    goToMove(replayMoves.length - 1);
  },

  // ... existing debugActions ...
}));
```

### 2. Update `loadPgn` Action in Zustand Store

The `loadPgn` action needs to be updated to properly set the `mode` and `replayMoves` state, and to ensure the board starts at the initial position for replay.

```typescript
// src/store/useGameStore.ts

// ... existing loadPgn action ...

  loadPgn: (pgn: string) => {
    try {
      // ... existing PGN cleaning ...

      const newGame = new Chess();
      const success = newGame.loadPgn(cleanedPgn);

      if (!success) {
        console.error('Failed to load PGN. Invalid format:', pgn);
        throw new Error('Invalid PGN format');
      }

      // Get all moves in verbose format for replay
      const moves = newGame.history({ verbose: true }) as Move[];

      // Reset to starting position for replay
      const replayGame = new Chess(); // This will be the initial board

      set({
        mode: 'replay', // Set mode to replay
        game: Object.assign(Object.create(Object.getPrototypeOf(replayGame)), replayGame), // Initial board state
        replayMoves: moves,
        replayIndex: -1, // Start before the first move

        // Clear live game state
        selectedSquare: null,
        validMoves: [],
        pendingMove: null,
        lastMove: null
      });
    } catch (error) {
      console.error('Failed to load PGN:', error);
      throw error; // Re-throw for component to handle
    }
  },
```

### 3. Disable Live Play Interactions in Replay Mode

Ensure that when in `replay` mode, users cannot make live moves or select pieces.

*   **`selectSquare` action:** Already has `if (mode === 'replay') return;`
*   **`handlePromotion` action:** Should also be guarded by `mode`.

```typescript
// src/store/useGameStore.ts

// ... in handlePromotion action ...
  handlePromotion: (piece: PieceType) => {
    const { game, pendingMove, mode } = get();
    if (mode === 'replay') return; // Disable in replay mode
    if (!pendingMove) return;
    // ... rest of logic ...
  },
```

---

## 2. Testing Checklist

1.  **Load PGN:**
    *   ✅ Import a game using the `GameImporter`.
    *   ✅ Verify the board resets to the initial position (before the first move).
    *   ✅ Verify `mode` is set to `'replay'`.
    *   ✅ Verify `replayMoves` is populated with the game's moves.
    *   ✅ Verify `replayIndex` is `-1`.

2.  **`nextMove` Functionality:**
    *   ✅ Call `nextMove()` repeatedly.
    *   ✅ Verify the board advances one move at a time.
    *   ✅ Verify `replayIndex` increments correctly.
    *   ✅ Verify `nextMove()` stops at the last move.

3.  **`prevMove` Functionality:**
    *   ✅ Advance several moves, then call `prevMove()` repeatedly.
    *   ✅ Verify the board goes back one move at a time.
    *   ✅ Verify `replayIndex` decrements correctly.
    *   ✅ Verify `prevMove()` stops at the initial board state (`replayIndex: -1`).

4.  **`goToMove` Functionality:**
    *   ✅ Test `goToMove(0)`: Should show the board after White's first move.
    *   ✅ Test `goToMove(replayMoves.length - 1)`: Should show the final board position.
    *   ✅ Test `goToMove(-1)`: Should show the initial board position.
    *   ✅ Test `goToMove` with an invalid index: Should log a warning and not change state.

5.  **`firstMove` and `lastMove` Functionality:**
    *   ✅ Test `firstMove()`: Should go to the initial board state.
    *   ✅ Test `lastMove()`: Should go to the final board position.

6.  **Mode Switching:**
    *   ✅ In replay mode, verify that clicking on squares does NOT select pieces or make moves.
    *   ✅ Verify that the promotion dialog does NOT appear in replay mode.

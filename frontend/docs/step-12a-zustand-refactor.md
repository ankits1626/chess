# Step 12a: State Management Refactor to Zustand

**Goal:** To replace the `GameController.tsx` render props pattern with a global Zustand store. This will simplify state management, improve performance, and prepare the application for more complex features like game replay and AI chat.

**This step corresponds to Phase 1 of the feedback on the original Step 12 plan.**

---

## 1. Implementation Plan

### 1. Install Zustand

Add Zustand as a project dependency.

```bash
pnpm add zustand
```

### 2. Create the Zustand Store (`useGameStore.ts`)

Create a new file at `src/store/useGameStore.ts`. This file will define the complete state of the application and all the actions that can modify it. It will contain all the logic currently in `GameController.tsx`.

**Key elements of the store:**
*   **State:** `game`, `selectedSquare`, `validMoves`, `pendingMove`, `lastMove`, and the `debugActions`.
*   **Actions:** `selectSquare`, `handlePromotion`, and `resetGame`.

**Example Implementation:**

```typescript
// src/store/useGameStore.ts
import { create } from 'zustand';
import { Chess } from 'chess.js';
import type { Square, PieceType, PieceColor } from '../types/chess';

// ... (Type definitions for LastMove, PendingPromotion, DebugActions) ...

interface GameState {
  game: Chess;
  selectedSquare: Square | null;
  validMoves: Square[];
  pendingMove: any; // Replace with PendingPromotion type
  lastMove: any; // Replace with LastMove type
  debugActions?: any; // Replace with DebugActions type

  selectSquare: (square: Square) => void;
  handlePromotion: (piece: PieceType) => void;
  resetGame: () => void;
}

export const useGameStore = create<GameState>((set, get) => ({
  // Initial State
  game: new Chess(),
  selectedSquare: null,
  validMoves: [],
  pendingMove: null,
  lastMove: null,

  // Actions
  selectSquare: (square) => {
    // ... (All logic from selectSquare in GameController.tsx) ...
    // Instead of calling React's set... functions, call Zustand's set()
    // e.g., set({ selectedSquare: square });
  },

  handlePromotion: (piece) => {
    // ... (All logic from handlePromotion in GameController.tsx) ...
    // e.g., set({ game: new Chess(get().game.fen()), ... });
  },

  resetGame: () => {
    // ... (All logic from resetGame in GameController.tsx) ...
    // e.g., set({ game: new Chess(), ... });
  },

  // Debug Actions (conditionally added)
  debugActions: import.meta.env.DEV ? { /* ... debug logic ... */ } : undefined,
}));
```

### 3. Refactor Components to Use the Store

Modify all components that previously received props from `GameController` to use the `useGameStore` hook instead.

*   **`App.tsx`:**
    *   Remove the `<GameController>` wrapper.
    *   The component will now be much simpler, primarily orchestrating the main components.

*   **`GameBoard.tsx`:**
    *   Will call `useGameStore` to get `game`, `selectedSquare`, `validMoves`, `lastMove`, and `selectSquare`.
    *   Props will no longer be passed down from `App.tsx`.

*   **`GameInfo.tsx`:**
    *   Will call `useGameStore` to get `game` and `resetGame`.

*   **`DebugPanel.tsx` (and related components):**
    *   Will call `useGameStore` to get `debugActions` and `game.fen()`.

### 4. Delete `GameController.tsx`

Once all components have been refactored and are using the Zustand store, the `GameController.tsx` file is no longer needed and can be safely deleted.

---

## 2. Detailed Steps

1.  **Run `pnpm add zustand`**.
2.  **Create `src/store/useGameStore.ts`** and populate it with the full state interface and the logic migrated from `GameController.tsx`.
3.  **Refactor `GameBoard.tsx`** to pull all its data and actions from the `useGameStore` hook.
4.  **Refactor `GameInfo.tsx`** to pull `game` and `resetGame` from the store.
5.  **Refactor `DebugPanel.tsx`** to pull `debugActions` from the store.
6.  **Refactor `App.tsx`** to remove the `GameController` and simplify its structure.
7.  **Delete `frontend/app/src/components/GameController.tsx`**.
8.  **Thoroughly test** the application to ensure that all functionality (piece selection, moving, promotion, reset, debug tools) works exactly as it did before the refactor. Check for any regressions.

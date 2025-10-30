# Step 12a: Zustand Refactor - Implementation Guide

## Dependency Verification ✅

**Current Installation:**
```bash
zustand 5.0.8  # Latest version as of October 2025
```

✅ **Correct version** - No upgrade needed
✅ **Import syntax in plan is correct:** `import { create } from 'zustand'`

---

## Plan Assessment

**Status: ⚠️ GOOD OUTLINE - Needs detailed implementation**

The plan provides correct structure but lacks critical implementation details. This document provides the complete, production-ready code.

---

## Complete Implementation

### Step 1: Install Zustand ✅ ALREADY DONE

```bash
pnpm add zustand
# Currently installed: zustand@5.0.8
```

---

### Step 2: Create Complete Zustand Store

**Create file:** `src/store/useGameStore.ts`

```typescript
import { create } from 'zustand';
import { Chess } from 'chess.js';
import type { Square, PieceType, PieceColor } from '../types/chess';

// Type definitions (export for use in other files)
export type LastMove = {
  from: Square;
  to: Square;
} | null;

export type PendingPromotion = {
  from: Square;
  to: Square;
  color: PieceColor;
} | null;

export interface DebugActions {
  loadFen: (fen: string) => void;
  resetGame: () => void;
}

// Main store interface
interface GameState {
  // Game state
  game: Chess;
  selectedSquare: Square | null;
  validMoves: Square[];
  pendingMove: PendingPromotion;
  lastMove: LastMove;

  // Actions
  selectSquare: (square: Square) => void;
  handlePromotion: (piece: PieceType) => void;
  resetGame: () => void;

  // Debug actions (dev-only)
  debugActions?: DebugActions;
}

export const useGameStore = create<GameState>((set, get) => ({
  // Initial State
  game: new Chess(),
  selectedSquare: null,
  validMoves: [],
  pendingMove: null,
  lastMove: null,

  // Actions
  selectSquare: (square: Square) => {
    const { game, selectedSquare, pendingMove } = get();

    // If a promotion is pending, don't allow other moves
    if (pendingMove) return;

    // If no square is selected, select this square (if it has a piece)
    if (!selectedSquare) {
      const piece = game.get(square);
      if (piece && piece.color === game.turn()) {
        const moves = game.moves({ square, verbose: true });
        set({
          selectedSquare: square,
          validMoves: moves.map(move => move.to)
        });
      }
      return;
    }

    // If clicking the same square, deselect it
    if (selectedSquare === square) {
      set({
        selectedSquare: null,
        validMoves: []
      });
      return;
    }

    // Try to make a move
    try {
      const piece = game.get(selectedSquare);

      // Check if this would be a promotion move
      if (piece?.type === 'p' && (square.endsWith('1') || square.endsWith('8'))) {
        const moves = game.moves({ square: selectedSquare, verbose: true });
        const isValidMove = moves.some(m => m.to === square);

        if (isValidMove) {
          // Valid promotion - show dialog
          set({
            pendingMove: {
              from: selectedSquare,
              to: square,
              color: piece.color
            },
            selectedSquare: null,
            validMoves: []
          });
          return;
        }
      }

      // Not a promotion, proceed with normal move
      const move = game.move({
        from: selectedSquare,
        to: square
      });

      if (move) {
        // Move was successful, update state
        set({
          game: new Chess(game.fen()),
          lastMove: { from: selectedSquare, to: square },
          selectedSquare: null,
          validMoves: []
        });
      } else {
        // Invalid move, check if clicking another piece of the same color
        const newPiece = game.get(square);
        if (newPiece && newPiece.color === game.turn()) {
          const moves = game.moves({ square, verbose: true });
          set({
            selectedSquare: square,
            validMoves: moves.map(move => move.to)
          });
        } else {
          set({
            selectedSquare: null,
            validMoves: []
          });
        }
      }
    } catch {
      set({
        selectedSquare: null,
        validMoves: []
      });
    }
  },

  handlePromotion: (piece: PieceType) => {
    const { game, pendingMove } = get();
    if (!pendingMove) return;

    game.move({
      from: pendingMove.from,
      to: pendingMove.to,
      promotion: piece
    });

    set({
      game: new Chess(game.fen()),
      lastMove: { from: pendingMove.from, to: pendingMove.to },
      pendingMove: null
    });
  },

  resetGame: () => {
    set({
      game: new Chess(),
      selectedSquare: null,
      validMoves: [],
      pendingMove: null,
      lastMove: null
    });
  },

  // Debug actions (conditionally added in dev mode)
  debugActions: import.meta.env.DEV ? {
    loadFen: (fen: string) => {
      try {
        const newGame = new Chess(fen);
        set({
          game: newGame,
          selectedSquare: null,
          validMoves: [],
          pendingMove: null,
          lastMove: null
        });
      } catch (error) {
        console.error('Invalid FEN:', error);
      }
    },
    resetGame: () => {
      get().resetGame();
    }
  } : undefined
}));
```

---

### Step 3: Refactor App.tsx

**Before (with GameController):**
```typescript
import { lazy, Suspense } from 'react';
import GameBoard from './components/GameBoard';
import GameInfo from './components/GameInfo';
import GameController from './components/GameController';
import PromotionDialog from './components/PromotionDialog';
import { useDebugPanel } from './hooks/useDebugPanel';

const DebugPanel = import.meta.env.DEV
  ? lazy(() => import('./components/debug/DebugPanel'))
  : null;

function App() {
  const { isPanelVisible, setIsPanelVisible } = useDebugPanel();

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-800 text-white p-8">
      <GameController>
        {(game, selectedSquare, validMoves, selectSquare, pendingMove, handlePromotion, lastMove, resetGame, debugActions) => (
          <>
            <div className="flex flex-col lg:flex-row gap-4 lg:gap-8 items-center">
              <GameBoard
                game={game}
                selectedSquare={selectedSquare}
                validMoves={validMoves}
                lastMove={lastMove}
                onSquareClick={selectSquare}
              />
              <GameInfo game={game} onNewGame={resetGame} />
            </div>

            {pendingMove && (
              <PromotionDialog
                color={pendingMove.color}
                onSelectPiece={handlePromotion}
              />
            )}

            {isPanelVisible && DebugPanel && debugActions && (
              <Suspense fallback={null}>
                <DebugPanel
                  actions={debugActions}
                  currentFen={game.fen()}
                  onClose={() => setIsPanelVisible(false)}
                />
              </Suspense>
            )}
          </>
        )}
      </GameController>
    </div>
  );
}

export default App;
```

**After (with Zustand):**
```typescript
import { lazy, Suspense } from 'react';
import GameBoard from './components/GameBoard';
import GameInfo from './components/GameInfo';
import PromotionDialog from './components/PromotionDialog';
import { useDebugPanel } from './hooks/useDebugPanel';
import { useGameStore } from './store/useGameStore';

const DebugPanel = import.meta.env.DEV
  ? lazy(() => import('./components/debug/DebugPanel'))
  : null;

function App() {
  const { isPanelVisible, setIsPanelVisible } = useDebugPanel();
  const pendingMove = useGameStore(state => state.pendingMove);
  const debugActions = useGameStore(state => state.debugActions);
  const game = useGameStore(state => state.game);

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-800 text-white p-8">
      <div className="flex flex-col lg:flex-row gap-4 lg:gap-8 items-center">
        <GameBoard />
        <GameInfo />
      </div>

      {pendingMove && (
        <PromotionDialog
          color={pendingMove.color}
          onSelectPiece={useGameStore.getState().handlePromotion}
        />
      )}

      {isPanelVisible && DebugPanel && debugActions && (
        <Suspense fallback={null}>
          <DebugPanel
            actions={debugActions}
            currentFen={game.fen()}
            onClose={() => setIsPanelVisible(false)}
          />
        </Suspense>
      )}
    </div>
  );
}

export default App;
```

**Key Changes:**
- ✅ Removed `<GameController>` wrapper
- ✅ Added `useGameStore` hook calls
- ✅ Simplified component - no render props
- ✅ Components get data directly from store

---

### Step 4: Refactor GameBoard.tsx

**Before (with props):**
```typescript
import type { Chess } from 'chess.js';
import Square from './Square';
import type { LastMove } from './GameController';
import type {
  SquareColor,
  Square as SquareType,
  ChessFile,
  ChessRank,
  ChessPiece,
  PieceType,
  PieceColor,
} from '../types/chess';

interface GameBoardProps {
  game: Chess;
  selectedSquare: SquareType | null;
  validMoves: SquareType[];
  lastMove: LastMove;
  onSquareClick: (square: SquareType) => void;
}

const GameBoard = ({ game, selectedSquare, validMoves, lastMove, onSquareClick }: GameBoardProps) => {
  // ... implementation
};

export default GameBoard;
```

**After (with Zustand):**
```typescript
import Square from './Square';
import { useGameStore } from '../store/useGameStore';
import type {
  SquareColor,
  Square as SquareType,
  ChessFile,
  ChessRank,
  ChessPiece,
  PieceType,
  PieceColor,
} from '../types/chess';

const GameBoard = () => {
  // Get data from store
  const game = useGameStore(state => state.game);
  const selectedSquare = useGameStore(state => state.selectedSquare);
  const validMoves = useGameStore(state => state.validMoves);
  const lastMove = useGameStore(state => state.lastMove);
  const selectSquare = useGameStore(state => state.selectSquare);

  const board = game.board();

  const files: ChessFile[] = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'];
  const ranks: ChessRank[] = ['1', '2', '3', '4', '5', '6', '7', '8'];

  const squares: Array<{
    name: SquareType;
    color: SquareColor;
    piece: ChessPiece | null;
  }> = [];

  // Iterate from rank 8 down to rank 1 (top to bottom visually)
  for (let rankIndex = 7; rankIndex >= 0; rankIndex--) {
    const rank = ranks[rankIndex];

    // Iterate from file 'a' to 'h' (left to right)
    for (let fileIndex = 0; fileIndex < 8; fileIndex++) {
      const file = files[fileIndex];
      const squareName = `${file}${rank}` as SquareType;

      // Calculate color
      const isLight = (rankIndex + fileIndex) % 2 !== 0;
      const squareColor: SquareColor = isLight ? 'light' : 'dark';

      // Get piece from board
      const chessJsPiece = board[7 - rankIndex][fileIndex];
      const pieceData = chessJsPiece
        ? {
            type: chessJsPiece.type as PieceType,
            color: chessJsPiece.color as PieceColor,
          }
        : null;

      squares.push({ name: squareName, color: squareColor, piece: pieceData });
    }
  }

  return (
    <div className="relative w-[min(100vw,calc(100vh-4rem))] h-[min(100vw,calc(100vh-4rem))] mx-auto transition-all duration-300">
      {/* Rank labels (8-1) on the left */}
      <div className="absolute -left-6 top-0 h-full flex flex-col justify-around text-gray-400 text-sm transition-all duration-300">
        {ranks
          .slice()
          .reverse()
          .map((rank) => (
            <span key={`rank-${rank}`} className="flex items-center justify-center h-full">
              {rank}
            </span>
          ))}
      </div>

      <div className="w-full h-full grid grid-cols-8 grid-rows-8 border-2 border-[#759656] transition-all duration-300">
        {squares.map((square) => (
          <Square
            key={square.name}
            squareColor={square.color}
            piece={square.piece}
            isSelected={selectedSquare === square.name}
            isValidMove={validMoves.includes(square.name)}
            isLastMoveFrom={lastMove?.from === square.name}
            isLastMoveTo={lastMove?.to === square.name}
            onClick={() => selectSquare(square.name)}
          />
        ))}
      </div>

      {/* File labels (a-h) at the bottom */}
      <div className="absolute -bottom-5 left-0 w-full flex justify-around text-gray-400 text-sm px-2 transition-all duration-300">
        {files.map((file) => (
          <span key={`file-${file}`} className="flex items-center justify-center w-full">
            {file}
          </span>
        ))}
      </div>
    </div>
  );
};

export default GameBoard;
```

**Key Changes:**
- ✅ No props interface needed
- ✅ Added `useGameStore` hook calls at top
- ✅ Changed `onSquareClick` to `selectSquare`
- ✅ Component is self-contained

---

### Step 5: Refactor GameInfo.tsx

**Before (with props):**
```typescript
import type { Chess } from 'chess.js';
import MoveHistory from './MoveHistory';

interface GameInfoProps {
  game: Chess;
  onNewGame: () => void;
}

const GameInfo = ({ game, onNewGame }: GameInfoProps) => {
  // ... implementation
};

export default GameInfo;
```

**After (with Zustand):**
```typescript
import MoveHistory from './MoveHistory';
import { useGameStore } from '../store/useGameStore';

const GameInfo = () => {
  const game = useGameStore(state => state.game);
  const resetGame = useGameStore(state => state.resetGame);

  const turn = game.turn() === 'w' ? 'White' : 'Black';
  const isCheck = game.isCheck();
  const isCheckmate = game.isCheckmate();
  const isDraw = game.isDraw();
  const isStalemate = game.isStalemate();
  const moveHistory = game.history();

  return (
    <div className="bg-gray-700 p-6 rounded-lg w-full lg:w-80">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold">Game Info</h2>
        <button
          onClick={resetGame}
          className="bg-blue-600 hover:bg-blue-500 text-white text-sm font-semibold px-3 py-1 rounded transition-colors"
        >
          New Game
        </button>
      </div>

      <div className="space-y-2">
        <p><span className="font-semibold">Turn:</span> {turn}</p>

        {isCheck && !isCheckmate && (
          <p className="text-yellow-400 font-semibold">Check!</p>
        )}

        {isCheckmate && (
          <p className="text-red-400 font-bold text-xl">
            Checkmate! {turn === 'White' ? 'Black' : 'White'} wins!
          </p>
        )}

        {isStalemate && (
          <p className="text-blue-400 font-bold">Stalemate - Draw!</p>
        )}

        {isDraw && !isStalemate && (
          <p className="text-blue-400 font-bold">Draw!</p>
        )}
      </div>

      <MoveHistory moves={moveHistory} />
    </div>
  );
};

export default GameInfo;
```

**Key Changes:**
- ✅ No props interface
- ✅ Gets `game` and `resetGame` from store
- ✅ Changed `onNewGame` to `resetGame`

---

### Step 6: Refactor DebugPanel.tsx

**Before (with props):**
```typescript
import FenLoader from './FenLoader';
import ScenarioPicker from './ScenarioPicker';
import type { DebugActions } from '../GameController';

interface DebugPanelProps {
  actions: DebugActions;
  currentFen: string;
  onClose: () => void;
}

const DebugPanel = ({ actions, currentFen, onClose }: DebugPanelProps) => {
  // ... implementation
};

export default DebugPanel;
```

**After (with Zustand):**
```typescript
import FenLoader from './FenLoader';
import ScenarioPicker from './ScenarioPicker';
import { useGameStore } from '../../store/useGameStore';

interface DebugPanelProps {
  onClose: () => void;
}

const DebugPanel = ({ onClose }: DebugPanelProps) => {
  const debugActions = useGameStore(state => state.debugActions);
  const game = useGameStore(state => state.game);

  if (!debugActions) return null;

  return (
    <div className="fixed bottom-4 right-4 bg-gray-800/90 backdrop-blur-sm border border-gray-700 rounded-lg shadow-2xl z-50 w-full max-w-md p-4">
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-lg font-bold text-white">🛠️ Debug Toolkit</h3>
        <button onClick={onClose} className="text-gray-400 hover:text-white">&times;</button>
      </div>
      <div className="space-y-4">
        <FenLoader onLoad={debugActions.loadFen} currentFen={game.fen()} />
        <ScenarioPicker onSelect={debugActions.loadFen} />
        <div>
          <button
            onClick={debugActions.resetGame}
            className="bg-red-600 hover:bg-red-500 text-white font-semibold px-4 py-2 rounded transition-colors w-full"
          >
            Reset to Start
          </button>
        </div>
      </div>
    </div>
  );
};

export default DebugPanel;
```

**Key Changes:**
- ✅ Only `onClose` prop remains
- ✅ Gets `debugActions` and `game` from store
- ✅ Removed `actions` and `currentFen` props

---

### Step 7: Delete GameController.tsx

```bash
rm src/components/GameController.tsx
```

✅ No longer needed - all logic is in the store

---

## Performance Optimization

### Zustand Selector Optimization

**Problem:** Naive selectors cause unnecessary re-renders.

**Bad (re-renders on any state change):**
```typescript
const state = useGameStore();  // ❌ Re-renders on ANY change
```

**Good (re-renders only when specific values change):**
```typescript
const game = useGameStore(state => state.game);  // ✅ Only when game changes
const selectedSquare = useGameStore(state => state.selectedSquare);  // ✅ Only when selectedSquare changes
```

**Best (combine related selectors with shallow comparison):**
```typescript
import { shallow } from 'zustand/shallow';

const { game, selectedSquare, validMoves } = useGameStore(
  state => ({
    game: state.game,
    selectedSquare: state.selectedSquare,
    validMoves: state.validMoves
  }),
  shallow
);
```

### Recommended Selector Patterns

**GameBoard (needs multiple pieces of state):**
```typescript
import { shallow } from 'zustand/shallow';

const GameBoard = () => {
  const { game, selectedSquare, validMoves, lastMove } = useGameStore(
    state => ({
      game: state.game,
      selectedSquare: state.selectedSquare,
      validMoves: state.validMoves,
      lastMove: state.lastMove
    }),
    shallow
  );

  const selectSquare = useGameStore(state => state.selectSquare);

  // ... rest of component
};
```

**GameInfo (needs only game):**
```typescript
const GameInfo = () => {
  const game = useGameStore(state => state.game);
  const resetGame = useGameStore(state => state.resetGame);

  // ... rest of component
};
```

---

## Testing Checklist

After refactoring, test **all** functionality:

### Core Gameplay ✅
1. ✅ Select white piece (e.g., e2 pawn) - yellow ring appears
2. ✅ Valid moves show (e.g., e3, e4 have yellow dots)
3. ✅ Click e4 - pawn moves, last move highlighting appears
4. ✅ Select black piece - only black pieces selectable
5. ✅ Make several moves for both sides - game progresses

### Pawn Promotion ✅
1. ✅ Set up pawn promotion (use debug panel or play to promotion)
2. ✅ Move pawn to 8th rank - dialog appears
3. ✅ Click Queen - pawn promotes, dialog closes
4. ✅ Last move highlighting appears on promotion

### Game Status ✅
1. ✅ Move history updates after each move
2. ✅ Turn indicator shows correct player
3. ✅ Check status shows when king in check
4. ✅ Checkmate detected and displayed

### UI Controls ✅
1. ✅ Click "New Game" - board resets, history clears
2. ✅ Last move highlighting clears on reset
3. ✅ Selected piece clears on reset

### Debug Panel (Dev Mode) ✅
1. ✅ Press Ctrl+Shift+D - panel appears
2. ✅ Load FEN - board updates
3. ✅ Click scenario - board loads position
4. ✅ Reset button works
5. ✅ Close panel - game continues

### Responsive Design ✅
1. ✅ Resize to mobile - components stack vertically
2. ✅ Resize to desktop - components side by side
3. ✅ Board resizes smoothly with transitions

---

## Common Migration Issues

### Issue 1: Action Functions Not Updating

**Problem:**
```typescript
const selectSquare = useGameStore(state => state.selectSquare);

// Later in component...
<button onClick={() => selectSquare(square)}>  // ❌ May use stale closure
```

**Solution:**
```typescript
// Option A: Get action on every render (recommended)
const selectSquare = useGameStore(state => state.selectSquare);

// Option B: Use getState() for stable reference
<button onClick={() => useGameStore.getState().selectSquare(square)}>
```

### Issue 2: Chess Instance Mutation

**Problem:**
```typescript
const game = useGameStore(state => state.game);
game.move('e4');  // ❌ Mutates store directly, no re-render
```

**Solution:**
```typescript
// Always create new instance after mutation
game.move('e4');
set({ game: new Chess(game.fen()) });  // ✅ Triggers re-render
```

### Issue 3: Debug Actions Type Error

**Problem:**
```typescript
const debugActions = useGameStore(state => state.debugActions);
debugActions.loadFen(fen);  // ❌ Type error: debugActions might be undefined
```

**Solution:**
```typescript
const debugActions = useGameStore(state => state.debugActions);
if (debugActions) {
  debugActions.loadFen(fen);  // ✅ Type guard
}
```

---

## Verification Commands

```bash
# Run TypeScript compiler
pnpm tsc --noEmit

# Run linter
pnpm lint

# Run dev server
pnpm dev

# Build for production (verify no errors)
pnpm build
```

---

## Summary of Changes

| File | Change | Status |
|------|--------|--------|
| `src/store/useGameStore.ts` | **Created** - Complete store | ✅ New |
| `src/components/App.tsx` | **Modified** - Remove GameController wrapper | ✅ Refactor |
| `src/components/GameBoard.tsx` | **Modified** - Use store hooks | ✅ Refactor |
| `src/components/GameInfo.tsx` | **Modified** - Use store hooks | ✅ Refactor |
| `src/components/debug/DebugPanel.tsx` | **Modified** - Use store hooks | ✅ Refactor |
| `src/components/GameController.tsx` | **Deleted** - Logic moved to store | ✅ Remove |

---

## Benefits of Zustand

**Before (Render Props):**
- ❌ 9 parameters in render prop function
- ❌ Prop drilling through multiple components
- ❌ Hard to add new state (update many files)
- ❌ Testing requires mocking render props

**After (Zustand):**
- ✅ Direct access to state from any component
- ✅ No prop drilling
- ✅ Easy to add new state (update store only)
- ✅ Testing with `useGameStore.setState()`

---

## Next Steps

After successful migration:

1. ✅ **Verify all functionality** - Complete testing checklist
2. ✅ **Commit changes** - "refactor: migrate to Zustand state management"
3. ✅ **Ready for Step 12b** - Add replay functionality
4. ✅ **Or proceed to AI Coach** - State management foundation ready

---

## Estimated Time

- **Store creation:** 30 minutes
- **Component refactoring:** 1 hour
- **Testing:** 30 minutes
- **Bug fixes:** 30 minutes

**Total: 2-3 hours** for complete, tested migration.

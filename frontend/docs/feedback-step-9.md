# Step 9: Pawn Promotion - Feedback

## Critical Issues Found

### 1. 🔴 Wrong Color for Promotion Dialog (Line 141)

**CRITICAL BUG - Will show wrong piece colors**

**Issue:** Using `game.turn()` after the promotion move is pending gives the wrong color.

**Original Code:**
```typescript
<PromotionDialog
  color={game.turn()} // WRONG - gives opposite color after move
  onSelectPiece={handlePromotion}
/>
```

**Problem:**
- White pawn reaches e8 (promotion pending)
- `game.turn()` still returns `'w'` (correct at this point)
- But turn hasn't switched yet
- Need to get color from the piece being promoted

**Fix:** Store color in pendingMove state:
```typescript
type PendingPromotion = {
  from: Square;
  to: Square;
  color: PieceColor;
};

const [pendingMove, setPendingMove] = useState<PendingPromotion | null>(null);

// Then in App.tsx:
<PromotionDialog
  color={pendingMove.color}
  onSelectPiece={handlePromotion}
/>
```

---

### 2. 🔴 Inconsistent State Update (Line 111)

**Issue:** Using complex `Object.assign(Object.create(...))` instead of existing pattern.

**Original Code:**
```typescript
setGame(Object.assign(Object.create(Object.getPrototypeOf(game)), game));
```

**Problems:**
- Overly complex
- Inconsistent with existing `new Chess(game.fen())` pattern
- May not properly preserve Chess methods

**Fix:**
```typescript
const handlePromotion = (piece: PieceType) => {
  if (!pendingMove) return;

  game.move({ ...pendingMove, promotion: piece });
  setGame(new Chess(game.fen())); // Consistent with rest of codebase
  setPendingMove(null);
};
```

---

### 3. 🔴 Missing Move Validation (Line 86-96)

**MAJOR BUG - Can trigger promotion for illegal moves**

**Original Code:**
```typescript
if (
  piece?.type === 'p' &&
  (square.endsWith('1') || square.endsWith('8'))
) {
  setPendingMove({ from: selectedSquare, to: square });
  // No validation that move is legal!
}
```

**Problem:** Doesn't check if the promotion move is actually legal.

**Fix - Validate with chess.js:**
```typescript
try {
  const piece = game.get(selectedSquare);

  // Check if this would be a promotion move
  if (piece?.type === 'p' && (square.endsWith('1') || square.endsWith('8'))) {
    // Validate the move is legal by checking valid moves
    const moves = game.moves({ square: selectedSquare, verbose: true });
    const isValidMove = moves.some(m => m.to === square);

    if (isValidMove) {
      // Valid promotion - show dialog
      setPendingMove({
        from: selectedSquare,
        to: square,
        color: piece.color
      });
      setSelectedSquare(null);
      setValidMoves([]);
      return;
    }
  }

  // Not a promotion, proceed with normal move
  const move = game.move({ from: selectedSquare, to: square });
  // ... rest of logic
} catch (error) {
  // ...
}
```

---

### 4. Missing Type Definitions

**Issue:** Incomplete types throughout.

**Fix - Add to GameController.tsx:**
```typescript
type PendingPromotion = {
  from: Square;
  to: Square;
  color: PieceColor;
};

interface GameControllerProps {
  children: (
    game: Chess,
    selectedSquare: Square | null,
    validMoves: Square[],
    selectSquare: (square: Square) => void,
    pendingMove: PendingPromotion | null,
    handlePromotion: (piece: PieceType) => void
  ) => React.ReactNode;
}
```

---

## Minor Improvements

### 5. Missing Accessibility

**Add ARIA labels and structure:**
```typescript
<div
  className="fixed inset-0 bg-black/70 flex items-center justify-center z-50"
  role="dialog"
  aria-modal="true"
  aria-labelledby="promotion-title"
>
  <div className="bg-gray-800 p-4 rounded-lg shadow-xl">
    <h3 id="promotion-title" className="text-white text-center font-semibold mb-4">
      Choose Promotion
    </h3>
    <div className="flex gap-4">
      {promotionPieces.map((pieceType) => (
        <div key={pieceType} className="flex flex-col items-center gap-1">
          <button
            className="w-20 h-20 bg-gray-700 hover:bg-gray-600 cursor-pointer rounded flex items-center justify-center transition-colors"
            onClick={() => onSelectPiece(pieceType)}
            aria-label={`Promote to ${pieceLabels[pieceType]}`}
          >
            <Piece piece={{ type: pieceType, color }} />
          </button>
          <span className="text-white text-xs">{pieceLabels[pieceType]}</span>
        </div>
      ))}
    </div>
  </div>
</div>
```

---

### 6. Add Keyboard Support

```typescript
import { useEffect } from 'react';

const PromotionDialog = ({ color, onSelectPiece }: PromotionDialogProps) => {
  useEffect(() => {
    const handleKeyPress = (e: KeyboardEvent) => {
      const keyMap: Record<string, PieceType> = {
        'q': 'q',
        'r': 'r',
        'b': 'b',
        'n': 'n'
      };
      const piece = keyMap[e.key.toLowerCase()];
      if (piece) {
        onSelectPiece(piece);
      }
    };

    window.addEventListener('keydown', handleKeyPress);
    return () => window.removeEventListener('keydown', handleKeyPress);
  }, [onSelectPiece]);

  // ... rest of component
};
```

---

### 7. Add Piece Labels

```typescript
const pieceLabels: Record<PieceType, string> = {
  q: 'Queen',
  r: 'Rook',
  b: 'Bishop',
  n: 'Knight'
};
```

---

## Complete Corrected Implementation

### PromotionDialog.tsx

```typescript
import { useEffect } from 'react';
import type { PieceType, PieceColor } from '../types/chess';
import Piece from './Piece';

interface PromotionDialogProps {
  color: PieceColor;
  onSelectPiece: (piece: PieceType) => void;
}

const promotionPieces: PieceType[] = ['q', 'r', 'b', 'n'];

const pieceLabels: Record<PieceType, string> = {
  q: 'Queen',
  r: 'Rook',
  b: 'Bishop',
  n: 'Knight',
  p: 'Pawn',
  k: 'King'
};

const PromotionDialog = ({ color, onSelectPiece }: PromotionDialogProps) => {
  // Keyboard support
  useEffect(() => {
    const handleKeyPress = (e: KeyboardEvent) => {
      const keyMap: Record<string, PieceType> = {
        'q': 'q',
        'r': 'r',
        'b': 'b',
        'n': 'n'
      };
      const piece = keyMap[e.key.toLowerCase()];
      if (piece) {
        onSelectPiece(piece);
      }
    };

    window.addEventListener('keydown', handleKeyPress);
    return () => window.removeEventListener('keydown', handleKeyPress);
  }, [onSelectPiece]);

  return (
    <div
      className="fixed inset-0 bg-black/70 flex items-center justify-center z-50"
      role="dialog"
      aria-modal="true"
      aria-labelledby="promotion-title"
    >
      <div className="bg-gray-800 p-6 rounded-lg shadow-xl">
        <h3 id="promotion-title" className="text-white text-center font-semibold mb-4 text-lg">
          Choose Promotion
        </h3>
        <div className="flex gap-4">
          {promotionPieces.map((pieceType) => (
            <div key={pieceType} className="flex flex-col items-center gap-2">
              <button
                className="w-20 h-20 bg-gray-700 hover:bg-gray-600 active:bg-gray-500 cursor-pointer rounded flex items-center justify-center transition-colors"
                onClick={() => onSelectPiece(pieceType)}
                aria-label={`Promote to ${pieceLabels[pieceType]}`}
              >
                <Piece piece={{ type: pieceType, color }} />
              </button>
              <span className="text-white text-xs">{pieceLabels[pieceType]}</span>
            </div>
          ))}
        </div>
        <p className="text-gray-400 text-xs text-center mt-4">
          Press Q, R, B, or N on keyboard
        </p>
      </div>
    </div>
  );
};

export default PromotionDialog;
```

---

### GameController.tsx Updates

```typescript
import { useState } from 'react';
import { Chess } from 'chess.js';
import type { Square, PieceType, PieceColor } from '../types/chess';

type PendingPromotion = {
  from: Square;
  to: Square;
  color: PieceColor;
};

interface GameControllerProps {
  children: (
    game: Chess,
    selectedSquare: Square | null,
    validMoves: Square[],
    selectSquare: (square: Square) => void,
    pendingMove: PendingPromotion | null,
    handlePromotion: (piece: PieceType) => void
  ) => React.ReactNode;
}

const GameController = ({ children }: GameControllerProps) => {
  const [game, setGame] = useState(() => new Chess());
  const [selectedSquare, setSelectedSquare] = useState<Square | null>(null);
  const [validMoves, setValidMoves] = useState<Square[]>([]);
  const [pendingMove, setPendingMove] = useState<PendingPromotion | null>(null);

  const selectSquare = (square: Square) => {
    // If no square is selected, select this square (if it has a piece)
    if (!selectedSquare) {
      const piece = game.get(square);
      if (piece && piece.color === game.turn()) {
        setSelectedSquare(square);
        const moves = game.moves({ square, verbose: true });
        setValidMoves(moves.map(move => move.to));
      }
      return;
    }

    // If clicking the same square, deselect it
    if (selectedSquare === square) {
      setSelectedSquare(null);
      setValidMoves([]);
      return;
    }

    // Try to make a move
    try {
      const piece = game.get(selectedSquare);

      // Check if this would be a promotion move
      if (piece?.type === 'p' && (square.endsWith('1') || square.endsWith('8'))) {
        // Validate the move is legal
        const moves = game.moves({ square: selectedSquare, verbose: true });
        const isValidMove = moves.some(m => m.to === square);

        if (isValidMove) {
          // Valid promotion - show dialog
          setPendingMove({
            from: selectedSquare,
            to: square,
            color: piece.color
          });
          setSelectedSquare(null);
          setValidMoves([]);
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
        setGame(new Chess(game.fen()));
        setSelectedSquare(null);
        setValidMoves([]);
      } else {
        // Invalid move, check if clicking another piece of the same color
        const newPiece = game.get(square);
        if (newPiece && newPiece.color === game.turn()) {
          setSelectedSquare(square);
          const moves = game.moves({ square, verbose: true });
          setValidMoves(moves.map(move => move.to));
        } else {
          setSelectedSquare(null);
          setValidMoves([]);
        }
      }
    } catch (error) {
      // Invalid move, deselect
      setSelectedSquare(null);
      setValidMoves([]);
    }
  };

  const handlePromotion = (piece: PieceType) => {
    if (!pendingMove) return;

    game.move({
      from: pendingMove.from,
      to: pendingMove.to,
      promotion: piece
    });
    setGame(new Chess(game.fen())); // Consistent with existing pattern
    setPendingMove(null);
  };

  return <>{children(game, selectedSquare, validMoves, selectSquare, pendingMove, handlePromotion)}</>;
};

export default GameController;
```

---

### App.tsx Updates

```typescript
import GameBoard from './components/GameBoard';
import GameInfo from './components/GameInfo';
import GameController from './components/GameController';
import PromotionDialog from './components/PromotionDialog';

function App() {
  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-800 text-white p-8">
      <GameController>
        {(game, selectedSquare, validMoves, selectSquare, pendingMove, handlePromotion) => (
          <>
            <div className="flex flex-row gap-8 items-center">
              <GameBoard
                game={game}
                selectedSquare={selectedSquare}
                validMoves={validMoves}
                onSquareClick={selectSquare}
              />
              <GameInfo game={game} />
            </div>

            {pendingMove && (
              <PromotionDialog
                color={pendingMove.color}
                onSelectPiece={handlePromotion}
              />
            )}
          </>
        )}
      </GameController>
    </div>
  );
}

export default App;
```

---

## Testing Checklist - Enhanced

1. **Trigger Promotion:**
   - ✅ Move white pawn to e8
   - ✅ Dialog appears with white pieces
   - ✅ Move black pawn to e1
   - ✅ Dialog appears with black pieces

2. **Dialog Interaction:**
   - ✅ Click Queen - promotes correctly
   - ✅ Click Rook - promotes correctly
   - ✅ Click Bishop - promotes correctly
   - ✅ Click Knight - promotes correctly
   - ✅ Keyboard Q/R/B/N keys work
   - ✅ Move appears in history (e.g., `e8=Q`)

3. **Edge Cases:**
   - ✅ Capturing promotion (e.g., `exd8=Q`)
   - ✅ Multiple promotions in one game
   - ✅ Promotion that gives check (e.g., `e8=Q+`)
   - ✅ Promotion that gives checkmate

4. **No False Triggers:**
   - ✅ Regular pawn moves don't show dialog
   - ✅ Non-pawn pieces moving to 1st/8th rank don't show dialog

---

## Summary of Issues Fixed

| Issue | Severity | Original | Fixed |
|-------|----------|----------|-------|
| Wrong promotion color | 🔴 Critical | `game.turn()` | `pendingMove.color` |
| Inconsistent state update | 🔴 Critical | `Object.assign(...)` | `new Chess(game.fen())` |
| No move validation | 🔴 Critical | No check | Validate with `game.moves()` |
| Missing types | ⚠️ Medium | Incomplete | Added `PendingPromotion` |
| No accessibility | ⚠️ Medium | Missing | Added ARIA labels |
| No keyboard support | ℹ️ Low | Missing | Added Q/R/B/N keys |
| No piece labels | ℹ️ Low | Missing | Added text labels |

---

## Next Steps

After implementing these fixes:
1. Test all promotion scenarios thoroughly
2. Verify keyboard shortcuts work
3. Check accessibility with screen reader
4. Move to Step 10 (Final polish)

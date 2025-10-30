# Step 7: Valid Move Highlighting

## Goal
When a piece is selected, visually highlight all valid destination squares to help users understand legal moves.

## Background
Currently, users can click pieces and make moves, but they don't know which moves are legal until they try. This step adds visual feedback showing all valid moves for the selected piece.

## Implementation Tasks

### 1. Update GameController to Calculate Valid Moves

Modify `frontend/app/src/components/GameController.tsx`:

```typescript
import { useState } from 'react';
import { Chess } from 'chess.js';
import type { Square } from '../types/chess';

interface GameControllerProps {
  children: (
    game: Chess,
    selectedSquare: Square | null,
    validMoves: Square[],
    selectSquare: (square: Square) => void
  ) => React.ReactNode;
}

const GameController = ({ children }: GameControllerProps) => {
  const [game, setGame] = useState(() => new Chess());
  const [selectedSquare, setSelectedSquare] = useState<Square | null>(null);
  const [validMoves, setValidMoves] = useState<Square[]>([]);

  const selectSquare = (square: Square) => {
    // If no square is selected, select this square (if it has a piece)
    if (!selectedSquare) {
      const piece = game.get(square);
      if (piece && piece.color === game.turn()) {
        setSelectedSquare(square);
        // Calculate valid moves for this piece
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
      const move = game.move({
        from: selectedSquare,
        to: square,
        promotion: 'q' // Always promote to queen for now (Step 10 will add dialog)
      });

      if (move) {
        // Move was successful, update state
        setGame(new Chess(game.fen())); // Create new instance to trigger re-render
        setSelectedSquare(null);
        setValidMoves([]);
      } else {
        // Invalid move, check if clicking another piece of the same color
        const piece = game.get(square);
        if (piece && piece.color === game.turn()) {
          setSelectedSquare(square);
          // Calculate valid moves for new piece
          const moves = game.moves({ square, verbose: true });
          setValidMoves(moves.map(move => move.to as Square));
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

  return <>{children(game, selectedSquare, validMoves, selectSquare)}</>;
};

export default GameController;
```

**Key Changes:**
- Added `validMoves` state: `Square[]`
- Calculate moves using `game.moves({ square, verbose: true })`
- `verbose: true` returns move objects with `to` property
- Map moves to extract destination squares
- Clear `validMoves` when deselecting or after move

### 2. Update App.tsx to Pass Valid Moves

Modify `frontend/app/src/App.tsx`:

```typescript
import GameBoard from './components/GameBoard';
import GameInfo from './components/GameInfo';
import GameController from './components/GameController';

function App() {
  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-800 text-white p-8">
      <GameController>
        {(game, selectedSquare, validMoves, selectSquare) => (
          <div className="flex flex-row gap-8 items-center">
            <GameBoard
              game={game}
              selectedSquare={selectedSquare}
              validMoves={validMoves}
              onSquareClick={selectSquare}
            />
            <GameInfo game={game} />
          </div>
        )}
      </GameController>
    </div>
  );
}

export default App;
```

### 3. Update GameBoard to Pass Valid Move Info

Modify `frontend/app/src/components/GameBoard.tsx`:

**Update interface:**
```typescript
interface GameBoardProps {
  game: Chess;
  selectedSquare: SquareType | null;
  validMoves: SquareType[];
  onSquareClick: (square: SquareType) => void;
}
```

**Update component signature:**
```typescript
const GameBoard = ({ game, selectedSquare, validMoves, onSquareClick }: GameBoardProps) => {
```

**Update Square components (around line 72-80):**
```typescript
<Square
  key={square.name}
  squareColor={square.color}
  squareName={square.name}
  piece={square.piece}
  isSelected={selectedSquare === square.name}
  isValidMove={validMoves.includes(square.name)}
  onClick={() => onSquareClick(square.name)}
/>
```

### 4. Update Square to Show Valid Move Indicators

Modify `frontend/app/src/components/Square.tsx`:

**Update interface:**
```typescript
interface SquareProps {
  squareColor: SquareColor;
  squareName: SquareType;
  piece: ChessPiece | null;
  isSelected: boolean;
  isValidMove: boolean;
  onClick: () => void;
}
```

**Update component:**
```typescript
const Square = ({ squareColor, squareName, piece, isSelected, isValidMove, onClick }: SquareProps) => {
  const bgColor = squareColor === 'light'
    ? 'bg-[#e8edd5]'
    : 'bg-[#759656]';

  return (
    <div
      className={`
        ${bgColor}
        ${isSelected ? 'ring-4 ring-yellow-400 ring-inset' : ''}
        flex items-center justify-center relative cursor-pointer
        hover:brightness-90 transition-all
      `}
      onClick={onClick}
    >
      {piece && <Piece piece={piece} />}

      {/* Valid move indicator */}
      {isValidMove && (
        <div className={`
          absolute rounded-full
          ${piece ? 'w-full h-full border-4 border-yellow-500/60' : 'w-4 h-4 bg-yellow-500/60'}
        `} />
      )}
    </div>
  );
};
```

**Visual Design:**
- **Empty square:** Small yellow circle (4x4, 60% opacity)
- **Square with opponent piece:** Yellow border ring (indicates capture)
- Uses `absolute` positioning to overlay indicators
- `/60` in Tailwind = 60% opacity

## Testing Checklist

After implementation, test the following:

1. **Valid Move Display:**
   - ✅ Select e2 pawn - should show dots on e3 and e4
   - ✅ Select knight - should show 2-3 L-shaped moves
   - ✅ Select rook (after moving pieces) - should show straight lines
   - ✅ Select bishop - should show diagonal moves

2. **Capture Indicators:**
   - ✅ Move to mid-game with opponent pieces nearby
   - ✅ Select piece that can capture - opponent squares show border ring
   - ✅ Empty squares show dot, capture squares show ring

3. **Visual Clarity:**
   - ✅ Yellow indicators are visible on both light and dark squares
   - ✅ Indicators don't obscure pieces
   - ✅ Selected square ring is distinct from valid move indicators

4. **Interaction:**
   - ✅ Click valid move square - move executes
   - ✅ Click non-valid square - nothing happens or piece switches
   - ✅ Deselect piece - indicators disappear
   - ✅ Switch between pieces - indicators update correctly

5. **Edge Cases:**
   - ✅ Pinned pieces show only legal moves (not moves that expose king)
   - ✅ King in check - only shows moves that escape check
   - ✅ Castling shows as valid move when legal
   - ✅ En passant shows as valid move when legal

## Visual Preview

```
   Before Selection          After Selecting e2 Pawn
┌─────────────────┐      ┌─────────────────┐
│ r n b q k b n r │      │ r n b q k b n r │
│ p p p p p p p p │      │ p p p p p p p p │
│ · · · · · · · · │      │ · · · · · · · · │
│ · · · · · · · · │      │ · · · · · · · · │
│ · · · · · · · · │      │ · · · ◉ · · · · │  ← Yellow dot (e4)
│ · · · · · · · · │      │ · · · ◉ · · · · │  ← Yellow dot (e3)
│ P P P P P P P P │      │ P P P [P] P P P │  ← Yellow ring (selected)
│ R N B Q K B N R │      │ R N B Q K B N R │
└─────────────────┘      └─────────────────┘

   Piece Can Capture         Castling Available
┌─────────────────┐      ┌─────────────────┐
│ · · · · k · · · │      │ r · · · k · · r │
│ · · · · · · · · │      │ · · · · · · · · │
│ · · · n · · · · │      │ · · · · · · · · │
│ · · · [Q] · · · │      │ · · · · · · · · │
│ · · ◉ · ◉ · · · │      │ · · · · · · · · │
│ · ◯ · · · ◯ · · │      │ · · · · · · · · │
│ · · ◉ · ◉ · · · │      │ · · · · · · · · │
│ · · · · K · · · │      │ · · · [K] ◉ · ◉ │  ← Dots on g1 (kingside) & c1 (queenside)
└─────────────────┘      └─────────────────┘
   ◯ = border (capture)
   ◉ = dot (empty)
```

## chess.js API Reference

**Getting valid moves for a piece:**
```typescript
// Get moves for specific square
const moves = game.moves({ square: 'e2', verbose: true });
// Returns: [{ from: 'e2', to: 'e3', ... }, { from: 'e2', to: 'e4', ... }]

// Extract destination squares
const destinations = moves.map(m => m.to); // ['e3', 'e4']
```

**Move object structure (verbose mode):**
```typescript
{
  color: 'w' | 'b',
  from: 'e2',
  to: 'e4',
  piece: 'p',
  captured?: 'p',  // Only present if capture
  promotion?: 'q', // Only present if pawn promotion
  flags: 'b',      // 'b' = big pawn, 'c' = capture, 'e' = en passant, etc.
  san: 'e4',       // Standard Algebraic Notation
  lan: 'e2e4'      // Long Algebraic Notation
}
```

## Next Steps

Once this step is complete:
- **Step 8:** Enhanced game info display (move history)
- **Step 9:** End game detection improvements
- **Step 10:** Pawn promotion dialog
- **Step 11:** Final polish and responsiveness

## SOLID Principles Applied

- **SRP:** Square handles display, GameController handles move calculation
- **OCP:** Adding new indicator styles doesn't require changing move logic
- **ISP:** Square receives only the boolean `isValidMove`, not entire move objects

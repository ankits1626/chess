# Step 6: Move Implementation (Click-to-Move)

## Goal
Implement interactive chess gameplay where users can select pieces and make moves by clicking on squares.

## Background
Currently, the board displays the initial position but is not interactive. In this step, we'll add:
1. Click handling to select pieces
2. Move execution when clicking destination squares
3. State management to track selected piece
4. Visual feedback for the selected piece

## Implementation Tasks

### 1. Add Move Logic to GameController

Modify `frontend/app/src/components/GameController.tsx`:

```typescript
import { useState } from 'react';
import { Chess } from 'chess.js';
import type { Square } from '../types/chess';

interface GameControllerProps {
  children: (
    game: Chess,
    selectedSquare: Square | null,
    selectSquare: (square: Square) => void
  ) => React.ReactNode;
}

const GameController = ({ children }: GameControllerProps) => {
  const [game, setGame] = useState(() => new Chess());
  const [selectedSquare, setSelectedSquare] = useState<Square | null>(null);

  const selectSquare = (square: Square) => {
    // If no square is selected, select this square (if it has a piece)
    if (!selectedSquare) {
      const piece = game.get(square);
      if (piece && piece.color === game.turn()) {
        setSelectedSquare(square);
      }
      return;
    }

    // If clicking the same square, deselect it
    if (selectedSquare === square) {
      setSelectedSquare(null);
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
      } else {
        // Invalid move, check if clicking another piece of the same color
        const piece = game.get(square);
        if (piece && piece.color === game.turn()) {
          setSelectedSquare(square);
        } else {
          setSelectedSquare(null);
        }
      }
    } catch (error) {
      // Invalid move, deselect
      setSelectedSquare(null);
    }
  };

  return <>{children(game, selectedSquare, selectSquare)}</>;
};

export default GameController;
```

**Key Design Decisions:**
- `selectedSquare` state tracks which piece is selected
- `selectSquare` function handles all click logic:
  - Select piece if none selected
  - Deselect if clicking same square
  - Attempt move if different square clicked
  - Allow switching between pieces of same color
- Always promote to queen (temporary - Step 10 will add dialog)
- Create new `Chess` instance from FEN to trigger React re-render

### 2. Update App.tsx to Pass New Props

Modify `frontend/app/src/App.tsx`:

```typescript
import GameBoard from './components/GameBoard';
import GameInfo from './components/GameInfo';
import GameController from './components/GameController';

function App() {
  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-800 text-white p-8">
      <GameController>
        {(game, selectedSquare, selectSquare) => (
          <div className="flex flex-row gap-8 items-center">
            <GameBoard
              game={game}
              selectedSquare={selectedSquare}
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

### 3. Update GameBoard to Handle Clicks

Modify `frontend/app/src/components/GameBoard.tsx`:

**Add to interface:**
```typescript
interface GameBoardProps {
  game: Chess;
  selectedSquare: Square | null;
  onSquareClick: (square: Square) => void;
}
```

**Update component signature:**
```typescript
const GameBoard = ({ game, selectedSquare, onSquareClick }: GameBoardProps) => {
```

**Pass props to Square components:**
```typescript
<Square
  key={square.name}
  squareColor={square.color}
  squareName={square.name}
  piece={square.piece}
  isSelected={selectedSquare === square.name}
  onClick={() => onSquareClick(square.name)}
/>
```

### 4. Update Square to Handle Selection

Modify `frontend/app/src/components/Square.tsx`:

**Add to interface:**
```typescript
interface SquareProps {
  squareColor: SquareColor;
  squareName: Square;
  piece: ChessPiece | null;
  isSelected: boolean;
  onClick: () => void;
}
```

**Update component:**
```typescript
const Square = ({ squareColor, squareName, piece, isSelected, onClick }: SquareProps) => {
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
    </div>
  );
};
```

**Key Visual Changes:**
- `ring-4 ring-yellow-400 ring-inset` - Yellow ring for selected square
- `cursor-pointer` - Shows clickable cursor
- `hover:brightness-90` - Subtle hover effect

### 5. Remove Unused Import from GameController

The `Square` type import in GameController.tsx line 3 is not used. Remove it:

```typescript
import { useState } from 'react';
import { Chess } from 'chess.js';
```

## Testing Checklist

After implementation, test the following:

1. **Piece Selection:**
   - ✅ Click white piece (e.g., e2 pawn) - should show yellow ring
   - ✅ Click empty square - nothing happens
   - ✅ Click black piece on white's turn - nothing happens

2. **Move Execution:**
   - ✅ Select e2 pawn, click e4 - pawn should move
   - ✅ Board should update to show new position
   - ✅ GameInfo should show "Black" turn after move

3. **Deselection:**
   - ✅ Select piece, click it again - deselects (ring disappears)
   - ✅ Select piece, click invalid square - deselects

4. **Piece Switching:**
   - ✅ Select e2 pawn, then click d2 pawn - should switch selection

5. **Full Game Flow:**
   - ✅ Make several moves for both white and black
   - ✅ Turn indicator updates correctly
   - ✅ Only current player's pieces can be selected

6. **Visual Feedback:**
   - ✅ Selected square has yellow ring
   - ✅ Hover effect works on all squares
   - ✅ Pieces render correctly after moves

## Known Limitations (To Be Addressed Later)

- No valid move highlighting (Step 7)
- Pawn promotion always defaults to queen (Step 10)
- No move history or undo (Future enhancement)
- No drag-and-drop support (Future enhancement)

## Next Steps

Once this step is complete:
- **Step 7:** Valid move highlighting (show legal moves for selected piece)
- **Step 8:** Enhanced game info display
- **Step 9:** End game detection improvements
- **Step 10:** Pawn promotion dialog

## SOLID Principles Applied

- **SRP:** Square handles display and click events, GameController handles game logic
- **OCP:** Move logic is centralized and can be extended without modifying Square
- **DIP:** Components depend on callbacks (abstractions) not concrete implementations

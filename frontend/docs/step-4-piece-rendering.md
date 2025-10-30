# Step 4: Piece Rendering Plan

This document outlines the plan for rendering chess pieces on the board in their initial positions. This involves implementing the `Piece.tsx` component and integrating `chess.js` to determine piece placement.

## 1. Goal

To display all 32 chess pieces in their standard starting positions on the static board created in Step 3. The piece data will be sourced from the `chess.js` library.

## 2. Component Implementation

### 2.1. Required Types (`src/types/chess.ts`)

We need to extend our type definitions to include pieces.

```typescript
// src/types/chess.ts

// ... existing types ...

export type PieceType = 'p' | 'n' | 'b' | 'r' | 'q' | 'k';
export type PieceColor = 'w' | 'b';

export interface ChessPiece {
  type: PieceType;
  color: PieceColor;
}
```

### 2.2. `src/components/Piece.tsx`

This component will render a single chess piece.

-   **Props:** It will accept a `piece: ChessPiece` object containing its `type` and `color`.
-   **Rendering:** It will use Unicode characters for the visual representation of the pieces. A mapping object will be used to select the correct character.
-   **Styling:** The piece will be a large, centered text character. The color will be set based on the `piece.color` prop (e.g., `text-gray-800` for black, `text-gray-100` for white) to ensure visibility on both light and dark squares.

**Example Implementation:**

```typescript
// src/components/Piece.tsx
import { ChessPiece } from '../types/chess';

interface PieceProps {
  piece: ChessPiece;
}

const pieceUnicode: Record<string, string> = {
  wk: '♔', wq: '♕', wr: '♖', wb: '♗', wn: '♘', wp: '♙',
  bk: '♚', bq: '♛', br: '♜', bb: '♝', bn: '♞', bp: '♟︎',
};

const Piece = ({ piece }: PieceProps) => {
  const pieceKey = `${piece.color}${piece.type}`;
  const textColor = piece.color === 'w' ? 'text-gray-100' : 'text-gray-800';

  return (
    <div className={`text-4xl ${textColor} relative z-10`}>
      {pieceUnicode[pieceKey]}
    </div>
  );
};

export default Piece;
```

### 2.3. `src/components/Square.tsx` (Update)

This component needs to be updated to accept and render a piece.

-   **Props:** The `SquareProps` interface will be updated to include an optional `piece: ChessPiece | null`.
-   **Rendering:** It will conditionally render the `<Piece />` component if the `piece` prop is not null. The square name will be rendered with a lower z-index to appear behind the piece.

**Example Update:**

```typescript
// src/components/Square.tsx
import { SquareColor, Square as SquareType, ChessPiece } from '../types/chess';
import Piece from './Piece';

interface SquareProps {
  squareColor: SquareColor;
  squareName: SquareType;
  piece: ChessPiece | null;
}

const Square = ({ squareColor, squareName, piece }: SquareProps) => {
  const bgColor = squareColor === 'light' ? 'bg-amber-100' : 'bg-amber-700';

  return (
    <div className={`${bgColor} flex items-center justify-center relative`}>
      {piece && <Piece piece={piece} />}
      <span className="absolute bottom-0 right-1 text-xs opacity-30 select-none z-0">
        {squareName}
      </span>
    </div>
  );
};

export default Square;
```

### 2.4. `src/components/GameBoard.tsx` (Update)

This component will be updated to use `chess.js` to get the board state and pass piece data down to the `Square` components.

-   **Logic:**
    1.  Import and instantiate the `Chess` object from `chess.js`.
    2.  Call the `game.board()` method to get a 2D array representing the board's state. Each element is either `null` or a piece object (`{ type, color, square }`).
    3.  Modify the existing rendering loop. For each square being rendered, look up the corresponding piece from the `game.board()` array.
-   **Rendering:** Pass the retrieved piece object (or `null`) to each `Square` component via the new `piece` prop.

**Example Update:**

```typescript
// src/components/GameBoard.tsx
import { Chess } from 'chess.js';
import Square from './Square';
import { SquareColor, Square as SquareType, ChessFile, ChessRank, ChessPiece } from '../types/chess';

const GameBoard = () => {
  const game = new Chess();
  const board = game.board(); // Get the board state

  const files: ChessFile[] = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'];
  const ranks: ChessRank[] = ['1', '2', '3', '4', '5', '6', '7', '8'];

  // ... (The existing loop for generating squares can be removed or repurposed)

  return (
    <div className="w-full max-w-lg mx-auto aspect-square grid grid-cols-8 border-2 border-gray-900">
      {board.flat().map((piece, index) => {
        const rankIndex = 7 - Math.floor(index / 8);
        const fileIndex = index % 8;

        const squareName = `${files[fileIndex]}${ranks[rankIndex]}` as SquareType;
        const isLight = (rankIndex + fileIndex) % 2 !== 0;
        const squareColor: SquareColor = isLight ? 'light' : 'dark';

        // The piece object from chess.js has the same shape as our ChessPiece type
        const pieceData = piece ? { type: piece.type, color: piece.color } as ChessPiece : null;

        return (
          <Square
            key={squareName}
            squareColor={squareColor}
            squareName={squareName}
            piece={pieceData}
          />
        );
      })}
    </div>
  );
};

export default GameBoard;
```

## 3. Development Steps

1.  **Update Types (`src/types/chess.ts`):** Add the `PieceType`, `PieceColor`, and `ChessPiece` types.
2.  **Implement `src/components/Piece.tsx`:** Create the component to render a piece using Unicode characters.
3.  **Update `src/components/Square.tsx`:**
    -   Modify the `SquareProps` to accept an optional `piece`.
    -   Conditionally render the `<Piece />` component.
    -   Adjust styling for the square name to avoid overlapping with the piece.
4.  **Update `src/components/GameBoard.tsx`:**
    -   Import and instantiate `chess.js`.
    -   Fetch the board state using `game.board()`.
    -   Refactor the rendering logic to map over the `board` array.
    -   For each square, determine its color and name, and pass the corresponding piece data to the `Square` component.

## 4. Verification Checklist

After implementation, run `pnpm dev` and verify the following:

-   [ ] All 32 pieces are rendered in their correct starting positions.
-   [ ] White pieces are on ranks 1 and 2.
-   [ ] Black pieces are on ranks 7 and 8.
-   [ ] The piece characters (e.g., ♙, ♖, ♚) are correct for each position.
-   [ ] White pieces have a light color, and black pieces have a dark color.
-   [ ] Pieces are clearly visible on both light and dark squares.
-   [ ] The underlying board grid and square colors are still correct.
-   [ ] There are no errors or warnings in the browser console.

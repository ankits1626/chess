# Step 3: Static Board Rendering Plan (Updated)

This document outlines the plan for rendering the static chessboard using `GameBoard.tsx` and `Square.tsx`, incorporating feedback for a more robust implementation.

## 1. Goal

To display an 8x8 chessboard grid with alternating light and dark squares, without any pieces or game logic yet. This serves as the foundational visual component of the application.

## 2. Component Implementation

### 2.1. Required Types (`src/types/chess.ts`)

First, ensure the necessary TypeScript types are defined. These provide type safety for our components.

```typescript
// src/types/chess.ts
export type SquareColor = 'light' | 'dark';
export type ChessFile = 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h';
export type ChessRank = '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8';
export type Square = `${ChessFile}${ChessRank}`;
```

### 2.2. `src/components/Square.tsx`

This component will represent a single square on the chessboard.

-   **Props:** It will receive `squareColor: SquareColor` and `squareName: Square` via a `SquareProps` interface.
-   **Styling:**
    -   The size will be determined by the `GameBoard`'s grid layout.
    -   Apply background colors conditionally based on `squareColor`. We will use `bg-amber-100` for light squares and `bg-amber-700` for dark squares.
-   **Content:** It will display its `squareName` with low opacity (`opacity-30`) for debugging purposes, which prevents visual clutter.

**Example Implementation:**

```typescript
// src/components/Square.tsx
import { SquareColor, Square as SquareType } from '../types/chess';

interface SquareProps {
  squareColor: SquareColor;
  squareName: SquareType;
}

const Square = ({ squareColor, squareName }: SquareProps) => {
  const bgColor = squareColor === 'light'
    ? 'bg-amber-100'
    : 'bg-amber-700';

  return (
    <div
      className={`${bgColor} flex items-center justify-center`}
    >
      <span className="text-xs opacity-30 select-none">
        {squareName}
      </span>
    </div>
  );
};

export default Square;
```

### 2.3. `src/components/GameBoard.tsx`

This component will render the entire 8x8 grid of squares.

-   **Structure:** It will use Tailwind's `w-full max-w-lg mx-auto aspect-square grid grid-cols-8` classes to create a responsive, square container that adapts to its parent's width while maintaining a maximum size and aspect ratio.
-   **Logic:**
    -   It will iterate to generate an array of 64 square objects, each with a `name` and `color`.
    -   **Board Orientation:** To follow chess conventions (rank 1 at the bottom), the iteration will loop ranks from **7 down to 0** (top to bottom visually) and files from **0 to 7** (left to right).
    -   **Color Calculation:** The color is determined by the rule: a square is light if `(rankIndex + fileIndex)` is odd, and dark if even. To ensure `a1` is dark, the logic is `(rankIndex + fileIndex) % 2 !== 0`.
-   **Rendering:** It will map over the generated squares array, rendering a `Square` component for each. A unique `key` prop (e.g., `key={square.name}`) is crucial for React's rendering process.

**Example Implementation:**

```typescript
// src/components/GameBoard.tsx
import Square from './Square';
import { SquareColor, Square as SquareType, ChessFile, ChessRank } from '../types/chess';

const GameBoard = () => {
  const files: ChessFile[] = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'];
  const ranks: ChessRank[] = ['1', '2', '3', '4', '5', '6', '7', '8'];

  const squares: Array<{ name: Square; color: SquareColor }> = [];

  // Iterate from rank 8 down to rank 1 (top to bottom visually)
  for (let rankIndex = 7; rankIndex >= 0; rankIndex--) {
    const rank = ranks[rankIndex];

    // Iterate from file 'a' to 'h' (left to right)
    for (let fileIndex = 0; fileIndex < 8; fileIndex++) {
      const file = files[fileIndex];
      const squareName = `${file}${rank}` as Square;

      // Correct color calculation: a1 should be dark.
      const isLight = (rankIndex + fileIndex) % 2 !== 0;
      const squareColor: SquareColor = isLight ? 'light' : 'dark';

      squares.push({ name: squareName, color: squareColor });
    }
  }

  return (
    <div className="w-full max-w-lg mx-auto aspect-square grid grid-cols-8 border-2 border-gray-900">
      {squares.map((square) => (
        <Square
          key={square.name}
          squareColor={square.color}
          squareName={square.name}
        />
      ))}
    </div>
  );
};

export default GameBoard;
```

## 3. Development Steps

1.  **Create Types (`src/types/chess.ts`):** If not already done, create the file and add the `SquareColor`, `ChessFile`, `ChessRank`, and `Square` types.
2.  **Update `src/components/Square.tsx`:**
    -   Import the types.
    -   Define the `SquareProps` interface.
    -   Implement the component with conditional background colors, removing any fixed size classes.
    -   Display the `squareName` with low opacity.
3.  **Update `src/components/GameBoard.tsx`:**
    -   Import the `Square` component and types.
    -   Implement the nested loops to generate the `squares` array with correct color calculation and board orientation.
    -   Render the squares inside a `div` with responsive `grid` layout (`w-full max-w-lg mx-auto aspect-square grid grid-cols-8`), ensuring each `Square` has a unique `key`.

## 4. Verification Checklist

After implementation, run `pnpm dev` and verify the following:

-   [ ] The board displays as a perfect 8×8 grid.
-   [ ] The square colors alternate correctly.
-   [ ] The bottom-left square (`a1`) is **dark**.
-   [ ] The bottom-right square (`h1`) is **light**.
-   [ ] The top-left square (`a8`) is **light**.
-   [ ] The top-right square (`h8`) is **dark**.
-   [ ] Each square displays its name faintly (e.g., "a1", "h8").
-   [ ] There are no errors or warnings in the browser console (especially no "missing key" warnings).

## 5. Common Pitfalls to Avoid

-   **Wrong Color on a1:** If `a1` is light, the color calculation is inverted. Use `(rankIndex + fileIndex) % 2 !== 0` for `isLight`.
-   **Upside-Down Board:** If rank 1 is at the top, the rank iteration is wrong. Iterate ranks from `7` down to `0`.
-   **Non-Square Shapes:** If squares appear rectangular, ensure the `GameBoard` container has `aspect-square` and the grid is correctly defined (`grid grid-cols-8`).
-   **Missing `key` Prop:** Always provide a unique `key` prop in `map()` loops to avoid React warnings and potential bugs.
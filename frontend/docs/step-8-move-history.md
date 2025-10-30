# Step 8: Move History Display

## Goal
To display a formatted list of all moves made during the game, providing users with a complete record of the game's progression.

## Background
The `GameInfo` panel currently shows the turn and game status. We will enhance it by adding a scrollable move history log, which is a standard feature in chess applications.

## Implementation Tasks

### 1. Create `MoveHistory.tsx` Component

Create a new component at `frontend/app/src/components/MoveHistory.tsx`. This component will be responsible for rendering the list of moves.

-   **Props:** It will accept a `moves: string[]` array.
-   **Rendering:** It will map over the `moves` array and display them in a numbered list, showing pairs of moves (White's move and Black's move) on each line.
-   **Styling:** It will be a scrollable container with a defined height.

**Example Implementation:**

```typescript
// src/components/MoveHistory.tsx
import { useEffect, useRef } from 'react';

interface MoveHistoryProps {
  moves: string[];
}

const MoveHistory = ({ moves }: MoveHistoryProps) => {
  const scrollContainerRef = useRef<HTMLDivElement>(null);

  // Auto-scroll to the latest move
  useEffect(() => {
    if (scrollContainerRef.current) {
      scrollContainerRef.current.scrollTop = scrollContainerRef.current.scrollHeight;
    }
  }, [moves]);

  return (
    <div className="mt-4">
      <h3 className="text-lg font-semibold mb-2">Move History</h3>
      <div
        ref={scrollContainerRef}
        className="h-48 bg-gray-800 p-2 rounded overflow-y-auto"
      >
        <ol className="text-white">
          {moves.reduce((acc, move, index) => {
            if (index % 2 === 0) {
              // Start of a new move pair
              acc.push([move]);
            } else {
              // Add black's move to the last pair
              acc[acc.length - 1].push(move);
            }
            return acc;
          }, [] as string[][]).map((pair, i) => (
            <li key={i} className="grid grid-cols-3 gap-2 py-1 px-2 rounded hover:bg-gray-700">
              <span className="text-gray-400">{i + 1}.</span>
              <span className="col-span-1">{pair[0]}</span>
              {pair[1] && <span className="col-span-1">{pair[1]}</span>}
            </li>
          ))}
        </ol>
      </div>
    </div>
  );
};

export default MoveHistory;
```

### 2. Update `GameInfo.tsx` to Use `MoveHistory`

Modify `frontend/app/src/components/GameInfo.tsx` to include the new `MoveHistory` component.

-   **Logic:** Get the move history from the `game` prop using `game.history()`.
-   **Rendering:** Pass the history array to the `<MoveHistory />` component.

**Example Update:**

```typescript
// src/components/GameInfo.tsx
import type { Chess } from 'chess.js';
import MoveHistory from './MoveHistory'; // Import the new component

interface GameInfoProps {
  game: Chess;
}

const GameInfo = ({ game }: GameInfoProps) => {
  const turn = game.turn() === 'w' ? 'White' : 'Black';
  const isCheck = game.isCheck();
  const isCheckmate = game.isCheckmate();
  const isDraw = game.isDraw();
  const isStalemate = game.isStalemate();
  const moveHistory = game.history(); // Get move history

  return (
    <div className="bg-gray-700 p-6 rounded-lg w-full lg:w-64">
      <h2 className="text-2xl font-bold mb-4">Game Info</h2>

      <div className="space-y-2">
        <p><span className="font-semibold">Turn:</span> {turn}</p>
        {/* ... existing status messages ... */}
      </div>

      {/* Add the MoveHistory component */}
      <MoveHistory moves={moveHistory} />
    </div>
  );
};

export default GameInfo;
```

## Testing Checklist

After implementation, test the following:

1.  **Initial State:**
    -   ✅ The "Move History" panel appears below the "Game Info".
    -   ✅ The move list is initially empty.

2.  **Making Moves:**
    -   ✅ Make a move for White (e.g., e4). The list should show "1. e4".
    -   ✅ Make a move for Black (e.g., e5). The list should update to "1. e4 e5".
    -   ✅ Continue making moves and verify the list populates correctly.

3.  **Scrolling:**
    -   ✅ Make enough moves to exceed the panel's height.
    -   ✅ A scrollbar should appear.
    -   ✅ The list should automatically scroll to show the most recent move.

4.  **Visuals:**
    -   ✅ The move list is easy to read.
    -   ✅ The styling matches the overall theme of the application.

## Next Steps

Once this step is complete:
- **Step 9:** End game detection improvements
- **Step 10:** Pawn promotion dialog
- **Step 11:** Final polish and responsiveness

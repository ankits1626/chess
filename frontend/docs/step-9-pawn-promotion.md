# Step 9: Pawn Promotion

## Goal
When a pawn reaches the opposite end of the board, present the user with a dialog to choose which piece to promote the pawn to (Queen, Rook, Bishop, or Knight).

## Background
Currently, in `GameController.tsx`, pawn promotion is hardcoded to always select a queen (`promotion: 'q'`). This step will replace that with a user-interactive dialog, making the game fully playable according to standard chess rules.

## Implementation Tasks

### 1. Create `PromotionDialog.tsx` Component

This component will be a modal dialog that appears when a promotion is possible.

-   **Props:** It will receive an `onSelectPiece` callback function and the `color` of the pawn being promoted.
-   **State:** It will be conditionally rendered based on state managed in `GameController`.
-   **Rendering:** It will display the four promotion piece options (Queen, Rook, Bishop, Knight) for the correct color.
-   **Interaction:** Clicking on a piece will call the `onSelectPiece` callback with the selected piece type (`'q'`, `'r'`, `'b'`, or `'n'`).

**Example Implementation (`frontend/app/src/components/PromotionDialog.tsx`):**

```typescript
import type { PieceType, PieceColor } from '../types/chess';
import Piece from './Piece';

interface PromotionDialogProps {
  color: PieceColor;
  onSelectPiece: (piece: PieceType) => void;
}

const promotionPieces: PieceType[] = ['q', 'r', 'b', 'n'];

const PromotionDialog = ({ color, onSelectPiece }: PromotionDialogProps) => {
  return (
    <div className="fixed inset-0 bg-black/70 flex items-center justify-center z-50">
      <div className="bg-gray-800 p-4 rounded-lg flex gap-4">
        {promotionPieces.map((pieceType) => (
          <div
            key={pieceType}
            className="w-20 h-20 bg-gray-700 hover:bg-gray-600 cursor-pointer rounded flex items-center justify-center"
            onClick={() => onSelectPiece(pieceType)}
          >
            <Piece piece={{ type: pieceType, color }} />
          </div>
        ))}
      </div>
    </div>
  );
};

export default PromotionDialog;
```

### 2. Update `GameController.tsx` to Manage Promotion State

We need to add state to track when a promotion is pending and handle the promotion logic.

**Key Changes:**

1.  **Add State:** Create state to hold the pending move details.
    ```typescript
    const [pendingMove, setPendingMove] = useState<{ from: Square; to: Square } | null>(null);
    ```
2.  **Modify `selectSquare` Logic:**
    -   When a pawn move to the 1st or 8th rank is detected, don't execute the move immediately.
    -   Instead, store the `from` and `to` squares in the `pendingMove` state.
3.  **Create a `handlePromotion` Function:**
    -   This function will be called by the `PromotionDialog`.
    -   It will execute the move using the stored `pendingMove` and the chosen promotion piece.
    -   It will then clear the `pendingMove` state.

**Example `GameController.tsx` Updates:**

```typescript
// ... imports

// Add state for the pending promotion move
const [pendingMove, setPendingMove] = useState<{ from: Square; to: Square } | null>(null);

const selectSquare = (square: Square) => {
  // ... existing selection logic ...

  // In the 'Try to make a move' block:
  try {
    const piece = game.get(selectedSquare);
    // Check if the move is a promotion
    if (
      piece?.type === 'p' &&
      (square.endsWith('1') || square.endsWith('8'))
    ) {
      // It's a promotion, so we set it as a pending move and wait for user input
      setPendingMove({ from: selectedSquare, to: square });
      setSelectedSquare(null);
      setValidMoves([]);
      return;
    }

    // If not a promotion, proceed as before
    const move = game.move({ from: selectedSquare, to: square });

    // ... rest of the move logic ...
  } catch (error) {
    // ...
  }
};

const handlePromotion = (piece: PieceType) => {
  if (!pendingMove) return;

  game.move({ ...pendingMove, promotion: piece });
  setGame(Object.assign(Object.create(Object.getPrototypeOf(game)), game));
  setPendingMove(null);
};

// Update the children render prop to include pendingMove and handlePromotion
return <>{children(game, selectedSquare, validMoves, selectSquare, pendingMove, handlePromotion)}</>;
```

### 3. Update `App.tsx` to Render the Dialog

The `App` component will now receive the `pendingMove` state and the `handlePromotion` callback and will conditionally render the `PromotionDialog`.

**Example `App.tsx` Update:**

```typescript
// ... imports
import PromotionDialog from './components/PromotionDialog';

function App() {
  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-800 text-white p-8">
      <GameController>
        {(game, selectedSquare, validMoves, selectSquare, pendingMove, handlePromotion) => (
          <>
            <div className="flex flex-row gap-8 items-center">
              {/* ... GameBoard and GameInfo ... */}
            </div>

            {pendingMove && (
              <PromotionDialog
                color={game.turn()} // The color of the player who is about to promote
                onSelectPiece={handlePromotion}
              />
            )}
          </>
        )}
      </GameController>
    </div>
  );
}
```

## Testing Checklist

1.  **Trigger Promotion:**
    -   ✅ Play a game until a pawn is one square away from the final rank.
    -   ✅ Move the pawn to the final rank.
    -   ✅ The `PromotionDialog` should appear, overlaying the game.

2.  **Dialog Interaction:**
    -   ✅ The dialog should show four piece options (Queen, Rook, Bishop, Knight) of the correct color.
    -   ✅ Click the Queen. The dialog should disappear, and the pawn should be replaced by a Queen on the board.
    -   ✅ The move should be recorded in the move history (e.g., `e8=Q`).
    -   ✅ The turn should switch to the other player.

3.  **Other Promotions:**
    -   ✅ Test promoting to a Rook, Bishop, and Knight (underpromotion) to ensure they work correctly.

4.  **No Dialog Case:**
    -   ✅ Play several non-promotion moves. The dialog should not appear.

## Next Steps

Once pawn promotion is fully functional, the core gameplay mechanics will be complete. The next steps will focus on polishing the application:

-   **Step 10:** Final polish and responsiveness improvements.
-   **(Future) Step 11:** Implement Drag-and-Drop functionality.

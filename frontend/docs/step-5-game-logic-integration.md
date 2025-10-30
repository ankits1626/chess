# Step 5: Game Logic Integration

## Goal
Implement `GameController.tsx` to manage the chess game state using `chess.js` and provide that state to the rest of the application.

## Background
Currently, `GameBoard.tsx` creates its own `Chess` instance and displays the initial board position. In this step, we'll centralize game state management in a dedicated component following the Single Responsibility Principle.

## Implementation Tasks

### 1. Create GameController Component

Create `frontend/app/src/components/GameController.tsx`:

```typescript
import { useState } from 'react';
import { Chess } from 'chess.js';
import type { Square } from '../types/chess';

interface GameControllerProps {
  children: (game: Chess) => React.ReactNode;
}

const GameController = ({ children }: GameControllerProps) => {
  const [game, setGame] = useState(() => new Chess());

  return <>{children(game)}</>;
};

export default GameController;
```

**Key Design Decisions:**
- Uses **render props pattern** (`children` as a function) to pass game state down
- Manages `Chess` instance in React state so updates trigger re-renders
- Component doesn't render any UI itself (SRP)
- Parent components receive `game` instance and can call methods on it

### 2. Update App.tsx to Use GameController

Modify `frontend/app/src/App.tsx`:

```typescript
import GameBoard from './components/GameBoard';
import GameInfo from './components/GameInfo';
import GameController from './components/GameController';

function App() {
  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-800 text-white p-8">
      <GameController>
        {(game) => (
          <div className="flex flex-row gap-8 items-center">
            <GameBoard game={game} />
            <GameInfo game={game} />
          </div>
        )}
      </GameController>
    </div>
  );
}

export default App;
```

### 3. Update GameBoard to Accept Game Prop

Modify `frontend/app/src/components/GameBoard.tsx`:

**Change this:**
```typescript
const GameBoard = () => {
  const game = new Chess();
  const board = game.board();
```

**To this:**
```typescript
import type { Chess } from 'chess.js';

interface GameBoardProps {
  game: Chess;
}

const GameBoard = ({ game }: GameBoardProps) => {
  const board = game.board();
```

### 4. Update GameInfo to Accept Game Prop

Modify `frontend/app/src/components/GameInfo.tsx`:

```typescript
import type { Chess } from 'chess.js';

interface GameInfoProps {
  game: Chess;
}

const GameInfo = ({ game }: GameInfoProps) => {
  const turn = game.turn() === 'w' ? 'White' : 'Black';
  const isCheck = game.isCheck();
  const isCheckmate = game.isCheckmate();
  const isDraw = game.isDraw();
  const isStalemate = game.isStalemate();

  return (
    <div className="bg-gray-700 p-6 rounded-lg">
      <h2 className="text-2xl font-bold mb-4">Game Info</h2>

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
    </div>
  );
};

export default GameInfo;
```

## Testing

After implementation:
1. The board should still display the initial chess position
2. GameInfo should show "White" as the current turn
3. No errors in the browser console
4. The application structure is now ready for move implementation in Step 6

## Next Steps

Once this step is complete, we'll implement:
- **Step 6:** Move implementation (click-to-move)
- **Step 7:** Valid move highlighting
- **Step 8:** Interactive game play

## SOLID Principles Applied

- **SRP:** GameController has one job - manage game state
- **DIP:** Components depend on abstractions (Chess interface) not concrete implementations
- **ISP:** Each component receives only the props it needs (game instance)

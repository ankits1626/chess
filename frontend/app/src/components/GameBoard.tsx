import Square from './Square';
import type { SquareColor, Square as SquareType, ChessFile, ChessRank } from '../types/chess';

const GameBoard = () => {
  const files: ChessFile[] = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'];
  const ranks: ChessRank[] = ['1', '2', '3', '4', '5', '6', '7', '8'];

  const squares: Array<{ name: SquareType; color: SquareColor }> = [];

  // Iterate from rank 8 down to rank 1 (top to bottom visually)
  for (let rankIndex = 7; rankIndex >= 0; rankIndex--) {
    const rank = ranks[rankIndex];

    // Iterate from file 'a' to 'h' (left to right)
    for (let fileIndex = 0; fileIndex < 8; fileIndex++) {
      const file = files[fileIndex];
      const squareName = `${file}${rank}` as SquareType;

      // Correct color calculation: a1 should be dark.
      const isLight = (rankIndex + fileIndex) % 2 !== 0;
      const squareColor: SquareColor = isLight ? 'light' : 'dark';

      squares.push({ name: squareName, color: squareColor });
    }
  }

  return (
    <div className="w-[min(100vw,100vh)] h-[min(100vw,100vh)] mx-auto grid grid-cols-8 border-2 border-gray-900">
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

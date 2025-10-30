import type { Chess } from 'chess.js';
import Square from './Square';
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
  onSquareClick: (square: SquareType) => void;
}

const GameBoard = ({ game, selectedSquare, onSquareClick }: GameBoardProps) => {
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
    <div className="relative">
      {/* Rank labels (8-1) on the left */}
      <div className="absolute -left-6 top-0 h-[min(100vw,calc(100vh-4rem))] flex flex-col justify-around text-gray-400 text-sm transition-all duration-300">
        {ranks
          .slice()
          .reverse()
          .map((rank) => (
            <span key={`rank-${rank}`} className="flex items-center justify-center h-full">
              {rank}
            </span>
          ))}
      </div>

      <div className="w-[min(100vw,calc(100vh-4rem))] h-[min(100vw,calc(100vh-4rem))] mx-auto grid grid-cols-8 grid-rows-8 border-2 border-gray-900 transition-all duration-300">
        {squares.map((square) => (
          <Square
            key={square.name}
            squareColor={square.color}
            squareName={square.name}
            piece={square.piece}
            isSelected={selectedSquare === square.name}
            onClick={() => onSquareClick(square.name)}
          />
        ))}
      </div>

      {/* File labels (a-h) at the bottom */}
      <div className="w-[min(100vw,calc(100vh-4rem))] mx-auto flex justify-around text-gray-400 text-sm mt-2 transition-all duration-300">
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

import type { SquareColor, Square as SquareType, ChessPiece } from '../types/chess';
import Piece from './Piece';

interface SquareProps {
  squareColor: SquareColor;
  squareName: SquareType;
  piece: ChessPiece | null;
}

const Square = ({ squareColor, squareName, piece }: SquareProps) => {
  const bgColor = squareColor === 'light'
    ? 'bg-[#e8edd5]'  // Light green (like chess.com)
    : 'bg-[#759656]'; // Dark green

  return (
    <div
      className={`${bgColor} flex items-center justify-center relative`}
    >
      {piece && <Piece piece={piece} />}
      <span className="absolute bottom-0 right-1 text-xs opacity-30 select-none z-0">
        {squareName}
      </span>
    </div>
  );
};

export default Square;

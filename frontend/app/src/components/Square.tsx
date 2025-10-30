import type { SquareColor, Square as SquareType, ChessPiece } from '../types/chess';
import Piece from './Piece';

interface SquareProps {
  squareColor: SquareColor;
  squareName: SquareType;
  piece: ChessPiece | null;
  isSelected: boolean;
  isValidMove: boolean;
  onClick: () => void;
}

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

export default Square;

import type { SquareColor, ChessPiece } from '@/types/chess';
import Piece from './Piece';

interface SquareProps {
  squareColor: SquareColor;
  piece: ChessPiece | null;
  isSelected: boolean;
  isValidMove: boolean;
  isLastMoveFrom: boolean;
  isLastMoveTo: boolean;
  onClick: () => void;
}

const Square = ({
  squareColor,
  piece,
  isSelected,
  isValidMove,
  isLastMoveFrom,
  isLastMoveTo,
  onClick
}: SquareProps) => {
  const bgColor = squareColor === 'light'
    ? 'bg-[#e8edd5]'
    : 'bg-[#759656]';

  // Last move highlighting (subtle yellow tint)
  const lastMoveHighlight = (isLastMoveFrom || isLastMoveTo)
    ? 'bg-yellow-200/30'
    : '';

  return (
    <div
      className={`
        ${bgColor}
        ${lastMoveHighlight}
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

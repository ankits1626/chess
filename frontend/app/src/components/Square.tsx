import type { SquareColor, Square as SquareType } from '../types/chess';

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

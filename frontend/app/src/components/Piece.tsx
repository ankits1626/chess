import type { ChessPiece } from '../types/chess';

interface PieceProps {
  piece: ChessPiece;
}

const Piece = ({ piece }: PieceProps) => {
  // Note: chess.js uses 'w' and 'b' for color, but the lichess SVGs use 'w' and 'b' followed by the uppercase piece initial.
  // The plan used lowercase (e.g., wk.svg), but the URLs used uppercase (e.g., wK.svg). I will assume uppercase filenames.
  const pieceKey = `${piece.color}${piece.type.toUpperCase()}`;
  const pieceImage = `/pieces/${pieceKey}.svg`;

  return (
    <img
      src={pieceImage}
      alt={`${piece.color} ${piece.type}`}
      className="w-full h-full p-1 cursor-pointer hover:scale-110 transition-transform"
      draggable="false"
    />
  );
};

export default Piece;

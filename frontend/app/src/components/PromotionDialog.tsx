import { useEffect } from 'react';
import type { PieceType, PieceColor } from '../types/chess';
import Piece from './Piece';

interface PromotionDialogProps {
  color: PieceColor;
  onSelectPiece: (piece: PieceType) => void;
}

const promotionPieces: PieceType[] = ['q', 'r', 'b', 'n'];

const pieceLabels: Record<PieceType, string> = {
  q: 'Queen',
  r: 'Rook',
  b: 'Bishop',
  n: 'Knight',
  p: 'Pawn', // Added for completeness, though not used in promotion
  k: 'King'  // Added for completeness
};

const PromotionDialog = ({ color, onSelectPiece }: PromotionDialogProps) => {
  // Keyboard support
  useEffect(() => {
    const handleKeyPress = (e: KeyboardEvent) => {
      const keyMap: Record<string, PieceType> = {
        'q': 'q',
        'r': 'r',
        'b': 'b',
        'n': 'n'
      };
      const piece = keyMap[e.key.toLowerCase()];
      if (piece) {
        onSelectPiece(piece);
      }
    };

    window.addEventListener('keydown', handleKeyPress);
    return () => window.removeEventListener('keydown', handleKeyPress);
  }, [onSelectPiece]);

  return (
    <div
      className="fixed inset-0 bg-black/70 flex items-center justify-center z-50"
      role="dialog"
      aria-modal="true"
      aria-labelledby="promotion-title"
    >
      <div className="bg-gray-800 p-6 rounded-lg shadow-xl">
        <h3 id="promotion-title" className="text-white text-center font-semibold mb-4 text-lg">
          Choose Promotion
        </h3>
        <div className="flex gap-4">
          {promotionPieces.map((pieceType) => (
            <div key={pieceType} className="flex flex-col items-center gap-2">
              <button
                className="w-20 h-20 bg-gray-700 hover:bg-gray-600 active:bg-gray-500 cursor-pointer rounded flex items-center justify-center transition-colors"
                onClick={() => onSelectPiece(pieceType)}
                aria-label={`Promote to ${pieceLabels[pieceType]}`}
              >
                <Piece piece={{ type: pieceType, color }} />
              </button>
              <span className="text-white text-xs">{pieceLabels[pieceType]}</span>
            </div>
          ))}
        </div>
        <p className="text-gray-400 text-xs text-center mt-4">
          Press Q, R, B, or N on keyboard
        </p>
      </div>
    </div>
  );
};

export default PromotionDialog;
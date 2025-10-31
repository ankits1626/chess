import { useRef, useEffect } from 'react';
import type { FC } from 'react';
import type { Move } from 'chess.js';

interface MoveListProps {
  moves: Move[];
  currentMoveIndex: number;
  onMoveClick: (moveIndex: number) => void;
}

const MoveList: FC<MoveListProps> = ({ moves, currentMoveIndex, onMoveClick }) => {
  const currentMoveRef = useRef<HTMLSpanElement>(null);

  useEffect(() => {
    // Automatically scroll to the current move
    currentMoveRef.current?.scrollIntoView({
      behavior: 'smooth',
      block: 'nearest',
    });
  }, [currentMoveIndex]);

  // Group moves into pairs for display (e.g., 1. e4 e5)
  const movePairs: { moveNumber: number; white: Move; black?: Move }[] = [];
  for (let i = 0; i < moves.length; i += 2) {
    movePairs.push({
      moveNumber: i / 2 + 1,
      white: moves[i],
      black: moves[i + 1],
    });
  }

  return (
    <div className="w-full h-96 bg-gray-900 rounded-lg p-4 overflow-y-auto">
      <h3 className="text-lg font-bold mb-2 text-white">Moves</h3>
      <div className="font-mono text-sm text-gray-300 space-y-1">
        {movePairs.map((pair, pairIndex) => {
          const whiteMoveIndex = pairIndex * 2;
          const blackMoveIndex = pairIndex * 2 + 1;

          const isWhiteCurrent = currentMoveIndex === whiteMoveIndex;
          const isBlackCurrent = currentMoveIndex === blackMoveIndex;

          return (
            <div key={pair.moveNumber} className="flex items-center gap-4">
              <span className="w-8 text-right text-gray-500">{pair.moveNumber}.</span>
              <span
                ref={isWhiteCurrent ? currentMoveRef : null}
                onClick={() => onMoveClick(whiteMoveIndex)}
                className={`cursor-pointer rounded px-2 py-0.5 transition-colors ${
                  isWhiteCurrent ? 'bg-blue-600 text-white' : 'hover:bg-gray-700'
                }`}
              >
                {pair.white.san}
              </span>
              {pair.black && (
                <span
                  ref={isBlackCurrent ? currentMoveRef : null}
                  onClick={() => onMoveClick(blackMoveIndex)}
                  className={`cursor-pointer rounded px-2 py-0.5 transition-colors ${
                    isBlackCurrent ? 'bg-blue-600 text-white' : 'hover:bg-gray-700'
                  }`}
                >
                  {pair.black.san}
                </span>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default MoveList;

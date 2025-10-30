import { useEffect, useRef, useMemo } from 'react';

interface MoveHistoryProps {
  moves: string[];
}

type MovePair = [white: string, black?: string];

const MoveHistory = ({ moves }: MoveHistoryProps) => {
  const scrollContainerRef = useRef<HTMLDivElement>(null);

  // Pair moves into [white, black] tuples
  // Only re-computes when moves array changes
  const movePairs = useMemo<MovePair[]>(() => {
    return moves.reduce((acc, move, index) => {
      if (index % 2 === 0) {
        acc.push([move]);
      } else {
        acc[acc.length - 1].push(move);
      }
      return acc;
    }, [] as MovePair[]);
  }, [moves]);

  // Auto-scroll to the latest move
  useEffect(() => {
    if (scrollContainerRef.current) {
      scrollContainerRef.current.scrollTop = scrollContainerRef.current.scrollHeight;
    }
  }, [moves]);

  return (
    <div className="mt-4">
      <h3 className="text-lg font-semibold mb-2">Move History</h3>
      <div
        ref={scrollContainerRef}
        className="h-48 bg-gray-800 p-2 rounded overflow-y-auto"
        role="log"
        aria-live="polite"
        aria-label="Move history"
      >
        {moves.length === 0 ? (
          <p className="text-gray-400 text-sm text-center py-4">No moves yet</p>
        ) : (
          <ol className="text-white">
            {movePairs.map((pair, i) => (
              <li
                key={i}
                className={`
                  grid grid-cols-[auto_1fr_1fr] gap-2 py-1 px-2 rounded
                  hover:bg-gray-700 transition-colors
                  ${i === movePairs.length - 1 ? 'bg-gray-700' : ''}
                `}
              >
                <span className="text-gray-400">{i + 1}.</span>
                <span>{pair[0]}</span>
                <span>{pair[1] || ''}</span>
              </li>
            ))}
          </ol>
        )}
      </div>
    </div>
  );
};

export default MoveHistory;
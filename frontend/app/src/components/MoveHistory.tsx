import { useEffect, useRef } from 'react';

interface MoveHistoryProps {
  moves: string[];
}

const MoveHistory = ({ moves }: MoveHistoryProps) => {
  const scrollContainerRef = useRef<HTMLDivElement>(null);

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
      >
        <ol className="text-white">
          {moves.reduce((acc, move, index) => {
            if (index % 2 === 0) {
              // Start of a new move pair
              acc.push([move]);
            } else {
              // Add black's move to the last pair
              acc[acc.length - 1].push(move);
            }
            return acc;
          }, [] as string[][]).map((pair, i) => (
            <li key={i} className="grid grid-cols-3 gap-2 py-1 px-2 rounded hover:bg-gray-700">
              <span className="text-gray-400">{i + 1}.</span>
              <span className="col-span-1">{pair[0]}</span>
              {pair[1] && <span className="col-span-1">{pair[1]}</span>}
            </li>
          ))}
        </ol>
      </div>
    </div>
  );
};

export default MoveHistory;

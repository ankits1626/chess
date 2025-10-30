import { useState } from 'react';

interface FenLoaderProps {
  onLoad: (fen: string) => void;
  currentFen: string;
}

const FenLoader = ({ onLoad, currentFen }: FenLoaderProps) => {
  const [fen, setFen] = useState(currentFen);

  const handleLoad = () => {
    if (fen.trim()) {
      onLoad(fen.trim());
    }
  };

  return (
    <div className="flex flex-col gap-2">
      <label htmlFor="fen-input" className="text-sm font-medium text-gray-300">Load FEN</label>
      <div className="flex gap-2">
        <input
          id="fen-input"
          type="text"
          value={fen}
          onChange={(e) => setFen(e.target.value)}
          className="flex-grow bg-gray-900 text-white p-2 rounded border border-gray-600 focus:ring-2 focus:ring-sky-500 outline-none"
          placeholder="Paste FEN string here..."
        />
        <button
          onClick={handleLoad}
          className="bg-sky-600 hover:bg-sky-500 text-white font-semibold px-4 py-2 rounded transition-colors"
        >
          Load
        </button>
      </div>
    </div>
  );
};

export default FenLoader;

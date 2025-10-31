import { useState } from 'react';
import { useGameStore } from '../store/useGameStore';
import { fetchLatestGamePgn, RateLimitError } from '../services/chesscomApi';
import { Chess } from 'chess.js'; // To parse PGN for preview

interface GamePreview {
  white: string;
  black: string;
  result: string;
  date: string;
  timeControl: string;
  rated: boolean;
  whiteElo?: number;
  blackElo?: number;
  pgn: string;
}

const isValidUsername = (username: string): boolean => {
  // Chess.com usernames: 3-20 characters, alphanumeric + underscore/dash
  return /^[a-zA-Z0-9_-]{3,20}$/.test(username);
};

const GameImporter = () => {
  const [username, setUsername] = useState('hikaru'); // Default for easy testing
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [gamePreview, setGamePreview] = useState<GamePreview | null>(null);
  const [fetchProgress, setFetchProgress] = useState<string>('');
  const [abortController, setAbortController] = useState<AbortController | null>(null);

  const loadPgn = useGameStore(state => state.loadPgn);

  const handleFetch = async () => {
    if (!isValidUsername(username)) {
      setError('Invalid username format. Use 3-20 characters (letters, numbers, -, _).');
      return;
    }

    const controller = new AbortController();
    setAbortController(controller);
    setIsLoading(true);
    setError(null);
    setGamePreview(null);
    setFetchProgress('Fetching archives...');

    try {
      const pgn = await fetchLatestGamePgn(username, controller.signal);
      setFetchProgress('Parsing game...');

      // Parse PGN for preview metadata
      const tempGame = new Chess();
      tempGame.loadPgn(pgn);
      const header = tempGame.header();
      const latestGame = {
        time_control: header.TimeControl || '',
        rated: header.Rated === 'True',
        white: { rating: header.WhiteElo ? parseInt(header.WhiteElo) : undefined, username: header.White || '' },
        black: { rating: header.BlackElo ? parseInt(header.BlackElo) : undefined, username: header.Black || '' },
      }; // Mocking structure from ChessComGamesResponse for metadata extraction

      setGamePreview({
        white: header.White || 'Unknown',
        black: header.Black || 'Unknown',
        result: header.Result || '*' ,
        date: header.Date || '',
        timeControl: latestGame.time_control || 'Unknown',
        rated: latestGame.rated ?? true,
        whiteElo: latestGame.white.rating,
        blackElo: latestGame.black.rating,
        pgn
      });
      setFetchProgress('');
    } catch (err) {
      if (err instanceof Error && err.name === 'AbortError') {
        setError('Fetch cancelled');
      } else if (err instanceof RateLimitError) {
        setError(`Rate limited. Please try again ${err.retryAfter ? `in ${err.retryAfter}s` : 'later'}.`);
      } else {
        setError(err instanceof Error ? err.message : 'Failed to fetch game');
      }
      setFetchProgress('');
    } finally {
      setIsLoading(false);
      setAbortController(null);
    }
  };

  const handleCancel = () => {
    abortController?.abort();
  };

  const handleLoad = () => {
    if (gamePreview) {
      loadPgn(gamePreview.pgn);
      setGamePreview(null); // Clear preview after loading
    }
  };

  const handleUsernameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUsername(e.target.value);
    setError(null); // Clear error on type
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !isLoading) {
      handleFetch();
    }
  };

  return (
    <div className="p-4 bg-gray-900 rounded-lg space-y-4">
      <h3 className="text-xl font-bold text-white">Import Game from Chess.com</h3>
      <div className="flex gap-2">
        <input
          type="text"
          value={username}
          onChange={handleUsernameChange}
          onKeyDown={handleKeyDown}
          placeholder="Chess.com username"
          className="bg-gray-700 p-2 rounded text-white flex-grow focus:ring-2 focus:ring-blue-500 outline-none"
          disabled={isLoading}
          pattern="[a-zA-Z0-9_-]{3,20}"
          title="3-20 characters (letters, numbers, -, _)"
        />
        <button
          onClick={handleFetch}
          disabled={isLoading}
          className="bg-blue-600 hover:bg-blue-500 disabled:bg-gray-600 text-white font-semibold px-4 py-2 rounded transition-colors"
        >
          {isLoading ? 'Fetching...' : 'Fetch Latest Game'}
        </button>
        {isLoading && (
          <button
            onClick={handleCancel}
            className="bg-red-600 hover:bg-red-500 text-white font-semibold px-4 py-2 rounded transition-colors"
          >
            Cancel
          </button>
        )}
      </div>

      {isLoading && fetchProgress && (
        <div className="text-blue-400 text-sm">
          {fetchProgress}
        </div>
      )}

      {error && (
        <div className="bg-red-900/50 border border-red-500 p-3 rounded text-red-200">
          Error: {error}
        </div>
      )}

      {gamePreview && (
        <div className="bg-gray-800 p-4 rounded space-y-2">
          <h4 className="font-bold text-white">Game Preview</h4>
          <p><span className="text-gray-400">White:</span> {gamePreview.white} {gamePreview.whiteElo ? `(${gamePreview.whiteElo})` : ''}</p>
          <p><span className="text-gray-400">Black:</span> {gamePreview.black} {gamePreview.blackElo ? `(${gamePreview.blackElo})` : ''}</p>
          <p><span className="text-gray-400">Result:</span> {gamePreview.result}</p>
          <p><span className="text-gray-400">Date:</span> {gamePreview.date}</p>
          <p><span className="text-gray-400">Time Control:</span> {gamePreview.timeControl}</p>
          <p><span className="text-gray-400">Rated:</span> {gamePreview.rated ? 'Yes' : 'No'}</p>
          <button
            onClick={handleLoad}
            className="bg-green-600 hover:bg-green-500 p-2 rounded text-white w-full font-semibold transition-colors mt-4"
          >
            Load Game
          </button>
        </div>
      )}
    </div>
  );
};

export default GameImporter;
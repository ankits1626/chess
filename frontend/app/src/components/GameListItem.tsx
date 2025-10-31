import type { FC } from 'react';
import type { ChessComGame } from '../services/chesscomApi';

interface GameListItemProps {
  game: ChessComGame;
  searchedUsername: string;
  onSelect: () => void;
}

const GameListItem: FC<GameListItemProps> = ({ game, searchedUsername, onSelect }) => {
  // The PGN header contains the result from White's perspective.
  // We need to determine if the searched player was White or Black to show the correct result.
  const pgnResult = game.pgn.match(/Result "(.+?)"/)?.[1] || '*';
  let playerResult: string;

  const isPlayerWhite = game.white.username.toLowerCase() === searchedUsername.toLowerCase();

  if (pgnResult === '1-0') {
    playerResult = isPlayerWhite ? 'Win' : 'Loss';
  } else if (pgnResult === '0-1') {
    playerResult = isPlayerWhite ? 'Loss' : 'Win';
  } else if (pgnResult === '1/2-1/2') {
    playerResult = 'Draw';
  } else {
    playerResult = '*'; // Unknown result
  }

  const date = new Date(game.end_time * 1000).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });

  return (
    <div className="bg-gray-900 p-3 rounded-lg flex justify-between items-center transition-colors hover:bg-gray-800">
      <div>
        <div className="font-semibold">
          <span className={isPlayerWhite ? 'text-blue-400' : ''}>
            {game.white.username} ({game.white.rating})
          </span>
          <span className="text-gray-400"> vs </span>
          <span className={!isPlayerWhite ? 'text-blue-400' : ''}>
            {game.black.username} ({game.black.rating})
          </span>
        </div>
        <div className="text-sm text-gray-400">
          {playerResult} • {date} • {game.time_control}
        </div>
      </div>
      <button
        onClick={onSelect}
        className="px-4 py-2 bg-blue-600 text-white font-semibold rounded-md hover:bg-blue-500 transition-colors focus:outline-none focus:ring-2 focus:ring-blue-400 focus:ring-opacity-75"
      >
        Load
      </button>
    </div>
  );
};

export default GameListItem;

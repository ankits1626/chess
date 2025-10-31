import type { FC } from 'react';
import type { ChessComGame } from '@/services/chesscomApi';
import GameListItem from './GameListItem';

interface GameListProps {
  games: ChessComGame[];
  searchedUsername: string;
  onSelectGame: (game: ChessComGame) => void;
  isLoading: boolean;
}

const GameList: FC<GameListProps> = ({ games, searchedUsername, onSelectGame, isLoading }) => {
  if (isLoading) {
    return <div className="text-center py-8 text-gray-400">Loading games...</div>;
  }

  if (games.length === 0) {
    return <div className="text-center py-8 text-gray-400">No games found in this archive.</div>;
  }

  return (
    <div className="space-y-2 max-h-96 overflow-y-auto p-1">
      {games.map((game) => (
        <GameListItem
          key={game.url}
          game={game}
          searchedUsername={searchedUsername}
          onSelect={() => onSelectGame(game)}
        />
      ))}
    </div>
  );
};

export default GameList;

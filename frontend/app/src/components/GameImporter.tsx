import { useState, useEffect } from 'react';
import { useGameStore } from '../store/useGameStore';
import GameList from './GameList';
import ArchivePaginator from './ArchivePaginator';

const isValidUsername = (username: string): boolean => {
  return /^[a-zA-Z0-9_-]{3,20}$/.test(username);
};

const GameImporter = () => {
  const [username, setUsername] = useState('hikaru');
  const [searchedUsername, setSearchedUsername] = useState('');
  const [validationError, setValidationError] = useState<string | null>(null);

  // Get new state and actions from the store
  const {
    gameArchives,
    importedGames,
    isGameListLoading,
    fetchGamesError,
    currentArchiveUrl,
    fetchGameArchives,
    fetchGamesForArchive,
    loadPgnFromGame,
    resetGameImporter,
  } = useGameStore();

  // Reset on unmount
  useEffect(() => {
    return () => {
      resetGameImporter();
    };
  }, [resetGameImporter]);

  const handleSearch = () => {
    if (!isValidUsername(username)) {
      setValidationError('Invalid username format. Use 3-20 characters (letters, numbers, -, _).');
      return;
    }
    setValidationError(null);
    setSearchedUsername(username);
    fetchGameArchives(username);
  };

  const handleUsernameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUsername(e.target.value);
    setValidationError(null); // Clear validation error on change
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !isGameListLoading) {
      handleSearch();
    }
  };

  return (
    <div className="p-4 bg-gray-900 rounded-lg space-y-4 w-full">
      <h3 className="text-xl font-bold text-white">Import Games from Chess.com</h3>
      <div className="flex gap-2">
        <input
          type="text"
          value={username}
          onChange={handleUsernameChange}
          onKeyDown={handleKeyDown}
          placeholder="Chess.com username"
          className="bg-gray-700 p-2 rounded text-white flex-grow focus:ring-2 focus:ring-blue-500 outline-none"
          disabled={isGameListLoading}
        />
        <button
          onClick={handleSearch}
          disabled={isGameListLoading}
          className="bg-blue-600 hover:bg-blue-500 disabled:bg-gray-600 text-white font-semibold px-4 py-2 rounded transition-colors"
        >
          {isGameListLoading ? 'Searching...' : 'Search'}
        </button>
      </div>

      {validationError && (
        <div className="bg-yellow-900/50 border border-yellow-500 p-3 rounded text-yellow-200">
          {validationError}
        </div>
      )}

      {fetchGamesError && (
        <div className="bg-red-900/50 border border-red-500 p-3 rounded text-red-200">
          Error: {fetchGamesError}
        </div>
      )}

      {/* Show loading indicator only on initial search */}
      {isGameListLoading && !importedGames && !fetchGamesError && (
         <div className="text-center py-8 text-gray-400">Fetching archives...</div>
      )}

      {gameArchives && currentArchiveUrl && importedGames && (
        <div className="space-y-4">
          <ArchivePaginator
            archives={gameArchives}
            currentArchiveUrl={currentArchiveUrl}
            onSelectArchive={fetchGamesForArchive}
          />
          <GameList
            games={importedGames}
            searchedUsername={searchedUsername}
            onSelectGame={loadPgnFromGame}
            isLoading={isGameListLoading}
          />
        </div>
      )}
    </div>
  );
};

export default GameImporter;
# Step 12b: Chess.com API Service & Game Importer UI

**Goal:** To enable users to import a game from chess.com by providing a username, fetching their latest game, and displaying a preview before loading it onto the board.

**This step corresponds to Phase 2 of the feedback on the original Step 12 plan.**

---

## 1. Implementation Plan

### 1. Create Chess.com API Service (`src/services/chesscomApi.ts`)

This file will contain functions to interact with the public chess.com API.

**Key functions:**
*   `fetchUserArchives(username: string)`: Fetches a list of monthly archive URLs for a given username.
*   `fetchMonthGames(archiveUrl: string)`: Fetches all games from a specific monthly archive URL.
*   `fetchLatestGamePgn(username: string)`: Orchestrates the above two functions to fetch the PGN of the most recent game played by the user.

**Example Implementation:**

```typescript
// src/services/chesscomApi.ts

interface ChessComArchivesResponse {
  archives: string[]; // List of archive URLs
}

interface ChessComGamesResponse {
  games: Array<{
    url: string;
    pgn: string;
    time_control: string;
    end_time: number;
    rated: boolean;
    // ... other fields
  }>;
}

export const fetchUserArchives = async (username: string): Promise<ChessComArchivesResponse> => {
  const response = await fetch(
    `https://api.chess.com/pub/player/${username}/games/archives`
  );

  if (!response.ok) {
    throw new Error(`Failed to fetch archives for ${username}: ${response.statusText}`);
  }

  return response.json();
};

export const fetchMonthGames = async (archiveUrl: string): Promise<ChessComGamesResponse> => {
  const response = await fetch(archiveUrl);

  if (!response.ok) {
    throw new Error(`Failed to fetch games from archive: ${response.statusText}`);
  }

  return response.json();
};

export const fetchLatestGamePgn = async (username: string): Promise<string> => {
  const archivesResponse = await fetchUserArchives(username);

  if (archivesResponse.archives.length === 0) {
    throw new Error('No game archives found for this user.');
  }

  // Get the most recent month's archive URL
  const latestArchiveUrl = archivesResponse.archives[archivesResponse.archives.length - 1];
  const monthGamesResponse = await fetchMonthGames(latestArchiveUrl);

  if (monthGamesResponse.games.length === 0) {
    throw new Error('No games found in the latest archive.');
  }

  // Get the PGN of the most recent game in that archive
  const latestGame = monthGamesResponse.games[monthGamesResponse.games.length - 1];
  return latestGame.pgn;
};
```

### 2. Create `GameImporter.tsx` Component

This component will provide the UI for users to input a chess.com username, fetch their latest game, and preview it before loading.

**Key elements:**
*   **State:** `username`, `isLoading`, `error`, `gamePreview` (to store parsed game metadata and PGN).
*   **UI:** Input field for username, "Fetch Latest Game" button, loading indicator, error display, game preview section, "Load Game" button.
*   **Integration:** Uses `fetchLatestGamePgn` from the API service and `loadPgn` from the Zustand store.

**Example Implementation:**

```typescript
// src/components/GameImporter.tsx
import { useState } from 'react';
import { useGameStore } from '../store/useGameStore';
import { fetchLatestGamePgn } from '../services/chesscomApi';
import { Chess } from 'chess.js'; // To parse PGN for preview

interface GamePreview {
  white: string;
  black: string;
  result: string;
  date: string;
  pgn: string;
}

const GameImporter = () => {
  const [username, setUsername] = useState('hikaru'); // Default for easy testing
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [gamePreview, setGamePreview] = useState<GamePreview | null>(null);

  const loadPgn = useGameStore(state => state.loadPgn);

  const handleFetch = async () => {
    setIsLoading(true);
    setError(null);
    setGamePreview(null);

    try {
      const pgn = await fetchLatestGamePgn(username);

      // Parse PGN for preview metadata
      const tempGame = new Chess();
      tempGame.loadPgn(pgn);
      const header = tempGame.header();

      setGamePreview({
        white: header.White || 'Unknown',
        black: header.Black || 'Unknown',
        result: header.Result || '*' ,
        date: header.Date || '',
        pgn
      });
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch game');
    } finally {
      setIsLoading(false);
    }
  };

  const handleLoad = () => {
    if (gamePreview) {
      loadPgn(gamePreview.pgn);
      setGamePreview(null); // Clear preview after loading
    }
  };

  return (
    <div className="p-4 bg-gray-900 rounded-lg space-y-4">
      <h3 className="text-xl font-bold text-white">Import Game from Chess.com</h3>
      <div className="flex gap-2">
        <input
          type="text"
          value={username}
          onChange={e => setUsername(e.target.value)}
          placeholder="Chess.com username"
          className="bg-gray-700 p-2 rounded text-white flex-grow focus:ring-2 focus:ring-blue-500 outline-none"
          disabled={isLoading}
        />
        <button
          onClick={handleFetch}
          disabled={isLoading}
          className="bg-blue-600 hover:bg-blue-500 disabled:bg-gray-600 text-white font-semibold px-4 py-2 rounded transition-colors"
        >
          {isLoading ? 'Fetching...' : 'Fetch Latest Game'}
        </button>
      </div>

      {error && (
        <div className="bg-red-900/50 border border-red-500 p-3 rounded text-red-200">
          Error: {error}
        </div>
      )}

      {gamePreview && (
        <div className="bg-gray-800 p-4 rounded space-y-2">
          <h4 className="font-bold text-white">Game Preview</h4>
          <p><span className="text-gray-400">White:</span> {gamePreview.white}</p>
          <p><span className="text-gray-400">Black:</span> {gamePreview.black}</p>
          <p><span className="text-gray-400">Result:</span> {gamePreview.result}</p>
          <p><span className="text-gray-400">Date:</span> {gamePreview.date}</p>
          <button
            onClick={handleLoad}
            className="bg-green-600 hover:bg-green-500 p-2 rounded text-white w-full font-semibold transition-colors"
          >
            Load Game
          </button>
        </div>
      )}
    </div>
  );
};

export default GameImporter;
```

### 3. Integrate `GameImporter.tsx` into `App.tsx`

Add the `GameImporter` component to the main application layout, likely in the `GameInfo` panel area or as a separate section.

```typescript
// App.tsx (simplified)
import GameImporter from './components/GameImporter';

function App() {
  // ... existing code ...

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-800 text-white p-8">
      <div className="flex flex-col lg:flex-row gap-4 lg:gap-8 items-start">
        <GameBoard />
        <GameInfo />
        {/* Add GameImporter here */}
        <GameImporter />
      </div>
      {/* ... PromotionDialog and DebugPanel ... */}
    </div>
  );
}
```

---

## 2. Testing Checklist

1.  **Successful Fetch & Preview:**
    *   ✅ Enter a valid chess.com username (e.g., "hikaru").
    *   ✅ Click "Fetch Latest Game".
    *   ✅ Verify "Fetching..." indicator appears.
    *   ✅ Verify game preview (White, Black, Result, Date) appears correctly.

2.  **Load Game:**
    *   ✅ Click "Load Game" in the preview.
    *   ✅ Verify the board updates to the final position of the imported game.
    *   ✅ Verify the move history panel shows the moves of the imported game.

3.  **Error Handling:**
    *   ✅ Enter an invalid/non-existent username.
    *   ✅ Click "Fetch Latest Game".
    *   ✅ Verify an error message is displayed (e.g., "Failed to fetch archives...").
    *   ✅ Enter a username with no games.
    *   ✅ Verify an appropriate error message is displayed (e.g., "No game archives found...").

4.  **Loading State:**
    *   ✅ Verify buttons are disabled and input is disabled while fetching.

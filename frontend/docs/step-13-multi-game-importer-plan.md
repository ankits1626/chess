### **Plan: Multi-Game Importer (Revised)**

This document outlines the plan to implement a multi-game importer, incorporating feedback from the initial review.

#### **1. Overview**

The goal is to allow users to enter a Chess.com username, browse a paginated list of their game archives, and select any game to load into the replay viewer.

**Phases:**
1.  **State Management:** Update the Zustand store to handle the new state and actions required for fetching and displaying game lists.
2.  **UI Components:** Create new React components for the game list, list items, and archive pagination.
3.  **Integration:** Combine the state and UI into a seamless user experience within the existing `GameImporter` component.

#### **2. Phase 1: State & API Management**

1.  **API Service (`chesscomApi.ts`):**
    *   The plan will use the **existing** API functions: `fetchUserArchives` and `fetchMonthGames`.
    *   The archive URLs have a format like `https://api.chess.com/pub/player/{username}/games/YYYY/MM`.
    *   A new `ChessComGame` type will be exported from this file for use in the store and components: `export type ChessComGame = ChessComGamesResponse['games'][0];`

2.  **Zustand Store (`useGameStore.ts`):**
    *   **New State:**
        *   `gameArchives: string[] | null`: Stores the list of monthly archive URLs.
        *   `importedGames: ChessComGame[] | null`: Stores games from the selected archive, using the proper type.
        *   `isGameListLoading: boolean`: Manages loading states for all fetch operations.
        *   `fetchGamesError: string | null`: Stores error messages.
        *   `currentArchiveUrl: string | null`: Tracks the currently viewed archive URL.
    *   **New Actions:**
        *   `fetchGameArchives(username: string)`: Fetches archives and then automatically fetches games from the latest month. Will include robust error handling and loading state management.
        *   `fetchGamesForArchive(url: string)`: Fetches games for a specific archive URL, handling loading and error states.
        *   `loadPgnFromGame(game: ChessComGame)`: Extracts the PGN from a selected game and calls the existing `loadPgn` action.
        *   `resetGameImporter()`: Clears all importer-related state (archives, games, errors) when starting a new search.

#### **3. Phase 2: UI Components**

**Component Hierarchy:**
```
GameImporter (container)
├── Username Input + Search Button
├── Loading State / Error Message
└── (If data loaded)
    ├── ArchivePaginator
    │   ├── "← Older" Button
    │   ├── Current Month Display
    │   └── "Newer →" Button
    └── GameList
        └── GameListItem (repeated)
            ├── Game Metadata
            └── Load Button
```

1.  **`GameImporter.tsx` (Enhancement):**
    *   Will be refactored to orchestrate the new components and state, including holding the `searchedUsername` in its local state.

2.  **`GameList.tsx` (New):**
    *   **Props:** `games: ChessComGame[]`, `searchedUsername: string`, `onSelectGame: (game: ChessComGame) => void`, `isLoading: boolean`.
    *   Displays a loading spinner or a list of `GameListItem`s.

3.  **`GameListItem.tsx` (New):**
    *   **Props:** `game: ChessComGame`, `searchedUsername: string`, `onSelect: () => void`.
    *   **Display Format:** Will show both players' usernames and ratings, the game result from the perspective of the `searchedUsername`, the date, and the time control. The searched player's name will be highlighted.

4.  **`ArchivePaginator.tsx` (New):**
    *   **Props:** `archives: string[]`, `currentArchive: string`, `onSelectArchive: (url: string) => void`.
    *   **Logic:** Will have "← Older Games" and "Newer Games →" buttons. The logic will correctly navigate the chronologically sorted `archives` array. For example, if viewing the last item, "Newer Games" will be disabled.

#### **4. Phase 3: Integration & User Flow**

1.  User enters a username in `GameImporter` and clicks "Search".
2.  `fetchGameArchives` is called. A loading state is shown.
3.  On success, games from the latest archive are fetched and displayed in `GameList`.
4.  `ArchivePaginator` displays the current month.
5.  User can click "Load" on a game, which calls `loadPgnFromGame`.
6.  User can navigate to older/newer months using the paginator, which calls `fetchGamesForArchive` and updates the `GameList`.
7.  If any fetch fails, `fetchGamesError` is set and an error message is displayed.
8.  If a monthly archive contains no games, the `GameList` will display a "No games found for this month" message.

#### **5. Testing Checklist**

- [ ] Search for a valid username loads archives and the latest month's games.
- [ ] Search for an invalid or user with no games displays an appropriate error message.
- [ ] An archive with zero games shows the correct message.
- [ ] Game list and pagination controls display correctly.
- [ ] Pagination buttons are correctly enabled/disabled at the boundaries of the archive list.
- [ ] Clicking "Load" on a game clears the importer UI and loads the game into the replay viewer.
- [ ] Loading spinners are shown during all fetch operations.
- [ ] Network errors are handled gracefully and a message is shown.
- [ ] The searched player is correctly highlighted in the `GameListItem`.

#### **6. Future Considerations**

*   **Rate Limiting:** The Chess.com public API has rate limits. If 429 "Too Many Requests" errors become an issue, a retry-with-exponential-backoff strategy could be implemented in the API service functions. For the initial implementation, this will be omitted.
*   **UI Enhancements:** Further improvements could include filtering games (by result, opponent), sorting, and a search history for recently viewed players.

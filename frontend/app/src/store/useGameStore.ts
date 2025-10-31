import { create } from 'zustand';
import { Chess, Move } from 'chess.js';
import type { Square, PieceType, PieceColor } from '../types/chess';
import {
  fetchUserArchives,
  fetchMonthGames,
  type ChessComGame,
} from '../services/chesscomApi';

// Type definitions
export type LastMove = { from: Square; to: Square } | null;
type PendingPromotion = { from: Square; to: Square; color: PieceColor };

export interface DebugActions {
  loadFen: (fen: string) => void;
  resetGame: () => void;
}

// The complete state of the game application
interface GameState {
  game: Chess;
  selectedSquare: Square | null;
  validMoves: Square[];
  pendingMove: PendingPromotion | null;
  lastMove: LastMove;
  debugActions?: DebugActions;

  // Replay state
  mode: 'live' | 'replay';
  replayMoves: Move[];
  replayIndex: number;
  isAutoplaying: boolean;
  autoplayIntervalId: NodeJS.Timeout | null;
  whitePlayer: string | null;
  blackPlayer: string | null;

  // Multi-game importer state
  gameArchives: string[] | null;
  importedGames: ChessComGame[] | null;
  isGameListLoading: boolean;
  fetchGamesError: string | null;
  currentArchiveUrl: string | null;

  // Actions
  selectSquare: (square: Square) => void;
  handlePromotion: (piece: PieceType) => void;
  resetGame: () => void;
  loadPgn: (pgn: string) => void;

  // Replay Actions
  goToMove: (index: number) => void;
  nextMove: () => void;
  prevMove: () => void;
  goToFirstMove: () => void;
  goToLastMove: () => void;
  toggleAutoplay: () => void;
  startAutoplay: () => void;
  stopAutoplay: () => void;

  // Multi-game importer actions
  fetchGameArchives: (username: string) => Promise<void>;
  fetchGamesForArchive: (url: string) => Promise<void>;
  loadPgnFromGame: (game: ChessComGame) => void;
  resetGameImporter: () => void;
}

export const useGameStore = create<GameState>((set, get) => ({
  // Initial State
  game: new Chess(),
  selectedSquare: null,
  validMoves: [],
  pendingMove: null,
  lastMove: null,
  mode: 'live',
  replayMoves: [],
  replayIndex: -1,
  isAutoplaying: false,
  autoplayIntervalId: null,
  whitePlayer: null,
  blackPlayer: null,
  gameArchives: null,
  importedGames: null,
  isGameListLoading: false,
  fetchGamesError: null,
  currentArchiveUrl: null,

  // Actions
  resetGame: () => {
    get().stopAutoplay();
    get().resetGameImporter();
    set({
      game: new Chess(),
      selectedSquare: null,
      validMoves: [],
      pendingMove: null,
      lastMove: null,
      mode: 'live',
      replayMoves: [],
      replayIndex: -1,
      whitePlayer: null,
      blackPlayer: null,
    });
  },

  loadPgn: (pgn: string) => {
    get().stopAutoplay();
    get().resetGameImporter();
    try {
      const whiteMatch = pgn.match(/\s*\[\s*White\s*"(.*?)"\s*\]\s*/);
      const blackMatch = pgn.match(/\s*\[\s*Black\s*"(.*?)"\s*\]\s*/);
      const whitePlayer = whiteMatch ? whiteMatch[1] : null;
      const blackPlayer = blackMatch ? blackMatch[1] : null;

      let cleanedPgn = pgn.replace(/\{[^}]*\}/g, '').replace(/\([^)]*\)/g, '');
      const moveTextMatch = cleanedPgn.match(/\n\n(.+)$/s);
      if (!moveTextMatch) throw new Error('No moves found in PGN');
      let moveText = moveTextMatch[1];

      const sanMovePattern = /([NBRQK]?[a-h]?[1-8]?x?[a-h][1-8](?:=[NBRQ])?[+#]?|O-O(?:-O)?)/g;
      const sanMoves = moveText.match(sanMovePattern);
      if (!sanMoves || sanMoves.length === 0) throw new Error('No valid moves found in PGN');

      const tempGame = new Chess();
      const moves: Move[] = sanMoves.map(san => {
        const move = tempGame.move(san.trim());
        if (!move) throw new Error(`Invalid move: ${san}`);
        return move;
      });

      const replayGame = new Chess();
      set({
        mode: 'replay',
        game: Object.assign(Object.create(Object.getPrototypeOf(replayGame)), replayGame),
        replayMoves: moves,
        replayIndex: -1,
        whitePlayer,
        blackPlayer,
        selectedSquare: null,
        validMoves: [],
        pendingMove: null,
        lastMove: null,
      });
    } catch (error) {
      console.error('Failed to load PGN:', error);
      throw error;
    }
  },

  goToMove: (index: number) => {
    const { replayMoves, stopAutoplay } = get();
    stopAutoplay();
    if (index < -1 || index >= replayMoves.length) return;

    const tempGame = new Chess();
    for (let i = 0; i <= index; i++) {
      tempGame.move(replayMoves[i].san);
    }

    set({ game: tempGame, replayIndex: index, selectedSquare: null, validMoves: [], pendingMove: null, lastMove: null });
  },

  nextMove: () => {
    const { isAutoplaying, stopAutoplay, replayIndex, replayMoves } = get();
    if (!isAutoplaying) stopAutoplay();
    if (replayIndex < replayMoves.length - 1) {
      const newIndex = replayIndex + 1;
      const tempGame = new Chess();
      for (let i = 0; i <= newIndex; i++) {
        tempGame.move(replayMoves[i].san);
      }
      set({ game: tempGame, replayIndex: newIndex, selectedSquare: null, validMoves: [], pendingMove: null, lastMove: null });
    }
  },

  prevMove: () => {
    const { replayIndex, goToMove } = get();
    if (replayIndex > -1) goToMove(replayIndex - 1);
  },

  goToFirstMove: () => get().goToMove(-1),
  goToLastMove: () => get().goToMove(get().replayMoves.length - 1),

  toggleAutoplay: () => {
    const { isAutoplaying, stopAutoplay, startAutoplay } = get();
    if (isAutoplaying) stopAutoplay();
    else startAutoplay();
  },

  startAutoplay: () => {
    const { replayIndex, replayMoves, nextMove } = get();
    if (replayIndex >= replayMoves.length - 1) return;
    const intervalId = setInterval(() => {
      const { replayIndex: currentIndex, replayMoves: currentMoves, stopAutoplay: currentStop } = get();
      if (currentIndex >= currentMoves.length - 1) {
        currentStop();
        return;
      }
      nextMove();
    }, 1000);
    set({ isAutoplaying: true, autoplayIntervalId: intervalId });
  },

  stopAutoplay: () => {
    const { autoplayIntervalId } = get();
    if (autoplayIntervalId) clearInterval(autoplayIntervalId);
    set({ isAutoplaying: false, autoplayIntervalId: null });
  },

  // Importer Actions
  fetchGameArchives: async (username: string) => {
    set({ isGameListLoading: true, fetchGamesError: null, gameArchives: null, importedGames: null, currentArchiveUrl: null });
    try {
      const archivesResponse = await fetchUserArchives(username);
      if (archivesResponse.archives.length === 0) throw new Error(`User "${username}" has no game archives`);
      const archives = archivesResponse.archives;
      set({ gameArchives: archives });
      const latestArchiveUrl = archives[archives.length - 1];
      await get().fetchGamesForArchive(latestArchiveUrl);
    } catch (error) {
      set({ fetchGamesError: error instanceof Error ? error.message : 'Failed to fetch archives', isGameListLoading: false });
    }
  },

  fetchGamesForArchive: async (url: string) => {
    set({ isGameListLoading: true, fetchGamesError: null, currentArchiveUrl: url });
    try {
      const gamesResponse = await fetchMonthGames(url);
      // The API returns games oldest first, so we reverse them to show newest first.
      const reversedGames = gamesResponse.games.reverse();
      set({ importedGames: reversedGames, isGameListLoading: false });
    } catch (error) {
      set({ fetchGamesError: error instanceof Error ? error.message : 'Failed to fetch games', isGameListLoading: false });
    }
  },

  loadPgnFromGame: (game: ChessComGame) => {
    get().loadPgn(game.pgn);
  },

  resetGameImporter: () => {
    set({ gameArchives: null, importedGames: null, isGameListLoading: false, fetchGamesError: null, currentArchiveUrl: null });
  },

  selectSquare: (square: Square) => {
    const { game, pendingMove, selectedSquare, mode } = get();
    if (mode === 'replay') return;
    if (pendingMove) return;
    if (!selectedSquare) {
      const piece = game.get(square);
      if (piece && piece.color === game.turn()) {
        const moves = game.moves({ square, verbose: true });
        set({ selectedSquare: square, validMoves: moves.map(move => move.to) });
      }
      return;
    }
    if (selectedSquare === square) {
      set({ selectedSquare: null, validMoves: [] });
      return;
    }
    try {
      const piece = game.get(selectedSquare);
      if (piece?.type === 'p' && (square.endsWith('1') || square.endsWith('8'))) {
        const moves = game.moves({ square: selectedSquare, verbose: true });
        const isValidMove = moves.some(m => m.to === square);
        if (isValidMove) {
          set({ pendingMove: { from: selectedSquare, to: square, color: piece.color }, selectedSquare: null, validMoves: [] });
          return;
        }
      }
      const move = game.move({ from: selectedSquare, to: square });
      if (move) {
        set({ game: Object.assign(Object.create(Object.getPrototypeOf(game)), game), lastMove: { from: selectedSquare, to: square }, selectedSquare: null, validMoves: [] });
      } else {
        const newPiece = game.get(square);
        if (newPiece && newPiece.color === game.turn()) {
          const moves = game.moves({ square, verbose: true });
          set({ selectedSquare: square, validMoves: moves.map(m => m.to) });
        } else {
          set({ selectedSquare: null, validMoves: [] });
        }
      }
    } catch {
      set({ selectedSquare: null, validMoves: [] });
    }
  },

  handlePromotion: (piece: PieceType) => {
    const { game, pendingMove, mode } = get();
    if (mode === 'replay') return;
    if (!pendingMove) return;
    game.move({ from: pendingMove.from, to: pendingMove.to, promotion: piece });
    set({ game: Object.assign(Object.create(Object.getPrototypeOf(game)), game), lastMove: { from: pendingMove.from, to: pendingMove.to }, pendingMove: null });
  },

  // Debug Actions
  debugActions: import.meta.env.DEV ? {
    loadFen: (fen: string) => {
      get().stopAutoplay();
      get().resetGameImporter();
      try {
        const newGame = new Chess(fen);
        set({ game: newGame, selectedSquare: null, validMoves: [], pendingMove: null, lastMove: null, mode: 'live', replayMoves: [], replayIndex: -1, whitePlayer: null, blackPlayer: null });
      } catch (error) {
        console.error('Invalid FEN:', error);
      }
    },
    resetGame: () => get().resetGame(),
  } : undefined,
}));


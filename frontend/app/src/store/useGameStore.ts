import { create } from 'zustand';
import { Chess, Move } from 'chess.js';
import type { Square, PieceType, PieceColor } from '../types/chess';

// Type definitions from the old GameController
export type LastMove = {
  from: Square;
  to: Square;
} | null;

type PendingPromotion = {
  from: Square;
  to: Square;
  color: PieceColor;
};

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
  replayMoves: Move[]; // `Move` type from chess.js
  replayIndex: number; // -1 for initial position, 0 for first move, etc.
  isAutoplaying: boolean;
  autoplayIntervalId: NodeJS.Timeout | null;
  whitePlayer: string | null;
  blackPlayer: string | null;

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
  replayIndex: -1, // -1 means initial board state
  isAutoplaying: false,
  autoplayIntervalId: null,
  whitePlayer: null,
  blackPlayer: null,

  // Actions
  resetGame: () => {
    get().stopAutoplay();
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
    try {
      // Parse player names from PGN header
      const whiteMatch = pgn.match(/\s*\[\s*White\s*"(.*?)"\s*\]\s*/);
      const blackMatch = pgn.match(/\s*\[\s*Black\s*"(.*?)"\s*\]\s*/);
      const whitePlayer = whiteMatch ? whiteMatch[1] : null;
      const blackPlayer = blackMatch ? blackMatch[1] : null;

      // Remove all comments and annotations
      let cleanedPgn = pgn
        .replace(/\{[^}]*\}/g, '')  // Remove {comments}
        .replace(/\([^)]*\)/g, ''); // Remove (variations)

      // Extract move text (everything after headers)
      const moveTextMatch = cleanedPgn.match(/\n\n(.+)$/s);
      if (!moveTextMatch) {
        throw new Error('No moves found in PGN');
      }

      let moveText = moveTextMatch[1];

      // Extract only the actual chess moves in SAN notation
      const sanMovePattern = /([NBRQK]?[a-h]?[1-8]?x?[a-h][1-8](?:=[NBRQ])?[+#]?|O-O(?:-O)?)/g;
      const sanMoves = moveText.match(sanMovePattern);

      if (!sanMoves || sanMoves.length === 0) {
        throw new Error('No valid moves found in PGN');
      }

      // Replay moves manually to build the move history
      const tempGame = new Chess();
      const moves: Move[] = [];

      for (let i = 0; i < sanMoves.length; i++) {
        const san = sanMoves[i].trim();
        try {
          const move = tempGame.move(san);
          if (!move) {
            throw new Error(`Invalid move: ${san} at position ${i + 1}`);
          }
          moves.push(move);
        } catch (error) {
          console.error(`Error playing move ${i + 1}: ${san}`, error);
          throw error;
        }
      }

      // Reset to starting position for replay
      const replayGame = new Chess();

      set({
        mode: 'replay',
        game: Object.assign(Object.create(Object.getPrototypeOf(replayGame)), replayGame),
        replayMoves: moves,
        replayIndex: -1,
        whitePlayer,
        blackPlayer,

        // Clear live game state
        selectedSquare: null,
        validMoves: [],
        pendingMove: null,
        lastMove: null
      });
    } catch (error) {
      console.error('Failed to load PGN:', error);
      throw error; // Re-throw for component to handle
    }
  },

  goToMove: (index: number) => {
    const { replayMoves, stopAutoplay } = get();
    stopAutoplay();

    // Ensure index is within bounds
    if (index < -1 || index >= replayMoves.length) {
      console.warn(`Attempted to go to invalid move index: ${index}`);
      return;
    }

    const tempGame = new Chess();
    // Replay moves up to the specified index
    for (let i = 0; i <= index; i++) {
      tempGame.move(replayMoves[i].san); // Use .san property for moves
    }

    set({
      game: Object.assign(Object.create(Object.getPrototypeOf(tempGame)), tempGame),
      replayIndex: index,
      // Clear live game state when navigating in replay mode
      selectedSquare: null,
      validMoves: [],
      pendingMove: null,
      lastMove: null,
    });
  },

  nextMove: () => {
    const { isAutoplaying, stopAutoplay, replayIndex, replayMoves } = get();
    if (!isAutoplaying) {
      stopAutoplay();
    }
    if (replayIndex < replayMoves.length - 1) {
      // Manually call the core logic of goToMove without the stopAutoplay side-effect
      const newIndex = replayIndex + 1;
      const tempGame = new Chess();
      for (let i = 0; i <= newIndex; i++) {
        tempGame.move(replayMoves[i].san);
      }
      set({
        game: Object.assign(Object.create(Object.getPrototypeOf(tempGame)), tempGame),
        replayIndex: newIndex,
        selectedSquare: null,
        validMoves: [],
        pendingMove: null,
        lastMove: null,
      });
    }
  },

  prevMove: () => {
    const { replayIndex, goToMove } = get();
    // goToMove already stops autoplay
    if (replayIndex > -1) {
      goToMove(replayIndex - 1);
    }
  },

  goToFirstMove: () => {
    get().goToMove(-1); // Go to initial board state
  },

  goToLastMove: () => {
    const { replayMoves, goToMove } = get();
    goToMove(replayMoves.length - 1);
  },

  toggleAutoplay: () => {
    const { isAutoplaying, stopAutoplay, startAutoplay } = get();
    if (isAutoplaying) {
      stopAutoplay();
    } else {
      startAutoplay();
    }
  },

  startAutoplay: () => {
    const { replayIndex, replayMoves, nextMove, stopAutoplay } = get();

    if (replayIndex >= replayMoves.length - 1) {
      return; // Don't start if already at the end
    }

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
    if (autoplayIntervalId) {
      clearInterval(autoplayIntervalId);
    }
    set({ isAutoplaying: false, autoplayIntervalId: null });
  },

  selectSquare: (square: Square) => {
    const { game, pendingMove, selectedSquare, mode } = get();

    // Disable clicks in replay mode
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
          set({
            pendingMove: { from: selectedSquare, to: square, color: piece.color },
            selectedSquare: null,
            validMoves: [],
          });
          return;
        }
      }

      const move = game.move({ from: selectedSquare, to: square });

      if (move) {
        set({
          game: Object.assign(Object.create(Object.getPrototypeOf(game)), game),
          lastMove: { from: selectedSquare, to: square },
          selectedSquare: null,
          validMoves: [],
        });
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
    if (mode === 'replay') return; // Disable in replay mode
    if (!pendingMove) return;

    game.move({ from: pendingMove.from, to: pendingMove.to, promotion: piece });
    set({
      game: Object.assign(Object.create(Object.getPrototypeOf(game)), game),
      lastMove: { from: pendingMove.from, to: pendingMove.to },
      pendingMove: null,
    });
  },

  // Debug Actions
  debugActions: import.meta.env.DEV ? {
    loadFen: (fen: string) => {
      get().stopAutoplay();
      try {
        const newGame = new Chess(fen);
        set({
          game: newGame,
          selectedSquare: null,
          validMoves: [],
          pendingMove: null,
          lastMove: null,
          mode: 'live',
          replayMoves: [],
          replayIndex: -1,
        });
      } catch (error) {
        console.error('Invalid FEN:', error);
      }
    },
    resetGame: () => get().resetGame(),
  } : undefined,
}));


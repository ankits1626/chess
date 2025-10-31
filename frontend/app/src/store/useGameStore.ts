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

  // Actions
  selectSquare: (square: Square) => void;
  handlePromotion: (piece: PieceType) => void;
  resetGame: () => void;
  loadPgn: (pgn: string) => void;

  // New Replay Actions
  goToMove: (index: number) => void;
  nextMove: () => void;
  prevMove: () => void;
  goToFirstMove: () => void;
  goToLastMove: () => void;
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

  // Actions
  resetGame: () => {
    set({
      game: new Chess(),
      selectedSquare: null,
      validMoves: [],
      pendingMove: null,
      lastMove: null,
      mode: 'live',
      replayMoves: [],
      replayIndex: -1,
    });
  },

  loadPgn: (pgn: string) => {
    try {
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
      // This regex matches: Nf3, e4, O-O, O-O-O, Bxf3+, exd5#, e8=Q, etc.
      const sanMovePattern = /([NBRQK]?[a-h]?[1-8]?x?[a-h][1-8](?:=[NBRQ])?[+#]?|O-O(?:-O)?)/g;
      const sanMoves = moveText.match(sanMovePattern);

      if (!sanMoves || sanMoves.length === 0) {
        throw new Error('No valid moves found in PGN');
      }

      console.log(`Extracted ${sanMoves.length} moves from PGN`);

      // Replay moves manually to build the move history
      const tempGame = new Chess();
      const moves: Move[] = [];

      for (let i = 0; i < sanMoves.length; i++) {
        const san = sanMoves[i].trim();
        try {
          const move = tempGame.move(san);
          if (!move) {
            console.error(`Failed to play move ${i + 1}: ${san}`);
            console.error('Current position:', tempGame.fen());
            throw new Error(`Invalid move: ${san} at position ${i + 1}`);
          }
          moves.push(move);
        } catch (error) {
          console.error(`Error playing move ${i + 1}: ${san}`, error);
          throw error;
        }
      }

      console.log(`Successfully loaded ${moves.length} moves`);

      // Reset to starting position for replay
      const replayGame = new Chess();

      set({
        mode: 'replay',
        game: Object.assign(Object.create(Object.getPrototypeOf(replayGame)), replayGame),
        replayMoves: moves,
        replayIndex: -1,

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
    const { replayMoves } = get();

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
    const { replayIndex, replayMoves, goToMove } = get();
    if (replayIndex < replayMoves.length - 1) {
      goToMove(replayIndex + 1);
    }
  },

  prevMove: () => {
    const { replayIndex, goToMove } = get();
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


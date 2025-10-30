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

  // Actions
  selectSquare: (square: Square) => void;
  handlePromotion: (piece: PieceType) => void;
  resetGame: () => void;
}

export const useGameStore = create<GameState>((set, get) => ({
  // Initial State
  game: new Chess(),
  selectedSquare: null,
  validMoves: [],
  pendingMove: null,
  lastMove: null,

  // Actions
  resetGame: () => {
    set({
      game: new Chess(),
      selectedSquare: null,
      validMoves: [],
      pendingMove: null,
      lastMove: null,
    });
  },

  selectSquare: (square: Square) => {
    const { game, pendingMove, selectedSquare } = get();

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
        // The `game` object is mutated by `move()`. We create a new object reference
        // to trigger a re-render in components that subscribe to the `game` state.
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
    const { game, pendingMove } = get();
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
          lastMove: null 
        });
      } catch (error) {
        console.error('Invalid FEN:', error);
      }
    },
    resetGame: () => get().resetGame(),
  } : undefined,
}));

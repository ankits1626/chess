import { useState } from 'react';
import { Chess } from 'chess.js';
import type { Square, PieceType, PieceColor } from '../types/chess';

type PendingPromotion = {
  from: Square;
  to: Square;
  color: PieceColor;
};

// Define the type for debug actions, which will only be available in dev mode.
export interface DebugActions {
  loadFen: (fen: string) => void;
  resetGame: () => void;
}

interface GameControllerProps {
  children: (
    game: Chess,
    selectedSquare: Square | null,
    validMoves: Square[],
    selectSquare: (square: Square) => void,
    pendingMove: PendingPromotion | null,
    handlePromotion: (piece: PieceType) => void,
    debugActions?: DebugActions
  ) => React.ReactNode;
}

const GameController = ({ children }: GameControllerProps) => {
  const [game, setGame] = useState(() => new Chess());
  const [selectedSquare, setSelectedSquare] = useState<Square | null>(null);
  const [validMoves, setValidMoves] = useState<Square[]>([]);
  const [pendingMove, setPendingMove] = useState<PendingPromotion | null>(null);

  // Debug actions are only created in development mode.
  const debugActions: DebugActions | undefined = import.meta.env.DEV ? {
    loadFen: (fen: string) => {
      try {
        const newGame = new Chess(fen);
        setGame(newGame);
        setSelectedSquare(null);
        setValidMoves([]);
        setPendingMove(null);
      } catch (error) {
        console.error('Invalid FEN:', error);
      }
    },
    resetGame: () => {
      setGame(new Chess());
      setSelectedSquare(null);
      setValidMoves([]);
      setPendingMove(null);
    },
  } : undefined;

  const selectSquare = (square: Square) => {
    // If a promotion is pending, don't allow other moves
    if (pendingMove) return;

    // If no square is selected, select this square (if it has a piece)
    if (!selectedSquare) {
      const piece = game.get(square);
      if (piece && piece.color === game.turn()) {
        setSelectedSquare(square);
        const moves = game.moves({ square, verbose: true });
        setValidMoves(moves.map(move => move.to));
      }
      return;
    }

    // If clicking the same square, deselect it
    if (selectedSquare === square) {
      setSelectedSquare(null);
      setValidMoves([]);
      return;
    }

    // Try to make a move
    try {
      const piece = game.get(selectedSquare);

      // Check if this would be a promotion move
      if (piece?.type === 'p' && (square.endsWith('1') || square.endsWith('8'))) {
        // Validate the move is legal
        const moves = game.moves({ square: selectedSquare, verbose: true });
        const isValidMove = moves.some(m => m.to === square);

        if (isValidMove) {
          // Valid promotion - show dialog
          setPendingMove({
            from: selectedSquare,
            to: square,
            color: piece.color
          });
          setSelectedSquare(null);
          setValidMoves([]);
          return;
        }
      }

      // Not a promotion, proceed with normal move
      const move = game.move({
        from: selectedSquare,
        to: square
      });

      if (move) {
        // Move was successful, update state
        setGame(new Chess(game.fen()));
        setSelectedSquare(null);
        setValidMoves([]);
      } else {
        // Invalid move, check if clicking another piece of the same color
        const newPiece = game.get(square);
        if (newPiece && newPiece.color === game.turn()) {
          setSelectedSquare(square);
          const moves = game.moves({ square, verbose: true });
          setValidMoves(moves.map(move => move.to));
        } else {
          setSelectedSquare(null);
          setValidMoves([]);
        }
      }
    } catch (error) {
      // Invalid move, deselect
      setSelectedSquare(null);
      setValidMoves([]);
    }
  };

  const handlePromotion = (piece: PieceType) => {
    if (!pendingMove) return;

    game.move({
      from: pendingMove.from,
      to: pendingMove.to,
      promotion: piece
    });
    setGame(new Chess(game.fen())); // Consistent with existing pattern
    setPendingMove(null);
  };

  return <>{children(game, selectedSquare, validMoves, selectSquare, pendingMove, handlePromotion, debugActions)}</>;
};

export default GameController;
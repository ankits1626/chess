import { useState } from 'react';
import { Chess } from 'chess.js';
import type { Square } from '../types/chess';

interface GameControllerProps {
  children: (
    game: Chess,
    selectedSquare: Square | null,
    validMoves: Square[],
    selectSquare: (square: Square) => void
  ) => React.ReactNode;
}

const GameController = ({ children }: GameControllerProps) => {
  const [game, setGame] = useState(() => new Chess());
  const [selectedSquare, setSelectedSquare] = useState<Square | null>(null);
  const [validMoves, setValidMoves] = useState<Square[]>([]);

  const selectSquare = (square: Square) => {
    // If no square is selected, select this square (if it has a piece)
    if (!selectedSquare) {
      const piece = game.get(square);
      if (piece && piece.color === game.turn()) {
        setSelectedSquare(square);
        // Calculate valid moves for this piece
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
      const move = game.move({
        from: selectedSquare,
        to: square,
        promotion: 'q' // Always promote to queen for now (Step 10 will add dialog)
      });

      if (move) {
        // Move was successful, create a new object reference to trigger re-render
        setGame(Object.assign(Object.create(Object.getPrototypeOf(game)), game));
        setSelectedSquare(null);
        setValidMoves([]);
      } else {
        // Invalid move, check if clicking another piece of the same color
        const piece = game.get(square);
        if (piece && piece.color === game.turn()) {
          setSelectedSquare(square);
          // Calculate valid moves for new piece
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

  return <>{children(game, selectedSquare, validMoves, selectSquare)}</>;
};

export default GameController;

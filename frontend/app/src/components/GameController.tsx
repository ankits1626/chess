import { useState } from 'react';
import { Chess } from 'chess.js';
import type { Square } from '../types/chess';

interface GameControllerProps {
  children: (
    game: Chess,
    selectedSquare: Square | null,
    selectSquare: (square: Square) => void
  ) => React.ReactNode;
}

const GameController = ({ children }: GameControllerProps) => {
  const [game, setGame] = useState(() => new Chess());
  const [selectedSquare, setSelectedSquare] = useState<Square | null>(null);

  const selectSquare = (square: Square) => {
    // If no square is selected, select this square (if it has a piece)
    if (!selectedSquare) {
      const piece = game.get(square);
      if (piece && piece.color === game.turn()) {
        setSelectedSquare(square);
      }
      return;
    }

    // If clicking the same square, deselect it
    if (selectedSquare === square) {
      setSelectedSquare(null);
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
        // Move was successful, update state
        setGame(new Chess(game.fen())); // Create new instance to trigger re-render
        setSelectedSquare(null);
      } else {
        // Invalid move, check if clicking another piece of the same color
        const piece = game.get(square);
        if (piece && piece.color === game.turn()) {
          setSelectedSquare(square);
        } else {
          setSelectedSquare(null);
        }
      }
    } catch (error) {
      // Invalid move, deselect
      setSelectedSquare(null);
    }
  };

  return <>{children(game, selectedSquare, selectSquare)}</>;
};

export default GameController;
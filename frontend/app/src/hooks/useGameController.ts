import { useState, useEffect } from 'react';
import { Chess } from 'chess.js';

const useGameController = () => {
  const [game, setGame] = useState(new Chess());
  const [board, setBoard] = useState(game.board());

  useEffect(() => {
    setBoard(game.board());
  }, [game]);

  const makeMove = (from: string, to: string) => {
    try {
      const newGame = new Chess(game.fen()); // Create a new instance to avoid direct state mutation
      const move = newGame.move({ from, to, promotion: 'q' }); // Always promote to queen for now
      if (move) {
        setGame(newGame);
      }
    } catch (e) {
      console.log(e);
    }
  };

  return {
    game,
    board,
    makeMove,
  };
};

export default useGameController;

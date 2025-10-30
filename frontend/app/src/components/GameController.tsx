import { useState } from 'react';
import { Chess } from 'chess.js';
import type { Square } from '../types/chess';

interface GameControllerProps {
  children: (game: Chess) => React.ReactNode;
}

const GameController = ({ children }: GameControllerProps) => {
  const [game, setGame] = useState(() => new Chess());

  // This component doesn't render UI itself, but provides state to its children.
  // The type assertion here is safe because React guarantees children will be present.
  return <>{children(game)}</>;
};

export default GameController;

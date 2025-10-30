import { useState } from 'react';
import { Chess } from 'chess.js';
import type { Square, PieceType, PieceColor } from '../types/chess';

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

interface GameControllerProps {
  children: (
    game: Chess,
    selectedSquare: Square | null,
    validMoves: Square[],
    selectSquare: (square: Square) => void,
    pendingMove: PendingPromotion | null,
    handlePromotion: (piece: PieceType) => void,
    lastMove: LastMove,
    resetGame: () => void,
    debugActions?: DebugActions
  ) => React.ReactNode;
}

const GameController = ({ children }: GameControllerProps) => {
  const [game, setGame] = useState(() => new Chess());
  const [selectedSquare, setSelectedSquare] = useState<Square | null>(null);
  const [validMoves, setValidMoves] = useState<Square[]>([]);
  const [pendingMove, setPendingMove] = useState<PendingPromotion | null>(null);
  const [lastMove, setLastMove] = useState<LastMove>(null);

  const resetGame = () => {
    setGame(new Chess());
    setSelectedSquare(null);
    setValidMoves([]);
    setPendingMove(null);
    setLastMove(null);
  };

  const debugActions: DebugActions | undefined = import.meta.env.DEV ? {
    loadFen: (fen: string) => {
      try {
        const newGame = new Chess(fen);
        setGame(newGame);
        setSelectedSquare(null);
        setValidMoves([]);
        setPendingMove(null);
        setLastMove(null);
      } catch (error) {
        console.error('Invalid FEN:', error);
      }
    },
    resetGame
  } : undefined;

  const selectSquare = (square: Square) => {
    if (pendingMove) return;

    if (!selectedSquare) {
      const piece = game.get(square);
      if (piece && piece.color === game.turn()) {
        setSelectedSquare(square);
        const moves = game.moves({ square, verbose: true });
        setValidMoves(moves.map(move => move.to));
      }
      return;
    }

    if (selectedSquare === square) {
      setSelectedSquare(null);
      setValidMoves([]);
      return;
    }

    try {
      const piece = game.get(selectedSquare);

      if (piece?.type === 'p' && (square.endsWith('1') || square.endsWith('8'))) {
        const moves = game.moves({ square: selectedSquare, verbose: true });
        const isValidMove = moves.some(m => m.to === square);

        if (isValidMove) {
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

      const move = game.move({
        from: selectedSquare,
        to: square
      });

      if (move) {
        setGame(new Chess(game.fen()));
        setLastMove({ from: selectedSquare, to: square });
        setSelectedSquare(null);
        setValidMoves([]);
      } else {
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
    } catch {
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
    setGame(new Chess(game.fen()));
    setLastMove({ from: pendingMove.from, to: pendingMove.to });
    setPendingMove(null);
  };

  return <>{children(game, selectedSquare, validMoves, selectSquare, pendingMove, handlePromotion, lastMove, resetGame, debugActions)}</>;
};

export default GameController;

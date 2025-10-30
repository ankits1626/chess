import type { Chess } from 'chess.js';

interface GameInfoProps {
  game: Chess;
}

const GameInfo = ({ game }: GameInfoProps) => {
  const turn = game.turn() === 'w' ? 'White' : 'Black';
  const isCheck = game.isCheck();
  const isCheckmate = game.isCheckmate();
  const isDraw = game.isDraw();
  const isStalemate = game.isStalemate();

  return (
    <div className="bg-gray-700 p-6 rounded-lg w-full lg:w-64">
      <h2 className="text-2xl font-bold mb-4">Game Info</h2>

      <div className="space-y-2">
        <p><span className="font-semibold">Turn:</span> {turn}</p>

        {isCheck && !isCheckmate && (
          <p className="text-yellow-400 font-semibold">Check!</p>
        )}

        {isCheckmate && (
          <p className="text-red-400 font-bold text-xl">
            Checkmate! {turn === 'White' ? 'Black' : 'White'} wins!
          </p>
        )}

        {isStalemate && (
          <p className="text-blue-400 font-bold">Stalemate - Draw!</p>
        )}

        {isDraw && !isStalemate && (
          <p className="text-blue-400 font-bold">Draw!</p>
        )}
      </div>
    </div>
  );
};

export default GameInfo;

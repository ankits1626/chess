import { useGameStore } from '../store/useGameStore';
import MoveHistory from './MoveHistory';

const GameInfo = () => {
  const game = useGameStore(state => state.game);
  const onNewGame = useGameStore(state => state.resetGame);

  const turn = game.turn() === 'w' ? 'White' : 'Black';
  const isCheck = game.isCheck();
  const isCheckmate = game.isCheckmate();
  const isDraw = game.isDraw();
  const isStalemate = game.isStalemate();
  const moveHistory = game.history();

  return (
    <div className="bg-gray-700 p-6 rounded-lg w-full lg:w-80">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold">Game Info</h2>
        <button
          onClick={onNewGame}
          className="bg-blue-600 hover:bg-blue-500 text-white text-sm font-semibold px-3 py-1 rounded transition-colors"
        >
          New Game
        </button>
      </div>

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

      <MoveHistory moves={moveHistory} />
    </div>
  );
};

export default GameInfo;

import { useGameStore } from '@/store/useGameStore';
import MoveHistory from './MoveHistory';

const GameInfo = () => {
  const game = useGameStore(state => state.game);
  const onNewGame = useGameStore(state => state.resetGame);
  const opponentType = useGameStore(state => state.opponentType);
  const computerDifficulty = useGameStore(state => state.computerDifficulty);
  const isComputerThinking = useGameStore(state => state.isComputerThinking);
  const gameError = useGameStore(state => state.gameError);
  const mode = useGameStore(state => state.mode);
  const toggleBoardFlip = useGameStore(state => state.toggleBoardFlip);
  const isBoardFlipped = useGameStore(state => state.isBoardFlipped);

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
        <div className="flex gap-2">
          <button
            onClick={toggleBoardFlip}
            className="bg-gray-600 hover:bg-gray-500 text-white text-sm font-semibold px-3 py-1 rounded transition-colors flex items-center gap-1"
            title="Flip board"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4" />
            </svg>
            {isBoardFlipped ? 'Reset' : 'Flip'}
          </button>
          <button
            onClick={onNewGame}
            className="bg-blue-600 hover:bg-blue-500 text-white text-sm font-semibold px-3 py-1 rounded transition-colors"
          >
            New Game
          </button>
        </div>
      </div>

      {/* Computer Game Info */}
      {mode === 'computer' && opponentType === 'computer' && (
        <div className="mb-4 p-3 bg-gray-800 rounded-lg border border-gray-600">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-semibold text-gray-300">Opponent</span>
            <span className="text-sm font-bold text-blue-400 capitalize">
              Computer ({computerDifficulty})
            </span>
          </div>
          {isComputerThinking && (
            <div className="flex items-center gap-2 text-sm text-yellow-400 animate-pulse">
              <div className="w-4 h-4 border-2 border-yellow-400 border-t-transparent rounded-full animate-spin"></div>
              <span>Computer is thinking...</span>
            </div>
          )}
        </div>
      )}

      {/* Error Display */}
      {gameError && (
        <div className="mb-4 p-3 bg-red-900 bg-opacity-50 border border-red-600 rounded-lg">
          <p className="text-sm text-red-200">{gameError}</p>
        </div>
      )}

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

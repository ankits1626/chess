import GameBoard from './components/GameBoard';
import GameInfo from './components/GameInfo';
import GameController from './components/GameController';

function App() {
  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-800 text-white p-8">
      <GameController>
        {(game, selectedSquare, selectSquare) => (
          <div className="flex flex-row gap-8 items-center">
            <GameBoard
              game={game}
              selectedSquare={selectedSquare}
              onSquareClick={selectSquare}
            />
            <GameInfo game={game} />
          </div>
        )}
      </GameController>
    </div>
  );
}

export default App;
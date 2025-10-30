import GameBoard from './components/GameBoard';
import GameInfo from './components/GameInfo';

function App() {
  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-800 text-white p-8">
      <div className="flex flex-row gap-8 items-center">
        <GameBoard />
        <GameInfo />
      </div>
    </div>
  );
}

export default App;
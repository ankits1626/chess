import { lazy, Suspense } from 'react';
import GameBoard from './components/GameBoard';
import GameInfo from './components/GameInfo';
import GameController from './components/GameController';
import PromotionDialog from './components/PromotionDialog';
import { useDebugPanel } from './hooks/useDebugPanel';

// Dynamically import the DebugPanel only in development mode.
const DebugPanel = import.meta.env.DEV
  ? lazy(() => import('./components/debug/DebugPanel'))
  : null;

function App() {
  const { isPanelVisible, setIsPanelVisible } = useDebugPanel();

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-800 text-white p-8">
      <GameController>
        {(game, selectedSquare, validMoves, selectSquare, pendingMove, handlePromotion, lastMove, resetGame, debugActions) => (
          <>
            <div className="flex flex-col lg:flex-row gap-4 lg:gap-8 items-center">
              <GameBoard
                game={game}
                selectedSquare={selectedSquare}
                validMoves={validMoves}
                lastMove={lastMove}
                onSquareClick={selectSquare}
              />
              <GameInfo game={game} onNewGame={resetGame} />
            </div>

            {pendingMove && (
              <PromotionDialog
                color={pendingMove.color}
                onSelectPiece={handlePromotion}
              />
            )}

            {/* Conditionally render the DebugPanel in dev mode */}
            {isPanelVisible && DebugPanel && debugActions && (
              <Suspense fallback={null}>
                <DebugPanel 
                  actions={debugActions} 
                  currentFen={game.fen()} 
                  onClose={() => setIsPanelVisible(false)} 
                />
              </Suspense>
            )}
          </>
        )}
      </GameController>
    </div>
  );
}

export default App;
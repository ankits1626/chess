import { lazy, Suspense } from 'react';
import { useGameStore } from './store/useGameStore';
import GameBoard from './components/GameBoard';
import GameInfo from './components/GameInfo';
import PromotionDialog from './components/PromotionDialog';
import GameImporter from './components/GameImporter';
import { useDebugPanel } from './hooks/useDebugPanel';

const DebugPanel = import.meta.env.DEV
  ? lazy(() => import('./components/debug/DebugPanel'))
  : null;

function App() {
  const { isPanelVisible, setIsPanelVisible } = useDebugPanel();
  const pendingMove = useGameStore(state => state.pendingMove);
  const handlePromotion = useGameStore(state => state.handlePromotion);

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-800 text-white p-8">
      <div className="flex flex-col lg:flex-row gap-4 lg:gap-8 items-start">
        <GameBoard />
        <GameInfo />
        <GameImporter />
      </div>

      {pendingMove && (
        <PromotionDialog
          color={pendingMove.color}
          onSelectPiece={handlePromotion}
        />
      )}

      {isPanelVisible && DebugPanel && (
        <Suspense fallback={null}>
          <DebugPanel onClose={() => setIsPanelVisible(false)} />
        </Suspense>
      )}
    </div>
  );
}

export default App;
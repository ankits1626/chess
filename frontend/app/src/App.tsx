import { lazy, Suspense, useEffect } from 'react';
import { useGameStore } from '@/store/useGameStore';
import GameBoard from '@/components/board/GameBoard';
import GameInfo from '@/components/game/GameInfo';
import PromotionDialog from '@/components/game/PromotionDialog';
import GameImporter from '@/components/importer/GameImporter';
import ReplayControls from '@/components/replay/ReplayControls';
import MoveList from '@/components/replay/MoveList';
import PlayerDisplay from '@/components/replay/PlayerDisplay';
import { useDebugPanel } from '@/hooks/useDebugPanel';

const DebugPanel = import.meta.env.DEV
  ? lazy(() => import('@/components/debug/DebugPanel'))
  : null;

function App() {
  const { isPanelVisible, setIsPanelVisible } = useDebugPanel();
  const pendingMove = useGameStore(state => state.pendingMove);
  const handlePromotion = useGameStore(state => state.handlePromotion);

  // Replay state and actions
  const mode = useGameStore(state => state.mode);
  const replayIndex = useGameStore(state => state.replayIndex);
  const replayMoves = useGameStore(state => state.replayMoves);
  const isAutoplaying = useGameStore(state => state.isAutoplaying);
  const whitePlayer = useGameStore(state => state.whitePlayer);
  const blackPlayer = useGameStore(state => state.blackPlayer);

  const goToFirstMove = useGameStore(state => state.goToFirstMove);
  const prevMove = useGameStore(state => state.prevMove);
  const nextMove = useGameStore(state => state.nextMove);
  const goToLastMove = useGameStore(state => state.goToLastMove);
  const toggleAutoplay = useGameStore(state => state.toggleAutoplay);
  const stopAutoplay = useGameStore(state => state.stopAutoplay);
  const goToMove = useGameStore(state => state.goToMove);

  const isFirstMove = replayIndex === -1;
  const isLastMove = replayIndex === replayMoves.length - 1;

  // Cleanup autoplay on unmount
  useEffect(() => {
    return () => {
      stopAutoplay();
    };
  }, [stopAutoplay]);

  // Keyboard shortcuts for replay controls
  useEffect(() => {
    if (mode !== 'replay') return;

    const handleKeyDown = (e: KeyboardEvent) => {
      switch (e.key) {
        case 'ArrowLeft':
          if (!isFirstMove) prevMove();
          break;
        case 'ArrowRight':
          if (!isLastMove) nextMove();
          break;
        case ' ': // Spacebar
          e.preventDefault(); // Prevent page scroll
          toggleAutoplay();
          break;
        case 'Home':
          goToFirstMove();
          break;
        case 'End':
          goToLastMove();
          break;
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [mode, isFirstMove, isLastMove, prevMove, nextMove, toggleAutoplay, goToFirstMove, goToLastMove]);

  return (
    <div className="min-h-screen bg-gray-800 text-white p-4 sm:p-6 lg:p-8">
      <div className="max-w-[1600px] mx-auto h-full">
        <div className="flex flex-col lg:flex-row gap-4 lg:gap-8 items-start lg:items-center justify-center min-h-[calc(100vh-2rem)] sm:min-h-[calc(100vh-3rem)] lg:min-h-[calc(100vh-4rem)]">
          <div className="flex flex-col gap-2 items-center w-full lg:w-auto">
            {mode === 'replay' && <PlayerDisplay name={blackPlayer} />}
            <GameBoard />
            {mode === 'replay' && <PlayerDisplay name={whitePlayer} />}
          </div>
          <div className="flex flex-col gap-4 w-full lg:w-80 pb-8 lg:pb-0">
            <GameInfo />
            <GameImporter />
            {mode === 'replay' && replayMoves.length > 0 && (
              <>
                <MoveList
                  moves={replayMoves}
                  currentMoveIndex={replayIndex}
                  onMoveClick={goToMove}
                />
                <ReplayControls
                  isFirstMove={isFirstMove}
                  isLastMove={isLastMove}
                  isAutoplaying={isAutoplaying}
                  onGoToFirst={goToFirstMove}
                  onPrevious={prevMove}
                  onToggleAutoplay={toggleAutoplay}
                  onNext={nextMove}
                  onGoToLast={goToLastMove}
                />
                {replayIndex >= 0 && (
                  <div className="text-center text-sm text-gray-400 mt-2">
                    Move {replayIndex + 1} of {replayMoves.length}
                  </div>
                )}
              </>
            )}
          </div>
        </div>
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

import { lazy, Suspense, useEffect, useState } from 'react';
import { useGameStore } from '@/store/useGameStore';
import type { Difficulty, PlayerColor } from '@/types/game';
import GameBoard from '@/components/board/GameBoard';
import GameInfo from '@/components/game/GameInfo';
import GameSetup from '@/components/game/GameSetup';
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
  const [showGameSetup, setShowGameSetup] = useState(false);
  const [isStartingGame, setIsStartingGame] = useState(false);

  const pendingMove = useGameStore(state => state.pendingMove);
  const handlePromotion = useGameStore(state => state.handlePromotion);
  const startComputerGame = useGameStore(state => state.startComputerGame);

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

  // Board flip state
  const playerColor = useGameStore(state => state.playerColor);
  const isBoardFlipped = useGameStore(state => state.isBoardFlipped);

  // Calculate if board is visually flipped (same logic as GameBoard/GameInfo)
  const autoFlip = playerColor === 'black';
  const isFlipped = playerColor ? (isBoardFlipped ? !autoFlip : autoFlip) : isBoardFlipped;

  // Determine player display order based on board orientation
  const topPlayer = isFlipped ? whitePlayer : blackPlayer;
  const bottomPlayer = isFlipped ? blackPlayer : whitePlayer;

  // Handle starting a computer game
  const handleStartComputerGame = async (color: PlayerColor, difficulty: Difficulty) => {
    setIsStartingGame(true);
    try {
      // TODO: Replace with actual user ID from auth system
      const userId = '25d30da5-0cc4-4f5a-8c88-69d0f90b004c';
      await startComputerGame(color, difficulty, userId);
      setShowGameSetup(false);
    } catch (error) {
      console.error('Failed to start computer game:', error);
      alert('Failed to start game. Please try again.');
    } finally {
      setIsStartingGame(false);
    }
  };

  // Cleanup autoplay and WebSocket connection on unmount
  useEffect(() => {
    const cleanup = () => {
      stopAutoplay();
      // GameStore handles WebSocket disconnection in resetGame
    };
    return cleanup;
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
            {mode === 'replay' && <PlayerDisplay name={topPlayer} />}
            <GameBoard />
            {mode === 'replay' && <PlayerDisplay name={bottomPlayer} />}
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

      {showGameSetup && (
        <GameSetup
          onStartGame={handleStartComputerGame}
          onClose={() => setShowGameSetup(false)}
          isLoading={isStartingGame}
        />
      )}

      {isPanelVisible && DebugPanel && (
        <Suspense fallback={null}>
          <DebugPanel onClose={() => setIsPanelVisible(false)} />
        </Suspense>
      )}

      {/* Floating button to start computer game */}
      {mode === 'live' && (
        <button
          onClick={() => setShowGameSetup(true)}
          className="fixed bottom-6 right-6 bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-500 hover:to-blue-500 text-white font-bold py-4 px-6 rounded-full shadow-lg transition-all transform hover:scale-105 flex items-center gap-2"
        >
          <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
          </svg>
          <span>Play vs Computer</span>
        </button>
      )}
    </div>
  );
}

export default App;

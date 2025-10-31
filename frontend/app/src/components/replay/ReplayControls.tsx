import type { FC } from 'react';

interface ReplayControlsProps {
  isFirstMove: boolean;
  isLastMove: boolean;
  isAutoplaying: boolean;
  onGoToFirst: () => void;
  onPrevious: () => void;
  onToggleAutoplay: () => void;
  onNext: () => void;
  onGoToLast: () => void;
}

const ReplayControls: FC<ReplayControlsProps> = ({
  isFirstMove,
  isLastMove,
  isAutoplaying,
  onGoToFirst,
  onPrevious,
  onToggleAutoplay,
  onNext,
  onGoToLast,
}) => {
  return (
    <div className="flex items-center justify-center gap-2 p-4">
      <button
        onClick={onGoToFirst}
        disabled={isFirstMove}
        aria-label="Go to first move"
        className="px-4 py-2 bg-gray-700 text-white rounded hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        &lt;&lt;
      </button>

      <button
        onClick={onPrevious}
        disabled={isFirstMove}
        aria-label="Previous move"
        className="px-4 py-2 bg-gray-700 text-white rounded hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        &lt;
      </button>

      <button
        onClick={onToggleAutoplay}
        aria-label={isAutoplaying ? 'Pause' : 'Play'}
        className="px-6 py-2 bg-blue-600 text-white rounded hover:bg-blue-500 transition-colors"
      >
        {isAutoplaying ? 'Pause' : 'Play'}
      </button>

      <button
        onClick={onNext}
        disabled={isLastMove}
        aria-label="Next move"
        className="px-4 py-2 bg-gray-700 text-white rounded hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        &gt;
      </button>

      <button
        onClick={onGoToLast}
        disabled={isLastMove}
        aria-label="Go to last move"
        className="px-4 py-2 bg-gray-700 text-white rounded hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        &gt;&gt;
      </button>
    </div>
  );
};

export default ReplayControls;

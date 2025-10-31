import type { FC } from 'react';

interface ArchivePaginatorProps {
  archives: string[];
  currentArchiveUrl: string;
  onSelectArchive: (url: string) => void;
}

const ArchivePaginator: FC<ArchivePaginatorProps> = ({ archives, currentArchiveUrl, onSelectArchive }) => {
  const currentIndex = archives.indexOf(currentArchiveUrl);
  const canGoPrevious = currentIndex > 0;
  const canGoNext = currentIndex < archives.length - 1;

  // Extract month/year from URL like "https://api.chess.com/pub/player/hikaru/games/2025/10"
  const extractMonthYear = (url: string) => {
    const match = url.match(/(\d{4})\/(\d{2})$/);
    if (!match) return url;
    const [, year, month] = match;
    const date = new Date(Number(year), Number(month) - 1);
    return date.toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
  };

  return (
    <div className="flex items-center justify-between bg-gray-900 p-3 rounded-lg">
      <button
        onClick={() => onSelectArchive(archives[currentIndex - 1])}
        disabled={!canGoPrevious}
        className="px-4 py-2 bg-gray-700 text-white font-semibold rounded-md hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        ← Older
      </button>
      <span className="font-semibold text-lg text-white">{extractMonthYear(currentArchiveUrl)}</span>
      <button
        onClick={() => onSelectArchive(archives[currentIndex + 1])}
        disabled={!canGoNext}
        className="px-4 py-2 bg-gray-700 text-white font-semibold rounded-md hover:bg-gray-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
      >
        Newer →
      </button>
    </div>
  );
};

export default ArchivePaginator;

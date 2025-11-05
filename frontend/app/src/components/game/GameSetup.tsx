import { useState } from 'react';
import type { Difficulty, PlayerColor } from '@/types/game';

interface GameSetupProps {
  onStartGame: (color: PlayerColor, difficulty: Difficulty) => void;
  onClose: () => void;
  isLoading?: boolean;
}

const GameSetup = ({ onStartGame, onClose, isLoading = false }: GameSetupProps) => {
  const [playerColor, setPlayerColor] = useState<PlayerColor>('white');
  const [difficulty, setDifficulty] = useState<Difficulty>('medium');

  const handleSubmit = () => {
    onStartGame(playerColor, difficulty);
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-gray-800 rounded-lg p-8 max-w-md w-full mx-4">
        <h2 className="text-2xl font-bold mb-6 text-white">Play vs Computer</h2>

        <div className="space-y-6">
          {/* Color Selection */}
          <div>
            <label className="block text-sm font-semibold mb-3 text-gray-300">
              Choose Your Color
            </label>
            <div className="grid grid-cols-2 gap-3">
              <button
                onClick={() => setPlayerColor('white')}
                className={`py-3 px-4 rounded-lg font-semibold transition-all ${
                  playerColor === 'white'
                    ? 'bg-blue-600 text-white ring-2 ring-blue-400'
                    : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                }`}
              >
                <div className="flex items-center justify-center gap-2">
                  <div className="w-6 h-6 bg-white rounded-full border-2 border-gray-400"></div>
                  <span>White</span>
                </div>
              </button>
              <button
                onClick={() => setPlayerColor('black')}
                className={`py-3 px-4 rounded-lg font-semibold transition-all ${
                  playerColor === 'black'
                    ? 'bg-blue-600 text-white ring-2 ring-blue-400'
                    : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                }`}
              >
                <div className="flex items-center justify-center gap-2">
                  <div className="w-6 h-6 bg-gray-900 rounded-full border-2 border-gray-400"></div>
                  <span>Black</span>
                </div>
              </button>
            </div>
          </div>

          {/* Difficulty Selection */}
          <div>
            <label className="block text-sm font-semibold mb-3 text-gray-300">
              Difficulty Level
            </label>
            <div className="space-y-2">
              {['easy', 'medium', 'hard'].map((level) => (
                <button
                  key={level}
                  onClick={() => setDifficulty(level as Difficulty)}
                  className={`w-full py-3 px-4 rounded-lg font-semibold transition-all text-left ${
                    difficulty === level
                      ? 'bg-blue-600 text-white ring-2 ring-blue-400'
                      : 'bg-gray-700 text-gray-300 hover:bg-gray-600'
                  }`}
                >
                  <div className="flex items-center justify-between">
                    <span className="capitalize">{level}</span>
                    <span className="text-sm opacity-70">
                      {level === 'easy' && '~100ms think time'}
                      {level === 'medium' && '~500ms think time'}
                      {level === 'hard' && '~2s think time'}
                    </span>
                  </div>
                </button>
              ))}
            </div>
          </div>

          {/* Action Buttons */}
          <div className="flex gap-3 mt-8">
            <button
              onClick={onClose}
              disabled={isLoading}
              className="flex-1 py-3 px-4 bg-gray-700 hover:bg-gray-600 disabled:bg-gray-800 disabled:text-gray-500 text-white font-semibold rounded-lg transition-colors"
            >
              Cancel
            </button>
            <button
              onClick={handleSubmit}
              disabled={isLoading}
              className="flex-1 py-3 px-4 bg-green-600 hover:bg-green-500 disabled:bg-gray-700 disabled:text-gray-500 text-white font-semibold rounded-lg transition-colors flex items-center justify-center gap-2"
            >
              {isLoading ? (
                <>
                  <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                  <span>Starting...</span>
                </>
              ) : (
                <span>Start Game</span>
              )}
            </button>
          </div>
        </div>

        {/* Info Note */}
        <p className="mt-6 text-sm text-gray-400 text-center">
          {playerColor === 'white'
            ? 'You will make the first move'
            : 'Computer will make the first move'
          }
        </p>
      </div>
    </div>
  );
};

export default GameSetup;

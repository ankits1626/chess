import FenLoader from './FenLoader';
import ScenarioPicker from './ScenarioPicker';
import type { DebugActions } from '../GameController'; // We will define this type in GameController

interface DebugPanelProps {
  actions: DebugActions;
  currentFen: string;
  onClose: () => void;
}

const DebugPanel = ({ actions, currentFen, onClose }: DebugPanelProps) => {
  return (
    <div className="fixed bottom-4 right-4 bg-gray-800/90 backdrop-blur-sm border border-gray-700 rounded-lg shadow-2xl z-50 w-full max-w-md p-4">
      <div className="flex justify-between items-center mb-4">
        <h3 className="text-lg font-bold text-white">🛠️ Debug Toolkit</h3>
        <button onClick={onClose} className="text-gray-400 hover:text-white">&times;</button>
      </div>
      <div className="space-y-4">
        <FenLoader onLoad={actions.loadFen} currentFen={currentFen} />
        <ScenarioPicker onSelect={actions.loadFen} />
        <div>
          <button
            onClick={actions.resetGame}
            className="bg-red-600 hover:bg-red-500 text-white font-semibold px-4 py-2 rounded transition-colors w-full"
          >
            Reset to Start
          </button>
        </div>
      </div>
    </div>
  );
};

export default DebugPanel;

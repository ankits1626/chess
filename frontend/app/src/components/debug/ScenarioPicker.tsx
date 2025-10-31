import { scenarios } from '@/constants/debugScenarios';

interface ScenarioPickerProps {
  onSelect: (fen: string) => void;
}

const ScenarioPicker = ({ onSelect }: ScenarioPickerProps) => {
  return (
    <div className="flex flex-col gap-2">
      <label className="text-sm font-medium text-gray-300">Quick Scenarios</label>
      <div className="flex flex-wrap gap-2">
        {Object.entries(scenarios).map(([name, fen]) => (
          <button
            key={name}
            onClick={() => onSelect(fen)}
            className="bg-gray-600 hover:bg-gray-500 text-white text-sm px-3 py-1 rounded-full transition-colors"
          >
            {name}
          </button>
        ))}
      </div>
    </div>
  );
};

export default ScenarioPicker;

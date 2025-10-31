import type { FC } from 'react';

interface PlayerDisplayProps {
  name: string | null;
}

const PlayerDisplay: FC<PlayerDisplayProps> = ({ name }) => {
  if (!name) return null;

  return (
    <div className="bg-gray-900 text-white text-lg rounded-md px-4 py-2">
      {name}
    </div>
  );
};

export default PlayerDisplay;

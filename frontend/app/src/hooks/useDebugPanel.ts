import { useEffect, useState } from 'react';

// This hook manages the visibility of the debug panel and listens for the shortcut.
export const useDebugPanel = () => {
  const [isPanelVisible, setIsPanelVisible] = useState(false);

  useEffect(() => {
    // The entire hook is a no-op in production.
    if (!import.meta.env.DEV) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      // Ctrl+Shift+D or Cmd+Shift+D
      if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.key.toUpperCase() === 'D') {
        e.preventDefault();
        setIsPanelVisible(prev => !prev);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, []);

  // In production, the panel is never visible.
  if (!import.meta.env.DEV) {
    return { isPanelVisible: false, setIsPanelVisible: () => {} };
  }

  return { isPanelVisible, setIsPanelVisible };
};

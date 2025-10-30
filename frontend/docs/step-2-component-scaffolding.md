# Step 2: Component Scaffolding Plan

This document outlines the plan for creating the initial file structure for our React components, incorporating modern best practices for React 19.

## 1. Directory Creation

To maintain a clean and organized codebase, the following directories will be created within the `src/` folder:

- `src/components/`: This will house all of our presentational (UI) React components.
- `src/hooks/`: This will be used for custom React hooks, starting with our main game logic controller.
- `src/types/`: This will contain shared TypeScript type definitions and interfaces.

## 2. Cleanup and File Removal

Before scaffolding new components, we will clean up the default Vite boilerplate files that won't be used.

### Files to Delete:
- `src/App.css` - Not needed (using Tailwind CSS).
- `src/assets/react.svg` - Default boilerplate asset.
- `public/vite.svg` - Default boilerplate asset.

### Files Already Cleaned:
- `src/index.css` - Already contains only the Tailwind import directive.

## 3. Component File Creation

The following files will be created. Each will contain basic placeholder content.

- **`src/components/GameBoard.tsx`**
- **`src/components/Square.tsx`**
- **`src/components/Piece.tsx`**
- **`src/components/GameInfo.tsx`**
- **`src/components/PromotionDialog.tsx`**
- **`src/hooks/useGameController.ts`**
- **`src/types/chess.ts`**

### Example Placeholder Content (React 19 Pattern):

```typescript
// src/components/GameBoard.tsx
const GameBoard = () => {
  return (
    <div className="w-96 h-96 bg-gray-500 flex items-center justify-center">
      <p className="text-white text-xl">GameBoard</p>
    </div>
  );
};

export default GameBoard;
```
**Note:** React 19 does not require `import React from 'react'` for JSX, and we will avoid the deprecated `React.FC` type.

## 4. App.tsx Cleanup

The main `App.tsx` file will be cleaned up and updated to render our primary components with an improved layout.

### After:

```typescript
// Updated App.tsx
import GameBoard from './components/GameBoard';
import GameInfo from './components/GameInfo';

function App() {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-800 text-white">
      <h1 className="text-4xl font-bold mb-8">Chess Coach</h1>
      <div className="flex flex-col lg:flex-row gap-8">
        <GameBoard />
        <GameInfo />
      </div>
    </div>
  );
}

export default App;
```

This updated plan establishes a clean, modern, and foundational structure for the application.
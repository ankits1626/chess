# Chess Coach App - File Structure Explained

This document explains every file in the `frontend/app/` directory, what it does, and why it exists.

---

## Root Configuration Files

### `package.json`
**What it is:** The project's manifest file
**What it does:**
- Lists all dependencies (libraries your app needs): `react`, `chess.js`, `tailwindcss`, etc.
- Defines scripts you can run: `pnpm dev`, `pnpm build`, `pnpm lint`
- Contains project metadata: name, version, type

**Think of it as:** A recipe card that tells `pnpm` what ingredients (packages) to download and how to cook (build/run) your app.

**Example:**
```json
{
  "scripts": {
    "dev": "vite"  // When you run `pnpm dev`, it runs `vite`
  },
  "dependencies": {
    "react": "^19.2.0"  // Your app needs React version 19.2.0
  }
}
```

---

### `vite.config.ts`
**What it is:** Configuration for Vite (the build tool)

**What is Vite?**
Vite (pronounced "veet", French for "fast") is a modern build tool that:
1. **Runs a dev server** - Serves your app at `http://localhost:5173/`
2. **Hot Module Replacement (HMR)** - Instantly updates the browser when you save files (no manual refresh!)
3. **Transforms code** - Converts TypeScript → JavaScript, JSX → React calls
4. **Bundles for production** - Creates optimized files for deployment

**Think of Vite as:** A smart assistant that watches your code and instantly shows changes in the browser.

**What it does:**
- Tells Vite to use React plugin (for JSX support)
- Tells Vite to use Tailwind CSS plugin (for styling)
- Configures how your code is bundled and served

**Current content:**
```typescript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [react(), tailwindcss()],
})
```

**Line-by-line:**
- `defineConfig()` - Helper for TypeScript autocomplete
- `react()` - Plugin that handles JSX and React features
- `tailwindcss()` - Plugin that processes Tailwind utility classes

**Translation:** "Use React and Tailwind CSS plugins when building/serving the app."

**Common commands:**
- `pnpm dev` - Starts Vite dev server with instant hot reload
- `pnpm build` - Bundles code for production deployment
- `pnpm preview` - Preview production build locally

**When you run `pnpm dev`:**
```bash
VITE v7.1.12  ready in 131 ms
➜  Local:   http://localhost:5173/
```
This means Vite started a web server on your computer at port 5173.

**Hot reload example:**
1. You change `bg-amber-100` to `bg-blue-500` in Square.tsx
2. Press Cmd+S (save)
3. **Instantly** (50-200ms) the browser shows blue squares - no refresh needed!
4. Console shows: `[vite] hmr update /src/components/Square.tsx`

**Why Vite is fast:**
- Only transforms files that are requested (on-demand)
- Uses esbuild (written in Go, extremely fast)
- Native ES modules support
- Caches transformed files

**Vite vs older tools (like Create React App):**
- ⚡ 10-100x faster startup
- ⚡ Instant hot reload vs 1-5 second refresh
- 🎯 Minimal configuration needed

---

### `tsconfig.json`
**What it is:** Base TypeScript configuration
**What it does:**
- References two other TypeScript configs: `tsconfig.app.json` and `tsconfig.node.json`
- Acts as a container that tells TypeScript to look at multiple configs

**Think of it as:** A table of contents pointing to different chapters (configs) for different parts of your project.

---

### `tsconfig.app.json`
**What it is:** TypeScript configuration for your application code (src/ folder)
**What it does:**
- Sets TypeScript compiler options for your React code
- Enables strict type checking (`"strict": true`)
- Configures module system (ESNext), JSX mode (react-jsx)
- Enables special features like `verbatimModuleSyntax` (requires `import type` for types)

**Key settings explained:**
```json
{
  "target": "ES2022",           // Compile to ES2022 JavaScript
  "jsx": "react-jsx",            // Use React 19's new JSX transform (no need to import React)
  "strict": true,                // Enable all strict type checks
  "verbatimModuleSyntax": true   // Require explicit `import type` for types
}
```

**Why it matters:** This is why you get TypeScript errors like "must use import type" - it's enforcing clean separation of types and values.

---

### `tsconfig.node.json`
**What it is:** TypeScript configuration for Node.js files (build scripts, config files)
**What it does:**
- Configures TypeScript for files like `vite.config.ts` that run in Node, not the browser
- Different settings than app code because Node has different capabilities

**Think of it as:** Separate instructions for the kitchen staff (build tools) vs the dining room (browser).

---

### `eslint.config.js`
**What it is:** Configuration for ESLint (code quality checker)
**What it does:**
- Defines coding rules: no unused variables, no missing dependencies in hooks, etc.
- Runs when you execute `pnpm lint`
- Shows warnings/errors in your IDE for code issues

**Think of it as:** A spell-checker and grammar-checker for your code.

**Example rules:**
- "Don't leave unused variables" → `noUnusedLocals: true`
- "React hooks must have correct dependencies" → `eslint-plugin-react-hooks`

---

### `index.html`
**What it is:** The main HTML file
**What it does:**
- The only actual HTML page in your app
- Contains a `<div id="root"></div>` where React mounts your app
- Loads `src/main.tsx` which starts the React app

**Think of it as:** The picture frame - React fills in the picture (your app).

**Key part:**
```html
<body>
  <div id="root"></div>          <!-- React app goes here -->
  <script type="module" src="/src/main.tsx"></script>  <!-- Starts React -->
</body>
```

---

## Source Code Files (`src/`)

### `src/main.tsx`
**What it is:** The entry point of your React application
**What it does:**
1. Imports React and ReactDOM
2. Imports your root `App` component
3. Imports global styles (`index.css`)
4. Finds the `<div id="root">` in `index.html`
5. Renders your `<App />` component inside it

**Think of it as:** The ignition key that starts your app.

**Code breakdown:**
```typescript
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'              // Global styles
import App from './App.tsx'       // Your app

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <App />                       // Your entire app starts here
  </StrictMode>,
)
```

**Why `StrictMode`?** It enables extra development checks to catch bugs early.

---

### `src/App.tsx`
**What it is:** The root component of your application
**What it does:**
- Defines the overall layout (dark background, centered content)
- Renders `<GameBoard />` and `<GameInfo />` components
- Acts as the top-level container for your entire chess app

**Think of it as:** The main stage where all your components perform.

**Current code:**
```typescript
import GameBoard from './components/GameBoard';
import GameInfo from './components/GameInfo';

function App() {
  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-800">
      <div className="flex flex-col lg:flex-row gap-8">
        <GameBoard />    // The chessboard
        <GameInfo />     // Game info panel
      </div>
    </div>
  );
}
```

**What it renders:**
- Full-screen dark gray background (`bg-gray-800`)
- Content centered on screen
- Board and info side-by-side on desktop, stacked on mobile

---

### `src/index.css`
**What it is:** Global CSS styles
**What it does:**
- Imports Tailwind CSS (`@import "tailwindcss"`)
- Sets global styles for `html`, `body`, and `#root`
- Ensures the app takes full viewport height
- Prevents scrolling (`overflow: hidden`)

**Think of it as:** The base coat of paint before you add details.

**Current code:**
```css
@import "tailwindcss";           /* Loads Tailwind CSS */

html, body, #root {
  height: 100%;                   /* Full height */
  margin: 0;                      /* No default margins */
  padding: 0;
  overflow: hidden;               /* No scrollbars */
}
```

---

## Components (`src/components/`)

### `src/components/GameBoard.tsx`
**What it is:** The chessboard component
**What it does:**
- Creates an 8×8 grid of squares
- Calculates which squares are light/dark
- Generates square names (a1, h8, etc.)
- Renders 64 `<Square>` components

**Think of it as:** The chessboard itself - the grid structure.

**How it works:**
1. **Define files and ranks:**
   ```typescript
   const files = ['a', 'b', 'c', 'd', 'e', 'f', 'g', 'h'];
   const ranks = ['1', '2', '3', '4', '5', '6', '7', '8'];
   ```

2. **Loop through ranks (8 to 1, top to bottom):**
   ```typescript
   for (let rankIndex = 7; rankIndex >= 0; rankIndex--) {  // rank 8, 7, 6, ..., 1
   ```

3. **Loop through files (a to h, left to right):**
   ```typescript
   for (let fileIndex = 0; fileIndex < 8; fileIndex++) {  // a, b, c, ..., h
   ```

4. **Calculate square color:**
   ```typescript
   // If (rank + file) is odd → light, even → dark
   const isLight = (rankIndex + fileIndex) % 2 !== 0;
   ```

   **Example math:**
   - a1: rank=0, file=0 → (0+0) = 0 (even) → dark ✓
   - a2: rank=1, file=0 → (1+0) = 1 (odd) → light ✓
   - b1: rank=0, file=1 → (0+1) = 1 (odd) → light ✓

5. **Render grid:**
   ```typescript
   <div className="grid grid-cols-8">  // 8 columns
     {squares.map(square => <Square key={square.name} ... />)}
   </div>
   ```

**Why the loops go this way:**
- Rank 8 first → appears at top (chess standard)
- Rank 1 last → appears at bottom
- File 'a' first → left side
- File 'h' last → right side

---

### `src/components/Square.tsx`
**What it is:** A single square on the chessboard
**What it does:**
- Receives props: `squareColor` (light/dark) and `squareName` (e.g., "e4")
- Applies background color based on squareColor
- Displays the square name with low opacity (for debugging)

**Think of it as:** One tile on the chessboard.

**Props interface:**
```typescript
interface SquareProps {
  squareColor: 'light' | 'dark';  // Which color is this square?
  squareName: 'a1' | 'a2' | ... | 'h8';  // Which square is this?
}
```

**How it styles itself:**
```typescript
const bgColor = squareColor === 'light'
  ? 'bg-amber-100'   // Light tan color
  : 'bg-amber-700';  // Dark brown color
```

**Rendering:**
```typescript
<div className={`${bgColor} flex items-center justify-center`}>
  <span className="text-xs opacity-30">
    {squareName}  // Shows "a1", "e4", etc. faintly
  </span>
</div>
```

**Why `opacity-30`?** So you can see square names for debugging but they don't clutter the board.

---

### `src/components/GameInfo.tsx`
**What it is:** Info panel next to the board
**What it does:**
- Currently just a placeholder showing "GameInfo"
- Will later show: whose turn, check status, captured pieces, etc.

**Think of it as:** The scoreboard/status display.

**Current code:**
```typescript
const GameInfo = () => {
  return (
    <div className="w-64 p-4 bg-gray-700">
      <p className="text-white text-lg">GameInfo</p>
    </div>
  );
};
```

**Future additions (later steps):**
- Current player's turn (White/Black)
- Check/Checkmate status
- Move history
- Captured pieces

---

### `src/components/Piece.tsx`
**What it is:** Component to render a chess piece (NOT USED YET)
**What it does:** Will display a chess piece (♔ ♕ ♖ ♗ ♘ ♙) on a square
**Status:** Placeholder from Step 2, will be implemented in Step 4

**Think of it as:** A costume waiting to be worn - created but not used yet.

---

### `src/components/PromotionDialog.tsx`
**What it is:** Dialog for pawn promotion (NOT USED YET)
**What it does:** When a pawn reaches the opposite end, this dialog lets the player choose which piece to promote to (Queen, Rook, Bishop, Knight)
**Status:** Placeholder from Step 2, will be implemented in Step 10

**Think of it as:** A menu that appears when needed - not shown yet.

---

## Hooks (`src/hooks/`)

### `src/hooks/useGameController.ts`
**What it is:** Custom React hook for game logic (NOT USED YET)
**What it does:** Will manage:
- Game state (current position, whose turn, etc.)
- Move validation (is this move legal?)
- Game rules (checkmate, castling, en passant)
- Integration with `chess.js` library

**Think of it as:** The referee and rulebook - knows all the rules and enforces them.

**Status:** Placeholder from Step 2, will be implemented in Step 5

**Future usage:**
```typescript
function GameBoard() {
  const { game, makeMove, validMoves } = useGameController();
  // Now GameBoard can access game state and make moves
}
```

---

## Types (`src/types/`)

### `src/types/chess.ts`
**What it is:** TypeScript type definitions for chess
**What it does:** Defines types to ensure type safety across the app

**Types defined:**

1. **SquareColor**
   ```typescript
   export type SquareColor = 'light' | 'dark';
   ```
   **Use:** Ensures you only pass 'light' or 'dark' to Square component

2. **ChessFile**
   ```typescript
   export type ChessFile = 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h';
   ```
   **Use:** Represents the columns (a-h) on a chessboard

3. **ChessRank**
   ```typescript
   export type ChessRank = '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8';
   ```
   **Use:** Represents the rows (1-8) on a chessboard

4. **Square**
   ```typescript
   export type Square = `${ChessFile}${ChessRank}`;
   ```
   **Use:** Represents a square name like 'e4', 'a1', 'h8'
   **How it works:** TypeScript template literal type combines file + rank
   **Result:** Only allows valid squares (64 possibilities)

**Why types matter:**
```typescript
// TypeScript catches errors at compile time
const square: Square = 'e4';  // ✓ Valid
const bad: Square = 'z9';     // ✗ Error! 'z9' is not a valid square
```

**Think of types as:** A safety net that catches mistakes before your code runs.

---

## How Everything Works Together

### The Flow (Startup to Screen)

1. **Browser loads** `index.html`
2. `index.html` runs `src/main.tsx`
3. `main.tsx` renders `<App />` into `<div id="root">`
4. `<App />` renders `<GameBoard />` and `<GameInfo />`
5. `<GameBoard />`:
   - Calculates 64 squares with colors and names
   - Renders 64 `<Square />` components
6. Each `<Square />`:
   - Receives its color and name as props
   - Displays itself with the right background color
7. **Result:** You see a chessboard on screen!

### Data Flow Diagram

```
index.html
    ↓
main.tsx (starts React)
    ↓
App.tsx (root component)
    ↓
    ├─→ GameBoard.tsx (creates 64 squares)
    │       ↓
    │   [64 × Square.tsx] (individual squares)
    │
    └─→ GameInfo.tsx (info panel)
```

### File Dependencies (Who Imports What)

```
main.tsx
  → App.tsx
      → GameBoard.tsx
          → Square.tsx
              → types/chess.ts
      → GameInfo.tsx

vite.config.ts (Vite configuration)
tsconfig.json (TypeScript settings)
eslint.config.js (Code quality rules)
package.json (Dependencies and scripts)
```

---

## Key Concepts Explained

### What is a Component?

A **component** is a reusable piece of UI. Think of it like a LEGO brick.

**Example:**
- `<Square />` is a component (one tile)
- `<GameBoard />` uses 64 `<Square />` components (the whole board)
- `<App />` uses `<GameBoard />` and `<GameInfo />` (the complete app)

**Code:**
```typescript
// Define a component
function Square({ color, name }) {
  return <div>{name}</div>;
}

// Use it
<Square color="light" name="e4" />
```

---

### What are Props?

**Props** = Properties = Data you pass to a component

**Think of it as:** Function parameters, but for components

**Example:**
```typescript
// Component definition
function Square({ squareColor, squareName }) {
  return <div>{squareName} is {squareColor}</div>;
}

// Using it (passing props)
<Square squareColor="light" squareName="e4" />

// Result: "e4 is light"
```

---

### What is TypeScript?

**TypeScript** = JavaScript + Types

**Why use it?** Catches errors before running code.

**Example:**
```typescript
// JavaScript (no error checking)
function add(a, b) {
  return a + b;
}
add(5, "hello");  // Oops! Returns "5hello" (string concatenation)

// TypeScript (catches error)
function add(a: number, b: number): number {
  return a + b;
}
add(5, "hello");  // ✗ Error: "hello" is not a number!
```

---

### What is Tailwind CSS?

**Tailwind** = CSS classes that style elements

**Instead of writing CSS:**
```css
.my-box {
  width: 100px;
  height: 100px;
  background-color: blue;
}
```

**Use Tailwind classes:**
```html
<div className="w-24 h-24 bg-blue-500"></div>
```

**Common classes:**
- `w-24` = width: 96px (24 × 4px)
- `h-24` = height: 96px
- `bg-blue-500` = blue background
- `flex` = display: flex
- `items-center` = align items vertically center
- `justify-center` = align items horizontally center

---

### What is Vite?

**Vite** = Build tool that:
1. Runs a development server (`pnpm dev`)
2. Hot-reloads when you save files (instant updates)
3. Bundles your code for production (`pnpm build`)

**Think of it as:** A smart assistant that watches your files and rebuilds your app instantly when you make changes.

---

## Current Implementation Status

### ✅ Completed (Step 1-3)
- Project setup with Vite, React, TypeScript
- Tailwind CSS v4 configured
- Component scaffolding (all files created)
- Static chessboard rendering (8×8 grid with colors)
- Type definitions for squares

### 🚧 Not Yet Implemented
- Piece rendering (Step 4)
- Game logic with chess.js (Step 5)
- Move implementation (Step 6-7)
- Game info display (Step 8-9)
- Pawn promotion (Step 10)
- Checkmate detection (Step 11)

---

## Quick Reference

### Run the App
```bash
pnpm dev          # Start development server
pnpm build        # Build for production
pnpm lint         # Check code quality
```

### File Naming Conventions
- `.tsx` = TypeScript + JSX (React components)
- `.ts` = TypeScript (no JSX)
- `.json` = Configuration files
- `.css` = Stylesheets

### Import Rules (with `verbatimModuleSyntax`)
```typescript
// ✓ Correct - types use `import type`
import type { Square } from './types/chess';

// ✗ Wrong - will cause TypeScript error
import { Square } from './types/chess';

// ✓ Correct - components use regular import
import GameBoard from './components/GameBoard';
```

---

## Learning Path

If you're new to this stack, learn in this order:

1. **HTML/CSS basics** → Understand structure and styling
2. **JavaScript (ES6+)** → Arrow functions, destructuring, modules
3. **React basics** → Components, props, JSX
4. **TypeScript** → Types, interfaces, type safety
5. **Tailwind CSS** → Utility-first styling
6. **Vite** → Build tool and dev server

**This project uses all of these together!**

---

## Troubleshooting

### "Module not found" error
- Check import path is correct
- Use `../` to go up one directory
- File extensions optional in imports

### TypeScript error about types
- Use `import type` for types when `verbatimModuleSyntax: true`

### Tailwind classes not working
- Check `@import "tailwindcss"` is in `index.css`
- Verify `tailwindcss` plugin in `vite.config.ts`

### Changes not showing in browser
- Hard refresh: Cmd+Shift+R (Mac) or Ctrl+Shift+R (Windows)
- Check Vite dev server is running
- Clear browser cache

---

**Next Steps:** Move to Step 4 (Piece Rendering) to add chess pieces to the board!

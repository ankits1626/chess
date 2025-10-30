# Step 10: Developer Debug Toolkit - Architecture & Design

## Goal
Create a non-intrusive, developer-focused debugging toolkit to test specific chess scenarios (promotion, checkmate, custom positions) without modifying source code or impacting production builds.

---

## Architectural Decisions

### 1. Activation & Access

**Recommended Approach: Multi-layered Protection**

```typescript
// Combine multiple activation methods for flexibility:

// Method 1: Environment-based (primary)
if (import.meta.env.DEV) {
  // Debug tools only available in dev mode
}

// Method 2: Keyboard shortcut (Ctrl+Shift+D)
// Easy for developers to toggle on/off

// Method 3: URL query parameter (backup)
// Useful for sharing debug sessions: ?debug=true
```

**Implementation:**
1. **Primary Guard:** `import.meta.env.DEV` ensures complete removal from production
2. **Toggle:** Keyboard shortcut shows/hides panel within dev environment
3. **Persistent State:** URL param allows linking to debug mode

**Why This Approach:**
- ✅ Zero production overhead (code-split, tree-shaken)
- ✅ Flexible for different workflows
- ✅ Can't be accidentally enabled in production
- ✅ Easy to use during development

---

### 2. Core Features

**Recommended Feature Set:**

#### Priority 1: Essential (MVP)
```typescript
interface DebugToolkit {
  // FEN Loader - Core feature
  loadFen: (fen: string) => void;

  // Pre-defined scenarios - Developer productivity
  loadScenario: (scenarioId: string) => void;

  // Reset - Quick recovery
  resetToStart: () => void;
}
```

**Pre-defined Scenarios:**
```typescript
const scenarios = {
  'pawn-promotion-white': '8/4k3/8/8/8/8/4P3/4K3 w - - 0 1',
  'pawn-promotion-black': '4k3/4p3/8/8/8/8/4K3/8 b - - 0 1',
  'checkmate-test': '6k1/5ppp/8/8/8/8/5PPP/R5K1 w - - 0 1', // Back rank mate
  'stalemate-test': '7k/5Q2/5K2/8/8/8/8/8 b - - 0 1',
  'castling-test': 'r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1',
  'en-passant-test': 'rnbqkbnr/ppp2ppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3',
  'complex-midgame': 'r1bqk2r/pppp1ppp/2n2n2/2b1p3/2B1P3/3P1N2/PPP2PPP/RNBQK2R w KQkq - 0 6'
};
```

#### Priority 2: Advanced (Future)
```typescript
interface AdvancedDebugToolkit extends DebugToolkit {
  // State inspection
  exportFen: () => string;
  exportPgn: () => string;

  // Move manipulation
  undoMove: () => void;
  redoMove: () => void;
  gotoMove: (moveNumber: number) => void;

  // Board manipulation
  flipBoard: () => void;
  changeTurn: () => void;
}
```

---

### 3. UI/UX Design

**Recommended: Floating Draggable Panel**

**Why:**
- ✅ Non-intrusive (doesn't push content)
- ✅ Flexible positioning
- ✅ Can be minimized
- ✅ Doesn't interfere with game testing

**Design Mockup:**
```
┌─────────────────────────────────────────┐
│  Chess Coach              [Debug] 🛠️   │
├───────────────────────────────────────┬─┤
│                           │           │━│ ← Minimize
│   [Chess Board]           │  Game     │ │
│                           │  Info     ├─┤
│                           │           │×│ ← Close
├───────────────────────────┴───────────┴─┤
│ ┌─ Debug Toolkit ──────────────────┐   │
│ │                                   │   │
│ │ FEN: [8/4k3/8/8/8/8/4P3/4K3...] │   │
│ │                          [Load]   │   │
│ │                                   │   │
│ │ Quick Scenarios:                  │   │
│ │ [Pawn Promo] [Checkmate] [Castle]│   │
│ │                                   │   │
│ │ [Reset] [Export FEN] [Export PGN]│   │
│ └───────────────────────────────────┘   │
└─────────────────────────────────────────┘
```

**Component Structure:**
```typescript
<DebugPanel>
  <DebugHeader>
    <h3>🛠️ Debug Toolkit</h3>
    <MinimizeButton />
    <CloseButton />
  </DebugHeader>

  <DebugContent>
    <FenLoader onLoad={loadFen} />
    <ScenarioPicker scenarios={scenarios} onSelect={loadScenario} />
    <QuickActions>
      <Button onClick={resetToStart}>Reset</Button>
      <Button onClick={exportFen}>Export FEN</Button>
      <Button onClick={exportPgn}>Export PGN</Button>
    </QuickActions>
  </DebugContent>
</DebugPanel>
```

**Styling:**
```typescript
// Matches app aesthetic (gray-800 theme)
// Semi-transparent backdrop when minimized
// Smooth animations (slide in/out)
// Drag handle for repositioning
```

---

### 4. Architectural Integration

**Recommended: Expose Control Functions via Render Props (Current Phase)**

**Why:**
- ✅ Consistent with existing render props pattern
- ✅ No new dependencies
- ✅ Simple to implement
- ✅ Easy to migrate to Zustand later

**Implementation:**

#### GameController.tsx Changes:
```typescript
interface GameControllerProps {
  children: (
    game: Chess,
    selectedSquare: Square | null,
    validMoves: Square[],
    selectSquare: (square: Square) => void,
    pendingMove: PendingPromotion | null,
    handlePromotion: (piece: PieceType) => void,
    // NEW: Debug functions (only in dev mode)
    debugActions?: DebugActions
  ) => React.ReactNode;
}

interface DebugActions {
  loadFen: (fen: string) => void;
  resetGame: () => void;
  exportFen: () => string;
  exportPgn: () => string;
}

const GameController = ({ children }: GameControllerProps) => {
  // ... existing state ...

  // Debug actions (only created in dev mode)
  const debugActions: DebugActions | undefined = import.meta.env.DEV ? {
    loadFen: (fen: string) => {
      try {
        const newGame = new Chess(fen);
        setGame(newGame);
        setSelectedSquare(null);
        setValidMoves([]);
        setPendingMove(null);
      } catch (error) {
        console.error('Invalid FEN:', error);
      }
    },

    resetGame: () => {
      setGame(new Chess());
      setSelectedSquare(null);
      setValidMoves([]);
      setPendingMove(null);
    },

    exportFen: () => game.fen(),

    exportPgn: () => game.pgn()
  } : undefined;

  return <>{children(
    game,
    selectedSquare,
    validMoves,
    selectSquare,
    pendingMove,
    handlePromotion,
    debugActions
  )}</>;
};
```

#### App.tsx Changes:
```typescript
import { lazy, Suspense } from 'react';

// Dynamic import - only loaded in dev mode
const DebugPanel = import.meta.env.DEV
  ? lazy(() => import('./components/DebugPanel'))
  : null;

function App() {
  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-800 text-white p-8">
      <GameController>
        {(game, selectedSquare, validMoves, selectSquare, pendingMove, handlePromotion, debugActions) => (
          <>
            <div className="flex flex-row gap-8 items-center">
              <GameBoard
                game={game}
                selectedSquare={selectedSquare}
                validMoves={validMoves}
                onSquareClick={selectSquare}
              />
              <GameInfo game={game} />
            </div>

            {pendingMove && (
              <PromotionDialog
                color={pendingMove.color}
                onSelectPiece={handlePromotion}
              />
            )}

            {/* Debug panel - only in dev mode */}
            {import.meta.env.DEV && DebugPanel && debugActions && (
              <Suspense fallback={null}>
                <DebugPanel actions={debugActions} currentFen={game.fen()} />
              </Suspense>
            )}
          </>
        )}
      </GameController>
    </div>
  );
}
```

**Migration Path to Zustand (Phase 2):**
```typescript
// Future Zustand store will expose debug actions directly:
interface ChessCoachStore {
  // ... game state ...

  // Debug actions (only in store during dev)
  debugActions: import.meta.env.DEV ? DebugActions : never;
}

// Debug panel will use store directly:
const DebugPanel = () => {
  const { loadFen, resetGame, exportFen } = useChessCoachStore(
    state => state.debugActions
  );
  // ...
};
```

---

### 5. Production Build Protection

**Multi-layered Protection Strategy:**

#### Layer 1: Compile-time Removal (Primary)
```typescript
// Vite automatically removes this code in production builds
if (import.meta.env.DEV) {
  // All debug code here is tree-shaken in production
}
```

#### Layer 2: Dynamic Imports (Secondary)
```typescript
// Component is never loaded in production
const DebugPanel = import.meta.env.DEV
  ? lazy(() => import('./components/DebugPanel'))
  : null;
```

#### Layer 3: Type Safety (Tertiary)
```typescript
// TypeScript prevents debug actions in production type
debugActions?: import.meta.env.DEV ? DebugActions : never;
```

**Verification:**
```bash
# Check production bundle doesn't include debug code
pnpm build
grep -r "DebugPanel" dist/  # Should be empty
grep -r "loadFen" dist/     # Should be empty
```

**Result:**
- ✅ Zero bytes of debug code in production
- ✅ Zero runtime overhead
- ✅ Impossible to enable accidentally
- ✅ Type-safe

---

## Implementation Plan

### Phase 1: MVP (Essential Features)

**Components to Create:**

1. **`DebugPanel.tsx`** - Main container
   ```typescript
   interface DebugPanelProps {
     actions: DebugActions;
     currentFen: string;
   }
   ```

2. **`FenLoader.tsx`** - FEN input with validation
   ```typescript
   interface FenLoaderProps {
     onLoad: (fen: string) => void;
     currentFen: string;
   }
   ```

3. **`ScenarioPicker.tsx`** - Pre-defined scenarios
   ```typescript
   interface ScenarioPickerProps {
     scenarios: Record<string, string>;
     onSelect: (fen: string) => void;
   }
   ```

**Files to Modify:**
1. `GameController.tsx` - Add `debugActions` to render props
2. `App.tsx` - Conditionally render `DebugPanel`

**Keyboard Hook:**
```typescript
// src/hooks/useDebugShortcut.ts
export const useDebugShortcut = (callback: () => void) => {
  useEffect(() => {
    if (!import.meta.env.DEV) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      // Ctrl+Shift+D or Cmd+Shift+D
      if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.key === 'D') {
        e.preventDefault();
        callback();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [callback]);
};
```

---

### Phase 2: Enhanced Features (Future)

After Phase 1 is stable, add:
1. **Move navigation** - Undo/redo, goto move
2. **Board manipulation** - Flip board, change turn
3. **PGN loader** - Load complete games
4. **Position analysis** - Show legal moves, threats
5. **Draggable positioning** - Use `react-draggable` library

---

## Testing Strategy

### Manual Testing Scenarios

1. **FEN Loader:**
   - ✅ Load valid FEN - board updates
   - ✅ Load invalid FEN - shows error, doesn't crash
   - ✅ Load pawn promotion scenario - test promotion
   - ✅ Load checkmate scenario - verify checkmate detection

2. **Scenario Picker:**
   - ✅ Click each scenario - loads correctly
   - ✅ Test promotion scenario - dialog appears
   - ✅ Test stalemate scenario - draw detected

3. **Production Build:**
   - ✅ Build for production - no debug code in bundle
   - ✅ Run production build - no debug panel visible
   - ✅ Check bundle size - no increase from debug code

4. **Keyboard Shortcut:**
   - ✅ Press Ctrl+Shift+D - panel toggles
   - ✅ Press again - panel closes
   - ✅ In production build - shortcut does nothing

---

## Security Considerations

**Concern:** Could debug tools be exploited in production?

**Mitigation:**
1. ✅ **Compile-time removal** - Code doesn't exist in production bundle
2. ✅ **Type safety** - TypeScript prevents accidental usage
3. ✅ **Tree shaking** - Vite removes unused imports
4. ✅ **Dynamic imports** - Component never loaded in production
5. ✅ **Environment checks** - Multiple layers of `import.meta.env.DEV` guards

**Result:** Impossible to enable in production, even with browser dev tools.

---

## Developer Experience Benefits

1. **Testing Promotion:**
   ```
   Before: Manually play 15-20 moves to get pawn to 8th rank
   After: Click "Pawn Promotion" scenario → immediate test
   Time saved: 2-3 minutes per test
   ```

2. **Testing Checkmate:**
   ```
   Before: Set up complex position manually or play full game
   After: Click "Checkmate Test" scenario → verify logic
   Time saved: 5-10 minutes per test
   ```

3. **Reproducing Bug Reports:**
   ```
   User reports: "Bug at position rnbqkbnr/pppppppp/..."
   Developer: Copy FEN → Load → Reproduce bug
   Time saved: Immediate reproduction vs manual setup
   ```

4. **QA Testing:**
   ```
   Test suite: Run through all scenarios systematically
   Coverage: Every edge case tested in seconds
   Confidence: All chess rules verified
   ```

---

## Comparison of Architectural Options

| Approach | Pros | Cons | Recommendation |
|----------|------|------|----------------|
| **Render Props** (Current) | ✅ Consistent pattern<br>✅ Simple<br>✅ No new deps | ⚠️ Prop drilling<br>⚠️ Verbose | ✅ Use for Phase 1 |
| **Zustand** (Future) | ✅ Clean access<br>✅ Scalable<br>✅ Phase 2 ready | ❌ Not implemented yet | ✅ Use for Phase 2 |
| **Event Bus** | ✅ Decoupled | ❌ Complex<br>❌ Hard to debug | ❌ Overkill |
| **Context API** | ✅ Built-in | ❌ Boilerplate<br>❌ Performance | ❌ Worse than Zustand |

---

## File Structure

```
src/
├── components/
│   ├── debug/                    # All debug components (dev only)
│   │   ├── DebugPanel.tsx       # Main container
│   │   ├── FenLoader.tsx        # FEN input
│   │   ├── ScenarioPicker.tsx   # Scenario dropdown
│   │   └── index.ts             # Barrel export
│   ├── GameController.tsx        # Modified for debug actions
│   └── ...
├── hooks/
│   └── useDebugShortcut.ts      # Keyboard shortcut hook
├── constants/
│   └── debugScenarios.ts        # Pre-defined positions
└── types/
    └── debug.ts                  # Debug-related types
```

---

## Summary

**Recommended Architecture:**

1. ✅ **Activation:** Keyboard shortcut (Ctrl+Shift+D) + `import.meta.env.DEV` guard
2. ✅ **Features:** FEN loader + Pre-defined scenarios + Reset (MVP)
3. ✅ **UI:** Floating draggable panel (minimizable, closable)
4. ✅ **Integration:** Render props for Phase 1, migrate to Zustand in Phase 2
5. ✅ **Production:** Multi-layered protection, complete removal from bundle

**Next Step:** Implement MVP (Step 10) with FEN loader and scenario picker.

---

## Example Usage Flow

```typescript
// Developer workflow:
1. Press Ctrl+Shift+D
2. Debug panel appears
3. Click "Pawn Promotion" button
4. Board shows: 8/4k3/8/8/8/8/4P3/4K3 w - - 0 1
5. Click e2 → e3 → e4 → e5 → e6 → e7 → e8
6. Promotion dialog appears
7. Test all 4 promotion options
8. Click "Reset" button
9. Back to starting position
10. Press Ctrl+Shift+D to close panel
```

**Total time:** 30 seconds vs 5+ minutes without toolkit.

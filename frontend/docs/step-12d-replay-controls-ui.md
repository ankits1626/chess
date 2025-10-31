# Step 12d: Replay Controls UI & UX Plan (Revised)

This document outlines the plan for creating the user interface for controlling game replays, incorporating feedback from the review on 2025-10-31.

## 1. Component Creation

*   A new component, `ReplayControls.tsx`, will be created in `frontend/app/src/components/`.
*   This component will be a presentational component responsible for rendering the UI and handling user interactions for replay control.

## 2. UI Elements

The `ReplayControls` component will include the following interactive elements:

*   **Go to Start:** Button to reset the replay to the initial board position.
*   **Previous Move:** Button to step back one move.
*   **Play/Pause:** Button to toggle automatic playback.
*   **Next Move:** Button to advance one move.
*   **Go to End:** Button to jump to the final position.

## 3. State Management & Props

The component will be driven by props from a parent container.

### Props Interface

```typescript
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
```

### Zustand Store (`useGameStore`) Modifications

The store will be updated to support the replay functionality.

*   **State:**
    *   `isAutoplaying: boolean`: Tracks if autoplay is active.
    *   `autoplayIntervalId: NodeJS.Timeout | null`: Stores the interval ID for cleanup.
    *   Selectors will determine `isFirstMove` and `isLastMove`.

*   **Actions:**
    *   The existing actions (`goToFirstMove`, `prevMove`, `nextMove`, `goToLastMove`) will be used.
    *   New actions for autoplay will be added:
        *   `toggleAutoplay()`: Switches between play and pause.
        *   `startAutoplay()`: Starts a `setInterval` that calls `nextMove()` every 1000ms. It will not start if the game is at the end. It will store the interval ID.
        *   `stopAutoplay()`: Clears the interval and updates the state. Autoplay should also stop when a user manually navigates moves or switches game modes.

## 4. Implementation Steps

1.  **Create `ReplayControls.tsx`:** Create the file with the revised props interface.
2.  **Build the UI:**
    *   Use Tailwind CSS for styling and layout.
    *   Use simple text labels (`<<`, `<`, `Play`, `>`, `>>`).
    *   Implement disabled states using the `disabled` attribute based on `isFirstMove` and `isLastMove` props.
3.  **Connect onClick Handlers:** Wire up button clicks to the corresponding prop functions (`onGoToFirst`, `onPrevious`, etc.).
4.  **Update `useGameStore`:**
    *   Add `isAutoplaying` and `autoplayIntervalId` to the state.
    *   Implement `toggleAutoplay`, `startAutoplay`, and `stopAutoplay` actions.
5.  **Integrate into `App.tsx`:**
    *   Render `<ReplayControls />` conditionally based on the existing `mode === 'replay'` state from `useGameStore`.
    *   Pass the required state and actions from the store to the component as props.

## 5. Styling

*   The control bar will be centered horizontally below the `GameBoard`.
*   Buttons will have a clean design with hover, active, and disabled states.

## 6. Move List Component (Optional Enhancement)

*   Create a `MoveList.tsx` component to display all moves in algebraic notation.
*   It will highlight the current move and allow clicking on a move to jump to that position in the replay.

## 7. Accessibility (a11y)

*   **ARIA Labels:** Add descriptive `aria-label` attributes to all control buttons (e.g., `aria-label="Next move"`).
*   **Keyboard Shortcuts:** Implement keyboard controls:
    *   `ArrowLeft`: Previous move
    *   `ArrowRight`: Next move
    *   `Space`: Toggle play/pause
    *   `Home`: First move
    *   `End`: Last move

## 8. Testing Checklist

- [ ] All buttons work correctly.
- [ ] Disabled states prevent clicks when appropriate.
- [ ] Autoplay starts and stops correctly.
- [ ] Autoplay stops automatically at the end of the game.
- [ ] Clicking navigation during autoplay stops autoplay.
- [ ] Switching from replay to live mode stops autoplay.
- [ ] Component unmounting clears the autoplay interval.
- [ ] Keyboard shortcuts work as expected.
- [ ] Screen readers can navigate and understand the controls.
- [ ] Play/Pause button label updates correctly.
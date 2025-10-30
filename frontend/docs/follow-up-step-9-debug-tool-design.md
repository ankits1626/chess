# Follow-up to Step 9: In-App Debugging & Testing Tools Design

This document captures the architectural design questions for creating a developer-focused toolkit within the application. The objective is to facilitate testing of specific game scenarios (e.g., pawn promotion, checkmate) without altering the source code temporarily.

---

### Questionnaire for Architect: In-App Debugging & Testing Tools

**Objective:** To design a non-intrusive, developer-focused toolkit within the application to facilitate testing of specific game scenarios, such as pawn promotion, checkmate, or custom positions.

**1. Activation & Access:**
*   How should the debug panel be activated?
    *   A specific keyboard shortcut (e.g., `Ctrl+Shift+D`)?
    *   A query parameter in the URL (e.g., `?debug=true`)?
    *   A small, unobtrusive button in a corner of the screen?
    *   Should this be a compile-time flag, completely removed from production builds?

**2. Core Features:**
*   **FEN String Loader:** This seems to be the most critical feature.
    *   Should it be a simple text input where a developer can paste a FEN string?
    *   Should we include a dropdown or list of pre-defined scenarios (e.g., "Pawn Promotion Test," "Checkmate Puzzle," "Stalemate Test")?
*   **Game State Control:**
    *   Should there be buttons for "Reset to Standard," "Flip Board," or "Change Turn"?
*   **Move History Manipulation:**
    *   Would it be useful to be able to clear the move history or load a game from a PGN string?

**3. UI/UX of the Debug Panel:**
*   How should the panel be presented?
    *   As a modal dialog that overlays the entire app?
    *   As a collapsible sidebar that pushes the main content?
    *   As a floating "widget" that can be moved around the screen?
*   What should the UI be? Simple HTML controls, or should it match the app's aesthetic?

**4. Architectural Integration:**
*   How should the debug panel communicate with the `GameController`?
    *   **Option A (Prop Drilling):** Should the `App` component manage the debug state and pass down a `setGame` function to the debug panel? This is simple but could lead to prop drilling.
    *   **Option B (State Management):** If/when we move to Zustand (as planned for Phase 2), the debug panel could directly interact with the global store (e.g., `useChessCoachStore.getState().setGame(newFen)`). This is cleaner and more decoupled.
    *   **Option C (Event Bus):** A simple event emitter/subscriber pattern could be used, where the debug panel emits a `LOAD_FEN` event and the `GameController` subscribes to it. This avoids direct coupling.
*   Given our current "render props" pattern, how would you recommend integrating this without making the `GameController`'s props overly complex?

**5. Production Builds:**
*   How do we ensure this entire feature is stripped from production builds to avoid impacting end-users and prevent potential cheating or state manipulation?
    *   Using environment variables (e.g., `if (import.meta.env.DEV) { ... }`)?
    *   Using a dynamic `import()` for the debug panel component so it's only loaded in development?

**Example Scenario Walkthrough:**
A developer wants to test pawn promotion.
1.  They press `Ctrl+Shift+D`.
2.  A small panel appears.
3.  They paste `8/4k3/8/8/8/8/4P3/4K3 w - - 0 1` into a text field and click "Load FEN".
4.  The debug panel calls a function (e.g., `loadFen(newFen)`).
5.  This function updates the `game` state in `GameController`, causing the board to re-render to the specified position.
6.  The developer can now test the e7->e8 move.

---

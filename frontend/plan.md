# Chess Coach Frontend Development Plan

This document outlines the plan for developing the frontend of the Chess Coach application. The project will be built with React and TypeScript, adhering to SOLID principles to ensure a maintainable and scalable codebase.

## 1. Guiding Principles (SOLID)

We will apply SOLID principles to the React component architecture:

*   **Single Responsibility Principle (SRP):** Each component or module will have one primary responsibility. For example, one component will render the board, while another will manage game logic.
*   **Open/Closed Principle (OCP):** Components will be open for extension but closed for modification. We can add new features (like different game modes) without rewriting existing components.
*   **Liskov Substitution Principle (LSP):** Component hierarchies will be designed to be substitutable. For instance, different piece components could be used interchangeably if they share a common interface.
*   **Interface Segregation Principle (ISP):** Components will not be forced to depend on props they do not use. We will keep our component props lean and specific.
*   **Dependency Inversion Principle (DIP):** High-level components will not depend on low-level components. Both will depend on abstractions. We will use hooks and context to invert dependencies.

## 2. Tech Stack

*   **Framework:** React
*   **Language:** TypeScript
*   **Build Tool:** Vite
*   **Chess Logic:** `chess.js`
*   **Styling:** Tailwind CSS

## 3. Dependency Versions

To ensure a stable and reproducible development environment, we will use the following latest production-ready versions for our core dependencies:

*   **Node.js:** `24.x` (LTS)
*   **React:** `19.2.0`
*   **Vite:** `7.1.12`
*   **TypeScript:** `5.9.3`
*   **chess.js:** `1.4.0`

The versions for other tools, such as ESLint and testing libraries, will be determined by the Vite `react-ts` template, and we will pin them in `package.json` after the initial setup.

## 4. Project Setup (Step 1)

*Note: Steps 1-3 are already complete.*

1.  **Scaffold the project with Vite using pnpm:** (Complete)

2.  **Navigate into the project directory:** (Complete)

3.  **Install dependencies with pnpm:** (Complete)

4.  **Remove outdated Tailwind dependencies:** The initial attempt installed packages for an older version of Tailwind. They must be removed.
    ```bash
    pnpm remove tailwindcss postcss autoprefixer
    ```
5.  **Install Tailwind CSS v4 for Vite:**
    ```bash
    pnpm add -D tailwindcss@next @tailwindcss/vite@next
    ```
6.  **Configure Vite for Tailwind:** Modify `frontend/app/vite.config.ts` to include the `@tailwindcss/vite` plugin.

7.  **Import Tailwind Styles:** Replace the content of `frontend/app/src/index.css` with the correct import directive.

## 5. Component Architecture (Step 2)

We will structure the application into the following components:

*   `App.tsx`: The root component, responsible for the overall layout and composition of the application.
*   `GameController.tsx`: A component that encapsulates the core game logic and state management. It will interact with `chess.js` and provide the game state to the rest of the application. This component will not render any UI itself.
*   `GameBoard.tsx`: Responsible for rendering the 8x8 grid of the chessboard. It will receive the board layout and piece positions as props from `GameController`.
*   `Square.tsx`: Represents a single square on the board. It will display its color and any piece on it. It will also handle user click events and delegate them to the `GameController`.
*   `Piece.tsx`: A presentational component that renders a chess piece based on its type (pawn, rook, etc.) and color (white or black).
*   `PromotionDialog.tsx`: A dialog that appears when a pawn reaches the promotion rank, allowing the user to select a new piece.
*   `GameInfo.tsx`: A component to display game information like the current turn, check status, and game over conditions.

## 6. Development Steps

The implementation will be done in the following order:

1.  **Project Setup:** Complete the steps outlined in section 4.
2.  **Component Scaffolding:** Create the initial files for all the components listed in section 5.
3.  **Static Board Rendering:** Implement `GameBoard.tsx` and `Square.tsx` to render a static, empty chessboard.
4.  **Piece Rendering:** Implement `Piece.tsx` and update `GameBoard.tsx` to render the initial position of the pieces on the board. The piece positions will be hardcoded initially.
5.  **Game Logic Integration:** Implement `GameController.tsx` to use `chess.js`. It will manage the game state and expose it to other components.
6.  **Connecting Logic and UI:** The `App.tsx` will use `GameController.tsx` and pass the game state down to `GameBoard.tsx`.
7.  **Move Implementation (Click-to-move):**
    *   When a square is clicked, the `Square.tsx` component will notify the `GameController.tsx`.
    *   `GameController.tsx` will handle the logic for selecting a piece and making a move.
    *   The game state will be updated, and the UI will re-render to show the new piece positions.
8.  **Valid Move Highlighting:** When a piece is selected, `GameController.tsx` will calculate the valid moves for that piece, and `GameBoard.tsx` will highlight the corresponding squares.
9.  **Game Information Display:** Implement `GameInfo.tsx` to display the current player's turn and other game state information.
10. **Pawn Promotion:** Implement the `PromotionDialog.tsx` and the logic in `GameController.tsx` to handle pawn promotion.
11. **End of Game Logic:** Implement logic in `GameController.tsx` to detect checkmate and stalemate, and display the result in `GameInfo.tsx`.
12. **Styling and Responsiveness:** Apply Tailwind CSS to all components to create a polished and responsive design.

This step-by-step plan will allow us to build the application incrementally and ensure that each part is well-designed and tested.
# Step 12: AI Coach Integration - Plan & Architectural Questions

**Goal:** To implement the core "Chess Coach" feature by integrating a chat interface that allows users to communicate with an AI assistant about the current game.

This marks the beginning of **Phase 2**. As decided, this is the appropriate time to refactor our state management to a more scalable solution (Zustand) and build out the UI and services required for the AI chat.

---

## 1. High-Level Plan

1.  **State Management Refactor:**
    *   Replace the `GameController.tsx` render props pattern with a global **Zustand** store.
    *   The store will manage the `chess.js` instance, selected square, valid moves, last move, and pending promotion state.
    *   It will also be designed to hold future chat-related state (messages, loading status).

2.  **Backend Service (Stub):**
    *   Create a simple backend service (e.g., a Node.js/Express server) with a single endpoint: `POST /api/chat`.
    *   Initially, this endpoint will not connect to a real LLM. It will receive the game state (FEN, PGN) and a user message, and return a hardcoded or simple rule-based response (e.g., "That's an interesting move for the {piece} on {square}!"). This allows us to build and test the full frontend-to-backend pipeline without needing live API keys.

3.  **Frontend Chat UI:**
    *   Create the necessary React components for a chat interface:
        *   `ChatPanel.tsx`: The main container for the chat window.
        *   `ChatMessage.tsx`: To display individual messages from the user and the AI.
        *   `ChatInput.tsx`: A text input field and send button.
    *   Integrate these components into `App.tsx`.

4.  **API Client Service:**
    *   Create a service file in the frontend (e.g., `src/services/api.ts`) to handle communication with the backend `POST /api/chat` endpoint.

---

## 2. Architectural Questions for the Architect

Before proceeding with implementation, we need guidance on the following key architectural decisions.

### Question 1: Zustand Store Structure

How should we structure the Zustand store for optimal organization and performance? The store needs to manage both game logic and chat logic.

*   **Option A: Monolithic Store:** A single store for everything.

    ```typescript
    interface ChessCoachStore {
      // Game State
      game: Chess;
      selectedSquare: Square | null;
      // ... etc.

      // Chat State
      messages: ChatMessage[];
      isAiLoading: boolean;

      // Game Actions
      selectSquare: (square: Square) => void;
      // ... etc.

      // Chat Actions
      sendMessage: (prompt: string) => Promise<void>;
    }
    ```

*   **Option B: Sliced Store:** Use Zustand's slice pattern to separate concerns. This is generally recommended for larger stores.

    ```typescript
    // src/store/gameSlice.ts
    interface GameSlice {
      game: Chess;
      // ... game state & actions
    }

    // src/store/chatSlice.ts
    interface ChatSlice {
      messages: ChatMessage[];
      // ... chat state & actions
    }

    // src/store/index.ts (combined store)
    const useChessCoachStore = create<GameSlice & ChatSlice>()((...a) => ({
      ...createGameSlice(...a),
      ...createChatSlice(...a),
    }));
    ```

**Recommendation Request:** Is the slice pattern the right approach here, or is a monolithic store sufficient for this application's scope?

### Question 2: Backend API Contract

What is the ideal data structure for the `POST /api/chat` request? We need to send enough context for the LLM to provide meaningful advice.

*   **Proposal A (Minimal):**
    ```json
    {
      "fen": "r1bqkbnr/pppp1ppp/...",
      "userMessage": "Why was that a good move?"
    }
    ```

*   **Proposal B (With History):**
    ```json
    {
      "fen": "...",
      "pgn": "1. e4 e5 2. Nf3 ...",
      "lastMove": { "from": "g1", "to": "f3" },
      "userMessage": "Why was that a good move?"
    }
    ```

*   **Proposal C (With Full Chat Context):**
    ```json
    {
      "pgn": "...",
      "chatHistory": [
        { "role": "user", "content": "..." },
        { "role": "assistant", "content": "..." }
      ],
      "userMessage": "Can you elaborate on that?"
    }
    ```

**Recommendation Request:** Which proposal provides the best balance of necessary context without being overly verbose? Should the frontend or backend be responsible for managing and trimming the `chatHistory` to fit context window limits?

### Question 3: Backend Technology Stack & Location

Where should the backend service live, and what stack should it use?

*   **Location:**
    *   **Option A: Monorepo:** Create a `backend/` directory at the root of the `chess-coach` project.
    *   **Option B: Separate Repository:** Create a new, separate Git repository for the backend service.

*   **Technology Stack:**
    *   **Node.js with Express/Fastify:** Good for TypeScript synergy and rapid development.
    *   **Python with FastAPI/Flask:** Excellent for AI/ML integrations and data science libraries.

**Recommendation Request:** What is the preferred location and technology stack for this backend service, keeping future LLM integration in mind?

### Question 4: Real-time UI Updates (Streaming Responses)

LLM responses can take several seconds to generate. To improve UX, we should stream the response word-by-word, like ChatGPT.

*   **Question:** Should we plan for this now or implement it later?
    *   **If now:** This requires the backend to use a streaming response format (e.g., Server-Sent Events or streaming over a standard `fetch` request). The frontend API client will need to be able to handle and parse this stream.
    *   **If later:** We can start with a simple request-response model and add streaming as a future enhancement.

**Recommendation Request:** Is it worth the initial complexity to build for streaming responses from the start, or should we stick to a simpler model for the MVP of Phase 2?

---

## 3. Implementation Plan (Post-Architect Feedback)

Once we have answers to the architectural questions, the implementation will proceed as follows:

1.  **Setup Backend:** Create the new backend directory/repo and scaffold the server with the chosen technology.
2.  **Implement Stub API:** Create the `POST /api/chat` endpoint that returns a hardcoded response.
3.  **Refactor to Zustand:** Install Zustand and create the store structure based on the architect's recommendation. Replace `GameController.tsx` with calls to the new `useChessCoachStore` hook in all relevant components (`App`, `GameBoard`, `GameInfo`).
4.  **Build Chat UI:** Create the `ChatPanel`, `ChatMessage`, and `ChatInput` components.
5.  **Integrate Chat:** Add the chat components to `App.tsx` and connect the `ChatInput` to a new `sendMessage` action in the Zustand store. This action will call the backend API and update the chat history with both the user's message and the AI's (stubbed) response.

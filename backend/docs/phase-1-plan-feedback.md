# Phase 1 Plan - Revised Review Feedback

**Overall Recommendation:** The initial plan is a good start, but it includes too much scope for a foundational Phase 1. The following feedback has been revised to recommend a **leaner initial phase** focused on getting the absolute minimum running, deferring framework choices and other features to Phase 2.

---

## High-Priority Corrections for Phase 1 Plan

These points must be addressed in the plan before implementation begins.

### 1. Docker Compose Location & Monorepo Strategy
*   **Correction:** The `docker-compose.yml` file should be located at the **project root (`/`)**, not within the `backend/` directory.
*   **Reasoning:** For a monorepo containing both a frontend and backend, a root-level `docker-compose.yml` is the best practice. It allows us to manage the *entire* application stack (frontend, backend, database, etc.) with a single `docker-compose up` command, which is ideal for local development.

### 2. Go Module Name
*   **Correction:** The `go mod init` command should use a simple, non-URL-based name appropriate for a monorepo.
*   **Recommendation:** Use `go mod init chess-coach/backend`.

### 3. Docker Volume Mount Path
*   **Correction:** The `docker-compose.yml` file must explicitly define the volume mount path from the root to the backend source code to enable hot-reloading.
*   **Required `docker-compose.yml` snippet:**
    ```yaml
    services:
      backend:
        build:
          context: ./backend
        volumes:
          - ./backend:/app
    ```

### 4. Environment Variable Loading
*   **Recommendation:** For the initial phase, use the lightweight `github.com/joho/godotenv` library instead of `viper`.
*   **Reasoning:** `viper` is a large dependency that is overkill for simply loading a `.env` file. `godotenv` is minimal and better aligns with a lean startup phase.

### 5. Makefile
*   **Recommendation:** A `Makefile` should be included in Phase 1 as a core developer experience improvement. It should be at the project root.
*   **Required `Makefile` example (Note: Makefiles require TABS, not spaces):**
    ```makefile
    .PHONY: up down logs

    # Starts all services defined in docker-compose.
    up:
    	docker-compose up --build -d

    # Stops all services.
    down:
    	docker-compose down

    # Tails logs for a specific service (e.g., make logs s=backend).
    logs:
    	docker-compose logs -f $(s)
    ```

---

## Recommended Scope for a Leaner Phase 1

To ensure a fast and focused start, the following items should constitute the **entirety of Phase 1**. All other items should be deferred.

*   **MUST-HAVE Deliverables for Phase 1:**
    1.  **Project Structure:** The basic `backend/cmd/server` directory.
    2.  **Go Module:** Initialized with `go mod init chess-coach/backend`.
    3.  **`.gitignore` file:** A basic gitignore for Go projects.
    4.  **Root `docker-compose.yml`:** Containing only the `backend` service for now, with the correct volume mount for hot-reloading.
    5.  **Minimal `main.go`:** A simple `net/http` server with a single `/health` endpoint. **Do not add Fiber or any other framework yet.**
    6.  **Production `Dockerfile`:** A multi-stage Dockerfile in the `backend/` directory.
    7.  **Hot-Reloading:** A working `air` setup via the Docker Compose command.
    8.  **Root `Makefile`:** With `up`, `down`, and `logs` commands.

## Defer to Phase 2: Framework & Features

The following items were correctly identified as necessary but should be moved out of the initial scaffolding phase to avoid scope creep:

*   **Web Framework:** Integration of **Fiber v3**.
*   **Structured Logging:** Setup of **`slog`**.
*   **CORS Middleware:** Adding and configuring CORS.
*   **Robust WebSocket Handler:** The simple echo handler is sufficient for Phase 1; the full Hub/Client pattern with heartbeats should be part of the feature implementation phase.
*   **Additional Docker Services:** The `postgres`, `redis`, and `nats` services can be added to `docker-compose.yml` as soon as they are needed by a feature.
*   **Testing Strategy:** A formal testing strategy with `testify` can be established once there is application logic to test.

By adopting this leaner approach, we can validate the foundational setup quickly and then iteratively add complexity in a more controlled manner.

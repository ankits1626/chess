# Phase 1 Plan: Foundational Go Backend

## 1. Objective

This document outlines the plan for Phase 1 of the backend development. The goal is to scaffold a new, production-ready Go backend project based on the decisions in the `tech-stack-recommendation.md` document. This phase will deliver a containerized, live-reloading development environment with a basic HTTP and WebSocket server, ready for future feature development.

---

## 2. Tech Stack

The following technologies will be used for this phase:

| Layer | Technology | Why? |
| :--- | :--- | :--- |
| **Language** | Go (1.23+) | High performance, low memory footprint, and excellent concurrency. |
| **Web Framework** | Fiber (v2) | Express.js-like API for rapid development and high performance. |
| **WebSocket** | `gorilla/websocket` | A robust and widely-used WebSocket library for Go. |
| **Local Env** | Docker & Docker Compose | For a consistent, one-command local development setup. |
| **Hot-Reloading** | `air` | Enables live-reloading of the Go application on file changes. |

---

## 3. Implementation Steps

The implementation will be broken down into the following steps:

### Step 1: Project Scaffolding

1.  Create the standard Go project directory structure inside the `backend/` folder:
    *   `cmd/server/main.go` (The application's entry point)
    *   `internal/` (For all core application logic like game management and WebSocket hubs)
    *   `pkg/` (For shared, non-business-specific code)
2.  Initialize the Go module within the `backend/` directory using `go mod init`.

### Step 2: Initial HTTP Server

1.  Add the `Fiber` framework as a project dependency.
2.  Implement a basic HTTP server in `main.go`.
3.  Create a `/health` endpoint that returns a `200 OK` status. This is crucial for container health checks in production environments.

### Step 3: Containerization & Local Environment

1.  Create a multi-stage `Dockerfile` in the `backend/` directory. This will compile the Go application and result in a minimal, production-ready container image (~15MB).
2.  Create a `docker-compose.yml` file at the project root. This file will define all the services required for local development:
    *   `backend`: The Go application.
    *   `postgres`: A PostgreSQL database instance.
    *   `redis`: A Redis instance for caching.
    *   `nats`: A NATS message broker for future real-time scaling.

### Step 4: Basic WebSocket Server

1.  Add the `gorilla/websocket` library as a project dependency.
2.  Create a middleware in Fiber to handle the HTTP-to-WebSocket upgrade request.
3.  Implement a WebSocket handler at the `/ws` endpoint.
4.  For this initial phase, the handler will simply "echo" any message it receives back to the same client. This serves as a functional test to confirm the real-time communication channel is working.

### Step 5: Hot-Reloading for Development

1.  Add instructions for installing the `air` live-reloading tool.
2.  Create an `.air.toml` configuration file in the `backend/` directory. This file will define the build commands and file paths for `air` to monitor.
3.  The `docker-compose.yml` service for the backend will be configured to use `air` as its startup command, providing a seamless and fast development loop.

---

## 4. Phase 1 Deliverables

Upon completion of this phase, the following assets will be in place:

*   A structured Go project within the `backend/` directory.
*   A running HTTP server with a `/health` endpoint.
*   A functioning WebSocket echo server at `/ws`.
*   A `Dockerfile` for building a minimal, production-ready Go container.
*   A `docker-compose.yml` file for a one-command local development environment.
*   A configuration file for `air` to enable live-reloading.

# Chess Coach - Project Structure Guide

## Monorepo Layout

```
chess-coach/                      # Root (monorepo)
├── Makefile                      # Dev commands (make up, make down, make logs)
├── docker-compose.yml            # Orchestrates backend (will add frontend/postgres/redis later)
├── backend/                      # Go backend (Phase 1 complete)
└── frontend/                     # React frontend (Steps 1-13 complete)
```

---

## Backend Structure (Go - Standard Layout)

```
backend/
├── cmd/                          # Application entry points
│   └── server/
│       └── main.go              # HTTP server with /health endpoint
│
├── internal/                     # Private application code
│   └── (empty for now)          # Will contain: websocket/, game/, ai/, db/
│
├── pkg/                          # Public libraries (reusable)
│   └── (empty for now)          # Will contain: types/, utils/
│
├── docs/                         # Architecture & planning docs
│   ├── go-tech-stack-2025.md    # Tech stack decisions
│   ├── aws-architecture.md      # AWS deployment plan
│   ├── messaging-architecture-2025.md
│   └── ...
│
├── tmp/                          # Air build artifacts (gitignored)
│   └── main                     # Hot-reloaded binary
│
├── .air.toml                    # Hot-reload configuration
├── .env.example                 # Environment variables template
├── .gitignore                   # Go-specific ignores
├── Dockerfile                   # Multi-stage production build
└── go.mod                       # Go module definition
```

---

## Key Go Concepts

### 1. `cmd/` Directory

**Purpose:** Application entry points (binaries)

**Pattern:** `cmd/<binary-name>/main.go`

**Example:** `cmd/server/main.go` → builds `./server` binary

**Why separate?** You can have multiple binaries in one project:
- `cmd/server/` - Main HTTP server
- `cmd/worker/` - Background job processor (future)
- `cmd/migration/` - Database migration tool (future)

### 2. `internal/` Directory

**Purpose:** Private application code (business logic)

**Special Go rule:** Code in `internal/` **cannot be imported** by other projects outside this repo. This enforces encapsulation.

**Future structure (Phase 2+):**
```
internal/
├── server/         # Fiber v3 HTTP server setup
├── websocket/      # WebSocket hub & connection management
│   ├── hub.go
│   └── client.go
├── game/           # Chess game logic, PGN handling
│   ├── manager.go
│   └── pgn.go
├── ai/             # Claude API integration
│   └── coach.go
├── db/             # Database queries (sqlc generated)
│   └── queries/
├── middleware/     # CORS, logging, auth
└── config/         # Environment loading
```

### 3. `pkg/` Directory

**Purpose:** Public libraries that could be reused by other projects

**Can be imported:** Unlike `internal/`, code in `pkg/` can be imported by external projects

**Example use cases:**
- `pkg/types/messages.go` - Shared message types
- `pkg/validator/` - Custom validation logic
- `pkg/chess/` - Pure chess utilities (if we wanted to publish a library)

**Current state:** Empty - we'll add code here only if we need truly reusable, public libraries

### 4. `tmp/` Directory

**Created by:** Air (hot-reload tool)

**Contains:** Compiled binary during development (`tmp/main`)

**Gitignored:** Never committed to version control

**Purpose:** Air compiles your code to `tmp/main` and runs it. On file changes, it rebuilds and restarts.

---

## File-by-File Breakdown

### Configuration Files

#### `.air.toml`
- **Purpose:** Configures Air hot-reload behavior
- **Key settings:**
  - `cmd`: Build command (`go build -o ./tmp/main ./cmd/server`)
  - `include_ext`: Watch `.go` files
  - `exclude_dir`: Ignore `tmp/`, `vendor/`
  - `delay`: Wait 1s before rebuilding (debounce)

#### `.env.example`
- **Purpose:** Template for environment variables
- **Not loaded automatically** - copy to `.env` for local overrides
- **Variables defined:**
  - `PORT=8080` - Server port
  - `POSTGRES_URL` - Database connection (Phase 2+)
  - `REDIS_URL` - Cache connection (Phase 2+)
  - `NATS_URL` - Message broker (Phase 2+)

#### `.gitignore`
- **Purpose:** Prevent committing build artifacts and secrets
- **Key patterns:**
  - `tmp/` - Air builds
  - `.env`, `.env.local` - Secrets
  - `*.exe`, `*.dylib` - Compiled binaries
  - `vendor/` - Go dependencies (if using vendoring)

#### `Dockerfile`
- **Purpose:** Multi-stage production build
- **Stage 1 (builder):** Compile Go binary on `golang:1.23-alpine`
- **Stage 2 (production):** Copy binary to minimal `alpine:latest` (~15MB final image)
- **NOT used in development** - dev uses `golang:1.23-alpine` directly for hot-reload

#### `go.mod`
- **Purpose:** Go module definition (like `package.json`)
- **Contents:**
  ```go
  module chess-coach/backend
  go 1.23
  ```
- **Module name:** `chess-coach/backend` - used for imports
- **Go version:** Requires Go 1.23+

---

## How It All Works Together

### Development Flow

```
1. Run: make up
   ↓
2. docker-compose.yml starts golang:1.23-alpine container
   ↓
3. Mounts ./backend → /app inside container
   ↓
4. Container runs: go install github.com/air-verse/air@v1.52.3
   ↓
5. Starts Air with config from .air.toml
   ↓
6. Air watches: cmd/, internal/, pkg/ for .go file changes
   ↓
7. Initial build: go build -o ./tmp/main ./cmd/server
   ↓
8. Runs: ./tmp/main (starts HTTP server on :8080)
   ↓
9. Docker port mapping: container:8080 → localhost:8080
   ↓
10. Server ready! Visit http://localhost:8080/health
   ↓
11. On file save: Air detects change → Rebuilds → Restarts (< 1 second)
```

### Request Flow (Current)

```
Browser: curl http://localhost:8080/health
  ↓
Docker port 8080:8080
  ↓
Container golang:1.23-alpine
  ↓
./tmp/main (running server)
  ↓
cmd/server/main.go: http.HandleFunc("/health", ...)
  ↓
Response: {"status":"ok","service":"chess-coach-backend"}
```

---

## Current Endpoints

| Method | Path | Response | Purpose |
|--------|------|----------|---------|
| `GET` | `/health` | `{"status":"ok","service":"chess-coach-backend"}` | Health checks (AWS ECS, load balancers) |
| `GET` | `/` | `{"message":"Chess Coach Backend API","version":"0.1.0"}` | API info |

---

## Why This Structure?

### 1. **Standard Go Convention**
- Any Go developer recognizes `cmd/`, `internal/`, `pkg/`
- Follows [golang-standards/project-layout](https://github.com/golang-standards/project-layout)

### 2. **Separation of Concerns**
- Entry points (`cmd/`) separate from business logic (`internal/`)
- Main function stays small, delegates to internal packages
- Easy to add new binaries (workers, CLIs) without touching core logic

### 3. **Monorepo Friendly**
- Backend self-contained in `backend/` folder
- Frontend in `frontend/` folder
- Can share root `docker-compose.yml` for full-stack dev
- Single repo = atomic commits across frontend + backend

### 4. **Hot-Reload Ready**
- Air watches source directories automatically
- Volume mounts enable instant file sync
- Rebuild + restart in < 1 second
- Mirrors production debugging workflow

### 5. **Production-Ready**
- Dockerfile creates minimal production image (15MB)
- Multi-stage build = small attack surface
- Can deploy same structure to AWS ECS, Kubernetes, etc.

---

## Comparison to Other Languages

| Go | Node.js/TypeScript | Python |
|----|-------------------|--------|
| `cmd/server/main.go` | `src/index.ts` | `main.py` |
| `internal/` | `src/lib/` or `src/services/` | `app/` |
| `pkg/` | `src/shared/` or published npm package | `lib/` or PyPI package |
| `go.mod` | `package.json` | `requirements.txt` / `pyproject.toml` |
| Air | nodemon / tsx --watch | uvicorn --reload |

---

## What's Next? (Phase 2 Preview)

### Planned Internal Structure

```
internal/
├── server/
│   └── server.go           # Fiber v3 app initialization
├── websocket/
│   ├── hub.go             # Connection registry, broadcast
│   └── client.go          # Per-client connection handler
├── game/
│   ├── manager.go         # Game state management
│   └── pgn.go            # PGN parsing & validation
├── ai/
│   └── coach.go          # Claude API integration
├── db/
│   ├── queries/          # SQL files for sqlc
│   └── models.go         # Generated database types
├── middleware/
│   ├── cors.go           # CORS for frontend
│   ├── logger.go         # Request logging
│   └── auth.go           # JWT validation (future)
└── config/
    └── config.go         # Load .env with godotenv
```

### New Dependencies (Phase 2)

```go
// go.mod after Phase 2
require (
    github.com/gofiber/fiber/v3 v3.0.0
    github.com/gorilla/websocket v1.5.1
    github.com/joho/godotenv v1.5.1
    // ... more as needed
)
```

---

## Common Commands

### Development

```bash
# Start backend
make up

# View logs
make logs
make backend-logs  # Backend only

# Stop backend
make down

# Restart after docker-compose.yml changes
make restart

# Stop and remove volumes
make clean
```

### Manual Docker Commands (if needed)

```bash
# Build and start
docker compose up --build -d

# View logs
docker compose logs -f backend

# Stop
docker compose down

# Enter container shell
docker compose exec backend sh

# Check Go version inside container
docker compose exec backend go version
```

### Testing Endpoints

```bash
# Health check
curl http://localhost:8080/health

# Root endpoint
curl http://localhost:8080/

# Pretty-print JSON
curl -s http://localhost:8080/health | jq
```

---

## Troubleshooting

### Container won't start
```bash
# Check logs
make backend-logs

# Common issues:
# 1. Port 8080 already in use → change PORT in docker-compose.yml
# 2. Air version incompatible → check go.mod requires go 1.23
# 3. Syntax error in main.go → check logs for compilation error
```

### Hot-reload not working
```bash
# 1. Check volume mount in docker-compose.yml
volumes:
  - ./backend:/app  # Must be present

# 2. Verify Air is watching
make backend-logs
# Should see: "watching ." "watching cmd" etc.

# 3. Make a change to main.go and check logs
# Should see: "building..." "running..."
```

### Can't access localhost:8080
```bash
# 1. Check container is running
docker ps | grep chess-coach-backend

# 2. Check port mapping
docker compose ps
# Should show: 0.0.0.0:8080->8080/tcp

# 3. Check firewall/VPN
# Some corporate VPNs block localhost ports
```

---

## Resources

- [Go Project Layout Standard](https://github.com/golang-standards/project-layout)
- [Air Documentation](https://github.com/air-verse/air)
- [Docker Compose v2 Docs](https://docs.docker.com/compose/)
- [Go Module Reference](https://go.dev/ref/mod)

---

**Last Updated:** Phase 1 Complete (Nov 1, 2025)
**Current State:** Minimal HTTP server with hot-reload working
**Next Phase:** Fiber v3 + WebSocket + Database integration

# Step 2: Install & Containerize Stockfish

**Duration**: 30 minutes (reduced - production Dockerfile already configured!)

**Status**: ⬜ Not Started

**Goal**: Install Stockfish locally for development and update development Docker setup

---

## 🎉 Good News!

Your **production Dockerfile already has Stockfish installed** at [Dockerfile:35](../../../Dockerfile#L35):

```dockerfile
RUN apk --no-cache add ca-certificates postgresql-client stockfish
```

✅ **Production deployment is ready!** No changes needed.

We only need to:
1. Install Stockfish locally for development
2. Add Stockfish to **Dockerfile.dev** (development)
3. Add environment variable to docker-compose.yml

---

## 📦 Part 1: Local Installation (Development)

### For macOS:

```bash
# Install Stockfish via Homebrew
brew install stockfish

# Verify installation
stockfish

# Expected output:
# Stockfish 16 by the Stockfish developers
# Type 'quit' to exit
```

### For Linux:

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install stockfish

# Verify
stockfish
```

### Find Stockfish Path

```bash
# macOS/Linux
which stockfish

# Common paths:
# macOS (Apple Silicon): /opt/homebrew/bin/stockfish
# macOS (Intel): /usr/local/bin/stockfish
# Linux: /usr/games/stockfish or /usr/bin/stockfish
```

**Save this path!** You'll use it in development.

---

## 🐳 Part 2: Update Development Docker Setup

Your production Dockerfile is already perfect! We just need to update the **development** setup.

---

## 📁 Step 1: Update Dockerfile.dev

### File: [backend/Dockerfile.dev](../../../Dockerfile.dev)

**Current line 6**:
```dockerfile
RUN apk add --no-cache git make postgresql-client
```

**Change to** (add `stockfish`):
```dockerfile
RUN apk add --no-cache git make postgresql-client stockfish
```

That's it! Just add `stockfish` to the existing `apk add` command.

---

## 📁 Step 2: Update docker-compose.yml

### File: [backend/docker-compose.yml](../../../docker-compose.yml)

**Find the `api` service environment section** (around line 31):

```yaml
  api:
    # ... existing config ...
    environment:
      - PORT=8080
      - GIN_MODE=debug
      - ENVIRONMENT=development
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=chess_coach
      - DB_PASSWORD=chess_coach_dev
      - DB_NAME=chess_coach_dev
      - DB_SSLMODE=disable
```

**Add this line** after `DB_SSLMODE`:
```yaml
      - STOCKFISH_PATH=/usr/games/stockfish
```

**Result**:
```yaml
  api:
    # ... existing config ...
    environment:
      - PORT=8080
      - GIN_MODE=debug
      - ENVIRONMENT=development
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=chess_coach
      - DB_PASSWORD=chess_coach_dev
      - DB_NAME=chess_coach_dev
      - DB_SSLMODE=disable
      - STOCKFISH_PATH=/usr/games/stockfish  # ← ADD THIS LINE
```

---

## 📁 Step 3: Update AI Service to Use Environment Variable

### File: `internal/websocket/handlers/ai/stockfish.go` (Preview)

```go
package ai

import (
	"fmt"
	"os"

	"github.com/notnil/chess/uci"
)

// NewStockfishAI creates Stockfish AI with path from environment
func NewStockfishAI() (*StockfishAI, error) {
	// Get path from environment variable
	stockfishPath := os.Getenv("STOCKFISH_PATH")

	// Default paths for different environments
	if stockfishPath == "" {
		// Try common paths
		paths := []string{
			"/usr/games/stockfish",        // Alpine/Docker
			"/usr/bin/stockfish",          // Linux
			"/opt/homebrew/bin/stockfish", // macOS Apple Silicon
			"/usr/local/bin/stockfish",    // macOS Intel
		}

		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				stockfishPath = path
				break
			}
		}

		if stockfishPath == "" {
			return nil, fmt.Errorf("stockfish not found, set STOCKFISH_PATH environment variable")
		}
	}

	// Create engine
	eng, err := uci.New(stockfishPath)
	if err != nil {
		return nil, fmt.Errorf("failed to start stockfish at %s: %w", stockfishPath, err)
	}

	// Initialize
	if err := eng.Run(uci.CmdUCI, uci.CmdIsReady, uci.CmdUCINewGame); err != nil {
		return nil, fmt.Errorf("failed to initialize stockfish: %w", err)
	}

	return &StockfishAI{
		engine: eng,
		path:   stockfishPath,
	}, nil
}
```

---

## 🧪 Part 3: Testing

### Test Local Installation

```bash
# Test Stockfish directly
stockfish

# Type these commands:
uci
isready
position startpos
go depth 10
quit

# Expected: You should see "readyok" and a bestmove
```

### Test Docker Build

```bash
cd /Users/ankit/code/learn/chess-coach/backend

# Build Docker image (development)
docker-compose build api

# Check if Stockfish is in the image
docker run --rm chess-coach-api which stockfish
# Should output: /usr/games/stockfish

# Test Stockfish in container
docker run --rm -it chess-coach-api stockfish
# Type 'uci' then 'quit' to test
```

### Test Full Development Stack

```bash
# Start everything
docker-compose up -d

# Check api service logs
docker-compose logs -f api

# Stop when done testing
docker-compose down
```

### Test Production Build

```bash
# Build production image
docker build -t chess-coach-backend -f Dockerfile .

# Verify Stockfish in production image
docker run --rm chess-coach-backend sh -c "which stockfish"
# Should output: /usr/games/stockfish

# Test Stockfish
docker run --rm -it chess-coach-backend stockfish
```

---

## 📝 Environment Variables

### Development (.env.local)
```bash
STOCKFISH_PATH=/opt/homebrew/bin/stockfish  # Your local path
DATABASE_URL=postgresql://localhost:5432/chess_coach
```

### Production (.env.production)
```bash
STOCKFISH_PATH=/usr/games/stockfish  # Docker path
DATABASE_URL=postgresql://prod-db:5432/chess_coach
```

### In Code (config/config.go)

```go
type Config struct {
	Database      DatabaseConfig
	Server        ServerConfig
	StockfishPath string `env:"STOCKFISH_PATH" envDefault:"/usr/games/stockfish"`
}
```

---

## 🚀 Deployment Options

### Option 1: Docker Compose (Simple)
```bash
# On your server
git clone <repo>
cd backend
docker-compose up -d
```

### Option 2: Kubernetes (Scalable)
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: chess-backend
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: backend
        image: your-registry/chess-backend:latest
        env:
        - name: STOCKFISH_PATH
          value: /usr/games/stockfish
```

### Option 3: Cloud Run / ECS / App Engine
All support Docker containers - just push your image!

---

## ✅ Verification Checklist

### Local Development
- [ ] Stockfish installed via Homebrew/apt
- [ ] `stockfish` command works in terminal
- [ ] Found Stockfish path (e.g., `/opt/homebrew/bin/stockfish`)
- [ ] Can test UCI protocol manually

### Docker
- [ ] Dockerfile created with Stockfish
- [ ] docker-compose.yml configured
- [ ] `docker-compose build` succeeds
- [ ] Stockfish found in container (`docker run ... which stockfish`)
- [ ] Can start Stockfish in container
- [ ] Environment variable `STOCKFISH_PATH` set

### Code
- [ ] AI service uses `os.Getenv("STOCKFISH_PATH")`
- [ ] Fallback paths for different environments
- [ ] Error message if Stockfish not found

---

## 🔧 Troubleshooting

### Issue: "stockfish not found" in Docker

**Solution**: Check Alpine package name
```bash
docker run --rm alpine:latest apk search stockfish
# Should show: stockfish-16.1-r0 (or similar)
```

### Issue: Stockfish path different in Docker

**Solution**: Find path in container
```bash
docker run --rm alpine:latest sh -c "apk add stockfish && which stockfish"
# Output: /usr/games/stockfish
```

### Issue: Permission denied

**Solution**: Make sure binary has execute permissions
```dockerfile
RUN chmod +x /usr/games/stockfish
```

---

## 📊 File Sizes

**Without optimization**:
- Image size: ~800MB (Go binary + Alpine + Stockfish)

**With multi-stage build** (our approach):
- Builder stage: ~1GB (doesn't matter, discarded)
- Final image: ~50MB ✅
  - Alpine base: ~7MB
  - Stockfish: ~20MB
  - Go binary: ~15MB
  - Dependencies: ~8MB

---

## 🎯 Benefits of This Approach

| Aspect | Local Only | Docker |
|--------|-----------|---------|
| **Dev Speed** | Fast ⚡ | Slightly slower |
| **Consistency** | ❌ Different on each machine | ✅ Same everywhere |
| **Deployment** | ❌ Manual setup needed | ✅ One command |
| **CI/CD** | ❌ Hard to test | ✅ Easy to test |
| **Scaling** | ❌ Manual | ✅ Automatic |
| **Version Lock** | ❌ Different versions | ✅ Locked version |

---

## 📁 Files Created/Modified

```
backend/
├── Dockerfile                 # NEW - Multi-stage build
├── docker-compose.yml         # NEW - Full stack setup
├── .dockerignore             # NEW - Ignore unnecessary files
├── .env.example              # NEW - Environment template
└── internal/
    └── websocket/handlers/ai/
        └── stockfish.go      # MODIFIED - Use env var
```

---

## 🎯 Next Step

Once Stockfish is installed and containerized:

→ **[Step 3: Implement Player Types](./03-implement-player-types.md)**

Now we'll create HumanPlayer and ComputerPlayer implementations!

---

## 💡 Pro Tips

1. **Development**: Use local Stockfish (faster iteration)
2. **CI/CD**: Use Docker (consistent tests)
3. **Production**: Use Docker (reliable deployment)
4. **Monitoring**: Add health check that verifies Stockfish is running

```go
// Health check endpoint
func (s *Server) healthCheck(c *gin.Context) {
    // Test Stockfish
    if err := s.aiService.HealthCheck(); err != nil {
        c.JSON(500, gin.H{"status": "unhealthy", "stockfish": "down"})
        return
    }
    c.JSON(200, gin.H{"status": "healthy", "stockfish": "up"})
}
```

---

**Estimated Time**: 45 minutes (15 min local + 30 min Docker)
**Actual Time**: ________

---

**Last Updated**: 2025-11-05
**Docker Image Size**: ~50MB

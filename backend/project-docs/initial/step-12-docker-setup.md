# Step 12: Docker Setup

## Goal
Containerize the Chess Coach backend with Docker and Docker Compose, including PostgreSQL database.

---

## Why Docker First?

Before adding database code:
- ✅ Consistent development environment
- ✅ PostgreSQL in container (no local install)
- ✅ Easy onboarding (`docker-compose up`)
- ✅ Production-ready from day 1
- ✅ Isolated dependencies

---

## Prerequisites

- ✅ SOLID-compliant backend structure
- ✅ Gin framework with Swagger
- Docker Desktop installed (or Docker Engine)

---

## Architecture

```
chess-coach/
├── backend/
│   ├── Dockerfile              # Go app container
│   ├── docker-compose.yml      # Multi-container setup
│   ├── .dockerignore           # Exclude files
│   └── cmd/server/main.go
└── database/
    └── init/
        └── 01-init.sql         # Database initialization
```

---

## Step-by-Step Implementation

### Step 1: Create Dockerfile
**What:** Multi-stage build for Go app
**Time:** 5 minutes

**File:** `backend/Dockerfile`

```dockerfile
# Build stage
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Copy Swagger docs
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./main"]
```

**Benefits:**
- Multi-stage build (small final image)
- Alpine Linux (minimal, secure)
- Only binary in final image
- ~20MB final image size

---

### Step 2: Create .dockerignore
**What:** Exclude unnecessary files
**Time:** 2 minutes

**File:** `backend/.dockerignore`

```
# Git
.git
.gitignore

# Documentation
project-docs/
*.md

# Tests
*_test.go

# Temporary files
tmp/
*.tmp
*.log

# IDE
.vscode/
.idea/

# Dependencies (will download in container)
vendor/

# OS
.DS_Store
```

---

### Step 3: Create docker-compose.yml
**What:** Multi-container orchestration
**Time:** 10 minutes

**File:** `backend/docker-compose.yml`

```yaml
version: '3.8'

services:
  # PostgreSQL Database
  postgres:
    image: postgres:16-alpine
    container_name: chess-coach-db
    restart: unless-stopped
    environment:
      POSTGRES_USER: chess_coach
      POSTGRES_PASSWORD: chess_coach_dev
      POSTGRES_DB: chess_coach_dev
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ../database/init:/docker-entrypoint-initdb.d
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U chess_coach"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Go Backend API
  api:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: chess-coach-api
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - GIN_MODE=debug
      - ENVIRONMENT=development
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_USER=chess_coach
      - DB_PASSWORD=chess_coach_dev
      - DB_NAME=chess_coach_dev
    depends_on:
      postgres:
        condition: service_healthy
    volumes:
      # Hot reload in development
      - .:/app
    command: air

  # pgAdmin (Optional - Database UI)
  pgadmin:
    image: dpage/pgadmin4:latest
    container_name: chess-coach-pgadmin
    restart: unless-stopped
    environment:
      PGADMIN_DEFAULT_EMAIL: admin@chesscoach.com
      PGADMIN_DEFAULT_PASSWORD: admin
    ports:
      - "5050:80"
    depends_on:
      - postgres

volumes:
  postgres_data:
```

**Services:**
1. **postgres** - Database (PostgreSQL 16)
2. **api** - Go backend
3. **pgadmin** - Database UI (optional)

---

### Step 4: Development docker-compose
**What:** Separate config for hot reload
**Time:** 5 minutes

**File:** `backend/docker-compose.dev.yml`

```yaml
version: '3.8'

services:
  api:
    build:
      context: .
      dockerfile: Dockerfile.dev
    volumes:
      - .:/app
      - /app/tmp  # Exclude tmp dir
    environment:
      - GIN_MODE=debug
    command: air
```

**File:** `backend/Dockerfile.dev`

```dockerfile
FROM golang:1.25.3-alpine

WORKDIR /app

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# Install swag for Swagger
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Generate Swagger docs
RUN swag init -g cmd/server/main.go

EXPOSE 8080

CMD ["air"]
```

---

### Step 5: Database Initialization
**What:** Initial database schema
**Time:** 5 minutes

**File:** `database/init/01-init.sql`

```sql
-- Initial database setup
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table (placeholder for now)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Health check (verify database is working)
CREATE TABLE IF NOT EXISTS health_check (
    id SERIAL PRIMARY KEY,
    checked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO health_check (checked_at) VALUES (CURRENT_TIMESTAMP);
```

---

### Step 6: Update Config Package
**What:** Read database connection from env
**Time:** 5 minutes

**File:** `internal/config/config.go` (update)

```go
// Config holds app settings.
type Config struct {
	// Server settings
	Port        string
	Environment string
	GinMode     string

	// Database settings
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

// Load reads config from env with defaults.
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		GinMode:     getEnv("GIN_MODE", "debug"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "chess_coach"),
		DBPassword: getEnv("DB_PASSWORD", "chess_coach_dev"),
		DBName:     getEnv("DB_NAME", "chess_coach_dev"),
	}
}
```

---

### Step 7: Update .gitignore
**What:** Ignore Docker volumes
**Time:** 1 minute

Add to `.gitignore`:
```
# Docker
*.env
.env.local
docker-compose.override.yml
```

---

### Step 8: Create Makefile (Optional)
**What:** Easy commands
**Time:** 5 minutes

**File:** `backend/Makefile`

```makefile
.PHONY: help build up down logs restart clean test

help:
	@echo "Chess Coach Backend - Docker Commands"
	@echo ""
	@echo "make build    - Build Docker images"
	@echo "make up       - Start all containers"
	@echo "make down     - Stop all containers"
	@echo "make logs     - View logs"
	@echo "make restart  - Restart containers"
	@echo "make clean    - Remove containers and volumes"
	@echo "make test     - Run tests"

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

restart:
	docker-compose restart

clean:
	docker-compose down -v
	docker system prune -f

test:
	docker-compose exec api go test -v ./...

dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
```

---

## Usage

### Start Development Environment
```bash
cd backend

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f api

# Or using Makefile
make up
make logs
```

### Access Services
- **API:** http://localhost:8080
- **Swagger:** http://localhost:8080/swagger/index.html
- **pgAdmin:** http://localhost:5050 (admin@chesscoach.com / admin)
- **PostgreSQL:** localhost:5432

### Stop Services
```bash
docker-compose down

# Or with Makefile
make down
```

### Rebuild After Code Changes
```bash
docker-compose up -d --build

# Or
make build
make up
```

---

## Verification

### Step 1: Build Images
```bash
cd backend
docker-compose build
```

**Expected:** No errors, images built successfully

### Step 2: Start Services
```bash
docker-compose up -d
```

**Expected:** 3 containers running
- chess-coach-db
- chess-coach-api
- chess-coach-pgadmin

### Step 3: Check Container Status
```bash
docker-compose ps
```

**Expected:**
```
NAME                 STATUS
chess-coach-db       Up (healthy)
chess-coach-api      Up
chess-coach-pgadmin  Up
```

### Step 4: Test API
```bash
curl http://localhost:8080/api/v1/health
```

**Expected:**
```json
{"status":"healthy","version":"1.0"}
```

### Step 5: Check Swagger
Open: http://localhost:8080/swagger/index.html

**Expected:** Swagger UI loads

### Step 6: Check Database
```bash
docker-compose exec postgres psql -U chess_coach -d chess_coach_dev -c "SELECT * FROM health_check;"
```

**Expected:** Returns 1 row

---

## Troubleshooting

### Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080

# Kill it or change port in docker-compose.yml
```

### Database Connection Failed
```bash
# Check postgres is healthy
docker-compose ps

# View postgres logs
docker-compose logs postgres

# Restart postgres
docker-compose restart postgres
```

### Air Not Working
```bash
# Check air is installed in container
docker-compose exec api which air

# Rebuild with dev dockerfile
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```

---

## Production Configuration

### Production docker-compose.prod.yml

```yaml
version: '3.8'

services:
  api:
    build:
      context: .
      dockerfile: Dockerfile
    environment:
      - GIN_MODE=release
      - ENVIRONMENT=production
    restart: always
    # No volume mounts
    # No Air (use compiled binary)
```

### Production Build
```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

---

## Next Steps

After Docker setup:
1. **Add database package** (Step 13)
2. **Add GORM models** (Step 13)
3. **Add migrations** (Step 13)
4. **Test database connectivity** (Step 13)

---

## Benefits Achieved

✅ **Consistent environment** across all machines
✅ **PostgreSQL containerized** (no local install)
✅ **Hot reload** in development
✅ **Production-ready** setup
✅ **Easy onboarding** (`docker-compose up`)
✅ **Database UI** (pgAdmin) for debugging
✅ **Isolated** from system dependencies

---

## File Checklist

- [ ] `backend/Dockerfile`
- [ ] `backend/Dockerfile.dev`
- [ ] `backend/.dockerignore`
- [ ] `backend/docker-compose.yml`
- [ ] `backend/docker-compose.dev.yml`
- [ ] `backend/docker-compose.prod.yml`
- [ ] `database/init/01-init.sql`
- [ ] `backend/Makefile` (optional)
- [ ] Updated `internal/config/config.go`
- [ ] Updated `.gitignore`

---

## Time Estimate

| Task | Time |
|------|------|
| Dockerfile | 5 min |
| .dockerignore | 2 min |
| docker-compose.yml | 10 min |
| Dev configs | 5 min |
| Database init | 5 min |
| Config update | 5 min |
| Makefile | 5 min |
| Testing | 10 min |
| **Total** | **45-60 min** |

---

**Ready to start? Say "start docker" and I'll create all files!**

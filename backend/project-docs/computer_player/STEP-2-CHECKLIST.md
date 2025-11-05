# Step 2: Install Stockfish - Quick Checklist

**Goal**: Install Stockfish locally for development and update Docker setup

**Estimated Time**: 30 minutes (production Dockerfile already has Stockfish!)

---

## 🎉 Update: Docker Solution Found!

**Issue Discovered**: Stockfish is NOT available in Alpine 3.22 stable repositories.

**Solution**: We build Stockfish 17.1 from source in the Docker image.

✅ **Dockerfile.dev already updated** - Builds Stockfish from source
✅ **docker-compose.yml already updated** - Uses `/usr/local/bin/stockfish`

You only need to:
1. Install Stockfish locally for development (optional but recommended)

---

## 📋 Part 1: Local Installation (15 minutes)

### For macOS:

```bash
# 1. Install Stockfish via Homebrew
brew install stockfish

# 2. Verify installation
stockfish
# Type 'quit' to exit

# 3. Find Stockfish path (SAVE THIS!)
which stockfish
# Expected: /opt/homebrew/bin/stockfish (Apple Silicon)
#       or: /usr/local/bin/stockfish (Intel)
```

### For Linux:

```bash
# 1. Install Stockfish
sudo apt-get update
sudo apt-get install stockfish

# 2. Verify installation
stockfish
# Type 'quit' to exit

# 3. Find Stockfish path (SAVE THIS!)
which stockfish
# Expected: /usr/games/stockfish or /usr/bin/stockfish
```

### ✅ Verification:

Test the UCI protocol manually:

```bash
stockfish

# Type these commands (one by one):
uci
isready
position startpos
go depth 10
quit

# You should see:
# - "readyok" after isready
# - Chess calculations after "go depth 10"
# - "bestmove e2e4" (or similar) at the end
```

**Save your Stockfish path**: ________________________________

---

## 📋 Part 2: Install Go UCI Package (5 minutes)

```bash
cd /Users/ankit/code/learn/chess-coach/backend

# Install the UCI wrapper (part of notnil/chess)
go get github.com/notnil/chess/uci

# Verify it's in go.mod
grep "github.com/notnil/chess" go.mod
```

---

## 📋 Part 3: Docker Setup - ✅ ALREADY COMPLETE!

### ✅ Dockerfile.dev - Already Updated

The Dockerfile now **builds Stockfish from source**:

```dockerfile
# Install development tools and Stockfish build dependencies
RUN apk add --no-cache curl g++ git make postgresql-client

# Build and install Stockfish from source
RUN cd /tmp && \
    git clone --depth 1 --branch sf_17.1 https://github.com/official-stockfish/Stockfish.git && \
    cd Stockfish/src && \
    make -j$(nproc) build ARCH=$(uname -m | sed 's/x86_64/x86-64/;s/aarch64/armv8/') && \
    cp stockfish /usr/local/bin/stockfish && \
    chmod +x /usr/local/bin/stockfish && \
    cd / && \
    rm -rf /tmp/Stockfish
```

**Why**: Stockfish is not in Alpine 3.22 repositories, so we build from source.

---

### ✅ docker-compose.yml - Already Updated

The environment variable is set to the correct path:

```yaml
  api:
    environment:
      # ... other vars ...
      - STOCKFISH_PATH=/usr/local/bin/stockfish  # ← Already added
```

**Note**: Path is `/usr/local/bin/stockfish` (where we build it), not `/usr/games/stockfish` (Alpine package location).

---

### Step 1: Verify Docker Build 🐳

```bash
cd /Users/ankit/code/learn/chess-coach/backend

# Build should already be complete, but verify
docker compose build api

# This takes ~2 minutes first time (compiles Stockfish)
# Subsequent builds are cached
```

---

### Step 2: Verify Stockfish in Container ✅

```bash
# Test that Stockfish is accessible
docker run --rm backend-api which stockfish
# Expected output: /usr/local/bin/stockfish

# Test Stockfish version
docker run --rm backend-api stockfish
# Expected output: Stockfish 17.1 by the Stockfish developers

# Test UCI protocol (interactive)
docker run --rm -it backend-api stockfish
# Type 'uci' - should see 'uciok'
# Type 'quit' to exit
```

---

### Step 3: Test Full Stack (Optional) 🚀

```bash
# Start everything
docker compose up -d

# Check api logs
docker compose logs -f api

# Stop when done
docker compose down
```

---

## 📋 Part 4: Environment Variable Setup (5 minutes)

### Create .env file (if you don't have one)

**File**: `/Users/ankit/code/learn/chess-coach/backend/.env`

```bash
# Development
STOCKFISH_PATH=/opt/homebrew/bin/stockfish  # Your local path from Part 1
DATABASE_URL=postgresql://postgres:password@localhost:5432/chess_coach
```

### Add to .gitignore

```bash
echo ".env" >> .gitignore
```

### Create .env.example (for other developers)

**File**: `/Users/ankit/code/learn/chess-coach/backend/.env.example`

```bash
# Stockfish Configuration
STOCKFISH_PATH=/opt/homebrew/bin/stockfish  # macOS Apple Silicon
# STOCKFISH_PATH=/usr/local/bin/stockfish   # macOS Intel
# STOCKFISH_PATH=/usr/games/stockfish       # Linux/Docker
# STOCKFISH_PATH=/usr/bin/stockfish         # Some Linux distributions

# Database
DATABASE_URL=postgresql://postgres:password@localhost:5432/chess_coach
```

---

## ✅ Final Verification Checklist

### Local Installation
- [ ] Stockfish installed locally (`brew install stockfish` or `apt-get install stockfish`)
- [ ] `stockfish` command works in terminal
- [ ] Found and saved Stockfish path (`which stockfish`)
- [ ] Tested UCI protocol manually (uci, isready, go depth 10)

### Go Dependencies
- [ ] Installed Go UCI package (`go get github.com/notnil/chess/uci`)
- [ ] Verified in go.mod (`grep "github.com/notnil/chess" go.mod`)

### Docker Configuration
- [x] **Dockerfile.dev** updated: Builds Stockfish 17.1 from source
- [x] **docker-compose.yml** updated: `STOCKFISH_PATH=/usr/local/bin/stockfish`
- [ ] Docker image builds successfully (`docker compose build api`)
- [ ] Stockfish found in Docker image (`docker run --rm backend-api which stockfish`)
- [ ] Stockfish version verified (`docker run --rm backend-api stockfish`)

### Environment Variables
- [ ] Created .env file with STOCKFISH_PATH (optional for local dev)
- [ ] Created .env.example for team (optional)
- [ ] Added .env to .gitignore (if created)

---

## 🚨 Common Issues & Solutions

### Issue: "stockfish: command not found" (Local)
**Solution**: Install Stockfish using the commands in Part 1

### Issue: "brew: command not found" (macOS)
**Solution**: Install Homebrew first:
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

### Issue: Docker build takes long time
**Solution**: Normal! Compiling Stockfish takes ~20 seconds. It's cached after first build.

### Issue: Build fails with "unable to select packages: stockfish"
**Solution**: This is why we build from source! Alpine 3.22 doesn't have Stockfish in repos.

### Issue: Wrong Stockfish path in container
**Solution**: We build to `/usr/local/bin/stockfish`, not `/usr/games/stockfish`

---

## 📊 What You'll Have After This Step

```
backend/
├── Dockerfile              # ⚠️  Needs same update as Dockerfile.dev
├── Dockerfile.dev          # ✅ Builds Stockfish 17.1 from source
├── docker-compose.yml      # ✅ STOCKFISH_PATH=/usr/local/bin/stockfish
├── .env                    # ✅ Local config (optional)
├── .env.example            # ✅ Template for team (optional)
└── .gitignore             # ✅ Excludes .env

Local Machine:
✅ Stockfish installed at /opt/homebrew/bin/stockfish (or similar)
✅ UCI protocol tested and working

Docker Development (Dockerfile.dev):
✅ Stockfish 17.1 compiled from source
✅ Installed at /usr/local/bin/stockfish
✅ Optimized for ARM64 (armv8) or x86_64
✅ Container builds successfully (~2 min first time)
✅ STOCKFISH_PATH environment variable configured

Docker Production (Dockerfile):
⚠️  Needs same build-from-source approach
📝 See DOCKER-SOLUTION.md for production Dockerfile updates
```

---

## 🎯 Next Step

Once all checkboxes are checked:

→ **[Step 3: Implement Player Types](./03-implement-player-types.md)**

You'll create:
- `HumanPlayer` (wraps WebSocket client)
- `ComputerPlayer` (wraps Stockfish AI)

And you'll apply the PlayerConfig refactor we discussed!

---

## 💡 Learning Notes

### What You Learned:
- **UCI Protocol**: Universal Chess Interface for chess engines
- **Docker Multi-stage Builds**: Smaller production images
- **Environment Variables**: Configuration management across environments
- **Local vs Production**: Different paths for development and deployment

### Key Takeaway:
**Docker ensures consistency** - "works on my machine" = "works in production"

---

**Questions?** Check the full guide: [02-install-stockfish.md](./02-install-stockfish.md)

**Ready?** Follow the checklist above step-by-step!

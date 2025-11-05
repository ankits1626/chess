# Chess Coach - Next Steps Roadmap

## Current State ✅

- ✅ SOLID-compliant Go backend
- ✅ Gin framework with API versioning
- ✅ Swagger documentation
- ✅ Hot reload with Air
- ✅ Test coverage (80%+)
- ✅ Graceful shutdown

---

## The Big Picture: What We're Building

**Chess Coach Application** where:
1. Players can play against an AI coach (live gameplay)
2. Coach analyzes games and provides feedback
3. Players can review past games
4. Coach gives personalized training recommendations

---

## Architecture Decision: Real-Time Requirements

Since you need **live play between player and coach**, we need:
- **WebSocket** for real-time move communication
- **Chess engine** (Stockfish) for AI coach moves
- **Game state management** for active games
- **Database** to store games and analysis

---

## Recommended Path Forward

### Phase 1: Foundation (Current - Next 2 weeks)
**Goal:** Core chess functionality without real-time yet

#### Step 1: Database Setup
**What:** Choose and setup database
**Options:**
- **PostgreSQL** (Recommended) - Relational, ACID compliant, JSON support
- **MongoDB** - Document-based, flexible schema

**Why PostgreSQL:**
- Structured data (games, users, moves)
- ACID transactions (important for game integrity)
- Great Go support (pgx, gorm)
- JSON columns for flexible data (PGN, analysis)

**Time:** 2-3 hours

#### Step 2: Chess Library Integration
**What:** Add chess logic library
**Options:**
- **notnil/chess** (Recommended) - Pure Go, complete chess implementation
- **andrewbackes/chess** - Another solid option

**Features:**
- Move validation
- Legal move generation
- Check/checkmate detection
- PGN import/export
- FEN position handling

**Time:** 1-2 hours

#### Step 3: Stockfish Integration
**What:** Integrate Stockfish chess engine
**Library:** github.com/notnil/chess/uci

**Features:**
- Position analysis
- Best move calculation
- Multi-depth analysis
- Position evaluation

**Time:** 2-3 hours

#### Step 4: Core REST API Endpoints
**What:** Build non-real-time endpoints first

```
POST   /api/v1/games              Create new game
GET    /api/v1/games              List games
GET    /api/v1/games/:id          Get game details
POST   /api/v1/games/:id/moves    Add move to game
POST   /api/v1/analyze            Analyze position
GET    /api/v1/analyze/:id        Get analysis results
```

**Time:** 1 week

---

### Phase 2: Real-Time Gameplay (Weeks 3-4)
**Goal:** Live player vs coach games

#### Step 5: WebSocket Setup
**What:** Add WebSocket support
**Library:** github.com/gorilla/websocket

**Features:**
- Persistent connections
- Real-time move broadcasting
- Connection management
- Heartbeat/ping-pong

**Time:** 2-3 days

#### Step 6: Game Room Management
**What:** Manage active game sessions

**Components:**
- Game room registry
- Player connection mapping
- Move queue
- Game state synchronization

**Time:** 2-3 days

#### Step 7: Coach AI Integration
**What:** Connect Stockfish to live games

**Features:**
- Real-time move generation
- Difficulty levels (Stockfish depth)
- Move timing (humanize coach)
- Hint system

**Time:** 2-3 days

---

### Phase 3: Enhanced Features (Weeks 5-6)
**Goal:** Polish and additional features

#### Step 8: Authentication
**What:** User accounts and sessions
**Library:** golang-jwt/jwt

**Features:**
- User registration/login
- JWT tokens
- Protected routes
- User profiles

**Time:** 3-4 days

#### Step 9: Game Analysis
**What:** Post-game analysis and insights

**Features:**
- Mistake detection
- Alternative move suggestions
- Opening identification
- Endgame tablebase

**Time:** 3-4 days

---

## Detailed Technology Stack

### Backend Core
```
Language: Go 1.25.3
Framework: Gin 1.11.0
API Docs: Swagger/OpenAPI
Testing: Go testing + testify
```

### Database
```
Primary: PostgreSQL 16+
ORM: GORM v2
Migrations: golang-migrate
```

### Chess Engine
```
Engine: Stockfish 16
Go Library: notnil/chess
UCI Protocol: notnil/chess/uci
```

### Real-Time
```
WebSocket: gorilla/websocket
State Management: Go channels + sync
Message Format: JSON
```

### Frontend (Future)
```
Framework: React/Next.js (Already setup)
Chess UI: react-chessboard
State: Zustand (Already setup)
WebSocket Client: native WebSocket API
```

---

## Immediate Next Step: Database Setup

### Why Database First?

Before adding more endpoints, we need persistent storage:
- Store games and moves
- User accounts
- Analysis history
- Coach recommendations

### PostgreSQL Setup Steps

**Step 1: Install PostgreSQL**
```bash
# macOS
brew install postgresql@16
brew services start postgresql@16
```

**Step 2: Create Database**
```bash
createdb chess_coach_dev
createdb chess_coach_test
```

**Step 3: Add Go Dependencies**
```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go get -u github.com/golang-migrate/migrate/v4
```

**Step 4: Database Models**
```go
// User
type User struct {
    ID        uuid.UUID
    Username  string
    Email     string
    Rating    int
    CreatedAt time.Time
}

// Game
type Game struct {
    ID        uuid.UUID
    WhiteID   uuid.UUID  // User or Coach
    BlackID   uuid.UUID
    PGN       string
    Result    string
    CreatedAt time.Time
}

// Move
type Move struct {
    ID       uuid.UUID
    GameID   uuid.UUID
    MoveNum  int
    Move     string  // e4, Nf3, etc.
    Position string  // FEN
    TimeUsed int     // milliseconds
}
```

---

## Alternative Approach: MVP First

If you want to see results faster, we can do **MVP-first approach**:

### MVP: Play Against Coach (No Database)

**Week 1: Minimal Viable Product**
1. Add chess library (notnil/chess)
2. Add Stockfish integration
3. Create in-memory game state
4. Add WebSocket endpoint
5. Single endpoint: `/ws/play`

**Features:**
- Connect via WebSocket
- Make moves in real-time
- Coach responds immediately
- No persistence (game lost on disconnect)
- No user accounts

**Pros:**
- See working game in 1 week
- Learn WebSocket + chess integration
- Validate core concept

**Cons:**
- Can't save games
- No analysis history
- Rebuild needed for persistence

---

## My Recommendation

### Hybrid Approach: Foundation + Quick Win

**Week 1-2: Foundation**
1. Setup PostgreSQL
2. Add database models
3. Add chess library
4. Create REST endpoints for games

**Week 3: Quick Win**
5. Add Stockfish
6. Build single WebSocket endpoint
7. Test live gameplay (in-memory first)

**Week 4+: Polish**
8. Connect WebSocket to database
9. Add authentication
10. Add analysis features

**Why this works:**
- Solid foundation (database, models)
- Quick validation (working game by week 3)
- Easier to iterate (proper structure)

---

## Questions to Decide Direction

1. **Do you want to see a working game ASAP** (MVP-first)?
   - Or build foundation properly first? (Recommended)

2. **Database preference?**
   - PostgreSQL (Recommended for chess)
   - MongoDB (More flexible, but less structured)

3. **Authentication priority?**
   - Need it from day 1?
   - Or add later after core gameplay works?

4. **Deployment target?**
   - Local development only for now?
   - Or plan for production (Docker, cloud)?

---

## What I Recommend Starting With

**Next Immediate Task: Database + Chess Library**

This gives you:
1. Storage for games
2. Chess move validation
3. PGN import/export
4. Foundation for everything else

**After that:**
5. Stockfish integration (coach brain)
6. WebSocket (real-time play)
7. Frontend connection

---

## Estimated Timeline

| Phase | Duration | Features |
|-------|----------|----------|
| Database Setup | 1-2 days | PostgreSQL, models, migrations |
| Chess Library | 1 day | Move validation, PGN |
| Stockfish | 2-3 days | Engine integration, analysis |
| REST Endpoints | 3-4 days | CRUD for games, moves |
| WebSocket | 2-3 days | Real-time communication |
| Game Rooms | 2-3 days | Session management |
| Coach AI | 2-3 days | Live gameplay logic |
| Authentication | 3-4 days | Users, JWT |
| **Total MVP** | **3-4 weeks** | Playable coach game |

---

## Ready to Start?

**I recommend:** Start with Database + Chess Library integration

**Say "start database" and I'll:**
1. Create database package
2. Setup PostgreSQL connection
3. Define models for User, Game, Move
4. Add migrations
5. Create repository layer (SOLID-compliant)

**Or tell me:**
- Your preferred approach (Foundation vs MVP-first)
- Any questions about the roadmap
- What timeline you're targeting

---

**Current Status:** ✅ Backend structure complete, ready for features!

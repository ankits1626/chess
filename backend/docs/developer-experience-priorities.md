# Developer Experience & Small Team Priorities

**Team Size:** 2-3 developers
**Context:** Cloud-native but locally-first development
**Philosophy:** Production-like locally, zero-friction deployment

---

## Core Principles

### 1. **Local Development Must Be Seamless**

Every service, every dependency, must run locally with a single command. No exceptions.

### 2. **Cloud-Native, Locally-First**

Use cloud services in production, but **everything must work offline** during development.

### 3. **Small Team Velocity**

With 2-3 developers, we optimize for:
- Time to first contribution
- Debugging speed
- Deploy confidence
- Operational simplicity

---

## Technology Re-Evaluation for Small Teams

### ❌ What We're Removing from Original Recommendation

#### 1. **Cloudflare Durable Objects** → Too Much Abstraction
**Problem:**
- Can't run locally (requires Cloudflare Workers)
- Debugging requires deploying to Cloudflare
- State inspection is difficult
- Small team can't afford vendor lock-in

**Better Alternative:**
- Use standard WebSocket server (Node.js/Bun)
- State in Redis (runs locally)
- Deploy to any cloud (Fly.io, Railway, Render)

#### 2. **Separate AI/Engine/MCP Services** → Over-Engineering for MVP
**Problem:**
- 3-4 separate services = 3-4 repos/deploys
- Microservices add complexity without scale benefits
- Debugging cross-service issues is painful

**Better Alternative:**
- Single backend monolith (initially)
- Clear module boundaries
- Split into services only when necessary (>50K users)

#### 3. **Kubernetes** → Operational Nightmare
**Problem:**
- 2-3 people can't maintain K8s
- Local K8s (minikube/kind) is slow
- YAML hell for small teams

**Better Alternative:**
- Docker Compose locally
- Fly.io or Railway for production (managed platform)
- Add K8s only if you have dedicated DevOps

---

## ✅ Revised Stack: Small Team Edition

### Architecture Philosophy

```
┌─────────────────────────────────────────┐
│  Simple, but production-ready           │
│  Works identically locally & in cloud   │
│  2-3 developers can maintain            │
└─────────────────────────────────────────┘
```

### The Stack

| Component | Technology | Why for Small Teams |
|-----------|-----------|---------------------|
| **Runtime** | Bun (not Node.js) | 3x faster, all-in-one (bundler, test, package manager) |
| **Backend** | TypeScript monolith | Single codebase, shared types with frontend |
| **Message Broker** | NATS (Docker) | Single binary, runs locally, zero config |
| **Database** | PostgreSQL + Drizzle ORM | SQL is simple, Drizzle is type-safe |
| **Cache/Realtime** | Redis | Runs locally, pub/sub for WebSocket |
| **AI Integration** | Direct SDK calls | No service abstraction needed yet |
| **Observability** | Pino (logs) + Prometheus | Simple, no external deps |
| **Deployment** | Fly.io | One command deploy, scales automatically |

---

## Local Development Setup

### The Golden Rule: `bun install && bun dev`

That's it. Everything should work.

### Docker Compose (Infrastructure Only)

```yaml
# docker-compose.yml
version: '3.8'

services:
  # NATS for messaging
  nats:
    image: nats:2.10-alpine
    ports:
      - "4222:4222"
      - "8222:8222"  # Monitoring UI
    command: --jetstream --http_port 8222
    healthcheck:
      test: ["CMD", "wget", "-q", "-O-", "http://localhost:8222/healthz"]
      interval: 5s
      timeout: 3s
      retries: 3

  # PostgreSQL for persistence
  postgres:
    image: postgres:16-alpine
    ports:
      - "5432:5432"
    environment:
      POSTGRES_DB: chess_coach
      POSTGRES_USER: dev
      POSTGRES_PASSWORD: dev
    volumes:
      - postgres-data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U dev"]
      interval: 5s
      timeout: 3s
      retries: 3

  # Redis for cache + WebSocket coordination
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 3

  # Optional: Prometheus for metrics
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus-data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'

volumes:
  postgres-data:
  prometheus-data:
```

### Backend Structure (Monolith)

```
backend/
├── src/
│   ├── index.ts                 # Server entry (Bun.serve)
│   ├── websocket.ts             # WebSocket handler
│   ├── nats.ts                  # NATS client
│   ├── db/
│   │   ├── schema.ts            # Drizzle schema
│   │   └── client.ts            # DB connection
│   ├── modules/
│   │   ├── chess/               # Chess game logic
│   │   │   ├── game.ts
│   │   │   ├── moves.ts
│   │   │   └── validation.ts
│   │   ├── ai/                  # AI coach
│   │   │   ├── claude.ts
│   │   │   └── streaming.ts
│   │   ├── engine/              # Stockfish
│   │   │   ├── uci.ts
│   │   │   └── pool.ts
│   │   └── multiplayer/         # Game rooms
│   │       ├── rooms.ts
│   │       └── coordination.ts
│   ├── types/
│   │   └── messages.ts          # Shared types
│   └── utils/
│       ├── logger.ts
│       └── metrics.ts
├── tests/
├── docker-compose.yml
├── Dockerfile
├── fly.toml                     # Fly.io config
└── package.json
```

### Single Command Start

```json
// package.json
{
  "scripts": {
    "dev": "bun run --watch src/index.ts",
    "docker": "docker-compose up -d",
    "setup": "bun run docker && bun run db:migrate",
    "db:migrate": "drizzle-kit push",
    "db:studio": "drizzle-kit studio",
    "test": "bun test",
    "deploy": "fly deploy"
  }
}
```

**Developer workflow:**
```bash
# First time setup
bun install
bun run setup

# Daily development
bun dev

# That's it! ✨
```

---

## Why Bun Instead of Node.js?

### For Small Teams, Bun Is a Game-Changer

| Feature | Bun | Node.js | Impact |
|---------|-----|---------|--------|
| **Startup time** | 3x faster | baseline | Dev iteration speed |
| **Test runner** | Built-in | Need Jest/Vitest | One less tool |
| **Package manager** | Built-in | Need pnpm/npm | One less tool |
| **Bundler** | Built-in | Need esbuild/webpack | One less tool |
| **TypeScript** | Native | Need ts-node/tsx | Zero config |
| **Watch mode** | Built-in | Need nodemon | One less tool |

**What this means:**
- No more `package.json` with 20 dev dependencies
- No configuration files (tsconfig, jest.config, webpack.config)
- Faster CI/CD (3x faster installs and tests)

**Example: Zero-config TypeScript Server**

```typescript
// src/index.ts
Bun.serve({
  port: 3000,
  async fetch(req) {
    return new Response("Hello from Bun!");
  },
  websocket: {
    open(ws) {
      console.log("Client connected");
    },
    message(ws, message) {
      ws.send(`Echo: ${message}`);
    },
  },
});

console.log("Server running on http://localhost:3000");
```

Run it:
```bash
bun run src/index.ts
```

No tsconfig, no build step, no webpack. Just works.

---

## Simplified Message Flow

### Before: Too Many Hops (Original Design)

```
Frontend → WebTransport → Durable Object → NATS → AI Service → NATS → Durable Object → Frontend
```

**Problems:**
- 6 network hops
- Debugging requires 4 services
- Durable Objects can't run locally

### After: Direct & Fast (Small Team Design)

```
Frontend → WebSocket → Bun Server → NATS → AI Handler → WebSocket → Frontend
```

**Benefits:**
- 3 network hops (2x faster)
- Single process to debug
- Runs identically locally

### Implementation

```typescript
// src/index.ts
import { connect } from "nats";
import Anthropic from "@anthropic-ai/sdk";

const nats = await connect({ servers: "localhost:4222" });
const claude = new Anthropic({ apiKey: process.env.ANTHROPIC_API_KEY });

// WebSocket server with message routing
Bun.serve({
  port: 3000,
  websocket: {
    async message(ws, msg) {
      const message = JSON.parse(msg as string);

      if (message.type === "chat") {
        // Stream AI response
        const stream = await claude.messages.stream({
          model: "claude-3-5-sonnet-20250122",
          messages: [{ role: "user", content: message.content }],
          max_tokens: 1024,
        });

        for await (const chunk of stream) {
          ws.send(JSON.stringify({
            type: "ai_chunk",
            content: chunk.delta?.text || "",
          }));
        }

        ws.send(JSON.stringify({ type: "ai_done" }));
      }

      if (message.type === "move") {
        // Publish to NATS for engine analysis
        nats.publish("chess.moves", JSON.stringify(message));
      }
    },
  },
});

// NATS subscriber for engine responses
const sub = nats.subscribe("chess.analysis");
for await (const msg of sub) {
  const analysis = JSON.parse(msg.data as string);
  // Broadcast to all connected clients
  broadcastToClients(analysis);
}
```

**This is 100x simpler than the original design, yet:**
- Still uses NATS for async processing
- Still scalable (Bun handles 100K+ concurrent WS)
- Runs perfectly locally

---

## Database: PostgreSQL + Drizzle

### Why PostgreSQL?

**For small teams:**
- ✅ SQL is universal (everyone knows it)
- ✅ ACID transactions (no data loss)
- ✅ JSON support (flexible schema)
- ✅ Full-text search built-in
- ✅ Works locally and in cloud

### Why Drizzle ORM?

Drizzle is a **type-safe SQL ORM** that feels like writing SQL but with TypeScript types.

**Example:**

```typescript
// db/schema.ts
import { pgTable, serial, text, timestamp, jsonb } from "drizzle-orm/pg-core";

export const games = pgTable("games", {
  id: serial("id").primaryKey(),
  pgn: text("pgn").notNull(),
  whitePlayer: text("white_player").notNull(),
  blackPlayer: text("black_player").notNull(),
  result: text("result"),
  metadata: jsonb("metadata"),
  createdAt: timestamp("created_at").defaultNow(),
});

export const chatHistory = pgTable("chat_history", {
  id: serial("id").primaryKey(),
  gameId: serial("game_id").references(() => games.id),
  role: text("role").notNull(), // 'user' | 'assistant'
  content: text("content").notNull(),
  timestamp: timestamp("timestamp").defaultNow(),
});
```

**Querying (fully type-safe):**

```typescript
import { db } from "./db/client";
import { games, chatHistory } from "./db/schema";
import { eq } from "drizzle-orm";

// Insert game
const newGame = await db.insert(games).values({
  pgn: "1. e4 e5 2. Nf3",
  whitePlayer: "Alice",
  blackPlayer: "Bob",
}).returning();

// Query with joins
const gameWithChat = await db
  .select()
  .from(games)
  .leftJoin(chatHistory, eq(games.id, chatHistory.gameId))
  .where(eq(games.id, 1));
```

**Benefits:**
- TypeScript autocomplete for columns
- Compile-time query validation
- Zero runtime overhead
- Migrations are just TypeScript files

---

## Deployment: Fly.io (Not Kubernetes)

### Why Fly.io for Small Teams?

**What Fly.io does:**
- Deploys Docker containers globally
- Manages scaling automatically
- Handles SSL/TLS certificates
- Provides databases (Postgres, Redis)
- **Zero Kubernetes knowledge required**

### One Command Deploy

```bash
# First time setup
fly launch

# Every deploy after
fly deploy
```

That's it. No Kubernetes YAML, no Load Balancers, no Ingress configs.

### fly.toml (Config File)

```toml
app = "chess-coach"
primary_region = "sjc"

[build]
  dockerfile = "Dockerfile"

[env]
  PORT = "3000"
  NODE_ENV = "production"

[[services]]
  internal_port = 3000
  protocol = "tcp"

  [[services.ports]]
    handlers = ["http"]
    port = 80
    force_https = true

  [[services.ports]]
    handlers = ["tls", "http"]
    port = 443

[http_service]
  auto_stop_machines = true
  auto_start_machines = true
  min_machines_running = 1
  max_machines_running = 10

[[vm]]
  memory = "1gb"
  cpu_kind = "shared"
  cpus = 1
```

### Add NATS and Redis

```bash
# Provision Postgres (managed)
fly postgres create --name chess-coach-db

# Add Redis (managed)
fly redis create --name chess-coach-redis

# NATS runs as separate Fly app
fly launch --config nats.toml
```

**Total infrastructure management time:** ~30 minutes
**Compared to Kubernetes:** ~3 days

---

## Observability: Keep It Simple

### Don't Over-Engineer Monitoring

For 2-3 developers, avoid:
- ❌ Datadog ($50-200/mo, complex setup)
- ❌ New Relic (overkill for MVP)
- ❌ Full ELK stack (operational nightmare)

### Use: Structured Logging + Basic Metrics

#### Pino for Logs (Built into Bun)

```typescript
import { pino } from "pino";

const logger = pino({
  level: process.env.LOG_LEVEL || "info",
  transport: {
    target: "pino-pretty",
    options: { colorize: true },
  },
});

logger.info({ userId: 123, gameId: 456 }, "Game started");
logger.error({ err }, "Failed to process move");
```

**Output:**
```
[2025-10-31 17:30:45] INFO  (12345): Game started
    userId: 123
    gameId: 456
```

**Centralize logs:**
```bash
# Production: Ship logs to Fly.io's log service (free)
fly logs

# Or use Betterstack (free tier)
```

#### Prometheus for Metrics

```typescript
import { register, Counter, Histogram } from "prom-client";

// Define metrics
const moveCounter = new Counter({
  name: "chess_moves_total",
  help: "Total chess moves processed",
});

const aiLatency = new Histogram({
  name: "ai_response_seconds",
  help: "AI response latency",
  buckets: [0.1, 0.5, 1, 2, 5],
});

// Use in code
moveCounter.inc();

const start = Date.now();
await callAI();
aiLatency.observe((Date.now() - start) / 1000);

// Expose metrics endpoint
Bun.serve({
  port: 9090,
  fetch() {
    return new Response(register.metrics(), {
      headers: { "Content-Type": register.contentType },
    });
  },
});
```

**Visualize with Grafana Cloud (free tier):**
- 10K metrics/month free
- Pre-built dashboards
- Alerts via Slack/email

---

## Testing Strategy for Small Teams

### Don't Write Too Many Tests (Controversial but True)

For 2-3 people:
- ❌ 80% test coverage is a waste of time
- ❌ Unit testing every function slows you down
- ✅ Focus on **critical path integration tests**

### What to Test

1. **Message protocol validation** (Zod schemas catch 90% of bugs)
2. **Move validation** (chess.js does this already)
3. **Critical user flows** (login → play game → save)

### Use Bun's Built-in Test Runner

```typescript
// tests/game.test.ts
import { test, expect } from "bun:test";
import { validateMove } from "../src/modules/chess/moves";

test("valid move is accepted", () => {
  const result = validateMove("e2", "e4", "startingFen");
  expect(result.valid).toBe(true);
});

test("invalid move is rejected", () => {
  const result = validateMove("e2", "e5", "startingFen");
  expect(result.valid).toBe(false);
});
```

Run tests:
```bash
bun test
```

**No Jest config, no Vitest setup. Just works.**

---

## Code Sharing: Frontend ↔ Backend

### Share TypeScript Types (No Code Generation)

With a monorepo, share types directly:

```
chess-coach/
├── frontend/
│   └── app/
│       └── src/
│           └── types/ → symlink to ../../shared/types
├── backend/
│   └── src/
│       └── types/ → symlink to ../../shared/types
└── shared/
    └── types/
        ├── messages.ts
        ├── game.ts
        └── user.ts
```

**Single source of truth:**

```typescript
// shared/types/messages.ts
export const MoveMessage = z.object({
  type: z.literal("move"),
  from: z.string(),
  to: z.string(),
});

export type MoveMessage = z.infer<typeof MoveMessage>;
```

**Frontend uses it:**
```typescript
import { MoveMessage } from "@/types/messages";

const move: MoveMessage = { type: "move", from: "e2", to: "e4" };
```

**Backend validates it:**
```typescript
import { MoveMessage } from "./types/messages";

const validated = MoveMessage.parse(message);
```

**No GraphQL codegen, no OpenAPI generation, no build step.**

---

## Development Workflow

### Typical Day for 2-3 Developers

**Morning:**
```bash
# Start infrastructure
docker-compose up -d

# Start backend (auto-reloads)
cd backend && bun dev

# Start frontend (in another terminal)
cd frontend/app && bun dev

# Both running in <5 seconds
```

**During development:**
- ✅ Change backend code → auto-reloads instantly
- ✅ Change frontend code → HMR updates browser
- ✅ All services work offline
- ✅ One breakpoint debugs everything

**Evening:**
```bash
# Commit and push
git add .
git commit -m "Add AI chat feature"
git push

# Deploy to production
fly deploy

# Live in 2 minutes
```

---

## Production Checklist (Small Team Edition)

### Week 1: MVP
- [ ] Docker Compose working locally
- [ ] Bun backend with WebSocket
- [ ] NATS pub/sub working
- [ ] Basic AI chat functional
- [ ] Deploy to Fly.io

### Week 2: Polish
- [ ] Add PostgreSQL persistence
- [ ] Stockfish engine integration
- [ ] Multiplayer game rooms (Redis)
- [ ] Structured logging (Pino)

### Week 3: Hardening
- [ ] Add Prometheus metrics
- [ ] Set up Grafana dashboards
- [ ] Add error tracking (Sentry free tier)
- [ ] Load test (100 concurrent users)

### Week 4: Launch
- [ ] Set up monitoring alerts
- [ ] Document deployment process
- [ ] Create backup strategy
- [ ] Soft launch to beta users

---

## Anti-Patterns to Avoid

### ❌ Don't Do This (Small Teams)

1. **Microservices from day 1**
   - Wait until you have >50K users
   - Monolith is faster to develop and debug

2. **Custom Docker images**
   - Use official images (nats:alpine, postgres:alpine)
   - Avoid "optimizing" until you measure

3. **Over-abstraction**
   - Don't create "repositories", "services", "factories"
   - Keep it simple: modules with functions

4. **Premature caching**
   - Don't add Redis caching until you measure
   - PostgreSQL is fast enough for MVP

5. **Building your own auth**
   - Use Clerk or Auth0 (free tiers)
   - Security is hard, let experts handle it

6. **Custom monitoring**
   - Use Fly.io logs + Prometheus
   - Don't build your own dashboards

---

## The "Works on My Machine" Guarantee

### Every Developer Should Be Able To:

```bash
# 1. Clone repo
git clone https://github.com/yourteam/chess-coach
cd chess-coach

# 2. Install dependencies
bun install

# 3. Start everything
bun run setup
bun dev

# 4. Open browser
open http://localhost:5173

# ✅ Full stack running in under 60 seconds
```

### No:
- ❌ "Install Kubernetes"
- ❌ "Set up VPN to staging"
- ❌ "Ask John for the AWS key"
- ❌ "Wait 30 minutes for Docker build"

### Yes:
- ✅ Works on Mac, Linux, Windows
- ✅ All services in Docker Compose
- ✅ Sample data included
- ✅ Environment variables in `.env.example`

---

## Cost Breakdown: Small Team Reality

### Development (Per Developer)

| Item | Cost |
|------|------|
| MacBook/PC | $0 (you own it) |
| Docker Desktop | Free |
| VS Code | Free |
| Git/GitHub | Free |
| **Total** | **$0/month** |

### Production (MVP: 0-10K users)

| Item | Service | Cost |
|------|---------|------|
| Backend hosting | Fly.io | $0-20/mo |
| Database | Fly Postgres | $0 (free tier) |
| Redis | Fly Redis | $0 (free tier) |
| AI API | Claude | ~$100/mo |
| Monitoring | Grafana Cloud | Free tier |
| Error tracking | Sentry | Free tier |
| **Total** | | **~$120/mo** |

### Production (Growth: 10K-100K users)

| Item | Service | Cost |
|------|---------|------|
| Backend hosting | Fly.io (scaled) | $100-200/mo |
| Database | Fly Postgres | $50/mo |
| Redis | Fly Redis | $25/mo |
| AI API | Claude | ~$500/mo |
| Monitoring | Grafana Cloud | $50/mo |
| **Total** | | **~$725-825/mo** |

**No Kubernetes, no DevOps team, no $10K/mo AWS bill.**

---

## When to Scale Up Architecture

### Monolith → Microservices

**Stay monolith until:**
- 50K+ concurrent users
- Backend uses >8GB RAM
- Deploy takes >5 minutes
- Team grows to 10+ engineers

**Then split into:**
- WebSocket service (stays fast)
- AI service (scales independently)
- Engine service (CPU-intensive)
- API service (CRUD operations)

### Fly.io → Kubernetes

**Move to K8s only when:**
- Multi-region required (beyond Fly's 30 regions)
- Need advanced networking (service mesh)
- Compliance requires self-hosting
- Budget >$5K/mo (at which point K8s is cheaper)

### Self-hosted NATS → NATS Cloud

**Use NATS Cloud when:**
- Managing NATS becomes time-consuming
- Need 99.99% uptime SLA
- Multi-region replication required

---

## Summary: Small Team Success Formula

### ✅ Do This

1. **Bun + TypeScript monolith** → Fast development
2. **Docker Compose** → Local development
3. **NATS + PostgreSQL + Redis** → Simple stack
4. **Fly.io** → One-command deploys
5. **Pino + Prometheus** → Basic observability

### ❌ Avoid This (For Now)

1. ~~Kubernetes~~ → Wait until >50K users
2. ~~Microservices~~ → Monolith is faster
3. ~~Cloudflare Durable Objects~~ → Can't debug locally
4. ~~Custom auth~~ → Use Clerk/Auth0
5. ~~Over-engineering~~ → Ship first, optimize later

### 🎯 Goals

- **First deploy:** Week 1
- **MVP with users:** Week 4
- **Break-even:** <$200/mo until revenue
- **Scale:** Proven before architectural complexity

---

## Next Steps

1. **Set up monorepo structure**
   ```bash
   chess-coach/
   ├── frontend/  (existing)
   ├── backend/   (new)
   └── shared/    (new - types)
   ```

2. **Create Docker Compose**
   - NATS
   - PostgreSQL
   - Redis

3. **Build Bun backend skeleton**
   - WebSocket server
   - NATS connection
   - Basic message routing

4. **Connect frontend**
   - WebSocket client
   - Message protocol (Zod)

5. **Deploy to Fly.io**
   - One command deploy
   - Verify works in production

---

**Philosophy:** Build for 2-3 developers today, scale for 100K users tomorrow.

**Reality Check:** Most startups fail before scaling becomes a problem. Optimize for speed and simplicity first.

---

**Document Version:** 1.0
**Last Updated:** October 31, 2025
**For:** 2-3 person development teams

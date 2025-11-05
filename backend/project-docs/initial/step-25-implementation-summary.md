# WebSocket Implementation - Quick Summary

**Full Plan**: [step-25-websocket-implementation-plan.md](step-25-websocket-implementation-plan.md)

---

## What We're Building

A production-ready WebSocket system for chess-coach with:

✅ Real-time game moves between players
✅ Game room management (multiple simultaneous games)
✅ Structured JSON message protocol
✅ Request/response pattern for actions
✅ Event broadcasting for real-time updates
✅ Connection management (from Lesson 2)
✅ Message validation and error handling

---

## Architecture Overview

```
Client (Browser)
    ↓
    ↓ WebSocket (ws://localhost:8080/ws)
    ↓
HTTP Server
    ├─ Upgrade Handler (JWT auth)
    ├─ Hub (Connection Registry)
    │   ├─ Room 1 (Game ABC) → Client A, Client B
    │   ├─ Room 2 (Game XYZ) → Client C, Client D
    │   └─ Message Handler
    │       ├─ createGame
    │       ├─ joinGame
    │       ├─ makeMove
    │       └─ chat
    └─ Database (Game State)
```

---

## Message Protocol

### Request/Response Example
```json
// Client → Server
{
  "id": "req-123",
  "type": "request",
  "action": "makeMove",
  "data": { "from": "e2", "to": "e4" }
}

// Server → Client
{
  "id": "req-123",
  "type": "response",
  "success": true
}
```

### Event Example
```json
// Server → All clients in room
{
  "type": "event",
  "event": "opponentMove",
  "data": {
    "from": "e7",
    "to": "e5",
    "timestamp": 1699000000
  }
}
```

---

## Implementation Phases

### Phase 1: Foundation (4-6 hours)
- Message protocol structs
- Client wrapper (readPump/writePump from Lesson 2)
- Hub with connection registry
- Room management

### Phase 2: Message Handlers (4-6 hours)
- createGame handler
- joinGame handler
- makeMove handler with validation
- getGameState handler
- chat handler

### Phase 3: HTTP Integration (2-3 hours)
- WebSocket upgrade endpoint
- JWT authentication
- Route registration

### Phase 4: Testing (3-4 hours)
- Unit tests (message, room, client)
- Integration tests (game flow)
- Manual testing

### Phase 5: Frontend (4-6 hours)
- React WebSocket hook
- Game component integration
- UI updates on events

### Phase 6: Production Polish (2-3 hours)
- Security hardening
- Monitoring/metrics
- Configuration
- Documentation

**Total: ~2-3 days**

---

## Key Files to Create

```
backend/internal/websocket/
├── message.go          # Message protocol + helpers
├── client.go           # Client struct (from Lesson 2)
├── hub.go              # Connection registry + rooms
├── room.go             # Game room management
└── handler.go          # Message routing + validation

backend/internal/api/
└── websocket.go        # HTTP upgrade handler

frontend/src/
├── hooks/
│   └── useWebSocket.ts # React hook
└── components/
    └── Game.tsx        # Game component
```

---

## Testing Strategy

### Unit Tests
```go
TestMessageCreation
TestRoomBroadcast
TestClientJoinRoom
TestMoveValidation
```

### Integration Tests
```go
TestCreateAndJoinGame
TestMakeMove
TestClientDisconnect
TestRoomCleanup
```

### Manual Tests
- Create game → Join game → Make moves
- Disconnect → Reconnect
- Invalid moves
- Multiple games simultaneously

---

## Quick Start Commands

```bash
# 1. Create directory structure
mkdir -p backend/internal/websocket
mkdir -p backend/internal/api

# 2. Copy code from plan
# (Use step-25-websocket-implementation-plan.md)

# 3. Run server
cd backend
go run cmd/server/main.go

# 4. Test with wscat
npm install -g wscat
wscat -c ws://localhost:8080/ws?token=test

# 5. Send test message
> {"type":"request","action":"createGame","data":{"timeControl":"5+0"}}
```

---

## Success Metrics

After implementation, you should be able to:

1. ✅ Open two browser tabs
2. ✅ Create game in tab 1
3. ✅ Join game in tab 2
4. ✅ Make move in tab 1 → Instantly see in tab 2
5. ✅ Make move in tab 2 → Instantly see in tab 1
6. ✅ Close tab 1 → Reconnect → Resume game
7. ✅ Handle 50+ concurrent games

---

## Integration with Existing Code

### Database
```go
// Reuse existing game models
type Game struct {
    ID          string
    WhiteUserID string
    BlackUserID string
    BoardState  string
    Status      string
    CreatedAt   time.Time
}
```

### Authentication
```go
// Reuse existing JWT validation
func authenticateRequest(r *http.Request) (string, error) {
    // Use your existing auth package
    return auth.ValidateToken(token)
}
```

### Game Logic
```go
// Reuse existing chess validation
func validateMove(from, to string, board Board) error {
    // Use your existing chess package
    return chess.ValidateMove(from, to, board)
}
```

---

## Common Pitfalls to Avoid

❌ **Don't** send game state in every message (it's large)
✅ **Do** send only move deltas

❌ **Don't** trust client-provided game state
✅ **Do** validate everything server-side

❌ **Don't** block the hub's event loop
✅ **Do** handle messages asynchronously

❌ **Don't** forget to clean up rooms
✅ **Do** remove empty rooms after players leave

---

## Resources

- **Lessons**: `/learnings/fundamentals/websocket/`
  - Lesson 1: Basics
  - Lesson 2: Connection Management
  - Lesson 3: Message Protocols

- **Full Plan**: `step-25-websocket-implementation-plan.md`

- **Gorilla WebSocket Docs**: https://pkg.go.dev/github.com/gorilla/websocket

---

## Next Actions

1. **Read** the full implementation plan
2. **Set up** project structure
3. **Implement** Phase 1 (Foundation)
4. **Test** as you go
5. **Iterate** based on learnings

**Let's build real-time chess! ♟️🚀**

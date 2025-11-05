# WebSocket Implementation Checklist

**Purpose**: Track implementation progress step-by-step

**How to use**: Check off items as you complete them

---

## 🎯 Overall Progress

- [ ] Phase 1: Foundation (0/15)
- [ ] Phase 2: Integration (0/12)
- [ ] Phase 3: Utilities (0/8)
- [ ] Phase 4: Handlers (0/16)
- [ ] Phase 5: Authentication (0/10)
- [ ] Phase 6: Testing (0/12)
- [ ] Phase 7: Production (0/8)

**Total**: 0/81 tasks complete

---

## Phase 1: Foundation (4-5 hours)

**Goal**: Create core WebSocket components that compile

**Document**: [`01-foundation.md`](./01-foundation.md)

### Setup
- [ ] Install `gorilla/websocket`: `go get github.com/gorilla/websocket`
- [ ] Install `google/uuid`: `go get github.com/google/uuid`
- [ ] Create directory: `mkdir -p internal/websocket`

### File: `internal/websocket/message.go`
- [ ] Create file with package declaration
- [ ] Add imports: `encoding/json`, `time`, `github.com/google/uuid`
- [ ] Define `MessageType` enum (request, response, event)
- [ ] Define `Message` struct with all fields
- [ ] Implement `NewRequest()` function
- [ ] Implement `NewResponse()` function
- [ ] Implement `NewErrorResponse()` function
- [ ] Implement `NewEvent()` function
- [ ] Implement `ToJSON()` method
- [ ] Implement `FromJSON()` function
- [ ] Implement `generateID()` utility function
- [ ] Test: Run `go build ./internal/websocket`

### File: `internal/websocket/client.go`
- [ ] Create file with package declaration
- [ ] Add imports: `fmt`, `log`, `sync`, `time`, `websocket`
- [ ] Define constants (writeWait, pongWait, pingPeriod, maxMessageSize)
- [ ] Define `Client` struct with all fields
- [ ] Implement `NewClient()` constructor
- [ ] Implement `readPump()` method with context
- [ ] Implement `writePump()` method
- [ ] Implement `SendMessage()` method
- [ ] Implement `JoinRoom()` method with error return
- [ ] Implement `LeaveRoom()` method
- [ ] Test: Run `go build ./internal/websocket`

### File: `internal/websocket/hub.go`
- [ ] Create file with package declaration
- [ ] Add imports: `log`, `sync`
- [ ] Define `Hub` struct with all fields
- [ ] Add `shutdown` channel to Hub
- [ ] Implement `NewHub()` constructor
- [ ] Implement `Run()` method with shutdown case
- [ ] Implement `GetOrCreateRoom()` method
- [ ] Implement `RemoveRoom()` method
- [ ] Implement `CleanupEmptyRooms()` method
- [ ] Implement `Shutdown()` method
- [ ] Implement `handleMessage()` routing method
- [ ] Test: Run `go build ./internal/websocket`

### File: `internal/websocket/room.go`
- [ ] Create file with package declaration
- [ ] Add imports: `log`, `sync`
- [ ] Define `Room` struct
- [ ] Implement `NewRoom()` constructor
- [ ] Implement `AddClient()` method
- [ ] Implement `RemoveClient()` method
- [ ] Implement `Broadcast()` method
- [ ] Implement `BroadcastExcept()` method
- [ ] Implement `ClientCount()` method
- [ ] Test: Run `go build ./internal/websocket`

### Verification
- [ ] All files compile without errors
- [ ] No import errors
- [ ] No missing function errors
- [ ] Run: `go vet ./internal/websocket`
- [ ] Run: `go build ./...`

---

## Phase 2: Integration (2-3 hours)

**Goal**: Connect WebSocket to existing app

**Document**: [`02-integration.md`](./02-integration.md)

### File: `internal/websocket/upgrade.go`
- [ ] Create file with package declaration
- [ ] Add imports
- [ ] Define `upgrader` with placeholder CORS
- [ ] Define `Handler` struct with dependencies
- [ ] Implement `NewHandler()` constructor
- [ ] Implement `ServeWS()` method (with placeholder auth)
- [ ] Implement message routing in `handleMessage()`
- [ ] Test: Run `go build ./internal/websocket`

### File: `internal/app/app.go` (Modify)
- [ ] Import `internal/websocket`
- [ ] Add `wsHub *websocket.Hub` field to `App` struct
- [ ] Create hub in `New()` function
- [ ] Pass hub to `server.New()`
- [ ] Add `go a.wsHub.Run()` in `Run()` method
- [ ] Add `a.wsHub.Shutdown()` in shutdown logic
- [ ] Test: Run `go build ./internal/app`

### File: `internal/server/server.go` (Modify)
- [ ] Add hub parameter to `New()` function signature
- [ ] Pass hub to `router.Setup()`
- [ ] Test: Run `go build ./internal/server`

### File: `internal/router/router.go` (Modify)
- [ ] Add hub parameter to `Setup()` function signature
- [ ] Import `internal/websocket`
- [ ] Create WebSocket handler
- [ ] Register `/ws` endpoint before API routes
- [ ] Test: Run `go build ./internal/router`

### File: `cmd/server/main.go` (Modify)
- [ ] No changes needed (hub created in app)
- [ ] Test: Run `go build ./cmd/server`

### Verification
- [ ] Server compiles and starts
- [ ] Navigate to `http://localhost:8080/ws` (should get "Bad Request" or similar)
- [ ] Check logs for "Hub shutting down" on Ctrl+C
- [ ] Run: `go build ./...`

---

## Phase 3: Utilities (1-2 hours)

**Goal**: Type-safe helper functions

**Document**: [`03-utilities.md`](./03-utilities.md)

### File: `internal/websocket/utils.go`
- [ ] Create file with package declaration
- [ ] Add imports: `fmt`, `github.com/google/uuid`, `github.com/jackc/pgx/v5/pgtype`
- [ ] Implement `UUIDToString()` function
- [ ] Implement `StringToUUID()` function
- [ ] Implement `ParseTimeControl()` function
- [ ] Add error handling for invalid inputs
- [ ] Test: Write unit tests for each function
- [ ] Test: Run `go test ./internal/websocket`

### Context Propagation
- [ ] Update `readPump()` to accept context
- [ ] Update `handleMessage()` to use context
- [ ] Verify context flows to all handlers
- [ ] Test: Run `go build ./internal/websocket`

---

## Phase 4: Handlers (6-8 hours)

**Goal**: Implement complete business logic

**Document**: [`04-message-handlers.md`](./04-message-handlers.md)

### File: `internal/websocket/handler.go` (Modify)
- [ ] Add repository fields to `MessageHandler` struct
- [ ] Update `NewMessageHandler()` to accept repositories
- [ ] Add database field
- [ ] Test: Run `go build ./internal/websocket`

### Handler: `handleCreateGame`
- [ ] Extract timeControl from message
- [ ] Validate input parameters
- [ ] Convert user ID to UUID
- [ ] Call `gameRepo.CreateGame()` with context
- [ ] Convert game ID back to string
- [ ] Create or join room (based on game type)
- [ ] Send success response to client
- [ ] Broadcast gameCreated event
- [ ] Handle errors at each step
- [ ] Test: Manual test with wscat

### Handler: `handleJoinGame`
- [ ] Extract gameId from message
- [ ] Validate game exists in database
- [ ] Check room capacity (max 2 players)
- [ ] Add client to room
- [ ] Update game in database (add black player)
- [ ] Send success response to joiner
- [ ] Broadcast playerJoined event to room
- [ ] Broadcast gameStarted if room full
- [ ] Handle errors
- [ ] Test: Manual test with two wscat connections

### Handler: `handleMakeMove`
- [ ] Extract from/to from message
- [ ] Validate client is in a game
- [ ] Get current game state from database
- [ ] Validate move is legal (placeholder for now)
- [ ] Create move in database
- [ ] Update game PGN
- [ ] Send success response to player
- [ ] Broadcast opponentMove event
- [ ] Handle errors
- [ ] Test: Make moves in test game

### Handler: `handleLeaveGame`
- [ ] Remove client from room
- [ ] Broadcast playerLeft event
- [ ] Clean up empty rooms
- [ ] Test: Disconnect and verify cleanup

### Error Handling
- [ ] Consistent error response format
- [ ] Log all errors server-side
- [ ] Don't expose internal errors to client
- [ ] Test: Send invalid messages, verify responses

---

## Phase 5: Authentication (2-3 hours)

**Goal**: Secure WebSocket connections

**Document**: [`05-authentication.md`](./05-authentication.md)

### JWT Implementation
- [ ] Install JWT library: `go get github.com/golang-jwt/jwt/v5`
- [ ] Define JWT secret in config
- [ ] Implement `authenticateRequest()` function
- [ ] Parse token from query or header
- [ ] Validate JWT signature
- [ ] Extract user ID from claims
- [ ] Return user ID or error
- [ ] Update `ServeWS()` to use real auth
- [ ] Test: Connect with invalid token (should reject)
- [ ] Test: Connect with valid token (should accept)

### CORS Configuration
- [ ] Define allowed origins list
- [ ] Implement proper `CheckOrigin` function
- [ ] Parse origin header
- [ ] Check against whitelist
- [ ] Test: Connect from allowed origin
- [ ] Test: Connect from disallowed origin (should reject)

---

## Phase 6: Testing (4-6 hours)

**Goal**: Comprehensive test coverage

**Document**: [`06-testing.md`](./06-testing.md)

### Unit Tests: `message_test.go`
- [ ] Test `NewRequest()` creates valid request
- [ ] Test `NewResponse()` creates valid response
- [ ] Test `NewEvent()` creates valid event
- [ ] Test `ToJSON()` serialization
- [ ] Test `FromJSON()` deserialization
- [ ] Run: `go test ./internal/websocket -v -run TestMessage`

### Unit Tests: `room_test.go`
- [ ] Test room creation
- [ ] Test adding clients
- [ ] Test removing clients
- [ ] Test broadcast to all
- [ ] Test broadcast except sender
- [ ] Test client count
- [ ] Run: `go test ./internal/websocket -v -run TestRoom`

### Unit Tests: `utils_test.go`
- [ ] Test UUID conversions
- [ ] Test time control parsing
- [ ] Test error cases
- [ ] Run: `go test ./internal/websocket -v -run TestUtils`

### Integration Tests: `integration_test.go`
- [ ] Test connect to WebSocket
- [ ] Test create game flow
- [ ] Test join game flow
- [ ] Test make move flow
- [ ] Test disconnect cleanup
- [ ] Run: `go test ./internal/websocket -v -run TestIntegration`

### Manual Testing with wscat
- [ ] Install: `npm install -g wscat`
- [ ] Connect: `wscat -c 'ws://localhost:8080/ws?token=test'`
- [ ] Send createGame message
- [ ] Send joinGame message
- [ ] Send makeMove message
- [ ] Verify responses
- [ ] Test with two connections

### Coverage Check
- [ ] Run: `go test ./internal/websocket -coverprofile=coverage.out`
- [ ] Run: `go tool cover -html=coverage.out`
- [ ] Verify 80%+ coverage

---

## Phase 7: Production (2-3 hours)

**Goal**: Production-ready deployment

**Document**: [`07-deployment.md`](./07-deployment.md)

### Configuration
- [ ] Move JWT secret to environment variable
- [ ] Move allowed origins to config
- [ ] Add WebSocket timeouts to config
- [ ] Add max connections limit
- [ ] Test: Start with different configs

### Logging
- [ ] Add structured logging for connections
- [ ] Log authentication attempts
- [ ] Log game creation/joining
- [ ] Log move broadcasts
- [ ] Log errors with context
- [ ] Test: Review logs for completeness

### Monitoring
- [ ] Add metrics for active connections
- [ ] Add metrics for messages sent/received
- [ ] Add metrics for errors
- [ ] Expose metrics endpoint (if applicable)

### Documentation
- [ ] Document WebSocket message protocol
- [ ] Document authentication flow
- [ ] Document error codes
- [ ] Add API examples
- [ ] Update main README

### Deployment Checklist
- [ ] Code review completed
- [ ] All tests passing
- [ ] Coverage 80%+
- [ ] No security vulnerabilities
- [ ] Configuration externalized
- [ ] Logging comprehensive
- [ ] Graceful shutdown working
- [ ] Ready for production

---

## 🎉 Completion Criteria

Mark this checklist complete when:

- [x] All 81 tasks checked off
- [x] Server starts without errors
- [x] WebSocket endpoint accessible
- [x] Authentication working
- [x] Create/join/move game flow working
- [x] All data persisted to database
- [x] Tests passing (80%+ coverage)
- [x] Code reviewed
- [x] Documentation updated
- [x] Ready for production deployment

---

## 📊 Progress Tracking

Update this after completing each phase:

| Phase | Started | Completed | Duration | Notes |
|-------|---------|-----------|----------|-------|
| Phase 1 | ____-__-__ | ____-__-__ | __h | |
| Phase 2 | ____-__-__ | ____-__-__ | __h | |
| Phase 3 | ____-__-__ | ____-__-__ | __h | |
| Phase 4 | ____-__-__ | ____-__-__ | __h | |
| Phase 5 | ____-__-__ | ____-__-__ | __h | |
| Phase 6 | ____-__-__ | ____-__-__ | __h | |
| Phase 7 | ____-__-__ | ____-__-__ | __h | |

**Total Time**: ___ hours

---

## 🆘 Troubleshooting

If stuck, check:

1. **Won't compile**: Check imports and missing functions
2. **Won't connect**: Check authentication and CORS
3. **Messages not received**: Check hub is running
4. **Database errors**: Check UUID conversions and context
5. **Tests failing**: Check test setup and mocks

Refer to phase-specific documents for detailed help.

---

**Last Updated**: 2025-11-05
**Status**: Ready for implementation

# Lesson 3: Message Protocols & State Management

**Goal**: Learn how to structure WebSocket communication with proper message protocols, implement request/response patterns, and manage application state for real-time multiplayer applications like chess-coach.

**Prerequisites**: Completed Lesson 1 (Basics) and Lesson 2 (Connection Management)

**Time**: 2-3 hours

---

## What You'll Learn

By the end of this lesson, you'll understand:

1. **Message Protocols** - How to structure WebSocket messages
2. **Request/Response Pattern** - Implementing RPC-style communication
3. **Event-Based Messaging** - Pub/sub pattern for broadcasting
4. **Message Validation** - Ensuring data integrity and security
5. **Room/Channel System** - Managing multiple game sessions
6. **State Synchronization** - Keeping clients in sync

---

## Why Message Protocols Matter

### The Problem

In Lesson 1 and 2, we sent simple text messages:

```javascript
// ❌ Unstructured - hard to handle different message types
ws.send("Hello");
ws.send("move e2e4");
ws.send("join game123");
```

How does the server know what these mean?
- Is "Hello" a chat message or a command?
- Is "move e2e4" valid chess notation?
- What if the client sends garbage?

### The Solution: Structured Messages

```javascript
// ✅ Structured - clear intent and data
ws.send(JSON.stringify({
    type: "chat",
    data: { message: "Hello" }
}));

ws.send(JSON.stringify({
    type: "move",
    data: { from: "e2", to: "e4" }
}));

ws.send(JSON.stringify({
    type: "join",
    data: { gameId: "game123" }
}));
```

Now the server can:
- ✅ Identify message type instantly
- ✅ Validate message structure
- ✅ Route to correct handler
- ✅ Return meaningful errors

---

## Message Protocol Patterns

### Pattern 1: Simple Event-Based (Pub/Sub)

**Use case**: Broadcasting updates to multiple clients

```javascript
// Client → Server
{
    "type": "move",
    "data": { "from": "e2", "to": "e4" }
}

// Server → All clients in game
{
    "type": "move",
    "data": {
        "player": "white",
        "from": "e2",
        "to": "e4",
        "timestamp": 1699000000
    }
}
```

**Characteristics**:
- One-way communication
- No response expected
- Used for notifications and broadcasts
- Simple and fast

**Example use cases**:
- Player makes chess move → broadcast to opponent
- Chat message → broadcast to all players in room
- Player joins → notify everyone

---

### Pattern 2: Request/Response (RPC-Style)

**Use case**: Operations that need confirmation or results

```javascript
// Client → Server (Request)
{
    "id": "req-123",
    "type": "request",
    "action": "createGame",
    "data": {
        "timeControl": "5+0"
    }
}

// Server → Client (Response)
{
    "id": "req-123",
    "type": "response",
    "success": true,
    "data": {
        "gameId": "game-abc-123",
        "url": "/game/game-abc-123"
    }
}
```

**Characteristics**:
- Two-way communication
- Response matches request via ID
- Can handle success or error
- Client can timeout if no response

**Example use cases**:
- Create game → get game ID
- Submit move → get validation result
- Request game state → get current board
- Authentication → get token

---

### Pattern 3: Hybrid (Event + Request/Response)

**Use case**: Most production applications (including chess-coach!)

```javascript
// Requests (need response)
{
    "id": "req-456",
    "type": "request",
    "action": "makeMove",
    "data": { "from": "e2", "to": "e4" }
}

// Events (no response needed)
{
    "type": "event",
    "event": "opponentMove",
    "data": { "from": "e7", "to": "e5" }
}

// Responses (reply to requests)
{
    "id": "req-456",
    "type": "response",
    "success": true
}
```

**Best of both worlds**:
- Requests for critical operations
- Events for real-time updates
- Flexible and scalable

---

## Message Structure Design

### Option 1: Minimal (What we'll use)

```json
{
    "type": "message_type",
    "data": { }
}
```

**Pros**: Simple, small payload, easy to understand
**Cons**: Limited metadata

---

### Option 2: Full-Featured (Production)

```json
{
    "id": "msg-uuid",
    "type": "request|response|event",
    "action": "specific_action",
    "timestamp": 1699000000,
    "userId": "user123",
    "data": { },
    "metadata": {
        "version": "1.0",
        "trace": "trace-id"
    }
}
```

**Pros**: Complete, traceable, versioned
**Cons**: Larger payload, more complex

---

### For Chess-Coach, We'll Use:

```json
{
    "id": "msg-uuid",           // For request/response matching
    "type": "request|response|event",
    "action": "specific_action",  // What to do
    "data": { },                // Payload
    "error": "error_message"    // Only in error responses
}
```

**Balanced**: Has what we need, not too complex

---

## Room/Channel System

### Why Rooms?

In chess-coach, you'll have multiple simultaneous games:

```
Server
├── Game Room 1 (Alice vs Bob)
│   ├── Alice's connection
│   └── Bob's connection
│
├── Game Room 2 (Carol vs Dave)
│   ├── Carol's connection
│   └── Dave's connection
│
└── Game Room 3 (Eve vs Frank)
    ├── Eve's connection
    └── Frank's connection
```

**Without rooms**: A move in Game 1 would broadcast to ALL players (chaos!)
**With rooms**: A move in Game 1 only goes to Alice and Bob (perfect!)

---

### Room Implementation Pattern

```go
type Room struct {
    ID      string
    Clients map[*Client]bool
    mu      sync.RWMutex
}

type Hub struct {
    Rooms map[string]*Room
    mu    sync.RWMutex
}

// Client joins room
func (h *Hub) JoinRoom(client *Client, roomID string) {
    h.mu.Lock()
    defer h.mu.Unlock()

    room, exists := h.Rooms[roomID]
    if !exists {
        room = &Room{
            ID:      roomID,
            Clients: make(map[*Client]bool),
        }
        h.Rooms[roomID] = room
    }

    room.mu.Lock()
    room.Clients[client] = true
    room.mu.Unlock()

    client.Room = room
}

// Broadcast to room only
func (r *Room) Broadcast(message []byte) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    for client := range r.Clients {
        select {
        case client.send <- message:
        default:
            // Client buffer full, skip
        }
    }
}
```

---

## Message Validation

### Why Validate?

**Security**: Prevent malicious input
**Stability**: Prevent crashes from bad data
**User Experience**: Give clear error messages

### Validation Layers

```
Message arrives
    ↓
1. JSON validation
    ↓ valid JSON?
2. Schema validation
    ↓ has required fields?
3. Business logic validation
    ↓ is move legal?
4. Process message
```

---

### Example: Chess Move Validation

```go
type MoveMessage struct {
    Type string `json:"type"`
    Data struct {
        From string `json:"from"`
        To   string `json:"to"`
    } `json:"data"`
}

func validateMove(msg MoveMessage) error {
    // 1. Check message type
    if msg.Type != "move" {
        return errors.New("invalid message type")
    }

    // 2. Check required fields
    if msg.Data.From == "" || msg.Data.To == "" {
        return errors.New("from and to are required")
    }

    // 3. Check chess notation (e.g., "e2" to "e4")
    if !isValidSquare(msg.Data.From) || !isValidSquare(msg.Data.To) {
        return errors.New("invalid chess notation")
    }

    // 4. Check if move is legal (requires game state)
    if !isLegalMove(msg.Data.From, msg.Data.To) {
        return errors.New("illegal move")
    }

    return nil
}
```

---

## State Management

### The Challenge

WebSocket is **stateless** - each message is independent. But games need **state**:

```
Current board position
Whose turn it is
Game status (active, checkmate, draw)
Move history
Time remaining
```

### Solution: Server-Side State

```go
type GameState struct {
    ID          string
    WhitePlayer *Client
    BlackPlayer *Client
    Board       [8][8]string
    Turn        string // "white" or "black"
    Status      string // "active", "checkmate", "draw"
    Moves       []Move
    CreatedAt   time.Time
}

var Games = make(map[string]*GameState)
var GamesMutex sync.RWMutex
```

---

### State Operations

**1. Create State**
```go
game := &GameState{
    ID:     generateID(),
    Board:  initializeBoard(),
    Turn:   "white",
    Status: "active",
}
Games[game.ID] = game
```

**2. Update State**
```go
game.Board[from] = ""
game.Board[to] = piece
game.Turn = oppositeColor(game.Turn)
game.Moves = append(game.Moves, move)
```

**3. Read State**
```go
game, exists := Games[gameID]
if !exists {
    return errors.New("game not found")
}
```

**4. Delete State** (game ends)
```go
delete(Games, gameID)
```

---

### State Synchronization Pattern

When a player makes a move:

```
1. Client sends move message
    ↓
2. Server validates move
    ↓
3. Server updates game state
    ↓
4. Server broadcasts new state to both players
    ↓
5. Both clients update their UI
```

**Key**: Server is the single source of truth!

---

## Error Handling

### Error Response Format

```json
{
    "id": "req-123",
    "type": "response",
    "success": false,
    "error": "Invalid move: piece cannot move that way"
}
```

### Error Categories

**1. Protocol Errors** (400-level)
```json
{ "error": "Invalid JSON" }
{ "error": "Missing required field: type" }
{ "error": "Unknown message type: xyz" }
```

**2. Business Logic Errors** (400-level)
```json
{ "error": "Game not found" }
{ "error": "Not your turn" }
{ "error": "Illegal move" }
```

**3. Server Errors** (500-level)
```json
{ "error": "Internal server error" }
{ "error": "Database unavailable" }
```

---

## Message Flow Examples

### Example 1: Creating and Joining a Game

```
Player 1 (Alice)                Server                Player 2 (Bob)
      │                           │                          │
      │  [REQUEST: create game]   │                          │
      │ ────────────────────────> │                          │
      │                           │                          │
      │                           ├─ Create game room        │
      │                           ├─ Generate game ID        │
      │                           ├─ Add Alice to room       │
      │                           │                          │
      │  [RESPONSE: game created] │                          │
      │ <──────────────────────── │                          │
      │   gameId: "game-123"      │                          │
      │                           │                          │
      │                           │  [REQUEST: join game-123]│
      │                           │ <────────────────────────│
      │                           │                          │
      │                           ├─ Add Bob to room         │
      │                           │                          │
      │  [EVENT: player joined]   │  [RESPONSE: joined]      │
      │ <──────────────────────── │ ────────────────────────>│
      │   player: "Bob"           │   gameId: "game-123"     │
      │                           │                          │
```

---

### Example 2: Making a Move

```
Player 1 (White)                Server                Player 2 (Black)
      │                           │                          │
      │  [REQUEST: move e2→e4]    │                          │
      │ ────────────────────────> │                          │
      │                           │                          │
      │                           ├─ Validate move          │
      │                           ├─ Update game state      │
      │                           │                          │
      │  [RESPONSE: success]      │  [EVENT: move made]      │
      │ <──────────────────────── │ ────────────────────────>│
      │                           │   from: "e2"             │
      │  [EVENT: move made]       │   to: "e4"               │
      │ <──────────────────────── │   turn: "black"          │
      │   (confirmation)          │                          │
      │                           │                          │
```

---

### Example 3: Invalid Move

```
Player 1                        Server                Player 2
      │                           │                          │
      │  [REQUEST: move e2→e5]    │                          │
      │ ────────────────────────> │                          │
      │   (illegal move!)         │                          │
      │                           │                          │
      │                           ├─ Validate move          │
      │                           ├─ ❌ Invalid!            │
      │                           │                          │
      │  [RESPONSE: error]        │                          │
      │ <──────────────────────── │                          │
      │   "Illegal move"          │                          │
      │                           │                          │
      │  (no state change)        │  (no notification)       │
```

---

## Production Considerations

### 1. Message Size Limits

```go
const MaxMessageSize = 1 * 1024 * 1024 // 1 MB

conn.SetReadLimit(MaxMessageSize)
```

For chess-coach:
- Move message: ~100 bytes
- Game state: ~5 KB (full board + history)
- 1 MB is plenty!

---

### 2. Rate Limiting

Prevent spam/abuse:

```go
type RateLimiter struct {
    lastMessage time.Time
    messageCount int
}

const MaxMessagesPerSecond = 10

func (r *RateLimiter) AllowMessage() bool {
    now := time.Now()
    if now.Sub(r.lastMessage) > time.Second {
        r.messageCount = 0
        r.lastMessage = now
    }

    r.messageCount++
    return r.messageCount <= MaxMessagesPerSecond
}
```

---

### 3. Message Queue Size

```go
// From Lesson 2
client.send = make(chan []byte, 256)

// For chess-coach, this is plenty:
// - Moves are infrequent (1-2 per second max)
// - 256 message buffer = 2+ minutes of buffering
```

---

### 4. Timeout Handling

```go
// Request timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

select {
case response := <-responseChan:
    // Got response
case <-ctx.Done():
    // Timeout - no response received
}
```

---

## Comparison: Message Protocols

| Pattern | Use Case | Pros | Cons |
|---------|----------|------|------|
| **Raw Text** | Very simple apps | Simple | No structure, hard to parse |
| **JSON** | Most web apps | Human-readable, flexible | Larger payload, slower parsing |
| **MessagePack** | Performance-critical | Smaller, faster than JSON | Binary, not human-readable |
| **Protocol Buffers** | High-scale, typed | Smallest, fastest, type-safe | Requires schema, compilation |
| **Custom Binary** | Specific needs | Optimized for use case | Hard to maintain |

**For chess-coach**: JSON is perfect!
- Readable during development
- Flexible for iteration
- Performance is fine (moves are infrequent)
- Small payload (moves are tiny)

---

## Chess-Coach Message Catalog

Here's what messages you'll need:

### Authentication
```json
{ "type": "request", "action": "login", "data": { "token": "..." } }
{ "type": "response", "success": true, "data": { "userId": "..." } }
```

### Game Management
```json
{ "type": "request", "action": "createGame", "data": { "timeControl": "5+0" } }
{ "type": "request", "action": "joinGame", "data": { "gameId": "..." } }
{ "type": "event", "event": "gameStarted", "data": { "color": "white" } }
```

### Moves
```json
{ "type": "request", "action": "move", "data": { "from": "e2", "to": "e4" } }
{ "type": "event", "event": "opponentMove", "data": { "from": "e7", "to": "e5" } }
```

### Game State
```json
{ "type": "request", "action": "getState" }
{ "type": "response", "data": { "board": [...], "turn": "white" } }
```

### Chat
```json
{ "type": "event", "event": "chat", "data": { "message": "Good game!" } }
```

### Game End
```json
{ "type": "event", "event": "gameOver", "data": { "result": "checkmate", "winner": "white" } }
```

---

## What You'll Build

In this lesson's code, you'll implement:

1. **Message Protocol** - JSON-based with types and validation
2. **Request/Response** - With message ID matching
3. **Room System** - Multiple isolated game rooms
4. **State Management** - Server-side game state
5. **Event Broadcasting** - Notify room members
6. **Error Handling** - Clear error messages
7. **Interactive Client** - Test all patterns

---

## Next Steps

1. Review this README thoroughly
2. Run the example code (coming next!)
3. Complete the experiments
4. Apply to your chess-coach project

**Ready to code?** Let's build it! 🚀

---

## Key Takeaways

✅ **Structure your messages** - Use JSON with clear types
✅ **Mix patterns** - Use both request/response and events
✅ **Validate everything** - Never trust client input
✅ **Use rooms** - Isolate game sessions
✅ **Server is truth** - Maintain state server-side
✅ **Handle errors gracefully** - Give clear feedback
✅ **Keep it simple** - Don't over-engineer

**You're ready to build production-ready WebSocket applications!** 🎉

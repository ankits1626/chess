# Lesson 1: What is WebSocket?

**Goal**: Understand what WebSocket is, why it exists, and when to use it

**Time**: 30-45 minutes

---

## The Problem WebSocket Solves

### Traditional HTTP Request-Response

Imagine your chess game using only REST APIs:

```
User makes move
    ↓
Frontend: POST /api/games/123/moves
    ↓
Backend: Saves move, returns success
    ↓
Frontend: Wait 2 seconds... GET /api/games/123/state
    ↓
Backend: Returns game state
    ↓
Frontend: Check if computer moved... No?
    ↓
Frontend: Wait 2 seconds... GET /api/games/123/state
    ↓
... (repeat polling)
```

**Problems:**
- ❌ Constant polling wastes bandwidth
- ❌ Delays in seeing opponent moves
- ❌ Server load from repeated requests
- ❌ Can't push notifications to client
- ❌ Terrible user experience

### WebSocket Solution

```
User connects once
    ↓
Persistent bidirectional connection established
    ↓
User makes move → Send instantly
    ↓
Computer responds → Receive instantly
    ↓
No polling needed!
```

**Benefits:**
- ✅ Real-time bidirectional communication
- ✅ One connection, unlimited messages
- ✅ Server can push to client anytime
- ✅ Low latency (perfect for games)
- ✅ Efficient (no HTTP overhead per message)

---

## What is WebSocket?

### Simple Definition

**WebSocket** is a communication protocol that provides **full-duplex** (two-way) communication channels over a single TCP connection.

### Key Characteristics

1. **Persistent Connection**
   - Opens once, stays open
   - Unlike HTTP: request → response → close

2. **Full-Duplex**
   - Both sides can send anytime
   - No request-response requirement
   - Like a phone call, not letters

3. **Low Overhead**
   - Initial handshake uses HTTP
   - After that, minimal framing
   - No HTTP headers on every message

4. **Real-Time**
   - Messages delivered immediately
   - No polling or waiting
   - Sub-millisecond latency

---

## HTTP vs WebSocket

### Visual Comparison

#### HTTP (Traditional)

```
Client                          Server
  |                               |
  |  GET /api/data   -----------> |
  |                               | (process)
  |  <----------- response        |
  |                               |
  | (connection closes)           |
  |                               |
  |  GET /api/data   -----------> | (new connection)
  |                               |
  |  <----------- response        |
  |                               |
```

**Each request = new connection**

#### WebSocket

```
Client                          Server
  |                               |
  |  HTTP Upgrade   ----------->  |
  |  <----------- 101 Switching   |
  |                               |
  |═══════════════════════════════| (persistent connection)
  |                               |
  |  message 1      ----------->  |
  |  <-----------   response      |
  |  message 2      ----------->  |
  |  <-----------   push msg      |
  |                               |
```

**One connection = many messages**

---

## The WebSocket Handshake

WebSocket starts as HTTP, then upgrades!

### Client Request

```http
GET /ws HTTP/1.1
Host: localhost:8080
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
Sec-WebSocket-Version: 13
```

Key headers:
- `Upgrade: websocket` - "I want to upgrade to WebSocket"
- `Connection: Upgrade` - "Don't close after response"
- `Sec-WebSocket-Key` - Random value for security
- `Sec-WebSocket-Version` - WebSocket protocol version

### Server Response

```http
HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
```

- `101 Switching Protocols` - "Okay, we're now WebSocket!"
- `Sec-WebSocket-Accept` - Proves server understands WebSocket

**After this, it's no longer HTTP! Pure WebSocket frames.**

---

## WebSocket Frame Structure

After handshake, messages are sent as frames:

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-------+-+-------------+-------------------------------+
|F|R|R|R| opcode|M| Payload len |    Extended payload length    |
|I|S|S|S|  (4)  |A|     (7)     |             (16/64)           |
|N|V|V|V|       |S|             |   (if payload len==126/127)   |
| |1|2|3|       |K|             |                               |
+-+-+-+-+-------+-+-------------+ - - - - - - - - - - - - - - - +
|     Extended payload length continued, if payload len == 127  |
+ - - - - - - - - - - - - - - - +-------------------------------+
|                               |Masking-key, if MASK set to 1  |
+-------------------------------+-------------------------------+
| Masking-key (continued)       |          Payload Data         |
+-------------------------------- - - - - - - - - - - - - - - - +
:                     Payload Data continued ...                :
+ - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - +
|                     Payload Data continued ...                |
+---------------------------------------------------------------+
```

**Don't worry!** Libraries like Gorilla WebSocket handle this for you.

**Key fields:**
- **FIN**: Is this the final frame?
- **Opcode**: Text (1), Binary (2), Close (8), Ping (9), Pong (10)
- **Payload**: Your actual message data

---

## WebSocket Message Types

### Text Messages (Opcode 1)
- UTF-8 encoded strings
- Typically JSON
- What we'll use most

```javascript
// Client sends
{"type": "move", "data": "e2e4"}

// Server responds
{"type": "move_response", "data": "e7e5"}
```

### Binary Messages (Opcode 2)
- Raw bytes
- Images, files, compressed data
- More efficient for non-text

### Control Frames
- **Close (8)**: Graceful shutdown
- **Ping (9)**: "Are you there?"
- **Pong (10)**: "Yes, I'm here!"

---

## When to Use WebSocket

### Perfect For:

✅ **Real-time games** (like chess)
- Instant move updates
- Live opponent actions

✅ **Chat applications**
- Messages appear instantly
- Typing indicators

✅ **Live feeds**
- Stock prices
- Sports scores
- Social media updates

✅ **Collaborative tools**
- Google Docs-style editing
- Shared whiteboards

✅ **IoT and monitoring**
- Sensor data streaming
- Live dashboards

### NOT Good For:

❌ **Simple CRUD operations**
- Creating a user? Use POST
- Reading data once? Use GET

❌ **File uploads**
- HTTP multipart is better

❌ **SEO-dependent pages**
- Search engines need HTTP

❌ **One-way server push**
- Server-Sent Events (SSE) simpler

---

## WebSocket vs Alternatives

### Long Polling
```
Client → Server (wait)
          ↓ (hold connection open)
          ↓ (event happens)
Client ← Server (respond)
Client → Server (new request)
```
- Simpler than WebSocket
- More overhead
- One direction per connection

### Server-Sent Events (SSE)
```
Client → Server (subscribe)
Client ← Server (event)
Client ← Server (event)
Client ← Server (event)
```
- Server to client only
- Automatic reconnection
- Simpler than WebSocket
- Good for notifications

### WebSocket
```
Client ↔ Server (bidirectional)
```
- Full duplex
- Most powerful
- More complex

---

## WebSocket in Your Chess-Coach Project

### Why Your Project Needs WebSocket

1. **Computer Move Generation**
   ```
   User makes move → WebSocket →
   Server calculates → WebSocket →
   Computer move appears instantly
   ```

2. **Multiplayer Games** (future)
   ```
   Player 1 moves → Server →
   Player 2 sees move instantly
   ```

3. **AI Coach** (future)
   ```
   User asks question → WebSocket →
   AI streams response → WebSocket →
   Tokens appear as typed
   ```

### Architecture Preview

```
Frontend                    Backend
   |                           |
   |  WS: /ws?user_id=123      |
   |  -----------------------> |
   |                           |
   |  ← 101 Switching          |
   |═══════════════════════════| (connected)
   |                           |
   |  {type: "game.move"}      |
   |  -----------------------> | → EngineService
   |                           |     (Stockfish)
   |  ← {type: "move.computer"}|
   |                           |
```

---

## Key Concepts Summary

### 1. Persistent Connection
- Opens once, stays open
- Survives multiple messages
- Must handle disconnects

### 2. Bidirectional
- Client can send anytime
- Server can send anytime
- No request-response requirement

### 3. Low Latency
- No connection overhead
- Minimal framing
- Perfect for real-time

### 4. Event-Driven
- Messages are events
- Asynchronous handling
- Concurrent processing

---

## Go WebSocket Libraries

### Gorilla WebSocket (What We'll Use)

```bash
go get github.com/gorilla/websocket
```

**Why Gorilla?**
- ✅ Most popular Go WebSocket library
- ✅ Production-tested
- ✅ Passes all Autobahn test suite
- ✅ Great documentation
- ✅ Used by your chess-coach project

**Alternatives:**
- `nhooyr.io/websocket` - Modern, context-aware
- `gobwas/ws` - Low-level, zero-copy
- Standard library? No native WebSocket support!

---

## What's Next?

In this lesson you learned:
- ✅ Why WebSocket exists
- ✅ HTTP vs WebSocket
- ✅ WebSocket handshake
- ✅ Message types and frames
- ✅ When to use WebSocket
- ✅ How it fits chess-coach

### Next Lesson Preview

**Lesson 2: Your First WebSocket Connection**

We'll write code:
- Create a WebSocket server
- Accept connections
- Send and receive messages
- Build an echo server

---

## Let's Code!

Ready to see WebSocket in action?

Run:
```bash
cd /Users/ankit/code/learn/chess-coach/learnings/fundamentals/websocket/lesson-01-basics
go run main.go
```

Then open `client.html` in your browser!

---

**Questions?** Write them in `QUESTIONS.md` and we'll discuss!

**Next**: When ready, say "Let me run the example" and I'll walk you through it!

# Lesson 1: WebSocket Basics - Summary

**Completed**: 2025-11-02

---

## What I Learned

### 1. Core WebSocket Concepts

#### **Persistent Bidirectional Connection**
- Unlike HTTP (request → response → close), WebSocket opens once and stays open
- Both client and server can send messages anytime without waiting for a request
- Perfect for real-time applications like chess games, chat, live feeds

#### **The WebSocket Handshake**
```
Client → Server: HTTP Upgrade request
Server → Client: 101 Switching Protocols
─────────────────────────────────────────
Now using WebSocket protocol (not HTTP!)
```

After the handshake, it's pure WebSocket frames - minimal overhead, maximum efficiency.

#### **Why WebSocket vs REST?**
**REST/HTTP Polling** (what I would have done without WebSocket):
```
User makes move
  ↓
POST /api/games/123/moves
  ↓
Every 2 seconds: GET /api/games/123/state (is computer move ready?)
  ↓
Repeat polling... wasteful!
```

**WebSocket** (what I'm doing now):
```
User makes move → WebSocket message
  ↓
Computer calculates
  ↓
Computer move → WebSocket message (instant!)
```

**Benefits for Chess-Coach:**
- ✅ Instant move updates (no polling!)
- ✅ Low latency (sub-millisecond)
- ✅ One connection for entire game
- ✅ Server can push updates anytime
- ✅ Efficient bandwidth usage

---

### 2. Text vs Binary Messages

#### **The Critical Difference**

**Text Messages** (Opcode 1):
- MUST be valid UTF-8 encoded strings
- Typically JSON for structured data
- Server/client validates UTF-8 automatically
- **Attempting to send binary data as text = error!**

**Binary Messages** (Opcode 2):
- Raw bytes (no encoding constraints)
- Images, files, compressed data
- No UTF-8 validation
- More efficient for non-text data

#### **The UTF-8 Error I Hit**

**What happened:**
```
I sent an image file → Server tried string(message)
Error: "Could not decode a text frame as UTF-8"
```

**Why it failed:**
Image bytes are not valid UTF-8 text. Attempting to treat binary data as text causes decoding errors.

**The Fix:**
```go
// Check message type BEFORE processing
if messageType == websocket.BinaryMessage {
    // Binary - don't convert to string!
    log.Printf("📨 Received BINARY message: %d bytes", len(message))
    err = conn.WriteMessage(websocket.BinaryMessage, message)
} else {
    // Text - safe to convert to string
    log.Printf("📨 Received: %s", string(message))
    responseMsg := fmt.Sprintf("Echo: %s", string(message))
    err = conn.WriteMessage(websocket.TextMessage, []byte(responseMsg))
}
```

**Key Takeaway:** Always check `messageType` before processing! Text and binary require different handling.

---

### 3. Image Detection and Rendering

#### **Magic Bytes for File Type Detection**

Learned that files have "magic bytes" - signature bytes at the beginning that identify the file type:

```javascript
// PNG: 89 50 4E 47
if (bytes[0] === 0x89 && bytes[1] === 0x50 && bytes[2] === 0x4E && bytes[3] === 0x47) {
    return 'image/png';
}

// JPEG: FF D8 FF
if (bytes[0] === 0xFF && bytes[1] === 0xD8 && bytes[2] === 0xFF) {
    return 'image/jpeg';
}

// And so on for GIF, WebP, BMP...
```

**Why this matters:**
- Can automatically identify image types without relying on file extensions
- More reliable than MIME type headers
- Used by operating systems, browsers, and file tools

#### **Rendering Binary Images**

Implemented automatic image rendering:

1. **Receive binary data** via WebSocket
2. **Detect if it's an image** using magic bytes
3. **Create Object URL** from binary data
4. **Display as `<img>` element**
5. **Clean up URL** after rendering to prevent memory leaks

```javascript
const img = document.createElement('img');
const url = URL.createObjectURL(new Blob([buffer], { type: mimeType }));
img.src = url;
// Clean up after loading
img.onload = () => URL.revokeObjectURL(url);
```

---

### 4. WebSocket Size Limits (Critical for Production!)

#### **Multi-Level Limits**

**Protocol Level:** 2^63 bytes (essentially unlimited)
- But this doesn't mean you SHOULD send huge messages!

**Browser Level:** ~100 MB - 1 GB
- Depends on available RAM
- Browsers may crash before reaching limit

**Network Infrastructure:** 1-16 MB
- Proxies, load balancers, firewalls
- Nginx default: 1 MB
- AWS ALB: 1 MB

**Memory Constraints:** Most Critical!
- WebSocket loads **entire message into RAM** before processing
- 100 connections × 10 MB = 1 GB RAM
- Memory exhaustion = server crash

#### **Security Implications**

**DoS Attack Scenario:**
```javascript
// Malicious client
const ws = new WebSocket('ws://yourserver.com');
ws.onopen = () => {
    ws.send('x'.repeat(100 * 1024 * 1024)); // 100 MB
};
```

Without limits:
```
100 attackers × 100 MB = 10 GB RAM consumed → Server crashes!
```

#### **The Solution: Always Set Limits!**

**Server-side protection:**
```go
const maxMessageSize = 10 * 1024 * 1024 // 10 MB

// Apply to each connection
conn.SetReadLimit(maxMessageSize)
```

**Client-side validation:**
```javascript
if (file.size > MAX_SIZE) {
    alert('File too large!');
    return;
}
```

#### **Recommended Limits by Use Case**

| Use Case | Limit | Rationale |
|----------|-------|-----------|
| Chat | 64 KB | Text messages rarely exceed 1-2 KB |
| JSON API | 1 MB | Complex nested data, still safe |
| Images | 5-10 MB | Compressed photos, quality balance |
| **Chess-Coach** | **1 MB** | Move messages ~200 bytes, huge safety margin |

#### **For Large Files: Use HTTP!**

**Wrong approach:**
```javascript
// ❌ Sending 50 MB video via WebSocket
ws.send(largeVideoFile); // Blocks connection!
```

**Right approach:**
```javascript
// 1. Upload via HTTP (with progress tracking)
const response = await fetch('/api/upload', {
    method: 'POST',
    body: formData
});

// 2. Notify via WebSocket
ws.send(JSON.stringify({
    type: 'file.uploaded',
    fileId: fileId
}));
```

**Benefits:**
- HTTP handles large files better (streaming, resume, caching)
- WebSocket stays responsive
- Memory-efficient
- Can show progress bars

---

### 5. Implementation Details

#### **Buffer Size vs Message Size Limit**

**Important distinction I learned:**

**Buffer Size** (`readBufferSize`, `writeBufferSize`):
```go
readBufferSize = 4096 // 4 KB I/O buffer
```
- How data is read/written in chunks
- **Does NOT limit message size**
- Just affects I/O performance

**Message Size Limit** (`maxMessageSize`):
```go
maxMessageSize = 10 * 1024 * 1024 // 10 MB
conn.SetReadLimit(maxMessageSize)
```
- Limits total message size
- Prevents memory exhaustion
- Rejects messages exceeding limit

Can receive 10 MB message with 4 KB buffer - just reads in multiple chunks!

#### **Gorilla WebSocket Library**

**Why Gorilla is the right choice:**
- ✅ Passes all 521 Autobahn test suite cases (100% RFC 6455 compliant)
- ✅ Production-tested and battle-hardened
- ✅ Handles all edge cases correctly
- ✅ Most popular Go WebSocket library
- ✅ Great documentation

**Autobahn Test Suite** = The "crash test" for WebSocket implementations
- Tests framing, UTF-8 validation, close handshake, compression, etc.
- 521 test cases covering every edge case
- Gorilla passes them all!

---

## Aha! Moments

### 1. **WebSocket is NOT just "faster HTTP"**
Initially thought WebSocket was just a speed optimization. Now I understand it's a fundamentally different communication model:
- HTTP: Request → Response (one-way, then close)
- WebSocket: Persistent bidirectional channel

For chess, this means computer moves can be pushed instantly, not polled.

### 2. **Binary ≠ Text at Protocol Level**
Thought binary was just "data". Learned that WebSocket has distinct frame types with different validation rules:
- Text frames = UTF-8 validated
- Binary frames = raw bytes
- Can't mix them up!

This explained the UTF-8 decode error immediately.

### 3. **Size Limits Are About Security, Not Protocol**
Protocol allows gigabyte messages, but that's dangerous! Size limits protect against:
- Memory exhaustion attacks
- Server crashes
- Excessive cloud costs

**Production rule:** Always set limits based on use case, not protocol maximum.

### 4. **Images via WebSocket Are Possible, But Not Always Best**
Learned that while WebSocket CAN send images, it's not always the best choice:
- Small images (< 1 MB): OK via WebSocket
- Large files: Use HTTP upload + WebSocket notification

Use the right tool for the job!

### 5. **Magic Bytes Are Universal**
File type detection via magic bytes is used everywhere:
- Operating systems
- File command on Unix
- Browser file handlers

This is a fundamental computer science concept I can apply beyond WebSocket.

---

## How This Applies to Chess-Coach

### Current Architecture Decisions

**WebSocket Use Cases:**
1. **Game Moves** - Player move → Computer response (instant!)
2. **Move Validation** - Real-time feedback on legal moves
3. **Game State Updates** - Board updates, captured pieces, check/checkmate

**Why 1 MB Message Limit is Perfect:**
```json
// Typical chess move message (~200 bytes)
{
  "type": "game.move",
  "payload": {
    "move": "e2e4",
    "fen": "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR",
    "evaluation": 15
  }
}
```

1 MB provides 5000x safety margin while protecting against abuse!

### Future Features Enabled

**AI Coach Streaming:**
```json
// AI coach responses (1-10 KB each)
{
  "type": "ai.response",
  "payload": {
    "token": "This move controls the center and opens lines for your bishop...",
    "context": { "move": "e2e4", "principle": "center_control" }
  }
}
```

**Multiplayer Games:**
```
Player 1 makes move → WebSocket → Server
                                     ↓
                          Player 2 sees move instantly (via WebSocket)
```

**What NOT to Use WebSocket For:**
- ❌ Loading game history (use REST GET)
- ❌ Creating new game (use REST POST)
- ❌ User profile updates (use REST PUT)
- ❌ Uploading PGN files (use HTTP multipart)

**Rule:** Use WebSocket for real-time updates, REST for CRUD operations.

---

## Surprising Discoveries

### 1. **WebSocket Works Over HTTP Ports**
Uses ports 80/443, so it works through corporate firewalls! Upgrades from HTTP, so no special port configuration needed.

### 2. **Browsers Handle WebSocket Automatically**
No need to manually handle framing, masking, or control frames. The browser's WebSocket API does it all!

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');
ws.onmessage = (e) => console.log(e.data);
ws.send('Hello!');
```

That's it! Browser handles:
- Frame masking (client → server)
- Fragmentation
- Ping/pong heartbeats
- Close handshake

### 3. **Message Size Affects Cost Dramatically**
```
1000 users × 10 MB messages × 100/day = 1 TB/day = $90/month
1000 users × 100 KB messages × 100/day = 10 GB/day = $0.90/month
```

**100x cost difference** just from message size limits!

Size limits aren't just security - they're cost optimization.

### 4. **Protocol Allows 9 Exabytes, Reality is ~1-10 MB**
The gap between protocol maximum (2^63 bytes) and practical limits (~1-10 MB) is HUGE. Real-world limits come from:
- Memory constraints
- Network infrastructure
- Security considerations
- Cost implications

### 5. **Object URLs Need Cleanup**
Creating Object URLs from binary data without cleanup causes memory leaks:
```javascript
const url = URL.createObjectURL(blob);
// Must clean up!
img.onload = () => URL.revokeObjectURL(url);
```

Small detail, but critical for long-running apps.

---

## Code I Wrote

### 1. **Binary Message Handler (Server)**
```go
// Handle different message types
if messageType == websocket.BinaryMessage {
    // Binary message (images, files, etc.)
    log.Printf("📨 Received BINARY message: %d bytes", len(message))

    // Echo binary data back as-is
    err = conn.WriteMessage(websocket.BinaryMessage, message)
    if err != nil {
        log.Printf("❌ Error sending binary message: %v", err)
        break
    }

    log.Printf("📤 Sent BINARY echo: %d bytes", len(message))
} else {
    // Text message - safe to convert to string
    log.Printf("📨 Received: %s", string(message))
    responseMsg := fmt.Sprintf("Echo: %s (received at %s)",
        string(message), time.Now().Format("15:04:05"))
    err = conn.WriteMessage(websocket.TextMessage, []byte(responseMsg))
}
```

### 2. **Size Limits (Server)**
```go
const (
    // Maximum message size allowed (10 MB)
    maxMessageSize = 10 * 1024 * 1024

    // Buffer sizes for reading/writing
    readBufferSize  = 4096 // 4 KB
    writeBufferSize = 4096 // 4 KB
)

// Apply limit to connection
conn.SetReadLimit(maxMessageSize)
```

### 3. **Magic Bytes Image Detection (Client)**
```javascript
function detectImageType(bytes) {
    if (bytes.length < 4) return null;

    // PNG: 89 50 4E 47
    if (bytes[0] === 0x89 && bytes[1] === 0x50 &&
        bytes[2] === 0x4E && bytes[3] === 0x47) {
        return 'image/png';
    }

    // JPEG: FF D8 FF
    if (bytes[0] === 0xFF && bytes[1] === 0xD8 && bytes[2] === 0xFF) {
        return 'image/jpeg';
    }

    // GIF, WebP, BMP...
    // (additional checks)

    return null;
}
```

### 4. **Image Rendering (Client)**
```javascript
function handleBinaryMessage(blob) {
    blob.arrayBuffer().then(buffer => {
        const bytes = new Uint8Array(buffer);
        const mimeType = detectImageType(bytes);

        if (mimeType) {
            // It's an image! Display it
            const img = document.createElement('img');
            const url = URL.createObjectURL(new Blob([buffer], { type: mimeType }));
            img.src = url;
            img.style.maxWidth = '400px';

            // Add to message area
            messagesDiv.lastElementChild.appendChild(img);

            // Clean up
            img.onload = () => URL.revokeObjectURL(url);
        } else {
            // Not an image, display as binary data
            addMessage(`[BINARY] ${bytes.length} bytes`, 'received');
        }
    });
}
```

### 5. **File Upload (Client)**
```javascript
function sendImageFile() {
    const reader = new FileReader();

    reader.onload = (event) => {
        // Send as binary ArrayBuffer
        const arrayBuffer = event.target.result;
        ws.send(arrayBuffer);

        const sizeMB = (selectedFile.size / 1024 / 1024).toFixed(2);
        addMessage(`[IMAGE] Sent: ${selectedFile.name} (${sizeMB} MB)`, 'sent');
    };

    // Read file as ArrayBuffer (binary data)
    reader.readAsArrayBuffer(selectedFile);
}
```

---

## Questions That Came Up

### Answered Questions (See QUESTIONS.md):
1. ✅ What are other similar technologies like WebSocket? (XMPP, MQTT, SSE, WebRTC, etc.)
2. ✅ Why not use GraphQL Subscriptions? (Overkill, adds complexity, we have REST already)
3. ✅ What's the difference between GraphQL and REST? (Query language vs resource-based)
4. ✅ What does "Passes all Autobahn test suite" mean? (100% RFC 6455 compliance, production-ready)
5. ✅ What are WebSocket message size limits? (Multi-level: protocol, browser, network, memory)

### Future Questions:
- How do I handle connection drops and reconnection?
- What are ping/pong frames and how do they work?
- How do I manage multiple WebSocket connections efficiently?
- What about compression (per-message deflate)?
- How do I secure WebSocket with authentication?
- Should I use one WebSocket for everything or multiple connections?

These will be covered in upcoming lessons!

---

## Key Takeaways for Production

### 1. **Always Set Message Size Limits**
```go
conn.SetReadLimit(maxMessageSize)
```
Protects against DoS attacks, memory exhaustion, and excessive costs.

### 2. **Check Message Type Before Processing**
```go
if messageType == websocket.BinaryMessage {
    // Handle binary
} else {
    // Handle text
}
```
Prevents UTF-8 decode errors and data corruption.

### 3. **Use HTTP for Large Files**
```
Large file → HTTP upload with progress
Small notification → WebSocket
```
Each protocol for its strength!

### 4. **Validate Data on Both Client and Server**
```javascript
// Client
if (file.size > MAX_SIZE) return;

// Server
conn.SetReadLimit(maxMessageSize)
```
Defense in depth - never trust just one layer.

### 5. **Clean Up Resources**
```javascript
URL.revokeObjectURL(url) // Clean up Object URLs
conn.Close()             // Close connections
defer conn.Close()       // Use defer in Go
```
Prevents memory leaks in long-running applications.

### 6. **Monitor in Production**
```go
log.Printf("Connections: %d, Memory: %d MB", connCount, memUsage)
```
Track connection count, memory usage, message rates.

---

## What I'll Apply to Chess-Coach

### Immediate Integration:
1. **WebSocket endpoint** at `/ws` for game moves
2. **Message size limit** of 1 MB (perfect for chess moves)
3. **JSON message format** for structured data:
   ```json
   {
     "type": "game.move",
     "payload": { "move": "e2e4", "fen": "..." }
   }
   ```
4. **Connection management** with proper error handling
5. **Binary message support** for future features (board images, etc.)

### Architecture Pattern:
```
Frontend                Backend
   |                       |
   |  WebSocket: /ws       |
   |  ------------------>  |
   |  {type: "game.move"}  |
   |                       | → EngineService
   |                       |   (Stockfish)
   |  <------------------  |
   |  {type: "computer.move"}
```

### REST + WebSocket Hybrid:
- **REST**: CRUD operations (create game, get history, update profile)
- **WebSocket**: Real-time updates (moves, state changes, AI streaming)

Best of both worlds!

---

## Next Steps

### Ready for Lesson 2: Connection Management
- Connection lifecycle
- Error handling and recovery
- Ping/pong heartbeats
- Graceful disconnection
- Reconnection strategies

### Ready to Integrate:
I now have enough knowledge to:
- Add WebSocket endpoint to chess-coach backend
- Handle game moves via WebSocket
- Implement real-time move updates
- Set appropriate size limits and security measures

### Confident About:
- ✅ When to use WebSocket vs REST
- ✅ Text vs binary message handling
- ✅ Size limits and security
- ✅ File type detection
- ✅ Real-world production considerations

### Still Learning:
- Connection reliability (Lesson 2)
- Scaling WebSocket servers (Lesson 3+)
- Advanced patterns (pub/sub, rooms, etc.)

---

## Reflection

### What Worked Well:
- **Hands-on experiments** - Actually hitting errors (UTF-8 decode) taught me more than just reading
- **Real examples** - Using images showed practical binary handling
- **Production focus** - Learning size limits and security from the start
- **Progressive complexity** - Started simple (echo server), added features incrementally

### What Surprised Me:
- How simple the basic WebSocket code is (just a few lines!)
- How many production considerations there are (limits, security, costs)
- The huge gap between protocol capabilities and real-world constraints
- How WebSocket complements REST rather than replacing it

### What I'd Do Differently:
- Could have experimented with chunking large data
- Could have tried breaking the server intentionally (DoS testing)
- Could have implemented basic authentication

### Confidence Level:
**Before Lesson 1:** 2/10 - No idea what WebSocket was
**After Lesson 1:** 7/10 - Can implement basic WebSocket for chess-coach

Still need to learn connection management, scaling, and advanced patterns, but I'm confident in the fundamentals!

---

## Files Modified

1. **main.go** - Added binary message handling and size limits
2. **client.html** - Added image upload, binary detection, and rendering
3. **QUESTIONS.md** - Documented size limits and answered questions

---

## Resources Used

1. **Gorilla WebSocket Documentation** - https://pkg.go.dev/github.com/gorilla/websocket
2. **RFC 6455** - WebSocket Protocol Specification
3. **Autobahn Test Suite** - https://github.com/crossbario/autobahn-testsuite
4. **MDN WebSocket API** - https://developer.mozilla.org/en-US/docs/Web/API/WebSocket

---

**Lesson 1 Status:** ✅ Completed

**Ready for:** Lesson 2 - Connection Management & Reliability

**Overall Feeling:** Excited to integrate this into chess-coach! WebSocket is simpler than I thought but has important production considerations. The hands-on experiments made everything click. 🎯

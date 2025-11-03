# Lesson 2: Connection Management & Reliability

**Goal**: Learn how to handle WebSocket connection lifecycle, detect failures, and maintain reliable connections

**Time**: 45-60 minutes

**Prerequisites**: Completed Lesson 1 (WebSocket Basics)

---

## What You'll Learn

By the end of this lesson, you'll understand:

1. **Connection Lifecycle**
   - Connection states (CONNECTING, OPEN, CLOSING, CLOSED)
   - Graceful vs ungraceful disconnection
   - Connection cleanup

2. **Heartbeats (Ping/Pong)**
   - Why connections need heartbeats
   - How ping/pong frames work
   - Implementing timeouts

3. **Error Detection**
   - Network failures
   - Server crashes
   - Client crashes
   - How to detect "zombie" connections

4. **Connection Recovery**
   - Detecting disconnections
   - Automatic reconnection strategies
   - Exponential backoff
   - Message queuing during disconnection

5. **Production Patterns**
   - Read/write deadlines
   - Graceful shutdown
   - Connection limits
   - Health monitoring

---

## The Problem: Connections Don't Last Forever

### Scenario 1: The Silent Disconnect

```
Client                          Server
  |═══════════════════════════════| Connected
  |                               |
  | (User closes laptop)          |
  |                               |
  |                               | Server thinks client is still connected!
  |                               | Sends messages to void...
  |                               | Resources wasted!
```

**Problem:** TCP doesn't always detect disconnections immediately. Server may hold resources for a "zombie" connection.

---

### Scenario 2: The Network Blip

```
Client                          Server
  |═══════════════════════════════| Connected
  |                               |
  | (WiFi disconnects briefly)    |
  |  X  X  X  X  X  X  X  X  X    |
  |                               |
  | (WiFi reconnects)             |
  |                               | Connection is DEAD but both sides don't know yet!
```

**Problem:** Network can fail temporarily. Connection appears open but is actually broken.

---

### Scenario 3: The Server Restart

```
Client                          Server
  |═══════════════════════════════|
  |                               |
  |                               | (Server restarts)
  |                               | X
  |                               |
  | Trying to send...             |
  | Error!                        |
  | What now?                     | (Server back up)
```

**Problem:** Client needs to detect failure and reconnect automatically.

---

## Solution 1: Heartbeats (Ping/Pong)

### How Ping/Pong Works

WebSocket has built-in **control frames** for checking connection health:

```
Client                          Server
  |                               |
  |  PING  -------------------->  |
  |                               | (Connection alive!)
  |  <-------------------- PONG   |
  |                               |
  |  (Wait 30 seconds)            |
  |                               |
  |  PING  -------------------->  |
  |                               |
  |  <-------------------- PONG   |
```

**If PONG doesn't arrive:** Connection is dead!

---

### Why You Need Heartbeats

#### 1. **Detect Zombie Connections**

Server can detect if client is actually alive:

```go
// Server sends ping every 30 seconds
// If no pong received within 60 seconds → close connection
```

#### 2. **Keep Connection Alive**

Some proxies/firewalls close "idle" connections:

```
Proxy sees no traffic for 60 seconds → Closes connection
Solution: Send ping every 30 seconds to keep connection "active"
```

#### 3. **Faster Failure Detection**

Without pings:
```
Network failure → Wait minutes for TCP timeout
With pings:
Network failure → No pong within 10 seconds → Immediate detection!
```

---

## WebSocket Connection States

The browser WebSocket API has 4 states:

```javascript
ws.readyState === WebSocket.CONNECTING  // 0 - Handshake in progress
ws.readyState === WebSocket.OPEN        // 1 - Connection established
ws.readyState === WebSocket.CLOSING     // 2 - Close handshake started
ws.readyState === WebSocket.CLOSED      // 3 - Connection closed
```

### State Transitions

```
CONNECTING  →  OPEN  →  CLOSING  →  CLOSED
    ↓                      ↓           ↓
    └─────→ CLOSED  ←──────┴───────────┘
```

### Checking Before Sending

```javascript
if (ws.readyState === WebSocket.OPEN) {
    ws.send(message); // Safe to send
} else {
    console.log('Cannot send - connection not open');
    // Queue message or reconnect
}
```

---

## Connection Lifecycle Events

### Client-Side Events

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => {
    console.log('✅ Connected');
    // Connection is ready - can send messages
};

ws.onmessage = (event) => {
    console.log('📨 Received:', event.data);
};

ws.onerror = (error) => {
    console.error('❌ Error:', error);
    // Something went wrong during connection or communication
};

ws.onclose = (event) => {
    console.log('🔌 Disconnected');
    console.log('Code:', event.code);
    console.log('Reason:', event.reason);
    console.log('Clean:', event.wasClean);

    // Decide: Reconnect? Show error?
};
```

---

### Close Event Codes

WebSocket defines standard close codes:

| Code | Name | Meaning |
|------|------|---------|
| 1000 | Normal Closure | Graceful disconnect |
| 1001 | Going Away | Browser tab closed, page navigated away |
| 1002 | Protocol Error | WebSocket protocol violation |
| 1003 | Unsupported Data | Received data type it can't handle |
| 1006 | Abnormal Closure | No close frame received (network failure) |
| 1008 | Policy Violation | Message violates policy (e.g., too large) |
| 1009 | Message Too Big | Message exceeded size limit |
| 1011 | Internal Error | Server encountered unexpected condition |

### Example: Detecting Abnormal Closure

```javascript
ws.onclose = (event) => {
    if (event.code === 1006) {
        console.log('⚠️ Connection lost (network issue)');
        // Attempt reconnection
        reconnect();
    } else if (event.code === 1000) {
        console.log('✅ Normal disconnect');
        // Don't reconnect
    } else {
        console.log(`❌ Closed with code ${event.code}: ${event.reason}`);
    }
};
```

---

## Server-Side Connection Management

### Setting Timeouts (Go/Gorilla)

```go
const (
    writeWait = 10 * time.Second    // Time to write message to client
    pongWait  = 60 * time.Second    // Time to wait for pong
    pingPeriod = (pongWait * 9) / 10 // Send pings at this interval (54 seconds)
)

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, _ := upgrader.Upgrade(w, r, nil)
    defer conn.Close()

    // Set read deadline - if no pong received within pongWait, read fails
    conn.SetReadDeadline(time.Now().Add(pongWait))

    // When pong received, reset deadline
    conn.SetPongHandler(func(string) error {
        conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    // Start ping ticker
    go func() {
        ticker := time.NewTicker(pingPeriod)
        defer ticker.Stop()

        for {
            <-ticker.C
            conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return // Connection dead
            }
        }
    }()

    // Read messages (blocks until message or timeout)
    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
                log.Printf("Error: %v", err)
            }
            break
        }

        // Process message...
    }
}
```

---

## Client-Side Automatic Reconnection

### Basic Reconnection Pattern

```javascript
class ReconnectingWebSocket {
    constructor(url) {
        this.url = url;
        this.ws = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 10;
        this.reconnectDelay = 1000; // Start with 1 second

        this.connect();
    }

    connect() {
        console.log('🔌 Connecting to', this.url);

        this.ws = new WebSocket(this.url);

        this.ws.onopen = () => {
            console.log('✅ Connected');
            this.reconnectAttempts = 0; // Reset on success
            this.reconnectDelay = 1000;
        };

        this.ws.onmessage = (event) => {
            // Handle message
            console.log('📨', event.data);
        };

        this.ws.onerror = (error) => {
            console.error('❌ Error:', error);
        };

        this.ws.onclose = (event) => {
            console.log('🔌 Disconnected');

            // Don't reconnect if normal closure
            if (event.code === 1000) {
                console.log('Normal close - not reconnecting');
                return;
            }

            // Attempt reconnection
            this.reconnect();
        };
    }

    reconnect() {
        if (this.reconnectAttempts >= this.maxReconnectAttempts) {
            console.error('❌ Max reconnection attempts reached');
            return;
        }

        this.reconnectAttempts++;

        console.log(`⏳ Reconnecting in ${this.reconnectDelay}ms (attempt ${this.reconnectAttempts})`);

        setTimeout(() => {
            this.connect();
        }, this.reconnectDelay);

        // Exponential backoff: 1s, 2s, 4s, 8s, 16s, 30s (max)
        this.reconnectDelay = Math.min(this.reconnectDelay * 2, 30000);
    }

    send(message) {
        if (this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(message);
        } else {
            console.warn('⚠️ Cannot send - not connected');
            // Could queue message here
        }
    }

    close() {
        this.ws.close(1000, 'Normal closure');
    }
}

// Usage
const ws = new ReconnectingWebSocket('ws://localhost:8080/ws');
ws.send('Hello!');
```

---

## Exponential Backoff

When reconnecting, use **exponential backoff** to avoid overwhelming server:

```
Attempt 1: Wait 1 second
Attempt 2: Wait 2 seconds
Attempt 3: Wait 4 seconds
Attempt 4: Wait 8 seconds
Attempt 5: Wait 16 seconds
Attempt 6+: Wait 30 seconds (max)
```

### Why Exponential Backoff?

**Without backoff:**
```
Server crashes
↓
1000 clients try to reconnect immediately every second
↓
Server comes back up
↓
1000 simultaneous reconnections
↓
Server crashes again from load!
```

**With exponential backoff:**
```
Server crashes
↓
Clients reconnect at different times: 1s, 2s, 4s, 8s...
↓
Server comes back up
↓
Connections spread over time
↓
Server handles load gracefully ✅
```

---

## Message Queuing During Disconnection

When disconnected, you can queue messages:

```javascript
class ReconnectingWebSocket {
    constructor(url) {
        this.url = url;
        this.ws = null;
        this.messageQueue = []; // Queue messages during disconnect

        this.connect();
    }

    send(message) {
        if (this.ws.readyState === WebSocket.OPEN) {
            // Connected - send immediately
            this.ws.send(message);
        } else {
            // Disconnected - queue for later
            console.log('⏳ Queuing message (not connected)');
            this.messageQueue.push(message);
        }
    }

    onConnectionEstablished() {
        // Connected! Send queued messages
        console.log(`📤 Sending ${this.messageQueue.length} queued messages`);

        while (this.messageQueue.length > 0) {
            const message = this.messageQueue.shift();
            this.ws.send(message);
        }
    }
}
```

**Important:** Only queue critical messages! Don't queue indefinitely - set limits.

---

## Graceful Shutdown

### Server-Side Graceful Shutdown

When server needs to restart, close connections gracefully:

```go
func main() {
    // Catch shutdown signals
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

    server := &http.Server{Addr: ":8080"}

    go func() {
        <-sigChan
        log.Println("🛑 Shutting down gracefully...")

        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        // Close all WebSocket connections with reason
        // (You'd track connections in a map/registry)
        for _, conn := range activeConnections {
            conn.WriteMessage(websocket.CloseMessage,
                websocket.FormatCloseMessage(1001, "Server restarting"))
        }

        server.Shutdown(ctx)
    }()

    log.Fatal(server.ListenAndServe())
}
```

### Client Receives Graceful Close

```javascript
ws.onclose = (event) => {
    if (event.code === 1001 && event.reason === "Server restarting") {
        console.log('⏳ Server restarting - will reconnect...');
        // Wait a bit longer before reconnecting
        setTimeout(() => reconnect(), 5000);
    }
};
```

---

## Connection Registry Pattern

Track all active connections on server:

```go
type ConnectionRegistry struct {
    connections map[*websocket.Conn]bool
    mu          sync.RWMutex
}

func (r *ConnectionRegistry) Register(conn *websocket.Conn) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.connections[conn] = true
}

func (r *ConnectionRegistry) Unregister(conn *websocket.Conn) {
    r.mu.Lock()
    defer r.mu.Unlock()
    delete(r.connections, conn)
}

func (r *ConnectionRegistry) Count() int {
    r.mu.RLock()
    defer r.mu.RUnlock()
    return len(r.connections)
}

func (r *ConnectionRegistry) Broadcast(message []byte) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    for conn := range r.connections {
        conn.WriteMessage(websocket.TextMessage, message)
    }
}
```

---

## For Your Chess-Coach Project

### Connection Patterns You'll Need

#### 1. **Heartbeats**
```go
// Server sends ping every 30 seconds
// Closes connection if no pong within 60 seconds
```

Why? Detect when user closed browser tab without proper disconnect.

#### 2. **Client Reconnection**
```javascript
// Reconnect automatically if connection drops
// Rejoin game with game ID
```

Why? Network blips shouldn't interrupt chess game!

#### 3. **Message Queuing**
```javascript
// If disconnected mid-move, queue the move
// Send when reconnected
```

Why? Don't lose player's move!

#### 4. **Connection Limits**
```go
// Max 10,000 concurrent connections
// Reject new connections if limit reached
```

Why? Prevent resource exhaustion.

---

## Key Concepts Summary

### 1. **Heartbeats Are Essential**
- Detect zombie connections
- Keep connections alive through proxies
- Fast failure detection

### 2. **Always Implement Reconnection**
- Use exponential backoff
- Don't overwhelm server
- Queue critical messages

### 3. **Handle All Close Codes**
- 1000 (normal) - Don't reconnect
- 1006 (abnormal) - Network issue, reconnect
- 1001 (going away) - User navigated, reconnect only if needed

### 4. **Set Deadlines**
- Read deadline - Detect stalled reads
- Write deadline - Detect stalled writes
- Pong deadline - Detect dead connections

### 5. **Graceful Shutdown**
- Send close frames with reason
- Give connections time to close
- Track active connections

---

## What's Next?

In this lesson, you'll implement:

1. **Server-side ping/pong** with timeouts
2. **Client-side reconnection** with exponential backoff
3. **Connection state management**
4. **Graceful disconnection handling**
5. **Connection monitoring** (count, health)

Ready to code? Let's build a reliable WebSocket server!

---

**Say "Let's start coding"** when you're ready to begin implementation!

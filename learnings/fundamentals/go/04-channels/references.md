# Channel `<-` Operator References in Chess Coach

Find these patterns in your actual codebase!

---

## Pattern 1: Hub Registration Channel

### Where to Find
Look for Hub struct definition and the `Run()` method.

### Typical Code
```go
type Hub struct {
    register   chan *Client
    unregister chan *Client
    clients    map[*Client]bool
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:  // ← RECEIVE from register channel
            h.clients[client] = true
            fmt.Printf("Client registered\n")
        }
    }
}
```

### What's Happening
- `h.register` is a channel of `*Client` pointers
- `<-h.register` **receives** a client from the channel
- This blocks until someone sends a client
- When a client connects, they send themselves: `hub.register <- client`

### Exercise
**Find this in your code:**
1. Open your websocket hub file
2. Find the `register chan *Client` declaration
3. Find where it's received: `client := <-h.register`
4. Find where clients send themselves to it

---

## Pattern 2: Hub Broadcast Channel

### Typical Code
```go
type Hub struct {
    broadcast chan []byte
}

func (h *Hub) Run() {
    for {
        select {
        case message := <-h.broadcast:  // ← RECEIVE broadcast message
            // Send to all clients
            for client := range h.clients {
                client.send <- message  // ← SEND to each client
            }
        }
    }
}
```

### What's Happening
- `<-h.broadcast` receives a message to broadcast
- Loop through all clients
- `client.send <- message` sends TO each client's queue
- Notice: **receive** from broadcast, **send** to clients

### Arrow Directions
```
Game Logic                Hub                   Clients
    │                      │                      │
    │  msg ─────────>      │                      │
    │  (broadcast <- msg)  │                      │
    │                      │   msg ───────────>   │
    │                      │  (client.send <- msg)│
```

---

## Pattern 3: Client Send Queue

### Typical Code
```go
type Client struct {
    send chan []byte  // Buffered channel (queue)
}

func (c *Client) writePump() {
    for {
        select {
        case message, ok := <-c.send:  // ← RECEIVE from queue
            if !ok {
                // Channel closed by hub
                return
            }
            // Send to actual websocket
            c.conn.WriteMessage(websocket.TextMessage, message)
        }
    }
}
```

### What's Happening
- `c.send` is **buffered** (e.g., size 256)
- Hub puts messages in: `client.send <- message`
- writePump takes messages out: `message := <-c.send`
- Acts as a queue between hub and websocket

### Why Buffered?
If hub tries to send but writePump is busy:
- **Unbuffered**: Hub blocks (waits)
- **Buffered**: Message goes in queue, hub continues

### The `ok` Pattern
```go
message, ok := <-c.send
if !ok {
    // Hub closed the channel
    // Client is being disconnected
    return
}
```

When hub unregisters a client:
```go
close(client.send)  // Close the channel
```

Next time writePump receives: `ok = false`

---

## Pattern 4: Ticker for Heartbeat

### Typical Code
```go
func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer ticker.Stop()

    for {
        select {
        case message := <-c.send:
            // Send message
            c.conn.WriteMessage(websocket.TextMessage, message)

        case <-ticker.C:  // ← RECEIVE from ticker (discard value)
            // Time to send ping (heartbeat)
            c.conn.WriteMessage(websocket.PingMessage, nil)
        }
    }
}
```

### What's Happening
- `time.NewTicker()` creates a channel that sends periodically
- `<-ticker.C` receives from ticker (but we don't store the value)
- This fires every X seconds (e.g., 54 seconds)
- Keeps connection alive with ping messages

### Why Discard?
```go
<-ticker.C  // We don't need the value, just the signal
```

We only care WHEN it fires, not what it sends.

---

## Pattern 5: Creating Client with Buffered Channel

### Typical Code
```go
func NewClient(hub *Hub, conn *websocket.Conn) *Client {
    return &Client{
        hub:  hub,
        conn: conn,
        send: make(chan []byte, 256),  // ← Create buffered channel
    }
}
```

### What's Happening
- `make(chan []byte, 256)` creates a buffered channel
- Can hold 256 messages before blocking
- Each client gets their own send queue

---

## Pattern 6: Hub Broadcast with Select

### Typical Code
```go
for client := range h.clients {
    select {
    case client.send <- message:  // ← Try to send
        // Sent successfully
    default:
        // Client buffer full (can't send without blocking)
        // Close and remove client
        close(client.send)
        delete(h.clients, client)
    }
}
```

### What's Happening
- Try to send: `client.send <- message`
- If buffer is full, `default` case executes
- This prevents hub from blocking on slow clients
- Slow clients get disconnected

### Select with Default
```go
select {
case ch <- value:
    // Sent successfully
default:
    // Would block, do something else
}
```

**Without default**: Blocks until send succeeds
**With default**: Never blocks, runs default if can't send

---

## Complete Flow Example

### When a move is made:

```go
// 1. Game logic creates move message
move := []byte(`{"type":"move","data":"e2e4"}`)

// 2. Send to hub's broadcast channel
hub.broadcast <- move  // ← SEND to hub

// 3. Hub receives and distributes
message := <-h.broadcast  // ← RECEIVE in hub
for client := range h.clients {
    client.send <- message  // ← SEND to each client queue
}

// 4. Each client's writePump receives
msg := <-c.send  // ← RECEIVE from queue

// 5. Send to actual websocket
c.conn.WriteMessage(websocket.TextMessage, msg)
```

### Arrow Flow
```
Game Logic → hub.broadcast → Hub → client.send → writePump → WebSocket
            (send)         (recv)  (send)       (recv)
```

---

## Channel Creation Patterns

### Unbuffered (Synchronous)
```go
register := make(chan *Client)
```

**When to use**: Registration/unregistration (rare events)

### Buffered (Asynchronous)
```go
send := make(chan []byte, 256)
broadcast := make(chan []byte, 10)
```

**When to use**: Message queues (frequent events)

---

## Exercise: Find in Your Code

### Task 1: Find All Channel Declarations
Look for:
```go
chan *Client
chan []byte
chan bool
```

### Task 2: Find All Receives
Look for:
```go
x := <-channel
x, ok := <-channel
<-channel
```

### Task 3: Find All Sends
Look for:
```go
channel <- value
```

### Task 4: Map the Flow
Draw the flow of a message from:
1. Game logic
2. Through hub
3. To all clients
4. To websocket

Mark each `<-` as send or receive!

---

## Common Patterns Summary

| Pattern | Example | Use Case |
|---------|---------|----------|
| Receive with ok | `msg, ok := <-ch` | Detect closure |
| Receive discard | `<-ticker.C` | Timing/signals |
| Send to queue | `client.send <- msg` | Buffer messages |
| Buffered creation | `make(chan T, N)` | Message queues |
| Close channel | `close(client.send)` | Signal shutdown |
| Range over channel | `for msg := range ch` | Process until close |

---

## Questions to Answer

After exploring your code:

1. **How many channels does the Hub have?**


2. **Which channels are buffered? What are their sizes?**


3. **Where does the hub RECEIVE from channels?**


4. **Where does the hub SEND to channels?**


5. **What happens when a client's send buffer is full?**


---

## Next Lesson Connection

**Lesson 5: Select Statement**

You've seen `select` used with `<-` in the hub:
```go
select {
case client := <-h.register:  // ← You understand this now!
case message := <-h.broadcast: // ← And this!
}
```

Next lesson: How does `select` CHOOSE which case to execute?

**Hint**: Whichever channel receives data first!

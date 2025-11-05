# Understanding Hub and Pump Patterns

**What**: Core architectural patterns for WebSocket servers
**Why**: Essential for managing multiple connections efficiently
**When**: Use these in any production WebSocket application

---

## Table of Contents

1. [The Problem We're Solving](#the-problem)
2. [The Pump Pattern](#pump-pattern)
3. [The Hub Pattern](#hub-pattern)
4. [How They Work Together](#working-together)
5. [Real-World Analogies](#analogies)
6. [Code Examples](#code-examples)
7. [Common Questions](#questions)

---

<a name="the-problem"></a>
## 1. The Problem We're Solving

### Naive Approach (Don't Do This!)

```go
// ❌ BAD: Blocking, can't send and receive simultaneously
func handleConnection(conn *websocket.Conn) {
    for {
        // Read message (BLOCKS here)
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            break
        }

        // Process message
        response := "Echo: " + string(message)

        // Write response
        conn.WriteMessage(messageType, []byte(response))
    }
}
```

**Problems**:
1. ❌ **Blocking**: Can't receive new messages while writing
2. ❌ **No pings**: Can't send heartbeat pings on schedule
3. ❌ **No concurrency**: One operation at a time
4. ❌ **No connection management**: Can't track active connections
5. ❌ **No broadcasting**: Can't send to multiple clients

---

<a name="pump-pattern"></a>
## 2. The Pump Pattern

**Concept**: Separate goroutines for reading and writing

Think of it like **two one-way streets** instead of a single two-way street that causes traffic jams.

### Visual Representation

```
┌─────────────────────────────────────────────────────────┐
│                    WebSocket Connection                 │
│                                                         │
│                    ┌──────────────┐                     │
│                    │   Socket     │                     │
│                    └──────┬───────┘                     │
│                           │                             │
│                    ┌──────┴──────┐                      │
│                    │             │                      │
│            ┌───────▼────┐   ┌────▼───────┐              │
│            │ readPump   │   │ writePump  │              │
│            │ goroutine  │   │ goroutine  │              │
│            │            │   │            │              │
│            │ Reads      │   │ Writes     │              │
│            │ messages   │   │ messages   │              │
│            │ from       │   │ to         │              │
│            │ client     │   │ client     │              │
│            └────────────┘   └────────────┘              │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### Why Two Separate Goroutines?

#### Without Separation (Sequential)
```go
// ❌ BAD: Can't do both at once
for {
    // Waiting for message...
    message := readMessage()  // BLOCKS HERE

    // Can't send ping while waiting!
    // Can't send other messages while waiting!

    sendResponse(message)
}
```

#### With Separation (Concurrent)
```go
// ✅ GOOD: Both happen simultaneously

// Goroutine 1: Always ready to read
go func() {
    for {
        message := readMessage()  // Blocking, but only for reading
        // Process message
    }
}()

// Goroutine 2: Always ready to write
go func() {
    for {
        select {
        case message := <-sendChannel:
            writeMessage(message)  // Write whenever ready
        case <-ticker.C:
            sendPing()  // Can send pings on schedule!
        }
    }
}()
```

---

### readPump() - The Reading Goroutine

**Job**: Continuously read messages from the WebSocket connection

```go
func (c *Client) readPump() {
    defer func() {
        // Cleanup when done
        c.hub.unregister <- c
        c.conn.Close()
    }()

    // Set up connection parameters
    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    // Read loop - runs forever until error
    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            // Connection closed or error
            break
        }

        // Process message (send to hub, etc.)
        c.hub.handleMessage(c, message)
    }
}
```

**Key Points**:
- Runs in its own goroutine: `go client.readPump()`
- **Blocks** on `ReadMessage()` - but that's OK! It's the only thing it does
- Handles pongs automatically
- Sets deadlines for timeout detection
- Runs until connection closes

---

### writePump() - The Writing Goroutine

**Job**: Write messages to the WebSocket AND send periodic pings

```go
func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
            // Got a message to send from the send channel
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                // Channel closed
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            // Write the message
            if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
                return
            }

        case <-ticker.C:
            // Time to send a ping!
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}
```

**Key Points**:
- Runs in its own goroutine: `go client.writePump()`
- Uses `select` to handle **multiple input sources**:
  1. Messages from `send` channel
  2. Ticker for periodic pings
- Never blocks indefinitely - always responsive
- Sets write deadlines for timeout detection

---

### Communication Between Pumps

**How do they coordinate?** Through a **channel**!

```go
type Client struct {
    conn *websocket.Conn
    send chan []byte  // ← This is the bridge!
}

// readPump receives message and wants to send response
func (c *Client) readPump() {
    for {
        _, message, _ := c.conn.ReadMessage()

        // Put response in send channel
        response := processMessage(message)
        c.send <- response  // writePump will pick this up!
    }
}

// writePump picks up messages from send channel
func (c *Client) writePump() {
    for {
        select {
        case message := <-c.send:  // ← Receives from readPump
            c.conn.WriteMessage(websocket.TextMessage, message)
        }
    }
}
```

**The Beauty**: readPump and writePump are completely independent but communicate safely via channel!

---

<a name="hub-pattern"></a>
## 3. The Hub Pattern

**Concept**: Central coordinator for all WebSocket connections

Think of it like a **switchboard operator** connecting multiple phone calls.

### Without Hub (Chaos!)

```go
// ❌ BAD: How do you manage multiple clients?
var clients []*Client  // Not thread-safe!

// How do you broadcast to all?
for _, client := range clients {  // Race condition!
    client.send <- message
}

// How do you remove disconnected clients?
// How do you handle concurrent access?
// MESSY!
```

### With Hub (Organized!)

```go
// ✅ GOOD: Hub manages everything
type Hub struct {
    clients    map[*Client]bool
    register   chan *Client      // New clients join here
    unregister chan *Client      // Disconnected clients leave here
    broadcast  chan []byte       // Messages to broadcast
}
```

---

### Hub Visual

```
                        ┌─────────────────────┐
                        │        HUB          │
                        │  (Goroutine Loop)   │
                        │                     │
    register ────────>  │  • clients map      │
                        │  • register chan    │
    unregister ──────>  │  • unregister chan  │
                        │  • broadcast chan   │
    broadcast ───────>  │                     │
                        └──────────┬──────────┘
                                   │
                         ┌─────────┴─────────┐
                         │                   │
                    ┌────▼────┐          ┌───▼─────┐
                    │ Client1 │          │ Client2 │
                    │         │          │         │
                    │ •send   │          │ •send   │
                    └─────────┘          └─────────┘
```

---

### Hub Implementation

```go
type Hub struct {
    // Active clients
    clients map[*Client]bool

    // Channels for managing clients
    register   chan *Client
    unregister chan *Client

    // Channel for broadcasting messages
    broadcast chan []byte

    // Mutex for thread-safe access
    mu sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        broadcast:  make(chan []byte, 256),
    }
}
```

---

### Hub Event Loop

This is the **heart** of the hub - a goroutine that runs forever:

```go
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            // New client connected
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()
            log.Printf("Client registered. Total: %d", len(h.clients))

        case client := <-h.unregister:
            // Client disconnected
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)  // Close their send channel
            }
            h.mu.Unlock()
            log.Printf("Client unregistered. Total: %d", len(h.clients))

        case message := <-h.broadcast:
            // Broadcast message to all clients
            h.mu.RLock()
            for client := range h.clients {
                select {
                case client.send <- message:
                    // Message queued successfully
                default:
                    // Client's send buffer full, close it
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mu.RUnlock()
        }
    }
}
```

**Key Points**:
- Runs in its own goroutine: `go hub.Run()`
- Uses `select` to handle multiple event types
- **Never blocks** - always responsive
- Thread-safe with mutex locks
- Single point of control for all connections

---

### Why Use Channels Instead of Direct Function Calls?

#### Without Channels (Dangerous!)
```go
// ❌ BAD: Race conditions!
func (h *Hub) Register(client *Client) {
    h.clients[client] = true  // Multiple goroutines writing = CRASH!
}

// From different goroutines:
hub.Register(client1)  // Goroutine 1
hub.Register(client2)  // Goroutine 2 (concurrent access!)
```

#### With Channels (Safe!)
```go
// ✅ GOOD: All access serialized through event loop
hub.register <- client1  // Goroutine 1 sends to channel
hub.register <- client2  // Goroutine 2 sends to channel

// Hub's event loop processes one at a time:
for {
    select {
    case client := <-h.register:
        h.clients[client] = true  // Only ONE goroutine accessing!
    }
}
```

**The Magic**: Channels serialize access - only the hub's goroutine modifies the map!

---

<a name="working-together"></a>
## 4. How Hub and Pumps Work Together

### The Complete Flow

```
┌─────────────────────────────────────────────────────────────┐
│                         SERVER                              │
│                                                             │
│                    ┌──────────────┐                         │
│                    │     HUB      │                         │
│                    │  (goroutine) │                         │
│                    └───────┬──────┘                         │
│                            │                                │
│            ┌───────────────┼───────────────┐                │
│            │               │                │               │
│    ┌───────▼──────┐  ┌────▼─────┐  ┌──────▼──────┐          │
│    │   Client 1   │  │ Client 2 │  │  Client 3   │          │
│    │              │  │          │  │             │          │
│    │  readPump()  │  │readPump()│  │ readPump()  │          │
│    │  writePump() │  │writePump()│  │ writePump() │         │
│    │              │  │          │  │             │          │
│    │  send chan   │  │send chan │  │  send chan  │          │
│    └──────┬───────┘  └────┬─────┘  └──────┬──────┘          │
│           │               │               │                 │
└───────────┼───────────────┼───────────────┼─--────────────-─┘
            │               │               │
            │               │               │
      ┌─────▼───┐     ┌────▼────┐      ┌────▼────┐
      │Browser 1│     │Browser 2│      │Browser 3│
      └─────────┘     └─────────┘      └─────────┘
```

---

### Step-by-Step: Complete Message Flow

Let's trace a message from Browser 1 to Browser 2:

```
1. Browser 1 sends: "Hello from Alice!"
   │
   ▼
2. Client 1's readPump() receives it
   │
   ▼
3. readPump sends to hub via hub.broadcast channel
   hub.broadcast <- message
   │
   ▼
4. Hub's event loop receives from broadcast channel
   case message := <-h.broadcast:
   │
   ▼
5. Hub sends to ALL clients' send channels
   for client := range h.clients {
       client.send <- message
   }
   │
   ▼
6. Each client's writePump() receives from their send channel
   case message := <-c.send:
   │
   ▼
7. writePump writes to WebSocket
   conn.WriteMessage(websocket.TextMessage, message)
   │
   ▼
8. Browser 2 receives: "Hello from Alice!"
```

**The Beauty**: Each step is non-blocking and thread-safe!

---

### Code Example: Complete System

```go
// 1. Define Client
type Client struct {
    hub  *Hub
    conn *websocket.Conn
    send chan []byte
}

// 2. Client's readPump (reads from WebSocket)
func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c  // Tell hub we're leaving
        c.conn.Close()
    }()

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }
        // Send to hub for broadcasting
        c.hub.broadcast <- message
    }
}

// 3. Client's writePump (writes to WebSocket)
func (c *Client) writePump() {
    ticker := time.NewTicker(54 * time.Second)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()

    for {
        select {
        case message := <-c.send:
            // Write message to WebSocket
            c.conn.WriteMessage(websocket.TextMessage, message)
        case <-ticker.C:
            // Send ping
            c.conn.WriteMessage(websocket.PingMessage, nil)
        }
    }
}

// 4. Hub manages all clients
type Hub struct {
    clients    map[*Client]bool
    register   chan *Client
    unregister chan *Client
    broadcast  chan []byte
}

// 5. Hub's event loop
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true

        case client := <-h.unregister:
            delete(h.clients, client)
            close(client.send)

        case message := <-h.broadcast:
            for client := range h.clients {
                client.send <- message
            }
        }
    }
}

// 6. HTTP handler creates clients
func handleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
    conn, _ := upgrader.Upgrade(w, r, nil)

    client := &Client{
        hub:  hub,
        conn: conn,
        send: make(chan []byte, 256),
    }

    hub.register <- client  // Register with hub

    go client.writePump()  // Start writing
    go client.readPump()   // Start reading (blocks here)
}

// 7. Main function
func main() {
    hub := NewHub()
    go hub.Run()  // Start hub event loop

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        handleWebSocket(hub, w, r)
    })

    http.ListenAndServe(":8080", nil)
}
```

---

<a name="analogies"></a>
## 5. Real-World Analogies

### Analogy 1: Restaurant

**Hub** = Restaurant Manager
- Manages all tables (clients)
- Coordinates between kitchen and tables
- Handles reservations (register) and checkouts (unregister)

**readPump** = Waiter taking orders
- Continuously listens to customers
- Takes orders to kitchen (hub)

**writePump** = Waiter delivering food
- Brings food from kitchen to table
- Also brings complimentary bread periodically (pings!)

**send channel** = Kitchen pass (food ready area)
- Kitchen puts food here
- Waiter picks it up

---

### Analogy 2: Post Office

**Hub** = Post Office Sorting Center
- Receives mail from many sources
- Sorts and distributes to correct mailboxes

**readPump** = Incoming Mail Truck
- Continuously brings mail to post office
- One truck per sender

**writePump** = Outgoing Mail Truck
- Continuously delivers mail from post office
- One truck per recipient

**send channel** = Mailbox
- Holds mail until delivery truck picks it up

---

### Analogy 3: Radio Station

**Hub** = Radio Station Control Room
- Receives input from DJs, callers, news feeds
- Broadcasts to all listeners

**readPump** = Microphone Input
- Continuously listening for sound
- Sends to control room

**writePump** = Radio Transmitter
- Continuously broadcasting
- Sends to all radios

**send channel** = Audio Buffer
- Holds audio before transmission

---

<a name="code-examples"></a>
## 6. Complete Code Example

Let's build a simple chat room using these patterns:

```go
package main

import (
    "log"
    "net/http"
    "time"

    "github.com/gorilla/websocket"
)

const (
    writeWait      = 10 * time.Second
    pongWait       = 60 * time.Second
    pingPeriod     = (pongWait * 9) / 10
    maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

// Client represents a connected user
type Client struct {
    hub  *Hub
    conn *websocket.Conn
    send chan []byte
}

// Hub manages all clients
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
}

func newHub() *Hub {
    return &Hub{
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
        clients:    make(map[*Client]bool),
    }
}

func (h *Hub) run() {
    for {
        select {
        case client := <-h.register:
            h.clients[client] = true
            log.Printf("New client. Total: %d", len(h.clients))

        case client := <-h.unregister:
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
                log.Printf("Client left. Total: %d", len(h.clients))
            }

        case message := <-h.broadcast:
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
        }
    }
}

func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()

    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }
        c.hub.broadcast <- message
    }
}

func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
                return
            }

        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }
    client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
    client.hub.register <- client

    go client.writePump()
    go client.readPump()
}

func main() {
    hub := newHub()
    go hub.run()

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        serveWs(hub, w, r)
    })

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

<a name="questions"></a>
## 7. Common Questions

### Q: Why not just use one goroutine for read and write?

**A**: Because reading blocks! You can't send pings on schedule if you're waiting for a message.

```go
// ❌ This can't send pings while waiting
for {
    message := conn.ReadMessage()  // STUCK HERE
    // Can't send ping until message arrives!
}
```

---

### Q: What if I want to send a message to just one specific client?

**A**: Send directly to that client's send channel:

```go
// Send to specific client
client.send <- message

// Instead of broadcasting
hub.broadcast <- message
```

---

### Q: How does the hub know which client sent a message?

**A**: The hub can track it:

```go
type MessageWithSender struct {
    sender  *Client
    message []byte
}

// Instead of:
broadcast chan []byte

// Use:
broadcast chan MessageWithSender
```

---

### Q: Can I have multiple hubs?

**A**: Yes! For example, one hub per game room:

```go
type Hub struct {
    roomID  string
    clients map[*Client]bool
    // ...
}

var gameRooms = make(map[string]*Hub)
```

---

### Q: What happens if a client's send channel fills up?

**A**: The default case in broadcast prevents blocking:

```go
select {
case client.send <- message:
    // Sent successfully
default:
    // Buffer full! Close this slow client
    close(client.send)
    delete(h.clients, client)
}
```

---

## Summary

### The Pump Pattern
- **Purpose**: Separate reading and writing for concurrency
- **readPump**: Reads from WebSocket, sends to hub
- **writePump**: Writes to WebSocket, sends pings
- **Communication**: Via `send` channel

### The Hub Pattern
- **Purpose**: Central coordinator for all connections
- **Event Loop**: Processes register/unregister/broadcast
- **Thread-Safe**: All access serialized through channels
- **Scalable**: Handles thousands of connections

### Key Takeaways
1. ✅ Use goroutines for concurrency
2. ✅ Use channels for communication
3. ✅ Separate read and write operations
4. ✅ Centralize connection management in hub
5. ✅ Never block the hub's event loop

---

## Next Steps

Now that you understand Hub and Pump:

1. **Review Lesson 2** - See these patterns in action
2. **Review Implementation Plan** - Apply to chess-coach
3. **Build it!** - Start with foundation, add features incrementally

**You're ready!** 🚀

# WebSocket Core Concepts Explained Simply

**Purpose**: Understand Hub, Pumps, Channels, and Clients before implementing WebSocket in chess-coach.

**Prerequisites**: Basic understanding of Go (variables, functions, structs)

---

## Table of Contents

1. [Starting with Go Basics](#go-basics)
2. [The Confusion Explained](#confusion)
3. [Building Blocks](#building-blocks)
4. [The Complete Picture](#complete-picture)
5. [Chess Game Example](#chess-example)
6. [Quick Reference](#quick-reference)

---

<a name="go-basics"></a>
## 1. Starting with Go Basics

### Goroutines - Running Things Concurrently

```go
// Normal function - runs and WAITS
func greet() {
    fmt.Println("Hello")
}
greet() // Program stops here until greet() finishes

// Goroutine - runs in BACKGROUND
go greet() // Program continues immediately
```

**Real-World Analogy**:
- **Normal function**: Asking someone a question in person - you wait for their answer
- **Goroutine**: Sending someone a text message - you continue with your day

**Why WebSocket needs this**:
WebSocket must do TWO things at once:
1. Listen for messages coming FROM the browser
2. Send messages going TO the browser

If we only had one thread, we'd get stuck waiting for incoming messages and couldn't send anything!

---

### Channels - Communication Between Goroutines

```go
// Create a mailbox
messages := make(chan string)

// Goroutine 1: Put message in mailbox
go func() {
    messages <- "Hello!" // Put in
}()

// Goroutine 2: Take message from mailbox
msg := <-messages // Take out
fmt.Println(msg) // "Hello!"
```

**Real-World Analogy**:
- **Channel** = Physical mailbox
- `messages <- "Hello"` = Putting a letter IN the mailbox
- `msg := <-messages` = Taking a letter OUT of the mailbox

**Why WebSocket needs this**:
When your code wants to send a message to a browser, it can't just write directly to the network connection (that would be unsafe if multiple goroutines do it). Instead:
1. Put the message in a **channel** (safe)
2. One dedicated goroutine reads from that channel and sends it (safe)

---

<a name="confusion"></a>
## 2. The Confusion Explained

### What's Confusing People

You'll see these terms used together and they sound similar:
- Browser
- Client
- Connection
- Channel
- Hub
- Pump

**Let's clarify each one:**

---

### Browser vs Client vs Connection

```
Your Chrome Tab (BROWSER)
        |
        | <-- This physical wire is the WEBSOCKET CONNECTION
        |
Server-side Go struct (CLIENT)
```

**Definitions**:

1. **Browser**: The actual Chrome/Firefox/Safari tab the user sees
2. **WebSocket Connection**: The network wire connecting browser to server
3. **Client (Go struct)**: The code on your server that manages that connection

**Key Insight**:
> One browser tab = One WebSocket connection = One Client struct

**Example**:
- Alice opens `http://chess-coach.com` in Chrome → Creates 1 Client
- Bob opens `http://chess-coach.com` in Firefox → Creates 1 Client
- Charlie opens TWO tabs in Safari → Creates 2 Clients

Total: 4 browser tabs = 4 connections = 4 Client structs on server

---

### Connection vs Channel (HUGE Difference!)

This is the **#1 source of confusion**:

```go
type Client struct {
    conn *websocket.Conn  // ← WebSocket CONNECTION (network wire to browser)
    send chan []byte      // ← Go CHANNEL (mailbox between goroutines)
}
```

**CONNECTION** (`conn`):
- Network wire to the browser
- Sends data over the internet
- Example: The Wi-Fi/Ethernet cable carrying data

**CHANNEL** (`send`):
- Mailbox inside your Go program
- Sends data between goroutines in your server
- Example: A physical inbox on your desk

**Why Both?**

```
Goroutine A: "I want to send this message to the browser"
    |
    v
Puts message in CHANNEL (mailbox)
    |
    v
Goroutine B (writePump): Reads from CHANNEL
    |
    v
Sends via CONNECTION (network wire)
    |
    v
Browser receives message
```

**The Flow**:
1. Any goroutine can safely put messages in the **channel**
2. One dedicated goroutine (writePump) reads from **channel**
3. writePump sends via the **connection** to browser

---

### Hub vs Client

```
HUB (Manager of everyone)
├── Client 1 (Alice's browser)
├── Client 2 (Bob's browser)
├── Client 3 (Charlie's browser tab 1)
└── Client 4 (Charlie's browser tab 2)
```

**CLIENT**:
- Manages ONE browser connection
- Has its own connection wire
- Has its own message mailbox
- Has its own two goroutines (read and write)

**HUB**:
- Knows about ALL clients
- Can send messages to any client
- Tracks who's connected and who disconnected

**Analogy**:
- **Client** = Individual phone line to one person
- **Hub** = Phone switchboard operator who can connect any two people

---

<a name="building-blocks"></a>
## 3. Building Blocks

### The Client Struct

```go
type Client struct {
    ID     string              // Unique ID for this client
    UserID string              // Which user is this? (Alice, Bob, etc.)
    conn   *websocket.Conn     // Network wire to browser
    send   chan []byte         // Mailbox for outgoing messages
    hub    *Hub                // Reference to the central Hub
}
```

**What it represents**: One browser tab's connection to your server

**Visual**:
```
Alice's Chrome Tab
        |
        | (conn - network wire)
        |
Client {
    ID: "client-abc-123"
    UserID: "alice"
    conn: <wire to Alice's browser>
    send: <mailbox for messages going to Alice>
}
```

---

### The Two Pumps

**Problem**: WebSocket needs to do TWO things at once:
1. Read messages FROM browser
2. Write messages TO browser

**Solution**: Two goroutines!

```go
// Goroutine 1: readPump - Listens for incoming messages
func (c *Client) readPump() {
    for {
        // BLOCKS here waiting for message from browser
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break // Browser disconnected
        }
        // Process the message
        handleMessage(message)
    }
}

// Goroutine 2: writePump - Sends outgoing messages
func (c *Client) writePump() {
    for {
        // BLOCKS here waiting for message in channel
        message := <-c.send
        // Send it via the connection
        c.conn.WriteMessage(websocket.TextMessage, message)
    }
}
```

**Starting both pumps**:
```go
go client.readPump()  // Start in background
go client.writePump() // Start in background
```

**Analogy**:
Think of a two-way radio:
- **readPump** = Your ear listening for incoming transmissions
- **writePump** = Your mouth speaking outgoing transmissions
- Both operate simultaneously

**Why not one function?**
```go
// ❌ This DOESN'T work:
func handleConnection() {
    for {
        msg := conn.ReadMessage() // STUCK waiting here
        // Can't send messages while blocked above!
        conn.WriteMessage(response)
    }
}
```

---

### The Hub

```go
type Hub struct {
    clients    map[*Client]bool  // List of all connected clients
    register   chan *Client      // Mailbox: "new client connected"
    unregister chan *Client      // Mailbox: "client disconnected"
    broadcast  chan []byte       // Mailbox: "send to everyone"
}
```

**Hub's Job**: Run an event loop that processes events

```go
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            // New browser connected
            h.clients[client] = true

        case client := <-h.unregister:
            // Browser disconnected
            delete(h.clients, client)
            close(client.send) // Close their mailbox

        case message := <-h.broadcast:
            // Send message to ALL browsers
            for client := range h.clients {
                client.send <- message // Put in each client's mailbox
            }
        }
    }
}
```

**Starting the Hub**:
```go
hub := &Hub{
    clients:    make(map[*Client]bool),
    register:   make(chan *Client),
    unregister: make(chan *Client),
    broadcast:  make(chan []byte),
}

go hub.Run() // Start the event loop in background
```

**Analogy**:
- **Hub** = School principal's office
- **register** = Students enrolling
- **unregister** = Students leaving
- **broadcast** = Announcements over PA system

---

<a name="complete-picture"></a>
## 4. The Complete Picture

### Connecting a Browser

```
Step 1: Browser opens WebSocket
Browser (Chrome) → ws://localhost:8080/ws

Step 2: Server creates Client struct
client := &Client{
    ID:     "client-123",
    UserID: "alice",
    conn:   <connection>,
    send:   make(chan []byte, 256),
}

Step 3: Server starts two goroutines
go client.readPump()  // Start listening
go client.writePump() // Start sending

Step 4: Server registers client with Hub
hub.register <- client

Step 5: Hub adds to its list
h.clients[client] = true
```

**Visual**:
```
Alice opens browser
        ↓
Creates WebSocket connection
        ↓
Server creates Client struct
        ↓
Starts readPump + writePump goroutines
        ↓
Registers with Hub
        ↓
Hub tracks this client
```

---

### Sending a Message

**Scenario**: Server wants to send "Hello" to Alice

```
Step 1: Find Alice's client
aliceClient := <find client where UserID == "alice">

Step 2: Put message in her mailbox
aliceClient.send <- []byte("Hello")

Step 3: writePump is waiting on that channel
message := <-c.send // Receives "Hello"

Step 4: writePump sends via connection
c.conn.WriteMessage(websocket.TextMessage, message)

Step 5: Alice's browser receives "Hello"
```

**Visual**:
```
Server code: "Send 'Hello' to Alice"
        ↓
Put in Alice's send channel (mailbox)
        ↓
Alice's writePump reads from channel
        ↓
writePump writes to connection (network wire)
        ↓
Alice's browser receives "Hello"
```

---

### Broadcasting to Everyone

**Scenario**: Server wants to tell all browsers "Game starting!"

```
Step 1: Put message in broadcast channel
hub.broadcast <- []byte("Game starting!")

Step 2: Hub's Run() loop receives it
message := <-h.broadcast

Step 3: Hub loops through all clients
for client := range h.clients {
    client.send <- message
}

Step 4: Each client's writePump sends it
(Alice's writePump sends to Alice)
(Bob's writePump sends to Bob)
(Charlie's writePump sends to Charlie)
```

**Visual**:
```
Server: "Game starting!"
        ↓
hub.broadcast channel
        ↓
Hub loops through all clients
        ↓
Puts message in each client's send channel
        ↓
Each writePump sends to their browser
        ↓
All browsers receive "Game starting!"
```

---

<a name="chess-example"></a>
## 5. Chess Game Example

### Setup: Alice and Bob Connect

```go
// Alice opens browser
aliceClient := &Client{
    ID:     "client-1",
    UserID: "alice",
    conn:   aliceWebSocket,
    send:   make(chan []byte, 256),
}
go aliceClient.readPump()
go aliceClient.writePump()
hub.register <- aliceClient

// Bob opens browser
bobClient := &Client{
    ID:     "client-2",
    UserID: "bob",
    conn:   bobWebSocket,
    send:   make(chan []byte, 256),
}
go bobClient.readPump()
go bobClient.writePump()
hub.register <- bobClient
```

**Current State**:
```
Hub.clients = {
    aliceClient: true,
    bobClient: true
}
```

---

### Alice Makes a Move

```
1. Alice's browser sends via WebSocket:
   {"action": "move", "from": "e2", "to": "e4"}

2. Alice's readPump receives it:
   _, message, _ := c.conn.ReadMessage()
   // message = {"action": "move", "from": "e2", "to": "e4"}

3. Server processes the move:
   game.MakeMove("e2", "e4")

4. Server creates responses:
   toAlice := []byte(`{"type": "moveAccepted"}`)
   toBob := []byte(`{"type": "opponentMove", "from": "e2", "to": "e4"}`)

5. Server puts messages in mailboxes:
   aliceClient.send <- toAlice
   bobClient.send <- toBob

6. Each writePump sends to their browser:
   // Alice's writePump
   msg := <-c.send // Gets "moveAccepted"
   c.conn.WriteMessage(websocket.TextMessage, msg)

   // Bob's writePump
   msg := <-c.send // Gets "opponentMove"
   c.conn.WriteMessage(websocket.TextMessage, msg)

7. Browsers receive updates:
   Alice sees: "Move accepted ✓"
   Bob sees: "Alice played e2 → e4"
```

**Timeline**:
```
Time  Alice Browser      Alice Client         Bob Client        Bob Browser
----  -------------      ------------         ----------        -----------
T0    Sends move e2→e4
T1                       readPump receives
T2                       Process move
T3                       Put in send ←----→  Put in send
T4                       writePump sends      writePump sends
T5    Sees "OK"                                                 Sees move!
```

---

<a name="quick-reference"></a>
## 6. Quick Reference

### Terminology

| Term | Meaning |
|------|---------|
| **Browser** | Chrome/Firefox/Safari tab the user sees |
| **WebSocket Connection** | Network wire between browser and server |
| **Client (struct)** | Server-side code managing ONE browser connection |
| **Channel** | Go's mailbox for goroutine communication |
| **Hub** | Central manager tracking ALL clients |
| **readPump** | Goroutine that reads FROM browser |
| **writePump** | Goroutine that writes TO browser |

---

### Key Insights

**One browser tab = One Client**
```
Browser Tab → WebSocket Connection → Client struct
```

**Connection vs Channel**
```
Connection: Network wire to browser (over internet)
Channel: Mailbox between goroutines (inside server)
```

**Why Two Pumps?**
```
readPump:  Waits for incoming messages (blocks on read)
writePump: Waits for outgoing messages (blocks on channel)

Both run simultaneously as goroutines
```

**Client vs Hub**
```
Client: Manages ONE browser connection
Hub: Manages ALL clients (can send to any client)
```

---

### The Data Flow

**Sending Message to Browser**:
```
Your code
    ↓
client.send <- message  (put in channel/mailbox)
    ↓
writePump reads from channel
    ↓
writePump writes to conn  (network wire)
    ↓
Browser receives
```

**Receiving Message from Browser**:
```
Browser sends
    ↓
readPump reads from conn  (network wire)
    ↓
Your code processes message
```

---

### Common Questions

**Q: Why can't readPump write directly to the connection?**

A: WebSocket connections are NOT thread-safe. If readPump and writePump both write, data gets corrupted. Solution: Only writePump writes.

**Q: Why does Hub use channels instead of function calls?**

A: Channels make it thread-safe. Multiple goroutines can send to `hub.register` simultaneously without race conditions.

**Q: How many goroutines are running?**

A:
- 1 Hub.Run() goroutine
- 2 per Client (readPump + writePump)
- 3 browser tabs = 1 + (3 × 2) = 7 total goroutines

---

## Visual Summary

```
BROWSER TABS                    SERVER

Chrome (Alice)                  Client 1
    |                           ├─ conn (wire to Alice)
    |───────conn─────────────→  ├─ send (Alice's mailbox)
                                ├─ readPump (listening)
                                └─ writePump (sending)
                                        ↓ registers with
Firefox (Bob)                   Client 2           HUB
    |                           ├─ conn            ├─ clients map
    |───────conn─────────────→  ├─ send            ├─ register chan
                                ├─ readPump ────→  ├─ unregister chan
                                └─ writePump       └─ broadcast chan

Safari (Charlie)                Client 3
    |                           ├─ conn
    |───────conn─────────────→  ├─ send
                                ├─ readPump
                                └─ writePump
```

---

## Next Steps

Now that you understand the core concepts:

1. ✅ **You understand**: Browser, Client, Connection, Channel, Hub, Pumps
2. ⏭️ **Re-read**: [HUB-AND-PUMP-EXPLAINED.md](../HUB-AND-PUMP-EXPLAINED.md) - Will make complete sense now
3. ⏭️ **Review**: [step-25-websocket-implementation-plan.md](../../../backend/project-docs/step-25-websocket-implementation-plan.md)
4. ⏭️ **Implement**: Phase 1 (Foundation) with confidence!

---

**Key Takeaway**:
> WebSocket architecture uses goroutines (for concurrency), channels (for safe communication), Clients (to manage individual connections), and a Hub (to coordinate everyone). Each piece has a clear job, and together they enable real-time communication.

**You're ready! 🚀**

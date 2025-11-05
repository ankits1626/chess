# Practical Go for Chess Coach WebSockets
## Learn Only What You Need, When You Need It

**Problem**: You're seeing confusing syntax like:
- `r.mu.Lock()` - What is `mu`?
- `case message, ok := <-c.send:` - What is `<-` and `:=` together?
- `&Client` - What does `&` mean?
- `<-sigChan` - Reading from channels?
- `time.NewTicker(pingPeriod)` - Tickers?

**Solution**: Learn these concepts step-by-step with simple, runnable examples.

---

## Fast Track: 8 Essential Lessons

### Lesson 1: Variables & Basic Syntax (30 min)
**You Need This For**: Understanding basic Go code

**Confusing Syntax**:
```go
message, ok := <-c.send  // What is :=?
var conn *websocket.Conn  // What is var?
```

**What You'll Learn**:
- `var name type` - explicit variable declaration
- `:=` - short declaration (auto type)
- Multiple assignment: `a, b := 1, 2`
- When to use `var` vs `:=`

**Simple Example**:
```go
// Explicit declaration
var username string
username = "player1"

// Short declaration (common in Go)
message := "hello"  // Go figures out it's a string

// Multiple values
name, score := "alice", 100
```

**Chess Coach Usage**:
```go
// From websocket code
conn, err := upgrader.Upgrade(w, r, nil)  // := creates conn and err
messageType, message, err := conn.ReadMessage()  // := creates all three
```

---

### Lesson 2: Pointers (The `&` and `*`) (45 min)
**You Need This For**: Understanding `&Client`, `*Hub`, pointer receivers

**Confusing Syntax**:
```go
&Client{...}        // What is &?
func (c *Client)    // What is * in function?
c.hub.register <- c // Sending pointer
```

**What You'll Learn**:
- `&` means "address of" (get pointer)
- `*` in type means "pointer to"
- `*` before variable means "value at pointer"
- Why websockets use pointers

**Simple Examples**:

```go
// Example 1: Creating and using pointers
package main
import "fmt"

type Player struct {
    Name  string
    Score int
}

func main() {
    // WITHOUT pointers (makes a copy)
    player1 := Player{Name: "Alice", Score: 0}
    increaseScore(player1)
    fmt.Println(player1.Score) // Still 0! Function got a copy

    // WITH pointers (modifies original)
    player2 := Player{Name: "Bob", Score: 0}
    increaseScorePointer(&player2)  // & means "send address"
    fmt.Println(player2.Score) // Now 10! Function modified original
}

func increaseScore(p Player) {
    p.Score += 10  // This only changes the copy
}

func increaseScorePointer(p *Player) {  // * means "expects a pointer"
    p.Score += 10  // This changes the original
}
```

**Chess Coach Usage**:
```go
// Creating a new client (returns pointer)
client := &Client{
    hub:  hub,
    conn: conn,
    send: make(chan []byte, 256),
}

// Why pointer? So hub can modify the same client
hub.register <- client  // Send pointer to hub
```

---

### Lesson 3: Structs & Methods (45 min)
**You Need This For**: Understanding `Client`, `Hub` structs

**Confusing Syntax**:
```go
func (c *Client) readPump() {}  // What is (c *Client)?
type Client struct { ... }       // How do structs work?
```

**What You'll Learn**:
- Defining structs (like classes in other languages)
- Methods on structs
- Pointer receivers vs value receivers
- Why websockets use pointer receivers

**Simple Example**:

```go
package main
import "fmt"

// Define a struct (like a class)
type GameSession struct {
    gameID   string
    player1  string
    player2  string
    moves    int
}

// Method with pointer receiver (can modify the struct)
func (g *GameSession) makeMove(player string) {
    g.moves++  // This actually changes g
    fmt.Printf("%s made move #%d\n", player, g.moves)
}

// Method with value receiver (can't modify, gets a copy)
func (g GameSession) getMoveCount() int {
    return g.moves  // Read-only
}

func main() {
    // Create a game session
    game := &GameSession{
        gameID:  "game123",
        player1: "Alice",
        player2: "Bob",
        moves:   0,
    }

    game.makeMove("Alice")  // moves becomes 1
    game.makeMove("Bob")    // moves becomes 2
    fmt.Printf("Total moves: %d\n", game.getMoveCount())
}
```

**Chess Coach Usage**:
```go
// Client struct holds connection data
type Client struct {
    hub  *Hub
    conn *websocket.Conn
    send chan []byte
}

// Method on Client (pointer receiver so it can modify)
func (c *Client) readPump() {
    // c can access c.hub, c.conn, c.send
    // pointer receiver so we work with the SAME client
}
```

---

### Lesson 4: Channels Basics (The `<-` operator) (60 min)
**You Need This For**: Understanding `c.send`, `hub.register`, message passing

**Confusing Syntax**:
```go
send := make(chan []byte, 256)  // What is chan?
c.send <- message               // What is <-?
msg := <-c.send                 // <- on the other side?
```

**What You'll Learn**:
- What channels are (pipes for goroutines)
- Creating channels with `make(chan Type)`
- Sending: `channel <- value`
- Receiving: `value := <-channel`
- Buffered vs unbuffered

**Simple Example**:

```go
package main
import (
    "fmt"
    "time"
)

func main() {
    // Create a channel (pipe) for messages
    messages := make(chan string)

    // Send in a goroutine (separate thread)
    go func() {
        time.Sleep(1 * time.Second)
        messages <- "ping"  // <- sends value INTO channel
    }()

    // Receive (waits until message arrives)
    msg := <-messages  // <- gets value FROM channel
    fmt.Println(msg)   // "ping"
}
```

**Practical Example - Chat Room**:

```go
package main
import "fmt"

func main() {
    // Buffered channel (can hold 3 messages before blocking)
    chatRoom := make(chan string, 3)

    // Send messages (doesn't block because buffer has space)
    chatRoom <- "Alice: Hello!"
    chatRoom <- "Bob: Hi there!"
    chatRoom <- "Charlie: Hey!"

    // Receive messages
    fmt.Println(<-chatRoom)  // "Alice: Hello!"
    fmt.Println(<-chatRoom)  // "Bob: Hi there!"
    fmt.Println(<-chatRoom)  // "Charlie: Hey!"
}
```

**Chess Coach Usage**:
```go
// Client has a send channel (buffered, holds 256 messages)
type Client struct {
    send chan []byte  // Channel for outgoing messages
}

// Sending to client
client.send <- []byte("your move")  // Put message in channel

// Reading from client (in writePump)
message := <-client.send  // Get message from channel
conn.WriteMessage(websocket.TextMessage, message)
```

---

### Lesson 5: Select Statement (Multiple Channels) (45 min)
**You Need This For**: Understanding `select` in readPump/writePump

**Confusing Syntax**:
```go
select {
case message := <-c.send:
    // do something
case <-ticker.C:
    // do something else
}
```

**What You'll Learn**:
- `select` waits on multiple channels
- Like `switch` but for channels
- First ready channel wins
- `case <-ticker.C:` means "when ticker sends"

**Simple Example**:

```go
package main
import (
    "fmt"
    "time"
)

func main() {
    messages := make(chan string)
    timeout := time.After(2 * time.Second)

    go func() {
        time.Sleep(1 * time.Second)
        messages <- "message arrived!"
    }()

    // Wait for EITHER message OR timeout
    select {
    case msg := <-messages:
        fmt.Println("Got:", msg)
    case <-timeout:
        fmt.Println("Timeout! No message received")
    }
}
```

**Chess Coach Usage**:
```go
// In writePump - send messages OR ping
select {
case message, ok := <-c.send:
    // Got a message to send to client
    if !ok {
        // Channel closed
        return
    }
    c.conn.WriteMessage(websocket.TextMessage, message)

case <-ticker.C:
    // Time to send a ping (heartbeat)
    c.conn.WriteMessage(websocket.PingMessage, nil)
}
```

---

### Lesson 6: Goroutines (The `go` keyword) (45 min)
**You Need This For**: Understanding `go c.writePump()`, `go c.readPump()`

**Confusing Syntax**:
```go
go c.writePump()     // What does 'go' do?
go func() { ... }()  // Anonymous function?
```

**What You'll Learn**:
- `go` starts a new lightweight thread
- Goroutines run concurrently
- Why websockets use 2 goroutines (read + write)
- Anonymous functions with `go func() { }`

**Simple Example**:

```go
package main
import (
    "fmt"
    "time"
)

func sayHello(name string) {
    for i := 0; i < 3; i++ {
        fmt.Printf("Hello %s (%d)\n", name, i)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    // Start goroutine (runs in background)
    go sayHello("Alice")
    go sayHello("Bob")

    // Main goroutine continues
    fmt.Println("Main: started both goroutines")

    // Wait for them to finish
    time.Sleep(500 * time.Millisecond)
    fmt.Println("Main: done")
}
```

**Chess Coach Usage**:
```go
// When client connects, start TWO goroutines
go client.writePump()  // Goroutine 1: sends messages to client
go client.readPump()   // Goroutine 2: receives messages from client

// They run concurrently!
// readPump receives from websocket → puts in hub
// writePump receives from hub → sends to websocket
```

---

### Lesson 7: Mutexes (The `mu.Lock()`) (45 min)
**You Need This For**: Understanding `r.mu.Lock()`, `r.mu.Unlock()`

**Confusing Syntax**:
```go
r.mu.Lock()
defer r.mu.Unlock()
r.clients[client] = true
```

**What You'll Learn**:
- Why we need locks (prevent race conditions)
- `sync.Mutex` for exclusive access
- `Lock()` and `Unlock()`
- `defer` ensures unlock happens

**Simple Example**:

```go
package main
import (
    "fmt"
    "sync"
)

type Counter struct {
    mu    sync.Mutex  // Lock for safety
    value int
}

func (c *Counter) Increment() {
    c.mu.Lock()           // Lock before modifying
    defer c.mu.Unlock()   // Unlock when function ends
    c.value++
}

func (c *Counter) GetValue() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

func main() {
    counter := &Counter{}

    // Simulate 1000 concurrent increments
    var wg sync.WaitGroup
    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Increment()
        }()
    }

    wg.Wait()
    fmt.Printf("Final value: %d\n", counter.GetValue()) // Always 1000 (safe!)
}
```

**Chess Coach Usage**:
```go
type Hub struct {
    clients    map[*Client]bool
    mu         sync.Mutex  // Protects clients map
    register   chan *Client
    unregister chan *Client
}

func (h *Hub) registerClient(client *Client) {
    h.mu.Lock()                  // Lock before modifying map
    defer h.mu.Unlock()          // Unlock when done
    h.clients[client] = true     // Safe to modify now
}
```

---

### Lesson 8: Putting It All Together (90 min)
**You Need This For**: Understanding the full websocket flow

**What You'll Learn**:
- How all concepts combine
- Hub pattern (central registry)
- Client lifecycle (connect → read/write → disconnect)
- Message flow through channels

**Complete Minimal Example**:

```go
package main

import (
    "fmt"
    "sync"
)

// Hub manages all clients
type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mu         sync.Mutex
}

// Client represents one connection
type Client struct {
    hub  *Hub
    send chan []byte
    id   string
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte, 256),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mu.Lock()
            h.clients[client] = true
            h.mu.Unlock()
            fmt.Printf("Client %s registered\n", client.id)

        case client := <-h.unregister:
            h.mu.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
                fmt.Printf("Client %s unregistered\n", client.id)
            }
            h.mu.Unlock()

        case message := <-h.broadcast:
            h.mu.Lock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mu.Unlock()
        }
    }
}

func (c *Client) readPump() {
    // Simulate reading messages
    for i := 0; i < 3; i++ {
        message := []byte(fmt.Sprintf("Message %d from %s", i, c.id))
        c.hub.broadcast <- message
    }
    c.hub.unregister <- c
}

func (c *Client) writePump() {
    for message := range c.send {
        fmt.Printf("[%s] Sending: %s\n", c.id, string(message))
    }
}

func main() {
    hub := NewHub()
    go hub.Run()

    // Simulate 2 clients
    client1 := &Client{hub: hub, send: make(chan []byte, 256), id: "Alice"}
    client2 := &Client{hub: hub, send: make(chan []byte, 256), id: "Bob"}

    hub.register <- client1
    hub.register <- client2

    go client1.writePump()
    go client2.writePump()
    go client1.readPump()
    go client2.readPump()

    // Keep running
    select {}
}
```

**Flow**:
1. Hub runs in goroutine, waiting for events
2. Clients register via channel
3. Client1 sends message → broadcast channel
4. Hub receives broadcast → sends to all clients
5. Clients receive via their `send` channel
6. writePump sends to websocket

---

## Learning Path

### Week 1: Core Syntax
- **Day 1**: Lesson 1 (Variables) + Lesson 2 (Pointers)
- **Day 2**: Lesson 3 (Structs & Methods)
- **Day 3**: Review + Practice

### Week 2: Concurrency
- **Day 4**: Lesson 4 (Channels)
- **Day 5**: Lesson 5 (Select) + Lesson 6 (Goroutines)
- **Day 6**: Review + Practice

### Week 3: Advanced + Integration
- **Day 7**: Lesson 7 (Mutexes)
- **Day 8**: Lesson 8 (Full Example)
- **Day 9**: Apply to chess-coach websockets

---

## Confusing Syntax Cheat Sheet

| Syntax | Meaning | Example |
|--------|---------|---------|
| `:=` | Short variable declaration | `x := 5` |
| `&` | Address of (get pointer) | `ptr := &myStruct` |
| `*Type` | Pointer to Type | `var p *int` |
| `*ptr` | Value at pointer | `val := *ptr` |
| `<-ch` | Receive from channel | `msg := <-messages` |
| `ch <-` | Send to channel | `messages <- "hi"` |
| `func (c *Client)` | Method with pointer receiver | Can modify `c` |
| `go func()` | Start goroutine | Runs concurrently |
| `mu.Lock()` | Acquire lock | Exclusive access |
| `defer` | Run when function exits | Cleanup |
| `select { case }` | Wait on channels | First ready wins |

---

## Next Steps

1. Start with **Lesson 1** (create folder with examples)
2. Run each example
3. Modify and experiment
4. Move to next lesson when comfortable
5. After Lesson 8, revisit your websocket code - it will make sense!

**Ready to create Lesson 1 with runnable examples?**

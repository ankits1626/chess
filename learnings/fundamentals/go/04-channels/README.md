# Lesson 4: Channels & The `<-` Operator

## Why This Lesson?

You're seeing this everywhere in websocket code:
```go
case message := <-c.send          // What is <-?
c.hub.register <- c               // Arrow on the right?
<-ticker.C                        // No variable?
message, ok := <-c.send           // What is ok?
```

By the end of this lesson, you'll master the **`<-` operator** and understand how goroutines communicate.

---

## What Are Channels?

**Channels are pipes that connect goroutines.**

Think of goroutines as workers, and channels as pipes between them:

```
Worker 1                Channel              Worker 2
   │                       ║                     │
   │  data ──────────>     ║   ──────────> reads data
   │                       ║                     │
```

**Why needed?** Goroutines run independently. Channels let them communicate safely.

---

## The `<-` Operator

### Official Name: **Channel Operator**

The `<-` is **not** a comparison operator like `<=`. It's a special operator for channels.

### Two Forms

#### Form 1: Receive (arrow before channel)
```go
value := <-channel   // Get value FROM channel
```

#### Form 2: Send (arrow after channel)
```go
channel <- value     // Put value INTO channel
```

**The arrow shows data flow direction!**

---

## Creating Channels

### Syntax: `make(chan TYPE)`

```go
// Channel for strings
messages := make(chan string)

// Channel for ints
numbers := make(chan int)

// Channel for byte slices (common in websockets)
data := make(chan []byte)

// Buffered channel (holds 10 values)
buffered := make(chan string, 10)
```

---

## Sending Data

**Syntax**: `channel <- value`

```go
messages := make(chan string)

// Send data INTO channel
messages <- "hello"
messages <- "world"
```

**Important**: Sending **blocks** (waits) until someone receives!

---

## Receiving Data

**Syntax**: `value := <-channel`

```go
messages := make(chan string)

// Receive FROM channel
msg := <-messages
fmt.Println(msg)
```

**Important**: Receiving **blocks** (waits) until data arrives!

---

## Complete Example

```go
package main

import "fmt"

func main() {
    // Create channel
    messages := make(chan string)

    // Sender (in goroutine)
    go func() {
        messages <- "ping"  // Send
    }()

    // Receiver (main goroutine)
    msg := <-messages  // Receive
    fmt.Println(msg)   // "ping"
}
```

**Why goroutine?** If we send without a goroutine, it blocks forever (no one to receive)!

---

## The Direction of `<-`

### Think Visually

```go
// RECEIVE: Arrow points FROM channel TO variable
msg := <-messages
   ↑      ↑
 variable channel

// SEND: Arrow points FROM value TO channel
messages <- "hello"
  ↑          ↑
channel    value
```

### More Examples

```go
// Receive
x := <-ch       // x gets value from ch
y := <-numbers  // y gets number from numbers

// Send
ch <- 42            // Send 42 into ch
names <- "alice"    // Send "alice" into names
data <- []byte{1,2} // Send bytes into data
```

---

## Receiving Without Storing

Sometimes you don't care about the value, just the signal:

```go
done := make(chan bool)

go func() {
    // Do work...
    done <- true  // Signal we're done
}()

<-done  // Wait for signal (don't store the value)
fmt.Println("Work complete!")
```

**Use case**: Synchronization, timing, waiting

---

## The `ok` Pattern

Channels can be **closed**. Check if closed with `ok`:

```go
message, ok := <-messages
if !ok {
    fmt.Println("Channel is closed!")
    return
}
fmt.Println("Got:", message)
```

- `ok = true` → Channel is open, message is valid
- `ok = false` → Channel is closed, no more data

---

## Closing Channels

**Syntax**: `close(channel)`

```go
messages := make(chan string)

// Sender
go func() {
    messages <- "hello"
    messages <- "world"
    close(messages)  // Signal: no more data
}()

// Receiver can detect closure
for {
    msg, ok := <-messages
    if !ok {
        break  // Channel closed
    }
    fmt.Println(msg)
}
```

**Rule**: Only the **sender** should close a channel, never the receiver.

---

## Buffered vs Unbuffered Channels

### Unbuffered (default)

```go
ch := make(chan int)  // No buffer
```

**Behavior**:
- Send blocks until someone receives
- Receive blocks until someone sends
- Synchronous communication

### Buffered

```go
ch := make(chan int, 3)  // Buffer size 3
```

**Behavior**:
- Can send 3 values without blocking
- Only blocks when buffer is full
- Asynchronous communication

### Example

```go
// Unbuffered - would deadlock without goroutine
unbuffered := make(chan int)
// unbuffered <- 1  // DEADLOCK! No receiver

// Buffered - works without goroutine
buffered := make(chan int, 2)
buffered <- 1  // OK, buffer has space
buffered <- 2  // OK, buffer still has space
// buffered <- 3  // Would block, buffer is full

fmt.Println(<-buffered)  // 1
fmt.Println(<-buffered)  // 2
```

---

## Websocket Patterns

### Pattern 1: Client Send Buffer

```go
type Client struct {
    send chan []byte  // Buffered channel for outgoing messages
}

// Create client
client := &Client{
    send: make(chan []byte, 256),  // Can queue 256 messages
}

// Hub sends to client
client.send <- message  // Put message in queue

// writePump receives and sends to websocket
msg := <-client.send  // Get message from queue
conn.WriteMessage(websocket.TextMessage, msg)
```

**Why buffered?** Prevents hub from blocking if client is slow.

### Pattern 2: Hub Register/Unregister

```go
type Hub struct {
    register   chan *Client  // Unbuffered
    unregister chan *Client  // Unbuffered
}

// Register client
hub.register <- client  // Send client to hub

// Hub receives registration
client := <-h.register  // Get client
h.clients[client] = true
```

### Pattern 3: Select with Multiple Channels

```go
select {
case message, ok := <-c.send:
    // Got message from send channel
    if !ok {
        return  // Channel closed
    }
    c.conn.WriteMessage(websocket.TextMessage, message)

case <-ticker.C:
    // Ticker fired (heartbeat)
    c.conn.WriteMessage(websocket.PingMessage, nil)
}
```

**Explained**:
- `case message, ok := <-c.send:` → Wait for message
- `case <-ticker.C:` → Wait for ticker (discard value)
- Whichever happens first wins

---

## Common Patterns

### Pattern 1: Worker Pool

```go
jobs := make(chan int, 100)
results := make(chan int, 100)

// Worker
go func() {
    for job := range jobs {
        result := job * 2
        results <- result
    }
}()

// Send jobs
for i := 1; i <= 5; i++ {
    jobs <- i
}
close(jobs)

// Receive results
for i := 1; i <= 5; i++ {
    fmt.Println(<-results)
}
```

### Pattern 2: Done Signal

```go
done := make(chan bool)

go func() {
    time.Sleep(2 * time.Second)
    done <- true
}()

fmt.Println("Waiting...")
<-done
fmt.Println("Done!")
```

### Pattern 3: Broadcast to Multiple

```go
clients := make(map[*Client]bool)
message := []byte("broadcast")

for client := range clients {
    select {
    case client.send <- message:
        // Sent successfully
    default:
        // Client buffer full, skip
        close(client.send)
        delete(clients, client)
    }
}
```

---

## Channel Directions (Advanced)

You can specify send-only or receive-only channels:

```go
// Send-only channel
func sender(ch chan<- string) {
    ch <- "hello"  // OK
    // msg := <-ch  // ERROR: receive from send-only
}

// Receive-only channel
func receiver(ch <-chan string) {
    msg := <-ch    // OK
    // ch <- "hi"  // ERROR: send to receive-only
}

func main() {
    messages := make(chan string)
    go sender(messages)
    receiver(messages)
}
```

**Use case**: Function signatures that enforce send or receive only.

---

## Blocking Behavior

**Critical concept**: Channels block!

```go
// This BLOCKS until someone receives
ch <- value

// This BLOCKS until someone sends
value := <-ch
```

**Good**: Forces synchronization
**Bad**: Can cause deadlocks if misused

---

## Deadlock Example

```go
// DEADLOCK!
ch := make(chan int)
ch <- 1  // Blocks forever - no receiver!
```

**Error**: `fatal error: all goroutines are asleep - deadlock!`

**Fix**: Use goroutine
```go
ch := make(chan int)
go func() {
    ch <- 1  // Send in goroutine
}()
fmt.Println(<-ch)  // Receive in main
```

---

## Range Over Channels

Receive until channel closes:

```go
messages := make(chan string)

go func() {
    messages <- "hello"
    messages <- "world"
    close(messages)
}()

// Receive all messages
for msg := range messages {
    fmt.Println(msg)
}
// Loop exits when channel closes
```

---

## Quick Reference

| Code | Meaning |
|------|---------|
| `make(chan T)` | Create unbuffered channel |
| `make(chan T, N)` | Create buffered channel (size N) |
| `ch <- v` | Send v to channel ch |
| `v := <-ch` | Receive from ch, store in v |
| `v, ok := <-ch` | Receive + check if open |
| `<-ch` | Receive and discard |
| `close(ch)` | Close channel |
| `for v := range ch` | Receive until closed |

---

## Chess Coach Examples

### Example 1: Hub Run Loop

```go
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            // Client wants to register
            h.clients[client] = true

        case client := <-h.unregister:
            // Client wants to unregister
            delete(h.clients, client)
            close(client.send)

        case message := <-h.broadcast:
            // Broadcast to all clients
            for client := range h.clients {
                client.send <- message
            }
        }
    }
}
```

### Example 2: Client writePump

```go
func (c *Client) writePump() {
    ticker := time.NewTicker(pingPeriod)
    defer ticker.Stop()

    for {
        select {
        case message, ok := <-c.send:
            // Receive from send channel
            if !ok {
                // Channel closed by hub
                return
            }
            c.conn.WriteMessage(websocket.TextMessage, message)

        case <-ticker.C:
            // Ticker fired - send ping
            c.conn.WriteMessage(websocket.PingMessage, nil)
        }
    }
}
```

---

## Exercises

I've created 6 progressive examples and exercises:

1. **Basic send/receive**
2. **Buffered channels**
3. **The ok pattern**
4. **Worker pool**
5. **Ping-pong game**
6. **Websocket hub simulation**

Run them in order!

---

## Next Steps

1. **Run examples** in `examples/` folder
2. **Complete exercises** in `exercises/` folder
3. **Find patterns** in your websocket code
4. **Fill out** `SUMMARY.md`

After this lesson, you'll understand every `<-` in your codebase!

---

## What's Next?

**Lesson 5**: Select Statement - Understanding `select { case ... }`

This builds on channels to handle multiple channel operations.

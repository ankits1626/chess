# Lesson 4: Channels & The `<-` Operator - START HERE

## What You'll Master

By the end of this lesson, you'll understand ALL of these:
```go
message := <-c.send              // ✅ You'll know this!
c.send <- []byte("data")         // ✅ And this!
case msg, ok := <-c.send:        // ✅ And this!
hub.register <- client           // ✅ And this!
<-ticker.C                       // ✅ And this!
```

---

## Why This Lesson Matters

**The `<-` operator is EVERYWHERE in your websocket code:**
- Hub communicating with clients
- Clients queueing messages
- Timer/ticker for heartbeats
- Registration/unregistration of clients

**After this lesson**, your websocket code will make complete sense!

---

## Learning Path (60 minutes)

### Step 1: Understand Channels (15 min)
Read: [README.md](README.md)

Focus on:
- What channels are (pipes between goroutines)
- The `<-` operator (send vs receive)
- Buffered vs unbuffered
- The `ok` pattern for closure detection

### Step 2: Run Examples (20 min)

```bash
cd examples

# Run in order:
go run 01-basic.go              # Send & receive basics
go run 02-buffered.go           # Buffered channels
go run 03-ok-pattern.go         # Detecting closure
go run 04-worker-pool.go        # Multiple workers
go run 05-ping-pong.go          # Two-way communication
go run 06-websocket-simulation.go  # ⭐ MUST RUN - This is your hub!
```

**Most Important**: `06-websocket-simulation.go` - This shows EXACTLY how your hub works!

### Step 3: Practice (20 min)

```bash
cd ../exercises
go run practice.go
```

Complete the TODOs in the exercises.

### Step 4: Apply to Websockets (15 min)

Open your websocket code and find:
1. Where hub receives registrations: `client := <-h.register`
2. Where hub broadcasts: `client.send <- message`
3. Where writePump receives: `msg := <-c.send`

Document these in your SUMMARY.md!

---

## Quick Reference Card

| You See | It Means |
|---------|----------|
| `ch <- v` | Send `v` INTO channel `ch` |
| `v := <-ch` | Receive FROM channel `ch` into `v` |
| `v, ok := <-ch` | Receive + check if channel is open |
| `<-ch` | Receive and discard (just wait) |
| `make(chan T)` | Create unbuffered channel |
| `make(chan T, N)` | Create buffered channel (size N) |
| `close(ch)` | Close channel (sender only!) |
| `for v := range ch` | Receive until closed |

---

## Arrow Direction Trick

**Think of the arrow as data flow:**

```
Sender              Channel              Receiver
   │                   │                     │
   │  value ────────>  │  ────────> result  │
   │  (ch <- v)        │  (v := <-ch)       │
```

**The arrow points in the direction data flows!**

---

## Expected Outcomes

After completing:
- ✅ Understand `<-` operator completely
- ✅ Know when to use buffered vs unbuffered
- ✅ Recognize hub registration pattern
- ✅ Understand client send queue pattern
- ✅ See how goroutines communicate safely

---

## Prerequisites

Before starting this lesson:
- ✅ Complete Lesson 1 (Variables & `:=`)
- ✅ Basic understanding of goroutines (we'll cover more in Lesson 6)

**Note**: You don't need to fully understand goroutines yet. Just know:
- `go func()` starts a background task
- Channels let these tasks communicate

---

## Common Mistakes

### Mistake 1: Deadlock
```go
ch := make(chan int)
ch <- 42  // ❌ DEADLOCK! No receiver
```

**Fix**: Use goroutine
```go
ch := make(chan int)
go func() { ch <- 42 }()
fmt.Println(<-ch)  // ✅ OK
```

### Mistake 2: Send on Closed Channel
```go
close(ch)
ch <- 42  // ❌ PANIC!
```

**Fix**: Don't send after close

### Mistake 3: Receive from Nil Channel
```go
var ch chan int  // nil
<-ch  // ❌ Blocks forever
```

**Fix**: Use `make(chan int)`

---

## How This Connects to Chess Coach

### Your Hub Pattern

```go
// In hub.go
func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:     // ← Receive registration
            h.clients[client] = true

        case client := <-h.unregister:   // ← Receive unregistration
            delete(h.clients, client)
            close(client.send)

        case message := <-h.broadcast:   // ← Receive broadcast
            for client := range h.clients {
                client.send <- message   // ← Send to each client
            }
        }
    }
}
```

### Your Client Pattern

```go
// In client.go - writePump
for {
    select {
    case message, ok := <-c.send:       // ← Receive from queue
        if !ok {
            return  // Channel closed
        }
        c.conn.WriteMessage(websocket.TextMessage, message)

    case <-ticker.C:                    // ← Receive ticker signal
        c.conn.WriteMessage(websocket.PingMessage, nil)
    }
}
```

**After this lesson, you'll understand EVERY `<-` in this code!**

---

## Time Estimates

- **Quick**: 30 min (read + run examples)
- **Thorough**: 60 min (+ exercises)
- **Deep dive**: 90 min (+ find all `<-` in your code)

---

## Let's Begin!

1. **Start**: Open [README.md](README.md)
2. **Run**: Examples in `examples/` folder
3. **Practice**: Exercises in `exercises/practice.go`
4. **Apply**: Find patterns in your websocket code
5. **Summarize**: Fill in `SUMMARY.md`

**Most important example**: `06-websocket-simulation.go` - Run this FIRST if you want to see the big picture!

Good luck! 🚀

---

## Next Lesson Preview

**Lesson 5**: Select Statement

You'll learn the `select` statement that makes the hub work:
```go
select {
case msg := <-channel1:  // ← You know this now!
case <-channel2:         // ← You know this too!
}
// Next lesson: How select CHOOSES which case to execute!
```

# Lesson 6: Goroutines - START HERE

## Why This Lesson BEFORE Channels?

**You were 100% right!** Channels don't make sense without understanding goroutines first.

**Goroutines** = Independent workers that run concurrently
**Channels** = Pipes that connect those workers

Without goroutines → No need for channels!

---

## What You'll Understand (45 min)

After this lesson, you'll know WHY you see:
```go
go c.readPump()     // ✅ You'll understand why!
go c.writePump()    // ✅ Why TWO goroutines!
go hub.Run()        // ✅ Why hub needs goroutine!
```

---

## The `go` Keyword

**One keyword, infinite power:**

```go
func doWork() {
    fmt.Println("Working...")
}

// Normal call (blocks/waits)
doWork()

// Goroutine (runs in background)
go doWork()
```

That's it! `go` makes it run concurrently.

---

## Learning Path

### Step 1: Understand Why (10 min)
Read: [README.md](README.md) - Focus on "Why You Need Goroutines in Websockets"

**The core problem**:
- Websockets need to READ and WRITE at the same time
- One function can't do both (they both block!)
- Solution: Two goroutines running simultaneously

### Step 2: Run Examples (20 min)

```bash
cd examples

# Run in order:
go run 01-basic.go          # The 'go' keyword
go run 02-anonymous.go      # go func() pattern
go run 03-waitgroup.go      # Proper synchronization
go run 04-worker-pool.go    # Controlled concurrency
go run 05-websocket-pumps.go    # ⭐ readPump + writePump!
go run 06-hub-simulation.go     # ⭐ Why go hub.Run()!
```

**Most Important**:
- `05-websocket-pumps.go` - Shows WHY you need two goroutines
- `06-hub-simulation.go` - Shows WHY hub runs in background

### Step 3: Practice (15 min)

```bash
cd ../exercises
go run practice.go
```

Complete the TODOs.

---

## Quick Visual

### Sequential (Normal)
```
Task 1 ━━━━━━━━━━
                  Task 2 ━━━━━━━━
                                  Task 3 ━━━━━
Total time: 25 seconds
```

### Concurrent (Goroutines)
```
Task 1 ━━━━━━━━━━
Task 2 ━━━━━━━━
Task 3 ━━━━━
Total time: 10 seconds (overlapping!)
```

---

## The Websocket Use Case

**Problem**:
```go
// Can't do both in one function!
for {
    msg := conn.ReadMessage()   // Blocks here waiting
    conn.WriteMessage(response) // Can't get here until message arrives!
}
```

**Solution**:
```go
// Two goroutines running simultaneously
go readPump()   // Goroutine 1: Waits for incoming messages
go writePump()  // Goroutine 2: Waits for outgoing messages
```

Now they can **both wait at the same time**!

---

## Common Patterns You'll See

### Pattern 1: Client Connection
```go
func handleWebSocket(conn *websocket.Conn) {
    client := &Client{conn: conn}

    // Start TWO goroutines per client
    go client.readPump()   // Read from websocket
    go client.writePump()  // Write to websocket
}
```

### Pattern 2: Hub in Background
```go
func main() {
    hub := newHub()
    go hub.Run()  // Run hub in background (infinite loop)

    // Main can continue to start HTTP server
    http.ListenAndServe(":8080", nil)
}
```

### Pattern 3: Anonymous Goroutine
```go
go func() {
    if err := server.Start(); err != nil {
        log.Fatal(err)
    }
}()
```

---

## Quick Reference

| Code | Meaning |
|------|---------|
| `go f()` | Start function as goroutine |
| `go func() {}()` | Start anonymous function as goroutine |
| `time.Sleep(d)` | Wait for duration |
| `wg.Add(1)` | Register goroutine with WaitGroup |
| `wg.Done()` | Mark goroutine complete |
| `wg.Wait()` | Wait for all goroutines |

---

## Expected Outcomes

After completing:
- ✅ Understand what `go` does
- ✅ Know why websockets need multiple goroutines
- ✅ Recognize readPump/writePump pattern
- ✅ Understand why hub runs in goroutine
- ✅ Be ready for Lesson 4 (Channels)!

---

## Time Estimate

- **Fast**: 30 min (read + run key examples)
- **Thorough**: 45 min (+ exercises)
- **Deep**: 60 min (+ find all `go` in your code)

---

## Common Mistakes to Avoid

### Mistake 1: Forgetting to Wait
```go
go doWork()
// Main exits immediately - goroutine killed!
```

**Fix**: Use `time.Sleep()` or `WaitGroup`

### Mistake 2: Loop Variable Bug
```go
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i)  // ❌ Might all print 3!
    }()
}
```

**Fix**: Pass as parameter
```go
for i := 0; i < 3; i++ {
    go func(id int) {
        fmt.Println(id)  // ✅ Each gets own copy
    }(i)
}
```

---

## After This Lesson

**Then do Lesson 4: Channels** - It will make perfect sense!

You'll understand:
```go
// Goroutines (this lesson)        Channels (Lesson 4)
go c.readPump()                    message := <-c.send
go c.writePump()                   c.send <- data
go hub.Run()                       hub.broadcast <- msg
```

---

## Let's Begin!

**Start here**: [README.md](README.md)

**Most important examples**:
- `05-websocket-pumps.go` - Why readPump + writePump
- `06-hub-simulation.go` - Why `go hub.Run()`

After this, Lesson 4 (Channels) will click instantly! 🚀

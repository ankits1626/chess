# Lesson 6 Summary: Goroutines

**Complete this after finishing the lesson!**

## What I Learned

### Key Concepts
- [ ] What goroutines are (lightweight threads)
- [ ] The `go` keyword starts a goroutine
- [ ] Goroutines run concurrently (simultaneously)
- [ ] Main goroutine kills all others when it exits
- [ ] WaitGroup for proper synchronization
- [ ] Anonymous goroutines: `go func() {}`

### The `go` Keyword - Now Clear!

**Before Lesson**: `go c.readPump()` ← What does this do??

**After Lesson**:
- `go` starts the function in a new goroutine
- Runs concurrently with main goroutine
- Doesn't block - main continues immediately
- Needed for readPump + writePump to run simultaneously!

---

## My Notes

1. **What is a goroutine in your own words?**


2. **Why do websockets need TWO goroutines per client?**


3. **What happens if main() exits before goroutines finish?**


4. **Any remaining questions?**


---

## Examples I Ran

- [ ] 01-basic.go (basic go keyword)
- [ ] 02-anonymous.go (go func() pattern)
- [ ] 03-waitgroup.go (synchronization)
- [ ] 04-worker-pool.go (controlled concurrency)
- [ ] 05-websocket-pumps.go ⭐ (readPump + writePump)
- [ ] 06-hub-simulation.go ⭐ (hub in background)

**Most helpful example**: ___________________________

---

## Exercises Completed

- [ ] Exercise 1: Basic goroutine
- [ ] Exercise 2: Multiple goroutines
- [ ] Exercise 3: Closure
- [ ] Exercise 4: WaitGroup
- [ ] Exercise 5: Read/Write pumps
- [ ] Exercise 6: Race condition
- [ ] Bonus: Worker pool

---

## Goroutine Patterns Found in My Code

### Pattern 1: Client Read/Write Pumps
**File**: `_______________________`
**Lines**: `_____ - _____`

```go
// Code with 'go' keyword:
go client.readPump()
go client.writePump()

// Why TWO goroutines?

```

### Pattern 2: Hub Running in Background
**File**: `_______________________`
**Line**: `_____`

```go
// Code with 'go' keyword:
go hub.Run()

// Why in goroutine?

```

### Pattern 3: Other Goroutines
**File**: `_______________________`
**Line**: `_____`

```go
// Code:


// Purpose:

```

---

## Why Websockets Need Goroutines

I now understand:

**readPump** (Goroutine 1):
- Continuously reads from websocket
- Blocks waiting for messages
- Can't do anything else while waiting

**writePump** (Goroutine 2):
- Continuously writes to websocket
- Blocks waiting for messages to send
- Can't do anything else while waiting

**They MUST run simultaneously** → Need separate goroutines!

---

## WaitGroup Pattern

I understand how to properly wait for goroutines:

```go
var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)              // Register goroutine
    go func(id int) {
        defer wg.Done()    // Mark complete when done
        // do work
    }(i)
}

wg.Wait()  // Block until all complete
```

---

## Common Mistakes I'll Avoid

- [ ] **Not waiting for goroutines**: Always use `time.Sleep()`, channels, or `WaitGroup`
- [ ] **Loop variable capture**: Pass loop var as parameter: `go func(i int){}(i)`
- [ ] **Race conditions**: Use mutex or channels when sharing data (Lesson 7)

---

## Connection to Channels (Lesson 4)

Now I understand WHY channels are needed:

**Goroutines** (this lesson):
- Independent workers running concurrently
- Need to communicate safely

**Channels** (Lesson 4):
- Safe way for goroutines to communicate
- Pipes connecting workers

Example:
```go
// Two goroutines
go sender()    // Needs to send data
go receiver()  // Needs to receive data

// How do they communicate?
// Answer: Channels! (Lesson 4)
ch := make(chan string)
```

---

## Aha Moments! 💡

What suddenly made sense:

1.

2.

3.

---

## Real Code I Now Understand

### Example 1: Client Pumps
```go
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
    conn, _ := upgrader.Upgrade(w, r, nil)
    client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}

    // I now understand these TWO goroutines!
    go client.writePump()  // Can wait for messages to send
    go client.readPump()   // Can wait for incoming messages
    // Both waiting simultaneously = both goroutines needed!
}
```

### Example 2: Hub Background
```go
func main() {
    hub := newHub()
    go hub.Run()  // I understand why this is in goroutine!
    // hub.Run() has infinite loop - would block forever
    // Goroutine lets main continue to start HTTP server

    http.ListenAndServe(":8080", nil)
}
```

---

## Questions to Review Later

1. How many goroutines can Go handle?


2. What's the difference between concurrency and parallelism?


3. When should I use a worker pool vs unlimited goroutines?


---

## Next Steps

- [ ] Review examples that were confusing
- [ ] Find all `go` keywords in websocket code
- [ ] Ready for Lesson 4: Channels!

---

**Ready for Lesson 4?** Check this when ready: [ ]

**Lesson 4 (Channels)** will show how goroutines communicate:
- `message := <-c.send` (receive from channel)
- `c.send <- data` (send to channel)
- Channels are pipes between goroutines!

After Lesson 6 + Lesson 4, the websocket code will be crystal clear!

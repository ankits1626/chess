# Lesson 4 Summary: Channels & `<-` Operator

**Complete this after finishing the lesson!**

## What I Learned

### Key Concepts
- [ ] What channels are (pipes between goroutines)
- [ ] Send operator: `ch <- value`
- [ ] Receive operator: `value := <-ch`
- [ ] Buffered vs unbuffered channels
- [ ] The `ok` pattern for detecting closure
- [ ] Closing channels with `close(ch)`

### The `<-` Operator - Now Clear!

**Before Lesson**: `message, ok := <-c.send` ← Confusing!

**After Lesson**:
- `<-c.send` receives FROM the channel
- `message` gets the value
- `ok` tells if channel is still open
- Arrow shows data flow direction!

## My Notes

Write your observations here:

1. **What clicked for you about channels?**


2. **What's the difference between `ch <- v` and `v := <-ch`?**


3. **When would you use a buffered channel?**


4. **Any remaining questions?**


---

## Examples I Ran

- [ ] 01-basic.go (send & receive)
- [ ] 02-buffered.go (buffered channels)
- [ ] 03-ok-pattern.go (closure detection)
- [ ] 04-worker-pool.go (multiple workers)
- [ ] 05-ping-pong.go (two-way communication)
- [ ] 06-websocket-simulation.go ⭐ (hub pattern)

**Most insightful example**: ___________________________

---

## Exercises Completed

- [ ] Exercise 1: Basic operations
- [ ] Exercise 2: Buffered channel
- [ ] Exercise 3: Closure detection
- [ ] Exercise 4: Range over channel
- [ ] Exercise 5: Worker pattern
- [ ] Exercise 6: Simple hub
- [ ] Bonus: Ping-pong

---

## Patterns Found in Chess Coach Code

### Pattern 1: Hub Registration
**File**: `_______________________`
**Line**: `_____`

```go
// Code with <- operator:


// What it does:

```

### Pattern 2: Hub Broadcast
**File**: `_______________________`
**Line**: `_____`

```go
// Code with <- operator:


// What it does:

```

### Pattern 3: Client Send Queue
**File**: `_______________________`
**Line**: `_____`

```go
// Code with <- operator:


// What it does:

```

---

## Arrow Direction Cheat Sheet

I understand that:
- [ ] `ch <- value` sends value INTO channel (arrow points right)
- [ ] `value := <-ch` receives FROM channel (arrow points left)
- [ ] Arrow shows the direction of data flow
- [ ] Sending blocks until someone receives
- [ ] Receiving blocks until someone sends

---

## Buffered vs Unbuffered

| Type | Created With | Behavior |
|------|-------------|----------|
| Unbuffered | `make(chan T)` | Blocks until receiver ready |
| Buffered | `make(chan T, N)` | Blocks only when buffer full |

**When I'd use buffered**: ___________________________

**When I'd use unbuffered**: ___________________________

---

## Questions for Review

1. Why does `ch <- value` without a goroutine cause deadlock?


2. What happens when you send to a closed channel?


3. What does the `ok` value tell you in `msg, ok := <-ch`?


4. How does the hub use channels to coordinate clients?


---

## Real Code I Now Understand

### Example 1: Hub Run Loop
I now understand this code from my websocket hub:

```go
select {
case client := <-h.register:
    // Receives client from register channel
    // Adds to clients map

case message := <-h.broadcast:
    // Receives broadcast message
    // Sends to all clients with:
    for client := range h.clients {
        client.send <- message  // ← I understand this!
    }
}
```

### Example 2: Client writePump
I now understand this code from my websocket client:

```go
case message, ok := <-c.send:
    // Receives from send channel (queue)
    // ok = false if hub closed the channel
    if !ok {
        return  // Client disconnected
    }
    // Send to actual websocket
    c.conn.WriteMessage(websocket.TextMessage, message)
```

---

## Aha Moments! 💡

What suddenly made sense:

1.

2.

3.

---

## Next Steps

- [ ] Review any confusing examples
- [ ] Find more `<-` patterns in websocket code
- [ ] Ready for Lesson 5: Select Statement

---

**Ready for Lesson 5?** Check this when ready: [ ]

Lesson 5 will cover: **Select Statement** - How `select { case ... }` chooses between multiple channels

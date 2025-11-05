# Lesson 6: Goroutines - The `go` Keyword

## Why This Lesson First?

You're seeing this in websocket code:
```go
go c.writePump()     // What does 'go' do?
go c.readPump()      // Why two of these?
go hub.Run()         // Why is this needed?
```

**You NEED to understand goroutines BEFORE channels make sense!**

Channels are pipes between goroutines. Without goroutines, channels are pointless.

---

## What Are Goroutines?

**Goroutines are lightweight threads managed by Go.**

Think of them as independent workers that can run at the same time:

```
Main Goroutine          Goroutine 1          Goroutine 2
      ║                      ║                     ║
   Print "A"             Print "X"            Print "1"
      ║                      ║                     ║
   Print "B"             Print "Y"            Print "2"
      ║                      ║                     ║
   Print "C"             Print "Z"            Print "3"
```

All three run **concurrently** (overlapping in time).

---

## The `go` Keyword

**Syntax**: `go functionName()`

**What it does**: Starts a new goroutine (lightweight thread)

### Simple Example

```go
package main

import (
    "fmt"
    "time"
)

func sayHello() {
    fmt.Println("Hello from goroutine!")
}

func main() {
    go sayHello()  // ← Starts new goroutine

    time.Sleep(100 * time.Millisecond)  // Wait for goroutine
    fmt.Println("Main function ending")
}
```

**Output**:
```
Hello from goroutine!
Main function ending
```

---

## How Goroutines Work

### Without `go` (Sequential)

```go
func main() {
    task1()  // Runs first, blocks until complete
    task2()  // Runs second, blocks until complete
    task3()  // Runs third
}
// Total time: task1 + task2 + task3
```

### With `go` (Concurrent)

```go
func main() {
    go task1()  // Starts in background
    go task2()  // Starts in background
    go task3()  // Starts in background

    time.Sleep(...)  // Wait for them
}
// Total time: max(task1, task2, task3) - they overlap!
```

---

## Why You Need Goroutines in Websockets

### Problem: Reading AND Writing Simultaneously

A websocket connection needs to:
1. **Read** messages from client (blocking operation)
2. **Write** messages to client (blocking operation)

**You can't do both in one thread!**

```go
// ❌ BAD - Can't work
for {
    msg := conn.ReadMessage()   // Blocks here waiting for message
    conn.WriteMessage(response) // Can't get here until message arrives!
}
```

### Solution: Two Goroutines

```go
// ✅ GOOD - Two independent goroutines
go readPump()   // Goroutine 1: continuously read from websocket
go writePump()  // Goroutine 2: continuously write to websocket

// Now they can both run at the same time!
```

---

## Creating Goroutines

### Method 1: Named Function

```go
func worker(id int) {
    fmt.Printf("Worker %d starting\n", id)
}

func main() {
    go worker(1)  // Start goroutine with argument
    go worker(2)
    go worker(3)
}
```

### Method 2: Anonymous Function

```go
func main() {
    go func() {
        fmt.Println("Anonymous goroutine")
    }()  // ← Note the () to call it!

    // With parameters
    go func(msg string) {
        fmt.Println(msg)
    }("Hello!")
}
```

---

## Goroutine Lifecycle

```go
func main() {
    fmt.Println("Main starts")

    go func() {
        fmt.Println("Goroutine running")
        time.Sleep(100 * time.Millisecond)
        fmt.Println("Goroutine done")
    }()

    fmt.Println("Main continues")
    time.Sleep(200 * time.Millisecond)  // Wait for goroutine
    fmt.Println("Main ends")
}
```

**Possible output** (order not guaranteed!):
```
Main starts
Main continues
Goroutine running
Goroutine done
Main ends
```

---

## Main Goroutine is Special

**When main() exits, ALL goroutines are killed!**

```go
func main() {
    go func() {
        time.Sleep(5 * time.Second)
        fmt.Println("This might not print!")
    }()

    // Main exits immediately - goroutine killed!
}
```

**Fix**: Wait for goroutines
```go
func main() {
    go func() {
        time.Sleep(1 * time.Second)
        fmt.Println("This will print!")
    }()

    time.Sleep(2 * time.Second)  // Wait longer than goroutine
}
```

---

## Goroutines vs Threads

| Goroutines | OS Threads |
|------------|------------|
| Lightweight (2KB stack) | Heavy (1-2MB stack) |
| Cheap to create (thousands) | Expensive (hundreds max) |
| Managed by Go runtime | Managed by OS |
| Fast context switching | Slow context switching |

**You can create 100,000 goroutines easily!**

---

## Common Patterns

### Pattern 1: Fire and Forget

```go
go logToFile(message)  // Don't wait for result
go sendEmail(user)     // Don't care when it completes
```

### Pattern 2: Multiple Workers

```go
for i := 1; i <= 10; i++ {
    go worker(i)  // Start 10 workers
}
```

### Pattern 3: Closure (Captures Variables)

```go
name := "Alice"
go func() {
    fmt.Println(name)  // Can access 'name' from outer scope
}()
```

**Warning**: Be careful with loop variables!

```go
// ❌ BAD
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i)  // All might print "3"!
    }()
}

// ✅ GOOD
for i := 0; i < 3; i++ {
    go func(id int) {
        fmt.Println(id)  // Each gets its own copy
    }(i)
}
```

---

## Waiting for Goroutines

### Method 1: Sleep (Simple but crude)

```go
go doWork()
time.Sleep(1 * time.Second)  // Hope 1 second is enough
```

### Method 2: Channel (Better - covered in Lesson 4)

```go
done := make(chan bool)
go func() {
    doWork()
    done <- true  // Signal completion
}()
<-done  // Wait for signal
```

### Method 3: WaitGroup (Best for multiple)

```go
var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)  // Register goroutine
    go func(id int) {
        defer wg.Done()  // Mark complete when done
        doWork(id)
    }(i)
}

wg.Wait()  // Wait for all goroutines
```

---

## Websocket Patterns

### Pattern 1: Client with Read/Write Pumps

```go
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
    conn, _ := upgrader.Upgrade(w, r, nil)

    client := &Client{
        conn: conn,
        send: make(chan []byte, 256),
    }

    // Start TWO goroutines for this client
    go client.readPump()   // Reads from websocket
    go client.writePump()  // Writes to websocket

    // Main function can continue or return
}
```

**Why two goroutines?**
- `readPump` blocks waiting for client messages
- `writePump` blocks waiting for messages to send
- They need to run **simultaneously**!

### Pattern 2: Hub Running Forever

```go
func main() {
    hub := newHub()
    go hub.Run()  // Start hub in background

    // Start HTTP server (also blocks)
    http.ListenAndServe(":8080", nil)
}
```

**Why goroutine for hub?**
- `hub.Run()` has infinite loop
- Would block forever if not in goroutine
- Main thread needs to start HTTP server

### Pattern 3: Broadcasting to Clients

```go
func (h *Hub) broadcastMessage(msg []byte) {
    for client := range h.clients {
        go func(c *Client) {
            c.send <- msg  // Each client gets message
        }(client)
    }
}
```

**Why goroutine per client?**
- If one client is slow, don't block others
- Each send can fail independently

---

## Race Conditions (Preview)

**Problem**: Multiple goroutines accessing same data

```go
var counter int

func increment() {
    counter++  // ❌ NOT SAFE with multiple goroutines!
}

func main() {
    for i := 0; i < 1000; i++ {
        go increment()
    }
    time.Sleep(1 * time.Second)
    fmt.Println(counter)  // Might not be 1000!
}
```

**Solutions** (covered in Lesson 7):
1. **Mutex**: Lock before accessing shared data
2. **Channels**: Communicate via channels instead

---

## Anonymous Goroutines in Websocket Code

You'll see this pattern:

```go
go func() {
    if err := a.server.Start(); err != nil {
        a.logger.Fatalf("Failed to start server: %v", err)
    }
}()
```

**Breakdown**:
- `go func() { ... }()` - Anonymous function as goroutine
- `()` at the end - Calls the function immediately
- Runs in background while main continues

---

## Goroutine Scheduling

Go runtime automatically schedules goroutines on CPU cores:

```
CPU Core 1          CPU Core 2          CPU Core 3
Goroutine 1         Goroutine 3         Goroutine 5
Goroutine 2         Goroutine 4         Goroutine 6
```

**You don't manage this!** Go runtime does it for you.

---

## Quick Reference

| Code | Meaning |
|------|---------|
| `go f()` | Start function `f` as goroutine |
| `go func() {}()` | Start anonymous function as goroutine |
| `time.Sleep(d)` | Pause current goroutine for duration |
| `runtime.NumGoroutine()` | Count active goroutines |
| `runtime.GOMAXPROCS(n)` | Set max CPU cores to use |

---

## Common Mistakes

### Mistake 1: Not Waiting for Goroutines

```go
func main() {
    go fmt.Println("Hello")
    // Main exits immediately, goroutine killed!
}
```

### Mistake 2: Loop Variable Capture

```go
for i := 0; i < 3; i++ {
    go func() {
        fmt.Println(i)  // ❌ All might print same value!
    }()
}
```

### Mistake 3: Too Many Goroutines

```go
for i := 0; i < 1000000; i++ {
    go expensiveTask()  // ❌ Might exhaust resources
}
```

**Better**: Use worker pool pattern

---

## Exercises

I've created 6 progressive examples:

1. **Basic goroutine** - Single `go` keyword
2. **Multiple goroutines** - Several running concurrently
3. **Anonymous functions** - Inline goroutines
4. **WaitGroup** - Proper synchronization
5. **Worker pool** - Controlled concurrency
6. **Websocket pattern** - readPump + writePump simulation

---

## Next Steps

1. **Run examples** in `examples/` folder
2. **Complete exercises** in `exercises/` folder
3. **See goroutines** in your websocket code
4. **Then** move to Lesson 4 (Channels) - it will make sense!

---

## Connection to Channels

**After understanding goroutines**, channels make perfect sense:

```go
// Two goroutines need to communicate
go sender()    // Sends data
go receiver()  // Receives data

// How do they communicate safely?
// Answer: CHANNELS! (Lesson 4)
```

**Channels are pipes between goroutines.**

Without goroutines → No need for channels
With goroutines → Channels are essential!

---

## What's Next?

**Recommended order**:
1. ✅ Lesson 1: Variables
2. **→ Lesson 6: Goroutines** (this lesson)
3. Lesson 4: Channels (makes sense now!)
4. Lesson 5: Select Statement
5. Lesson 7: Mutexes

After this lesson + Lesson 4, you'll understand:
```go
go c.writePump()              // Goroutine (this lesson)
message := <-c.send           // Channel (Lesson 4)
go hub.Run()                  // Goroutine (this lesson)
hub.broadcast <- msg          // Channel (Lesson 4)
```

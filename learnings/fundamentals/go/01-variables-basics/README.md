# Lesson 1: Variables & Basic Syntax

## Why This Lesson?

You're seeing code like this in websockets and it's confusing:
```go
message, ok := <-c.send        // What is :=?
var conn *websocket.Conn       // What is var?
messageType, message, err := conn.ReadMessage()  // Multiple variables?
```

By the end of this lesson, you'll understand **how Go creates and uses variables**.

---

## Core Concepts

### 1. Go is Statically Typed

Unlike JavaScript or Python, Go requires you to specify types:
```javascript
// JavaScript - type changes freely
let x = 5;      // x is a number
x = "hello";    // now x is a string (OK!)
```

```go
// Go - type is fixed
var x int = 5
x = "hello"  // ERROR! x is an int, can't assign string
```

### 2. Two Ways to Declare Variables

#### Method 1: `var` (Explicit Declaration)

**Syntax**: `var name type`

```go
var username string
var age int
var isActive bool
```

**When to use**:
- When you want to be explicit about the type
- When declaring without initial value
- At package level (outside functions)

#### Method 2: `:=` (Short Declaration)

**Syntax**: `name := value`

```go
username := "alice"    // Go knows it's a string
age := 25              // Go knows it's an int
isActive := true       // Go knows it's a bool
```

**When to use**:
- Inside functions (most common!)
- When initial value makes type obvious
- Go figures out the type automatically

---

## Key Rules

### Rule 1: `:=` Only Works Inside Functions

```go
// ❌ WRONG - outside function
package main
message := "hello"  // ERROR!

// ✅ CORRECT
package main
var message = "hello"  // OK at package level

func main() {
    msg := "hello"  // OK inside function
}
```

### Rule 2: `:=` Creates New Variables

```go
// First time - creates the variable
count := 0

// Can't use := again for same variable
count := 1  // ERROR! count already exists

// Use = for assignment
count = 1   // OK - assigns new value
```

### Rule 3: Multiple Assignment

```go
// Create multiple variables at once
name, age := "alice", 25

// Or with var
var x, y int = 10, 20

// Common pattern: function returns multiple values
result, err := doSomething()
```

**Special case**: `:=` can reuse variables if at least ONE is new:
```go
name, age := "alice", 25    // Both new
name, score := "bob", 100   // name exists, but score is new - OK!
```

---

## Common Patterns in Chess Coach

### Pattern 1: Websocket Upgrade
```go
// From your backend code
conn, err := upgrader.Upgrade(w, r, nil)
```

**Breakdown**:
- `:=` creates TWO variables: `conn` and `err`
- `upgrader.Upgrade()` returns two values
- `conn` is the websocket connection
- `err` is the error (if any)

### Pattern 2: Reading Messages
```go
messageType, message, err := conn.ReadMessage()
```

**Breakdown**:
- `:=` creates THREE variables
- `messageType` - int (1=text, 2=binary, etc.)
- `message` - []byte (the actual message)
- `err` - error (if read failed)

### Pattern 3: Channel Receive
```go
message, ok := <-c.send
```

**Breakdown**:
- `:=` creates TWO variables
- `<-c.send` receives from channel
- `message` - the value from channel
- `ok` - bool (false if channel is closed)

---

## Zero Values

If you declare without initializing, Go gives **zero values**:

```go
var count int        // 0
var name string      // "" (empty string)
var isActive bool    // false
var ptr *int         // nil
```

---

## Type Inference with `:=`

Go is smart about figuring out types:

```go
x := 42              // int
y := 3.14            // float64
name := "alice"      // string
isActive := true     // bool

// Be careful with numbers!
a := 5       // int
b := 5.0     // float64 (has decimal point)
```

---

## Examples to Run

I've created 5 progressive examples in the `examples/` folder:

1. **01-basic.go** - Simple var declarations
2. **02-short-declaration.go** - Using :=
3. **03-multiple.go** - Multiple assignments
4. **04-zero-values.go** - Default values
5. **05-websocket-pattern.go** - Real websocket examples

---

## Quick Reference

| Syntax | Usage | Example |
|--------|-------|---------|
| `var x int` | Declare with type, zero value | `var count int` (count = 0) |
| `var x = 5` | Declare with value, type inferred | `var count = 5` (count is int) |
| `x := 5` | Short declaration (inside functions) | `count := 5` |
| `x, y := 1, 2` | Multiple declaration | `name, age := "alice", 25` |
| `var x int = 5` | Explicit type and value | `var count int = 5` |

---

## Next Steps

1. **Run Examples**: Go to `examples/` folder and run each file
   ```bash
   cd examples
   go run 01-basic.go
   go run 02-short-declaration.go
   # ... etc
   ```

2. **Complete Exercise**: Try the exercise in `exercises/practice.go`

3. **Ask Questions**: Anything unclear? Ask before moving on!

4. **Create Summary**: After completing, write your understanding in `SUMMARY.md`

---

## What's Next?

**Lesson 2**: Pointers (`&` and `*`) - Understanding `&Client` and pointer receivers

The confusing syntax builds on what you learned here!

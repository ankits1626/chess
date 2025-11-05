# References to Chess Coach Backend

Find these patterns in your actual codebase to reinforce learning!

## Pattern 1: Short Declaration with Error Handling

### Example in Backend:
**File**: `backend/internal/database/connection.go`

```go
pool, err := pgxpool.New(ctx, dsn)
if err != nil {
    return nil, fmt.Errorf("failed to create connection pool: %w", err)
}
```

**What's happening**:
- `:=` creates TWO variables: `pool` and `err`
- `pgxpool.New()` returns `(*pgxpool.Pool, error)`
- Always check `err` before using `pool`

---

### Example in Handlers:
**File**: `backend/internal/handler/v1/user/handler.go`

```go
var req CreateRequest
if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}
```

**What's happening**:
- `var req CreateRequest` declares the struct
- `:=` creates `err` in the if statement scope
- If binding fails, return error immediately (early return pattern)

---

## Pattern 2: Multiple Return Values

### Example in Repository:
**File**: `backend/internal/repository/user_repository.go`

```go
user, err := r.db.GetUser(ctx, pgtype.UUID{Bytes: id, Valid: true})
if err != nil {
    return database.User{}, err
}
return user, nil
```

**What's happening**:
- `GetUser` returns `(database.User, error)`
- `:=` creates both `user` and `err`
- Return zero value `database.User{}` on error

---

### Example in Websocket:
**File**: `learnings/fundamentals/websocket/lesson-01-basics/main.go`

```go
messageType, message, err := conn.ReadMessage()
if err != nil {
    log.Printf("Error reading message: %v", err)
    break
}
```

**What's happening**:
- `ReadMessage()` returns THREE values
- `:=` creates `messageType`, `message`, and `err`
- Break the loop if error occurs

---

## Pattern 3: Reusing Variables

### Example in Handler:
**File**: `backend/internal/handler/v1/user/handler.go`

```go
id, err := uuid.Parse(c.Param("id"))
if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
    return
}

user, err := h.repo.GetByID(c.Request.Context(), id)
// ↑ 'err' is reused, 'user' is new
```

**What's happening**:
- First `:=` creates `id` and `err`
- Second `:=` creates `user` and reuses `err` (legal because `user` is new!)
- Very common pattern in Go

---

## Pattern 4: Zero Values

### Example in Structs:
**File**: `backend/internal/handler/v1/user/dto.go`

```go
type CreateRequest struct {
    Username string `json:"username" binding:"required"`
    Email    string `json:"email" binding:"required"`
    Rating   int32  `json:"rating,omitempty"` // Can be 0 (zero value)
}
```

**What's happening**:
- If client doesn't send `rating`, it's 0 (zero value for int32)
- `omitempty` means don't include in JSON if zero value
- Zero values are useful for optional fields

---

## Pattern 5: Package-Level Variables

### Example in Main:
**File**: `learnings/fundamentals/websocket/lesson-01-basics/main.go`

```go
var upgrader = websocket.Upgrader{
    ReadBufferSize:  readBufferSize,
    WriteBufferSize: writeBufferSize,
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}
```

**What's happening**:
- `var` used at package level (outside functions)
- `:=` would NOT work here (only inside functions)
- `upgrader` is available to all functions in the package

---

## Pattern 6: Ignoring Values with `_`

### Common Pattern:
```go
// Only care about error, not the result
_, err := somethingThatReturnsTwoValues()

// Only care about the value, not if channel is closed
message := <-c.send  // Don't need the 'ok' value
```

**When to use**:
- Use `_` when you don't need a return value
- Compiler won't complain about unused variables

---

## Exercises

1. **Find 3 examples** of `:=` in your backend code
2. **Find 2 examples** of `var` at package level
3. **Find 1 example** of multiple return values with error handling

Write them in your SUMMARY.md!

---

## Next Lesson Preview

In Lesson 2, you'll see patterns like:
```go
client := &Client{hub: hub, conn: conn}  // What is &?
func (c *Client) readPump() {}           // What is *?
```

These use **pointers** - we'll demystify them next!

# Go Language Fundamentals

**A comprehensive guide to understanding Go semantics and best practices**

---

## Table of Contents
1. [Basic Syntax & Types](#basic-syntax--types)
2. [Pointers vs Values](#pointers-vs-values)
3. [Structs](#structs)
4. [Interfaces](#interfaces)
5. [Methods & Receivers](#methods--receivers)
6. [Error Handling](#error-handling)
7. [Packages & Imports](#packages--imports)
8. [Goroutines & Concurrency](#goroutines--concurrency)
9. [Channels](#channels)
10. [Defer, Panic, Recover](#defer-panic-recover)
11. [Context](#context)
12. [Common Patterns](#common-patterns)

---

## Basic Syntax & Types

### Variable Declaration

```go
// Explicit type
var name string = "John"

// Type inference
var age = 30

// Short declaration (inside functions only)
email := "john@example.com"

// Multiple variables
var (
    host string = "localhost"
    port int    = 8080
)

// Zero values (default values)
var count int        // 0
var enabled bool     // false
var message string   // ""
var ptr *int         // nil
```

### Basic Types

```go
// Numeric types
var i int = 42           // Platform-dependent (32 or 64 bit)
var i8 int8 = 127        // -128 to 127
var i16 int16            // -32768 to 32767
var i32 int32            // -2^31 to 2^31-1
var i64 int64            // -2^63 to 2^63-1

var u uint = 42          // Unsigned, platform-dependent
var u8 uint8             // 0 to 255 (also called byte)
var u32 uint32           // 0 to 2^32-1

var f32 float32 = 3.14
var f64 float64 = 3.14159265

// String
var str string = "Hello, World!"

// Boolean
var isActive bool = true

// Rune (Unicode code point, alias for int32)
var char rune = 'A'
```

### Constants

```go
// Constants cannot be changed
const Pi = 3.14159
const (
    StatusOK       = 200
    StatusNotFound = 404
    StatusError    = 500
)

// iota: auto-incrementing constant
const (
    Sunday = iota    // 0
    Monday           // 1
    Tuesday          // 2
    Wednesday        // 3
)
```

---

## Pointers vs Values

**Key Concept**: Understanding when to use pointers (`*T`) vs values (`T`) is critical in Go.

### Value Semantics

```go
type User struct {
    Name string
    Age  int
}

// Pass by VALUE - creates a copy
func updateAge(u User, newAge int) {
    u.Age = newAge  // This modifies the COPY, not the original
}

func main() {
    user := User{Name: "John", Age: 30}
    updateAge(user, 31)
    fmt.Println(user.Age)  // Still 30! Original unchanged
}
```

### Pointer Semantics

```go
// Pass by POINTER - shares the original
func updateAge(u *User, newAge int) {
    u.Age = newAge  // This modifies the ORIGINAL
}

func main() {
    user := User{Name: "John", Age: 30}
    updateAge(&user, 31)  // &user creates pointer to user
    fmt.Println(user.Age)  // Now 31! Original changed
}
```

### When to Use Pointers

✅ **Use Pointers When**:
- You need to modify the original value
- The struct is large (avoid expensive copies)
- You need to represent "absence" (nil pointer)
- Working with methods that modify state

❌ **Use Values When**:
- The data is small (< 64 bytes)
- The data is immutable
- Simplicity is preferred
- Working with basic types (int, string, bool)

### Pointer Example from Your Code

```go
// internal/repository/user_repository.go
type UserRepository struct {
    db *database.DB  // POINTER: Repository shares DB connection
}

// Returns concrete instance (not pointer)
func NewUserRepository(db *database.DB) *UserRepository {
    return &UserRepository{db: db}  // &UserRepository creates pointer
}

// Method receiver is pointer (can modify repository)
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (database.User, error) {
    // Returns VALUE (not pointer) because User is copied from DB
    return r.db.GetUser(ctx, pgtype.UUID{Bytes: id, Valid: true})
}
```

**Why these choices?**
- `db *database.DB` - Pointer because DB connection is shared across app
- `*UserRepository` - Pointer because repositories are stateful objects
- `database.User` - Value because it's a data snapshot from DB

---

## Structs

Structs are custom types that group related data.

### Struct Definition

```go
// Basic struct
type User struct {
    ID       uuid.UUID
    Username string
    Email    string
    Age      int
}

// Struct with tags (used for JSON, DB mapping)
type UserJSON struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email" validate:"email"`
    Age      int    `json:"age,omitempty"`  // omitempty: skip if zero value
}

// Anonymous/embedded struct
var config = struct {
    Host string
    Port int
}{
    Host: "localhost",
    Port: 8080,
}
```

### Creating Structs

```go
// Method 1: Field names (preferred)
user := User{
    ID:       uuid.New(),
    Username: "john",
    Email:    "john@example.com",
    Age:      30,
}

// Method 2: Positional (NOT recommended - fragile)
user := User{uuid.New(), "john", "john@example.com", 30}

// Method 3: Zero value then assign
var user User
user.Username = "john"
user.Email = "john@example.com"

// Method 4: Pointer with new
user := new(User)  // Returns *User with zero values
user.Username = "john"
```

### Struct Embedding (Composition)

```go
// Base struct
type Person struct {
    Name string
    Age  int
}

// Employee "embeds" Person (composition, not inheritance)
type Employee struct {
    Person      // Embedded struct
    EmployeeID  string
    Department  string
}

func main() {
    emp := Employee{
        Person: Person{
            Name: "John",
            Age:  30,
        },
        EmployeeID: "E123",
        Department: "Engineering",
    }

    // Can access embedded fields directly
    fmt.Println(emp.Name)  // "John" (promoted from Person)
    fmt.Println(emp.Person.Name)  // Also valid
}
```

---

## Interfaces

**Key Concept**: Interfaces in Go are IMPLICIT. You don't declare "implements" - if a type has the methods, it satisfies the interface.

### Interface Definition

```go
// Interface: defines behavior (methods)
type Logger interface {
    Info(msg string)
    Error(msg string)
}

// Any type with these methods automatically satisfies Logger
type StdLogger struct {}

func (l *StdLogger) Info(msg string) {
    fmt.Println("INFO:", msg)
}

func (l *StdLogger) Error(msg string) {
    fmt.Println("ERROR:", msg)
}

// StdLogger now implements Logger (implicitly!)
```

### Empty Interface

```go
// interface{} or any - can hold ANY type
var data interface{}
data = 42
data = "hello"
data = User{Name: "John"}

// Type assertion
str, ok := data.(string)  // ok is true if data is string
if ok {
    fmt.Println("String:", str)
}

// Type switch
switch v := data.(type) {
case int:
    fmt.Println("Integer:", v)
case string:
    fmt.Println("String:", v)
default:
    fmt.Println("Unknown type")
}
```

### Interface Example from Your Code

```go
// internal/server/server.go
type Server interface {
    Start() error
    Shutdown(ctx context.Context) error
}

// Concrete implementation
type server struct {
    httpServer *http.Server
    config     *config.Config
    db         *database.DB
}

// These methods make server satisfy Server interface
func (s *server) Start() error {
    return s.httpServer.ListenAndServe()
}

func (s *server) Shutdown(ctx context.Context) error {
    return s.httpServer.Shutdown(ctx)
}

// main.go uses the interface, not the concrete type
func main() {
    var srv Server = server.New(cfg, db)  // Polymorphism!
}
```

**Benefits**:
- Testability (can create mock Server)
- Decoupling (main.go doesn't know about server internals)
- Flexibility (can swap implementations)

---

## Methods & Receivers

Methods are functions attached to types.

### Value Receiver

```go
type Rectangle struct {
    Width  float64
    Height float64
}

// Value receiver: operates on a COPY
func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    area := rect.Area()  // rect is NOT modified
}
```

### Pointer Receiver

```go
// Pointer receiver: operates on the ORIGINAL
func (r *Rectangle) Scale(factor float64) {
    r.Width *= factor   // Modifies original
    r.Height *= factor  // Modifies original
}

func main() {
    rect := Rectangle{Width: 10, Height: 5}
    rect.Scale(2)  // rect is now {20, 10}
}
```

### When to Use Pointer Receivers

✅ **Use Pointer Receiver (`*T`) When**:
- Method modifies the receiver
- Receiver is large (avoid copying)
- Consistency: if ANY method uses `*T`, use `*T` for ALL methods

❌ **Use Value Receiver (`T`) When**:
- Method doesn't modify receiver
- Receiver is small (int, bool, small struct)
- Receiver is a map, func, or chan (already reference types)

### Example from Your Code

```go
// internal/repository/user_repository.go
type UserRepository struct {
    db *database.DB
}

// Pointer receiver - repository methods use *UserRepository consistently
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (database.User, error) {
    return r.db.GetUser(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *UserRepository) Create(ctx context.Context, username, email string, rating int32) (database.User, error) {
    return r.db.CreateUser(ctx, database.CreateUserParams{
        Username: username,
        Email:    email,
        Rating:   pgtype.Int4{Int32: rating, Valid: true},
    })
}
```

**Why pointer receiver?**
- Consistency (all methods use `*UserRepository`)
- Future-proofing (if we add state like cache, can modify it)
- Convention for repository pattern

---

## Error Handling

Go doesn't have exceptions. Errors are VALUES.

### Basic Error Handling

```go
// Functions return error as last return value
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

func main() {
    result, err := divide(10, 0)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("Result:", result)
}
```

### Creating Errors

```go
import (
    "errors"
    "fmt"
)

// Method 1: errors.New
err := errors.New("something went wrong")

// Method 2: fmt.Errorf (with formatting)
err := fmt.Errorf("user %s not found", username)

// Method 3: Wrap errors (Go 1.13+)
err := fmt.Errorf("failed to get user: %w", originalErr)

// Method 4: Custom error type
type ValidationError struct {
    Field string
    Err   error
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %v", e.Field, e.Err)
}
```

### Error Handling Patterns

```go
// Pattern 1: Early return
func GetUser(id string) (*User, error) {
    if id == "" {
        return nil, errors.New("id cannot be empty")
    }

    user, err := db.Find(id)
    if err != nil {
        return nil, fmt.Errorf("failed to find user: %w", err)
    }

    return user, nil
}

// Pattern 2: Multiple error checks
func ProcessOrder(orderID string) error {
    order, err := getOrder(orderID)
    if err != nil {
        return err
    }

    if err := validateOrder(order); err != nil {
        return err
    }

    if err := chargePayment(order); err != nil {
        return err
    }

    return nil
}
```

### Example from Your Code

```go
// internal/router/v1/user_handler.go
func (h *UserHandler) Create(c *gin.Context) {
    var req CreateUserRequest

    // Error handling: validation
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return  // Early return on error
    }

    // Error handling: business logic
    user, err := h.repo.Create(c.Request.Context(), req.Username, req.Email, req.Rating)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
        return  // Early return on error
    }

    c.JSON(http.StatusCreated, ToUserResponse(user))
}
```

---

## Packages & Imports

### Package Declaration

```go
// Every .go file starts with package declaration
package main  // Executable package (has main function)
package user  // Library package
```

### Import Statements

```go
// Single import
import "fmt"

// Multiple imports (preferred)
import (
    "context"
    "fmt"
    "log"
)

// Aliased import
import (
    stdlog "log"  // Rename to avoid conflict
    "github.com/sirupsen/logrus"
)

// Blank import (for side effects only)
import _ "github.com/lib/pq"  // Registers PostgreSQL driver
```

### Package Visibility

```go
// EXPORTED (public): starts with capital letter
type User struct {
    ID   string  // Exported field
    Name string  // Exported field
}

func NewUser() *User {  // Exported function
    return &User{}
}

// UNEXPORTED (private): starts with lowercase letter
type config struct {
    host string  // Unexported field
}

func loadConfig() *config {  // Unexported function
    return &config{}
}
```

### Package Structure Example

```
backend/
├── cmd/
│   └── server/
│       └── main.go           # package main
├── internal/                 # internal packages (not importable outside project)
│   ├── config/
│   │   └── config.go         # package config
│   ├── database/
│   │   ├── connection.go     # package database
│   │   └── db.go             # package database (generated)
│   ├── repository/
│   │   └── user_repository.go # package repository
│   └── router/
│       └── v1/
│           └── user_handler.go # package v1
└── go.mod
```

---

## Goroutines & Concurrency

Goroutines are lightweight threads managed by Go runtime.

### Creating Goroutines

```go
// Sequential execution
func main() {
    doWork()
    doMoreWork()
}

// Concurrent execution
func main() {
    go doWork()       // Runs in background
    go doMoreWork()   // Runs in background

    time.Sleep(time.Second)  // Wait for goroutines (BAD practice)
}
```

### WaitGroup (Proper Synchronization)

```go
import "sync"

func main() {
    var wg sync.WaitGroup

    wg.Add(2)  // We're launching 2 goroutines

    go func() {
        defer wg.Done()  // Signal completion
        doWork()
    }()

    go func() {
        defer wg.Done()
        doMoreWork()
    }()

    wg.Wait()  // Wait for all goroutines to complete
}
```

### Example from Your Code

```go
// cmd/server/main.go
func main() {
    srv := server.New(cfg, db)

    // Start server in background goroutine
    go func() {
        if err := srv.Start(); err != nil {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    // Main goroutine waits for shutdown signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit  // Block until signal received
}
```

---

## Channels

Channels are pipes for communication between goroutines.

### Creating and Using Channels

```go
// Create channel
ch := make(chan int)          // Unbuffered channel
ch := make(chan int, 10)      // Buffered channel (capacity 10)

// Send to channel
ch <- 42

// Receive from channel
value := <-ch

// Close channel
close(ch)
```

### Channel Patterns

```go
// Pattern 1: Worker pool
func worker(jobs <-chan int, results chan<- int) {
    for job := range jobs {  // Receive until channel closed
        results <- job * 2
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)

    // Start 3 workers
    for w := 1; w <= 3; w++ {
        go worker(jobs, results)
    }

    // Send jobs
    for j := 1; j <= 9; j++ {
        jobs <- j
    }
    close(jobs)

    // Collect results
    for r := 1; r <= 9; r++ {
        fmt.Println(<-results)
    }
}

// Pattern 2: Select (multiplex channels)
func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "one"
    }()

    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "two"
    }()

    for i := 0; i < 2; i++ {
        select {
        case msg1 := <-ch1:
            fmt.Println("Received", msg1)
        case msg2 := <-ch2:
            fmt.Println("Received", msg2)
        }
    }
}
```

---

## Defer, Panic, Recover

### Defer

```go
// Defer: executes function when surrounding function returns
func main() {
    defer fmt.Println("Third")   // Executes last
    defer fmt.Println("Second")  // Executes second
    fmt.Println("First")         // Executes first
}
// Output: First, Second, Third

// Common use: cleanup
func processFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()  // Ensures file closes even if error occurs

    // Process file...
    return nil
}
```

### Panic and Recover

```go
// Panic: unrecoverable error (like exception)
func divide(a, b int) int {
    if b == 0 {
        panic("division by zero")  // Crashes program
    }
    return a / b
}

// Recover: catch panic (like try-catch)
func safeDivide(a, b int) (result int) {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
            result = 0
        }
    }()

    result = divide(a, b)
    return result
}
```

### Example from Your Code

```go
// cmd/server/main.go
func main() {
    db, err := database.NewDB(ctx, cfg.DSN())
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()  // Ensures DB closes when main exits

    application := app.New(cfg, db, appLogger)
    defer application.Close()  // Ensures cleanup on exit

    // ...
}
```

---

## Context

Context carries deadlines, cancellation signals, and request-scoped values.

### Context Basics

```go
import "context"

// Create contexts
ctx := context.Background()  // Root context (never canceled)
ctx := context.TODO()        // Placeholder when unsure which context to use

// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()  // Always call cancel to release resources

// With deadline
deadline := time.Now().Add(10 * time.Second)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()

// With cancellation
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
```

### Context Usage

```go
// Passing context through call chain
func ProcessRequest(ctx context.Context, userID string) error {
    user, err := GetUser(ctx, userID)
    if err != nil {
        return err
    }

    return SaveUser(ctx, user)
}

func GetUser(ctx context.Context, userID string) (*User, error) {
    // Check if context canceled
    select {
    case <-ctx.Done():
        return nil, ctx.Err()  // context.Canceled or context.DeadlineExceeded
    default:
        // Proceed with work
    }

    // Database query with context
    return db.QueryUser(ctx, userID)
}
```

### Example from Your Code

```go
// internal/repository/user_repository.go
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (database.User, error) {
    // Context passed to database query
    // If request canceled/timeout, query is aborted
    return r.db.GetUser(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

// internal/router/v1/user_handler.go
func (h *UserHandler) Create(c *gin.Context) {
    // Extract request context from Gin
    ctx := c.Request.Context()

    // Pass context down the call chain
    user, err := h.repo.Create(ctx, req.Username, req.Email, req.Rating)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
        return
    }

    c.JSON(http.StatusCreated, ToUserResponse(user))
}
```

---

## Common Patterns

### Constructor Pattern

```go
// Factory function: creates and initializes struct
type UserRepository struct {
    db *database.DB
}

// New* prefix is idiomatic for constructors
func NewUserRepository(db *database.DB) *UserRepository {
    return &UserRepository{db: db}
}
```

### Options Pattern (Functional Options)

```go
// For structs with many optional parameters
type Server struct {
    host    string
    port    int
    timeout time.Duration
}

type Option func(*Server)

func WithHost(host string) Option {
    return func(s *Server) {
        s.host = host
    }
}

func WithTimeout(timeout time.Duration) Option {
    return func(s *Server) {
        s.timeout = timeout
    }
}

func NewServer(port int, opts ...Option) *Server {
    s := &Server{
        host:    "localhost",  // Default
        port:    port,
        timeout: 30 * time.Second,  // Default
    }

    for _, opt := range opts {
        opt(s)
    }

    return s
}

// Usage
srv := NewServer(8080, WithHost("0.0.0.0"), WithTimeout(60*time.Second))
```

### Builder Pattern

```go
type QueryBuilder struct {
    table  string
    where  []string
    limit  int
    offset int
}

func NewQueryBuilder(table string) *QueryBuilder {
    return &QueryBuilder{table: table}
}

func (b *QueryBuilder) Where(condition string) *QueryBuilder {
    b.where = append(b.where, condition)
    return b  // Return self for chaining
}

func (b *QueryBuilder) Limit(limit int) *QueryBuilder {
    b.limit = limit
    return b
}

func (b *QueryBuilder) Build() string {
    // Build SQL query
    return fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT %d",
        b.table, strings.Join(b.where, " AND "), b.limit)
}

// Usage (method chaining)
query := NewQueryBuilder("users").
    Where("age > 18").
    Where("active = true").
    Limit(10).
    Build()
```

### Repository Pattern (from your code)

```go
// Abstraction layer between business logic and data access
type UserRepository struct {
    db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (database.User, error) {
    return r.db.GetUser(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *UserRepository) Create(ctx context.Context, username, email string, rating int32) (database.User, error) {
    return r.db.CreateUser(ctx, database.CreateUserParams{
        Username: username,
        Email:    email,
        Rating:   pgtype.Int4{Int32: rating, Valid: true},
    })
}
```

**Benefits**:
- Separates business logic from data access
- Easier to test (can mock repository)
- Centralized data access logic

---

## Key Takeaways

### 1. **Pointers vs Values**
- Use pointers for large structs or when you need to modify
- Use values for small, immutable data
- Be consistent within a type's methods

### 2. **Interfaces**
- Define behavior, not implementation
- Keep interfaces small (1-3 methods)
- Implicit satisfaction (no "implements" keyword)

### 3. **Error Handling**
- Errors are values, not exceptions
- Check errors immediately: `if err != nil`
- Wrap errors with context: `fmt.Errorf("context: %w", err)`

### 4. **Concurrency**
- "Don't communicate by sharing memory; share memory by communicating"
- Use channels for communication
- Use mutexes for shared state protection
- Always call `defer wg.Done()` and `defer cancel()`

### 5. **Context**
- Always first parameter: `func Do(ctx context.Context, ...)`
- Pass down the call chain
- Use for cancellation, timeouts, request-scoped values

### 6. **Package Design**
- Exported (public): starts with uppercase
- Unexported (private): starts with lowercase
- Use `internal/` for private packages

### 7. **Idiomatic Go**
- Simple is better than clever
- Explicit is better than implicit
- Errors are values
- Composition over inheritance
- Accept interfaces, return structs

---

## Further Reading

- **Official Tour**: https://go.dev/tour/
- **Effective Go**: https://go.dev/doc/effective_go
- **Go Proverbs**: https://go-proverbs.github.io/
- **Go by Example**: https://gobyexample.com/
- **Dave Cheney's Blog**: https://dave.cheney.net/

---

**Last Updated**: 2025-11-01

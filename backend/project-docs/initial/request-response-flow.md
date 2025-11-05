# Request-Response Flow in Chess Coach Backend

**Understanding how HTTP requests flow through the application**

---

## Table of Contents
1. [High-Level Overview](#high-level-overview)
2. [Step-by-Step Flow](#step-by-step-flow)
3. [Example: Create User Request](#example-create-user-request)
4. [Layer Responsibilities](#layer-responsibilities)
5. [Complete Code Walkthrough](#complete-code-walkthrough)
6. [Error Flow](#error-flow)
7. [Visual Diagrams](#visual-diagrams)

---

## High-Level Overview

```
Client Request
     ↓
[1] main.go - Application Entry Point
     ↓
[2] server.go - HTTP Server Setup
     ↓
[3] router.go - Route Registration
     ↓
[4] v1_routes.go - Version-Specific Routes
     ↓
[5] handler.go - HTTP Request Handler
     ↓
[6] repository.go - Data Access Layer
     ↓
[7] database (PostgreSQL)
     ↓
[6] repository.go - Returns Data
     ↓
[5] handler.go - Converts to DTO
     ↓
Client Response (JSON)
```

---

## Step-by-Step Flow

### 1️⃣ Application Starts

**File**: `cmd/server/main.go`

```go
func main() {
    // Load configuration (DB credentials, port, etc.)
    cfg := config.Load()

    // Initialize logger
    appLogger := logger.NewStdLogger()

    // Connect to database
    ctx := context.Background()
    db, err := database.NewDB(ctx, cfg.DSN())
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Create application (server + routes)
    application := app.New(cfg, db, appLogger)
    defer application.Close()

    // Start server (blocks until shutdown signal)
    if err := application.Run(ctx); err != nil {
        log.Fatalf("Application error: %v", err)
    }
}
```

**What happens:**
- Configuration loaded from environment variables
- Database connection established
- HTTP server starts listening on port 8080
- Waits for incoming HTTP requests

---

### 2️⃣ HTTP Server Setup

**File**: `internal/app/app.go`

```go
type App struct {
    config *config.Config
    db     *database.DB
    server server.Server  // HTTP server
    logger logger.Logger
}

func New(cfg *config.Config, db *database.DB, log logger.Logger) *App {
    return &App{
        config: cfg,
        db:     db,
        server: server.New(cfg, db),  // Creates Gin router
        logger: log,
    }
}

func (a *App) Run(ctx context.Context) error {
    go func() {
        if err := a.server.Start(); err != nil {
            a.logger.Fatalf("Failed to start server: %v", err)
        }
    }()

    // ... wait for shutdown signal
}
```

**File**: `internal/server/server.go`

```go
func New(cfg *config.Config, db *database.DB) Server {
    gin.SetMode(cfg.GinMode)
    r := router.Setup(db)  // Setup routes

    httpServer := &http.Server{
        Addr:    fmt.Sprintf(":%s", cfg.Port),
        Handler: r,  // Gin router
    }

    return &server{
        httpServer: httpServer,
        config:     cfg,
        db:         db,
    }
}
```

**What happens:**
- Gin web framework initialized
- Routes registered
- HTTP server listens on `:8080`

---

### 3️⃣ Route Registration

**File**: `internal/router/router.go`

```go
func Setup(db *database.DB) *gin.Engine {
    r := gin.Default()  // Create Gin router

    // Swagger endpoint
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Register API v1 routes
    RegisterV1Routes(r, db)

    return r
}
```

**File**: `internal/router/v1_routes.go`

```go
func RegisterV1Routes(r *gin.Engine, db *database.DB) {
    v1 := r.Group("/api/v1")  // All routes start with /api/v1
    {
        // Health check
        v1.GET("/health", health.Check)

        // User routes
        userRepo := repository.NewUserRepository(db)
        userHandler := user.NewHandler(userRepo)

        users := v1.Group("/users")
        {
            users.GET("", userHandler.List)           // GET /api/v1/users
            users.GET("/:id", userHandler.Get)        // GET /api/v1/users/123
            users.POST("", userHandler.Create)        // POST /api/v1/users
            users.PUT("/:id", userHandler.Update)     // PUT /api/v1/users/123
            users.DELETE("/:id", userHandler.Delete)  // DELETE /api/v1/users/123
        }
    }
}
```

**What happens:**
- Routes mapped to handler functions
- Dependencies injected (repository → handler)
- Pattern: `HTTP Method + URL Path → Handler Function`

---

### 4️⃣ Request Arrives

**Example Request:**
```bash
POST http://localhost:8080/api/v1/users
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "rating": 1500
}
```

**Gin Router:**
1. Receives HTTP request
2. Matches route: `POST /api/v1/users`
3. Calls handler: `userHandler.Create(c *gin.Context)`

---

### 5️⃣ Handler Processes Request

**File**: `internal/handler/v1/user/handler.go`

```go
// Step 1: Handler receives Gin context
func (h *Handler) Create(c *gin.Context) {
    var req CreateRequest

    // Step 2: Parse JSON body into Go struct
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return  // Early return on validation error
    }

    // Step 3: Apply business logic (default rating)
    if req.Rating == 0 {
        req.Rating = 1200  // Default chess rating
    }

    // Step 4: Call repository to save data
    user, err := h.repo.Create(
        c.Request.Context(),
        req.Username,
        req.Email,
        req.Rating,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
        return  // Early return on database error
    }

    // Step 5: Convert database model to DTO
    response := ToResponse(user)

    // Step 6: Send JSON response
    c.JSON(http.StatusCreated, response)
}
```

**File**: `internal/handler/v1/user/dto.go`

```go
// Request DTO (Data Transfer Object)
type CreateRequest struct {
    Username string `json:"username" binding:"required,min=3,max=50"`
    Email    string `json:"email" binding:"required,email"`
    Rating   int32  `json:"rating,omitempty"`
}

// Response DTO
type Response struct {
    ID        uuid.UUID `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Rating    int32     `json:"rating"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// Converter function
func ToResponse(u database.User) Response {
    return Response{
        ID:        uuid.UUID(u.ID.Bytes),
        Username:  u.Username,
        Email:     u.Email,
        Rating:    u.Rating.Int32,
        CreatedAt: u.CreatedAt.Time,
        UpdatedAt: u.UpdatedAt.Time,
    }
}
```

**What happens:**
- Request JSON parsed and validated
- Business logic applied (defaults, validation)
- Repository called to persist data
- Database model converted to API response format
- JSON response sent to client

---

### 6️⃣ Repository Layer

**File**: `internal/repository/user_repository.go`

```go
type UserRepository struct {
    db *database.DB  // Database connection
}

func NewUserRepository(db *database.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(
    ctx context.Context,
    username, email string,
    rating int32,
) (database.User, error) {
    // Convert Go types to PostgreSQL types
    return r.db.CreateUser(ctx, database.CreateUserParams{
        Username: username,
        Email:    email,
        Rating:   pgtype.Int4{Int32: rating, Valid: true},
    })
}
```

**What happens:**
- Business data types converted to database types
- SQL query executed via sqlc-generated code
- Transaction management (if needed)
- Database row returned as Go struct

---

### 7️⃣ Database Layer (sqlc-generated)

**File**: `internal/database/users.sql.go` (AUTO-GENERATED)

```go
// Generated from SQL query in db/queries/users.sql
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
    row := q.db.QueryRow(ctx, createUser,
        arg.Username,
        arg.Email,
        arg.Rating,
    )

    var user User
    err := row.Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.Rating,
        &user.CreatedAt,
        &user.UpdatedAt,
    )

    return user, err
}
```

**SQL Query** (`db/queries/users.sql`):
```sql
-- name: CreateUser :one
INSERT INTO users (username, email, rating)
VALUES ($1, $2, $3)
RETURNING *;
```

**What happens:**
- SQL query executed against PostgreSQL
- Parameters safely bound (prevents SQL injection)
- Result row scanned into Go struct
- Returns database.User struct

---

### 8️⃣ Response Sent

**Response travels back up the stack:**

```
Database → Repository → Handler → Gin → HTTP Response
```

**Final JSON Response:**
```json
HTTP/1.1 201 Created
Content-Type: application/json

{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "johndoe",
  "email": "john@example.com",
  "rating": 1500,
  "created_at": "2025-11-01T10:30:00Z",
  "updated_at": "2025-11-01T10:30:00Z"
}
```

---

## Layer Responsibilities

### 🎯 Presentation Layer (Handler)

**Files**: `internal/handler/v1/user/handler.go`, `dto.go`

**Responsibilities:**
- ✅ HTTP-specific logic (status codes, headers)
- ✅ Request parsing and validation
- ✅ Response formatting (JSON serialization)
- ✅ DTO conversions (database → API format)
- ❌ NO business logic
- ❌ NO database queries

**Example:**
```go
func (h *Handler) Create(c *gin.Context) {
    // Parse request
    var req CreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Call repository
    user, err := h.repo.Create(c.Request.Context(), req.Username, req.Email, req.Rating)

    // Format response
    c.JSON(201, ToResponse(user))
}
```

---

### 🔧 Business Logic Layer (Repository)

**Files**: `internal/repository/user_repository.go`

**Responsibilities:**
- ✅ Business rules and validation
- ✅ Coordinate multiple database operations
- ✅ Transaction management
- ✅ Caching logic
- ✅ Type conversions (Go ↔ Database types)
- ❌ NO HTTP concerns
- ❌ NO direct SQL (uses sqlc)

**Example:**
```go
func (r *UserRepository) Create(ctx context.Context, username, email string, rating int32) (database.User, error) {
    // Type conversion
    return r.db.CreateUser(ctx, database.CreateUserParams{
        Username: username,
        Email:    email,
        Rating:   pgtype.Int4{Int32: rating, Valid: true},
    })
}
```

---

### 💾 Data Access Layer (Database/sqlc)

**Files**: `internal/database/*.sql.go` (GENERATED)

**Responsibilities:**
- ✅ Execute SQL queries
- ✅ Parameter binding
- ✅ Result mapping
- ✅ Connection pooling
- ❌ NO business logic
- ❌ NO validation

**Example:**
```go
// Auto-generated by sqlc
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
    row := q.db.QueryRow(ctx, createUser, arg.Username, arg.Email, arg.Rating)
    var user User
    err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Rating, &user.CreatedAt, &user.UpdatedAt)
    return user, err
}
```

---

## Complete Code Walkthrough

### Example: GET /api/v1/users/123

**1. Client sends request:**
```bash
curl http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000
```

**2. Router matches route:**
```go
// internal/router/v1_routes.go
users.GET("/:id", userHandler.Get)
```

**3. Handler extracts ID and calls repository:**
```go
// internal/handler/v1/user/handler.go
func (h *Handler) Get(c *gin.Context) {
    // Extract ID from URL path
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }

    // Call repository
    user, err := h.repo.GetByID(c.Request.Context(), id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
        return
    }

    // Convert and send response
    c.JSON(http.StatusOK, ToResponse(user))
}
```

**4. Repository converts types and queries database:**
```go
// internal/repository/user_repository.go
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (database.User, error) {
    return r.db.GetUser(ctx, pgtype.UUID{Bytes: id, Valid: true})
}
```

**5. sqlc executes SQL:**
```go
// internal/database/users.sql.go (generated)
func (q *Queries) GetUser(ctx context.Context, id pgtype.UUID) (User, error) {
    row := q.db.QueryRow(ctx, getUser, id)
    var user User
    err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Rating, &user.CreatedAt, &user.UpdatedAt)
    return user, err
}
```

**SQL executed:**
```sql
SELECT id, username, email, rating, created_at, updated_at
FROM users
WHERE id = $1;
```

**6. Response sent:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "username": "johndoe",
  "email": "john@example.com",
  "rating": 1500,
  "created_at": "2025-11-01T10:30:00Z",
  "updated_at": "2025-11-01T10:30:00Z"
}
```

---

## Error Flow

### Validation Error (400 Bad Request)

```
Client: POST /api/v1/users {"username": "ab"}  // Too short
    ↓
Handler: c.ShouldBindJSON(&req)  // Validation fails
    ↓
Response: 400 Bad Request
{
  "error": "Key: 'CreateRequest.Username' Error:Field validation for 'Username' failed on the 'min' tag"
}
```

### Database Error (500 Internal Server Error)

```
Client: POST /api/v1/users {...}
    ↓
Handler: h.repo.Create(...)
    ↓
Repository: r.db.CreateUser(...)
    ↓
Database: ERROR - duplicate email
    ↓
Repository: returns error
    ↓
Handler: c.JSON(500, gin.H{"error": "failed to create user"})
    ↓
Response: 500 Internal Server Error
{
  "error": "failed to create user"
}
```

### Not Found Error (404 Not Found)

```
Client: GET /api/v1/users/nonexistent-id
    ↓
Handler: h.repo.GetByID(id)
    ↓
Repository: r.db.GetUser(id)
    ↓
Database: no rows in result set
    ↓
Repository: returns error
    ↓
Handler: c.JSON(404, gin.H{"error": "user not found"})
    ↓
Response: 404 Not Found
{
  "error": "user not found"
}
```

---

## Visual Diagrams

### Architecture Layers

```
┌─────────────────────────────────────────────┐
│           Client (Browser/Mobile)           │
└─────────────────────────────────────────────┘
                     ▲
                     │ HTTP Request/Response
                     ▼
┌─────────────────────────────────────────────┐
│      Presentation Layer (Handlers)          │
│  - Parse HTTP requests                      │
│  - Validate input                           │
│  - Format JSON responses                    │
│  - Convert DTOs                             │
│  Files: handler/v1/user/*.go                │
└─────────────────────────────────────────────┘
                     ▲
                     │ Go Types
                     ▼
┌─────────────────────────────────────────────┐
│    Business Logic Layer (Repository)        │
│  - Business rules                           │
│  - Type conversions                         │
│  - Transaction management                   │
│  Files: repository/user_repository.go       │
└─────────────────────────────────────────────┘
                     ▲
                     │ Database Types
                     ▼
┌─────────────────────────────────────────────┐
│     Data Access Layer (sqlc/Database)       │
│  - Execute SQL queries                      │
│  - Connection pooling                       │
│  - Result mapping                           │
│  Files: database/*.sql.go (generated)       │
└─────────────────────────────────────────────┘
                     ▲
                     │ SQL
                     ▼
┌─────────────────────────────────────────────┐
│          PostgreSQL Database                │
│  - Store data                               │
│  - Enforce constraints                      │
│  - Execute queries                          │
└─────────────────────────────────────────────┘
```

### Request Flow Sequence

```
Client                Handler              Repository         Database
  │                      │                     │                 │
  │  POST /api/v1/users  │                     │                 │
  ├─────────────────────>│                     │                 │
  │                      │                     │                 │
  │                      │  Parse JSON         │                 │
  │                      │  Validate           │                 │
  │                      │                     │                 │
  │                      │  Create(...)        │                 │
  │                      ├────────────────────>│                 │
  │                      │                     │                 │
  │                      │                     │  INSERT query   │
  │                      │                     ├────────────────>│
  │                      │                     │                 │
  │                      │                     │  User row       │
  │                      │                     │<────────────────┤
  │                      │                     │                 │
  │                      │  database.User      │                 │
  │                      │<────────────────────┤                 │
  │                      │                     │                 │
  │                      │  ToResponse()       │                 │
  │                      │                     │                 │
  │  201 Created (JSON)  │                     │                 │
  │<─────────────────────┤                     │                 │
  │                      │                     │                 │
```

### Context Flow

```go
// Context carries request-scoped values through the entire flow

main.go
  └─> ctx := context.Background()
        └─> app.Run(ctx)
              └─> server.Start()
                    └─> Handler receives: c *gin.Context
                          └─> c.Request.Context() passed to Repository
                                └─> Repository passes to Database query
                                      └─> Database uses context for:
                                            - Cancellation
                                            - Timeouts
                                            - Deadlines
```

---

## Key Takeaways

### 🎯 Request Flow Summary

1. **Entry**: HTTP request → Gin router
2. **Routing**: Router matches path to handler
3. **Parsing**: Handler parses JSON into DTO
4. **Validation**: Gin validates struct tags
5. **Business Logic**: Repository applies rules
6. **Database**: sqlc executes SQL query
7. **Conversion**: Database model → DTO
8. **Response**: JSON sent to client

### 🔐 Type Safety

- **Request**: JSON → `CreateRequest` struct (validated)
- **Handler**: `CreateRequest` → Go primitives
- **Repository**: Go primitives → `pgtype` types
- **Database**: `pgtype` → SQL types
- **Response**: `database.User` → `Response` struct → JSON

### 🎨 Separation of Concerns

- **Handler**: HTTP concerns ONLY
- **Repository**: Business logic ONLY
- **Database**: SQL queries ONLY
- **No mixing of responsibilities**

### 🔄 Context Propagation

- Request context flows through all layers
- Enables cancellation and timeouts
- Always first parameter: `func Do(ctx context.Context, ...)`

---

**Last Updated**: 2025-11-01

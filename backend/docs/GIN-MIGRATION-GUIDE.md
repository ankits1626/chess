# Gin Migration Guide - From Standard Library to Gin

## What We Just Did

We migrated your simple Go server from standard library to **Gin framework** with API versioning.

---

## Before vs After Comparison

### Before (Standard Library)

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Chess Coach API - Server Running!")
    })

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprintf(w, `{"status": "healthy"}`)
    })

    port := ":8080"
    fmt.Printf("Server starting on port %s\n", port)
    log.Fatal(http.ListenAndServe(port, nil))
}
```

### After (Gin Framework)

```go
package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func main() {
    // Create Gin router
    router := gin.Default()

    // API v1 group
    v1 := router.Group("/api/v1")
    {
        v1.GET("/health", healthHandler)
    }

    // Root endpoint
    router.GET("/", rootHandler)

    // Start server
    router.Run(":8080")
}

func rootHandler(c *gin.Context) {
    c.String(http.StatusOK, "Chess Coach API - Server Running!")
}

func healthHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":  "healthy",
        "version": "1.0",
    })
}
```

---

## Key Changes Explained

### 1. **Imports**

```go
// Before
import (
    "fmt"
    "log"
    "net/http"
)

// After
import (
    "net/http"  // Still needed for status codes
    "github.com/gin-gonic/gin"  // Gin framework
)
```

**What happened:**
- Removed `fmt` and `log` - Gin handles this
- Added `github.com/gin-gonic/gin`

---

### 2. **Router Creation**

```go
// Before
// No router - used http.HandleFunc directly

// After
router := gin.Default()
```

**What `gin.Default()` gives you:**
- Logger middleware (logs all requests)
- Recovery middleware (recovers from panics)
- Automatic request/response handling

---

### 3. **API Versioning** ⭐

```go
// After
v1 := router.Group("/api/v1")
{
    v1.GET("/health", healthHandler)
}
```

**What this does:**
- Creates a route group with `/api/v1` prefix
- All routes inside `{ }` automatically get this prefix
- `v1.GET("/health", ...)` becomes `/api/v1/health`

**Benefits:**
- Easy to add v2, v3 later
- Industry standard pattern
- Clean separation

---

### 4. **Handler Functions**

```go
// Before
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"status": "healthy"}`)
})

// After
func healthHandler(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":  "healthy",
        "version": "1.0",
    })
}
```

**Key differences:**
| Standard Library | Gin |
|------------------|-----|
| `w http.ResponseWriter, r *http.Request` | `c *gin.Context` |
| `w.Header().Set(...)` | Automatic headers |
| `fmt.Fprintf(w, ...)` | `c.JSON(...)` |
| Manual JSON string | Automatic JSON encoding |

**Why Gin is better:**
- One parameter (`c`) instead of two (`w`, `r`)
- `c.JSON()` handles content-type and encoding
- `gin.H{}` is shorthand for `map[string]interface{}`

---

### 5. **Server Startup**

```go
// Before
port := ":8080"
fmt.Printf("Server starting on port %s\n", port)
log.Fatal(http.ListenAndServe(port, nil))

// After
router.Run(":8080")
```

**What `router.Run()` does:**
- Starts the server
- Logs startup automatically
- Handles errors gracefully

---

## Current API Endpoints

After migration, your API has:

| Endpoint | Method | Description | Version |
|----------|--------|-------------|---------|
| `/` | GET | Root welcome message | None |
| `/api/v1/health` | GET | Health check | v1 |

---

## How API Versioning Works

### Current Structure

```
Your API
├── / (root, not versioned)
└── /api/v1/
    └── /health
```

### When You Add More Endpoints

```go
v1 := router.Group("/api/v1")
{
    v1.GET("/health", healthHandler)
    v1.POST("/analyze", analyzeHandler)      // New!
    v1.GET("/games", listGamesHandler)       // New!
    v1.GET("/games/:id", getGameHandler)     // New!
}
```

This creates:
- `POST /api/v1/analyze`
- `GET /api/v1/games`
- `GET /api/v1/games/:id`

### Future: Adding v2

```go
v1 := router.Group("/api/v1")
{
    v1.GET("/health", healthHandlerV1)
}

v2 := router.Group("/api/v2")
{
    v2.GET("/health", healthHandlerV2)  // Different implementation
    v2.POST("/analyze", analyzeHandlerV2)
}
```

---

## Testing Your New API

### 1. Run the server

```bash
cd /Users/ankit/code/learn/chess-coach/backend
air
```

### 2. Test endpoints

```bash
# Test root (not versioned)
curl http://localhost:8080/

# Test health check (v1)
curl http://localhost:8080/api/v1/health
```

**Expected responses:**

```bash
# Root
Chess Coach API - Server Running!

# Health
{
  "status": "healthy",
  "version": "1.0"
}
```

---

## Gin Context (`c *gin.Context`) - The Power Tool

The `c *gin.Context` parameter gives you access to:

### 1. **JSON Responses**
```go
c.JSON(200, gin.H{
    "message": "success",
    "data": someData,
})
```

### 2. **URL Parameters**
```go
// Route: /games/:id
id := c.Param("id")  // Gets "123" from /games/123
```

### 3. **Query Parameters**
```go
// URL: /games?page=1&limit=10
page := c.Query("page")    // "1"
limit := c.Query("limit")  // "10"
```

### 4. **Request Body (JSON)**
```go
var request PositionRequest
if err := c.BindJSON(&request); err != nil {
    c.JSON(400, gin.H{"error": err.Error()})
    return
}
```

### 5. **Headers**
```go
authHeader := c.GetHeader("Authorization")
c.Header("X-Custom-Header", "value")
```

### 6. **Status Codes**
```go
c.JSON(200, data)  // OK
c.JSON(400, err)   // Bad Request
c.JSON(404, err)   // Not Found
c.JSON(500, err)   // Internal Server Error
```

---

## Comparison: Gin vs FastAPI

Since you know FastAPI, here's a direct comparison:

### Route Definition

**FastAPI:**
```python
@app.get("/api/v1/health")
def health():
    return {"status": "healthy"}
```

**Gin:**
```go
v1.GET("/health", func(c *gin.Context) {
    c.JSON(200, gin.H{"status": "healthy"})
})
```

### Request Body

**FastAPI:**
```python
@app.post("/api/v1/analyze")
def analyze(request: PositionRequest):
    return AnalysisResponse(...)
```

**Gin:**
```go
v1.POST("/analyze", func(c *gin.Context) {
    var req PositionRequest
    c.BindJSON(&req)
    c.JSON(200, AnalysisResponse{...})
})
```

### Path Parameters

**FastAPI:**
```python
@app.get("/games/{game_id}")
def get_game(game_id: str):
    return game
```

**Gin:**
```go
v1.GET("/games/:id", func(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, game)
})
```

---

## What You Gain with Gin

✅ **Cleaner code** - Less boilerplate
✅ **Automatic JSON** - No manual encoding
✅ **Built-in middleware** - Logging, recovery, CORS
✅ **Route grouping** - API versioning made easy
✅ **Better error handling** - Panic recovery
✅ **Parameter binding** - Automatic validation
✅ **Better performance** - Optimized routing

---

## Next Steps

Now that you have Gin set up:

1. ✅ **Done:** Migrated to Gin
2. ✅ **Done:** Added API versioning (`/api/v1`)
3. **Next:** Add Swagger documentation
4. **Next:** Add your first chess endpoint
5. **Later:** Add database, authentication, etc.

---

## Common Gin Patterns

### Multiple Route Groups

```go
func main() {
    router := gin.Default()

    // Public API
    v1 := router.Group("/api/v1")
    {
        v1.GET("/health", healthHandler)
        v1.POST("/analyze", analyzeHandler)
    }

    // Admin API (could add auth middleware)
    admin := router.Group("/admin")
    {
        admin.GET("/stats", adminStatsHandler)
    }

    router.Run(":8080")
}
```

### Middleware Example

```go
// Custom middleware
func authMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "unauthorized"})
            c.Abort()
            return
        }
        c.Next()
    }
}

// Use it
authorized := router.Group("/api/v1")
authorized.Use(authMiddleware())
{
    authorized.POST("/analyze", analyzeHandler)
}
```

---

## Quick Reference

### Common Gin Methods

```go
// Response methods
c.JSON(code, obj)           // JSON response
c.String(code, str)         // Plain text
c.HTML(code, name, data)    // HTML template
c.File(filepath)            // Serve file
c.Redirect(code, url)       // Redirect

// Request methods
c.Param("id")               // URL parameter
c.Query("key")              // Query parameter
c.BindJSON(&obj)            // Parse JSON body
c.GetHeader("key")          // Get header

// Context methods
c.Next()                    // Call next middleware
c.Abort()                   // Stop execution
c.Set("key", value)         // Store data
c.Get("key")                // Retrieve data
```

---

**You're now using Gin!** 🎉

Ready to add Swagger next?

# Swagger Setup - Incremental Learning Guide

## Understanding the Pieces

### 1. What is Swagger?
- **Purpose:** Creates interactive API documentation
- **Works with:** Any Go web server (standard library, Gin, Echo, Fiber, etc.)
- **You get:** A web page where you can see and test all your APIs

### 2. What is Gin?
- **Purpose:** A web framework that makes building REST APIs easier
- **NOT required for Swagger**
- **Think of it as:** Express.js for Node.js, but for Go
- **Standard library vs Gin:**
  ```go
  // Standard library (what you have now)
  http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
      w.Header().Set("Content-Type", "application/json")
      fmt.Fprintf(w, `{"status": "healthy"}`)
  })

  // With Gin (simpler)
  router.GET("/health", func(c *gin.Context) {
      c.JSON(200, gin.H{"status": "healthy"})
  })
  ```

---

## Decision: Two Paths Forward

### Path A: Keep Standard Library + Add Swagger (Simpler)
- ✅ Less new concepts to learn
- ✅ Uses what you already know
- ✅ Still get Swagger documentation
- ❌ More boilerplate code as app grows

### Path B: Switch to Gin + Add Swagger (Industry Standard)
- ✅ Cleaner code for complex APIs
- ✅ Better middleware support
- ✅ What most Go projects use
- ❌ One more thing to learn

**Recommendation:** Start with **Path A** (standard library), understand Swagger, then optionally migrate to Gin later.

---

## Path A: Swagger with Standard Library

### Step 1: Current State (You are here)

Your [cmd/server/main.go](../cmd/server/main.go) is simple:
```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    http.HandleFunc("/", rootHandler)
    http.HandleFunc("/health", healthHandler)

    port := ":8080"
    fmt.Printf("Server starting on port %s\n", port)
    log.Fatal(http.ListenAndServe(port, nil))
}
```

**Endpoints:**
- `GET /` - Root message
- `GET /health` - Health check

---

### Step 2: Add Swagger Annotations (Just Comments!)

Swagger works by reading special comments in your code. Let's add them:

**File: cmd/server/main.go** (enhanced version)

```go
package main

import (
    "fmt"
    "log"
    "net/http"

    httpSwagger "github.com/swaggo/http-swagger"
    _ "github.com/ankits1626/chess-coach-backend/docs"
)

// @title Chess Coach API
// @version 1.0
// @description Simple chess coaching API
// @host localhost:8080
// @BasePath /
func main() {
    // Root endpoint
    http.HandleFunc("/", rootHandler)

    // Health check endpoint
    http.HandleFunc("/health", healthHandler)

    // Swagger UI endpoint (this is new!)
    http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

    port := ":8080"
    fmt.Printf("Server starting on port %s\n", port)
    log.Fatal(http.ListenAndServe(port, nil))
}

// rootHandler handles the root endpoint
// @Summary Root endpoint
// @Description Returns a welcome message
// @Produce plain
// @Success 200 {string} string "Welcome message"
// @Router / [get]
func rootHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Chess Coach API - Server Running!")
}

// healthHandler handles the health check
// @Summary Health check
// @Description Returns API health status
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func healthHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"status": "healthy"}`)
}
```

**What changed?**
1. Added `httpSwagger` import (for standard library)
2. Added `@` comment annotations above functions
3. Added `/swagger/` endpoint to view docs
4. Extracted handlers to named functions (cleaner for annotations)

---

### Step 3: Install Dependencies

```bash
# Install swag CLI tool (generates documentation)
go install github.com/swaggo/swag/cmd/swag@latest

# Install http-swagger (for standard library, NOT gin-swagger)
go get -u github.com/swaggo/http-swagger
```

**Note the difference:**
- `gin-swagger` = For Gin framework ❌ (not what you need)
- `http-swagger` = For standard library ✅ (what you need)

---

### Step 4: Generate Swagger Documentation

```bash
cd /Users/ankit/code/learn/chess-coach/backend
swag init -g cmd/server/main.go
```

This creates:
```
backend/
├── docs/
│   ├── docs.go       # Generated Go code
│   ├── swagger.json  # API specification
│   └── swagger.yaml  # API specification (YAML)
```

---

### Step 5: Run and Test

```bash
# Run your server
go run cmd/server/main.go

# Visit in browser:
http://localhost:8080/swagger/index.html
```

You'll see an interactive API documentation page!

---

## Path B: Migrate to Gin (Optional, Later)

If you want to learn Gin later, I can help you migrate. But let's get Swagger working with standard library first.

---

## Summary of Concepts

| Concept | What It Is | Do You Need It? |
|---------|-----------|-----------------|
| **Standard Library** | Go's built-in HTTP server | ✅ Yes (you have this) |
| **Gin** | Web framework (optional) | ❌ No (nice to have) |
| **Swagger/swag** | Documentation generator | ✅ Yes (for docs) |
| **http-swagger** | Swagger UI for standard lib | ✅ Yes (to view docs) |
| **gin-swagger** | Swagger UI for Gin | ❌ No (only if using Gin) |

---

## Next Steps - Your Choice

**Option 1: Go with Standard Library** (Recommended for learning)
1. I'll update main.go with annotations (standard library)
2. Install `http-swagger` (not gin-swagger)
3. Generate docs
4. Test Swagger UI

**Option 2: Learn Gin First, Then Swagger**
1. Understand Gin basics
2. Convert your server to Gin
3. Add Swagger

**Which path do you want to take?**

I recommend **Option 1** - keep it simple, learn Swagger, then optionally learn Gin later.

---

**Key Takeaway:**
- Swagger = Documentation tool (works with anything)
- Gin = Web framework (optional, makes code cleaner)
- They are independent - you can use Swagger WITHOUT Gin!

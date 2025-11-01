# Swagger Integration Guide

## Goal
Add Swagger/OpenAPI documentation to our SOLID-compliant Go backend with crisp annotations.

---

## Why Swagger?

- **Interactive API docs** - Test endpoints directly in browser
- **Auto-generated** - Docs stay in sync with code
- **Industry standard** - OpenAPI 3.0 specification
- **Client generation** - Can generate API clients for frontend

---

## Prerequisites

✅ SOLID-compliant backend structure (completed)
✅ Gin framework installed
✅ All tests passing

---

## Step-by-Step Integration

### Step 1: Install Swag CLI
**What:** Install the swag command-line tool
**Time:** 1 minute

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

**Verify:**
```bash
swag --version
```

---

### Step 2: Install Swagger Dependencies
**What:** Add gin-swagger and swag files packages
**Time:** 1 minute

```bash
cd /Users/ankit/code/learn/chess-coach/backend
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```

---

### Step 3: Add Main Swagger Annotations
**What:** Annotate main.go with API metadata
**File:** `cmd/server/main.go`
**Time:** 2 minutes

Add at the top of main.go (before package declaration):

```go
// @title Chess Coach API
// @version 1.0
// @description API for chess game analysis and coaching
// @host localhost:8080
// @BasePath /api/v1
```

**Full example:**
```go
// @title Chess Coach API
// @version 1.0
// @description API for chess game analysis and coaching
// @host localhost:8080
// @BasePath /api/v1

// Package main bootstraps the Chess Coach API.
package main
```

---

### Step 4: Annotate Health Endpoint
**What:** Add Swagger annotations to health handler
**File:** `internal/router/v1/health.go`
**Time:** 2 minutes

Update HealthHandler function:

```go
// HealthHandler handles health check endpoint.
// @Summary Health check
// @Description Returns API health status
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:  "healthy",
		Version: "1.0",
	})
}
```

**Annotation explanation:**
- `@Summary` - Brief description (shows in endpoint list)
- `@Description` - Detailed description
- `@Tags` - Groups endpoints together
- `@Produce` - Response content type
- `@Success` - Success response with schema
- `@Router` - Route path (relative to BasePath)

---

### Step 5: Generate Swagger Docs
**What:** Run swag to generate documentation
**Time:** 1 minute

```bash
cd /Users/ankit/code/learn/chess-coach/backend
swag init -g cmd/server/main.go
```

**What it creates:**
```
backend/
├── docs/
│   ├── docs.go          # Generated Go code
│   ├── swagger.json     # OpenAPI JSON spec
│   └── swagger.yaml     # OpenAPI YAML spec
```

---

### Step 6: Add Swagger Route
**What:** Register Swagger UI endpoint
**File:** `internal/router/router.go`
**Time:** 2 minutes

**Update imports:**
```go
import (
	"net/http"

	"github.com/ankits1626/chess-coach-backend/internal/router/v1"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/ankits1626/chess-coach-backend/docs" // Import generated docs
)
```

**Add Swagger route in Setup():**
```go
func Setup() *gin.Engine {
	r := gin.Default()

	// Register v1 API routes
	v1.RegisterRoutes(r)

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Root endpoint
	r.GET("/", rootHandler)

	return r
}
```

---

### Step 7: Update .gitignore
**What:** Ignore generated docs (they're auto-generated)
**File:** `.gitignore`
**Time:** 1 minute

Add to the Go section:
```
# Swagger generated docs
backend/docs/
```

**Note:** Some teams commit docs/, others regenerate. We'll ignore it.

---

### Step 8: Test Swagger UI
**What:** Start server and access Swagger
**Time:** 2 minutes

**Start server:**
```bash
air
```

**Access Swagger UI:**
```
http://localhost:8080/swagger/index.html
```

**Test the endpoint:**
1. Click on "health" endpoint
2. Click "Try it out"
3. Click "Execute"
4. See the response!

---

### Step 9: Update Air Config (Optional)
**What:** Auto-regenerate docs on file changes
**File:** `.air.toml`
**Time:** 1 minute

Update the build command:
```toml
[build]
  cmd = "swag init -g cmd/server/main.go && go build -o ./tmp/main ./cmd/server"
  bin = "tmp/main"
```

This regenerates Swagger docs before each build.

---

## Complete File Examples

### cmd/server/main.go (with Swagger)

```go
// @title Chess Coach API
// @version 1.0
// @description API for chess game analysis and coaching
// @host localhost:8080
// @BasePath /api/v1

// Package main bootstraps the Chess Coach API.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ankits1626/chess-coach-backend/internal/config"
	"github.com/ankits1626/chess-coach-backend/internal/server"
	_ "github.com/ankits1626/chess-coach-backend/docs" // Swagger docs
)

// main starts server with graceful shutdown.
func main() {
	// Load configuration
	cfg := config.Load()

	// Create server
	srv := server.New(cfg)

	// Start server in goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.Port)
	log.Printf("Swagger UI: http://localhost:%s/swagger/index.html", cfg.Port)

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
```

### internal/router/router.go (with Swagger)

```go
// Package router configures HTTP routes.
package router

import (
	"net/http"

	"github.com/ankits1626/chess-coach-backend/internal/router/v1"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/ankits1626/chess-coach-backend/docs"
)

// Setup creates configured Gin router.
func Setup() *gin.Engine {
	r := gin.Default()

	// Register v1 API routes
	v1.RegisterRoutes(r)

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Root endpoint
	r.GET("/", rootHandler)

	return r
}

// rootHandler handles the root endpoint.
func rootHandler(c *gin.Context) {
	c.String(http.StatusOK, "Chess Coach API - Server Running!")
}
```

### internal/router/v1/health.go (with Swagger)

```go
// Package v1 provides API v1 endpoints.
package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthResponse represents health check response.
type HealthResponse struct {
	// Status indicates if service is healthy
	Status string `json:"status" example:"healthy"`
	// Version is the API version
	Version string `json:"version" example:"1.0"`
}

// HealthHandler handles health check endpoint.
// @Summary Health check
// @Description Returns API health status
// @Tags health
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /health [get]
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:  "healthy",
		Version: "1.0",
	})
}
```

---

## Swagger Annotation Cheat Sheet

### Common Annotations

```go
// @Summary Brief description
// @Description Detailed description
// @Tags group-name
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param body body RequestType true "Request body"
// @Success 200 {object} ResponseType
// @Failure 400 {object} ErrorResponse
// @Router /endpoint [get]
```

### Parameter Types
- `path` - URL path parameter (e.g., `/users/:id`)
- `query` - Query string (e.g., `?name=value`)
- `body` - Request body
- `header` - HTTP header

### HTTP Methods
```go
// @Router /endpoint [get]
// @Router /endpoint [post]
// @Router /endpoint [put]
// @Router /endpoint [delete]
// @Router /endpoint [patch]
```

---

## Common Swagger Patterns

### POST Endpoint with Body

```go
// CreateGame creates a new chess game.
// @Summary Create game
// @Description Creates new game from PGN
// @Tags games
// @Accept json
// @Produce json
// @Param game body CreateGameRequest true "Game data"
// @Success 201 {object} GameResponse
// @Failure 400 {object} ErrorResponse
// @Router /games [post]
func CreateGame(c *gin.Context) {
	// Handler code
}
```

### GET with Path Parameter

```go
// GetGame retrieves a game by ID.
// @Summary Get game
// @Description Get game details by ID
// @Tags games
// @Produce json
// @Param id path string true "Game ID"
// @Success 200 {object} GameResponse
// @Failure 404 {object} ErrorResponse
// @Router /games/{id} [get]
func GetGame(c *gin.Context) {
	// Handler code
}
```

### GET with Query Parameters

```go
// ListGames lists all games.
// @Summary List games
// @Description List games with pagination
// @Tags games
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {array} GameResponse
// @Router /games [get]
func ListGames(c *gin.Context) {
	// Handler code
}
```

---

## Workflow After Setup

### Every time you add/modify an endpoint:

1. **Add Swagger annotations** to the handler
2. **Regenerate docs:** `swag init -g cmd/server/main.go`
3. **Restart server:** Air does this automatically if configured
4. **Test in Swagger UI:** http://localhost:8080/swagger/index.html

**Or if Air is configured (Step 9):**
Just save the file - Air regenerates docs and restarts!

---

## Troubleshooting

### Issue: "docs not found"
**Solution:** Run `swag init -g cmd/server/main.go`

### Issue: Changes not reflected
**Solution:**
1. Stop server
2. Run `swag init -g cmd/server/main.go`
3. Restart server

### Issue: "404 Not Found" on /swagger
**Solution:** Check router.go has the swagger route registered

### Issue: Endpoint not showing
**Solution:**
1. Check annotations syntax
2. Ensure handler is exported (starts with capital letter)
3. Regenerate docs

---

## Next Steps After Swagger

1. Add more endpoints (analyze, games)
2. Add request/response validation
3. Add authentication endpoints
4. Add error response schemas
5. Group endpoints with tags

---

## Summary

**After integration you'll have:**
- ✅ Interactive API documentation
- ✅ Auto-generated OpenAPI spec
- ✅ Easy testing in browser
- ✅ Client code generation capability
- ✅ Professional API docs

**Access:**
- Swagger UI: http://localhost:8080/swagger/index.html
- JSON spec: http://localhost:8080/swagger/doc.json

---

## Time Estimate

| Step | Time | Cumulative |
|------|------|------------|
| 1. Install swag | 1 min | 1 min |
| 2. Install deps | 1 min | 2 min |
| 3. Main annotations | 2 min | 4 min |
| 4. Health annotations | 2 min | 6 min |
| 5. Generate docs | 1 min | 7 min |
| 6. Add route | 2 min | 9 min |
| 7. Update gitignore | 1 min | 10 min |
| 8. Test | 2 min | 12 min |
| 9. Air config | 1 min | 13 min |

**Total: ~13 minutes**

---

**Ready to start? Say "start" and I'll begin with Step 1!**

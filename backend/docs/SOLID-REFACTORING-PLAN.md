# SOLID Principles Refactoring Plan

## Current State Analysis

### Problems with Current `main.go`

#### 1. Single Responsibility Principle (SRP) Violations
```go
func main() {
    router := gin.Default()           // Server creation
    v1 := router.Group("/api/v1")    // Route configuration
    v1.GET("/health", healthHandler) // Route registration
    router.Run(":8080")              // Server startup
}
```

**Issues:**
- `main()` does 4 different things
- Route handlers mixed with setup
- Configuration hardcoded

#### 2. Missing Documentation
- No package-level documentation
- No function documentation (godoc format)
- No examples

#### 3. Poor Testability
- Can't test route setup separately
- Can't mock dependencies
- No interfaces

---

## Proposed Structure (SOLID Compliant)

### Directory Structure
```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Entry point only
├── internal/
│   ├── server/
│   │   ├── server.go            # Server interface & implementation
│   │   └── server_test.go
│   ├── router/
│   │   ├── router.go            # Route setup
│   │   ├── router_test.go
│   │   └── v1/
│   │       ├── health.go        # Health handler
│   │       ├── health_test.go
│   │       └── routes.go        # v1 routes
│   ├── config/
│   │   └── config.go            # Configuration
│   └── handlers/
│       └── response.go          # Standard responses
└── pkg/
    └── logger/
        └── logger.go            # Logging interface
```

### Applying SOLID Principles

#### 1. Single Responsibility Principle (SRP)
Each file has ONE job:
- `main.go` - Only bootstraps the application
- `server.go` - Only manages HTTP server lifecycle
- `router.go` - Only sets up routes
- `health.go` - Only handles health endpoint

#### 2. Open/Closed Principle (OCP)
- Routes can be added without modifying existing code
- Handlers follow a common interface
- Middleware can be added/removed easily

#### 3. Liskov Substitution Principle (LSP)
- All handlers implement same interface
- Can swap implementations without breaking

#### 4. Interface Segregation Principle (ISP)
- Small, focused interfaces
- Handlers don't depend on things they don't use

#### 5. Dependency Inversion Principle (DIP)
- Depend on interfaces, not concrete types
- Easy to mock for testing

---

## Refactored Code Structure

### 1. Entry Point (`cmd/server/main.go`)

```go
// Package main is the entry point for the Chess Coach API server.
//
// This package bootstraps the application by:
//   - Loading configuration
//   - Setting up the HTTP server
//   - Starting the server
//
// Usage:
//   go run cmd/server/main.go
//
// Environment variables:
//   PORT - Server port (default: 8080)
//   GIN_MODE - Gin mode: debug, release, test (default: debug)
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
)

// main is the application entry point.
// It handles configuration loading, server initialization, and graceful shutdown.
func main() {
	// Load configuration
	cfg := config.Load()

	// Create server
	srv := server.New(cfg)

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", cfg.Port)

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

**SOLID compliance:**
- ✅ SRP: Only bootstraps application
- ✅ DIP: Depends on `server.Server` interface
- ✅ Documentation: Package and function docs
- ✅ Error handling: Proper error management
- ✅ Graceful shutdown: Production-ready

---

### 2. Configuration (`internal/config/config.go`)

```go
// Package config provides application configuration management.
//
// Configuration is loaded from environment variables with sensible defaults.
package config

import (
	"os"
)

// Config holds all application configuration.
// It follows the 12-factor app methodology by using environment variables.
type Config struct {
	// Port is the HTTP server port (default: "8080")
	Port string

	// Environment is the application environment: development, staging, production
	Environment string

	// GinMode is the Gin framework mode: debug, release, test
	GinMode string
}

// Load creates a new Config instance with values from environment variables.
// If an environment variable is not set, a default value is used.
//
// Environment variables:
//   PORT - HTTP server port (default: "8080")
//   ENVIRONMENT - Application environment (default: "development")
//   GIN_MODE - Gin mode (default: "debug")
//
// Example:
//   cfg := config.Load()
//   fmt.Println(cfg.Port) // "8080" or value from PORT env var
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		GinMode:     getEnv("GIN_MODE", "debug"),
	}
}

// getEnv retrieves an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
```

**SOLID compliance:**
- ✅ SRP: Only handles configuration
- ✅ Documentation: Full godoc comments
- ✅ Testability: Pure functions

---

### 3. Server (`internal/server/server.go`)

```go
// Package server provides HTTP server management with graceful shutdown support.
package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ankits1626/chess-coach-backend/internal/config"
	"github.com/ankits1626/chess-coach-backend/internal/router"
)

// Server represents an HTTP server with lifecycle management.
// It abstracts the underlying HTTP server implementation.
type Server interface {
	// Start begins listening for HTTP requests.
	// It blocks until the server is shut down or an error occurs.
	Start() error

	// Shutdown gracefully shuts down the server.
	// It waits for active connections to close within the context deadline.
	Shutdown(ctx context.Context) error
}

// server is the concrete implementation of the Server interface.
type server struct {
	httpServer *http.Server
	config     *config.Config
}

// New creates a new Server instance with the given configuration.
//
// The server is configured with:
//   - Gin router with all application routes
//   - Default middleware (logging, recovery)
//   - Configured port from config
//
// Example:
//   cfg := config.Load()
//   srv := server.New(cfg)
//   srv.Start()
func New(cfg *config.Config) Server {
	// Set Gin mode
	gin.SetMode(cfg.GinMode)

	// Create router with all routes
	r := router.Setup()

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	return &server{
		httpServer: httpServer,
		config:     cfg,
	}
}

// Start implements Server.Start.
func (s *server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown implements Server.Shutdown.
func (s *server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
```

**SOLID compliance:**
- ✅ SRP: Only manages server lifecycle
- ✅ OCP: Can extend without modifying
- ✅ LSP: Implements Server interface
- ✅ ISP: Server interface is minimal
- ✅ DIP: Depends on config interface
- ✅ Documentation: Complete godoc

---

### 4. Router (`internal/router/router.go`)

```go
// Package router provides HTTP route configuration for the application.
//
// Routes are organized by API version (v1, v2, etc.).
package router

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/ankits1626/chess-coach-backend/internal/router/v1"
)

// Setup creates and configures the main application router.
//
// The router includes:
//   - API v1 routes under /api/v1
//   - Root endpoint at /
//   - Default middleware (logger, recovery)
//
// Route structure:
//   / - Root welcome endpoint
//   /api/v1/health - Health check
//   /api/v1/analyze - Chess position analysis (future)
//
// Example:
//   r := router.Setup()
//   r.Run(":8080")
func Setup() *gin.Engine {
	r := gin.Default()

	// Register v1 API routes
	v1.RegisterRoutes(r)

	// Root endpoint (not versioned)
	r.GET("/", rootHandler)

	return r
}

// rootHandler handles the root endpoint.
// It returns a welcome message to verify the API is running.
//
// Route: GET /
// Response: 200 OK with plain text message
func rootHandler(c *gin.Context) {
	c.String(200, "Chess Coach API - Server Running!")
}
```

**SOLID compliance:**
- ✅ SRP: Only sets up routes
- ✅ OCP: Can add routes without modifying
- ✅ Documentation: Full godoc

---

### 5. V1 Health Handler (`internal/router/v1/health.go`)

```go
// Package v1 provides API version 1 handlers and routes.
package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthResponse represents the response for the health check endpoint.
// It follows the standard health check format used in microservices.
type HealthResponse struct {
	// Status indicates if the service is healthy
	Status string `json:"status" example:"healthy"`

	// Version is the API version
	Version string `json:"version" example:"1.0"`

	// Timestamp is when the health check was performed
	Timestamp string `json:"timestamp,omitempty" example:"2025-11-01T10:00:00Z"`
}

// HealthHandler handles the health check endpoint.
//
// It returns the current health status of the API.
// This endpoint is typically used by load balancers and monitoring systems.
//
// Route: GET /api/v1/health
// Response: 200 OK with HealthResponse JSON
//
// Example response:
//   {
//     "status": "healthy",
//     "version": "1.0"
//   }
//
// @Summary Health check
// @Description Returns the health status of the API
// @Tags general
// @Produce json
// @Success 200 {object} HealthResponse
// @Router /api/v1/health [get]
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:  "healthy",
		Version: "1.0",
	})
}
```

**SOLID compliance:**
- ✅ SRP: Only handles health check
- ✅ Documentation: Complete with examples
- ✅ Type safety: Proper response struct
- ✅ Swagger annotations: Ready for docs

---

### 6. V1 Routes (`internal/router/v1/routes.go`)

```go
// Package v1 provides API version 1 routes registration.
package v1

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all v1 API routes.
//
// All routes are prefixed with /api/v1.
// This function is called by the main router setup.
//
// Registered routes:
//   GET /api/v1/health - Health check endpoint
//
// Future routes:
//   POST /api/v1/analyze - Chess position analysis
//   GET /api/v1/games - List games
//   GET /api/v1/games/:id - Get game details
//
// Example:
//   r := gin.Default()
//   v1.RegisterRoutes(r)
func RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", HealthHandler)

		// Future endpoints
		// v1.POST("/analyze", AnalyzeHandler)
		// v1.GET("/games", ListGamesHandler)
		// v1.GET("/games/:id", GetGameHandler)
	}
}
```

**SOLID compliance:**
- ✅ SRP: Only registers routes
- ✅ OCP: Easy to add new routes
- ✅ Documentation: Clear structure

---

## Go Documentation Standards (godoc)

### Package Documentation
```go
// Package name provides a description of what the package does.
//
// More detailed explanation if needed.
// Can span multiple lines.
//
// Example:
//   // Usage example
package name
```

### Function Documentation
```go
// FunctionName does something specific.
//
// It explains what the function does in detail.
// Parameters and return values are described.
//
// Example:
//   result := FunctionName(param)
func FunctionName(param string) (string, error) {
```

### Type Documentation
```go
// TypeName represents something specific.
// Each field is documented below.
type TypeName struct {
	// Field1 is the first field
	Field1 string `json:"field1"`

	// Field2 is the second field
	Field2 int `json:"field2"`
}
```

---

## Benefits of This Structure

### 1. Testability
Each component can be tested in isolation:
```go
// server_test.go
func TestServer_Start(t *testing.T) {
	cfg := &config.Config{Port: "8080"}
	srv := server.New(cfg)
	// Test server
}

// health_test.go
func TestHealthHandler(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	HealthHandler(c)
	// Assert response
}
```

### 2. Maintainability
- Clear separation of concerns
- Easy to find code
- Self-documenting structure

### 3. Scalability
- Easy to add new versions (v2, v3)
- Can add middleware easily
- Can swap implementations

### 4. Production-Ready
- Graceful shutdown
- Proper error handling
- Configuration management
- Logging support

---

## Next Steps

1. ✅ Review this refactoring plan
2. Implement the structure incrementally
3. Add tests for each component
4. Add Swagger documentation
5. Add logging middleware
6. Add error handling middleware

Would you like me to implement this structure?

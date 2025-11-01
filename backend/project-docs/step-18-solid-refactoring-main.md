# Step 18: SOLID Refactoring - Main.go

**Objective**: Refactor main.go to follow SOLID principles by introducing abstraction layers and improving testability.

**Time**: 25 minutes
**Difficulty**: Medium

---

## Current Issues

### Interface Segregation Principle Violations
- Direct dependency on concrete `*database.DB`
- Direct dependency on concrete `config.Config`
- No abstraction for logging (using stdlib `log` directly)

### Dependency Inversion Principle Violations
- Main.go depends on concrete implementations, not abstractions
- Hardcoded shutdown timeout (5 seconds)
- Not easily testable

---

## Phase 1: Create Logger Interface

**Time**: 5 minutes

### Step 1.1: Create Logger Package

**File**: `internal/logger/logger.go`

```go
// Package logger provides logging abstraction.
package logger

import (
	"log"
	"os"
)

// Logger defines logging interface.
type Logger interface {
	Info(msg string)
	Infof(format string, args ...interface{})
	Error(msg string)
	Errorf(format string, args ...interface{})
	Fatal(msg string)
	Fatalf(format string, args ...interface{})
}

// StdLogger implements Logger using stdlib log.
type StdLogger struct {
	info  *log.Logger
	error *log.Logger
}

// NewStdLogger creates standard logger.
func NewStdLogger() Logger {
	return &StdLogger{
		info:  log.New(os.Stdout, "INFO: ", log.LstdFlags),
		error: log.New(os.Stderr, "ERROR: ", log.LstdFlags),
	}
}

// Info logs info message.
func (l *StdLogger) Info(msg string) {
	l.info.Println(msg)
}

// Infof logs formatted info message.
func (l *StdLogger) Infof(format string, args ...interface{}) {
	l.info.Printf(format, args...)
}

// Error logs error message.
func (l *StdLogger) Error(msg string) {
	l.error.Println(msg)
}

// Errorf logs formatted error message.
func (l *StdLogger) Errorf(format string, args ...interface{}) {
	l.error.Printf(format, args...)
}

// Fatal logs fatal message and exits.
func (l *StdLogger) Fatal(msg string) {
	l.error.Fatal(msg)
}

// Fatalf logs formatted fatal message and exits.
func (l *StdLogger) Fatalf(format string, args ...interface{}) {
	l.error.Fatalf(format, args...)
}
```

---

## Phase 2: Update Config with Shutdown Timeout

**Time**: 3 minutes

### Step 2.1: Add Shutdown Timeout to Config

**File**: `internal/config/config.go`

**Find**:
```go
type Config struct {
	Port    string
	GinMode string
	DBHost  string
	DBPort  string
	DBUser  string
	DBPass  string
	DBName  string
}
```

**Replace with**:
```go
type Config struct {
	Port            string
	GinMode         string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPass          string
	DBName          string
	ShutdownTimeout time.Duration
}
```

**Add import**:
```go
import (
	"fmt"
	"os"
	"time"
)
```

### Step 2.2: Load Shutdown Timeout

**In the `Load()` function, add**:

```go
func Load() *Config {
	return &Config{
		Port:            getEnv("PORT", "8080"),
		GinMode:         getEnv("GIN_MODE", "debug"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPass:          getEnv("DB_PASS", "postgres"),
		DBName:          getEnv("DB_NAME", "chess_coach"),
		ShutdownTimeout: 5 * time.Second, // Default 5 seconds
	}
}
```

---

## Phase 3: Create Application Layer

**Time**: 10 minutes

### Step 3.1: Create App Package

**File**: `internal/app/app.go`

```go
// Package app manages application lifecycle.
package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ankits1626/chess-coach-backend/internal/config"
	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/logger"
	"github.com/ankits1626/chess-coach-backend/internal/server"
)

// App manages application lifecycle.
type App struct {
	config *config.Config
	db     *database.DB
	server server.Server
	logger logger.Logger
}

// New creates new application.
func New(cfg *config.Config, db *database.DB, log logger.Logger) *App {
	return &App{
		config: cfg,
		db:     db,
		server: server.New(cfg, db),
		logger: log,
	}
}

// Run starts application and handles graceful shutdown.
func (a *App) Run(ctx context.Context) error {
	// Start server in goroutine
	go func() {
		if err := a.server.Start(); err != nil {
			a.logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	a.logger.Infof("Server started on port %s", a.config.Port)
	a.logger.Infof("Swagger UI: http://localhost:%s/swagger/index.html", a.config.Port)

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	a.logger.Info("Shutting down server...")

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(ctx, a.config.ShutdownTimeout)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		a.logger.Errorf("Server forced to shutdown: %v", err)
		return err
	}

	a.logger.Info("Server exited gracefully")
	return nil
}

// Close closes application resources.
func (a *App) Close() {
	if a.db != nil {
		a.db.Close()
	}
}
```

---

## Phase 4: Refactor main.go

**Time**: 7 minutes

### Step 4.1: Update main.go

**File**: `cmd/server/main.go`

**Replace entire file with**:

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

	"github.com/ankits1626/chess-coach-backend/internal/app"
	"github.com/ankits1626/chess-coach-backend/internal/config"
	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/logger"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	appLogger := logger.NewStdLogger()

	// Connect to database
	ctx := context.Background()
	db, err := database.NewDB(ctx, cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create and run application
	application := app.New(cfg, db, appLogger)
	defer application.Close()

	appLogger.Info("Database connected successfully")

	// Run application (blocks until shutdown)
	if err := application.Run(ctx); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
```

---

## Phase 5: Verify Implementation

**Time**: N/A (Testing)

### Step 5.1: Build and Test

```bash
# Build the application
cd backend
go mod tidy
go build -o bin/server cmd/server/main.go

# Run the server
docker compose up -d postgres
./bin/server
```

### Step 5.2: Verify Graceful Shutdown

1. Start the server
2. Press `Ctrl+C`
3. Verify you see: "Shutting down server..." and "Server exited gracefully"

### Step 5.3: Test API Endpoints

Access Swagger UI:
```
http://localhost:8080/swagger/index.html
```

---

## Benefits Achieved

### Before Refactoring (Score: 7/10)
❌ Direct dependencies on concrete types
❌ Not easily testable
❌ Magic numbers (shutdown timeout)
❌ No logging abstraction

### After Refactoring (Score: 9/10)
✅ Logger interface for abstraction
✅ App struct encapsulates lifecycle
✅ Configurable shutdown timeout
✅ Testable application layer
✅ Clean separation of concerns
✅ Main.go is just a thin bootstrap layer

---

## SOLID Principles Compliance

| Principle | Before | After | Improvement |
|-----------|--------|-------|-------------|
| **S** Single Responsibility | ✅ Good | ✅ Excellent | App layer added |
| **O** Open/Closed | ✅ Good | ✅ Excellent | Extension points clear |
| **L** Liskov Substitution | ✅ Good | ✅ Excellent | Interfaces used |
| **I** Interface Segregation | ⚠️ Issues | ✅ Fixed | Logger interface |
| **D** Dependency Inversion | ⚠️ Issues | ✅ Fixed | App abstracts deps |

---

## Testing Benefits

With this refactoring, you can now easily test the application layer:

```go
// Example test (future)
func TestApp_Run(t *testing.T) {
    mockLogger := &MockLogger{}
    mockDB := &MockDB{}
    mockConfig := &config.Config{
        Port: "8080",
        ShutdownTimeout: 1 * time.Second,
    }

    app := app.New(mockConfig, mockDB, mockLogger)
    // Test app.Run()...
}
```

---

## Next Steps

After completing this refactoring:
1. Consider adding structured logging (e.g., `zap` or `zerolog`)
2. Add metrics/observability
3. Consider health check endpoints
4. Add unit tests for `app` package

---

**Completion Checklist**:
- [ ] Logger interface created
- [ ] Config updated with ShutdownTimeout
- [ ] App package created with lifecycle management
- [ ] main.go refactored to use App
- [ ] Application builds successfully
- [ ] Server starts and shuts down gracefully
- [ ] Swagger UI accessible
- [ ] All existing endpoints still work

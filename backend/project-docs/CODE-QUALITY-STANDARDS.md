# Code Quality Standards for Chess Coach Backend

## Overview

This document defines the code quality standards for the Chess Coach backend project. All code must follow these standards.

---

## SOLID Principles

### ✅ We Follow SOLID Throughout

1. **Single Responsibility Principle (SRP)**
   - Each file has ONE purpose
   - Each function does ONE thing
   - Each package has ONE responsibility

2. **Open/Closed Principle (OCP)**
   - Open for extension
   - Closed for modification
   - Use interfaces for extensibility

3. **Liskov Substitution Principle (LSP)**
   - Implementations are interchangeable
   - Interfaces define contracts

4. **Interface Segregation Principle (ISP)**
   - Small, focused interfaces
   - No fat interfaces

5. **Dependency Inversion Principle (DIP)**
   - Depend on abstractions (interfaces)
   - Not on concrete implementations

---

## Documentation Standards (godoc)

### Every Package Must Have Documentation

```go
// Package config provides application configuration management.
//
// Configuration is loaded from environment variables with sensible defaults.
// This follows the 12-factor app methodology.
//
// Example:
//   cfg := config.Load()
//   fmt.Println(cfg.Port)
package config
```

### Every Exported Function Must Have Documentation

```go
// Load creates a new Config instance with values from environment variables.
// If an environment variable is not set, a default value is used.
//
// Environment variables:
//   PORT - HTTP server port (default: "8080")
//   ENVIRONMENT - Application environment (default: "development")
//
// Example:
//   cfg := config.Load()
//   fmt.Println(cfg.Port) // "8080"
func Load() *Config {
```

### Every Exported Type Must Have Documentation

```go
// Config holds all application configuration.
// It follows the 12-factor app methodology.
type Config struct {
	// Port is the HTTP server port (default: "8080")
	Port string

	// Environment is the application environment
	Environment string
}
```

### Every Struct Field Must Have Documentation

```go
type HealthResponse struct {
	// Status indicates if the service is healthy
	Status string `json:"status" example:"healthy"`

	// Version is the API version
	Version string `json:"version" example:"1.0"`
}
```

---

## File Organization Standards

### Rule: One Responsibility Per File

❌ **Bad:** Everything in one file
```
main.go (500 lines - does everything)
```

✅ **Good:** Separated by responsibility
```
cmd/server/main.go        (30 lines - bootstrap only)
internal/server/server.go  (80 lines - server management)
internal/router/router.go  (50 lines - route setup)
internal/router/v1/health.go (40 lines - health handler)
```

### Rule: Tests Next to Code

```
internal/
├── server/
│   ├── server.go
│   └── server_test.go      ✅ Test file next to implementation
├── router/
│   ├── router.go
│   └── router_test.go
```

---

## Error Handling Standards

### Always Handle Errors Explicitly

❌ **Bad:**
```go
srv.Start() // Ignoring error
```

✅ **Good:**
```go
if err := srv.Start(); err != nil {
	log.Fatalf("Failed to start server: %v", err)
}
```

### Return Errors, Don't Panic

❌ **Bad:**
```go
func Load() *Config {
	if err != nil {
		panic(err) // Don't panic in library code
	}
}
```

✅ **Good:**
```go
func Load() (*Config, error) {
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return cfg, nil
}
```

---

## Testing Standards

### Every Package Must Have Tests

```
package_test.go    ✅ Required
```

### Minimum Coverage: 80%

```bash
go test -cover ./...
# Should show at least 80% coverage
```

### Test Function Naming

```go
func TestFunctionName_Scenario(t *testing.T) {
	// Test implementation
}

// Examples:
func TestHealthHandler_ReturnsHealthy(t *testing.T)
func TestServer_Start_FailsOnInvalidPort(t *testing.T)
func TestConfig_Load_UsesDefaults(t *testing.T)
```

---

## Code Style Standards

### Use gofmt and golangci-lint

```bash
# Format code
go fmt ./...

# Lint code
golangci-lint run
```

### Variable Naming

```go
// ✅ Good: Clear, descriptive names
config := config.Load()
healthResponse := HealthResponse{Status: "healthy"}

// ❌ Bad: Unclear abbreviations
cfg := config.Load()
hr := HealthResponse{Status: "healthy"}
```

### Function Naming

```go
// ✅ Good: Verb + Noun
func CreateUser()
func ValidateInput()
func FetchGameData()

// ❌ Bad: Unclear purpose
func DoStuff()
func Handle()
func Process()
```

---

## Dependency Management

### Use Interfaces for Dependencies

❌ **Bad:** Direct dependency
```go
type Service struct {
	db *sql.DB  // Concrete type
}
```

✅ **Good:** Interface dependency
```go
type Database interface {
	Query(query string) ([]Row, error)
}

type Service struct {
	db Database  // Interface - mockable
}
```

---

## Configuration Standards

### Use Environment Variables

```go
// ✅ Good: Configurable
port := os.Getenv("PORT")

// ❌ Bad: Hardcoded
port := "8080"
```

### Provide Defaults

```go
// ✅ Good: Sensible defaults
port := getEnv("PORT", "8080")

// ❌ Bad: Fails if not set
port := os.Getenv("PORT") // Empty string if not set
```

---

## Security Standards

### Never Commit Secrets

```bash
# .gitignore
.env
*.key
*.pem
secrets/
```

### Use Environment Variables for Secrets

```go
// ✅ Good
apiKey := os.Getenv("API_KEY")

// ❌ Bad
apiKey := "sk-1234567890abcdef" // Never hardcode secrets
```

---

## API Design Standards

### Use Proper HTTP Status Codes

```go
c.JSON(200, data)  // OK
c.JSON(201, data)  // Created
c.JSON(400, err)   // Bad Request
c.JSON(401, err)   // Unauthorized
c.JSON(404, err)   // Not Found
c.JSON(500, err)   // Internal Server Error
```

### Use Consistent Response Formats

```go
// ✅ Good: Consistent structure
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ❌ Bad: Inconsistent
// Sometimes: {"result": ...}
// Sometimes: {"data": ...}
// Sometimes: {"response": ...}
```

### Version Your APIs

```
✅ /api/v1/health
❌ /health
```

---

## Logging Standards

### Use Structured Logging

```go
// ✅ Good: Structured
log.Printf("Server started on port %s", port)

// ❌ Bad: Unstructured
fmt.Println("Server started")
```

### Log Levels

```go
log.Debug("Detailed debug information")
log.Info("Normal operational messages")
log.Warn("Warning conditions")
log.Error("Error conditions")
log.Fatal("Critical - application cannot continue")
```

---

## Performance Standards

### Use Context for Cancellation

```go
func (s *server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
```

### Avoid Goroutine Leaks

```go
// ✅ Good: Controlled shutdown
go func() {
	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}()

// Wait for signal
<-quit

// Cleanup goroutine
srv.Shutdown(ctx)
```

---

## Checklist for Every File

- [ ] Package documentation exists
- [ ] All exported functions have documentation
- [ ] All exported types have documentation
- [ ] All struct fields have documentation
- [ ] Follows SOLID principles
- [ ] Has corresponding test file
- [ ] Test coverage > 80%
- [ ] No hardcoded values (use config)
- [ ] Proper error handling
- [ ] Passes `go fmt`
- [ ] Passes `golangci-lint run`
- [ ] No secrets in code

---

## Code Review Standards

### Before Submitting PR

1. Run tests: `go test ./...`
2. Check coverage: `go test -cover ./...`
3. Format: `go fmt ./...`
4. Lint: `golangci-lint run`
5. Review documentation
6. Check for SOLID violations

### Code Review Checklist

- [ ] Code follows SOLID principles
- [ ] All functions are documented
- [ ] Tests are included
- [ ] No hardcoded values
- [ ] Error handling is proper
- [ ] No security issues
- [ ] Performance is acceptable

---

## Summary

**Every file must:**
- Follow SOLID principles
- Have complete documentation
- Have tests (>80% coverage)
- Pass linting
- Handle errors properly
- Use interfaces for dependencies
- Be production-ready

**No exceptions!**

---

**Next:** See [SOLID-REFACTORING-PLAN.md](./SOLID-REFACTORING-PLAN.md) for implementation details.

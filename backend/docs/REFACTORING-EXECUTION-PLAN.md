# Refactoring Execution Plan

## Goal
Transform current monolithic `main.go` into SOLID-compliant structure with crisp documentation.

---

## Documentation Style

### ❌ Overly Verbose (Don't do this)
```go
// Load creates a new Config instance with values from environment variables.
// If an environment variable is not set, a default value is used.
// This function is the primary way to initialize configuration in the application.
// It follows the 12-factor app methodology by using environment variables.
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

### ✅ Crisp & Precise (Do this)
```go
// Load reads config from env vars with defaults.
func Load() *Config {
```

---

## Atomic Steps

### Step 1: Create Config Package
**What:** Extract hardcoded port to configuration
**Files:** `internal/config/config.go`
**Changes:** Create config struct, environment loading
**Test:** Can we load config?

### Step 2: Create Server Package
**What:** Extract server lifecycle management
**Files:** `internal/server/server.go`
**Changes:** Server interface, graceful shutdown
**Test:** Can we start/stop server?

### Step 3: Create Router Package
**What:** Extract route setup logic
**Files:** `internal/router/router.go`
**Changes:** Route registration separated from main
**Test:** Are routes registered correctly?

### Step 4: Create V1 Routes Package
**What:** Organize v1 endpoints
**Files:** `internal/router/v1/routes.go`, `internal/router/v1/health.go`
**Changes:** v1 routes grouped, health handler separated
**Test:** Does `/api/v1/health` work?

### Step 5: Refactor Main
**What:** Simplify main to bootstrap only
**Files:** `cmd/server/main.go`
**Changes:** Use new packages, add graceful shutdown
**Test:** Does server start and stop gracefully?

### Step 6: Add Tests
**What:** Test each component
**Files:** `*_test.go` for each package
**Changes:** Unit tests with >80% coverage
**Test:** Do all tests pass?

### Step 7: Update Air Config
**What:** Ensure hot reload still works
**Files:** `.air.toml`
**Changes:** Update build command if needed
**Test:** Does `air` work?

---

## Detailed Atomic Steps

### Step 1: Config Package
```
Action: Create config package
Input: None
Output: internal/config/config.go
Dependencies: None
Time: 2 minutes
```

**Create:**
- `internal/config/config.go` - Config struct + Load function

**Documentation:**
- Package: "Package config manages application settings."
- Struct: "Config holds app settings."
- Function: "Load reads config from env with defaults."

**No changes to main.go yet.**

---

### Step 2: Server Package
```
Action: Create server abstraction
Input: Config from Step 1
Output: internal/server/server.go
Dependencies: config package
Time: 3 minutes
```

**Create:**
- `internal/server/server.go` - Server interface + implementation

**Documentation:**
- Package: "Package server manages HTTP server lifecycle."
- Interface: "Server handles HTTP server start/stop."
- Functions: One-line each

**No changes to main.go yet.**

---

### Step 3: Router Package
```
Action: Extract route setup
Input: None
Output: internal/router/router.go
Dependencies: gin
Time: 2 minutes
```

**Create:**
- `internal/router/router.go` - Setup function

**Documentation:**
- Package: "Package router configures HTTP routes."
- Function: "Setup creates configured Gin router."

**No changes to main.go yet.**

---

### Step 4: V1 Handlers
```
Action: Organize v1 endpoints
Input: None
Output: internal/router/v1/{routes.go, health.go}
Dependencies: gin
Time: 3 minutes
```

**Create:**
- `internal/router/v1/routes.go` - RegisterRoutes function
- `internal/router/v1/health.go` - HealthHandler

**Documentation:**
- Package: "Package v1 provides API v1 endpoints."
- Each handler: One-line purpose

**No changes to main.go yet.**

---

### Step 5: Refactor Main
```
Action: Simplify main.go
Input: All packages from Steps 1-4
Output: Updated cmd/server/main.go
Dependencies: All above packages
Time: 3 minutes
```

**Update:**
- `cmd/server/main.go` - Use new packages, add shutdown

**Documentation:**
- Package: "Package main bootstraps the Chess Coach API."
- Function: "main starts server with graceful shutdown."

**This is when everything comes together.**

---

### Step 6: Add Tests
```
Action: Create unit tests
Input: All packages
Output: *_test.go files
Dependencies: testing package
Time: 5 minutes
```

**Create:**
- `internal/config/config_test.go`
- `internal/server/server_test.go`
- `internal/router/v1/health_test.go`

**Skip router_test.go for now** (integration test)

---

### Step 7: Verify Air
```
Action: Test hot reload
Input: Refactored code
Output: Working air setup
Dependencies: None
Time: 1 minute
```

**Test:**
- Run `air`
- Make a change
- Verify auto-reload

---

## Execution Order

```
Step 1: Config        →  Can run independently
Step 2: Server        →  Needs config
Step 3: Router        →  Can run independently
Step 4: V1 Handlers   →  Can run independently
Step 5: Main          →  Needs all above
Step 6: Tests         →  Needs all above
Step 7: Air           →  Verification only
```

---

## After Each Step

### Checklist
- [ ] Files created
- [ ] Code compiles (`go build ./...`)
- [ ] Docs are crisp (1 line per item)
- [ ] No errors
- [ ] Ready for next step

### If Something Breaks
**STOP and ask me!** Don't proceed to next step.

---

## File Tree After Refactoring

```
backend/
├── cmd/
│   └── server/
│       └── main.go              (30 lines - bootstrap only)
├── internal/
│   ├── config/
│   │   ├── config.go            (40 lines)
│   │   └── config_test.go       (30 lines)
│   ├── server/
│   │   ├── server.go            (70 lines)
│   │   └── server_test.go       (40 lines)
│   └── router/
│       ├── router.go            (30 lines)
│       └── v1/
│           ├── routes.go        (20 lines)
│           ├── health.go        (30 lines)
│           └── health_test.go   (40 lines)
├── go.mod
├── go.sum
└── .air.toml
```

**Total new files:** 9
**Lines per file:** ~30-70 (small, focused)

---

## Time Estimate

| Step | Time | Cumulative |
|------|------|------------|
| Step 1: Config | 2 min | 2 min |
| Step 2: Server | 3 min | 5 min |
| Step 3: Router | 2 min | 7 min |
| Step 4: V1 | 3 min | 10 min |
| Step 5: Main | 3 min | 13 min |
| Step 6: Tests | 5 min | 18 min |
| Step 7: Air | 1 min | 19 min |

**Total:** ~20 minutes (with explanations)

---

## Current vs After

### Current (main.go)
```go
43 lines
- Creates router
- Defines routes
- Defines handlers
- Starts server
- No graceful shutdown
- No tests
- Hardcoded config
```

### After (main.go)
```go
30 lines
- Loads config
- Creates server
- Starts server
- Graceful shutdown
- Fully tested
- No hardcoded values
```

---

## Questions You Might Have

### Q: Can we skip a step?
**A:** No. Each step builds on previous ones.

### Q: Can we change the order?
**A:** Only Steps 1-4 can be reordered. Step 5+ must be in order.

### Q: What if I want different structure?
**A:** Tell me before we start Step 1.

### Q: How do we test as we go?
**A:** After each step, run `go build ./...` to verify compilation.

### Q: What about Swagger?
**A:** We'll add it AFTER refactoring is complete.

---

## Ready to Start?

**Say "start" and I'll begin with Step 1: Config Package**

Or ask any questions about the plan!

---

## Notes

- One step at a time
- I'll explain each change
- You can ask questions anytime
- We can pause between steps
- Crisp docs (1 line per item)
- SOLID principles throughout

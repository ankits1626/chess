# Step 19: Package Reorganization - Feature-Based Structure

**Objective**: Reorganize codebase from layer-based to feature-based structure for better scalability and maintainability.

**Time**: 30 minutes
**Difficulty**: Medium

---

## Current vs Target Structure

### Current Structure (Layer-Based)
```
internal/
├── router/
│   └── v1/
│       ├── models.go          # Mixed: Only User DTOs
│       ├── user_handler.go    # User HTTP handlers
│       └── routes.go          # Route registration
└── repository/
    ├── user_repository.go
    └── game_repository.go
```

**Issues:**
- ❌ `models.go` will become bloated with User, Game, Move DTOs
- ❌ No clear feature boundaries
- ❌ Handler and DTOs separated (cognitive load)
- ❌ Hard to navigate as project grows

### Target Structure (Feature-Based with Versioning)
```
internal/
├── handler/                    # HTTP layer
│   └── v1/                    # Version 1 namespace
│       ├── user/
│       │   ├── handler.go     # UserHandler + methods
│       │   └── dto.go         # User DTOs (request/response)
│       ├── game/
│       │   ├── handler.go     # GameHandler + methods
│       │   └── dto.go         # Game DTOs
│       └── health/
│           └── handler.go     # Health check
├── repository/                 # Data access layer (shared across versions)
│   ├── user_repository.go
│   └── game_repository.go
└── router/
    ├── router.go              # Main router setup
    └── v1_routes.go           # v1 route registration
```

**Benefits:**
- ✅ Clear feature boundaries (everything User-related in `handler/v1/user/`)
- ✅ DTOs live next to handlers that use them
- ✅ Version namespace for future API versions (add `handler/v2/` when needed)
- ✅ Repository layer shared across versions (no duplication)
- ✅ Easy to add new features (just add `handler/v1/game/`)
- ✅ Follows Go standard library patterns (net/http, database/sql)
- ✅ Team-friendly (work on different features independently)

---

## Phase 1: Create New Directory Structure

**Time**: 2 minutes

### Step 1.1: Create Handler Directories

```bash
cd backend/internal

# Create versioned feature-based handler directories
mkdir -p handler/v1/user
mkdir -p handler/v1/game
mkdir -p handler/v1/health
```

**Verify:**
```bash
ls -la handler/v1/
# Should show: user/ game/ health/
```

---

## Phase 2: Move User Handler Code

**Time**: 8 minutes

### Step 2.1: Create User Handler File

**File**: `internal/handler/v1/user/handler.go`

**Action**: Move content from `internal/router/v1/user_handler.go`

```go
// Package user handles user-related HTTP requests for API v1.
package user

import (
	"net/http"

	"github.com/ankits1626/chess-coach-backend/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles user endpoints.
type Handler struct {
	repo *repository.UserRepository
}

// NewHandler creates user handler.
func NewHandler(repo *repository.UserRepository) *Handler {
	return &Handler{repo: repo}
}

// List lists users with pagination.
// @Summary List users
// @Description Get paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} Response
// @Router /users [get]
func (h *Handler) List(c *gin.Context) {
	var params struct {
		Limit  int32 `form:"limit" binding:"omitempty,min=1,max=100"`
		Offset int32 `form:"offset" binding:"omitempty,min=0"`
	}

	if params.Limit == 0 {
		params.Limit = 10
	}

	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := h.repo.List(c.Request.Context(), params.Limit, params.Offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, ToResponses(users))
}

// Get retrieves user by ID.
// @Summary Get user
// @Description Get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} Response
// @Router /users/{id} [get]
func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, ToResponse(user))
}

// Create creates new user.
// @Summary Create user
// @Description Create new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateRequest true "User data"
// @Success 201 {object} Response
// @Router /users [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Rating == 0 {
		req.Rating = 1200 // Default rating
	}

	user, err := h.repo.Create(c.Request.Context(), req.Username, req.Email, req.Rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, ToResponse(user))
}

// Update updates existing user.
// @Summary Update user
// @Description Update user details
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body UpdateRequest true "User data"
// @Success 200 {object} Response
// @Router /users/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.repo.Update(c.Request.Context(), id, req.Username, req.Email, req.Rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, ToResponse(user))
}

// Delete deletes user.
// @Summary Delete user
// @Description Delete user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204
// @Router /users/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}
```

**Key Changes:**
- Renamed `UserHandler` → `Handler` (package name provides context)
- Renamed `UserResponse` → `Response` (in DTOs)
- Renamed `ToUserResponse` → `ToResponse`
- Renamed `ToUserResponses` → `ToResponses`

---

### Step 2.2: Create User DTO File

**File**: `internal/handler/v1/user/dto.go`

**Action**: Move content from `internal/router/v1/models.go`

```go
// Package user handles user-related HTTP requests for API v1.
package user

import (
	"time"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/google/uuid"
)

// CreateRequest represents user creation request.
type CreateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Rating   int32  `json:"rating,omitempty"`
}

// UpdateRequest represents user update request.
type UpdateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Rating   int32  `json:"rating,omitempty"`
}

// Response represents user in API responses.
type Response struct {
	ID        uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username  string    `json:"username" example:"johndoe"`
	Email     string    `json:"email" example:"john@example.com"`
	Rating    int32     `json:"rating" example:"1500"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// ToResponse converts database.User to Response.
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

// ToResponses converts slice of database.User to Response slice.
func ToResponses(users []database.User) []Response {
	responses := make([]Response, len(users))
	for i, u := range users {
		responses[i] = ToResponse(u)
	}
	return responses
}
```

**Naming Convention:**
- Request types: `CreateRequest`, `UpdateRequest`
- Response types: `Response` (singular)
- Converters: `ToResponse()`, `ToResponses()`

---

## Phase 3: Create Health Handler

**Time**: 3 minutes

### Step 3.1: Move Health Check Handler

**File**: `internal/handler/v1/health/handler.go`

```go
// Package health provides health check endpoints.
package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents health check response.
type Response struct {
	Status  string `json:"status" example:"ok"`
	Service string `json:"service" example:"chess-coach-api"`
}

// Check handles health check endpoint.
// @Summary Health check
// @Description Check if service is healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router /health [get]
func Check(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Status:  "ok",
		Service: "chess-coach-api",
	})
}
```

---

## Phase 4: Update Router

**Time**: 7 minutes

### Step 4.1: Update Routes File

**File**: `internal/router/v1_routes.go` (NEW FILE)

**Create this file:**

```go
// Package router provides HTTP routing setup.
package router

import (
	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/handler/v1/health"
	"github.com/ankits1626/chess-coach-backend/internal/handler/v1/user"
	"github.com/ankits1626/chess-coach-backend/internal/repository"
	"github.com/gin-gonic/gin"
)

// RegisterV1Routes registers API v1 routes.
func RegisterV1Routes(r *gin.Engine, db *database.DB) {
	v1 := r.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", health.Check)

		// User routes
		userRepo := repository.NewUserRepository(db)
		userHandler := user.NewHandler(userRepo)

		users := v1.Group("/users")
		{
			users.GET("", userHandler.List)
			users.GET("/:id", userHandler.Get)
			users.POST("", userHandler.Create)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}
	}
}
```

**Key Changes:**
- Import `internal/handler/v1/user` instead of `internal/router/v1`
- Import `internal/handler/v1/health`
- Use `user.NewHandler()` instead of `v1.NewUserHandler()`
- Use `health.Check` instead of `v1.HealthHandler`
- Renamed function to `RegisterV1Routes` (version-specific)

---

### Step 4.2: Update Router Setup File

**File**: `internal/router/router.go`

**Update to:**

```go
// Package router provides HTTP routing setup.
package router

import (
	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/ankits1626/chess-coach-backend/docs"
)

// Setup creates and configures the Gin router.
func Setup(db *database.DB) *gin.Engine {
	r := gin.Default()

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Register API v1 routes
	RegisterV1Routes(r, db)

	return r
}
```

**Key Change:**
- Changed `RegisterRoutes(r, db)` to `RegisterV1Routes(r, db)`

---

## Phase 5: Clean Up Old Files

**Time**: 2 minutes

### Step 5.1: Remove Old Files

```bash
cd backend/internal

# Remove old v1 directory
rm -rf router/v1/

# Verify
ls -la router/
# Should only show: router.go, routes.go
```

---

## Phase 6: Update Imports and Rebuild

**Time**: 8 minutes

### Step 6.1: Tidy Dependencies

```bash
cd backend
go mod tidy
```

### Step 6.2: Regenerate Swagger Docs

```bash
cd backend
swag init -g cmd/server/main.go
```

**Expected output:**
```
Generate swagger docs....
Generating health.Response
Generating user.Response
Generating user.CreateRequest
Generating user.UpdateRequest
create docs.go at docs/docs.go
```

### Step 6.3: Build and Test

```bash
# Build
go build -o bin/server cmd/server/main.go

# Should compile without errors
```

### Step 6.4: Run Application

```bash
# Start database (if not running)
docker compose up -d postgres

# Run server
./bin/server
```

**Expected output:**
```
INFO: Database connected successfully
INFO: Server started on port 8080
INFO: Swagger UI: http://localhost:8080/swagger/index.html
```

---

## Phase 7: Verify Endpoints

**Time**: N/A (Testing)

### Step 7.1: Test Swagger UI

1. Open browser: `http://localhost:8080/swagger/index.html`
2. Verify you see:
   - ✅ `GET /api/v1/health`
   - ✅ `GET /api/v1/users`
   - ✅ `GET /api/v1/users/{id}`
   - ✅ `POST /api/v1/users`
   - ✅ `PUT /api/v1/users/{id}`
   - ✅ `DELETE /api/v1/users/{id}`

### Step 7.2: Test Health Check

```bash
curl http://localhost:8080/api/v1/health
```

**Expected:**
```json
{
  "status": "ok",
  "service": "chess-coach-api"
}
```

### Step 7.3: Test User Endpoints

```bash
# Create user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "rating": 1500
  }'

# List users
curl http://localhost:8080/api/v1/users
```

---

## Final Directory Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── app/
│   │   └── app.go
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── connection.go
│   │   ├── db.go (generated)
│   │   ├── games.sql.go (generated)
│   │   ├── models.go (generated)
│   │   └── users.sql.go (generated)
│   ├── handler/              # ✨ NEW - Versioned handlers
│   │   └── v1/               # ✨ Version 1 namespace
│   │       ├── health/
│   │       │   └── handler.go
│   │       └── user/
│   │           ├── dto.go
│   │           └── handler.go
│   ├── logger/
│   │   └── logger.go
│   ├── repository/           # Shared across versions
│   │   ├── game_repository.go
│   │   └── user_repository.go
│   ├── router/
│   │   ├── router.go
│   │   └── v1_routes.go      # ✨ Version-specific routes
│   └── server/
│       └── server.go
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
└── go.mod
```

---

## Benefits Achieved

### Before (Layer-Based)
❌ Mixed concerns in `models.go`
❌ Handlers and DTOs separated
❌ Hard to navigate
❌ Doesn't scale well

### After (Feature-Based with Versioning)
✅ Clear feature boundaries
✅ DTOs next to handlers
✅ Version namespaces for backward compatibility
✅ Easy to find user-related code (`handler/v1/user/`)
✅ Scales with new features (just add `handler/v1/game/`)
✅ Scales with new versions (just add `handler/v2/` when needed)
✅ Repository layer shared across versions (DRY principle)
✅ Team-friendly (parallel development)

---

## Adding New Features

To add a new feature (e.g., Game handler) to v1:

```bash
# 1. Create directory
mkdir -p internal/handler/v1/game

# 2. Create handler.go
touch internal/handler/v1/game/handler.go

# 3. Create dto.go
touch internal/handler/v1/game/dto.go

# 4. Implement handlers and DTOs

# 5. Register routes in internal/router/v1_routes.go
```

---

## API Versioning

### Why Versioning Matters

This structure is designed for API versioning from day one:

- **Current**: All handlers in `handler/v1/`
- **Future**: Add `handler/v2/` when breaking changes needed
- **Routes**: Separate route files (`v1_routes.go`, `v2_routes.go`)
- **Shared Logic**: Repository and service layers work across versions

### When to Create v2

Create a new version when you have **breaking changes**:

- ❌ Field removed from response
- ❌ Field renamed or type changed
- ❌ Endpoint URL changed
- ❌ Required parameter added
- ✅ New optional field (add to v1)
- ✅ New endpoint (add to v1)

### Adding v2 (Future)

```bash
# 1. Create v2 handler directory
mkdir -p internal/handler/v2/user

# 2. Copy or create new handlers
cp -r internal/handler/v1/user/* internal/handler/v2/user/

# 3. Modify v2 DTOs with breaking changes
# Edit internal/handler/v2/user/dto.go

# 4. Create v2 routes
touch internal/router/v2_routes.go

# 5. Register v2 routes in router.go
# Add: RegisterV2Routes(r, db)
```

**See [api-versioning-strategy.md](./api-versioning-strategy.md) for comprehensive versioning guide.**

---

## Naming Conventions Summary

### Package Structure
- `internal/handler/v1/{feature}/handler.go` - HTTP handlers (v1)
- `internal/handler/v1/{feature}/dto.go` - Request/Response DTOs (v1)
- `internal/repository/{feature}_repository.go` - Data access (shared)

### Type Naming (within feature package)
- Handler: `Handler` (not `UserHandler` - package provides context)
- Requests: `CreateRequest`, `UpdateRequest`, `ListRequest`
- Response: `Response` (singular - represents single entity)
- Converters: `ToResponse()`, `ToResponses()`

### Import Convention
```go
// Use aliased imports for clarity
import (
    userv1 "github.com/ankits1626/chess-coach-backend/internal/handler/v1/user"
    userv2 "github.com/ankits1626/chess-coach-backend/internal/handler/v2/user"
)

// Usage
v1Handler := userv1.NewHandler(repo)
v2Handler := userv2.NewHandler(repo)
```

### Example:
```go
// internal/handler/v1/user/handler.go
package user

type Handler struct { ... }
func NewHandler(...) *Handler { ... }

// internal/handler/v1/user/dto.go
package user

type CreateRequest struct { ... }
type Response struct { ... }
func ToResponse(...) Response { ... }
```

---

## Troubleshooting

### Issue: Import errors after moving files

**Solution:**
```bash
cd backend
go mod tidy
go clean -cache
```

### Issue: Swagger not showing new structure

**Solution:**
```bash
cd backend
swag init -g cmd/server/main.go --parseDependency --parseInternal
```

### Issue: Routes not working

**Check:**
1. Verify imports in `internal/router/routes.go`
2. Ensure handler methods are exported (capitalized)
3. Check router registration in `RegisterRoutes()`

---

## Completion Checklist

- [ ] Created `internal/handler/user/` directory
- [ ] Created `internal/handler/health/` directory
- [ ] Moved handler code to `handler.go`
- [ ] Moved DTOs to `dto.go`
- [ ] Updated `internal/router/routes.go`
- [ ] Removed old `internal/router/v1/` directory
- [ ] Ran `go mod tidy`
- [ ] Regenerated Swagger docs
- [ ] Application builds successfully
- [ ] All endpoints working in Swagger UI
- [ ] Health check returns correct response
- [ ] User CRUD operations work

---

**Next Steps:**
- Add Game handler in `internal/handler/game/`
- Add Move handler in `internal/handler/move/`
- Consider adding middleware in `internal/middleware/`
- Add service layer for complex business logic

---

**Last Updated**: 2025-11-01

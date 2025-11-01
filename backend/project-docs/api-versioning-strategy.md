# API Versioning Strategy

**A comprehensive guide to API versioning in Go REST APIs**

---

## Table of Contents
1. [Why Version APIs?](#why-version-apis)
2. [Versioning Strategies](#versioning-strategies)
3. [Recommended Approach](#recommended-approach)
4. [Implementation Options](#implementation-options)
5. [Migration Path](#migration-path)
6. [Best Practices](#best-practices)

---

## Why Version APIs?

### The Problem
```
Your Chess Coach API v1.0:
- Users have { username, email, rating }

Six months later, you want to add:
- Users have { username, email, elo_rating, rapid_rating, blitz_rating }

Breaking change! Old clients expect "rating", new API has multiple ratings.
```

### Without Versioning
❌ Break existing mobile apps
❌ Force all clients to update simultaneously
❌ No rollback path
❌ Integration nightmares

### With Versioning
✅ Old clients use `/api/v1/users` (still works)
✅ New clients use `/api/v2/users` (new features)
✅ Gradual migration path
✅ Backward compatibility maintained

---

## Versioning Strategies

### 1. URI Versioning (Most Common)

**Format:** `/api/v1/users`, `/api/v2/users`

```
GET /api/v1/users/123
GET /api/v2/users/123
```

**Pros:**
- ✅ Clear and explicit
- ✅ Easy to test (just change URL)
- ✅ Simple routing
- ✅ Cache-friendly (different URLs)

**Cons:**
- ❌ Pollutes URI space
- ❌ Less RESTful (resource identity changes)

**Best for:** Public APIs, mobile apps, third-party integrations

---

### 2. Header Versioning

**Format:** `Accept: application/vnd.chess-coach.v1+json`

```
GET /api/users/123
Accept: application/vnd.chess-coach.v1+json
```

**Pros:**
- ✅ Clean URIs
- ✅ More RESTful (same resource, different representation)
- ✅ Follows HTTP standards

**Cons:**
- ❌ Harder to test (need to set headers)
- ❌ Not browser-friendly
- ❌ Cache complications

**Best for:** Internal APIs, sophisticated clients

---

### 3. Query Parameter Versioning

**Format:** `/api/users?version=1`

```
GET /api/users/123?version=1
```

**Pros:**
- ✅ Clean base URIs
- ✅ Easy to test

**Cons:**
- ❌ Query params have other semantic meanings
- ❌ Easy to forget
- ❌ Not standard practice

**Best for:** Rarely used (not recommended)

---

### 4. Content Negotiation

**Format:** Custom media types

```
GET /api/users/123
Accept: application/json; version=1
```

**Pros:**
- ✅ HTTP-compliant
- ✅ Flexible

**Cons:**
- ❌ Complex implementation
- ❌ Not widely understood

**Best for:** Advanced use cases

---

## Recommended Approach

### For Chess Coach Backend: **URI Versioning**

**Why?**
1. **Clarity** - Version is immediately visible in URL
2. **Simplicity** - Easy to implement and understand
3. **Testing** - Just change URL in tests/Postman/Swagger
4. **Mobile-friendly** - Apps can hardcode version
5. **Industry standard** - Used by Stripe, Twitter, GitHub APIs

---

## Implementation Options

### Option A: Version in Package Structure (Recommended)

```
internal/
├── handler/
│   ├── v1/                    # Version 1 handlers
│   │   ├── user/
│   │   │   ├── handler.go
│   │   │   └── dto.go
│   │   ├── game/
│   │   │   ├── handler.go
│   │   │   └── dto.go
│   │   └── health/
│   │       └── handler.go
│   └── v2/                    # Version 2 handlers (when needed)
│       ├── user/
│       │   ├── handler.go     # Can have different logic
│       │   └── dto.go         # Can have different fields
│       └── game/
│           ├── handler.go
│           └── dto.go
├── repository/
│   ├── user_repository.go     # Shared across versions
│   └── game_repository.go
└── router/
    ├── router.go
    ├── v1_routes.go           # v1 route registration
    └── v2_routes.go           # v2 route registration (future)
```

**Benefits:**
- ✅ Clear version boundaries
- ✅ Can have completely different logic per version
- ✅ No version conditionals in code
- ✅ Easy to deprecate (just remove v1 folder)
- ✅ Repository layer shared (no duplication)

**Routing:**
```go
// internal/router/v1_routes.go
func RegisterV1Routes(r *gin.Engine, db *database.DB) {
    v1 := r.Group("/api/v1")
    {
        // User routes
        userRepo := repository.NewUserRepository(db)
        userHandler := userv1.NewHandler(userRepo)

        v1.GET("/users", userHandler.List)
        v1.POST("/users", userHandler.Create)
    }
}

// internal/router/v2_routes.go (future)
func RegisterV2Routes(r *gin.Engine, db *database.DB) {
    v2 := r.Group("/api/v2")
    {
        // User routes with new features
        userRepo := repository.NewUserRepository(db)
        userHandler := userv2.NewHandler(userRepo)

        v2.GET("/users", userHandler.List)  // Different response format
        v2.POST("/users", userHandler.Create)
    }
}

// internal/router/router.go
func Setup(db *database.DB) *gin.Engine {
    r := gin.Default()

    RegisterV1Routes(r, db)  // /api/v1/*
    RegisterV2Routes(r, db)  // /api/v2/* (when available)

    return r
}
```

---

### Option B: Version in Handler Methods (Not Recommended)

```go
// Single handler with version logic
func (h *UserHandler) Get(c *gin.Context) {
    version := c.Param("version")  // from /api/:version/users

    if version == "v1" {
        // v1 logic
        c.JSON(200, UserV1Response{...})
    } else if version == "v2" {
        // v2 logic
        c.JSON(200, UserV2Response{...})
    }
}
```

**Problems:**
- ❌ Conditional logic everywhere
- ❌ Hard to maintain
- ❌ Difficult to test
- ❌ Messy code

---

## Migration Path

### Current State
```
/api/v1/health
/api/v1/users
/api/v1/users/{id}
```

### Recommended Structure After Reorganization

```
backend/internal/
├── handler/
│   └── v1/                        # Version 1 namespace
│       ├── user/
│       │   ├── handler.go
│       │   └── dto.go
│       ├── game/
│       │   ├── handler.go
│       │   └── dto.go
│       └── health/
│           └── handler.go
├── repository/                    # Shared (version-agnostic)
│   ├── user_repository.go
│   └── game_repository.go
└── router/
    ├── router.go                  # Main setup
    ├── v1_routes.go               # v1 registration
    └── v2_routes.go               # v2 registration (future)
```

### Imports
```go
// Clean imports
import (
    userv1 "github.com/ankits1626/chess-coach-backend/internal/handler/v1/user"
    userv2 "github.com/ankits1626/chess-coach-backend/internal/handler/v2/user"
)

// Usage
v1Handler := userv1.NewHandler(repo)
v2Handler := userv2.NewHandler(repo)
```

---

## Best Practices

### 1. Semantic Versioning

Use **major versions only** in API URLs:
- ✅ `/api/v1/users` - Major version
- ❌ `/api/v1.2.3/users` - Too granular

**When to bump version:**
- **v1 → v2**: Breaking changes (field removed, renamed, type changed)
- **v1**: Non-breaking changes (new optional fields, new endpoints)

### 2. Version Support Policy

Define clear support windows:
```
v1: Supported until 2026-01-01 (2 years)
v2: Current version (active development)
v3: Beta (testing phase)
```

**Communication:**
- Announce deprecation 6-12 months in advance
- Use `Deprecation` HTTP header
- Document migration guide

### 3. Default Version

```go
// Redirect unversioned requests to latest stable
r.GET("/api/users", func(c *gin.Context) {
    c.Redirect(http.StatusMovedPermanently, "/api/v2/users")
})
```

### 4. Version in Response Headers

```go
func (h *Handler) List(c *gin.Context) {
    c.Header("API-Version", "v1")
    c.Header("Deprecation", "Sun, 01 Jan 2026 00:00:00 GMT")

    // ... handler logic
}
```

### 5. Shared Business Logic

**Keep versions DRY:**
```
internal/
├── handler/
│   ├── v1/user/
│   │   ├── handler.go          # HTTP layer only
│   │   └── dto.go              # v1 DTOs
│   └── v2/user/
│       ├── handler.go          # HTTP layer only
│       └── dto.go              # v2 DTOs (different fields)
├── service/                     # ✨ Shared business logic
│   └── user_service.go         # Reused by v1 and v2
└── repository/                  # ✨ Shared data access
    └── user_repository.go
```

**Example:**
```go
// internal/service/user_service.go
type UserService struct {
    repo *repository.UserRepository
}

func (s *UserService) GetUser(ctx context.Context, id uuid.UUID) (database.User, error) {
    // Business logic here (validation, caching, etc.)
    return s.repo.GetByID(ctx, id)
}

// internal/handler/v1/user/handler.go
func (h *Handler) Get(c *gin.Context) {
    user, err := h.service.GetUser(c.Request.Context(), id)
    // Convert to v1 response
    c.JSON(200, toV1Response(user))
}

// internal/handler/v2/user/handler.go
func (h *Handler) Get(c *gin.Context) {
    user, err := h.service.GetUser(c.Request.Context(), id)
    // Convert to v2 response (different fields)
    c.JSON(200, toV2Response(user))
}
```

---

## Example: Breaking Change Scenario

### V1 API (Current)
```go
// GET /api/v1/users/123
{
  "id": "123",
  "username": "john",
  "rating": 1500         // Single rating
}
```

### V2 API (New - Multiple Ratings)
```go
// GET /api/v2/users/123
{
  "id": "123",
  "username": "john",
  "ratings": {           // Breaking: Changed structure
    "classical": 1500,
    "rapid": 1400,
    "blitz": 1600
  }
}
```

### Implementation
```go
// internal/handler/v1/user/dto.go
type Response struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Rating   int32  `json:"rating"`
}

func ToResponse(u database.User) Response {
    return Response{
        ID:       u.ID.String(),
        Username: u.Username,
        Rating:   u.ClassicalRating,  // Use classical as default
    }
}

// internal/handler/v2/user/dto.go
type Response struct {
    ID       string  `json:"id"`
    Username string  `json:"username"`
    Ratings  Ratings `json:"ratings"`
}

type Ratings struct {
    Classical int32 `json:"classical"`
    Rapid     int32 `json:"rapid"`
    Blitz     int32 `json:"blitz"`
}

func ToResponse(u database.User) Response {
    return Response{
        ID:       u.ID.String(),
        Username: u.Username,
        Ratings: Ratings{
            Classical: u.ClassicalRating,
            Rapid:     u.RapidRating,
            Blitz:     u.BlitzRating,
        },
    }
}
```

**Result:**
- Old mobile apps keep working with `/api/v1/users`
- New web app uses `/api/v2/users` with rich ratings
- Shared repository, different presentations

---

## Swagger Documentation

### Document All Versions
```go
// cmd/server/main.go

// @title Chess Coach API
// @version 2.0
// @description API for chess game analysis and coaching

// @contact.name API Support
// @contact.email support@chess-coach.com

// @host localhost:8080
// @BasePath /api/v2

// @tag.name v1
// @tag.description Deprecated - Version 1 (supported until 2026-01-01)

// @tag.name v2
// @tag.description Current - Version 2 (active development)
```

### Version-Specific Annotations
```go
// internal/handler/v1/user/handler.go

// List lists users with pagination (v1).
// @Summary List users (v1 - DEPRECATED)
// @Description Get paginated list of users - Version 1
// @Tags users,v1
// @Deprecated true
// @Router /v1/users [get]

// internal/handler/v2/user/handler.go

// List lists users with pagination (v2).
// @Summary List users (v2)
// @Description Get paginated list of users with enhanced ratings
// @Tags users,v2
// @Router /v2/users [get]
```

---

## Testing Multiple Versions

```go
// tests/user_test.go
func TestUserV1Get(t *testing.T) {
    resp := makeRequest("GET", "/api/v1/users/123")

    var user userv1.Response
    json.Unmarshal(resp.Body, &user)

    assert.NotNil(t, user.Rating)  // v1 has single rating
}

func TestUserV2Get(t *testing.T) {
    resp := makeRequest("GET", "/api/v2/users/123")

    var user userv2.Response
    json.Unmarshal(resp.Body, &user)

    assert.NotNil(t, user.Ratings.Classical)  // v2 has ratings object
}
```

---

## Summary

### Recommended Approach for Chess Coach

1. **Use URI versioning**: `/api/v1/`, `/api/v2/`
2. **Package structure**: `internal/handler/v1/`, `internal/handler/v2/`
3. **Share repositories**: Keep data access layer version-agnostic
4. **Separate route files**: `v1_routes.go`, `v2_routes.go`
5. **Clear deprecation policy**: 2-year support window

### Folder Structure
```
internal/
├── handler/
│   └── v1/              # Version namespace
│       ├── user/        # Feature
│       ├── game/
│       └── health/
├── repository/          # Shared
└── router/
    ├── router.go        # Main
    └── v1_routes.go     # Version-specific
```

### When to Create v2
- Breaking changes needed
- Major feature overhaul
- API redesign required
- Sufficient v1 adoption

### Migration Checklist
- [ ] Announce deprecation (6-12 months ahead)
- [ ] Document differences
- [ ] Provide migration guide
- [ ] Create v2 handlers/DTOs
- [ ] Register v2 routes
- [ ] Update Swagger docs
- [ ] Update client SDKs
- [ ] Monitor v1 usage
- [ ] Deprecate v1 after grace period

---

**Last Updated**: 2025-11-01

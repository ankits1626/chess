# Package Reorganization & API Versioning - Summary

**Quick reference guide for the Chess Coach backend restructuring**

---

## What Changed?

### Before
```
internal/router/v1/
├── models.go          # Mixed concerns
├── user_handler.go
└── routes.go
```

### After
```
internal/
├── handler/
│   └── v1/           # Version namespace
│       ├── user/     # Feature grouping
│       │   ├── handler.go
│       │   └── dto.go
│       ├── game/
│       └── health/
└── router/
    ├── router.go
    └── v1_routes.go  # Version-specific routes
```

---

## Key Benefits

1. **Feature-Based Organization** - Everything related to Users is in `handler/v1/user/`
2. **API Versioning Ready** - Can add `handler/v2/` when needed
3. **Cleaner Imports** - `import userv1 "internal/handler/v1/user"`
4. **Scalable** - Add new features without touching existing code
5. **DRY Principle** - Repository layer shared across all versions

---

## Quick Commands

### Add New Feature to v1
```bash
mkdir -p internal/handler/v1/game
touch internal/handler/v1/game/{handler,dto}.go
# Implement and register in v1_routes.go
```

### Add v2 API (Future)
```bash
mkdir -p internal/handler/v2/user
cp -r internal/handler/v1/user/* internal/handler/v2/user/
# Modify v2 DTOs with breaking changes
# Create v2_routes.go
```

---

## Import Pattern

```go
import (
    userv1 "github.com/ankits1626/chess-coach-backend/internal/handler/v1/user"
    gamev1 "github.com/ankits1626/chess-coach-backend/internal/handler/v1/game"
)

userHandler := userv1.NewHandler(userRepo)
gameHandler := gamev1.NewHandler(gameRepo)
```

---

## Naming Conventions

### Within Feature Packages
- `Handler` (not `UserHandler` - package provides context)
- `CreateRequest`, `UpdateRequest` (request DTOs)
- `Response` (response DTO)
- `ToResponse()`, `ToResponses()` (converters)

### Files
- `handler.go` - HTTP handlers
- `dto.go` - Request/response types

---

## API Versioning Strategy

### When to Create New Version

**Breaking Changes (need v2):**
- Remove field from response
- Rename field or change type
- Change endpoint URL
- Add required parameter

**Non-Breaking Changes (stay in v1):**
- Add optional field
- Add new endpoint
- Add query parameter with default

### Version Support
- Support each version for **2 years**
- Announce deprecation **6-12 months** in advance
- Document migration path clearly

---

## Documentation

- **Full Migration Guide**: [step-19-package-reorganization.md](./step-19-package-reorganization.md)
- **Versioning Strategy**: [api-versioning-strategy.md](./api-versioning-strategy.md)
- **Go Fundamentals**: [golang-fundamentals.md](./golang-fundamentals.md)

---

## Example: User Handler Structure

```
internal/handler/v1/user/
├── handler.go
│   ├── type Handler struct
│   ├── func NewHandler() *Handler
│   ├── func (h *Handler) List(c *gin.Context)
│   ├── func (h *Handler) Get(c *gin.Context)
│   ├── func (h *Handler) Create(c *gin.Context)
│   ├── func (h *Handler) Update(c *gin.Context)
│   └── func (h *Handler) Delete(c *gin.Context)
└── dto.go
    ├── type CreateRequest struct
    ├── type UpdateRequest struct
    ├── type Response struct
    ├── func ToResponse(u database.User) Response
    └── func ToResponses(users []database.User) []Response
```

---

## Testing

```bash
# Build
go build -o bin/server cmd/server/main.go

# Run
docker compose up -d postgres
./bin/server

# Test
curl http://localhost:8080/api/v1/health
curl http://localhost:8080/api/v1/users
open http://localhost:8080/swagger/index.html
```

---

## Checklist

- [ ] Read [step-19-package-reorganization.md](./step-19-package-reorganization.md)
- [ ] Understand [api-versioning-strategy.md](./api-versioning-strategy.md)
- [ ] Create `internal/handler/v1/` directories
- [ ] Move user handler to new structure
- [ ] Move health handler to new structure
- [ ] Update router files
- [ ] Remove old `internal/router/v1/` directory
- [ ] Run `go mod tidy`
- [ ] Regenerate Swagger docs
- [ ] Test all endpoints
- [ ] Update team documentation

---

**Last Updated**: 2025-11-01

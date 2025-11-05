# Go Fundamentals Learning Plan
## Project-Focused: Chess Coach Backend Mastery

**Goal**: Learn Go fundamentals through progressive, hands-on lessons that build toward understanding the chess-coach backend project.

**Learning Method**:
- Each lesson has its own folder with runnable examples
- Complete exercises and ask questions as you go
- Create a summary after each lesson
- Build progressively toward backend concepts

---

## Phase 1: Go Basics (Lessons 1-5)
*Foundation for reading and writing Go code*

### Lesson 1: Hello World & Basic Syntax
**Folder**: `01-hello-world/`

**Concepts**:
- Package declaration and `main` function
- Imports and standard library
- Variables, constants, and types (string, int, bool)
- `fmt` package for printing
- Comments and documentation

**Why This Matters**:
- Every Go file starts with `package`
- Backend uses `package main` in `cmd/server/main.go`
- Understanding imports is critical for reading the codebase

**Exercise**: Write programs that print formatted output

---

### Lesson 2: Functions & Error Handling
**Folder**: `02-functions-errors/`

**Concepts**:
- Function declarations and parameters
- Multiple return values
- Error return pattern (`result, err := doSomething()`)
- Error checking with `if err != nil`
- `defer` for cleanup
- Error wrapping with `fmt.Errorf("%w", err)`

**Why This Matters**:
- Go doesn't have exceptions - uses explicit error returns
- Backend has error handling in every handler and repository
- `defer` used for database connection cleanup
- Example: `config/config.go`, all handler functions

**Exercise**: Write functions that return errors and handle them properly

---

### Lesson 3: Structs & Methods
**Folder**: `03-structs-methods/`

**Concepts**:
- Defining structs
- Creating struct instances
- Pointer vs value receivers
- Methods on structs
- Struct embedding (composition)
- Zero values

**Why This Matters**:
- Backend uses structs everywhere: `Config`, `Handler`, `Repository`, `DB`
- Pointer receivers allow mutation
- Embedding used in `database.DB` struct
- Example: `handler/v1/user/handler.go`, `repository/user_repository.go`

**Exercise**: Create structs representing domain models (User, Game, Move)

---

### Lesson 4: Interfaces & Abstraction
**Folder**: `04-interfaces/`

**Concepts**:
- Interface declarations
- Implicit interface satisfaction
- Empty interface `interface{}`
- Type assertions
- Dependency injection pattern

**Why This Matters**:
- Backend uses interfaces for `Logger` and `Server`
- Enables testing with mocks
- Repository pattern uses interfaces
- Example: `logger/logger.go`, `server/server.go`

**Exercise**: Create interfaces and multiple implementations

---

### Lesson 5: Packages & Modules
**Folder**: `05-packages-modules/`

**Concepts**:
- Package organization
- Exported vs unexported (public vs private)
- go.mod and go.sum
- Importing local packages
- Third-party dependencies
- `internal/` directory convention

**Why This Matters**:
- Backend organized into packages: `handler`, `repository`, `database`
- Capital letters export symbols (public API)
- `internal/` prevents external imports
- Understanding `go.mod` dependencies

**Exercise**: Create multi-package project with internal dependencies

---

## Phase 2: Intermediate Go (Lessons 6-10)
*Building blocks for concurrent and web applications*

### Lesson 6: Pointers & Memory
**Folder**: `06-pointers/`

**Concepts**:
- Pointer basics (`*` and `&`)
- Pass by value vs pass by reference
- Pointer receivers vs value receivers
- nil pointers
- When to use pointers

**Why This Matters**:
- Backend uses pointer receivers extensively
- Understanding `*Handler`, `*Repository`, `*DB`
- Pointer parameters for optional fields
- Example: All handler and repository methods

**Exercise**: Understand pointer behavior and memory implications

---

### Lesson 7: Slices, Arrays & Maps
**Folder**: `07-collections/`

**Concepts**:
- Arrays (fixed size)
- Slices (dynamic arrays)
- Slice operations (append, slice syntax)
- Maps (key-value stores)
- Range loops
- Make and capacity

**Why This Matters**:
- Backend returns slices of users, games, moves
- Pagination returns `[]database.User`
- Maps used for configuration
- Example: `handler/v1/user/handler.go` List methods

**Exercise**: Work with collections, implement search and filter

---

### Lesson 8: Goroutines & Channels
**Folder**: `08-concurrency/`

**Concepts**:
- Goroutines with `go` keyword
- Channels for communication
- Channel operations (send, receive)
- Buffered vs unbuffered channels
- Select statement
- WaitGroups

**Why This Matters**:
- Backend starts server in goroutine
- Channels used for signal handling (SIGINT, SIGTERM)
- Graceful shutdown pattern
- Example: `app/app.go` Run method

**Exercise**: Create concurrent programs with goroutines and channels

---

### Lesson 9: Context Package
**Folder**: `09-context/`

**Concepts**:
- What is context
- `context.Background()` and `context.TODO()`
- `context.WithTimeout()` and `context.WithCancel()`
- Context propagation through functions
- Context values
- Cancellation signals

**Why This Matters**:
- Every database operation takes `ctx context.Context`
- Graceful shutdown uses context with timeout
- Request cancellation in HTTP handlers
- Example: All repository methods, `app/app.go` shutdown

**Exercise**: Implement timeout and cancellation with context

---

### Lesson 10: JSON & Struct Tags
**Folder**: `10-json-tags/`

**Concepts**:
- JSON marshaling and unmarshaling
- Struct tags (`json:"field_name"`)
- Validation tags (`binding:"required"`)
- omitempty and custom field names
- Nested structs
- Pointer fields for optional values

**Why This Matters**:
- Backend DTOs use extensive struct tags
- `CreateRequest`, `Response` structs
- Gin validation uses binding tags
- Example: `handler/v1/user/dto.go`

**Exercise**: Create API request/response models with proper tags

---

## Phase 3: Web Development (Lessons 11-15)
*HTTP services and API development*

### Lesson 11: HTTP Server Basics
**Folder**: `11-http-basics/`

**Concepts**:
- `net/http` package
- HTTP handlers
- Request and response
- Status codes
- ServeMux (router)
- Starting and stopping server

**Why This Matters**:
- Foundation before learning Gin framework
- Understanding underlying HTTP primitives
- Server lifecycle management
- Example: `server/http_server.go`

**Exercise**: Build a basic HTTP server with multiple routes

---

### Lesson 12: Gin Framework
**Folder**: `12-gin-framework/`

**Concepts**:
- Gin router and engine
- Route handlers with `gin.Context`
- Request binding (`ShouldBindJSON`, `ShouldBindQuery`)
- Response writing (`c.JSON()`)
- Path parameters (`c.Param()`)
- Query parameters
- Route groups

**Why This Matters**:
- Backend uses Gin for all HTTP handling
- Every handler uses `gin.Context`
- Understanding middleware and route groups
- Example: All handlers, `router/v1_routes.go`

**Exercise**: Build REST API with Gin (CRUD operations)

---

### Lesson 13: Middleware & Request Flow
**Folder**: `13-middleware/`

**Concepts**:
- What is middleware
- Gin middleware pattern
- Chain of responsibility
- Request/response interception
- CORS, logging, recovery middleware
- Custom middleware

**Why This Matters**:
- Backend uses CORS middleware
- Middleware for logging requests
- Understanding request pipeline
- Example: `server/http_server.go` setupMiddleware

**Exercise**: Create custom middleware (auth, logging, timing)

---

### Lesson 14: Input Validation & DTOs
**Folder**: `14-validation-dtos/`

**Concepts**:
- DTO (Data Transfer Object) pattern
- Request validation with tags
- Custom validators
- Error responses
- Request/Response transformation
- Separation of concerns (DB model vs API response)

**Why This Matters**:
- Backend separates request/response DTOs from database models
- Validation prevents invalid data
- Conversion functions between layers
- Example: `handler/v1/user/dto.go`, ToResponse functions

**Exercise**: Implement validated DTOs with conversion logic

---

### Lesson 15: Graceful Shutdown
**Folder**: `15-graceful-shutdown/`

**Concepts**:
- Signal handling (`os.Signal`)
- Graceful server shutdown
- Context timeout for shutdown
- Resource cleanup
- Goroutine coordination

**Why This Matters**:
- Backend implements proper graceful shutdown
- Prevents data loss or connection leaks
- Production-ready pattern
- Example: `app/app.go` Run method

**Exercise**: Implement server with graceful shutdown on SIGINT

---

## Phase 4: Database Integration (Lessons 16-20)
*Working with PostgreSQL and data persistence*

### Lesson 16: Database Basics & SQL
**Folder**: `16-database-sql/`

**Concepts**:
- `database/sql` package
- Driver registration
- Connection string (DSN)
- Query and Exec
- Scanning results
- Prepared statements
- SQL injection prevention

**Why This Matters**:
- Foundation for database operations
- Understanding pgx driver
- Safe query parameterization
- Example: `database/connection.go`

**Exercise**: Connect to PostgreSQL and run queries

---

### Lesson 17: Connection Pooling with pgx
**Folder**: `17-pgx-pool/`

**Concepts**:
- `pgxpool.Pool` for connection pooling
- Pool configuration
- Context-aware queries
- Ping for health checks
- Pool cleanup
- Transaction management

**Why This Matters**:
- Backend uses pgxpool for all database operations
- Connection reuse for performance
- Understanding pool lifecycle
- Example: `database/connection.go` New function

**Exercise**: Set up connection pool and execute queries

---

### Lesson 18: SQLC Code Generation
**Folder**: `18-sqlc/`

**Concepts**:
- What is sqlc
- Writing SQL queries for generation
- sqlc.yaml configuration
- Generated code structure
- Type-safe queries
- Query parameters

**Why This Matters**:
- Backend uses sqlc for all database access
- No ORM overhead, pure SQL
- Type safety without boilerplate
- Example: `internal/database/` generated files, `db/queries/`

**Exercise**: Write SQL queries and generate Go code with sqlc

---

### Lesson 19: Repository Pattern
**Folder**: `19-repository-pattern/`

**Concepts**:
- What is repository pattern
- Abstracting database layer
- Business logic in repositories
- Wrapping generated sqlc code
- Constructor injection
- Testing repositories

**Why This Matters**:
- Backend wraps sqlc queries in repositories
- Separation of data access and business logic
- Enables mocking for tests
- Example: `repository/user_repository.go`

**Exercise**: Implement repository pattern wrapping sqlc

---

### Lesson 20: Migrations & Schema Management
**Folder**: `20-migrations/`

**Concepts**:
- Database migrations
- golang-migrate tool
- Up and down migrations
- Schema versioning
- Migration best practices

**Why This Matters**:
- Backend uses migrations for schema changes
- Version control for database
- Deployment automation
- Example: `db/migrations/`

**Exercise**: Create and apply migrations

---

## Phase 5: Advanced Topics (Lessons 21-25)
*Testing, architecture, and production readiness*

### Lesson 21: Unit Testing
**Folder**: `21-unit-testing/`

**Concepts**:
- `testing` package
- Test functions and naming
- Table-driven tests
- Assertions and error checking
- Test setup and teardown
- `go test` command
- Test coverage

**Why This Matters**:
- Backend has unit tests for config
- Testing handlers and repositories
- TDD best practices
- Example: `config/config_test.go`

**Exercise**: Write unit tests for functions and methods

---

### Lesson 22: Mocking & Interface Testing
**Folder**: `22-mocking/`

**Concepts**:
- Why mock dependencies
- Interface-based mocking
- Manual mocks vs tools
- Testing with fake implementations
- Dependency injection for testability

**Why This Matters**:
- Backend interfaces enable mocking
- Test handlers without real database
- Test repositories without database
- Example: Testing patterns with Logger, Repository interfaces

**Exercise**: Create mocks and test components in isolation

---

### Lesson 23: Integration Testing
**Folder**: `23-integration-testing/`

**Concepts**:
- Testing with real database
- Test database setup
- Transaction rollback pattern
- Testing HTTP handlers end-to-end
- Seeding test data

**Why This Matters**:
- Verify full request flow
- Database integration tests
- API contract validation
- Example: Testing full handler → repository → database flow

**Exercise**: Write integration tests with test database

---

### Lesson 24: Dependency Injection & Architecture
**Folder**: `24-dependency-injection/`

**Concepts**:
- Constructor injection pattern
- Dependency graph
- Layered architecture (handler → repository → database)
- Hexagonal architecture principles
- Separation of concerns

**Why This Matters**:
- Backend architecture is layered
- Dependencies injected through constructors
- Testable, maintainable design
- Example: `cmd/server/main.go` wiring, `app/app.go`

**Exercise**: Build layered application with DI

---

### Lesson 25: Production Readiness
**Folder**: `25-production/`

**Concepts**:
- Configuration management
- Environment variables
- Logging best practices
- Health check endpoints
- Docker deployment
- Build optimization
- Security considerations

**Why This Matters**:
- Backend is production-ready
- 12-factor app principles
- Deployment patterns
- Example: `config/config.go`, `handler/v1/health/`, Dockerfile

**Exercise**: Deploy complete application with Docker

---

## Final Project: Mini Chess Coach
**Folder**: `final-project/`

**Objective**: Build a simplified version of chess-coach backend from scratch

**Requirements**:
1. User CRUD with Gin and PostgreSQL
2. Game management (create, list, get)
3. Repository pattern with sqlc
4. Unit and integration tests
5. Docker deployment
6. Graceful shutdown

**This consolidates all 25 lessons into one project!**

---

## How to Use This Plan

### For Each Lesson:

1. **Read** - Review the concept description
2. **Code** - Work through examples in the lesson folder
3. **Run** - Execute programs and see output
4. **Experiment** - Modify code and observe behavior
5. **Ask Questions** - Clarify anything unclear
6. **Exercise** - Complete the hands-on exercise
7. **Summarize** - Create `SUMMARY.md` in the lesson folder
8. **Relate** - Find examples in the backend codebase

### Lesson Folder Structure:
```
01-hello-world/
├── README.md          # Lesson content and concepts
├── examples/          # Runnable example code
│   ├── 01-simple.go
│   ├── 02-advanced.go
│   └── ...
├── exercises/         # Your practice code
│   └── solution.go
├── SUMMARY.md         # Your notes and learnings
└── references.md      # Links to backend examples
```

### Progress Tracking:
- [ ] Phase 1: Go Basics (Lessons 1-5)
- [ ] Phase 2: Intermediate Go (Lessons 6-10)
- [ ] Phase 3: Web Development (Lessons 11-15)
- [ ] Phase 4: Database Integration (Lessons 16-20)
- [ ] Phase 5: Advanced Topics (Lessons 21-25)
- [ ] Final Project: Mini Chess Coach

---

## Estimated Timeline

**Intensive Track** (2-3 weeks):
- 2 lessons per day
- 2-3 hours per lesson

**Standard Track** (4-6 weeks):
- 1 lesson per day
- 2 hours per lesson

**Relaxed Track** (8-10 weeks):
- 3-4 lessons per week
- 1-2 hours per session

---

## Learning Tips

1. **Type Code**: Don't copy-paste, type examples yourself
2. **Break Things**: Modify code to see what breaks
3. **Read Errors**: Go error messages are helpful
4. **Use Backend**: Reference chess-coach code frequently
5. **Ask Questions**: Document what you don't understand
6. **Build Projects**: Apply concepts immediately
7. **Review Summaries**: Your lesson summaries are your reference

---

## Next Steps

1. Create `01-hello-world/` folder
2. Start with Lesson 1
3. Work through progressively
4. Build toward backend understanding

**Let's begin! Ready to start Lesson 1?**

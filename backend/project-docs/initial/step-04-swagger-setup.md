# Swagger API Documentation Setup Guide

## Overview

This guide explains how to set up and use Swagger (OpenAPI) documentation for the Chess Coach backend API.

---

## What is Swagger?

Swagger (OpenAPI) is a standard way to document REST APIs. It provides:
- **Interactive API Documentation** - Test APIs directly from your browser
- **Auto-generated UI** - Beautiful documentation interface
- **API Discovery** - Easy for frontend developers to understand available endpoints
- **Code Generation** - Can generate client SDKs automatically

---

## Tools Used

- **swaggo/swag** - Generates Swagger docs from Go code annotations
- **gin-swagger** - Integrates Swagger UI with Gin framework
- **swaggo/files** - Serves Swagger UI static files

---

## Installation

### Step 1: Install Swag CLI Tool

```bash
# Install swag command-line tool
go install github.com/swaggo/swag/cmd/swag@latest

# Verify installation
swag --version
```

Expected output: `swag version v1.x.x`

### Step 2: Install Go Dependencies

```bash
# These should already be in your go.mod
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
go get -u github.com/gin-gonic/gin
```

---

## How It Works

### 1. Annotation Format

Swagger documentation is generated from special comments in your Go code:

```go
// @Summary Short description
// @Description Detailed description
// @Tags category
// @Accept json
// @Produce json
// @Param name location type required "description"
// @Success 200 {object} ResponseType
// @Failure 400 {object} ErrorType
// @Router /path [method]
func handlerFunction(c *gin.Context) {
    // Handler code
}
```

### 2. Main Annotations (in main.go)

```go
// @title Chess Coach API
// @version 1.0
// @description API for chess game analysis and coaching
// @host localhost:8080
// @BasePath /api/v1
```

These define the overall API metadata.

---

## Generating Swagger Documentation

### Command

Run this from your **backend** directory:

```bash
swag init -g cmd/server/main.go
```

### What This Does

1. Scans your code for `// @` annotations
2. Generates `docs/` directory with:
   - `docs.go` - Go code
   - `swagger.json` - OpenAPI JSON spec
   - `swagger.yaml` - OpenAPI YAML spec

### After Running

You should see:
```
backend/
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
```

---

## Running Your API with Swagger

### Step 1: Generate docs

```bash
cd backend
swag init -g cmd/server/main.go
```

### Step 2: Run your server

```bash
# Using Air (hot reload)
air

# Or using go run
go run cmd/server/main.go
```

### Step 3: Access Swagger UI

Open your browser and visit:

```
http://localhost:8080/swagger/index.html
```

You should see an interactive API documentation page!

---

## Common Swagger Annotations

### Endpoint Annotations

```go
// @Summary Get health status
// @Description Returns the health status of the API
// @Tags general
// @Produce json
// @Success 200 {object} map[string]string
// @Router /api/v1/health [get]
func getHealth(c *gin.Context) {
    c.JSON(200, gin.H{"status": "healthy"})
}
```

### POST Endpoint with Body

```go
// @Summary Analyze chess position
// @Description Analyzes a chess position and returns evaluation
// @Tags analysis
// @Accept json
// @Produce json
// @Param position body PositionRequest true "Chess position"
// @Success 200 {object} AnalysisResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/analyze [post]
func analyzePosition(c *gin.Context) {
    // Implementation
}
```

### GET Endpoint with Path Parameter

```go
// @Summary Get game by ID
// @Description Retrieves a chess game by its ID
// @Tags games
// @Produce json
// @Param id path string true "Game ID"
// @Success 200 {object} Game
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/games/{id} [get]
func getGame(c *gin.Context) {
    // Implementation
}
```

### GET Endpoint with Query Parameters

```go
// @Summary List games
// @Description Get a list of games with pagination
// @Tags games
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {array} Game
// @Router /api/v1/games [get]
func listGames(c *gin.Context) {
    // Implementation
}
```

---

## Model Documentation

Define your request/response models with JSON tags and examples:

```go
// PositionRequest represents a chess position to analyze
type PositionRequest struct {
    FEN string `json:"fen" binding:"required" example:"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"`
}

// AnalysisResponse represents the analysis result
type AnalysisResponse struct {
    Evaluation float64 `json:"evaluation" example:"0.5"`
    BestMove   string  `json:"best_move" example:"e2e4"`
    Depth      int     `json:"depth" example:"20"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
    Error   string `json:"error" example:"Invalid FEN notation"`
    Message string `json:"message" example:"The provided FEN string is malformed"`
}
```

---

## Annotation Reference

### Main Annotations

| Annotation | Description | Example |
|------------|-------------|---------|
| `@title` | API title | `@title Chess Coach API` |
| `@version` | API version | `@version 1.0` |
| `@description` | API description | `@description Chess analysis API` |
| `@host` | API host | `@host localhost:8080` |
| `@BasePath` | Base path | `@BasePath /api/v1` |
| `@schemes` | Supported schemes | `@schemes http https` |

### Endpoint Annotations

| Annotation | Description | Example |
|------------|-------------|---------|
| `@Summary` | Short description | `@Summary Get health status` |
| `@Description` | Detailed description | `@Description Returns API health` |
| `@Tags` | Group endpoints | `@Tags analysis` |
| `@Accept` | Request content type | `@Accept json` |
| `@Produce` | Response content type | `@Produce json` |
| `@Param` | Parameter definition | `@Param id path string true "ID"` |
| `@Success` | Success response | `@Success 200 {object} Response` |
| `@Failure` | Error response | `@Failure 400 {object} Error` |
| `@Router` | Route path and method | `@Router /health [get]` |

### Parameter Types

- `path` - URL path parameter (`/games/:id`)
- `query` - Query string parameter (`?page=1`)
- `header` - HTTP header
- `body` - Request body
- `formData` - Form data

---

## Development Workflow

### 1. Write Handler with Annotations

```go
// @Summary Create new game
// @Description Creates a new chess game
// @Tags games
// @Accept json
// @Produce json
// @Param game body GameCreateRequest true "Game details"
// @Success 201 {object} Game
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/games [post]
func createGame(c *gin.Context) {
    // Your code
}
```

### 2. Regenerate Swagger Docs

```bash
swag init -g cmd/server/main.go
```

### 3. Test in Swagger UI

1. Start your server: `air` or `go run cmd/server/main.go`
2. Open: `http://localhost:8080/swagger/index.html`
3. Click "Try it out" on your endpoint
4. Test the API directly from the UI

### 4. Commit Generated Docs

```bash
git add docs/
git commit -m "Update Swagger documentation"
```

---

## Troubleshooting

### Issue: `swag: command not found`

**Solution:** Install swag CLI tool
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Issue: `cannot import docs package`

**Solution:** Generate the docs first
```bash
swag init -g cmd/server/main.go
```

### Issue: Swagger UI shows old API

**Solution:** Regenerate docs and restart server
```bash
swag init -g cmd/server/main.go
air  # or go run cmd/server/main.go
```

### Issue: Model not showing in Swagger

**Solution:** Make sure the model is:
1. Exported (starts with capital letter)
2. Used in a `@Param` or `@Success` annotation
3. Has proper JSON tags

---

## Best Practices

### 1. Always Regenerate Before Committing

```bash
# Before git commit
swag init -g cmd/server/main.go
git add docs/
```

### 2. Use Meaningful Tags

Group related endpoints:
```go
// @Tags authentication
// @Tags games
// @Tags analysis
// @Tags users
```

### 3. Document All Responses

```go
// @Success 200 {object} Game
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 404 {object} ErrorResponse "Game not found"
// @Failure 500 {object} ErrorResponse "Internal error"
```

### 4. Use Examples in Models

```go
type Game struct {
    ID   string `json:"id" example:"game-123"`
    FEN  string `json:"fen" example:"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"`
}
```

### 5. Keep Docs in Sync with Code

After changing an endpoint, update the annotations and regenerate:
```bash
swag init -g cmd/server/main.go
```

---

## Air Integration

Add this to your `.air.toml` to auto-regenerate Swagger docs:

```toml
[build]
  # After build, regenerate swagger docs
  post_cmd = ["swag", "init", "-g", "cmd/server/main.go"]
```

Or create a simple script `scripts/dev.sh`:

```bash
#!/bin/bash
swag init -g cmd/server/main.go && air
```

---

## Production Considerations

### Disable Swagger in Production

```go
func main() {
    r := gin.Default()

    // Only enable Swagger in development
    if gin.Mode() != gin.ReleaseMode {
        r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    }

    r.Run(":8080")
}
```

### Set Release Mode

```bash
export GIN_MODE=release
go run cmd/server/main.go
```

---

## Example: Full Chess Analysis Endpoint

```go
// PositionRequest represents a chess position to analyze
type PositionRequest struct {
    FEN    string `json:"fen" binding:"required" example:"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"`
    Depth  int    `json:"depth" example:"20"`
    Time   int    `json:"time" example:"5000" description:"Analysis time in milliseconds"`
}

// AnalysisResponse represents the analysis result
type AnalysisResponse struct {
    Evaluation float64     `json:"evaluation" example:"0.5"`
    BestMove   string      `json:"best_move" example:"e2e4"`
    PV         []string    `json:"pv" example:"e2e4,e7e5,g1f3"`
    Depth      int         `json:"depth" example:"20"`
}

// @Summary Analyze chess position
// @Description Analyzes a chess position using Stockfish engine
// @Tags analysis
// @Accept json
// @Produce json
// @Param position body PositionRequest true "Chess position and analysis parameters"
// @Success 200 {object} AnalysisResponse "Analysis result"
// @Failure 400 {object} ErrorResponse "Invalid FEN or parameters"
// @Failure 500 {object} ErrorResponse "Analysis engine error"
// @Router /api/v1/analyze [post]
func analyzePosition(c *gin.Context) {
    var req PositionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, ErrorResponse{
            Error: "Invalid request",
            Message: err.Error(),
        })
        return
    }

    // Analysis logic here

    c.JSON(200, AnalysisResponse{
        Evaluation: 0.5,
        BestMove: "e2e4",
        PV: []string{"e2e4", "e7e5", "g1f3"},
        Depth: 20,
    })
}
```

---

## Quick Reference Commands

```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
swag init -g cmd/server/main.go

# Run server
air

# Access Swagger UI
open http://localhost:8080/swagger/index.html

# Format generated docs
swag fmt

# Check swag version
swag --version
```

---

## Resources

- [Swaggo GitHub](https://github.com/swaggo/swag)
- [Gin-Swagger](https://github.com/swaggo/gin-swagger)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger UI Demo](https://petstore.swagger.io/)

---

**Version:** 1.0
**Last Updated:** 2025-11-01
**For:** Chess Coach Backend

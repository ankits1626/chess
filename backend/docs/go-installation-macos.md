# Go Installation & Setup Guide for macOS

This document captures the step-by-step process of installing Go and setting up the development environment on macOS for the Chess Coach backend project.

---

## Prerequisites

- macOS with Homebrew installed
- VSCode installed
- Terminal (zsh shell)

---

## Installation Steps

### 1. Install Go via Homebrew

```bash
brew install go
```

**Installed version:** Go 1.25.3

**Verify installation:**
```bash
go version
# Output: go version go1.25.3 darwin/arm64
```

---

### 2. Install Go Development Tools

Install the essential Go tools for VSCode integration:

```bash
# Install all tools in one command
go install golang.org/x/tools/gopls@latest && \
go install github.com/go-delve/delve/cmd/dlv@latest && \
go install golang.org/x/tools/cmd/goimports@latest && \
go install honnef.co/go/tools/cmd/staticcheck@latest
```

**Installed tools:**
- `gopls@v0.20.0` - Go language server (IntelliSense, navigation)
- `delve` - Go debugger
- `goimports` - Auto-format and organize imports
- `staticcheck` - Advanced linter

**Installation location:** `~/go/bin/`

---

### 3. Update PATH Environment Variable

Add Go binaries to your shell PATH:

```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc
```

**Verify gopls is accessible:**
```bash
gopls version
# Output: golang.org/x/tools/gopls v0.20.0
```

---

### 4. Install VSCode Go Extension

```bash
code --install-extension golang.go
```

**Installed:** VSCode Go extension v0.50.0

---

### 5. Initialize Go Modules in Backend

```bash
cd backend
go mod tidy
```

This ensures `go.mod` is properly configured and dependencies are resolved.

---

## Verification Steps

### Test 1: Go Command Works

```bash
go version
# Should output: go version go1.25.3 darwin/arm64

go env GOPATH
# Should output: /Users/<username>/go
```

### Test 2: Go Tools Installed

```bash
which gopls
# Should output: /Users/<username>/go/bin/gopls

which dlv
# Should output: /Users/<username>/go/bin/dlv

which goimports
# Should output: /Users/<username>/go/bin/goimports
```

### Test 3: VSCode Integration

1. **Reload VSCode:**
   ```
   Cmd + Shift + P → "Developer: Reload Window"
   ```

2. **Open Go file:**
   - Navigate to `backend/cmd/server/main.go`

3. **Test IntelliSense:**
   - Hover over `http.HandleFunc` → Should show documentation popup
   - Type `http.` → Should see autocomplete suggestions
   - `Cmd + Click` on `HandleFunc` → Should jump to Go stdlib source

4. **Test Auto-Format:**
   - Add some extra spaces in main.go
   - Save file (`Cmd + S`)
   - Should auto-format with proper indentation

5. **Test Go Commands:**
   - `Cmd + Shift + P` → Type "Go: Install/Update Tools"
   - Should see all tools marked as installed

### Test 4: Build and Run

```bash
cd backend
go build -o bin/server ./cmd/server
./bin/server
```

Should output:
```json
{"time":"...","level":"INFO","msg":"starting server","port":"8080","addr":":8080"}
```

Test endpoint:
```bash
curl http://localhost:8080/health
# Output: {"status":"ok","service":"chess-coach-backend"}
```

---

## VSCode Configuration Applied

### Settings (.vscode/settings.json)
- Go language server enabled
- Format on save enabled (goimports)
- Lint on save enabled (golangci-lint)
- Organize imports automatically
- Tab size: 4 spaces (Go convention)

### Launch Configurations (.vscode/launch.json)
1. **Launch Backend Server** - Debug locally without Docker
2. **Attach to Docker Container** - Debug inside running container
3. **Test Current File** - Run tests with debugger

### Tasks (.vscode/tasks.json)
- `go: build` - Build all packages
- `go: test` - Run all tests
- `go: mod tidy` - Clean up dependencies
- `docker: up` - Start backend
- `docker: down` - Stop backend
- `docker: logs` - View logs

---

## Common Commands Reference

### Go Commands

```bash
# Build project
go build -o bin/server ./cmd/server

# Run directly
go run ./cmd/server

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Format code
go fmt ./...

# Organize imports
goimports -w .

# Lint code
staticcheck ./...

# Update dependencies
go mod tidy

# Download dependencies
go mod download

# Verify dependencies
go mod verify

# View module info
go list -m all
```

### VSCode Commands (Cmd + Shift + P)

- `Go: Install/Update Tools` - Update Go tools
- `Go: Test Function At Cursor` - Run test under cursor
- `Go: Test Package` - Run all tests in package
- `Go: Add Import` - Add missing import
- `Go: Generate Unit Tests For Function` - Auto-generate tests
- `Go: Fill Struct` - Auto-fill struct fields
- `Go: Add Tags To Struct Fields` - Add JSON/DB tags

---

## Troubleshooting

### Issue: "command not found: go"

**Solution:**
```bash
# Check if Go is installed
brew list go

# If not installed
brew install go

# Add to PATH
echo 'export PATH=$PATH:/opt/homebrew/bin' >> ~/.zshrc
source ~/.zshrc
```

---

### Issue: "command not found: gopls"

**Solution:**
```bash
# Install gopls
go install golang.org/x/tools/gopls@latest

# Add to PATH
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc

# Verify
gopls version
```

---

### Issue: VSCode doesn't recognize Go tools

**Solution:**
1. Check PATH in VSCode:
   - Open integrated terminal in VSCode (`Ctrl + ~`)
   - Run: `echo $PATH | grep go`
   - Should see `~/go/bin`

2. If not, reload VSCode:
   - `Cmd + Shift + P` → "Developer: Reload Window"

3. Manually install tools from VSCode:
   - `Cmd + Shift + P` → "Go: Install/Update Tools"
   - Select all tools and install

---

### Issue: "gopls was not able to find modules in your workspace"

**Solution:**
```bash
cd backend
go mod tidy
```

Then reload VSCode window.

---

### Issue: Import suggestions not working

**Solution:**
```bash
# Update gopls
go install golang.org/x/tools/gopls@latest

# Verify gopls is running
ps aux | grep gopls

# If not running, reload VSCode
```

---

### Issue: Format on save not working

**Solution:**
1. Check settings.json has:
```json
"[go]": {
  "editor.formatOnSave": true
}
```

2. Ensure goimports is installed:
```bash
go install golang.org/x/tools/cmd/goimports@latest
```

3. Restart gopls:
   - `Cmd + Shift + P` → "Go: Restart Language Server"

---

## Environment Variables

### Standard Go Environment

```bash
# View all Go environment variables
go env

# Key variables
GOPATH=/Users/<username>/go       # Go workspace
GOROOT=/opt/homebrew/Cellar/go/1.25.3/libexec  # Go installation
GOBIN=/Users/<username>/go/bin    # Installed binaries
GOCACHE=/Users/<username>/Library/Caches/go-build  # Build cache
GOMODCACHE=/Users/<username>/go/pkg/mod  # Module cache
```

### Project-Specific Variables

Set in `backend/.env` (not tracked in git):
```bash
PORT=8080
APP_ENV=development
```

---

## File Structure After Setup

```
~/go/                              # GOPATH
├── bin/                           # Installed Go binaries
│   ├── gopls
│   ├── dlv
│   ├── goimports
│   └── staticcheck
├── pkg/                           # Compiled packages (cache)
│   └── mod/                       # Downloaded modules
└── src/                           # Legacy GOPATH projects (unused)

chess-coach/
├── .vscode/
│   ├── settings.json              # Go editor settings
│   ├── launch.json                # Debug configurations
│   └── tasks.json                 # Build/test tasks
├── backend/
│   ├── cmd/server/main.go
│   ├── go.mod                     # Module definition
│   ├── go.sum                     # Dependency checksums (auto-generated)
│   └── tmp/                       # Air build artifacts (gitignored)
└── Makefile
```

---

## Additional Tools (Optional)

### golangci-lint (Advanced Linter)

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run
golangci-lint run
```

### Air (Hot Reload - Already in Docker)

For local development without Docker:
```bash
go install github.com/air-verse/air@v1.52.3

# Run
cd backend
air
```

### go-migrate (Database Migrations - Phase 2)

```bash
brew install golang-migrate

# Usage (future)
migrate -path db/migrations -database "postgres://..." up
```

---

## Performance Tips

### Speed Up Module Downloads

```bash
# Enable Go proxy (China users)
go env -w GOPROXY=https://goproxy.cn,direct

# Or use Google's proxy
go env -w GOPROXY=https://proxy.golang.org,direct

# Disable proxy
go env -w GOPROXY=direct
```

### Clean Build Cache

```bash
# Clear build cache
go clean -cache

# Clear module cache
go clean -modcache

# Clear test cache
go clean -testcache
```

---

## Next Steps

1. **Explore Go Tour:** https://go.dev/tour/
2. **Read Effective Go:** https://go.dev/doc/effective_go
3. **Review Go by Example:** https://gobyexample.com
4. **Setup Phase 2:** Add Fiber v3, WebSocket, Database

---

## Resources

- [Go Official Docs](https://go.dev/doc/)
- [VSCode Go Extension](https://github.com/golang/vscode-go)
- [gopls Documentation](https://github.com/golang/tools/blob/master/gopls/README.md)
- [Go Module Reference](https://go.dev/ref/mod)
- [Delve Debugger](https://github.com/go-delve/delve)

---

**Installation Date:** Nov 1, 2025
**Go Version:** 1.25.3
**macOS:** Sequoia (arm64)
**Status:** ✅ Verified and working

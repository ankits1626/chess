# VSCode Setup Guide for Go Development

This guide will help you configure VSCode for optimal Go development experience with the Chess Coach backend.

---

## Required Extensions

### 1. Go Extension (Official)
**Extension ID:** `golang.go`

**Install:**
```bash
code --install-extension golang.go
```

**Features:**
- IntelliSense (autocomplete)
- Code navigation (Go to Definition, Find References)
- Formatting (gofmt/goimports)
- Linting (staticcheck, golangci-lint)
- Debugging
- Test runner

---

## Recommended VSCode Settings

Create or update `.vscode/settings.json` in project root:

```json
{
  // Go Extension Settings
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "workspace",
  "go.formatTool": "goimports",
  "go.formatOnSave": true,
  "go.vetOnSave": "workspace",

  // Organize imports automatically
  "go.buildOnSave": "off",
  "go.coverOnSave": false,
  "go.testOnSave": false,

  // IntelliSense settings
  "go.autocompleteUnimportedPackages": true,
  "go.gocodeAutoBuild": false,

  // Gopls (Go Language Server) settings
  "gopls": {
    "ui.semanticTokens": true,
    "ui.completion.usePlaceholders": true,
    "ui.diagnostic.analyses": {
      "composites": false,
      "unusedparams": true,
      "unusedwrite": true,
      "useany": true
    },
    "ui.codelenses": {
      "gc_details": true,
      "generate": true,
      "test": true,
      "tidy": true,
      "upgrade_dependency": true,
      "vendor": true
    }
  },

  // Editor settings for Go files
  "[go]": {
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
      "source.organizeImports": "explicit"
    },
    "editor.insertSpaces": false,
    "editor.tabSize": 4,
    "editor.rulers": [100],
    "editor.bracketPairColorization.enabled": true
  },

  "[go.mod]": {
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
      "source.organizeImports": "explicit"
    }
  },

  // File associations
  "files.associations": {
    "*.tmpl": "html",
    "*.toml": "toml"
  },

  // File exclusions (keep explorer clean)
  "files.exclude": {
    "**/.git": true,
    "**/.DS_Store": true,
    "**/tmp": true
  },

  // Search exclusions
  "search.exclude": {
    "**/node_modules": true,
    "**/tmp": true,
    "**/.next": true,
    "**/dist": true,
    "**/vendor": true
  },

  // Terminal settings
  "terminal.integrated.cwd": "${workspaceFolder}",
  "terminal.integrated.defaultProfile.osx": "zsh"
}
```

---

## Workspace-Specific Settings

Create `.vscode/settings.json` for project-specific overrides:

```json
{
  "go.toolsManagement.autoUpdate": true,
  "go.gopath": "",
  "go.goroot": "",
  "go.inferGopath": false,

  // Point to backend directory for Go files
  "go.alternateTools": {},

  // Exclude frontend files from Go tooling
  "gopls": {
    "build.directoryFilters": [
      "-frontend",
      "-node_modules",
      "-tmp"
    ]
  }
}
```

---

## Recommended Launch Configuration

Create `.vscode/launch.json` for debugging:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Backend Server",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/backend/cmd/server",
      "env": {
        "PORT": "8080",
        "APP_ENV": "development"
      },
      "args": [],
      "showLog": true,
      "console": "integratedTerminal",
      "cwd": "${workspaceFolder}/backend"
    },
    {
      "name": "Attach to Docker Container",
      "type": "go",
      "request": "attach",
      "mode": "remote",
      "remotePath": "/app",
      "port": 2345,
      "host": "localhost",
      "showLog": true,
      "cwd": "${workspaceFolder}/backend"
    },
    {
      "name": "Test Current File",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${file}",
      "env": {},
      "args": ["-v"],
      "showLog": true
    }
  ]
}
```

---

## Recommended Tasks

Create `.vscode/tasks.json` for common operations:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "go: build",
      "type": "shell",
      "command": "go",
      "args": ["build", "-v", "./..."],
      "options": {
        "cwd": "${workspaceFolder}/backend"
      },
      "group": {
        "kind": "build",
        "isDefault": true
      },
      "problemMatcher": ["$go"]
    },
    {
      "label": "go: test",
      "type": "shell",
      "command": "go",
      "args": ["test", "-v", "./..."],
      "options": {
        "cwd": "${workspaceFolder}/backend"
      },
      "group": {
        "kind": "test",
        "isDefault": true
      },
      "problemMatcher": ["$go"]
    },
    {
      "label": "go: mod tidy",
      "type": "shell",
      "command": "go",
      "args": ["mod", "tidy"],
      "options": {
        "cwd": "${workspaceFolder}/backend"
      },
      "problemMatcher": []
    },
    {
      "label": "docker: up",
      "type": "shell",
      "command": "make up",
      "options": {
        "cwd": "${workspaceFolder}"
      },
      "problemMatcher": []
    },
    {
      "label": "docker: down",
      "type": "shell",
      "command": "make down",
      "options": {
        "cwd": "${workspaceFolder}"
      },
      "problemMatcher": []
    },
    {
      "label": "docker: logs",
      "type": "shell",
      "command": "make backend-logs",
      "options": {
        "cwd": "${workspaceFolder}"
      },
      "problemMatcher": [],
      "isBackground": true
    }
  ]
}
```

---

## Additional Recommended Extensions

### Development Productivity

1. **Error Lens** - `usernamehw.errorlens`
   - Shows errors inline in editor
   ```bash
   code --install-extension usernamehw.errorlens
   ```

2. **Better Comments** - `aaron-bond.better-comments`
   - Highlight TODO, FIXME, etc.
   ```bash
   code --install-extension aaron-bond.better-comments
   ```

3. **REST Client** - `humao.rest-client`
   - Test HTTP endpoints from VSCode
   ```bash
   code --install-extension humao.rest-client
   ```

4. **Docker** - `ms-azuretools.vscode-docker`
   - Manage containers, view logs
   ```bash
   code --install-extension ms-azuretools.vscode-docker
   ```

### Code Quality

5. **SonarLint** - `SonarSource.sonarlint-vscode`
   - Detect code smells and bugs
   ```bash
   code --install-extension SonarSource.sonarlint-vscode
   ```

6. **Code Spell Checker** - `streetsidesoftware.code-spell-checker`
   - Catch typos in code and comments
   ```bash
   code --install-extension streetsidesoftware.code-spell-checker
   ```

### Git Integration

7. **GitLens** - `eamodio.gitlens`
   - Enhanced Git features
   ```bash
   code --install-extension eamodio.gitlens
   ```

---

## Go Tools Installation

The Go extension will prompt you to install these tools. Accept the prompt or run manually:

```bash
# Install all recommended Go tools
go install golang.org/x/tools/gopls@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

**Note:** These will be installed to `$GOPATH/bin` (usually `~/go/bin`). Ensure it's in your PATH:

```bash
# Add to ~/.zshrc or ~/.bashrc
export PATH=$PATH:$(go env GOPATH)/bin
```

---

## Keyboard Shortcuts (macOS)

Useful shortcuts for Go development:

| Shortcut | Action |
|----------|--------|
| `Cmd + Shift + P` | Command Palette |
| `F12` | Go to Definition |
| `Cmd + Click` | Go to Definition (alternative) |
| `Shift + F12` | Find All References |
| `F2` | Rename Symbol |
| `Cmd + .` | Quick Fix |
| `Ctrl + Space` | Trigger IntelliSense |
| `Cmd + Shift + O` | Go to Symbol in File |
| `Cmd + T` | Go to Symbol in Workspace |
| `Ctrl + \`` | Toggle Terminal |
| `Cmd + B` | Toggle Sidebar |
| `Cmd + Shift + F` | Search in Files |

### Go-Specific Commands

Open Command Palette (`Cmd + Shift + P`) and type:

- `Go: Install/Update Tools` - Install/update Go tools
- `Go: Test Function At Cursor` - Run test under cursor
- `Go: Test Package` - Run all tests in package
- `Go: Generate Unit Tests For Function` - Auto-generate test
- `Go: Add Import` - Add missing import
- `Go: Add Tags To Struct Fields` - Add JSON/DB tags
- `Go: Fill Struct` - Auto-fill struct fields
- `Go: Extract to Function` - Refactor code to function

---

## Testing HTTP Endpoints with REST Client

Create `backend/api.http` for testing:

```http
### Health Check
GET http://localhost:8080/health
Content-Type: application/json

### Root Endpoint
GET http://localhost:8080/
Content-Type: application/json

### Future: Create Game (Phase 2)
# POST http://localhost:8080/api/games
# Content-Type: application/json
#
# {
#   "pgn": "1. e4 e5 2. Nf3 Nc6"
# }

### Future: WebSocket Connection Test (Phase 2)
# @host = localhost:8080
# GET ws://{{host}}/ws
# Upgrade: websocket
# Connection: Upgrade
```

Click "Send Request" above each `###` section to test endpoints.

---

## Snippets for Common Go Patterns

Create `.vscode/go.code-snippets`:

```json
{
  "HTTP Handler": {
    "prefix": "httphandler",
    "body": [
      "func ${1:handlerName}(w http.ResponseWriter, r *http.Request) {",
      "\tw.Header().Set(\"Content-Type\", \"application/json\")",
      "\tw.WriteHeader(http.StatusOK)",
      "\tfmt.Fprintf(w, `{\"message\":\"${2:response}\"})`)",
      "}"
    ],
    "description": "Create HTTP handler"
  },
  "Struct with JSON tags": {
    "prefix": "structjson",
    "body": [
      "type ${1:StructName} struct {",
      "\t${2:Field} ${3:string} `json:\"${4:field}\"`",
      "}"
    ],
    "description": "Create struct with JSON tags"
  },
  "Error Check": {
    "prefix": "iferr",
    "body": [
      "if err != nil {",
      "\treturn ${1:err}",
      "}"
    ],
    "description": "Basic error check"
  },
  "Log Error": {
    "prefix": "logerr",
    "body": [
      "if err != nil {",
      "\tslog.Error(\"${1:operation failed}\", \"error\", err)",
      "\treturn ${2:err}",
      "}"
    ],
    "description": "Log and return error"
  },
  "Test Function": {
    "prefix": "testfunc",
    "body": [
      "func Test${1:FunctionName}(t *testing.T) {",
      "\t// Arrange",
      "\t${2:input} := ${3:value}",
      "\t",
      "\t// Act",
      "\t${4:result} := ${5:FunctionToTest}(${2:input})",
      "\t",
      "\t// Assert",
      "\tif ${4:result} != ${6:expected} {",
      "\t\tt.Errorf(\"expected %v, got %v\", ${6:expected}, ${4:result})",
      "\t}",
      "}"
    ],
    "description": "Create test function"
  }
}
```

**Usage:** Type prefix (e.g., `httphandler`) and press Tab.

---

## File Watching and Auto-Formatting

### Format on Save
Already configured in `.vscode/settings.json`:
```json
"[go]": {
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": "explicit"
  }
}
```

### Manual Format
- **Format Document:** `Shift + Option + F`
- **Organize Imports:** `Cmd + Shift + P` → "Go: Add Import"

---

## Debugging in VSCode

### Debug Local Server (without Docker)

1. Set breakpoint in `backend/cmd/server/main.go` (click left of line number)
2. Press `F5` or go to Run & Debug sidebar
3. Select "Launch Backend Server"
4. Server starts with debugger attached
5. Hit `http://localhost:8080/health` in browser
6. Execution pauses at breakpoint

### Debug Inside Docker Container (Advanced)

1. Update `docker-compose.yml` to add Delve debugger:
```yaml
services:
  backend:
    # ... existing config
    ports:
      - "8080:8080"
      - "2345:2345"  # Delve debugger port
    command: |
      sh -c "
        go install github.com/go-delve/delve/cmd/dlv@latest &&
        dlv debug ./cmd/server --headless --listen=:2345 --api-version=2 --accept-multiclient
      "
```

2. Attach debugger with "Attach to Docker Container" launch config
3. Set breakpoints and debug

---

## Troubleshooting

### Issue: "gopls was not able to find modules in your workspace"

**Fix:**
```bash
cd backend
go mod tidy
```

Then reload VSCode: `Cmd + Shift + P` → "Developer: Reload Window"

### Issue: Import autocomplete not working

**Fix:**
1. `Cmd + Shift + P` → "Go: Install/Update Tools"
2. Select `gopls` and install
3. Reload window

### Issue: Formatting not working

**Fix:**
```bash
go install golang.org/x/tools/cmd/goimports@latest
```

Check PATH includes `~/go/bin`:
```bash
echo $PATH | grep go/bin
```

### Issue: "cannot find package" errors

**Fix:**
```bash
cd backend
go mod download
go mod tidy
```

### Issue: Extension slowing down VSCode

**Fix:** Exclude large directories from gopls:
```json
"gopls": {
  "build.directoryFilters": [
    "-node_modules",
    "-tmp",
    "-vendor",
    "-frontend"
  ]
}
```

---

## Multi-Root Workspace Setup (Optional)

If you want separate VSCode windows for frontend/backend:

Create `chess-coach.code-workspace`:

```json
{
  "folders": [
    {
      "name": "Backend (Go)",
      "path": "./backend"
    },
    {
      "name": "Frontend (React)",
      "path": "./frontend/app"
    }
  ],
  "settings": {
    "go.gopath": "",
    "go.goroot": ""
  }
}
```

Open: `File → Open Workspace from File` → Select `chess-coach.code-workspace`

---

## Quick Setup Script

Save as `backend/scripts/setup-vscode.sh`:

```bash
#!/bin/bash

# Install Go tools
echo "Installing Go tools..."
go install golang.org/x/tools/gopls@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install honnef.co/go/tools/cmd/staticcheck@latest
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install VSCode extensions
echo "Installing VSCode extensions..."
code --install-extension golang.go
code --install-extension usernamehw.errorlens
code --install-extension aaron-bond.better-comments
code --install-extension humao.rest-client
code --install-extension ms-azuretools.vscode-docker

echo "✓ Setup complete! Reload VSCode."
```

Run:
```bash
chmod +x backend/scripts/setup-vscode.sh
./backend/scripts/setup-vscode.sh
```

---

## Resources

- [Official Go VSCode Extension Docs](https://github.com/golang/vscode-go/wiki)
- [gopls Settings Reference](https://github.com/golang/tools/blob/master/gopls/doc/settings.md)
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- [Effective Go](https://go.dev/doc/effective_go)

---

**Last Updated:** Phase 1 Complete
**Next:** After Phase 2 setup, update debugging configs for Fiber + WebSocket

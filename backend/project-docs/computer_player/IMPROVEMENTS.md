# Improvements & Refactoring Notes

**Date**: 2025-11-05

---

## 🎯 Issue #1: GameMode Default Case

### Current Code (game_mode.go:36-40):
```go
func (gm GameMode) GetPlayerTypes() (PlayerType, PlayerType) {
    switch gm {
    // ... cases ...
    default:
        return PlayerTypeHuman, PlayerTypeHuman  // ❌ Wrong default
    }
}
```

### Problem:
- Default returns `HumanVsHuman` but our primary use case is `HumanVsComputer`
- If invalid mode is passed, we get wrong player types

### Solution:
```go
func (gm GameMode) GetPlayerTypes() (PlayerType, PlayerType) {
    switch gm {
    case GameModeHumanVsHuman:
        return PlayerTypeHuman, PlayerTypeHuman
    case GameModeHumanVsComputer:
        return PlayerTypeHuman, PlayerTypeComputer
    case GameModeComputerVsComputer:
        return PlayerTypeComputer, PlayerTypeComputer
    case GameModeHumanVsLLM:
        return PlayerTypeHuman, PlayerTypeLLM
    default:
        // Return HumanVsComputer as most common case
        return PlayerTypeHuman, PlayerTypeComputer  // ✅ Better default
    }
}
```

**Status**: 🟡 To be fixed in next iteration

---

## 🎯 Issue #2: PlayerConfig Violates ISP (Interface Segregation)

### Current Code (factory.go:20-35):
```go
// ❌ VIOLATES Interface Segregation Principle
type PlayerConfig struct {
    // For Human players
    Client *websocket.Client

    // For Computer players
    Difficulty string

    // For LLM players (future)
    APIKey string
    Model  string

    // For Network players (future)
    RemoteURL string
}
```

### Problem:
1. **Interface Segregation Violation**: Each player type must know about ALL config fields
2. **Confusing API**: Which fields are required for which player type?
3. **Error-prone**: Easy to pass wrong config (e.g., APIKey to HumanPlayer)
4. **Not extensible**: Adding new player type pollutes config for all

### SOLID Solution:

#### Option A: Separate Config Structs (Recommended) ⭐

```go
// Base config interface
type PlayerConfig interface {
    Validate() error
}

// Human player config
type HumanPlayerConfig struct {
    Client *websocket.Client
}

func (c HumanPlayerConfig) Validate() error {
    if c.Client == nil {
        return fmt.Errorf("client required for human player")
    }
    return nil
}

// Computer player config
type ComputerPlayerConfig struct {
    Difficulty string // "easy", "medium", "hard"
}

func (c ComputerPlayerConfig) Validate() error {
    if c.Difficulty == "" {
        c.Difficulty = "medium"
    }
    valid := c.Difficulty == "easy" || c.Difficulty == "medium" || c.Difficulty == "hard"
    if !valid {
        return fmt.Errorf("invalid difficulty: %s", c.Difficulty)
    }
    return nil
}

// LLM player config (future)
type LLMPlayerConfig struct {
    APIKey string
    Model  string
}

func (c LLMPlayerConfig) Validate() error {
    if c.APIKey == "" {
        return fmt.Errorf("API key required for LLM player")
    }
    if c.Model == "" {
        c.Model = "gpt-4" // default
    }
    return nil
}

// Network player config (future)
type NetworkPlayerConfig struct {
    RemoteURL string
}

func (c NetworkPlayerConfig) Validate() error {
    if c.RemoteURL == "" {
        return fmt.Errorf("remote URL required for network player")
    }
    return nil
}
```

#### Updated Factory:

```go
// CreatePlayer creates a player of the specified type
func (f *PlayerFactory) CreatePlayer(playerType PlayerType, config PlayerConfig) (Player, error) {
    if !playerType.IsValid() {
        return nil, fmt.Errorf("invalid player type: %s", playerType)
    }

    // Validate config
    if err := config.Validate(); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }

    switch playerType {
    case PlayerTypeHuman:
        humanConfig, ok := config.(HumanPlayerConfig)
        if !ok {
            return nil, fmt.Errorf("HumanPlayerConfig required for human player")
        }
        return NewHumanPlayer(humanConfig.Client), nil

    case PlayerTypeComputer:
        computerConfig, ok := config.(ComputerPlayerConfig)
        if !ok {
            return nil, fmt.Errorf("ComputerPlayerConfig required for computer player")
        }
        if f.aiService == nil {
            return nil, fmt.Errorf("AI service not configured")
        }
        return NewComputerPlayer(f.aiService, computerConfig.Difficulty), nil

    case PlayerTypeLLM:
        llmConfig, ok := config.(LLMPlayerConfig)
        if !ok {
            return nil, fmt.Errorf("LLMPlayerConfig required for LLM player")
        }
        return NewLLMPlayer(llmConfig.APIKey, llmConfig.Model), nil

    case PlayerTypeNetwork:
        networkConfig, ok := config.(NetworkPlayerConfig)
        if !ok {
            return nil, fmt.Errorf("NetworkPlayerConfig required for network player")
        }
        return NewNetworkPlayer(networkConfig.RemoteURL), nil

    default:
        return nil, fmt.Errorf("unsupported player type: %s", playerType)
    }
}
```

#### Usage Example:

```go
// ✅ CLEAN - Each config has only what it needs

// Create human player
humanPlayer, err := factory.CreatePlayer(
    PlayerTypeHuman,
    HumanPlayerConfig{Client: client},
)

// Create computer player
computerPlayer, err := factory.CreatePlayer(
    PlayerTypeComputer,
    ComputerPlayerConfig{Difficulty: "hard"},
)

// Create LLM player (future)
llmPlayer, err := factory.CreatePlayer(
    PlayerTypeLLM,
    LLMPlayerConfig{
        APIKey: "sk-...",
        Model:  "gpt-4",
    },
)
```

---

#### Option B: Functional Options Pattern (Alternative)

```go
// PlayerOption configures a player
type PlayerOption func(*playerBuilder) error

type playerBuilder struct {
    playerType PlayerType
    client     *websocket.Client
    difficulty string
    apiKey     string
    model      string
    remoteURL  string
}

// Options for different player types
func WithClient(client *websocket.Client) PlayerOption {
    return func(b *playerBuilder) error {
        b.client = client
        return nil
    }
}

func WithDifficulty(difficulty string) PlayerOption {
    return func(b *playerBuilder) error {
        b.difficulty = difficulty
        return nil
    }
}

func WithAPIKey(apiKey string) PlayerOption {
    return func(b *playerBuilder) error {
        b.apiKey = apiKey
        return nil
    }
}

// Usage
humanPlayer := factory.CreatePlayer(
    PlayerTypeHuman,
    WithClient(client),
)

computerPlayer := factory.CreatePlayer(
    PlayerTypeComputer,
    WithDifficulty("hard"),
)
```

---

## 📊 Comparison

| Approach | Pros | Cons | Recommendation |
|----------|------|------|----------------|
| **Current (Monolithic)** | Simple | Violates ISP, confusing | ❌ Refactor |
| **Separate Structs** | Clean, type-safe, SOLID | More types | ⭐ **Recommended** |
| **Functional Options** | Flexible, extensible | Less type-safe | Good alternative |

---

## 🎯 Recommended Changes

### Priority 1: Fix GameMode Default
- **File**: `game_mode.go`
- **Change**: Update default case to return `PlayerTypeHuman, PlayerTypeComputer`
- **Impact**: Low (1 line change)
- **When**: Before Step 3

### Priority 2: Refactor PlayerConfig
- **File**: `factory.go`
- **Change**: Implement separate config structs with interface
- **Impact**: Medium (affects factory usage)
- **When**: During Step 3 implementation

---

## ✅ Benefits of Refactoring

### Before (Current):
```go
// ❌ Confusing - what fields do I need?
config := PlayerConfig{
    Client:     client,     // Is this needed?
    Difficulty: "hard",     // Is this needed?
    APIKey:     "",         // Is this needed?
}
```

### After (Refactored):
```go
// ✅ Clear - only what I need
config := ComputerPlayerConfig{
    Difficulty: "hard",
}
```

**SOLID Principles Satisfied**:
- ✅ **S**: Each config has single responsibility
- ✅ **O**: Can add new player types without modifying existing configs
- ✅ **L**: All configs implement same PlayerConfig interface
- ✅ **I**: Each config has only what it needs (ISP)
- ✅ **D**: Factory depends on PlayerConfig interface

---

## 🚀 Implementation Plan

### When to Apply These Changes:

1. **GameMode Default**: Apply before Step 3 (5 minutes)
2. **PlayerConfig Refactor**: Apply during Step 3 (30 minutes)

### Step-by-Step:

1. Update `game_mode.go` default case
2. Create `player_config.go` with separate config structs
3. Update `factory.go` to use new configs
4. Update any usage in Step 3 implementation
5. Add validation to each config type
6. Test with different player types

---

**Excellent catch!** 🎯 This is exactly the kind of thinking that leads to maintainable, extensible code.

**Question**: Would you like to apply these fixes now, or continue to Step 2 and apply them during Step 3 implementation?

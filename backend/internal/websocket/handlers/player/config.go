package player

import (
	"fmt"

	"github.com/ankits1626/chess-coach-backend/internal/websocket"
)

// PlayerConfig is the base interface for all player configurations
type PlayerConfig interface {
	Validate() error
	GetPlayerType() PlayerType
}

// HumanPlayerConfig configures a human player
type HumanPlayerConfig struct {
	Client *websocket.Client
}

func (c HumanPlayerConfig) Validate() error {
	if c.Client == nil {
		return fmt.Errorf("client required for human player")
	}
	return nil
}

func (c HumanPlayerConfig) GetPlayerType() PlayerType {
	return PlayerTypeHuman
}

// ComputerPlayerConfig configures a computer player
type ComputerPlayerConfig struct {
	Difficulty string // "easy", "medium", "hard"
}

func (c ComputerPlayerConfig) Validate() error {
	if c.Difficulty == "" {
		c.Difficulty = "medium" // default
	}

	valid := c.Difficulty == "easy" ||
		c.Difficulty == "medium" ||
		c.Difficulty == "hard"

	if !valid {
		return fmt.Errorf("invalid difficulty: %s (must be easy, medium, or hard)", c.Difficulty)
	}

	return nil
}

func (c ComputerPlayerConfig) GetPlayerType() PlayerType {
	return PlayerTypeComputer
}

// Future configs can be added here...

// LLMPlayerConfig configures an LLM-powered player (future)
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

func (c LLMPlayerConfig) GetPlayerType() PlayerType {
	return PlayerTypeLLM
}

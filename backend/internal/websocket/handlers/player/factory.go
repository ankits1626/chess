package player

import (
	"fmt"
)

// PlayerFactory creates player instances based on type
type PlayerFactory struct {
	aiService AIService
}

// NewPlayerFactory creates a new player factory
func NewPlayerFactory(aiService AIService) *PlayerFactory {
	return &PlayerFactory{
		aiService: aiService,
	}
}

// CreatePlayer creates a player of the specified type with the given config
func (f *PlayerFactory) CreatePlayer(config PlayerConfig) (Player, error) {
	// Validate config
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// Create player based on type
	switch config.GetPlayerType() {
	case PlayerTypeHuman:
		humanConfig, ok := config.(HumanPlayerConfig)
		if !ok {
			return nil, fmt.Errorf("expected HumanPlayerConfig, got %T", config)
		}
		return NewHumanPlayer(humanConfig.Client), nil

	case PlayerTypeComputer:
		computerConfig, ok := config.(ComputerPlayerConfig)
		if !ok {
			return nil, fmt.Errorf("expected ComputerPlayerConfig, got %T", config)
		}
		if f.aiService == nil {
			return nil, fmt.Errorf("AI service not configured")
		}
		return NewComputerPlayer(f.aiService, computerConfig.Difficulty), nil

	case PlayerTypeLLM:
		return nil, fmt.Errorf("LLM player not yet implemented")

	case PlayerTypeNetwork:
		return nil, fmt.Errorf("network player not yet implemented")

	default:
		return nil, fmt.Errorf("unsupported player type: %s", config.GetPlayerType())
	}
}

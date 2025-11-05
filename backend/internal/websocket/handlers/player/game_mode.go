package player

import "fmt"

// GameMode defines who plays against whom
type GameMode string

const (
	GameModeHumanVsHuman       GameMode = "human_vs_human"
	GameModeHumanVsComputer    GameMode = "human_vs_computer"
	GameModeComputerVsComputer GameMode = "computer_vs_computer"
	GameModeHumanVsLLM         GameMode = "human_vs_llm"
)

// String returns string representation
func (gm GameMode) String() string {
	return string(gm)
}

// IsValid checks if game mode is valid
func (gm GameMode) IsValid() bool {
	switch gm {
	case GameModeHumanVsHuman, GameModeHumanVsComputer,
		GameModeComputerVsComputer, GameModeHumanVsLLM:
		return true
	}
	return false
}

// Validate returns an error if game mode is invalid
func (gm GameMode) Validate() error {
	if !gm.IsValid() {
		return fmt.Errorf("invalid game mode: %s", gm)
	}
	return nil
}

// GetPlayerTypes returns the two player types for this mode
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
		return PlayerTypeHuman, PlayerTypeComputer
	}
}

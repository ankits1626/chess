package player

// PlayerType represents the type of player
type PlayerType string

const (
	PlayerTypeHuman    PlayerType = "human"
	PlayerTypeComputer PlayerType = "computer"
	PlayerTypeLLM      PlayerType = "llm"
	PlayerTypeNetwork  PlayerType = "network" // Future: remote player via API
)

// String returns string representation
func (pt PlayerType) String() string {
	return string(pt)
}

// IsValid checks if player type is valid
func (pt PlayerType) IsValid() bool {
	switch pt {
	case PlayerTypeHuman, PlayerTypeComputer, PlayerTypeLLM, PlayerTypeNetwork:
		return true
	}
	return false
}

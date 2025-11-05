package player

import (
	"context"
	"errors"
)

// Player represents any chess player in the system.
// This interface allows for different player implementations:
// - HumanPlayer (WebSocket client)
// - ComputerPlayer (AI engine)
// - LLMPlayer (GPT-4, Claude, etc.)
type Player interface {
	// Core identification
	GetID() string
	GetType() PlayerType

	// Move handling
	// RequestMove asks the player for their next move given current position
	// - For Human: waits for WebSocket message
	// - For Computer: runs AI calculation
	// - For LLM: calls API
	RequestMove(ctx context.Context, fen string) (string, error)

	// SendMove notifies the player about a move (their own or opponent's)
	SendMove(move MoveNotification) error

	// Game events
	NotifyGameStart(info GameStartInfo) error
	NotifyGameEnd(result GameResult) error

	// Lifecycle
	IsConnected() bool
	Cleanup() error
}

// MoveNotification contains information about a move
type MoveNotification struct {
	GameID     string
	MoveSAN    string
	MoveUCI    string
	FEN        string
	MoveNumber int
	IsGameOver bool
	Result     *GameResult // nil if game not over
}

// GameStartInfo contains game start information
type GameStartInfo struct {
	GameID      string
	Mode        GameMode
	Opponent    string // Opponent ID or "Computer" or "GPT-4"
	YourSide    string // "white" or "black"
	TimeControl string
	StartingFEN string
}

// GameResult contains game result information
type GameResult struct {
	Winner string // "white", "black", "draw"
	Method string // "checkmate", "stalemate", "resignation", etc.
	PGN    string
}

// ErrPlayerDisconnected is returned when a player is not connected
var ErrPlayerDisconnected = errors.New("player disconnected")

// ErrMoveTimeout is returned when RequestMove times out
var ErrMoveTimeout = errors.New("move request timeout")

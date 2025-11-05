package player

import (
	"context"
	"fmt"
	"log"

	"github.com/notnil/chess"
)

// ComputerPlayer represents an AI player powered by a chess engine
type ComputerPlayer struct {
	id          string
	aiService   AIService
	difficulty  string
	currentGame *chess.Game
}

// NewComputerPlayer creates a new computer player
func NewComputerPlayer(aiService AIService, difficulty string) *ComputerPlayer {
	return &ComputerPlayer{
		id:         fmt.Sprintf("computer-%s", difficulty),
		aiService:  aiService,
		difficulty: difficulty,
	}
}

// GetID returns the player's unique identifier
func (c *ComputerPlayer) GetID() string {
	return c.id
}

// GetType returns the player type
func (c *ComputerPlayer) GetType() PlayerType {
	return PlayerTypeComputer
}

// RequestMove calculates the best move using the AI engine
func (c *ComputerPlayer) RequestMove(ctx context.Context, fen string) (string, error) {
	// Parse FEN to create game state
	fenFunc, err := chess.FEN(fen)
	if err != nil {
		return "", fmt.Errorf("invalid FEN: %w", err)
	}

	game := chess.NewGame(fenFunc)
	c.currentGame = game

	// Check if game is over
	if game.Outcome() != chess.NoOutcome {
		return "", fmt.Errorf("game is already over")
	}

	// Get best move from AI
	move, err := c.aiService.GetBestMove(ctx, game, c.difficulty)
	if err != nil {
		return "", fmt.Errorf("AI move calculation failed: %w", err)
	}

	// Return move in UCI format (e.g., "e2e4")
	return move.String(), nil
}

// SendMove receives move notifications (computer doesn't need to see moves)
func (c *ComputerPlayer) SendMove(move MoveNotification) error {
	// Computer player doesn't need move notifications
	// But we log for debugging
	log.Printf("[ComputerPlayer %s] Received move: %s", c.id, move.MoveSAN)
	return nil
}

// NotifyGameStart receives game start notification
func (c *ComputerPlayer) NotifyGameStart(info GameStartInfo) error {
	log.Printf("[ComputerPlayer %s] Game started: %s vs %s", c.id, info.YourSide, info.Opponent)
	return nil
}

// NotifyGameEnd receives game end notification
func (c *ComputerPlayer) NotifyGameEnd(result GameResult) error {
	log.Printf("[ComputerPlayer %s] Game ended: %s by %s", c.id, result.Winner, result.Method)
	return nil
}

// IsConnected always returns true for computer players
func (c *ComputerPlayer) IsConnected() bool {
	return true
}

// Cleanup performs any necessary cleanup
func (c *ComputerPlayer) Cleanup() error {
	c.currentGame = nil
	return nil
}

package player

import (
	"context"
	"time"

	"github.com/ankits1626/chess-coach-backend/internal/websocket"
)

// HumanPlayer represents a human player connected via WebSocket
type HumanPlayer struct {
	client *websocket.Client
}

// NewHumanPlayer creates a new human player
func NewHumanPlayer(client *websocket.Client) *HumanPlayer {
	return &HumanPlayer{
		client: client,
	}
}

// GetID returns the player's unique identifier
func (h *HumanPlayer) GetID() string {
	return h.client.ID
}

// GetType returns the player type
func (h *HumanPlayer) GetType() PlayerType {
	return PlayerTypeHuman
}

// GetClient returns the underlying WebSocket client (for HumanPlayer-specific operations)
func (h *HumanPlayer) GetClient() *websocket.Client {
	return h.client
}

// RequestMove waits for the human player to make a move via WebSocket
func (h *HumanPlayer) RequestMove(ctx context.Context, fen string) (string, error) {
	// Create a channel to receive the move
	moveCh := make(chan string, 1)
	errCh := make(chan error, 1)

	// Set up a temporary handler for move responses
	// This would integrate with your existing WebSocket message handling
	// For now, this is a placeholder that shows the pattern

	// In practice, you'd:
	// 1. Register a pending move request for this client
	// 2. When the client sends a move message, fulfill this request
	// 3. Use the context for timeout handling

	select {
	case move := <-moveCh:
		return move, nil
	case err := <-errCh:
		return "", err
	case <-ctx.Done():
		return "", ErrMoveTimeout
	case <-time.After(5 * time.Minute): // Generous timeout for human
		return "", ErrMoveTimeout
	}
}

// SendMove sends a move notification to the human player
func (h *HumanPlayer) SendMove(move MoveNotification) error {
	// Send move to client via WebSocket
	msg := &websocket.Message{
		Type:  websocket.TypeEvent,
		Event: "move",
		Data: map[string]interface{}{
			"gameId":     move.GameID,
			"moveSAN":    move.MoveSAN,
			"moveUCI":    move.MoveUCI,
			"fen":        move.FEN,
			"moveNumber": move.MoveNumber,
			"isGameOver": move.IsGameOver,
			"result":     move.Result,
		},
	}
	return h.client.SendMessage(msg)
}

// NotifyGameStart sends game start notification
func (h *HumanPlayer) NotifyGameStart(info GameStartInfo) error {
	msg := &websocket.Message{
		Type:  websocket.TypeEvent,
		Event: "game_start",
		Data: map[string]interface{}{
			"gameId":      info.GameID,
			"mode":        info.Mode,
			"opponent":    info.Opponent,
			"yourSide":    info.YourSide,
			"timeControl": info.TimeControl,
			"startingFEN": info.StartingFEN,
		},
	}
	return h.client.SendMessage(msg)
}

// NotifyGameEnd sends game end notification
func (h *HumanPlayer) NotifyGameEnd(result GameResult) error {
	msg := &websocket.Message{
		Type:  websocket.TypeEvent,
		Event: "game_end",
		Data: map[string]interface{}{
			"winner": result.Winner,
			"method": result.Method,
			"pgn":    result.PGN,
		},
	}
	return h.client.SendMessage(msg)
}

// IsConnected checks if the player is still connected
func (h *HumanPlayer) IsConnected() bool {
	// Since Client doesn't expose IsConnected(), we just check if client exists
	// The actual connection state is managed by the WebSocket infrastructure
	return h.client != nil
}

// Cleanup performs any necessary cleanup
func (h *HumanPlayer) Cleanup() error {
	// Human player cleanup is handled by WebSocket disconnect
	return nil
}

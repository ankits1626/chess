package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/websocket"
)

// ResignHandler handles game resignations.
type ResignHandler struct {
	manager *GameManager
}

// Handle implements ActionHandler.Handle for resignation.
func (h *ResignHandler) Handle(ctx context.Context, client *websocket.Client, msg *websocket.Message) error {
	gameIDStr, err := websocket.RequireString(msg.Data, "gameId")
	if err != nil {
		return err
	}

	log.Printf("Resign: User %s resigning from game %s", client.UserID, gameIDStr)

	activeGame, exists := h.manager.GetActiveGame(gameIDStr)
	if !exists {
		return fmt.Errorf("game not found or not active")
	}

	return h.handleForfeit(ctx, client, activeGame, msg.ID)
}

// handleForfeit processes game forfeit.
func (h *ResignHandler) handleForfeit(ctx context.Context, client *websocket.Client, activeGame *ActiveGame, requestID string) error {
	gameIDStr := activeGame.GameID
	gameID := websocket.StringToUUID(gameIDStr)

	// Determine result
	playerSide, err := h.manager.GetPlayerSide(gameIDStr, client)
	if err != nil {
		return err
	}

	var result string
	var winner *websocket.Client
	var loser *websocket.Client

	if playerSide == "white" {
		result = "0-1" // Black wins
		winner = activeGame.BlackPlayer
		loser = activeGame.WhitePlayer
	} else {
		result = "1-0" // White wins
		winner = activeGame.WhitePlayer
		loser = activeGame.BlackPlayer
	}

	log.Printf("Forfeit: Game %s - %s resigned, result: %s", gameIDStr, client.UserID, result)

	// Update game in database
	_, err = h.manager.GetDB().UpdateGame(ctx, database.UpdateGameParams{
		ID:     gameID,
		Pgn:    "", // TODO: Get PGN from game state
		Result: websocket.StringToText(result),
	})
	if err != nil {
		log.Printf("Forfeit: Failed to update game: %v", err)
	}

	// Send response to resigner
	websocket.SendSuccessResponse(loser, requestID, map[string]interface{}{
		"message":  "You resigned",
		"result":   result,
		"gameOver": true,
	})

	// Notify winner
	opponentResignedEvent := websocket.NewEvent("opponentResigned", map[string]interface{}{
		"gameId":   gameIDStr,
		"result":   result,
		"gameOver": true,
		"message":  "Opponent resigned",
	})
	winner.SendMessage(opponentResignedEvent)

	// Clean up
	h.manager.RemoveActiveGame(gameIDStr)
	loser.SetGameID("")
	winner.SetGameID("")

	return nil
}

// LeaveGameHandler handles leaving games.
type LeaveGameHandler struct {
	manager *GameManager
}

// Handle implements ActionHandler.Handle for leaving.
func (h *LeaveGameHandler) Handle(ctx context.Context, client *websocket.Client, msg *websocket.Message) error {
	gameIDStr, err := websocket.RequireString(msg.Data, "gameId")
	if err != nil {
		return err
	}

	log.Printf("LeaveGame: User %s leaving game %s", client.UserID, gameIDStr)

	// Check pending game
	if pending, exists := h.manager.GetPendingGame(gameIDStr); exists {
		if pending.WhitePlayer.ID == client.ID {
			h.manager.RemovePendingGame(gameIDStr)
			client.SetGameID("")

			websocket.SendSuccessResponse(client, msg.ID, map[string]interface{}{
				"message": "Game cancelled",
			})
			return nil
		}
	}

	// Check active game - treat as resignation
	if activeGame, exists := h.manager.GetActiveGame(gameIDStr); exists {
		resignHandler := &ResignHandler{manager: h.manager}
		return resignHandler.handleForfeit(ctx, client, activeGame, msg.ID)
	}

	return fmt.Errorf("game not found")
}

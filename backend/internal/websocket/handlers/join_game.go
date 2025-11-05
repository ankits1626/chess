package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/websocket"
	"github.com/ankits1626/chess-coach-backend/internal/websocket/handlers/player"
)

// JoinGameHandler handles joining existing games.
type JoinGameHandler struct {
	manager *GameManager
}

// Handle implements ActionHandler.Handle.
//
// Request format:
//
//	{
//	  "id": "req-456",
//	  "type": "request",
//	  "action": "joinGame",
//	  "data": {
//	    "gameId": "uuid"
//	  }
//	}
func (h *JoinGameHandler) Handle(ctx context.Context, client *websocket.Client, msg *websocket.Message) error {
	// 1. Extract and validate game ID
	gameIDStr, err := websocket.RequireString(msg.Data, "gameId")
	if err != nil {
		return err
	}

	if !websocket.IsValidUUID(gameIDStr) {
		return fmt.Errorf("invalid game ID format")
	}

	log.Printf("JoinGame: User %s attempting to join game %s", client.UserID, gameIDStr)

	// 2. Get pending game
	pending, exists := h.manager.GetPendingGame(gameIDStr)
	if !exists {
		return fmt.Errorf("game not found or already started")
	}

	// 3. Only allow joining human vs human games
	if pending.Mode != player.GameModeHumanVsHuman {
		return fmt.Errorf("cannot join this game mode")
	}

	// 4. Get white player's client (must be human player)
	whiteHuman, ok := pending.WhitePlayer.(*player.HumanPlayer)
	if !ok {
		return fmt.Errorf("white player is not a human player")
	}

	// 5. Verify not joining own game
	whiteClient := whiteHuman.GetClient()
	if whiteClient.UserID == client.UserID {
		return fmt.Errorf("cannot join your own game")
	}

	// 4. Update game in database
	gameID := websocket.StringToUUID(gameIDStr)
	blackPlayerUUID := websocket.StringToUUID(client.UserID)
	if !blackPlayerUUID.Valid {
		return fmt.Errorf("invalid user ID")
	}

	_, err = h.manager.GetDB().UpdateGame(ctx, database.UpdateGameParams{
		ID:     gameID,
		Pgn:    "",
		Result: websocket.StringToText("*"),
	})
	if err != nil {
		log.Printf("JoinGame: Database error: %v", err)
		return fmt.Errorf("failed to join game")
	}

	// 5. Activate game
	activeGame, err := h.manager.ActivateGame(gameIDStr, client)
	if err != nil {
		return fmt.Errorf("failed to activate game: %w", err)
	}

	// 6. Set game IDs
	client.SetGameID(gameIDStr)
	whiteClient.SetGameID(gameIDStr)

	// 7. Create room and add both players
	room := websocket.NewRoom(gameIDStr)
	whiteClient.JoinRoom(room)
	client.JoinRoom(room)

	log.Printf("JoinGame: Game %s started - White: %s, Black: %s",
		gameIDStr, whiteClient.UserID, client.UserID)

	// 8. Send success response to joiner
	websocket.SendSuccessResponse(client, msg.ID, map[string]interface{}{
		"gameId":      gameIDStr,
		"status":      "active",
		"side":        "black",
		"opponent":    whiteClient.UserID,
		"timeControl": pending.TimeControl,
		"fen":         activeGame.CurrentFEN,
	})

	// 9. Notify white player
	opponentJoinedEvent := websocket.NewEvent("opponentJoined", map[string]interface{}{
		"gameId":   gameIDStr,
		"opponent": client.UserID,
		"side":     "black",
		"fen":      activeGame.CurrentFEN,
	})
	whiteClient.SendMessage(opponentJoinedEvent)

	return nil
}

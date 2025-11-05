package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/ankits1626/chess-coach-backend/internal/websocket"
	"github.com/jackc/pgx/v5/pgtype"
)

// MakeMoveHandler handles move processing.
type MakeMoveHandler struct {
	manager *GameManager
}

// Handle implements ActionHandler.Handle.
//
// Request format:
//
//	{
//	  "id": "req-789",
//	  "type": "request",
//	  "action": "makeMove",
//	  "data": {
//	    "gameId": "uuid",
//	    "move": "e4"
//	  }
//	}
func (h *MakeMoveHandler) Handle(ctx context.Context, client *websocket.Client, msg *websocket.Message) error {
	// 1. Extract and validate inputs
	gameIDStr, err := websocket.RequireString(msg.Data, "gameId")
	if err != nil {
		return err
	}

	moveSAN, err := websocket.RequireString(msg.Data, "move")
	if err != nil {
		return err
	}

	if !websocket.ValidateMoveSAN(moveSAN) {
		return fmt.Errorf("invalid move format")
	}

	log.Printf("MakeMove: User %s attempting move '%s' in game %s",
		client.UserID, moveSAN, gameIDStr)

	// 2. Get active game
	activeGame, exists := h.manager.GetActiveGame(gameIDStr)
	if !exists {
		return fmt.Errorf("game not found or not active")
	}

	// 3. Verify it's player's turn
	playerSide, err := h.manager.GetPlayerSide(gameIDStr, client)
	if err != nil {
		return err
	}

	currentTurn := h.getCurrentTurn(activeGame.CurrentFEN)
	if playerSide != currentTurn {
		return fmt.Errorf("not your turn")
	}

	// 4. Apply move using chess service
	gameState, err := h.manager.GetChessService().ApplyMove(activeGame.CurrentFEN, moveSAN)
	if err != nil {
		return fmt.Errorf("invalid move: %w", err)
	}

	// 5. Update game state
	if err := h.manager.UpdateGameState(gameIDStr, gameState.FEN); err != nil {
		return fmt.Errorf("failed to update game state: %w", err)
	}

	log.Printf("MakeMove: Move applied - #%d %s (%s) by %s, gameOver=%v",
		gameState.MoveNumber, gameState.SAN, gameState.UCI, playerSide, gameState.IsGameOver)

	// 6. Save move to database
	gameID := websocket.StringToUUID(gameIDStr)
	_, err = h.manager.GetDB().CreateMove(ctx, database.CreateMoveParams{
		GameID:     gameID,
		MoveNumber: int32(gameState.MoveNumber),
		Side:       playerSide,
		MoveSan:    gameState.SAN,
		MoveUci:    gameState.UCI,
		Fen:        gameState.FEN,
		TimeTaken:  pgtype.Int4{Valid: false}, // TODO: Phase 5
	})
	if err != nil {
		log.Printf("MakeMove: Failed to save move: %v", err)
	}

	// 7. Update game if finished
	if gameState.IsGameOver {
		result := h.getResult(gameState.Result)
		_, err := h.manager.GetDB().UpdateGame(ctx, database.UpdateGameParams{
			ID:     gameID,
			Pgn:    gameState.Result.PGN,
			Result: websocket.StringToText(result),
		})
		if err != nil {
			log.Printf("MakeMove: Failed to update game result: %v", err)
		}
	}

	// 8. Prepare response data
	responseData := map[string]interface{}{
		"moveNumber": gameState.MoveNumber,
		"san":        gameState.SAN,
		"uci":        gameState.UCI,
		"fen":        gameState.FEN,
		"gameOver":   gameState.IsGameOver,
	}

	if gameState.IsGameOver {
		responseData["result"] = h.getResult(gameState.Result)
		responseData["method"] = gameState.Result.Method
	}

	// 9. Send success response
	websocket.SendSuccessResponse(client, msg.ID, responseData)

	// 10. Broadcast to opponent
	opponent, err := h.manager.GetOpponent(gameIDStr, client)
	if err == nil {
		moveMadeEvent := websocket.NewEvent("moveMade", responseData)
		opponent.SendMessage(moveMadeEvent)
	}

	// 11. Clean up if game over
	if gameState.IsGameOver {
		h.manager.RemoveActiveGame(gameIDStr)
	}

	return nil
}

// getCurrentTurn determines whose turn it is from FEN.
func (h *MakeMoveHandler) getCurrentTurn(fen string) string {
	// FEN format: "pieces active_color ..."
	// Second field is 'w' or 'b'
	parts := strings.Fields(fen)
	if len(parts) >= 2 {
		if parts[1] == "w" {
			return "white"
		}
		return "black"
	}
	return "white" // Default
}

// getResult converts GameResult to result string.
func (h *MakeMoveHandler) getResult(result GameResult) string {
	switch result.Winner {
	case "white":
		return "1-0"
	case "black":
		return "0-1"
	case "draw":
		return "1/2-1/2"
	default:
		return "*"
	}
}

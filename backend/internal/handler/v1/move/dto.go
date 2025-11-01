// Package move handles move-related HTTP requests for API v1.
package move

import (
	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/google/uuid"
)

// CreateRequest represents move creation request.
type CreateRequest struct {
	MoveNumber int32  `json:"move_number" binding:"required,min=1"`
	Side       string `json:"side" binding:"required,oneof=white black"`
	MoveSan    string `json:"move_san" binding:"required"`
	MoveUci    string `json:"move_uci" binding:"required"`
	Fen        string `json:"fen" binding:"required"`
	TimeTaken  *int32 `json:"time_taken,omitempty"` // Optional: milliseconds
}

// Response represents move in API responses.
type Response struct {
	ID         uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	GameID     uuid.UUID `json:"game_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	MoveNumber int32     `json:"move_number" example:"1"`
	Side       string    `json:"side" example:"white"`
	MoveSan    string    `json:"move_san" example:"e4"`
	MoveUci    string    `json:"move_uci" example:"e2e4"`
	Fen        string    `json:"fen" example:"rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 1"`
	TimeTaken  *int32    `json:"time_taken,omitempty" example:"1500"`
}

// ToResponse converts database.Move to Response.
func ToResponse(m database.Move) Response {
	resp := Response{
		ID:         uuid.UUID(m.ID.Bytes),
		GameID:     uuid.UUID(m.GameID.Bytes),
		MoveNumber: m.MoveNumber,
		Side:       m.Side,
		MoveSan:    m.MoveSan,
		MoveUci:    m.MoveUci,
		Fen:        m.Fen,
	}

	// Handle nullable TimeTaken
	if m.TimeTaken.Valid {
		timeTaken := m.TimeTaken.Int32
		resp.TimeTaken = &timeTaken
	}

	return resp
}

// ToResponses converts slice of database.Move to Response slice.
func ToResponses(moves []database.Move) []Response {
	responses := make([]Response, len(moves))
	for i, m := range moves {
		responses[i] = ToResponse(m)
	}
	return responses
}

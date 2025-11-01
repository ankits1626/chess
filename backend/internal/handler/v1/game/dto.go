// Package game handles game-related HTTP requests for API v1.
package game

import (
	"time"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/google/uuid"
)

// CreateRequest represents game creation request.
type CreateRequest struct {
	WhitePlayerID uuid.UUID `json:"white_player_id" binding:"required"`
	BlackPlayerID uuid.UUID `json:"black_player_id" binding:"required"`
	Pgn           string    `json:"pgn" binding:"required"`
	Result        string    `json:"result,omitempty"`       // Optional: "1-0", "0-1", "1/2-1/2"
	TimeControl   int32     `json:"time_control,omitempty"` // Optional: seconds
}

// UpdateRequest represents game update request.
type UpdateRequest struct {
	Pgn         string `json:"pgn" binding:"required"`
	Result      string `json:"result,omitempty"`
	TimeControl int32  `json:"time_control,omitempty"`
}

// Response represents game in API responses.
type Response struct {
	ID            uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	WhitePlayerID uuid.UUID `json:"white_player_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	BlackPlayerID uuid.UUID `json:"black_player_id" example:"550e8400-e29b-41d4-a716-446655440002"`
	Pgn           string    `json:"pgn" example:"1. e4 e5 2. Nf3 Nc6"`
	Result        *string   `json:"result,omitempty" example:"1-0"`
	TimeControl   *int32    `json:"time_control,omitempty" example:"600"`
	CreatedAt     time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt     time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// ToResponse converts database.Game to Response.
func ToResponse(g database.Game) Response {
	resp := Response{
		ID:            uuid.UUID(g.ID.Bytes),
		WhitePlayerID: uuid.UUID(g.WhitePlayerID.Bytes),
		BlackPlayerID: uuid.UUID(g.BlackPlayerID.Bytes),
		Pgn:           g.Pgn,
		CreatedAt:     g.CreatedAt.Time,
		UpdatedAt:     g.UpdatedAt.Time,
	}

	// Handle nullable Result
	if g.Result.Valid {
		result := g.Result.String
		resp.Result = &result
	}

	// Handle nullable TimeControl
	if g.TimeControl.Valid {
		timeControl := g.TimeControl.Int32
		resp.TimeControl = &timeControl
	}

	return resp
}

// ToResponses converts slice of database.Game to Response slice.
func ToResponses(games []database.Game) []Response {
	responses := make([]Response, len(games))
	for i, g := range games {
		responses[i] = ToResponse(g)
	}
	return responses
}

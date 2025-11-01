// Package user handles user-related HTTP requests for API v1.
package user

import (
	"time"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/google/uuid"
)

// CreateRequest represents user creation request.
type CreateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Rating   int32  `json:"rating,omitempty"`
}

// UpdateRequest represents user update request.
type UpdateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Rating   int32  `json:"rating,omitempty"`
}

// Response represents user in API responses.
type Response struct {
	ID        uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Username  string    `json:"username" example:"johndoe"`
	Email     string    `json:"email" example:"john@example.com"`
	Rating    int32     `json:"rating" example:"1500"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// ToResponse converts database.User to Response.
func ToResponse(u database.User) Response {
	return Response{
		ID:        uuid.UUID(u.ID.Bytes),
		Username:  u.Username,
		Email:     u.Email,
		Rating:    u.Rating.Int32,
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.UpdatedAt.Time,
	}
}

// ToResponses converts slice of database.User to Response slice.
func ToResponses(users []database.User) []Response {
	responses := make([]Response, len(users))
	for i, u := range users {
		responses[i] = ToResponse(u)
	}
	return responses
}

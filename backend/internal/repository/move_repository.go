package repository

import (
	"context"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// MoveRepository handles move business logic.
type MoveRepository struct {
	db *database.DB
}

// NewMoveRepository creates move repository.
func NewMoveRepository(db *database.DB) *MoveRepository {
	return &MoveRepository{db: db}
}

// GetByID retrieves move by ID.
func (r *MoveRepository) GetByID(ctx context.Context, id uuid.UUID) (database.Move, error) {
	return r.db.GetMove(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

// ListByGame retrieves all moves for a game.
func (r *MoveRepository) ListByGame(ctx context.Context, gameID uuid.UUID) ([]database.Move, error) {
	return r.db.ListGameMoves(ctx, pgtype.UUID{Bytes: gameID, Valid: true})
}

// Create creates new move.
func (r *MoveRepository) Create(ctx context.Context, params database.CreateMoveParams) (database.Move, error) {
	return r.db.CreateMove(ctx, params)
}

// DeleteByGame deletes all moves for a game.
func (r *MoveRepository) DeleteByGame(ctx context.Context, gameID uuid.UUID) error {
	return r.db.DeleteGameMoves(ctx, pgtype.UUID{Bytes: gameID, Valid: true})
}

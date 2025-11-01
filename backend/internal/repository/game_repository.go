package repository

import (
	"context"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// GameRepository handles game business logic.
type GameRepository struct {
	db *database.DB
}

// NewGameRepository creates game repository.
func NewGameRepository(db *database.DB) *GameRepository {
	return &GameRepository{db: db}
}

// GetByID retrieves game by ID.
func (r *GameRepository) GetByID(ctx context.Context, id uuid.UUID) (database.Game, error) {
	return r.db.GetGame(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

// List retrieves paginated games.
func (r *GameRepository) List(ctx context.Context, limit, offset int32) ([]database.Game, error) {
	return r.db.ListGames(ctx, database.ListGamesParams{
		Limit:  limit,
		Offset: offset,
	})
}

// ListByUser retrieves user's games.
func (r *GameRepository) ListByUser(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]database.Game, error) {
	return r.db.ListUserGames(ctx, database.ListUserGamesParams{
		WhitePlayerID: pgtype.UUID{Bytes: userID, Valid: true},
		Limit:         limit,
		Offset:        offset,
	})
}

// Create creates new game.
func (r *GameRepository) Create(ctx context.Context, params database.CreateGameParams) (database.Game, error) {
	return r.db.CreateGame(ctx, params)
}

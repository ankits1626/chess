// Package repository implements business logic layer.
package repository

import (
	"context"

	"github.com/ankits1626/chess-coach-backend/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// UserRepository handles user business logic.
type UserRepository struct {
	db *database.DB
}

// NewUserRepository creates user repository.
func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByID retrieves user by ID.
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (database.User, error) {
	return r.db.GetUser(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

// GetByUsername retrieves user by username.
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (database.User, error) {
	return r.db.GetUserByUsername(ctx, username)
}

// GetByEmail retrieves user by email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (database.User, error) {
	return r.db.GetUserByEmail(ctx, email)
}

// List retrieves paginated users.
func (r *UserRepository) List(ctx context.Context, limit, offset int32) ([]database.User, error) {
	return r.db.ListUsers(ctx, database.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
}

// Create creates new user.
func (r *UserRepository) Create(ctx context.Context, username, email string, rating int32) (database.User, error) {
	return r.db.CreateUser(ctx, database.CreateUserParams{
		Username: username,
		Email:    email,
		Rating:   pgtype.Int4{Int32: rating, Valid: true},
	})
}

// Update updates user details.
func (r *UserRepository) Update(ctx context.Context, id uuid.UUID, username, email string, rating int32) (database.User, error) {
	return r.db.UpdateUser(ctx, database.UpdateUserParams{
		ID:       pgtype.UUID{Bytes: id, Valid: true},
		Username: username,
		Email:    email,
		Rating:   pgtype.Int4{Int32: rating, Valid: true},
	})
}

// Delete deletes user by ID.
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.DeleteUser(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

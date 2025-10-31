package game

import (
    "context"
    "database/sql"
    "fmt"

    "chess-coach/backend/internal/db"
    "github.com/google/uuid"
    "github.com/notnil/chess"
)

type Manager struct {
    queries *db.Queries
    parser  *PGNParser
}

func NewManager(queries *db.Queries) *Manager {
    return &Manager{
        queries: queries,
        parser:  NewPGNParser(),
    }
}

func (m *Manager) CreateGame(ctx context.Context, userID uuid.UUID, pgnString string) (*db.Game, error) {
    // Validate PGN
    if err := m.parser.Validate(pgnString); err != nil {
        return nil, fmt.Errorf("invalid PGN: %w", err)
    }

    // Create game in database
    game, err := m.queries.CreateGame(ctx, db.CreateGameParams{
        UserID: userID,
        Pgn:    pgnString,
        Title: sql.NullString{
            String: "New Game",
            Valid:  true,
        },
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create game: %w", err)
    }

    return &game, nil
}

func (m *Manager) GetGame(ctx context.Context, gameID uuid.UUID) (*chess.Game, error) {
    // Fetch from database
    dbGame, err := m.queries.GetGame(ctx, gameID)
    if err != nil {
        return nil, fmt.Errorf("game not found: %w", err)
    }

    // Parse PGN
    game, err := m.parser.Parse(dbGame.Pgn)
    if err != nil {
        return nil, fmt.Errorf("failed to parse game: %w", err)
    }

    return game, nil
}

func (m *Manager) UpdateGame(ctx context.Context, gameID uuid.UUID, pgnString string) error {
    // Validate PGN
    if err := m.parser.Validate(pgnString); err != nil {
        return fmt.Errorf("invalid PGN: %w", err)
    }

    // Update database
    err := m.queries.UpdateGamePGN(ctx, db.UpdateGamePGNParams{
        ID:  gameID,
        Pgn: pgnString,
    })
    if err != nil {
        return fmt.Errorf("failed to update game: %w", err)
    }

    return nil
}

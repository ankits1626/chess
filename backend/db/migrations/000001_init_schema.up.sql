-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    rating INTEGER DEFAULT 1200,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Games table
CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    white_player_id UUID REFERENCES users(id),
    black_player_id UUID REFERENCES users(id),
    pgn TEXT NOT NULL,
    result VARCHAR(10), -- 1-0, 0-1, 1/2-1/2, *
    time_control INTEGER, -- in seconds
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT valid_result CHECK (result IN ('1-0', '0-1', '1/2-1/2', '*'))
);

-- Moves table (for detailed analysis)
CREATE TABLE moves (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    game_id UUID NOT NULL REFERENCES games(id) ON DELETE CASCADE,
    move_number INTEGER NOT NULL,
    side VARCHAR(5) NOT NULL, -- 'white' or 'black'
    move_san VARCHAR(20) NOT NULL, -- Standard Algebraic Notation (e.g., 'e4', 'Nf3')
    move_uci VARCHAR(10) NOT NULL, -- UCI format (e.g., 'e2e4')
    fen TEXT NOT NULL, -- Board position after move
    time_taken INTEGER, -- milliseconds

    CONSTRAINT valid_side CHECK (side IN ('white', 'black'))
);

-- Indexes for performance
CREATE INDEX idx_games_white_player ON games(white_player_id);
CREATE INDEX idx_games_black_player ON games(black_player_id);
CREATE INDEX idx_moves_game_id ON moves(game_id);
CREATE INDEX idx_moves_game_move ON moves(game_id, move_number);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);

-- Updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply triggers
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_games_updated_at BEFORE UPDATE ON games
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
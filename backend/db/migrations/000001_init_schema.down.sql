-- Rollback migration (undo changes)
DROP TRIGGER IF EXISTS update_games_updated_at ON games;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP FUNCTION IF EXISTS update_updated_at_column();

DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_moves_game_move;
DROP INDEX IF EXISTS idx_moves_game_id;
DROP INDEX IF EXISTS idx_games_black_player;
DROP INDEX IF EXISTS idx_games_white_player;

DROP TABLE IF EXISTS moves;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "uuid-ossp";
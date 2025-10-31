-- name: CreateGame :one
INSERT INTO games (user_id, pgn, title)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetGame :one
SELECT * FROM games
WHERE id = $1 LIMIT 1;

-- name: ListUserGames :many
SELECT * FROM games
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateGamePGN :exec
UPDATE games
SET pgn = $2, updated_at = NOW()
WHERE id = $1;

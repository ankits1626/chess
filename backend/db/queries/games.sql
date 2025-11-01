-- name: GetGame :one
SELECT * FROM games
WHERE id = $1 LIMIT 1;

-- name: ListGames :many
SELECT * FROM games
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListUserGames :many
SELECT * FROM games
WHERE white_player_id = $1 OR black_player_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateGame :one
INSERT INTO games (white_player_id, black_player_id, pgn, result, time_control)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateGame :one
UPDATE games
SET pgn = $2, result = $3
WHERE id = $1
RETURNING *;

-- name: DeleteGame :exec
DELETE FROM games
WHERE id = $1;
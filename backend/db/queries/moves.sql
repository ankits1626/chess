-- name: GetMove :one
SELECT * FROM moves
WHERE id = $1 LIMIT 1;

-- name: ListGameMoves :many
SELECT * FROM moves
WHERE game_id = $1
ORDER BY move_number ASC;

-- name: CreateMove :one
INSERT INTO moves (game_id, move_number, side, move_san, move_uci, fen, time_taken)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: DeleteGameMoves :exec
DELETE FROM moves
WHERE game_id = $1;
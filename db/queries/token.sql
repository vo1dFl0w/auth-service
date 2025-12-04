-- name: FindRefreshToken :one
SELECT id, user_id, refresh_token_hash, created_at, expires_at
FROM tokens
WHERE refresh_token_hash = $1;

-- name: SaveHashedRefreshToken :one
INSERT INTO tokens (user_id, refresh_token_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING id, user_id, refresh_token_hash, created_at, expires_at;

-- name: DeleteRefreshToken :one
DELETE FROM tokens
WHERE refresh_token_hash = $1
RETURNING user_id, refresh_token_hash, expires_at;
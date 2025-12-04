-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES($1, $2)
RETURNING user_id, email, created_at, is_active;

-- name: GetUserInfo :one
SELECT user_id, email, created_at, is_active
FROM users
WHERE user_id = $1;

-- name: FindUserByEmail :one
SELECT user_id, email, password_hash, created_at, is_active
FROM users
WHERE email = $1;
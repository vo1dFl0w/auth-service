-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES($1, $2)
RETURNING user_id, email, created_at, is_active, provider_id, oauth_provider;

-- name: CreateOAuthUser :one
INSERT INTO users (email, provider_id, oauth_provider)
VALUES ($1, $2, $3)
RETURNING user_id, email, created_at, is_active, provider_id, oauth_provider;

-- name: GetUserInfo :one
SELECT user_id, email, created_at, is_active, provider_id, oauth_provider
FROM users
WHERE user_id = $1;

-- name: FindUserByEmail :one
SELECT user_id, email, password_hash, created_at, is_active, provider_id, oauth_provider
FROM users
WHERE email = $1;

-- name: FindUserByProviderID :one
SELECT user_id, email, password_hash, created_at, is_active, provider_id, oauth_provider
FROM users
WHERE provider_id = $1 AND oauth_provider = $2;
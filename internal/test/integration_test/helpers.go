package integrationtest

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	pggen "github.com/vo1dFl0w/auth-service/internal/adapters/storage/postgres/pggen"
)

func createUserHelper(t *testing.T, q *pggen.Queries, email string, passwordHash string) pggen.CreateUserRow {
	t.Helper()

	if email == "" {
		email = "user@example.org"
	}

	u, err := q.CreateUser(context.Background(), pggen.CreateUserParams{
		Email:        email,
		PasswordHash: sql.NullString{String: passwordHash, Valid: passwordHash != ""},
	})

	assert.NoError(t, err)
	assert.NotNil(t, u)

	return u
}

func saveHashedRefreshTokenHelper(t *testing.T, q *pggen.Queries, userID uuid.UUID, hash string, expiresAt time.Time) pggen.Token {
	t.Helper()

	if hash == "" {
		hash = "password-hash"
	}

	token, err := q.SaveHashedRefreshToken(context.Background(), pggen.SaveHashedRefreshTokenParams{
		UserID:           userID,
		RefreshTokenHash: hash,
		ExpiresAt:        expiresAt,
	})

	assert.NoError(t, err)
	assert.NotNil(t, token)

	return token
}

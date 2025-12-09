package integrationtest

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vo1dFl0w/auth-service/internal/gen"
)

func createUserHelper(t *testing.T, q *gen.Queries, email string, passwordHash string) gen.CreateUserRow {
	t.Helper()

	if email == "" {
		email = "user@example.org"
	}

	u, err := q.CreateUser(context.Background(), gen.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
	})

	assert.NoError(t, err)
	assert.NotNil(t, u)

	return u
}

func saveHashedRefreshTokenHelper(t *testing.T, q *gen.Queries, userID uuid.UUID, hash string, expiresAt time.Time) gen.Token {
	t.Helper()

	if hash == "" {
		hash = "password-hash"
	}

	token, err := q.SaveHashedRefreshToken(context.Background(), gen.SaveHashedRefreshTokenParams{
		UserID:           userID,
		RefreshTokenHash: hash,
		ExpiresAt:        expiresAt,
	})

	assert.NoError(t, err)
	assert.NotNil(t, token)

	return token
}

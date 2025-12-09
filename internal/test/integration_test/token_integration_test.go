package integrationtest

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vo1dFl0w/auth-service/internal/gen"
)

func TestFindRefreshToken(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := TestDB.BeginTx(ctx, nil)
	assert.NoError(t, err)
	defer tx.Rollback()

	q := gen.New(tx)

	email := "user@example.org"
	passwordHash := "password-hash"

	u := createUserHelper(t, q, email, passwordHash)

	hash := "refresh-token-hash"
	expiresAt := time.Now().UTC().Add(time.Hour * 24 * 7)

	token := saveHashedRefreshTokenHelper(t, q, u.UserID, hash, expiresAt)

	res, err := q.FindRefreshToken(ctx, hash)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, token, res)

	_, err = q.FindRefreshToken(ctx, "not-exists-hash")
	assert.Error(t, err)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestSaveHashedRefreshToken(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := TestDB.BeginTx(ctx, nil)
	assert.NoError(t, err)
	defer tx.Rollback()

	q := gen.New(tx)

	email := "user@example.org"
	passwordHash := "password-hash"

	u := createUserHelper(t, q, email, passwordHash)

	hash := "refresh-token-hash"
	h, err := q.SaveHashedRefreshToken(ctx, gen.SaveHashedRefreshTokenParams{
		UserID:           u.UserID,
		RefreshTokenHash: hash,
		ExpiresAt:        time.Now().UTC().Add(time.Hour * 24 * 7),
	})
	assert.NoError(t, err)
	assert.NotNil(t, h)
}

func TestDeleteRefreshToken(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := TestDB.BeginTx(ctx, nil)
	assert.NoError(t, err)
	defer tx.Rollback()

	q := gen.New(tx)

	email := "user@example.org"
	passwordHash := "password-hash"

	u := createUserHelper(t, q, email, passwordHash)

	hash := "refresh-token-hash"
	expiresAt := time.Now().UTC().Add(time.Hour * 24 * 7)

	token := saveHashedRefreshTokenHelper(t, q, u.UserID, hash, expiresAt)

	res, err := q.DeleteRefreshToken(ctx, hash)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, token.UserID, res.UserID)
	assert.Equal(t, token.RefreshTokenHash, res.RefreshTokenHash)
	assert.Equal(t, token.ExpiresAt, res.ExpiresAt)

	_, err = q.DeleteRefreshToken(ctx, "not-exists-hash")
	assert.Error(t, err)
}

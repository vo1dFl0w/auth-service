package integrationtest

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vo1dFl0w/auth-service/internal/gen"
)

func TestCreateUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := TestDB.BeginTx(ctx, nil)
	assert.NoError(t, err)
	defer tx.Rollback()

	email := "user@example.org"
	passwordHash := "refresh-token-hash"

	q := gen.New(tx)
	res, err := q.CreateUser(ctx, gen.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)

	u, err := q.FindUserByEmail(ctx, res.Email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, res.UserID, u.UserID)
	assert.Equal(t, res.Email, u.Email)
	assert.Equal(t, passwordHash, u.PasswordHash)
	assert.Equal(t, res.CreatedAt, u.CreatedAt)
	assert.Equal(t, res.IsActive, u.IsActive)
}

func TestGetUserInfo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := TestDB.BeginTx(ctx, nil)
	assert.NoError(t, err)
	defer tx.Rollback()

	q := gen.New(tx)

	email := "user@example.org"
	passwordHash := "refresh-token-hash"

	u := createUserHelper(t, q, email, passwordHash)

	res, err := q.GetUserInfo(ctx, u.UserID)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, u.UserID, res.UserID)
	assert.Equal(t, u.Email, res.Email)
	assert.Equal(t, u.CreatedAt, res.CreatedAt)
	assert.Equal(t, u.IsActive, res.IsActive)
}

func TestFindUserByEmail(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := TestDB.BeginTx(ctx, nil)
	assert.NoError(t, err)
	defer tx.Rollback()

	q := gen.New(tx)

	email := "user@example.org"
	passwordHash := "refresh-token-hash"

	u := createUserHelper(t, q, email, passwordHash)

	res, err := q.FindUserByEmail(ctx, u.Email)
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, u.UserID, res.UserID)
	assert.Equal(t, u.Email, res.Email)
	assert.Equal(t, passwordHash, res.PasswordHash)
	assert.Equal(t, u.CreatedAt, res.CreatedAt)
	assert.Equal(t, u.IsActive, res.IsActive)
}

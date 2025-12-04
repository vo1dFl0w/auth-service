package usecase_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vo1dFl0w/auth-service/internal/app/domain"
	"github.com/vo1dFl0w/auth-service/internal/app/usecase"
	"github.com/vo1dFl0w/auth-service/internal/test/mocks"
)

func TestTokenRepository_GenerateAccessToken(t *testing.T) {
	tokenRepo := &mocks.TokenRepositoryMock{}
	tokenService := usecase.NewTokenService([]byte("very-secret-key"), tokenRepo)

	userID := uuid.New()

	accessToken, err := tokenService.GenerateAccessToken(userID)
	assert.NoError(t, err)
	assert.NotNil(t, accessToken)
}

func TestTokenRepository_GenerateRefreshToken(t *testing.T) {
	tokenRepo := &mocks.TokenRepositoryMock{}
	tokenService := usecase.NewTokenService([]byte("very-secret-key"), tokenRepo)

	refreshToken, err := tokenService.GenerateRefreshToken()
	assert.NoError(t, err)
	assert.NotNil(t, refreshToken)
}

func TestTokenRepository_HashRefreshTokent(t *testing.T) {
	tokenRepo := &mocks.TokenRepositoryMock{}
	tokenService := usecase.NewTokenService([]byte("very-secret-key"), tokenRepo)

	refreshToken, err := tokenService.GenerateRefreshToken()
	assert.NoError(t, err)
	assert.NotNil(t, refreshToken)

	hash, expiresAt := tokenService.HashRefreshToken(refreshToken)
	assert.NotNil(t, hash)
	assert.NotNil(t, expiresAt)
}

func TestTokenRepository_ValidateAccessToken(t *testing.T) {
	tokenRepo := &mocks.TokenRepositoryMock{}
	tokenService := usecase.NewTokenService([]byte("very-secret-key"), tokenRepo)

	userID := uuid.New()

	accessToken, err := tokenService.GenerateAccessToken(userID)
	assert.NoError(t, err)
	assert.NotNil(t, accessToken)

	res, err := tokenService.ValidateAccessToken(accessToken)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	_, err = tokenService.ValidateAccessToken("fake token")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidAccessToken)
}

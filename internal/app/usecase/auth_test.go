package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vo1dFl0w/auth-service/internal/app/domain"
	"github.com/vo1dFl0w/auth-service/internal/app/repository"
	"github.com/vo1dFl0w/auth-service/internal/app/usecase"
	"github.com/vo1dFl0w/auth-service/internal/test/mocks"
)

func TestAuthRepository_Register(t *testing.T) {
	authRepo := &mocks.AuthRepositoryMock{}
	tokenRepo := &mocks.TokenRepositoryMock{}
	tokenService := usecase.NewTokenService([]byte("very-secret-key"), tokenRepo)
	authService := usecase.NewAuthService(authRepo, tokenRepo, tokenService)

	email := "user@example.org"
	password := "password"

	u := &domain.User{
		UserID:    uuid.New(),
		Email:     email,
		CreatedAt: time.Now().UTC(),
		IsActive:  true,
	}

	authRepo.On("CreateUser", mock.Anything, email, mock.Anything).Return(u, nil).Once()

	res, err := authService.Register(context.Background(), email, password)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, u, res)

	_, err = authService.Register(context.Background(), "invalid-email", password)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidEmail)

	_, err = authService.Register(context.Background(), email, "")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidPassword)

	authRepo.AssertExpectations(t)
}

func TestAuthRepository_Login(t *testing.T) {
	authRepo := &mocks.AuthRepositoryMock{}
	tokenRepo := &mocks.TokenRepositoryMock{}
	tokenService := &mocks.TokenServiceMock{}
	authService := usecase.NewAuthService(authRepo, tokenRepo, tokenService)

	email := "user@example.org"
	password := "password"
	accessToken := "access-token"
	refreshToken := "refresh-token"
	refreshTokenHash := "hashed-refresh-token"
	expiresAt := time.Now().UTC().Add(time.Hour * 24 * 7)

	hash, _ := usecase.HashPassword(password)

	u := &domain.UserWithPassword{
		UserID:       uuid.New(),
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    time.Now().UTC(),
		IsActive:     true,
	}

	authRepo.On("FindUserByEmail", mock.Anything, email).Return(u, nil).Once()
	tokenService.On("GenerateAccessToken", u.UserID).Return(accessToken, nil).Once()
	tokenService.On("GenerateRefreshToken").Return(refreshToken, nil).Once()
	tokenService.On("HashRefreshToken", refreshToken).Return(refreshTokenHash, expiresAt).Once()
	tokenRepo.On("SaveHashedRefreshToken", mock.Anything, u.UserID, refreshTokenHash, expiresAt).Return(nil).Once()

	res, err := authService.Login(context.Background(), email, password)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, refreshToken, res.RefreshToken)
	assert.Equal(t, accessToken, res.AccessToken)
	assert.Equal(t, expiresAt, res.RefreshTokenExpiresAt)

	_, err = authService.Login(context.Background(), "invalid email", password)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrWrongEmailOrPassword)

	_, err = authService.Login(context.Background(), email, "")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrWrongEmailOrPassword)

	authRepo.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
	tokenService.AssertExpectations(t)
}

func TestAuthRepository_Logout(t *testing.T) {
	authRepo := &mocks.AuthRepositoryMock{}
	tokenRepo := &mocks.TokenRepositoryMock{}
	tokenService := &mocks.TokenServiceMock{}
	authService := usecase.NewAuthService(authRepo, tokenRepo, tokenService)

	refreshToken := "refresh-token"
	hash := usecase.HashRefreshTokenFunc(refreshToken)

	tokenRepo.On("DeleteRefreshToken", mock.Anything, hash).Return(nil, nil).Once()

	err := authService.Logout(context.Background(), refreshToken)
	assert.NoError(t, err)

	tokenRepo.AssertExpectations(t)
}

func TestAuthRepository_UserInfo(t *testing.T) {
	authRepo := &mocks.AuthRepositoryMock{}
	tokenRepo := &mocks.TokenRepositoryMock{}
	tokenService := &mocks.TokenServiceMock{}
	authService := usecase.NewAuthService(authRepo, tokenRepo, tokenService)

	userID := uuid.New()
	u := &domain.User{
		UserID:    userID,
		Email:     "user@example.org",
		CreatedAt: time.Now().UTC(),
		IsActive:  true,
	}

	authRepo.On("GetUserInfo", mock.Anything, userID).Return(u, nil).Once()

	res, err := authService.UserInfo(context.Background(), userID)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, u, res)

	authRepo.On("GetUserInfo", mock.Anything, uuid.Nil).Return(nil, repository.ErrNotFound).Once()

	_, err = authService.UserInfo(context.Background(), uuid.Nil)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrWrongUserID)

	authRepo.AssertExpectations(t)
}

func TestAuthRepository_RefreshTokens(t *testing.T) {
	authRepo := &mocks.AuthRepositoryMock{}
	tokenRepo := &mocks.TokenRepositoryMock{}
	tokenService := &mocks.TokenServiceMock{}
	authService := usecase.NewAuthService(authRepo, tokenRepo, tokenService)

	refreshToken := "refresh-token"
	hash := usecase.HashRefreshTokenFunc(refreshToken)
	expiresAt := time.Now().UTC().Add(time.Hour * 24)

	accessToken := "access-token"
	newRefreshToken := "new-refresh-token"
	newHash := usecase.HashRefreshTokenFunc(newRefreshToken)
	newExpiresAt := time.Now().UTC().Add(time.Hour * 24 * 7)

	ref := &domain.RefreshToken{
		UserID:       uuid.New(),
		RefreshToken: hash,
		ExpiresAt:    expiresAt,
	}

	tokens := &domain.Tokens{
		AccessToken:           accessToken,
		RefreshToken:          newRefreshToken,
		RefreshTokenExpiresAt: newExpiresAt,
	}

	tokenRepo.On("DeleteRefreshToken", mock.Anything, hash).Return(ref, nil).Once()
	tokenService.On("GenerateAccessToken", ref.UserID).Return(accessToken, nil).Once()
	tokenService.On("GenerateRefreshToken").Return(newRefreshToken, nil).Once()
	tokenService.On("HashRefreshToken", newRefreshToken).Return(newHash, newExpiresAt).Once()
	tokenRepo.On("SaveHashedRefreshToken", mock.Anything, ref.UserID, newHash, newExpiresAt).Return(nil).Once()

	res, err := authService.RefreshTokens(context.Background(), refreshToken)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, tokens, res)

	_, err = authService.RefreshTokens(context.Background(), "")
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrEmptyRefreshToken)
	
	fakeToken := "fake-token"
	fakeHash := usecase.HashRefreshTokenFunc(fakeToken)
	tokenRepo.On("DeleteRefreshToken", mock.Anything, fakeHash).Return(nil, repository.ErrNoRowDeleted).Once()

	_, err = authService.RefreshTokens(context.Background(), fakeToken)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidOrExpiredRefreshToken)

	tokenRepo.AssertExpectations(t)
	tokenService.AssertExpectations(t)

}

package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/auth-service/internal/app/domain"
	"github.com/vo1dFl0w/auth-service/internal/app/repository"
)

type AuthService interface {
	Register(ctx context.Context, email string, password string) (*domain.User, error)
	Login(ctx context.Context, email string, password string) (*domain.Tokens, error)
	Logout(ctx context.Context, token string) error
	UserInfo(ctx context.Context, user_id uuid.UUID) (*domain.User, error)
	RefreshTokens(ctx context.Context, token string) (*domain.Tokens, error)
}

type authService struct {
	authRepo     repository.AuthRepository
	tokenRepo    repository.TokenRepository
	tokenService TokenService
}

func NewAuthService(authRepo repository.AuthRepository, tokenRepo repository.TokenRepository, tokenService TokenService) *authService {
	return &authService{
		authRepo:     authRepo,
		tokenRepo:    tokenRepo,
		tokenService: tokenService,
	}
}

func (s *authService) Register(ctx context.Context, email string, password string) (*domain.User, error) {
	u := &domain.User{
		Email:    email,
		Password: password,
	}

	if err := validateEmail(u); err != nil {
		return nil, domain.ErrInvalidEmail
	}

	if err := validatePassword(u); err != nil {
		return nil, domain.ErrInvalidPassword
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		if errors.Is(err, domain.ErrEmptyPassword) {
			return nil, domain.ErrEmptyPassword
		}
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u.Password = ""

	res, err := s.authRepo.CreateUser(ctx, email, hashedPassword)
	if err != nil {
		if errors.Is(err, repository.ErrEmailAlreadyExists) {
			return nil, domain.ErrEmailAlreadyExists
		} else {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	}

	return res, nil
}

func (s *authService) Login(ctx context.Context, email string, password string) (*domain.Tokens, error) {
	u := &domain.User{
		Email:    email,
		Password: password,
	}

	if err := validateEmail(u); err != nil {
		return nil, domain.ErrWrongEmailOrPassword
	}

	if err := validatePassword(u); err != nil {
		return nil, domain.ErrWrongEmailOrPassword
	}

	res, err := s.authRepo.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, domain.ErrWrongEmailOrPassword
		} else {
			return nil, fmt.Errorf("find user by email: %w", err)
		}
	}

	if !comparePasswords(password, res.PasswordHash) {
		return nil, domain.ErrWrongEmailOrPassword
	}

	accessToken, err := s.tokenService.GenerateAccessToken(res.UserID)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshToken, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	h, expiresAt := s.tokenService.HashRefreshToken(refreshToken)
	if err := s.tokenRepo.SaveHashedRefreshToken(ctx, res.UserID, h, expiresAt); err != nil {
		return nil, fmt.Errorf("save hashed refresh token: %w", err)
	}

	return &domain.Tokens{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: expiresAt,
	}, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	if token == "" {
		return domain.ErrEmptyRefreshToken
	}

	hash := HashRefreshTokenFunc(token)

	_, err := s.tokenRepo.DeleteRefreshToken(ctx, hash)
	if err != nil {
		if errors.Is(err, repository.ErrNoRowDeleted) {
			return nil
		} else {
			return fmt.Errorf("delete refresh token: %w", err)
		}
	}

	return nil
}

func (s *authService) UserInfo(ctx context.Context, user_id uuid.UUID) (*domain.User, error) {
	u, err := s.authRepo.GetUserInfo(ctx, user_id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, domain.ErrWrongUserID
		}
		return nil, fmt.Errorf("get user info: %w", err)
	}

	return u, nil
}

func (s *authService) RefreshTokens(ctx context.Context, token string) (*domain.Tokens, error) {
	if token == "" {
		return nil, domain.ErrEmptyRefreshToken
	}

	hashToken := HashRefreshTokenFunc(token)
	res, err := s.tokenRepo.DeleteRefreshToken(ctx, hashToken)
	if err != nil {
		if errors.Is(err, repository.ErrNoRowDeleted) {
			return nil, domain.ErrInvalidOrExpiredRefreshToken
		} else {
			return nil, fmt.Errorf("delete refresh token: %w", err)
		}
	}

	if time.Now().After(res.ExpiresAt) {
		return nil, domain.ErrInvalidOrExpiredRefreshToken
	}

	accessToken, err := s.tokenService.GenerateAccessToken(res.UserID)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}
	refreshToken, err := s.tokenService.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	h, expiresAt := s.tokenService.HashRefreshToken(refreshToken)

	if err := s.tokenRepo.SaveHashedRefreshToken(ctx, res.UserID, h, expiresAt); err != nil {
		return nil, fmt.Errorf("save hashed refresh token: %w", err)
	}

	return &domain.Tokens{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: expiresAt,
	}, nil
}

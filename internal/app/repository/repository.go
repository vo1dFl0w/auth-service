package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/auth-service/internal/app/domain"
)


type AuthRepository interface {
	CreateUser(ctx context.Context, email string, passwordHash string) (*domain.User, error)
	GetUserInfo(ctx context.Context, user_id uuid.UUID) (*domain.User, error)
	FindUserByEmail(ctx context.Context, email string) (*domain.UserWithPassword, error)
}

type TokenRepository interface {
	FindRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error)
	SaveHashedRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	DeleteRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error)
}

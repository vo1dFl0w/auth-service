package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vo1dFl0w/auth-service/internal/domain"
)


type AuthRepository interface {
	CreateUser(ctx context.Context, email string, passwordHash string) (*domain.User, error)
	CreateOAuthUser(ctx context.Context, email string, providerID string, oauthProvider string) (*domain.User, error)
	GetUserInfo(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	FindUserByEmail(ctx context.Context, email string) (*domain.UserWithPassword, error)
	FindUserByProviderID(ctx context.Context, providerID string, oauthProvider string) (*domain.User, error)
}

type TokenRepository interface {
	FindRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error)
	SaveHashedRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	DeleteRefreshToken(ctx context.Context, tokenHash string) (*domain.RefreshToken, error)
}

type OAuthRepository interface {
	AuthCodeURL(ctx context.Context, state string) string
	GenerateState(ctx context.Context) (string, error)
	GetUserFromCode(ctx context.Context, code string) (*domain.ProviderUser, error)
}

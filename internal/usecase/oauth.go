package usecase

import (
	"context"
	"fmt"

	"github.com/vo1dFl0w/auth-service/internal/domain"
	"github.com/vo1dFl0w/auth-service/internal/repository"
)

type OAuthService interface {
	GetAuthCodeURL(ctx context.Context, state string) string
	GetGenerateState(ctx context.Context) (string, error)
	GetUserFromCode(ctx context.Context, code string, state string) (*domain.ProviderUser, error)
}

type oauthService struct {
	oauthRepo repository.OAuthRepository
}

func NewOAuthService(oauthRepo repository.OAuthRepository) *oauthService {
	return &oauthService{oauthRepo: oauthRepo}
}

func (s *oauthService) GetAuthCodeURL(ctx context.Context, state string) string {
	return s.oauthRepo.AuthCodeURL(ctx, state)
}

func (s *oauthService) GetGenerateState(ctx context.Context) (string, error) {
	return s.oauthRepo.GenerateState(ctx)
}

func (s *oauthService) GetUserFromCode(ctx context.Context, code string, state string) (*domain.ProviderUser, error) {
	if err := s.oauthRepo.ValidateCallbackParams(ctx, code, state); err != nil {
		return nil, fmt.Errorf("validate callback params: %w", err)
	}

	u, err := s.oauthRepo.GetUserFromCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("get user from code: %w", err)
	}

	return u, err
}

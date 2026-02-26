package google

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/vo1dFl0w/auth-service/internal/config"
	"github.com/vo1dFl0w/auth-service/internal/domain"
	"github.com/vo1dFl0w/auth-service/internal/repository"
)

const (
	googleProvider = "google"
)

type Google struct {
	cfg             *Config
	client          Client
	baseUserInfoURL string
}

func NewGoogleRepository(cfg *config.Config) *Google {
	return &Google{
		cfg:             NewOAuthGoogleConfig(cfg),
		client:          NewClient(cfg),
		baseUserInfoURL: cfg.GoogleOAuth.BaseUserInfoURL,
	}
}

func (g *Google) AuthCodeURL(ctx context.Context, state string) string {
	return g.cfg.Cfg.AuthCodeURL(state)
}

func (g *Google) GenerateState(ctx context.Context) (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func (g *Google) GetUserFromCode(ctx context.Context, code string) (*domain.ProviderUser, error) {
	token, err := g.cfg.Cfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchange code: %w", err)
	}

	if token == nil || token.AccessToken == "" {
		return nil, repository.ErrEmptyOAuthAccessToken
	}

	ui, err := g.client.FetchUserProfile(ctx, token.AccessToken, g.baseUserInfoURL)
	if err != nil {
		return nil, fmt.Errorf("fetch user profile: %w", err)
	}

	providerID := ui.Sub
	if providerID == "" {
		providerID = ui.ID
	}

	return &domain.ProviderUser{
		Provider:     googleProvider,
		ProviderID:   providerID,
		Email:        ui.Email,
		Name:         ui.Name,
		FirstName:    ui.GivenName,
		LastName:     ui.FamilyName,
		Picture:      ui.Picture,
		Verified:     ui.EmailVerified,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenExpiry:  token.Expiry,
	}, nil
}

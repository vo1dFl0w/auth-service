package google

import (
	"github.com/vo1dFl0w/auth-service/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	Cfg *oauth2.Config
}

func NewOAuthGoogleConfig(cfg *config.Config) *Config {
	return &Config{
		Cfg: &oauth2.Config{
			RedirectURL:  cfg.GoogleOAuth.RedirectURL,
			ClientID:     cfg.GoogleOAuth.ClientID,
			ClientSecret: cfg.GoogleOAuth.ClientSecret,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
	}
}

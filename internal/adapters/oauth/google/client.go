package google

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vo1dFl0w/auth-service/internal/config"
)

type Client interface {
	FetchUserProfile(ctx context.Context, accessToken string, url string) (*UserInfo, error)
}

type client struct {
	httpClient *http.Client
}

func NewClient(cfg *config.Config) *client {
	return &client{httpClient: &http.Client{
		Timeout: cfg.Server.RequestTimeout,
	}}
}

func (c *client) FetchUserProfile(ctx context.Context, accessToken string, url string) (*UserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("new request with context: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client do request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	ui := &UserInfo{}

	if err := json.NewDecoder(resp.Body).Decode(&ui); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return ui, nil
}

package google

import "github.com/vo1dFl0w/auth-service/internal/repository"

type OAuth struct {
	oauthRepo repository.OAuthRepository
}

func NewOAuth(oauthRepo repository.OAuthRepository) *OAuth {
	return &OAuth{oauthRepo: oauthRepo}
}

func (oa *OAuth) Google() repository.OAuthRepository {
	return oa.oauthRepo
}
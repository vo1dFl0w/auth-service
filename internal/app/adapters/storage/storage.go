package storage

import "github.com/vo1dFl0w/auth-service/internal/app/repository"

type Storage interface {
	Auth() repository.AuthRepository
	Token() repository.TokenRepository
}
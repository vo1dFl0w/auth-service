package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID
	Email     string
	Password  string
	CreatedAt time.Time
	IsActive  bool
}

type UserWithPassword struct {
	UserID       uuid.UUID
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	IsActive     bool
}

type Tokens struct {
	AccessToken           string
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

type RefreshToken struct {
	UserID       uuid.UUID
	RefreshToken string
	ExpiresAt    time.Time
}

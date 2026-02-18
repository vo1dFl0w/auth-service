package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID        uuid.UUID
	Email         string
	Password      string
	CreatedAt     time.Time
	IsActive      bool
	ProviderID    string
	OAuthProvider string
}

type UserWithPassword struct {
	UserID        uuid.UUID
	Email         string
	PasswordHash  string
	CreatedAt     time.Time
	IsActive      bool
	ProviderID    string
	OAuthProvider string
}

type ProviderUser struct {
	Provider     string
	ProviderID   string
	Email        string
	Name         string
	FirstName    string
	LastName     string
	Picture      string
	Verified     bool
	RawJSON      []byte
	AccessToken  string
	RefreshToken string
	TokenExpiry  time.Time
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

package usecase

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/vo1dFl0w/auth-service/internal/app/domain"
	"github.com/vo1dFl0w/auth-service/internal/app/repository"
)

type TokenService interface {
	GenerateAccessToken(userID uuid.UUID) (string, error)
	GenerateRefreshToken() (string, error)
	HashRefreshToken(refreshToken string) (string, time.Time)
	ValidateAccessToken(accessToken string) (*jwt.RegisteredClaims, error)
}

type tokenService struct {
	jwtSecret []byte
	tokenRepo repository.TokenRepository
}

func NewTokenService(jwtSecret []byte, tokenRepo repository.TokenRepository) TokenService {
	return &tokenService{
		jwtSecret: jwtSecret,
		tokenRepo: tokenRepo,
	}
}

func (s *tokenService) GenerateAccessToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	return token.SignedString(s.jwtSecret)
}

func (s *tokenService) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (s *tokenService) ValidateAccessToken(accessToken string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}

	t, err := jwt.ParseWithClaims(accessToken, claims, func(t *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired){
			return nil, domain.ErrExpiredAccessToken
		} else {
			return nil, domain.ErrInvalidAccessToken
		}
	}

	if !t.Valid {
		return nil, domain.ErrInvalidAccessToken
	}

	return claims, nil
}

func (s *tokenService) HashRefreshToken(refreshToken string) (string, time.Time) {
	h := HashRefreshTokenFunc(refreshToken)
	expiry := time.Now().Add(time.Hour * 24 * 7)

	return h, expiry
}

func HashRefreshTokenFunc(refreshToken string) string {
	h := sha256.Sum256([]byte(refreshToken))
	return hex.EncodeToString(h[:])
}

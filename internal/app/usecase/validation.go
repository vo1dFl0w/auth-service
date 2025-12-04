package usecase

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/vo1dFl0w/auth-service/internal/app/domain"
	"golang.org/x/crypto/bcrypt"
)

func validateEmail(u *domain.User) error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
	)
}

func validatePassword(u *domain.User) error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Password, validation.Required, validation.Length(8, 100)),
	)
}

func HashPassword(password string) (string, error) {
	if len(password) == 0 {
		return "", domain.ErrEmptyPassword
	}

	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate hashed password: %w", err)
	}

	return string(h), nil
}

func comparePasswords(password, hashedPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

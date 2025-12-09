package usecase

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vo1dFl0w/auth-service/internal/app/domain"
)

func TestUsecase_ValidateEmail(t *testing.T) {
	testCases := []struct {
		name   string
		u      *domain.User
		expErr bool
	}{
		{
			name: "valid-email",
			u: &domain.User{
				Email: "user@example.org",
			},
			expErr: false,
		},
		{
			name: "invalid-email",
			u: &domain.User{
				Email: "invalid-email",
			},
			expErr: true,
		},
		{
			name: "empty-email",
			u: &domain.User{
				Email: "",
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.expErr {
				err := validateEmail(tc.u)
				assert.NoError(t, err)
			} else {
				err := validateEmail(tc.u)
				assert.Error(t, err)
			}
		})
	}
}

func TestUsecase_ValidatePassword(t *testing.T) {
	buffer := make([]byte, 101)
	_, _ = rand.Read(buffer)

	testCases := []struct {
		name   string
		u      *domain.User
		expErr bool
	}{
		{
			name: "valid-password",
			u: &domain.User{
				Password: "password",
			},
			expErr: false,
		},
		{
			name: "to-short-password",
			u: &domain.User{
				Password: "short",
			},
			expErr: true,
		},
		{
			name: "to-long-password",
			u: &domain.User{
				Password: string(buffer),
			},
			expErr: true,
		},
		{
			name: "empty-password",
			u: &domain.User{
				Password: "",
			},
			expErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.expErr {
				err := validatePassword(tc.u)
				assert.NoError(t, err)
			} else {
				err := validatePassword(tc.u)
				assert.Error(t, err)
			}
		})
	}
}

func TestUsecase_HashPassword(t *testing.T) {
	password := "password"

	hash, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEqual(t, "", hash)

	err = comparePasswords(password, hash)
	assert.NoError(t, err)

	_, err = HashPassword("")
	assert.Error(t, err)
}

func TestUsecase_ComparePasswords(t *testing.T) {
	password := "password"

	hash, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEqual(t, "", hash)

	err = comparePasswords(password, hash)
	assert.NoError(t, err)

	err = comparePasswords("wrong-password", hash)
	assert.Error(t, err)
}

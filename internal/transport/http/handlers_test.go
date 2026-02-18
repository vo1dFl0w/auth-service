package http_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vo1dFl0w/auth-service/internal/config"
	"github.com/vo1dFl0w/auth-service/internal/domain"
	"github.com/vo1dFl0w/auth-service/internal/test/mocks"
	ht "github.com/vo1dFl0w/auth-service/internal/transport/http"
	"github.com/vo1dFl0w/auth-service/internal/transport/http/httpgen"
)

func TestHandlers_APIV1AuthLoginPost(t *testing.T) {
	cfg := &config.Config{}
	log := &mocks.LoggerMock{}

	testCases := []struct {
		name     string
		email    string
		password string
		expErr   bool
	}{
		{
			name:     "valid",
			email:    "user@example.org",
			password: "password",
			expErr:   false,
		},
		{
			name:     "invalid email",
			email:    "",
			password: "password",
			expErr:   true,
		},
		{
			name:     "invalid password",
			email:    "user@example.org",
			password: "",
			expErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authService := &mocks.AuthServiceMock{}
			oauthService := &mocks.OAuthServiceMock{}

			handler := ht.NewHandler(cfg, log, authService, oauthService)

			if !tc.expErr {
				tokens := &domain.Tokens{
					AccessToken:           "access-token",
					RefreshToken:          "refresh-token",
					RefreshTokenExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 7),
				}

				authService.On("Login", mock.Anything, tc.email, tc.password).Return(tokens, nil).Once()
				res, err := handler.APIV1AuthLoginPost(context.Background(), &httpgen.LoginRequest{
					Email:    tc.email,
					Password: tc.password,
				})
				assert.NoError(t, err)
				assert.NotNil(t, res)

				headers, ok := res.(*httpgen.AccessTokenHeaders)
				assert.True(t, ok)
				assert.Equal(t, tokens.AccessToken, headers.Response.AccessToken)

				cookieStr, ok := headers.SetCookie.Get()
				assert.True(t, ok)

				cookie, err := http.ParseSetCookie(cookieStr)
				assert.NoError(t, err)
				assert.NotNil(t, cookie)
				assert.WithinDuration(t, tokens.RefreshTokenExpiresAt, cookie.Expires, time.Second)
				assert.Equal(t, tokens.RefreshToken, cookie.Value)

				authService.AssertExpectations(t)
			} else {
				authService.On("Login", mock.Anything, tc.email, tc.password).Return(nil, domain.ErrWrongEmailOrPassword).Once()
				log.On("Warn", mock.Anything, mock.Anything).Once()
				res, err := handler.APIV1AuthLoginPost(context.Background(), &httpgen.LoginRequest{
					Email:    tc.email,
					Password: tc.password,
				})
				assert.NoError(t, err)
				assert.NotNil(t, res)

				_, ok := res.(*httpgen.APIV1AuthLoginPostUnauthorized)
				assert.True(t, ok)

				authService.AssertExpectations(t)
				log.AssertExpectations(t)
			}
		})
	}
}

func TestHandlers_APIV1AuthLogoutPost(t *testing.T) {
	cfg := &config.Config{}
	log := &mocks.LoggerMock{}

	authService := &mocks.AuthServiceMock{}
	oauthService := &mocks.OAuthServiceMock{}

	handler := ht.NewHandler(cfg, log, authService, oauthService)

	refreshToken := "refresh-token"

	authService.On("Logout", mock.Anything, refreshToken).Return(nil).Once()

	res, err := handler.APIV1AuthLogoutPost(context.Background(), httpgen.APIV1AuthLogoutPostParams{
		RefreshToken: refreshToken,
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)

	_, ok := res.(*httpgen.APIV1AuthLogoutPostNoContent)
	assert.True(t, ok)

	authService.AssertExpectations(t)
}

func TestHandlers_APIV1AuthMeGet(t *testing.T) {
	cfg := &config.Config{}
	log := &mocks.LoggerMock{}

	authService := &mocks.AuthServiceMock{}
	oauthService := &mocks.OAuthServiceMock{}

	handler := ht.NewHandler(cfg, log, authService, oauthService)

	userID := uuid.New()
	u := &domain.User{
		UserID:    userID,
		Email:     "user@example.org",
		CreatedAt: time.Now().UTC(),
		IsActive:  true,
	}

	authService.On("UserInfo", mock.Anything, userID).Return(u, nil).Once()

	ctx := context.WithValue(context.Background(), ht.CtxKeyUserID, userID.String())

	res, err := handler.APIV1AuthMeGet(ctx)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	resp, ok := res.(*httpgen.UserInfoResponse)
	assert.True(t, ok)

	assert.Equal(t, u.UserID.String(), resp.UserID)
	assert.Equal(t, u.Email, resp.Email)

	authService.AssertExpectations(t)
}

func TestHandlers_APIV1AuthRefreshPost(t *testing.T) {
	cfg := &config.Config{}
	log := &mocks.LoggerMock{}

	authService := &mocks.AuthServiceMock{}
	oauthService := &mocks.OAuthServiceMock{}

	handler := ht.NewHandler(cfg, log, authService, oauthService)

	refreshToken := "refresh-token"
	accessToken := "access-token"
	newRefreshToken := "new-refresh-token"
	expiresAt := time.Now().UTC().Add(7 * 24 * time.Hour)

	tokens := &domain.Tokens{
		AccessToken:           accessToken,
		RefreshToken:          newRefreshToken,
		RefreshTokenExpiresAt: expiresAt,
	}

	authService.On("RefreshTokens", mock.Anything, refreshToken).Return(tokens, nil).Once()

	res, err := handler.APIV1AuthRefreshPost(context.Background(), httpgen.APIV1AuthRefreshPostParams{
		RefreshToken: refreshToken,
	})
	assert.NoError(t, err)
	assert.NotNil(t, res)

	headers, ok := res.(*httpgen.AccessTokenHeaders)
	assert.True(t, ok)
	assert.Equal(t, accessToken, headers.Response.AccessToken)

	cookieStr, ok := headers.SetCookie.Get()
	assert.True(t, ok)

	cookie, err := http.ParseSetCookie(cookieStr)
	assert.NoError(t, err)
	assert.NotNil(t, cookie)

	assert.Equal(t, tokens.RefreshToken, cookie.Value)
	assert.WithinDuration(t, tokens.RefreshTokenExpiresAt, cookie.Expires, time.Second)

	authService.AssertExpectations(t)
}

func TestHandlers_APIV1AuthRegisterPost(t *testing.T) {
	cfg := &config.Config{}
	log := &mocks.LoggerMock{}

	testCases := []struct {
		name     string
		email    string
		password string
		expErr   bool
	}{
		{
			name:     "valid",
			email:    "user@example.org",
			password: "password",
			expErr:   false,
		},
		{
			name:     "invalid email",
			email:    "",
			password: "password",
			expErr:   true,
		},
		{
			name:     "invalid password",
			email:    "user@example.org",
			password: "",
			expErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			authService := &mocks.AuthServiceMock{}
			oauthService := &mocks.OAuthServiceMock{}

			handler := ht.NewHandler(cfg, log, authService, oauthService)

			if !tc.expErr {
				userID := uuid.New()
				createdAt := time.Now().UTC()

				u := &domain.User{
					UserID:    userID,
					Email:     tc.email,
					CreatedAt: createdAt,
					IsActive:  true,
				}

				authService.On("Register", mock.Anything, tc.email, tc.password).Return(u, nil).Once()

				res, err := handler.APIV1AuthRegisterPost(context.Background(), &httpgen.RegisterRequest{
					Email:    tc.email,
					Password: tc.password,
				})

				assert.NoError(t, err)
				assert.NotNil(t, res)

				resp, ok := res.(*httpgen.RegisterResponse)
				assert.True(t, ok)
				assert.Equal(t, u.UserID.String(), resp.UserID)
				assert.Equal(t, u.Email, resp.Email)

				assert.WithinDuration(t, u.CreatedAt, resp.CreatedAt, time.Second)

				authService.AssertExpectations(t)
			} else {
				var err error
				if tc.email == "" {
					err = domain.ErrInvalidPassword
				} else if tc.password == "" {
					err = domain.ErrInvalidPassword
				}

				authService.On("Register", mock.Anything, tc.email, tc.password).Return(nil, err).Once()
				log.On("Warn", mock.Anything, mock.Anything).Once()

				res, err := handler.APIV1AuthRegisterPost(context.Background(), &httpgen.RegisterRequest{
					Email:    tc.email,
					Password: tc.password,
				})

				_, ok := res.(*httpgen.APIV1AuthRegisterPostBadRequest)
				assert.True(t, ok)
				authService.AssertExpectations(t)
				log.AssertExpectations(t)
			}
		})
	}
}

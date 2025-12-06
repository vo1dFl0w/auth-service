package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/cors"
	"github.com/vo1dFl0w/auth-service/internal/app/domain"
	"github.com/vo1dFl0w/auth-service/internal/app/usecase"
	"github.com/vo1dFl0w/auth-service/internal/config"
	"github.com/vo1dFl0w/auth-service/internal/gen"
)

var (
	ctxDuration = time.Second * 5
)

type Handler struct {
	cfg          *config.Config
	log          *slog.Logger
	cors         *cors.Cors
	authService  usecase.AuthService
	cookieSecure bool
}

func NewHandler(cfg *config.Config, log *slog.Logger, authService usecase.AuthService) *Handler {
	opts := cors.Options{
		AllowedOrigins:   cfg.Cors.AllowedOrigins,
		AllowedMethods:   cfg.Cors.AllowedMethods,
		AllowedHeaders:   cfg.Cors.AllowedHeaders,
		ExposedHeaders:   cfg.Cors.ExposedHeaders,
		AllowCredentials: cfg.Cors.AllowCredentials,
		MaxAge:           cfg.Cors.MaxAge,
	}

	c := cors.New(opts)

	return &Handler{
		cfg:          cfg,
		log:          log,
		cors:         c,
		authService:  authService,
		cookieSecure: cfg.Cookie.CookieSecure,
	}
}

func (h *Handler) APIV1AuthRegisterPost(ctx context.Context, req *gen.RegisterRequest) (gen.APIV1AuthRegisterPostRes, error) {
	u, err := h.authService.Register(ctx, string(req.Email), req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			return &gen.APIV1AuthRegisterPostConflict{
				Message: domain.ErrEmailAlreadyExists.Error(),
				Status:  http.StatusConflict,
			}, nil
		} else if errors.Is(err, domain.ErrInvalidEmail) || errors.Is(err, domain.ErrInvalidPassword) {
			return &gen.APIV1AuthRegisterPostBadRequest{
				Message: err.Error(),
				Status:  http.StatusBadRequest,
			}, nil
		} else {
			return &gen.APIV1AuthRegisterPostInternalServerError{
				Message: ErrInternalError.Error(),
				Status:  http.StatusInternalServerError,
			}, nil
		}
	}

	return &gen.RegisterResponse{
		UserID:    u.UserID.String(),
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}, nil
}

func (h *Handler) APIV1AuthLoginPost(ctx context.Context, req *gen.LoginRequest) (gen.APIV1AuthLoginPostRes, error) {
	tokens, err := h.authService.Login(ctx, string(req.Email), req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrWrongEmailOrPassword) {
			return &gen.APIV1AuthLoginPostUnauthorized{
				Message: domain.ErrWrongEmailOrPassword.Error(),
				Status:  http.StatusUnauthorized,
			}, nil
		} else {
			return &gen.APIV1AuthLoginPostInternalServerError{
				Message: ErrInternalError.Error(),
				Status:  http.StatusInternalServerError,
			}, nil
		}
	}

	cookie := h.formCookieString(tokens.RefreshToken, tokens.RefreshTokenExpiresAt)

	resp := &gen.AccessTokenHeaders{
		SetCookie: gen.NewOptString(cookie),
		Response: gen.AccessToken{
			AccessToken: tokens.AccessToken,
		},
	}

	return resp, nil
}

func (h *Handler) APIV1AuthMeGet(ctx context.Context) (gen.APIV1AuthMeGetRes, error) {
	id, err := getUserID(ctx)
	if err != nil {
		return &gen.APIV1AuthMeGetUnauthorized{
			Message: ErrAccessDenied.Error(),
			Status:  http.StatusUnauthorized,
		}, nil
	}

	u, err := h.authService.UserInfo(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrWrongUserID) {
			return &gen.APIV1AuthMeGetUnauthorized{
				Message: ErrAccessDenied.Error(),
				Status:  http.StatusUnauthorized,
			}, nil
		} else {
			return &gen.APIV1AuthMeGetInternalServerError{
				Message: ErrInternalError.Error(),
				Status:  http.StatusInternalServerError,
			}, nil
		}
	}

	return &gen.UserInfoResponse{
		UserID:    u.UserID.String(),
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}, nil
}

func (h *Handler) APIV1AuthLogoutPost(ctx context.Context, params gen.APIV1AuthLogoutPostParams) (gen.APIV1AuthLogoutPostRes, error) {
	token := params.RefreshToken
	if token == "" {
		return &gen.APIV1AuthLogoutPostUnauthorized{
			Message: ErrEmptyRefreshToken.Error(),
			Status:  http.StatusUnauthorized,
		}, nil
	}

	if err := h.authService.Logout(ctx, token); err != nil {
		if errors.Is(err, domain.ErrEmptyRefreshToken) {
			return &gen.APIV1AuthLogoutPostUnauthorized{
				Message: ErrEmptyRefreshToken.Error(),
				Status:  http.StatusUnauthorized,
			}, nil
		} else {
			return &gen.APIV1AuthLogoutPostInternalServerError{
				Message: ErrInternalError.Error(),
				Status:  http.StatusInternalServerError,
			}, nil
		}
	}

	clearCookie := h.clearRefreshTokenCookie()

	return &gen.APIV1AuthLogoutPostNoContent{
		SetCookie: gen.NewOptString(clearCookie),
	}, nil
}

func (h *Handler) APIV1AuthRefreshPost(ctx context.Context, params gen.APIV1AuthRefreshPostParams) (gen.APIV1AuthRefreshPostRes, error) {
	token := params.RefreshToken
	if token == "" {
		return &gen.APIV1AuthRefreshPostUnauthorized{
			Message: ErrEmptyRefreshToken.Error(),
			Status:  http.StatusUnauthorized,
		}, nil
	}

	t, err := h.authService.RefreshTokens(ctx, token)
	if err != nil {
		if errors.Is(err, domain.ErrEmptyRefreshToken) {
			return &gen.APIV1AuthRefreshPostUnauthorized{
				Message: ErrEmptyRefreshToken.Error(),
				Status:  http.StatusUnauthorized,
			}, nil
		} else if errors.Is(err, domain.ErrInvalidOrExpiredRefreshToken) {
			return &gen.APIV1AuthRefreshPostUnauthorized{
				Message: ErrInvalidOrExpiredRefreshToken.Error(),
				Status:  http.StatusUnauthorized,
			}, nil
		} else {
			return &gen.APIV1AuthRefreshPostInternalServerError{
				Message: ErrInternalError.Error(),
				Status:  http.StatusInternalServerError,
			}, nil
		}
	}

	cookie := h.formCookieString(t.RefreshToken, t.RefreshTokenExpiresAt)

	resp := &gen.AccessTokenHeaders{
		SetCookie: gen.NewOptString(cookie),
		Response: gen.AccessToken{
			AccessToken: t.AccessToken,
		},
	}

	return resp, nil
}

func (h *Handler) formCookieString(token string, expiresAt time.Time) string {
	c := &http.Cookie{
		Name:     string(CtxKeyRefreshToken),
		Value:    token,
		Path:     "/api/v1/auth/",
		Expires:  expiresAt,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
		Secure:   h.cookieSecure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	return c.String()
}

func (h *Handler) clearRefreshTokenCookie() string {
	c := &http.Cookie{
		Name:     string(CtxKeyRefreshToken),
		Value:    "",
		Path:     "/api/v1/auth/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   h.cookieSecure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	return c.String()
}

func getUserID(ctx context.Context) (uuid.UUID, error) {
	v := ctx.Value(CtxKeyUserID)
	idStr, ok := v.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("failed to convert user id: %v", v)
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse user_id: %w", err)
	}

	return id, nil
}

package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/cors"
	"github.com/vo1dFl0w/auth-service/internal/app/usecase"
	"github.com/vo1dFl0w/auth-service/internal/config"
	"github.com/vo1dFl0w/auth-service/internal/gen"
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
		errHttp := MapError(err)
		h.LogHTTPError(err, errHttp)
		return errHttp.ToRegisterErrResp(), nil
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
		errHttp := MapError(err)
		h.LogHTTPError(err, errHttp)
		return errHttp.ToLoginErrResp(), nil
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
		errHttp := MapError(err)
		h.LogHTTPError(err, errHttp)
		return errHttp.ToMeErrResp(), nil
	}

	u, err := h.authService.UserInfo(ctx, id)
	if err != nil {
		errHttp := MapError(err)
		h.LogHTTPError(err, errHttp)
		return errHttp.ToMeErrResp(), nil
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
		errHttp := MapError(ErrEmptyRefreshToken)
		h.LogHTTPError(ErrEmptyRefreshToken, errHttp)
		return errHttp.ToLogoutErrResp(), nil
	}

	if err := h.authService.Logout(ctx, token); err != nil {
		errHttp := MapError(err)
		h.LogHTTPError(err, errHttp)
		return errHttp.ToLogoutErrResp(), nil
	}

	clearCookie := h.clearRefreshTokenCookie()

	return &gen.APIV1AuthLogoutPostNoContent{
		SetCookie: gen.NewOptString(clearCookie),
	}, nil
}

func (h *Handler) APIV1AuthRefreshPost(ctx context.Context, params gen.APIV1AuthRefreshPostParams) (gen.APIV1AuthRefreshPostRes, error) {
	token := params.RefreshToken
	if token == "" {
		errHttp := MapError(ErrEmptyRefreshToken)
		h.LogHTTPError(ErrEmptyRefreshToken, errHttp)
		return errHttp.ToRefreshErrResp(), nil
	}

	t, err := h.authService.RefreshTokens(ctx, token)
	if err != nil {
		errHttp := MapError(err)
		h.LogHTTPError(err, errHttp)
		return errHttp.ToRefreshErrResp(), nil
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

func (h *Handler) LogHTTPError(err error, httpErr *HTTPError) {
	attrs := []any{
		"error", err.Error(),
		"status", httpErr.Status,
		"message", httpErr.Message,
	}

	switch {
	case httpErr.Status >= 500:
		switch httpErr.Status {
		case http.StatusGatewayTimeout:
			h.log.Error("http_request_failed", append(attrs, "reason", "dependency_timeout")...)
		default:
			h.log.Error("http_request_failed", append(attrs, "reason", "internal_server_error")...)
		}
	case httpErr.Status >= 400:
		h.log.Warn("http_request_rejected", append(attrs, "reason", "client_error")...)
	}
}

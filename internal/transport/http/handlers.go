package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/cors"
	"github.com/vo1dFl0w/auth-service/internal/config"
	"github.com/vo1dFl0w/auth-service/internal/transport/http/httpgen"
	"github.com/vo1dFl0w/auth-service/internal/usecase"
	"github.com/vo1dFl0w/auth-service/pkg/logger"
)

type Handler struct {
	cfg          *config.Config
	logger       logger.Logger
	cors         *cors.Cors
	authService  usecase.AuthService
	oauthService usecase.OAuthService
}

func NewHandler(cfg *config.Config, logger logger.Logger, authService usecase.AuthService, oauthService usecase.OAuthService) *Handler {
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
		logger:       logger,
		cors:         c,
		authService:  authService,
		oauthService: oauthService,
	}
}

func (h *Handler) APIV1AuthRegisterPost(ctx context.Context, req *httpgen.RegisterRequest) (httpgen.APIV1AuthRegisterPostRes, error) {
	u, err := h.authService.Register(ctx, string(req.Email), req.Password)
	if err != nil {
		errHttp := MapError(err)
		h.LogHTTPError(ctx, err, errHttp)
		return errHttp.ToRegisterErrResp(), nil
	}

	return &httpgen.RegisterResponse{
		UserID:    u.UserID.String(),
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}, nil
}

func (h *Handler) APIV1AuthLoginPost(ctx context.Context, req *httpgen.LoginRequest) (httpgen.APIV1AuthLoginPostRes, error) {
	tokens, err := h.authService.Login(ctx, string(req.Email), req.Password)
	if err != nil {
		errHttp := MapError(err)
		h.LogHTTPError(ctx, err, errHttp)
		return errHttp.ToLoginErrResp(), nil
	}

	cookie := h.formCookieString(tokens.RefreshToken, tokens.RefreshTokenExpiresAt)

	resp := &httpgen.AccessTokenHeaders{
		SetCookie: httpgen.NewOptString(cookie),
		Response: httpgen.AccessToken{
			AccessToken: tokens.AccessToken,
		},
	}

	return resp, nil
}

func (h *Handler) APIV1AuthMeGet(ctx context.Context) (httpgen.APIV1AuthMeGetRes, error) {
	id, err := getUserID(ctx)
	if err != nil {
		errHttp := MapError(err)
		h.LogHTTPError(ctx, err, errHttp)
		return errHttp.ToMeErrResp(), nil
	}

	u, err := h.authService.UserInfo(ctx, id)
	if err != nil {
		errHttp := MapError(err)
		h.LogHTTPError(ctx, err, errHttp)
		return errHttp.ToMeErrResp(), nil
	}

	return &httpgen.UserInfoResponse{
		UserID:    u.UserID.String(),
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}, nil
}

func (h *Handler) APIV1AuthLogoutPost(ctx context.Context, params httpgen.APIV1AuthLogoutPostParams) (httpgen.APIV1AuthLogoutPostRes, error) {
	token := params.RefreshToken
	if token == "" {
		errHttp := MapError(ErrEmptyRefreshToken)
		h.LogHTTPError(ctx, ErrEmptyRefreshToken, errHttp)
		return errHttp.ToLogoutErrResp(), nil
	}

	if err := h.authService.Logout(ctx, token); err != nil {
		errHttp := MapError(err)
		h.LogHTTPError(ctx, err, errHttp)
		return errHttp.ToLogoutErrResp(), nil
	}

	clearCookie := h.clearRefreshTokenCookie()

	return &httpgen.APIV1AuthLogoutPostNoContent{
		SetCookie: httpgen.NewOptString(clearCookie),
	}, nil
}

func (h *Handler) APIV1AuthRefreshPost(ctx context.Context, params httpgen.APIV1AuthRefreshPostParams) (httpgen.APIV1AuthRefreshPostRes, error) {
	token := params.RefreshToken
	if token == "" {
		errHttp := MapError(ErrEmptyRefreshToken)
		h.LogHTTPError(ctx, ErrEmptyRefreshToken, errHttp)
		return errHttp.ToRefreshErrResp(), nil
	}

	t, err := h.authService.RefreshTokens(ctx, token)
	if err != nil {
		errHttp := MapError(err)
		h.LogHTTPError(ctx, err, errHttp)
		return errHttp.ToRefreshErrResp(), nil
	}

	cookie := h.formCookieString(t.RefreshToken, t.RefreshTokenExpiresAt)

	resp := &httpgen.AccessTokenHeaders{
		SetCookie: httpgen.NewOptString(cookie),
		Response: httpgen.AccessToken{
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
		Secure:   h.cfg.Cookie.CookieSecure,
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
		Secure:   h.cfg.Cookie.CookieSecure,
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

func (h *Handler) LogHTTPError(ctx context.Context, err error, httpErr *HTTPError) {
	attrs := []any{
		"request_id", ctx.Value(CtxKeyRequestID),
		"error", err,
		"status", httpErr.Status,
		"message", httpErr.Message,
	}

	switch {
	case httpErr.Status >= 500:
		switch httpErr.Status {
		case http.StatusGatewayTimeout:
			h.logger.Error("http_request_failed", append(attrs, "reason", "dependency_timeout")...)
		default:
			h.logger.Error("http_request_failed", append(attrs, "reason", "internal_server_error")...)
		}
	case httpErr.Status >= 400:
		h.logger.Warn("http_request_rejected", append(attrs, "reason", "client_error")...)
	}
}

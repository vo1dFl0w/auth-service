package http

import (
	"context"
	"net/http"
	"time"

	"github.com/vo1dFl0w/auth-service/internal/transport/http/httpgen"
)

func (h *Handler) AuthGoogleLogin(ctx context.Context) (httpgen.AuthGoogleLoginRes, error) {
	state, err := h.oauthService.GetGenerateState(ctx)
	if err != nil {
		errHttp := MapError(err)
		h.LogHTTPError(ctx, err, errHttp)
		return errHttp.ToAuthGoogleLoginErrResp(), nil
	}

	authURL := h.oauthService.GetAuthCodeURL(ctx, state)
	cookie := h.formOAuthStateCookie(state)


	resp := &httpgen.AuthGoogleLoginFound{
		SetCookie: httpgen.NewOptString(cookie),
		Location:  httpgen.NewOptString(authURL),
	}

	return resp, nil
}

func (h *Handler) AuthGoogleCallback(ctx context.Context, params httpgen.AuthGoogleCallbackParams) (httpgen.AuthGoogleCallbackRes, error) {
	u, err := h.oauthService.GetUserFromCode(ctx, params.Code, params.State)
	if err != nil {
		errHttp := MapError(err)
		h.LogHTTPError(ctx, err, errHttp)
		return errHttp.ToAuthGoogleCallbackErrResp(), nil
	}

	res, err := h.authService.LoginWithOAuthProvider(ctx, u.Email, u.ProviderID, u.Provider)
	if err != nil {
		errHttp := MapError(err)
		h.LogHTTPError(ctx, err, errHttp)
		return errHttp.ToAuthGoogleCallbackErrResp(), nil
	}

	cookie := h.formCookieString(res.RefreshToken, res.RefreshTokenExpiresAt)

	resp := &httpgen.AccessTokenHeaders{
		SetCookie: httpgen.NewOptString(cookie),
		Response: httpgen.AccessToken{
			AccessToken: res.AccessToken,
		},
	}

	return resp, nil
}

func (h *Handler) formOAuthStateCookie(state string) string {
	c := &http.Cookie{
		Name:     string(CtxKeyOAuthState),
		Value:    state,
		Path:     "/api/v1/auth/",
		Expires:  time.Now().Add(time.Second * time.Duration(h.cfg.Cookie.MaxAge)),
		MaxAge:   h.cfg.Cookie.MaxAge,
		Secure:   h.cfg.Cookie.CookieSecure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	return c.String()
}

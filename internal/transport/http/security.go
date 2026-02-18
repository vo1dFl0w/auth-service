package http

import (
    "context"
    "fmt"

    "github.com/vo1dFl0w/auth-service/internal/transport/http/httpgen"
    "github.com/vo1dFl0w/auth-service/internal/usecase"
)

type SecuredHandler struct {
    tokenService usecase.TokenService
}

func NewSecuredHandler(tokenService usecase.TokenService) *SecuredHandler {
    return &SecuredHandler{tokenService: tokenService}
}

func (h *SecuredHandler) HandleBearerAuth(ctx context.Context, operationName httpgen.OperationName, t httpgen.BearerAuth) (context.Context, error) {
    if t.Token == "" {
        return ctx, fmt.Errorf("missing bearer token")
    }

    claims, err := h.tokenService.ValidateAccessToken(t.Token)
    if err != nil {
        return ctx, err
    }

    ctx = context.WithValue(ctx, CtxKeyUserID, claims.Subject)
    return ctx, nil
}
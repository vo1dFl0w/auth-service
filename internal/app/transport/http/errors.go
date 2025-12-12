package http

import "errors"

var (
	ErrAccessDenied                 = errors.New("access denied")
	ErrBadRequest                   = errors.New("bad request")
	ErrEmptyRefreshToken            = errors.New("empty refresh token")
	ErrGatewayTimeout               = errors.New("gateway timeout")
	ErrInternalError                = errors.New("internal error")
	ErrInvalidAuthorizationHeader   = errors.New("invalid authorization header")
	ErrInvalidOrExpiredRefreshToken = errors.New("invalid or expired refresh token")
)

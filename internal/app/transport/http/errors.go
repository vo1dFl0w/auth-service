package http

import (
	"errors"
	"net/http"

	"github.com/vo1dFl0w/auth-service/internal/app/domain"
	"github.com/vo1dFl0w/auth-service/internal/gen"
)

var (
	ErrAccessDenied                 = errors.New("access denied")
	ErrBadRequest                   = errors.New("bad request")
	ErrEmptyRefreshToken            = errors.New("empty refresh token")
	ErrGatewayTimeout               = errors.New("gateway timeout")
	ErrInternalError                = errors.New("internal error")
	ErrInvalidAuthorizationHeader   = errors.New("invalid authorization header")
	ErrInvalidOrExpiredRefreshToken = errors.New("invalid or expired refresh token")
)

type HTTPError struct {
	Message string
	Status  int
}

func (e *HTTPError) Error() string {
	return e.Message
}

func MapError(err error) *HTTPError {
	switch {
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		return &HTTPError{
			Message: domain.ErrEmailAlreadyExists.Error(),
			Status:  http.StatusConflict,
		}
	case errors.Is(err, domain.ErrInvalidEmail) || errors.Is(err, domain.ErrInvalidPassword):
		return &HTTPError{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
	case errors.Is(err, domain.ErrWrongEmailOrPassword):
		return &HTTPError{
			Message: domain.ErrWrongEmailOrPassword.Error(),
			Status:  http.StatusUnauthorized,
		}
	case errors.Is(err, domain.ErrGatewayTimeout):
		return &HTTPError{
			Message: ErrGatewayTimeout.Error(),
			Status:  http.StatusGatewayTimeout,
		}
	case errors.Is(err, domain.ErrWrongUserID):
		return &HTTPError{
			Message: ErrAccessDenied.Error(),
			Status:  http.StatusUnauthorized,
		}
	case errors.Is(err, domain.ErrEmptyRefreshToken):
		return &HTTPError{
			Message: ErrEmptyRefreshToken.Error(),
			Status:  http.StatusUnauthorized,
		}
	case errors.Is(err, domain.ErrInvalidOrExpiredRefreshToken):
		return &HTTPError{
			Message: ErrInvalidOrExpiredRefreshToken.Error(),
			Status:  http.StatusUnauthorized,
		}
	default:
		return &HTTPError{
			Message: ErrInternalError.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
}

func (e *HTTPError) ToRegisterErrResp() gen.APIV1AuthRegisterPostRes {
	switch e.Status {
	case http.StatusConflict:
		return &gen.APIV1AuthRegisterPostConflict{
			Message: e.Message,
			Status:  e.Status,
		}
	case http.StatusBadRequest:
		return &gen.APIV1AuthRegisterPostBadRequest{
			Message: e.Message,
			Status:  e.Status,
		}
	case http.StatusGatewayTimeout:
		return &gen.APIV1AuthRegisterPostGatewayTimeout{
			Message: e.Message,
			Status:  e.Status,
		}
	default:
		return &gen.APIV1AuthRegisterPostInternalServerError{
			Message: e.Message,
			Status:  e.Status,
		}
	}
}

func (e *HTTPError) ToLoginErrResp() gen.APIV1AuthLoginPostRes {
	switch e.Status {
	case http.StatusUnauthorized:
		return &gen.APIV1AuthLoginPostUnauthorized{
			Message: e.Message,
			Status:  e.Status,
		}
	case http.StatusGatewayTimeout:
		return &gen.APIV1AuthLoginPostGatewayTimeout{
			Message: e.Message,
			Status:  e.Status,
		}
	default:
		return &gen.APIV1AuthLoginPostInternalServerError{
			Message: e.Message,
			Status:  e.Status,
		}
	}
}

func (e *HTTPError) ToMeErrResp() gen.APIV1AuthMeGetRes {
	switch e.Status {
	case http.StatusUnauthorized:
		return &gen.APIV1AuthMeGetUnauthorized{
			Message: e.Message,
			Status:  e.Status,
		}
	case http.StatusGatewayTimeout:
		return &gen.APIV1AuthMeGetGatewayTimeout{
			Message: e.Message,
			Status:  e.Status,
		}
	default:
		return &gen.APIV1AuthMeGetInternalServerError{
			Message: e.Message,
			Status:  e.Status,
		}
	}
}

func (e *HTTPError) ToLogoutErrResp() gen.APIV1AuthLogoutPostRes {
	switch e.Status {
	case http.StatusUnauthorized:
		return &gen.APIV1AuthLogoutPostUnauthorized{
			Message: e.Message,
			Status:  e.Status,
		}
	case http.StatusGatewayTimeout:
		return &gen.APIV1AuthLogoutPostGatewayTimeout{
			Message: e.Message,
			Status:  e.Status,
		}
	default:
		return &gen.APIV1AuthLogoutPostInternalServerError{
			Message: e.Message,
			Status:  e.Status,
		}
	}
}

func (e *HTTPError) ToRefreshErrResp() gen.APIV1AuthRefreshPostRes {
	switch e.Status {
	case http.StatusUnauthorized:
		return &gen.APIV1AuthRefreshPostUnauthorized{
			Message: e.Message,
			Status:  e.Status,
		}
	case http.StatusGatewayTimeout:
		return &gen.APIV1AuthRefreshPostGatewayTimeout{
			Message: e.Message,
			Status:  e.Status,
		}
	default:
		return &gen.APIV1AuthRefreshPostInternalServerError{
			Message: e.Message,
			Status:  e.Status,
		}
	}
}

package http

import (
	"errors"
	"net/http"

	"github.com/vo1dFl0w/auth-service/internal/domain"
	httpgen "github.com/vo1dFl0w/auth-service/internal/transport/http/httpgen"
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

func (e *HTTPError) ToRegisterErrResp() httpgen.APIV1AuthRegisterPostRes {
	switch e.Status {
	case http.StatusConflict:
		return &httpgen.APIV1AuthRegisterPostConflict{
			Message: e.Message,
			Status:  e.Status,
		}
	case http.StatusBadRequest:
		return &httpgen.APIV1AuthRegisterPostBadRequest{
			Message: e.Message,
			Status:  e.Status,
		}
	case http.StatusGatewayTimeout:
		return &httpgen.APIV1AuthRegisterPostGatewayTimeout{
			Message: e.Message,
			Status:  e.Status,
		}
	default:
		return &httpgen.APIV1AuthRegisterPostInternalServerError{
			Message: e.Message,
			Status:  e.Status,
		}
	}
}

func (e *HTTPError) ToLoginErrResp() httpgen.APIV1AuthLoginPostRes {
	switch e.Status {
	case http.StatusUnauthorized:
		return &httpgen.APIV1AuthLoginPostUnauthorized{
			Message: e.Message,
			Status:  e.Status,
		}
	case http.StatusGatewayTimeout:
		return &httpgen.APIV1AuthLoginPostGatewayTimeout{
			Message: e.Message,
			Status:  e.Status,
		}
	default:
		return &httpgen.APIV1AuthLoginPostInternalServerError{
			Message: e.Message,
			Status:  e.Status,
		}
	}
}

func (e *HTTPError) ToMeErrResp() httpgen.APIV1AuthMeGetRes {
	switch e.Status {
	case http.StatusUnauthorized:
		return &httpgen.APIV1AuthMeGetUnauthorized{
			Message: e.Message,
			Status:  e.Status,
		}
	case http.StatusGatewayTimeout:
		return &httpgen.APIV1AuthMeGetGatewayTimeout{
			Message: e.Message,
			Status:  e.Status,
		}
	default:
		return &httpgen.APIV1AuthMeGetInternalServerError{
			Message: e.Message,
			Status:  e.Status,
		}
	}
}

func (e *HTTPError) ToLogoutErrResp() httpgen.APIV1AuthLogoutPostRes {
	switch e.Status {
	case http.StatusUnauthorized:
		return &httpgen.APIV1AuthLogoutPostUnauthorized{
			Message: e.Message,
			Status:  e.Status,
		}
	case http.StatusGatewayTimeout:
		return &httpgen.APIV1AuthLogoutPostGatewayTimeout{
			Message: e.Message,
			Status:  e.Status,
		}
	default:
		return &httpgen.APIV1AuthLogoutPostInternalServerError{
			Message: e.Message,
			Status:  e.Status,
		}
	}
}

func (e *HTTPError) ToRefreshErrResp() httpgen.APIV1AuthRefreshPostRes {
	switch e.Status {
	case http.StatusUnauthorized:
		return &httpgen.APIV1AuthRefreshPostUnauthorized{
			Message: e.Message,
			Status:  e.Status,
		}
	case http.StatusGatewayTimeout:
		return &httpgen.APIV1AuthRefreshPostGatewayTimeout{
			Message: e.Message,
			Status:  e.Status,
		}
	default:
		return &httpgen.APIV1AuthRefreshPostInternalServerError{
			Message: e.Message,
			Status:  e.Status,
		}
	}
}

func (e *HTTPError) ToAuthGoogleLoginErrResp() httpgen.AuthGoogleLoginRes {
	switch e.Status {
	case http.StatusGatewayTimeout:
		return &httpgen.AuthGoogleLoginGatewayTimeout{
			Message: e.Message,
			Status:  e.Status,
		}
	default:
		return &httpgen.AuthGoogleLoginInternalServerError{
			Message: e.Message,
			Status:  e.Status,
		}
	}
}

func (e *HTTPError) ToAuthGoogleCallbackErrResp() httpgen.AuthGoogleCallbackRes {
	switch e.Status {
	case http.StatusGatewayTimeout:
		return &httpgen.AuthGoogleCallbackGatewayTimeout{
			Message: e.Message,
			Status:  e.Status,
		}
	default:
		return &httpgen.AuthGoogleCallbackInternalServerError{
			Message: e.Message,
			Status:  e.Status,
		}
	}
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

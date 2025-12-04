package domain

import "errors"

var (
	ErrEmailAlreadyExists           = errors.New("email already exists")
	ErrEmptyPassword                = errors.New("empty password")
	ErrEmptyRefreshToken            = errors.New("empty refresh token")
	ErrExpiredAccessToken           = errors.New("expired access token")
	ErrInvalidPassword              = errors.New("invalid password")
	ErrInvalidEmail                 = errors.New("invalid email")
	ErrInvalidAccessToken           = errors.New("invalid access token")
	ErrInvalidOrExpiredRefreshToken = errors.New("invalid or expired refresh token")
	ErrNotFound                     = errors.New("not found")
	ErrWrongEmailOrPassword         = errors.New("wrong email or password")
	ErrWrongUserID                  = errors.New("wrong user id")
)

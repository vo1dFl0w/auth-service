package repository

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("already exists")
	ErrGatewayTimeout            = errors.New("gateway timeout")
	ErrNotFound           = errors.New("not found")
	ErrNoRowDeleted       = errors.New("no row deleted")
)

package repository

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrNoRowDeleted = errors.New("no row deleted")
	ErrEmailAlreadyExists = errors.New("already exists")
)
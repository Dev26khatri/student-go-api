package utils

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("user email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

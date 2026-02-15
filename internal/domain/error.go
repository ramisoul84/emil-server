package domain

import "errors"

// JWT Errors
var (
	ErrTokenGenerate = errors.New("failed to generate token")
	ErrTokenInvalid  = errors.New("invalid token")
	ErrTokenExpired  = errors.New("expired token")
)

// DB Errors
var (
	ErrNotFound = errors.New("not found")
	ErrInternal = errors.New("internal error")
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

package domain

import "errors"

var (
	ErrInternalServer     = errors.New("internal server error")
	ErrNotFound           = errors.New("not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrExpiredAccessToken = errors.New("access token expired")
	ErrBadRequest         = errors.New("bad request")
)

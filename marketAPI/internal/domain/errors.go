package domain

import (
	"errors"
)

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUserExists           = errors.New("user already exists")
	ErrInvalidUsername      = errors.New("username must be 3-50 characters long")
	ErrInvalidPassword      = errors.New("password must be at least 8 characters long")
	ErrInvalidToken         = errors.New("invalid token")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrInvalidAdTitle       = errors.New("ad title must be 5-72 characters long")
	ErrInvalidAdDescription = errors.New("ad description must be 10-1000 characters long")
	ErrInvalidAdPrice       = errors.New("price must be positive")
	ErrInvalidImageURL      = errors.New("invalid image URL")
)

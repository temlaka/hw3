package domain

import "errors"

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidInput       = errors.New("invalid input")
	ErrUserNotFound       = errors.New("user not found")
	ErrCryptoExists       = errors.New("crypto tracked already")
	ErrUnknownSymbol      = errors.New("unknown symbol")
	ErrCryptoNotFound     = errors.New("crypto not found")
)

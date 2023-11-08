package services

import (
	"errors"
)

var (
	ErrInvalidOrderNumberFormat = errors.New("invalid order number format")
	ErrInvalidOrderNumber       = errors.New("invalid order number")
	ErrOrderAlreadyExists       = errors.New("order already exists")

	ErrInvalidToken             = errors.New("invalid token")
	ErrGenerateToken            = errors.New("failed to generate token")
	ErrGenerateHashFromPassword = errors.New("failed to generate hash from password")
	ErrIncorrectPassword        = errors.New("incorrect password")

	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

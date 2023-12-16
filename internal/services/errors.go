package services

import (
	"errors"
)

var (
	ErrInvalidOrderNumberFormat  = errors.New("invalid order number format")
	ErrInvalidOrderNumber        = errors.New("invalid order number")
	ErrOrderAlreadyExists        = errors.New("order already exists")
	ErrOrderAlreadyExistsForUser = errors.New("order already exists for this user")

	ErrGenerateToken            = errors.New("failed to generate token")
	ErrGenerateHashFromPassword = errors.New("failed to generate hash from password")
	ErrIncorrectPassword        = errors.New("incorrect password")

	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidLoginFormat = errors.New("invalid login format")
	ErrInsufficientFunds  = errors.New("insufficient funds")
)

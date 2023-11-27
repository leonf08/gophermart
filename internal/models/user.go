package models

import "github.com/golang-jwt/jwt/v5"

type (
	User struct {
		UserID   string `json:"-" db:"user_id"`
		Login    string `json:"login" db:"login"`
		Password string `json:"password,omitempty" db:"password"`
	}

	UserAccount struct {
		UserID    string  `json:"-" db:"user_id"`
		Current   float64 `json:"current" db:"current"`
		Withdrawn float64 `json:"withdrawn" db:"withdrawn"`
	}

	CustomJWTClaims struct {
		jwt.RegisteredClaims
		UserID string `json:"user_id"`
	}
)

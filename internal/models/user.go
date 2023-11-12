package models

import "github.com/golang-jwt/jwt/v5"

type (
	User struct {
		UserID   int64  `json:"user_id,omitempty"`
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	UserAccount struct {
		UserID    int64 `json:"-"`
		Current   int64 `json:"current"`
		Withdrawn int64 `json:"withdrawn"`
	}

	CustomJWTClaims struct {
		jwt.RegisteredClaims
		Login string `json:"login"`
	}
)

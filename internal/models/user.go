package models

import "github.com/golang-jwt/jwt/v5"

type (
	// swagger:model
	User struct {
		// the id for the user
		//
		// required: false
		// min length: 16
		// max length: 16
		UserID string `json:"user_id,omitempty" db:"user_id"`

		// the login for the user
		//
		// required: true
		// min length: 3
		// max length: 16
		Login string `json:"login" db:"login"`

		// the password for the user
		//
		// required: true
		// min length: 6
		Password string `json:"password" db:"password"`
	}

	// swagger:model
	UserAccount struct {
		// the id for the user
		//
		// required: false
		// min length: 16
		// max length: 16
		UserID string `json:"-" db:"user_id"`

		// the current balance for the user
		//
		// required: true
		// min: 0
		Current int64 `json:"current" db:"current"`

		// the withdrawn amount for the user
		//
		// required: true
		// min: 0
		Withdrawn int64 `json:"withdrawn" db:"withdrawn"`
	}

	CustomJWTClaims struct {
		jwt.RegisteredClaims
		UserID string `json:"user-id"`
	}
)

package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/leonf08/gophermart.git/internal/models"
	errs "github.com/leonf08/gophermart.git/internal/services"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Authenticator struct {
	key string
}

func NewAuthenticator(key string) *Authenticator {
	return &Authenticator{
		key: key,
	}
}

// GenerateToken generates a JWT token for a given user
// and returns it as a string.
// The token is signed with HMAC-SHA256 algorithm.
// The token contains the following claims:
// - issuer: "gophermart"
// - expiration time: 24 hours
// - issued at: current time
// - login: user login
// The token is signed with the key provided to the Authenticator.
// If the token generation fails, an error is returned.
// If the token generation succeeds, the token is returned.
func (a *Authenticator) GenerateToken(user *models.User) (string, error) {
	claims := &models.CustomJWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "gophermart",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Login: user.Login,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(a.key))
	if err != nil {
		return "", errs.ErrGenerateToken
	}

	return signedToken, nil
}

// ValidateTokenAndExtractClaims validates a JWT token
// and returns the claims if the token is valid.
// The token is signed with HMAC-SHA256 algorithm.
// If the token validation fails, an error is returned.
// If the token validation succeeds, the claims are returned.
func (a *Authenticator) ValidateTokenAndExtractClaims(token string) (*models.CustomJWTClaims, error) {
	claims := &models.CustomJWTClaims{}

	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.key), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	if err != nil {
		return nil, errs.ErrInvalidToken
	}

	return claims, nil
}

func (a *Authenticator) GenerateHashFromPassword(user *models.User) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", errs.ErrGenerateHashFromPassword
	}

	return string(hashedPassword), nil
}

func (a *Authenticator) CheckPasswordHash(user, storedUser *models.User) error {
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		return errs.ErrIncorrectPassword
	}

	return nil
}

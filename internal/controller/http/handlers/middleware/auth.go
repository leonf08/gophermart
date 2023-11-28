package middleware

import (
	"context"
	"github.com/leonf08/gophermart.git/internal/services"
	"net/http"
	"strings"
)

type KeyUserID struct{}

func Auth(auth services.Authenticator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := strings.Split(r.Header.Get("Authorization"), " ")
			if len(token) != 2 {
				http.Error(w, "invalid token format", http.StatusUnauthorized)
				return
			}

			claims, err := auth.ValidateTokenAndExtractClaims(token[1])
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), KeyUserID{}, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

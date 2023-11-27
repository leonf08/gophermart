package handlers

import (
	"compress/flate"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/leonf08/gophermart.git/internal/controller/http/handlers/middleware"
	"github.com/leonf08/gophermart.git/internal/services"
	"log/slog"
)

func NewRouter(users services.Users, orders services.Orders, auth services.Authenticator, log *slog.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"POST", "GET"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
			AllowCredentials: true,
			ExposedHeaders:   []string{"Authorization"},
			MaxAge:           300,
		}),
		chiMiddleware.Compress(flate.BestCompression),
		chiMiddleware.RequestID,
		middleware.Logging(log),
	)

	r.Route("/api/user/", func(r chi.Router) {
		newHandler(r, users, orders, auth, log)
	})

	return r
}

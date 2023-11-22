package handlers

import (
	"compress/flate"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/leonf08/gophermart.git/internal/controller/http/handlers/middleware"
	"github.com/leonf08/gophermart.git/internal/services"
	"log/slog"
)

func NewRouter(users services.Users, orders services.Orders, auth services.Authenticator, log *slog.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(
		chiMiddleware.Compress(flate.BestCompression),
		chiMiddleware.RequestID,
		middleware.Logging(log),
	)

	r.Route("/api/user/", func(r chi.Router) {
		newHandler(r, users, orders, auth, log)
	})

	return r
}

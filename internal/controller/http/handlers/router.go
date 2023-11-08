package handlers

import (
	"compress/flate"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/leonf08/gophermart.git/internal/controller/http/handlers/middleware"
	"log/slog"
)

func NewRouter(logger *slog.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(
		chiMiddleware.Compress(flate.BestCompression),
		chiMiddleware.RequestID,
		middleware.Logging(logger),
	)

	r.Route("/api/user/", func(r chi.Router) {
		r.Post("/register/", userSignUp)
		r.Post("/login/", userLogIn)
		r.Post("/orders/", uploadOrder)
		r.Get("/orders/", getOrders)
		r.Get("/balance/", getUserBalance)
		r.Post("/balance/withdraw/", withdraw)
		r.Get("/withdrawals/", getWithdrawals)
	})

	return r
}

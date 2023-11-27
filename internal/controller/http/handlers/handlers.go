package handlers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/leonf08/gophermart.git/internal/controller/http/handlers/middleware"
	"github.com/leonf08/gophermart.git/internal/models"
	"github.com/leonf08/gophermart.git/internal/services"
	"io"
	"log/slog"
	"net/http"
	"sort"
)

type handler struct {
	users  services.Users
	orders services.Orders
	log    services.Logger
}

func newHandler(r chi.Router, users services.Users, orders services.Orders, auth services.Authenticator, log services.Logger) {
	h := &handler{
		users:  users,
		orders: orders,
		log:    log,
	}

	r.Post("/register/", h.userSignUp)
	r.Post("/login/", h.userLogIn)
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(auth))
		r.Post("/orders/", h.uploadOrder)
		r.Get("/orders/", h.getOrders)
		r.Get("/balance/", h.getUserBalance)
		r.Post("/balance/withdraw/", h.withdraw)
		r.Get("/withdrawals/", h.getWithdrawals)
	})
}

func (h *handler) userSignUp(w http.ResponseWriter, r *http.Request) {
	entry := logEntry(h.log, r)

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		entry.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.users.RegisterUser(r.Context(), user)
	if err != nil {
		entry.Error(err.Error())
		switch {
		case errors.Is(err, services.ErrUserAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	token, err := h.users.GetToken(user)
	if err != nil {
		entry.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}

func (h *handler) userLogIn(w http.ResponseWriter, r *http.Request) {
	entry := logEntry(h.log, r)

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		entry.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.users.LoginUser(r.Context(), user)
	if err != nil {
		entry.Error(err.Error())
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := h.users.GetToken(user)
	if err != nil {
		entry.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}

func (h *handler) uploadOrder(w http.ResponseWriter, r *http.Request) {
	entry := logEntry(h.log, r)

	userId := r.Context().Value("userId").(string)

	orderNum, err := io.ReadAll(r.Body)
	if err != nil {
		entry.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.orders.CreateNewOrder(r.Context(), userId, string(orderNum))
	if err != nil {
		entry.Error(err.Error())
		switch {
		case errors.Is(err, services.ErrInvalidOrderNumber):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, services.ErrInvalidOrderNumberFormat):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		case errors.Is(err, services.ErrOrderAlreadyExists):
			http.Error(w, err.Error(), http.StatusConflict)
		case errors.Is(err, services.ErrOrderAlreadyExistsForUser):
			http.Error(w, err.Error(), http.StatusOK)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusAccepted)
}

func (h *handler) getOrders(w http.ResponseWriter, r *http.Request) {
	entry := logEntry(h.log, r)

	userId := r.Context().Value("userId").(string)

	orders, err := h.orders.GetOrdersForUser(r.Context(), userId)
	if err != nil {
		entry.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	sort.SliceStable(orders, func(i, j int) bool {
		return orders[i].UploadedAt.Before(orders[j].UploadedAt)
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(orders); err != nil {
		entry.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) getUserBalance(w http.ResponseWriter, r *http.Request) {
	entry := logEntry(h.log, r)

	userId := r.Context().Value("userId").(string)

	balance, err := h.users.GetUserAccount(r.Context(), userId)
	if err != nil {
		entry.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(balance); err != nil {
		entry.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *handler) withdraw(w http.ResponseWriter, r *http.Request) {
	entry := logEntry(h.log, r)

	userId := r.Context().Value("userId").(string)

	withdrawal := &models.Withdrawal{}
	err := json.NewDecoder(r.Body).Decode(withdrawal)
	if err != nil {
		entry.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	withdrawal.UserID = userId
	err = h.users.WithdrawFromAccount(r.Context(), withdrawal)
	if err != nil {
		entry.Error(err.Error())
		switch {
		case errors.Is(err, services.ErrInsufficientFunds):
			http.Error(w, err.Error(), http.StatusPaymentRequired)
		case errors.Is(err, services.ErrInvalidOrderNumber):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, services.ErrInvalidOrderNumberFormat):
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) getWithdrawals(w http.ResponseWriter, r *http.Request) {
	entry := logEntry(h.log, r)

	userId := r.Context().Value("userId").(string)

	withdrawals, err := h.users.GetWithdrawals(r.Context(), userId)
	if err != nil {
		entry.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	sort.SliceStable(withdrawals, func(i, j int) bool {
		return withdrawals[i].ProcessedAt.Before(withdrawals[j].ProcessedAt)
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(withdrawals); err != nil {
		entry.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func logEntry(l services.Logger, r *http.Request) services.Logger {
	log, ok := l.(*slog.Logger)
	if !ok {
		return l
	}

	return log.With(
		slog.String("component", "handler"),
		slog.String("method", r.Method),
		slog.String("url", r.URL.Path),
	)
}

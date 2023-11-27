package handlers

import (
	"context"
	"errors"
	"github.com/leonf08/gophermart.git/internal/models"
	"github.com/leonf08/gophermart.git/internal/services"
	"github.com/leonf08/gophermart.git/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type mockLogger struct{}

func (m *mockLogger) Info(msg string, args ...any) {}

func (m *mockLogger) Error(msg string, args ...any) {}

func Test_handler_getOrders(t *testing.T) {
	users := mocks.NewUsers(t)
	orders := mocks.NewOrders(t)
	log := &mockLogger{}

	h := &handler{
		users:  users,
		orders: orders,
		log:    log,
	}

	type want struct {
		contentType string
		status      int
	}

	tests := []struct {
		name   string
		userId string
		want   want
	}{
		{
			name:   "1. get orders success",
			userId: "user",
			want: want{
				contentType: "application/json",
				status:      http.StatusOK,
			},
		},
		{
			name:   "2. get orders, no content",
			userId: "admin",
			want: want{
				contentType: "text/plain; charset=utf-8",
				status:      http.StatusNoContent,
			},
		},
		{
			name: "3. get orders, internal server error",
			want: want{
				contentType: "text/plain; charset=utf-8",
				status:      http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders.
				On("GetOrdersForUser", mock.Anything, mock.Anything).
				Return(func(ctx context.Context, userId string) ([]*models.Order, error) {
					if userId == "user" {
						return []*models.Order{
							{
								Number:     "123456789",
								UploadedAt: time.Now(),
							},
						}, nil
					} else if userId == "admin" {
						return []*models.Order{}, nil
					}

					return nil, errors.New("internal server error")
				})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			resp := httptest.NewRecorder()
			h.getOrders(resp, req.WithContext(context.WithValue(req.Context(), "userId", tt.userId)))

			orders.AssertExpectations(t)

			assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"))
			assert.Equal(t, tt.want.status, resp.Code)
		})
	}
}

func Test_handler_getUserBalance(t *testing.T) {
	users := mocks.NewUsers(t)
	orders := mocks.NewOrders(t)
	log := &mockLogger{}

	h := &handler{
		users:  users,
		orders: orders,
		log:    log,
	}

	type want struct {
		contentType string
		status      int
	}

	tests := []struct {
		name   string
		userId string
		want   want
	}{
		{
			name:   "1. get user balance success",
			userId: "user",
			want: want{
				contentType: "application/json",
				status:      http.StatusOK,
			},
		},
		{
			name: "2. get user balance, internal server error",
			want: want{
				contentType: "text/plain; charset=utf-8",
				status:      http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users.
				On("GetUserAccount", mock.Anything, mock.Anything).
				Return(func(ctx context.Context, userId string) (*models.UserAccount, error) {
					if userId == "user" {
						return &models.UserAccount{
							Current:   100,
							Withdrawn: 0,
						}, nil
					}

					return nil, errors.New("internal server error")
				})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			resp := httptest.NewRecorder()
			h.getUserBalance(resp, req.WithContext(context.WithValue(req.Context(), "userId", tt.userId)))

			users.AssertExpectations(t)

			assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"))
			assert.Equal(t, tt.want.status, resp.Code)
		})
	}
}

func Test_handler_getWithdrawals(t *testing.T) {
	users := mocks.NewUsers(t)
	orders := mocks.NewOrders(t)
	log := &mockLogger{}

	h := &handler{
		users:  users,
		orders: orders,
		log:    log,
	}

	type want struct {
		contentType string
		status      int
	}

	tests := []struct {
		name   string
		userId string
		want   want
	}{
		{
			name:   "1. get withdrawals success",
			userId: "user",
			want: want{
				contentType: "application/json",
				status:      http.StatusOK,
			},
		},
		{
			name:   "2. get withdrawals, no content",
			userId: "admin",
			want: want{
				contentType: "text/plain; charset=utf-8",
				status:      http.StatusNoContent,
			},
		},
		{
			name: "3. get withdrawals, internal server error",
			want: want{
				contentType: "text/plain; charset=utf-8",
				status:      http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users.
				On("GetWithdrawals", mock.Anything, mock.Anything).
				Return(func(ctx context.Context, userId string) ([]*models.Withdrawal, error) {
					if userId == "user" {
						return []*models.Withdrawal{
							{
								OrderNumber: "123456789",
								Sum:         100,
								ProcessedAt: time.Now(),
							},
						}, nil
					} else if userId == "admin" {
						return []*models.Withdrawal{}, nil
					}

					return nil, errors.New("internal server error")
				})

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			resp := httptest.NewRecorder()
			h.getWithdrawals(resp, req.WithContext(context.WithValue(req.Context(), "userId", tt.userId)))

			users.AssertExpectations(t)

			assert.Equal(t, tt.want.contentType, resp.Header().Get("Content-Type"))
			assert.Equal(t, tt.want.status, resp.Code)
		})
	}
}

func Test_handler_uploadOrder(t *testing.T) {
	users := mocks.NewUsers(t)
	orders := mocks.NewOrders(t)
	log := &mockLogger{}

	h := &handler{
		users:  users,
		orders: orders,
		log:    log,
	}

	type want struct {
		status int
	}

	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "1. upload order success",
			body: `123456`,
			want: want{
				status: http.StatusAccepted,
			},
		},
		{
			name: "2. upload order fail, invalid order number",
			body: `12345hfjfh`,
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "3. upload order fail, invalid order number format",
			body: `1`,
			want: want{
				status: http.StatusUnprocessableEntity,
			},
		},
		{
			name: "4. upload order fail, order already exists",
			body: "1234567",
			want: want{
				status: http.StatusConflict,
			},
		},
		{
			name: "5. upload order fail, order already exists for user",
			body: "12345678",
			want: want{
				status: http.StatusOK,
			},
		},
		{
			name: "6. upload order fail, internal server error",
			body: "123456789",
			want: want{
				status: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orders.
				On("CreateNewOrder", mock.Anything, mock.Anything, mock.Anything).
				Return(func(ctx context.Context, userId, orderNum string) error {
					if orderNum == "1234567" {
						return services.ErrOrderAlreadyExists
					} else if orderNum == "12345678" {
						return services.ErrOrderAlreadyExistsForUser
					} else if orderNum == "123456789" {
						return errors.New("internal server error")
					} else if orderNum == "12345hfjfh" {
						return services.ErrInvalidOrderNumber
					} else if orderNum == "1" {
						return services.ErrInvalidOrderNumberFormat
					}

					return nil
				})

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			resp := httptest.NewRecorder()
			h.uploadOrder(resp, req.WithContext(context.WithValue(req.Context(), "userId", "user")))

			orders.AssertExpectations(t)

			assert.Equal(t, tt.want.status, resp.Code)
		})
	}
}

func Test_handler_userLogIn(t *testing.T) {
	users := mocks.NewUsers(t)
	orders := mocks.NewOrders(t)
	log := &mockLogger{}

	h := &handler{
		users:  users,
		orders: orders,
		log:    log,
	}

	type want struct {
		status int
	}

	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "1. login success",
			body: `{"login":"user","password":"user"}`,
			want: want{
				status: http.StatusOK,
			},
		},
		{
			name: "2. login fail, bad request",
			body: `{"login":"user","pass":"user"`,
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "3. login fail, unauthorized",
			body: `{"login":"user","password":"test"}`,
			want: want{
				status: http.StatusUnauthorized,
			},
		},
		{
			name: "4. login fail, internal server error",
			body: `{"login":"admin","password":"admin"}`,
			want: want{
				status: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users.
				On("LoginUser", mock.Anything, mock.Anything).
				Return(func(ctx context.Context, user *models.User) error {
					if user.Password == "test" {
						return errors.New("unauthorized")
					}

					return nil
				})

			users.
				On("GetToken", mock.Anything).
				Return(func(user *models.User) (string, error) {
					if user.Login == "admin" {
						return "", errors.New("internal server error")
					}

					return "token", nil
				})

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			resp := httptest.NewRecorder()
			h.userLogIn(resp, req)

			users.AssertExpectations(t)

			assert.Equal(t, tt.want.status, resp.Code)
		})
	}
}

func Test_handler_userSignUp(t *testing.T) {
	users := mocks.NewUsers(t)
	orders := mocks.NewOrders(t)
	log := &mockLogger{}

	h := &handler{
		users:  users,
		orders: orders,
		log:    log,
	}

	type want struct {
		status int
	}

	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "1. sign up success",
			body: `{"login":"test","password":"test"}`,
			want: want{
				status: http.StatusOK,
			},
		},
		{
			name: "2. sign up fail, bad request",
			body: `{"login":"test","pass":"test"`,
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name: "3. sign up fail, conflict",
			body: `{"login":"user","password":"user"}`,
			want: want{
				status: http.StatusConflict,
			},
		},
		{
			name: "4. sign up fail, internal server error",
			body: `{"login":"admin","password":"test"}`,
			want: want{
				status: http.StatusInternalServerError,
			},
		},
		{
			name: "5. sign up fail, get token internal server error",
			body: `{"login":"login","password":"password"}`,
			want: want{
				status: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users.
				On("RegisterUser", mock.Anything, mock.Anything).
				Return(func(ctx context.Context, user *models.User) error {
					if user.Login == "user" {
						return services.ErrUserAlreadyExists
					} else if user.Login == "admin" {
						return errors.New("internal server error")
					}
					return nil
				})

			users.
				On("GetToken", mock.Anything).
				Return(func(user *models.User) (string, error) {
					if user.Login == "login" {
						return "", errors.New("internal server error")
					}

					return "token", nil
				})

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			resp := httptest.NewRecorder()
			h.userSignUp(resp, req)

			users.AssertExpectations(t)

			assert.Equal(t, tt.want.status, resp.Code)
		})
	}
}

func Test_handler_withdraw(t *testing.T) {
	users := mocks.NewUsers(t)
	orders := mocks.NewOrders(t)
	log := &mockLogger{}

	h := &handler{
		users:  users,
		orders: orders,
		log:    log,
	}

	type want struct {
		status int
	}

	tests := []struct {
		name   string
		userId string
		body   string
		want   want
	}{
		{
			name:   "1. withdraw success",
			userId: "user",
			body:   `{"order":"123456789","sum":100}`,
			want: want{
				status: http.StatusOK,
			},
		},
		{
			name:   "2. withdraw fail, insufficient funds",
			userId: "user",
			body:   `{"order":"123456789","sum":1000}`,
			want: want{
				status: http.StatusPaymentRequired,
			},
		},
		{
			name:   "3. withdraw fail, invalid order number",
			userId: "user",
			body:   `{"order":"12345678dfg","sum":100}`,
			want: want{
				status: http.StatusBadRequest,
			},
		},
		{
			name:   "4. withdraw fail, invalid order number format",
			userId: "user",
			body:   `{"order":"1","sum":100}`,
			want: want{
				status: http.StatusUnprocessableEntity,
			},
		},
		{
			name: "5. withdraw fail, internal server error",
			body: `{"order":"123456789","sum":100}`,
			want: want{
				status: http.StatusInternalServerError,
			},
		},
		{
			name:   "6. withdraw fail, bad request",
			userId: "user",
			want: want{
				status: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users.
				On("WithdrawFromAccount", mock.Anything, mock.Anything).
				Return(func(ctx context.Context, withdrawal *models.Withdrawal) error {
					if withdrawal.Sum == 1000 {
						return services.ErrInsufficientFunds
					} else if withdrawal.OrderNumber == "12345678dfg" {
						return services.ErrInvalidOrderNumber
					} else if withdrawal.OrderNumber == "1" {
						return services.ErrInvalidOrderNumberFormat
					} else if withdrawal.UserID == "" {
						return errors.New("internal server error")
					}

					return nil
				})

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			resp := httptest.NewRecorder()
			h.withdraw(resp, req.WithContext(context.WithValue(req.Context(), "userId", tt.userId)))

			users.AssertExpectations(t)

			assert.Equal(t, tt.want.status, resp.Code)
		})
	}
}

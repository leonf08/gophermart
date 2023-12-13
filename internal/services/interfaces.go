package services

import (
	"context"
	"github.com/leonf08/gophermart.git/internal/models"
)

//go:generate mockery --name Users --output ./mocks --filename users_mock.go
//go:generate mockery --name Orders --output ./mocks --filename orders_mock.go
//go:generate mockery --name Authenticator --output ./mocks --filename authenticator_mock.go
//go:generate mockery --name UserRepo --output ./mocks --filename user_repo_mock.go
//go:generate mockery --name OrderRepo --output ./mocks --filename order_repo_mock.go
type (
	// UserRepo is an interface for working with the user repository.
	UserRepo interface {
		CreateUser(ctx context.Context, login, hashedPasswd string) (int64, error)
		GetUserByLogin(ctx context.Context, login string) (*models.User, error)
		GetUserAccount(ctx context.Context, userID int64) (*models.UserAccount, error)
		DoWithdrawal(ctx context.Context, w *models.Withdrawal) error
		GetWithdrawalList(ctx context.Context, userID int64) ([]*models.Withdrawal, error)
	}

	// OrderRepo is an interface for working with the order repository.
	OrderRepo interface {
		CreateOrder(ctx context.Context, order models.Order) error
		GetOrderByNumber(ctx context.Context, orderNum string) (*models.Order, error)
		GetOrderList(ctx context.Context, userID int64) ([]*models.Order, error)
		UpdateOrder(ctx context.Context, order *models.Order) error
	}

	// AccrualRepo is an interface for working with the accrual repository.
	AccrualRepo interface {
		UserRepo
		OrderRepo
	}

	// Authenticator is an interface for working with the authenticator service.
	Authenticator interface {
		GenerateHashFromPassword(user *models.User) (string, error)
		CheckPasswordHash(user, storedUser *models.User) error
		GenerateToken(user *models.User) (string, error)
		ValidateTokenAndExtractClaims(token string) (*models.CustomJWTClaims, error)
	}

	// Users is an interface for working with the user service.
	Users interface {
		RegisterUser(ctx context.Context, user *models.User) error
		LoginUser(ctx context.Context, user *models.User) error
		WithdrawFromAccount(ctx context.Context, w *models.Withdrawal) error
		GetUserAccount(ctx context.Context, userID int64) (*models.UserAccount, error)
		GetWithdrawals(ctx context.Context, userID int64) ([]*models.Withdrawal, error)
		GetToken(user *models.User) (string, error)
	}

	// Orders is an interface for working with the order service.
	Orders interface {
		CreateNewOrder(ctx context.Context, userID int64, orderNum string) error
		GetOrdersForUser(ctx context.Context, userID int64) ([]*models.Order, error)
	}

	// Logger is an interface for working with the logging tools
	Logger interface {
		Info(msg string, args ...any)
		Error(msg string, args ...any)
	}

	// Accrual is an interface for working with the accrual service.
	Accrual interface {
		SendOrderAccrual(orderNum string)
	}
)

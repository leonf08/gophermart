package services

import (
	"context"
	"github.com/leonf08/gophermart.git/internal/models"
)

// UserRepo is an interface for working with the user repository.
type UserRepo interface {
	CreateUser(ctx context.Context, login, hashedPasswd string) error
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	GetUserAccount(ctx context.Context, userId int64) (*models.UserAccount, error)
	DoWithdrawal(ctx context.Context, withdrawal models.Withdrawal) error
	UpdateUserAccount(ctx context.Context, userAccount *models.UserAccount) error
}

type Authenticator interface {
	GenerateHashFromPassword(user *models.User) (string, error)
	CheckPasswordHash(user, storedUser *models.User) error
	GenerateToken(user *models.User) (string, error)
}

// OrderRepo is an interface for working with the order repository.
type OrderRepo interface {
	CreateOrder(ctx context.Context, order models.Order) error
	GetOrderByNumber(ctx context.Context, orderNum string) (*models.Order, error)
	GetOrderList(ctx context.Context, userId int64) ([]*models.Order, error)
	UpdateOrder(ctx context.Context, order *models.Order) error
}

// Accrual is an interface for working with the accrual service.
type Accrual interface {
	SendOrderAccrual(orderNum string)
}

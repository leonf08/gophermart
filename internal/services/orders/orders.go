package orders

import (
	"context"
	"github.com/leonf08/gophermart.git/internal/models"
	errs "github.com/leonf08/gophermart.git/internal/services"
	"github.com/leonf08/gophermart.git/internal/services/utils"
	"time"
)

// OrderRepo is an interface for working with the order repository.
type OrderRepo interface {
	CreateOrder(ctx context.Context, order models.Order) error
	CheckOrder(ctx context.Context, orderNum string) (*models.Order, error)
	GetOrders(ctx context.Context, userId int64) ([]*models.Order, error)
}

// Accrual is an interface for working with the accrual service.
type Accrual interface {
	SendOrderAccrual(orderNum string)
}

// OrderManager is a service for working with orders.
type OrderManager struct {
	repo    OrderRepo
	accrual Accrual
}

// NewOrderManager creates a new order manager.
func NewOrderManager(repo OrderRepo, accrual Accrual) *OrderManager {
	return &OrderManager{
		repo:    repo,
		accrual: accrual,
	}
}

// CreateNewOrder creates a new order.
// If the order creation fails, an error is returned.
// If the order creation succeeds, nil is returned.
func (o *OrderManager) CreateNewOrder(ctx context.Context, userId int64, orderNum string) error {
	// Check if the order number is valid.
	if !utils.IsNumberValid(orderNum) {
		return errs.ErrInvalidOrderNumberFormat
	}

	// Check if the order number is valid by luhn algorithm.
	if !utils.LuhnValidate(orderNum) {
		return errs.ErrInvalidOrderNumber
	}

	// Check if the order already exists.
	_, err := o.repo.CheckOrder(ctx, orderNum)
	if err == nil {
		return errs.ErrOrderAlreadyExists
	}

	// Create order.
	err = o.repo.CreateOrder(ctx, models.Order{
		UserID:     userId,
		Number:     orderNum,
		Status:     models.OrderStatusNew,
		UploadedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	// Register order in accrual service.
	o.accrual.SendOrderAccrual(orderNum)

	return nil
}

// GetOrdersForUser returns all orders for a given user.
// If the order retrieval fails, an error is returned.
// If the order retrieval succeeds, the orders are returned.
func (o *OrderManager) GetOrdersForUser(ctx context.Context, userId int64) ([]*models.Order, error) {
	// Retrieve orders.
	orders, err := o.repo.GetOrders(ctx, userId)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

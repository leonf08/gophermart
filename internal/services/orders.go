package services

import (
	"context"
	"github.com/leonf08/gophermart.git/internal/models"
	"github.com/leonf08/gophermart.git/internal/services/utils"
	"time"
)

// OrderManager is a service for working with orders.
type OrderManager struct {
	repo    OrderRepo
	accrual AccrualService
}

// NewOrderManager creates a new order manager.
func NewOrderManager(repo OrderRepo, accrual AccrualService) *OrderManager {
	return &OrderManager{
		repo:    repo,
		accrual: accrual,
	}
}

// CreateNewOrder creates a new order.
// If the order creation fails, an error is returned.
// If the order creation succeeds, nil is returned.
func (o *OrderManager) CreateNewOrder(ctx context.Context, userId, orderNum string) error {
	// Check if the order number is valid.
	if !utils.IsNumber(orderNum) {
		return ErrInvalidOrderNumber
	}

	// Check if the order number is valid by luhn algorithm.
	if !utils.LuhnValidate(orderNum) {
		return ErrInvalidOrderNumberFormat
	}

	// Check if the order already exists.
	order, err := o.repo.GetOrderByNumber(ctx, orderNum)
	if err == nil {
		if order.UserID == userId {
			return ErrOrderAlreadyExistsForUser
		}
		return ErrOrderAlreadyExists
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
func (o *OrderManager) GetOrdersForUser(ctx context.Context, userId string) ([]*models.Order, error) {
	// Retrieve orders.
	orders, err := o.repo.GetOrderList(ctx, userId)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

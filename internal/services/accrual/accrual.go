package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/leonf08/gophermart.git/internal/models"
	"github.com/leonf08/gophermart.git/internal/services"
	"net/http"
	"strconv"
	"time"
)

// Accrual is a service for working with the accrual system.
type Accrual struct {
	address   string
	orderRepo services.OrderRepo
	userRepo  services.UserRepo
	orderNum  chan string
}

// NewAccrual creates a new accrual service.
func NewAccrual(orderRepo services.OrderRepo, userRepo services.UserRepo) *Accrual {
	return &Accrual{
		orderRepo: orderRepo,
		userRepo:  userRepo,
		orderNum:  make(chan string),
	}
}

// SendOrderAccrual sends an order number to the accrual system.
func (a *Accrual) SendOrderAccrual(orderNum string) {
	a.orderNum <- orderNum
}

// Run starts the accrual service.
func (a *Accrual) Run(ctx context.Context) error {
	orderNum := <-a.orderNum
	url := fmt.Sprintf("%s/api/orders/%s", a.address, orderNum)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		accrualResp := &models.AccrualResponse{}
		if err = json.NewDecoder(resp.Body).Decode(accrualResp); err != nil {
			return err
		}

		switch accrualResp.Status {
		case models.OrderStatusInvalid:
			if err = a.orderRepo.UpdateOrder(ctx, &models.Order{
				Number: orderNum,
				Status: models.OrderStatusInvalid,
			}); err != nil {
				return err
			}
		case models.OrderStatusProcessed:
			order, err := a.orderRepo.GetOrderByNumber(ctx, orderNum)
			if err != nil {
				return err
			}

			order.Status, order.Accrual = models.OrderStatusProcessed, accrualResp.Accrual
			if err = a.orderRepo.UpdateOrder(ctx, order); err != nil {
				return err
			}

			userAccount, err := a.userRepo.GetUserAccount(ctx, order.UserID)
			if err != nil {
				return err
			}

			userAccount.Current += accrualResp.Accrual
			if err = a.userRepo.UpdateUserAccount(ctx, userAccount); err != nil {
				return err
			}
		case models.OrderStatusProcessing:
			if err = a.orderRepo.UpdateOrder(ctx, &models.Order{
				Number: orderNum,
				Status: models.OrderStatusProcessing,
			}); err != nil {
				return err
			}

			a.SendOrderAccrual(orderNum)
		case models.OrderStatusRegistered:
			a.SendOrderAccrual(orderNum)
		}
	case http.StatusTooManyRequests:
		pause, err := strconv.Atoi(resp.Header.Get("Retry-After"))
		if err != nil {
			return err
		}

		time.Sleep(time.Duration(pause) * time.Second)
		a.SendOrderAccrual(orderNum)
	case http.StatusInternalServerError:
		return services.ErrAccrualInternalError
	}

	return nil
}

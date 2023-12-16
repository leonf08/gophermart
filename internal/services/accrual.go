package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/leonf08/gophermart.git/internal/models"
	"net/http"
	"strconv"
	"time"
)

// AccrualService is a service for working with the accrual system.
type AccrualService struct {
	address  string
	repo     AccrualRepo
	log      Logger
	orderNum chan string
}

// NewAccrual creates a new accrual service.
func NewAccrual(address string, repo AccrualRepo, log Logger) *AccrualService {
	a := &AccrualService{
		address:  address,
		repo:     repo,
		log:      log,
		orderNum: make(chan string, 10),
	}

	go a.run(context.Background())

	return a
}

// SendOrderAccrual sends an order number to the accrual system.
func (a *AccrualService) SendOrderAccrual(orderNum string) {
	a.orderNum <- orderNum
}

// Run starts the accrual service.
func (a *AccrualService) run(ctx context.Context) {
	for {
		orderNum := <-a.orderNum
		if err := a.process(ctx, orderNum); err != nil {
			a.log.Error("accrual - run - a.process", "error", err)
		}
	}
}

func (a *AccrualService) process(ctx context.Context, orderNum string) error {
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
			if err = a.repo.UpdateOrder(ctx, &models.Order{
				Number: orderNum,
				Status: models.OrderStatusInvalid,
			}); err != nil {
				return err
			}
		case models.OrderStatusProcessed:
			order, err := a.repo.GetOrderByNumber(ctx, orderNum)
			if err != nil {
				return err
			}

			order.Status, order.Accrual = models.OrderStatusProcessed, accrualResp.Accrual*100
			if err = a.repo.UpdateOrder(ctx, order); err != nil {
				return err
			}
		case models.OrderStatusProcessing:
			if err = a.repo.UpdateOrder(ctx, &models.Order{
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
		return fmt.Errorf("accrual system internal error")
	}

	return nil
}

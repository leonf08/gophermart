package services

import (
	"context"
	"errors"
	"github.com/leonf08/gophermart.git/internal/models"
	"github.com/leonf08/gophermart.git/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

type mockAccrual struct{}

func (m *mockAccrual) SendOrderAccrual(orderNum string) {}

func TestNewOrderManager(t *testing.T) {
	type args struct {
		repo    OrderRepo
		accrual Accrual
	}
	tests := []struct {
		name string
		args args
		want *OrderManager
	}{
		{
			name: "NewOrderManager",
			args: args{
				repo:    nil,
				accrual: nil,
			},
			want: &OrderManager{
				repo:    nil,
				accrual: nil,
			},
		},
		{
			name: "NewOrderManager",
			args: args{
				repo:    mocks.NewOrderRepo(t),
				accrual: &mockAccrual{},
			},
			want: &OrderManager{
				repo:    mocks.NewOrderRepo(t),
				accrual: &mockAccrual{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOrderManager(tt.args.repo, tt.args.accrual); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrderManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrderManager_CreateNewOrder(t *testing.T) {
	repo := mocks.NewOrderRepo(t)
	accr := &mockAccrual{}

	type args struct {
		userID   string
		orderNum string
	}
	type want struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "CreateNewOrder_no_error",
			args: args{
				userID:   "1",
				orderNum: "2030",
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "CreateNewOrder_invalid_order_number",
			args: args{
				userID:   "1",
				orderNum: "2030a",
			},
			want: want{
				err: ErrInvalidOrderNumber,
			},
		},
		{
			name: "CreateNewOrder_invalid_order_number_format",
			args: args{
				userID:   "1",
				orderNum: "2040",
			},
			want: want{
				err: ErrInvalidOrderNumberFormat,
			},
		},
		{
			name: "CreateNewOrder_order_already_exists_for_user",
			args: args{
				userID:   "1",
				orderNum: "4010",
			},
			want: want{
				err: ErrOrderAlreadyExistsForUser,
			},
		},
		{
			name: "CreateNewOrder_order_already_exists",
			args: args{
				userID:   "2",
				orderNum: "4010",
			},
			want: want{
				err: ErrOrderAlreadyExists,
			},
		},
		{
			name: "CreateNewOrder_error",
			args: args{
				userID:   "2",
				orderNum: "2030",
			},
			want: want{
				err: errors.New("error"),
			},
		},
	}

	repo.
		On("GetOrderByNumber", context.Background(), mock.Anything).
		Return(func(ctx context.Context, num string) (*models.Order, error) {
			if num == "2030" {
				return nil, errors.New("error")
			}

			if num == "4010" {
				return &models.Order{
					UserID: "1",
				}, nil
			}

			return nil, nil
		})

	repo.
		On("CreateOrder", context.Background(), mock.Anything).
		Return(func(ctx context.Context, order models.Order) error {
			if order.UserID == "2" {
				return errors.New("error")
			}

			return nil
		})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrderManager{
				repo:    repo,
				accrual: accr,
			}

			err := o.CreateNewOrder(context.Background(), tt.args.userID, tt.args.orderNum)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestOrderManager_GetOrdersForUser(t *testing.T) {
	repo := mocks.NewOrderRepo(t)
	accr := &mockAccrual{}

	type args struct {
		userID string
	}
	type want struct {
		orders []*models.Order
		err    bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "GetOrdersForUser_no_error",
			args: args{
				userID: "1",
			},
			want: want{
				orders: []*models.Order{
					{
						UserID: "1",
					},
				},
				err: false,
			},
		},
		{
			name: "GetOrdersForUser_error",
			args: args{
				userID: "2",
			},
			want: want{
				orders: nil,
				err:    true,
			},
		},
	}

	repo.
		On("GetOrderList", context.Background(), mock.Anything).
		Return(func(ctx context.Context, userID string) ([]*models.Order, error) {
			if userID == "1" {
				return []*models.Order{
					{
						UserID: "1",
					},
				}, nil
			}

			return nil, errors.New("error")
		})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrderManager{
				repo:    repo,
				accrual: accr,
			}
			got, err := o.GetOrdersForUser(context.Background(), tt.args.userID)
			assert.Equal(t, tt.want.orders, got)
			assert.Equal(t, tt.want.err, err != nil)
		})
	}
}

// Code generated by mockery v2.37.1. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/leonf08/gophermart.git/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// UserRepo is an autogenerated mock type for the UserRepo type
type UserRepo struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: ctx, login, hashedPasswd
func (_m *UserRepo) CreateUser(ctx context.Context, login, hashedPasswd string) (int64, error) {
	ret := _m.Called(ctx, login, hashedPasswd)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) error); ok {
		r0 = rf(ctx, login, hashedPasswd)
	} else {
		r0 = ret.Error(0)
	}

	return 0, r0
}

// DoWithdrawal provides a mock function with given fields: ctx, w
func (_m *UserRepo) DoWithdrawal(ctx context.Context, w *models.Withdrawal) error {
	ret := _m.Called(ctx, w)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Withdrawal) error); ok {
		r0 = rf(ctx, w)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUserAccount provides a mock function with given fields: ctx, userID
func (_m *UserRepo) GetUserAccount(ctx context.Context, userID int64) (*models.UserAccount, error) {
	ret := _m.Called(ctx, userID)

	var r0 *models.UserAccount
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*models.UserAccount, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *models.UserAccount); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.UserAccount)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByLogin provides a mock function with given fields: ctx, login
func (_m *UserRepo) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	ret := _m.Called(ctx, login)

	var r0 *models.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.User, error)); ok {
		return rf(ctx, login)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.User); ok {
		r0 = rf(ctx, login)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, login)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWithdrawalList provides a mock function with given fields: ctx, userID
func (_m *UserRepo) GetWithdrawalList(ctx context.Context, userID int64) ([]*models.Withdrawal, error) {
	ret := _m.Called(ctx, userID)

	var r0 []*models.Withdrawal
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) ([]*models.Withdrawal, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) []*models.Withdrawal); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Withdrawal)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewUserRepo creates a new instance of UserRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepo {
	mock := &UserRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

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

func TestNewUserManager(t *testing.T) {
	type args struct {
		repo UserRepo
		auth Authenticator
	}
	tests := []struct {
		name string
		args args
		want *UserManager
	}{
		{
			name: "TestNewUserManager",
			args: args{
				repo: nil,
				auth: nil,
			},
			want: &UserManager{},
		},
		{
			name: "TestNewUserManager",
			args: args{
				repo: mocks.NewUserRepo(t),
				auth: mocks.NewAuthenticator(t),
			},
			want: &UserManager{mocks.NewUserRepo(t), mocks.NewAuthenticator(t)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserManager(tt.args.repo, tt.args.auth); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserManager_GetToken(t *testing.T) {
	repo := mocks.NewUserRepo(t)
	auth := mocks.NewAuthenticator(t)

	type args struct {
		user *models.User
	}
	type want struct {
		token string
		err   error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "TestUserManager_GetToken_no_error",
			args: args{
				user: &models.User{
					Login:    "test",
					Password: "test",
				},
			},
			want: want{
				token: "test",
				err:   nil,
			},
		},
		{
			name: "TestUserManager_GetToken_error",
			args: args{
				user: &models.User{
					Login:    "user",
					Password: "test",
				},
			},
			want: want{
				token: "",
				err:   ErrGenerateToken,
			},
		},
	}

	auth.
		On("GenerateToken", mock.Anything).
		Return(func(user *models.User) (string, error) {
			if user.Login != "test" {
				return "", ErrGenerateToken
			}

			return "test", nil
		})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserManager{
				repo: repo,
				auth: auth,
			}
			got, err := u.GetToken(tt.args.user)

			assert.Equal(t, tt.want.token, got)
			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestUserManager_GetUserAccount(t *testing.T) {
	repo := mocks.NewUserRepo(t)
	auth := mocks.NewAuthenticator(t)

	type args struct {
		userID int64
	}
	type want struct {
		userAccount *models.UserAccount
		err         bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "TestUserManager_GetUserAccount_no_error",
			args: args{
				userID: 1,
			},
			want: want{
				userAccount: &models.UserAccount{
					UserID: 1,
				},
				err: false,
			},
		},
		{
			name: "TestUserManager_GetUserAccount_error",
			args: args{
				userID: 2,
			},
			want: want{
				userAccount: nil,
				err:         true,
			},
		},
	}

	repo.
		On("GetUserAccount", mock.Anything, mock.Anything).
		Return(func(ctx context.Context, userID int64) (*models.UserAccount, error) {
			if userID != 1 {
				return nil, errors.New("error")
			}

			return &models.UserAccount{
				UserID: userID,
			}, nil
		})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserManager{
				repo: repo,
				auth: auth,
			}
			got, err := u.GetUserAccount(context.Background(), tt.args.userID)

			assert.Equal(t, tt.want.userAccount, got)
			assert.Equal(t, tt.want.err, err != nil)
		})
	}
}

func TestUserManager_GetWithdrawals(t *testing.T) {
	repo := mocks.NewUserRepo(t)
	auth := mocks.NewAuthenticator(t)

	type args struct {
		userID int64
	}
	type want struct {
		withdrawals []*models.Withdrawal
		err         bool
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "TestUserManager_GetWithdrawals_no_error",
			args: args{
				userID: 1,
			},
			want: want{
				withdrawals: []*models.Withdrawal{
					{
						UserID: 1,
					},
				},
				err: false,
			},
		},
		{
			name: "TestUserManager_GetWithdrawals_error",
			args: args{
				userID: 2,
			},
			want: want{
				withdrawals: nil,
				err:         true,
			},
		},
	}

	repo.
		On("GetWithdrawalList", mock.Anything, mock.Anything).
		Return(func(ctx context.Context, userID int64) ([]*models.Withdrawal, error) {
			if userID != 1 {
				return nil, errors.New("error")
			}

			return []*models.Withdrawal{
				{
					UserID: userID,
				},
			}, nil
		})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserManager{
				repo: repo,
				auth: auth,
			}
			got, err := u.GetWithdrawals(context.Background(), tt.args.userID)

			assert.Equal(t, tt.want.withdrawals, got)
			assert.Equal(t, tt.want.err, err != nil)
		})
	}
}

func TestUserManager_LoginUser(t *testing.T) {
	repo := mocks.NewUserRepo(t)
	auth := mocks.NewAuthenticator(t)

	type args struct {
		user *models.User
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
			name: "TestUserManager_LoginUser_no_error",
			args: args{
				user: &models.User{
					Login:    "test",
					Password: "test",
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "TestUserManager_LoginUser_invalid_login_format",
			args: args{
				user: &models.User{
					Login: "",
				},
			},
			want: want{
				err: ErrInvalidLoginFormat,
			},
		},
		{
			name: "TestUserManager_LoginUser_user_not_found",
			args: args{
				user: &models.User{
					Login: "user",
				},
			},
			want: want{
				err: ErrUserNotFound,
			},
		},
		{
			name: "TestUserManager_LoginUser_incorrect_password",
			args: args{
				user: &models.User{
					Login:    "test",
					Password: "user",
				},
			},
			want: want{
				err: ErrIncorrectPassword,
			},
		},
	}

	repo.
		On("GetUserByLogin", mock.Anything, mock.Anything).
		Return(func(ctx context.Context, login string) (*models.User, error) {
			if login != "test" {
				return nil, ErrUserNotFound
			}

			return &models.User{
				Login:    login,
				Password: "test",
			}, nil
		})

	auth.
		On("CheckPasswordHash", mock.Anything, mock.Anything).
		Return(func(user, storedUser *models.User) error {
			if user.Password != "test" {
				return ErrIncorrectPassword
			}

			return nil
		})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserManager{
				repo: repo,
				auth: auth,
			}
			err := u.LoginUser(context.Background(), tt.args.user)

			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestUserManager_RegisterUser(t *testing.T) {
	repo := mocks.NewUserRepo(t)
	auth := mocks.NewAuthenticator(t)

	type args struct {
		user *models.User
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
			name: "TestUserManager_RegisterUser_no_error",
			args: args{
				user: &models.User{
					Login:    "test",
					Password: "test",
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "TestUserManager_RegisterUser_generate_hash_from_password_error",
			args: args{
				user: &models.User{
					Login:    "test",
					Password: "",
				},
			},
			want: want{
				err: ErrGenerateHashFromPassword,
			},
		},
		{
			name: "TestUserManager_RegisterUser_create_user_error",
			args: args{
				user: &models.User{
					Login:    "admin",
					Password: "test",
				},
			},
			want: want{
				err: errors.New("error"),
			},
		},
	}

	repo.
		On("CreateUser", mock.Anything, mock.Anything, mock.Anything).
		Return(func(ctx context.Context, login, hashedPasswd string) error {
			if login == "admin" {
				return errors.New("error")
			}

			return nil
		})

	auth.
		On("GenerateHashFromPassword", mock.Anything).
		Return(func(user *models.User) (string, error) {
			if user.Password == "" {
				return "", ErrGenerateHashFromPassword
			}

			return "test", nil
		})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserManager{
				repo: repo,
				auth: auth,
			}

			err := u.RegisterUser(context.Background(), tt.args.user)

			assert.Equal(t, tt.want.err, err)
		})
	}
}

func TestUserManager_WithdrawFromAccount(t *testing.T) {
	repo := mocks.NewUserRepo(t)
	auth := mocks.NewAuthenticator(t)

	type args struct {
		w *models.Withdrawal
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
			name: "TestUserManager_WithdrawFromAccount_no_error",
			args: args{
				w: &models.Withdrawal{
					UserID:      1,
					Sum:         1,
					OrderNumber: "2030",
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "TestUserManager_WithdrawFromAccount_invalid_order_number_format",
			args: args{
				w: &models.Withdrawal{
					UserID:      2,
					Sum:         1,
					OrderNumber: "test",
				},
			},
			want: want{
				err: ErrInvalidOrderNumberFormat,
			},
		},
		{
			name: "TestUserManager_WithdrawFromAccount_invalid_order_number",
			args: args{
				w: &models.Withdrawal{
					UserID:      2,
					Sum:         1,
					OrderNumber: "2031",
				},
			},
			want: want{
				err: ErrInvalidOrderNumber,
			},
		},
		{
			name: "TestUserManager_WithdrawFromAccount_insufficient_funds",
			args: args{
				w: &models.Withdrawal{
					UserID:      2,
					Sum:         100000,
					OrderNumber: "2030",
				},
			},
			want: want{
				err: ErrInsufficientFunds,
			},
		},
		{
			name: "TestUserManager_WithdrawFromAccount_get_user_account_error",
			args: args{
				w: &models.Withdrawal{
					UserID:      3,
					Sum:         1,
					OrderNumber: "2030",
				},
			},
			want: want{
				err: errors.New("error"),
			},
		},
		{
			name: "TestUserManager_WithdrawFromAccount_do_withdrawal_error",
			args: args{
				w: &models.Withdrawal{
					UserID:      2,
					Sum:         1,
					OrderNumber: "2030",
				},
			},
			want: want{
				err: errors.New("error"),
			},
		},
	}

	repo.
		On("GetUserAccount", mock.Anything, mock.Anything).
		Return(func(ctx context.Context, userID int64) (*models.UserAccount, error) {
			if userID == 3 {
				return nil, errors.New("error")
			}

			return &models.UserAccount{
				UserID:  userID,
				Current: 1000,
			}, nil
		})

	repo.
		On("DoWithdrawal", mock.Anything, mock.Anything).
		Return(func(ctx context.Context, withdrawal *models.Withdrawal) error {
			if withdrawal.UserID == 2 {
				return errors.New("error")
			}

			return nil
		})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserManager{
				repo: repo,
				auth: auth,
			}

			err := u.WithdrawFromAccount(context.Background(), tt.args.w)

			assert.Equal(t, tt.want.err, err)
		})
	}
}

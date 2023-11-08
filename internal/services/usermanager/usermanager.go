package usermanager

import (
	"context"
	"fmt"
	"github.com/leonf08/gophermart.git/internal/models"
	errs "github.com/leonf08/gophermart.git/internal/services"
	"github.com/leonf08/gophermart.git/internal/services/utils"
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

type UserManager struct {
	repo UserRepo
	auth Authenticator
}

func NewUserManager(repo UserRepo, auth Authenticator) *UserManager {
	return &UserManager{
		repo: repo,
		auth: auth,
	}
}

// RegisterUser registers a new user.
// If the user registration fails, an error is returned.
// If the user registration succeeds, nil is returned.
// The user registration fails if the user already exists.
func (u *UserManager) RegisterUser(ctx context.Context, user *models.User) error {
	// Check if the user already exists.
	_, err := u.repo.GetUserByLogin(ctx, user.Login)
	if err == nil {
		return errs.ErrUserAlreadyExists
	}

	// Generate hash from password.
	hashedPasswd, err := u.auth.GenerateHashFromPassword(user)
	if err != nil {
		return err
	}

	// Create user.
	err = u.repo.CreateUser(ctx, user.Login, hashedPasswd)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// LoginUser logs in a user.
// If the user login fails, an error is returned.
// If the user login succeeds, nil is returned.
// The user login fails if the user does not exist or the password is incorrect.
// The user login succeeds if the user exists and the password is correct.
func (u *UserManager) LoginUser(ctx context.Context, user *models.User) error {
	// Check if the user exists.
	storedUser, err := u.repo.GetUserByLogin(ctx, user.Login)
	if err != nil {
		return err
	}

	// Check if the password is correct.
	err = u.auth.CheckPasswordHash(user, storedUser)
	if err != nil {
		return err
	}

	return nil
}

// GetToken generates a JWT token for a given user.
// If the token generation fails, an error is returned.
// If the token generation succeeds, the token is returned.
// The token is signed with HMAC-SHA256 algorithm.
func (u *UserManager) GetToken(user *models.User) (string, error) {
	token, err := u.auth.GenerateToken(user)

	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserAccount returns a user account.
// If the user account is found, it returns the user account and nil.
// If the user account is not found, it returns nil and an error.
func (u *UserManager) GetUserAccount(ctx context.Context, userId int64) (*models.UserAccount, error) {
	userAccount, err := u.repo.GetUserAccount(ctx, userId)
	if err != nil {
		return nil, err
	}

	return userAccount, nil
}

// WithdrawFromAccount withdraws a given sum from a user account.
// If the withdrawal succeeds, it returns nil.
// If the withdrawal fails, it returns an error.
// The withdrawal fails if the sum is greater than the current balance.
func (u *UserManager) WithdrawFromAccount(ctx context.Context, w models.Withdrawal) error {
	// Check if the orderNumber is valid.
	if !utils.IsNumberValid(w.OrderNumber) {
		return errs.ErrInvalidOrderNumberFormat
	}

	// Check if the orderNumber is valid by luhn algorithm.
	if !utils.LuhnValidate(w.OrderNumber) {
		return errs.ErrInvalidOrderNumber
	}

	// Get user account.
	userAccount, err := u.repo.GetUserAccount(ctx, w.UserID)
	if err != nil {
		return err
	}

	// Check if the sum is greater than the current balance.
	if userAccount.Current < w.Sum {
		return errs.ErrInsufficientFunds
	}

	// Withdraw from account.
	err = u.repo.DoWithdrawal(ctx, w)
	if err != nil {
		return err
	}

	// Update user account.
	userAccount.Current -= w.Sum
	userAccount.Withdrawn += w.Sum
	err = u.repo.UpdateUserAccount(ctx, userAccount)
	if err != nil {
		return err
	}

	return nil
}

package services

import (
	"context"
	"github.com/leonf08/gophermart.git/internal/models"
	"github.com/leonf08/gophermart.git/internal/services/utils"
	"time"
)

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
	// Generate hash from password.
	hashedPasswd, err := u.auth.GenerateHashFromPassword(user)
	if err != nil {
		return err
	}

	// Create user.
	userID, err := u.repo.CreateUser(ctx, user.Login, hashedPasswd)
	if err != nil {
		return err
	}

	user.UserID = userID

	return nil
}

// LoginUser logs in a user.
// If the user login fails, an error is returned.
// If the user login succeeds, nil is returned.
// The user login fails if the user does not exist or the password is incorrect.
// The user login succeeds if the user exists and the password is correct.
func (u *UserManager) LoginUser(ctx context.Context, user *models.User) error {
	// Check if the login is valid.
	if !utils.IsLoginValid(user.Login) {
		return ErrInvalidLoginFormat
	}
	// Check if the user exists.
	storedUser, err := u.repo.GetUserByLogin(ctx, user.Login)
	if err != nil {
		return ErrUserNotFound
	}

	// Check if the password is correct.
	err = u.auth.CheckPasswordHash(user, storedUser)
	if err != nil {
		return ErrIncorrectPassword
	}

	user.UserID = storedUser.UserID

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

// GetUserAccount returns details about user account.
// If the user account is found, it returns the user account and nil.
// If the user account is not found, it returns nil and an error.
func (u *UserManager) GetUserAccount(ctx context.Context, userID int64) (*models.UserAccount, error) {
	userAccount, err := u.repo.GetUserAccount(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert integer sum to float sum
	userAccount.Current /= 100
	userAccount.Withdrawn /= 100

	return userAccount, nil
}

// WithdrawFromAccount withdraws a given sum from a user account.
// If the withdrawal succeeds, it returns nil.
// If the withdrawal fails, it returns an error.
// The withdrawal fails if the sum is greater than the current balance.
func (u *UserManager) WithdrawFromAccount(ctx context.Context, w *models.Withdrawal) error {
	// Convert float sum to integer sum
	w.Sum *= 100
	// Check if the orderNumber is valid.
	if !utils.IsNumber(w.OrderNumber) {
		return ErrInvalidOrderNumberFormat
	}

	// Check if the orderNumber is valid by luhn algorithm.
	if !utils.LuhnValidate(w.OrderNumber) {
		return ErrInvalidOrderNumber
	}

	// Get user account.
	userAccount, err := u.repo.GetUserAccount(ctx, w.UserID)
	if err != nil {
		return err
	}

	// Check if the sum is greater than the current balance.
	if userAccount.Current < w.Sum {
		return ErrInsufficientFunds
	}

	// Withdraw from account.
	w.ProcessedAt = time.Now()
	err = u.repo.DoWithdrawal(ctx, w)
	if err != nil {
		return err
	}

	return nil
}

// GetWithdrawals returns a list of withdrawals.
// If the list of withdrawals is found, it returns the list of withdrawals and nil.
// If the list of withdrawals is not found, it returns nil and an error.
func (u *UserManager) GetWithdrawals(ctx context.Context, userID int64) ([]*models.Withdrawal, error) {
	withdrawals, err := u.repo.GetWithdrawalList(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Convert integer sum to float sum
	for _, w := range withdrawals {
		w.Sum /= 100
	}

	return withdrawals, nil
}

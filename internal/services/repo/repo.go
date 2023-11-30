package repo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/leonf08/gophermart.git/internal/models"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// CreateUser creates a new user in database.
// If user creation fails, returns error.
// If user creation succeeds, returns nil.
func (r *Repository) CreateUser(ctx context.Context, login, hashedPasswd string) error {
	query := `INSERT INTO users (user_id, login, password) VALUES (gen_random_uuid(), $1, $2)`
	_, err := r.db.ExecContext(ctx, query, login, hashedPasswd)
	if err != nil {
		return err
	}

	return nil
}

// GetUserByLogin gets a user from database by login.
// If user does not exist, returns error.
// If user exists, returns nil.
func (r *Repository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	query := `SELECT * FROM users WHERE login = $1`
	user := &models.User{}
	err := r.db.GetContext(ctx, user, query, login)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserID gets a user id from database by login.
// If user id does not exist, returns error.
// If user id exists, returns nil.
func (r *Repository) GetUserID(ctx context.Context, login string) (string, error) {
	query := `SELECT user_id FROM users WHERE login = $1`
	var userID string
	err := r.db.GetContext(ctx, &userID, query, login)
	if err != nil {
		return "", err
	}

	return userID, nil
}

// GetUserAccount gets a user account from database by user id.
// If user account does not exist, returns error.
// If user account exists, returns nil.
func (r *Repository) GetUserAccount(ctx context.Context, userID string) (*models.UserAccount, error) {
	query := `SELECT * FROM user_accounts WHERE user_id = $1`
	userAccount := &models.UserAccount{}
	err := r.db.GetContext(ctx, userAccount, query, userID)
	if err != nil {
		return nil, err
	}

	return userAccount, nil
}

// DoWithdrawal does a withdrawal in withdrawal table.
// If withdrawal fails, returns error.
// If withdrawal succeeds, returns nil.
func (r *Repository) DoWithdrawal(ctx context.Context, withdrawal *models.Withdrawal) error {
	query := `INSERT INTO withdrawals (user_id, order_number, sum, processed_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, withdrawal.UserID, withdrawal.OrderNumber,
		withdrawal.Sum, withdrawal.ProcessedAt)
	if err != nil {
		return err
	}

	return nil
}

// CreateUserAccount creates a new user account in database.
// If user account creation fails, returns error.
// If user account creation succeeds, returns nil.
func (r *Repository) CreateUserAccount(ctx context.Context, userID string) error {
	query := `INSERT INTO user_accounts (user_id) VALUES ($1)`
	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUserAccount updates a user account.
// If update fails, returns error.
// If update succeeds, returns nil.
func (r *Repository) UpdateUserAccount(ctx context.Context, userAccount *models.UserAccount) error {
	query := `UPDATE user_accounts SET current = $1, withdrawn = $2 WHERE user_id = $3`
	_, err := r.db.ExecContext(ctx, query, userAccount.Current, userAccount.Withdrawn, userAccount.UserID)
	if err != nil {
		return err
	}

	return nil
}

// GetWithdrawalList gets a list of withdrawals from database by user id.
// If list of withdrawals does not exist, returns error.
// If list of withdrawals exists, returns nil.
func (r *Repository) GetWithdrawalList(ctx context.Context, userID string) ([]*models.Withdrawal, error) {
	query := `SELECT * FROM withdrawals WHERE user_id = $1`
	withdrawals := make([]*models.Withdrawal, 0)
	err := r.db.SelectContext(ctx, &withdrawals, query, userID)
	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}

// CreateOrder creates a new order in database.
// If order creation fails, returns error.
// If order creation succeeds, returns nil.
func (r *Repository) CreateOrder(ctx context.Context, order models.Order) error {
	query := `INSERT INTO orders (user_id, number, status, uploaded_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, order.UserID, order.Number,
		order.Status, order.UploadedAt)
	if err != nil {
		return err
	}

	return nil
}

// GetOrderByNumber gets an order from database by order number.
// If order does not exist, returns error.
// If order exists, returns nil.
func (r *Repository) GetOrderByNumber(ctx context.Context, orderNum string) (*models.Order, error) {
	query := `SELECT * FROM orders WHERE number = $1`
	order := &models.Order{}
	err := r.db.GetContext(ctx, order, query, orderNum)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrderList gets a list of orders from database by user id.
// If list of orders does not exist, returns error.
// If list of orders exists, returns nil.
func (r *Repository) GetOrderList(ctx context.Context, userID string) ([]*models.Order, error) {
	query := `SELECT * FROM orders WHERE user_id = $1`
	orders := make([]*models.Order, 0)
	err := r.db.SelectContext(ctx, &orders, query, userID)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

// UpdateOrder updates an order.
// If update fails, returns error.
// If update succeeds, returns nil.
func (r *Repository) UpdateOrder(ctx context.Context, order *models.Order) error {
	if order.Accrual == 0 {
		query := `UPDATE orders SET status = $1 WHERE number = $2`
		_, err := r.db.ExecContext(ctx, query, order.Status, order.Number)
		if err != nil {
			return err
		}

	} else {
		query := `UPDATE orders SET status = $1, accrual = $2 WHERE number = $3`
		_, err := r.db.ExecContext(ctx, query, order.Status, order.Accrual, order.Number)
		if err != nil {
			return err
		}

	}

	return nil
}

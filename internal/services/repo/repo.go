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
func (r *Repository) CreateUser(ctx context.Context, login, hashedPasswd string) (int64, error) {
	query := `INSERT INTO users (login, password) VALUES ($1, $2) RETURNING user_id`

	var userID int64
	err := r.db.QueryRowContext(ctx, query, login, hashedPasswd).Scan(&userID)

	return userID, err
}

// GetUserByLogin gets a user from database by login.
// If user does not exist, returns error.
// If user exists, returns nil.
func (r *Repository) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	query := `SELECT user_id, login, password FROM users WHERE login = $1`
	user := &models.User{}
	err := r.db.GetContext(ctx, user, query, login)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserAccount gets a user account from database by user id.
// If user account does not exist, returns error.
// If user account exists, returns nil.
func (r *Repository) GetUserAccount(ctx context.Context, userID int64) (*models.UserAccount, error) {
	query := `SELECT user_id, current, withdrawn FROM users WHERE user_id = $1`
	userAcc := &models.UserAccount{}
	err := r.db.GetContext(ctx, userAcc, query, userID)
	if err != nil {
		return nil, err
	}

	return userAcc, nil
}

// DoWithdrawal does a withdrawal and updates user account.
// If withdrawal fails, returns error.
// If withdrawal succeeds, returns nil.
func (r *Repository) DoWithdrawal(ctx context.Context, w *models.Withdrawal) error {
	queryWithdraw := `INSERT INTO withdrawals (user_id, order_number, sum, updated_at) VALUES ($1, $2, $3, $4)`
	queryUpdateAcc := `UPDATE users SET current = current - $1, withdrawn = withdrawn + $1 WHERE user_id = $2`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, queryWithdraw, w.UserID, w.OrderNumber, w.Sum, w.ProcessedAt)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, queryUpdateAcc, w.Sum, w.UserID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetWithdrawalList gets a list of withdrawals from database by user id.
// If list of withdrawals does not exist, returns error.
// If list of withdrawals exists, returns nil.
func (r *Repository) GetWithdrawalList(ctx context.Context, userID int64) ([]*models.Withdrawal, error) {
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
	query := `INSERT INTO orders (user_id, number, status, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, order.UserID, order.Number, order.Status, order.UploadedAt)
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
func (r *Repository) GetOrderList(ctx context.Context, userID int64) ([]*models.Order, error) {
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
		queryUpdateOrder := `UPDATE orders SET status = $1, accrual = $2 WHERE number = $3`
		queryUpdateAcc := `UPDATE users SET current = current + $1 WHERE user_id = $2`

		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		defer tx.Rollback()

		_, err = tx.ExecContext(ctx, queryUpdateOrder, order.Status, order.Accrual, order.Number)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, queryUpdateAcc, order.Accrual, order.UserID)
		if err != nil {
			return err
		}

		return tx.Commit()
	}

	return nil
}

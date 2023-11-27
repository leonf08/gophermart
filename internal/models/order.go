package models

import "time"

const (
	OrderStatusNew        = "NEW"
	OrderStatusRegistered = "REGISTERED"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessed  = "PROCESSED"
)

type (
	Order struct {
		UserID     string    `json:"-" db:"user_id"`
		Number     string    `json:"number" db:"number"`
		Status     string    `json:"status" db:"status"`
		Accrual    float64   `json:"accrual,omitempty" db:"accrual"`
		UploadedAt time.Time `json:"uploaded_at" db:"uploaded_at"`
	}

	Withdrawal struct {
		UserID      string    `json:"-" db:"user_id"`
		OrderNumber string    `json:"order" db:"order_number"`
		Sum         float64   `json:"sum" db:"sum"`
		ProcessedAt time.Time `json:"processed_at,omitempty" db:"processed_at"`
	}

	AccrualResponse struct {
		OrderNumber string  `json:"order"`
		Status      string  `json:"status"`
		Accrual     float64 `json:"accrual,omitempty"`
	}
)

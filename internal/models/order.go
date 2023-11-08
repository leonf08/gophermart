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
		UserID     int64     `json:"-"`
		Number     string    `json:"number"`
		Status     string    `json:"status"`
		Accrual    int64     `json:"accrual,omitempty"`
		UploadedAt time.Time `json:"uploaded_at"`
	}

	Withdrawal struct {
		UserID      int64     `json:"-"`
		OrderNumber string    `json:"order"`
		Sum         int64     `json:"sum"`
		ProcessedAt time.Time `json:"processed_at,omitempty"`
	}

	OrderStatusResponse struct {
		OrderNumber string `json:"order"`
		Status      string `json:"status"`
		Accrual     int64  `json:"accrual,omitempty"`
	}
)

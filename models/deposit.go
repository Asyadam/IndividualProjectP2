package models

import "time"

type DepositHistory struct {
	ID            int       `json:"id"`
	UserID        int       `json:"user_id"`
	Amount        int       `json:"amount"`
	PaymentLink   string    `json:"payment_link"`
	PaymentStatus string    `json:"payment_status"`
	ExternalID    string    `json:"external_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

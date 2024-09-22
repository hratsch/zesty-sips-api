package models

import (
	"time"
)

type LoyaltyPoints struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Points    int       `json:"points"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LoyaltyTransaction struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	OrderID   int64     `json:"order_id"`
	Points    int       `json:"points"`
	Type      string    `json:"type"` // "earn" or "redeem"
	CreatedAt time.Time `json:"created_at"`
}

package models

import (
	"time"
)

type Order struct {
	ID              int64       `json:"id"`
	UserID          int64       `json:"user_id"`
	TotalAmount     float64     `json:"total_amount"`
	Status          string      `json:"status"`
	OrderType       string      `json:"order_type"`
	DeliveryAddress string      `json:"delivery_address,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	Items           []OrderItem `json:"items"`
}

type OrderItem struct {
	ID        int64   `json:"id"`
	OrderID   int64   `json:"order_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

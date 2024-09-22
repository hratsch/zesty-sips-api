package models

import (
	"time"
)

type Product struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Size          string    `json:"size"`
	Price         float64   `json:"price"`
	StockQuantity int       `json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

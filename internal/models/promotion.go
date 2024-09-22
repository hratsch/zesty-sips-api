package models

import (
	"time"
)

type Promotion struct {
	ID              int64     `json:"id"`
	Code            string    `json:"code"`
	Description     string    `json:"description"`
	DiscountPercent float64   `json:"discount_percent"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

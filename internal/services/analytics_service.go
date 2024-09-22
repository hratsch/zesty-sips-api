package services

import (
	"database/sql"
	"time"
)

type AnalyticsService struct {
	DB *sql.DB
}

func NewAnalyticsService(db *sql.DB) *AnalyticsService {
	return &AnalyticsService{DB: db}
}

type SalesReport struct {
	TotalSales        float64 `json:"total_sales"`
	OrderCount        int     `json:"order_count"`
	AverageOrderValue float64 `json:"average_order_value"`
}

func (s *AnalyticsService) GetSalesReport(startDate, endDate time.Time) (*SalesReport, error) {
	query := `
		SELECT 
			COALESCE(SUM(total_amount), 0) as total_sales,
			COUNT(*) as order_count
		FROM orders
		WHERE created_at BETWEEN $1 AND $2
	`
	var report SalesReport
	err := s.DB.QueryRow(query, startDate, endDate).Scan(&report.TotalSales, &report.OrderCount)
	if err != nil {
		return nil, err
	}

	if report.OrderCount > 0 {
		report.AverageOrderValue = report.TotalSales / float64(report.OrderCount)
	}

	return &report, nil
}

type TopProduct struct {
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	TotalSales  float64 `json:"total_sales"`
	Quantity    int     `json:"quantity"`
}

func (s *AnalyticsService) GetTopProducts(limit int) ([]TopProduct, error) {
	query := `
		SELECT 
			p.id, 
			p.name, 
			SUM(oi.quantity * oi.unit_price) as total_sales,
			SUM(oi.quantity) as quantity
		FROM order_items oi
		JOIN products p ON oi.product_id = p.id
		GROUP BY p.id, p.name
		ORDER BY total_sales DESC
		LIMIT $1
	`
	rows, err := s.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []TopProduct
	for rows.Next() {
		var p TopProduct
		if err := rows.Scan(&p.ProductID, &p.ProductName, &p.TotalSales, &p.Quantity); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

type LoyaltyStats struct {
	TotalPointsEarned   int `json:"total_points_earned"`
	TotalPointsRedeemed int `json:"total_points_redeemed"`
	ActiveUsers         int `json:"active_users"`
}

func (s *AnalyticsService) GetLoyaltyStats() (*LoyaltyStats, error) {
	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN type = 'earn' THEN points ELSE 0 END), 0) as points_earned,
			COALESCE(SUM(CASE WHEN type = 'redeem' THEN points ELSE 0 END), 0) as points_redeemed,
			COUNT(DISTINCT user_id) as active_users
		FROM loyalty_transactions
	`
	var stats LoyaltyStats
	err := s.DB.QueryRow(query).Scan(&stats.TotalPointsEarned, &stats.TotalPointsRedeemed, &stats.ActiveUsers)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

package services

import (
	"database/sql"
	"errors"

	"github.com/hratsch/zesty-sips-api/internal/models"
)

type LoyaltyService struct {
	DB *sql.DB
}

func NewLoyaltyService(db *sql.DB) *LoyaltyService {
	return &LoyaltyService{DB: db}
}

func (s *LoyaltyService) GetLoyaltyPoints(userID int64) (int, error) {
	var points int
	query := `SELECT points FROM loyalty_points WHERE user_id = $1`
	err := s.DB.QueryRow(query, userID).Scan(&points)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return points, err
}

func (s *LoyaltyService) AddLoyaltyPoints(userID, orderID int64, points int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update or insert loyalty points
	query := `
		INSERT INTO loyalty_points (user_id, points)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE
		SET points = loyalty_points.points + $2,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err = tx.Exec(query, userID, points)
	if err != nil {
		return err
	}

	// Record loyalty transaction
	query = `INSERT INTO loyalty_transactions (user_id, order_id, points, type) VALUES ($1, $2, $3, 'earn')`
	_, err = tx.Exec(query, userID, orderID, points)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *LoyaltyService) RedeemLoyaltyPoints(userID, orderID int64, points int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if user has enough points
	var currentPoints int
	query := `SELECT points FROM loyalty_points WHERE user_id = $1 FOR UPDATE`
	err = tx.QueryRow(query, userID).Scan(&currentPoints)
	if err != nil {
		return err
	}

	if currentPoints < points {
		return errors.New("insufficient loyalty points")
	}

	// Update loyalty points
	query = `
		UPDATE loyalty_points
		SET points = points - $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1
	`
	_, err = tx.Exec(query, userID, points)
	if err != nil {
		return err
	}

	// Record loyalty transaction
	query = `INSERT INTO loyalty_transactions (user_id, order_id, points, type) VALUES ($1, $2, $3, 'redeem')`
	_, err = tx.Exec(query, userID, orderID, points)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *LoyaltyService) GetLoyaltyTransactions(userID int64) ([]models.LoyaltyTransaction, error) {
	query := `
		SELECT id, user_id, order_id, points, type, created_at
		FROM loyalty_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.LoyaltyTransaction
	for rows.Next() {
		var t models.LoyaltyTransaction
		err := rows.Scan(&t.ID, &t.UserID, &t.OrderID, &t.Points, &t.Type, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}

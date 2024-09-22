package services

import (
	"database/sql"
	"errors"
	"time"

	"github.com/hratsch/zesty-sips-api/internal/models"
)

type PromotionService struct {
	DB *sql.DB
}

func NewPromotionService(db *sql.DB) *PromotionService {
	return &PromotionService{DB: db}
}

func (s *PromotionService) CreatePromotion(promotion *models.Promotion) error {
	query := `
		INSERT INTO promotions (code, description, discount_percent, start_date, end_date, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`
	err := s.DB.QueryRow(
		query,
		promotion.Code,
		promotion.Description,
		promotion.DiscountPercent,
		promotion.StartDate,
		promotion.EndDate,
		promotion.IsActive,
	).Scan(&promotion.ID, &promotion.CreatedAt, &promotion.UpdatedAt)

	return err
}

func (s *PromotionService) GetPromotion(id int64) (*models.Promotion, error) {
	promotion := &models.Promotion{}
	query := `
		SELECT id, code, description, discount_percent, start_date, end_date, is_active, created_at, updated_at
		FROM promotions
		WHERE id = $1
	`
	err := s.DB.QueryRow(query, id).Scan(
		&promotion.ID,
		&promotion.Code,
		&promotion.Description,
		&promotion.DiscountPercent,
		&promotion.StartDate,
		&promotion.EndDate,
		&promotion.IsActive,
		&promotion.CreatedAt,
		&promotion.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("promotion not found")
		}
		return nil, err
	}

	return promotion, nil
}

func (s *PromotionService) ListActivePromotions() ([]*models.Promotion, error) {
	query := `
		SELECT id, code, description, discount_percent, start_date, end_date, is_active, created_at, updated_at
		FROM promotions
		WHERE is_active = true AND start_date <= $1 AND end_date >= $1
		ORDER BY created_at DESC
	`
	rows, err := s.DB.Query(query, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var promotions []*models.Promotion
	for rows.Next() {
		promotion := &models.Promotion{}
		err := rows.Scan(
			&promotion.ID,
			&promotion.Code,
			&promotion.Description,
			&promotion.DiscountPercent,
			&promotion.StartDate,
			&promotion.EndDate,
			&promotion.IsActive,
			&promotion.CreatedAt,
			&promotion.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		promotions = append(promotions, promotion)
	}

	return promotions, nil
}

func (s *PromotionService) UpdatePromotion(promotion *models.Promotion) error {
	query := `
		UPDATE promotions
		SET code = $1, description = $2, discount_percent = $3, start_date = $4, end_date = $5, is_active = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
		RETURNING updated_at
	`
	err := s.DB.QueryRow(
		query,
		promotion.Code,
		promotion.Description,
		promotion.DiscountPercent,
		promotion.StartDate,
		promotion.EndDate,
		promotion.IsActive,
		promotion.ID,
	).Scan(&promotion.UpdatedAt)

	return err
}

func (s *PromotionService) DeletePromotion(id int64) error {
	query := `DELETE FROM promotions WHERE id = $1`
	_, err := s.DB.Exec(query, id)
	return err
}

func (s *PromotionService) ApplyPromotion(code string, totalAmount float64) (float64, error) {
	query := `
		SELECT discount_percent
		FROM promotions
		WHERE code = $1 AND is_active = true AND start_date <= $2 AND end_date >= $2
	`
	var discountPercent float64
	err := s.DB.QueryRow(query, code, time.Now()).Scan(&discountPercent)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("invalid or expired promotion code")
		}
		return 0, err
	}

	discountAmount := totalAmount * (discountPercent / 100)
	return discountAmount, nil
}

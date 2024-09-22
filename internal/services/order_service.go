package services

import (
	"database/sql"
	"errors"

	"github.com/hratsch/zesty-sips-api/internal/models"
)

type OrderService struct {
	DB               *sql.DB
	ProductService   *ProductService
	LoyaltyService   *LoyaltyService
	PromotionService *PromotionService
}

func NewOrderService(db *sql.DB, productService *ProductService, loyaltyService *LoyaltyService, promotionService *PromotionService) *OrderService {
	return &OrderService{
		DB:               db,
		ProductService:   productService,
		LoyaltyService:   loyaltyService,
		PromotionService: promotionService,
	}
}

func (s *OrderService) CreateOrder(order *models.Order, promotionCode string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check and update stock for each item
	for _, item := range order.Items {
		err := s.ProductService.UpdateStock(tx, item.ProductID, item.Quantity)
		if err != nil {
			return err
		}
	}

	// Apply promotion if a code is provided
	if promotionCode != "" {
		discountAmount, err := s.PromotionService.ApplyPromotion(promotionCode, order.TotalAmount)
		if err != nil {
			return err
		}
		order.TotalAmount -= discountAmount
	}

	// Insert order
	query := `INSERT INTO orders (user_id, total_amount, status, order_type, delivery_address) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`

	err = tx.QueryRow(query, order.UserID, order.TotalAmount, order.Status, order.OrderType, order.DeliveryAddress).
		Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return err
	}

	// Insert order items
	for i := range order.Items {
		query := `INSERT INTO order_items (order_id, product_id, quantity, unit_price) 
                  VALUES ($1, $2, $3, $4) RETURNING id`
		err = tx.QueryRow(query, order.ID, order.Items[i].ProductID, order.Items[i].Quantity, order.Items[i].UnitPrice).
			Scan(&order.Items[i].ID)
		if err != nil {
			return err
		}
	}

	// Calculate and add loyalty points (1 point per $1 spent)
	loyaltyPoints := int(order.TotalAmount)
	err = s.LoyaltyService.AddLoyaltyPoints(order.UserID, order.ID, loyaltyPoints)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *OrderService) GetOrder(id int64) (*models.Order, error) {
	order := &models.Order{}
	query := `SELECT id, user_id, total_amount, status, order_type, delivery_address, created_at, updated_at 
              FROM orders WHERE id = $1`

	err := s.DB.QueryRow(query, id).Scan(
		&order.ID, &order.UserID, &order.TotalAmount, &order.Status,
		&order.OrderType, &order.DeliveryAddress, &order.CreatedAt, &order.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	// Get order items
	itemsQuery := `SELECT id, product_id, quantity, unit_price FROM order_items WHERE order_id = $1`
	rows, err := s.DB.Query(itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.UnitPrice)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (s *OrderService) ListOrders(userID int64) ([]*models.Order, error) {
	query := `SELECT id, user_id, total_amount, status, order_type, delivery_address, created_at, updated_at 
              FROM orders WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := s.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		order := &models.Order{}
		err := rows.Scan(
			&order.ID, &order.UserID, &order.TotalAmount, &order.Status,
			&order.OrderType, &order.DeliveryAddress, &order.CreatedAt, &order.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (s *OrderService) UpdateOrderStatus(id int64, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := s.DB.Exec(query, status, id)
	return err
}

func (s *OrderService) CancelOrder(orderID int64) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get order items
	query := `SELECT product_id, quantity FROM order_items WHERE order_id = $1`
	rows, err := tx.Query(query, orderID)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Restock items
	for rows.Next() {
		var productID int64
		var quantity int
		err := rows.Scan(&productID, &quantity)
		if err != nil {
			return err
		}
		err = s.ProductService.RestockProduct(productID, quantity)
		if err != nil {
			return err
		}
	}

	// Update order status to cancelled
	query = `UPDATE orders SET status = 'cancelled', updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err = tx.Exec(query, orderID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

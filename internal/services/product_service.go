package services

import (
	"database/sql"
	"errors"

	"github.com/hratsch/zesty-sips-api/internal/models"
)

type ProductService struct {
	DB *sql.DB
}

func NewProductService(db *sql.DB) *ProductService {
	return &ProductService{DB: db}
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	query := `INSERT INTO products (name, description, size, price, stock_quantity) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`

	err := s.DB.QueryRow(query, product.Name, product.Description, product.Size, product.Price, product.StockQuantity).
		Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	return err
}

func (s *ProductService) GetProduct(id int64) (*models.Product, error) {
	product := &models.Product{}
	query := `SELECT id, name, description, size, price, stock_quantity, created_at, updated_at 
              FROM products WHERE id = $1`

	err := s.DB.QueryRow(query, id).Scan(
		&product.ID, &product.Name, &product.Description, &product.Size,
		&product.Price, &product.StockQuantity, &product.CreatedAt, &product.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return product, nil
}

func (s *ProductService) ListProducts() ([]*models.Product, error) {
	query := `SELECT id, name, description, size, price, stock_quantity, created_at, updated_at 
              FROM products ORDER BY name`

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		product := &models.Product{}
		err := rows.Scan(
			&product.ID, &product.Name, &product.Description, &product.Size,
			&product.Price, &product.StockQuantity, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (s *ProductService) UpdateProduct(product *models.Product) error {
	query := `UPDATE products SET name = $1, description = $2, size = $3, price = $4, 
              stock_quantity = $5, updated_at = CURRENT_TIMESTAMP 
              WHERE id = $6 RETURNING updated_at`

	err := s.DB.QueryRow(query, product.Name, product.Description, product.Size,
		product.Price, product.StockQuantity, product.ID).Scan(&product.UpdatedAt)

	return err
}

func (s *ProductService) DeleteProduct(id int64) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := s.DB.Exec(query, id)
	return err
}

func (s *ProductService) UpdateStock(tx *sql.Tx, productID int64, quantity int) error {
	query := `
		UPDATE products
		SET stock_quantity = stock_quantity - $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND stock_quantity >= $1
		RETURNING stock_quantity
	`
	var newStockQuantity int
	err := tx.QueryRow(query, quantity, productID).Scan(&newStockQuantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("insufficient stock")
		}
		return err
	}
	return nil
}

func (s *ProductService) GetStockQuantity(productID int64) (int, error) {
	query := `SELECT stock_quantity FROM products WHERE id = $1`
	var stockQuantity int
	err := s.DB.QueryRow(query, productID).Scan(&stockQuantity)
	if err != nil {
		return 0, err
	}
	return stockQuantity, nil
}

func (s *ProductService) RestockProduct(productID int64, quantity int) error {
	query := `
		UPDATE products
		SET stock_quantity = stock_quantity + $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING stock_quantity
	`
	var newStockQuantity int
	err := s.DB.QueryRow(query, quantity, productID).Scan(&newStockQuantity)
	if err != nil {
		return err
	}
	return nil
}

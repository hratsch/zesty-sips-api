package services

import (
	"database/sql"
	"errors"

	"github.com/hratsch/zesty-sips-api/internal/models"
)

type UserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{DB: db}
}

func (s *UserService) CreateUser(user *models.User) error {
	if err := validateEmail(user.Email); err != nil {
		log.WithFields(log.Fields{
			"email": user.Email,
		}).Errorf("Invalid email: %v", err)
		return fmt.Errorf("invalid email: %w", err)
	}
	if err := validatePassword(user.Password); err != nil {
		log.WithFields(log.Fields{
			"password": user.Password,
		}).Errorf("Invalid password: %v", err)
		return fmt.Errorf("invalid password: %w", err)
	}

	if err := user.HashPassword(); err != nil {
		return err
	}

	query := `INSERT INTO users (email, password, first_name, last_name, phone, role) 
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at, updated_at`

	err := s.DB.QueryRow(query, user.Email, user.Password, user.FirstName, user.LastName, user.Phone, user.Role).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		log.WithFields(log.Fields{
			"email": user.Email,
		}).Errorf("Failed to create user: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

var userCache = make(map[string]*models.User)

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	if user, ok := userCache[email]; ok {
		return user, nil
	}

	user := &models.User{}
	query := `SELECT id, email, password, first_name, last_name, phone, role, created_at, updated_at 
              FROM users WHERE email = $1`

	err := s.DB.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
		&user.Phone, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	userCache[email] = user
	return user, nil
}

func (s *UserService) GetUser(userID int64) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, first_name, last_name, phone, role, created_at, updated_at 
              FROM users WHERE id = $1`

	err := s.DB.QueryRow(query, userID).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName,
		&user.Phone, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

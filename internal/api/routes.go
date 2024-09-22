package api

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/hratsch/zesty-sips-api/internal/api/handlers"
	"github.com/hratsch/zesty-sips-api/internal/api/middleware"
	"github.com/hratsch/zesty-sips-api/internal/services"
)

func NewRouter(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	// Middleware
	r.Use(middleware.Logging)

	// Services
	userService := services.NewUserService(db)
	productService := services.NewProductService(db)
	loyaltyService := services.NewLoyaltyService(db)
	promotionService := services.NewPromotionService(db)
	orderService := services.NewOrderService(db, productService, loyaltyService, promotionService)
	analyticsService := services.NewAnalyticsService(db)

	// Handlers
	userHandler := handlers.NewUserHandler(userService)
	productHandler := handlers.NewProductHandler(productService)
	orderHandler := handlers.NewOrderHandler(orderService)
	loyaltyHandler := handlers.NewLoyaltyHandler(loyaltyService)
	promotionHandler := handlers.NewPromotionHandler(promotionService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)

	// Public routes
	r.HandleFunc("/health", handlers.HealthCheck).Methods("GET")
	r.HandleFunc("/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")

	// Protected routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(middleware.Auth)

	// User routes
	api.HandleFunc("/users/{id}", userHandler.GetUser).Methods("GET")

	// Product routes
	api.HandleFunc("/products", productHandler.ListProducts).Methods("GET")
	api.HandleFunc("/products", productHandler.CreateProduct).Methods("POST")
	api.HandleFunc("/products/{id}", productHandler.GetProduct).Methods("GET")
	api.HandleFunc("/products/{id}", productHandler.UpdateProduct).Methods("PUT")
	api.HandleFunc("/products/{id}", productHandler.DeleteProduct).Methods("DELETE")

	// Order routes
	api.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST")
	api.HandleFunc("/orders", orderHandler.ListOrders).Methods("GET")
	api.HandleFunc("/orders/{id}", orderHandler.GetOrder).Methods("GET")
	api.HandleFunc("/orders/{id}/status", orderHandler.UpdateOrderStatus).Methods("PATCH")
	api.HandleFunc("/orders/{id}/cancel", orderHandler.CancelOrder).Methods("POST")

	// Loyalty routes
	api.HandleFunc("/loyalty/points", loyaltyHandler.GetLoyaltyPoints).Methods("GET")
	api.HandleFunc("/loyalty/transactions", loyaltyHandler.GetLoyaltyTransactions).Methods("GET")
	api.HandleFunc("/loyalty/redeem", loyaltyHandler.RedeemPoints).Methods("POST")

	// Promotion routes
	api.HandleFunc("/promotions", promotionHandler.CreatePromotion).Methods("POST")
	api.HandleFunc("/promotions", promotionHandler.ListActivePromotions).Methods("GET")
	api.HandleFunc("/promotions/{id}", promotionHandler.GetPromotion).Methods("GET")
	api.HandleFunc("/promotions/{id}", promotionHandler.UpdatePromotion).Methods("PUT")
	api.HandleFunc("/promotions/{id}", promotionHandler.DeletePromotion).Methods("DELETE")
	api.HandleFunc("/promotions/apply", promotionHandler.ApplyPromotion).Methods("POST")

	// Analytics routes
	api.HandleFunc("/analytics/sales", analyticsHandler.GetSalesReport).Methods("GET")
	api.HandleFunc("/analytics/top-products", analyticsHandler.GetTopProducts).Methods("GET")
	api.HandleFunc("/analytics/loyalty", analyticsHandler.GetLoyaltyStats).Methods("GET")

	return r
}

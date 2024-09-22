package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hratsch/zesty-sips-api/internal/services"
)

type LoyaltyHandler struct {
	LoyaltyService *services.LoyaltyService
}

func NewLoyaltyHandler(loyaltyService *services.LoyaltyService) *LoyaltyHandler {
	return &LoyaltyHandler{LoyaltyService: loyaltyService}
}

func (h *LoyaltyHandler) GetLoyaltyPoints(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)

	points, err := h.LoyaltyService.GetLoyaltyPoints(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"points": points})
}

func (h *LoyaltyHandler) GetLoyaltyTransactions(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)

	transactions, err := h.LoyaltyService.GetLoyaltyTransactions(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transactions)
}

func (h *LoyaltyHandler) RedeemPoints(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int64)

	var redeemRequest struct {
		OrderID int64 `json:"order_id"`
		Points  int   `json:"points"`
	}

	if err := json.NewDecoder(r.Body).Decode(&redeemRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.LoyaltyService.RedeemLoyaltyPoints(userID, redeemRequest.OrderID, redeemRequest.Points)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Points redeemed successfully"})
}

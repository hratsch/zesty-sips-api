package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hratsch/zesty-sips-api/internal/models"
	"github.com/hratsch/zesty-sips-api/internal/services"
)

type PromotionHandler struct {
	PromotionService *services.PromotionService
}

func NewPromotionHandler(promotionService *services.PromotionService) *PromotionHandler {
	return &PromotionHandler{PromotionService: promotionService}
}

func (h *PromotionHandler) CreatePromotion(w http.ResponseWriter, r *http.Request) {
	var promotion models.Promotion
	if err := json.NewDecoder(r.Body).Decode(&promotion); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.PromotionService.CreatePromotion(&promotion); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(promotion)
}

func (h *PromotionHandler) GetPromotion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid promotion ID", http.StatusBadRequest)
		return
	}

	promotion, err := h.PromotionService.GetPromotion(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(promotion)
}

func (h *PromotionHandler) ListActivePromotions(w http.ResponseWriter, r *http.Request) {
	promotions, err := h.PromotionService.ListActivePromotions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(promotions)
}

func (h *PromotionHandler) UpdatePromotion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid promotion ID", http.StatusBadRequest)
		return
	}

	var promotion models.Promotion
	if err := json.NewDecoder(r.Body).Decode(&promotion); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	promotion.ID = id

	if err := h.PromotionService.UpdatePromotion(&promotion); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(promotion)
}

func (h *PromotionHandler) DeletePromotion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid promotion ID", http.StatusBadRequest)
		return
	}

	if err := h.PromotionService.DeletePromotion(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *PromotionHandler) ApplyPromotion(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Code        string  `json:"code"`
		TotalAmount float64 `json:"total_amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	discountAmount, err := h.PromotionService.ApplyPromotion(request.Code, request.TotalAmount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := struct {
		DiscountAmount float64 `json:"discount_amount"`
		FinalAmount    float64 `json:"final_amount"`
	}{
		DiscountAmount: discountAmount,
		FinalAmount:    request.TotalAmount - discountAmount,
	}

	json.NewEncoder(w).Encode(response)
}

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/hratsch/zesty-sips-api/internal/services"
)

type AnalyticsHandler struct {
	AnalyticsService *services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{AnalyticsService: analyticsService}
}

func (h *AnalyticsHandler) GetSalesReport(w http.ResponseWriter, r *http.Request) {
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		http.Error(w, "Invalid start date format", http.StatusBadRequest)
		return
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		http.Error(w, "Invalid end date format", http.StatusBadRequest)
		return
	}

	report, err := h.AnalyticsService.GetSalesReport(start, end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(report)
}

func (h *AnalyticsHandler) GetTopProducts(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // Default to top 10 if not specified or invalid
	}

	products, err := h.AnalyticsService.GetTopProducts(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}

func (h *AnalyticsHandler) GetLoyaltyStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.AnalyticsService.GetLoyaltyStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(stats)
}

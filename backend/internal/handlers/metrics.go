package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Jdavid77/kubernetes-dashboard/internal/kubernetes"
)

// MetricsHandler handles metrics requests
type MetricsHandler struct {
	aggregator kubernetes.MetricsSource
}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler(aggregator kubernetes.MetricsSource) *MetricsHandler {
	return &MetricsHandler{
		aggregator: aggregator,
	}
}

// ServeHTTP handles metrics requests
func (h *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	metrics, err := h.aggregator.GetClusterMetrics(ctx)
	if err != nil {
		log.Printf("Error getting cluster metrics: %v", err)
		http.Error(w, "Failed to get cluster metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(metrics)
}

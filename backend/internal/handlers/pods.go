package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Jdavid77/kubernetes-dashboard/internal/kubernetes"
)

// PodsHandler handles pod-related HTTP requests
type PodsHandler struct {
	aggregator *kubernetes.Aggregator
}

// NewPodsHandler creates a new pods handler
func NewPodsHandler(aggregator *kubernetes.Aggregator) *PodsHandler {
	return &PodsHandler{aggregator: aggregator}
}

// ListPods handles GET /api/pods
func (h *PodsHandler) ListPods(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get namespace filter from query parameter
	namespace := r.URL.Query().Get("namespace")

	pods, err := h.aggregator.GetPods(ctx)
	if err != nil {
		log.Printf("Error getting pods: %v", err)
		http.Error(w, "Failed to get pods", http.StatusInternalServerError)
		return
	}

	// Filter by namespace if specified
	if namespace != "" && namespace != "all" {
		filtered := make([]interface{}, 0)
		for _, pod := range pods {
			if pod.Namespace == namespace {
				filtered = append(filtered, pod)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(filtered)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pods)
}

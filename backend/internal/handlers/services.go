package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Jdavid77/kubernetes-dashboard/internal/kubernetes"
)

// ServicesHandler handles services requests
type ServicesHandler struct {
	aggregator kubernetes.MetricsSource
}

// NewServicesHandler creates a new services handler
func NewServicesHandler(aggregator kubernetes.MetricsSource) *ServicesHandler {
	return &ServicesHandler{
		aggregator: aggregator,
	}
}

// ListServices handles listing all services
func (h *ServicesHandler) ListServices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get namespace filter from query parameter
	namespace := r.URL.Query().Get("namespace")

	services, err := h.aggregator.GetServices(ctx)
	if err != nil {
		log.Printf("Error getting services: %v", err)
		http.Error(w, "Failed to get services", http.StatusInternalServerError)
		return
	}

	// Filter by namespace if specified
	if namespace != "" && namespace != "all" {
		filtered := make([]interface{}, 0)
		for _, svc := range services {
			if svc.Namespace == namespace {
				filtered = append(filtered, svc)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(filtered); err != nil {
			log.Printf("Error encoding filtered services: %v", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(services); err != nil {
		log.Printf("Error encoding services: %v", err)
	}
}

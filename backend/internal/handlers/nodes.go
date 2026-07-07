package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Jdavid77/kubernetes-dashboard/internal/kubernetes"
)

// NodesHandler handles node-related HTTP requests
type NodesHandler struct {
	aggregator kubernetes.MetricsSource
}

// NewNodesHandler creates a new nodes handler
func NewNodesHandler(aggregator kubernetes.MetricsSource) *NodesHandler {
	return &NodesHandler{aggregator: aggregator}
}

// ListNodes handles GET /api/nodes
func (h *NodesHandler) ListNodes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	nodes, err := h.aggregator.GetNodes(ctx)
	if err != nil {
		log.Printf("Error getting nodes: %v", err)
		http.Error(w, "Failed to get nodes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(nodes)
}

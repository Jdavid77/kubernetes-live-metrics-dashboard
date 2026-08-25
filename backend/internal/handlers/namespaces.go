package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Jdavid77/kubernetes-dashboard/internal/kubernetes"
)

// NamespacesHandler handles namespace requests
type NamespacesHandler struct {
	source kubernetes.MetricsSource
}

// NewNamespacesHandler creates a new namespaces handler
func NewNamespacesHandler(source kubernetes.MetricsSource) *NamespacesHandler {
	return &NamespacesHandler{source: source}
}

// ListNamespaces handles listing all namespaces
func (h *NamespacesHandler) ListNamespaces(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	namespaces, err := h.source.GetNamespaces(ctx)
	if err != nil {
		log.Printf("Error getting namespaces: %v", err)
		http.Error(w, "Failed to get namespaces", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(namespaces); err != nil {
		log.Printf("Error encoding namespaces: %v", err)
	}
}

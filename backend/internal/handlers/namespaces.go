package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Jdavid77/kubernetes-dashboard/internal/kubernetes"
)

// NamespacesHandler handles namespace requests
type NamespacesHandler struct {
	client *kubernetes.Client
}

// NewNamespacesHandler creates a new namespaces handler
func NewNamespacesHandler(client *kubernetes.Client) *NamespacesHandler {
	return &NamespacesHandler{
		client: client,
	}
}

// ListNamespaces handles listing all namespaces
func (h *NamespacesHandler) ListNamespaces(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	namespaces, err := h.client.GetNamespaces(ctx)
	if err != nil {
		log.Printf("Error getting namespaces: %v", err)
		http.Error(w, "Failed to get namespaces", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(namespaces)
}

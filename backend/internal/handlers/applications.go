package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Jdavid77/kubernetes-dashboard/internal/kubernetes"
	"github.com/gorilla/mux"
)

// ApplicationsHandler handles application requests
type ApplicationsHandler struct {
	aggregator kubernetes.MetricsSource
}

// NewApplicationsHandler creates a new applications handler
func NewApplicationsHandler(aggregator kubernetes.MetricsSource) *ApplicationsHandler {
	return &ApplicationsHandler{
		aggregator: aggregator,
	}
}

// ListApplications handles listing all applications
func (h *ApplicationsHandler) ListApplications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Get namespace filter from query parameter
	namespace := r.URL.Query().Get("namespace")

	apps, err := h.aggregator.GetApplications(ctx)
	if err != nil {
		log.Printf("Error getting applications: %v", err)
		http.Error(w, "Failed to get applications", http.StatusInternalServerError)
		return
	}

	// Filter by namespace if specified
	if namespace != "" && namespace != "all" {
		filtered := make([]interface{}, 0)
		for _, app := range apps {
			if app.Namespace == namespace {
				filtered = append(filtered, app)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(filtered); err != nil {
			log.Printf("Error encoding filtered applications: %v", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(apps); err != nil {
		log.Printf("Error encoding applications: %v", err)
	}
}

// GetApplicationDetail handles getting detailed information about an application
func (h *ApplicationsHandler) GetApplicationDetail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	namespace := r.URL.Query().Get("namespace")
	name := vars["name"]

	if namespace == "" {
		http.Error(w, "namespace parameter is required", http.StatusBadRequest)
		return
	}

	if name == "" {
		http.Error(w, "name parameter is required", http.StatusBadRequest)
		return
	}

	detail, err := h.aggregator.GetApplicationDetail(ctx, namespace, name)
	if err != nil {
		log.Printf("Error getting application detail: %v", err)
		http.Error(w, "Failed to get application detail", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(detail); err != nil {
		log.Printf("Error encoding application detail: %v", err)
	}
}

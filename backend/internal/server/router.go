package server

import (
	"github.com/Jdavid77/kubernetes-dashboard/internal/handlers"
	"github.com/Jdavid77/kubernetes-dashboard/internal/middleware"
	"github.com/gorilla/mux"
)

// SetupRouter sets up the HTTP router with all routes and middleware
func SetupRouter(
	healthHandler *handlers.HealthHandler,
	metricsHandler *handlers.MetricsHandler,
	appsHandler *handlers.ApplicationsHandler,
	podsHandler *handlers.PodsHandler,
	servicesHandler *handlers.ServicesHandler,
	namespacesHandler *handlers.NamespacesHandler,
	nodesHandler *handlers.NodesHandler,
	wsHandler *handlers.WebSocketHandler,
	corsOrigin string,
) *mux.Router {
	r := mux.NewRouter()

	// Apply middleware
	r.Use(middleware.Logger)
	r.Use(middleware.CORS(corsOrigin))

	// Health check endpoint
	r.Handle("/api/health", healthHandler).Methods("GET")

	// Metrics endpoints
	r.Handle("/api/metrics", metricsHandler).Methods("GET")

	// Applications endpoints
	r.HandleFunc("/api/applications", appsHandler.ListApplications).Methods("GET")
	r.HandleFunc("/api/applications/{name}", appsHandler.GetApplicationDetail).Methods("GET")

	// Pods endpoints
	r.HandleFunc("/api/pods", podsHandler.ListPods).Methods("GET")

	// Services endpoints
	r.HandleFunc("/api/services", servicesHandler.ListServices).Methods("GET")

	// Namespaces endpoints
	r.HandleFunc("/api/namespaces", namespacesHandler.ListNamespaces).Methods("GET")

	// Nodes endpoints
	r.HandleFunc("/api/nodes", nodesHandler.ListNodes).Methods("GET", "OPTIONS")

	// WebSocket endpoint
	r.Handle("/ws/metrics", wsHandler)

	return r
}

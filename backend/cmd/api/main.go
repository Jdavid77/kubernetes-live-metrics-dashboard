package main

import (
	"log"

	"github.com/Jdavid77/kubernetes-dashboard/internal/config"
	"github.com/Jdavid77/kubernetes-dashboard/internal/handlers"
	"github.com/Jdavid77/kubernetes-dashboard/internal/kubernetes"
	"github.com/Jdavid77/kubernetes-dashboard/internal/server"
)

func main() {
	log.Println("Starting Kubernetes Dashboard API...")

	// Load configuration
	cfg := config.Load()
	log.Printf("Configuration loaded: Port=%s, CORS=%s", cfg.Port, cfg.CORSOrigin)

	// Initialize Kubernetes client
	k8sClient, err := kubernetes.NewClient(cfg.KubeConfig)
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	// Create aggregator
	aggregator := kubernetes.NewAggregator(k8sClient, cfg.MetricsRefreshInterval)

	// Create handlers
	healthHandler := handlers.NewHealthHandler()
	metricsHandler := handlers.NewMetricsHandler(aggregator)
	appsHandler := handlers.NewApplicationsHandler(aggregator)
	podsHandler := handlers.NewPodsHandler(aggregator)
	servicesHandler := handlers.NewServicesHandler(aggregator)
	namespacesHandler := handlers.NewNamespacesHandler(k8sClient)
	nodesHandler := handlers.NewNodesHandler(aggregator)
	wsHandler := handlers.NewWebSocketHandler(aggregator, cfg.MetricsRefreshInterval)

	// Setup router
	router := server.SetupRouter(
		healthHandler,
		metricsHandler,
		appsHandler,
		podsHandler,
		servicesHandler,
		namespacesHandler,
		nodesHandler,
		wsHandler,
		cfg.CORSOrigin,
	)

	// Create and start server
	srv := server.NewServer(cfg.Port, router)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

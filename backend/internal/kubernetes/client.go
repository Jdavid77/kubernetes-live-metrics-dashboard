package kubernetes

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

// Client holds Kubernetes clients
type Client struct {
	Clientset       *kubernetes.Clientset
	MetricsClientset *versioned.Clientset
}

// NewClient creates a new Kubernetes client with auto-detection
func NewClient(kubeconfigPath string) (*Client, error) {
	config, err := getKubeConfig(kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	// Create standard clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	// Create metrics clientset
	metricsClientset, err := versioned.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics clientset: %w", err)
	}

	log.Println("Successfully connected to Kubernetes cluster")

	return &Client{
		Clientset:       clientset,
		MetricsClientset: metricsClientset,
	}, nil
}

// getKubeConfig auto-detects and returns the Kubernetes configuration
func getKubeConfig(kubeconfigPath string) (*rest.Config, error) {
	// Try in-cluster config first
	config, err := rest.InClusterConfig()
	if err == nil {
		log.Println("Using in-cluster Kubernetes configuration")
		return config, nil
	}

	log.Println("Not running in cluster, attempting to use kubeconfig file")

	// Fallback to kubeconfig file
	if kubeconfigPath == "" {
		// Try default kubeconfig location
		if home := homeDir(); home != "" {
			kubeconfigPath = filepath.Join(home, ".kube", "config")
		}
	}

	// Check if kubeconfig exists
	if _, err := os.Stat(kubeconfigPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("kubeconfig file not found at %s", kubeconfigPath)
	}

	// Build config from kubeconfig file
	config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build config from kubeconfig: %w", err)
	}

	log.Printf("Using kubeconfig file: %s\n", kubeconfigPath)
	return config, nil
}

// homeDir returns the home directory for the current user
func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

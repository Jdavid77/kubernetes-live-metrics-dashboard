package config

import (
	"os"
	"time"
)

// Config holds all application configuration
type Config struct {
	Port                   string
	KubeConfig             string
	CORSOrigin             string
	MetricsRefreshInterval time.Duration
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Port:                   getEnv("PORT", "8080"),
		KubeConfig:             getEnv("KUBECONFIG", ""),
		CORSOrigin:             getEnv("CORS_ORIGIN", "http://localhost:3000"),
		MetricsRefreshInterval: parseDuration(getEnv("METRICS_REFRESH_INTERVAL", "5s")),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseDuration parses a duration string or returns default
func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 5 * time.Second
	}
	return d
}

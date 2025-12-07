package models

import "time"

// ClusterMetrics represents the overall cluster metrics
type ClusterMetrics struct {
	CPUUsage       string    `json:"cpu_usage"`
	MemoryUsage    string    `json:"memory_usage"`
	PodsRunning    int       `json:"pods_running"`
	GatewayCount   int       `json:"gateway_count"`
	HTTPRouteCount int       `json:"httproute_count"`
	ServicesCount  int       `json:"services_count"`
	Timestamp      time.Time `json:"timestamp"`
}

// NodeMetrics represents metrics for a single node
type NodeMetrics struct {
	Name        string `json:"name"`
	CPUUsage    string `json:"cpu_usage"`
	MemoryUsage string `json:"memory_usage"`
}

// Node represents detailed node information with metrics
type Node struct {
	Name           string  `json:"name"`
	Status         string  `json:"status"` // "Ready" or "NotReady"
	Role           string  `json:"role"`   // "control-plane" or "worker"
	CPUUsage       float64 `json:"cpu_usage"`
	MemoryUsage    float64 `json:"memory_usage"`
	CPUCapacity    string  `json:"cpu_capacity"`
	MemoryCapacity string  `json:"memory_capacity"`
	PodCount       int     `json:"pod_count"`
}

// PodMetrics represents metrics for a single pod
type PodMetrics struct {
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	CPUUsage    string `json:"cpu_usage"`
	MemoryUsage string `json:"memory_usage"`
}

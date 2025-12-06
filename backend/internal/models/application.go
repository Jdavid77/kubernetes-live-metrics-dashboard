package models

import "time"

// Application represents a deployed application/service
type Application struct {
	Name       string            `json:"name"`
	Namespace  string            `json:"namespace"`
	Status     string            `json:"status"` // Healthy, Degraded, Down
	PodCount   int               `json:"pod_count"`
	Technology []string          `json:"technology"`
	Image      string            `json:"image"`
	CreatedAt  time.Time         `json:"created_at"`
	Labels     map[string]string `json:"labels"`
}

// ApplicationDetail represents detailed information about an application
type ApplicationDetail struct {
	Application
	Pods        []PodInfo         `json:"pods"`
	Environment map[string]string `json:"environment"`
	Resources   ResourceInfo      `json:"resources"`
}

// PodInfo represents information about a pod
type PodInfo struct {
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	Ready     string    `json:"ready"`
	Restarts  int32     `json:"restarts"`
	Age       string    `json:"age"`
	CreatedAt time.Time `json:"created_at"`
}

// ResourceInfo represents resource requests and limits
type ResourceInfo struct {
	CPURequest    string `json:"cpu_request"`
	CPULimit      string `json:"cpu_limit"`
	MemoryRequest string `json:"memory_request"`
	MemoryLimit   string `json:"memory_limit"`
}

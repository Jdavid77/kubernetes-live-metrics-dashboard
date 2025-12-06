package models

import "time"

// Service represents a Kubernetes service
type Service struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	Type        string            `json:"type"`
	ClusterIP   string            `json:"cluster_ip"`
	ExternalIP  string            `json:"external_ip"`
	Ports       []ServicePort     `json:"ports"`
	Selector    map[string]string `json:"selector"`
	CreatedAt   time.Time         `json:"created_at"`
}

// ServicePort represents a service port
type ServicePort struct {
	Name       string `json:"name"`
	Protocol   string `json:"protocol"`
	Port       int32  `json:"port"`
	TargetPort string `json:"target_port"`
	NodePort   int32  `json:"node_port,omitempty"`
}

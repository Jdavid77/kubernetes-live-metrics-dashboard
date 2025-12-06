package models

import "time"

// Pod represents a Kubernetes pod
type Pod struct {
	Name      string    `json:"name"`
	Namespace string    `json:"namespace"`
	Status    string    `json:"status"`
	Ready     string    `json:"ready"`
	Restarts  int32     `json:"restarts"`
	Age       string    `json:"age"`
	Node      string    `json:"node"`
	IP        string    `json:"ip"`
	CreatedAt time.Time `json:"created_at"`
}

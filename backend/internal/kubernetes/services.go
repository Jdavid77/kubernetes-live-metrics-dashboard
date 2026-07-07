package kubernetes

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetServices retrieves all services across all namespaces
func (c *Client) GetServices(ctx context.Context) (*corev1.ServiceList, error) {
	services, err := c.Clientset.CoreV1().Services("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}
	return services, nil
}


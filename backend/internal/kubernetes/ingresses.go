package kubernetes

import (
	"context"
	"fmt"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetIngresses retrieves all ingresses across all namespaces
func (c *Client) GetIngresses(ctx context.Context) (*networkingv1.IngressList, error) {
	ingresses, err := c.Clientset.NetworkingV1().Ingresses("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list ingresses: %w", err)
	}
	return ingresses, nil
}

// GetIngressesByNamespace retrieves ingresses in a specific namespace
func (c *Client) GetIngressesByNamespace(ctx context.Context, namespace string) (*networkingv1.IngressList, error) {
	ingresses, err := c.Clientset.NetworkingV1().Ingresses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list ingresses in namespace %s: %w", namespace, err)
	}
	return ingresses, nil
}

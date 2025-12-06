package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetGateways retrieves all Gateway API gateways across all namespaces
func (c *Client) GetGateways(ctx context.Context) (int, error) {
	// Use dynamic client to query Gateway API
	list, err := c.Clientset.Discovery().ServerResourcesForGroupVersion("gateway.networking.k8s.io/v1")
	if err != nil {
		// Gateway API not installed, return 0
		return 0, nil
	}

	// Check if gateways resource exists
	found := false
	for _, resource := range list.APIResources {
		if resource.Name == "gateways" {
			found = true
			break
		}
	}

	if !found {
		return 0, nil
	}

	// Try to list gateways using REST client
	result := c.Clientset.RESTClient().
		Get().
		AbsPath("/apis/gateway.networking.k8s.io/v1/gateways").
		Do(ctx)

	if result.Error() != nil {
		// If error, Gateway API might not be installed
		return 0, nil
	}

	var list2 metav1.List
	if err := result.Into(&list2); err != nil {
		return 0, nil
	}

	return len(list2.Items), nil
}

// GetHTTPRoutes retrieves all HTTPRoutes across all namespaces
func (c *Client) GetHTTPRoutes(ctx context.Context) (int, error) {
	// Try to list HTTPRoutes
	result := c.Clientset.RESTClient().
		Get().
		AbsPath("/apis/gateway.networking.k8s.io/v1/httproutes").
		Do(ctx)

	if result.Error() != nil {
		return 0, nil
	}

	var list metav1.List
	if err := result.Into(&list); err != nil {
		return 0, nil
	}

	return len(list.Items), nil
}

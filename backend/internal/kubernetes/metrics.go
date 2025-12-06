package kubernetes

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

// GetNodeMetrics retrieves metrics for all nodes
func (c *Client) GetNodeMetrics(ctx context.Context) (*metricsv1beta1.NodeMetricsList, error) {
	metrics, err := c.MetricsClientset.MetricsV1beta1().NodeMetricses().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get node metrics: %w", err)
	}
	return metrics, nil
}

// GetPodMetrics retrieves metrics for all pods across all namespaces
func (c *Client) GetPodMetrics(ctx context.Context) (*metricsv1beta1.PodMetricsList, error) {
	metrics, err := c.MetricsClientset.MetricsV1beta1().PodMetricses("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod metrics: %w", err)
	}
	return metrics, nil
}

// GetPodMetricsByNamespace retrieves metrics for pods in a specific namespace
func (c *Client) GetPodMetricsByNamespace(ctx context.Context, namespace string) (*metricsv1beta1.PodMetricsList, error) {
	metrics, err := c.MetricsClientset.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod metrics in namespace %s: %w", namespace, err)
	}
	return metrics, nil
}

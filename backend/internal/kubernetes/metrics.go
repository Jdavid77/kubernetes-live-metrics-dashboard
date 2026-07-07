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


package kubernetes

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetPods retrieves all pods across all namespaces
func (c *Client) GetPods(ctx context.Context) (*corev1.PodList, error) {
	pods, err := c.Clientset.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}
	return pods, nil
}

// GetPodsByNamespace retrieves pods in a specific namespace
func (c *Client) GetPodsByNamespace(ctx context.Context, namespace string) (*corev1.PodList, error) {
	pods, err := c.Clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods in namespace %s: %w", namespace, err)
	}
	return pods, nil
}

// CountPodsByStatus counts pods by their status
func CountPodsByStatus(pods *corev1.PodList) (running, pending, failed int) {
	for _, pod := range pods.Items {
		switch pod.Status.Phase {
		case corev1.PodRunning:
			running++
		case corev1.PodPending:
			pending++
		case corev1.PodFailed:
			failed++
		}
	}
	return
}

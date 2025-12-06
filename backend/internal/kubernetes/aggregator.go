package kubernetes

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Jdavid77/kubernetes-dashboard/internal/models"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// Aggregator aggregates Kubernetes metrics with caching
type Aggregator struct {
	client      *Client
	cache       *metricsCache
	cacheTTL    time.Duration
	mu          sync.RWMutex
}

type metricsCache struct {
	metrics   *models.ClusterMetrics
	apps      []models.Application
	timestamp time.Time
}

// NewAggregator creates a new metrics aggregator
func NewAggregator(client *Client, cacheTTL time.Duration) *Aggregator {
	return &Aggregator{
		client:   client,
		cacheTTL: cacheTTL,
		cache:    &metricsCache{},
	}
}

// GetClusterMetrics retrieves and aggregates cluster metrics
func (a *Aggregator) GetClusterMetrics(ctx context.Context) (*models.ClusterMetrics, error) {
	a.mu.RLock()
	if a.cache.metrics != nil && time.Since(a.cache.timestamp) < a.cacheTTL {
		metrics := a.cache.metrics
		a.mu.RUnlock()
		return metrics, nil
	}
	a.mu.RUnlock()

	// Fetch fresh data
	metrics, err := a.collectMetrics(ctx)
	if err != nil {
		return nil, err
	}

	// Update cache
	a.mu.Lock()
	a.cache.metrics = metrics
	a.cache.timestamp = time.Now()
	a.mu.Unlock()

	return metrics, nil
}

// collectMetrics collects metrics from Kubernetes
func (a *Aggregator) collectMetrics(ctx context.Context) (*models.ClusterMetrics, error) {
	// Get pods
	pods, err := a.client.GetPods(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %w", err)
	}

	// Get gateways (Gateway API)
	gatewayCount, _ := a.client.GetGateways(ctx)

	// Get HTTPRoutes (Gateway API)
	httpRouteCount, _ := a.client.GetHTTPRoutes(ctx)

	// Get services
	services, err := a.client.GetServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}

	// Count running pods
	running, _, _ := CountPodsByStatus(pods)

	// Get CPU and Memory metrics
	cpuUsage, memoryUsage := "0%", "0GB"

	nodeMetrics, err := a.client.GetNodeMetrics(ctx)
	if err == nil && len(nodeMetrics.Items) > 0 {
		var totalCPU, totalMemory int64

		for _, node := range nodeMetrics.Items {
			// Get CPU usage from the node metrics
			if cpu, ok := node.Usage[corev1.ResourceCPU]; ok {
				totalCPU += cpu.MilliValue()
			}
			// Get Memory usage from the node metrics
			if memory, ok := node.Usage[corev1.ResourceMemory]; ok {
				totalMemory += memory.Value()
			}
		}

		// Calculate CPU usage in cores (1000 millicores = 1 core)
		cpuCores := float64(totalCPU) / 1000.0
		cpuUsage = fmt.Sprintf("%.2f cores", cpuCores)

		// Convert memory to GB
		memoryGB := float64(totalMemory) / (1024 * 1024 * 1024)
		memoryUsage = fmt.Sprintf("%.1fGB", memoryGB)
	}

	return &models.ClusterMetrics{
		CPUUsage:       cpuUsage,
		MemoryUsage:    memoryUsage,
		PodsRunning:    running,
		GatewayCount:   gatewayCount,
		HTTPRouteCount: httpRouteCount,
		ServicesCount:  len(services.Items),
		Timestamp:      time.Now(),
	}, nil
}

// GetApplications retrieves all applications
func (a *Aggregator) GetApplications(ctx context.Context) ([]models.Application, error) {
	a.mu.RLock()
	if a.cache.apps != nil && time.Since(a.cache.timestamp) < a.cacheTTL {
		apps := a.cache.apps
		a.mu.RUnlock()
		return apps, nil
	}
	a.mu.RUnlock()

	// Fetch fresh data
	apps, err := a.collectApplications(ctx)
	if err != nil {
		return nil, err
	}

	// Update cache
	a.mu.Lock()
	a.cache.apps = apps
	a.mu.Unlock()

	return apps, nil
}

// collectApplications collects application data from deployments
func (a *Aggregator) collectApplications(ctx context.Context) ([]models.Application, error) {
	deployments, err := a.client.GetDeployments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployments: %w", err)
	}

	apps := make([]models.Application, 0, len(deployments.Items))
	for _, dep := range deployments.Items {
		app := a.deploymentToApplication(&dep)
		apps = append(apps, app)
	}

	return apps, nil
}

// deploymentToApplication converts a deployment to an application model
func (a *Aggregator) deploymentToApplication(dep *appsv1.Deployment) models.Application {
	// Determine status based on ready replicas
	status := "Healthy"
	replicas := int32(1) // Default to 1 if not specified
	if dep.Spec.Replicas != nil {
		replicas = *dep.Spec.Replicas
	}

	if dep.Status.ReadyReplicas == 0 {
		status = "Down"
	} else if dep.Status.ReadyReplicas < replicas {
		status = "Degraded"
	}

	// Extract technology from labels or annotations
	technology := extractTechnology(dep)

	// Get container image
	image := ""
	if len(dep.Spec.Template.Spec.Containers) > 0 {
		image = dep.Spec.Template.Spec.Containers[0].Image
	}

	return models.Application{
		Name:       dep.Name,
		Namespace:  dep.Namespace,
		Status:     status,
		PodCount:   int(dep.Status.ReadyReplicas),
		Technology: technology,
		Image:      image,
		CreatedAt:  dep.CreationTimestamp.Time,
		Labels:     dep.Labels,
	}
}

// extractTechnology extracts technology stack from deployment labels/annotations
func extractTechnology(dep *appsv1.Deployment) []string {
	tech := []string{}

	// Check labels for common technology indicators
	if app, ok := dep.Labels["app"]; ok {
		tech = append(tech, app)
	}

	if lang, ok := dep.Labels["language"]; ok {
		tech = append(tech, lang)
	}

	if framework, ok := dep.Labels["framework"]; ok {
		tech = append(tech, framework)
	}

	// Check annotations
	if techAnnotation, ok := dep.Annotations["technology"]; ok {
		tech = append(tech, techAnnotation)
	}

	// If no technology info, try to infer from image
	if len(tech) == 0 && len(dep.Spec.Template.Spec.Containers) > 0 {
		image := dep.Spec.Template.Spec.Containers[0].Image
		tech = inferTechnologyFromImage(image)
	}

	return tech
}

// inferTechnologyFromImage attempts to infer technology from container image
func inferTechnologyFromImage(image string) []string {
	tech := []string{}

	// Simple pattern matching
	if contains(image, "node") || contains(image, "nodejs") {
		tech = append(tech, "Node.js")
	}
	if contains(image, "python") {
		tech = append(tech, "Python")
	}
	if contains(image, "golang") || contains(image, "go:") {
		tech = append(tech, "Go")
	}
	if contains(image, "nginx") {
		tech = append(tech, "nginx")
	}
	if contains(image, "postgres") {
		tech = append(tech, "PostgreSQL")
	}
	if contains(image, "redis") {
		tech = append(tech, "Redis")
	}
	if contains(image, "mongo") {
		tech = append(tech, "MongoDB")
	}

	return tech
}

// GetApplicationDetail retrieves detailed information about an application
func (a *Aggregator) GetApplicationDetail(ctx context.Context, namespace, name string) (*models.ApplicationDetail, error) {
	// Get deployment
	deployment, err := a.client.GetDeployment(ctx, namespace, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	// Get pods for this deployment
	pods, err := a.client.GetPodsByNamespace(ctx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %w", err)
	}

	// Filter pods belonging to this deployment
	podInfos := []models.PodInfo{}
	for _, pod := range pods.Items {
		if isOwnedByDeployment(&pod, deployment) {
			podInfos = append(podInfos, podToInfo(&pod))
		}
	}

	// Extract environment variables
	env := make(map[string]string)
	if len(deployment.Spec.Template.Spec.Containers) > 0 {
		for _, envVar := range deployment.Spec.Template.Spec.Containers[0].Env {
			env[envVar.Name] = envVar.Value
		}
	}

	// Extract resources
	resources := models.ResourceInfo{}
	if len(deployment.Spec.Template.Spec.Containers) > 0 {
		container := deployment.Spec.Template.Spec.Containers[0]
		if req, ok := container.Resources.Requests[corev1.ResourceCPU]; ok {
			resources.CPURequest = req.String()
		}
		if limit, ok := container.Resources.Limits[corev1.ResourceCPU]; ok {
			resources.CPULimit = limit.String()
		}
		if req, ok := container.Resources.Requests[corev1.ResourceMemory]; ok {
			resources.MemoryRequest = formatMemory(req)
		}
		if limit, ok := container.Resources.Limits[corev1.ResourceMemory]; ok {
			resources.MemoryLimit = formatMemory(limit)
		}
	}

	app := a.deploymentToApplication(deployment)

	return &models.ApplicationDetail{
		Application: app,
		Pods:        podInfos,
		Environment: env,
		Resources:   resources,
	}, nil
}

// isOwnedByDeployment checks if a pod is owned by a deployment
func isOwnedByDeployment(pod *corev1.Pod, deployment *appsv1.Deployment) bool {
	for _, owner := range pod.OwnerReferences {
		// Check if owned by a ReplicaSet
		if owner.Kind == "ReplicaSet" {
			// ReplicaSet names start with deployment name
			if len(owner.Name) > len(deployment.Name) &&
				owner.Name[:len(deployment.Name)] == deployment.Name {
				return true
			}
		}
	}
	return false
}

// podToInfo converts a pod to a PodInfo model
func podToInfo(pod *corev1.Pod) models.PodInfo {
	// Count ready containers
	readyContainers := 0
	totalContainers := len(pod.Status.ContainerStatuses)
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.Ready {
			readyContainers++
		}
	}

	// Count restarts
	var restarts int32
	for _, cs := range pod.Status.ContainerStatuses {
		restarts += cs.RestartCount
	}

	// Calculate age
	age := time.Since(pod.CreationTimestamp.Time).Round(time.Second).String()

	return models.PodInfo{
		Name:      pod.Name,
		Status:    string(pod.Status.Phase),
		Ready:     fmt.Sprintf("%d/%d", readyContainers, totalContainers),
		Restarts:  restarts,
		Age:       age,
		CreatedAt: pod.CreationTimestamp.Time,
	}
}

// GetPods retrieves all pods across all namespaces
func (a *Aggregator) GetPods(ctx context.Context) ([]models.Pod, error) {
	pods, err := a.client.GetPods(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %w", err)
	}

	result := make([]models.Pod, 0, len(pods.Items))
	for _, pod := range pods.Items {
		// Count ready containers
		readyContainers := 0
		totalContainers := len(pod.Status.ContainerStatuses)
		for _, cs := range pod.Status.ContainerStatuses {
			if cs.Ready {
				readyContainers++
			}
		}

		// Count restarts
		var restarts int32
		for _, cs := range pod.Status.ContainerStatuses {
			restarts += cs.RestartCount
		}

		// Calculate age
		age := time.Since(pod.CreationTimestamp.Time).Round(time.Second).String()

		result = append(result, models.Pod{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Status:    string(pod.Status.Phase),
			Ready:     fmt.Sprintf("%d/%d", readyContainers, totalContainers),
			Restarts:  restarts,
			Age:       age,
			Node:      pod.Spec.NodeName,
			IP:        pod.Status.PodIP,
			CreatedAt: pod.CreationTimestamp.Time,
		})
	}

	return result, nil
}

// GetServices retrieves all services
func (a *Aggregator) GetServices(ctx context.Context) ([]models.Service, error) {
	services, err := a.client.GetServices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %w", err)
	}

	result := make([]models.Service, 0, len(services.Items))
	for _, svc := range services.Items {
		// Get external IP
		externalIP := ""
		if len(svc.Status.LoadBalancer.Ingress) > 0 {
			if svc.Status.LoadBalancer.Ingress[0].IP != "" {
				externalIP = svc.Status.LoadBalancer.Ingress[0].IP
			} else if svc.Status.LoadBalancer.Ingress[0].Hostname != "" {
				externalIP = svc.Status.LoadBalancer.Ingress[0].Hostname
			}
		}

		// Get ports
		ports := make([]models.ServicePort, 0, len(svc.Spec.Ports))
		for _, port := range svc.Spec.Ports {
			ports = append(ports, models.ServicePort{
				Name:       port.Name,
				Protocol:   string(port.Protocol),
				Port:       port.Port,
				TargetPort: port.TargetPort.String(),
				NodePort:   port.NodePort,
			})
		}

		result = append(result, models.Service{
			Name:       svc.Name,
			Namespace:  svc.Namespace,
			Type:       string(svc.Spec.Type),
			ClusterIP:  svc.Spec.ClusterIP,
			ExternalIP: externalIP,
			Ports:      ports,
			Selector:   svc.Spec.Selector,
			CreatedAt:  svc.CreationTimestamp.Time,
		})
	}

	return result, nil
}

// GetNodes retrieves all nodes with their metrics and status
func (a *Aggregator) GetNodes(ctx context.Context) ([]models.Node, error) {
	// Get nodes
	nodeList, err := a.client.GetNodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}

	// Get node metrics
	nodeMetrics, err := a.client.GetNodeMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get node metrics: %w", err)
	}

	// Get pods to count per node
	pods, err := a.client.GetPods(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %w", err)
	}

	// Create a map of node name to pod count
	podCountByNode := make(map[string]int)
	for _, pod := range pods.Items {
		podCountByNode[pod.Spec.NodeName]++
	}

	// Create a map of node metrics for quick lookup
	metricsMap := make(map[string]corev1.ResourceList)
	for _, nm := range nodeMetrics.Items {
		metricsMap[nm.Name] = nm.Usage
	}

	// Build the result
	result := make([]models.Node, 0, len(nodeList.Items))
	for _, node := range nodeList.Items {
		// Determine node status
		status := "NotReady"
		for _, condition := range node.Status.Conditions {
			if condition.Type == corev1.NodeReady {
				if condition.Status == corev1.ConditionTrue {
					status = "Ready"
				}
				break
			}
		}

		// Get capacity
		cpuCapacity := node.Status.Capacity[corev1.ResourceCPU]
		memCapacity := node.Status.Capacity[corev1.ResourceMemory]

		// Calculate usage percentages
		var cpuUsage, memUsage float64
		if usage, ok := metricsMap[node.Name]; ok {
			// CPU usage percentage
			if cpuUsed, ok := usage[corev1.ResourceCPU]; ok {
				cpuUsageMillis := cpuUsed.MilliValue()
				cpuCapacityMillis := cpuCapacity.MilliValue()
				if cpuCapacityMillis > 0 {
					cpuUsage = float64(cpuUsageMillis) / float64(cpuCapacityMillis) * 100
				}
			}

			// Memory usage percentage
			if memUsed, ok := usage[corev1.ResourceMemory]; ok {
				memUsageBytes := memUsed.Value()
				memCapacityBytes := memCapacity.Value()
				if memCapacityBytes > 0 {
					memUsage = float64(memUsageBytes) / float64(memCapacityBytes) * 100
				}
			}
		}

		result = append(result, models.Node{
			Name:           node.Name,
			Status:         status,
			CPUUsage:       cpuUsage,
			MemoryUsage:    memUsage,
			CPUCapacity:    cpuCapacity.String(),
			MemoryCapacity: formatMemory(memCapacity),
			PodCount:       podCountByNode[node.Name],
		})
	}

	return result, nil
}

// formatMemory formats a memory quantity to a readable string
func formatMemory(q resource.Quantity) string {
	bytes := q.Value()
	if bytes >= 1024*1024*1024 {
		return fmt.Sprintf("%.1fGi", float64(bytes)/(1024*1024*1024))
	} else if bytes >= 1024*1024 {
		return fmt.Sprintf("%.1fMi", float64(bytes)/(1024*1024))
	}
	return fmt.Sprintf("%dKi", bytes/1024)
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
		 (len(s) > len(substr) &&
		  (s[:len(substr)] == substr ||
		   s[len(s)-len(substr):] == substr)))
}

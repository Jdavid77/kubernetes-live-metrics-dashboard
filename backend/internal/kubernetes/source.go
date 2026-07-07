package kubernetes

import (
	"context"

	"github.com/Jdavid77/kubernetes-dashboard/internal/models"
)

// MetricsSource is the interface all handlers use to retrieve cluster data.
// *Aggregator satisfies it; a fake can satisfy it in tests.
type MetricsSource interface {
	GetClusterMetrics(ctx context.Context) (*models.ClusterMetrics, error)
	GetApplications(ctx context.Context) ([]models.Application, error)
	GetApplicationDetail(ctx context.Context, namespace, name string) (*models.ApplicationDetail, error)
	GetPods(ctx context.Context) ([]models.Pod, error)
	GetServices(ctx context.Context) ([]models.Service, error)
	GetNodes(ctx context.Context) ([]models.Node, error)
	GetNamespaces(ctx context.Context) ([]string, error)
}

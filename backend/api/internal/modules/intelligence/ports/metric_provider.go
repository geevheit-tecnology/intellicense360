package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/domain"
)

type MetricProvider interface {
	Metrics(ctx context.Context, tenantID string, q Query) (Page[domain.IntelligenceMetric], error)
}

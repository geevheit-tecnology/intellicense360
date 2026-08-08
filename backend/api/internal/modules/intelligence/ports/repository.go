package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/domain"
)

type RecommendationRepository interface {
	Save(ctx context.Context, recommendation domain.Recommendation) error
	FindActiveByTenant(ctx context.Context, tenantID string) ([]domain.Recommendation, error)
}

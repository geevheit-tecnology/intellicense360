package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/application"
)

type RecommendationService interface {
	CurrentOperationalSummary(ctx context.Context) application.OperationalSummary
}

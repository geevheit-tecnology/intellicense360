package ports

import (
	"context"

	coreevents "github.com/geevheit/intelligence360/backend/api/internal/events"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/domain"
)

type MetricService interface {
	Create(ctx context.Context, item domain.IntelligenceMetric) (domain.IntelligenceMetric, error)
	FindByID(ctx context.Context, tenantID, id string) (domain.IntelligenceMetric, error)
	Search(ctx context.Context, tenantID string, q Query) (Page[domain.IntelligenceMetric], error)
}

type AnalysisService interface {
	DetectFuelConsumptionDeviation(ctx context.Context, tenantID string, historical []float64, current float64, period domain.IntelligencePeriod) (domain.AnalysisResult, error)
	DetectRepeatedFailure(ctx context.Context, tenantID, category, referenceID string, occurrences int, period domain.IntelligencePeriod) (domain.Risk, error)
	DetectMaintenanceRecurrence(ctx context.Context, tenantID, referenceID string, occurrences int, period domain.IntelligencePeriod) (domain.Risk, error)
}

type IntelligenceService interface {
	CreateAnomaly(ctx context.Context, item domain.Anomaly) (domain.Anomaly, error)
	SearchAnomalies(ctx context.Context, tenantID string, q Query) (Page[domain.Anomaly], error)
	CreateRisk(ctx context.Context, item domain.Risk) (domain.Risk, error)
	SearchRisks(ctx context.Context, tenantID string, q Query) (Page[domain.Risk], error)
	CreateOpportunity(ctx context.Context, item domain.Opportunity) (domain.Opportunity, error)
	SearchOpportunities(ctx context.Context, tenantID string, q Query) (Page[domain.Opportunity], error)
	CreateRecommendation(ctx context.Context, item domain.Recommendation) (domain.Recommendation, error)
	SearchRecommendations(ctx context.Context, tenantID string, q Query) (Page[domain.Recommendation], error)
	CreateInsight(ctx context.Context, item domain.Insight) (domain.Insight, error)
	FindInsight(ctx context.Context, tenantID, id string) (domain.Insight, error)
	SearchInsights(ctx context.Context, tenantID string, q Query) (Page[domain.Insight], error)
	AcknowledgeInsight(ctx context.Context, tenantID, id string) (domain.Insight, error)
	ResolveInsight(ctx context.Context, tenantID, id string) (domain.Insight, error)
	DismissInsight(ctx context.Context, tenantID, id string) (domain.Insight, error)
	CurrentOperationalSummary(ctx context.Context) domain.OperationalSummary
}

type EventProjectionService interface {
	Handle(ctx context.Context, event coreevents.DomainEvent) error
}

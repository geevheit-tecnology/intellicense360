package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/domain"
)

type Query struct {
	Search  string
	Page    int
	PerPage int
	Sort    string
	Filters map[string]string
}

type Page[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type MetricRepository interface {
	Create(ctx context.Context, item domain.IntelligenceMetric) (domain.IntelligenceMetric, error)
	FindByID(ctx context.Context, tenantID, id string) (domain.IntelligenceMetric, error)
	Search(ctx context.Context, tenantID string, q Query) (Page[domain.IntelligenceMetric], error)
}

type AnomalyRepository interface {
	Create(ctx context.Context, item domain.Anomaly) (domain.Anomaly, error)
	Search(ctx context.Context, tenantID string, q Query) (Page[domain.Anomaly], error)
}

type RiskRepository interface {
	Create(ctx context.Context, item domain.Risk) (domain.Risk, error)
	Search(ctx context.Context, tenantID string, q Query) (Page[domain.Risk], error)
}

type OpportunityRepository interface {
	Create(ctx context.Context, item domain.Opportunity) (domain.Opportunity, error)
	Search(ctx context.Context, tenantID string, q Query) (Page[domain.Opportunity], error)
}

type RecommendationRepository interface {
	Create(ctx context.Context, item domain.Recommendation) (domain.Recommendation, error)
	Search(ctx context.Context, tenantID string, q Query) (Page[domain.Recommendation], error)
}

type InsightRepository interface {
	Create(ctx context.Context, item domain.Insight) (domain.Insight, error)
	FindByID(ctx context.Context, tenantID, id string) (domain.Insight, error)
	FindByDeduplicationKey(ctx context.Context, tenantID, key string) (domain.Insight, error)
	Search(ctx context.Context, tenantID string, q Query) (Page[domain.Insight], error)
	Update(ctx context.Context, item domain.Insight) (domain.Insight, error)
}

type RuleRepository interface {
	Create(ctx context.Context, item domain.IntelligenceRule) (domain.IntelligenceRule, error)
	Search(ctx context.Context, tenantID string, q Query) (Page[domain.IntelligenceRule], error)
}

type ReadModelRepository interface {
	Upsert(ctx context.Context, item domain.ReadModel) (domain.ReadModel, error)
	Search(ctx context.Context, tenantID string, q Query) (Page[domain.ReadModel], error)
	HasProcessedEvent(ctx context.Context, tenantID, eventID string) (bool, error)
	MarkProcessedEvent(ctx context.Context, tenantID, eventID string) error
}

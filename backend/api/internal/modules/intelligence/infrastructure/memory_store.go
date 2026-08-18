package infrastructure

import (
	"context"
	"strings"
	"sync"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/ports"
)

type MemoryStore struct {
	mu              sync.RWMutex
	metrics         map[string]domain.IntelligenceMetric
	anomalies       map[string]domain.Anomaly
	risks           map[string]domain.Risk
	opportunities   map[string]domain.Opportunity
	recommendations map[string]domain.Recommendation
	insights        map[string]domain.Insight
	rules           map[string]domain.IntelligenceRule
	readModels      map[string]domain.ReadModel
	processedEvents map[string]bool
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		metrics: map[string]domain.IntelligenceMetric{}, anomalies: map[string]domain.Anomaly{}, risks: map[string]domain.Risk{},
		opportunities: map[string]domain.Opportunity{}, recommendations: map[string]domain.Recommendation{}, insights: map[string]domain.Insight{},
		rules: map[string]domain.IntelligenceRule{}, readModels: map[string]domain.ReadModel{}, processedEvents: map[string]bool{},
	}
}

func (s *MemoryStore) Metrics() MetricRepository            { return MetricRepository{s: s} }
func (s *MemoryStore) Anomalies() AnomalyRepository         { return AnomalyRepository{s: s} }
func (s *MemoryStore) Risks() RiskRepository                { return RiskRepository{s: s} }
func (s *MemoryStore) Opportunities() OpportunityRepository { return OpportunityRepository{s: s} }
func (s *MemoryStore) Recommendations() RecommendationRepository {
	return RecommendationRepository{s: s}
}
func (s *MemoryStore) Insights() InsightRepository     { return InsightRepository{s: s} }
func (s *MemoryStore) Rules() RuleRepository           { return RuleRepository{s: s} }
func (s *MemoryStore) ReadModels() ReadModelRepository { return ReadModelRepository{s: s} }

type MetricRepository struct{ s *MemoryStore }
type AnomalyRepository struct{ s *MemoryStore }
type RiskRepository struct{ s *MemoryStore }
type OpportunityRepository struct{ s *MemoryStore }
type RecommendationRepository struct{ s *MemoryStore }
type InsightRepository struct{ s *MemoryStore }
type RuleRepository struct{ s *MemoryStore }
type ReadModelRepository struct{ s *MemoryStore }

func (r MetricRepository) Create(_ context.Context, item domain.IntelligenceMetric) (domain.IntelligenceMetric, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.metrics[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r MetricRepository) FindByID(_ context.Context, tenantID, id string) (domain.IntelligenceMetric, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.metrics[key(tenantID, id)]
	if !ok {
		return domain.IntelligenceMetric{}, application.ErrNotFound
	}
	return item, nil
}
func (r MetricRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.IntelligenceMetric], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.IntelligenceMetric{}
	for _, item := range r.s.metrics {
		if item.TenantID == tenantID && match(item.Name+" "+string(item.MetricType)+" "+item.DimensionValue, q) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}

func (r AnomalyRepository) Create(_ context.Context, item domain.Anomaly) (domain.Anomaly, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.anomalies[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r AnomalyRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.Anomaly], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Anomaly{}
	for _, item := range r.s.anomalies {
		if item.TenantID == tenantID && match(item.Type+" "+string(item.Severity)+" "+string(item.Status), q) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}

func (r RiskRepository) Create(_ context.Context, item domain.Risk) (domain.Risk, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.risks[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r RiskRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.Risk], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Risk{}
	for _, item := range r.s.risks {
		if item.TenantID == tenantID && match(item.Category+" "+string(item.Severity)+" "+string(item.Status), q) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}

func (r OpportunityRepository) Create(_ context.Context, item domain.Opportunity) (domain.Opportunity, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.opportunities[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r OpportunityRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.Opportunity], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Opportunity{}
	for _, item := range r.s.opportunities {
		if item.TenantID == tenantID && match(item.Category+" "+string(item.Status), q) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}

func (r RecommendationRepository) Create(_ context.Context, item domain.Recommendation) (domain.Recommendation, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.recommendations[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r RecommendationRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.Recommendation], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Recommendation{}
	for _, item := range r.s.recommendations {
		if item.TenantID == tenantID && match(item.Title+" "+item.ImpactArea+" "+string(item.Priority), q) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}

func (r InsightRepository) Create(_ context.Context, item domain.Insight) (domain.Insight, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.insights[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r InsightRepository) FindByID(_ context.Context, tenantID, id string) (domain.Insight, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.insights[key(tenantID, id)]
	if !ok {
		return domain.Insight{}, application.ErrNotFound
	}
	return item, nil
}
func (r InsightRepository) FindByDeduplicationKey(_ context.Context, tenantID, dedupe string) (domain.Insight, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	for _, item := range r.s.insights {
		if item.TenantID == tenantID && item.DeduplicationKey == dedupe {
			return item, nil
		}
	}
	return domain.Insight{}, application.ErrNotFound
}
func (r InsightRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.Insight], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Insight{}
	for _, item := range r.s.insights {
		if item.TenantID == tenantID && match(item.Title+" "+item.Category+" "+string(item.Severity)+" "+string(item.Status), q) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r InsightRepository) Update(_ context.Context, item domain.Insight) (domain.Insight, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(item.TenantID, item.ID)
	if _, ok := r.s.insights[k]; !ok {
		return domain.Insight{}, application.ErrNotFound
	}
	r.s.insights[k] = item
	return item, nil
}

func (r RuleRepository) Create(_ context.Context, item domain.IntelligenceRule) (domain.IntelligenceRule, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.rules[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r RuleRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.IntelligenceRule], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.IntelligenceRule{}
	for _, item := range r.s.rules {
		if (item.TenantID == tenantID || item.TenantID == "") && match(item.Name+" "+string(item.RuleType), q) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}

func (r ReadModelRepository) Upsert(_ context.Context, item domain.ReadModel) (domain.ReadModel, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.readModels[key(item.TenantID, item.SourceType+":"+item.ReferenceID+":"+item.MetricName)] = item
	return item, nil
}
func (r ReadModelRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.ReadModel], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.ReadModel{}
	for _, item := range r.s.readModels {
		if item.TenantID == tenantID && match(item.Source+" "+item.SourceType+" "+item.MetricName, q) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r ReadModelRepository) HasProcessedEvent(_ context.Context, tenantID, eventID string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	return r.s.processedEvents[key(tenantID, eventID)], nil
}
func (r ReadModelRepository) MarkProcessedEvent(_ context.Context, tenantID, eventID string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, eventID)
	if r.s.processedEvents[k] {
		return nil
	}
	r.s.processedEvents[k] = true
	return nil
}

func key(tenantID, id string) string { return tenantID + ":" + id }

func match(text string, q ports.Query) bool {
	if q.Filters != nil {
		if status := q.Filters["status"]; status != "" && !strings.Contains(strings.ToLower(text), strings.ToLower(status)) {
			return false
		}
		if severity := q.Filters["severity"]; severity != "" && !strings.Contains(strings.ToLower(text), strings.ToLower(severity)) {
			return false
		}
		if category := q.Filters["category"]; category != "" && !strings.Contains(strings.ToLower(text), strings.ToLower(category)) {
			return false
		}
	}
	if q.Search == "" {
		return true
	}
	return strings.Contains(strings.ToLower(text), strings.ToLower(q.Search))
}

func page[T any](items []T, q ports.Query) ports.Page[T] {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PerPage <= 0 {
		q.PerPage = 25
	}
	total := len(items)
	start := (q.Page - 1) * q.PerPage
	if start > total {
		start = total
	}
	end := start + q.PerPage
	if end > total {
		end = total
	}
	totalPages := 0
	if total > 0 {
		totalPages = (total + q.PerPage - 1) / q.PerPage
	}
	return ports.Page[T]{Data: items[start:end], Page: q.Page, PerPage: q.PerPage, Total: total, TotalPages: totalPages}
}

package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
	"time"

	coreevents "github.com/geevheit/intelligence360/backend/api/internal/events"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/ports"
)

type MetricService struct{ repo ports.MetricRepository }

func NewMetricService(repo ports.MetricRepository) MetricService { return MetricService{repo: repo} }

func (s MetricService) Create(ctx context.Context, item domain.IntelligenceMetric) (domain.IntelligenceMetric, error) {
	if item.TenantID == "" || item.Name == "" || item.MetricType == "" {
		return domain.IntelligenceMetric{}, ErrInvalidData
	}
	now := time.Now().UTC()
	item.ID = newID("met")
	item.CalculatedAt = zeroTime(item.CalculatedAt, now)
	item.CreatedAt, item.UpdatedAt, item.Version = now, now, 1
	item.Confidence = clamp01(item.Confidence)
	if item.Metadata == nil {
		item.Metadata = map[string]string{}
	}
	return s.repo.Create(ctx, item)
}

func (s MetricService) FindByID(ctx context.Context, tenantID, id string) (domain.IntelligenceMetric, error) {
	return s.repo.FindByID(ctx, tenantID, id)
}

func (s MetricService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.IntelligenceMetric], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(q))
}

type IntelligenceService struct {
	anomalies       ports.AnomalyRepository
	risks           ports.RiskRepository
	opportunities   ports.OpportunityRepository
	recommendations ports.RecommendationRepository
	insights        ports.InsightRepository
}

func NewIntelligenceService(anomalies ports.AnomalyRepository, risks ports.RiskRepository, opportunities ports.OpportunityRepository, recommendations ports.RecommendationRepository, insights ports.InsightRepository) IntelligenceService {
	return IntelligenceService{anomalies: anomalies, risks: risks, opportunities: opportunities, recommendations: recommendations, insights: insights}
}

func (s IntelligenceService) CreateAnomaly(ctx context.Context, item domain.Anomaly) (domain.Anomaly, error) {
	if item.TenantID == "" || item.Type == "" {
		return domain.Anomaly{}, ErrInvalidData
	}
	now := time.Now().UTC()
	item.ID, item.DetectedAt, item.CreatedAt, item.UpdatedAt, item.Version = newID("ano"), zeroTime(item.DetectedAt, now), now, now, 1
	item.Confidence = clamp01(item.Confidence)
	if item.Status == "" {
		item.Status = domain.StatusNew
	}
	return s.anomalies.Create(ctx, item)
}
func (s IntelligenceService) SearchAnomalies(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Anomaly], error) {
	return s.anomalies.Search(ctx, tenantID, normalizeQuery(q))
}
func (s IntelligenceService) CreateRisk(ctx context.Context, item domain.Risk) (domain.Risk, error) {
	if item.TenantID == "" || item.Category == "" {
		return domain.Risk{}, ErrInvalidData
	}
	now := time.Now().UTC()
	item.ID, item.DetectedAt, item.CreatedAt, item.UpdatedAt, item.Version = newID("ris"), zeroTime(item.DetectedAt, now), now, now, 1
	item.Confidence, item.Probability = clamp01(item.Confidence), clamp01(item.Probability)
	if item.Status == "" {
		item.Status = domain.StatusNew
	}
	return s.risks.Create(ctx, item)
}
func (s IntelligenceService) SearchRisks(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Risk], error) {
	return s.risks.Search(ctx, tenantID, normalizeQuery(q))
}
func (s IntelligenceService) CreateOpportunity(ctx context.Context, item domain.Opportunity) (domain.Opportunity, error) {
	if item.TenantID == "" || item.Category == "" {
		return domain.Opportunity{}, ErrInvalidData
	}
	now := time.Now().UTC()
	item.ID, item.DetectedAt, item.CreatedAt, item.UpdatedAt, item.Version = newID("opp"), zeroTime(item.DetectedAt, now), now, now, 1
	item.Confidence = clamp01(item.Confidence)
	if item.Status == "" {
		item.Status = domain.StatusNew
	}
	return s.opportunities.Create(ctx, item)
}
func (s IntelligenceService) SearchOpportunities(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Opportunity], error) {
	return s.opportunities.Search(ctx, tenantID, normalizeQuery(q))
}
func (s IntelligenceService) CreateRecommendation(ctx context.Context, item domain.Recommendation) (domain.Recommendation, error) {
	if item.TenantID == "" || item.Title == "" || item.SuggestedAction == "" {
		return domain.Recommendation{}, ErrInvalidData
	}
	now := time.Now().UTC()
	item.ID, item.CreatedAt, item.UpdatedAt, item.Version = newID("rec"), now, now, 1
	item.Confidence = clamp01(item.Confidence)
	if item.Status == "" {
		item.Status = domain.StatusNew
	}
	return s.recommendations.Create(ctx, item)
}
func (s IntelligenceService) SearchRecommendations(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Recommendation], error) {
	return s.recommendations.Search(ctx, tenantID, normalizeQuery(q))
}
func (s IntelligenceService) CreateInsight(ctx context.Context, item domain.Insight) (domain.Insight, error) {
	if item.TenantID == "" || item.Title == "" || item.Type == "" {
		return domain.Insight{}, ErrInvalidData
	}
	if item.DeduplicationKey == "" {
		item.DeduplicationKey = DeduplicationKey(item.TenantID, item.Type, item.Category, item.MetricID, item.CreatedAt.Format("2006-01"), "v1")
	}
	if _, err := s.insights.FindByDeduplicationKey(ctx, item.TenantID, item.DeduplicationKey); err == nil {
		return domain.Insight{}, ErrDuplicateInsight
	}
	now := time.Now().UTC()
	item.ID, item.CreatedAt, item.UpdatedAt, item.Version = newID("ins"), now, now, 1
	item.Confidence = clamp01(item.Confidence)
	if item.Status == "" {
		item.Status = domain.StatusNew
	}
	return s.insights.Create(ctx, item)
}
func (s IntelligenceService) FindInsight(ctx context.Context, tenantID, id string) (domain.Insight, error) {
	return s.insights.FindByID(ctx, tenantID, id)
}
func (s IntelligenceService) SearchInsights(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Insight], error) {
	return s.insights.Search(ctx, tenantID, normalizeQuery(q))
}
func (s IntelligenceService) AcknowledgeInsight(ctx context.Context, tenantID, id string) (domain.Insight, error) {
	return s.transitionInsight(ctx, tenantID, id, domain.StatusAcknowledged)
}
func (s IntelligenceService) ResolveInsight(ctx context.Context, tenantID, id string) (domain.Insight, error) {
	return s.transitionInsight(ctx, tenantID, id, domain.StatusResolved)
}
func (s IntelligenceService) DismissInsight(ctx context.Context, tenantID, id string) (domain.Insight, error) {
	return s.transitionInsight(ctx, tenantID, id, domain.StatusDismissed)
}
func (s IntelligenceService) transitionInsight(ctx context.Context, tenantID, id string, status domain.InsightStatus) (domain.Insight, error) {
	item, err := s.insights.FindByID(ctx, tenantID, id)
	if err != nil {
		return domain.Insight{}, err
	}
	item.Status, item.UpdatedAt, item.Version = status, time.Now().UTC(), item.Version+1
	return s.insights.Update(ctx, item)
}
func (s IntelligenceService) CurrentOperationalSummary(ctx context.Context) domain.OperationalSummary {
	_ = ctx
	return domain.OperationalSummary{
		CurrentOperation: "Intelligence Engine deterministico ativo",
		CurrentRisks:     []string{"Aguardando volume historico suficiente para alta confianca"},
		CurrentSavings:   "Estimativas serao exibidas apenas com evidencias suficientes",
		Priorities:       []string{"Coletar eventos operacionais", "Construir baselines por tenant", "Validar regras deterministicas"},
		Recommendations: []domain.Recommendation{{
			ID: "bootstrap-intelligence-engine", Title: "Consolidar baselines operacionais", Description: "A inteligencia deterministica precisa de historico para comparar consumo, custo e recorrencia.", Priority: domain.PriorityHigh, ImpactArea: "mission-control", Confidence: 0.55,
		}},
	}
}

type AnalysisService struct {
	metrics ports.MetricRepository
	core    IntelligenceService
}

func NewAnalysisService(metrics ports.MetricRepository, core IntelligenceService) AnalysisService {
	return AnalysisService{metrics: metrics, core: core}
}

func (s AnalysisService) DetectFuelConsumptionDeviation(ctx context.Context, tenantID string, historical []float64, current float64, period domain.IntelligencePeriod) (domain.AnalysisResult, error) {
	baseline, confidence, ok := MovingAverage(historical)
	metric := domain.IntelligenceMetric{TenantID: tenantID, MetricType: domain.MetricFuel, Name: "Fuel Consumption", Value: current, Unit: "km/L", PeriodStart: period.Start, PeriodEnd: period.End, Source: "intelligence-rule", Confidence: confidence}
	savedMetric, err := NewMetricService(s.metrics).Create(ctx, metric)
	if err != nil {
		return domain.AnalysisResult{}, err
	}
	result := domain.AnalysisResult{Metric: savedMetric, Trend: Trend(historical, period), Confidence: confidence, InsufficientData: !ok}
	if !ok {
		return result, nil
	}
	deviation := current - baseline
	deviationPct := 0.0
	if baseline != 0 {
		deviationPct = deviation / baseline * 100
	}
	if math.Abs(deviationPct) < 10 {
		return result, nil
	}
	severity := severityFromDeviation(math.Abs(deviationPct))
	evidence := domain.InsightEvidence{Source: "fuel-read-model", SourceType: "metric", Metric: "fuel_consumption", ObservedValue: current, ExpectedValue: baseline, Period: period, Timestamp: time.Now().UTC(), Explanation: fmt.Sprintf("Consumo atual %.2f km/L comparado ao baseline %.2f km/L", current, baseline)}
	anomaly, err := s.core.CreateAnomaly(ctx, domain.Anomaly{TenantID: tenantID, Type: "fuel_consumption_deviation", Severity: severity, MetricID: savedMetric.ID, ObservedValue: current, ExpectedValue: baseline, Deviation: deviation, DeviationPercentage: deviationPct, Period: period, Evidence: []domain.InsightEvidence{evidence}, Confidence: confidence})
	if err != nil {
		return domain.AnalysisResult{}, err
	}
	rec, err := s.core.CreateRecommendation(ctx, domain.Recommendation{TenantID: tenantID, Title: "Investigar desvio de consumo de combustivel", WhatHappened: "Consumo atual abaixo do baseline historico.", WhyItMatters: "Desvios persistentes podem indicar perda de eficiencia operacional ou problema mecanico.", Evidence: []domain.InsightEvidence{evidence}, SuggestedAction: "Investigar eventos recentes de abastecimento, condicoes operacionais e comportamento do ativo.", ExpectedImpact: domain.InsightImpact{OperationalImpact: "reduzir consumo anormal", Confidence: confidence, Currency: "BRL"}, Confidence: confidence, Priority: priorityFromSeverity(severity), ImpactArea: "fuel"})
	if err != nil {
		return domain.AnalysisResult{}, err
	}
	insight, err := s.core.CreateInsight(ctx, domain.Insight{TenantID: tenantID, Type: "anomaly", Title: "Desvio de consumo de combustivel detectado", Summary: fmt.Sprintf("Consumo %.1f%% diferente do baseline historico.", deviationPct), Category: "fuel_efficiency", Severity: severity, Evidence: []domain.InsightEvidence{evidence}, MetricID: savedMetric.ID, AnomalyID: anomaly.ID, RecommendationID: rec.ID, EstimatedImpact: rec.ExpectedImpact, Confidence: confidence, Priority: rec.Priority, DeduplicationKey: DeduplicationKey(tenantID, "fuel_consumption_deviation", "fuel_efficiency", savedMetric.DimensionValue, period.Start.Format("2006-01"), "v1")})
	if err != nil && err != ErrDuplicateInsight {
		return domain.AnalysisResult{}, err
	}
	result.Anomaly, result.Recommendation = &anomaly, &rec
	if err == nil {
		result.Insight = &insight
	}
	return result, nil
}

func (s AnalysisService) DetectRepeatedFailure(ctx context.Context, tenantID, category, referenceID string, occurrences int, period domain.IntelligencePeriod) (domain.Risk, error) {
	if tenantID == "" || category == "" || occurrences <= 0 {
		return domain.Risk{}, ErrInvalidData
	}
	confidence := Confidence(0.8, occurrences, 0.7, 0.8)
	severity := domain.SeverityMedium
	if occurrences >= 5 {
		severity = domain.SeverityHigh
	}
	evidence := domain.InsightEvidence{Source: "read-model", SourceType: "failure_pattern", ReferenceID: referenceID, ObservedValue: float64(occurrences), ExpectedValue: 1, Period: period, Timestamp: time.Now().UTC(), Explanation: "Ocorrencia repetida no periodo analisado"}
	return s.core.CreateRisk(ctx, domain.Risk{TenantID: tenantID, Category: category, Severity: severity, Probability: clamp01(float64(occurrences) / 6), Impact: domain.InsightImpact{OperationalImpact: "recorrencia operacional", Confidence: confidence}, Confidence: confidence, Evidence: []domain.InsightEvidence{evidence}, Status: domain.StatusNew})
}

func (s AnalysisService) DetectMaintenanceRecurrence(ctx context.Context, tenantID, referenceID string, occurrences int, period domain.IntelligencePeriod) (domain.Risk, error) {
	return s.DetectRepeatedFailure(ctx, tenantID, "maintenance_recurrence", referenceID, occurrences, period)
}

type ProjectionService struct{ readModels ports.ReadModelRepository }

func NewProjectionService(repo ports.ReadModelRepository) ProjectionService {
	return ProjectionService{readModels: repo}
}

func (s ProjectionService) Handle(ctx context.Context, event coreevents.DomainEvent) error {
	if event.EventID == "" || event.TenantID == "" {
		return ErrInvalidData
	}
	done, err := s.readModels.HasProcessedEvent(ctx, event.TenantID, event.EventID)
	if err != nil || done {
		return err
	}
	now := time.Now().UTC()
	value := 1.0
	if raw, ok := event.Payload["value"].(float64); ok {
		value = raw
	}
	_, err = s.readModels.Upsert(ctx, domain.ReadModel{ID: newID("irm"), TenantID: event.TenantID, Source: "event", SourceType: event.EventType, ReferenceID: event.AggregateID, MetricName: event.EventType, Value: value, PeriodStart: event.OccurredAt, PeriodEnd: event.OccurredAt, Metadata: event.Metadata, LastEventID: event.EventID, CreatedAt: now, UpdatedAt: now, Version: 1})
	if err != nil {
		return err
	}
	return s.readModels.MarkProcessedEvent(ctx, event.TenantID, event.EventID)
}

func MovingAverage(values []float64) (float64, float64, bool) {
	if len(values) < 3 {
		return 0, Confidence(0.4, len(values), 0.2, 0.5), false
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values)), Confidence(0.9, len(values), 0.8, 0.8), true
}

func Trend(values []float64, period domain.IntelligencePeriod) domain.Trend {
	if len(values) < 3 {
		return domain.Trend{Direction: domain.TrendInsufficientData, Period: period, Confidence: Confidence(0.4, len(values), 0.2, 0.5)}
	}
	first, last := values[0], values[len(values)-1]
	magnitude := last - first
	direction := domain.TrendStable
	if math.Abs(magnitude) > math.Abs(first)*0.1 {
		if magnitude > 0 {
			direction = domain.TrendIncreasing
		} else {
			direction = domain.TrendDecreasing
		}
	}
	return domain.Trend{Direction: direction, Magnitude: magnitude, Period: period, Confidence: Confidence(0.8, len(values), 0.7, 0.8)}
}

func Confidence(dataQuality float64, sampleSize int, historicalConsistency float64, ruleConfidence float64) float64 {
	sizeScore := math.Min(float64(sampleSize)/10, 1)
	return clamp01((dataQuality * 0.3) + (sizeScore * 0.25) + (historicalConsistency * 0.25) + (ruleConfidence * 0.2))
}

func DeduplicationKey(parts ...string) string {
	clean := []string{}
	for _, part := range parts {
		clean = append(clean, strings.ToLower(strings.TrimSpace(part)))
	}
	return strings.Join(clean, ":")
}

func severityFromDeviation(pct float64) domain.Severity {
	switch {
	case pct >= 30:
		return domain.SeverityCritical
	case pct >= 20:
		return domain.SeverityHigh
	case pct >= 10:
		return domain.SeverityMedium
	default:
		return domain.SeverityLow
	}
}

func priorityFromSeverity(severity domain.Severity) domain.Priority {
	switch severity {
	case domain.SeverityCritical:
		return domain.PriorityCritical
	case domain.SeverityHigh:
		return domain.PriorityHigh
	case domain.SeverityMedium:
		return domain.PriorityMedium
	default:
		return domain.PriorityLow
	}
}

func normalizeQuery(q ports.Query) ports.Query {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PerPage <= 0 {
		q.PerPage = 25
	}
	if q.PerPage > 100 {
		q.PerPage = 100
	}
	if q.Filters == nil {
		q.Filters = map[string]string{}
	}
	return q
}

func zeroTime(value time.Time, fallback time.Time) time.Time {
	if value.IsZero() {
		return fallback
	}
	return value
}

func clamp01(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}

func newID(prefix string) string {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return prefix + "_" + time.Now().UTC().Format("20060102150405")
	}
	return prefix + "_" + hex.EncodeToString(b[:])
}

package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreevents "github.com/geevheit/intelligence360/backend/api/internal/events"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/infrastructure"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/ports"
)

func services() (*infrastructure.MemoryStore, application.IntelligenceService, application.AnalysisService) {
	store := infrastructure.NewMemoryStore()
	core := application.NewIntelligenceService(store.Anomalies(), store.Risks(), store.Opportunities(), store.Recommendations(), store.Insights())
	analysis := application.NewAnalysisService(store.Metrics(), core)
	return store, core, analysis
}

func TestFuelConsumptionDeviationCreatesAnomalyEvidenceRecommendation(t *testing.T) {
	_, _, analysis := services()
	period := domain.IntelligencePeriod{Start: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2026, 8, 31, 0, 0, 0, 0, time.UTC)}
	result, err := analysis.DetectFuelConsumptionDeviation(context.Background(), "tenant-1", []float64{2.8, 2.7, 2.9, 2.8}, 2.35, period)
	if err != nil {
		t.Fatalf("detect deviation: %v", err)
	}
	if result.Anomaly == nil || result.Recommendation == nil || result.Insight == nil {
		t.Fatalf("expected anomaly, recommendation and insight: %+v", result)
	}
	if result.Anomaly.Severity != domain.SeverityMedium {
		t.Fatalf("severity = %s", result.Anomaly.Severity)
	}
	if len(result.Anomaly.Evidence) == 0 || result.Anomaly.Confidence <= 0.5 {
		t.Fatalf("missing evidence/confidence: %+v", result.Anomaly)
	}
}

func TestRepeatedChecklistFailureCreatesRisk(t *testing.T) {
	_, _, analysis := services()
	risk, err := analysis.DetectRepeatedFailure(context.Background(), "tenant-1", "checklist_failure_pattern", "nc-1", 5, domain.IntelligencePeriod{Start: time.Now().AddDate(0, 0, -7), End: time.Now()})
	if err != nil {
		t.Fatalf("detect failure: %v", err)
	}
	if risk.Category != "checklist_failure_pattern" || risk.Severity != domain.SeverityHigh || risk.Confidence <= 0 {
		t.Fatalf("risk = %+v", risk)
	}
}

func TestMaintenanceRecurrenceDoesNotCreateOperationalOrder(t *testing.T) {
	_, _, analysis := services()
	risk, err := analysis.DetectMaintenanceRecurrence(context.Background(), "tenant-1", "asset-1", 4, domain.IntelligencePeriod{Start: time.Now().AddDate(0, -1, 0), End: time.Now()})
	if err != nil {
		t.Fatalf("detect recurrence: %v", err)
	}
	if risk.Category != "maintenance_recurrence" {
		t.Fatalf("risk = %+v", risk)
	}
}

func TestInsufficientDataProducesLowConfidenceNoRecommendation(t *testing.T) {
	_, _, analysis := services()
	result, err := analysis.DetectFuelConsumptionDeviation(context.Background(), "tenant-1", []float64{2.8}, 2.35, domain.IntelligencePeriod{Start: time.Now().AddDate(0, 0, -1), End: time.Now()})
	if err != nil {
		t.Fatalf("insufficient data should not error: %v", err)
	}
	if !result.InsufficientData || result.Recommendation != nil || result.Confidence >= 0.6 {
		t.Fatalf("result = %+v", result)
	}
}

func TestInsightDeduplicationAndTenantIsolation(t *testing.T) {
	_, core, _ := services()
	ctx := context.Background()
	insight := domain.Insight{TenantID: "tenant-1", Type: "risk", Title: "Risco", Summary: "Resumo", Category: "fuel", Severity: domain.SeverityHigh, Priority: domain.PriorityHigh, DeduplicationKey: "tenant-1:risk:fuel:asset-1:2026-08:v1"}
	if _, err := core.CreateInsight(ctx, insight); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := core.CreateInsight(ctx, insight); !errors.Is(err, application.ErrDuplicateInsight) {
		t.Fatalf("duplicate err = %v", err)
	}
	insight.TenantID = "tenant-2"
	if _, err := core.CreateInsight(ctx, insight); err != nil {
		t.Fatalf("same key another tenant should pass: %v", err)
	}
}

func TestEventProjectionIsIdempotent(t *testing.T) {
	store, _, _ := services()
	projection := application.NewProjectionService(store.ReadModels())
	event, _ := coreevents.NewDomainEvent("evt-1", coreevents.FuelTransactionCompletedV1, "tenant-1", "fuel-1", "fuel_transaction", coreevents.Payload{"value": 1.0})
	if err := projection.Handle(context.Background(), event); err != nil {
		t.Fatalf("first handle: %v", err)
	}
	if err := projection.Handle(context.Background(), event); err != nil {
		t.Fatalf("duplicate handle should be ignored: %v", err)
	}
	page, err := store.ReadModels().Search(context.Background(), "tenant-1", ports.Query{Page: 1, PerPage: 25})
	if err != nil {
		t.Fatalf("read models: %v", err)
	}
	if page.Total != 1 {
		t.Fatalf("read model total = %d", page.Total)
	}
}

func TestConfidenceAndTrend(t *testing.T) {
	avg, confidence, ok := application.MovingAverage([]float64{2.8, 2.7, 2.9, 2.8})
	if !ok || avg < 2.79 || avg > 2.81 || confidence <= 0.5 {
		t.Fatalf("avg=%f confidence=%f ok=%v", avg, confidence, ok)
	}
	trend := application.Trend([]float64{10, 11, 12}, domain.IntelligencePeriod{})
	if trend.Direction != domain.TrendIncreasing {
		t.Fatalf("trend = %+v", trend)
	}
}

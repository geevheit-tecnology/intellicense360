package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/infrastructure"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/ports"
)

func newTestService() application.CommandService {
	store := infrastructure.NewMemoryStore()
	return application.NewCommandService(store.Items(), store.Events(), store.Actions(), store.Snapshots(), store.Idempotency())
}

func validItem(tenantID string) domain.CommandItem {
	return domain.CommandItem{
		TenantID:    tenantID,
		Type:        domain.TypeRisk,
		Category:    domain.CategoryOperational,
		Severity:    domain.SeverityHigh,
		Priority:    domain.PriorityUrgent,
		Title:       "High operational risk",
		Confidence:  0.9,
		ImpactScore: 0.8,
		RiskScore:   0.85,
	}
}

func TestCommandItemLifecycle(t *testing.T) {
	service := newTestService()
	ctx := context.Background()

	item, err := service.Create(ctx, validItem("tenant-a"), "actor-a")
	if err != nil {
		t.Fatalf("create item: %v", err)
	}

	item, err = service.Acknowledge(ctx, item.TenantID, item.ID, "actor-a")
	if err != nil {
		t.Fatalf("acknowledge item: %v", err)
	}
	if item.Status != domain.StatusAcknowledged || item.AcknowledgedAt == nil {
		t.Fatalf("unexpected acknowledged state: %#v", item)
	}

	item, err = service.Start(ctx, item.TenantID, item.ID, "actor-a")
	if err != nil {
		t.Fatalf("start item: %v", err)
	}
	item, err = service.Resolve(ctx, item.TenantID, item.ID, "actor-a")
	if err != nil {
		t.Fatalf("resolve item: %v", err)
	}
	if item.Status != domain.StatusResolved || item.ResolvedAt == nil {
		t.Fatalf("unexpected resolved state: %#v", item)
	}

	if _, err := service.Acknowledge(ctx, item.TenantID, item.ID, "actor-a"); !errors.Is(err, application.ErrInvalidStatusTransition) {
		t.Fatalf("expected invalid transition after resolved, got %v", err)
	}
}

func TestCreateDedupeAndTenantIsolation(t *testing.T) {
	service := newTestService()
	ctx := context.Background()
	item := validItem("tenant-a")
	item.Fingerprint = "same-risk"

	if _, err := service.Create(ctx, item, "actor-a"); err != nil {
		t.Fatalf("create first item: %v", err)
	}
	if _, err := service.Create(ctx, item, "actor-a"); !errors.Is(err, application.ErrDuplicateCommandItem) {
		t.Fatalf("expected duplicate fingerprint, got %v", err)
	}

	item.TenantID = "tenant-b"
	if _, err := service.Create(ctx, item, "actor-b"); err != nil {
		t.Fatalf("same fingerprint must be allowed for another tenant: %v", err)
	}
}

func TestIdempotencyKeyRejectsDuplicateCreate(t *testing.T) {
	service := newTestService()
	ctx := context.Background()
	item := validItem("tenant-a")
	item.IdempotencyKey = "request-1"

	if _, err := service.Create(ctx, item, "actor-a"); err != nil {
		t.Fatalf("create first request: %v", err)
	}
	item.Title = "Another risk"
	item.Fingerprint = "another-risk"
	if _, err := service.Create(ctx, item, "actor-a"); !errors.Is(err, application.ErrDuplicateCommandItem) {
		t.Fatalf("expected duplicate idempotency key, got %v", err)
	}
}

func TestSLAAndRiskBoundaries(t *testing.T) {
	now := time.Now().UTC()
	dueSoon := now.Add(time.Hour)
	dueLater := now.Add(24 * time.Hour)
	overdue := now.Add(-time.Minute)

	if got := domain.SLA(nil, now); got != domain.SLANotApplicable {
		t.Fatalf("nil due_at SLA = %s", got)
	}
	if got := domain.SLA(&dueSoon, now); got != domain.SLAAtRisk {
		t.Fatalf("due soon SLA = %s", got)
	}
	if got := domain.SLA(&dueLater, now); got != domain.SLAWithin {
		t.Fatalf("due later SLA = %s", got)
	}
	if got := domain.SLA(&overdue, now); got != domain.SLABreached {
		t.Fatalf("overdue SLA = %s", got)
	}
	if got := domain.RiskLevel(0.81); got != "critical" {
		t.Fatalf("risk level = %s", got)
	}
	if got := domain.HealthLevel(39); got != "critical" {
		t.Fatalf("health level = %s", got)
	}
}

func TestRankingOrdersMostUrgentItemsFirst(t *testing.T) {
	items := []domain.CommandItem{
		{ID: "low", Severity: domain.SeverityLow, Priority: domain.PriorityCritical, RiskScore: 1, ImpactScore: 1, SLAStatus: domain.SLABreached, DetectedAt: time.Now()},
		{ID: "critical", Severity: domain.SeverityCritical, Priority: domain.PriorityLow, RiskScore: 0.1, ImpactScore: 0.1, SLAStatus: domain.SLANotApplicable, DetectedAt: time.Now()},
		{ID: "high", Severity: domain.SeverityHigh, Priority: domain.PriorityUrgent, RiskScore: 0.9, ImpactScore: 0.8, SLAStatus: domain.SLAAtRisk, DetectedAt: time.Now()},
	}

	application.SortCommandItems(items)
	if items[0].ID != "critical" || items[1].ID != "high" {
		t.Fatalf("unexpected ranking: %#v", items)
	}
}

func TestSummarySnapshotAndRecommendationsAreDeterministic(t *testing.T) {
	service := newTestService()
	ctx := context.Background()
	item := validItem("tenant-a")
	due := time.Now().UTC().Add(-time.Hour)
	item.DueAt = &due
	item.Severity = domain.SeverityCritical
	item.Priority = domain.PriorityCritical

	created, err := service.Create(ctx, item, "actor-a")
	if err != nil {
		t.Fatalf("create item: %v", err)
	}
	summary, err := service.Summary(ctx, "tenant-a")
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if summary.TotalOpen != 1 || summary.Critical != 1 || summary.BreachedSLA != 1 {
		t.Fatalf("unexpected summary: %#v", summary)
	}

	snapshot, err := service.RebuildSnapshot(ctx, "tenant-a")
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if snapshot.OpenItems != 1 || snapshot.HealthScore >= 100 || snapshot.RiskScore <= 0 {
		t.Fatalf("unexpected snapshot: %#v", snapshot)
	}

	actions, err := service.EvaluateRecommendations(ctx, "tenant-a")
	if err != nil {
		t.Fatalf("recommendations: %v", err)
	}
	if len(actions) != 1 || actions[0].CommandItemID != created.ID || actions[0].Type != domain.ActionReview {
		t.Fatalf("unexpected recommendation actions: %#v", actions)
	}

	history, err := service.History(ctx, "tenant-a", created.ID, ports.Query{})
	if err != nil {
		t.Fatalf("history: %v", err)
	}
	if history.Total != 1 || history.Data[0].EventType != "created" {
		t.Fatalf("unexpected history: %#v", history)
	}
}

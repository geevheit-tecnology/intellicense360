package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/ciot/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/ciot/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/ciot/infrastructure"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/ciot/ports"
)

func TestCIOTLifecycleAllowsOnlyDocumentedTransitions(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	service := application.NewCIOTService(store.CIOTs(), store.History(), store.Errors())
	ctx := context.Background()

	ciot, err := service.Create(ctx, domain.CIOT{TenantID: "tenant-1", Type: domain.TypeTACAgregado, StartDate: time.Now().UTC()})
	if err != nil {
		t.Fatalf("create ciot: %v", err)
	}
	steps := []struct {
		name string
		run  func() (domain.CIOT, error)
		want domain.CIOTStatus
	}{
		{"submit", func() (domain.CIOT, error) { return service.Submit(ctx, "tenant-1", ciot.ID, "actor-1") }, domain.StatusPending},
		{"generated", func() (domain.CIOT, error) {
			return service.MarkGenerated(ctx, "tenant-1", ciot.ID, "PROTO-1", "actor-1")
		}, domain.StatusGenerated},
		{"activate", func() (domain.CIOT, error) { return service.Activate(ctx, "tenant-1", ciot.ID, "actor-1") }, domain.StatusActive},
		{"suspend", func() (domain.CIOT, error) { return service.Suspend(ctx, "tenant-1", ciot.ID, "actor-1") }, domain.StatusSuspended},
		{"reactivate", func() (domain.CIOT, error) { return service.Reactivate(ctx, "tenant-1", ciot.ID, "actor-1") }, domain.StatusActive},
		{"close", func() (domain.CIOT, error) { return service.Close(ctx, "tenant-1", ciot.ID, "actor-1") }, domain.StatusClosed},
	}
	for _, step := range steps {
		t.Run(step.name, func(t *testing.T) {
			saved, err := step.run()
			if err != nil {
				t.Fatalf("transition failed: %v", err)
			}
			if saved.Status != step.want {
				t.Fatalf("status = %s, want %s", saved.Status, step.want)
			}
		})
	}
	history, err := store.History().ListByCIOT(ctx, "tenant-1", ciot.ID, ports.Query{Page: 1, PerPage: 25})
	if err != nil {
		t.Fatalf("history: %v", err)
	}
	if history.Total != 7 {
		t.Fatalf("history total = %d, want 7", history.Total)
	}
}

func TestCIOTRejectsInvalidTransition(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	service := application.NewCIOTService(store.CIOTs(), store.History(), store.Errors())
	ctx := context.Background()
	ciot, err := service.Create(ctx, domain.CIOT{TenantID: "tenant-1", Type: domain.TypeTACIndependente})
	if err != nil {
		t.Fatalf("create ciot: %v", err)
	}
	if _, err := service.Activate(ctx, "tenant-1", ciot.ID, "actor-1"); !errors.Is(err, application.ErrInvalidTransition) {
		t.Fatalf("activate from draft error = %v, want ErrInvalidTransition", err)
	}
}

func TestCIOTIdempotencyKeyIsTenantScoped(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	service := application.NewCIOTService(store.CIOTs(), store.History(), store.Errors())
	ctx := context.Background()
	input := domain.CIOT{TenantID: "tenant-1", Type: domain.TypeOther, IdempotencyKey: "req-1"}
	if _, err := service.Create(ctx, input); err != nil {
		t.Fatalf("create first: %v", err)
	}
	if _, err := service.Create(ctx, input); !errors.Is(err, application.ErrDuplicateRequest) {
		t.Fatalf("duplicate error = %v, want ErrDuplicateRequest", err)
	}
	if _, err := service.Create(ctx, domain.CIOT{TenantID: "tenant-2", Type: domain.TypeOther, IdempotencyKey: "req-1"}); err != nil {
		t.Fatalf("same key in another tenant should be accepted: %v", err)
	}
}

func TestTACAgregadoCanRemainActiveWithoutActualEndDate(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	service := application.NewCIOTService(store.CIOTs(), store.History(), store.Errors())
	ctx := context.Background()
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	expected := time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC)
	ciot, err := service.Create(ctx, domain.CIOT{TenantID: "tenant-1", Type: domain.TypeTACAgregado, StartDate: start, ExpectedEndDate: &expected, OperationalPeriod: "2026-08/2026-12", ContractReference: "CTR-2026-001"})
	if err != nil {
		t.Fatalf("create ciot: %v", err)
	}
	ciot, _ = service.Submit(ctx, "tenant-1", ciot.ID, "actor-1")
	ciot, _ = service.MarkGenerated(ctx, "tenant-1", ciot.ID, "PROTO-TAC", "actor-1")
	ciot, err = service.Activate(ctx, "tenant-1", ciot.ID, "actor-1")
	if err != nil {
		t.Fatalf("activate tac agregado: %v", err)
	}
	if ciot.ActualEndDate != nil {
		t.Fatalf("actual end date should remain open for active TAC Agregado")
	}
}

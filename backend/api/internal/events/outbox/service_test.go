package outbox_test

import (
	"context"
	"errors"
	"testing"
	"time"

	coreevents "github.com/geevheit/intelligence360/backend/api/internal/events"
	"github.com/geevheit/intelligence360/backend/api/internal/events/inmemory"
	"github.com/geevheit/intelligence360/backend/api/internal/events/outbox"
)

func TestOutboxPersistenceMarkPublishedAndTenantIsolation(t *testing.T) {
	repo := outbox.NewMemoryRepository()
	bus := inmemory.NewBus(nil)
	service := outbox.NewService(repo, bus, coreevents.RetryConfig{MaxAttempts: 3, Delay: time.Second})
	event, _ := coreevents.NewDomainEvent("evt-1", coreevents.CIOTCreatedV1, "tenant-1", "ciot-1", "ciot", nil)
	if _, err := service.Save(context.Background(), event); err != nil {
		t.Fatalf("save: %v", err)
	}
	if _, err := service.Save(context.Background(), event); !errors.Is(err, coreevents.ErrDuplicateEvent) {
		t.Fatalf("duplicate err = %v", err)
	}
	if _, err := repo.FindByEventID(context.Background(), "tenant-2", "evt-1"); !errors.Is(err, coreevents.ErrOutboxNotFound) {
		t.Fatalf("tenant isolation err = %v", err)
	}
	processed, err := service.ProcessPending(context.Background(), time.Now().UTC(), 10)
	if err != nil {
		t.Fatalf("process: %v", err)
	}
	if len(processed) != 1 || processed[0].Status != outbox.StatusPublished || processed[0].PublishedAt == nil {
		t.Fatalf("processed = %+v", processed)
	}
}

func TestOutboxRetryAndDeadLetter(t *testing.T) {
	repo := outbox.NewMemoryRepository()
	bus := inmemory.NewBus(nil)
	_ = bus.Subscribe(coreevents.HandlerFunc{HandlerName: "always-fails", EventTypes: []string{coreevents.CIOTClosedV1}, Fn: func(context.Context, coreevents.DomainEvent) error {
		return errors.New("handler failed")
	}})
	service := outbox.NewService(repo, bus, coreevents.RetryConfig{MaxAttempts: 2, Delay: time.Millisecond})
	event, _ := coreevents.NewDomainEvent("evt-2", coreevents.CIOTClosedV1, "tenant-1", "ciot-1", "ciot", nil)
	if _, err := service.Save(context.Background(), event); err != nil {
		t.Fatalf("save: %v", err)
	}
	processed, err := service.ProcessPending(context.Background(), time.Now().UTC(), 10)
	if err != nil {
		t.Fatalf("first process: %v", err)
	}
	if len(processed) != 1 || processed[0].Status != outbox.StatusFailed || processed[0].Attempts != 1 {
		t.Fatalf("first processed = %+v", processed)
	}
	time.Sleep(time.Millisecond)
	if _, err := service.ProcessPending(context.Background(), time.Now().UTC(), 10); err != nil {
		t.Fatalf("second process: %v", err)
	}
	record, err := repo.FindByEventID(context.Background(), "tenant-1", "evt-2")
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if record.Status != outbox.StatusDeadLetter || record.Attempts != 2 {
		t.Fatalf("record = %+v", record)
	}
}

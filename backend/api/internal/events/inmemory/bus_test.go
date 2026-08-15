package inmemory_test

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	coreevents "github.com/geevheit/intelligence360/backend/api/internal/events"
	"github.com/geevheit/intelligence360/backend/api/internal/events/inmemory"
)

func TestDomainEventMetadataSerializationAndVersioning(t *testing.T) {
	event, err := coreevents.NewDomainEvent("evt-1", coreevents.TireRemovedV1, "tenant-1", "tire-1", "tire", coreevents.Payload{"reason": "worn"})
	if err != nil {
		t.Fatalf("create event: %v", err)
	}
	event.CorrelationID = "corr-1"
	event.CausationID = "cause-1"
	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded coreevents.DomainEvent
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if decoded.EventVersion != 1 || decoded.EventType != coreevents.TireRemovedV1 || decoded.CorrelationID != "corr-1" || decoded.CausationID != "cause-1" {
		t.Fatalf("decoded event mismatch: %+v", decoded)
	}
}

func TestInMemoryBusMultipleHandlersFailureIsolationAndOrdering(t *testing.T) {
	bus := inmemory.NewBus(nil)
	calls := []string{}
	event, _ := coreevents.NewDomainEvent("evt-1", coreevents.TireRemovedV1, "tenant-1", "tire-1", "tire", nil)
	_ = bus.Subscribe(coreevents.HandlerFunc{HandlerName: "first", EventTypes: []string{coreevents.TireRemovedV1}, Fn: func(context.Context, coreevents.DomainEvent) error {
		calls = append(calls, "first")
		return nil
	}})
	_ = bus.Subscribe(coreevents.HandlerFunc{HandlerName: "failing", EventTypes: []string{coreevents.TireRemovedV1}, Fn: func(context.Context, coreevents.DomainEvent) error {
		calls = append(calls, "failing")
		return errors.New("boom")
	}})
	_ = bus.Subscribe(coreevents.HandlerFunc{HandlerName: "third", EventTypes: []string{coreevents.TireRemovedV1}, Fn: func(context.Context, coreevents.DomainEvent) error {
		calls = append(calls, "third")
		return nil
	}})
	if err := bus.Publish(context.Background(), event); err == nil {
		t.Fatalf("expected handler failure")
	}
	if !reflect.DeepEqual(calls, []string{"first", "failing", "third"}) {
		t.Fatalf("calls = %v", calls)
	}
}

func TestConsumerStoreIdempotencyAndTenantAgnosticConsumerKey(t *testing.T) {
	store := inmemory.NewConsumerStore()
	ctx := context.Background()
	processed, err := store.HasProcessed(ctx, "consumer-1", "evt-1")
	if err != nil || processed {
		t.Fatalf("initial processed = %v err=%v", processed, err)
	}
	if err := store.MarkProcessed(ctx, "consumer-1", "evt-1"); err != nil {
		t.Fatalf("mark processed: %v", err)
	}
	if err := store.MarkProcessed(ctx, "consumer-1", "evt-1"); !errors.Is(err, coreevents.ErrDuplicateConsumer) {
		t.Fatalf("duplicate mark err = %v", err)
	}
}

func TestIntegrationTireRemovedAndFuelCompletedHandlers(t *testing.T) {
	bus := inmemory.NewBus(nil)
	seen := []string{}
	_ = bus.Subscribe(coreevents.HandlerFunc{HandlerName: "tire-removed-observer", EventTypes: []string{coreevents.TireRemovedV1}, Fn: func(_ context.Context, event coreevents.DomainEvent) error {
		seen = append(seen, event.EventType)
		return nil
	}})
	_ = bus.Subscribe(coreevents.HandlerFunc{HandlerName: "fuel-completed-observer", EventTypes: []string{coreevents.FuelTransactionCompletedV1}, Fn: func(_ context.Context, event coreevents.DomainEvent) error {
		seen = append(seen, event.EventType)
		return nil
	}})
	tireEvent, _ := coreevents.NewDomainEvent("evt-tire", coreevents.TireRemovedV1, "tenant-1", "tire-1", "tire", nil)
	fuelEvent, _ := coreevents.NewDomainEvent("evt-fuel", coreevents.FuelTransactionCompletedV1, "tenant-1", "fuel-1", "fuel_transaction", nil)
	if err := bus.PublishBatch(context.Background(), []coreevents.DomainEvent{tireEvent, fuelEvent}); err != nil {
		t.Fatalf("publish batch: %v", err)
	}
	if !reflect.DeepEqual(seen, []string{coreevents.TireRemovedV1, coreevents.FuelTransactionCompletedV1}) {
		t.Fatalf("seen = %v", seen)
	}
}

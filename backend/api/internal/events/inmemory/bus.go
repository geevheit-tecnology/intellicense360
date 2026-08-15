package inmemory

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	coreevents "github.com/geevheit/intelligence360/backend/api/internal/events"
)

type Bus struct {
	mu       sync.RWMutex
	handlers map[string]coreevents.EventHandler
	order    []string
	logger   *slog.Logger
}

func NewBus(logger *slog.Logger) *Bus {
	if logger == nil {
		logger = slog.Default()
	}
	return &Bus{handlers: map[string]coreevents.EventHandler{}, order: []string{}, logger: logger}
}

func (b *Bus) Subscribe(handler coreevents.EventHandler) error {
	if handler == nil || handler.Name() == "" || len(handler.Handles()) == 0 {
		return coreevents.ErrHandlerFailed
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, exists := b.handlers[handler.Name()]; !exists {
		b.order = append(b.order, handler.Name())
	}
	b.handlers[handler.Name()] = handler
	return nil
}

func (b *Bus) Unsubscribe(handlerName string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.handlers, handlerName)
	next := []string{}
	for _, name := range b.order {
		if name != handlerName {
			next = append(next, name)
		}
	}
	b.order = next
	return nil
}

func (b *Bus) Publish(ctx context.Context, event coreevents.DomainEvent) error {
	if err := event.Validate(); err != nil {
		return err
	}
	handlers := b.matchingHandlers(event.EventType)
	var combined error
	for _, handler := range handlers {
		start := time.Now()
		err := handler.Handle(ctx, event)
		duration := time.Since(start)
		if err != nil {
			b.logger.Error("event failed", "event_id", event.EventID, "event_type", event.EventType, "tenant_id", event.TenantID, "correlation_id", event.CorrelationID, "handler", handler.Name(), "duration", duration.String(), "result", "failed")
			combined = errors.Join(combined, err)
			continue
		}
		b.logger.Info("event handled", "event_id", event.EventID, "event_type", event.EventType, "tenant_id", event.TenantID, "correlation_id", event.CorrelationID, "handler", handler.Name(), "duration", duration.String(), "result", "success")
	}
	if combined != nil {
		return combined
	}
	b.logger.Info("event published", "event_id", event.EventID, "event_type", event.EventType, "tenant_id", event.TenantID, "correlation_id", event.CorrelationID, "result", "success")
	return nil
}

func (b *Bus) PublishBatch(ctx context.Context, batch []coreevents.DomainEvent) error {
	var combined error
	for _, event := range batch {
		if err := b.Publish(ctx, event); err != nil {
			combined = errors.Join(combined, err)
		}
	}
	return combined
}

func (b *Bus) matchingHandlers(eventType string) []coreevents.EventHandler {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := []coreevents.EventHandler{}
	for _, name := range b.order {
		handler, ok := b.handlers[name]
		if !ok {
			continue
		}
		for _, handled := range handler.Handles() {
			if handled == eventType {
				out = append(out, handler)
				break
			}
		}
	}
	return out
}

type ConsumerStore struct {
	mu        sync.RWMutex
	processed map[string]time.Time
}

func NewConsumerStore() *ConsumerStore {
	return &ConsumerStore{processed: map[string]time.Time{}}
}

func (s *ConsumerStore) HasProcessed(_ context.Context, consumerName string, eventID string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.processed[consumerName+":"+eventID]
	return ok, nil
}

func (s *ConsumerStore) MarkProcessed(_ context.Context, consumerName string, eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := consumerName + ":" + eventID
	if _, ok := s.processed[key]; ok {
		return coreevents.ErrDuplicateConsumer
	}
	s.processed[key] = time.Now().UTC()
	return nil
}

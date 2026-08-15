package events

import "context"

type EventBus interface {
	Publish(ctx context.Context, event DomainEvent) error
	PublishBatch(ctx context.Context, batch []DomainEvent) error
	Subscribe(handler EventHandler) error
	Unsubscribe(handlerName string) error
}

type DispatchResult struct {
	EventID   string `json:"event_id"`
	EventType string `json:"event_type"`
	Handler   string `json:"handler"`
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
}

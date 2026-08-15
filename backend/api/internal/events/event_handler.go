package events

import "context"

type EventHandler interface {
	Name() string
	Handles() []string
	Handle(ctx context.Context, event DomainEvent) error
}

type HandlerFunc struct {
	HandlerName string
	EventTypes  []string
	Fn          func(context.Context, DomainEvent) error
}

func (h HandlerFunc) Name() string      { return h.HandlerName }
func (h HandlerFunc) Handles() []string { return h.EventTypes }
func (h HandlerFunc) Handle(ctx context.Context, event DomainEvent) error {
	return h.Fn(ctx, event)
}

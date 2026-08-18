package ports

import (
	"context"

	coreevents "github.com/geevheit/intelligence360/backend/api/internal/events"
)

type EventSource interface {
	Subscribe(handler coreevents.EventHandler) error
	Publish(ctx context.Context, event coreevents.DomainEvent) error
}

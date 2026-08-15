package outbox

import (
	"context"
	"time"

	coreevents "github.com/geevheit/intelligence360/backend/api/internal/events"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusPublished  Status = "published"
	StatusFailed     Status = "failed"
	StatusDeadLetter Status = "dead_letter"
)

type Record struct {
	ID            string                 `json:"id"`
	TenantID      string                 `json:"tenant_id"`
	EventID       string                 `json:"event_id"`
	EventType     string                 `json:"event_type"`
	AggregateID   string                 `json:"aggregate_id"`
	AggregateType string                 `json:"aggregate_type"`
	Payload       coreevents.DomainEvent `json:"payload"`
	OccurredAt    time.Time              `json:"occurred_at"`
	Status        Status                 `json:"status"`
	Attempts      int                    `json:"attempts"`
	AvailableAt   time.Time              `json:"available_at"`
	PublishedAt   *time.Time             `json:"published_at,omitempty"`
	LastError     string                 `json:"last_error,omitempty"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

type DeadLetter struct {
	ID            string                 `json:"id"`
	TenantID      string                 `json:"tenant_id"`
	EventID       string                 `json:"event_id"`
	EventType     string                 `json:"event_type"`
	Payload       coreevents.DomainEvent `json:"payload"`
	FailureReason string                 `json:"failure_reason"`
	Attempts      int                    `json:"attempts"`
	LastError     string                 `json:"last_error"`
	CreatedAt     time.Time              `json:"created_at"`
}

type Repository interface {
	Save(ctx context.Context, record Record) (Record, error)
	GetPending(ctx context.Context, now time.Time, limit int) ([]Record, error)
	MarkProcessing(ctx context.Context, id string) (Record, error)
	MarkPublished(ctx context.Context, id string) (Record, error)
	MarkFailed(ctx context.Context, id string, err error, nextAttempt time.Time) (Record, error)
	MoveToDeadLetter(ctx context.Context, id string, reason string) (DeadLetter, error)
	FindByEventID(ctx context.Context, tenantID string, eventID string) (Record, error)
}

type ConsumerStore interface {
	HasProcessed(ctx context.Context, consumerName string, eventID string) (bool, error)
	MarkProcessed(ctx context.Context, consumerName string, eventID string) error
}

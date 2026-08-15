package events

import (
	"time"

	"github.com/google/uuid"
)

type AggregateMetadata struct {
	ID        uuid.UUID  `json:"id"`
	TenantID  uuid.UUID  `json:"tenant_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Version   int64      `json:"version"`
}

type EventMetadata struct {
	EventID       uuid.UUID         `json:"event_id"`
	EventType     string            `json:"event_type,omitempty"`
	EventVersion  int               `json:"event_version,omitempty"`
	TenantID      uuid.UUID         `json:"tenant_id"`
	AggregateID   uuid.UUID         `json:"aggregate_id"`
	AggregateType string            `json:"aggregate_type,omitempty"`
	OccurredAt    time.Time         `json:"occurred_at"`
	CorrelationID string            `json:"correlation_id,omitempty"`
	CausationID   string            `json:"causation_id,omitempty"`
	ActorID       uuid.UUID         `json:"actor_id,omitempty"`
	ActorType     string            `json:"actor_type,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	Version       int64             `json:"version"`
}

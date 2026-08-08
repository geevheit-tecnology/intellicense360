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
	EventID     uuid.UUID `json:"event_id"`
	TenantID    uuid.UUID `json:"tenant_id"`
	AggregateID uuid.UUID `json:"aggregate_id"`
	OccurredAt  time.Time `json:"occurred_at"`
	Version     int64     `json:"version"`
}

package events

import (
	"encoding/json"
	"time"
)

type Metadata map[string]string
type Payload map[string]any

type DomainEvent struct {
	EventID       string    `json:"event_id"`
	EventType     string    `json:"event_type"`
	EventVersion  int       `json:"event_version"`
	OccurredAt    time.Time `json:"occurred_at"`
	TenantID      string    `json:"tenant_id"`
	AggregateID   string    `json:"aggregate_id"`
	AggregateType string    `json:"aggregate_type"`
	CorrelationID string    `json:"correlation_id,omitempty"`
	CausationID   string    `json:"causation_id,omitempty"`
	ActorID       string    `json:"actor_id,omitempty"`
	ActorType     string    `json:"actor_type,omitempty"`
	Metadata      Metadata  `json:"metadata,omitempty"`
	Payload       Payload   `json:"payload,omitempty"`
}

func NewDomainEvent(eventID, eventType, tenantID, aggregateID, aggregateType string, payload Payload) (DomainEvent, error) {
	if eventID == "" || eventType == "" || tenantID == "" || aggregateID == "" || aggregateType == "" {
		return DomainEvent{}, ErrInvalidEvent
	}
	version, err := VersionFromType(eventType)
	if err != nil {
		return DomainEvent{}, err
	}
	return DomainEvent{EventID: eventID, EventType: eventType, EventVersion: version, OccurredAt: time.Now().UTC(), TenantID: tenantID, AggregateID: aggregateID, AggregateType: aggregateType, Metadata: Metadata{}, Payload: payload}, nil
}

func (e DomainEvent) Validate() error {
	if e.EventID == "" || e.EventType == "" || e.EventVersion <= 0 || e.TenantID == "" || e.AggregateID == "" || e.AggregateType == "" || e.OccurredAt.IsZero() {
		return ErrInvalidEvent
	}
	if _, err := json.Marshal(e.Payload); err != nil {
		return err
	}
	return nil
}

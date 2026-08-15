package events

import "github.com/google/uuid"

const PartCreatedEventName = "inventory.part.created.v1"

type PartCreated struct {
	EventMetadata
	Part         AggregateMetadata `json:"part"`
	SKU          string            `json:"sku"`
	InternalCode string            `json:"internal_code"`
	Name         string            `json:"name"`
	CategoryID   uuid.UUID         `json:"category_id,omitempty"`
	UnitID       uuid.UUID         `json:"unit_id,omitempty"`
	Status       string            `json:"status"`
}

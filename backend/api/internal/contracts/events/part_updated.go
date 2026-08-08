package events

import "github.com/google/uuid"

const PartUpdatedEventName = "part.updated"

type PartUpdated struct {
	EventMetadata
	Part          AggregateMetadata `json:"part"`
	SKU           string            `json:"sku"`
	InternalCode  string            `json:"internal_code"`
	Name          string            `json:"name"`
	CategoryID    uuid.UUID         `json:"category_id,omitempty"`
	UnitID        uuid.UUID         `json:"unit_id,omitempty"`
	Status        string            `json:"status"`
	ChangedFields []string          `json:"changed_fields,omitempty"`
}

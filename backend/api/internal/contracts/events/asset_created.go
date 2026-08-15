package events

import "github.com/google/uuid"

const AssetCreatedEventName = "asset.created.v1"

type AssetCreated struct {
	EventMetadata
	Asset      AggregateMetadata `json:"asset"`
	Code       string            `json:"code"`
	Name       string            `json:"name"`
	CategoryID uuid.UUID         `json:"category_id,omitempty"`
	TypeID     uuid.UUID         `json:"type_id,omitempty"`
	Status     string            `json:"status"`
}

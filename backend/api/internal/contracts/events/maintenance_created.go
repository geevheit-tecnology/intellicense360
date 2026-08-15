package events

import "github.com/google/uuid"

const MaintenanceCreatedEventName = "maintenance.order.created.v1"

type MaintenanceCreated struct {
	EventMetadata
	Maintenance AggregateMetadata `json:"maintenance"`
	Code        string            `json:"code"`
	AssetID     uuid.UUID         `json:"asset_id,omitempty"`
	VehicleID   uuid.UUID         `json:"vehicle_id,omitempty"`
	Kind        string            `json:"kind"`
	Status      string            `json:"status"`
	Priority    string            `json:"priority"`
}

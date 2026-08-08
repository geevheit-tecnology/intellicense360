package events

import (
	"time"

	"github.com/google/uuid"
)

const MaintenanceClosedEventName = "maintenance.closed"

type MaintenanceClosed struct {
	EventMetadata
	Maintenance AggregateMetadata `json:"maintenance"`
	Code        string            `json:"code"`
	AssetID     uuid.UUID         `json:"asset_id,omitempty"`
	VehicleID   uuid.UUID         `json:"vehicle_id,omitempty"`
	Status      string            `json:"status"`
	ClosedAt    time.Time         `json:"closed_at"`
	ClosedByID  uuid.UUID         `json:"closed_by_id,omitempty"`
}

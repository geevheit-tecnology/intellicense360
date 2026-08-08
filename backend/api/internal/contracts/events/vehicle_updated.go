package events

import "github.com/google/uuid"

const VehicleUpdatedEventName = "vehicle.updated"

type VehicleUpdated struct {
	EventMetadata
	Vehicle       AggregateMetadata `json:"vehicle"`
	Plate         string            `json:"plate"`
	Renavam       string            `json:"renavam,omitempty"`
	BrandID       uuid.UUID         `json:"brand_id,omitempty"`
	ModelID       uuid.UUID         `json:"model_id,omitempty"`
	Status        string            `json:"status"`
	ChangedFields []string          `json:"changed_fields,omitempty"`
}

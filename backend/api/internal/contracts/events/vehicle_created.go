package events

import "github.com/google/uuid"

const VehicleCreatedEventName = "vehicle.created.v1"

type VehicleCreated struct {
	EventMetadata
	Vehicle AggregateMetadata `json:"vehicle"`
	Plate   string            `json:"plate"`
	Renavam string            `json:"renavam,omitempty"`
	BrandID uuid.UUID         `json:"brand_id,omitempty"`
	ModelID uuid.UUID         `json:"model_id,omitempty"`
	Status  string            `json:"status"`
}

package events

import "github.com/google/uuid"

const FuelPriceRecordedEventName = "fuel.price.recorded.v1"

type FuelPriceRecorded struct {
	EventMetadata
	FuelTypeID uuid.UUID `json:"fuel_type_id,omitempty"`
	FuelKind   string    `json:"fuel_kind"`
	UnitPrice  float64   `json:"unit_price"`
	StationID  uuid.UUID `json:"station_id,omitempty"`
	Source     string    `json:"source,omitempty"`
}

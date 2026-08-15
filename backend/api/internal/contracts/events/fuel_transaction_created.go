package events

import "github.com/google/uuid"

const FuelTransactionCreatedEventName = "fuel.transaction.created.v1"

type FuelTransactionCreated struct {
	EventMetadata
	Transaction      AggregateMetadata `json:"transaction"`
	FuelTypeID       uuid.UUID         `json:"fuel_type_id,omitempty"`
	FuelKind         string            `json:"fuel_kind"`
	Quantity         float64           `json:"quantity"`
	UnitPrice        float64           `json:"unit_price"`
	TotalAmount      float64           `json:"total_amount"`
	StationID        uuid.UUID         `json:"station_id,omitempty"`
	DriverReference  string            `json:"driver_reference,omitempty"`
	VehicleReference string            `json:"vehicle_reference,omitempty"`
	AssetReference   string            `json:"asset_reference,omitempty"`
	Status           string            `json:"status"`
}

package events

const FuelStationCreatedEventName = "fuel.station.created.v1"

type FuelStationCreated struct {
	EventMetadata
	Station AggregateMetadata `json:"station"`
	Name    string            `json:"name"`
	CNPJ    string            `json:"cnpj,omitempty"`
	Active  bool              `json:"active"`
}

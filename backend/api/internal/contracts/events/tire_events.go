package events

const (
	TireCreatedEventName   = "tire.created.v1"
	TireInstalledEventName = "tire.installed.v1"
	TireRemovedEventName   = "tire.removed.v1"
	TireInspectedEventName = "tire.inspected.v1"
	TireRetreadedEventName = "tire.retreaded.v1"
	TireDamagedEventName   = "tire.damaged.v1"
	TireDisposedEventName  = "tire.disposed.v1"
)

type TireCreated struct {
	EventMetadata
	Tire         AggregateMetadata `json:"tire"`
	SerialNumber string            `json:"serial_number"`
	Status       string            `json:"status"`
}

type TireInstalled struct {
	EventMetadata
	Tire             AggregateMetadata `json:"tire"`
	VehicleReference string            `json:"vehicle_reference,omitempty"`
	Position         string            `json:"position,omitempty"`
	KM               int64             `json:"km,omitempty"`
}

type TireRemoved struct {
	EventMetadata
	Tire   AggregateMetadata `json:"tire"`
	KM     int64             `json:"km,omitempty"`
	Reason string            `json:"reason,omitempty"`
}

type TireInspected struct {
	EventMetadata
	Tire     AggregateMetadata `json:"tire"`
	Result   string            `json:"result,omitempty"`
	Severity string            `json:"severity,omitempty"`
}

type TireRetreaded struct {
	EventMetadata
	Tire   AggregateMetadata `json:"tire"`
	Status string            `json:"status"`
}

type TireDamaged struct {
	EventMetadata
	Tire     AggregateMetadata `json:"tire"`
	Severity string            `json:"severity,omitempty"`
	Reason   string            `json:"reason,omitempty"`
}

type TireDisposed struct {
	EventMetadata
	Tire   AggregateMetadata `json:"tire"`
	Reason string            `json:"reason,omitempty"`
}

package events

const DriverCreatedEventName = "driver.created"

type DriverCreated struct {
	EventMetadata
	Driver   AggregateMetadata `json:"driver"`
	Name     string            `json:"name"`
	Document string            `json:"document,omitempty"`
	Email    string            `json:"email,omitempty"`
	Phone    string            `json:"phone,omitempty"`
	Status   string            `json:"status"`
}

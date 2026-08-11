package events

const FuelReadingRecordedEventName = "fuel.reading.recorded"

type FuelReadingRecorded struct {
	EventMetadata
	ReadingType string  `json:"reading_type"`
	ReferenceID string  `json:"reference_id,omitempty"`
	Value       float64 `json:"value"`
	Source      string  `json:"source,omitempty"`
}

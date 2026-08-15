package events

const FuelTransactionCompletedEventName = "fuel.transaction.completed.v1"

type FuelTransactionCompleted struct {
	EventMetadata
	Transaction AggregateMetadata `json:"transaction"`
	Status      string            `json:"status"`
}

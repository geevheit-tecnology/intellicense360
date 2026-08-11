package events

const FuelTransactionCompletedEventName = "fuel.transaction.completed"

type FuelTransactionCompleted struct {
	EventMetadata
	Transaction AggregateMetadata `json:"transaction"`
	Status      string            `json:"status"`
}

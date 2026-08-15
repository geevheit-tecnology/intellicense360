package events

const FuelTransactionCanceledEventName = "fuel.transaction.canceled.v1"

type FuelTransactionCanceled struct {
	EventMetadata
	Transaction AggregateMetadata `json:"transaction"`
	Reason      string            `json:"reason"`
	Status      string            `json:"status"`
}

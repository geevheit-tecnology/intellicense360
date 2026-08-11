package events

const FuelTransactionCanceledEventName = "fuel.transaction.canceled"

type FuelTransactionCanceled struct {
	EventMetadata
	Transaction AggregateMetadata `json:"transaction"`
	Reason      string            `json:"reason"`
	Status      string            `json:"status"`
}

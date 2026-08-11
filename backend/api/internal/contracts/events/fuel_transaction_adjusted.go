package events

const FuelTransactionAdjustedEventName = "fuel.transaction.adjusted"

type FuelTransactionAdjusted struct {
	EventMetadata
	TransactionID       string  `json:"transaction_id"`
	AdjustmentType      string  `json:"adjustment_type"`
	Reason              string  `json:"reason"`
	AdjustedQuantity    float64 `json:"adjusted_quantity,omitempty"`
	AdjustedUnitPrice   float64 `json:"adjusted_unit_price,omitempty"`
	AdjustedTotalAmount float64 `json:"adjusted_total_amount,omitempty"`
}

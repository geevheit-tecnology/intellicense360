package events

const (
	MaintenanceOrderCreatedEventName   = "maintenance.order.created.v1"
	MaintenanceOrderCompletedEventName = "maintenance.order.completed.v1"
	MaintenanceOrderCanceledEventName  = "maintenance.order.canceled.v1"
)

type MaintenanceOrderCreated struct {
	EventMetadata
	WorkOrder AggregateMetadata `json:"work_order"`
	Number    string            `json:"number"`
	Status    string            `json:"status"`
	Priority  string            `json:"priority,omitempty"`
}

type MaintenanceOrderCompleted struct {
	EventMetadata
	WorkOrder AggregateMetadata `json:"work_order"`
	Status    string            `json:"status"`
}

type MaintenanceOrderCanceled struct {
	EventMetadata
	WorkOrder AggregateMetadata `json:"work_order"`
	Status    string            `json:"status"`
	Reason    string            `json:"reason,omitempty"`
}

package events

const (
	SupplierCreatedEventName = "supplier.created.v1"
	SupplierUpdatedEventName = "supplier.updated.v1"
)

type SupplierCreated struct {
	EventMetadata
	Supplier AggregateMetadata `json:"supplier"`
	Document string            `json:"document,omitempty"`
	Name     string            `json:"name"`
	Status   string            `json:"status"`
}

type SupplierUpdated struct {
	EventMetadata
	Supplier      AggregateMetadata `json:"supplier"`
	Document      string            `json:"document,omitempty"`
	Name          string            `json:"name"`
	Status        string            `json:"status"`
	ChangedFields []string          `json:"changed_fields,omitempty"`
}

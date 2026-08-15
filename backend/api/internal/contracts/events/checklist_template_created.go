package events

const ChecklistTemplateCreatedEventName = "checklist.template.created.v1"

type ChecklistTemplateCreated struct {
	EventMetadata
	Template AggregateMetadata `json:"template"`
	Name     string            `json:"name"`
	Type     string            `json:"type,omitempty"`
	Status   string            `json:"status"`
}

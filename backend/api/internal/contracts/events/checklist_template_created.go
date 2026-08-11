package events

const ChecklistTemplateCreatedEventName = "checklist.template.created"

type ChecklistTemplateCreated struct {
	EventMetadata
	Template AggregateMetadata `json:"template"`
	Name     string            `json:"name"`
	Type     string            `json:"type,omitempty"`
	Status   string            `json:"status"`
}

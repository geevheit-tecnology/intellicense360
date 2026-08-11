package events

const ChecklistTemplateArchivedEventName = "checklist.template.archived"

type ChecklistTemplateArchived struct {
	EventMetadata
	TemplateID string `json:"template_id"`
	Status     string `json:"status"`
}

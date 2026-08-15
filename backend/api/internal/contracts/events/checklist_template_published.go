package events

const ChecklistTemplatePublishedEventName = "checklist.template.published.v1"

type ChecklistTemplatePublished struct {
	EventMetadata
	TemplateID    string `json:"template_id"`
	VersionID     string `json:"version_id"`
	VersionNumber int    `json:"version_number"`
}

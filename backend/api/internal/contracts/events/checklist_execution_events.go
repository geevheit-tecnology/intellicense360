package events

const (
	ChecklistExecutionStartedEventName     = "checklist.execution.started"
	ChecklistResponseRecordedEventName     = "checklist.response.recorded"
	ChecklistEvidenceAddedEventName        = "checklist.evidence.added"
	ChecklistNonConformityCreatedEventName = "checklist.non_conformity.created"
	ChecklistExecutionCompletedEventName   = "checklist.execution.completed"
	ChecklistExecutionCanceledEventName    = "checklist.execution.canceled"
	ChecklistExecutionInvalidatedEventName = "checklist.execution.invalidated"
)

type ChecklistExecutionStarted struct {
	EventMetadata
	Execution         AggregateMetadata `json:"execution"`
	TemplateVersionID string            `json:"template_version_id"`
}

type ChecklistResponseRecorded struct {
	EventMetadata
	ExecutionID string `json:"execution_id"`
	ItemID      string `json:"item_id"`
	Result      string `json:"result,omitempty"`
}

type ChecklistEvidenceAdded struct {
	EventMetadata
	ExecutionID  string `json:"execution_id"`
	ResponseID   string `json:"response_id,omitempty"`
	EvidenceType string `json:"evidence_type"`
}

type ChecklistNonConformityCreated struct {
	EventMetadata
	ExecutionID string `json:"execution_id"`
	Severity    string `json:"severity"`
	Status      string `json:"status"`
}

type ChecklistExecutionCompleted struct{ EventMetadata }
type ChecklistExecutionCanceled struct{ EventMetadata }
type ChecklistExecutionInvalidated struct{ EventMetadata }

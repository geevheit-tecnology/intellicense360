package events

const (
	ChecklistExecutionStartedEventName     = "checklist.execution.started.v1"
	ChecklistResponseRecordedEventName     = "checklist.response.recorded.v1"
	ChecklistEvidenceAddedEventName        = "checklist.evidence.added.v1"
	ChecklistNonConformityCreatedEventName = "checklist.non_conformity.created.v1"
	ChecklistExecutionCompletedEventName   = "checklist.execution.completed.v1"
	ChecklistExecutionCanceledEventName    = "checklist.execution.canceled.v1"
	ChecklistExecutionInvalidatedEventName = "checklist.execution.invalidated.v1"
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

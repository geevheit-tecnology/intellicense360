package domain

import "time"

type ChecklistStatus string
type AnswerType string
type TemplateStatus string
type ExecutionStatus string
type ChecklistSeverity string

const (
	ChecklistStatusDraft      ChecklistStatus = "draft"
	ChecklistStatusInProgress ChecklistStatus = "in_progress"
	ChecklistStatusCompleted  ChecklistStatus = "completed"
	ChecklistStatusCancelled  ChecklistStatus = "cancelled"

	AnswerTypeBoolean   AnswerType = "boolean"
	AnswerTypeText      AnswerType = "text"
	AnswerTypeNumber    AnswerType = "number"
	AnswerTypePhoto     AnswerType = "photo"
	AnswerTypeSignature AnswerType = "signature"
	AnswerTypeSelect    AnswerType = "select"

	TemplateStatusDraft     TemplateStatus = "draft"
	TemplateStatusPublished TemplateStatus = "published"
	TemplateStatusArchived  TemplateStatus = "archived"

	ExecutionStatusDraft       ExecutionStatus = "draft"
	ExecutionStatusInProgress  ExecutionStatus = "in_progress"
	ExecutionStatusCompleted   ExecutionStatus = "completed"
	ExecutionStatusCanceled    ExecutionStatus = "canceled"
	ExecutionStatusInvalidated ExecutionStatus = "invalidated"

	SeverityInfo     ChecklistSeverity = "info"
	SeverityLow      ChecklistSeverity = "low"
	SeverityMedium   ChecklistSeverity = "medium"
	SeverityHigh     ChecklistSeverity = "high"
	SeverityCritical ChecklistSeverity = "critical"
)

type Checklist struct {
	ID             string          `json:"id"`
	TenantID       string          `json:"tenant_id"`
	VehicleID      string          `json:"vehicle_id"`
	Title          string          `json:"title"`
	Description    string          `json:"description,omitempty"`
	Type           string          `json:"type"`
	Status         ChecklistStatus `json:"status"`
	StartedAt      *time.Time      `json:"started_at,omitempty"`
	FinishedAt     *time.Time      `json:"finished_at,omitempty"`
	DriverName     string          `json:"driver_name,omitempty"`
	DriverDocument string          `json:"driver_document,omitempty"`
	CreatedBy      string          `json:"created_by,omitempty"`
	UpdatedBy      string          `json:"updated_by,omitempty"`
	DeletedAt      *time.Time      `json:"deleted_at,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type ChecklistItem struct {
	ID            string     `json:"id"`
	TenantID      string     `json:"tenant_id"`
	ChecklistID   string     `json:"checklist_id"`
	Title         string     `json:"title"`
	Description   string     `json:"description,omitempty"`
	Category      string     `json:"category,omitempty"`
	Required      bool       `json:"required"`
	OrderIndex    int        `json:"order_index"`
	AnswerType    AnswerType `json:"answer_type"`
	ExpectedValue string     `json:"expected_value,omitempty"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type ChecklistAnswer struct {
	ID              string    `json:"id"`
	TenantID        string    `json:"tenant_id"`
	ChecklistID     string    `json:"checklist_id"`
	ChecklistItemID string    `json:"checklist_item_id"`
	Answer          string    `json:"answer"`
	Notes           string    `json:"notes,omitempty"`
	PhotoURL        string    `json:"photo_url,omitempty"`
	AnsweredBy      string    `json:"answered_by,omitempty"`
	AnsweredAt      time.Time `json:"answered_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ChecklistTemplate struct {
	ID                   string         `json:"id"`
	TenantID             string         `json:"tenant_id"`
	Name                 string         `json:"name"`
	Description          string         `json:"description,omitempty"`
	TypeID               string         `json:"type_id,omitempty"`
	Type                 string         `json:"type,omitempty"`
	Status               TemplateStatus `json:"status"`
	Active               bool           `json:"active"`
	CurrentVersionNumber int            `json:"current_version_number"`
	CreatedBy            string         `json:"created_by,omitempty"`
	UpdatedBy            string         `json:"updated_by,omitempty"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            *time.Time     `json:"deleted_at,omitempty"`
	Version              int64          `json:"version"`
}

type ChecklistTemplateVersion struct {
	ID                string         `json:"id"`
	TenantID          string         `json:"tenant_id"`
	TemplateID        string         `json:"template_id"`
	VersionNumber     int            `json:"version_number"`
	Status            TemplateStatus `json:"status"`
	Instructions      string         `json:"instructions,omitempty"`
	ScoringConfig     string         `json:"scoring_config,omitempty"`
	SeverityConfig    string         `json:"severity_config,omitempty"`
	EvidenceRequired  bool           `json:"evidence_required"`
	SignatureRequired bool           `json:"signature_required"`
	PublishedAt       *time.Time     `json:"published_at,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         *time.Time     `json:"deleted_at,omitempty"`
	Version           int64          `json:"version"`
}

type ChecklistType struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	Description string     `json:"description,omitempty"`
	Active      bool       `json:"active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Version     int64      `json:"version"`
}

type ChecklistSection struct {
	ID                string     `json:"id"`
	TenantID          string     `json:"tenant_id"`
	TemplateVersionID string     `json:"template_version_id"`
	Name              string     `json:"name"`
	Description       string     `json:"description,omitempty"`
	OrderIndex        int        `json:"order_index"`
	Active            bool       `json:"active"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
	Version           int64      `json:"version"`
}

type ChecklistEngineItem struct {
	ID                string            `json:"id"`
	TenantID          string            `json:"tenant_id"`
	TemplateVersionID string            `json:"template_version_id"`
	SectionID         string            `json:"section_id,omitempty"`
	Question          string            `json:"question"`
	Description       string            `json:"description,omitempty"`
	ItemType          string            `json:"item_type"`
	Required          bool              `json:"required"`
	OrderIndex        int               `json:"order_index"`
	HelpText          string            `json:"help_text,omitempty"`
	Severity          ChecklistSeverity `json:"severity"`
	EvidenceRequired  bool              `json:"evidence_required"`
	Active            bool              `json:"active"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	DeletedAt         *time.Time        `json:"deleted_at,omitempty"`
	Version           int64             `json:"version"`
}

type ChecklistItemOption struct {
	ID         string     `json:"id"`
	TenantID   string     `json:"tenant_id"`
	ItemID     string     `json:"item_id"`
	Label      string     `json:"label"`
	Value      string     `json:"value"`
	OrderIndex int        `json:"order_index"`
	Active     bool       `json:"active"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	Version    int64      `json:"version"`
}

type ChecklistExecution struct {
	ID                string          `json:"id"`
	TenantID          string          `json:"tenant_id"`
	TemplateVersionID string          `json:"template_version_id"`
	Status            ExecutionStatus `json:"status"`
	PerformedBy       string          `json:"performed_by,omitempty"`
	LocationReference string          `json:"location_reference,omitempty"`
	Notes             string          `json:"notes,omitempty"`
	Score             float64         `json:"score,omitempty"`
	FinalResult       string          `json:"final_result,omitempty"`
	StartedAt         *time.Time      `json:"started_at,omitempty"`
	EndedAt           *time.Time      `json:"ended_at,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	DeletedAt         *time.Time      `json:"deleted_at,omitempty"`
	Version           int64           `json:"version"`
}

type ChecklistResponse struct {
	ID          string            `json:"id"`
	TenantID    string            `json:"tenant_id"`
	ExecutionID string            `json:"execution_id"`
	ItemID      string            `json:"item_id"`
	Value       string            `json:"value"`
	Result      string            `json:"result,omitempty"`
	Notes       string            `json:"notes,omitempty"`
	Responder   string            `json:"responder,omitempty"`
	Severity    ChecklistSeverity `json:"severity,omitempty"`
	RespondedAt time.Time         `json:"responded_at"`
	CreatedAt   time.Time         `json:"created_at"`
}

type ChecklistEvidence struct {
	ID           string    `json:"id"`
	TenantID     string    `json:"tenant_id"`
	ExecutionID  string    `json:"execution_id"`
	ResponseID   string    `json:"response_id,omitempty"`
	EvidenceType string    `json:"evidence_type"`
	Reference    string    `json:"reference"`
	FileName     string    `json:"file_name,omitempty"`
	MimeType     string    `json:"mime_type,omitempty"`
	Notes        string    `json:"notes,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type ChecklistNonConformity struct {
	ID             string            `json:"id"`
	TenantID       string            `json:"tenant_id"`
	ExecutionID    string            `json:"execution_id"`
	ResponseID     string            `json:"response_id,omitempty"`
	Title          string            `json:"title"`
	Description    string            `json:"description,omitempty"`
	Severity       ChecklistSeverity `json:"severity"`
	Status         string            `json:"status"`
	Recommendation string            `json:"recommendation,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	ResolvedAt     *time.Time        `json:"resolved_at,omitempty"`
}

type ChecklistSignature struct {
	ID                 string    `json:"id"`
	TenantID           string    `json:"tenant_id"`
	ExecutionID        string    `json:"execution_id"`
	Signer             string    `json:"signer"`
	SignatureReference string    `json:"signature_reference"`
	SignatureType      string    `json:"signature_type"`
	SignedAt           time.Time `json:"signed_at"`
	CreatedAt          time.Time `json:"created_at"`
}

type ChecklistAssignment struct {
	ID                  string     `json:"id"`
	TenantID            string     `json:"tenant_id"`
	TemplateID          string     `json:"template_id"`
	AssignedToReference string     `json:"assigned_to_reference,omitempty"`
	TargetReference     string     `json:"target_reference,omitempty"`
	Status              string     `json:"status"`
	Notes               string     `json:"notes,omitempty"`
	DueAt               *time.Time `json:"due_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
	Version             int64      `json:"version"`
}

type ChecklistHistory struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	ExecutionID string    `json:"execution_id"`
	Event       string    `json:"event"`
	ActorID     string    `json:"actor_id,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

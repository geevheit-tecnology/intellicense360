package domain

import "time"

type ChecklistStatus string
type AnswerType string

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

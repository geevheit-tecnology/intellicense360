package domain

import "time"

type WorkOrderStatus string
type WorkOrderPriority string
type MaintenanceKind string
type PlanFrequency string
type AttachmentKind string

const (
	WorkOrderDraft      WorkOrderStatus = "draft"
	WorkOrderOpen       WorkOrderStatus = "open"
	WorkOrderWaiting    WorkOrderStatus = "waiting"
	WorkOrderApproved   WorkOrderStatus = "approved"
	WorkOrderExecuting  WorkOrderStatus = "executing"
	WorkOrderInProgress WorkOrderStatus = "in_progress"
	WorkOrderPaused     WorkOrderStatus = "paused"
	WorkOrderCompleted  WorkOrderStatus = "completed"
	WorkOrderCancelled  WorkOrderStatus = "cancelled"
	WorkOrderCanceled   WorkOrderStatus = "canceled"

	PriorityLow      WorkOrderPriority = "low"
	PriorityMedium   WorkOrderPriority = "medium"
	PriorityHigh     WorkOrderPriority = "high"
	PriorityCritical WorkOrderPriority = "critical"

	KindPreventive MaintenanceKind = "preventive"
	KindCorrective MaintenanceKind = "corrective"
	KindPredictive MaintenanceKind = "predictive"
	KindInspection MaintenanceKind = "inspection"
	KindEmergency  MaintenanceKind = "emergency"
	KindWarranty   MaintenanceKind = "warranty"
	KindExternal   MaintenanceKind = "external"
	KindInternal   MaintenanceKind = "internal"

	FrequencyKm    PlanFrequency = "km"
	FrequencyDays  PlanFrequency = "days"
	FrequencyHours PlanFrequency = "hours"

	AttachmentPhoto    AttachmentKind = "photo"
	AttachmentInvoice  AttachmentKind = "invoice"
	AttachmentReport   AttachmentKind = "report"
	AttachmentWarranty AttachmentKind = "warranty"
	AttachmentManual   AttachmentKind = "manual"
)

type WorkOrder struct {
	ID             string            `json:"id"`
	TenantID       string            `json:"tenant_id"`
	Code           string            `json:"code"`
	AssetID        string            `json:"asset_id,omitempty"`
	VehicleID      string            `json:"vehicle_id,omitempty"`
	Title          string            `json:"title"`
	Description    string            `json:"description,omitempty"`
	Kind           MaintenanceKind   `json:"kind"`
	Status         WorkOrderStatus   `json:"status"`
	Priority       WorkOrderPriority `json:"priority"`
	CategoryID     string            `json:"category_id,omitempty"`
	WorkshopID     string            `json:"workshop_id,omitempty"`
	TechnicianID   string            `json:"technician_id,omitempty"`
	ReasonID       string            `json:"reason_id,omitempty"`
	ServiceTypeID  string            `json:"service_type_id,omitempty"`
	OpenedAt       time.Time         `json:"opened_at"`
	ScheduledAt    *time.Time        `json:"scheduled_at,omitempty"`
	StartedAt      *time.Time        `json:"started_at,omitempty"`
	CompletedAt    *time.Time        `json:"completed_at,omitempty"`
	CancelledAt    *time.Time        `json:"cancelled_at,omitempty"`
	EstimatedHours float64           `json:"estimated_hours,omitempty"`
	ActualHours    float64           `json:"actual_hours,omitempty"`
	CreatedBy      string            `json:"created_by,omitempty"`
	UpdatedBy      string            `json:"updated_by,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	DeletedAt      *time.Time        `json:"deleted_at,omitempty"`
	Version        int64             `json:"version"`
}

type MaintenanceOrder = WorkOrder
type MaintenanceStatus = WorkOrderStatus
type MaintenancePriority = WorkOrderPriority
type MaintenanceType = MaintenanceKind

type PreventivePlan struct {
	ID            string        `json:"id"`
	TenantID      string        `json:"tenant_id"`
	Name          string        `json:"name"`
	Description   string        `json:"description,omitempty"`
	AssetID       string        `json:"asset_id,omitempty"`
	VehicleID     string        `json:"vehicle_id,omitempty"`
	ServiceTypeID string        `json:"service_type_id,omitempty"`
	Frequency     PlanFrequency `json:"frequency"`
	IntervalValue int64         `json:"interval_value"`
	NextDueAt     *time.Time    `json:"next_due_at,omitempty"`
	NextDueValue  int64         `json:"next_due_value,omitempty"`
	Active        bool          `json:"active"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	DeletedAt     *time.Time    `json:"deleted_at,omitempty"`
	Version       int64         `json:"version"`
}

type CorrectiveRecord struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	WorkOrderID string     `json:"work_order_id"`
	FailureMode string     `json:"failure_mode"`
	RootCause   string     `json:"root_cause,omitempty"`
	Resolution  string     `json:"resolution,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type ServiceType struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	Description string     `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Version     int64      `json:"version"`
}

type MaintenanceCatalog struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	Description string     `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Version     int64      `json:"version"`
}

type MaintenanceCategory = MaintenanceCatalog
type MaintenanceReason = MaintenanceCatalog
type Priority = MaintenanceCatalog

type Workshop struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	Name      string     `json:"name"`
	Document  string     `json:"document,omitempty"`
	Phone     string     `json:"phone,omitempty"`
	Email     string     `json:"email,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Version   int64      `json:"version"`
}

type Technician struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	Name      string     `json:"name"`
	Document  string     `json:"document,omitempty"`
	Phone     string     `json:"phone,omitempty"`
	Email     string     `json:"email,omitempty"`
	Active    bool       `json:"active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Version   int64      `json:"version"`
}

type ServiceItem struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	WorkOrderID string     `json:"work_order_id"`
	Description string     `json:"description"`
	Quantity    float64    `json:"quantity"`
	UnitCost    float64    `json:"unit_cost,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type MaintenanceSchedule struct {
	ID               string     `json:"id"`
	TenantID         string     `json:"tenant_id"`
	WorkOrderID      string     `json:"work_order_id,omitempty"`
	PreventivePlanID string     `json:"preventive_plan_id,omitempty"`
	ScheduledAt      time.Time  `json:"scheduled_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

type MaintenanceAttachment struct {
	ID          string         `json:"id"`
	TenantID    string         `json:"tenant_id"`
	WorkOrderID string         `json:"work_order_id"`
	Kind        AttachmentKind `json:"kind"`
	FileName    string         `json:"file_name"`
	MimeType    string         `json:"mime_type,omitempty"`
	URI         string         `json:"uri,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
}

type LaborEntry struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	WorkOrderID string     `json:"work_order_id"`
	Technician  string     `json:"technician"`
	Hours       float64    `json:"hours"`
	Cost        float64    `json:"cost,omitempty"`
	StartedAt   time.Time  `json:"started_at"`
	FinishedAt  *time.Time `json:"finished_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type Downtime struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	WorkOrderID string     `json:"work_order_id"`
	AssetID     string     `json:"asset_id,omitempty"`
	VehicleID   string     `json:"vehicle_id,omitempty"`
	StartedAt   time.Time  `json:"started_at"`
	EndedAt     *time.Time `json:"ended_at,omitempty"`
	Reason      string     `json:"reason,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type AvailabilitySnapshot struct {
	ID              string    `json:"id"`
	TenantID        string    `json:"tenant_id"`
	AssetID         string    `json:"asset_id,omitempty"`
	VehicleID       string    `json:"vehicle_id,omitempty"`
	Available       bool      `json:"available"`
	AvailabilityPct float64   `json:"availability_pct"`
	CapturedAt      time.Time `json:"captured_at"`
	CreatedAt       time.Time `json:"created_at"`
}

type MaintenanceHistory struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	WorkOrderID string    `json:"work_order_id"`
	Event       string    `json:"event"`
	ActorID     string    `json:"actor_id,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

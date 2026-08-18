package domain

import "time"

type CenterStatus string
type CommandItemType string
type Category string
type Severity string
type Priority string
type CommandStatus string
type SLAStatus string
type ActionType string
type ActionStatus string

const (
	CenterActive   CenterStatus = "active"
	CenterPaused   CenterStatus = "paused"
	CenterArchived CenterStatus = "archived"

	TypeAlert          CommandItemType = "alert"
	TypeRisk           CommandItemType = "risk"
	TypeIncident       CommandItemType = "incident"
	TypeOpportunity    CommandItemType = "opportunity"
	TypeRecommendation CommandItemType = "recommendation"
	TypeTask           CommandItemType = "task"
	TypeAnomaly        CommandItemType = "anomaly"
	TypeWarning        CommandItemType = "warning"
	TypeInsight        CommandItemType = "insight"

	CategoryOperational Category = "operational"
	CategoryMaintenance Category = "maintenance"
	CategoryFleet       Category = "fleet"
	CategoryTire        Category = "tire"
	CategoryFuel        Category = "fuel"
	CategoryInventory   Category = "inventory"
	CategoryFinancial   Category = "financial"
	CategoryCompliance  Category = "compliance"
	CategoryDriver      Category = "driver"
	CategoryDocument    Category = "document"
	CategoryCIOT        Category = "ciot"
	CategorySafety      Category = "safety"
	CategoryPerformance Category = "performance"
	CategoryCost        Category = "cost"

	SeverityInfo     Severity = "info"
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"

	PriorityLow      Priority = "low"
	PriorityNormal   Priority = "normal"
	PriorityHigh     Priority = "high"
	PriorityUrgent   Priority = "urgent"
	PriorityCritical Priority = "critical"

	StatusOpen         CommandStatus = "open"
	StatusAcknowledged CommandStatus = "acknowledged"
	StatusInProgress   CommandStatus = "in_progress"
	StatusResolved     CommandStatus = "resolved"
	StatusDismissed    CommandStatus = "dismissed"
	StatusExpired      CommandStatus = "expired"

	SLAWithin        SLAStatus = "within_sla"
	SLAAtRisk        SLAStatus = "at_risk"
	SLABreached      SLAStatus = "breached"
	SLANotApplicable SLAStatus = "not_applicable"

	ActionAcknowledge ActionType = "acknowledge"
	ActionAssign      ActionType = "assign"
	ActionStart       ActionType = "start"
	ActionResolve     ActionType = "resolve"
	ActionDismiss     ActionType = "dismiss"
	ActionEscalate    ActionType = "escalate"
	ActionReview      ActionType = "review"

	ActionPending   ActionStatus = "pending"
	ActionCompleted ActionStatus = "completed"
	ActionCanceled  ActionStatus = "canceled"
)

type CommandCenter struct {
	ID               string       `json:"id"`
	TenantID         string       `json:"tenant_id"`
	Name             string       `json:"name"`
	Status           CenterStatus `json:"status"`
	LastCalculatedAt *time.Time   `json:"last_calculated_at,omitempty"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
}

type CommandItem struct {
	ID             string            `json:"id"`
	TenantID       string            `json:"tenant_id"`
	Type           CommandItemType   `json:"type"`
	Category       Category          `json:"category"`
	Severity       Severity          `json:"severity"`
	Priority       Priority          `json:"priority"`
	Status         CommandStatus     `json:"status"`
	Title          string            `json:"title"`
	Description    string            `json:"description,omitempty"`
	Source         string            `json:"source,omitempty"`
	SourceType     string            `json:"source_type,omitempty"`
	SourceID       string            `json:"source_id,omitempty"`
	Confidence     float64           `json:"confidence"`
	ImpactScore    float64           `json:"impact_score"`
	RiskScore      float64           `json:"risk_score"`
	UrgencyScore   float64           `json:"urgency_score"`
	DueAt          *time.Time        `json:"due_at,omitempty"`
	SLAStatus      SLAStatus         `json:"sla_status"`
	SLAHours       int               `json:"sla_hours,omitempty"`
	DetectedAt     time.Time         `json:"detected_at"`
	AcknowledgedAt *time.Time        `json:"acknowledged_at,omitempty"`
	ResolvedAt     *time.Time        `json:"resolved_at,omitempty"`
	AssignedTo     string            `json:"assigned_to,omitempty"`
	Fingerprint    string            `json:"fingerprint"`
	IdempotencyKey string            `json:"idempotency_key,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	Version        int64             `json:"version"`
}

type CommandAction struct {
	ID            string            `json:"id"`
	TenantID      string            `json:"tenant_id"`
	CommandItemID string            `json:"command_item_id"`
	Type          ActionType        `json:"type"`
	Label         string            `json:"label"`
	Status        ActionStatus      `json:"status"`
	Priority      Priority          `json:"priority"`
	Payload       map[string]string `json:"payload,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type CommandEvent struct {
	ID             string            `json:"id"`
	TenantID       string            `json:"tenant_id"`
	CommandItemID  string            `json:"command_item_id"`
	EventType      string            `json:"event_type"`
	PreviousStatus CommandStatus     `json:"previous_status,omitempty"`
	NewStatus      CommandStatus     `json:"new_status,omitempty"`
	ActorID        string            `json:"actor_id,omitempty"`
	Payload        map[string]string `json:"payload,omitempty"`
	OccurredAt     time.Time         `json:"occurred_at"`
}

type OperationalSnapshot struct {
	ID                    string    `json:"id"`
	TenantID              string    `json:"tenant_id"`
	SnapshotAt            time.Time `json:"snapshot_at"`
	OpenItems             int       `json:"open_items"`
	CriticalItems         int       `json:"critical_items"`
	HighPriorityItems     int       `json:"high_priority_items"`
	ActiveRisks           int       `json:"active_risks"`
	ActiveAlerts          int       `json:"active_alerts"`
	OpenIncidents         int       `json:"open_incidents"`
	Opportunities         int       `json:"opportunities"`
	BreachedSLAs          int       `json:"breached_slas"`
	AverageResolutionTime float64   `json:"average_resolution_time"`
	OperationalScore      float64   `json:"operational_score"`
	RiskScore             float64   `json:"risk_score"`
	HealthScore           float64   `json:"health_score"`
}

type MissionControlSummary struct {
	TotalOpen        int       `json:"total_open"`
	Critical         int       `json:"critical"`
	High             int       `json:"high"`
	Medium           int       `json:"medium"`
	Risks            int       `json:"risks"`
	Alerts           int       `json:"alerts"`
	Incidents        int       `json:"incidents"`
	Opportunities    int       `json:"opportunities"`
	Recommendations  int       `json:"recommendations"`
	BreachedSLA      int       `json:"breached_sla"`
	OperationalScore float64   `json:"operational_score"`
	RiskScore        float64   `json:"risk_score"`
	HealthScore      float64   `json:"health_score"`
	LastUpdated      time.Time `json:"last_updated"`
}

type RiskAggregation struct {
	OverallRiskScore float64       `json:"overall_risk_score"`
	RiskLevel        string        `json:"risk_level"`
	TopRisks         []CommandItem `json:"top_risks"`
}

type MissionControlRule struct {
	ID          string   `json:"id"`
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Enabled     bool     `json:"enabled"`
	Priority    Priority `json:"priority"`
	Version     int      `json:"version"`
}

type Signal struct {
	TenantID   string            `json:"tenant_id"`
	Type       CommandItemType   `json:"type"`
	Category   Category          `json:"category"`
	SourceType string            `json:"source_type"`
	SourceID   string            `json:"source_id"`
	Title      string            `json:"title"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

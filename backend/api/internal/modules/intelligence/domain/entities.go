package domain

import "time"

type MetricType string
type Severity string
type InsightStatus string
type Priority string
type TrendDirection string
type RuleType string

const (
	MetricOperational MetricType = "operational"
	MetricFinancial   MetricType = "financial"
	MetricFuel        MetricType = "fuel"
	MetricTire        MetricType = "tire"
	MetricMaintenance MetricType = "maintenance"
	MetricChecklist   MetricType = "checklist"
	MetricInventory   MetricType = "inventory"
	MetricCIOT        MetricType = "ciot"
	MetricAsset       MetricType = "asset"
	MetricEfficiency  MetricType = "efficiency"
	MetricRisk        MetricType = "risk"
	MetricCost        MetricType = "cost"

	SeverityInfo     Severity = "info"
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"

	StatusNew          InsightStatus = "new"
	StatusAcknowledged InsightStatus = "acknowledged"
	StatusInProgress   InsightStatus = "in_progress"
	StatusResolved     InsightStatus = "resolved"
	StatusDismissed    InsightStatus = "dismissed"
	StatusExpired      InsightStatus = "expired"

	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"

	TrendIncreasing       TrendDirection = "increasing"
	TrendDecreasing       TrendDirection = "decreasing"
	TrendStable           TrendDirection = "stable"
	TrendVolatile         TrendDirection = "volatile"
	TrendInsufficientData TrendDirection = "insufficient_data"

	RuleThreshold RuleType = "threshold"
	RuleDeviation RuleType = "deviation"
	RuleTrend     RuleType = "trend"
	RuleFrequency RuleType = "frequency"
	RulePattern   RuleType = "pattern"
	RuleCost      RuleType = "cost"
)

type IntelligencePeriod struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

type IntelligenceMetric struct {
	ID             string            `json:"id"`
	TenantID       string            `json:"tenant_id"`
	MetricType     MetricType        `json:"metric_type"`
	Name           string            `json:"name"`
	Value          float64           `json:"value"`
	Unit           string            `json:"unit"`
	PeriodStart    time.Time         `json:"period_start"`
	PeriodEnd      time.Time         `json:"period_end"`
	Dimension      string            `json:"dimension,omitempty"`
	DimensionValue string            `json:"dimension_value,omitempty"`
	Source         string            `json:"source,omitempty"`
	CalculatedAt   time.Time         `json:"calculated_at"`
	Confidence     float64           `json:"confidence"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	Version        int64             `json:"version"`
}

type MetricSnapshot struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	MetricID  string    `json:"metric_id"`
	Value     float64   `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}

type MetricDefinition struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	MetricType  MetricType `json:"metric_type"`
	Name        string     `json:"name"`
	Unit        string     `json:"unit"`
	Description string     `json:"description,omitempty"`
	Active      bool       `json:"active"`
	Version     int64      `json:"version"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type InsightEvidence struct {
	ID            string             `json:"id,omitempty"`
	Source        string             `json:"source"`
	SourceType    string             `json:"source_type"`
	ReferenceID   string             `json:"reference_id,omitempty"`
	Metric        string             `json:"metric,omitempty"`
	ObservedValue float64            `json:"observed_value,omitempty"`
	ExpectedValue float64            `json:"expected_value,omitempty"`
	Period        IntelligencePeriod `json:"period"`
	Timestamp     time.Time          `json:"timestamp"`
	Explanation   string             `json:"explanation"`
}

type InsightImpact struct {
	EstimatedCost        float64 `json:"estimated_cost,omitempty"`
	EstimatedSaving      float64 `json:"estimated_saving,omitempty"`
	OperationalImpact    string  `json:"operational_impact,omitempty"`
	RiskImpact           string  `json:"risk_impact,omitempty"`
	Confidence           float64 `json:"confidence"`
	Currency             string  `json:"currency,omitempty"`
	EstimatedImpactRange string  `json:"estimated_impact_range,omitempty"`
}

type Anomaly struct {
	ID                  string             `json:"id"`
	TenantID            string             `json:"tenant_id"`
	Type                string             `json:"type"`
	Severity            Severity           `json:"severity"`
	MetricID            string             `json:"metric_id,omitempty"`
	ObservedValue       float64            `json:"observed_value"`
	ExpectedValue       float64            `json:"expected_value"`
	Deviation           float64            `json:"deviation"`
	DeviationPercentage float64            `json:"deviation_percentage"`
	DetectedAt          time.Time          `json:"detected_at"`
	Period              IntelligencePeriod `json:"period"`
	Evidence            []InsightEvidence  `json:"evidence"`
	Confidence          float64            `json:"confidence"`
	Status              InsightStatus      `json:"status"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
	Version             int64              `json:"version"`
}

type Risk struct {
	ID          string            `json:"id"`
	TenantID    string            `json:"tenant_id"`
	Category    string            `json:"category"`
	Severity    Severity          `json:"severity"`
	Probability float64           `json:"probability"`
	Impact      InsightImpact     `json:"impact"`
	Confidence  float64           `json:"confidence"`
	Evidence    []InsightEvidence `json:"evidence"`
	DetectedAt  time.Time         `json:"detected_at"`
	Status      InsightStatus     `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Version     int64             `json:"version"`
}

type Opportunity struct {
	ID              string            `json:"id"`
	TenantID        string            `json:"tenant_id"`
	Category        string            `json:"category"`
	EstimatedImpact InsightImpact     `json:"estimated_impact"`
	EstimatedSaving float64           `json:"estimated_saving,omitempty"`
	Confidence      float64           `json:"confidence"`
	Evidence        []InsightEvidence `json:"evidence"`
	DetectedAt      time.Time         `json:"detected_at"`
	ValidUntil      *time.Time        `json:"valid_until,omitempty"`
	Status          InsightStatus     `json:"status"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	Version         int64             `json:"version"`
}

type Recommendation struct {
	ID              string            `json:"id"`
	TenantID        string            `json:"tenant_id,omitempty"`
	Title           string            `json:"title"`
	Description     string            `json:"description,omitempty"`
	WhatHappened    string            `json:"what_happened,omitempty"`
	WhyItMatters    string            `json:"why_it_matters,omitempty"`
	Evidence        []InsightEvidence `json:"evidence,omitempty"`
	SuggestedAction string            `json:"suggested_action,omitempty"`
	ExpectedImpact  InsightImpact     `json:"expected_impact"`
	Confidence      float64           `json:"confidence"`
	Priority        Priority          `json:"priority"`
	ImpactArea      string            `json:"impact_area,omitempty"`
	Status          InsightStatus     `json:"status,omitempty"`
	CreatedAt       time.Time         `json:"created_at,omitempty"`
	UpdatedAt       time.Time         `json:"updated_at,omitempty"`
	Version         int64             `json:"version,omitempty"`
}

type Insight struct {
	ID               string            `json:"id"`
	TenantID         string            `json:"tenant_id"`
	Type             string            `json:"type"`
	Title            string            `json:"title"`
	Summary          string            `json:"summary"`
	Category         string            `json:"category"`
	Severity         Severity          `json:"severity"`
	Evidence         []InsightEvidence `json:"evidence"`
	MetricID         string            `json:"metric_id,omitempty"`
	AnomalyID        string            `json:"anomaly_id,omitempty"`
	RiskID           string            `json:"risk_id,omitempty"`
	OpportunityID    string            `json:"opportunity_id,omitempty"`
	RecommendationID string            `json:"recommendation_id,omitempty"`
	EstimatedImpact  InsightImpact     `json:"estimated_impact"`
	Confidence       float64           `json:"confidence"`
	Priority         Priority          `json:"priority"`
	DeduplicationKey string            `json:"deduplication_key"`
	Status           InsightStatus     `json:"status"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	Version          int64             `json:"version"`
}

type IntelligenceRule struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id,omitempty"`
	Name        string    `json:"name"`
	RuleType    RuleType  `json:"rule_type"`
	Version     int       `json:"version"`
	Threshold   float64   `json:"threshold,omitempty"`
	Window      string    `json:"window,omitempty"`
	Active      bool      `json:"active"`
	Explanation string    `json:"explanation,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type IntelligenceThreshold struct {
	ID         string     `json:"id"`
	TenantID   string     `json:"tenant_id,omitempty"`
	RuleID     string     `json:"rule_id"`
	MetricType MetricType `json:"metric_type"`
	Warning    float64    `json:"warning"`
	Critical   float64    `json:"critical"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type Trend struct {
	Direction  TrendDirection     `json:"direction"`
	Magnitude  float64            `json:"magnitude"`
	Period     IntelligencePeriod `json:"period"`
	Confidence float64            `json:"confidence"`
}

type ReadModel struct {
	ID          string            `json:"id"`
	TenantID    string            `json:"tenant_id"`
	Source      string            `json:"source"`
	SourceType  string            `json:"source_type"`
	ReferenceID string            `json:"reference_id"`
	MetricName  string            `json:"metric_name"`
	Value       float64           `json:"value"`
	Unit        string            `json:"unit,omitempty"`
	PeriodStart time.Time         `json:"period_start"`
	PeriodEnd   time.Time         `json:"period_end"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	LastEventID string            `json:"last_event_id,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Version     int64             `json:"version"`
}

type AnalysisResult struct {
	Metric           IntelligenceMetric `json:"metric"`
	Anomaly          *Anomaly           `json:"anomaly,omitempty"`
	Risk             *Risk              `json:"risk,omitempty"`
	Opportunity      *Opportunity       `json:"opportunity,omitempty"`
	Recommendation   *Recommendation    `json:"recommendation,omitempty"`
	Insight          *Insight           `json:"insight,omitempty"`
	Trend            Trend              `json:"trend"`
	Confidence       float64            `json:"confidence"`
	InsufficientData bool               `json:"insufficient_data"`
}

type OperationalSummary struct {
	CurrentOperation string           `json:"current_operation"`
	CurrentRisks     []string         `json:"current_risks"`
	CurrentSavings   string           `json:"current_savings"`
	Priorities       []string         `json:"priorities"`
	Recommendations  []Recommendation `json:"recommendations"`
}

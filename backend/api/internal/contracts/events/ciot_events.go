package events

import "time"

type CIOTCreated struct {
	EventID    string    `json:"event_id"`
	TenantID   string    `json:"tenant_id"`
	CIOTID     string    `json:"ciot_id"`
	CIOTType   string    `json:"ciot_type"`
	Status     string    `json:"status"`
	OccurredAt time.Time `json:"occurred_at"`
}

type CIOTSubmitted struct {
	EventID    string    `json:"event_id"`
	TenantID   string    `json:"tenant_id"`
	CIOTID     string    `json:"ciot_id"`
	OccurredAt time.Time `json:"occurred_at"`
}

type CIOTGenerated struct {
	EventID          string    `json:"event_id"`
	TenantID         string    `json:"tenant_id"`
	CIOTID           string    `json:"ciot_id"`
	ExternalProtocol string    `json:"external_protocol,omitempty"`
	OccurredAt       time.Time `json:"occurred_at"`
}

type CIOTActivated struct {
	EventID    string    `json:"event_id"`
	TenantID   string    `json:"tenant_id"`
	CIOTID     string    `json:"ciot_id"`
	OccurredAt time.Time `json:"occurred_at"`
}

type CIOTSuspended struct {
	EventID    string    `json:"event_id"`
	TenantID   string    `json:"tenant_id"`
	CIOTID     string    `json:"ciot_id"`
	OccurredAt time.Time `json:"occurred_at"`
}

type CIOTClosed struct {
	EventID    string    `json:"event_id"`
	TenantID   string    `json:"tenant_id"`
	CIOTID     string    `json:"ciot_id"`
	OccurredAt time.Time `json:"occurred_at"`
}

type CIOTCanceled struct {
	EventID    string    `json:"event_id"`
	TenantID   string    `json:"tenant_id"`
	CIOTID     string    `json:"ciot_id"`
	OccurredAt time.Time `json:"occurred_at"`
}

type CIOTProviderAttemptRecorded struct {
	EventID       string    `json:"event_id"`
	TenantID      string    `json:"tenant_id"`
	CIOTID        string    `json:"ciot_id"`
	AttemptID     string    `json:"attempt_id"`
	Provider      string    `json:"provider"`
	AttemptNumber int       `json:"attempt_number"`
	Status        string    `json:"status"`
	OccurredAt    time.Time `json:"occurred_at"`
}

type CIOTErrorRecorded struct {
	EventID    string    `json:"event_id"`
	TenantID   string    `json:"tenant_id"`
	CIOTID     string    `json:"ciot_id"`
	Code       string    `json:"code"`
	OccurredAt time.Time `json:"occurred_at"`
}

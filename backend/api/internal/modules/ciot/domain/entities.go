package domain

import "time"

type CIOTStatus string
type CIOTType string
type CIOTHistoryEvent string

const (
	StatusDraft     CIOTStatus = "draft"
	StatusPending   CIOTStatus = "pending"
	StatusGenerated CIOTStatus = "generated"
	StatusActive    CIOTStatus = "active"
	StatusSuspended CIOTStatus = "suspended"
	StatusClosed    CIOTStatus = "closed"
	StatusCanceled  CIOTStatus = "canceled"
	StatusError     CIOTStatus = "error"

	TypeTACAgregado     CIOTType = "tac_agregado"
	TypeTACIndependente CIOTType = "tac_independente"
	TypeOther           CIOTType = "other"

	EventCreated           CIOTHistoryEvent = "created"
	EventSubmitted         CIOTHistoryEvent = "submitted"
	EventGenerated         CIOTHistoryEvent = "generated"
	EventActivated         CIOTHistoryEvent = "activated"
	EventSuspended         CIOTHistoryEvent = "suspended"
	EventReactivated       CIOTHistoryEvent = "reactivated"
	EventClosed            CIOTHistoryEvent = "closed"
	EventCanceled          CIOTHistoryEvent = "canceled"
	EventProviderAttempted CIOTHistoryEvent = "provider_attempted"
	EventProviderSucceeded CIOTHistoryEvent = "provider_succeeded"
	EventProviderFailed    CIOTHistoryEvent = "provider_failed"
	EventPaymentRecorded   CIOTHistoryEvent = "payment_recorded"
	EventErrorOccurred     CIOTHistoryEvent = "error_occurred"
)

type CIOT struct {
	ID                 string     `json:"id"`
	TenantID           string     `json:"tenant_id"`
	CIOTNumber         string     `json:"ciot_number,omitempty"`
	Type               CIOTType   `json:"type"`
	Status             CIOTStatus `json:"status"`
	ContractID         string     `json:"contract_id,omitempty"`
	CarrierID          string     `json:"carrier_id,omitempty"`
	TransporterID      string     `json:"transporter_id,omitempty"`
	OperationID        string     `json:"operation_id,omitempty"`
	VehicleReferenceID string     `json:"vehicle_reference_id,omitempty"`
	DriverReferenceID  string     `json:"driver_reference_id,omitempty"`
	AmountID           string     `json:"amount_id,omitempty"`
	StartDate          time.Time  `json:"start_date"`
	ExpectedEndDate    *time.Time `json:"expected_end_date,omitempty"`
	ActualEndDate      *time.Time `json:"actual_end_date,omitempty"`
	OperationalPeriod  string     `json:"operational_period,omitempty"`
	ContractReference  string     `json:"contract_reference,omitempty"`
	ExternalProtocol   string     `json:"external_protocol,omitempty"`
	IdempotencyKey     string     `json:"idempotency_key,omitempty"`
	RequestHash        string     `json:"request_hash,omitempty"`
	Notes              string     `json:"notes,omitempty"`
	ErrorCode          string     `json:"error_code,omitempty"`
	ErrorMessage       string     `json:"error_message,omitempty"`
	CreatedBy          string     `json:"created_by,omitempty"`
	UpdatedBy          string     `json:"updated_by,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
	Version            int64      `json:"version"`
}

type CIOTContract struct {
	ID             string     `json:"id"`
	TenantID       string     `json:"tenant_id"`
	ContractNumber string     `json:"contract_number"`
	ContractType   string     `json:"contract_type,omitempty"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	Status         string     `json:"status"`
	Notes          string     `json:"notes,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	Version        int64      `json:"version"`
}

type CIOTCarrier struct {
	ID               string     `json:"id"`
	TenantID         string     `json:"tenant_id"`
	Document         string     `json:"document"`
	LegalName        string     `json:"legal_name"`
	TradeName        string     `json:"trade_name,omitempty"`
	Registration     string     `json:"registration,omitempty"`
	ContactReference string     `json:"contact_reference,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
	Version          int64      `json:"version"`
}

type CIOTTransporter struct {
	ID                string     `json:"id"`
	TenantID          string     `json:"tenant_id"`
	Document          string     `json:"document"`
	Name              string     `json:"name"`
	Registration      string     `json:"registration,omitempty"`
	ContractReference string     `json:"contract_reference,omitempty"`
	ContactReference  string     `json:"contact_reference,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
	Version           int64      `json:"version"`
}

type CIOTOperation struct {
	ID              string     `json:"id"`
	TenantID        string     `json:"tenant_id"`
	OperationNumber string     `json:"operation_number"`
	Origin          string     `json:"origin,omitempty"`
	Destination     string     `json:"destination,omitempty"`
	StartDate       time.Time  `json:"start_date"`
	ExpectedEndDate *time.Time `json:"expected_end_date,omitempty"`
	ActualEndDate   *time.Time `json:"actual_end_date,omitempty"`
	CargoReference  string     `json:"cargo_reference,omitempty"`
	Weight          float64    `json:"weight,omitempty"`
	Distance        float64    `json:"distance,omitempty"`
	Notes           string     `json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at,omitempty"`
	Version         int64      `json:"version"`
}

type CIOTVehicleReference struct {
	ID           string     `json:"id"`
	TenantID     string     `json:"tenant_id"`
	VehicleID    string     `json:"vehicle_id,omitempty"`
	LicensePlate string     `json:"license_plate,omitempty"`
	VehicleType  string     `json:"vehicle_type,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	Version      int64      `json:"version"`
}

type CIOTDriverReference struct {
	ID                string     `json:"id"`
	TenantID          string     `json:"tenant_id"`
	DriverID          string     `json:"driver_id,omitempty"`
	NameReference     string     `json:"name_reference,omitempty"`
	DocumentReference string     `json:"document_reference,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
	Version           int64      `json:"version"`
}

type CIOTAmount struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenant_id"`
	CIOTID        string    `json:"ciot_id,omitempty"`
	FreightAmount float64   `json:"freight_amount"`
	AdvanceAmount float64   `json:"advance_amount,omitempty"`
	BalanceAmount float64   `json:"balance_amount,omitempty"`
	TollAmount    float64   `json:"toll_amount,omitempty"`
	OtherAmount   float64   `json:"other_amount,omitempty"`
	TotalAmount   float64   `json:"total_amount"`
	Currency      string    `json:"currency"`
	CreatedAt     time.Time `json:"created_at"`
}

type CIOTPayment struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	CIOTID      string     `json:"ciot_id"`
	PaymentType string     `json:"payment_type"`
	Amount      float64    `json:"amount"`
	DueDate     time.Time  `json:"due_date"`
	PaymentDate *time.Time `json:"payment_date,omitempty"`
	Status      string     `json:"status"`
	Reference   string     `json:"reference,omitempty"`
	Notes       string     `json:"notes,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Version     int64      `json:"version"`
}

type CIOTStatusHistory struct {
	ID         string            `json:"id"`
	TenantID   string            `json:"tenant_id"`
	CIOTID     string            `json:"ciot_id"`
	Event      CIOTHistoryEvent  `json:"event"`
	FromStatus CIOTStatus        `json:"from_status,omitempty"`
	ToStatus   CIOTStatus        `json:"to_status,omitempty"`
	Reason     string            `json:"reason,omitempty"`
	ActorID    string            `json:"actor_id,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
}

type CIOTProviderAttempt struct {
	ID                string     `json:"id"`
	TenantID          string     `json:"tenant_id"`
	CIOTID            string     `json:"ciot_id"`
	Provider          string     `json:"provider"`
	AttemptNumber     int        `json:"attempt_number"`
	RequestedAt       time.Time  `json:"requested_at"`
	CompletedAt       *time.Time `json:"completed_at,omitempty"`
	Status            string     `json:"status"`
	RequestReference  string     `json:"request_reference,omitempty"`
	ResponseReference string     `json:"response_reference,omitempty"`
	ErrorCode         string     `json:"error_code,omitempty"`
	ErrorMessage      string     `json:"error_message,omitempty"`
	HTTPStatus        int        `json:"http_status,omitempty"`
	LatencyMS         int64      `json:"latency_ms,omitempty"`
	IdempotencyKey    string     `json:"idempotency_key,omitempty"`
	RequestHash       string     `json:"request_hash,omitempty"`
	ResultReference   string     `json:"result_reference,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
}

type CIOTExternalReference struct {
	ID                 string     `json:"id"`
	TenantID           string     `json:"tenant_id"`
	CIOTID             string     `json:"ciot_id"`
	Provider           string     `json:"provider"`
	ExternalCIOTNumber string     `json:"external_ciot_number,omitempty"`
	ExternalProtocol   string     `json:"external_protocol,omitempty"`
	ExternalStatus     string     `json:"external_status,omitempty"`
	GeneratedAt        *time.Time `json:"generated_at,omitempty"`
	LastSynchronizedAt *time.Time `json:"last_synchronized_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	Version            int64      `json:"version"`
}

type CIOTDocument struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	CIOTID    string    `json:"ciot_id"`
	Type      string    `json:"type"`
	Number    string    `json:"number,omitempty"`
	Reference string    `json:"reference,omitempty"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type CIOTError struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	CIOTID    string    `json:"ciot_id"`
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	Details   string    `json:"details,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

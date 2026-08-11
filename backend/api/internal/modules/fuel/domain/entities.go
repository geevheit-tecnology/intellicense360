package domain

import "time"

type FuelKind string
type FuelTransactionStatus string
type FuelTankStatus string
type FuelNozzleStatus string
type FuelReadingType string

const (
	FuelKindDieselS10  FuelKind = "diesel_s10"
	FuelKindDieselS500 FuelKind = "diesel_s500"
	FuelKindGasoline   FuelKind = "gasoline"
	FuelKindEthanol    FuelKind = "ethanol"
	FuelKindGNV        FuelKind = "gnv"
	FuelKindARLA32     FuelKind = "arla_32"
	FuelKindOther      FuelKind = "other"

	FuelTransactionDraft     FuelTransactionStatus = "draft"
	FuelTransactionCompleted FuelTransactionStatus = "completed"
	FuelTransactionCanceled  FuelTransactionStatus = "canceled"
	FuelTransactionAdjusted  FuelTransactionStatus = "adjusted"
	FuelTransactionRejected  FuelTransactionStatus = "rejected"

	FuelTankActive      FuelTankStatus = "active"
	FuelTankInactive    FuelTankStatus = "inactive"
	FuelTankMaintenance FuelTankStatus = "maintenance"

	FuelNozzleActive      FuelNozzleStatus = "active"
	FuelNozzleInactive    FuelNozzleStatus = "inactive"
	FuelNozzleMaintenance FuelNozzleStatus = "maintenance"

	FuelReadingOdometer   FuelReadingType = "odometer"
	FuelReadingEngineHour FuelReadingType = "engine_hours"
	FuelReadingTank       FuelReadingType = "fuel_tank"
)

type FuelTransaction struct {
	ID                 string                `json:"id"`
	TenantID           string                `json:"tenant_id"`
	TransactionDate    time.Time             `json:"transaction_date"`
	FuelTypeID         string                `json:"fuel_type_id"`
	FuelKind           FuelKind              `json:"fuel_kind"`
	Quantity           float64               `json:"quantity"`
	UnitPrice          float64               `json:"unit_price"`
	TotalAmount        float64               `json:"total_amount"`
	OdometerReading    float64               `json:"odometer_reading,omitempty"`
	EngineHourReading  float64               `json:"engine_hour_reading,omitempty"`
	StationID          string                `json:"station_id,omitempty"`
	TankID             string                `json:"tank_id,omitempty"`
	NozzleID           string                `json:"nozzle_id,omitempty"`
	ReceiptID          string                `json:"receipt_id,omitempty"`
	ReceiptNumber      string                `json:"receipt_number,omitempty"`
	DriverReference    string                `json:"driver_reference,omitempty"`
	VehicleReference   string                `json:"vehicle_reference,omitempty"`
	AssetReference     string                `json:"asset_reference,omitempty"`
	PaymentMethod      string                `json:"payment_method,omitempty"`
	Notes              string                `json:"notes,omitempty"`
	Status             FuelTransactionStatus `json:"status"`
	CancellationReason string                `json:"cancellation_reason,omitempty"`
	CompletedAt        *time.Time            `json:"completed_at,omitempty"`
	CanceledAt         *time.Time            `json:"canceled_at,omitempty"`
	CreatedBy          string                `json:"created_by,omitempty"`
	UpdatedBy          string                `json:"updated_by,omitempty"`
	CreatedAt          time.Time             `json:"created_at"`
	UpdatedAt          time.Time             `json:"updated_at"`
	DeletedAt          *time.Time            `json:"deleted_at,omitempty"`
	Version            int64                 `json:"version"`
}

type FuelType struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	Name        string     `json:"name"`
	Kind        FuelKind   `json:"kind"`
	Code        string     `json:"code,omitempty"`
	Description string     `json:"description,omitempty"`
	Active      bool       `json:"active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Version     int64      `json:"version"`
}

type FuelStation struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	Name      string     `json:"name"`
	LegalName string     `json:"legal_name,omitempty"`
	CNPJ      string     `json:"cnpj,omitempty"`
	Address   string     `json:"address,omitempty"`
	City      string     `json:"city,omitempty"`
	State     string     `json:"state,omitempty"`
	Country   string     `json:"country,omitempty"`
	Latitude  float64    `json:"latitude,omitempty"`
	Longitude float64    `json:"longitude,omitempty"`
	Active    bool       `json:"active"`
	Notes     string     `json:"notes,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Version   int64      `json:"version"`
}

type FuelTank struct {
	ID             string         `json:"id"`
	TenantID       string         `json:"tenant_id"`
	Code           string         `json:"code"`
	Name           string         `json:"name,omitempty"`
	Capacity       float64        `json:"capacity"`
	CurrentReading float64        `json:"current_reading"`
	FuelTypeID     string         `json:"fuel_type_id,omitempty"`
	FuelKind       FuelKind       `json:"fuel_kind"`
	LocationRef    string         `json:"location_reference,omitempty"`
	Status         FuelTankStatus `json:"status"`
	Notes          string         `json:"notes,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      *time.Time     `json:"deleted_at,omitempty"`
	Version        int64          `json:"version"`
}

type FuelNozzle struct {
	ID           string           `json:"id"`
	TenantID     string           `json:"tenant_id"`
	Code         string           `json:"code"`
	FuelTypeID   string           `json:"fuel_type_id,omitempty"`
	FuelKind     FuelKind         `json:"fuel_kind"`
	TankID       string           `json:"tank_id,omitempty"`
	Status       FuelNozzleStatus `json:"status"`
	MeterReading float64          `json:"meter_reading"`
	Notes        string           `json:"notes,omitempty"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	DeletedAt    *time.Time       `json:"deleted_at,omitempty"`
	Version      int64            `json:"version"`
}

type FuelReading struct {
	ID          string          `json:"id"`
	TenantID    string          `json:"tenant_id"`
	ReadingType FuelReadingType `json:"reading_type"`
	ReferenceID string          `json:"reference_id,omitempty"`
	Value       float64         `json:"value"`
	ReadingDate time.Time       `json:"reading_date"`
	Source      string          `json:"source,omitempty"`
	Notes       string          `json:"notes,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
}

type FuelPrice struct {
	ID            string     `json:"id"`
	TenantID      string     `json:"tenant_id"`
	FuelTypeID    string     `json:"fuel_type_id,omitempty"`
	FuelKind      FuelKind   `json:"fuel_kind"`
	UnitPrice     float64    `json:"unit_price"`
	EffectiveDate time.Time  `json:"effective_date"`
	StationID     string     `json:"station_id,omitempty"`
	Source        string     `json:"source,omitempty"`
	Notes         string     `json:"notes,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
	Version       int64      `json:"version"`
}

type FuelReceipt struct {
	ID                  string     `json:"id"`
	TenantID            string     `json:"tenant_id"`
	ReceiptNumber       string     `json:"receipt_number"`
	ReceiptDate         time.Time  `json:"receipt_date"`
	Amount              float64    `json:"amount"`
	AttachmentReference string     `json:"attachment_reference,omitempty"`
	Notes               string     `json:"notes,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
	Version             int64      `json:"version"`
}

type FuelAttachment struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	ReferenceID string    `json:"reference_id,omitempty"`
	Kind        string    `json:"kind"`
	FileName    string    `json:"file_name"`
	MimeType    string    `json:"mime_type,omitempty"`
	URI         string    `json:"uri,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type FuelAdjustment struct {
	ID                  string    `json:"id"`
	TenantID            string    `json:"tenant_id"`
	TransactionID       string    `json:"transaction_id"`
	AdjustmentType      string    `json:"adjustment_type"`
	Reason              string    `json:"reason"`
	OriginalReference   string    `json:"original_reference,omitempty"`
	AdjustedQuantity    float64   `json:"adjusted_quantity,omitempty"`
	AdjustedUnitPrice   float64   `json:"adjusted_unit_price,omitempty"`
	AdjustedTotalAmount float64   `json:"adjusted_total_amount,omitempty"`
	Notes               string    `json:"notes,omitempty"`
	CreatedBy           string    `json:"created_by,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
}

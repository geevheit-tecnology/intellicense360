package domain

import "time"

type TireStatus string
type TireCondition string
type MovementType string

const (
	TireStatusNew             TireStatus = "new"
	TireStatusInStock         TireStatus = "in_stock"
	TireStatusInstalled       TireStatus = "installed"
	TireStatusRemoved         TireStatus = "removed"
	TireStatusUnderInspection TireStatus = "under_inspection"
	TireStatusUnderRetread    TireStatus = "under_retread"
	TireStatusRetreaded       TireStatus = "retreaded"
	TireStatusReserved        TireStatus = "reserved"
	TireStatusDamaged         TireStatus = "damaged"
	TireStatusEndOfLife       TireStatus = "end_of_life"
	TireStatusDisposed        TireStatus = "disposed"
	TireStatusLost            TireStatus = "lost"
	TireStatusRecapping       TireStatus = "recapping"

	ConditionExcellent TireCondition = "excellent"
	ConditionGood      TireCondition = "good"
	ConditionAttention TireCondition = "attention"
	ConditionHeavyWear TireCondition = "heavy_wear"
	ConditionAtLimit   TireCondition = "at_limit"
	ConditionCritical  TireCondition = "critical"
	ConditionDamaged   TireCondition = "damaged"
	ConditionEndOfLife TireCondition = "end_of_life"

	MovementPurchase     MovementType = "purchase"
	MovementReceipt      MovementType = "receipt"
	MovementInstallation MovementType = "installation"
	MovementRotation     MovementType = "rotation"
	MovementRemoval      MovementType = "removal"
	MovementTransfer     MovementType = "transfer"
	MovementRetread      MovementType = "retread"
	MovementReturn       MovementType = "return"
	MovementRecapping    MovementType = "recapping"
	MovementRepair       MovementType = "repair"
	MovementDisposal     MovementType = "disposal"
	MovementAdjustment   MovementType = "adjustment"
	MovementLoss         MovementType = "loss"
)

type Tire struct {
	ID                string        `json:"id"`
	TenantID          string        `json:"tenant_id"`
	SerialNumber      string        `json:"serial_number"`
	FireNumber        string        `json:"fire_number"`
	Brand             string        `json:"brand"`
	BrandID           string        `json:"brand_id,omitempty"`
	Model             string        `json:"model"`
	ModelID           string        `json:"model_id,omitempty"`
	Size              string        `json:"size"`
	Dimension         string        `json:"dimension,omitempty"`
	Construction      string        `json:"construction"`
	LoadIndex         string        `json:"load_index,omitempty"`
	SpeedRating       string        `json:"speed_rating,omitempty"`
	TireType          string        `json:"tire_type"`
	PositionType      string        `json:"position_type"`
	ManufacturingDate *time.Time    `json:"manufacturing_date,omitempty"`
	PurchaseDate      *time.Time    `json:"purchase_date,omitempty"`
	PurchaseValue     float64       `json:"purchase_value,omitempty"`
	Supplier          string        `json:"supplier,omitempty"`
	SupplierReference string        `json:"supplier_reference,omitempty"`
	Warranty          string        `json:"warranty,omitempty"`
	DOT               string        `json:"dot"`
	CurrentTreadMM    float64       `json:"current_tread_mm"`
	OriginalTreadMM   float64       `json:"original_tread_mm"`
	MinimumTreadMM    float64       `json:"minimum_tread_mm"`
	Status            TireStatus    `json:"status"`
	Condition         TireCondition `json:"condition,omitempty"`
	VehicleID         string        `json:"vehicle_id,omitempty"`
	Position          string        `json:"position,omitempty"`
	CurrentKM         int64         `json:"current_km,omitempty"`
	TotalKM           int64         `json:"total_km,omitempty"`
	RecapCount        int           `json:"recap_count,omitempty"`
	Notes             string        `json:"notes,omitempty"`
	CreatedBy         string        `json:"created_by,omitempty"`
	UpdatedBy         string        `json:"updated_by,omitempty"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
	DeletedAt         *time.Time    `json:"deleted_at,omitempty"`
	Version           int64         `json:"version"`
}

type TireBrand struct {
	ID, TenantID, Name, Code, Description string
	CreatedAt, UpdatedAt                  time.Time
	DeletedAt                             *time.Time
	Version                               int64
}
type TireModel struct {
	ID, TenantID, BrandID, Name, Code, Description string
	CreatedAt, UpdatedAt                           time.Time
	DeletedAt                                      *time.Time
	Version                                        int64
}
type TireSpecification struct {
	ID, TenantID, Dimension, Construction, LoadIndex, SpeedRating string
	OriginalTreadDepthMM                                          float64
	CreatedAt, UpdatedAt                                          time.Time
	DeletedAt                                                     *time.Time
	Version                                                       int64
}
type TirePosition struct {
	ID, TenantID, Axle, Side, Position, InnerOuter, PositionCode, PositionLabel string
	CreatedAt, UpdatedAt                                                        time.Time
	DeletedAt                                                                   *time.Time
	Version                                                                     int64
}

type TireInstallation struct {
	ID                   string       `json:"id"`
	TenantID             string       `json:"tenant_id"`
	TireID               string       `json:"tire_id"`
	InstallationDate     time.Time    `json:"installation_date"`
	InstallationPosition TirePosition `json:"installation_position"`
	InitialKM            int64        `json:"initial_km"`
	InitialTreadDepth    float64      `json:"initial_tread_depth"`
	InstallationReason   string       `json:"installation_reason,omitempty"`
	Notes                string       `json:"notes,omitempty"`
	CreatedAt            time.Time    `json:"created_at"`
}

type TireRemoval struct {
	ID                  string        `json:"id"`
	TenantID            string        `json:"tenant_id"`
	TireID              string        `json:"tire_id"`
	RemovalDate         time.Time     `json:"removal_date"`
	RemovalPosition     TirePosition  `json:"removal_position"`
	RemovalKM           int64         `json:"removal_km"`
	RemainingTreadDepth float64       `json:"remaining_tread_depth"`
	RemovalReason       string        `json:"removal_reason,omitempty"`
	Condition           TireCondition `json:"condition"`
	Notes               string        `json:"notes,omitempty"`
	CreatedAt           time.Time     `json:"created_at"`
}

type TireInspection struct {
	ID                string     `json:"id"`
	TenantID          string     `json:"tenant_id"`
	TireID            string     `json:"tire_id"`
	InspectionDate    time.Time  `json:"inspection_date"`
	TreadMM           float64    `json:"tread_mm"`
	Pressure          float64    `json:"pressure,omitempty"`
	Temperature       float64    `json:"temperature,omitempty"`
	Condition         string     `json:"condition"`
	VisualCondition   string     `json:"visual_condition,omitempty"`
	SidewallCondition string     `json:"sidewall_condition,omitempty"`
	ShoulderCondition string     `json:"shoulder_condition,omitempty"`
	CrownCondition    string     `json:"crown_condition,omitempty"`
	IrregularWear     bool       `json:"irregular_wear,omitempty"`
	Damage            bool       `json:"damage,omitempty"`
	Cracks            bool       `json:"cracks,omitempty"`
	Cuts              bool       `json:"cuts,omitempty"`
	ObjectPenetration bool       `json:"object_penetration,omitempty"`
	Recommendation    string     `json:"recommendation,omitempty"`
	Observations      string     `json:"observations,omitempty"`
	Inspector         string     `json:"inspector,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

type TireMeasurement struct {
	ID                  string    `json:"id"`
	TenantID            string    `json:"tenant_id"`
	TireID              string    `json:"tire_id"`
	TreadDepthMM        float64   `json:"tread_depth_mm"`
	Pressure            float64   `json:"pressure,omitempty"`
	MeasurementPosition string    `json:"measurement_position,omitempty"`
	MeasurementDate     time.Time `json:"measurement_date"`
	CreatedAt           time.Time `json:"created_at"`
}

type TireRetread struct {
	ID                      string    `json:"id"`
	TenantID                string    `json:"tenant_id"`
	TireID                  string    `json:"tire_id"`
	RetreadNumber           int       `json:"retread_number"`
	RetreadDate             time.Time `json:"retread_date"`
	ProviderReference       string    `json:"provider_reference,omitempty"`
	TreadBrand              string    `json:"tread_brand,omitempty"`
	TreadModel              string    `json:"tread_model,omitempty"`
	Cost                    float64   `json:"cost,omitempty"`
	BeforeRetreadTreadDepth float64   `json:"before_retread_tread_depth,omitempty"`
	AfterRetreadTreadDepth  float64   `json:"after_retread_tread_depth,omitempty"`
	Warranty                string    `json:"warranty,omitempty"`
	Result                  string    `json:"result,omitempty"`
	Status                  string    `json:"status,omitempty"`
	Notes                   string    `json:"notes,omitempty"`
	CreatedAt               time.Time `json:"created_at"`
}

type TireRetreadEvent struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	TireID    string    `json:"tire_id"`
	RetreadID string    `json:"retread_id"`
	Event     string    `json:"event"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
type TireHistory struct {
	ID         string     `json:"id"`
	TenantID   string     `json:"tenant_id"`
	TireID     string     `json:"tire_id"`
	Event      string     `json:"event"`
	FromStatus TireStatus `json:"from_status,omitempty"`
	ToStatus   TireStatus `json:"to_status,omitempty"`
	ActorID    string     `json:"actor_id,omitempty"`
	Notes      string     `json:"notes,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}
type TireCost struct {
	ID                                                                     string `json:"id"`
	TenantID                                                               string `json:"tenant_id"`
	TireID                                                                 string `json:"tire_id"`
	PurchaseCost, RetreadCost, RepairCost, OtherCost, TotalCost, CostPerKM float64
	CreatedAt, UpdatedAt                                                   time.Time
	DeletedAt                                                              *time.Time
	Version                                                                int64
}
type TireDisposal struct {
	ID                  string    `json:"id"`
	TenantID            string    `json:"tenant_id"`
	TireID              string    `json:"tire_id"`
	DisposalDate        time.Time `json:"disposal_date"`
	Reason              string    `json:"reason,omitempty"`
	AttachmentReference string    `json:"attachment_reference,omitempty"`
	CreatedAt           time.Time `json:"created_at"`
}
type TireAttachment struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	TireID    string    `json:"tire_id"`
	Kind      string    `json:"kind"`
	FileName  string    `json:"file_name"`
	MimeType  string    `json:"mime_type,omitempty"`
	URI       string    `json:"uri,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type TireMovement struct {
	ID           string       `json:"id"`
	TenantID     string       `json:"tenant_id"`
	TireID       string       `json:"tire_id"`
	MovementType MovementType `json:"movement_type"`
	VehicleID    string       `json:"vehicle_id,omitempty"`
	Position     string       `json:"position,omitempty"`
	KM           int64        `json:"km,omitempty"`
	Reason       string       `json:"reason,omitempty"`
	PerformedBy  string       `json:"performed_by,omitempty"`
	MovementDate time.Time    `json:"movement_date"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	DeletedAt    *time.Time   `json:"deleted_at,omitempty"`
}

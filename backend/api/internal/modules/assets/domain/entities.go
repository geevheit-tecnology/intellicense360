package domain

import "time"

type AssetStatus string
type Ownership string
type AttachmentKind string

const (
	AssetStatusDraft       AssetStatus = "draft"
	AssetStatusAvailable   AssetStatus = "available"
	AssetStatusAssigned    AssetStatus = "assigned"
	AssetStatusInOperation AssetStatus = "in_operation"
	AssetStatusMaintenance AssetStatus = "maintenance"
	AssetStatusInactive    AssetStatus = "inactive"
	AssetStatusSold        AssetStatus = "sold"
	AssetStatusDisposed    AssetStatus = "disposed"

	OwnershipOwned      Ownership = "owned"
	OwnershipLeased     Ownership = "leased"
	OwnershipRented     Ownership = "rented"
	OwnershipThirdParty Ownership = "third_party"

	AttachmentImage    AttachmentKind = "image"
	AttachmentPDF      AttachmentKind = "pdf"
	AttachmentDocument AttachmentKind = "document"
	AttachmentWarranty AttachmentKind = "warranty"
	AttachmentInvoice  AttachmentKind = "invoice"
	AttachmentManual   AttachmentKind = "manual"
)

type Asset struct {
	ID             string            `json:"id"`
	TenantID       string            `json:"tenant_id"`
	AuditID        string            `json:"audit_id,omitempty"`
	InternalCode   string            `json:"internal_code"`
	SerialNumber   string            `json:"serial_number,omitempty"`
	AssetTag       string            `json:"asset_tag"`
	Name           string            `json:"name"`
	Description    string            `json:"description,omitempty"`
	CategoryID     string            `json:"category_id"`
	TypeID         string            `json:"type_id"`
	ManufacturerID string            `json:"manufacturer_id,omitempty"`
	ModelID        string            `json:"model_id,omitempty"`
	Status         AssetStatus       `json:"status"`
	Ownership      Ownership         `json:"ownership"`
	Location       Location          `json:"location"`
	Depreciation   Depreciation      `json:"depreciation"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	DeletedAt      *time.Time        `json:"deleted_at,omitempty"`
	Version        int64             `json:"version"`
}

type AssetCategory struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	Name      string     `json:"name"`
	Code      string     `json:"code"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Version   int64      `json:"version"`
}

type AssetType struct {
	ID         string     `json:"id"`
	TenantID   string     `json:"tenant_id"`
	CategoryID string     `json:"category_id"`
	Name       string     `json:"name"`
	Code       string     `json:"code"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	Version    int64      `json:"version"`
}

type Manufacturer struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Version   int64      `json:"version"`
}

type Model struct {
	ID             string     `json:"id"`
	TenantID       string     `json:"tenant_id"`
	ManufacturerID string     `json:"manufacturer_id"`
	Name           string     `json:"name"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	Version        int64      `json:"version"`
}

type Equipment struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	AssetID   string     `json:"asset_id"`
	Category  string     `json:"category"`
	Type      string     `json:"type"`
	Capacity  string     `json:"capacity,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Version   int64      `json:"version"`
}

type Implement struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	AssetID   string     `json:"asset_id"`
	Type      string     `json:"type"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Version   int64      `json:"version"`
}

type Depreciation struct {
	PurchaseValue    float64 `json:"purchase_value,omitempty"`
	ResidualValue    float64 `json:"residual_value,omitempty"`
	UsefulLifeMonths int     `json:"useful_life_months,omitempty"`
}

type Location struct {
	Site    string `json:"site,omitempty"`
	Area    string `json:"area,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Country string `json:"country,omitempty"`
}

type Attachment struct {
	ID        string         `json:"id"`
	TenantID  string         `json:"tenant_id"`
	AssetID   string         `json:"asset_id"`
	Kind      AttachmentKind `json:"kind"`
	FileName  string         `json:"file_name"`
	MimeType  string         `json:"mime_type,omitempty"`
	URI       string         `json:"uri,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
}

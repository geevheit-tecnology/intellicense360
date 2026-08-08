package domain

import "time"

type StockStatus string

const (
	StockStatusActive   StockStatus = "active"
	StockStatusInactive StockStatus = "inactive"
)

type Part struct {
	ID           string            `json:"id"`
	TenantID     string            `json:"tenant_id"`
	AuditID      string            `json:"audit_id,omitempty"`
	SKU          string            `json:"sku"`
	InternalCode string            `json:"internal_code"`
	Name         string            `json:"name"`
	Description  string            `json:"description,omitempty"`
	CategoryID   string            `json:"category_id,omitempty"`
	BrandID      string            `json:"brand_id,omitempty"`
	ModelID      string            `json:"model_id,omitempty"`
	UnitID       string            `json:"unit_id"`
	Status       StockStatus       `json:"status"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	DeletedAt    *time.Time        `json:"deleted_at,omitempty"`
	Version      int64             `json:"version"`
}

type Catalog struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	AuditID     string     `json:"audit_id,omitempty"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	Description string     `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Version     int64      `json:"version"`
}

type PartCategory = Catalog
type PartBrand = Catalog
type PartModel = Catalog
type UnitOfMeasure = Catalog

type SupplierReference struct {
	ID           string     `json:"id"`
	TenantID     string     `json:"tenant_id"`
	PartID       string     `json:"part_id"`
	SupplierID   string     `json:"supplier_id"`
	SupplierCode string     `json:"supplier_code,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	Version      int64      `json:"version"`
}

type Warehouse struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	AuditID   string     `json:"audit_id,omitempty"`
	Name      string     `json:"name"`
	Code      string     `json:"code"`
	Address   string     `json:"address,omitempty"`
	Active    bool       `json:"active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Version   int64      `json:"version"`
}

type WarehouseLocation struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	WarehouseID string     `json:"warehouse_id"`
	Name        string     `json:"name"`
	Code        string     `json:"code"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Version     int64      `json:"version"`
}

type StockItem struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	PartID      string     `json:"part_id"`
	WarehouseID string     `json:"warehouse_id"`
	LocationID  string     `json:"location_id,omitempty"`
	Quantity    float64    `json:"quantity"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Version     int64      `json:"version"`
}

type StockBatch struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	PartID    string     `json:"part_id"`
	BatchCode string     `json:"batch_code"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Version   int64      `json:"version"`
}

type StockLevel struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	PartID      string     `json:"part_id"`
	WarehouseID string     `json:"warehouse_id"`
	Quantity    float64    `json:"quantity"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Version     int64      `json:"version"`
}

type MinimumStock struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	PartID    string     `json:"part_id"`
	Quantity  float64    `json:"quantity"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Version   int64      `json:"version"`
}

type MaximumStock struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	PartID    string     `json:"part_id"`
	Quantity  float64    `json:"quantity"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Version   int64      `json:"version"`
}

type Attachment struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	PartID    string    `json:"part_id"`
	FileName  string    `json:"file_name"`
	MimeType  string    `json:"mime_type,omitempty"`
	URI       string    `json:"uri,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

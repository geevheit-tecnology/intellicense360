package domain

import "time"

type SupplierStatus string
type SupplierTypeCode string

const (
	StatusDraft    SupplierStatus = "draft"
	StatusActive   SupplierStatus = "active"
	StatusInactive SupplierStatus = "inactive"
	StatusBlocked  SupplierStatus = "blocked"
	StatusArchived SupplierStatus = "archived"

	TypePartsSupplier      SupplierTypeCode = "parts_supplier"
	TypeTireSupplier       SupplierTypeCode = "tire_supplier"
	TypeFuelSupplier       SupplierTypeCode = "fuel_supplier"
	TypeWorkshop           SupplierTypeCode = "workshop"
	TypeServiceProvider    SupplierTypeCode = "service_provider"
	TypeTechnologyProvider SupplierTypeCode = "technology_provider"
	TypeInsuranceProvider  SupplierTypeCode = "insurance_provider"
	TypeTransportProvider  SupplierTypeCode = "transport_provider"
	TypeOther              SupplierTypeCode = "other"
)

type Supplier struct {
	ID                    string            `json:"id"`
	TenantID              string            `json:"tenant_id"`
	AuditID               string            `json:"audit_id,omitempty"`
	LegalName             string            `json:"legal_name"`
	TradeName             string            `json:"trade_name,omitempty"`
	CNPJ                  string            `json:"cnpj,omitempty"`
	StateRegistration     string            `json:"state_registration,omitempty"`
	MunicipalRegistration string            `json:"municipal_registration,omitempty"`
	Email                 string            `json:"email,omitempty"`
	Phone                 string            `json:"phone,omitempty"`
	Website               string            `json:"website,omitempty"`
	Notes                 string            `json:"notes,omitempty"`
	Status                SupplierStatus    `json:"status"`
	CategoryID            string            `json:"category_id,omitempty"`
	Type                  SupplierTypeCode  `json:"type"`
	Metadata              map[string]string `json:"metadata,omitempty"`
	CreatedAt             time.Time         `json:"created_at"`
	UpdatedAt             time.Time         `json:"updated_at"`
	DeletedAt             *time.Time        `json:"deleted_at,omitempty"`
	Version               int64             `json:"version"`
}

type SupplierCategory struct {
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

type SupplierType struct {
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

type SupplierContact struct {
	ID         string     `json:"id"`
	TenantID   string     `json:"tenant_id"`
	SupplierID string     `json:"supplier_id"`
	Name       string     `json:"name"`
	Role       string     `json:"role,omitempty"`
	Email      string     `json:"email,omitempty"`
	Phone      string     `json:"phone,omitempty"`
	Mobile     string     `json:"mobile,omitempty"`
	Primary    bool       `json:"primary"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	Version    int64      `json:"version"`
}

type SupplierAddress struct {
	ID           string     `json:"id"`
	TenantID     string     `json:"tenant_id"`
	SupplierID   string     `json:"supplier_id"`
	Street       string     `json:"street"`
	Number       string     `json:"number,omitempty"`
	Complement   string     `json:"complement,omitempty"`
	Neighborhood string     `json:"neighborhood,omitempty"`
	City         string     `json:"city"`
	State        string     `json:"state,omitempty"`
	PostalCode   string     `json:"postal_code,omitempty"`
	Country      string     `json:"country"`
	AddressType  string     `json:"address_type"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	Version      int64      `json:"version"`
}

type SupplierBankAccount struct {
	ID             string     `json:"id"`
	TenantID       string     `json:"tenant_id"`
	SupplierID     string     `json:"supplier_id"`
	Bank           string     `json:"bank"`
	Branch         string     `json:"branch,omitempty"`
	Account        string     `json:"account,omitempty"`
	AccountType    string     `json:"account_type,omitempty"`
	PIXKey         string     `json:"pix_key,omitempty"`
	HolderName     string     `json:"holder_name,omitempty"`
	HolderDocument string     `json:"holder_document,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
	Version        int64      `json:"version"`
}

type SupplierDocument struct {
	ID                  string     `json:"id"`
	TenantID            string     `json:"tenant_id"`
	SupplierID          string     `json:"supplier_id"`
	DocumentType        string     `json:"document_type"`
	DocumentNumber      string     `json:"document_number,omitempty"`
	IssueDate           *time.Time `json:"issue_date,omitempty"`
	ExpirationDate      *time.Time `json:"expiration_date,omitempty"`
	Status              string     `json:"status"`
	AttachmentReference string     `json:"attachment_reference,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
	Version             int64      `json:"version"`
}

type SupplierRating struct {
	ID           string     `json:"id"`
	TenantID     string     `json:"tenant_id"`
	SupplierID   string     `json:"supplier_id"`
	Quality      float64    `json:"quality"`
	Price        float64    `json:"price"`
	Delivery     float64    `json:"delivery"`
	Service      float64    `json:"service"`
	Reliability  float64    `json:"reliability"`
	OverallScore float64    `json:"overall_score"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	Version      int64      `json:"version"`
}

type SupplierContract struct {
	ID                  string     `json:"id"`
	TenantID            string     `json:"tenant_id"`
	SupplierID          string     `json:"supplier_id"`
	ContractNumber      string     `json:"contract_number"`
	StartDate           *time.Time `json:"start_date,omitempty"`
	EndDate             *time.Time `json:"end_date,omitempty"`
	Status              string     `json:"status"`
	Notes               string     `json:"notes,omitempty"`
	AttachmentReference string     `json:"attachment_reference,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty"`
	Version             int64      `json:"version"`
}

type SupplierRepresentative struct {
	ID         string     `json:"id"`
	TenantID   string     `json:"tenant_id"`
	SupplierID string     `json:"supplier_id"`
	Name       string     `json:"name"`
	Document   string     `json:"document,omitempty"`
	Email      string     `json:"email,omitempty"`
	Phone      string     `json:"phone,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	Version    int64      `json:"version"`
}

package domain

import (
	"strings"
	"time"
)

type VehicleCategory string
type VehicleType string
type VehicleStatus string
type OwnershipType string
type FuelType string
type Transmission string
type EmissionStandard string

const (
	VehicleCategoryTruck       VehicleCategory = "truck"
	VehicleCategoryTrailer     VehicleCategory = "trailer"
	VehicleCategorySemiTrailer VehicleCategory = "semi_trailer"
	VehicleCategoryVan         VehicleCategory = "van"
	VehicleCategoryCar         VehicleCategory = "car"
	VehicleCategoryMotorcycle  VehicleCategory = "motorcycle"
	VehicleCategoryMachinery   VehicleCategory = "machinery"
	VehicleCategoryEquipment   VehicleCategory = "equipment"

	VehicleStatusDraft       VehicleStatus = "draft"
	VehicleStatusActive      VehicleStatus = "active"
	VehicleStatusInactive    VehicleStatus = "inactive"
	VehicleStatusMaintenance VehicleStatus = "maintenance"
	VehicleStatusSold        VehicleStatus = "sold"
	VehicleStatusDisposed    VehicleStatus = "disposed"

	OwnershipOwned  OwnershipType = "owned"
	OwnershipLeased OwnershipType = "leased"
	OwnershipRented OwnershipType = "rented"
	OwnershipThird  OwnershipType = "third_party"
)

type LicensePlate string
type Renavam string
type Chassis string
type Engine string
type Color string
type AxleConfiguration string

func (p LicensePlate) Normalized() string {
	return strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(string(p)), "-", ""))
}

type Vehicle struct {
	ID                string            `json:"id"`
	TenantID          string            `json:"tenant_id"`
	AuditID           string            `json:"audit_id,omitempty"`
	CategoryID        string            `json:"category_id"`
	TypeID            string            `json:"type_id"`
	BrandID           string            `json:"brand_id"`
	ModelID           string            `json:"model_id"`
	AssetID           string            `json:"asset_id,omitempty"`
	LicensePlate      LicensePlate      `json:"license_plate"`
	Renavam           Renavam           `json:"renavam"`
	Chassis           Chassis           `json:"chassis"`
	Engine            Engine            `json:"engine,omitempty"`
	Color             Color             `json:"color,omitempty"`
	FuelType          FuelType          `json:"fuel_type,omitempty"`
	Transmission      Transmission      `json:"transmission,omitempty"`
	AxleConfiguration AxleConfiguration `json:"axle_configuration,omitempty"`
	EmissionStandard  EmissionStandard  `json:"emission_standard,omitempty"`
	OwnershipType     OwnershipType     `json:"ownership_type"`
	Status            VehicleStatus     `json:"status"`
	YearManufacture   int               `json:"year_manufacture,omitempty"`
	YearModel         int               `json:"year_model,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
	DeletedAt         *time.Time        `json:"deleted_at,omitempty"`
	Version           int64             `json:"version"`
}

type VehicleBrand struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Version   int64      `json:"version"`
}

type VehicleModel struct {
	ID         string     `json:"id"`
	TenantID   string     `json:"tenant_id"`
	BrandID    string     `json:"brand_id"`
	Name       string     `json:"name"`
	CategoryID string     `json:"category_id,omitempty"`
	TypeID     string     `json:"type_id,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	Version    int64      `json:"version"`
}

type VehicleCategoryEntity struct {
	ID        string          `json:"id"`
	TenantID  string          `json:"tenant_id"`
	Name      string          `json:"name"`
	Code      VehicleCategory `json:"code"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *time.Time      `json:"deleted_at,omitempty"`
	Version   int64           `json:"version"`
}

type VehicleTypeEntity struct {
	ID         string     `json:"id"`
	TenantID   string     `json:"tenant_id"`
	CategoryID string     `json:"category_id"`
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	Version    int64      `json:"version"`
}

type Asset struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Version     int64      `json:"version"`
}

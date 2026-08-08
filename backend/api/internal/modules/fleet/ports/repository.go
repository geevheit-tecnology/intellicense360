package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/domain"
)

type Query struct {
	Search    string
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
	Filters   map[string]string
}

type Page[T any] struct {
	Items      []T `json:"items"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type VehicleRepository interface {
	Create(ctx context.Context, vehicle domain.Vehicle) (domain.Vehicle, error)
	FindByID(ctx context.Context, tenantID string, vehicleID string) (domain.Vehicle, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.Vehicle], error)
	Update(ctx context.Context, vehicle domain.Vehicle) (domain.Vehicle, error)
	Delete(ctx context.Context, tenantID string, vehicleID string) error
	Exists(ctx context.Context, tenantID string, vehicleID string) (bool, error)
	ExistsLicensePlate(ctx context.Context, tenantID string, plate domain.LicensePlate, exceptVehicleID string) (bool, error)
	ExistsChassis(ctx context.Context, tenantID string, chassis domain.Chassis, exceptVehicleID string) (bool, error)
	ExistsRenavam(ctx context.Context, tenantID string, renavam domain.Renavam, exceptVehicleID string) (bool, error)
}

type VehicleBrandRepository interface {
	Create(ctx context.Context, brand domain.VehicleBrand) (domain.VehicleBrand, error)
	FindByID(ctx context.Context, tenantID string, id string) (domain.VehicleBrand, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.VehicleBrand], error)
	Update(ctx context.Context, brand domain.VehicleBrand) (domain.VehicleBrand, error)
	Delete(ctx context.Context, tenantID string, id string) error
}

type VehicleModelRepository interface {
	Create(ctx context.Context, model domain.VehicleModel) (domain.VehicleModel, error)
	FindByID(ctx context.Context, tenantID string, id string) (domain.VehicleModel, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.VehicleModel], error)
	Update(ctx context.Context, model domain.VehicleModel) (domain.VehicleModel, error)
	Delete(ctx context.Context, tenantID string, id string) error
}

type VehicleCategoryRepository interface {
	Create(ctx context.Context, category domain.VehicleCategoryEntity) (domain.VehicleCategoryEntity, error)
	FindByID(ctx context.Context, tenantID string, id string) (domain.VehicleCategoryEntity, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.VehicleCategoryEntity], error)
	Update(ctx context.Context, category domain.VehicleCategoryEntity) (domain.VehicleCategoryEntity, error)
	Delete(ctx context.Context, tenantID string, id string) error
}

type VehicleTypeRepository interface {
	Create(ctx context.Context, vehicleType domain.VehicleTypeEntity) (domain.VehicleTypeEntity, error)
	FindByID(ctx context.Context, tenantID string, id string) (domain.VehicleTypeEntity, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.VehicleTypeEntity], error)
	Update(ctx context.Context, vehicleType domain.VehicleTypeEntity) (domain.VehicleTypeEntity, error)
	Delete(ctx context.Context, tenantID string, id string) error
}

type AssetRepository interface {
	Create(ctx context.Context, asset domain.Asset) (domain.Asset, error)
	FindByID(ctx context.Context, tenantID string, id string) (domain.Asset, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.Asset], error)
	Update(ctx context.Context, asset domain.Asset) (domain.Asset, error)
	Delete(ctx context.Context, tenantID string, id string) error
}

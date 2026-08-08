package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/domain"
)

type VehicleService interface {
	Create(ctx context.Context, vehicle domain.Vehicle) (domain.Vehicle, error)
	FindByID(ctx context.Context, tenantID string, vehicleID string) (domain.Vehicle, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.Vehicle], error)
	Update(ctx context.Context, vehicle domain.Vehicle) (domain.Vehicle, error)
	Delete(ctx context.Context, tenantID string, vehicleID string) error
	ValidateVehicleAccess(ctx context.Context, tenantID string, vehicleID string) error
}

type VehicleBrandService interface {
	Create(ctx context.Context, brand domain.VehicleBrand) (domain.VehicleBrand, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.VehicleBrand], error)
	Update(ctx context.Context, brand domain.VehicleBrand) (domain.VehicleBrand, error)
	Delete(ctx context.Context, tenantID string, id string) error
}

type VehicleModelService interface {
	Create(ctx context.Context, model domain.VehicleModel) (domain.VehicleModel, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.VehicleModel], error)
	Update(ctx context.Context, model domain.VehicleModel) (domain.VehicleModel, error)
	Delete(ctx context.Context, tenantID string, id string) error
}

type VehicleCategoryService interface {
	Create(ctx context.Context, category domain.VehicleCategoryEntity) (domain.VehicleCategoryEntity, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.VehicleCategoryEntity], error)
	Update(ctx context.Context, category domain.VehicleCategoryEntity) (domain.VehicleCategoryEntity, error)
	Delete(ctx context.Context, tenantID string, id string) error
}

type VehicleTypeService interface {
	Create(ctx context.Context, vehicleType domain.VehicleTypeEntity) (domain.VehicleTypeEntity, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.VehicleTypeEntity], error)
	Update(ctx context.Context, vehicleType domain.VehicleTypeEntity) (domain.VehicleTypeEntity, error)
	Delete(ctx context.Context, tenantID string, id string) error
}

type AssetService interface {
	Create(ctx context.Context, asset domain.Asset) (domain.Asset, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.Asset], error)
	Update(ctx context.Context, asset domain.Asset) (domain.Asset, error)
	Delete(ctx context.Context, tenantID string, id string) error
}

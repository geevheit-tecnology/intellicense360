package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/assets/domain"
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

type AssetRepository interface {
	Create(ctx context.Context, asset domain.Asset) (domain.Asset, error)
	FindByID(ctx context.Context, tenantID string, assetID string) (domain.Asset, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.Asset], error)
	Update(ctx context.Context, asset domain.Asset) (domain.Asset, error)
	Delete(ctx context.Context, tenantID string, assetID string) error
	ExistsInternalCode(ctx context.Context, tenantID string, internalCode string, exceptAssetID string) (bool, error)
	ExistsSerialNumber(ctx context.Context, tenantID string, serialNumber string, exceptAssetID string) (bool, error)
	ExistsAssetTag(ctx context.Context, tenantID string, assetTag string, exceptAssetID string) (bool, error)
}

type CategoryRepository interface {
	Create(ctx context.Context, category domain.AssetCategory) (domain.AssetCategory, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.AssetCategory], error)
	Update(ctx context.Context, category domain.AssetCategory) (domain.AssetCategory, error)
	Delete(ctx context.Context, tenantID string, id string) error
}

type TypeRepository interface {
	Create(ctx context.Context, assetType domain.AssetType) (domain.AssetType, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.AssetType], error)
	Update(ctx context.Context, assetType domain.AssetType) (domain.AssetType, error)
	Delete(ctx context.Context, tenantID string, id string) error
}

type ManufacturerRepository interface {
	Create(ctx context.Context, manufacturer domain.Manufacturer) (domain.Manufacturer, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.Manufacturer], error)
	Update(ctx context.Context, manufacturer domain.Manufacturer) (domain.Manufacturer, error)
	Delete(ctx context.Context, tenantID string, id string) error
}

type ModelRepository interface {
	Create(ctx context.Context, model domain.Model) (domain.Model, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.Model], error)
	Update(ctx context.Context, model domain.Model) (domain.Model, error)
	Delete(ctx context.Context, tenantID string, id string) error
}

type EquipmentRepository interface {
	Create(ctx context.Context, equipment domain.Equipment) (domain.Equipment, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.Equipment], error)
	Update(ctx context.Context, equipment domain.Equipment) (domain.Equipment, error)
	Delete(ctx context.Context, tenantID string, id string) error
}

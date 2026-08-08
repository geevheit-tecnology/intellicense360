package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/assets/domain"
)

type AssetService interface {
	Create(ctx context.Context, asset domain.Asset) (domain.Asset, error)
	FindByID(ctx context.Context, tenantID string, assetID string) (domain.Asset, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.Asset], error)
	Update(ctx context.Context, asset domain.Asset) (domain.Asset, error)
	Delete(ctx context.Context, tenantID string, assetID string) error
}

type CategoryService interface {
	Create(ctx context.Context, category domain.AssetCategory) (domain.AssetCategory, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.AssetCategory], error)
	Update(ctx context.Context, category domain.AssetCategory) (domain.AssetCategory, error)
	Delete(ctx context.Context, tenantID string, id string) error
}

type TypeService interface {
	Create(ctx context.Context, assetType domain.AssetType) (domain.AssetType, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.AssetType], error)
	Update(ctx context.Context, assetType domain.AssetType) (domain.AssetType, error)
	Delete(ctx context.Context, tenantID string, id string) error
}

type ManufacturerService interface {
	Create(ctx context.Context, manufacturer domain.Manufacturer) (domain.Manufacturer, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.Manufacturer], error)
	Update(ctx context.Context, manufacturer domain.Manufacturer) (domain.Manufacturer, error)
	Delete(ctx context.Context, tenantID string, id string) error
}

type ModelService interface {
	Create(ctx context.Context, model domain.Model) (domain.Model, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.Model], error)
	Update(ctx context.Context, model domain.Model) (domain.Model, error)
	Delete(ctx context.Context, tenantID string, id string) error
}

type EquipmentService interface {
	Create(ctx context.Context, equipment domain.Equipment) (domain.Equipment, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.Equipment], error)
	Update(ctx context.Context, equipment domain.Equipment) (domain.Equipment, error)
	Delete(ctx context.Context, tenantID string, id string) error
}

type AuditRecorder interface {
	RecordAssetEvent(ctx context.Context, tenantID string, actorID string, action string, resourceID string) error
}

type ImportPort interface {
	PrepareImport(ctx context.Context, tenantID string, fileName string) error
}

type ExportPort interface {
	PrepareExport(ctx context.Context, tenantID string, format string) error
}

type AttachmentPort interface {
	PrepareAttachment(ctx context.Context, attachment domain.Attachment) error
}

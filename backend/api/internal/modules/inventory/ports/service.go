package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/domain"
)

type PartService interface {
	Create(context.Context, domain.Part) (domain.Part, error)
	FindByID(context.Context, string, string) (domain.Part, error)
	Search(context.Context, string, Query) (Page[domain.Part], error)
	Update(context.Context, domain.Part) (domain.Part, error)
	Delete(context.Context, string, string) error
	ValidateInventoryItemAccess(context.Context, string, string) error
}

type CatalogService interface {
	Create(context.Context, domain.Catalog) (domain.Catalog, error)
	Search(context.Context, string, Query) (Page[domain.Catalog], error)
	Update(context.Context, domain.Catalog) (domain.Catalog, error)
	Delete(context.Context, string, string) error
}

type WarehouseService interface {
	Create(context.Context, domain.Warehouse) (domain.Warehouse, error)
	Search(context.Context, string, Query) (Page[domain.Warehouse], error)
	Update(context.Context, domain.Warehouse) (domain.Warehouse, error)
	Delete(context.Context, string, string) error
}

type LocationService interface {
	Create(context.Context, domain.WarehouseLocation) (domain.WarehouseLocation, error)
	Search(context.Context, string, Query) (Page[domain.WarehouseLocation], error)
	Update(context.Context, domain.WarehouseLocation) (domain.WarehouseLocation, error)
	Delete(context.Context, string, string) error
}

type AttachmentPort interface {
	PrepareAttachment(context.Context, domain.Attachment) error
}

package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/domain"
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

type PartRepository interface {
	Create(context.Context, domain.Part) (domain.Part, error)
	FindByID(context.Context, string, string) (domain.Part, error)
	Search(context.Context, string, Query) (Page[domain.Part], error)
	Update(context.Context, domain.Part) (domain.Part, error)
	Delete(context.Context, string, string) error
	ExistsSKU(context.Context, string, string, string) (bool, error)
	ExistsInternalCode(context.Context, string, string, string) (bool, error)
}

type CatalogRepository interface {
	Create(context.Context, domain.Catalog) (domain.Catalog, error)
	Search(context.Context, string, Query) (Page[domain.Catalog], error)
	Update(context.Context, domain.Catalog) (domain.Catalog, error)
	Delete(context.Context, string, string) error
	ExistsCode(context.Context, string, string, string) (bool, error)
}

type WarehouseRepository interface {
	Create(context.Context, domain.Warehouse) (domain.Warehouse, error)
	Search(context.Context, string, Query) (Page[domain.Warehouse], error)
	Update(context.Context, domain.Warehouse) (domain.Warehouse, error)
	Delete(context.Context, string, string) error
}

type LocationRepository interface {
	Create(context.Context, domain.WarehouseLocation) (domain.WarehouseLocation, error)
	Search(context.Context, string, Query) (Page[domain.WarehouseLocation], error)
	Update(context.Context, domain.WarehouseLocation) (domain.WarehouseLocation, error)
	Delete(context.Context, string, string) error
}

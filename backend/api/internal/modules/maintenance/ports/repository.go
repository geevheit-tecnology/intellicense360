package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/domain"
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

type WorkOrderRepository interface {
	Create(context.Context, domain.WorkOrder) (domain.WorkOrder, error)
	FindByID(context.Context, string, string) (domain.WorkOrder, error)
	Search(context.Context, string, Query) (Page[domain.WorkOrder], error)
	Update(context.Context, domain.WorkOrder) (domain.WorkOrder, error)
	Delete(context.Context, string, string) error
	Exists(context.Context, string, string) (bool, error)
	ExistsCode(context.Context, string, string, string) (bool, error)
}

type PreventivePlanRepository interface {
	Create(context.Context, domain.PreventivePlan) (domain.PreventivePlan, error)
	Search(context.Context, string, Query) (Page[domain.PreventivePlan], error)
	Update(context.Context, domain.PreventivePlan) (domain.PreventivePlan, error)
	Delete(context.Context, string, string) error
}

type ServiceTypeRepository interface {
	Create(context.Context, domain.ServiceType) (domain.ServiceType, error)
	Search(context.Context, string, Query) (Page[domain.ServiceType], error)
	Update(context.Context, domain.ServiceType) (domain.ServiceType, error)
	Delete(context.Context, string, string) error
}

type LaborRepository interface {
	Create(context.Context, domain.LaborEntry) (domain.LaborEntry, error)
	List(context.Context, string, string, Query) (Page[domain.LaborEntry], error)
	Delete(context.Context, string, string, string) error
}

type DowntimeRepository interface {
	Create(context.Context, domain.Downtime) (domain.Downtime, error)
	List(context.Context, string, string, Query) (Page[domain.Downtime], error)
	End(context.Context, string, string, string) (domain.Downtime, error)
}

type HistoryRepository interface {
	Create(context.Context, domain.MaintenanceHistory) (domain.MaintenanceHistory, error)
	List(context.Context, string, string, Query) (Page[domain.MaintenanceHistory], error)
}

type CatalogRepository interface {
	Create(context.Context, domain.MaintenanceCatalog) (domain.MaintenanceCatalog, error)
	Search(context.Context, string, Query) (Page[domain.MaintenanceCatalog], error)
	Update(context.Context, domain.MaintenanceCatalog) (domain.MaintenanceCatalog, error)
	Delete(context.Context, string, string) error
}

type WorkshopRepository interface {
	Create(context.Context, domain.Workshop) (domain.Workshop, error)
	Search(context.Context, string, Query) (Page[domain.Workshop], error)
	Update(context.Context, domain.Workshop) (domain.Workshop, error)
	Delete(context.Context, string, string) error
}

type TechnicianRepository interface {
	Create(context.Context, domain.Technician) (domain.Technician, error)
	Search(context.Context, string, Query) (Page[domain.Technician], error)
	Update(context.Context, domain.Technician) (domain.Technician, error)
	Delete(context.Context, string, string) error
}

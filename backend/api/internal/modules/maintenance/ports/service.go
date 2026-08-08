package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/domain"
)

type WorkOrderService interface {
	Create(context.Context, domain.WorkOrder) (domain.WorkOrder, error)
	FindByID(context.Context, string, string) (domain.WorkOrder, error)
	Search(context.Context, string, Query) (Page[domain.WorkOrder], error)
	Update(context.Context, domain.WorkOrder) (domain.WorkOrder, error)
	Delete(context.Context, string, string) error
	Start(context.Context, string, string, string) (domain.WorkOrder, error)
	Complete(context.Context, string, string, string) (domain.WorkOrder, error)
	Cancel(context.Context, string, string, string) (domain.WorkOrder, error)
	ValidateMaintenanceAccess(context.Context, string, string) error
}

type PreventivePlanService interface {
	Create(context.Context, domain.PreventivePlan) (domain.PreventivePlan, error)
	Search(context.Context, string, Query) (Page[domain.PreventivePlan], error)
	Update(context.Context, domain.PreventivePlan) (domain.PreventivePlan, error)
	Delete(context.Context, string, string) error
}

type ServiceTypeService interface {
	Create(context.Context, domain.ServiceType) (domain.ServiceType, error)
	Search(context.Context, string, Query) (Page[domain.ServiceType], error)
	Update(context.Context, domain.ServiceType) (domain.ServiceType, error)
	Delete(context.Context, string, string) error
}

type LaborService interface {
	Add(context.Context, domain.LaborEntry) (domain.LaborEntry, error)
	List(context.Context, string, string, Query) (Page[domain.LaborEntry], error)
	Delete(context.Context, string, string, string) error
}

type DowntimeService interface {
	Start(context.Context, domain.Downtime) (domain.Downtime, error)
	End(context.Context, string, string, string) (domain.Downtime, error)
	List(context.Context, string, string, Query) (Page[domain.Downtime], error)
}

type HistoryService interface {
	List(context.Context, string, string, Query) (Page[domain.MaintenanceHistory], error)
}

type CatalogService interface {
	Create(context.Context, domain.MaintenanceCatalog) (domain.MaintenanceCatalog, error)
	Search(context.Context, string, Query) (Page[domain.MaintenanceCatalog], error)
	Update(context.Context, domain.MaintenanceCatalog) (domain.MaintenanceCatalog, error)
	Delete(context.Context, string, string) error
}

type WorkshopService interface {
	Create(context.Context, domain.Workshop) (domain.Workshop, error)
	Search(context.Context, string, Query) (Page[domain.Workshop], error)
	Update(context.Context, domain.Workshop) (domain.Workshop, error)
	Delete(context.Context, string, string) error
}

type TechnicianService interface {
	Create(context.Context, domain.Technician) (domain.Technician, error)
	Search(context.Context, string, Query) (Page[domain.Technician], error)
	Update(context.Context, domain.Technician) (domain.Technician, error)
	Delete(context.Context, string, string) error
}

type AttachmentPort interface {
	PrepareAttachment(context.Context, domain.MaintenanceAttachment) error
}

type AuditRecorder interface {
	RecordMaintenanceEvent(context.Context, string, string, string, string) error
}

type ReportPort interface {
	PrepareReport(context.Context, string, string, string) error
}

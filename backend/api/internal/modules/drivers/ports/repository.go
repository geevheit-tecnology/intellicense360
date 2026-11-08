package ports

import (
	"context"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/drivers/domain"
)

type Query struct {
	Search    string
	Status    string
	VehicleID string
	Page      int
	PageSize  int
}
type Page struct {
	Items      []domain.Driver `json:"items"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalItems int             `json:"total_items"`
	TotalPages int             `json:"total_pages"`
}
type DriverRepository interface {
	Create(context.Context, domain.Driver) (domain.Driver, error)
	FindByID(context.Context, string, string) (domain.Driver, error)
	Search(context.Context, string, Query) (Page, error)
	Update(context.Context, domain.Driver) (domain.Driver, error)
	Delete(context.Context, string, string, int64) error
	Exists(context.Context, string, string) (bool, error)
}

// These small local contracts avoid importing other modules' implementations.
type VehicleAccess interface {
	ValidateVehicleAccess(context.Context, string, string) error
}
type Transactions interface {
	Run(context.Context, func(context.Context) error) error
}
type AuditEvent struct {
	TenantID   string
	ActorID    string
	Action     string
	ResourceID string
	Metadata   map[string]string
}
type Audit interface {
	RecordDriverEvent(context.Context, AuditEvent) error
}

type VehicleAccessFunc func(context.Context, string, string) error

func (f VehicleAccessFunc) ValidateVehicleAccess(ctx context.Context, tenant, id string) error {
	return f(ctx, tenant, id)
}

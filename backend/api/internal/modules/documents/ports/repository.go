package ports

import (
	"context"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/documents/domain"
	"time"
)

type Query struct {
	Search       string
	Kind         string
	Status       string
	DriverID     string
	VehicleID    string
	ExpiryStatus string
	Page         int
	PageSize     int
	AsOf         time.Time
}
type Page struct {
	Items      []domain.Document `json:"items"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalItems int               `json:"total_items"`
	TotalPages int               `json:"total_pages"`
}
type DocumentRepository interface {
	Create(context.Context, domain.Document) (domain.Document, error)
	FindByID(context.Context, string, string) (domain.Document, error)
	Search(context.Context, string, Query) (Page, error)
	Update(context.Context, domain.Document) (domain.Document, error)
	Delete(context.Context, string, string, int64) error
	Exists(context.Context, string, string) (bool, error)
}
type SubjectAccess interface {
	ValidateDriverAccess(context.Context, string, string) error
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
	RecordDocumentEvent(context.Context, AuditEvent) error
}

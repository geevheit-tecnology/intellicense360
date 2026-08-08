package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/tires/domain"
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

type TireRepository interface {
	Create(ctx context.Context, tire domain.Tire) (domain.Tire, error)
	FindByID(ctx context.Context, tenantID string, tireID string) (domain.Tire, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.Tire], error)
	Update(ctx context.Context, tire domain.Tire) (domain.Tire, error)
	Delete(ctx context.Context, tenantID string, tireID string) error
	Exists(ctx context.Context, tenantID string, tireID string) (bool, error)
	ExistsSerialNumber(ctx context.Context, tenantID string, serialNumber string, exceptTireID string) (bool, error)
	ExistsFireNumber(ctx context.Context, tenantID string, fireNumber string, exceptTireID string) (bool, error)
}

type TireInspectionRepository interface {
	Create(ctx context.Context, inspection domain.TireInspection) (domain.TireInspection, error)
	FindByID(ctx context.Context, tenantID string, tireID string, inspectionID string) (domain.TireInspection, error)
	List(ctx context.Context, tenantID string, tireID string, query Query) (Page[domain.TireInspection], error)
	Update(ctx context.Context, inspection domain.TireInspection) (domain.TireInspection, error)
	Delete(ctx context.Context, tenantID string, tireID string, inspectionID string) error
}

type TireMovementRepository interface {
	Create(ctx context.Context, movement domain.TireMovement) (domain.TireMovement, error)
	List(ctx context.Context, tenantID string, tireID string, query Query) (Page[domain.TireMovement], error)
}

type TireHistoryRepository interface {
	Create(ctx context.Context, history domain.TireHistory) (domain.TireHistory, error)
	List(ctx context.Context, tenantID string, tireID string, query Query) (Page[domain.TireHistory], error)
}

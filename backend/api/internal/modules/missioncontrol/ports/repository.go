package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/domain"
)

type Query struct {
	Search  string
	Page    int
	PerPage int
	Sort    string
	Filters map[string]string
}

type Page[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type CommandItemRepository interface {
	Create(ctx context.Context, item domain.CommandItem) (domain.CommandItem, error)
	GetByID(ctx context.Context, tenantID, id string) (domain.CommandItem, error)
	Update(ctx context.Context, item domain.CommandItem) (domain.CommandItem, error)
	List(ctx context.Context, tenantID string, q Query) (Page[domain.CommandItem], error)
	FindByFingerprint(ctx context.Context, tenantID, fingerprint string) (domain.CommandItem, error)
	AllOpen(ctx context.Context, tenantID string) ([]domain.CommandItem, error)
}

type CommandEventRepository interface {
	Create(ctx context.Context, event domain.CommandEvent) (domain.CommandEvent, error)
	ListByItem(ctx context.Context, tenantID, itemID string, q Query) (Page[domain.CommandEvent], error)
}

type CommandActionRepository interface {
	Create(ctx context.Context, action domain.CommandAction) (domain.CommandAction, error)
	ListByItem(ctx context.Context, tenantID, itemID string, q Query) (Page[domain.CommandAction], error)
}

type OperationalSnapshotRepository interface {
	Create(ctx context.Context, snapshot domain.OperationalSnapshot) (domain.OperationalSnapshot, error)
	Latest(ctx context.Context, tenantID string) (domain.OperationalSnapshot, error)
}

type IdempotencyRepository interface {
	Exists(ctx context.Context, tenantID, key string) (bool, error)
	Save(ctx context.Context, tenantID, key, resourceID string) error
}

type OperationalSignalProvider interface {
	GetSignals(ctx context.Context, tenantID string) ([]domain.Signal, error)
}

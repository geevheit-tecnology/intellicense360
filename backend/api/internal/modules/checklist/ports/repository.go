package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/domain"
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

type ChecklistRepository interface {
	Create(ctx context.Context, checklist domain.Checklist) (domain.Checklist, error)
	FindByID(ctx context.Context, tenantID string, checklistID string) (domain.Checklist, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.Checklist], error)
	Update(ctx context.Context, checklist domain.Checklist) (domain.Checklist, error)
	Delete(ctx context.Context, tenantID string, checklistID string) error
	Exists(ctx context.Context, tenantID string, checklistID string) (bool, error)
}

type ChecklistItemRepository interface {
	Create(ctx context.Context, item domain.ChecklistItem) (domain.ChecklistItem, error)
	FindByID(ctx context.Context, tenantID string, checklistID string, itemID string) (domain.ChecklistItem, error)
	List(ctx context.Context, tenantID string, checklistID string, query Query) (Page[domain.ChecklistItem], error)
	Update(ctx context.Context, item domain.ChecklistItem) (domain.ChecklistItem, error)
	Delete(ctx context.Context, tenantID string, checklistID string, itemID string) error
}

type ChecklistAnswerRepository interface {
	Create(ctx context.Context, answer domain.ChecklistAnswer) (domain.ChecklistAnswer, error)
	List(ctx context.Context, tenantID string, checklistID string, query Query) (Page[domain.ChecklistAnswer], error)
}

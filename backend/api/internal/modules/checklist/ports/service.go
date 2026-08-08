package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/domain"
)

type ChecklistService interface {
	Create(ctx context.Context, checklist domain.Checklist) (domain.Checklist, error)
	Update(ctx context.Context, checklist domain.Checklist) (domain.Checklist, error)
	Delete(ctx context.Context, tenantID string, checklistID string) error
	FindByID(ctx context.Context, tenantID string, checklistID string) (domain.Checklist, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.Checklist], error)
	Start(ctx context.Context, tenantID string, checklistID string, actorID string) (domain.Checklist, error)
	Finish(ctx context.Context, tenantID string, checklistID string, actorID string) (domain.Checklist, error)
	Cancel(ctx context.Context, tenantID string, checklistID string, actorID string) (domain.Checklist, error)
	ValidateChecklistAccess(ctx context.Context, tenantID string, checklistID string) error
}

type ChecklistItemService interface {
	Add(ctx context.Context, item domain.ChecklistItem) (domain.ChecklistItem, error)
	Update(ctx context.Context, item domain.ChecklistItem) (domain.ChecklistItem, error)
	Delete(ctx context.Context, tenantID string, checklistID string, itemID string) error
	List(ctx context.Context, tenantID string, checklistID string, query Query) (Page[domain.ChecklistItem], error)
}

type ChecklistAnswerService interface {
	AnswerItem(ctx context.Context, answer domain.ChecklistAnswer) (domain.ChecklistAnswer, error)
	List(ctx context.Context, tenantID string, checklistID string, query Query) (Page[domain.ChecklistAnswer], error)
}

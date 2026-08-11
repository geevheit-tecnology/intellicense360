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

type ChecklistTemplateService interface {
	Create(context.Context, domain.ChecklistTemplate) (domain.ChecklistTemplate, error)
	Update(context.Context, domain.ChecklistTemplate) (domain.ChecklistTemplate, error)
	Archive(context.Context, string, string, string) (domain.ChecklistTemplate, error)
	FindByID(context.Context, string, string) (domain.ChecklistTemplate, error)
	Search(context.Context, string, Query) (Page[domain.ChecklistTemplate], error)
}

type ChecklistTemplateVersionService interface {
	Create(context.Context, domain.ChecklistTemplateVersion) (domain.ChecklistTemplateVersion, error)
	Publish(context.Context, string, string, string) (domain.ChecklistTemplateVersion, error)
	ListByTemplate(context.Context, string, string, Query) (Page[domain.ChecklistTemplateVersion], error)
}

type ChecklistExecutionService interface {
	Start(context.Context, domain.ChecklistExecution) (domain.ChecklistExecution, error)
	Complete(context.Context, string, string, string) (domain.ChecklistExecution, error)
	Cancel(context.Context, string, string, string) (domain.ChecklistExecution, error)
	Invalidate(context.Context, string, string, string) (domain.ChecklistExecution, error)
	FindByID(context.Context, string, string) (domain.ChecklistExecution, error)
	Search(context.Context, string, Query) (Page[domain.ChecklistExecution], error)
}

type ChecklistEngineResponseService interface {
	Record(context.Context, domain.ChecklistResponse) (domain.ChecklistResponse, error)
	ListByExecution(context.Context, string, string, Query) (Page[domain.ChecklistResponse], error)
}

type ChecklistEvidenceService interface {
	Add(context.Context, domain.ChecklistEvidence) (domain.ChecklistEvidence, error)
	ListByExecution(context.Context, string, string, Query) (Page[domain.ChecklistEvidence], error)
}

type ChecklistNonConformityService interface {
	Create(context.Context, domain.ChecklistNonConformity) (domain.ChecklistNonConformity, error)
	ListByExecution(context.Context, string, string, Query) (Page[domain.ChecklistNonConformity], error)
}

type ChecklistHistoryService interface {
	ListByExecution(context.Context, string, string, Query) (Page[domain.ChecklistHistory], error)
}

type ChecklistTypeService interface {
	Create(context.Context, domain.ChecklistType) (domain.ChecklistType, error)
	Search(context.Context, string, Query) (Page[domain.ChecklistType], error)
}

type ChecklistSectionService interface {
	Create(context.Context, domain.ChecklistSection) (domain.ChecklistSection, error)
	Search(context.Context, string, Query) (Page[domain.ChecklistSection], error)
}

type ChecklistEngineItemService interface {
	Create(context.Context, domain.ChecklistEngineItem) (domain.ChecklistEngineItem, error)
	ListByVersion(context.Context, string, string, Query) (Page[domain.ChecklistEngineItem], error)
}

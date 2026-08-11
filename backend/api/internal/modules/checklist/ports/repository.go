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

type ChecklistTemplateRepository interface {
	Create(context.Context, domain.ChecklistTemplate) (domain.ChecklistTemplate, error)
	FindByID(context.Context, string, string) (domain.ChecklistTemplate, error)
	Search(context.Context, string, Query) (Page[domain.ChecklistTemplate], error)
	Update(context.Context, domain.ChecklistTemplate) (domain.ChecklistTemplate, error)
	Delete(context.Context, string, string) error
}

type ChecklistTemplateVersionRepository interface {
	Create(context.Context, domain.ChecklistTemplateVersion) (domain.ChecklistTemplateVersion, error)
	FindByID(context.Context, string, string) (domain.ChecklistTemplateVersion, error)
	ListByTemplate(context.Context, string, string, Query) (Page[domain.ChecklistTemplateVersion], error)
	Update(context.Context, domain.ChecklistTemplateVersion) (domain.ChecklistTemplateVersion, error)
}

type ChecklistTypeRepository interface {
	Create(context.Context, domain.ChecklistType) (domain.ChecklistType, error)
	Search(context.Context, string, Query) (Page[domain.ChecklistType], error)
}

type ChecklistSectionRepository interface {
	Create(context.Context, domain.ChecklistSection) (domain.ChecklistSection, error)
	Search(context.Context, string, Query) (Page[domain.ChecklistSection], error)
}

type ChecklistEngineItemRepository interface {
	Create(context.Context, domain.ChecklistEngineItem) (domain.ChecklistEngineItem, error)
	FindByID(context.Context, string, string) (domain.ChecklistEngineItem, error)
	ListByVersion(context.Context, string, string, Query) (Page[domain.ChecklistEngineItem], error)
}

type ChecklistExecutionRepository interface {
	Create(context.Context, domain.ChecklistExecution) (domain.ChecklistExecution, error)
	FindByID(context.Context, string, string) (domain.ChecklistExecution, error)
	Search(context.Context, string, Query) (Page[domain.ChecklistExecution], error)
	Update(context.Context, domain.ChecklistExecution) (domain.ChecklistExecution, error)
}

type ChecklistResponseRepository interface {
	Create(context.Context, domain.ChecklistResponse) (domain.ChecklistResponse, error)
	ListByExecution(context.Context, string, string, Query) (Page[domain.ChecklistResponse], error)
	ExistsForItem(context.Context, string, string, string) (bool, error)
}

type ChecklistEvidenceRepository interface {
	Create(context.Context, domain.ChecklistEvidence) (domain.ChecklistEvidence, error)
	ListByExecution(context.Context, string, string, Query) (Page[domain.ChecklistEvidence], error)
	ExistsForResponse(context.Context, string, string, string) (bool, error)
}

type ChecklistNonConformityRepository interface {
	Create(context.Context, domain.ChecklistNonConformity) (domain.ChecklistNonConformity, error)
	ListByExecution(context.Context, string, string, Query) (Page[domain.ChecklistNonConformity], error)
}

type ChecklistSignatureRepository interface {
	Create(context.Context, domain.ChecklistSignature) (domain.ChecklistSignature, error)
	ExistsForExecution(context.Context, string, string) (bool, error)
}

type ChecklistHistoryRepository interface {
	Create(context.Context, domain.ChecklistHistory) (domain.ChecklistHistory, error)
	ListByExecution(context.Context, string, string, Query) (Page[domain.ChecklistHistory], error)
}

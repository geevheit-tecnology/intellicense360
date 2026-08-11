package application

import (
	"context"
	"strings"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/ports"
)

type ChecklistTemplateService struct {
	repo ports.ChecklistTemplateRepository
}

func NewChecklistTemplateService(repo ports.ChecklistTemplateRepository) ChecklistTemplateService {
	return ChecklistTemplateService{repo: repo}
}
func (s ChecklistTemplateService) Create(ctx context.Context, item domain.ChecklistTemplate) (domain.ChecklistTemplate, error) {
	if item.TenantID == "" || strings.TrimSpace(item.Name) == "" {
		return domain.ChecklistTemplate{}, ErrValidation
	}
	now := time.Now().UTC()
	item.ID = newID("ctpl")
	item.Status = domain.TemplateStatusDraft
	item.Active = true
	item.CreatedAt, item.UpdatedAt, item.Version = now, now, 1
	return s.repo.Create(ctx, item)
}
func (s ChecklistTemplateService) Update(ctx context.Context, item domain.ChecklistTemplate) (domain.ChecklistTemplate, error) {
	current, err := s.repo.FindByID(ctx, item.TenantID, item.ID)
	if err != nil {
		return domain.ChecklistTemplate{}, err
	}
	if current.Status == domain.TemplateStatusArchived {
		return domain.ChecklistTemplate{}, ErrInvalidTransition
	}
	if strings.TrimSpace(item.Name) == "" {
		return domain.ChecklistTemplate{}, ErrValidation
	}
	item.Status = current.Status
	item.Active = current.Active
	item.CurrentVersionNumber = current.CurrentVersionNumber
	item.CreatedAt = current.CreatedAt
	item.UpdatedAt = time.Now().UTC()
	item.Version = current.Version + 1
	return s.repo.Update(ctx, item)
}
func (s ChecklistTemplateService) Archive(ctx context.Context, tenantID string, id string, actorID string) (domain.ChecklistTemplate, error) {
	item, err := s.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return domain.ChecklistTemplate{}, err
	}
	item.Status = domain.TemplateStatusArchived
	item.Active = false
	item.UpdatedBy = actorID
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}
func (s ChecklistTemplateService) FindByID(ctx context.Context, tenantID string, id string) (domain.ChecklistTemplate, error) {
	return s.repo.FindByID(ctx, tenantID, id)
}
func (s ChecklistTemplateService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.ChecklistTemplate], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(q))
}

type ChecklistTemplateVersionService struct {
	templates ports.ChecklistTemplateRepository
	versions  ports.ChecklistTemplateVersionRepository
}

func NewChecklistTemplateVersionService(templates ports.ChecklistTemplateRepository, versions ports.ChecklistTemplateVersionRepository) ChecklistTemplateVersionService {
	return ChecklistTemplateVersionService{templates: templates, versions: versions}
}
func (s ChecklistTemplateVersionService) Create(ctx context.Context, item domain.ChecklistTemplateVersion) (domain.ChecklistTemplateVersion, error) {
	template, err := s.templates.FindByID(ctx, item.TenantID, item.TemplateID)
	if err != nil {
		return domain.ChecklistTemplateVersion{}, err
	}
	if template.Status == domain.TemplateStatusArchived {
		return domain.ChecklistTemplateVersion{}, ErrInvalidTransition
	}
	now := time.Now().UTC()
	item.ID = newID("ctv")
	item.Status = domain.TemplateStatusDraft
	item.VersionNumber = template.CurrentVersionNumber + 1
	item.CreatedAt, item.UpdatedAt, item.Version = now, now, 1
	saved, err := s.versions.Create(ctx, item)
	if err != nil {
		return domain.ChecklistTemplateVersion{}, err
	}
	template.CurrentVersionNumber = saved.VersionNumber
	template.UpdatedAt = now
	template.Version++
	_, _ = s.templates.Update(ctx, template)
	return saved, nil
}
func (s ChecklistTemplateVersionService) Publish(ctx context.Context, tenantID string, id string, actorID string) (domain.ChecklistTemplateVersion, error) {
	item, err := s.versions.FindByID(ctx, tenantID, id)
	if err != nil {
		return domain.ChecklistTemplateVersion{}, err
	}
	if item.Status != domain.TemplateStatusDraft {
		return domain.ChecklistTemplateVersion{}, ErrPublishedVersionImmutable
	}
	now := time.Now().UTC()
	item.Status = domain.TemplateStatusPublished
	item.PublishedAt = &now
	item.UpdatedAt = now
	item.Version++
	saved, err := s.versions.Update(ctx, item)
	if err != nil {
		return domain.ChecklistTemplateVersion{}, err
	}
	template, err := s.templates.FindByID(ctx, tenantID, item.TemplateID)
	if err == nil {
		template.Status = domain.TemplateStatusPublished
		template.UpdatedBy = actorID
		template.UpdatedAt = now
		template.Version++
		_, _ = s.templates.Update(ctx, template)
	}
	return saved, nil
}
func (s ChecklistTemplateVersionService) ListByTemplate(ctx context.Context, tenantID string, templateID string, q ports.Query) (ports.Page[domain.ChecklistTemplateVersion], error) {
	return s.versions.ListByTemplate(ctx, tenantID, templateID, normalizeQuery(q))
}

type ChecklistExecutionService struct {
	versions   ports.ChecklistTemplateVersionRepository
	items      ports.ChecklistEngineItemRepository
	executions ports.ChecklistExecutionRepository
	responses  ports.ChecklistResponseRepository
	evidence   ports.ChecklistEvidenceRepository
	signatures ports.ChecklistSignatureRepository
	history    ports.ChecklistHistoryRepository
}

func NewChecklistExecutionService(versions ports.ChecklistTemplateVersionRepository, items ports.ChecklistEngineItemRepository, executions ports.ChecklistExecutionRepository, responses ports.ChecklistResponseRepository, evidence ports.ChecklistEvidenceRepository, signatures ports.ChecklistSignatureRepository, history ports.ChecklistHistoryRepository) ChecklistExecutionService {
	return ChecklistExecutionService{versions: versions, items: items, executions: executions, responses: responses, evidence: evidence, signatures: signatures, history: history}
}
func (s ChecklistExecutionService) Start(ctx context.Context, item domain.ChecklistExecution) (domain.ChecklistExecution, error) {
	version, err := s.versions.FindByID(ctx, item.TenantID, item.TemplateVersionID)
	if err != nil {
		return domain.ChecklistExecution{}, err
	}
	if version.Status != domain.TemplateStatusPublished {
		return domain.ChecklistExecution{}, ErrInvalidTransition
	}
	now := time.Now().UTC()
	item.ID = newID("cex")
	item.Status = domain.ExecutionStatusInProgress
	item.StartedAt = &now
	item.CreatedAt, item.UpdatedAt, item.Version = now, now, 1
	saved, err := s.executions.Create(ctx, item)
	if err == nil {
		_, _ = s.history.Create(ctx, history(item.TenantID, saved.ID, "started", item.PerformedBy, ""))
	}
	return saved, err
}
func (s ChecklistExecutionService) Complete(ctx context.Context, tenantID string, id string, actorID string) (domain.ChecklistExecution, error) {
	execution, err := s.executions.FindByID(ctx, tenantID, id)
	if err != nil {
		return domain.ChecklistExecution{}, err
	}
	if execution.Status != domain.ExecutionStatusInProgress {
		return domain.ChecklistExecution{}, ErrInvalidTransition
	}
	version, err := s.versions.FindByID(ctx, tenantID, execution.TemplateVersionID)
	if err != nil {
		return domain.ChecklistExecution{}, err
	}
	items, err := s.items.ListByVersion(ctx, tenantID, execution.TemplateVersionID, ports.Query{PageSize: 100})
	if err != nil {
		return domain.ChecklistExecution{}, err
	}
	for _, item := range items.Items {
		answered, err := s.responses.ExistsForItem(ctx, tenantID, execution.ID, item.ID)
		if err != nil {
			return domain.ChecklistExecution{}, err
		}
		if item.Required && !answered {
			return domain.ChecklistExecution{}, ErrRequiredItemsUnanswered
		}
		if item.EvidenceRequired && answered {
			hasEvidence, err := s.evidence.ExistsForResponse(ctx, tenantID, execution.ID, item.ID)
			if err != nil {
				return domain.ChecklistExecution{}, err
			}
			if !hasEvidence {
				return domain.ChecklistExecution{}, ErrRequiredEvidenceMissing
			}
		}
	}
	if version.SignatureRequired {
		ok, err := s.signatures.ExistsForExecution(ctx, tenantID, execution.ID)
		if err != nil {
			return domain.ChecklistExecution{}, err
		}
		if !ok {
			return domain.ChecklistExecution{}, ErrRequiredSignatureMissing
		}
	}
	now := time.Now().UTC()
	execution.Status = domain.ExecutionStatusCompleted
	execution.EndedAt = &now
	execution.UpdatedAt = now
	execution.Version++
	saved, err := s.executions.Update(ctx, execution)
	if err == nil {
		_, _ = s.history.Create(ctx, history(tenantID, id, "completed", actorID, ""))
	}
	return saved, err
}
func (s ChecklistExecutionService) Cancel(ctx context.Context, tenantID string, id string, actorID string) (domain.ChecklistExecution, error) {
	return s.transition(ctx, tenantID, id, actorID, domain.ExecutionStatusCanceled, "canceled")
}
func (s ChecklistExecutionService) Invalidate(ctx context.Context, tenantID string, id string, actorID string) (domain.ChecklistExecution, error) {
	return s.transition(ctx, tenantID, id, actorID, domain.ExecutionStatusInvalidated, "invalidated")
}
func (s ChecklistExecutionService) FindByID(ctx context.Context, tenantID string, id string) (domain.ChecklistExecution, error) {
	return s.executions.FindByID(ctx, tenantID, id)
}
func (s ChecklistExecutionService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.ChecklistExecution], error) {
	return s.executions.Search(ctx, tenantID, normalizeQuery(q))
}
func (s ChecklistExecutionService) transition(ctx context.Context, tenantID string, id string, actorID string, status domain.ExecutionStatus, event string) (domain.ChecklistExecution, error) {
	execution, err := s.executions.FindByID(ctx, tenantID, id)
	if err != nil {
		return domain.ChecklistExecution{}, err
	}
	if status == domain.ExecutionStatusCanceled && execution.Status != domain.ExecutionStatusInProgress {
		return domain.ChecklistExecution{}, ErrInvalidTransition
	}
	if status == domain.ExecutionStatusInvalidated && execution.Status != domain.ExecutionStatusCompleted {
		return domain.ChecklistExecution{}, ErrInvalidTransition
	}
	now := time.Now().UTC()
	execution.Status = status
	execution.UpdatedAt = now
	execution.Version++
	saved, err := s.executions.Update(ctx, execution)
	if err == nil {
		_, _ = s.history.Create(ctx, history(tenantID, id, event, actorID, ""))
	}
	return saved, err
}

type ChecklistEngineResponseService struct {
	executions ports.ChecklistExecutionRepository
	items      ports.ChecklistEngineItemRepository
	responses  ports.ChecklistResponseRepository
	history    ports.ChecklistHistoryRepository
}

func NewChecklistEngineResponseService(executions ports.ChecklistExecutionRepository, items ports.ChecklistEngineItemRepository, responses ports.ChecklistResponseRepository, history ports.ChecklistHistoryRepository) ChecklistEngineResponseService {
	return ChecklistEngineResponseService{executions: executions, items: items, responses: responses, history: history}
}
func (s ChecklistEngineResponseService) Record(ctx context.Context, item domain.ChecklistResponse) (domain.ChecklistResponse, error) {
	execution, err := s.executions.FindByID(ctx, item.TenantID, item.ExecutionID)
	if err != nil {
		return domain.ChecklistResponse{}, err
	}
	if execution.Status != domain.ExecutionStatusInProgress {
		return domain.ChecklistResponse{}, ErrInvalidTransition
	}
	if _, err := s.items.FindByID(ctx, item.TenantID, item.ItemID); err != nil {
		return domain.ChecklistResponse{}, err
	}
	if item.Value == "" {
		return domain.ChecklistResponse{}, ErrInvalidAnswer
	}
	now := time.Now().UTC()
	item.ID = newID("crs")
	item.RespondedAt, item.CreatedAt = now, now
	saved, err := s.responses.Create(ctx, item)
	if err == nil {
		_, _ = s.history.Create(ctx, history(item.TenantID, item.ExecutionID, "response_recorded", item.Responder, item.ItemID))
	}
	return saved, err
}
func (s ChecklistEngineResponseService) ListByExecution(ctx context.Context, tenantID string, executionID string, q ports.Query) (ports.Page[domain.ChecklistResponse], error) {
	return s.responses.ListByExecution(ctx, tenantID, executionID, normalizeQuery(q))
}

type ChecklistEvidenceService struct {
	repo    ports.ChecklistEvidenceRepository
	history ports.ChecklistHistoryRepository
}
type ChecklistNonConformityService struct {
	repo    ports.ChecklistNonConformityRepository
	history ports.ChecklistHistoryRepository
}
type ChecklistHistoryService struct {
	repo ports.ChecklistHistoryRepository
}

func NewChecklistEvidenceService(repo ports.ChecklistEvidenceRepository, history ports.ChecklistHistoryRepository) ChecklistEvidenceService {
	return ChecklistEvidenceService{repo: repo, history: history}
}
func NewChecklistNonConformityService(repo ports.ChecklistNonConformityRepository, history ports.ChecklistHistoryRepository) ChecklistNonConformityService {
	return ChecklistNonConformityService{repo: repo, history: history}
}
func NewChecklistHistoryService(repo ports.ChecklistHistoryRepository) ChecklistHistoryService {
	return ChecklistHistoryService{repo: repo}
}

func (s ChecklistEvidenceService) Add(ctx context.Context, item domain.ChecklistEvidence) (domain.ChecklistEvidence, error) {
	if item.TenantID == "" || item.ExecutionID == "" || item.EvidenceType == "" || item.Reference == "" {
		return domain.ChecklistEvidence{}, ErrValidation
	}
	item.ID, item.CreatedAt = newID("cev"), time.Now().UTC()
	saved, err := s.repo.Create(ctx, item)
	if err == nil {
		_, _ = s.history.Create(ctx, history(item.TenantID, item.ExecutionID, "evidence_added", "", item.ResponseID))
	}
	return saved, err
}
func (s ChecklistEvidenceService) ListByExecution(ctx context.Context, tenantID string, executionID string, q ports.Query) (ports.Page[domain.ChecklistEvidence], error) {
	return s.repo.ListByExecution(ctx, tenantID, executionID, normalizeQuery(q))
}
func (s ChecklistNonConformityService) Create(ctx context.Context, item domain.ChecklistNonConformity) (domain.ChecklistNonConformity, error) {
	if item.TenantID == "" || item.ExecutionID == "" || strings.TrimSpace(item.Title) == "" {
		return domain.ChecklistNonConformity{}, ErrValidation
	}
	if item.Status == "" {
		item.Status = "open"
	}
	if item.Severity == "" {
		item.Severity = domain.SeverityMedium
	}
	item.ID, item.CreatedAt = newID("cnc"), time.Now().UTC()
	saved, err := s.repo.Create(ctx, item)
	if err == nil {
		_, _ = s.history.Create(ctx, history(item.TenantID, item.ExecutionID, "non_conformity_created", "", item.Title))
	}
	return saved, err
}
func (s ChecklistNonConformityService) ListByExecution(ctx context.Context, tenantID string, executionID string, q ports.Query) (ports.Page[domain.ChecklistNonConformity], error) {
	return s.repo.ListByExecution(ctx, tenantID, executionID, normalizeQuery(q))
}
func (s ChecklistHistoryService) ListByExecution(ctx context.Context, tenantID string, executionID string, q ports.Query) (ports.Page[domain.ChecklistHistory], error) {
	return s.repo.ListByExecution(ctx, tenantID, executionID, normalizeQuery(q))
}
func history(tenantID string, executionID string, event string, actorID string, notes string) domain.ChecklistHistory {
	return domain.ChecklistHistory{ID: newID("chh"), TenantID: tenantID, ExecutionID: executionID, Event: event, ActorID: actorID, Notes: notes, CreatedAt: time.Now().UTC()}
}

type ChecklistTypeService struct{ repo ports.ChecklistTypeRepository }
type ChecklistSectionService struct {
	repo ports.ChecklistSectionRepository
}
type ChecklistEngineItemService struct {
	versions ports.ChecklistTemplateVersionRepository
	repo     ports.ChecklistEngineItemRepository
}

func NewChecklistTypeService(repo ports.ChecklistTypeRepository) ChecklistTypeService {
	return ChecklistTypeService{repo: repo}
}
func NewChecklistSectionService(repo ports.ChecklistSectionRepository) ChecklistSectionService {
	return ChecklistSectionService{repo: repo}
}
func NewChecklistEngineItemService(versions ports.ChecklistTemplateVersionRepository, repo ports.ChecklistEngineItemRepository) ChecklistEngineItemService {
	return ChecklistEngineItemService{versions: versions, repo: repo}
}
func (s ChecklistTypeService) Create(ctx context.Context, item domain.ChecklistType) (domain.ChecklistType, error) {
	if item.TenantID == "" || strings.TrimSpace(item.Name) == "" || strings.TrimSpace(item.Code) == "" {
		return domain.ChecklistType{}, ErrValidation
	}
	now := time.Now().UTC()
	item.ID, item.Active, item.CreatedAt, item.UpdatedAt, item.Version = newID("cty"), true, now, now, 1
	return s.repo.Create(ctx, item)
}
func (s ChecklistTypeService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.ChecklistType], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(q))
}
func (s ChecklistSectionService) Create(ctx context.Context, item domain.ChecklistSection) (domain.ChecklistSection, error) {
	if item.TenantID == "" || item.TemplateVersionID == "" || strings.TrimSpace(item.Name) == "" {
		return domain.ChecklistSection{}, ErrValidation
	}
	now := time.Now().UTC()
	item.ID, item.Active, item.CreatedAt, item.UpdatedAt, item.Version = newID("cse"), true, now, now, 1
	return s.repo.Create(ctx, item)
}
func (s ChecklistSectionService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.ChecklistSection], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(q))
}
func (s ChecklistEngineItemService) Create(ctx context.Context, item domain.ChecklistEngineItem) (domain.ChecklistEngineItem, error) {
	version, err := s.versions.FindByID(ctx, item.TenantID, item.TemplateVersionID)
	if err != nil {
		return domain.ChecklistEngineItem{}, err
	}
	if version.Status == domain.TemplateStatusPublished {
		return domain.ChecklistEngineItem{}, ErrPublishedVersionImmutable
	}
	if strings.TrimSpace(item.Question) == "" || strings.TrimSpace(item.ItemType) == "" {
		return domain.ChecklistEngineItem{}, ErrValidation
	}
	now := time.Now().UTC()
	if item.Severity == "" {
		item.Severity = domain.SeverityMedium
	}
	item.ID, item.Active, item.CreatedAt, item.UpdatedAt, item.Version = newID("cit"), true, now, now, 1
	return s.repo.Create(ctx, item)
}
func (s ChecklistEngineItemService) ListByVersion(ctx context.Context, tenantID string, versionID string, q ports.Query) (ports.Page[domain.ChecklistEngineItem], error) {
	return s.repo.ListByVersion(ctx, tenantID, versionID, normalizeQuery(q))
}

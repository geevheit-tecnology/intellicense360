package infrastructure

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/ports"
)

type MemoryStore struct {
	mu              sync.RWMutex
	checklists      map[string]domain.Checklist
	items           map[string]domain.ChecklistItem
	answers         map[string]domain.ChecklistAnswer
	templates       map[string]domain.ChecklistTemplate
	versions        map[string]domain.ChecklistTemplateVersion
	types           map[string]domain.ChecklistType
	sections        map[string]domain.ChecklistSection
	engineItems     map[string]domain.ChecklistEngineItem
	executions      map[string]domain.ChecklistExecution
	responses       map[string]domain.ChecklistResponse
	evidence        map[string]domain.ChecklistEvidence
	nonConformities map[string]domain.ChecklistNonConformity
	signatures      map[string]domain.ChecklistSignature
	history         map[string]domain.ChecklistHistory
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		checklists:      map[string]domain.Checklist{},
		items:           map[string]domain.ChecklistItem{},
		answers:         map[string]domain.ChecklistAnswer{},
		templates:       map[string]domain.ChecklistTemplate{},
		versions:        map[string]domain.ChecklistTemplateVersion{},
		types:           map[string]domain.ChecklistType{},
		sections:        map[string]domain.ChecklistSection{},
		engineItems:     map[string]domain.ChecklistEngineItem{},
		executions:      map[string]domain.ChecklistExecution{},
		responses:       map[string]domain.ChecklistResponse{},
		evidence:        map[string]domain.ChecklistEvidence{},
		nonConformities: map[string]domain.ChecklistNonConformity{},
		signatures:      map[string]domain.ChecklistSignature{},
		history:         map[string]domain.ChecklistHistory{},
	}
}

func (s *MemoryStore) Checklists() ChecklistRepository          { return ChecklistRepository{s: s} }
func (s *MemoryStore) Items() ChecklistItemRepository           { return ChecklistItemRepository{s: s} }
func (s *MemoryStore) Answers() ChecklistAnswerRepository       { return ChecklistAnswerRepository{s: s} }
func (s *MemoryStore) Templates() TemplateRepository            { return TemplateRepository{s: s} }
func (s *MemoryStore) Versions() TemplateVersionRepository      { return TemplateVersionRepository{s: s} }
func (s *MemoryStore) Types() TypeRepository                    { return TypeRepository{s: s} }
func (s *MemoryStore) Sections() SectionRepository              { return SectionRepository{s: s} }
func (s *MemoryStore) EngineItems() EngineItemRepository        { return EngineItemRepository{s: s} }
func (s *MemoryStore) Executions() ExecutionRepository          { return ExecutionRepository{s: s} }
func (s *MemoryStore) Responses() ResponseRepository            { return ResponseRepository{s: s} }
func (s *MemoryStore) Evidence() EvidenceRepository             { return EvidenceRepository{s: s} }
func (s *MemoryStore) NonConformities() NonConformityRepository { return NonConformityRepository{s: s} }
func (s *MemoryStore) Signatures() SignatureRepository          { return SignatureRepository{s: s} }
func (s *MemoryStore) History() HistoryRepository               { return HistoryRepository{s: s} }

type ChecklistRepository struct{ s *MemoryStore }

func (r ChecklistRepository) Create(_ context.Context, checklist domain.Checklist) (domain.Checklist, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.checklists[key(checklist.TenantID, checklist.ID)] = checklist
	return checklist, nil
}

func (r ChecklistRepository) FindByID(_ context.Context, tenantID string, checklistID string) (domain.Checklist, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	checklist, ok := r.s.checklists[key(tenantID, checklistID)]
	if !ok || checklist.DeletedAt != nil {
		return domain.Checklist{}, application.ErrNotFound
	}
	return checklist, nil
}

func (r ChecklistRepository) Search(_ context.Context, tenantID string, query ports.Query) (ports.Page[domain.Checklist], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Checklist{}
	for _, checklist := range r.s.checklists {
		if checklist.TenantID == tenantID && checklist.DeletedAt == nil && matchesChecklist(checklist, query) {
			items = append(items, checklist)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if query.SortOrder == "desc" {
			return items[i].CreatedAt.After(items[j].CreatedAt)
		}
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})
	return page(items, query), nil
}

func (r ChecklistRepository) Update(_ context.Context, checklist domain.Checklist) (domain.Checklist, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(checklist.TenantID, checklist.ID)
	current, ok := r.s.checklists[k]
	if !ok || current.DeletedAt != nil {
		return domain.Checklist{}, application.ErrNotFound
	}
	checklist.DeletedAt = current.DeletedAt
	r.s.checklists[k] = checklist
	return checklist, nil
}

func (r ChecklistRepository) Delete(_ context.Context, tenantID string, checklistID string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, checklistID)
	checklist, ok := r.s.checklists[k]
	if !ok || checklist.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	checklist.DeletedAt = &now
	checklist.UpdatedAt = now
	r.s.checklists[k] = checklist
	return nil
}

func (r ChecklistRepository) Exists(_ context.Context, tenantID string, checklistID string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	checklist, ok := r.s.checklists[key(tenantID, checklistID)]
	return ok && checklist.DeletedAt == nil, nil
}

type ChecklistItemRepository struct{ s *MemoryStore }

func (r ChecklistItemRepository) Create(_ context.Context, item domain.ChecklistItem) (domain.ChecklistItem, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.items[itemKey(item.TenantID, item.ChecklistID, item.ID)] = item
	return item, nil
}

func (r ChecklistItemRepository) FindByID(_ context.Context, tenantID string, checklistID string, itemID string) (domain.ChecklistItem, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.items[itemKey(tenantID, checklistID, itemID)]
	if !ok || item.DeletedAt != nil {
		return domain.ChecklistItem{}, application.ErrNotFound
	}
	return item, nil
}

func (r ChecklistItemRepository) List(_ context.Context, tenantID string, checklistID string, query ports.Query) (ports.Page[domain.ChecklistItem], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.ChecklistItem{}
	for _, item := range r.s.items {
		if item.TenantID == tenantID && item.ChecklistID == checklistID && item.DeletedAt == nil && matchesItem(item, query) {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].OrderIndex < items[j].OrderIndex })
	return page(items, query), nil
}

func (r ChecklistItemRepository) Update(_ context.Context, item domain.ChecklistItem) (domain.ChecklistItem, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := itemKey(item.TenantID, item.ChecklistID, item.ID)
	current, ok := r.s.items[k]
	if !ok || current.DeletedAt != nil {
		return domain.ChecklistItem{}, application.ErrNotFound
	}
	item.DeletedAt = current.DeletedAt
	r.s.items[k] = item
	return item, nil
}

func (r ChecklistItemRepository) Delete(_ context.Context, tenantID string, checklistID string, itemID string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := itemKey(tenantID, checklistID, itemID)
	item, ok := r.s.items[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	r.s.items[k] = item
	return nil
}

type ChecklistAnswerRepository struct{ s *MemoryStore }

func (r ChecklistAnswerRepository) Create(_ context.Context, answer domain.ChecklistAnswer) (domain.ChecklistAnswer, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.answers[itemKey(answer.TenantID, answer.ChecklistID, answer.ID)] = answer
	return answer, nil
}

func (r ChecklistAnswerRepository) List(_ context.Context, tenantID string, checklistID string, query ports.Query) (ports.Page[domain.ChecklistAnswer], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.ChecklistAnswer{}
	for _, answer := range r.s.answers {
		if answer.TenantID == tenantID && answer.ChecklistID == checklistID {
			if query.Filters["checklist_item_id"] == "" || answer.ChecklistItemID == query.Filters["checklist_item_id"] {
				items = append(items, answer)
			}
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].AnsweredAt.Before(items[j].AnsweredAt) })
	return page(items, query), nil
}

func matchesChecklist(checklist domain.Checklist, query ports.Query) bool {
	if query.Filters["vehicle_id"] != "" && checklist.VehicleID != query.Filters["vehicle_id"] {
		return false
	}
	if query.Filters["status"] != "" && string(checklist.Status) != query.Filters["status"] {
		return false
	}
	if query.Filters["type"] != "" && checklist.Type != query.Filters["type"] {
		return false
	}
	return query.Search == "" || contains(checklist.Title, query.Search) || contains(checklist.Description, query.Search) || contains(checklist.DriverName, query.Search)
}

func matchesItem(item domain.ChecklistItem, query ports.Query) bool {
	if query.Filters["category"] != "" && item.Category != query.Filters["category"] {
		return false
	}
	return query.Search == "" || contains(item.Title, query.Search) || contains(item.Description, query.Search)
}

func contains(value string, search string) bool {
	return search == "" || strings.Contains(strings.ToLower(value), strings.ToLower(search))
}

func page[T any](items []T, query ports.Query) ports.Page[T] {
	pageNumber := query.Page
	if pageNumber <= 0 {
		pageNumber = 1
	}
	pageSize := query.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	total := len(items)
	start := (pageNumber - 1) * pageSize
	if start > total {
		start = total
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	totalPages := 0
	if total > 0 {
		totalPages = (total + pageSize - 1) / pageSize
	}
	return ports.Page[T]{Items: items[start:end], Page: pageNumber, PageSize: pageSize, TotalItems: total, TotalPages: totalPages}
}

func key(tenantID string, id string) string { return tenantID + ":" + id }
func itemKey(tenantID string, checklistID string, id string) string {
	return tenantID + ":" + checklistID + ":" + id
}

type TemplateRepository struct{ s *MemoryStore }
type TemplateVersionRepository struct{ s *MemoryStore }
type TypeRepository struct{ s *MemoryStore }
type SectionRepository struct{ s *MemoryStore }
type EngineItemRepository struct{ s *MemoryStore }
type ExecutionRepository struct{ s *MemoryStore }
type ResponseRepository struct{ s *MemoryStore }
type EvidenceRepository struct{ s *MemoryStore }
type NonConformityRepository struct{ s *MemoryStore }
type SignatureRepository struct{ s *MemoryStore }
type HistoryRepository struct{ s *MemoryStore }

func (r TemplateRepository) Create(_ context.Context, item domain.ChecklistTemplate) (domain.ChecklistTemplate, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.templates[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r TemplateRepository) FindByID(_ context.Context, tenantID string, id string) (domain.ChecklistTemplate, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.templates[key(tenantID, id)]
	if !ok || item.DeletedAt != nil {
		return domain.ChecklistTemplate{}, application.ErrNotFound
	}
	return item, nil
}
func (r TemplateRepository) Search(_ context.Context, tenantID string, query ports.Query) (ports.Page[domain.ChecklistTemplate], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.ChecklistTemplate{}
	for _, item := range r.s.templates {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name+" "+item.Description+" "+item.Type, query.Search) {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.Before(items[j].CreatedAt) })
	return page(items, query), nil
}
func (r TemplateRepository) Update(_ context.Context, item domain.ChecklistTemplate) (domain.ChecklistTemplate, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(item.TenantID, item.ID)
	current, ok := r.s.templates[k]
	if !ok || current.DeletedAt != nil {
		return domain.ChecklistTemplate{}, application.ErrNotFound
	}
	item.DeletedAt = current.DeletedAt
	r.s.templates[k] = item
	return item, nil
}
func (r TemplateRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, id)
	item, ok := r.s.templates[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	item.Version++
	r.s.templates[k] = item
	return nil
}

func (r TemplateVersionRepository) Create(_ context.Context, item domain.ChecklistTemplateVersion) (domain.ChecklistTemplateVersion, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.versions[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r TemplateVersionRepository) FindByID(_ context.Context, tenantID string, id string) (domain.ChecklistTemplateVersion, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.versions[key(tenantID, id)]
	if !ok || item.DeletedAt != nil {
		return domain.ChecklistTemplateVersion{}, application.ErrNotFound
	}
	return item, nil
}
func (r TemplateVersionRepository) ListByTemplate(_ context.Context, tenantID string, templateID string, query ports.Query) (ports.Page[domain.ChecklistTemplateVersion], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.ChecklistTemplateVersion{}
	for _, item := range r.s.versions {
		if item.TenantID == tenantID && item.TemplateID == templateID && item.DeletedAt == nil {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].VersionNumber < items[j].VersionNumber })
	return page(items, query), nil
}
func (r TemplateVersionRepository) Update(_ context.Context, item domain.ChecklistTemplateVersion) (domain.ChecklistTemplateVersion, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(item.TenantID, item.ID)
	current, ok := r.s.versions[k]
	if !ok || current.DeletedAt != nil {
		return domain.ChecklistTemplateVersion{}, application.ErrNotFound
	}
	item.CreatedAt, item.DeletedAt = current.CreatedAt, current.DeletedAt
	r.s.versions[k] = item
	return item, nil
}

func (r TypeRepository) Create(_ context.Context, item domain.ChecklistType) (domain.ChecklistType, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.types[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r TypeRepository) Search(_ context.Context, tenantID string, query ports.Query) (ports.Page[domain.ChecklistType], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.ChecklistType{}
	for _, item := range r.s.types {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name+" "+item.Code, query.Search) {
			items = append(items, item)
		}
	}
	return page(items, query), nil
}
func (r SectionRepository) Create(_ context.Context, item domain.ChecklistSection) (domain.ChecklistSection, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.sections[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r SectionRepository) Search(_ context.Context, tenantID string, query ports.Query) (ports.Page[domain.ChecklistSection], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.ChecklistSection{}
	for _, item := range r.s.sections {
		if item.TenantID == tenantID && item.DeletedAt == nil && (query.Filters["template_version_id"] == "" || item.TemplateVersionID == query.Filters["template_version_id"]) {
			items = append(items, item)
		}
	}
	return page(items, query), nil
}
func (r EngineItemRepository) Create(_ context.Context, item domain.ChecklistEngineItem) (domain.ChecklistEngineItem, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.engineItems[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r EngineItemRepository) FindByID(_ context.Context, tenantID string, id string) (domain.ChecklistEngineItem, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.engineItems[key(tenantID, id)]
	if !ok || item.DeletedAt != nil {
		return domain.ChecklistEngineItem{}, application.ErrNotFound
	}
	return item, nil
}
func (r EngineItemRepository) ListByVersion(_ context.Context, tenantID string, versionID string, query ports.Query) (ports.Page[domain.ChecklistEngineItem], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.ChecklistEngineItem{}
	for _, item := range r.s.engineItems {
		if item.TenantID == tenantID && item.TemplateVersionID == versionID && item.DeletedAt == nil {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].OrderIndex < items[j].OrderIndex })
	return page(items, query), nil
}

func (r ExecutionRepository) Create(_ context.Context, item domain.ChecklistExecution) (domain.ChecklistExecution, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.executions[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r ExecutionRepository) FindByID(_ context.Context, tenantID string, id string) (domain.ChecklistExecution, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.executions[key(tenantID, id)]
	if !ok || item.DeletedAt != nil {
		return domain.ChecklistExecution{}, application.ErrNotFound
	}
	return item, nil
}
func (r ExecutionRepository) Search(_ context.Context, tenantID string, query ports.Query) (ports.Page[domain.ChecklistExecution], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.ChecklistExecution{}
	for _, item := range r.s.executions {
		if item.TenantID == tenantID && item.DeletedAt == nil && (query.Filters["status"] == "" || string(item.Status) == query.Filters["status"]) {
			items = append(items, item)
		}
	}
	return page(items, query), nil
}
func (r ExecutionRepository) Update(_ context.Context, item domain.ChecklistExecution) (domain.ChecklistExecution, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(item.TenantID, item.ID)
	current, ok := r.s.executions[k]
	if !ok || current.DeletedAt != nil {
		return domain.ChecklistExecution{}, application.ErrNotFound
	}
	item.CreatedAt, item.DeletedAt = current.CreatedAt, current.DeletedAt
	r.s.executions[k] = item
	return item, nil
}

func (r ResponseRepository) Create(_ context.Context, item domain.ChecklistResponse) (domain.ChecklistResponse, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.responses[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r ResponseRepository) ListByExecution(_ context.Context, tenantID string, executionID string, query ports.Query) (ports.Page[domain.ChecklistResponse], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.ChecklistResponse{}
	for _, item := range r.s.responses {
		if item.TenantID == tenantID && item.ExecutionID == executionID {
			items = append(items, item)
		}
	}
	return page(items, query), nil
}
func (r ResponseRepository) ExistsForItem(_ context.Context, tenantID string, executionID string, itemID string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	for _, item := range r.s.responses {
		if item.TenantID == tenantID && item.ExecutionID == executionID && item.ItemID == itemID {
			return true, nil
		}
	}
	return false, nil
}
func (r EvidenceRepository) Create(_ context.Context, item domain.ChecklistEvidence) (domain.ChecklistEvidence, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.evidence[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r EvidenceRepository) ListByExecution(_ context.Context, tenantID string, executionID string, query ports.Query) (ports.Page[domain.ChecklistEvidence], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.ChecklistEvidence{}
	for _, item := range r.s.evidence {
		if item.TenantID == tenantID && item.ExecutionID == executionID {
			items = append(items, item)
		}
	}
	return page(items, query), nil
}
func (r EvidenceRepository) ExistsForResponse(_ context.Context, tenantID string, executionID string, responseID string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	for _, item := range r.s.evidence {
		if item.TenantID == tenantID && item.ExecutionID == executionID && item.ResponseID == responseID {
			return true, nil
		}
	}
	return false, nil
}
func (r NonConformityRepository) Create(_ context.Context, item domain.ChecklistNonConformity) (domain.ChecklistNonConformity, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.nonConformities[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r NonConformityRepository) ListByExecution(_ context.Context, tenantID string, executionID string, query ports.Query) (ports.Page[domain.ChecklistNonConformity], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.ChecklistNonConformity{}
	for _, item := range r.s.nonConformities {
		if item.TenantID == tenantID && item.ExecutionID == executionID {
			items = append(items, item)
		}
	}
	return page(items, query), nil
}
func (r SignatureRepository) Create(_ context.Context, item domain.ChecklistSignature) (domain.ChecklistSignature, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.signatures[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r SignatureRepository) ExistsForExecution(_ context.Context, tenantID string, executionID string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	for _, item := range r.s.signatures {
		if item.TenantID == tenantID && item.ExecutionID == executionID {
			return true, nil
		}
	}
	return false, nil
}
func (r HistoryRepository) Create(_ context.Context, item domain.ChecklistHistory) (domain.ChecklistHistory, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.history[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r HistoryRepository) ListByExecution(_ context.Context, tenantID string, executionID string, query ports.Query) (ports.Page[domain.ChecklistHistory], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.ChecklistHistory{}
	for _, item := range r.s.history {
		if item.TenantID == tenantID && item.ExecutionID == executionID {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.Before(items[j].CreatedAt) })
	return page(items, query), nil
}

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
	mu         sync.RWMutex
	checklists map[string]domain.Checklist
	items      map[string]domain.ChecklistItem
	answers    map[string]domain.ChecklistAnswer
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		checklists: map[string]domain.Checklist{},
		items:      map[string]domain.ChecklistItem{},
		answers:    map[string]domain.ChecklistAnswer{},
	}
}

func (s *MemoryStore) Checklists() ChecklistRepository    { return ChecklistRepository{s: s} }
func (s *MemoryStore) Items() ChecklistItemRepository     { return ChecklistItemRepository{s: s} }
func (s *MemoryStore) Answers() ChecklistAnswerRepository { return ChecklistAnswerRepository{s: s} }

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

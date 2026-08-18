package infrastructure

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/ports"
)

type MemoryStore struct {
	mu          sync.RWMutex
	items       map[string]domain.CommandItem
	events      map[string]domain.CommandEvent
	actions     map[string]domain.CommandAction
	snapshots   map[string]domain.OperationalSnapshot
	idempotency map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{items: map[string]domain.CommandItem{}, events: map[string]domain.CommandEvent{}, actions: map[string]domain.CommandAction{}, snapshots: map[string]domain.OperationalSnapshot{}, idempotency: map[string]string{}}
}

func (s *MemoryStore) Items() ItemRepository              { return ItemRepository{s: s} }
func (s *MemoryStore) Events() EventRepository            { return EventRepository{s: s} }
func (s *MemoryStore) Actions() ActionRepository          { return ActionRepository{s: s} }
func (s *MemoryStore) Snapshots() SnapshotRepository      { return SnapshotRepository{s: s} }
func (s *MemoryStore) Idempotency() IdempotencyRepository { return IdempotencyRepository{s: s} }

type ItemRepository struct{ s *MemoryStore }
type EventRepository struct{ s *MemoryStore }
type ActionRepository struct{ s *MemoryStore }
type SnapshotRepository struct{ s *MemoryStore }
type IdempotencyRepository struct{ s *MemoryStore }

func (r ItemRepository) Create(_ context.Context, item domain.CommandItem) (domain.CommandItem, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.items[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r ItemRepository) GetByID(_ context.Context, tenantID, id string) (domain.CommandItem, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.items[key(tenantID, id)]
	if !ok {
		return domain.CommandItem{}, application.ErrCommandItemNotFound
	}
	return item, nil
}
func (r ItemRepository) Update(_ context.Context, item domain.CommandItem) (domain.CommandItem, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(item.TenantID, item.ID)
	if _, ok := r.s.items[k]; !ok {
		return domain.CommandItem{}, application.ErrCommandItemNotFound
	}
	r.s.items[k] = item
	return item, nil
}
func (r ItemRepository) List(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.CommandItem], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.CommandItem{}
	for _, item := range r.s.items {
		if item.TenantID == tenantID && match(item, q) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r ItemRepository) FindByFingerprint(_ context.Context, tenantID, fingerprint string) (domain.CommandItem, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	for _, item := range r.s.items {
		if item.TenantID == tenantID && item.Fingerprint == fingerprint && item.Status != domain.StatusDismissed && item.Status != domain.StatusResolved {
			return item, nil
		}
	}
	return domain.CommandItem{}, application.ErrCommandItemNotFound
}
func (r ItemRepository) AllOpen(_ context.Context, tenantID string) ([]domain.CommandItem, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	out := []domain.CommandItem{}
	for _, item := range r.s.items {
		if item.TenantID == tenantID && item.Status != domain.StatusResolved && item.Status != domain.StatusDismissed && item.Status != domain.StatusExpired {
			out = append(out, item)
		}
	}
	return out, nil
}

func (r EventRepository) Create(_ context.Context, event domain.CommandEvent) (domain.CommandEvent, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.events[key(event.TenantID, event.ID)] = event
	return event, nil
}
func (r EventRepository) ListByItem(_ context.Context, tenantID, itemID string, q ports.Query) (ports.Page[domain.CommandEvent], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	out := []domain.CommandEvent{}
	for _, event := range r.s.events {
		if event.TenantID == tenantID && event.CommandItemID == itemID {
			out = append(out, event)
		}
	}
	return page(out, q), nil
}

func (r ActionRepository) Create(_ context.Context, action domain.CommandAction) (domain.CommandAction, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.actions[key(action.TenantID, action.ID)] = action
	return action, nil
}
func (r ActionRepository) ListByItem(_ context.Context, tenantID, itemID string, q ports.Query) (ports.Page[domain.CommandAction], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	out := []domain.CommandAction{}
	for _, action := range r.s.actions {
		if action.TenantID == tenantID && action.CommandItemID == itemID {
			out = append(out, action)
		}
	}
	return page(out, q), nil
}

func (r SnapshotRepository) Create(_ context.Context, snapshot domain.OperationalSnapshot) (domain.OperationalSnapshot, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.snapshots[key(snapshot.TenantID, snapshot.ID)] = snapshot
	return snapshot, nil
}
func (r SnapshotRepository) Latest(_ context.Context, tenantID string) (domain.OperationalSnapshot, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	var latest domain.OperationalSnapshot
	found := false
	for _, snapshot := range r.s.snapshots {
		if snapshot.TenantID == tenantID && (!found || snapshot.SnapshotAt.After(latest.SnapshotAt)) {
			latest = snapshot
			found = true
		}
	}
	if !found {
		return domain.OperationalSnapshot{}, application.ErrCommandItemNotFound
	}
	return latest, nil
}

func (r IdempotencyRepository) Exists(_ context.Context, tenantID, idempotencyKey string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	_, ok := r.s.idempotency[key(tenantID, idempotencyKey)]
	return ok, nil
}
func (r IdempotencyRepository) Save(_ context.Context, tenantID, idempotencyKey, resourceID string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.idempotency[key(tenantID, idempotencyKey)] = resourceID
	return nil
}

func key(tenantID, id string) string { return tenantID + ":" + id }

func match(item domain.CommandItem, q ports.Query) bool {
	f := q.Filters
	if f != nil {
		if f["type"] != "" && string(item.Type) != f["type"] {
			return false
		}
		if f["category"] != "" && string(item.Category) != f["category"] {
			return false
		}
		if f["severity"] != "" && string(item.Severity) != f["severity"] {
			return false
		}
		if f["priority"] != "" && string(item.Priority) != f["priority"] {
			return false
		}
		if f["status"] != "" && string(item.Status) != f["status"] {
			return false
		}
		if f["source_type"] != "" && item.SourceType != f["source_type"] {
			return false
		}
		if f["assigned_to"] != "" && item.AssignedTo != f["assigned_to"] {
			return false
		}
		if f["sla_status"] != "" && string(item.SLAStatus) != f["sla_status"] {
			return false
		}
	}
	if q.Search == "" {
		return true
	}
	text := strings.ToLower(item.Title + " " + item.Description + " " + item.SourceType + " " + item.SourceID)
	return strings.Contains(text, strings.ToLower(q.Search))
}

func page[T any](items []T, q ports.Query) ports.Page[T] {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PerPage <= 0 {
		q.PerPage = 25
	}
	total := len(items)
	start := (q.Page - 1) * q.PerPage
	if start > total {
		start = total
	}
	end := start + q.PerPage
	if end > total {
		end = total
	}
	totalPages := 0
	if total > 0 {
		totalPages = (total + q.PerPage - 1) / q.PerPage
	}
	_ = time.Now
	return ports.Page[T]{Data: items[start:end], Page: q.Page, PerPage: q.PerPage, Total: total, TotalPages: totalPages}
}

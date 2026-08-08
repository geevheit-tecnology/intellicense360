package infrastructure

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/tires/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/tires/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/tires/ports"
)

type MemoryStore struct {
	mu          sync.RWMutex
	tires       map[string]domain.Tire
	inspections map[string]domain.TireInspection
	movements   map[string]domain.TireMovement
	history     map[string]domain.TireHistory
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{tires: map[string]domain.Tire{}, inspections: map[string]domain.TireInspection{}, movements: map[string]domain.TireMovement{}, history: map[string]domain.TireHistory{}}
}

func (s *MemoryStore) Tires() TireRepository                 { return TireRepository{s: s} }
func (s *MemoryStore) Inspections() TireInspectionRepository { return TireInspectionRepository{s: s} }
func (s *MemoryStore) Movements() TireMovementRepository     { return TireMovementRepository{s: s} }
func (s *MemoryStore) History() TireHistoryRepository        { return TireHistoryRepository{s: s} }

type TireRepository struct{ s *MemoryStore }

func (r TireRepository) Create(_ context.Context, tire domain.Tire) (domain.Tire, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.tires[key(tire.TenantID, tire.ID)] = tire
	return tire, nil
}
func (r TireRepository) FindByID(_ context.Context, tenantID string, tireID string) (domain.Tire, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	tire, ok := r.s.tires[key(tenantID, tireID)]
	if !ok || tire.DeletedAt != nil {
		return domain.Tire{}, application.ErrNotFound
	}
	return tire, nil
}
func (r TireRepository) Search(_ context.Context, tenantID string, query ports.Query) (ports.Page[domain.Tire], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Tire{}
	for _, tire := range r.s.tires {
		if tire.TenantID == tenantID && tire.DeletedAt == nil && matchesTire(tire, query) {
			items = append(items, tire)
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
func (r TireRepository) Update(_ context.Context, tire domain.Tire) (domain.Tire, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tire.TenantID, tire.ID)
	current, ok := r.s.tires[k]
	if !ok || current.DeletedAt != nil {
		return domain.Tire{}, application.ErrNotFound
	}
	tire.DeletedAt = current.DeletedAt
	r.s.tires[k] = tire
	return tire, nil
}
func (r TireRepository) Delete(_ context.Context, tenantID string, tireID string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, tireID)
	tire, ok := r.s.tires[k]
	if !ok || tire.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	tire.DeletedAt = &now
	tire.UpdatedAt = now
	r.s.tires[k] = tire
	return nil
}
func (r TireRepository) Exists(_ context.Context, tenantID string, tireID string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	tire, ok := r.s.tires[key(tenantID, tireID)]
	return ok && tire.DeletedAt == nil, nil
}
func (r TireRepository) ExistsSerialNumber(_ context.Context, tenantID string, serialNumber string, exceptTireID string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	value := strings.ToUpper(strings.TrimSpace(serialNumber))
	for _, tire := range r.s.tires {
		if tire.TenantID == tenantID && tire.ID != exceptTireID && tire.DeletedAt == nil && strings.ToUpper(strings.TrimSpace(tire.SerialNumber)) == value {
			return true, nil
		}
	}
	return false, nil
}
func (r TireRepository) ExistsFireNumber(_ context.Context, tenantID string, fireNumber string, exceptTireID string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	value := strings.ToUpper(strings.TrimSpace(fireNumber))
	for _, tire := range r.s.tires {
		if tire.TenantID == tenantID && tire.ID != exceptTireID && tire.DeletedAt == nil && strings.ToUpper(strings.TrimSpace(tire.FireNumber)) == value {
			return true, nil
		}
	}
	return false, nil
}

type TireInspectionRepository struct{ s *MemoryStore }

func (r TireInspectionRepository) Create(_ context.Context, inspection domain.TireInspection) (domain.TireInspection, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.inspections[itemKey(inspection.TenantID, inspection.TireID, inspection.ID)] = inspection
	return inspection, nil
}
func (r TireInspectionRepository) FindByID(_ context.Context, tenantID string, tireID string, inspectionID string) (domain.TireInspection, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.inspections[itemKey(tenantID, tireID, inspectionID)]
	if !ok || item.DeletedAt != nil {
		return domain.TireInspection{}, application.ErrNotFound
	}
	return item, nil
}
func (r TireInspectionRepository) List(_ context.Context, tenantID string, tireID string, query ports.Query) (ports.Page[domain.TireInspection], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.TireInspection{}
	for _, item := range r.s.inspections {
		if item.TenantID == tenantID && item.TireID == tireID && item.DeletedAt == nil {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].InspectionDate.Before(items[j].InspectionDate) })
	return page(items, query), nil
}
func (r TireInspectionRepository) Update(_ context.Context, inspection domain.TireInspection) (domain.TireInspection, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := itemKey(inspection.TenantID, inspection.TireID, inspection.ID)
	current, ok := r.s.inspections[k]
	if !ok || current.DeletedAt != nil {
		return domain.TireInspection{}, application.ErrNotFound
	}
	inspection.DeletedAt = current.DeletedAt
	r.s.inspections[k] = inspection
	return inspection, nil
}
func (r TireInspectionRepository) Delete(_ context.Context, tenantID string, tireID string, inspectionID string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := itemKey(tenantID, tireID, inspectionID)
	item, ok := r.s.inspections[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	r.s.inspections[k] = item
	return nil
}

type TireMovementRepository struct{ s *MemoryStore }

func (r TireMovementRepository) Create(_ context.Context, movement domain.TireMovement) (domain.TireMovement, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.movements[itemKey(movement.TenantID, movement.TireID, movement.ID)] = movement
	return movement, nil
}

type TireHistoryRepository struct{ s *MemoryStore }

func (r TireHistoryRepository) Create(_ context.Context, history domain.TireHistory) (domain.TireHistory, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.history[itemKey(history.TenantID, history.TireID, history.ID)] = history
	return history, nil
}
func (r TireHistoryRepository) List(_ context.Context, tenantID string, tireID string, query ports.Query) (ports.Page[domain.TireHistory], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.TireHistory{}
	for _, item := range r.s.history {
		if item.TenantID == tenantID && item.TireID == tireID {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.Before(items[j].CreatedAt) })
	return page(items, query), nil
}
func (r TireMovementRepository) List(_ context.Context, tenantID string, tireID string, query ports.Query) (ports.Page[domain.TireMovement], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.TireMovement{}
	for _, item := range r.s.movements {
		if item.TenantID == tenantID && item.TireID == tireID && item.DeletedAt == nil {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].MovementDate.Before(items[j].MovementDate) })
	return page(items, query), nil
}

func matchesTire(tire domain.Tire, query ports.Query) bool {
	if query.Filters["vehicle_id"] != "" && tire.VehicleID != query.Filters["vehicle_id"] {
		return false
	}
	if query.Filters["status"] != "" && string(tire.Status) != query.Filters["status"] {
		return false
	}
	return query.Search == "" || contains(tire.SerialNumber, query.Search) || contains(tire.FireNumber, query.Search) || contains(tire.Brand, query.Search) || contains(tire.Model, query.Search) || contains(tire.DOT, query.Search)
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
func itemKey(tenantID string, tireID string, id string) string {
	return tenantID + ":" + tireID + ":" + id
}

package infrastructure

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/ports"
)

type MemoryStore struct {
	mu           sync.RWMutex
	workOrders   map[string]domain.WorkOrder
	plans        map[string]domain.PreventivePlan
	serviceTypes map[string]domain.ServiceType
	categories   map[string]domain.MaintenanceCatalog
	types        map[string]domain.MaintenanceCatalog
	priorities   map[string]domain.MaintenanceCatalog
	reasons      map[string]domain.MaintenanceCatalog
	workshops    map[string]domain.Workshop
	technicians  map[string]domain.Technician
	labor        map[string]domain.LaborEntry
	downtime     map[string]domain.Downtime
	history      map[string]domain.MaintenanceHistory
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		workOrders: map[string]domain.WorkOrder{}, plans: map[string]domain.PreventivePlan{}, serviceTypes: map[string]domain.ServiceType{},
		categories: map[string]domain.MaintenanceCatalog{}, types: map[string]domain.MaintenanceCatalog{}, priorities: map[string]domain.MaintenanceCatalog{}, reasons: map[string]domain.MaintenanceCatalog{},
		workshops: map[string]domain.Workshop{}, technicians: map[string]domain.Technician{},
		labor: map[string]domain.LaborEntry{}, downtime: map[string]domain.Downtime{}, history: map[string]domain.MaintenanceHistory{},
	}
}

func (s *MemoryStore) WorkOrders() WorkOrderRepository { return WorkOrderRepository{s: s} }
func (s *MemoryStore) PreventivePlans() PreventivePlanRepository {
	return PreventivePlanRepository{s: s}
}
func (s *MemoryStore) ServiceTypes() ServiceTypeRepository { return ServiceTypeRepository{s: s} }
func (s *MemoryStore) Categories() CatalogRepository {
	return CatalogRepository{s: s, items: s.categories}
}
func (s *MemoryStore) Types() CatalogRepository { return CatalogRepository{s: s, items: s.types} }
func (s *MemoryStore) Priorities() CatalogRepository {
	return CatalogRepository{s: s, items: s.priorities}
}
func (s *MemoryStore) Reasons() CatalogRepository        { return CatalogRepository{s: s, items: s.reasons} }
func (s *MemoryStore) Workshops() WorkshopRepository     { return WorkshopRepository{s: s} }
func (s *MemoryStore) Technicians() TechnicianRepository { return TechnicianRepository{s: s} }
func (s *MemoryStore) Labor() LaborRepository            { return LaborRepository{s: s} }
func (s *MemoryStore) Downtime() DowntimeRepository      { return DowntimeRepository{s: s} }
func (s *MemoryStore) History() HistoryRepository        { return HistoryRepository{s: s} }

type WorkOrderRepository struct{ s *MemoryStore }

func (r WorkOrderRepository) Create(_ context.Context, item domain.WorkOrder) (domain.WorkOrder, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.workOrders[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r WorkOrderRepository) FindByID(_ context.Context, tenantID string, id string) (domain.WorkOrder, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.workOrders[key(tenantID, id)]
	if !ok || item.DeletedAt != nil {
		return domain.WorkOrder{}, application.ErrNotFound
	}
	return item, nil
}
func (r WorkOrderRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.WorkOrder], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.WorkOrder{}
	for _, item := range r.s.workOrders {
		if item.TenantID == tenantID && item.DeletedAt == nil && matchesWO(item, q) {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if q.SortOrder == "desc" {
			return items[i].CreatedAt.After(items[j].CreatedAt)
		}
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})
	return page(items, q), nil
}
func (r WorkOrderRepository) Update(_ context.Context, item domain.WorkOrder) (domain.WorkOrder, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(item.TenantID, item.ID)
	current, ok := r.s.workOrders[k]
	if !ok || current.DeletedAt != nil {
		return domain.WorkOrder{}, application.ErrNotFound
	}
	item.DeletedAt = current.DeletedAt
	r.s.workOrders[k] = item
	return item, nil
}
func (r WorkOrderRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, id)
	item, ok := r.s.workOrders[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	r.s.workOrders[k] = item
	return nil
}
func (r WorkOrderRepository) Exists(_ context.Context, tenantID string, id string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.workOrders[key(tenantID, id)]
	return ok && item.DeletedAt == nil, nil
}

func (r WorkOrderRepository) ExistsCode(_ context.Context, tenantID string, code string, exceptID string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	value := strings.ToUpper(strings.TrimSpace(code))
	for _, item := range r.s.workOrders {
		if item.TenantID == tenantID && item.ID != exceptID && item.DeletedAt == nil && strings.ToUpper(strings.TrimSpace(item.Code)) == value {
			return true, nil
		}
	}
	return false, nil
}

type PreventivePlanRepository struct{ s *MemoryStore }

func (r PreventivePlanRepository) Create(_ context.Context, item domain.PreventivePlan) (domain.PreventivePlan, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.plans[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r PreventivePlanRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.PreventivePlan], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.PreventivePlan{}
	for _, item := range r.s.plans {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name, q.Search) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r PreventivePlanRepository) Update(_ context.Context, item domain.PreventivePlan) (domain.PreventivePlan, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.plans[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r PreventivePlanRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	item := r.s.plans[key(tenantID, id)]
	now := time.Now().UTC()
	item.DeletedAt = &now
	r.s.plans[key(tenantID, id)] = item
	return nil
}

type ServiceTypeRepository struct{ s *MemoryStore }

func (r ServiceTypeRepository) Create(_ context.Context, item domain.ServiceType) (domain.ServiceType, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.serviceTypes[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r ServiceTypeRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.ServiceType], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.ServiceType{}
	for _, item := range r.s.serviceTypes {
		if item.TenantID == tenantID && item.DeletedAt == nil && (contains(item.Name, q.Search) || contains(item.Code, q.Search)) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r ServiceTypeRepository) Update(_ context.Context, item domain.ServiceType) (domain.ServiceType, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.serviceTypes[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r ServiceTypeRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	item := r.s.serviceTypes[key(tenantID, id)]
	now := time.Now().UTC()
	item.DeletedAt = &now
	r.s.serviceTypes[key(tenantID, id)] = item
	return nil
}

type LaborRepository struct{ s *MemoryStore }

func (r LaborRepository) Create(_ context.Context, item domain.LaborEntry) (domain.LaborEntry, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.labor[itemKey(item.TenantID, item.WorkOrderID, item.ID)] = item
	return item, nil
}
func (r LaborRepository) List(_ context.Context, tenantID string, workOrderID string, q ports.Query) (ports.Page[domain.LaborEntry], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.LaborEntry{}
	for _, item := range r.s.labor {
		if item.TenantID == tenantID && item.WorkOrderID == workOrderID && item.DeletedAt == nil {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r LaborRepository) Delete(_ context.Context, tenantID string, workOrderID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := itemKey(tenantID, workOrderID, id)
	item := r.s.labor[k]
	now := time.Now().UTC()
	item.DeletedAt = &now
	r.s.labor[k] = item
	return nil
}

type DowntimeRepository struct{ s *MemoryStore }

func (r DowntimeRepository) Create(_ context.Context, item domain.Downtime) (domain.Downtime, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.downtime[itemKey(item.TenantID, item.WorkOrderID, item.ID)] = item
	return item, nil
}
func (r DowntimeRepository) List(_ context.Context, tenantID string, workOrderID string, q ports.Query) (ports.Page[domain.Downtime], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Downtime{}
	for _, item := range r.s.downtime {
		if item.TenantID == tenantID && item.WorkOrderID == workOrderID && item.DeletedAt == nil {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r DowntimeRepository) End(_ context.Context, tenantID string, workOrderID string, id string) (domain.Downtime, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := itemKey(tenantID, workOrderID, id)
	item, ok := r.s.downtime[k]
	if !ok || item.DeletedAt != nil {
		return domain.Downtime{}, application.ErrNotFound
	}
	now := time.Now().UTC()
	item.EndedAt = &now
	item.UpdatedAt = now
	r.s.downtime[k] = item
	return item, nil
}

type HistoryRepository struct{ s *MemoryStore }

func (r HistoryRepository) Create(_ context.Context, item domain.MaintenanceHistory) (domain.MaintenanceHistory, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.history[itemKey(item.TenantID, item.WorkOrderID, item.ID)] = item
	return item, nil
}
func (r HistoryRepository) List(_ context.Context, tenantID string, workOrderID string, q ports.Query) (ports.Page[domain.MaintenanceHistory], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.MaintenanceHistory{}
	for _, item := range r.s.history {
		if item.TenantID == tenantID && item.WorkOrderID == workOrderID {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}

type CatalogRepository struct {
	s     *MemoryStore
	items map[string]domain.MaintenanceCatalog
}

func (r CatalogRepository) Create(_ context.Context, item domain.MaintenanceCatalog) (domain.MaintenanceCatalog, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.items[key(item.TenantID, item.ID)] = item
	return item, nil
}

func (r CatalogRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.MaintenanceCatalog], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.MaintenanceCatalog{}
	for _, item := range r.items {
		if item.TenantID == tenantID && item.DeletedAt == nil && (contains(item.Name, q.Search) || contains(item.Code, q.Search)) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}

func (r CatalogRepository) Update(_ context.Context, item domain.MaintenanceCatalog) (domain.MaintenanceCatalog, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.items[key(item.TenantID, item.ID)] = item
	return item, nil
}

func (r CatalogRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	item := r.items[key(tenantID, id)]
	now := time.Now().UTC()
	item.DeletedAt = &now
	r.items[key(tenantID, id)] = item
	return nil
}

type WorkshopRepository struct{ s *MemoryStore }

func (r WorkshopRepository) Create(_ context.Context, item domain.Workshop) (domain.Workshop, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.workshops[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r WorkshopRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.Workshop], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Workshop{}
	for _, item := range r.s.workshops {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name, q.Search) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r WorkshopRepository) Update(_ context.Context, item domain.Workshop) (domain.Workshop, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.workshops[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r WorkshopRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	item := r.s.workshops[key(tenantID, id)]
	now := time.Now().UTC()
	item.DeletedAt = &now
	r.s.workshops[key(tenantID, id)] = item
	return nil
}

type TechnicianRepository struct{ s *MemoryStore }

func (r TechnicianRepository) Create(_ context.Context, item domain.Technician) (domain.Technician, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.technicians[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r TechnicianRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.Technician], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Technician{}
	for _, item := range r.s.technicians {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name, q.Search) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r TechnicianRepository) Update(_ context.Context, item domain.Technician) (domain.Technician, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.technicians[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r TechnicianRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	item := r.s.technicians[key(tenantID, id)]
	now := time.Now().UTC()
	item.DeletedAt = &now
	r.s.technicians[key(tenantID, id)] = item
	return nil
}

func matchesWO(item domain.WorkOrder, q ports.Query) bool {
	if q.Filters["status"] != "" && string(item.Status) != q.Filters["status"] {
		return false
	}
	if q.Filters["kind"] != "" && string(item.Kind) != q.Filters["kind"] {
		return false
	}
	if q.Filters["vehicle_id"] != "" && item.VehicleID != q.Filters["vehicle_id"] {
		return false
	}
	if q.Filters["asset_id"] != "" && item.AssetID != q.Filters["asset_id"] {
		return false
	}
	return q.Search == "" || contains(item.Title, q.Search) || contains(item.Description, q.Search)
}

func page[T any](items []T, q ports.Query) ports.Page[T] {
	p := q.Page
	if p <= 0 {
		p = 1
	}
	ps := q.PageSize
	if ps <= 0 {
		ps = 20
	}
	total := len(items)
	start := (p - 1) * ps
	if start > total {
		start = total
	}
	end := start + ps
	if end > total {
		end = total
	}
	totalPages := 0
	if total > 0 {
		totalPages = (total + ps - 1) / ps
	}
	return ports.Page[T]{Items: items[start:end], Page: p, PageSize: ps, TotalItems: total, TotalPages: totalPages}
}

func contains(v string, s string) bool {
	return s == "" || strings.Contains(strings.ToLower(v), strings.ToLower(s))
}
func key(tenantID string, id string) string { return tenantID + ":" + id }
func itemKey(tenantID string, parentID string, id string) string {
	return tenantID + ":" + parentID + ":" + id
}

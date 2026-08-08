package infrastructure

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/ports"
)

type MemoryStore struct {
	mu         sync.RWMutex
	parts      map[string]domain.Part
	categories map[string]domain.Catalog
	brands     map[string]domain.Catalog
	models     map[string]domain.Catalog
	units      map[string]domain.Catalog
	warehouses map[string]domain.Warehouse
	locations  map[string]domain.WarehouseLocation
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		parts: map[string]domain.Part{}, categories: map[string]domain.Catalog{}, brands: map[string]domain.Catalog{},
		models: map[string]domain.Catalog{}, units: map[string]domain.Catalog{}, warehouses: map[string]domain.Warehouse{}, locations: map[string]domain.WarehouseLocation{},
	}
}

func (s *MemoryStore) Parts() PartRepository { return PartRepository{s: s} }
func (s *MemoryStore) Categories() CatalogRepository {
	return CatalogRepository{s: s, items: s.categories}
}
func (s *MemoryStore) Brands() CatalogRepository       { return CatalogRepository{s: s, items: s.brands} }
func (s *MemoryStore) Models() CatalogRepository       { return CatalogRepository{s: s, items: s.models} }
func (s *MemoryStore) Units() CatalogRepository        { return CatalogRepository{s: s, items: s.units} }
func (s *MemoryStore) Warehouses() WarehouseRepository { return WarehouseRepository{s: s} }
func (s *MemoryStore) Locations() LocationRepository   { return LocationRepository{s: s} }

type PartRepository struct{ s *MemoryStore }

func (r PartRepository) Create(_ context.Context, item domain.Part) (domain.Part, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.parts[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r PartRepository) FindByID(_ context.Context, tenantID string, id string) (domain.Part, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.parts[key(tenantID, id)]
	if !ok || item.DeletedAt != nil {
		return domain.Part{}, application.ErrNotFound
	}
	return item, nil
}
func (r PartRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.Part], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Part{}
	for _, item := range r.s.parts {
		if item.TenantID == tenantID && item.DeletedAt == nil && matchesPart(item, q) {
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
func (r PartRepository) Update(_ context.Context, item domain.Part) (domain.Part, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(item.TenantID, item.ID)
	current, ok := r.s.parts[k]
	if !ok || current.DeletedAt != nil {
		return domain.Part{}, application.ErrNotFound
	}
	item.DeletedAt = current.DeletedAt
	r.s.parts[k] = item
	return item, nil
}
func (r PartRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, id)
	item, ok := r.s.parts[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	item.Version++
	r.s.parts[k] = item
	return nil
}
func (r PartRepository) ExistsSKU(_ context.Context, tenantID string, sku string, exceptID string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	value := norm(sku)
	for _, item := range r.s.parts {
		if item.TenantID == tenantID && item.ID != exceptID && item.DeletedAt == nil && norm(item.SKU) == value {
			return true, nil
		}
	}
	return false, nil
}
func (r PartRepository) ExistsInternalCode(_ context.Context, tenantID string, code string, exceptID string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	value := norm(code)
	for _, item := range r.s.parts {
		if item.TenantID == tenantID && item.ID != exceptID && item.DeletedAt == nil && norm(item.InternalCode) == value {
			return true, nil
		}
	}
	return false, nil
}

type CatalogRepository struct {
	s     *MemoryStore
	items map[string]domain.Catalog
}

func (r CatalogRepository) Create(_ context.Context, item domain.Catalog) (domain.Catalog, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.items[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r CatalogRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.Catalog], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Catalog{}
	for _, item := range r.items {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name+" "+item.Code, q.Search) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r CatalogRepository) Update(_ context.Context, item domain.Catalog) (domain.Catalog, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	current, ok := r.items[key(item.TenantID, item.ID)]
	if !ok || current.DeletedAt != nil {
		return domain.Catalog{}, application.ErrNotFound
	}
	item.CreatedAt = current.CreatedAt
	item.DeletedAt = current.DeletedAt
	r.items[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r CatalogRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, id)
	item, ok := r.items[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	item.Version++
	r.items[k] = item
	return nil
}
func (r CatalogRepository) ExistsCode(_ context.Context, tenantID string, code string, exceptID string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	value := norm(code)
	for _, item := range r.items {
		if item.TenantID == tenantID && item.ID != exceptID && item.DeletedAt == nil && norm(item.Code) == value {
			return true, nil
		}
	}
	return false, nil
}

type WarehouseRepository struct{ s *MemoryStore }

func (r WarehouseRepository) Create(_ context.Context, item domain.Warehouse) (domain.Warehouse, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.warehouses[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r WarehouseRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.Warehouse], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Warehouse{}
	for _, item := range r.s.warehouses {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name+" "+item.Code, q.Search) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r WarehouseRepository) Update(_ context.Context, item domain.Warehouse) (domain.Warehouse, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	current, ok := r.s.warehouses[key(item.TenantID, item.ID)]
	if !ok || current.DeletedAt != nil {
		return domain.Warehouse{}, application.ErrNotFound
	}
	item.CreatedAt = current.CreatedAt
	item.DeletedAt = current.DeletedAt
	r.s.warehouses[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r WarehouseRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, id)
	item, ok := r.s.warehouses[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	item.Version++
	r.s.warehouses[k] = item
	return nil
}

type LocationRepository struct{ s *MemoryStore }

func (r LocationRepository) Create(_ context.Context, item domain.WarehouseLocation) (domain.WarehouseLocation, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.locations[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r LocationRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.WarehouseLocation], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.WarehouseLocation{}
	for _, item := range r.s.locations {
		if item.TenantID == tenantID && item.DeletedAt == nil && matchesLocation(item, q) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r LocationRepository) Update(_ context.Context, item domain.WarehouseLocation) (domain.WarehouseLocation, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	current, ok := r.s.locations[key(item.TenantID, item.ID)]
	if !ok || current.DeletedAt != nil {
		return domain.WarehouseLocation{}, application.ErrNotFound
	}
	item.CreatedAt = current.CreatedAt
	item.DeletedAt = current.DeletedAt
	r.s.locations[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r LocationRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, id)
	item, ok := r.s.locations[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	item.Version++
	r.s.locations[k] = item
	return nil
}

func matchesPart(item domain.Part, q ports.Query) bool {
	if !contains(item.Name+" "+item.SKU+" "+item.InternalCode, q.Search) {
		return false
	}
	if q.Filters != nil {
		if v := q.Filters["status"]; v != "" && string(item.Status) != v {
			return false
		}
		if v := q.Filters["category_id"]; v != "" && item.CategoryID != v {
			return false
		}
	}
	return true
}
func matchesLocation(item domain.WarehouseLocation, q ports.Query) bool {
	if !contains(item.Name+" "+item.Code, q.Search) {
		return false
	}
	return q.Filters == nil || q.Filters["warehouse_id"] == "" || item.WarehouseID == q.Filters["warehouse_id"]
}
func key(tenantID, id string) string { return tenantID + ":" + id }
func norm(value string) string       { return strings.ToLower(strings.TrimSpace(value)) }
func contains(value, search string) bool {
	return search == "" || strings.Contains(norm(value), norm(search))
}
func page[T any](items []T, q ports.Query) ports.Page[T] {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 {
		q.PageSize = 20
	}
	total := len(items)
	start := (q.Page - 1) * q.PageSize
	if start > total {
		start = total
	}
	end := start + q.PageSize
	if end > total {
		end = total
	}
	totalPages := 0
	if total > 0 {
		totalPages = (total + q.PageSize - 1) / q.PageSize
	}
	return ports.Page[T]{Items: items[start:end], Page: q.Page, PageSize: q.PageSize, TotalItems: total, TotalPages: totalPages}
}

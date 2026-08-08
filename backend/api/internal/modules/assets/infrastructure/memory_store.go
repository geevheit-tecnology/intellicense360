package infrastructure

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/assets/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/assets/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/assets/ports"
)

type MemoryStore struct {
	mu            sync.RWMutex
	assets        map[string]domain.Asset
	categories    map[string]domain.AssetCategory
	types         map[string]domain.AssetType
	manufacturers map[string]domain.Manufacturer
	models        map[string]domain.Model
	equipment     map[string]domain.Equipment
	auditEvents   []AuditEvent
}

type AuditEvent struct {
	TenantID   string
	ActorID    string
	Action     string
	ResourceID string
	CreatedAt  time.Time
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		assets: map[string]domain.Asset{}, categories: map[string]domain.AssetCategory{}, types: map[string]domain.AssetType{},
		manufacturers: map[string]domain.Manufacturer{}, models: map[string]domain.Model{}, equipment: map[string]domain.Equipment{},
	}
}

func (s *MemoryStore) Assets() AssetRepository               { return AssetRepository{s: s} }
func (s *MemoryStore) Categories() CategoryRepository        { return CategoryRepository{s: s} }
func (s *MemoryStore) Types() TypeRepository                 { return TypeRepository{s: s} }
func (s *MemoryStore) Manufacturers() ManufacturerRepository { return ManufacturerRepository{s: s} }
func (s *MemoryStore) Models() ModelRepository               { return ModelRepository{s: s} }
func (s *MemoryStore) Equipment() EquipmentRepository        { return EquipmentRepository{s: s} }
func (s *MemoryStore) Audit() AuditRecorder                  { return AuditRecorder{s: s} }

type AuditRecorder struct{ s *MemoryStore }

func (r AuditRecorder) RecordAssetEvent(_ context.Context, tenantID string, actorID string, action string, resourceID string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.auditEvents = append(r.s.auditEvents, AuditEvent{TenantID: tenantID, ActorID: actorID, Action: action, ResourceID: resourceID, CreatedAt: time.Now().UTC()})
	return nil
}

type AssetRepository struct{ s *MemoryStore }

func (r AssetRepository) Create(_ context.Context, asset domain.Asset) (domain.Asset, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.assets[key(asset.TenantID, asset.ID)] = asset
	return asset, nil
}
func (r AssetRepository) FindByID(_ context.Context, tenantID string, assetID string) (domain.Asset, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.assets[key(tenantID, assetID)]
	if !ok || item.DeletedAt != nil {
		return domain.Asset{}, application.ErrNotFound
	}
	return item, nil
}
func (r AssetRepository) Search(_ context.Context, tenantID string, query ports.Query) (ports.Page[domain.Asset], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Asset{}
	for _, item := range r.s.assets {
		if item.TenantID == tenantID && item.DeletedAt == nil && matchesAsset(item, query) {
			items = append(items, item)
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
func (r AssetRepository) Update(_ context.Context, asset domain.Asset) (domain.Asset, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(asset.TenantID, asset.ID)
	current, ok := r.s.assets[k]
	if !ok || current.DeletedAt != nil {
		return domain.Asset{}, application.ErrNotFound
	}
	asset.DeletedAt = current.DeletedAt
	r.s.assets[k] = asset
	return asset, nil
}
func (r AssetRepository) Delete(_ context.Context, tenantID string, assetID string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, assetID)
	item, ok := r.s.assets[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	item.Version++
	r.s.assets[k] = item
	return nil
}
func (r AssetRepository) ExistsInternalCode(_ context.Context, tenantID string, internalCode string, exceptAssetID string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	value := norm(internalCode)
	for _, item := range r.s.assets {
		if item.TenantID == tenantID && item.ID != exceptAssetID && item.DeletedAt == nil && norm(item.InternalCode) == value {
			return true, nil
		}
	}
	return false, nil
}
func (r AssetRepository) ExistsSerialNumber(_ context.Context, tenantID string, serialNumber string, exceptAssetID string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	value := norm(serialNumber)
	for _, item := range r.s.assets {
		if item.TenantID == tenantID && item.ID != exceptAssetID && item.DeletedAt == nil && norm(item.SerialNumber) == value {
			return true, nil
		}
	}
	return false, nil
}
func (r AssetRepository) ExistsAssetTag(_ context.Context, tenantID string, assetTag string, exceptAssetID string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	value := norm(assetTag)
	for _, item := range r.s.assets {
		if item.TenantID == tenantID && item.ID != exceptAssetID && item.DeletedAt == nil && norm(item.AssetTag) == value {
			return true, nil
		}
	}
	return false, nil
}

type CategoryRepository struct{ s *MemoryStore }

func (r CategoryRepository) Create(_ context.Context, item domain.AssetCategory) (domain.AssetCategory, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.categories[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r CategoryRepository) Search(_ context.Context, tenantID string, query ports.Query) (ports.Page[domain.AssetCategory], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.AssetCategory{}
	for _, item := range r.s.categories {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name, query.Search) {
			items = append(items, item)
		}
	}
	return page(items, query), nil
}
func (r CategoryRepository) Update(_ context.Context, item domain.AssetCategory) (domain.AssetCategory, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.categories[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r CategoryRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	item := r.s.categories[key(tenantID, id)]
	now := time.Now().UTC()
	item.DeletedAt = &now
	r.s.categories[key(tenantID, id)] = item
	return nil
}

type TypeRepository struct{ s *MemoryStore }

func (r TypeRepository) Create(_ context.Context, item domain.AssetType) (domain.AssetType, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.types[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r TypeRepository) Search(_ context.Context, tenantID string, query ports.Query) (ports.Page[domain.AssetType], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.AssetType{}
	for _, item := range r.s.types {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name, query.Search) {
			items = append(items, item)
		}
	}
	return page(items, query), nil
}
func (r TypeRepository) Update(_ context.Context, item domain.AssetType) (domain.AssetType, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.types[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r TypeRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	item := r.s.types[key(tenantID, id)]
	now := time.Now().UTC()
	item.DeletedAt = &now
	r.s.types[key(tenantID, id)] = item
	return nil
}

type ManufacturerRepository struct{ s *MemoryStore }

func (r ManufacturerRepository) Create(_ context.Context, item domain.Manufacturer) (domain.Manufacturer, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.manufacturers[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r ManufacturerRepository) Search(_ context.Context, tenantID string, query ports.Query) (ports.Page[domain.Manufacturer], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Manufacturer{}
	for _, item := range r.s.manufacturers {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name, query.Search) {
			items = append(items, item)
		}
	}
	return page(items, query), nil
}
func (r ManufacturerRepository) Update(_ context.Context, item domain.Manufacturer) (domain.Manufacturer, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.manufacturers[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r ManufacturerRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	item := r.s.manufacturers[key(tenantID, id)]
	now := time.Now().UTC()
	item.DeletedAt = &now
	r.s.manufacturers[key(tenantID, id)] = item
	return nil
}

type ModelRepository struct{ s *MemoryStore }

func (r ModelRepository) Create(_ context.Context, item domain.Model) (domain.Model, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.models[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r ModelRepository) Search(_ context.Context, tenantID string, query ports.Query) (ports.Page[domain.Model], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Model{}
	for _, item := range r.s.models {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name, query.Search) {
			items = append(items, item)
		}
	}
	return page(items, query), nil
}
func (r ModelRepository) Update(_ context.Context, item domain.Model) (domain.Model, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.models[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r ModelRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	item := r.s.models[key(tenantID, id)]
	now := time.Now().UTC()
	item.DeletedAt = &now
	r.s.models[key(tenantID, id)] = item
	return nil
}

type EquipmentRepository struct{ s *MemoryStore }

func (r EquipmentRepository) Create(_ context.Context, item domain.Equipment) (domain.Equipment, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.equipment[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r EquipmentRepository) Search(_ context.Context, tenantID string, query ports.Query) (ports.Page[domain.Equipment], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Equipment{}
	for _, item := range r.s.equipment {
		if item.TenantID == tenantID && item.DeletedAt == nil && (contains(item.Category, query.Search) || contains(item.Type, query.Search)) {
			items = append(items, item)
		}
	}
	return page(items, query), nil
}
func (r EquipmentRepository) Update(_ context.Context, item domain.Equipment) (domain.Equipment, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.equipment[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r EquipmentRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	item := r.s.equipment[key(tenantID, id)]
	now := time.Now().UTC()
	item.DeletedAt = &now
	r.s.equipment[key(tenantID, id)] = item
	return nil
}

func matchesAsset(item domain.Asset, query ports.Query) bool {
	if query.Filters["status"] != "" && string(item.Status) != query.Filters["status"] {
		return false
	}
	if query.Filters["category_id"] != "" && item.CategoryID != query.Filters["category_id"] {
		return false
	}
	if query.Filters["type_id"] != "" && item.TypeID != query.Filters["type_id"] {
		return false
	}
	return query.Search == "" || contains(item.InternalCode, query.Search) || contains(item.SerialNumber, query.Search) || contains(item.AssetTag, query.Search) || contains(item.Name, query.Search)
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

func contains(value string, search string) bool {
	return search == "" || strings.Contains(strings.ToLower(value), strings.ToLower(search))
}
func norm(value string) string              { return strings.ToUpper(strings.TrimSpace(value)) }
func key(tenantID string, id string) string { return tenantID + ":" + id }

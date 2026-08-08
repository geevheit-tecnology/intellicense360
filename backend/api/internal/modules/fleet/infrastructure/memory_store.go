package infrastructure

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/ports"
)

type MemoryStore struct {
	mu         sync.RWMutex
	vehicles   map[string]domain.Vehicle
	brands     map[string]domain.VehicleBrand
	models     map[string]domain.VehicleModel
	categories map[string]domain.VehicleCategoryEntity
	types      map[string]domain.VehicleTypeEntity
	assets     map[string]domain.Asset
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		vehicles: map[string]domain.Vehicle{}, brands: map[string]domain.VehicleBrand{}, models: map[string]domain.VehicleModel{},
		categories: map[string]domain.VehicleCategoryEntity{}, types: map[string]domain.VehicleTypeEntity{}, assets: map[string]domain.Asset{},
	}
}

func (s *MemoryStore) Vehicles() VehicleRepository           { return VehicleRepository{s: s} }
func (s *MemoryStore) Brands() VehicleBrandRepository        { return VehicleBrandRepository{s: s} }
func (s *MemoryStore) Models() VehicleModelRepository        { return VehicleModelRepository{s: s} }
func (s *MemoryStore) Categories() VehicleCategoryRepository { return VehicleCategoryRepository{s: s} }
func (s *MemoryStore) Types() VehicleTypeRepository          { return VehicleTypeRepository{s: s} }
func (s *MemoryStore) Assets() AssetRepository               { return AssetRepository{s: s} }

type VehicleRepository struct{ s *MemoryStore }

func (r VehicleRepository) Create(_ context.Context, vehicle domain.Vehicle) (domain.Vehicle, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.vehicles[key(vehicle.TenantID, vehicle.ID)] = vehicle
	return vehicle, nil
}
func (r VehicleRepository) FindByID(_ context.Context, tenantID string, vehicleID string) (domain.Vehicle, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.vehicles[key(tenantID, vehicleID)]
	if !ok || item.DeletedAt != nil {
		return domain.Vehicle{}, application.ErrNotFound
	}
	return item, nil
}
func (r VehicleRepository) Search(_ context.Context, tenantID string, query ports.Query) (ports.Page[domain.Vehicle], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Vehicle{}
	for _, item := range r.s.vehicles {
		if item.TenantID == tenantID && item.DeletedAt == nil && matchesVehicle(item, query) {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].UpdatedAt.Before(items[j].UpdatedAt) })
	return page(items, query), nil
}
func (r VehicleRepository) Update(_ context.Context, vehicle domain.Vehicle) (domain.Vehicle, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(vehicle.TenantID, vehicle.ID)
	current, ok := r.s.vehicles[k]
	if !ok || current.DeletedAt != nil {
		return domain.Vehicle{}, application.ErrNotFound
	}
	vehicle.CreatedAt = current.CreatedAt
	r.s.vehicles[k] = vehicle
	return vehicle, nil
}
func (r VehicleRepository) Delete(_ context.Context, tenantID string, vehicleID string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, vehicleID)
	item, ok := r.s.vehicles[k]
	if !ok {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	item.Version++
	r.s.vehicles[k] = item
	return nil
}
func (r VehicleRepository) Exists(_ context.Context, tenantID string, vehicleID string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.vehicles[key(tenantID, vehicleID)]
	return ok && item.DeletedAt == nil, nil
}
func (r VehicleRepository) ExistsLicensePlate(_ context.Context, tenantID string, plate domain.LicensePlate, exceptVehicleID string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	normalized := plate.Normalized()
	for _, item := range r.s.vehicles {
		if item.TenantID == tenantID && item.ID != exceptVehicleID && item.DeletedAt == nil && item.LicensePlate.Normalized() == normalized {
			return true, nil
		}
	}
	return false, nil
}
func (r VehicleRepository) ExistsChassis(_ context.Context, tenantID string, chassis domain.Chassis, exceptVehicleID string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	value := strings.ToUpper(strings.TrimSpace(string(chassis)))
	for _, item := range r.s.vehicles {
		if item.TenantID == tenantID && item.ID != exceptVehicleID && item.DeletedAt == nil && strings.ToUpper(strings.TrimSpace(string(item.Chassis))) == value {
			return true, nil
		}
	}
	return false, nil
}
func (r VehicleRepository) ExistsRenavam(_ context.Context, tenantID string, renavam domain.Renavam, exceptVehicleID string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	value := strings.TrimSpace(string(renavam))
	for _, item := range r.s.vehicles {
		if item.TenantID == tenantID && item.ID != exceptVehicleID && item.DeletedAt == nil && strings.TrimSpace(string(item.Renavam)) == value {
			return true, nil
		}
	}
	return false, nil
}

type VehicleBrandRepository struct{ s *MemoryStore }

func (r VehicleBrandRepository) Create(_ context.Context, item domain.VehicleBrand) (domain.VehicleBrand, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.brands[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r VehicleBrandRepository) FindByID(_ context.Context, tenantID string, id string) (domain.VehicleBrand, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.brands[key(tenantID, id)]
	if !ok || item.DeletedAt != nil {
		return domain.VehicleBrand{}, application.ErrNotFound
	}
	return item, nil
}
func (r VehicleBrandRepository) Search(_ context.Context, tenantID string, query ports.Query) (ports.Page[domain.VehicleBrand], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.VehicleBrand{}
	for _, item := range r.s.brands {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name, query.Search) {
			items = append(items, item)
		}
	}
	return page(items, query), nil
}
func (r VehicleBrandRepository) Update(_ context.Context, item domain.VehicleBrand) (domain.VehicleBrand, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.brands[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r VehicleBrandRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	item := r.s.brands[key(tenantID, id)]
	now := time.Now().UTC()
	item.DeletedAt = &now
	r.s.brands[key(tenantID, id)] = item
	return nil
}

type VehicleModelRepository struct{ s *MemoryStore }

func (r VehicleModelRepository) Create(_ context.Context, item domain.VehicleModel) (domain.VehicleModel, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.models[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r VehicleModelRepository) FindByID(_ context.Context, tenantID string, id string) (domain.VehicleModel, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.models[key(tenantID, id)]
	if !ok || item.DeletedAt != nil {
		return domain.VehicleModel{}, application.ErrNotFound
	}
	return item, nil
}
func (r VehicleModelRepository) Search(_ context.Context, tenantID string, query ports.Query) (ports.Page[domain.VehicleModel], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.VehicleModel{}
	for _, item := range r.s.models {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name, query.Search) {
			items = append(items, item)
		}
	}
	return page(items, query), nil
}
func (r VehicleModelRepository) Update(_ context.Context, item domain.VehicleModel) (domain.VehicleModel, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.models[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r VehicleModelRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	item := r.s.models[key(tenantID, id)]
	now := time.Now().UTC()
	item.DeletedAt = &now
	r.s.models[key(tenantID, id)] = item
	return nil
}

type VehicleCategoryRepository struct{ s *MemoryStore }

func (r VehicleCategoryRepository) Create(_ context.Context, item domain.VehicleCategoryEntity) (domain.VehicleCategoryEntity, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.categories[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r VehicleCategoryRepository) FindByID(_ context.Context, tenantID string, id string) (domain.VehicleCategoryEntity, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.categories[key(tenantID, id)]
	if !ok || item.DeletedAt != nil {
		return domain.VehicleCategoryEntity{}, application.ErrNotFound
	}
	return item, nil
}
func (r VehicleCategoryRepository) Search(_ context.Context, tenantID string, query ports.Query) (ports.Page[domain.VehicleCategoryEntity], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.VehicleCategoryEntity{}
	for _, item := range r.s.categories {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name, query.Search) {
			items = append(items, item)
		}
	}
	return page(items, query), nil
}
func (r VehicleCategoryRepository) Update(_ context.Context, item domain.VehicleCategoryEntity) (domain.VehicleCategoryEntity, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.categories[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r VehicleCategoryRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	item := r.s.categories[key(tenantID, id)]
	now := time.Now().UTC()
	item.DeletedAt = &now
	r.s.categories[key(tenantID, id)] = item
	return nil
}

type VehicleTypeRepository struct{ s *MemoryStore }

func (r VehicleTypeRepository) Create(_ context.Context, item domain.VehicleTypeEntity) (domain.VehicleTypeEntity, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.types[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r VehicleTypeRepository) FindByID(_ context.Context, tenantID string, id string) (domain.VehicleTypeEntity, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.types[key(tenantID, id)]
	if !ok || item.DeletedAt != nil {
		return domain.VehicleTypeEntity{}, application.ErrNotFound
	}
	return item, nil
}
func (r VehicleTypeRepository) Search(_ context.Context, tenantID string, query ports.Query) (ports.Page[domain.VehicleTypeEntity], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.VehicleTypeEntity{}
	for _, item := range r.s.types {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name, query.Search) {
			items = append(items, item)
		}
	}
	return page(items, query), nil
}
func (r VehicleTypeRepository) Update(_ context.Context, item domain.VehicleTypeEntity) (domain.VehicleTypeEntity, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.types[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r VehicleTypeRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	item := r.s.types[key(tenantID, id)]
	now := time.Now().UTC()
	item.DeletedAt = &now
	r.s.types[key(tenantID, id)] = item
	return nil
}

type AssetRepository struct{ s *MemoryStore }

func (r AssetRepository) Create(_ context.Context, item domain.Asset) (domain.Asset, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.assets[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r AssetRepository) FindByID(_ context.Context, tenantID string, id string) (domain.Asset, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.assets[key(tenantID, id)]
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
		if item.TenantID == tenantID && item.DeletedAt == nil && (contains(item.Name, query.Search) || contains(item.Code, query.Search)) {
			items = append(items, item)
		}
	}
	return page(items, query), nil
}
func (r AssetRepository) Update(_ context.Context, item domain.Asset) (domain.Asset, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.assets[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r AssetRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	item := r.s.assets[key(tenantID, id)]
	now := time.Now().UTC()
	item.DeletedAt = &now
	r.s.assets[key(tenantID, id)] = item
	return nil
}

func matchesVehicle(item domain.Vehicle, query ports.Query) bool {
	if query.Filters["status"] != "" && string(item.Status) != query.Filters["status"] {
		return false
	}
	if query.Filters["category_id"] != "" && item.CategoryID != query.Filters["category_id"] {
		return false
	}
	return query.Search == "" || contains(string(item.LicensePlate), query.Search) || contains(string(item.Chassis), query.Search) || contains(string(item.Renavam), query.Search)
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

package infrastructure

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/ports"
)

type MemoryStore struct {
	mu           sync.RWMutex
	transactions map[string]domain.FuelTransaction
	types        map[string]domain.FuelType
	stations     map[string]domain.FuelStation
	tanks        map[string]domain.FuelTank
	nozzles      map[string]domain.FuelNozzle
	readings     map[string]domain.FuelReading
	prices       map[string]domain.FuelPrice
	receipts     map[string]domain.FuelReceipt
	adjustments  map[string]domain.FuelAdjustment
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		transactions: map[string]domain.FuelTransaction{},
		types:        map[string]domain.FuelType{},
		stations:     map[string]domain.FuelStation{},
		tanks:        map[string]domain.FuelTank{},
		nozzles:      map[string]domain.FuelNozzle{},
		readings:     map[string]domain.FuelReading{},
		prices:       map[string]domain.FuelPrice{},
		receipts:     map[string]domain.FuelReceipt{},
		adjustments:  map[string]domain.FuelAdjustment{},
	}
}

func (s *MemoryStore) Transactions() TransactionRepository { return TransactionRepository{s: s} }
func (s *MemoryStore) Types() TypeRepository               { return TypeRepository{s: s} }
func (s *MemoryStore) Stations() StationRepository         { return StationRepository{s: s} }
func (s *MemoryStore) Tanks() TankRepository               { return TankRepository{s: s} }
func (s *MemoryStore) Nozzles() NozzleRepository           { return NozzleRepository{s: s} }
func (s *MemoryStore) Readings() ReadingRepository         { return ReadingRepository{s: s} }
func (s *MemoryStore) Prices() PriceRepository             { return PriceRepository{s: s} }
func (s *MemoryStore) Receipts() ReceiptRepository         { return ReceiptRepository{s: s} }
func (s *MemoryStore) Adjustments() AdjustmentRepository   { return AdjustmentRepository{s: s} }

type TransactionRepository struct{ s *MemoryStore }

func (r TransactionRepository) Create(_ context.Context, item domain.FuelTransaction) (domain.FuelTransaction, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.transactions[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r TransactionRepository) FindByID(_ context.Context, tenantID string, id string) (domain.FuelTransaction, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.transactions[key(tenantID, id)]
	if !ok || item.DeletedAt != nil {
		return domain.FuelTransaction{}, application.ErrNotFound
	}
	return item, nil
}
func (r TransactionRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelTransaction], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.FuelTransaction{}
	for _, item := range r.s.transactions {
		if item.TenantID == tenantID && item.DeletedAt == nil && matchesTransaction(item, q) {
			items = append(items, item)
		}
	}
	sortByCreated(items, q.SortOrder, func(i int) time.Time { return items[i].CreatedAt })
	return page(items, q), nil
}
func (r TransactionRepository) Update(_ context.Context, item domain.FuelTransaction) (domain.FuelTransaction, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(item.TenantID, item.ID)
	current, ok := r.s.transactions[k]
	if !ok || current.DeletedAt != nil {
		return domain.FuelTransaction{}, application.ErrNotFound
	}
	item.CreatedAt = current.CreatedAt
	item.DeletedAt = current.DeletedAt
	r.s.transactions[k] = item
	return item, nil
}
func (r TransactionRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, id)
	item, ok := r.s.transactions[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	item.Version++
	r.s.transactions[k] = item
	return nil
}

type TypeRepository struct{ s *MemoryStore }

func (r TypeRepository) Create(_ context.Context, item domain.FuelType) (domain.FuelType, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.types[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r TypeRepository) FindByID(_ context.Context, tenantID string, id string) (domain.FuelType, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.types[key(tenantID, id)]
	if !ok || item.DeletedAt != nil {
		return domain.FuelType{}, application.ErrNotFound
	}
	return item, nil
}
func (r TypeRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelType], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.FuelType{}
	for _, item := range r.s.types {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name+" "+item.Code+" "+string(item.Kind), q.Search) {
			items = append(items, item)
		}
	}
	sortByCreated(items, q.SortOrder, func(i int) time.Time { return items[i].CreatedAt })
	return page(items, q), nil
}
func (r TypeRepository) Update(_ context.Context, item domain.FuelType) (domain.FuelType, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(item.TenantID, item.ID)
	current, ok := r.s.types[k]
	if !ok || current.DeletedAt != nil {
		return domain.FuelType{}, application.ErrNotFound
	}
	item.CreatedAt = current.CreatedAt
	item.DeletedAt = current.DeletedAt
	r.s.types[k] = item
	return item, nil
}
func (r TypeRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, id)
	item, ok := r.s.types[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	item.Version++
	r.s.types[k] = item
	return nil
}

type StationRepository struct{ s *MemoryStore }

func (r StationRepository) Create(_ context.Context, item domain.FuelStation) (domain.FuelStation, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.stations[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r StationRepository) FindByID(_ context.Context, tenantID string, id string) (domain.FuelStation, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.stations[key(tenantID, id)]
	if !ok || item.DeletedAt != nil {
		return domain.FuelStation{}, application.ErrNotFound
	}
	return item, nil
}
func (r StationRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelStation], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.FuelStation{}
	for _, item := range r.s.stations {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name+" "+item.LegalName+" "+item.CNPJ+" "+item.City, q.Search) {
			items = append(items, item)
		}
	}
	sortByCreated(items, q.SortOrder, func(i int) time.Time { return items[i].CreatedAt })
	return page(items, q), nil
}
func (r StationRepository) Update(_ context.Context, item domain.FuelStation) (domain.FuelStation, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(item.TenantID, item.ID)
	current, ok := r.s.stations[k]
	if !ok || current.DeletedAt != nil {
		return domain.FuelStation{}, application.ErrNotFound
	}
	item.CreatedAt = current.CreatedAt
	item.DeletedAt = current.DeletedAt
	r.s.stations[k] = item
	return item, nil
}
func (r StationRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, id)
	item, ok := r.s.stations[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	item.Version++
	r.s.stations[k] = item
	return nil
}
func (r StationRepository) ExistsCNPJ(_ context.Context, tenantID string, cnpj string, exceptID string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	for _, item := range r.s.stations {
		if item.TenantID == tenantID && item.ID != exceptID && item.DeletedAt == nil && item.CNPJ == cnpj {
			return true, nil
		}
	}
	return false, nil
}

type TankRepository struct{ s *MemoryStore }
type NozzleRepository struct{ s *MemoryStore }
type ReadingRepository struct{ s *MemoryStore }
type PriceRepository struct{ s *MemoryStore }
type ReceiptRepository struct{ s *MemoryStore }
type AdjustmentRepository struct{ s *MemoryStore }

func (r TankRepository) Create(_ context.Context, item domain.FuelTank) (domain.FuelTank, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.tanks[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r TankRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelTank], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.FuelTank{}
	for _, item := range r.s.tanks {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Code+" "+item.Name+" "+string(item.FuelKind), q.Search) {
			items = append(items, item)
		}
	}
	sortByCreated(items, q.SortOrder, func(i int) time.Time { return items[i].CreatedAt })
	return page(items, q), nil
}
func (r TankRepository) Update(_ context.Context, item domain.FuelTank) (domain.FuelTank, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(item.TenantID, item.ID)
	current, ok := r.s.tanks[k]
	if !ok || current.DeletedAt != nil {
		return domain.FuelTank{}, application.ErrNotFound
	}
	item.CreatedAt, item.DeletedAt = current.CreatedAt, current.DeletedAt
	if item.Version <= current.Version {
		item.Version = current.Version + 1
	}
	r.s.tanks[k] = item
	return item, nil
}
func (r TankRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, id)
	item, ok := r.s.tanks[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt, item.UpdatedAt, item.Version = &now, now, item.Version+1
	r.s.tanks[k] = item
	return nil
}

func (r NozzleRepository) Create(_ context.Context, item domain.FuelNozzle) (domain.FuelNozzle, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.nozzles[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r NozzleRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelNozzle], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.FuelNozzle{}
	for _, item := range r.s.nozzles {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Code+" "+string(item.FuelKind), q.Search) {
			items = append(items, item)
		}
	}
	sortByCreated(items, q.SortOrder, func(i int) time.Time { return items[i].CreatedAt })
	return page(items, q), nil
}
func (r NozzleRepository) Update(_ context.Context, item domain.FuelNozzle) (domain.FuelNozzle, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(item.TenantID, item.ID)
	current, ok := r.s.nozzles[k]
	if !ok || current.DeletedAt != nil {
		return domain.FuelNozzle{}, application.ErrNotFound
	}
	item.CreatedAt, item.DeletedAt = current.CreatedAt, current.DeletedAt
	if item.Version <= current.Version {
		item.Version = current.Version + 1
	}
	r.s.nozzles[k] = item
	return item, nil
}
func (r NozzleRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, id)
	item, ok := r.s.nozzles[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt, item.UpdatedAt, item.Version = &now, now, item.Version+1
	r.s.nozzles[k] = item
	return nil
}

func (r ReadingRepository) Create(_ context.Context, item domain.FuelReading) (domain.FuelReading, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.readings[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r ReadingRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelReading], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.FuelReading{}
	for _, item := range r.s.readings {
		if item.TenantID == tenantID && contains(string(item.ReadingType)+" "+item.ReferenceID+" "+item.Source, q.Search) {
			items = append(items, item)
		}
	}
	sortByCreated(items, q.SortOrder, func(i int) time.Time { return items[i].CreatedAt })
	return page(items, q), nil
}

func (r PriceRepository) Create(_ context.Context, item domain.FuelPrice) (domain.FuelPrice, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.prices[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r PriceRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelPrice], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.FuelPrice{}
	for _, item := range r.s.prices {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(string(item.FuelKind)+" "+item.StationID+" "+item.Source, q.Search) {
			items = append(items, item)
		}
	}
	sortByCreated(items, q.SortOrder, func(i int) time.Time { return items[i].CreatedAt })
	return page(items, q), nil
}
func (r PriceRepository) Update(_ context.Context, item domain.FuelPrice) (domain.FuelPrice, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(item.TenantID, item.ID)
	current, ok := r.s.prices[k]
	if !ok || current.DeletedAt != nil {
		return domain.FuelPrice{}, application.ErrNotFound
	}
	item.CreatedAt, item.DeletedAt = current.CreatedAt, current.DeletedAt
	if item.Version <= current.Version {
		item.Version = current.Version + 1
	}
	r.s.prices[k] = item
	return item, nil
}
func (r PriceRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, id)
	item, ok := r.s.prices[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt, item.UpdatedAt, item.Version = &now, now, item.Version+1
	r.s.prices[k] = item
	return nil
}

func (r ReceiptRepository) Create(_ context.Context, item domain.FuelReceipt) (domain.FuelReceipt, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.receipts[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r ReceiptRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelReceipt], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.FuelReceipt{}
	for _, item := range r.s.receipts {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.ReceiptNumber+" "+item.AttachmentReference, q.Search) {
			items = append(items, item)
		}
	}
	sortByCreated(items, q.SortOrder, func(i int) time.Time { return items[i].CreatedAt })
	return page(items, q), nil
}
func (r ReceiptRepository) Update(_ context.Context, item domain.FuelReceipt) (domain.FuelReceipt, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(item.TenantID, item.ID)
	current, ok := r.s.receipts[k]
	if !ok || current.DeletedAt != nil {
		return domain.FuelReceipt{}, application.ErrNotFound
	}
	item.CreatedAt, item.DeletedAt = current.CreatedAt, current.DeletedAt
	if item.Version <= current.Version {
		item.Version = current.Version + 1
	}
	r.s.receipts[k] = item
	return item, nil
}
func (r ReceiptRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, id)
	item, ok := r.s.receipts[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt, item.UpdatedAt, item.Version = &now, now, item.Version+1
	r.s.receipts[k] = item
	return nil
}

func (r AdjustmentRepository) Create(_ context.Context, item domain.FuelAdjustment) (domain.FuelAdjustment, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.adjustments[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r AdjustmentRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelAdjustment], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.FuelAdjustment{}
	for _, item := range r.s.adjustments {
		if item.TenantID == tenantID && contains(item.TransactionID+" "+item.AdjustmentType+" "+item.Reason, q.Search) {
			items = append(items, item)
		}
	}
	sortByCreated(items, q.SortOrder, func(i int) time.Time { return items[i].CreatedAt })
	return page(items, q), nil
}

func key(tenantID string, id string) string { return tenantID + ":" + id }
func contains(value string, search string) bool {
	search = strings.TrimSpace(strings.ToLower(search))
	if search == "" {
		return true
	}
	return strings.Contains(strings.ToLower(value), search)
}
func matchesTransaction(item domain.FuelTransaction, q ports.Query) bool {
	if q.Filters["status"] != "" && string(item.Status) != q.Filters["status"] {
		return false
	}
	if q.Filters["vehicle_reference"] != "" && item.VehicleReference != q.Filters["vehicle_reference"] {
		return false
	}
	return contains(item.ReceiptNumber+" "+item.DriverReference+" "+item.VehicleReference+" "+item.AssetReference+" "+string(item.FuelKind), q.Search)
}
func sortByCreated[T any](items []T, order string, createdAt func(int) time.Time) {
	sort.Slice(items, func(i, j int) bool {
		if order == "desc" {
			return createdAt(i).After(createdAt(j))
		}
		return createdAt(i).Before(createdAt(j))
	})
}
func page[T any](items []T, q ports.Query) ports.Page[T] {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
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

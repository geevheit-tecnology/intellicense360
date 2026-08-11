package infrastructure

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/ciot/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/ciot/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/ciot/ports"
)

type MemoryStore struct {
	mu                 sync.RWMutex
	ciots              map[string]domain.CIOT
	contracts          map[string]domain.CIOTContract
	carriers           map[string]domain.CIOTCarrier
	transporters       map[string]domain.CIOTTransporter
	operations         map[string]domain.CIOTOperation
	vehicles           map[string]domain.CIOTVehicleReference
	drivers            map[string]domain.CIOTDriverReference
	amounts            map[string]domain.CIOTAmount
	payments           map[string]domain.CIOTPayment
	history            map[string]domain.CIOTStatusHistory
	attempts           map[string]domain.CIOTProviderAttempt
	externalReferences map[string]domain.CIOTExternalReference
	documents          map[string]domain.CIOTDocument
	errors             map[string]domain.CIOTError
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		ciots:              map[string]domain.CIOT{},
		contracts:          map[string]domain.CIOTContract{},
		carriers:           map[string]domain.CIOTCarrier{},
		transporters:       map[string]domain.CIOTTransporter{},
		operations:         map[string]domain.CIOTOperation{},
		vehicles:           map[string]domain.CIOTVehicleReference{},
		drivers:            map[string]domain.CIOTDriverReference{},
		amounts:            map[string]domain.CIOTAmount{},
		payments:           map[string]domain.CIOTPayment{},
		history:            map[string]domain.CIOTStatusHistory{},
		attempts:           map[string]domain.CIOTProviderAttempt{},
		externalReferences: map[string]domain.CIOTExternalReference{},
		documents:          map[string]domain.CIOTDocument{},
		errors:             map[string]domain.CIOTError{},
	}
}

func (s *MemoryStore) CIOTs() CIOTRepository               { return CIOTRepository{s: s} }
func (s *MemoryStore) Contracts() ContractRepository       { return ContractRepository{s: s} }
func (s *MemoryStore) Carriers() CarrierRepository         { return CarrierRepository{s: s} }
func (s *MemoryStore) Transporters() TransporterRepository { return TransporterRepository{s: s} }
func (s *MemoryStore) Operations() OperationRepository     { return OperationRepository{s: s} }
func (s *MemoryStore) Vehicles() VehicleRepository         { return VehicleRepository{s: s} }
func (s *MemoryStore) Drivers() DriverRepository           { return DriverRepository{s: s} }
func (s *MemoryStore) Amounts() AmountRepository           { return AmountRepository{s: s} }
func (s *MemoryStore) Payments() PaymentRepository         { return PaymentRepository{s: s} }
func (s *MemoryStore) History() HistoryRepository          { return HistoryRepository{s: s} }
func (s *MemoryStore) Attempts() ProviderAttemptRepository { return ProviderAttemptRepository{s: s} }
func (s *MemoryStore) ExternalReferences() ExternalReferenceRepository {
	return ExternalReferenceRepository{s: s}
}
func (s *MemoryStore) Documents() DocumentRepository { return DocumentRepository{s: s} }
func (s *MemoryStore) Errors() ErrorRepository       { return ErrorRepository{s: s} }

type CIOTRepository struct{ s *MemoryStore }

func (r CIOTRepository) Create(_ context.Context, item domain.CIOT) (domain.CIOT, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.ciots[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r CIOTRepository) FindByID(_ context.Context, tenantID, id string) (domain.CIOT, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.ciots[key(tenantID, id)]
	if !ok || item.DeletedAt != nil {
		return domain.CIOT{}, application.ErrNotFound
	}
	return item, nil
}
func (r CIOTRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.CIOT], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.CIOT{}
	for _, item := range r.s.ciots {
		if item.TenantID == tenantID && item.DeletedAt == nil && matchCIOT(item, q) {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.Before(items[j].CreatedAt) })
	return page(items, q), nil
}
func (r CIOTRepository) Update(_ context.Context, item domain.CIOT) (domain.CIOT, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(item.TenantID, item.ID)
	current, ok := r.s.ciots[k]
	if !ok || current.DeletedAt != nil {
		return domain.CIOT{}, application.ErrNotFound
	}
	item.CreatedAt, item.DeletedAt = current.CreatedAt, current.DeletedAt
	r.s.ciots[k] = item
	return item, nil
}
func (r CIOTRepository) Delete(_ context.Context, tenantID, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, id)
	item, ok := r.s.ciots[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	item.Version++
	r.s.ciots[k] = item
	return nil
}
func (r CIOTRepository) ExistsIdempotencyKey(_ context.Context, tenantID string, idempotencyKey string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	for _, item := range r.s.ciots {
		if item.TenantID == tenantID && item.DeletedAt == nil && item.IdempotencyKey == idempotencyKey {
			return true, nil
		}
	}
	return false, nil
}

type ContractRepository struct{ s *MemoryStore }
type CarrierRepository struct{ s *MemoryStore }
type TransporterRepository struct{ s *MemoryStore }
type OperationRepository struct{ s *MemoryStore }
type VehicleRepository struct{ s *MemoryStore }
type DriverRepository struct{ s *MemoryStore }
type AmountRepository struct{ s *MemoryStore }

func (r ContractRepository) Create(_ context.Context, item domain.CIOTContract) (domain.CIOTContract, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.contracts[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r ContractRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.CIOTContract], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.CIOTContract{}
	for _, item := range r.s.contracts {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.ContractNumber+" "+item.ContractType, q.Search) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r CarrierRepository) Create(_ context.Context, item domain.CIOTCarrier) (domain.CIOTCarrier, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.carriers[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r CarrierRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.CIOTCarrier], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.CIOTCarrier{}
	for _, item := range r.s.carriers {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Document+" "+item.LegalName+" "+item.TradeName, q.Search) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r TransporterRepository) Create(_ context.Context, item domain.CIOTTransporter) (domain.CIOTTransporter, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.transporters[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r TransporterRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.CIOTTransporter], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.CIOTTransporter{}
	for _, item := range r.s.transporters {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Document+" "+item.Name, q.Search) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r OperationRepository) Create(_ context.Context, item domain.CIOTOperation) (domain.CIOTOperation, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.operations[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r OperationRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.CIOTOperation], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.CIOTOperation{}
	for _, item := range r.s.operations {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.OperationNumber+" "+item.Origin+" "+item.Destination, q.Search) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r VehicleRepository) Create(_ context.Context, item domain.CIOTVehicleReference) (domain.CIOTVehicleReference, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.vehicles[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r VehicleRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.CIOTVehicleReference], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.CIOTVehicleReference{}
	for _, item := range r.s.vehicles {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.VehicleID+" "+item.LicensePlate, q.Search) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r DriverRepository) Create(_ context.Context, item domain.CIOTDriverReference) (domain.CIOTDriverReference, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.drivers[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r DriverRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.CIOTDriverReference], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.CIOTDriverReference{}
	for _, item := range r.s.drivers {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.DriverID+" "+item.NameReference+" "+item.DocumentReference, q.Search) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r AmountRepository) Create(_ context.Context, item domain.CIOTAmount) (domain.CIOTAmount, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.amounts[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r AmountRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.CIOTAmount], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.CIOTAmount{}
	for _, item := range r.s.amounts {
		if item.TenantID == tenantID {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}

type PaymentRepository struct{ s *MemoryStore }

func (r PaymentRepository) Create(_ context.Context, item domain.CIOTPayment) (domain.CIOTPayment, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.payments[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r PaymentRepository) ListByCIOT(_ context.Context, tenantID string, ciotID string, q ports.Query) (ports.Page[domain.CIOTPayment], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.CIOTPayment{}
	for _, item := range r.s.payments {
		if item.TenantID == tenantID && item.CIOTID == ciotID && item.DeletedAt == nil {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}

type HistoryRepository struct{ s *MemoryStore }

func (r HistoryRepository) Create(_ context.Context, item domain.CIOTStatusHistory) (domain.CIOTStatusHistory, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.history[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r HistoryRepository) ListByCIOT(_ context.Context, tenantID string, ciotID string, q ports.Query) (ports.Page[domain.CIOTStatusHistory], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.CIOTStatusHistory{}
	for _, item := range r.s.history {
		if item.TenantID == tenantID && item.CIOTID == ciotID {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.Before(items[j].CreatedAt) })
	return page(items, q), nil
}

type ProviderAttemptRepository struct{ s *MemoryStore }

func (r ProviderAttemptRepository) Create(_ context.Context, item domain.CIOTProviderAttempt) (domain.CIOTProviderAttempt, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.attempts[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r ProviderAttemptRepository) ListByCIOT(_ context.Context, tenantID string, ciotID string, q ports.Query) (ports.Page[domain.CIOTProviderAttempt], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.CIOTProviderAttempt{}
	for _, item := range r.s.attempts {
		if item.TenantID == tenantID && item.CIOTID == ciotID {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r ProviderAttemptRepository) ExistsIdempotencyKey(_ context.Context, tenantID string, idempotencyKey string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	for _, item := range r.s.attempts {
		if item.TenantID == tenantID && item.IdempotencyKey == idempotencyKey {
			return true, nil
		}
	}
	return false, nil
}

type ExternalReferenceRepository struct{ s *MemoryStore }

func (r ExternalReferenceRepository) Upsert(_ context.Context, item domain.CIOTExternalReference) (domain.CIOTExternalReference, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	for k, current := range r.s.externalReferences {
		if current.TenantID == item.TenantID && current.CIOTID == item.CIOTID {
			item.ID = current.ID
			item.CreatedAt = current.CreatedAt
			item.Version = current.Version + 1
			r.s.externalReferences[k] = item
			return item, nil
		}
	}
	r.s.externalReferences[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r ExternalReferenceRepository) FindByCIOT(_ context.Context, tenantID string, ciotID string) (domain.CIOTExternalReference, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	for _, item := range r.s.externalReferences {
		if item.TenantID == tenantID && item.CIOTID == ciotID {
			return item, nil
		}
	}
	return domain.CIOTExternalReference{}, application.ErrNotFound
}

type DocumentRepository struct{ s *MemoryStore }

func (r DocumentRepository) Create(_ context.Context, item domain.CIOTDocument) (domain.CIOTDocument, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.documents[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r DocumentRepository) ListByCIOT(_ context.Context, tenantID string, ciotID string, q ports.Query) (ports.Page[domain.CIOTDocument], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.CIOTDocument{}
	for _, item := range r.s.documents {
		if item.TenantID == tenantID && item.CIOTID == ciotID {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}

type ErrorRepository struct{ s *MemoryStore }

func (r ErrorRepository) Create(_ context.Context, item domain.CIOTError) (domain.CIOTError, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.errors[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r ErrorRepository) ListByCIOT(_ context.Context, tenantID string, ciotID string, q ports.Query) (ports.Page[domain.CIOTError], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.CIOTError{}
	for _, item := range r.s.errors {
		if item.TenantID == tenantID && item.CIOTID == ciotID {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}

func matchCIOT(item domain.CIOT, q ports.Query) bool {
	if status := q.Filters["status"]; status != "" && string(item.Status) != status {
		return false
	}
	if ciotType := q.Filters["type"]; ciotType != "" && string(item.Type) != ciotType {
		return false
	}
	return contains(item.CIOTNumber+" "+item.ContractReference+" "+item.ExternalProtocol+" "+item.OperationalPeriod, q.Search)
}

func key(tenantID, id string) string { return tenantID + ":" + id }

func contains(source, search string) bool {
	if search == "" {
		return true
	}
	return strings.Contains(strings.ToLower(source), strings.ToLower(search))
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
	return ports.Page[T]{Data: items[start:end], Page: q.Page, PerPage: q.PerPage, Total: total, TotalPages: totalPages}
}

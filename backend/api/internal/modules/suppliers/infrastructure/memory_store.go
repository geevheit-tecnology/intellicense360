package infrastructure

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/ports"
)

type MemoryStore struct {
	mu         sync.RWMutex
	suppliers  map[string]domain.Supplier
	categories map[string]domain.SupplierCategory
	types      map[string]domain.SupplierType
	contacts   map[string]domain.SupplierContact
	addresses  map[string]domain.SupplierAddress
	documents  map[string]domain.SupplierDocument
	contracts  map[string]domain.SupplierContract
	ratings    map[string]domain.SupplierRating
	audit      []AuditEvent
}

type AuditEvent struct {
	TenantID, ActorID, Action, ResourceID string
	CreatedAt                             time.Time
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		suppliers: map[string]domain.Supplier{}, categories: map[string]domain.SupplierCategory{}, types: map[string]domain.SupplierType{},
		contacts: map[string]domain.SupplierContact{}, addresses: map[string]domain.SupplierAddress{}, documents: map[string]domain.SupplierDocument{},
		contracts: map[string]domain.SupplierContract{}, ratings: map[string]domain.SupplierRating{},
	}
}

func (s *MemoryStore) Suppliers() SupplierRepository  { return SupplierRepository{s: s} }
func (s *MemoryStore) Categories() CategoryRepository { return CategoryRepository{s: s} }
func (s *MemoryStore) Types() TypeRepository          { return TypeRepository{s: s} }
func (s *MemoryStore) Contacts() ContactRepository    { return ContactRepository{s: s} }
func (s *MemoryStore) Addresses() AddressRepository   { return AddressRepository{s: s} }
func (s *MemoryStore) Documents() DocumentRepository  { return DocumentRepository{s: s} }
func (s *MemoryStore) Contracts() ContractRepository  { return ContractRepository{s: s} }
func (s *MemoryStore) Ratings() RatingRepository      { return RatingRepository{s: s} }
func (s *MemoryStore) Audit() AuditRecorder           { return AuditRecorder{s: s} }

type AuditRecorder struct{ s *MemoryStore }

func (r AuditRecorder) RecordSupplierEvent(_ context.Context, tenantID string, actorID string, action string, resourceID string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.audit = append(r.s.audit, AuditEvent{TenantID: tenantID, ActorID: actorID, Action: action, ResourceID: resourceID, CreatedAt: time.Now().UTC()})
	return nil
}

type SupplierRepository struct{ s *MemoryStore }

func (r SupplierRepository) Create(_ context.Context, item domain.Supplier) (domain.Supplier, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.suppliers[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r SupplierRepository) FindByID(_ context.Context, tenantID string, id string) (domain.Supplier, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.suppliers[key(tenantID, id)]
	if !ok || item.DeletedAt != nil {
		return domain.Supplier{}, application.ErrNotFound
	}
	return item, nil
}
func (r SupplierRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.Supplier], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Supplier{}
	for _, item := range r.s.suppliers {
		if item.TenantID == tenantID && item.DeletedAt == nil && matchesSupplier(item, q) {
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
func (r SupplierRepository) Update(_ context.Context, item domain.Supplier) (domain.Supplier, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(item.TenantID, item.ID)
	current, ok := r.s.suppliers[k]
	if !ok || current.DeletedAt != nil {
		return domain.Supplier{}, application.ErrNotFound
	}
	item.DeletedAt = current.DeletedAt
	r.s.suppliers[k] = item
	return item, nil
}
func (r SupplierRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, id)
	item, ok := r.s.suppliers[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	item.Version++
	r.s.suppliers[k] = item
	return nil
}
func (r SupplierRepository) ExistsCNPJ(_ context.Context, tenantID string, cnpj string, exceptID string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	for _, item := range r.s.suppliers {
		if item.TenantID == tenantID && item.ID != exceptID && item.DeletedAt == nil && normDigits(item.CNPJ) == cnpj {
			return true, nil
		}
	}
	return false, nil
}
func (r SupplierRepository) Exists(_ context.Context, tenantID string, id string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.suppliers[key(tenantID, id)]
	return ok && item.DeletedAt == nil, nil
}

type CategoryRepository struct{ s *MemoryStore }

func (r CategoryRepository) Create(_ context.Context, item domain.SupplierCategory) (domain.SupplierCategory, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.categories[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r CategoryRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.SupplierCategory], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.SupplierCategory{}
	for _, item := range r.s.categories {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name+" "+item.Code, q.Search) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r CategoryRepository) Update(_ context.Context, item domain.SupplierCategory) (domain.SupplierCategory, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	current, ok := r.s.categories[key(item.TenantID, item.ID)]
	if !ok || current.DeletedAt != nil {
		return domain.SupplierCategory{}, application.ErrNotFound
	}
	item.CreatedAt = current.CreatedAt
	item.DeletedAt = current.DeletedAt
	r.s.categories[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r CategoryRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, id)
	item, ok := r.s.categories[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	item.Version++
	r.s.categories[k] = item
	return nil
}
func (r CategoryRepository) Exists(_ context.Context, tenantID string, id string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.categories[key(tenantID, id)]
	return ok && item.DeletedAt == nil, nil
}

type TypeRepository struct{ s *MemoryStore }

func (r TypeRepository) Create(_ context.Context, item domain.SupplierType) (domain.SupplierType, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.types[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r TypeRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.SupplierType], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.SupplierType{}
	for _, item := range r.s.types {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name+" "+item.Code, q.Search) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r TypeRepository) Update(_ context.Context, item domain.SupplierType) (domain.SupplierType, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	current, ok := r.s.types[key(item.TenantID, item.ID)]
	if !ok || current.DeletedAt != nil {
		return domain.SupplierType{}, application.ErrNotFound
	}
	item.CreatedAt = current.CreatedAt
	item.DeletedAt = current.DeletedAt
	r.s.types[key(item.TenantID, item.ID)] = item
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

type ContactRepository struct{ s *MemoryStore }

func (r ContactRepository) Create(_ context.Context, item domain.SupplierContact) (domain.SupplierContact, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.contacts[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r ContactRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.SupplierContact], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.SupplierContact{}
	for _, item := range r.s.contacts {
		if item.TenantID == tenantID && item.DeletedAt == nil && matchesChild(item.SupplierID, item.Name, q) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r ContactRepository) Update(_ context.Context, item domain.SupplierContact) (domain.SupplierContact, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	current, ok := r.s.contacts[key(item.TenantID, item.ID)]
	if !ok || current.DeletedAt != nil {
		return domain.SupplierContact{}, application.ErrNotFound
	}
	item.CreatedAt = current.CreatedAt
	item.DeletedAt = current.DeletedAt
	r.s.contacts[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r ContactRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, id)
	item, ok := r.s.contacts[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	item.Version++
	r.s.contacts[k] = item
	return nil
}

type AddressRepository struct{ s *MemoryStore }

func (r AddressRepository) Create(_ context.Context, item domain.SupplierAddress) (domain.SupplierAddress, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.addresses[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r AddressRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.SupplierAddress], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.SupplierAddress{}
	for _, item := range r.s.addresses {
		if item.TenantID == tenantID && item.DeletedAt == nil && matchesChild(item.SupplierID, item.City+" "+item.Street, q) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r AddressRepository) Update(_ context.Context, item domain.SupplierAddress) (domain.SupplierAddress, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	current, ok := r.s.addresses[key(item.TenantID, item.ID)]
	if !ok || current.DeletedAt != nil {
		return domain.SupplierAddress{}, application.ErrNotFound
	}
	item.CreatedAt = current.CreatedAt
	item.DeletedAt = current.DeletedAt
	r.s.addresses[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r AddressRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, id)
	item, ok := r.s.addresses[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	item.Version++
	r.s.addresses[k] = item
	return nil
}

type DocumentRepository struct{ s *MemoryStore }

func (r DocumentRepository) Create(_ context.Context, item domain.SupplierDocument) (domain.SupplierDocument, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.documents[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r DocumentRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.SupplierDocument], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.SupplierDocument{}
	for _, item := range r.s.documents {
		if item.TenantID == tenantID && item.DeletedAt == nil && matchesChild(item.SupplierID, item.DocumentType+" "+item.DocumentNumber, q) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r DocumentRepository) Update(_ context.Context, item domain.SupplierDocument) (domain.SupplierDocument, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	current, ok := r.s.documents[key(item.TenantID, item.ID)]
	if !ok || current.DeletedAt != nil {
		return domain.SupplierDocument{}, application.ErrNotFound
	}
	item.CreatedAt = current.CreatedAt
	item.DeletedAt = current.DeletedAt
	r.s.documents[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r DocumentRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, id)
	item, ok := r.s.documents[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	item.Version++
	r.s.documents[k] = item
	return nil
}

type ContractRepository struct{ s *MemoryStore }

func (r ContractRepository) Create(_ context.Context, item domain.SupplierContract) (domain.SupplierContract, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.contracts[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r ContractRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.SupplierContract], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.SupplierContract{}
	for _, item := range r.s.contracts {
		if item.TenantID == tenantID && item.DeletedAt == nil && matchesChild(item.SupplierID, item.ContractNumber, q) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r ContractRepository) Update(_ context.Context, item domain.SupplierContract) (domain.SupplierContract, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	current, ok := r.s.contracts[key(item.TenantID, item.ID)]
	if !ok || current.DeletedAt != nil {
		return domain.SupplierContract{}, application.ErrNotFound
	}
	item.CreatedAt = current.CreatedAt
	item.DeletedAt = current.DeletedAt
	r.s.contracts[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r ContractRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, id)
	item, ok := r.s.contracts[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	item.Version++
	r.s.contracts[k] = item
	return nil
}

type RatingRepository struct{ s *MemoryStore }

func (r RatingRepository) Create(_ context.Context, item domain.SupplierRating) (domain.SupplierRating, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.ratings[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r RatingRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.SupplierRating], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.SupplierRating{}
	for _, item := range r.s.ratings {
		if item.TenantID == tenantID && item.DeletedAt == nil && matchesChild(item.SupplierID, "", q) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r RatingRepository) Update(_ context.Context, item domain.SupplierRating) (domain.SupplierRating, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	current, ok := r.s.ratings[key(item.TenantID, item.ID)]
	if !ok || current.DeletedAt != nil {
		return domain.SupplierRating{}, application.ErrNotFound
	}
	item.CreatedAt = current.CreatedAt
	item.DeletedAt = current.DeletedAt
	r.s.ratings[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r RatingRepository) Delete(_ context.Context, tenantID string, id string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(tenantID, id)
	item, ok := r.s.ratings[k]
	if !ok || item.DeletedAt != nil {
		return application.ErrNotFound
	}
	now := time.Now().UTC()
	item.DeletedAt = &now
	item.UpdatedAt = now
	item.Version++
	r.s.ratings[k] = item
	return nil
}

func matchesSupplier(item domain.Supplier, q ports.Query) bool {
	if !contains(item.LegalName+" "+item.TradeName+" "+item.CNPJ, q.Search) {
		return false
	}
	if q.Filters != nil {
		if v := q.Filters["status"]; v != "" && string(item.Status) != v {
			return false
		}
		if v := q.Filters["type"]; v != "" && string(item.Type) != v {
			return false
		}
		if v := q.Filters["category_id"]; v != "" && item.CategoryID != v {
			return false
		}
	}
	return true
}
func matchesChild(supplierID string, text string, q ports.Query) bool {
	if !contains(text, q.Search) {
		return false
	}
	return q.Filters == nil || q.Filters["supplier_id"] == "" || supplierID == q.Filters["supplier_id"]
}
func key(tenantID, id string) string { return tenantID + ":" + id }
func norm(value string) string       { return strings.ToLower(strings.TrimSpace(value)) }
func normDigits(value string) string {
	b := strings.Builder{}
	for _, r := range value {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
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

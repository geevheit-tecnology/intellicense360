package infrastructure

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/financial/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/financial/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/financial/ports"
)

type MemoryStore struct {
	mu           sync.RWMutex
	transactions map[string]domain.FinancialTransaction
	categories   map[string]domain.CostCategory
	types        map[string]domain.CostType
	centers      map[string]domain.CostCenter
	accounts     map[string]domain.Account
	methods      map[string]domain.PaymentMethod
	periods      map[string]domain.FinancialPeriod
	budgets      map[string]domain.Budget
	adjustments  map[string]domain.FinancialAdjustment
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{transactions: map[string]domain.FinancialTransaction{}, categories: map[string]domain.CostCategory{}, types: map[string]domain.CostType{}, centers: map[string]domain.CostCenter{}, accounts: map[string]domain.Account{}, methods: map[string]domain.PaymentMethod{}, periods: map[string]domain.FinancialPeriod{}, budgets: map[string]domain.Budget{}, adjustments: map[string]domain.FinancialAdjustment{}}
}
func (s *MemoryStore) Transactions() TransactionRepository { return TransactionRepository{s: s} }
func (s *MemoryStore) Categories() CategoryRepository      { return CategoryRepository{s: s} }
func (s *MemoryStore) Types() TypeRepository               { return TypeRepository{s: s} }
func (s *MemoryStore) Centers() CenterRepository           { return CenterRepository{s: s} }
func (s *MemoryStore) Accounts() AccountRepository         { return AccountRepository{s: s} }
func (s *MemoryStore) Methods() MethodRepository           { return MethodRepository{s: s} }
func (s *MemoryStore) Periods() PeriodRepository           { return PeriodRepository{s: s} }
func (s *MemoryStore) Budgets() BudgetRepository           { return BudgetRepository{s: s} }
func (s *MemoryStore) Adjustments() AdjustmentRepository   { return AdjustmentRepository{s: s} }

type TransactionRepository struct{ s *MemoryStore }

func (r TransactionRepository) Create(_ context.Context, item domain.FinancialTransaction) (domain.FinancialTransaction, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.transactions[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r TransactionRepository) FindByID(_ context.Context, tenantID, id string) (domain.FinancialTransaction, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.transactions[key(tenantID, id)]
	if !ok || item.DeletedAt != nil {
		return domain.FinancialTransaction{}, application.ErrNotFound
	}
	return item, nil
}
func (r TransactionRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.FinancialTransaction], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.FinancialTransaction{}
	for _, item := range r.s.transactions {
		if item.TenantID == tenantID && item.DeletedAt == nil && matchTx(item, q) {
			items = append(items, item)
		}
	}
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.Before(items[j].CreatedAt) })
	return page(items, q), nil
}
func (r TransactionRepository) Update(_ context.Context, item domain.FinancialTransaction) (domain.FinancialTransaction, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(item.TenantID, item.ID)
	cur, ok := r.s.transactions[k]
	if !ok || cur.DeletedAt != nil {
		return domain.FinancialTransaction{}, application.ErrNotFound
	}
	item.CreatedAt, item.DeletedAt = cur.CreatedAt, cur.DeletedAt
	r.s.transactions[k] = item
	return item, nil
}
func (r TransactionRepository) Delete(_ context.Context, tenantID, id string) error {
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

type CategoryRepository struct{ s *MemoryStore }
type TypeRepository struct{ s *MemoryStore }
type CenterRepository struct{ s *MemoryStore }
type AccountRepository struct{ s *MemoryStore }
type MethodRepository struct{ s *MemoryStore }

func (r CategoryRepository) Create(_ context.Context, item domain.CostCategory) (domain.CostCategory, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.categories[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r CategoryRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.CostCategory], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.CostCategory{}
	for _, item := range r.s.categories {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name+" "+item.Code, q.Search) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r TypeRepository) Create(_ context.Context, item domain.CostType) (domain.CostType, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.types[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r TypeRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.CostType], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.CostType{}
	for _, item := range r.s.types {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name+" "+item.Code, q.Search) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r CenterRepository) Create(_ context.Context, item domain.CostCenter) (domain.CostCenter, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.centers[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r CenterRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.CostCenter], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.CostCenter{}
	for _, item := range r.s.centers {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name+" "+item.Code, q.Search) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r AccountRepository) Create(_ context.Context, item domain.Account) (domain.Account, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.accounts[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r AccountRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.Account], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Account{}
	for _, item := range r.s.accounts {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name+" "+item.AccountCode, q.Search) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r MethodRepository) Create(_ context.Context, item domain.PaymentMethod) (domain.PaymentMethod, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.methods[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r MethodRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.PaymentMethod], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.PaymentMethod{}
	for _, item := range r.s.methods {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name+" "+item.Code, q.Search) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}

type PeriodRepository struct{ s *MemoryStore }

func (r PeriodRepository) Create(_ context.Context, item domain.FinancialPeriod) (domain.FinancialPeriod, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.periods[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r PeriodRepository) FindByID(_ context.Context, tenantID, id string) (domain.FinancialPeriod, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	item, ok := r.s.periods[key(tenantID, id)]
	if !ok || item.DeletedAt != nil {
		return domain.FinancialPeriod{}, application.ErrNotFound
	}
	return item, nil
}
func (r PeriodRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.FinancialPeriod], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.FinancialPeriod{}
	for _, item := range r.s.periods {
		if item.TenantID == tenantID && item.DeletedAt == nil {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}
func (r PeriodRepository) Update(_ context.Context, item domain.FinancialPeriod) (domain.FinancialPeriod, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	k := key(item.TenantID, item.ID)
	cur, ok := r.s.periods[k]
	if !ok || cur.DeletedAt != nil {
		return domain.FinancialPeriod{}, application.ErrNotFound
	}
	item.CreatedAt, item.DeletedAt = cur.CreatedAt, cur.DeletedAt
	r.s.periods[k] = item
	return item, nil
}
func (r PeriodRepository) FindClosedForDate(_ context.Context, tenantID, date string) (bool, error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false, err
	}
	for _, p := range r.s.periods {
		if p.TenantID == tenantID && p.DeletedAt == nil && p.Status == domain.StatusClosed && !d.Before(p.StartDate) && !d.After(p.EndDate) {
			return true, nil
		}
	}
	return false, nil
}

type BudgetRepository struct{ s *MemoryStore }

func (r BudgetRepository) Create(_ context.Context, item domain.Budget) (domain.Budget, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.budgets[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r BudgetRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.Budget], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.Budget{}
	for _, item := range r.s.budgets {
		if item.TenantID == tenantID && item.DeletedAt == nil && contains(item.Name, q.Search) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}

type AdjustmentRepository struct{ s *MemoryStore }

func (r AdjustmentRepository) Create(_ context.Context, item domain.FinancialAdjustment) (domain.FinancialAdjustment, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	r.s.adjustments[key(item.TenantID, item.ID)] = item
	return item, nil
}
func (r AdjustmentRepository) Search(_ context.Context, tenantID string, q ports.Query) (ports.Page[domain.FinancialAdjustment], error) {
	r.s.mu.RLock()
	defer r.s.mu.RUnlock()
	items := []domain.FinancialAdjustment{}
	for _, item := range r.s.adjustments {
		if item.TenantID == tenantID && contains(item.Reason+" "+item.TransactionID, q.Search) {
			items = append(items, item)
		}
	}
	return page(items, q), nil
}

func matchTx(item domain.FinancialTransaction, q ports.Query) bool {
	if q.Filters["kind"] != "" && string(item.Kind) != q.Filters["kind"] {
		return false
	}
	if q.Filters["status"] != "" && string(item.Status) != q.Filters["status"] {
		return false
	}
	return contains(item.Description+" "+item.DocumentNumber+" "+item.SupplierReference, q.Search)
}
func key(tenantID, id string) string { return tenantID + ":" + id }
func contains(value, search string) bool {
	search = strings.ToLower(strings.TrimSpace(search))
	return search == "" || strings.Contains(strings.ToLower(value), search)
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
	pages := 0
	if total > 0 {
		pages = (total + q.PageSize - 1) / q.PageSize
	}
	return ports.Page[T]{Items: items[start:end], Page: q.Page, PageSize: q.PageSize, TotalItems: total, TotalPages: pages}
}

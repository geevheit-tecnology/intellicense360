package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/financial/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/financial/ports"
)

type TransactionService struct {
	tx          ports.TransactionRepository
	periods     ports.PeriodRepository
	adjustments ports.AdjustmentRepository
}

func NewTransactionService(tx ports.TransactionRepository, periods ports.PeriodRepository, adjustments ports.AdjustmentRepository) TransactionService {
	return TransactionService{tx: tx, periods: periods, adjustments: adjustments}
}
func (s TransactionService) CreateExpense(ctx context.Context, item domain.FinancialTransaction) (domain.FinancialTransaction, error) {
	item.Kind = domain.KindExpense
	return s.create(ctx, item)
}
func (s TransactionService) CreateRevenue(ctx context.Context, item domain.FinancialTransaction) (domain.FinancialTransaction, error) {
	item.Kind = domain.KindRevenue
	return s.create(ctx, item)
}
func (s TransactionService) create(ctx context.Context, item domain.FinancialTransaction) (domain.FinancialTransaction, error) {
	if err := s.validate(ctx, item); err != nil {
		return domain.FinancialTransaction{}, err
	}
	now := time.Now().UTC()
	item.ID = newID("fin")
	if item.Date.IsZero() {
		item.Date = now
	}
	if item.Status == "" {
		item.Status = domain.StatusDraft
	}
	item.CreatedAt, item.UpdatedAt, item.Version = now, now, 1
	return s.tx.Create(ctx, item)
}
func (s TransactionService) Update(ctx context.Context, item domain.FinancialTransaction) (domain.FinancialTransaction, error) {
	current, err := s.tx.FindByID(ctx, item.TenantID, item.ID)
	if err != nil {
		return domain.FinancialTransaction{}, err
	}
	if isFinal(current.Status) {
		return domain.FinancialTransaction{}, ErrFinalizedImmutable
	}
	if err := s.validate(ctx, item); err != nil {
		return domain.FinancialTransaction{}, err
	}
	item.Kind, item.Status = current.Kind, current.Status
	item.CreatedAt, item.UpdatedAt, item.Version = current.CreatedAt, time.Now().UTC(), current.Version+1
	return s.tx.Update(ctx, item)
}
func (s TransactionService) Approve(ctx context.Context, tenantID string, id string, actorID string) (domain.FinancialTransaction, error) {
	return s.transition(ctx, tenantID, id, actorID, domain.StatusApproved)
}
func (s TransactionService) Pay(ctx context.Context, tenantID string, id string, actorID string) (domain.FinancialTransaction, error) {
	return s.transition(ctx, tenantID, id, actorID, domain.StatusPaid)
}
func (s TransactionService) Receive(ctx context.Context, tenantID string, id string, actorID string) (domain.FinancialTransaction, error) {
	return s.transition(ctx, tenantID, id, actorID, domain.StatusReceived)
}
func (s TransactionService) Cancel(ctx context.Context, tenantID string, id string, actorID string) (domain.FinancialTransaction, error) {
	return s.transition(ctx, tenantID, id, actorID, domain.StatusCanceled)
}
func (s TransactionService) transition(ctx context.Context, tenantID string, id string, actorID string, to domain.FinancialStatus) (domain.FinancialTransaction, error) {
	item, err := s.tx.FindByID(ctx, tenantID, id)
	if err != nil {
		return domain.FinancialTransaction{}, err
	}
	if err := s.checkOpenPeriod(ctx, item); err != nil {
		return domain.FinancialTransaction{}, err
	}
	if !canTransition(item.Kind, item.Status, to) {
		return domain.FinancialTransaction{}, ErrInvalidTransition
	}
	now := time.Now().UTC()
	item.Status, item.UpdatedBy, item.UpdatedAt, item.Version = to, actorID, now, item.Version+1
	if to == domain.StatusPaid || to == domain.StatusReceived {
		item.SettlementDate = &now
	}
	return s.tx.Update(ctx, item)
}
func (s TransactionService) Adjust(ctx context.Context, item domain.FinancialAdjustment) (domain.FinancialAdjustment, error) {
	tx, err := s.tx.FindByID(ctx, item.TenantID, item.TransactionID)
	if err != nil {
		return domain.FinancialAdjustment{}, err
	}
	if tx.Status != domain.StatusPaid && tx.Status != domain.StatusReceived {
		return domain.FinancialAdjustment{}, ErrInvalidTransition
	}
	if item.AdjustedAmount <= 0 || strings.TrimSpace(item.Reason) == "" {
		return domain.FinancialAdjustment{}, ErrInvalidFinancialData
	}
	now := time.Now().UTC()
	item.ID, item.CreatedAt = newID("faj"), now
	saved, err := s.adjustments.Create(ctx, item)
	if err != nil {
		return domain.FinancialAdjustment{}, err
	}
	tx.Status, tx.UpdatedBy, tx.UpdatedAt, tx.Version = domain.StatusAdjusted, item.CreatedBy, now, tx.Version+1
	_, err = s.tx.Update(ctx, tx)
	return saved, err
}
func (s TransactionService) FindByID(ctx context.Context, tenantID string, id string) (domain.FinancialTransaction, error) {
	return s.tx.FindByID(ctx, tenantID, id)
}
func (s TransactionService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FinancialTransaction], error) {
	return s.tx.Search(ctx, tenantID, normalizeQuery(q))
}
func (s TransactionService) validate(ctx context.Context, item domain.FinancialTransaction) error {
	if item.TenantID == "" || strings.TrimSpace(item.Description) == "" || item.Amount <= 0 {
		return ErrInvalidFinancialData
	}
	if item.Date.IsZero() {
		item.Date = time.Now().UTC()
	}
	return s.checkOpenPeriod(ctx, item)
}
func (s TransactionService) checkOpenPeriod(ctx context.Context, item domain.FinancialTransaction) error {
	if s.periods == nil {
		return nil
	}
	closed, err := s.periods.FindClosedForDate(ctx, item.TenantID, item.Date.Format("2006-01-02"))
	if err != nil {
		return err
	}
	if closed {
		return ErrClosedPeriodImmutable
	}
	return nil
}

type PeriodService struct{ repo ports.PeriodRepository }

func NewPeriodService(repo ports.PeriodRepository) PeriodService { return PeriodService{repo: repo} }
func (s PeriodService) Create(ctx context.Context, item domain.FinancialPeriod) (domain.FinancialPeriod, error) {
	if item.TenantID == "" || item.Year <= 0 || item.Month < 1 || item.Month > 12 {
		return domain.FinancialPeriod{}, ErrInvalidFinancialData
	}
	now := time.Now().UTC()
	item.ID, item.Status, item.CreatedAt, item.UpdatedAt, item.Version = newID("fpe"), domain.StatusPending, now, now, 1
	return s.repo.Create(ctx, item)
}
func (s PeriodService) Close(ctx context.Context, tenantID string, id string) (domain.FinancialPeriod, error) {
	item, err := s.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return domain.FinancialPeriod{}, err
	}
	item.Status, item.UpdatedAt, item.Version = domain.StatusClosed, time.Now().UTC(), item.Version+1
	return s.repo.Update(ctx, item)
}
func (s PeriodService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FinancialPeriod], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(q))
}

type BudgetService struct{ repo ports.BudgetRepository }

func NewBudgetService(repo ports.BudgetRepository) BudgetService { return BudgetService{repo: repo} }
func (s BudgetService) Create(ctx context.Context, item domain.Budget) (domain.Budget, error) {
	if item.TenantID == "" || item.Name == "" || item.PlannedAmount < 0 {
		return domain.Budget{}, ErrInvalidFinancialData
	}
	now := time.Now().UTC()
	item.ID, item.Status, item.CreatedAt, item.UpdatedAt, item.Version = newID("fbu"), "draft", now, now, 1
	return s.repo.Create(ctx, item)
}
func (s BudgetService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Budget], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(q))
}

type AdjustmentService struct{ repo ports.AdjustmentRepository }

func NewAdjustmentService(repo ports.AdjustmentRepository) AdjustmentService {
	return AdjustmentService{repo: repo}
}
func (s AdjustmentService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FinancialAdjustment], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(q))
}

type CatalogService[T any] struct {
	repo ports.CatalogRepository[T]
	init func(*T) error
}

func NewCatalogService[T any](repo ports.CatalogRepository[T], init func(*T) error) CatalogService[T] {
	return CatalogService[T]{repo: repo, init: init}
}
func (s CatalogService[T]) Create(ctx context.Context, item T) (T, error) {
	if s.init != nil {
		if err := s.init(&item); err != nil {
			var zero T
			return zero, err
		}
	}
	return s.repo.Create(ctx, item)
}
func (s CatalogService[T]) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[T], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(q))
}

func InitCategory(item *domain.CostCategory) error {
	if item.TenantID == "" || item.Name == "" {
		return ErrInvalidFinancialData
	}
	now := time.Now().UTC()
	item.ID, item.Active, item.CreatedAt, item.UpdatedAt, item.Version = newID("fcc"), true, now, now, 1
	return nil
}
func InitType(item *domain.CostType) error {
	if item.TenantID == "" || item.Name == "" {
		return ErrInvalidFinancialData
	}
	now := time.Now().UTC()
	item.ID, item.Active, item.CreatedAt, item.UpdatedAt, item.Version = newID("fct"), true, now, now, 1
	return nil
}
func InitCenter(item *domain.CostCenter) error {
	if item.TenantID == "" || item.Name == "" {
		return ErrInvalidFinancialData
	}
	now := time.Now().UTC()
	item.ID, item.Active, item.CreatedAt, item.UpdatedAt, item.Version = newID("fce"), true, now, now, 1
	return nil
}
func InitAccount(item *domain.Account) error {
	if item.TenantID == "" || item.Name == "" {
		return ErrInvalidFinancialData
	}
	now := time.Now().UTC()
	item.ID, item.Status, item.CreatedAt, item.UpdatedAt, item.Version = newID("fac"), "active", now, now, 1
	return nil
}
func InitPaymentMethod(item *domain.PaymentMethod) error {
	if item.TenantID == "" || item.Name == "" {
		return ErrInvalidFinancialData
	}
	now := time.Now().UTC()
	item.ID, item.Active, item.CreatedAt, item.UpdatedAt, item.Version = newID("fpm"), true, now, now, 1
	return nil
}

func canTransition(kind domain.TransactionKind, from domain.FinancialStatus, to domain.FinancialStatus) bool {
	if from == domain.StatusDraft && to == domain.StatusPending {
		return true
	}
	if from == domain.StatusDraft && to == domain.StatusApproved {
		return true
	}
	if from == domain.StatusPending && (to == domain.StatusApproved || to == domain.StatusCanceled) {
		return true
	}
	if from == domain.StatusApproved && (to == domain.StatusCanceled || (kind == domain.KindExpense && to == domain.StatusPaid) || (kind == domain.KindRevenue && to == domain.StatusReceived)) {
		return true
	}
	return false
}
func isFinal(status domain.FinancialStatus) bool {
	return status == domain.StatusPaid || status == domain.StatusReceived || status == domain.StatusAdjusted || status == domain.StatusCanceled
}
func normalizeQuery(q ports.Query) ports.Query {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 20
	}
	return q
}
func newID(prefix string) string {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return prefix + "_" + time.Now().UTC().Format("20060102150405")
	}
	return prefix + "_" + hex.EncodeToString(b[:])
}

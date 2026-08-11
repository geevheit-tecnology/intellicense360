package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"math"
	"strings"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/ports"
)

type FuelTransactionService struct {
	transactions ports.FuelTransactionRepository
	adjustments  ports.FuelAdjustmentRepository
}

func NewFuelTransactionService(transactions ports.FuelTransactionRepository, adjustments ports.FuelAdjustmentRepository) FuelTransactionService {
	return FuelTransactionService{transactions: transactions, adjustments: adjustments}
}

func (s FuelTransactionService) Create(ctx context.Context, item domain.FuelTransaction) (domain.FuelTransaction, error) {
	if err := validateTransaction(item); err != nil {
		return domain.FuelTransaction{}, err
	}
	now := time.Now().UTC()
	item.ID = newID("ftr")
	if item.TransactionDate.IsZero() {
		item.TransactionDate = now
	}
	if item.Status == "" {
		item.Status = domain.FuelTransactionDraft
	}
	if !validTransactionStatus(item.Status) || item.Status == domain.FuelTransactionAdjusted {
		return domain.FuelTransaction{}, ErrInvalidStatus
	}
	if item.Status == domain.FuelTransactionCompleted {
		item.CompletedAt = &now
	}
	item.TotalAmount = normalizeMoney(item.TotalAmount)
	item.CreatedAt = now
	item.UpdatedAt = now
	item.Version = 1
	return s.transactions.Create(ctx, item)
}

func (s FuelTransactionService) Update(ctx context.Context, item domain.FuelTransaction) (domain.FuelTransaction, error) {
	if item.TenantID == "" {
		return domain.FuelTransaction{}, ErrTenantIDRequired
	}
	current, err := s.transactions.FindByID(ctx, item.TenantID, item.ID)
	if err != nil {
		return domain.FuelTransaction{}, err
	}
	if isFinalTransaction(current.Status) {
		return domain.FuelTransaction{}, ErrCompletedImmutable
	}
	if item.Status == "" {
		item.Status = current.Status
	}
	if item.Status != current.Status {
		return domain.FuelTransaction{}, ErrInvalidTransition
	}
	if err := validateTransaction(item); err != nil {
		return domain.FuelTransaction{}, err
	}
	item.CreatedAt = current.CreatedAt
	item.UpdatedAt = time.Now().UTC()
	item.Version = current.Version + 1
	return s.transactions.Update(ctx, item)
}

func (s FuelTransactionService) Complete(ctx context.Context, tenantID string, id string, actorID string) (domain.FuelTransaction, error) {
	item, err := s.transactions.FindByID(ctx, tenantID, id)
	if err != nil {
		return domain.FuelTransaction{}, err
	}
	if item.Status != domain.FuelTransactionDraft {
		return domain.FuelTransaction{}, ErrInvalidTransition
	}
	now := time.Now().UTC()
	item.Status = domain.FuelTransactionCompleted
	item.CompletedAt = &now
	item.UpdatedBy = actorID
	item.UpdatedAt = now
	item.Version++
	return s.transactions.Update(ctx, item)
}

func (s FuelTransactionService) Cancel(ctx context.Context, tenantID string, id string, reason string, actorID string) (domain.FuelTransaction, error) {
	if strings.TrimSpace(reason) == "" {
		return domain.FuelTransaction{}, ErrReasonRequired
	}
	item, err := s.transactions.FindByID(ctx, tenantID, id)
	if err != nil {
		return domain.FuelTransaction{}, err
	}
	if item.Status != domain.FuelTransactionDraft && item.Status != domain.FuelTransactionCompleted {
		return domain.FuelTransaction{}, ErrInvalidTransition
	}
	now := time.Now().UTC()
	item.Status = domain.FuelTransactionCanceled
	item.CancellationReason = reason
	item.CanceledAt = &now
	item.UpdatedBy = actorID
	item.UpdatedAt = now
	item.Version++
	return s.transactions.Update(ctx, item)
}

func (s FuelTransactionService) Adjust(ctx context.Context, adjustment domain.FuelAdjustment) (domain.FuelAdjustment, error) {
	if adjustment.TenantID == "" {
		return domain.FuelAdjustment{}, ErrTenantIDRequired
	}
	if adjustment.TransactionID == "" {
		return domain.FuelAdjustment{}, ErrTransactionRequired
	}
	if strings.TrimSpace(adjustment.Reason) == "" || strings.TrimSpace(adjustment.AdjustmentType) == "" {
		return domain.FuelAdjustment{}, ErrReasonRequired
	}
	transaction, err := s.transactions.FindByID(ctx, adjustment.TenantID, adjustment.TransactionID)
	if err != nil {
		return domain.FuelAdjustment{}, err
	}
	if transaction.Status != domain.FuelTransactionCompleted {
		return domain.FuelAdjustment{}, ErrInvalidTransition
	}
	if adjustment.AdjustedQuantity < 0 || adjustment.AdjustedUnitPrice < 0 || adjustment.AdjustedTotalAmount < 0 {
		return domain.FuelAdjustment{}, ErrInvalidFuelData
	}
	now := time.Now().UTC()
	adjustment.ID = newID("faj")
	adjustment.CreatedAt = now
	saved, err := s.adjustments.Create(ctx, adjustment)
	if err != nil {
		return domain.FuelAdjustment{}, err
	}
	transaction.Status = domain.FuelTransactionAdjusted
	transaction.UpdatedBy = adjustment.CreatedBy
	transaction.UpdatedAt = now
	transaction.Version++
	_, err = s.transactions.Update(ctx, transaction)
	return saved, err
}

func (s FuelTransactionService) Delete(ctx context.Context, tenantID string, id string) error {
	item, err := s.transactions.FindByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if isFinalTransaction(item.Status) {
		return ErrFuelTransactionFinal
	}
	return s.transactions.Delete(ctx, tenantID, id)
}
func (s FuelTransactionService) FindByID(ctx context.Context, tenantID string, id string) (domain.FuelTransaction, error) {
	return s.transactions.FindByID(ctx, tenantID, id)
}
func (s FuelTransactionService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelTransaction], error) {
	return s.transactions.Search(ctx, tenantID, normalizeQuery(q))
}

type FuelTypeService struct{ repo ports.FuelTypeRepository }

func NewFuelTypeService(repo ports.FuelTypeRepository) FuelTypeService {
	return FuelTypeService{repo: repo}
}
func (s FuelTypeService) Create(ctx context.Context, item domain.FuelType) (domain.FuelType, error) {
	if item.TenantID == "" || strings.TrimSpace(item.Name) == "" || !validFuelKind(item.Kind) {
		return domain.FuelType{}, ErrInvalidFuelData
	}
	now := time.Now().UTC()
	item.ID = newID("fty")
	item.Active = true
	item.CreatedAt, item.UpdatedAt, item.Version = now, now, 1
	return s.repo.Create(ctx, item)
}
func (s FuelTypeService) Update(ctx context.Context, item domain.FuelType) (domain.FuelType, error) {
	current, err := s.repo.FindByID(ctx, item.TenantID, item.ID)
	if err != nil {
		return domain.FuelType{}, err
	}
	if strings.TrimSpace(item.Name) == "" || !validFuelKind(item.Kind) {
		return domain.FuelType{}, ErrInvalidFuelData
	}
	item.CreatedAt = current.CreatedAt
	item.UpdatedAt = time.Now().UTC()
	item.Version = current.Version + 1
	return s.repo.Update(ctx, item)
}
func (s FuelTypeService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}
func (s FuelTypeService) FindByID(ctx context.Context, tenantID string, id string) (domain.FuelType, error) {
	return s.repo.FindByID(ctx, tenantID, id)
}
func (s FuelTypeService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelType], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(q))
}

type FuelStationService struct{ repo ports.FuelStationRepository }

func NewFuelStationService(repo ports.FuelStationRepository) FuelStationService {
	return FuelStationService{repo: repo}
}
func (s FuelStationService) Create(ctx context.Context, item domain.FuelStation) (domain.FuelStation, error) {
	if item.TenantID == "" || strings.TrimSpace(item.Name) == "" {
		return domain.FuelStation{}, ErrInvalidFuelData
	}
	if cnpj := normalizeDigits(item.CNPJ); cnpj != "" {
		exists, err := s.repo.ExistsCNPJ(ctx, item.TenantID, cnpj, "")
		if err != nil || exists {
			if exists {
				return domain.FuelStation{}, ErrDuplicateCNPJ
			}
			return domain.FuelStation{}, err
		}
		item.CNPJ = cnpj
	}
	now := time.Now().UTC()
	item.ID, item.Active = newID("fst"), true
	item.CreatedAt, item.UpdatedAt, item.Version = now, now, 1
	return s.repo.Create(ctx, item)
}
func (s FuelStationService) Update(ctx context.Context, item domain.FuelStation) (domain.FuelStation, error) {
	current, err := s.repo.FindByID(ctx, item.TenantID, item.ID)
	if err != nil {
		return domain.FuelStation{}, err
	}
	if strings.TrimSpace(item.Name) == "" {
		return domain.FuelStation{}, ErrInvalidFuelData
	}
	item.CNPJ = normalizeDigits(item.CNPJ)
	if item.CNPJ != "" {
		exists, err := s.repo.ExistsCNPJ(ctx, item.TenantID, item.CNPJ, item.ID)
		if err != nil || exists {
			if exists {
				return domain.FuelStation{}, ErrDuplicateCNPJ
			}
			return domain.FuelStation{}, err
		}
	}
	item.CreatedAt = current.CreatedAt
	item.UpdatedAt = time.Now().UTC()
	item.Version = current.Version + 1
	return s.repo.Update(ctx, item)
}
func (s FuelStationService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}
func (s FuelStationService) FindByID(ctx context.Context, tenantID string, id string) (domain.FuelStation, error) {
	return s.repo.FindByID(ctx, tenantID, id)
}
func (s FuelStationService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelStation], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(q))
}

type FuelTankService struct{ repo ports.FuelTankRepository }
type FuelNozzleService struct{ repo ports.FuelNozzleRepository }

func NewFuelTankService(repo ports.FuelTankRepository) FuelTankService {
	return FuelTankService{repo: repo}
}
func NewFuelNozzleService(repo ports.FuelNozzleRepository) FuelNozzleService {
	return FuelNozzleService{repo: repo}
}

func (s FuelTankService) Create(ctx context.Context, item domain.FuelTank) (domain.FuelTank, error) {
	if item.TenantID == "" || item.Code == "" || item.Capacity < 0 || item.CurrentReading < 0 || !validFuelKind(item.FuelKind) {
		return domain.FuelTank{}, ErrInvalidFuelData
	}
	now := time.Now().UTC()
	item.ID, item.Status, item.CreatedAt, item.UpdatedAt, item.Version = newID("ftk"), defaultTankStatus(item.Status), now, now, 1
	return s.repo.Create(ctx, item)
}
func (s FuelTankService) Update(ctx context.Context, item domain.FuelTank) (domain.FuelTank, error) {
	if item.TenantID == "" || item.ID == "" || item.Code == "" || item.Capacity < 0 || item.CurrentReading < 0 || !validFuelKind(item.FuelKind) {
		return domain.FuelTank{}, ErrInvalidFuelData
	}
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}
func (s FuelTankService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}
func (s FuelTankService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelTank], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(q))
}

func (s FuelNozzleService) Create(ctx context.Context, item domain.FuelNozzle) (domain.FuelNozzle, error) {
	if item.TenantID == "" || item.Code == "" || item.MeterReading < 0 || !validFuelKind(item.FuelKind) {
		return domain.FuelNozzle{}, ErrInvalidFuelData
	}
	now := time.Now().UTC()
	item.ID, item.Status, item.CreatedAt, item.UpdatedAt, item.Version = newID("fnz"), defaultNozzleStatus(item.Status), now, now, 1
	return s.repo.Create(ctx, item)
}
func (s FuelNozzleService) Update(ctx context.Context, item domain.FuelNozzle) (domain.FuelNozzle, error) {
	if item.TenantID == "" || item.ID == "" || item.Code == "" || item.MeterReading < 0 || !validFuelKind(item.FuelKind) {
		return domain.FuelNozzle{}, ErrInvalidFuelData
	}
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}
func (s FuelNozzleService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}
func (s FuelNozzleService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelNozzle], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(q))
}

type FuelReadingService struct{ repo ports.FuelReadingRepository }
type FuelPriceService struct{ repo ports.FuelPriceRepository }
type FuelReceiptService struct{ repo ports.FuelReceiptRepository }
type FuelAdjustmentService struct {
	repo ports.FuelAdjustmentRepository
}

func NewFuelReadingService(repo ports.FuelReadingRepository) FuelReadingService {
	return FuelReadingService{repo: repo}
}
func NewFuelPriceService(repo ports.FuelPriceRepository) FuelPriceService {
	return FuelPriceService{repo: repo}
}
func NewFuelReceiptService(repo ports.FuelReceiptRepository) FuelReceiptService {
	return FuelReceiptService{repo: repo}
}
func NewFuelAdjustmentService(repo ports.FuelAdjustmentRepository) FuelAdjustmentService {
	return FuelAdjustmentService{repo: repo}
}

func (s FuelReadingService) Record(ctx context.Context, item domain.FuelReading) (domain.FuelReading, error) {
	if item.TenantID == "" || item.Value < 0 || !validReadingType(item.ReadingType) {
		return domain.FuelReading{}, ErrInvalidFuelData
	}
	now := time.Now().UTC()
	if item.ReadingDate.IsZero() {
		item.ReadingDate = now
	}
	item.ID, item.CreatedAt = newID("frd"), now
	return s.repo.Create(ctx, item)
}
func (s FuelReadingService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelReading], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(q))
}

func (s FuelPriceService) Record(ctx context.Context, item domain.FuelPrice) (domain.FuelPrice, error) {
	if item.TenantID == "" || item.UnitPrice < 0 || !validFuelKind(item.FuelKind) {
		return domain.FuelPrice{}, ErrInvalidFuelData
	}
	now := time.Now().UTC()
	if item.EffectiveDate.IsZero() {
		item.EffectiveDate = now
	}
	item.ID, item.CreatedAt, item.UpdatedAt, item.Version = newID("fpr"), now, now, 1
	return s.repo.Create(ctx, item)
}
func (s FuelPriceService) Update(ctx context.Context, item domain.FuelPrice) (domain.FuelPrice, error) {
	if item.TenantID == "" || item.ID == "" || item.UnitPrice < 0 || !validFuelKind(item.FuelKind) {
		return domain.FuelPrice{}, ErrInvalidFuelData
	}
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}
func (s FuelPriceService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}
func (s FuelPriceService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelPrice], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(q))
}

func (s FuelReceiptService) Create(ctx context.Context, item domain.FuelReceipt) (domain.FuelReceipt, error) {
	if item.TenantID == "" || item.ReceiptNumber == "" || item.Amount < 0 {
		return domain.FuelReceipt{}, ErrInvalidFuelData
	}
	now := time.Now().UTC()
	if item.ReceiptDate.IsZero() {
		item.ReceiptDate = now
	}
	item.ID, item.CreatedAt, item.UpdatedAt, item.Version = newID("frc"), now, now, 1
	return s.repo.Create(ctx, item)
}
func (s FuelReceiptService) Update(ctx context.Context, item domain.FuelReceipt) (domain.FuelReceipt, error) {
	if item.TenantID == "" || item.ID == "" || item.ReceiptNumber == "" || item.Amount < 0 {
		return domain.FuelReceipt{}, ErrInvalidFuelData
	}
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}
func (s FuelReceiptService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}
func (s FuelReceiptService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelReceipt], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(q))
}
func (s FuelAdjustmentService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelAdjustment], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(q))
}

func validateTransaction(item domain.FuelTransaction) error {
	if item.TenantID == "" {
		return ErrTenantIDRequired
	}
	if item.Quantity <= 0 || item.UnitPrice < 0 || item.TotalAmount < 0 || item.OdometerReading < 0 || item.EngineHourReading < 0 {
		return ErrInvalidFuelData
	}
	if !validFuelKind(item.FuelKind) {
		return ErrInvalidFuelKind
	}
	expected := normalizeMoney(item.Quantity * item.UnitPrice)
	if math.Abs(expected-normalizeMoney(item.TotalAmount)) > 0.01 {
		return ErrInvalidFuelData
	}
	return nil
}

func validFuelKind(kind domain.FuelKind) bool {
	switch kind {
	case domain.FuelKindDieselS10, domain.FuelKindDieselS500, domain.FuelKindGasoline, domain.FuelKindEthanol, domain.FuelKindGNV, domain.FuelKindARLA32, domain.FuelKindOther:
		return true
	default:
		return false
	}
}
func validTransactionStatus(status domain.FuelTransactionStatus) bool {
	switch status {
	case domain.FuelTransactionDraft, domain.FuelTransactionCompleted, domain.FuelTransactionCanceled, domain.FuelTransactionAdjusted, domain.FuelTransactionRejected:
		return true
	default:
		return false
	}
}
func validReadingType(kind domain.FuelReadingType) bool {
	return kind == domain.FuelReadingOdometer || kind == domain.FuelReadingEngineHour || kind == domain.FuelReadingTank
}
func isFinalTransaction(status domain.FuelTransactionStatus) bool {
	return status == domain.FuelTransactionCompleted || status == domain.FuelTransactionCanceled || status == domain.FuelTransactionAdjusted || status == domain.FuelTransactionRejected
}
func defaultTankStatus(status domain.FuelTankStatus) domain.FuelTankStatus {
	if status == "" {
		return domain.FuelTankActive
	}
	return status
}
func defaultNozzleStatus(status domain.FuelNozzleStatus) domain.FuelNozzleStatus {
	if status == "" {
		return domain.FuelNozzleActive
	}
	return status
}
func normalizeQuery(q ports.Query) ports.Query {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = 20
	}
	if q.PageSize > 100 {
		q.PageSize = 100
	}
	q.SortOrder = strings.ToLower(q.SortOrder)
	return q
}
func normalizeMoney(value float64) float64 { return math.Round(value*100) / 100 }
func normalizeDigits(value string) string {
	var b strings.Builder
	for _, r := range value {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
func newID(prefix string) string {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return prefix + "_" + time.Now().UTC().Format("20060102150405")
	}
	return prefix + "_" + hex.EncodeToString(b[:])
}

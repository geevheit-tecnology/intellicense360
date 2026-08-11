package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/ciot/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/ciot/ports"
)

type CIOTService struct {
	repo    ports.CIOTRepository
	history ports.StatusHistoryRepository
	errors  ports.ErrorRepository
}

func NewCIOTService(repo ports.CIOTRepository, history ports.StatusHistoryRepository, errors ports.ErrorRepository) CIOTService {
	return CIOTService{repo: repo, history: history, errors: errors}
}

func (s CIOTService) Create(ctx context.Context, item domain.CIOT) (domain.CIOT, error) {
	if err := validateCIOT(item); err != nil {
		return domain.CIOT{}, err
	}
	if strings.TrimSpace(item.IdempotencyKey) != "" {
		exists, err := s.repo.ExistsIdempotencyKey(ctx, item.TenantID, item.IdempotencyKey)
		if err != nil {
			return domain.CIOT{}, err
		}
		if exists {
			return domain.CIOT{}, ErrDuplicateRequest
		}
	}
	now := time.Now().UTC()
	item.ID = newID("cio")
	item.Status = domain.StatusDraft
	item.CreatedAt, item.UpdatedAt, item.Version = now, now, 1
	if item.StartDate.IsZero() {
		item.StartDate = now
	}
	saved, err := s.repo.Create(ctx, item)
	if err != nil {
		return domain.CIOT{}, err
	}
	_, err = s.appendHistory(ctx, saved, "", domain.StatusDraft, domain.EventCreated, item.CreatedBy, "")
	return saved, err
}

func (s CIOTService) Update(ctx context.Context, item domain.CIOT) (domain.CIOT, error) {
	current, err := s.repo.FindByID(ctx, item.TenantID, item.ID)
	if err != nil {
		return domain.CIOT{}, err
	}
	if current.Status == domain.StatusClosed || current.Status == domain.StatusCanceled {
		return domain.CIOT{}, ErrFinalizedImmutable
	}
	if err := validateCIOT(item); err != nil {
		return domain.CIOT{}, err
	}
	item.Status = current.Status
	item.IdempotencyKey = current.IdempotencyKey
	item.RequestHash = current.RequestHash
	item.CreatedAt = current.CreatedAt
	item.UpdatedAt = time.Now().UTC()
	item.Version = current.Version + 1
	return s.repo.Update(ctx, item)
}

func (s CIOTService) FindByID(ctx context.Context, tenantID string, id string) (domain.CIOT, error) {
	return s.repo.FindByID(ctx, tenantID, id)
}

func (s CIOTService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.CIOT], error) {
	return s.repo.Search(ctx, tenantID, normalizeQuery(q))
}

func (s CIOTService) Submit(ctx context.Context, tenantID string, id string, actorID string) (domain.CIOT, error) {
	return s.transition(ctx, tenantID, id, domain.StatusPending, domain.EventSubmitted, actorID, "")
}

func (s CIOTService) MarkGenerated(ctx context.Context, tenantID string, id string, externalProtocol string, actorID string) (domain.CIOT, error) {
	item, err := s.transition(ctx, tenantID, id, domain.StatusGenerated, domain.EventGenerated, actorID, "")
	if err != nil {
		return domain.CIOT{}, err
	}
	if strings.TrimSpace(externalProtocol) != "" {
		item.ExternalProtocol = strings.TrimSpace(externalProtocol)
		item.UpdatedAt = time.Now().UTC()
		item.Version++
		return s.repo.Update(ctx, item)
	}
	return item, nil
}

func (s CIOTService) Activate(ctx context.Context, tenantID string, id string, actorID string) (domain.CIOT, error) {
	return s.transition(ctx, tenantID, id, domain.StatusActive, domain.EventActivated, actorID, "")
}

func (s CIOTService) Suspend(ctx context.Context, tenantID string, id string, actorID string) (domain.CIOT, error) {
	return s.transition(ctx, tenantID, id, domain.StatusSuspended, domain.EventSuspended, actorID, "")
}

func (s CIOTService) Reactivate(ctx context.Context, tenantID string, id string, actorID string) (domain.CIOT, error) {
	return s.transition(ctx, tenantID, id, domain.StatusActive, domain.EventReactivated, actorID, "")
}

func (s CIOTService) Close(ctx context.Context, tenantID string, id string, actorID string) (domain.CIOT, error) {
	item, err := s.transition(ctx, tenantID, id, domain.StatusClosed, domain.EventClosed, actorID, "")
	if err != nil {
		return domain.CIOT{}, err
	}
	now := time.Now().UTC()
	item.ActualEndDate = &now
	item.UpdatedAt = now
	item.Version++
	return s.repo.Update(ctx, item)
}

func (s CIOTService) Cancel(ctx context.Context, tenantID string, id string, actorID string) (domain.CIOT, error) {
	return s.transition(ctx, tenantID, id, domain.StatusCanceled, domain.EventCanceled, actorID, "")
}

func (s CIOTService) RecordError(ctx context.Context, tenantID string, id string, code string, message string, actorID string) (domain.CIOT, error) {
	item, err := s.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return domain.CIOT{}, err
	}
	now := time.Now().UTC()
	from := item.Status
	item.Status = domain.StatusError
	item.ErrorCode = strings.TrimSpace(code)
	item.ErrorMessage = strings.TrimSpace(message)
	item.UpdatedBy = actorID
	item.UpdatedAt = now
	item.Version++
	saved, err := s.repo.Update(ctx, item)
	if err != nil {
		return domain.CIOT{}, err
	}
	_, err = s.appendHistory(ctx, saved, from, domain.StatusError, domain.EventErrorOccurred, actorID, message)
	if err != nil {
		return domain.CIOT{}, err
	}
	if s.errors != nil {
		_, err = s.errors.Create(ctx, domain.CIOTError{ID: newID("cie"), TenantID: tenantID, CIOTID: id, Code: code, Message: message, CreatedAt: now})
	}
	return saved, err
}

func (s CIOTService) RetryFromError(ctx context.Context, tenantID string, id string, actorID string) (domain.CIOT, error) {
	return s.transition(ctx, tenantID, id, domain.StatusPending, domain.EventSubmitted, actorID, "")
}

func (s CIOTService) Delete(ctx context.Context, tenantID string, id string) error {
	item, err := s.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if item.Status == domain.StatusClosed || item.Status == domain.StatusCanceled {
		return ErrFinalizedImmutable
	}
	return s.repo.Delete(ctx, tenantID, id)
}

func (s CIOTService) transition(ctx context.Context, tenantID string, id string, to domain.CIOTStatus, event domain.CIOTHistoryEvent, actorID string, reason string) (domain.CIOT, error) {
	item, err := s.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return domain.CIOT{}, err
	}
	from := item.Status
	if !canTransition(from, to) {
		return domain.CIOT{}, ErrInvalidTransition
	}
	now := time.Now().UTC()
	item.Status = to
	item.UpdatedBy = actorID
	item.UpdatedAt = now
	item.Version++
	saved, err := s.repo.Update(ctx, item)
	if err != nil {
		return domain.CIOT{}, err
	}
	_, err = s.appendHistory(ctx, saved, from, to, event, actorID, reason)
	return saved, err
}

func (s CIOTService) appendHistory(ctx context.Context, ciot domain.CIOT, from domain.CIOTStatus, to domain.CIOTStatus, event domain.CIOTHistoryEvent, actorID string, reason string) (domain.CIOTStatusHistory, error) {
	if s.history == nil {
		return domain.CIOTStatusHistory{}, nil
	}
	return s.history.Create(ctx, domain.CIOTStatusHistory{ID: newID("cih"), TenantID: ciot.TenantID, CIOTID: ciot.ID, Event: event, FromStatus: from, ToStatus: to, ActorID: actorID, Reason: reason, CreatedAt: time.Now().UTC()})
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

type StatusHistoryService struct{ repo ports.StatusHistoryRepository }

func NewStatusHistoryService(repo ports.StatusHistoryRepository) StatusHistoryService {
	return StatusHistoryService{repo: repo}
}
func (s StatusHistoryService) ListByCIOT(ctx context.Context, tenantID string, ciotID string, q ports.Query) (ports.Page[domain.CIOTStatusHistory], error) {
	return s.repo.ListByCIOT(ctx, tenantID, ciotID, normalizeQuery(q))
}

type PaymentService struct {
	repo    ports.PaymentRepository
	history ports.StatusHistoryRepository
}

func NewPaymentService(repo ports.PaymentRepository, history ports.StatusHistoryRepository) PaymentService {
	return PaymentService{repo: repo, history: history}
}
func (s PaymentService) Create(ctx context.Context, item domain.CIOTPayment) (domain.CIOTPayment, error) {
	if item.TenantID == "" || item.CIOTID == "" || item.Amount <= 0 || strings.TrimSpace(item.PaymentType) == "" {
		return domain.CIOTPayment{}, ErrInvalidData
	}
	now := time.Now().UTC()
	item.ID = newID("cip")
	item.CreatedAt, item.UpdatedAt, item.Version = now, now, 1
	if item.Status == "" {
		item.Status = "pending"
	}
	saved, err := s.repo.Create(ctx, item)
	if err != nil {
		return domain.CIOTPayment{}, err
	}
	if s.history != nil {
		_, err = s.history.Create(ctx, domain.CIOTStatusHistory{ID: newID("cih"), TenantID: item.TenantID, CIOTID: item.CIOTID, Event: domain.EventPaymentRecorded, CreatedAt: now})
	}
	return saved, err
}
func (s PaymentService) ListByCIOT(ctx context.Context, tenantID string, ciotID string, q ports.Query) (ports.Page[domain.CIOTPayment], error) {
	return s.repo.ListByCIOT(ctx, tenantID, ciotID, normalizeQuery(q))
}

type ProviderAttemptService struct {
	repo    ports.ProviderAttemptRepository
	history ports.StatusHistoryRepository
}

func NewProviderAttemptService(repo ports.ProviderAttemptRepository, history ports.StatusHistoryRepository) ProviderAttemptService {
	return ProviderAttemptService{repo: repo, history: history}
}
func (s ProviderAttemptService) Create(ctx context.Context, item domain.CIOTProviderAttempt) (domain.CIOTProviderAttempt, error) {
	if item.TenantID == "" || item.CIOTID == "" || strings.TrimSpace(item.Provider) == "" {
		return domain.CIOTProviderAttempt{}, ErrInvalidData
	}
	if strings.TrimSpace(item.IdempotencyKey) != "" {
		exists, err := s.repo.ExistsIdempotencyKey(ctx, item.TenantID, item.IdempotencyKey)
		if err != nil {
			return domain.CIOTProviderAttempt{}, err
		}
		if exists {
			return domain.CIOTProviderAttempt{}, ErrDuplicateRequest
		}
	}
	now := time.Now().UTC()
	item.ID = newID("cia")
	item.CreatedAt = now
	if item.RequestedAt.IsZero() {
		item.RequestedAt = now
	}
	if item.AttemptNumber == 0 {
		item.AttemptNumber = 1
	}
	saved, err := s.repo.Create(ctx, item)
	if err != nil {
		return domain.CIOTProviderAttempt{}, err
	}
	event := domain.EventProviderAttempted
	if item.Status == "succeeded" {
		event = domain.EventProviderSucceeded
	}
	if item.Status == "failed" {
		event = domain.EventProviderFailed
	}
	if s.history != nil {
		_, err = s.history.Create(ctx, domain.CIOTStatusHistory{ID: newID("cih"), TenantID: item.TenantID, CIOTID: item.CIOTID, Event: event, CreatedAt: now})
	}
	return saved, err
}
func (s ProviderAttemptService) ListByCIOT(ctx context.Context, tenantID string, ciotID string, q ports.Query) (ports.Page[domain.CIOTProviderAttempt], error) {
	return s.repo.ListByCIOT(ctx, tenantID, ciotID, normalizeQuery(q))
}

type ExternalReferenceService struct {
	repo ports.ExternalReferenceRepository
}

func NewExternalReferenceService(repo ports.ExternalReferenceRepository) ExternalReferenceService {
	return ExternalReferenceService{repo: repo}
}
func (s ExternalReferenceService) Upsert(ctx context.Context, item domain.CIOTExternalReference) (domain.CIOTExternalReference, error) {
	if item.TenantID == "" || item.CIOTID == "" || strings.TrimSpace(item.Provider) == "" {
		return domain.CIOTExternalReference{}, ErrInvalidData
	}
	now := time.Now().UTC()
	if item.ID == "" {
		item.ID = newID("cir")
		item.CreatedAt = now
		item.Version = 1
	}
	item.UpdatedAt = now
	return s.repo.Upsert(ctx, item)
}
func (s ExternalReferenceService) FindByCIOT(ctx context.Context, tenantID string, ciotID string) (domain.CIOTExternalReference, error) {
	return s.repo.FindByCIOT(ctx, tenantID, ciotID)
}

type DocumentService struct{ repo ports.DocumentRepository }

func NewDocumentService(repo ports.DocumentRepository) DocumentService {
	return DocumentService{repo: repo}
}
func (s DocumentService) Create(ctx context.Context, item domain.CIOTDocument) (domain.CIOTDocument, error) {
	if item.TenantID == "" || item.CIOTID == "" || item.Type == "" {
		return domain.CIOTDocument{}, ErrInvalidData
	}
	item.ID, item.CreatedAt = newID("cid"), time.Now().UTC()
	return s.repo.Create(ctx, item)
}
func (s DocumentService) ListByCIOT(ctx context.Context, tenantID string, ciotID string, q ports.Query) (ports.Page[domain.CIOTDocument], error) {
	return s.repo.ListByCIOT(ctx, tenantID, ciotID, normalizeQuery(q))
}

type ErrorService struct{ repo ports.ErrorRepository }

func NewErrorService(repo ports.ErrorRepository) ErrorService { return ErrorService{repo: repo} }
func (s ErrorService) Create(ctx context.Context, item domain.CIOTError) (domain.CIOTError, error) {
	if item.TenantID == "" || item.CIOTID == "" || item.Code == "" || item.Message == "" {
		return domain.CIOTError{}, ErrInvalidData
	}
	item.ID, item.CreatedAt = newID("cie"), time.Now().UTC()
	return s.repo.Create(ctx, item)
}
func (s ErrorService) ListByCIOT(ctx context.Context, tenantID string, ciotID string, q ports.Query) (ports.Page[domain.CIOTError], error) {
	return s.repo.ListByCIOT(ctx, tenantID, ciotID, normalizeQuery(q))
}

func InitContract(item *domain.CIOTContract) error {
	if item.TenantID == "" || item.ContractNumber == "" {
		return ErrInvalidData
	}
	now := time.Now().UTC()
	item.ID, item.CreatedAt, item.UpdatedAt, item.Version = newID("cic"), now, now, 1
	if item.StartDate.IsZero() {
		item.StartDate = now
	}
	if item.Status == "" {
		item.Status = "active"
	}
	return nil
}
func InitCarrier(item *domain.CIOTCarrier) error {
	if item.TenantID == "" || item.Document == "" || item.LegalName == "" {
		return ErrInvalidData
	}
	now := time.Now().UTC()
	item.ID, item.CreatedAt, item.UpdatedAt, item.Version = newID("crc"), now, now, 1
	return nil
}
func InitTransporter(item *domain.CIOTTransporter) error {
	if item.TenantID == "" || item.Document == "" || item.Name == "" {
		return ErrInvalidData
	}
	now := time.Now().UTC()
	item.ID, item.CreatedAt, item.UpdatedAt, item.Version = newID("ctr"), now, now, 1
	return nil
}
func InitOperation(item *domain.CIOTOperation) error {
	if item.TenantID == "" || item.OperationNumber == "" {
		return ErrInvalidData
	}
	now := time.Now().UTC()
	item.ID, item.CreatedAt, item.UpdatedAt, item.Version = newID("cop"), now, now, 1
	if item.StartDate.IsZero() {
		item.StartDate = now
	}
	return nil
}
func InitVehicleReference(item *domain.CIOTVehicleReference) error {
	if item.TenantID == "" || (item.VehicleID == "" && item.LicensePlate == "") {
		return ErrInvalidData
	}
	now := time.Now().UTC()
	item.ID, item.CreatedAt, item.UpdatedAt, item.Version = newID("cvr"), now, now, 1
	return nil
}
func InitDriverReference(item *domain.CIOTDriverReference) error {
	if item.TenantID == "" || (item.DriverID == "" && item.DocumentReference == "" && item.NameReference == "") {
		return ErrInvalidData
	}
	now := time.Now().UTC()
	item.ID, item.CreatedAt, item.UpdatedAt, item.Version = newID("cdr"), now, now, 1
	return nil
}
func InitAmount(item *domain.CIOTAmount) error {
	if item.TenantID == "" || item.TotalAmount < 0 {
		return ErrInvalidData
	}
	item.ID, item.CreatedAt = newID("cam"), time.Now().UTC()
	if item.Currency == "" {
		item.Currency = "BRL"
	}
	return nil
}

func validateCIOT(item domain.CIOT) error {
	if item.TenantID == "" {
		return ErrInvalidData
	}
	if !validType(item.Type) {
		return ErrInvalidData
	}
	if item.ExpectedEndDate != nil && !item.StartDate.IsZero() && item.ExpectedEndDate.Before(item.StartDate) {
		return ErrInvalidData
	}
	if item.ActualEndDate != nil && !item.StartDate.IsZero() && item.ActualEndDate.Before(item.StartDate) {
		return ErrInvalidData
	}
	return nil
}

func validType(value domain.CIOTType) bool {
	return value == domain.TypeTACAgregado || value == domain.TypeTACIndependente || value == domain.TypeOther
}

func canTransition(from domain.CIOTStatus, to domain.CIOTStatus) bool {
	switch from {
	case domain.StatusDraft:
		return to == domain.StatusPending
	case domain.StatusPending:
		return to == domain.StatusGenerated || to == domain.StatusCanceled
	case domain.StatusGenerated:
		return to == domain.StatusActive || to == domain.StatusCanceled
	case domain.StatusActive:
		return to == domain.StatusSuspended || to == domain.StatusClosed || to == domain.StatusCanceled
	case domain.StatusSuspended:
		return to == domain.StatusActive
	case domain.StatusError:
		return to == domain.StatusPending
	default:
		return false
	}
}

func normalizeQuery(q ports.Query) ports.Query {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PerPage <= 0 {
		q.PerPage = 25
	}
	if q.PerPage > 100 {
		q.PerPage = 100
	}
	if q.Filters == nil {
		q.Filters = map[string]string{}
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

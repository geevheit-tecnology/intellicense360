package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/tires/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/tires/ports"
)

type TireService struct {
	tires     ports.TireRepository
	movements ports.TireMovementRepository
	history   ports.TireHistoryRepository
}

func NewTireService(tires ports.TireRepository, movements ports.TireMovementRepository) TireService {
	return TireService{tires: tires, movements: movements}
}

func NewTireLifecycleService(tires ports.TireRepository, movements ports.TireMovementRepository, history ports.TireHistoryRepository) TireService {
	return TireService{tires: tires, movements: movements, history: history}
}

func (s TireService) Create(ctx context.Context, tire domain.Tire) (domain.Tire, error) {
	if err := s.validate(ctx, tire, ""); err != nil {
		return domain.Tire{}, err
	}
	now := time.Now().UTC()
	tire.ID = newID("tir")
	if tire.Status == "" {
		tire.Status = domain.TireStatusNew
	}
	if !validStatus(tire.Status) {
		return domain.Tire{}, ErrInvalidStatus
	}
	if tire.Condition != "" && !validCondition(tire.Condition) {
		return domain.Tire{}, ErrInvalidCondition
	}
	tire.CreatedAt = now
	tire.UpdatedAt = now
	tire.Version = 1
	saved, err := s.tires.Create(ctx, tire)
	if err == nil {
		_ = s.recordHistory(ctx, saved, "created", "", saved.Status, tire.CreatedBy, "")
	}
	return saved, err
}

func (s TireService) Update(ctx context.Context, tire domain.Tire) (domain.Tire, error) {
	if err := s.validate(ctx, tire, tire.ID); err != nil {
		return domain.Tire{}, err
	}
	current, err := s.tires.FindByID(ctx, tire.TenantID, tire.ID)
	if err != nil {
		return domain.Tire{}, err
	}
	if tire.Status == "" {
		tire.Status = current.Status
	}
	if tire.Status != current.Status {
		return domain.Tire{}, ErrInvalidTransition
	}
	tire.CreatedAt = current.CreatedAt
	tire.UpdatedAt = time.Now().UTC()
	tire.Version = current.Version + 1
	return s.tires.Update(ctx, tire)
}

func (s TireService) Delete(ctx context.Context, tenantID string, tireID string) error {
	return s.tires.Delete(ctx, tenantID, tireID)
}

func (s TireService) FindByID(ctx context.Context, tenantID string, tireID string) (domain.Tire, error) {
	return s.tires.FindByID(ctx, tenantID, tireID)
}

func (s TireService) Search(ctx context.Context, tenantID string, query ports.Query) (ports.Page[domain.Tire], error) {
	return s.tires.Search(ctx, tenantID, normalizeQuery(query))
}

func (s TireService) Receive(ctx context.Context, tenantID string, tireID string, reason string, actorID string) (domain.Tire, error) {
	tire, err := s.tires.FindByID(ctx, tenantID, tireID)
	if err != nil {
		return domain.Tire{}, err
	}
	if !canTransition(tire.Status, domain.TireStatusInStock) {
		return domain.Tire{}, ErrInvalidTransition
	}
	from := tire.Status
	tire.Status = domain.TireStatusInStock
	tire.UpdatedBy = actorID
	tire.UpdatedAt = time.Now().UTC()
	tire.Version++
	saved, err := s.tires.Update(ctx, tire)
	if err != nil {
		return domain.Tire{}, err
	}
	_ = s.recordMovement(ctx, saved, domain.MovementReceipt, tire.CurrentKM, reason, actorID)
	_ = s.recordHistory(ctx, saved, "received", from, saved.Status, actorID, reason)
	return saved, nil
}

func (s TireService) Install(ctx context.Context, tenantID string, tireID string, vehicleID string, position string, km int64, actorID string) (domain.Tire, error) {
	tire, err := s.loadMovable(ctx, tenantID, tireID)
	if err != nil {
		return domain.Tire{}, err
	}
	if !canTransition(tire.Status, domain.TireStatusInstalled) {
		return domain.Tire{}, ErrInvalidTransition
	}
	from := tire.Status
	tire.VehicleID = vehicleID
	tire.Position = position
	tire.CurrentKM = km
	tire.Status = domain.TireStatusInstalled
	tire.UpdatedBy = actorID
	tire.UpdatedAt = time.Now().UTC()
	tire.Version++
	saved, err := s.tires.Update(ctx, tire)
	if err != nil {
		return domain.Tire{}, err
	}
	_ = s.recordMovement(ctx, saved, domain.MovementInstallation, km, "", actorID)
	_ = s.recordHistory(ctx, saved, "installed", from, saved.Status, actorID, position)
	return saved, nil
}

func (s TireService) Remove(ctx context.Context, tenantID string, tireID string, km int64, reason string, actorID string) (domain.Tire, error) {
	tire, err := s.loadMovable(ctx, tenantID, tireID)
	if err != nil {
		return domain.Tire{}, err
	}
	if !canTransition(tire.Status, domain.TireStatusRemoved) {
		return domain.Tire{}, ErrInvalidTransition
	}
	from := tire.Status
	tire.VehicleID = ""
	tire.Position = ""
	tire.CurrentKM = km
	tire.Status = domain.TireStatusRemoved
	tire.UpdatedBy = actorID
	tire.UpdatedAt = time.Now().UTC()
	tire.Version++
	saved, err := s.tires.Update(ctx, tire)
	if err != nil {
		return domain.Tire{}, err
	}
	_ = s.recordMovement(ctx, saved, domain.MovementRemoval, km, reason, actorID)
	_ = s.recordHistory(ctx, saved, "removed", from, saved.Status, actorID, reason)
	return saved, nil
}

func (s TireService) Rotate(ctx context.Context, tenantID string, tireID string, position string, km int64, actorID string) (domain.Tire, error) {
	tire, err := s.loadMovable(ctx, tenantID, tireID)
	if err != nil {
		return domain.Tire{}, err
	}
	if tire.Status != domain.TireStatusInstalled {
		return domain.Tire{}, ErrInvalidTransition
	}
	tire.Position = position
	tire.CurrentKM = km
	tire.Status = domain.TireStatusInstalled
	tire.UpdatedBy = actorID
	tire.UpdatedAt = time.Now().UTC()
	tire.Version++
	saved, err := s.tires.Update(ctx, tire)
	if err != nil {
		return domain.Tire{}, err
	}
	_ = s.recordMovement(ctx, saved, domain.MovementRotation, km, "", actorID)
	_ = s.recordHistory(ctx, saved, "moved", domain.TireStatusInstalled, saved.Status, actorID, position)
	return saved, nil
}

func (s TireService) SendToRecap(ctx context.Context, tenantID string, tireID string, reason string, actorID string) (domain.Tire, error) {
	tire, err := s.loadMovable(ctx, tenantID, tireID)
	if err != nil {
		return domain.Tire{}, err
	}
	if !canTransition(tire.Status, domain.TireStatusUnderRetread) {
		return domain.Tire{}, ErrInvalidTransition
	}
	from := tire.Status
	tire.VehicleID = ""
	tire.Position = ""
	tire.Status = domain.TireStatusUnderRetread
	tire.UpdatedBy = actorID
	tire.UpdatedAt = time.Now().UTC()
	tire.Version++
	saved, err := s.tires.Update(ctx, tire)
	if err != nil {
		return domain.Tire{}, err
	}
	_ = s.recordMovement(ctx, saved, domain.MovementRetread, tire.CurrentKM, reason, actorID)
	_ = s.recordHistory(ctx, saved, "retread_sent", from, saved.Status, actorID, reason)
	return saved, nil
}

func (s TireService) ReturnFromRecap(ctx context.Context, tenantID string, tireID string, actorID string) (domain.Tire, error) {
	tire, err := s.loadMovable(ctx, tenantID, tireID)
	if err != nil {
		return domain.Tire{}, err
	}
	if !canTransition(tire.Status, domain.TireStatusRetreaded) {
		return domain.Tire{}, ErrInvalidTransition
	}
	from := tire.Status
	tire.Status = domain.TireStatusRetreaded
	tire.UpdatedBy = actorID
	tire.UpdatedAt = time.Now().UTC()
	tire.Version++
	retreaded, err := s.tires.Update(ctx, tire)
	if err != nil {
		return domain.Tire{}, err
	}
	_ = s.recordHistory(ctx, retreaded, "retreaded", from, retreaded.Status, actorID, "")
	if !canTransition(retreaded.Status, domain.TireStatusInStock) {
		return retreaded, nil
	}
	tire = retreaded
	tire.Status = domain.TireStatusInStock
	tire.RecapCount++
	tire.UpdatedBy = actorID
	tire.UpdatedAt = time.Now().UTC()
	tire.Version++
	saved, err := s.tires.Update(ctx, tire)
	if err == nil {
		_ = s.recordMovement(ctx, saved, domain.MovementReturn, saved.CurrentKM, "", actorID)
		_ = s.recordHistory(ctx, saved, "returned_to_stock", retreaded.Status, saved.Status, actorID, "")
	}
	return saved, err
}

func (s TireService) Dispose(ctx context.Context, tenantID string, tireID string, reason string, actorID string) (domain.Tire, error) {
	tire, err := s.tires.FindByID(ctx, tenantID, tireID)
	if err != nil {
		return domain.Tire{}, err
	}
	if !canTransition(tire.Status, domain.TireStatusDisposed) {
		if canTransition(tire.Status, domain.TireStatusEndOfLife) {
			from := tire.Status
			tire.Status = domain.TireStatusEndOfLife
			tire.UpdatedBy = actorID
			tire.UpdatedAt = time.Now().UTC()
			tire.Version++
			var eolErr error
			tire, eolErr = s.tires.Update(ctx, tire)
			if eolErr != nil {
				return domain.Tire{}, eolErr
			}
			_ = s.recordHistory(ctx, tire, "end_of_life", from, tire.Status, actorID, reason)
		} else {
			return domain.Tire{}, ErrInvalidTransition
		}
	}
	from := tire.Status
	tire.VehicleID = ""
	tire.Position = ""
	tire.Status = domain.TireStatusDisposed
	tire.UpdatedBy = actorID
	tire.UpdatedAt = time.Now().UTC()
	tire.Version++
	saved, err := s.tires.Update(ctx, tire)
	if err != nil {
		return domain.Tire{}, err
	}
	_ = s.recordMovement(ctx, saved, domain.MovementDisposal, tire.CurrentKM, reason, actorID)
	_ = s.recordHistory(ctx, saved, "disposed", from, saved.Status, actorID, reason)
	return saved, nil
}

func (s TireService) ValidateTireAccess(ctx context.Context, tenantID string, tireID string) error {
	exists, err := s.tires.Exists(ctx, tenantID, tireID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}
	return nil
}

func (s TireService) validate(ctx context.Context, tire domain.Tire, exceptID string) error {
	if tire.TenantID == "" || strings.TrimSpace(tire.DOT) == "" || strings.TrimSpace(tire.SerialNumber) == "" || strings.TrimSpace(tire.FireNumber) == "" {
		return ErrValidation
	}
	if tire.CurrentTreadMM < 0 || tire.OriginalTreadMM < 0 || tire.MinimumTreadMM < 0 || tire.MinimumTreadMM > tire.OriginalTreadMM {
		return ErrValidation
	}
	if ok, err := s.tires.ExistsSerialNumber(ctx, tire.TenantID, tire.SerialNumber, exceptID); err != nil || ok {
		if err != nil {
			return err
		}
		return ErrSerialNumberTaken
	}
	if ok, err := s.tires.ExistsFireNumber(ctx, tire.TenantID, tire.FireNumber, exceptID); err != nil || ok {
		if err != nil {
			return err
		}
		return ErrFireNumberTaken
	}
	return nil
}

func (s TireService) loadMovable(ctx context.Context, tenantID string, tireID string) (domain.Tire, error) {
	tire, err := s.tires.FindByID(ctx, tenantID, tireID)
	if err != nil {
		return domain.Tire{}, err
	}
	if tire.Status == domain.TireStatusDisposed || tire.Status == domain.TireStatusLost || tire.Status == domain.TireStatusEndOfLife {
		return domain.Tire{}, ErrDisposedTire
	}
	return tire, nil
}

func (s TireService) recordMovement(ctx context.Context, tire domain.Tire, movementType domain.MovementType, km int64, reason string, actorID string) error {
	now := time.Now().UTC()
	_, err := s.movements.Create(ctx, domain.TireMovement{
		ID: newID("tmv"), TenantID: tire.TenantID, TireID: tire.ID, MovementType: movementType,
		VehicleID: tire.VehicleID, Position: tire.Position, KM: km, Reason: reason, PerformedBy: actorID,
		MovementDate: now, CreatedAt: now, UpdatedAt: now,
	})
	return err
}

func (s TireService) recordHistory(ctx context.Context, tire domain.Tire, event string, from domain.TireStatus, to domain.TireStatus, actorID string, notes string) error {
	if s.history == nil {
		return nil
	}
	now := time.Now().UTC()
	_, err := s.history.Create(ctx, domain.TireHistory{ID: newID("thi"), TenantID: tire.TenantID, TireID: tire.ID, Event: event, FromStatus: from, ToStatus: to, ActorID: actorID, Notes: notes, CreatedAt: now})
	return err
}

func validStatus(status domain.TireStatus) bool {
	switch status {
	case domain.TireStatusNew, domain.TireStatusInStock, domain.TireStatusInstalled, domain.TireStatusRemoved, domain.TireStatusUnderInspection, domain.TireStatusUnderRetread, domain.TireStatusRetreaded, domain.TireStatusReserved, domain.TireStatusDamaged, domain.TireStatusEndOfLife, domain.TireStatusDisposed, domain.TireStatusLost, domain.TireStatusRecapping:
		return true
	default:
		return false
	}
}

func validCondition(condition domain.TireCondition) bool {
	switch condition {
	case domain.ConditionExcellent, domain.ConditionGood, domain.ConditionAttention, domain.ConditionHeavyWear, domain.ConditionAtLimit, domain.ConditionCritical, domain.ConditionDamaged, domain.ConditionEndOfLife:
		return true
	default:
		return false
	}
}

func canTransition(from domain.TireStatus, to domain.TireStatus) bool {
	if from == to {
		return false
	}
	allowed := map[domain.TireStatus][]domain.TireStatus{
		domain.TireStatusNew:          {domain.TireStatusInStock},
		domain.TireStatusInStock:      {domain.TireStatusInstalled, domain.TireStatusReserved, domain.TireStatusLost, domain.TireStatusEndOfLife},
		domain.TireStatusInstalled:    {domain.TireStatusRemoved, domain.TireStatusDamaged, domain.TireStatusUnderInspection},
		domain.TireStatusRemoved:      {domain.TireStatusInStock, domain.TireStatusUnderRetread, domain.TireStatusDamaged, domain.TireStatusEndOfLife},
		domain.TireStatusUnderRetread: {domain.TireStatusRetreaded, domain.TireStatusDamaged},
		domain.TireStatusRetreaded:    {domain.TireStatusInStock},
		domain.TireStatusDamaged:      {domain.TireStatusEndOfLife, domain.TireStatusUnderRetread},
		domain.TireStatusEndOfLife:    {domain.TireStatusDisposed},
		domain.TireStatusRecapping:    {domain.TireStatusRetreaded, domain.TireStatusInStock},
	}
	for _, candidate := range allowed[from] {
		if candidate == to {
			return true
		}
	}
	return false
}

type TireInspectionService struct {
	tires       ports.TireRepository
	inspections ports.TireInspectionRepository
}

func NewTireInspectionService(tires ports.TireRepository, inspections ports.TireInspectionRepository) TireInspectionService {
	return TireInspectionService{tires: tires, inspections: inspections}
}

func (s TireInspectionService) Register(ctx context.Context, inspection domain.TireInspection) (domain.TireInspection, error) {
	if err := s.validate(ctx, inspection); err != nil {
		return domain.TireInspection{}, err
	}
	now := time.Now().UTC()
	inspection.ID = newID("tin")
	if inspection.InspectionDate.IsZero() {
		inspection.InspectionDate = now
	}
	inspection.CreatedAt = now
	inspection.UpdatedAt = now
	return s.inspections.Create(ctx, inspection)
}

func (s TireInspectionService) Update(ctx context.Context, inspection domain.TireInspection) (domain.TireInspection, error) {
	if err := s.validate(ctx, inspection); err != nil {
		return domain.TireInspection{}, err
	}
	current, err := s.inspections.FindByID(ctx, inspection.TenantID, inspection.TireID, inspection.ID)
	if err != nil {
		return domain.TireInspection{}, err
	}
	inspection.CreatedAt = current.CreatedAt
	inspection.UpdatedAt = time.Now().UTC()
	return s.inspections.Update(ctx, inspection)
}

func (s TireInspectionService) Delete(ctx context.Context, tenantID string, tireID string, inspectionID string) error {
	return s.inspections.Delete(ctx, tenantID, tireID, inspectionID)
}

func (s TireInspectionService) List(ctx context.Context, tenantID string, tireID string, query ports.Query) (ports.Page[domain.TireInspection], error) {
	return s.inspections.List(ctx, tenantID, tireID, normalizeQuery(query))
}

func (s TireInspectionService) validate(ctx context.Context, inspection domain.TireInspection) error {
	if inspection.TenantID == "" || inspection.TireID == "" || inspection.TreadMM < 0 {
		return ErrValidation
	}
	if inspection.Pressure < 0 {
		return ErrValidation
	}
	exists, err := s.tires.Exists(ctx, inspection.TenantID, inspection.TireID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotFound
	}
	return nil
}

type TireHistoryService struct{ history ports.TireHistoryRepository }

func NewTireHistoryService(history ports.TireHistoryRepository) TireHistoryService {
	return TireHistoryService{history: history}
}

func (s TireHistoryService) List(ctx context.Context, tenantID string, tireID string, query ports.Query) (ports.Page[domain.TireHistory], error) {
	return s.history.List(ctx, tenantID, tireID, normalizeQuery(query))
}

type TireMovementService struct{ movements ports.TireMovementRepository }

func NewTireMovementService(movements ports.TireMovementRepository) TireMovementService {
	return TireMovementService{movements: movements}
}

func (s TireMovementService) Register(ctx context.Context, movement domain.TireMovement) (domain.TireMovement, error) {
	now := time.Now().UTC()
	movement.ID = newID("tmv")
	if movement.MovementDate.IsZero() {
		movement.MovementDate = now
	}
	movement.CreatedAt = now
	movement.UpdatedAt = now
	return s.movements.Create(ctx, movement)
}

func (s TireMovementService) List(ctx context.Context, tenantID string, tireID string, query ports.Query) (ports.Page[domain.TireMovement], error) {
	return s.movements.List(ctx, tenantID, tireID, normalizeQuery(query))
}

func normalizeQuery(query ports.Query) ports.Query {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 20
	}
	if query.SortOrder != "desc" {
		query.SortOrder = "asc"
	}
	return query
}

func newID(prefix string) string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return prefix + "_" + hex.EncodeToString(b[:])
}

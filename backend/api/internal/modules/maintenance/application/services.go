package application

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/ports"
)

type WorkOrderService struct {
	repo    ports.WorkOrderRepository
	history ports.HistoryRepository
}

func NewWorkOrderService(repo ports.WorkOrderRepository, history ports.HistoryRepository) WorkOrderService {
	return WorkOrderService{repo: repo, history: history}
}

func (s WorkOrderService) Create(ctx context.Context, wo domain.WorkOrder) (domain.WorkOrder, error) {
	if wo.TenantID == "" || strings.TrimSpace(wo.Title) == "" {
		return domain.WorkOrder{}, ErrValidation
	}
	if wo.Status != "" && !validWorkOrderStatus(wo.Status) {
		return domain.WorkOrder{}, ErrValidation
	}
	if wo.Priority != "" && !validWorkOrderPriority(wo.Priority) {
		return domain.WorkOrder{}, ErrValidation
	}
	if wo.Kind != "" && !validMaintenanceKind(wo.Kind) {
		return domain.WorkOrder{}, ErrValidation
	}
	if wo.Code != "" {
		if ok, err := s.repo.ExistsCode(ctx, wo.TenantID, wo.Code, ""); err != nil || ok {
			if err != nil {
				return domain.WorkOrder{}, err
			}
			return domain.WorkOrder{}, ErrValidation
		}
	}
	now := time.Now().UTC()
	wo.ID = newID("mwo")
	if wo.Code == "" {
		wo.Code = "MO-" + wo.ID[4:12]
	}
	if wo.Status == "" {
		wo.Status = domain.WorkOrderOpen
	}
	if wo.Priority == "" {
		wo.Priority = domain.PriorityMedium
	}
	if wo.Kind == "" {
		wo.Kind = domain.KindCorrective
	}
	wo.OpenedAt = now
	wo.CreatedAt = now
	wo.UpdatedAt = now
	wo.Version = 1
	saved, err := s.repo.Create(ctx, wo)
	if err == nil {
		_ = s.record(ctx, saved, "work_order.created", wo.CreatedBy, "")
	}
	return saved, err
}

func (s WorkOrderService) FindByID(ctx context.Context, tenantID string, id string) (domain.WorkOrder, error) {
	return s.repo.FindByID(ctx, tenantID, id)
}

func (s WorkOrderService) Search(ctx context.Context, tenantID string, query ports.Query) (ports.Page[domain.WorkOrder], error) {
	return s.repo.Search(ctx, tenantID, normalize(query))
}

func (s WorkOrderService) Update(ctx context.Context, wo domain.WorkOrder) (domain.WorkOrder, error) {
	current, err := s.repo.FindByID(ctx, wo.TenantID, wo.ID)
	if err != nil {
		return domain.WorkOrder{}, err
	}
	if strings.TrimSpace(wo.Title) == "" {
		return domain.WorkOrder{}, ErrValidation
	}
	if wo.Status != "" && !validWorkOrderStatus(wo.Status) {
		return domain.WorkOrder{}, ErrValidation
	}
	if wo.Priority != "" && !validWorkOrderPriority(wo.Priority) {
		return domain.WorkOrder{}, ErrValidation
	}
	if wo.Kind != "" && !validMaintenanceKind(wo.Kind) {
		return domain.WorkOrder{}, ErrValidation
	}
	if wo.Code != "" {
		if ok, err := s.repo.ExistsCode(ctx, wo.TenantID, wo.Code, wo.ID); err != nil || ok {
			if err != nil {
				return domain.WorkOrder{}, err
			}
			return domain.WorkOrder{}, ErrValidation
		}
	}
	if wo.Code == "" {
		wo.Code = current.Code
	}
	wo.OpenedAt = current.OpenedAt
	wo.CreatedAt = current.CreatedAt
	wo.StartedAt = current.StartedAt
	wo.CompletedAt = current.CompletedAt
	wo.CancelledAt = current.CancelledAt
	if wo.Status == "" {
		wo.Status = current.Status
	}
	wo.Version = current.Version + 1
	wo.UpdatedAt = time.Now().UTC()
	saved, err := s.repo.Update(ctx, wo)
	if err == nil {
		_ = s.record(ctx, saved, "work_order.updated", wo.UpdatedBy, "")
	}
	return saved, err
}

func (s WorkOrderService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

func (s WorkOrderService) Start(ctx context.Context, tenantID string, id string, actorID string) (domain.WorkOrder, error) {
	return s.transition(ctx, tenantID, id, actorID, domain.WorkOrderInProgress, "work_order.started")
}

func (s WorkOrderService) Complete(ctx context.Context, tenantID string, id string, actorID string) (domain.WorkOrder, error) {
	return s.transition(ctx, tenantID, id, actorID, domain.WorkOrderCompleted, "work_order.completed")
}

func (s WorkOrderService) Cancel(ctx context.Context, tenantID string, id string, actorID string) (domain.WorkOrder, error) {
	return s.transition(ctx, tenantID, id, actorID, domain.WorkOrderCanceled, "work_order.canceled")
}

func (s WorkOrderService) ValidateMaintenanceAccess(ctx context.Context, tenantID string, id string) error {
	ok, err := s.repo.Exists(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}

func (s WorkOrderService) transition(ctx context.Context, tenantID string, id string, actorID string, status domain.WorkOrderStatus, event string) (domain.WorkOrder, error) {
	wo, err := s.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return domain.WorkOrder{}, err
	}
	if isTerminalStatus(wo.Status) {
		return domain.WorkOrder{}, ErrInvalidStatus
	}
	now := time.Now().UTC()
	wo.Status = status
	wo.UpdatedBy = actorID
	wo.UpdatedAt = now
	wo.Version++
	if status == domain.WorkOrderInProgress {
		wo.StartedAt = &now
	}
	if status == domain.WorkOrderCompleted {
		wo.CompletedAt = &now
	}
	if status == domain.WorkOrderCanceled || status == domain.WorkOrderCancelled {
		wo.CancelledAt = &now
	}
	saved, err := s.repo.Update(ctx, wo)
	if err == nil {
		_ = s.record(ctx, saved, event, actorID, "")
	}
	return saved, err
}

func (s WorkOrderService) record(ctx context.Context, wo domain.WorkOrder, event string, actorID string, notes string) error {
	if s.history == nil {
		return nil
	}
	now := time.Now().UTC()
	_, err := s.history.Create(ctx, domain.MaintenanceHistory{ID: newID("mhi"), TenantID: wo.TenantID, WorkOrderID: wo.ID, Event: event, ActorID: actorID, Notes: notes, CreatedAt: now})
	return err
}

func validWorkOrderStatus(status domain.WorkOrderStatus) bool {
	switch status {
	case domain.WorkOrderDraft, domain.WorkOrderOpen, domain.WorkOrderWaiting, domain.WorkOrderApproved, domain.WorkOrderExecuting, domain.WorkOrderInProgress, domain.WorkOrderPaused, domain.WorkOrderCompleted, domain.WorkOrderCanceled, domain.WorkOrderCancelled:
		return true
	default:
		return false
	}
}

func validWorkOrderPriority(priority domain.WorkOrderPriority) bool {
	switch priority {
	case domain.PriorityLow, domain.PriorityMedium, domain.PriorityHigh, domain.PriorityCritical:
		return true
	default:
		return false
	}
}

func validMaintenanceKind(kind domain.MaintenanceKind) bool {
	switch kind {
	case domain.KindPreventive, domain.KindCorrective, domain.KindPredictive, domain.KindInspection, domain.KindEmergency, domain.KindWarranty, domain.KindExternal, domain.KindInternal:
		return true
	default:
		return false
	}
}

func isTerminalStatus(status domain.WorkOrderStatus) bool {
	return status == domain.WorkOrderCompleted || status == domain.WorkOrderCanceled || status == domain.WorkOrderCancelled
}

type PreventivePlanService struct {
	repo ports.PreventivePlanRepository
}

func NewPreventivePlanService(repo ports.PreventivePlanRepository) PreventivePlanService {
	return PreventivePlanService{repo: repo}
}
func (s PreventivePlanService) Create(ctx context.Context, p domain.PreventivePlan) (domain.PreventivePlan, error) {
	if p.TenantID == "" || p.Name == "" || p.IntervalValue <= 0 {
		return domain.PreventivePlan{}, ErrValidation
	}
	stamp(&p.ID, &p.CreatedAt, &p.UpdatedAt, &p.Version, "mpp")
	p.Active = true
	return s.repo.Create(ctx, p)
}
func (s PreventivePlanService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.PreventivePlan], error) {
	return s.repo.Search(ctx, tenantID, normalize(q))
}
func (s PreventivePlanService) Update(ctx context.Context, p domain.PreventivePlan) (domain.PreventivePlan, error) {
	p.UpdatedAt = time.Now().UTC()
	p.Version++
	return s.repo.Update(ctx, p)
}
func (s PreventivePlanService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

type ServiceTypeService struct{ repo ports.ServiceTypeRepository }

func NewServiceTypeService(repo ports.ServiceTypeRepository) ServiceTypeService {
	return ServiceTypeService{repo: repo}
}
func (s ServiceTypeService) Create(ctx context.Context, st domain.ServiceType) (domain.ServiceType, error) {
	if st.TenantID == "" || st.Name == "" || st.Code == "" {
		return domain.ServiceType{}, ErrValidation
	}
	stamp(&st.ID, &st.CreatedAt, &st.UpdatedAt, &st.Version, "mst")
	return s.repo.Create(ctx, st)
}
func (s ServiceTypeService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.ServiceType], error) {
	return s.repo.Search(ctx, tenantID, normalize(q))
}
func (s ServiceTypeService) Update(ctx context.Context, st domain.ServiceType) (domain.ServiceType, error) {
	st.UpdatedAt = time.Now().UTC()
	st.Version++
	return s.repo.Update(ctx, st)
}
func (s ServiceTypeService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

type LaborService struct {
	repo       ports.LaborRepository
	workOrders ports.WorkOrderRepository
}

func NewLaborService(repo ports.LaborRepository, workOrders ports.WorkOrderRepository) LaborService {
	return LaborService{repo: repo, workOrders: workOrders}
}
func (s LaborService) Add(ctx context.Context, l domain.LaborEntry) (domain.LaborEntry, error) {
	if l.TenantID == "" || l.WorkOrderID == "" || l.Technician == "" || l.Hours <= 0 {
		return domain.LaborEntry{}, ErrValidation
	}
	if ok, _ := s.workOrders.Exists(ctx, l.TenantID, l.WorkOrderID); !ok {
		return domain.LaborEntry{}, ErrNotFound
	}
	now := time.Now().UTC()
	l.ID = newID("mla")
	l.StartedAt = now
	l.CreatedAt = now
	l.UpdatedAt = now
	return s.repo.Create(ctx, l)
}
func (s LaborService) List(ctx context.Context, tenantID string, workOrderID string, q ports.Query) (ports.Page[domain.LaborEntry], error) {
	return s.repo.List(ctx, tenantID, workOrderID, normalize(q))
}
func (s LaborService) Delete(ctx context.Context, tenantID string, workOrderID string, id string) error {
	return s.repo.Delete(ctx, tenantID, workOrderID, id)
}

type DowntimeService struct {
	repo       ports.DowntimeRepository
	workOrders ports.WorkOrderRepository
}

func NewDowntimeService(repo ports.DowntimeRepository, workOrders ports.WorkOrderRepository) DowntimeService {
	return DowntimeService{repo: repo, workOrders: workOrders}
}
func (s DowntimeService) Start(ctx context.Context, d domain.Downtime) (domain.Downtime, error) {
	if d.TenantID == "" || d.WorkOrderID == "" {
		return domain.Downtime{}, ErrValidation
	}
	if ok, _ := s.workOrders.Exists(ctx, d.TenantID, d.WorkOrderID); !ok {
		return domain.Downtime{}, ErrNotFound
	}
	now := time.Now().UTC()
	d.ID = newID("mdt")
	d.StartedAt = now
	d.CreatedAt = now
	d.UpdatedAt = now
	return s.repo.Create(ctx, d)
}
func (s DowntimeService) End(ctx context.Context, tenantID string, workOrderID string, id string) (domain.Downtime, error) {
	return s.repo.End(ctx, tenantID, workOrderID, id)
}
func (s DowntimeService) List(ctx context.Context, tenantID string, workOrderID string, q ports.Query) (ports.Page[domain.Downtime], error) {
	return s.repo.List(ctx, tenantID, workOrderID, normalize(q))
}

type HistoryService struct{ repo ports.HistoryRepository }

func NewHistoryService(repo ports.HistoryRepository) HistoryService {
	return HistoryService{repo: repo}
}
func (s HistoryService) List(ctx context.Context, tenantID string, workOrderID string, q ports.Query) (ports.Page[domain.MaintenanceHistory], error) {
	return s.repo.List(ctx, tenantID, workOrderID, normalize(q))
}

type CatalogService struct{ repo ports.CatalogRepository }

func NewCatalogService(repo ports.CatalogRepository) CatalogService {
	return CatalogService{repo: repo}
}

func (s CatalogService) Create(ctx context.Context, item domain.MaintenanceCatalog) (domain.MaintenanceCatalog, error) {
	if item.TenantID == "" || item.Name == "" || item.Code == "" {
		return domain.MaintenanceCatalog{}, ErrValidation
	}
	stamp(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.Version, "mca")
	return s.repo.Create(ctx, item)
}

func (s CatalogService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.MaintenanceCatalog], error) {
	return s.repo.Search(ctx, tenantID, normalize(q))
}

func (s CatalogService) Update(ctx context.Context, item domain.MaintenanceCatalog) (domain.MaintenanceCatalog, error) {
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}

func (s CatalogService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

type WorkshopService struct{ repo ports.WorkshopRepository }

func NewWorkshopService(repo ports.WorkshopRepository) WorkshopService {
	return WorkshopService{repo: repo}
}

func (s WorkshopService) Create(ctx context.Context, item domain.Workshop) (domain.Workshop, error) {
	if item.TenantID == "" || item.Name == "" {
		return domain.Workshop{}, ErrValidation
	}
	stamp(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.Version, "mws")
	return s.repo.Create(ctx, item)
}

func (s WorkshopService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Workshop], error) {
	return s.repo.Search(ctx, tenantID, normalize(q))
}

func (s WorkshopService) Update(ctx context.Context, item domain.Workshop) (domain.Workshop, error) {
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}

func (s WorkshopService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

type TechnicianService struct{ repo ports.TechnicianRepository }

func NewTechnicianService(repo ports.TechnicianRepository) TechnicianService {
	return TechnicianService{repo: repo}
}

func (s TechnicianService) Create(ctx context.Context, item domain.Technician) (domain.Technician, error) {
	if item.TenantID == "" || item.Name == "" {
		return domain.Technician{}, ErrValidation
	}
	stamp(&item.ID, &item.CreatedAt, &item.UpdatedAt, &item.Version, "mte")
	item.Active = true
	return s.repo.Create(ctx, item)
}

func (s TechnicianService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Technician], error) {
	return s.repo.Search(ctx, tenantID, normalize(q))
}

func (s TechnicianService) Update(ctx context.Context, item domain.Technician) (domain.Technician, error) {
	item.UpdatedAt = time.Now().UTC()
	item.Version++
	return s.repo.Update(ctx, item)
}

func (s TechnicianService) Delete(ctx context.Context, tenantID string, id string) error {
	return s.repo.Delete(ctx, tenantID, id)
}

func normalize(q ports.Query) ports.Query {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.PageSize <= 0 || q.PageSize > 100 {
		q.PageSize = 20
	}
	if q.SortOrder != "desc" {
		q.SortOrder = "asc"
	}
	return q
}
func stamp(id *string, createdAt *time.Time, updatedAt *time.Time, version *int64, prefix string) {
	now := time.Now().UTC()
	*id = newID(prefix)
	*createdAt = now
	*updatedAt = now
	*version = 1
}
func newID(prefix string) string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return prefix + "_" + hex.EncodeToString(b[:])
}

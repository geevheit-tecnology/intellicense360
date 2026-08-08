package application_test

import (
	"context"
	"testing"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/infrastructure"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/ports"
)

func TestWorkOrderLifecycleTenantIsolationAndHistory(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	service := application.NewWorkOrderService(store.WorkOrders(), store.History())
	history := application.NewHistoryService(store.History())
	created, err := service.Create(context.Background(), validWorkOrder("tenant-a"))
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	_, _ = service.Create(context.Background(), validWorkOrder("tenant-b"))
	started, err := service.Start(context.Background(), "tenant-a", created.ID, "actor")
	if err != nil || started.Status != domain.WorkOrderInProgress {
		t.Fatalf("start failed: %+v err=%v", started, err)
	}
	completed, err := service.Complete(context.Background(), "tenant-a", created.ID, "actor")
	if err != nil || completed.Status != domain.WorkOrderCompleted {
		t.Fatalf("complete failed: %+v err=%v", completed, err)
	}
	page, err := service.Search(context.Background(), "tenant-a", ports.Query{Page: 1, PageSize: 10})
	if err != nil || page.TotalItems != 1 {
		t.Fatalf("tenant isolation failed: %+v err=%v", page, err)
	}
	events, err := history.List(context.Background(), "tenant-a", created.ID, ports.Query{Page: 1, PageSize: 10})
	if err != nil || events.TotalItems < 3 {
		t.Fatalf("expected history events: %+v err=%v", events, err)
	}
}

func TestPreventivePlanServiceTypeLaborAndDowntime(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	workOrders := application.NewWorkOrderService(store.WorkOrders(), store.History())
	plans := application.NewPreventivePlanService(store.PreventivePlans())
	serviceTypes := application.NewServiceTypeService(store.ServiceTypes())
	labor := application.NewLaborService(store.Labor(), store.WorkOrders())
	downtime := application.NewDowntimeService(store.Downtime(), store.WorkOrders())
	wo, _ := workOrders.Create(context.Background(), validWorkOrder("tenant-a"))
	st, err := serviceTypes.Create(context.Background(), domain.ServiceType{TenantID: "tenant-a", Name: "Oil change", Code: "oil_change"})
	if err != nil {
		t.Fatalf("service type failed: %v", err)
	}
	if _, err := plans.Create(context.Background(), domain.PreventivePlan{TenantID: "tenant-a", Name: "Every 10k km", ServiceTypeID: st.ID, Frequency: domain.FrequencyKm, IntervalValue: 10000}); err != nil {
		t.Fatalf("plan failed: %v", err)
	}
	if _, err := labor.Add(context.Background(), domain.LaborEntry{TenantID: "tenant-a", WorkOrderID: wo.ID, Technician: "Tech", Hours: 2}); err != nil {
		t.Fatalf("labor failed: %v", err)
	}
	down, err := downtime.Start(context.Background(), domain.Downtime{TenantID: "tenant-a", WorkOrderID: wo.ID, Reason: "repair"})
	if err != nil {
		t.Fatalf("downtime failed: %v", err)
	}
	ended, err := downtime.End(context.Background(), "tenant-a", wo.ID, down.ID)
	if err != nil || ended.EndedAt == nil {
		t.Fatalf("end downtime failed: %+v err=%v", ended, err)
	}
}

func TestWorkOrderSoftDelete(t *testing.T) {
	service := application.NewWorkOrderService(infrastructure.NewMemoryStore().WorkOrders(), nil)
	created, _ := service.Create(context.Background(), validWorkOrder("tenant-a"))
	if err := service.Delete(context.Background(), "tenant-a", created.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if _, err := service.FindByID(context.Background(), "tenant-a", created.ID); err == nil {
		t.Fatal("expected deleted work order hidden")
	}
}

func TestWorkOrderValidationAndUniqueCode(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	service := application.NewWorkOrderService(store.WorkOrders(), store.History())
	created, err := service.Create(context.Background(), domain.WorkOrder{TenantID: "tenant-a", Code: "MO-001", Title: "Inspection", Kind: domain.KindInspection, Priority: domain.PriorityLow})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if created.Code != "MO-001" {
		t.Fatalf("expected explicit code preserved, got %q", created.Code)
	}
	if _, err := service.Create(context.Background(), domain.WorkOrder{TenantID: "tenant-a", Code: "MO-001", Title: "Duplicate", Kind: domain.KindCorrective, Priority: domain.PriorityMedium}); err == nil {
		t.Fatal("expected duplicate code validation error")
	}
	if _, err := service.Create(context.Background(), domain.WorkOrder{TenantID: "tenant-a", Title: "Invalid priority", Kind: domain.KindCorrective, Priority: "urgent"}); err == nil {
		t.Fatal("expected invalid priority validation error")
	}
	if _, err := service.Create(context.Background(), domain.WorkOrder{TenantID: "tenant-a", Title: "Invalid kind", Kind: "other", Priority: domain.PriorityMedium}); err == nil {
		t.Fatal("expected invalid kind validation error")
	}
}

func TestWorkOrderCancelUsesPromptStatusAndBlocksTerminalTransition(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	service := application.NewWorkOrderService(store.WorkOrders(), store.History())
	created, _ := service.Create(context.Background(), validWorkOrder("tenant-a"))
	canceled, err := service.Cancel(context.Background(), "tenant-a", created.ID, "actor")
	if err != nil {
		t.Fatalf("cancel failed: %v", err)
	}
	if canceled.Status != domain.WorkOrderCanceled {
		t.Fatalf("expected canceled status, got %q", canceled.Status)
	}
	if canceled.CancelledAt == nil {
		t.Fatal("expected cancelled timestamp populated")
	}
	if _, err := service.Start(context.Background(), "tenant-a", created.ID, "actor"); err == nil {
		t.Fatal("expected terminal status transition error")
	}
}

func TestMaintenanceCatalogWorkshopAndTechnician(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	categories := application.NewCatalogService(store.Categories())
	workshops := application.NewWorkshopService(store.Workshops())
	technicians := application.NewTechnicianService(store.Technicians())
	if _, err := categories.Create(context.Background(), domain.MaintenanceCatalog{TenantID: "tenant-a", Name: "Engine", Code: "engine"}); err != nil {
		t.Fatalf("category failed: %v", err)
	}
	if _, err := workshops.Create(context.Background(), domain.Workshop{TenantID: "tenant-a", Name: "Internal shop"}); err != nil {
		t.Fatalf("workshop failed: %v", err)
	}
	tech, err := technicians.Create(context.Background(), domain.Technician{TenantID: "tenant-a", Name: "Maria"})
	if err != nil {
		t.Fatalf("technician failed: %v", err)
	}
	if !tech.Active {
		t.Fatal("expected technician active by default")
	}
}

func validWorkOrder(tenantID string) domain.WorkOrder {
	return domain.WorkOrder{TenantID: tenantID, Title: "Brake repair", Kind: domain.KindCorrective, Priority: domain.PriorityHigh}
}

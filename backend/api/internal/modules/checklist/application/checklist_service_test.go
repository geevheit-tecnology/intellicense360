package application_test

import (
	"context"
	"testing"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/infrastructure"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/ports"
)

func TestChecklistLifecycle(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	service := application.NewChecklistService(store.Checklists())
	created, err := service.Create(context.Background(), validChecklist("tenant-a", "vehicle-1"))
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	started, err := service.Start(context.Background(), "tenant-a", created.ID, "actor-1")
	if err != nil {
		t.Fatalf("start failed: %v", err)
	}
	if started.Status != domain.ChecklistStatusInProgress || started.StartedAt == nil {
		t.Fatalf("unexpected started checklist: %+v", started)
	}
	finished, err := service.Finish(context.Background(), "tenant-a", created.ID, "actor-1")
	if err != nil {
		t.Fatalf("finish failed: %v", err)
	}
	if finished.Status != domain.ChecklistStatusCompleted || finished.FinishedAt == nil {
		t.Fatalf("unexpected finished checklist: %+v", finished)
	}
}

func TestChecklistTenantIsolationAndPagination(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	service := application.NewChecklistService(store.Checklists())
	_, _ = service.Create(context.Background(), validChecklist("tenant-a", "vehicle-1"))
	_, _ = service.Create(context.Background(), validChecklist("tenant-a", "vehicle-2"))
	_, _ = service.Create(context.Background(), validChecklist("tenant-b", "vehicle-1"))

	page, err := service.Search(context.Background(), "tenant-a", ports.Query{Page: 1, PageSize: 1, Filters: map[string]string{"vehicle_id": "vehicle-1"}})
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if page.TotalItems != 1 || len(page.Items) != 1 || page.Items[0].TenantID != "tenant-a" {
		t.Fatalf("unexpected page: %+v", page)
	}
}

func TestChecklistItemsAnswersAndSoftDelete(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	checklists := application.NewChecklistService(store.Checklists())
	items := application.NewChecklistItemService(store.Checklists(), store.Items())
	answers := application.NewChecklistAnswerService(store.Items(), store.Answers())
	checklist, _ := checklists.Create(context.Background(), validChecklist("tenant-a", "vehicle-1"))

	item, err := items.Add(context.Background(), domain.ChecklistItem{TenantID: "tenant-a", ChecklistID: checklist.ID, Title: "Lights", AnswerType: domain.AnswerTypeBoolean})
	if err != nil {
		t.Fatalf("add item failed: %v", err)
	}
	answer, err := answers.AnswerItem(context.Background(), domain.ChecklistAnswer{TenantID: "tenant-a", ChecklistID: checklist.ID, ChecklistItemID: item.ID, Answer: "true"})
	if err != nil {
		t.Fatalf("answer failed: %v", err)
	}
	if answer.AnsweredAt.IsZero() {
		t.Fatal("expected answered_at")
	}
	if err := checklists.Delete(context.Background(), "tenant-a", checklist.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if _, err := checklists.FindByID(context.Background(), "tenant-a", checklist.ID); err == nil {
		t.Fatal("expected deleted checklist to be hidden")
	}
}

func validChecklist(tenantID string, vehicleID string) domain.Checklist {
	return domain.Checklist{
		TenantID:  tenantID,
		VehicleID: vehicleID,
		Title:     "Daily inspection",
		Type:      "daily",
	}
}

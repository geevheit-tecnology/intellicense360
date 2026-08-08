package application_test

import (
	"context"
	"testing"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/tires/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/tires/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/tires/infrastructure"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/tires/ports"
)

func TestTireCRUDValidationTenantIsolationAndSoftDelete(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	service := application.NewTireService(store.Tires(), store.Movements())
	created, err := service.Create(context.Background(), validTire("tenant-a", "SN1", "F1"))
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if _, err := service.Create(context.Background(), validTire("tenant-a", "SN1", "F2")); err == nil {
		t.Fatal("expected duplicate serial validation")
	}
	if _, err := service.Create(context.Background(), validTire("tenant-b", "SN1", "F1")); err != nil {
		t.Fatalf("same serial/fire must be allowed in another tenant: %v", err)
	}
	page, err := service.Search(context.Background(), "tenant-a", ports.Query{Page: 1, PageSize: 10, Search: "SN1"})
	if err != nil || page.TotalItems != 1 {
		t.Fatalf("unexpected search: page=%+v err=%v", page, err)
	}
	if err := service.Delete(context.Background(), "tenant-a", created.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if _, err := service.FindByID(context.Background(), "tenant-a", created.ID); err == nil {
		t.Fatal("expected soft-deleted tire to be hidden")
	}
}

func TestTireOperationsAndMovements(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	tires := application.NewTireService(store.Tires(), store.Movements())
	movements := application.NewTireMovementService(store.Movements())
	tire, _ := tires.Create(context.Background(), validTire("tenant-a", "SN1", "F1"))

	received, err := tires.Receive(context.Background(), "tenant-a", tire.ID, "receipt", "actor")
	if err != nil || received.Status != domain.TireStatusInStock {
		t.Fatalf("receive failed: %+v err=%v", received, err)
	}
	installed, err := tires.Install(context.Background(), "tenant-a", tire.ID, "vehicle-1", "1L", 1000, "actor")
	if err != nil {
		t.Fatalf("install failed: %v", err)
	}
	if installed.Status != domain.TireStatusInstalled || installed.VehicleID != "vehicle-1" {
		t.Fatalf("unexpected installed tire: %+v", installed)
	}
	rotated, err := tires.Rotate(context.Background(), "tenant-a", tire.ID, "2R", 1200, "actor")
	if err != nil || rotated.Position != "2R" {
		t.Fatalf("rotate failed: %+v err=%v", rotated, err)
	}
	removed, err := tires.Remove(context.Background(), "tenant-a", tire.ID, 1250, "retread candidate", "actor")
	if err != nil || removed.Status != domain.TireStatusRemoved {
		t.Fatalf("remove failed: %+v err=%v", removed, err)
	}
	recapping, err := tires.SendToRecap(context.Background(), "tenant-a", tire.ID, "low tread", "actor")
	if err != nil || recapping.Status != domain.TireStatusUnderRetread {
		t.Fatalf("recap failed: %+v err=%v", recapping, err)
	}
	returned, err := tires.ReturnFromRecap(context.Background(), "tenant-a", tire.ID, "actor")
	if err != nil || returned.RecapCount != 1 || returned.Status != domain.TireStatusInStock {
		t.Fatalf("return failed: %+v err=%v", returned, err)
	}
	disposed, err := tires.Dispose(context.Background(), "tenant-a", tire.ID, "end of life", "actor")
	if err != nil || disposed.Status != domain.TireStatusDisposed {
		t.Fatalf("dispose failed: %+v err=%v", disposed, err)
	}
	if _, err := tires.Remove(context.Background(), "tenant-a", tire.ID, 1300, "", "actor"); err == nil {
		t.Fatal("expected disposed tire movement to be rejected")
	}
	page, err := movements.List(context.Background(), "tenant-a", tire.ID, ports.Query{Page: 1, PageSize: 20})
	if err != nil || page.TotalItems < 4 {
		t.Fatalf("expected operation movements: page=%+v err=%v", page, err)
	}
}

func TestTireInspectionCRUD(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	tires := application.NewTireService(store.Tires(), store.Movements())
	inspections := application.NewTireInspectionService(store.Tires(), store.Inspections())
	tire, _ := tires.Create(context.Background(), validTire("tenant-a", "SN1", "F1"))
	inspection, err := inspections.Register(context.Background(), domain.TireInspection{TenantID: "tenant-a", TireID: tire.ID, TreadMM: 12, Condition: "ok"})
	if err != nil {
		t.Fatalf("inspection failed: %v", err)
	}
	inspection.TreadMM = 11
	updated, err := inspections.Update(context.Background(), inspection)
	if err != nil || updated.TreadMM != 11 {
		t.Fatalf("update inspection failed: %+v err=%v", updated, err)
	}
	if err := inspections.Delete(context.Background(), "tenant-a", tire.ID, inspection.ID); err != nil {
		t.Fatalf("delete inspection failed: %v", err)
	}
	page, _ := inspections.List(context.Background(), "tenant-a", tire.ID, ports.Query{Page: 1, PageSize: 10})
	if page.TotalItems != 0 {
		t.Fatalf("expected deleted inspection hidden: %+v", page)
	}
}

func validTire(tenantID string, serial string, fire string) domain.Tire {
	return domain.Tire{TenantID: tenantID, SerialNumber: serial, FireNumber: fire, Brand: "Michelin", Model: "X", Size: "295/80R22.5", DOT: "DOT1234", CurrentTreadMM: 14, OriginalTreadMM: 16, MinimumTreadMM: 3}
}

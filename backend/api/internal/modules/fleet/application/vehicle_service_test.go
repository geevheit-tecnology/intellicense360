package application_test

import (
	"context"
	"testing"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/infrastructure"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/ports"
)

func TestVehicleCreateValidatesUniqueLicensePlatePerTenant(t *testing.T) {
	service := application.NewVehicleService(infrastructure.NewMemoryStore().Vehicles())
	vehicle := validVehicle("tenant-a", "ABC1D23", "REN1", "CHA1")
	if _, err := service.Create(context.Background(), vehicle); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if _, err := service.Create(context.Background(), validVehicle("tenant-a", "ABC-1D23", "REN2", "CHA2")); err == nil {
		t.Fatal("expected duplicate plate error")
	}
	if _, err := service.Create(context.Background(), validVehicle("tenant-b", "ABC1D23", "REN1", "CHA1")); err != nil {
		t.Fatalf("same identifiers must be allowed in another tenant: %v", err)
	}
}

func TestVehicleSearchFiltersTenantAndPaginates(t *testing.T) {
	service := application.NewVehicleService(infrastructure.NewMemoryStore().Vehicles())
	_, _ = service.Create(context.Background(), validVehicle("tenant-a", "AAA1A11", "REN1", "CHA1"))
	_, _ = service.Create(context.Background(), validVehicle("tenant-a", "BBB2B22", "REN2", "CHA2"))
	_, _ = service.Create(context.Background(), validVehicle("tenant-b", "CCC3C33", "REN3", "CHA3"))

	result, err := service.Search(context.Background(), "tenant-a", ports.Query{Page: 1, PageSize: 1})
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if result.TotalItems != 2 || len(result.Items) != 1 || result.TotalPages != 2 {
		t.Fatalf("unexpected page: %+v", result)
	}
}

func TestVehicleDeleteSoftDeletes(t *testing.T) {
	service := application.NewVehicleService(infrastructure.NewMemoryStore().Vehicles())
	created, _ := service.Create(context.Background(), validVehicle("tenant-a", "AAA1A11", "REN1", "CHA1"))
	if err := service.Delete(context.Background(), "tenant-a", created.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if _, err := service.FindByID(context.Background(), "tenant-a", created.ID); err == nil {
		t.Fatal("expected deleted vehicle to be hidden")
	}
}

func validVehicle(tenantID string, plate string, renavam string, chassis string) domain.Vehicle {
	return domain.Vehicle{
		TenantID:     tenantID,
		CategoryID:   "cat",
		TypeID:       "type",
		BrandID:      "brand",
		ModelID:      "model",
		LicensePlate: domain.LicensePlate(plate),
		Renavam:      domain.Renavam(renavam),
		Chassis:      domain.Chassis(chassis),
		Status:       domain.VehicleStatusActive,
	}
}

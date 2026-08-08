package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/infrastructure"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/ports"
)

func TestInventoryPartValidationTenantIsolationAndSearch(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	units := application.NewCatalogService(store.Units())
	parts := application.NewPartService(store.Parts(), store.Units())
	unit, err := units.Create(context.Background(), domain.Catalog{TenantID: "tenant-a", Name: "Unit", Code: "UN"})
	if err != nil {
		t.Fatalf("unit failed: %v", err)
	}
	created, err := parts.Create(context.Background(), validPart("tenant-a", unit.Code))
	if err != nil {
		t.Fatalf("part failed: %v", err)
	}
	if _, err := parts.Create(context.Background(), validPart("tenant-a", unit.Code)); !errors.Is(err, application.ErrSKUTaken) {
		t.Fatalf("expected sku taken, got %v", err)
	}
	_, _ = parts.Create(context.Background(), domain.Part{TenantID: "tenant-b", SKU: "SKU-1", InternalCode: "INT-1", Name: "Tenant B", UnitID: unit.Code})
	page, err := parts.Search(context.Background(), "tenant-a", ports.Query{Search: "Filter", Page: 1, PageSize: 10})
	if err != nil || page.TotalItems != 1 || page.Items[0].ID != created.ID {
		t.Fatalf("tenant search failed: %+v err=%v", page, err)
	}
}

func TestInventoryPartUnitAndStockRangeValidation(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	parts := application.NewPartService(store.Parts(), store.Units())
	if _, err := parts.Create(context.Background(), validPart("tenant-a", "UN")); !errors.Is(err, application.ErrInvalidUnit) {
		t.Fatalf("expected invalid unit, got %v", err)
	}
	units := application.NewCatalogService(store.Units())
	_, _ = units.Create(context.Background(), domain.Catalog{TenantID: "tenant-a", Name: "Unit", Code: "UN"})
	part := validPart("tenant-a", "UN")
	part.Metadata = map[string]string{"minimum_stock": "10", "maximum_stock": "5"}
	if _, err := parts.Create(context.Background(), part); !errors.Is(err, application.ErrInvalidStockRange) {
		t.Fatalf("expected invalid stock range, got %v", err)
	}
}

func TestInventoryWarehouseLocationAndSoftDelete(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	warehouses := application.NewWarehouseService(store.Warehouses())
	locations := application.NewLocationService(store.Locations())
	warehouse, err := warehouses.Create(context.Background(), domain.Warehouse{TenantID: "tenant-a", Name: "Main", Code: "MAIN"})
	if err != nil {
		t.Fatalf("warehouse failed: %v", err)
	}
	location, err := locations.Create(context.Background(), domain.WarehouseLocation{TenantID: "tenant-a", WarehouseID: warehouse.ID, Name: "Aisle A", Code: "A"})
	if err != nil {
		t.Fatalf("location failed: %v", err)
	}
	if err := locations.Delete(context.Background(), "tenant-a", location.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	page, _ := locations.Search(context.Background(), "tenant-a", ports.Query{Page: 1, PageSize: 10})
	if page.TotalItems != 0 {
		t.Fatalf("expected deleted location hidden, got %+v", page)
	}
}

func validPart(tenantID string, unit string) domain.Part {
	return domain.Part{TenantID: tenantID, SKU: "SKU-1", InternalCode: "INT-1", Name: "Filter element", UnitID: unit}
}

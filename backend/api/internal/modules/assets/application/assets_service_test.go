package application_test

import (
	"context"
	"testing"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/assets/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/assets/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/assets/infrastructure"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/assets/ports"
)

func TestAssetValidationUniquenessTenantIsolationAndSoftDelete(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	service := application.NewAssetService(store.Assets(), store.Audit())
	created, err := service.Create(context.Background(), validAsset("tenant-a", "A-001", "SN-001", "TAG-001"))
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if _, err := service.Create(context.Background(), validAsset("tenant-a", "A-001", "SN-002", "TAG-002")); err == nil {
		t.Fatal("expected duplicate internal code")
	}
	if _, err := service.Create(context.Background(), validAsset("tenant-b", "A-001", "SN-001", "TAG-001")); err != nil {
		t.Fatalf("same identifiers must be allowed in another tenant: %v", err)
	}
	page, err := service.Search(context.Background(), "tenant-a", ports.Query{Search: "A-001", Page: 1, PageSize: 1})
	if err != nil || page.TotalItems != 1 || len(page.Items) != 1 {
		t.Fatalf("unexpected search page=%+v err=%v", page, err)
	}
	if err := service.Delete(context.Background(), "tenant-a", created.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if _, err := service.FindByID(context.Background(), "tenant-a", created.ID); err == nil {
		t.Fatal("expected soft deleted asset to be hidden")
	}
}

func TestCatalogAndEquipmentServices(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	categories := application.NewCategoryService(store.Categories())
	types := application.NewTypeService(store.Types())
	manufacturers := application.NewManufacturerService(store.Manufacturers())
	models := application.NewModelService(store.Models())
	equipment := application.NewEquipmentService(store.Equipment())

	category, err := categories.Create(context.Background(), domain.AssetCategory{TenantID: "tenant-a", Name: "Truck", Code: "truck"})
	if err != nil {
		t.Fatalf("category failed: %v", err)
	}
	assetType, err := types.Create(context.Background(), domain.AssetType{TenantID: "tenant-a", CategoryID: category.ID, Name: "Semi Trailer", Code: "semi_trailer"})
	if err != nil {
		t.Fatalf("type failed: %v", err)
	}
	manufacturer, err := manufacturers.Create(context.Background(), domain.Manufacturer{TenantID: "tenant-a", Name: "Volvo"})
	if err != nil {
		t.Fatalf("manufacturer failed: %v", err)
	}
	if _, err := models.Create(context.Background(), domain.Model{TenantID: "tenant-a", ManufacturerID: manufacturer.ID, Name: "FH"}); err != nil {
		t.Fatalf("model failed: %v", err)
	}
	if _, err := equipment.Create(context.Background(), domain.Equipment{TenantID: "tenant-a", AssetID: "asset-1", Category: category.Code, Type: assetType.Code}); err != nil {
		t.Fatalf("equipment failed: %v", err)
	}
}

func TestInvalidStatusRejected(t *testing.T) {
	service := application.NewAssetService(infrastructure.NewMemoryStore().Assets(), nil)
	asset := validAsset("tenant-a", "A-001", "SN-001", "TAG-001")
	asset.Status = "invalid"
	if _, err := service.Create(context.Background(), asset); err == nil {
		t.Fatal("expected invalid status")
	}
}

func validAsset(tenantID string, internalCode string, serial string, tag string) domain.Asset {
	return domain.Asset{
		TenantID:     tenantID,
		InternalCode: internalCode,
		SerialNumber: serial,
		AssetTag:     tag,
		Name:         "Asset " + internalCode,
		Status:       domain.AssetStatusAvailable,
		Ownership:    domain.OwnershipOwned,
	}
}

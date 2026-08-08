package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/infrastructure"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/ports"
)

func TestSupplierValidationUniqueCNPJAndTenantIsolation(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	service := application.NewSupplierService(store.Suppliers(), store.Categories(), store.Audit())
	created, err := service.Create(context.Background(), validSupplier("tenant-a"))
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if _, err := service.Create(context.Background(), validSupplier("tenant-a")); !errors.Is(err, application.ErrCNPJTaken) {
		t.Fatalf("expected duplicate cnpj, got %v", err)
	}
	_, err = service.Create(context.Background(), validSupplier("tenant-b"))
	if err != nil {
		t.Fatalf("same cnpj in other tenant should pass: %v", err)
	}
	page, err := service.Search(context.Background(), "tenant-a", ports.Query{Search: "Alpha", Page: 1, PageSize: 10})
	if err != nil || page.TotalItems != 1 || page.Items[0].ID != created.ID {
		t.Fatalf("tenant isolation failed: %+v err=%v", page, err)
	}
}

func TestSupplierFormatValidationAndSoftDelete(t *testing.T) {
	store := infrastructure.NewMemoryStore()
	service := application.NewSupplierService(store.Suppliers(), store.Categories(), nil)
	if _, err := service.Create(context.Background(), domain.Supplier{TenantID: "tenant-a", LegalName: "Invalid", CNPJ: "123", Type: domain.TypeOther}); !errors.Is(err, application.ErrInvalidCNPJ) {
		t.Fatalf("expected invalid cnpj, got %v", err)
	}
	if _, err := service.Create(context.Background(), domain.Supplier{TenantID: "tenant-a", LegalName: "Invalid", Email: "bad", Type: domain.TypeOther}); !errors.Is(err, application.ErrInvalidEmail) {
		t.Fatalf("expected invalid email, got %v", err)
	}
	created, _ := service.Create(context.Background(), validSupplier("tenant-a"))
	if err := service.Delete(context.Background(), "tenant-a", created.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if _, err := service.FindByID(context.Background(), "tenant-a", created.ID); !errors.Is(err, application.ErrNotFound) {
		t.Fatalf("expected soft deleted hidden, got %v", err)
	}
}

func validSupplier(tenantID string) domain.Supplier {
	return domain.Supplier{TenantID: tenantID, LegalName: "Alpha Parts LTDA", TradeName: "Alpha", CNPJ: "11222333000181", Email: "alpha@example.com", Status: domain.StatusActive, Type: domain.TypePartsSupplier}
}

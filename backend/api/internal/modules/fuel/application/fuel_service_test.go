package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/infrastructure"
)

func TestFuelTransactionLifecycleRequiresAdjustmentAfterCompleted(t *testing.T) {
	ctx := context.Background()
	store := infrastructure.NewMemoryStore()
	service := application.NewFuelTransactionService(store.Transactions(), store.Adjustments())

	created, err := service.Create(ctx, validTransaction("tenant-a"))
	if err != nil {
		t.Fatalf("create transaction: %v", err)
	}
	completed, err := service.Complete(ctx, "tenant-a", created.ID, "user-a")
	if err != nil {
		t.Fatalf("complete transaction: %v", err)
	}
	if completed.Status != domain.FuelTransactionCompleted {
		t.Fatalf("expected completed status, got %s", completed.Status)
	}

	completed.Notes = "silent change"
	_, err = service.Update(ctx, completed)
	if !errors.Is(err, application.ErrCompletedImmutable) {
		t.Fatalf("expected immutable completed transaction error, got %v", err)
	}

	adjustment, err := service.Adjust(ctx, domain.FuelAdjustment{
		TenantID:            "tenant-a",
		TransactionID:       created.ID,
		AdjustmentType:      "receipt_correction",
		Reason:              "receipt corrected by station",
		AdjustedQuantity:    50,
		AdjustedUnitPrice:   6,
		AdjustedTotalAmount: 300,
		OriginalReference:   created.ReceiptNumber,
		CreatedBy:           "user-a",
	})
	if err != nil {
		t.Fatalf("adjust transaction: %v", err)
	}
	if adjustment.ID == "" {
		t.Fatal("expected adjustment id")
	}
	adjusted, err := service.FindByID(ctx, "tenant-a", created.ID)
	if err != nil {
		t.Fatalf("find adjusted transaction: %v", err)
	}
	if adjusted.Status != domain.FuelTransactionAdjusted {
		t.Fatalf("expected adjusted status, got %s", adjusted.Status)
	}
}

func TestFuelTransactionTenantIsolation(t *testing.T) {
	ctx := context.Background()
	store := infrastructure.NewMemoryStore()
	service := application.NewFuelTransactionService(store.Transactions(), store.Adjustments())

	created, err := service.Create(ctx, validTransaction("tenant-a"))
	if err != nil {
		t.Fatalf("create transaction: %v", err)
	}
	_, err = service.FindByID(ctx, "tenant-b", created.ID)
	if !errors.Is(err, application.ErrNotFound) {
		t.Fatalf("expected not found across tenant boundary, got %v", err)
	}
}

func TestFuelValidation(t *testing.T) {
	ctx := context.Background()
	store := infrastructure.NewMemoryStore()
	service := application.NewFuelTransactionService(store.Transactions(), store.Adjustments())

	item := validTransaction("tenant-a")
	item.Quantity = 0
	_, err := service.Create(ctx, item)
	if !errors.Is(err, application.ErrInvalidFuelData) {
		t.Fatalf("expected invalid quantity error, got %v", err)
	}

	item = validTransaction("tenant-a")
	item.TotalAmount = 12
	_, err = service.Create(ctx, item)
	if !errors.Is(err, application.ErrInvalidFuelData) {
		t.Fatalf("expected invalid total error, got %v", err)
	}
}

func TestFuelStationCNPJUniquePerTenant(t *testing.T) {
	ctx := context.Background()
	store := infrastructure.NewMemoryStore()
	service := application.NewFuelStationService(store.Stations())

	_, err := service.Create(ctx, domain.FuelStation{TenantID: "tenant-a", Name: "Posto A", CNPJ: "12.345.678/0001-99"})
	if err != nil {
		t.Fatalf("create station: %v", err)
	}
	_, err = service.Create(ctx, domain.FuelStation{TenantID: "tenant-a", Name: "Posto B", CNPJ: "12345678000199"})
	if !errors.Is(err, application.ErrDuplicateCNPJ) {
		t.Fatalf("expected duplicate cnpj error, got %v", err)
	}
	_, err = service.Create(ctx, domain.FuelStation{TenantID: "tenant-b", Name: "Posto B", CNPJ: "12345678000199"})
	if err != nil {
		t.Fatalf("same cnpj must be allowed for another tenant: %v", err)
	}
}

func validTransaction(tenantID string) domain.FuelTransaction {
	return domain.FuelTransaction{
		TenantID:         tenantID,
		FuelKind:         domain.FuelKindDieselS10,
		Quantity:         50,
		UnitPrice:        6,
		TotalAmount:      300,
		OdometerReading:  12000,
		ReceiptNumber:    "NF-100",
		DriverReference:  "driver-1",
		VehicleReference: "vehicle-1",
		PaymentMethod:    "fuel_card",
	}
}

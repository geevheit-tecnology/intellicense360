package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/financial/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/financial/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/financial/infrastructure"
)

func TestExpenseLifecycleAndImmutability(t *testing.T) {
	ctx := context.Background()
	store := infrastructure.NewMemoryStore()
	service := application.NewTransactionService(store.Transactions(), store.Periods(), store.Adjustments())

	expense, err := service.CreateExpense(ctx, validTransaction("tenant-a"))
	if err != nil {
		t.Fatal(err)
	}
	expense, err = service.Approve(ctx, "tenant-a", expense.ID, "user-a")
	if err != nil {
		t.Fatal(err)
	}
	expense, err = service.Pay(ctx, "tenant-a", expense.ID, "user-a")
	if err != nil {
		t.Fatal(err)
	}
	if expense.Status != domain.StatusPaid {
		t.Fatalf("expected paid, got %s", expense.Status)
	}

	expense.Description = "silent mutation"
	_, err = service.Update(ctx, expense)
	if !errors.Is(err, application.ErrFinalizedImmutable) {
		t.Fatalf("expected immutable error, got %v", err)
	}

	adjustment, err := service.Adjust(ctx, domain.FinancialAdjustment{TenantID: "tenant-a", TransactionID: expense.ID, AdjustmentType: "correction", Reason: "document correction", AdjustedAmount: 100, CreatedBy: "user-a"})
	if err != nil {
		t.Fatal(err)
	}
	if adjustment.ID == "" {
		t.Fatal("expected adjustment id")
	}
}

func TestRevenueLifecycle(t *testing.T) {
	ctx := context.Background()
	store := infrastructure.NewMemoryStore()
	service := application.NewTransactionService(store.Transactions(), store.Periods(), store.Adjustments())

	revenue, err := service.CreateRevenue(ctx, validTransaction("tenant-a"))
	if err != nil {
		t.Fatal(err)
	}
	revenue, err = service.Approve(ctx, "tenant-a", revenue.ID, "user-a")
	if err != nil {
		t.Fatal(err)
	}
	revenue, err = service.Receive(ctx, "tenant-a", revenue.ID, "user-a")
	if err != nil {
		t.Fatal(err)
	}
	if revenue.Status != domain.StatusReceived {
		t.Fatalf("expected received, got %s", revenue.Status)
	}
}

func TestClosedPeriodBlocksTransactionMutation(t *testing.T) {
	ctx := context.Background()
	store := infrastructure.NewMemoryStore()
	periods := application.NewPeriodService(store.Periods())
	service := application.NewTransactionService(store.Transactions(), store.Periods(), store.Adjustments())

	period, err := periods.Create(ctx, domain.FinancialPeriod{TenantID: "tenant-a", Year: 2026, Month: 8, StartDate: mustDate("2026-08-01"), EndDate: mustDate("2026-08-31")})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := periods.Close(ctx, "tenant-a", period.ID); err != nil {
		t.Fatal(err)
	}
	_, err = service.CreateExpense(ctx, domain.FinancialTransaction{TenantID: "tenant-a", Description: "Closed period expense", Amount: 10, Date: mustDate("2026-08-10")})
	if !errors.Is(err, application.ErrClosedPeriodImmutable) {
		t.Fatalf("expected closed period error, got %v", err)
	}
}

func TestTenantIsolation(t *testing.T) {
	ctx := context.Background()
	store := infrastructure.NewMemoryStore()
	service := application.NewTransactionService(store.Transactions(), store.Periods(), store.Adjustments())
	expense, err := service.CreateExpense(ctx, validTransaction("tenant-a"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = service.FindByID(ctx, "tenant-b", expense.ID)
	if !errors.Is(err, application.ErrNotFound) {
		t.Fatalf("expected not found across tenant, got %v", err)
	}
}

func validTransaction(tenantID string) domain.FinancialTransaction {
	return domain.FinancialTransaction{TenantID: tenantID, Description: "Operational expense", Amount: 100, Date: mustDate("2026-09-10")}
}
func mustDate(value string) time.Time {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		panic(err)
	}
	return parsed
}

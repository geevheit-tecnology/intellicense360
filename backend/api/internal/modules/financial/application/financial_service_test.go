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

func TestCashFlowAndOperationalReconciliation(t *testing.T) {
	ctx := context.Background()
	store := infrastructure.NewMemoryStore()
	txService := application.NewTransactionService(store.Transactions(), store.Periods(), store.Adjustments())
	cashFlowService := application.NewCashFlowService(store.Transactions(), store.Closings())

	// 1. Reconcile Fleet Fuel Expense with requisition and invoice
	due := mustDate("2026-09-25")
	fuelTx, err := cashFlowService.ReconcileOperationalExpense(ctx, "tenant-a", "user-admin", domain.OperationalExpenseReconciliation{
		SourceType:        "fuel",
		SourceID:          "fuel_001",
		VehicleID:         "veh_001",
		DriverID:          "drv_001",
		InvoiceNumber:     "NF-10492",
		InvoiceKey:        "35260900000000000000550010000104921000104925",
		RequisitionNumber: "REQ-2026-0081",
		IsThirdParty:      false,
		Amount:            850.50,
		Date:              mustDate("2026-09-10"),
		DueDate:           &due,
		PaymentMethod:     "boleto",
		Barcode:           "341917900101043510047910201500081890000085050",
		Description:       "Diesel S10 - Frota Própria",
	})
	if err != nil {
		t.Fatalf("unexpected error reconciling fuel: %v", err)
	}
	if fuelTx.Status != domain.StatusPending {
		t.Fatalf("expected pending status, got %s", fuelTx.Status)
	}
	if fuelTx.RequisitionNumber != "REQ-2026-0081" {
		t.Fatalf("expected requisition REQ-2026-0081, got %s", fuelTx.RequisitionNumber)
	}

	// 2. Reconcile Third-Party Maintenance Expense
	maintTx, err := cashFlowService.ReconcileOperationalExpense(ctx, "tenant-a", "user-admin", domain.OperationalExpenseReconciliation{
		SourceType:           "maintenance",
		SourceID:             "wo_002",
		InvoiceNumber:        "NFS-5512",
		IsThirdParty:         true,
		ThirdPartyIdentifier: "ABC-1234 (Agregado Silva)",
		Amount:               1200.00,
		Date:                 mustDate("2026-09-10"),
		PaymentMethod:        "pix",
		Description:          "Troca de lona de freio - Terceiro",
	})
	if err != nil {
		t.Fatalf("unexpected error reconciling maintenance: %v", err)
	}
	if !maintTx.IsThirdParty {
		t.Fatal("expected is_third_party to be true")
	}

	// 3. Create a Revenue to balance cash flow
	rev, err := txService.CreateRevenue(ctx, domain.FinancialTransaction{
		TenantID:    "tenant-a",
		Description: "Frete Concluído - CTE 9941",
		Amount:      5000.00,
		Date:        mustDate("2026-09-10"),
	})
	if err != nil {
		t.Fatal(err)
	}
	rev, err = txService.Approve(ctx, "tenant-a", rev.ID, "user-admin")
	if err != nil {
		t.Fatal(err)
	}
	_, err = txService.Receive(ctx, "tenant-a", rev.ID, "user-admin")
	if err != nil {
		t.Fatal(err)
	}

	// 4. Pay the fleet fuel expense
	fuelTx, err = txService.Approve(ctx, "tenant-a", fuelTx.ID, "user-admin")
	if err != nil {
		t.Fatal(err)
	}
	fuelTx, err = txService.Pay(ctx, "tenant-a", fuelTx.ID, "user-admin")
	if err != nil {
		t.Fatal(err)
	}

	// 5. Query Cash Flow Summary
	summary, err := cashFlowService.GetSummary(ctx, "tenant-a", "2026-09-01", "2026-09-30")
	if err != nil {
		t.Fatalf("unexpected error getting cash flow summary: %v", err)
	}
	if summary.TotalRevenues != 5000.00 {
		t.Fatalf("expected 5000 total revenues, got %.2f", summary.TotalRevenues)
	}
	if summary.TotalFleetExpenses != 850.50 {
		t.Fatalf("expected 850.50 fleet expenses, got %.2f", summary.TotalFleetExpenses)
	}
	if summary.PendingPayables != 1200.00 {
		t.Fatalf("expected 1200 pending payables (third-party maintenance), got %.2f", summary.PendingPayables)
	}
	if summary.NetCashFlow != 4149.50 { // 5000 - 850.50
		t.Fatalf("expected 4149.50 net cash flow, got %.2f", summary.NetCashFlow)
	}

	// 6. Close Cash for 2026-09-10
	closing, err := cashFlowService.CloseCash(ctx, "tenant-a", "user-admin", domain.CashClosing{
		ClosingDate:    mustDate("2026-09-10"),
		OpeningBalance: 1000.00,
		Notes:          "Fechamento regular do dia",
	})
	if err != nil {
		t.Fatalf("unexpected error closing cash: %v", err)
	}
	if closing.Status != "closed" {
		t.Fatalf("expected closed status, got %s", closing.Status)
	}
	// opening (1000) + rev (5000) - exp (850.50) = 5149.50
	if closing.ClosingBalance != 5149.50 {
		t.Fatalf("expected 5149.50 closing balance, got %.2f", closing.ClosingBalance)
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

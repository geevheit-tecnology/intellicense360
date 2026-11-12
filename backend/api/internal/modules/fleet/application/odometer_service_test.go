package application_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/infrastructure"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/ports"
	"github.com/google/uuid"
)

func setupOdometerTest() (application.OdometerService, ports.VehicleService, string, string, string) {
	store := infrastructure.NewMemoryStore()
	vehicles := application.NewVehicleService(store.Vehicles())
	odometer := application.NewOdometerService(store.Odometer(), store.Vehicles())

	tenantID := uuid.NewString()
	actorID := uuid.NewString()

	v, err := vehicles.Create(context.Background(), domain.Vehicle{
		TenantID:     tenantID,
		CategoryID:   uuid.NewString(),
		TypeID:       uuid.NewString(),
		BrandID:      uuid.NewString(),
		ModelID:      uuid.NewString(),
		LicensePlate: "ABC1D23",
		Renavam:      "12345678901",
		Chassis:      "9BWZZZ377VT004251",
	})
	if err != nil {
		panic(err)
	}

	return odometer, vehicles, tenantID, v.ID, actorID
}

func TestOdometerReadingRecordSuccess(t *testing.T) {
	odometer, _, tenantID, vehicleID, actorID := setupOdometerTest()
	ctx := context.Background()

	// Initial reading: 10,000 km
	first, err := odometer.RecordReading(ctx, domain.OdometerReading{
		TenantID:  tenantID,
		VehicleID: vehicleID,
		ActorID:   actorID,
		Reading:   10000,
		Source:    "manual",
	})
	if err != nil {
		t.Fatalf("expected first reading to succeed: %v", err)
	}
	if first.Reading != 10000 || first.ID == "" {
		t.Fatalf("unexpected reading payload: %+v", first)
	}

	// Subsequent reading: 10,500 km
	second, err := odometer.RecordReading(ctx, domain.OdometerReading{
		TenantID:  tenantID,
		VehicleID: vehicleID,
		ActorID:   actorID,
		Reading:   10500,
		Source:    "fueling",
	})
	if err != nil {
		t.Fatalf("expected progressive reading to succeed: %v", err)
	}
	if second.Reading != 10500 {
		t.Fatalf("expected 10500, got %d", second.Reading)
	}
}

func TestOdometerReadingAntiRegressionRejectsLowerKM(t *testing.T) {
	odometer, _, tenantID, vehicleID, actorID := setupOdometerTest()
	ctx := context.Background()

	_, _ = odometer.RecordReading(ctx, domain.OdometerReading{
		TenantID:  tenantID,
		VehicleID: vehicleID,
		ActorID:   actorID,
		Reading:   50000,
		Source:    "manual",
	})

	// Lower reading without correction flag: regression error
	_, err := odometer.RecordReading(ctx, domain.OdometerReading{
		TenantID:  tenantID,
		VehicleID: vehicleID,
		ActorID:   actorID,
		Reading:   48000,
		Source:    "manual",
	})
	if !errors.Is(err, application.ErrOdometerRegression) {
		t.Fatalf("expected ErrOdometerRegression, got %v", err)
	}
}

func TestOdometerReadingCorrectionAcceptsLowerKMWithReason(t *testing.T) {
	odometer, _, tenantID, vehicleID, actorID := setupOdometerTest()
	ctx := context.Background()

	_, _ = odometer.RecordReading(ctx, domain.OdometerReading{
		TenantID:  tenantID,
		VehicleID: vehicleID,
		ActorID:   actorID,
		Reading:   80000,
		Source:    "manual",
	})

	// Valid correction: cluster replacement with >= 10 chars explanation
	corrected, err := odometer.RecordReading(ctx, domain.OdometerReading{
		TenantID:         tenantID,
		VehicleID:        vehicleID,
		ActorID:          actorID,
		Reading:          1200,
		Source:           "correction",
		IsCorrection:     true,
		CorrectionReason: "Substituição do painel de instrumentos por defeito elétrico",
	})
	if err != nil {
		t.Fatalf("expected justified correction to succeed: %v", err)
	}
	if corrected.Reading != 1200 || !corrected.IsCorrection {
		t.Fatalf("unexpected corrected reading: %+v", corrected)
	}
}

func TestOdometerReadingCorrectionRejectsShortReason(t *testing.T) {
	odometer, _, tenantID, vehicleID, actorID := setupOdometerTest()
	ctx := context.Background()

	_, err := odometer.RecordReading(ctx, domain.OdometerReading{
		TenantID:         tenantID,
		VehicleID:        vehicleID,
		ActorID:          actorID,
		Reading:          1000,
		Source:           "correction",
		IsCorrection:     true,
		CorrectionReason: "erro", // < 10 chars
	})
	if !errors.Is(err, application.ErrOdometerCorrectionRequired) {
		t.Fatalf("expected ErrOdometerCorrectionRequired, got %v", err)
	}
}

func TestOdometerReadingRejectsUnknownVehicle(t *testing.T) {
	odometer, _, tenantID, _, actorID := setupOdometerTest()
	ctx := context.Background()

	_, err := odometer.RecordReading(ctx, domain.OdometerReading{
		TenantID:  tenantID,
		VehicleID: uuid.NewString(), // unknown
		ActorID:   actorID,
		Reading:   1000,
		Source:    "manual",
	})
	if !errors.Is(err, application.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestOdometerReadingHistoryAndSummary(t *testing.T) {
	odometer, _, tenantID, vehicleID, actorID := setupOdometerTest()
	ctx := context.Background()

	t1 := time.Now().UTC().Add(-2 * time.Hour)
	t2 := time.Now().UTC().Add(-1 * time.Hour)

	_, _ = odometer.RecordReading(ctx, domain.OdometerReading{
		TenantID:  tenantID,
		VehicleID: vehicleID,
		ActorID:   actorID,
		Reading:   20000,
		ReadAt:    t1,
		Source:    "manual",
	})

	_, _ = odometer.RecordReading(ctx, domain.OdometerReading{
		TenantID:  tenantID,
		VehicleID: vehicleID,
		ActorID:   actorID,
		Reading:   20300,
		ReadAt:    t2,
		Source:    "fueling",
	})

	summary, err := odometer.GetSummary(ctx, tenantID, vehicleID)
	if err != nil {
		t.Fatalf("summary failed: %v", err)
	}
	if summary.CurrentKM != 20300 {
		t.Fatalf("expected current KM 20300, got %d", summary.CurrentKM)
	}

	history, err := odometer.ListHistory(ctx, tenantID, vehicleID, ports.Query{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("history failed: %v", err)
	}
	if history.TotalItems != 2 || len(history.Items) != 2 {
		t.Fatalf("unexpected history count: %+v", history)
	}
	if history.Items[0].Reading != 20300 || history.Items[1].Reading != 20000 {
		t.Fatalf("expected descending order by read_at, got: %d then %d", history.Items[0].Reading, history.Items[1].Reading)
	}
}

func TestOdometerReadingTenantIsolation(t *testing.T) {
	odometer, _, tenantA, vehicleA, actorID := setupOdometerTest()
	ctx := context.Background()
	tenantB := uuid.NewString()

	_, _ = odometer.RecordReading(ctx, domain.OdometerReading{
		TenantID:  tenantA,
		VehicleID: vehicleA,
		ActorID:   actorID,
		Reading:   50000,
		Source:    "manual",
	})

	// Tenant B cannot access Tenant A's vehicle odometer
	if _, err := odometer.GetSummary(ctx, tenantB, vehicleA); !errors.Is(err, application.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for cross-tenant access, got %v", err)
	}

	// Tenant B cannot list Tenant A's history
	if _, err := odometer.ListHistory(ctx, tenantB, vehicleA, ports.Query{}); !errors.Is(err, application.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for cross-tenant history, got %v", err)
	}
}


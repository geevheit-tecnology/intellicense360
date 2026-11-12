package application

import (
	"context"
	"testing"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/drivers/domain"
)

type memoryJourneyRepo struct {
	records []domain.DriverJourneyRecord
}

func (r *memoryJourneyRepo) Create(_ context.Context, item domain.DriverJourneyRecord) (domain.DriverJourneyRecord, error) {
	r.records = append(r.records, item)
	return item, nil
}

func (r *memoryJourneyRepo) FindByID(_ context.Context, tenantID, id string) (domain.DriverJourneyRecord, error) {
	for _, rec := range r.records {
		if rec.TenantID == tenantID && rec.ID == id {
			return rec, nil
		}
	}
	return domain.DriverJourneyRecord{}, ErrNotFound
}

func (r *memoryJourneyRepo) ListByDriver(_ context.Context, tenantID, driverID string, _ int) ([]domain.DriverJourneyRecord, error) {
	var res []domain.DriverJourneyRecord
	for _, rec := range r.records {
		if rec.TenantID == tenantID && rec.DriverID == driverID {
			res = append(res, rec)
		}
	}
	return res, nil
}

func (r *memoryJourneyRepo) GetLastActiveRecord(_ context.Context, tenantID, driverID string) (*domain.DriverJourneyRecord, error) {
	for i := len(r.records) - 1; i >= 0; i-- {
		if r.records[i].TenantID == tenantID && r.records[i].DriverID == driverID && r.records[i].EndAt == nil {
			return &r.records[i], nil
		}
	}
	return nil, nil
}

func (r *memoryJourneyRepo) CloseActiveRecord(_ context.Context, tenantID, id string, endAt time.Time, durationMinutes int, isViolation bool, violationReason string) error {
	for i := range r.records {
		if r.records[i].TenantID == tenantID && r.records[i].ID == id {
			r.records[i].EndAt = &endAt
			r.records[i].DurationMinutes = durationMinutes
			r.records[i].IsViolation = isViolation
			r.records[i].ViolationReason = violationReason
			return nil
		}
	}
	return ErrNotFound
}

func (r *memoryJourneyRepo) GetSummary(_ context.Context, tenantID, driverID string) (domain.JourneySummary, error) {
	var s domain.JourneySummary
	for _, rec := range r.records {
		if rec.TenantID == tenantID && rec.DriverID == driverID {
			if rec.EventType == domain.EventDriving {
				s.TotalDrivingMinutes += rec.DurationMinutes
			} else if rec.EventType == domain.EventBreakRest || rec.EventType == domain.EventMeal || rec.EventType == domain.EventDailyRest {
				s.TotalRestMinutes += rec.DurationMinutes
			} else if rec.EventType == domain.EventWaiting {
				s.TotalWaitingMinutes += rec.DurationMinutes
			}
			if rec.IsViolation {
				s.ViolationsCount++
			}
			s.TotalAllowanceBRL += rec.DailyAllowanceAmount
		}
	}
	return s, nil
}

func TestJourneyService_ContinuousDrivingViolation(t *testing.T) {
	repo := &memoryJourneyRepo{}
	svc := NewJourneyService(repo)
	ctx := context.Background()
	tenantID := "tenant-1"
	driverID := "driver-1"

	// 1. Driver starts driving at 08:00
	startDriving := time.Date(2026, 11, 12, 8, 0, 0, 0, time.UTC)
	_, err := svc.RecordEvent(ctx, domain.DriverJourneyRecord{
		TenantID:  tenantID,
		DriverID:  driverID,
		EventType: domain.EventDriving,
		StartAt:   startDriving,
		Location:  "São Paulo - SP",
	})
	if err != nil {
		t.Fatalf("failed to record driving start: %v", err)
	}

	// 2. Driver stops for meal at 14:00 (6 hours = 360 mins later without break! Exceeds 5h30 limit)
	startMeal := time.Date(2026, 11, 12, 14, 0, 0, 0, time.UTC)
	_, err = svc.RecordEvent(ctx, domain.DriverJourneyRecord{
		TenantID:  tenantID,
		DriverID:  driverID,
		EventType: domain.EventMeal,
		StartAt:   startMeal,
		Location:  "Posto Graal - Registro/SP",
	})
	if err != nil {
		t.Fatalf("failed to record meal stop: %v", err)
	}

	// 3. Verify that the previous driving record was automatically closed with violation
	records, err := svc.ListJourney(ctx, tenantID, driverID, 10)
	if err != nil {
		t.Fatalf("failed to list records: %v", err)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}

	drivingRec := records[0]
	if !drivingRec.IsViolation {
		t.Fatalf("expected violation on continuous driving > 5h30, got false")
	}
	if drivingRec.DurationMinutes != 360 {
		t.Fatalf("expected 360 minutes duration, got %d", drivingRec.DurationMinutes)
	}

	// 4. Verify summary
	summary, err := svc.GetSummary(ctx, tenantID, driverID)
	if err != nil {
		t.Fatalf("failed to get summary: %v", err)
	}
	if summary.ViolationsCount != 1 {
		t.Fatalf("expected 1 violation, got %d", summary.ViolationsCount)
	}
	if summary.TotalDrivingMinutes != 360 {
		t.Fatalf("expected 360 total driving minutes, got %d", summary.TotalDrivingMinutes)
	}
	if summary.TotalAllowanceBRL != 45.00 {
		t.Fatalf("expected 45.00 meal allowance, got %f", summary.TotalAllowanceBRL)
	}
}


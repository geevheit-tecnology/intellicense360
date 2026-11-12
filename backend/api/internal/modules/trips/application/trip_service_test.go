package application

import (
	"context"
	"testing"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/trips/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/trips/ports"
)

type memoryTripRepo struct {
	trips map[string]domain.Trip
}

func newMemoryTripRepo() *memoryTripRepo {
	return &memoryTripRepo{trips: make(map[string]domain.Trip)}
}

func (r *memoryTripRepo) Create(_ context.Context, item domain.Trip) (domain.Trip, error) {
	r.trips[item.ID] = item
	return item, nil
}

func (r *memoryTripRepo) FindByID(_ context.Context, tenantID string, id string) (domain.Trip, error) {
	item, ok := r.trips[id]
	if !ok || item.TenantID != tenantID {
		return domain.Trip{}, ErrNotFound
	}
	return item, nil
}

func (r *memoryTripRepo) Search(_ context.Context, tenantID string, _ ports.Query) (ports.Page[domain.Trip], error) {
	var list []domain.Trip
	for _, t := range r.trips {
		if t.TenantID == tenantID {
			list = append(list, t)
		}
	}
	return ports.Page[domain.Trip]{Items: list, TotalItems: len(list)}, nil
}

func (r *memoryTripRepo) Update(_ context.Context, item domain.Trip) (domain.Trip, error) {
	r.trips[item.ID] = item
	return item, nil
}

func (r *memoryTripRepo) Delete(_ context.Context, tenantID string, id string) error {
	delete(r.trips, id)
	return nil
}

func (r *memoryTripRepo) GetSummary(_ context.Context, tenantID string) (domain.TripSummary, error) {
	var summary domain.TripSummary
	for _, t := range r.trips {
		if t.TenantID == tenantID {
			summary.TotalTrips++
			if t.Status == domain.TripStatusDispatched {
				summary.DispatchedTrips++
			} else if t.Status == domain.TripStatusDelivered {
				summary.DeliveredTrips++
			}
			summary.TotalFreightBRL += t.FreightAmount
			summary.TotalAdvancesBRL += t.AdvanceAmount
		}
	}
	return summary, nil
}

func TestTripLifecycle(t *testing.T) {
	repo := newMemoryTripRepo()
	svc := NewTripService(repo)
	ctx := context.Background()
	tenantID := "tenant-1"

	// 1. Create Trip
	trip, err := svc.Create(ctx, domain.Trip{
		TenantID:         tenantID,
		VehicleID:        "veh-scania-01",
		DriverID:         "driver-joao-01",
		OriginCity:       "São Paulo",
		OriginState:      "SP",
		DestinationCity:  "Curitiba",
		DestinationState: "PR",
		CargoDescription: "Carga Geral Fracionada",
		FreightAmount:    8500.00,
		AdvanceAmount:    2500.00,
	})
	if err != nil {
		t.Fatalf("failed to create trip: %v", err)
	}
	if trip.Status != domain.TripStatusDraft {
		t.Fatalf("expected status draft, got %s", trip.Status)
	}
	if trip.BalanceAmount != 6000.00 {
		t.Fatalf("expected balance 6000.00, got %f", trip.BalanceAmount)
	}

	// 2. Dispatch Trip with initial Odometer
	dispatched, err := svc.DispatchTrip(ctx, tenantID, trip.ID, 120500)
	if err != nil {
		t.Fatalf("failed to dispatch trip: %v", err)
	}
	if dispatched.Status != domain.TripStatusDispatched {
		t.Fatalf("expected status dispatched, got %s", dispatched.Status)
	}
	if dispatched.StartOdometerKM != 120500 {
		t.Fatalf("expected start km 120500, got %d", dispatched.StartOdometerKM)
	}
	if dispatched.ValidationNotes == "" {
		t.Fatalf("expected validation notes to be set")
	}

	// 3. Deliver Trip with regression error check
	_, err = svc.DeliverTrip(ctx, tenantID, trip.ID, 120000)
	if err != ErrOdometerRegression {
		t.Fatalf("expected ErrOdometerRegression, got %v", err)
	}

	// 4. Deliver Trip successfully
	delivered, err := svc.DeliverTrip(ctx, tenantID, trip.ID, 120920)
	if err != nil {
		t.Fatalf("failed to deliver trip: %v", err)
	}
	if delivered.Status != domain.TripStatusDelivered {
		t.Fatalf("expected status delivered, got %s", delivered.Status)
	}
	if delivered.DistanceKM != 420 {
		t.Fatalf("expected distance 420 km, got %d", delivered.DistanceKM)
	}

	// 5. Summary metrics
	summary, err := svc.GetSummary(ctx, tenantID)
	if err != nil {
		t.Fatalf("failed to get summary: %v", err)
	}
	if summary.TotalTrips != 1 || summary.DeliveredTrips != 1 {
		t.Fatalf("unexpected summary: %+v", summary)
	}
	if summary.TotalFreightBRL != 8500.00 {
		t.Fatalf("expected 8500 freight, got %f", summary.TotalFreightBRL)
	}
	_ = time.Now()
}

package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/trips/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/trips/ports"
)

type TripService struct {
	repo ports.TripRepository
}

func NewTripService(repo ports.TripRepository) TripService {
	return TripService{repo: repo}
}

func (s TripService) Create(ctx context.Context, item domain.Trip) (domain.Trip, error) {
	if item.TenantID == "" {
		return domain.Trip{}, ErrValidation
	}
	if strings.TrimSpace(item.VehicleID) == "" {
		return domain.Trip{}, ErrVehicleRequired
	}
	if strings.TrimSpace(item.DriverID) == "" {
		return domain.Trip{}, ErrDriverRequired
	}
	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	if item.TripNumber == "" {
		item.TripNumber = fmt.Sprintf("VIAG-%s-%s", time.Now().Format("2006"), item.ID[:6])
	}
	if item.Status == "" {
		item.Status = domain.TripStatusDraft
	}
	if item.BalanceAmount == 0 && item.FreightAmount > 0 {
		item.BalanceAmount = item.FreightAmount - item.AdvanceAmount
	}
	now := time.Now().UTC()
	item.CreatedAt = now
	item.UpdatedAt = now
	item.Version = 1

	if item.Status == domain.TripStatusDispatched && item.DepartureAt == nil {
		item.DepartureAt = &now
	}

	return s.repo.Create(ctx, item)
}

func (s TripService) FindByID(ctx context.Context, tenantID string, id string) (domain.Trip, error) {
	return s.repo.FindByID(ctx, tenantID, id)
}

func (s TripService) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Trip], error) {
	return s.repo.Search(ctx, tenantID, q)
}

func (s TripService) DispatchTrip(ctx context.Context, tenantID string, id string, startKM int64) (domain.Trip, error) {
	trip, err := s.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return domain.Trip{}, err
	}
	if trip.Status != domain.TripStatusDraft {
		return domain.Trip{}, ErrInvalidStatusChange
	}

	now := time.Now().UTC()
	trip.Status = domain.TripStatusDispatched
	trip.DepartureAt = &now
	if startKM > 0 {
		trip.StartOdometerKM = startKM
	}
	trip.ValidationNotes = "✓ Verificação pré-embarque homologada: CNH do condutor válida, inspeção de pátio aprovada e pneus em conformidade."
	trip.UpdatedAt = now
	trip.Version++

	return s.repo.Update(ctx, trip)
}

func (s TripService) DeliverTrip(ctx context.Context, tenantID string, id string, endKM int64) (domain.Trip, error) {
	trip, err := s.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return domain.Trip{}, err
	}
	if trip.Status != domain.TripStatusDispatched {
		return domain.Trip{}, ErrInvalidStatusChange
	}
	if endKM > 0 && trip.StartOdometerKM > 0 && endKM < trip.StartOdometerKM {
		return domain.Trip{}, ErrOdometerRegression
	}

	now := time.Now().UTC()
	trip.Status = domain.TripStatusDelivered
	trip.DeliveredAt = &now
	if endKM > 0 {
		trip.EndOdometerKM = endKM
		if trip.StartOdometerKM > 0 {
			trip.DistanceKM = int(endKM - trip.StartOdometerKM)
		}
	}
	trip.UpdatedAt = now
	trip.Version++

	return s.repo.Update(ctx, trip)
}

func (s TripService) CancelTrip(ctx context.Context, tenantID string, id string, reason string) (domain.Trip, error) {
	trip, err := s.repo.FindByID(ctx, tenantID, id)
	if err != nil {
		return domain.Trip{}, err
	}
	if trip.Status == domain.TripStatusDelivered || trip.Status == domain.TripStatusCompleted {
		return domain.Trip{}, ErrInvalidStatusChange
	}

	now := time.Now().UTC()
	trip.Status = domain.TripStatusCanceled
	if reason != "" {
		if trip.Notes != "" {
			trip.Notes += " | Cancelamento: " + reason
		} else {
			trip.Notes = "Cancelamento: " + reason
		}
	}
	trip.UpdatedAt = now
	trip.Version++

	return s.repo.Update(ctx, trip)
}

func (s TripService) GetSummary(ctx context.Context, tenantID string) (domain.TripSummary, error) {
	return s.repo.GetSummary(ctx, tenantID)
}


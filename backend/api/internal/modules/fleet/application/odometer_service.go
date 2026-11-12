package application

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/ports"
	"github.com/geevheit/intelligence360/backend/api/internal/shared/contextkeys"
)

type OdometerService struct {
	repo     ports.OdometerRepository
	vehicles ports.VehicleRepository
}

func NewOdometerService(repo ports.OdometerRepository, vehicles ports.VehicleRepository) OdometerService {
	return OdometerService{
		repo:     repo,
		vehicles: vehicles,
	}
}

var validOdometerSources = map[string]bool{
	"manual":      true,
	"fueling":     true,
	"maintenance": true,
	"checklist":   true,
	"telemetry":   true,
	"correction":  true,
}

func (s OdometerService) RecordReading(ctx context.Context, reading domain.OdometerReading) (domain.OdometerReading, error) {
	if reading.TenantID == "" || reading.VehicleID == "" || reading.Reading < 0 {
		return domain.OdometerReading{}, ErrValidation
	}
	if _, err := uuid.Parse(reading.TenantID); err != nil {
		return domain.OdometerReading{}, ErrValidation
	}
	if _, err := uuid.Parse(reading.VehicleID); err != nil {
		return domain.OdometerReading{}, ErrValidation
	}
	if reading.DriverID != nil && *reading.DriverID != "" {
		if _, err := uuid.Parse(*reading.DriverID); err != nil {
			return domain.OdometerReading{}, ErrValidation
		}
	}

	if reading.ActorID == "" {
		if actorID, ok := contextkeys.ActorIDFromContext(ctx); ok && actorID != "" {
			reading.ActorID = actorID
		}
	}
	if reading.ActorID == "" {
		return domain.OdometerReading{}, ErrValidation
	}
	if _, err := uuid.Parse(reading.ActorID); err != nil {
		return domain.OdometerReading{}, ErrValidation
	}

	reading.Source = strings.ToLower(strings.TrimSpace(reading.Source))
	if reading.Source == "" {
		reading.Source = "manual"
	}
	if !validOdometerSources[reading.Source] {
		return domain.OdometerReading{}, ErrValidation
	}

	// Verify vehicle existence and tenant ownership
	if _, err := s.vehicles.FindByID(ctx, reading.TenantID, reading.VehicleID); err != nil {
		return domain.OdometerReading{}, err
	}

	// Anti-regression check against latest reading
	latest, err := s.repo.GetLatest(ctx, reading.TenantID, reading.VehicleID)
	if err != nil && err != ErrNotFound {
		return domain.OdometerReading{}, err
	}

	if latest != nil {
		if reading.Reading < latest.Reading {
			if !reading.IsCorrection {
				return domain.OdometerReading{}, ErrOdometerRegression
			}
			if len(strings.TrimSpace(reading.CorrectionReason)) < 10 {
				return domain.OdometerReading{}, ErrOdometerCorrectionRequired
			}
		}
	}

	if reading.IsCorrection {
		if len(strings.TrimSpace(reading.CorrectionReason)) < 10 {
			return domain.OdometerReading{}, ErrOdometerCorrectionRequired
		}
	}

	now := time.Now().UTC()
	if reading.ReadAt.IsZero() {
		reading.ReadAt = now
	}
	reading.ID = uuid.NewString()
	reading.CreatedAt = now

	return s.repo.Record(ctx, reading)
}

func (s OdometerService) GetLatest(ctx context.Context, tenantID, vehicleID string) (*domain.OdometerReading, error) {
	if tenantID == "" || vehicleID == "" {
		return nil, ErrValidation
	}
	if _, err := s.vehicles.FindByID(ctx, tenantID, vehicleID); err != nil {
		return nil, err
	}
	return s.repo.GetLatest(ctx, tenantID, vehicleID)
}

func (s OdometerService) GetSummary(ctx context.Context, tenantID, vehicleID string) (domain.VehicleOdometerSummary, error) {
	if tenantID == "" || vehicleID == "" {
		return domain.VehicleOdometerSummary{}, ErrValidation
	}
	if _, err := s.vehicles.FindByID(ctx, tenantID, vehicleID); err != nil {
		return domain.VehicleOdometerSummary{}, err
	}

	latest, err := s.repo.GetLatest(ctx, tenantID, vehicleID)
	if err != nil {
		return domain.VehicleOdometerSummary{}, err
	}

	summary := domain.VehicleOdometerSummary{
		VehicleID: vehicleID,
	}
	if latest != nil {
		summary.CurrentKM = latest.Reading
		readAt := latest.ReadAt
		summary.LastReadAt = &readAt
		readingID := latest.ID
		summary.LastReadingID = &readingID
	}
	return summary, nil
}

func (s OdometerService) ListHistory(ctx context.Context, tenantID, vehicleID string, query ports.Query) (ports.Page[domain.OdometerReading], error) {
	if tenantID == "" || vehicleID == "" {
		return ports.Page[domain.OdometerReading]{}, ErrValidation
	}
	if _, err := s.vehicles.FindByID(ctx, tenantID, vehicleID); err != nil {
		return ports.Page[domain.OdometerReading]{}, err
	}
	return s.repo.ListByVehicle(ctx, tenantID, vehicleID, normalizeQuery(query))
}


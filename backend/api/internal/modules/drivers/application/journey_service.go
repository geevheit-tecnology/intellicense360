package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/drivers/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/drivers/ports"
)

type JourneyService struct {
	repo ports.JourneyRepository
}

func NewJourneyService(repo ports.JourneyRepository) JourneyService {
	return JourneyService{repo: repo}
}

func (s JourneyService) RecordEvent(ctx context.Context, item domain.DriverJourneyRecord) (domain.DriverJourneyRecord, error) {
	if item.TenantID == "" || item.DriverID == "" || item.EventType == "" {
		return domain.DriverJourneyRecord{}, ErrValidation
	}

	now := time.Now().UTC()
	if item.StartAt.IsZero() {
		item.StartAt = now
	}

	// Check if there is an active previous event to close
	active, err := s.repo.GetLastActiveRecord(ctx, item.TenantID, item.DriverID)
	if err == nil && active != nil {
		duration := int(item.StartAt.Sub(active.StartAt).Minutes())
		if duration < 0 {
			duration = 0
		}
		isViolation := false
		violationReason := ""

		// Lei 13.103/2015 Art. 67-C: Maximum 5h30 (330 min) continuous driving
		if active.EventType == domain.EventDriving && duration > 330 {
			isViolation = true
			hours := duration / 60
			mins := duration % 60
			violationReason = fmt.Sprintf("Excesso de direção contínua: %dh%02dm sem pausa regulamentar de 30 min (Lei 13.103/2015)", hours, mins)
		}

		_ = s.repo.CloseActiveRecord(ctx, item.TenantID, active.ID, item.StartAt, duration, isViolation, violationReason)
	}

	if item.ID == "" {
		item.ID = uuid.NewString()
	}
	item.CreatedAt = now
	item.UpdatedAt = now
	item.Version = 1

	// Daily food allowance calculation default
	if item.DailyAllowanceAmount == 0 {
		if item.EventType == domain.EventMeal {
			item.DailyAllowanceAmount = 45.00 // Diária alimentação padrão convenção
		} else if item.EventType == domain.EventDailyRest {
			item.DailyAllowanceAmount = 90.00 // Pernoite / estadia
		}
	}

	return s.repo.Create(ctx, item)
}

func (s JourneyService) ListJourney(ctx context.Context, tenantID, driverID string, limit int) ([]domain.DriverJourneyRecord, error) {
	if tenantID == "" || driverID == "" {
		return nil, ErrValidation
	}
	return s.repo.ListByDriver(ctx, tenantID, driverID, limit)
}

func (s JourneyService) GetSummary(ctx context.Context, tenantID, driverID string) (domain.JourneySummary, error) {
	if tenantID == "" || driverID == "" {
		return domain.JourneySummary{}, ErrValidation
	}
	return s.repo.GetSummary(ctx, tenantID, driverID)
}


package ports

import (
	"context"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/drivers/domain"
)

type DriversService interface {
	ValidateDriverAccess(context.Context, string, string) error
}
type ManagementService interface {
	DriversService
	Create(context.Context, domain.Driver) (domain.Driver, error)
	FindByID(context.Context, string, string) (domain.Driver, error)
	Search(context.Context, string, Query) (Page, error)
	Update(context.Context, domain.Driver) (domain.Driver, error)
	Delete(context.Context, string, string, int64) error
}

type JourneyService interface {
	RecordEvent(ctx context.Context, item domain.DriverJourneyRecord) (domain.DriverJourneyRecord, error)
	ListJourney(ctx context.Context, tenantID, driverID string, limit int) ([]domain.DriverJourneyRecord, error)
	GetSummary(ctx context.Context, tenantID, driverID string) (domain.JourneySummary, error)
}

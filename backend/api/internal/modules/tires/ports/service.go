package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/tires/domain"
)

type TireService interface {
	Create(ctx context.Context, tire domain.Tire) (domain.Tire, error)
	Update(ctx context.Context, tire domain.Tire) (domain.Tire, error)
	Delete(ctx context.Context, tenantID string, tireID string) error
	FindByID(ctx context.Context, tenantID string, tireID string) (domain.Tire, error)
	Search(ctx context.Context, tenantID string, query Query) (Page[domain.Tire], error)
	Receive(ctx context.Context, tenantID string, tireID string, reason string, actorID string) (domain.Tire, error)
	Install(ctx context.Context, tenantID string, tireID string, vehicleID string, position string, km int64, actorID string) (domain.Tire, error)
	Remove(ctx context.Context, tenantID string, tireID string, km int64, reason string, actorID string) (domain.Tire, error)
	Rotate(ctx context.Context, tenantID string, tireID string, position string, km int64, actorID string) (domain.Tire, error)
	SendToRecap(ctx context.Context, tenantID string, tireID string, reason string, actorID string) (domain.Tire, error)
	ReturnFromRecap(ctx context.Context, tenantID string, tireID string, actorID string) (domain.Tire, error)
	Dispose(ctx context.Context, tenantID string, tireID string, reason string, actorID string) (domain.Tire, error)
	ValidateTireAccess(ctx context.Context, tenantID string, tireID string) error
}

type TireInspectionService interface {
	Register(ctx context.Context, inspection domain.TireInspection) (domain.TireInspection, error)
	Update(ctx context.Context, inspection domain.TireInspection) (domain.TireInspection, error)
	Delete(ctx context.Context, tenantID string, tireID string, inspectionID string) error
	List(ctx context.Context, tenantID string, tireID string, query Query) (Page[domain.TireInspection], error)
}

type TireMovementService interface {
	Register(ctx context.Context, movement domain.TireMovement) (domain.TireMovement, error)
	List(ctx context.Context, tenantID string, tireID string, query Query) (Page[domain.TireMovement], error)
}

type TireHistoryService interface {
	List(ctx context.Context, tenantID string, tireID string, query Query) (Page[domain.TireHistory], error)
}

type CostCalculator interface {
	CalculateCostPerKM(ctx context.Context, cost domain.TireCost, km int64) (float64, error)
}

type AttachmentPort interface {
	PrepareAttachment(ctx context.Context, attachment domain.TireAttachment) error
}

package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/trips/domain"
)

type TripService interface {
	Create(ctx context.Context, item domain.Trip) (domain.Trip, error)
	FindByID(ctx context.Context, tenantID string, id string) (domain.Trip, error)
	Search(ctx context.Context, tenantID string, q Query) (Page[domain.Trip], error)
	DispatchTrip(ctx context.Context, tenantID string, id string, startKM int64) (domain.Trip, error)
	DeliverTrip(ctx context.Context, tenantID string, id string, endKM int64) (domain.Trip, error)
	CancelTrip(ctx context.Context, tenantID string, id string, reason string) (domain.Trip, error)
	GetSummary(ctx context.Context, tenantID string) (domain.TripSummary, error)
}


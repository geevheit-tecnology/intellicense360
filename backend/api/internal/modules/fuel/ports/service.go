package ports

import "context"

type FuelService interface {
	ValidateFuelEventAccess(ctx context.Context, tenantID string, fuelEventID string) error
}

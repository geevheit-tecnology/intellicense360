package ports

import "context"

type FuelEventRepository interface {
	Exists(ctx context.Context, tenantID string, fuelEventID string) (bool, error)
}

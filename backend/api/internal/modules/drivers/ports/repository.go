package ports

import "context"

type DriverRepository interface {
	Exists(ctx context.Context, tenantID string, driverID string) (bool, error)
}

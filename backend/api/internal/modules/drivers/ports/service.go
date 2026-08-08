package ports

import "context"

type DriversService interface {
	ValidateDriverAccess(ctx context.Context, tenantID string, driverID string) error
}

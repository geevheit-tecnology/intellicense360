package ports

import "context"

type CIOTService interface {
	ValidateCIOTAccess(ctx context.Context, tenantID string, ciotID string) error
}

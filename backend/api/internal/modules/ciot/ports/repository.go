package ports

import "context"

type CIOTRepository interface {
	Exists(ctx context.Context, tenantID string, ciotID string) (bool, error)
}

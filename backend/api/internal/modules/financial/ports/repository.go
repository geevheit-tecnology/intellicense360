package ports

import "context"

type FinancialEventRepository interface {
	Exists(ctx context.Context, tenantID string, financialEventID string) (bool, error)
}

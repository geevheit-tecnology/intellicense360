package ports

import "context"

type FinancialService interface {
	ValidateFinancialEventAccess(ctx context.Context, tenantID string, financialEventID string) error
}

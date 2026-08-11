package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/financial/domain"
)

type TransactionService interface {
	CreateExpense(context.Context, domain.FinancialTransaction) (domain.FinancialTransaction, error)
	CreateRevenue(context.Context, domain.FinancialTransaction) (domain.FinancialTransaction, error)
	Update(context.Context, domain.FinancialTransaction) (domain.FinancialTransaction, error)
	Approve(context.Context, string, string, string) (domain.FinancialTransaction, error)
	Pay(context.Context, string, string, string) (domain.FinancialTransaction, error)
	Receive(context.Context, string, string, string) (domain.FinancialTransaction, error)
	Cancel(context.Context, string, string, string) (domain.FinancialTransaction, error)
	Adjust(context.Context, domain.FinancialAdjustment) (domain.FinancialAdjustment, error)
	FindByID(context.Context, string, string) (domain.FinancialTransaction, error)
	Search(context.Context, string, Query) (Page[domain.FinancialTransaction], error)
}
type CatalogService[T any] interface {
	Create(context.Context, T) (T, error)
	Search(context.Context, string, Query) (Page[T], error)
}
type PeriodService interface {
	Create(context.Context, domain.FinancialPeriod) (domain.FinancialPeriod, error)
	Close(context.Context, string, string) (domain.FinancialPeriod, error)
	Search(context.Context, string, Query) (Page[domain.FinancialPeriod], error)
}
type BudgetService interface {
	Create(context.Context, domain.Budget) (domain.Budget, error)
	Search(context.Context, string, Query) (Page[domain.Budget], error)
}
type AdjustmentService interface {
	Search(context.Context, string, Query) (Page[domain.FinancialAdjustment], error)
}

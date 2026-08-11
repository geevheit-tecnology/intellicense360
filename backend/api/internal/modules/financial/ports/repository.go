package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/financial/domain"
)

type Query struct {
	Search            string
	Page, PageSize    int
	SortBy, SortOrder string
	Filters           map[string]string
}
type Page[T any] struct {
	Items      []T `json:"items"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type TransactionRepository interface {
	Create(context.Context, domain.FinancialTransaction) (domain.FinancialTransaction, error)
	FindByID(context.Context, string, string) (domain.FinancialTransaction, error)
	Search(context.Context, string, Query) (Page[domain.FinancialTransaction], error)
	Update(context.Context, domain.FinancialTransaction) (domain.FinancialTransaction, error)
	Delete(context.Context, string, string) error
}
type CatalogRepository[T any] interface {
	Create(context.Context, T) (T, error)
	Search(context.Context, string, Query) (Page[T], error)
}
type PeriodRepository interface {
	Create(context.Context, domain.FinancialPeriod) (domain.FinancialPeriod, error)
	FindByID(context.Context, string, string) (domain.FinancialPeriod, error)
	Search(context.Context, string, Query) (Page[domain.FinancialPeriod], error)
	Update(context.Context, domain.FinancialPeriod) (domain.FinancialPeriod, error)
	FindClosedForDate(context.Context, string, string) (bool, error)
}
type BudgetRepository interface {
	Create(context.Context, domain.Budget) (domain.Budget, error)
	Search(context.Context, string, Query) (Page[domain.Budget], error)
}
type AdjustmentRepository interface {
	Create(context.Context, domain.FinancialAdjustment) (domain.FinancialAdjustment, error)
	Search(context.Context, string, Query) (Page[domain.FinancialAdjustment], error)
}

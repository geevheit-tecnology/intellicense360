package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/domain"
)

type Query struct {
	Search    string
	Page      int
	PageSize  int
	SortBy    string
	SortOrder string
	Filters   map[string]string
}

type Page[T any] struct {
	Items      []T `json:"items"`
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

type FuelTransactionRepository interface {
	Create(context.Context, domain.FuelTransaction) (domain.FuelTransaction, error)
	FindByID(context.Context, string, string) (domain.FuelTransaction, error)
	Search(context.Context, string, Query) (Page[domain.FuelTransaction], error)
	Update(context.Context, domain.FuelTransaction) (domain.FuelTransaction, error)
	Delete(context.Context, string, string) error
}

type FuelTypeRepository interface {
	Create(context.Context, domain.FuelType) (domain.FuelType, error)
	FindByID(context.Context, string, string) (domain.FuelType, error)
	Search(context.Context, string, Query) (Page[domain.FuelType], error)
	Update(context.Context, domain.FuelType) (domain.FuelType, error)
	Delete(context.Context, string, string) error
}

type FuelStationRepository interface {
	Create(context.Context, domain.FuelStation) (domain.FuelStation, error)
	FindByID(context.Context, string, string) (domain.FuelStation, error)
	Search(context.Context, string, Query) (Page[domain.FuelStation], error)
	Update(context.Context, domain.FuelStation) (domain.FuelStation, error)
	Delete(context.Context, string, string) error
	ExistsCNPJ(context.Context, string, string, string) (bool, error)
}

type FuelTankRepository interface {
	Create(context.Context, domain.FuelTank) (domain.FuelTank, error)
	Search(context.Context, string, Query) (Page[domain.FuelTank], error)
	Update(context.Context, domain.FuelTank) (domain.FuelTank, error)
	Delete(context.Context, string, string) error
}

type FuelNozzleRepository interface {
	Create(context.Context, domain.FuelNozzle) (domain.FuelNozzle, error)
	Search(context.Context, string, Query) (Page[domain.FuelNozzle], error)
	Update(context.Context, domain.FuelNozzle) (domain.FuelNozzle, error)
	Delete(context.Context, string, string) error
}

type FuelReadingRepository interface {
	Create(context.Context, domain.FuelReading) (domain.FuelReading, error)
	Search(context.Context, string, Query) (Page[domain.FuelReading], error)
}

type FuelPriceRepository interface {
	Create(context.Context, domain.FuelPrice) (domain.FuelPrice, error)
	Search(context.Context, string, Query) (Page[domain.FuelPrice], error)
	Update(context.Context, domain.FuelPrice) (domain.FuelPrice, error)
	Delete(context.Context, string, string) error
}

type FuelReceiptRepository interface {
	Create(context.Context, domain.FuelReceipt) (domain.FuelReceipt, error)
	Search(context.Context, string, Query) (Page[domain.FuelReceipt], error)
	Update(context.Context, domain.FuelReceipt) (domain.FuelReceipt, error)
	Delete(context.Context, string, string) error
}

type FuelAdjustmentRepository interface {
	Create(context.Context, domain.FuelAdjustment) (domain.FuelAdjustment, error)
	Search(context.Context, string, Query) (Page[domain.FuelAdjustment], error)
}

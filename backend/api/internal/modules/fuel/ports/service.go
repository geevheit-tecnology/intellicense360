package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/domain"
)

type FuelTransactionService interface {
	Create(context.Context, domain.FuelTransaction) (domain.FuelTransaction, error)
	Update(context.Context, domain.FuelTransaction) (domain.FuelTransaction, error)
	Complete(context.Context, string, string, string) (domain.FuelTransaction, error)
	Cancel(context.Context, string, string, string, string) (domain.FuelTransaction, error)
	Adjust(context.Context, domain.FuelAdjustment) (domain.FuelAdjustment, error)
	Delete(context.Context, string, string) error
	FindByID(context.Context, string, string) (domain.FuelTransaction, error)
	Search(context.Context, string, Query) (Page[domain.FuelTransaction], error)
}

type FuelTypeService interface {
	Create(context.Context, domain.FuelType) (domain.FuelType, error)
	Update(context.Context, domain.FuelType) (domain.FuelType, error)
	Delete(context.Context, string, string) error
	FindByID(context.Context, string, string) (domain.FuelType, error)
	Search(context.Context, string, Query) (Page[domain.FuelType], error)
}

type FuelStationService interface {
	Create(context.Context, domain.FuelStation) (domain.FuelStation, error)
	Update(context.Context, domain.FuelStation) (domain.FuelStation, error)
	Delete(context.Context, string, string) error
	FindByID(context.Context, string, string) (domain.FuelStation, error)
	Search(context.Context, string, Query) (Page[domain.FuelStation], error)
}

type FuelTankService interface {
	Create(context.Context, domain.FuelTank) (domain.FuelTank, error)
	Update(context.Context, domain.FuelTank) (domain.FuelTank, error)
	Delete(context.Context, string, string) error
	Search(context.Context, string, Query) (Page[domain.FuelTank], error)
}

type FuelNozzleService interface {
	Create(context.Context, domain.FuelNozzle) (domain.FuelNozzle, error)
	Update(context.Context, domain.FuelNozzle) (domain.FuelNozzle, error)
	Delete(context.Context, string, string) error
	Search(context.Context, string, Query) (Page[domain.FuelNozzle], error)
}

type FuelReadingService interface {
	Record(context.Context, domain.FuelReading) (domain.FuelReading, error)
	Search(context.Context, string, Query) (Page[domain.FuelReading], error)
}

type FuelPriceService interface {
	Record(context.Context, domain.FuelPrice) (domain.FuelPrice, error)
	Update(context.Context, domain.FuelPrice) (domain.FuelPrice, error)
	Delete(context.Context, string, string) error
	Search(context.Context, string, Query) (Page[domain.FuelPrice], error)
}

type FuelReceiptService interface {
	Create(context.Context, domain.FuelReceipt) (domain.FuelReceipt, error)
	Update(context.Context, domain.FuelReceipt) (domain.FuelReceipt, error)
	Delete(context.Context, string, string) error
	Search(context.Context, string, Query) (Page[domain.FuelReceipt], error)
}

type FuelAdjustmentService interface {
	Search(context.Context, string, Query) (Page[domain.FuelAdjustment], error)
}

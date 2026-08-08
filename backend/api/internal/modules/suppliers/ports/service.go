package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/domain"
)

type SupplierService interface {
	Create(context.Context, domain.Supplier) (domain.Supplier, error)
	FindByID(context.Context, string, string) (domain.Supplier, error)
	Search(context.Context, string, Query) (Page[domain.Supplier], error)
	Update(context.Context, domain.Supplier) (domain.Supplier, error)
	Delete(context.Context, string, string) error
}

type CategoryService interface {
	Create(context.Context, domain.SupplierCategory) (domain.SupplierCategory, error)
	Search(context.Context, string, Query) (Page[domain.SupplierCategory], error)
	Update(context.Context, domain.SupplierCategory) (domain.SupplierCategory, error)
	Delete(context.Context, string, string) error
}

type TypeService interface {
	Create(context.Context, domain.SupplierType) (domain.SupplierType, error)
	Search(context.Context, string, Query) (Page[domain.SupplierType], error)
	Update(context.Context, domain.SupplierType) (domain.SupplierType, error)
	Delete(context.Context, string, string) error
}

type ContactService interface {
	Create(context.Context, domain.SupplierContact) (domain.SupplierContact, error)
	Search(context.Context, string, Query) (Page[domain.SupplierContact], error)
	Update(context.Context, domain.SupplierContact) (domain.SupplierContact, error)
	Delete(context.Context, string, string) error
}

type AddressService interface {
	Create(context.Context, domain.SupplierAddress) (domain.SupplierAddress, error)
	Search(context.Context, string, Query) (Page[domain.SupplierAddress], error)
	Update(context.Context, domain.SupplierAddress) (domain.SupplierAddress, error)
	Delete(context.Context, string, string) error
}

type DocumentService interface {
	Create(context.Context, domain.SupplierDocument) (domain.SupplierDocument, error)
	Search(context.Context, string, Query) (Page[domain.SupplierDocument], error)
	Update(context.Context, domain.SupplierDocument) (domain.SupplierDocument, error)
	Delete(context.Context, string, string) error
}

type ContractService interface {
	Create(context.Context, domain.SupplierContract) (domain.SupplierContract, error)
	Search(context.Context, string, Query) (Page[domain.SupplierContract], error)
	Update(context.Context, domain.SupplierContract) (domain.SupplierContract, error)
	Delete(context.Context, string, string) error
}

type RatingService interface {
	Create(context.Context, domain.SupplierRating) (domain.SupplierRating, error)
	Search(context.Context, string, Query) (Page[domain.SupplierRating], error)
	Update(context.Context, domain.SupplierRating) (domain.SupplierRating, error)
	Delete(context.Context, string, string) error
}

type AuditRecorder interface {
	RecordSupplierEvent(context.Context, string, string, string, string) error
}

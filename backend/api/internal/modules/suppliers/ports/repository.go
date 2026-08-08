package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/domain"
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

type SupplierRepository interface {
	Create(context.Context, domain.Supplier) (domain.Supplier, error)
	FindByID(context.Context, string, string) (domain.Supplier, error)
	Search(context.Context, string, Query) (Page[domain.Supplier], error)
	Update(context.Context, domain.Supplier) (domain.Supplier, error)
	Delete(context.Context, string, string) error
	ExistsCNPJ(context.Context, string, string, string) (bool, error)
	Exists(context.Context, string, string) (bool, error)
}

type CategoryRepository interface {
	Create(context.Context, domain.SupplierCategory) (domain.SupplierCategory, error)
	Search(context.Context, string, Query) (Page[domain.SupplierCategory], error)
	Update(context.Context, domain.SupplierCategory) (domain.SupplierCategory, error)
	Delete(context.Context, string, string) error
	Exists(context.Context, string, string) (bool, error)
}

type TypeRepository interface {
	Create(context.Context, domain.SupplierType) (domain.SupplierType, error)
	Search(context.Context, string, Query) (Page[domain.SupplierType], error)
	Update(context.Context, domain.SupplierType) (domain.SupplierType, error)
	Delete(context.Context, string, string) error
}

type ContactRepository interface {
	Create(context.Context, domain.SupplierContact) (domain.SupplierContact, error)
	Search(context.Context, string, Query) (Page[domain.SupplierContact], error)
	Update(context.Context, domain.SupplierContact) (domain.SupplierContact, error)
	Delete(context.Context, string, string) error
}

type AddressRepository interface {
	Create(context.Context, domain.SupplierAddress) (domain.SupplierAddress, error)
	Search(context.Context, string, Query) (Page[domain.SupplierAddress], error)
	Update(context.Context, domain.SupplierAddress) (domain.SupplierAddress, error)
	Delete(context.Context, string, string) error
}

type DocumentRepository interface {
	Create(context.Context, domain.SupplierDocument) (domain.SupplierDocument, error)
	Search(context.Context, string, Query) (Page[domain.SupplierDocument], error)
	Update(context.Context, domain.SupplierDocument) (domain.SupplierDocument, error)
	Delete(context.Context, string, string) error
}

type ContractRepository interface {
	Create(context.Context, domain.SupplierContract) (domain.SupplierContract, error)
	Search(context.Context, string, Query) (Page[domain.SupplierContract], error)
	Update(context.Context, domain.SupplierContract) (domain.SupplierContract, error)
	Delete(context.Context, string, string) error
}

type RatingRepository interface {
	Create(context.Context, domain.SupplierRating) (domain.SupplierRating, error)
	Search(context.Context, string, Query) (Page[domain.SupplierRating], error)
	Update(context.Context, domain.SupplierRating) (domain.SupplierRating, error)
	Delete(context.Context, string, string) error
}

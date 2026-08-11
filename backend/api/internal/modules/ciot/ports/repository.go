package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/ciot/domain"
)

type Query struct {
	Search  string
	Page    int
	PerPage int
	Sort    string
	Filters map[string]string
}

type Page[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type CIOTRepository interface {
	Create(ctx context.Context, item domain.CIOT) (domain.CIOT, error)
	FindByID(ctx context.Context, tenantID string, id string) (domain.CIOT, error)
	Search(ctx context.Context, tenantID string, q Query) (Page[domain.CIOT], error)
	Update(ctx context.Context, item domain.CIOT) (domain.CIOT, error)
	Delete(ctx context.Context, tenantID string, id string) error
	ExistsIdempotencyKey(ctx context.Context, tenantID string, key string) (bool, error)
}

type CatalogRepository[T any] interface {
	Create(ctx context.Context, item T) (T, error)
	Search(ctx context.Context, tenantID string, q Query) (Page[T], error)
}

type StatusHistoryRepository interface {
	Create(ctx context.Context, item domain.CIOTStatusHistory) (domain.CIOTStatusHistory, error)
	ListByCIOT(ctx context.Context, tenantID string, ciotID string, q Query) (Page[domain.CIOTStatusHistory], error)
}

type PaymentRepository interface {
	Create(ctx context.Context, item domain.CIOTPayment) (domain.CIOTPayment, error)
	ListByCIOT(ctx context.Context, tenantID string, ciotID string, q Query) (Page[domain.CIOTPayment], error)
}

type ProviderAttemptRepository interface {
	Create(ctx context.Context, item domain.CIOTProviderAttempt) (domain.CIOTProviderAttempt, error)
	ListByCIOT(ctx context.Context, tenantID string, ciotID string, q Query) (Page[domain.CIOTProviderAttempt], error)
	ExistsIdempotencyKey(ctx context.Context, tenantID string, key string) (bool, error)
}

type ExternalReferenceRepository interface {
	Upsert(ctx context.Context, item domain.CIOTExternalReference) (domain.CIOTExternalReference, error)
	FindByCIOT(ctx context.Context, tenantID string, ciotID string) (domain.CIOTExternalReference, error)
}

type DocumentRepository interface {
	Create(ctx context.Context, item domain.CIOTDocument) (domain.CIOTDocument, error)
	ListByCIOT(ctx context.Context, tenantID string, ciotID string, q Query) (Page[domain.CIOTDocument], error)
}

type ErrorRepository interface {
	Create(ctx context.Context, item domain.CIOTError) (domain.CIOTError, error)
	ListByCIOT(ctx context.Context, tenantID string, ciotID string, q Query) (Page[domain.CIOTError], error)
}

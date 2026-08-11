package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/ciot/domain"
)

type CIOTService interface {
	Create(ctx context.Context, item domain.CIOT) (domain.CIOT, error)
	Update(ctx context.Context, item domain.CIOT) (domain.CIOT, error)
	FindByID(ctx context.Context, tenantID string, id string) (domain.CIOT, error)
	Search(ctx context.Context, tenantID string, q Query) (Page[domain.CIOT], error)
	Submit(ctx context.Context, tenantID string, id string, actorID string) (domain.CIOT, error)
	MarkGenerated(ctx context.Context, tenantID string, id string, externalProtocol string, actorID string) (domain.CIOT, error)
	Activate(ctx context.Context, tenantID string, id string, actorID string) (domain.CIOT, error)
	Suspend(ctx context.Context, tenantID string, id string, actorID string) (domain.CIOT, error)
	Reactivate(ctx context.Context, tenantID string, id string, actorID string) (domain.CIOT, error)
	Close(ctx context.Context, tenantID string, id string, actorID string) (domain.CIOT, error)
	Cancel(ctx context.Context, tenantID string, id string, actorID string) (domain.CIOT, error)
	RecordError(ctx context.Context, tenantID string, id string, code string, message string, actorID string) (domain.CIOT, error)
	RetryFromError(ctx context.Context, tenantID string, id string, actorID string) (domain.CIOT, error)
	Delete(ctx context.Context, tenantID string, id string) error
}

type CatalogService[T any] interface {
	Create(ctx context.Context, item T) (T, error)
	Search(ctx context.Context, tenantID string, q Query) (Page[T], error)
}

type StatusHistoryService interface {
	ListByCIOT(ctx context.Context, tenantID string, ciotID string, q Query) (Page[domain.CIOTStatusHistory], error)
}

type PaymentService interface {
	Create(ctx context.Context, item domain.CIOTPayment) (domain.CIOTPayment, error)
	ListByCIOT(ctx context.Context, tenantID string, ciotID string, q Query) (Page[domain.CIOTPayment], error)
}

type ProviderAttemptService interface {
	Create(ctx context.Context, item domain.CIOTProviderAttempt) (domain.CIOTProviderAttempt, error)
	ListByCIOT(ctx context.Context, tenantID string, ciotID string, q Query) (Page[domain.CIOTProviderAttempt], error)
}

type ExternalReferenceService interface {
	Upsert(ctx context.Context, item domain.CIOTExternalReference) (domain.CIOTExternalReference, error)
	FindByCIOT(ctx context.Context, tenantID string, ciotID string) (domain.CIOTExternalReference, error)
}

type DocumentService interface {
	Create(ctx context.Context, item domain.CIOTDocument) (domain.CIOTDocument, error)
	ListByCIOT(ctx context.Context, tenantID string, ciotID string, q Query) (Page[domain.CIOTDocument], error)
}

type ErrorService interface {
	Create(ctx context.Context, item domain.CIOTError) (domain.CIOTError, error)
	ListByCIOT(ctx context.Context, tenantID string, ciotID string, q Query) (Page[domain.CIOTError], error)
}

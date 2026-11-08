package ports

import (
	"context"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/documents/domain"
)

type DocumentsService interface {
	ValidateDocumentAccess(context.Context, string, string) error
}
type ManagementService interface {
	DocumentsService
	Create(context.Context, domain.Document) (domain.Document, error)
	FindByID(context.Context, string, string) (domain.Document, error)
	Search(context.Context, string, Query) (Page, error)
	Update(context.Context, domain.Document) (domain.Document, error)
	Delete(context.Context, string, string, int64) error
}

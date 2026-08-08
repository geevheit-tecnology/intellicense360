package ports

import "context"

type DocumentsService interface {
	ValidateDocumentAccess(ctx context.Context, tenantID string, documentID string) error
}

package ports

import "context"

type DocumentRepository interface {
	Exists(ctx context.Context, tenantID string, documentID string) (bool, error)
}

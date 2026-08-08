package ports

import "context"

type NotificationRepository interface {
	Exists(ctx context.Context, tenantID string, notificationID string) (bool, error)
}

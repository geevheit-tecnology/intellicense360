package ports

import "context"

type NotificationsService interface {
	ValidateNotificationAccess(ctx context.Context, tenantID string, notificationID string) error
}

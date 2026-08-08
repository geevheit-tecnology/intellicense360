package contextkeys

import "context"

type Key string

const (
	TenantID    Key = "tenant_id"
	ActorID     Key = "actor_id"
	SessionID   Key = "session_id"
	Permissions Key = "permissions"
)

func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, TenantID, tenantID)
}

func TenantIDFromContext(ctx context.Context) (string, bool) {
	tenantID, ok := ctx.Value(TenantID).(string)
	return tenantID, ok && tenantID != ""
}

func WithActorID(ctx context.Context, actorID string) context.Context {
	return context.WithValue(ctx, ActorID, actorID)
}

func ActorIDFromContext(ctx context.Context) (string, bool) {
	actorID, ok := ctx.Value(ActorID).(string)
	return actorID, ok && actorID != ""
}

package ports

import (
	"context"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/domain"
)

type CommandService interface {
	Create(ctx context.Context, item domain.CommandItem, actorID string) (domain.CommandItem, error)
	Update(ctx context.Context, item domain.CommandItem, actorID string) (domain.CommandItem, error)
	Get(ctx context.Context, tenantID, id string) (domain.CommandItem, error)
	List(ctx context.Context, tenantID string, q Query) (Page[domain.CommandItem], error)
	Acknowledge(ctx context.Context, tenantID, id, actorID string) (domain.CommandItem, error)
	Start(ctx context.Context, tenantID, id, actorID string) (domain.CommandItem, error)
	Resolve(ctx context.Context, tenantID, id, actorID string) (domain.CommandItem, error)
	Dismiss(ctx context.Context, tenantID, id, actorID string) (domain.CommandItem, error)
	CreateAction(ctx context.Context, action domain.CommandAction) (domain.CommandAction, error)
	ListActions(ctx context.Context, tenantID, itemID string, q Query) (Page[domain.CommandAction], error)
	History(ctx context.Context, tenantID, itemID string, q Query) (Page[domain.CommandEvent], error)
	Summary(ctx context.Context, tenantID string) (domain.MissionControlSummary, error)
	RebuildSnapshot(ctx context.Context, tenantID string) (domain.OperationalSnapshot, error)
	LatestSnapshot(ctx context.Context, tenantID string) (domain.OperationalSnapshot, error)
	EvaluateRecommendations(ctx context.Context, tenantID string) ([]domain.CommandAction, error)
}

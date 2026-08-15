package outbox

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	coreevents "github.com/geevheit/intelligence360/backend/api/internal/events"
)

type Service struct {
	repo   Repository
	bus    coreevents.EventBus
	config coreevents.RetryConfig
}

func NewService(repo Repository, bus coreevents.EventBus, config coreevents.RetryConfig) Service {
	if config.MaxAttempts <= 0 {
		config = coreevents.DefaultRetryConfig()
	}
	return Service{repo: repo, bus: bus, config: config}
}

func (s Service) Save(ctx context.Context, event coreevents.DomainEvent) (Record, error) {
	if err := event.Validate(); err != nil {
		return Record{}, err
	}
	if _, err := s.repo.FindByEventID(ctx, event.TenantID, event.EventID); err == nil {
		return Record{}, coreevents.ErrDuplicateEvent
	}
	now := time.Now().UTC()
	return s.repo.Save(ctx, Record{ID: newID("out"), TenantID: event.TenantID, EventID: event.EventID, EventType: event.EventType, AggregateID: event.AggregateID, AggregateType: event.AggregateType, Payload: event, OccurredAt: event.OccurredAt, Status: StatusPending, AvailableAt: now, CreatedAt: now, UpdatedAt: now})
}

func (s Service) ProcessPending(ctx context.Context, now time.Time, limit int) ([]Record, error) {
	records, err := s.repo.GetPending(ctx, now, limit)
	if err != nil {
		return nil, err
	}
	processed := []Record{}
	for _, record := range records {
		current, err := s.repo.MarkProcessing(ctx, record.ID)
		if err != nil {
			return processed, err
		}
		if err := s.bus.Publish(ctx, current.Payload); err != nil {
			if current.Attempts+1 >= s.config.MaxAttempts {
				_, err = s.repo.MoveToDeadLetter(ctx, current.ID, err.Error())
				if err != nil {
					return processed, err
				}
				continue
			}
			nextAttempt := now.Add(s.config.Delay)
			current, err = s.repo.MarkFailed(ctx, current.ID, err, nextAttempt)
			if err != nil {
				return processed, err
			}
			processed = append(processed, current)
			continue
		}
		current, err = s.repo.MarkPublished(ctx, current.ID)
		if err != nil {
			return processed, err
		}
		processed = append(processed, current)
	}
	return processed, nil
}

func newID(prefix string) string {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		return prefix + "_" + time.Now().UTC().Format("20060102150405")
	}
	return prefix + "_" + hex.EncodeToString(b[:])
}

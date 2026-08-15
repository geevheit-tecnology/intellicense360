package outbox

import (
	"context"
	"sync"
	"time"

	coreevents "github.com/geevheit/intelligence360/backend/api/internal/events"
)

type MemoryRepository struct {
	mu          sync.RWMutex
	records     map[string]Record
	deadLetters map[string]DeadLetter
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{records: map[string]Record{}, deadLetters: map[string]DeadLetter{}}
}

func (r *MemoryRepository) Save(_ context.Context, record Record) (Record, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records[record.ID] = record
	return record, nil
}

func (r *MemoryRepository) GetPending(_ context.Context, now time.Time, limit int) ([]Record, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := []Record{}
	for _, record := range r.records {
		if len(out) == limit && limit > 0 {
			break
		}
		if (record.Status == StatusPending || record.Status == StatusFailed) && !record.AvailableAt.After(now) {
			out = append(out, record)
		}
	}
	return out, nil
}

func (r *MemoryRepository) MarkProcessing(_ context.Context, id string) (Record, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	record, ok := r.records[id]
	if !ok {
		return Record{}, coreevents.ErrOutboxNotFound
	}
	record.Status = StatusProcessing
	record.UpdatedAt = time.Now().UTC()
	r.records[id] = record
	return record, nil
}

func (r *MemoryRepository) MarkPublished(_ context.Context, id string) (Record, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	record, ok := r.records[id]
	if !ok {
		return Record{}, coreevents.ErrOutboxNotFound
	}
	now := time.Now().UTC()
	record.Status = StatusPublished
	record.PublishedAt = &now
	record.UpdatedAt = now
	r.records[id] = record
	return record, nil
}

func (r *MemoryRepository) MarkFailed(_ context.Context, id string, failure error, nextAttempt time.Time) (Record, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	record, ok := r.records[id]
	if !ok {
		return Record{}, coreevents.ErrOutboxNotFound
	}
	record.Status = StatusFailed
	record.Attempts++
	record.LastError = failure.Error()
	record.AvailableAt = nextAttempt
	record.UpdatedAt = time.Now().UTC()
	r.records[id] = record
	return record, nil
}

func (r *MemoryRepository) MoveToDeadLetter(_ context.Context, id string, reason string) (DeadLetter, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	record, ok := r.records[id]
	if !ok {
		return DeadLetter{}, coreevents.ErrOutboxNotFound
	}
	record.Status = StatusDeadLetter
	record.Attempts++
	record.LastError = reason
	record.UpdatedAt = time.Now().UTC()
	r.records[id] = record
	dead := DeadLetter{ID: newID("dlq"), TenantID: record.TenantID, EventID: record.EventID, EventType: record.EventType, Payload: record.Payload, FailureReason: reason, Attempts: record.Attempts, LastError: reason, CreatedAt: time.Now().UTC()}
	r.deadLetters[dead.ID] = dead
	return dead, nil
}

func (r *MemoryRepository) FindByEventID(_ context.Context, tenantID string, eventID string) (Record, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, record := range r.records {
		if record.TenantID == tenantID && record.EventID == eventID {
			return record, nil
		}
	}
	return Record{}, coreevents.ErrOutboxNotFound
}

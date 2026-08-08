package types

import (
	"time"

	"github.com/google/uuid"
)

type BaseEntity struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Version   int64
	AuditID   uuid.UUID
}

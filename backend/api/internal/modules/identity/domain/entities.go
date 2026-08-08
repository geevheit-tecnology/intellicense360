package domain

import "time"

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
)

type User struct {
	ID           string     `json:"id"`
	TenantID     string     `json:"tenant_id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Status       UserStatus `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type Role struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Permission struct {
	ID          string    `json:"id"`
	Key         string    `json:"key"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Group struct {
	ID          string    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Tenant struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Session struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	UserID    string     `json:"user_id"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type RefreshToken struct {
	ID           string     `json:"id"`
	TenantID     string     `json:"tenant_id"`
	UserID       string     `json:"user_id"`
	SessionID    string     `json:"session_id"`
	TokenHash    string     `json:"-"`
	ExpiresAt    time.Time  `json:"expires_at"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty"`
	ReplacedByID string     `json:"replaced_by_id,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

type PasswordPolicy struct {
	MinLength      int  `json:"min_length"`
	RequireNumber  bool `json:"require_number"`
	RequireSpecial bool `json:"require_special"`
}

type AuditLog struct {
	ID           string            `json:"id"`
	TenantID     string            `json:"tenant_id"`
	ActorID      string            `json:"actor_id,omitempty"`
	Action       string            `json:"action"`
	ResourceType string            `json:"resource_type"`
	ResourceID   string            `json:"resource_id,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
}

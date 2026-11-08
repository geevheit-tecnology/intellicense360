package domain

import "time"

// Driver is the operational profile; license expiry is held by Documents.
// VehicleID is a current allocation, not a trip or an exclusive reservation.
type Driver struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	Code      string     `json:"code"`
	Name      string     `json:"name"`
	TaxID     string     `json:"tax_id,omitempty"`
	Phone     string     `json:"phone,omitempty"`
	Email     string     `json:"email,omitempty"`
	VehicleID string     `json:"vehicle_id,omitempty"`
	Status    string     `json:"status"`
	Notes     string     `json:"notes,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Version   int64      `json:"version"`
}

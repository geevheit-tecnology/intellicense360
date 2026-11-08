package domain

import "time"

// Document records operator-entered metadata. No binary upload or external
// verification is implied by its presence or by its calculated expiry status.
type Document struct {
	ID           string     `json:"id"`
	TenantID     string     `json:"tenant_id"`
	DriverID     string     `json:"driver_id,omitempty"`
	VehicleID    string     `json:"vehicle_id,omitempty"`
	Kind         string     `json:"kind"`
	Title        string     `json:"title"`
	Number       string     `json:"number"`
	Issuer       string     `json:"issuer,omitempty"`
	IssuedOn     *string    `json:"issued_on,omitempty"`
	ExpiresOn    *string    `json:"expires_on,omitempty"`
	Status       string     `json:"status"`
	ExpiryStatus string     `json:"expiry_status"`
	Notes        string     `json:"notes,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
	Version      int64      `json:"version"`
}

// ExpiryAt uses calendar dates (UTC) and an inclusive 30-day attention window.
// A document expiring today is due_soon; it becomes expired the following day.
func (d Document) ExpiryAt(now time.Time) string {
	if d.Status == "archived" {
		return "archived"
	}
	if d.ExpiresOn == nil {
		return "no_expiry"
	}
	expiry, err := time.Parse(time.DateOnly, *d.ExpiresOn)
	if err != nil {
		return "unknown"
	}
	today, _ := time.Parse(time.DateOnly, now.UTC().Format(time.DateOnly))
	if expiry.Before(today) {
		return "expired"
	}
	if !expiry.After(today.AddDate(0, 0, 30)) {
		return "due_soon"
	}
	return "valid"
}

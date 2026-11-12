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

type JourneyEventType string

const (
	EventDriving   JourneyEventType = "driving"
	EventBreakRest JourneyEventType = "break_rest"
	EventMeal      JourneyEventType = "meal"
	EventDailyRest JourneyEventType = "daily_rest"
	EventWaiting   JourneyEventType = "waiting"
	EventStandby   JourneyEventType = "standby"
)

type DriverJourneyRecord struct {
	ID                   string           `json:"id"`
	TenantID             string           `json:"tenant_id"`
	DriverID             string           `json:"driver_id"`
	VehicleID            *string          `json:"vehicle_id,omitempty"`
	TripID               *string          `json:"trip_id,omitempty"`
	EventType            JourneyEventType `json:"event_type"`
	StartAt              time.Time        `json:"start_at"`
	EndAt                *time.Time       `json:"end_at,omitempty"`
	DurationMinutes      int              `json:"duration_minutes"`
	Location             string           `json:"location"`
	OdometerKM           int64            `json:"odometer_km"`
	IsViolation          bool             `json:"is_violation"`
	ViolationReason      string           `json:"violation_reason,omitempty"`
	DailyAllowanceAmount float64          `json:"daily_allowance_amount"`
	Notes                string           `json:"notes,omitempty"`
	CreatedAt            time.Time        `json:"created_at"`
	UpdatedAt            time.Time        `json:"updated_at"`
	Version              int64            `json:"version"`
}

type JourneySummary struct {
	TotalDrivingMinutes int     `json:"total_driving_minutes"`
	TotalRestMinutes    int     `json:"total_rest_minutes"`
	TotalWaitingMinutes int     `json:"total_waiting_minutes"`
	ViolationsCount     int     `json:"violations_count"`
	TotalAllowanceBRL   float64 `json:"total_allowance_amount"`
	ActiveEvent         *string `json:"active_event,omitempty"`
}

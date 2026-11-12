package domain

import (
	"time"
)

type TripStatus string

const (
	TripStatusDraft      TripStatus = "draft"
	TripStatusDispatched TripStatus = "dispatched"
	TripStatusDelivered  TripStatus = "delivered"
	TripStatusCompleted  TripStatus = "completed"
	TripStatusCanceled   TripStatus = "canceled"
)

type Trip struct {
	ID                string     `json:"id"`
	TenantID          string     `json:"tenant_id"`
	TripNumber        string     `json:"trip_number"`
	Status            TripStatus `json:"status"`
	VehicleID         string     `json:"vehicle_id"`
	TrailerID         *string    `json:"trailer_id,omitempty"`
	DriverID          string     `json:"driver_id"`
	OriginCity        string     `json:"origin_city"`
	OriginState       string     `json:"origin_state"`
	DestinationCity   string     `json:"destination_city"`
	DestinationState  string     `json:"destination_state"`
	DistanceKM        int        `json:"distance_km"`
	DepartureAt       *time.Time `json:"departure_at,omitempty"`
	ExpectedArrivalAt *time.Time `json:"expected_arrival_at,omitempty"`
	DeliveredAt       *time.Time `json:"delivered_at,omitempty"`
	CargoDescription  string     `json:"cargo_description"`
	CargoWeightKG     float64    `json:"cargo_weight_kg"`
	CargoValue        float64    `json:"cargo_value"`
	CustomerName      string     `json:"customer_name"`
	MDFEKey           string     `json:"mdfe_key,omitempty"`
	CTEKey            string     `json:"cte_key,omitempty"`
	CIOTNumber        string     `json:"ciot_number,omitempty"`
	FreightAmount     float64    `json:"freight_amount"`
	AdvanceAmount     float64    `json:"advance_amount"`
	TollAmount        float64    `json:"toll_amount"`
	BalanceAmount     float64    `json:"balance_amount"`
	StartOdometerKM   int64      `json:"start_odometer_km"`
	EndOdometerKM     int64      `json:"end_odometer_km"`
	ValidationNotes   string     `json:"validation_notes,omitempty"`
	Notes             string     `json:"notes,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
	Version           int64      `json:"version"`
}

type TripSummary struct {
	TotalTrips       int     `json:"total_trips"`
	DispatchedTrips  int     `json:"dispatched_trips"`
	DeliveredTrips   int     `json:"delivered_trips"`
	TotalFreightBRL  float64 `json:"total_freight_amount"`
	TotalAdvancesBRL float64 `json:"total_advances_amount"`
	TotalDistanceKM  int     `json:"total_distance_km"`
}

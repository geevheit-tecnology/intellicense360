package dto

import "time"

type TripRequest struct {
	TripNumber        string     `json:"trip_number"`
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
	StartOdometerKM   int64      `json:"start_odometer_km"`
	Notes             string     `json:"notes,omitempty"`
}

type DispatchRequest struct {
	StartOdometerKM int64 `json:"start_odometer_km"`
}

type DeliverRequest struct {
	EndOdometerKM int64 `json:"end_odometer_km"`
}

type CancelRequest struct {
	Reason string `json:"reason"`
}


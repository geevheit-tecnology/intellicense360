package dto

type TireRequest struct {
	SerialNumber    string  `json:"serial_number"`
	FireNumber      string  `json:"fire_number"`
	Brand           string  `json:"brand"`
	Model           string  `json:"model"`
	Size            string  `json:"size"`
	Construction    string  `json:"construction"`
	TireType        string  `json:"tire_type"`
	PositionType    string  `json:"position_type"`
	PurchaseValue   float64 `json:"purchase_value"`
	Supplier        string  `json:"supplier"`
	DOT             string  `json:"dot"`
	CurrentTreadMM  float64 `json:"current_tread_mm"`
	OriginalTreadMM float64 `json:"original_tread_mm"`
	MinimumTreadMM  float64 `json:"minimum_tread_mm"`
	Status          string  `json:"status"`
	VehicleID       string  `json:"vehicle_id"`
	Position        string  `json:"position"`
	CurrentKM       int64   `json:"current_km"`
	TotalKM         int64   `json:"total_km"`
	RecapCount      int     `json:"recap_count"`
	Notes           string  `json:"notes"`
}

type TireInspectionRequest struct {
	TireID       string  `json:"tire_id"`
	TreadMM      float64 `json:"tread_mm"`
	Pressure     float64 `json:"pressure"`
	Temperature  float64 `json:"temperature"`
	Condition    string  `json:"condition"`
	Observations string  `json:"observations"`
	Inspector    string  `json:"inspector"`
}

type TireMovementRequest struct {
	MovementType string `json:"movement_type"`
	VehicleID    string `json:"vehicle_id"`
	Position     string `json:"position"`
	KM           int64  `json:"km"`
	Reason       string `json:"reason"`
}

type TireOperationRequest struct {
	VehicleID string `json:"vehicle_id"`
	Position  string `json:"position"`
	KM        int64  `json:"km"`
	Reason    string `json:"reason"`
}

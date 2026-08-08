package dto

type WorkOrderRequest struct {
	Code           string  `json:"code"`
	AssetID        string  `json:"asset_id"`
	VehicleID      string  `json:"vehicle_id"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	Kind           string  `json:"kind"`
	Status         string  `json:"status"`
	Priority       string  `json:"priority"`
	ServiceTypeID  string  `json:"service_type_id"`
	EstimatedHours float64 `json:"estimated_hours"`
}

type PreventivePlanRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	AssetID       string `json:"asset_id"`
	VehicleID     string `json:"vehicle_id"`
	ServiceTypeID string `json:"service_type_id"`
	Frequency     string `json:"frequency"`
	IntervalValue int64  `json:"interval_value"`
}

type ServiceTypeRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type CatalogRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type WorkshopRequest struct {
	Name     string `json:"name"`
	Document string `json:"document"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type TechnicianRequest struct {
	Name     string `json:"name"`
	Document string `json:"document"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Active   bool   `json:"active"`
}
type LaborRequest struct {
	Technician string  `json:"technician"`
	Hours      float64 `json:"hours"`
	Cost       float64 `json:"cost"`
}
type DowntimeRequest struct {
	AssetID   string `json:"asset_id"`
	VehicleID string `json:"vehicle_id"`
	Reason    string `json:"reason"`
}

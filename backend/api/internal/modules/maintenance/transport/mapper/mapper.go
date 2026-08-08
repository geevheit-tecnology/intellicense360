package mapper

import (
	"github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/transport/dto"
)

func WorkOrderFromRequest(req dto.WorkOrderRequest) domain.WorkOrder {
	return domain.WorkOrder{Code: req.Code, AssetID: req.AssetID, VehicleID: req.VehicleID, Title: req.Title, Description: req.Description, Kind: domain.MaintenanceKind(req.Kind), Status: domain.WorkOrderStatus(req.Status), Priority: domain.WorkOrderPriority(req.Priority), ServiceTypeID: req.ServiceTypeID, EstimatedHours: req.EstimatedHours}
}

func PreventivePlanFromRequest(req dto.PreventivePlanRequest) domain.PreventivePlan {
	return domain.PreventivePlan{Name: req.Name, Description: req.Description, AssetID: req.AssetID, VehicleID: req.VehicleID, ServiceTypeID: req.ServiceTypeID, Frequency: domain.PlanFrequency(req.Frequency), IntervalValue: req.IntervalValue}
}

func ServiceTypeFromRequest(req dto.ServiceTypeRequest) domain.ServiceType {
	return domain.ServiceType{Name: req.Name, Code: req.Code, Description: req.Description}
}
func CatalogFromRequest(req dto.CatalogRequest) domain.MaintenanceCatalog {
	return domain.MaintenanceCatalog{Name: req.Name, Code: req.Code, Description: req.Description}
}
func WorkshopFromRequest(req dto.WorkshopRequest) domain.Workshop {
	return domain.Workshop{Name: req.Name, Document: req.Document, Phone: req.Phone, Email: req.Email}
}
func TechnicianFromRequest(req dto.TechnicianRequest) domain.Technician {
	return domain.Technician{Name: req.Name, Document: req.Document, Phone: req.Phone, Email: req.Email, Active: req.Active}
}
func LaborFromRequest(req dto.LaborRequest) domain.LaborEntry {
	return domain.LaborEntry{Technician: req.Technician, Hours: req.Hours, Cost: req.Cost}
}
func DowntimeFromRequest(req dto.DowntimeRequest) domain.Downtime {
	return domain.Downtime{AssetID: req.AssetID, VehicleID: req.VehicleID, Reason: req.Reason}
}

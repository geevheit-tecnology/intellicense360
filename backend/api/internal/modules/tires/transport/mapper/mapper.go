package mapper

import (
	"github.com/geevheit/intelligence360/backend/api/internal/modules/tires/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/tires/transport/dto"
)

func TireFromRequest(req dto.TireRequest) domain.Tire {
	return domain.Tire{
		SerialNumber: req.SerialNumber, FireNumber: req.FireNumber, Brand: req.Brand, Model: req.Model,
		Size: req.Size, Construction: req.Construction, TireType: req.TireType, PositionType: req.PositionType,
		PurchaseValue: req.PurchaseValue, Supplier: req.Supplier, DOT: req.DOT, CurrentTreadMM: req.CurrentTreadMM,
		OriginalTreadMM: req.OriginalTreadMM, MinimumTreadMM: req.MinimumTreadMM, Status: domain.TireStatus(req.Status),
		VehicleID: req.VehicleID, Position: req.Position, CurrentKM: req.CurrentKM, TotalKM: req.TotalKM,
		RecapCount: req.RecapCount, Notes: req.Notes,
	}
}

func InspectionFromRequest(req dto.TireInspectionRequest) domain.TireInspection {
	return domain.TireInspection{TreadMM: req.TreadMM, Pressure: req.Pressure, Temperature: req.Temperature, Condition: req.Condition, Observations: req.Observations, Inspector: req.Inspector}
}

func MovementFromRequest(req dto.TireMovementRequest) domain.TireMovement {
	return domain.TireMovement{MovementType: domain.MovementType(req.MovementType), VehicleID: req.VehicleID, Position: req.Position, KM: req.KM, Reason: req.Reason}
}

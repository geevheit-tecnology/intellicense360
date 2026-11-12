package mapper

import (
	"github.com/geevheit/intelligence360/backend/api/internal/modules/trips/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/trips/transport/dto"
)

func TripFromRequest(req dto.TripRequest) domain.Trip {
	return domain.Trip{
		TripNumber:        req.TripNumber,
		VehicleID:         req.VehicleID,
		TrailerID:         req.TrailerID,
		DriverID:          req.DriverID,
		OriginCity:        req.OriginCity,
		OriginState:       req.OriginState,
		DestinationCity:   req.DestinationCity,
		DestinationState:  req.DestinationState,
		DistanceKM:        req.DistanceKM,
		DepartureAt:       req.DepartureAt,
		ExpectedArrivalAt: req.ExpectedArrivalAt,
		CargoDescription:  req.CargoDescription,
		CargoWeightKG:     req.CargoWeightKG,
		CargoValue:        req.CargoValue,
		CustomerName:      req.CustomerName,
		MDFEKey:           req.MDFEKey,
		CTEKey:            req.CTEKey,
		CIOTNumber:        req.CIOTNumber,
		FreightAmount:     req.FreightAmount,
		AdvanceAmount:     req.AdvanceAmount,
		TollAmount:        req.TollAmount,
		StartOdometerKM:   req.StartOdometerKM,
		Notes:             req.Notes,
	}
}


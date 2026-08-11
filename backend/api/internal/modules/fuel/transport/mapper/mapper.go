package mapper

import (
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/transport/dto"
)

func TransactionFromRequest(req dto.FuelTransactionRequest) domain.FuelTransaction {
	return domain.FuelTransaction{
		TransactionDate:   parseTime(req.TransactionDate),
		FuelTypeID:        req.FuelTypeID,
		FuelKind:          domain.FuelKind(req.FuelKind),
		Quantity:          req.Quantity,
		UnitPrice:         req.UnitPrice,
		TotalAmount:       req.TotalAmount,
		OdometerReading:   req.OdometerReading,
		EngineHourReading: req.EngineHourReading,
		StationID:         req.StationID,
		TankID:            req.TankID,
		NozzleID:          req.NozzleID,
		ReceiptID:         req.ReceiptID,
		ReceiptNumber:     req.ReceiptNumber,
		DriverReference:   req.DriverReference,
		VehicleReference:  req.VehicleReference,
		AssetReference:    req.AssetReference,
		PaymentMethod:     req.PaymentMethod,
		Notes:             req.Notes,
		Status:            domain.FuelTransactionStatus(req.Status),
	}
}
func TypeFromRequest(req dto.FuelTypeRequest) domain.FuelType {
	return domain.FuelType{Name: req.Name, Kind: domain.FuelKind(req.Kind), Code: req.Code, Description: req.Description, Active: req.Active}
}
func StationFromRequest(req dto.FuelStationRequest) domain.FuelStation {
	return domain.FuelStation{Name: req.Name, LegalName: req.LegalName, CNPJ: req.CNPJ, Address: req.Address, City: req.City, State: req.State, Country: req.Country, Latitude: req.Latitude, Longitude: req.Longitude, Active: req.Active, Notes: req.Notes}
}
func TankFromRequest(req dto.FuelTankRequest) domain.FuelTank {
	return domain.FuelTank{Code: req.Code, Name: req.Name, Capacity: req.Capacity, CurrentReading: req.CurrentReading, FuelTypeID: req.FuelTypeID, FuelKind: domain.FuelKind(req.FuelKind), LocationRef: req.LocationRef, Status: domain.FuelTankStatus(req.Status), Notes: req.Notes}
}
func NozzleFromRequest(req dto.FuelNozzleRequest) domain.FuelNozzle {
	return domain.FuelNozzle{Code: req.Code, FuelTypeID: req.FuelTypeID, FuelKind: domain.FuelKind(req.FuelKind), TankID: req.TankID, Status: domain.FuelNozzleStatus(req.Status), MeterReading: req.MeterReading, Notes: req.Notes}
}
func ReadingFromRequest(req dto.FuelReadingRequest) domain.FuelReading {
	return domain.FuelReading{ReadingType: domain.FuelReadingType(req.ReadingType), ReferenceID: req.ReferenceID, Value: req.Value, ReadingDate: parseTime(req.ReadingDate), Source: req.Source, Notes: req.Notes}
}
func PriceFromRequest(req dto.FuelPriceRequest) domain.FuelPrice {
	return domain.FuelPrice{FuelTypeID: req.FuelTypeID, FuelKind: domain.FuelKind(req.FuelKind), UnitPrice: req.UnitPrice, EffectiveDate: parseTime(req.EffectiveDate), StationID: req.StationID, Source: req.Source, Notes: req.Notes}
}
func ReceiptFromRequest(req dto.FuelReceiptRequest) domain.FuelReceipt {
	return domain.FuelReceipt{ReceiptNumber: req.ReceiptNumber, ReceiptDate: parseTime(req.ReceiptDate), Amount: req.Amount, AttachmentReference: req.AttachmentReference, Notes: req.Notes}
}
func AdjustmentFromRequest(req dto.FuelAdjustmentRequest) domain.FuelAdjustment {
	return domain.FuelAdjustment{TransactionID: req.TransactionID, AdjustmentType: req.AdjustmentType, Reason: req.Reason, OriginalReference: req.OriginalReference, AdjustedQuantity: req.AdjustedQuantity, AdjustedUnitPrice: req.AdjustedUnitPrice, AdjustedTotalAmount: req.AdjustedTotalAmount, Notes: req.Notes}
}
func parseTime(value string) time.Time {
	if value == "" {
		return time.Time{}
	}
	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return parsed
	}
	if parsed, err := time.Parse("2006-01-02", value); err == nil {
		return parsed
	}
	return time.Time{}
}

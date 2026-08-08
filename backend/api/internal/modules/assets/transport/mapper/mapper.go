package mapper

import (
	"github.com/geevheit/intelligence360/backend/api/internal/modules/assets/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/assets/transport/dto"
)

func AssetFromRequest(req dto.AssetRequest) domain.Asset {
	return domain.Asset{InternalCode: req.InternalCode, SerialNumber: req.SerialNumber, AssetTag: req.AssetTag, Name: req.Name, Description: req.Description, CategoryID: req.CategoryID, TypeID: req.TypeID, ManufacturerID: req.ManufacturerID, ModelID: req.ModelID, Status: domain.AssetStatus(req.Status), Ownership: domain.Ownership(req.Ownership), Metadata: req.Metadata}
}

func CategoryFromRequest(req dto.CategoryRequest) domain.AssetCategory {
	return domain.AssetCategory{Name: req.Name, Code: req.Code}
}
func TypeFromRequest(req dto.TypeRequest) domain.AssetType {
	return domain.AssetType{CategoryID: req.CategoryID, Name: req.Name, Code: req.Code}
}
func ManufacturerFromRequest(req dto.ManufacturerRequest) domain.Manufacturer {
	return domain.Manufacturer{Name: req.Name}
}
func ModelFromRequest(req dto.ModelRequest) domain.Model {
	return domain.Model{ManufacturerID: req.ManufacturerID, Name: req.Name}
}
func EquipmentFromRequest(req dto.EquipmentRequest) domain.Equipment {
	return domain.Equipment{AssetID: req.AssetID, Category: req.Category, Type: req.Type, Capacity: req.Capacity}
}

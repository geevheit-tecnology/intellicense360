package mapper

import (
	"github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/transport/dto"
)

func PartFromRequest(req dto.PartRequest) domain.Part {
	return domain.Part{SKU: req.SKU, InternalCode: req.InternalCode, Name: req.Name, Description: req.Description, CategoryID: req.CategoryID, BrandID: req.BrandID, ModelID: req.ModelID, UnitID: req.UnitID, Status: domain.StockStatus(req.Status), Metadata: req.Metadata}
}

func CatalogFromRequest(req dto.CatalogRequest) domain.Catalog {
	return domain.Catalog{Name: req.Name, Code: req.Code, Description: req.Description}
}

func WarehouseFromRequest(req dto.WarehouseRequest) domain.Warehouse {
	return domain.Warehouse{Name: req.Name, Code: req.Code, Address: req.Address, Active: req.Active}
}

func LocationFromRequest(req dto.LocationRequest) domain.WarehouseLocation {
	return domain.WarehouseLocation{WarehouseID: req.WarehouseID, Name: req.Name, Code: req.Code}
}

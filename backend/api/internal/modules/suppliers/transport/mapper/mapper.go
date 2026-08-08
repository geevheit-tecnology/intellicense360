package mapper

import (
	"github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/transport/dto"
)

func SupplierFromRequest(req dto.SupplierRequest) domain.Supplier {
	return domain.Supplier{LegalName: req.LegalName, TradeName: req.TradeName, CNPJ: req.CNPJ, StateRegistration: req.StateRegistration, MunicipalRegistration: req.MunicipalRegistration, Email: req.Email, Phone: req.Phone, Website: req.Website, Notes: req.Notes, Status: domain.SupplierStatus(req.Status), CategoryID: req.CategoryID, Type: domain.SupplierTypeCode(req.Type), Metadata: req.Metadata}
}
func CategoryFromRequest(req dto.CategoryRequest) domain.SupplierCategory {
	return domain.SupplierCategory{Name: req.Name, Code: req.Code, Description: req.Description}
}
func TypeFromRequest(req dto.TypeRequest) domain.SupplierType {
	return domain.SupplierType{Name: req.Name, Code: req.Code, Description: req.Description}
}
func ContactFromRequest(req dto.ContactRequest) domain.SupplierContact {
	return domain.SupplierContact{SupplierID: req.SupplierID, Name: req.Name, Role: req.Role, Email: req.Email, Phone: req.Phone, Mobile: req.Mobile, Primary: req.Primary}
}
func AddressFromRequest(req dto.AddressRequest) domain.SupplierAddress {
	return domain.SupplierAddress{SupplierID: req.SupplierID, Street: req.Street, Number: req.Number, Complement: req.Complement, Neighborhood: req.Neighborhood, City: req.City, State: req.State, PostalCode: req.PostalCode, Country: req.Country, AddressType: req.AddressType}
}
func DocumentFromRequest(req dto.DocumentRequest) domain.SupplierDocument {
	return domain.SupplierDocument{SupplierID: req.SupplierID, DocumentType: req.DocumentType, DocumentNumber: req.DocumentNumber, Status: req.Status, AttachmentReference: req.AttachmentReference}
}
func ContractFromRequest(req dto.ContractRequest) domain.SupplierContract {
	return domain.SupplierContract{SupplierID: req.SupplierID, ContractNumber: req.ContractNumber, Status: req.Status, Notes: req.Notes, AttachmentReference: req.AttachmentReference}
}
func RatingFromRequest(req dto.RatingRequest) domain.SupplierRating {
	return domain.SupplierRating{SupplierID: req.SupplierID, Quality: req.Quality, Price: req.Price, Delivery: req.Delivery, Service: req.Service, Reliability: req.Reliability, OverallScore: req.OverallScore}
}

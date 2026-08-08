package dto

type AssetRequest struct {
	InternalCode   string            `json:"internal_code"`
	SerialNumber   string            `json:"serial_number"`
	AssetTag       string            `json:"asset_tag"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	CategoryID     string            `json:"category_id"`
	TypeID         string            `json:"type_id"`
	ManufacturerID string            `json:"manufacturer_id"`
	ModelID        string            `json:"model_id"`
	Status         string            `json:"status"`
	Ownership      string            `json:"ownership"`
	Metadata       map[string]string `json:"metadata"`
}

type CategoryRequest struct {
	Name string `json:"name"`
	Code string `json:"code"`
}
type TypeRequest struct {
	CategoryID string `json:"category_id"`
	Name       string `json:"name"`
	Code       string `json:"code"`
}
type ManufacturerRequest struct {
	Name string `json:"name"`
}
type ModelRequest struct {
	ManufacturerID string `json:"manufacturer_id"`
	Name           string `json:"name"`
}
type EquipmentRequest struct {
	AssetID  string `json:"asset_id"`
	Category string `json:"category"`
	Type     string `json:"type"`
	Capacity string `json:"capacity"`
}

package dto

type PartRequest struct {
	SKU          string            `json:"sku"`
	InternalCode string            `json:"internal_code"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	CategoryID   string            `json:"category_id"`
	BrandID      string            `json:"brand_id"`
	ModelID      string            `json:"model_id"`
	UnitID       string            `json:"unit_id"`
	Status       string            `json:"status"`
	Metadata     map[string]string `json:"metadata"`
}

type CatalogRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type WarehouseRequest struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	Address string `json:"address"`
	Active  bool   `json:"active"`
}

type LocationRequest struct {
	WarehouseID string `json:"warehouse_id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
}

package dto

type SupplierRequest struct {
	LegalName             string            `json:"legal_name"`
	TradeName             string            `json:"trade_name"`
	CNPJ                  string            `json:"cnpj"`
	StateRegistration     string            `json:"state_registration"`
	MunicipalRegistration string            `json:"municipal_registration"`
	Email                 string            `json:"email"`
	Phone                 string            `json:"phone"`
	Website               string            `json:"website"`
	Notes                 string            `json:"notes"`
	Status                string            `json:"status"`
	CategoryID            string            `json:"category_id"`
	Type                  string            `json:"type"`
	Metadata              map[string]string `json:"metadata"`
}

type CategoryRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type TypeRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type ContactRequest struct {
	SupplierID string `json:"supplier_id"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Mobile     string `json:"mobile"`
	Primary    bool   `json:"primary"`
}

type AddressRequest struct {
	SupplierID   string `json:"supplier_id"`
	Street       string `json:"street"`
	Number       string `json:"number"`
	Complement   string `json:"complement"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
	AddressType  string `json:"address_type"`
}

type DocumentRequest struct {
	SupplierID          string `json:"supplier_id"`
	DocumentType        string `json:"document_type"`
	DocumentNumber      string `json:"document_number"`
	Status              string `json:"status"`
	AttachmentReference string `json:"attachment_reference"`
}

type ContractRequest struct {
	SupplierID          string `json:"supplier_id"`
	ContractNumber      string `json:"contract_number"`
	Status              string `json:"status"`
	Notes               string `json:"notes"`
	AttachmentReference string `json:"attachment_reference"`
}

type RatingRequest struct {
	SupplierID   string  `json:"supplier_id"`
	Quality      float64 `json:"quality"`
	Price        float64 `json:"price"`
	Delivery     float64 `json:"delivery"`
	Service      float64 `json:"service"`
	Reliability  float64 `json:"reliability"`
	OverallScore float64 `json:"overall_score"`
}

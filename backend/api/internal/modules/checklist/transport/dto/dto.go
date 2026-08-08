package dto

type ChecklistRequest struct {
	VehicleID      string `json:"vehicle_id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Type           string `json:"type"`
	Status         string `json:"status"`
	DriverName     string `json:"driver_name"`
	DriverDocument string `json:"driver_document"`
}

type ChecklistItemRequest struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	Category      string `json:"category"`
	Required      bool   `json:"required"`
	OrderIndex    int    `json:"order_index"`
	AnswerType    string `json:"answer_type"`
	ExpectedValue string `json:"expected_value"`
}

type ChecklistAnswerRequest struct {
	ChecklistItemID string `json:"checklist_item_id"`
	Answer          string `json:"answer"`
	Notes           string `json:"notes"`
	PhotoURL        string `json:"photo_url"`
}

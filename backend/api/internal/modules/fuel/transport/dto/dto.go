package dto

type FuelTransactionRequest struct {
	TransactionDate   string  `json:"transaction_date"`
	FuelTypeID        string  `json:"fuel_type_id"`
	FuelKind          string  `json:"fuel_kind"`
	Quantity          float64 `json:"quantity"`
	UnitPrice         float64 `json:"unit_price"`
	TotalAmount       float64 `json:"total_amount"`
	OdometerReading   float64 `json:"odometer_reading"`
	EngineHourReading float64 `json:"engine_hour_reading"`
	StationID         string  `json:"station_id"`
	TankID            string  `json:"tank_id"`
	NozzleID          string  `json:"nozzle_id"`
	ReceiptID         string  `json:"receipt_id"`
	ReceiptNumber     string  `json:"receipt_number"`
	DriverReference   string  `json:"driver_reference"`
	VehicleReference  string  `json:"vehicle_reference"`
	AssetReference    string  `json:"asset_reference"`
	PaymentMethod     string  `json:"payment_method"`
	Notes             string  `json:"notes"`
	Status            string  `json:"status"`
}

type FuelTypeRequest struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type FuelStationRequest struct {
	Name      string  `json:"name"`
	LegalName string  `json:"legal_name"`
	CNPJ      string  `json:"cnpj"`
	Address   string  `json:"address"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Active    bool    `json:"active"`
	Notes     string  `json:"notes"`
}

type FuelTankRequest struct {
	Code           string  `json:"code"`
	Name           string  `json:"name"`
	Capacity       float64 `json:"capacity"`
	CurrentReading float64 `json:"current_reading"`
	FuelTypeID     string  `json:"fuel_type_id"`
	FuelKind       string  `json:"fuel_kind"`
	LocationRef    string  `json:"location_reference"`
	Status         string  `json:"status"`
	Notes          string  `json:"notes"`
}

type FuelNozzleRequest struct {
	Code         string  `json:"code"`
	FuelTypeID   string  `json:"fuel_type_id"`
	FuelKind     string  `json:"fuel_kind"`
	TankID       string  `json:"tank_id"`
	Status       string  `json:"status"`
	MeterReading float64 `json:"meter_reading"`
	Notes        string  `json:"notes"`
}

type FuelReadingRequest struct {
	ReadingType string  `json:"reading_type"`
	ReferenceID string  `json:"reference_id"`
	Value       float64 `json:"value"`
	ReadingDate string  `json:"reading_date"`
	Source      string  `json:"source"`
	Notes       string  `json:"notes"`
}

type FuelPriceRequest struct {
	FuelTypeID    string  `json:"fuel_type_id"`
	FuelKind      string  `json:"fuel_kind"`
	UnitPrice     float64 `json:"unit_price"`
	EffectiveDate string  `json:"effective_date"`
	StationID     string  `json:"station_id"`
	Source        string  `json:"source"`
	Notes         string  `json:"notes"`
}

type FuelReceiptRequest struct {
	ReceiptNumber       string  `json:"receipt_number"`
	ReceiptDate         string  `json:"receipt_date"`
	Amount              float64 `json:"amount"`
	AttachmentReference string  `json:"attachment_reference"`
	Notes               string  `json:"notes"`
}

type FuelAdjustmentRequest struct {
	TransactionID       string  `json:"transaction_id"`
	AdjustmentType      string  `json:"adjustment_type"`
	Reason              string  `json:"reason"`
	OriginalReference   string  `json:"original_reference"`
	AdjustedQuantity    float64 `json:"adjusted_quantity"`
	AdjustedUnitPrice   float64 `json:"adjusted_unit_price"`
	AdjustedTotalAmount float64 `json:"adjusted_total_amount"`
	Notes               string  `json:"notes"`
}

type FuelCancelRequest struct {
	Reason string `json:"reason"`
}

package domain

import "time"

type TransactionKind string
type CostClassification string
type FinancialStatus string

const (
	KindExpense    TransactionKind = "expense"
	KindRevenue    TransactionKind = "revenue"
	KindAdjustment TransactionKind = "adjustment"
	KindTransfer   TransactionKind = "transfer"
	KindRefund     TransactionKind = "refund"
	KindOther      TransactionKind = "other"

	CostFixed          CostClassification = "fixed"
	CostVariable       CostClassification = "variable"
	CostOperational    CostClassification = "operational"
	CostAdministrative CostClassification = "administrative"
	CostFinancial      CostClassification = "financial"
	CostExtraordinary  CostClassification = "extraordinary"
	CostOther          CostClassification = "other"

	StatusDraft    FinancialStatus = "draft"
	StatusPending  FinancialStatus = "pending"
	StatusApproved FinancialStatus = "approved"
	StatusPaid     FinancialStatus = "paid"
	StatusReceived FinancialStatus = "received"
	StatusCanceled FinancialStatus = "canceled"
	StatusOverdue  FinancialStatus = "overdue"
	StatusAdjusted FinancialStatus = "adjusted"
	StatusClosed   FinancialStatus = "closed"
)

type FinancialTransaction struct {
	ID                  string          `json:"id"`
	TenantID            string          `json:"tenant_id"`
	Kind                TransactionKind `json:"kind"`
	Description         string          `json:"description"`
	Amount              float64         `json:"amount"`
	Date                time.Time       `json:"date"`
	DueDate             *time.Time      `json:"due_date,omitempty"`
	SettlementDate      *time.Time      `json:"settlement_date,omitempty"`
	CategoryID          string          `json:"category_id,omitempty"`
	CostTypeID          string          `json:"cost_type_id,omitempty"`
	CostCenterID        string          `json:"cost_center_id,omitempty"`
	AccountID           string          `json:"account_id,omitempty"`
	SupplierReference   string          `json:"supplier_reference,omitempty"`
	OperationReference  string          `json:"operation_reference,omitempty"`
	DocumentNumber      string          `json:"document_number,omitempty"`
	PaymentMethodID     string          `json:"payment_method_id,omitempty"`
	Status              FinancialStatus `json:"status"`
	Notes               string          `json:"notes,omitempty"`
	AttachmentReference string          `json:"attachment_reference,omitempty"`
	CreatedBy           string          `json:"created_by,omitempty"`
	UpdatedBy           string          `json:"updated_by,omitempty"`
	CreatedAt           time.Time       `json:"created_at"`
	UpdatedAt           time.Time       `json:"updated_at"`
	DeletedAt           *time.Time      `json:"deleted_at,omitempty"`
	Version             int64           `json:"version"`
}

type Expense = FinancialTransaction
type Revenue = FinancialTransaction

type CostCategory struct {
	ID             string             `json:"id"`
	TenantID       string             `json:"tenant_id"`
	Name           string             `json:"name"`
	Code           string             `json:"code,omitempty"`
	Description    string             `json:"description,omitempty"`
	Classification CostClassification `json:"classification,omitempty"`
	Active         bool               `json:"active"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	DeletedAt      *time.Time         `json:"deleted_at,omitempty"`
	Version        int64              `json:"version"`
}

type CostType struct {
	ID             string             `json:"id"`
	TenantID       string             `json:"tenant_id"`
	Name           string             `json:"name"`
	Code           string             `json:"code,omitempty"`
	Description    string             `json:"description,omitempty"`
	Classification CostClassification `json:"classification,omitempty"`
	Active         bool               `json:"active"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	DeletedAt      *time.Time         `json:"deleted_at,omitempty"`
	Version        int64              `json:"version"`
}

type CostCenter struct {
	ID        string     `json:"id"`
	TenantID  string     `json:"tenant_id"`
	ParentID  string     `json:"parent_id,omitempty"`
	Name      string     `json:"name"`
	Code      string     `json:"code,omitempty"`
	Notes     string     `json:"notes,omitempty"`
	Active    bool       `json:"active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
	Version   int64      `json:"version"`
}

type Account struct {
	ID          string     `json:"id"`
	TenantID    string     `json:"tenant_id"`
	ParentID    string     `json:"parent_id,omitempty"`
	AccountCode string     `json:"account_code,omitempty"`
	Name        string     `json:"name"`
	Type        string     `json:"type,omitempty"`
	Status      string     `json:"status"`
	Notes       string     `json:"notes,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Version     int64      `json:"version"`
}

type PaymentMethod struct {
	ID         string     `json:"id"`
	TenantID   string     `json:"tenant_id"`
	Name       string     `json:"name"`
	Code       string     `json:"code,omitempty"`
	MethodType string     `json:"method_type,omitempty"`
	Notes      string     `json:"notes,omitempty"`
	Active     bool       `json:"active"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	Version    int64      `json:"version"`
}

type FinancialPeriod struct {
	ID        string          `json:"id"`
	TenantID  string          `json:"tenant_id"`
	Year      int             `json:"year"`
	Month     int             `json:"month"`
	StartDate time.Time       `json:"start_date"`
	EndDate   time.Time       `json:"end_date"`
	Status    FinancialStatus `json:"status"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
	DeletedAt *time.Time      `json:"deleted_at,omitempty"`
	Version   int64           `json:"version"`
}

type Budget struct {
	ID            string     `json:"id"`
	TenantID      string     `json:"tenant_id"`
	PeriodID      string     `json:"period_id,omitempty"`
	CostCenterID  string     `json:"cost_center_id,omitempty"`
	CategoryID    string     `json:"category_id,omitempty"`
	Name          string     `json:"name"`
	Status        string     `json:"status"`
	Notes         string     `json:"notes,omitempty"`
	PlannedAmount float64    `json:"planned_amount"`
	ActualAmount  float64    `json:"actual_amount"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
	Version       int64      `json:"version"`
}

type BudgetItem struct {
	ID            string     `json:"id"`
	TenantID      string     `json:"tenant_id"`
	BudgetID      string     `json:"budget_id"`
	CategoryID    string     `json:"category_id,omitempty"`
	CostCenterID  string     `json:"cost_center_id,omitempty"`
	Description   string     `json:"description"`
	PlannedAmount float64    `json:"planned_amount"`
	ActualAmount  float64    `json:"actual_amount"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`
	Version       int64      `json:"version"`
}

type FinancialAttachment struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenant_id"`
	TransactionID string    `json:"transaction_id"`
	Kind          string    `json:"kind"`
	FileName      string    `json:"file_name"`
	MimeType      string    `json:"mime_type,omitempty"`
	Reference     string    `json:"reference"`
	CreatedAt     time.Time `json:"created_at"`
}

type FinancialAdjustment struct {
	ID             string    `json:"id"`
	TenantID       string    `json:"tenant_id"`
	TransactionID  string    `json:"transaction_id"`
	AdjustmentType string    `json:"adjustment_type"`
	Reason         string    `json:"reason"`
	Notes          string    `json:"notes,omitempty"`
	CreatedBy      string    `json:"created_by,omitempty"`
	AdjustedAmount float64   `json:"adjusted_amount"`
	CreatedAt      time.Time `json:"created_at"`
}

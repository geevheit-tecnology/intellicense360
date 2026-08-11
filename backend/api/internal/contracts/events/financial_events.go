package events

const (
	FinancialTransactionCreatedEventName = "financial.transaction.created"
	ExpenseCreatedEventName              = "financial.expense.created"
	ExpenseApprovedEventName             = "financial.expense.approved"
	ExpensePaidEventName                 = "financial.expense.paid"
	RevenueCreatedEventName              = "financial.revenue.created"
	RevenueReceivedEventName             = "financial.revenue.received"
	FinancialAdjustmentCreatedEventName  = "financial.adjustment.created"
	FinancialPeriodClosedEventName       = "financial.period.closed"
	BudgetCreatedEventName               = "financial.budget.created"
)

type FinancialTransactionCreated struct {
	EventMetadata
	Transaction AggregateMetadata `json:"transaction"`
	Kind        string            `json:"kind"`
	Amount      float64           `json:"amount"`
	Status      string            `json:"status"`
}

type ExpenseCreated struct{ FinancialTransactionCreated }
type ExpenseApproved struct{ EventMetadata }
type ExpensePaid struct{ EventMetadata }
type RevenueCreated struct{ FinancialTransactionCreated }
type RevenueReceived struct{ EventMetadata }

type FinancialAdjustmentCreated struct {
	EventMetadata
	TransactionID  string  `json:"transaction_id"`
	AdjustmentType string  `json:"adjustment_type"`
	AdjustedAmount float64 `json:"adjusted_amount"`
}

type FinancialPeriodClosed struct {
	EventMetadata
	PeriodID string `json:"period_id"`
	Year     int    `json:"year"`
	Month    int    `json:"month"`
}

type BudgetCreated struct {
	EventMetadata
	BudgetID      string  `json:"budget_id"`
	PlannedAmount float64 `json:"planned_amount"`
}

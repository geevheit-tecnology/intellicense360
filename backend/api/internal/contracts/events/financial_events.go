package events

const (
	FinancialTransactionCreatedEventName = "financial.transaction.created.v1"
	ExpenseCreatedEventName              = "financial.expense.created.v1"
	ExpenseApprovedEventName             = "financial.expense.approved.v1"
	ExpensePaidEventName                 = "financial.expense.paid.v1"
	RevenueCreatedEventName              = "financial.revenue.created.v1"
	RevenueReceivedEventName             = "financial.revenue.received.v1"
	FinancialAdjustmentCreatedEventName  = "financial.adjustment.created.v1"
	FinancialPeriodClosedEventName       = "financial.period.closed.v1"
	BudgetCreatedEventName               = "financial.budget.created.v1"
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

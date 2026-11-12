package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/financial/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/financial/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/financial/ports"
	platform "github.com/geevheit/intelligence360/backend/api/internal/platform/postgres"
)

type Store struct{ db *sql.DB }

func NewStore(db *sql.DB) *Store { return &Store{db: db} }

func mapError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return application.ErrNotFound
	}
	var pg *pq.Error
	if errors.As(err, &pg) {
		if pg.Code == "23505" {
			return application.ErrDuplicateCode
		}
		if pg.Code == "23503" || pg.Code == "23502" || pg.Code == "23514" || pg.Code == "22P02" {
			return application.ErrInvalidFinancialData
		}
	}
	return err
}

func read[T any](row *sql.Row) (T, error) {
	var item T
	var body []byte
	if err := row.Scan(&body); err != nil {
		return item, mapError(err)
	}
	err := json.Unmarshal(body, &item)
	return item, err
}

func changed(result sql.Result, err error) error {
	if err != nil {
		return mapError(err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return application.ErrNotFound
	}
	return nil
}

func (s *Store) Transactions() ports.TransactionRepository {
	return TransactionRepository{db: s.db}
}
func (s *Store) Categories() ports.CatalogRepository[domain.CostCategory] {
	return CategoryRepository{db: s.db}
}
func (s *Store) Types() ports.CatalogRepository[domain.CostType] {
	return TypeRepository{db: s.db}
}
func (s *Store) Centers() ports.CatalogRepository[domain.CostCenter] {
	return CenterRepository{db: s.db}
}
func (s *Store) Accounts() ports.CatalogRepository[domain.Account] {
	return AccountRepository{db: s.db}
}
func (s *Store) Methods() ports.CatalogRepository[domain.PaymentMethod] {
	return MethodRepository{db: s.db}
}
func (s *Store) Periods() ports.PeriodRepository {
	return PeriodRepository{db: s.db}
}
func (s *Store) Budgets() ports.BudgetRepository {
	return BudgetRepository{db: s.db}
}
func (s *Store) Adjustments() ports.AdjustmentRepository {
	return AdjustmentRepository{db: s.db}
}
func (s *Store) Closings() ports.CashClosingRepository {
	return CashClosingRepository{db: s.db}
}

// TransactionRepository implements ports.TransactionRepository
type TransactionRepository struct{ db *sql.DB }

func (r TransactionRepository) Create(ctx context.Context, item domain.FinancialTransaction) (domain.FinancialTransaction, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	query := `INSERT INTO financial_transactions (
		id, tenant_id, kind, description, amount, date, due_date, settlement_date,
		category_id, cost_type_id, cost_center_id, account_id, supplier_reference,
		operation_reference, document_number, payment_method_id, vehicle_id, driver_id,
		fuel_transaction_id, work_order_id, invoice_key, invoice_number, requisition_number,
		is_third_party, third_party_identifier, payment_method, barcode, bank_account,
		cash_flow_nature, status, notes, attachment_reference, version, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8,
		nullif($9,'')::uuid, nullif($10,'')::uuid, nullif($11,'')::uuid, nullif($12,'')::uuid, $13,
		$14, $15, nullif($16,'')::uuid, nullif($17,'')::uuid, nullif($18,'')::uuid,
		nullif($19,'')::uuid, nullif($20,'')::uuid, $21, $22, $23,
		$24, $25, $26, $27, $28,
		$29, $30, $31, $32, $33, $34, $35
	)`
	_, err := exec.ExecContext(ctx, query,
		item.ID, item.TenantID, string(item.Kind), item.Description, item.Amount, item.Date, item.DueDate, item.SettlementDate,
		item.CategoryID, item.CostTypeID, item.CostCenterID, item.AccountID, item.SupplierReference,
		item.OperationReference, item.DocumentNumber, item.PaymentMethodID, item.VehicleID, item.DriverID,
		item.FuelTransactionID, item.WorkOrderID, item.InvoiceKey, item.InvoiceNumber, item.RequisitionNumber,
		item.IsThirdParty, item.ThirdPartyIdentifier, item.PaymentMethod, item.Barcode, item.BankAccount,
		item.CashFlowNature, string(item.Status), item.Notes, item.AttachmentReference, item.Version, item.CreatedAt, item.UpdatedAt,
	)
	if err != nil {
		return domain.FinancialTransaction{}, mapError(err)
	}
	return item, nil
}

func (r TransactionRepository) FindByID(ctx context.Context, tenantID, id string) (domain.FinancialTransaction, error) {
	query := `SELECT row_to_json(t) FROM financial_transactions t WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL`
	return read[domain.FinancialTransaction](r.db.QueryRowContext(ctx, query, tenantID, id))
}

func (r TransactionRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FinancialTransaction], error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 20
	}
	args := []any{tenantID}
	where := "tenant_id = $1 AND deleted_at IS NULL"
	if q.Search != "" {
		args = append(args, q.Search)
		where += fmt.Sprintf(" AND (description ILIKE '%%' || $%d || '%%' OR invoice_number ILIKE '%%' || $%d || '%%' OR requisition_number ILIKE '%%' || $%d || '%%')", len(args), len(args), len(args))
	}
	for k, v := range q.Filters {
		if v == "" {
			continue
		}
		switch k {
		case "kind":
			args = append(args, v)
			where += fmt.Sprintf(" AND kind = $%d", len(args))
		case "status":
			args = append(args, v)
			where += fmt.Sprintf(" AND status = $%d", len(args))
		case "vehicle_id":
			args = append(args, v)
			where += fmt.Sprintf(" AND vehicle_id = $%d::uuid", len(args))
		case "is_third_party":
			if v == "true" {
				where += " AND is_third_party = true"
			} else if v == "false" {
				where += " AND is_third_party = false"
			}
		}
	}

	var count int
	if err := r.db.QueryRowContext(ctx, "SELECT count(*) FROM financial_transactions WHERE "+where, args...).Scan(&count); err != nil {
		return ports.Page[domain.FinancialTransaction]{}, mapError(err)
	}

	limit := q.PageSize
	offset := (q.Page - 1) * q.PageSize
	order := "date DESC, created_at DESC"
	query := fmt.Sprintf("SELECT coalesce(json_agg(row_to_json(t)), '[]'::json) FROM (SELECT * FROM financial_transactions WHERE %s ORDER BY %s LIMIT %d OFFSET %d) t", where, order, limit, offset)

	var body []byte
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&body); err != nil {
		return ports.Page[domain.FinancialTransaction]{}, mapError(err)
	}
	items := []domain.FinancialTransaction{}
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.FinancialTransaction]{}, err
	}
	pages := (count + q.PageSize - 1) / q.PageSize
	return ports.Page[domain.FinancialTransaction]{Items: items, Page: q.Page, PageSize: q.PageSize, TotalItems: count, TotalPages: pages}, nil
}

func (r TransactionRepository) Update(ctx context.Context, item domain.FinancialTransaction) (domain.FinancialTransaction, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	query := `UPDATE financial_transactions SET
		description = $3, amount = $4, date = $5, due_date = $6, settlement_date = $7,
		category_id = nullif($8,'')::uuid, cost_type_id = nullif($9,'')::uuid,
		cost_center_id = nullif($10,'')::uuid, account_id = nullif($11,'')::uuid,
		supplier_reference = $12, operation_reference = $13, document_number = $14,
		payment_method_id = nullif($15,'')::uuid, vehicle_id = nullif($16,'')::uuid,
		driver_id = nullif($17,'')::uuid, invoice_key = $18, invoice_number = $19,
		requisition_number = $20, is_third_party = $21, third_party_identifier = $22,
		payment_method = $23, barcode = $24, bank_account = $25, cash_flow_nature = $26,
		status = $27, notes = $28, attachment_reference = $29, version = version + 1, updated_at = $30
	WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL`
	res, err := exec.ExecContext(ctx, query,
		item.TenantID, item.ID, item.Description, item.Amount, item.Date, item.DueDate, item.SettlementDate,
		item.CategoryID, item.CostTypeID, item.CostCenterID, item.AccountID,
		item.SupplierReference, item.OperationReference, item.DocumentNumber,
		item.PaymentMethodID, item.VehicleID, item.DriverID, item.InvoiceKey, item.InvoiceNumber,
		item.RequisitionNumber, item.IsThirdParty, item.ThirdPartyIdentifier,
		item.PaymentMethod, item.Barcode, item.BankAccount, item.CashFlowNature,
		string(item.Status), item.Notes, item.AttachmentReference, time.Now().UTC(),
	)
	if err := changed(res, err); err != nil {
		return domain.FinancialTransaction{}, err
	}
	return item, nil
}

func (r TransactionRepository) Delete(ctx context.Context, tenantID, id string) error {
	exec := platform.ExecutorFor(ctx, r.db)
	res, err := exec.ExecContext(ctx, `UPDATE financial_transactions SET deleted_at = now(), version = version + 1 WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL`, tenantID, id)
	return changed(res, err)
}

func (r TransactionRepository) GetCashFlowSummary(ctx context.Context, tenantID, startDate, endDate string) (domain.CashFlowSummary, error) {
	query := `SELECT
		coalesce(sum(case when kind = 'revenue' and status = 'received' then amount else 0 end), 0),
		coalesce(sum(case when kind = 'expense' and status = 'paid' then amount else 0 end), 0),
		coalesce(sum(case when kind = 'expense' and status = 'paid' and is_third_party = false then amount else 0 end), 0),
		coalesce(sum(case when kind = 'expense' and status = 'paid' and is_third_party = true then amount else 0 end), 0),
		coalesce(sum(case when kind = 'expense' and status in ('pending', 'approved') then amount else 0 end), 0),
		coalesce(sum(case when kind = 'expense' and status in ('pending', 'approved') and due_date < current_date then amount else 0 end), 0),
		coalesce(count(case when kind = 'expense' and status in ('pending', 'approved') then 1 end), 0)
	FROM financial_transactions
	WHERE tenant_id = $1 AND deleted_at IS NULL
		AND ($2 = '' OR date >= $2::date)
		AND ($3 = '' OR date <= $3::date)`

	summary := domain.CashFlowSummary{}
	err := r.db.QueryRowContext(ctx, query, tenantID, startDate, endDate).Scan(
		&summary.TotalRevenues,
		&summary.TotalExpenses,
		&summary.TotalFleetExpenses,
		&summary.TotalThirdPartyExpenses,
		&summary.PendingPayables,
		&summary.OverduePayables,
		&summary.PendingInvoicesCount,
	)
	if err != nil {
		return domain.CashFlowSummary{}, mapError(err)
	}
	summary.NetCashFlow = summary.TotalRevenues - summary.TotalExpenses
	if startDate != "" {
		if t, err := time.Parse("2006-01-02", startDate); err == nil {
			summary.StartDate = t
		}
	}
	if endDate != "" {
		if t, err := time.Parse("2006-01-02", endDate); err == nil {
			summary.EndDate = t
		}
	}
	return summary, nil
}

// CashClosingRepository implements ports.CashClosingRepository
type CashClosingRepository struct{ db *sql.DB }

func (r CashClosingRepository) Create(ctx context.Context, item domain.CashClosing) (domain.CashClosing, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	query := `INSERT INTO financial_cash_closings (
		id, tenant_id, closing_date, opening_balance, total_revenues, total_fleet_expenses,
		total_third_party_expenses, total_other_expenses, closing_balance, status, notes,
		closed_by, closed_at, version, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6,
		$7, $8, $9, $10, $11,
		nullif($12,'')::uuid, $13, $14, $15, $16
	)`
	_, err := exec.ExecContext(ctx, query,
		item.ID, item.TenantID, item.ClosingDate, item.OpeningBalance, item.TotalRevenues, item.TotalFleetExpenses,
		item.TotalThirdPartyExpenses, item.TotalOtherExpenses, item.ClosingBalance, item.Status, item.Notes,
		item.ClosedBy, item.ClosedAt, item.Version, item.CreatedAt, item.UpdatedAt,
	)
	if err != nil {
		return domain.CashClosing{}, mapError(err)
	}
	return item, nil
}

func (r CashClosingRepository) FindByDate(ctx context.Context, tenantID, dateStr string) (domain.CashClosing, error) {
	query := `SELECT row_to_json(c) FROM financial_cash_closings c WHERE tenant_id = $1 AND closing_date = $2::date AND deleted_at IS NULL`
	return read[domain.CashClosing](r.db.QueryRowContext(ctx, query, tenantID, dateStr))
}

func (r CashClosingRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.CashClosing], error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 20
	}
	args := []any{tenantID}
	where := "tenant_id = $1 AND deleted_at IS NULL"
	var count int
	if err := r.db.QueryRowContext(ctx, "SELECT count(*) FROM financial_cash_closings WHERE "+where, args...).Scan(&count); err != nil {
		return ports.Page[domain.CashClosing]{}, mapError(err)
	}

	limit := q.PageSize
	offset := (q.Page - 1) * q.PageSize
	order := "closing_date DESC"
	query := fmt.Sprintf("SELECT coalesce(json_agg(row_to_json(c)), '[]'::json) FROM (SELECT * FROM financial_cash_closings WHERE %s ORDER BY %s LIMIT %d OFFSET %d) c", where, order, limit, offset)

	var body []byte
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&body); err != nil {
		return ports.Page[domain.CashClosing]{}, mapError(err)
	}
	items := []domain.CashClosing{}
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.CashClosing]{}, err
	}
	pages := (count + q.PageSize - 1) / q.PageSize
	return ports.Page[domain.CashClosing]{Items: items, Page: q.Page, PageSize: q.PageSize, TotalItems: count, TotalPages: pages}, nil
}

// Catalogs: Categories, Types, Centers, Accounts, Methods
type CategoryRepository struct{ db *sql.DB }

func (r CategoryRepository) Create(ctx context.Context, item domain.CostCategory) (domain.CostCategory, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	query := `INSERT INTO financial_categories (id, tenant_id, name, code, description, classification, active, version, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := exec.ExecContext(ctx, query, item.ID, item.TenantID, item.Name, item.Code, item.Description, string(item.Classification), item.Active, item.Version, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return domain.CostCategory{}, mapError(err)
	}
	return item, nil
}
func (r CategoryRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.CostCategory], error) {
	return searchCatalog[domain.CostCategory](ctx, r.db, "financial_categories", tenantID, q)
}

type TypeRepository struct{ db *sql.DB }

func (r TypeRepository) Create(ctx context.Context, item domain.CostType) (domain.CostType, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	query := `INSERT INTO financial_cost_types (id, tenant_id, name, code, description, classification, active, version, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := exec.ExecContext(ctx, query, item.ID, item.TenantID, item.Name, item.Code, item.Description, string(item.Classification), item.Active, item.Version, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return domain.CostType{}, mapError(err)
	}
	return item, nil
}
func (r TypeRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.CostType], error) {
	return searchCatalog[domain.CostType](ctx, r.db, "financial_cost_types", tenantID, q)
}

type CenterRepository struct{ db *sql.DB }

func (r CenterRepository) Create(ctx context.Context, item domain.CostCenter) (domain.CostCenter, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	query := `INSERT INTO financial_cost_centers (id, tenant_id, parent_id, name, code, notes, active, version, created_at, updated_at)
	VALUES ($1, $2, nullif($3,'')::uuid, $4, $5, $6, $7, $8, $9, $10)`
	_, err := exec.ExecContext(ctx, query, item.ID, item.TenantID, item.ParentID, item.Name, item.Code, item.Notes, item.Active, item.Version, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return domain.CostCenter{}, mapError(err)
	}
	return item, nil
}
func (r CenterRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.CostCenter], error) {
	return searchCatalog[domain.CostCenter](ctx, r.db, "financial_cost_centers", tenantID, q)
}

type AccountRepository struct{ db *sql.DB }

func (r AccountRepository) Create(ctx context.Context, item domain.Account) (domain.Account, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	query := `INSERT INTO financial_accounts (id, tenant_id, parent_id, account_code, name, type, status, notes, version, created_at, updated_at)
	VALUES ($1, $2, nullif($3,'')::uuid, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := exec.ExecContext(ctx, query, item.ID, item.TenantID, item.ParentID, item.AccountCode, item.Name, item.Type, item.Status, item.Notes, item.Version, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return domain.Account{}, mapError(err)
	}
	return item, nil
}
func (r AccountRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Account], error) {
	return searchCatalog[domain.Account](ctx, r.db, "financial_accounts", tenantID, q)
}

type MethodRepository struct{ db *sql.DB }

func (r MethodRepository) Create(ctx context.Context, item domain.PaymentMethod) (domain.PaymentMethod, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	query := `INSERT INTO financial_payment_methods (id, tenant_id, name, code, method_type, notes, active, version, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := exec.ExecContext(ctx, query, item.ID, item.TenantID, item.Name, item.Code, item.MethodType, item.Notes, item.Active, item.Version, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return domain.PaymentMethod{}, mapError(err)
	}
	return item, nil
}
func (r MethodRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.PaymentMethod], error) {
	return searchCatalog[domain.PaymentMethod](ctx, r.db, "financial_payment_methods", tenantID, q)
}

func searchCatalog[T any](ctx context.Context, db *sql.DB, table, tenantID string, q ports.Query) (ports.Page[T], error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 20
	}
	args := []any{tenantID}
	where := "tenant_id = $1 AND deleted_at IS NULL"
	if q.Search != "" {
		args = append(args, q.Search)
		where += fmt.Sprintf(" AND (name ILIKE '%%' || $%d || '%%')", len(args))
	}
	var count int
	if err := db.QueryRowContext(ctx, fmt.Sprintf("SELECT count(*) FROM %s WHERE %s", table, where), args...).Scan(&count); err != nil {
		return ports.Page[T]{}, mapError(err)
	}
	limit := q.PageSize
	offset := (q.Page - 1) * q.PageSize
	query := fmt.Sprintf("SELECT coalesce(json_agg(row_to_json(t)), '[]'::json) FROM (SELECT * FROM %s WHERE %s ORDER BY name ASC LIMIT %d OFFSET %d) t", table, where, limit, offset)
	var body []byte
	if err := db.QueryRowContext(ctx, query, args...).Scan(&body); err != nil {
		return ports.Page[T]{}, mapError(err)
	}
	items := []T{}
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[T]{}, err
	}
	pages := (count + q.PageSize - 1) / q.PageSize
	return ports.Page[T]{Items: items, Page: q.Page, PageSize: q.PageSize, TotalItems: count, TotalPages: pages}, nil
}

// PeriodRepository implements ports.PeriodRepository
type PeriodRepository struct{ db *sql.DB }

func (r PeriodRepository) Create(ctx context.Context, item domain.FinancialPeriod) (domain.FinancialPeriod, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	query := `INSERT INTO financial_periods (id, tenant_id, year, month, start_date, end_date, status, version, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := exec.ExecContext(ctx, query, item.ID, item.TenantID, item.Year, item.Month, item.StartDate, item.EndDate, string(item.Status), item.Version, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return domain.FinancialPeriod{}, mapError(err)
	}
	return item, nil
}
func (r PeriodRepository) FindByID(ctx context.Context, tenantID, id string) (domain.FinancialPeriod, error) {
	query := `SELECT row_to_json(p) FROM financial_periods p WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL`
	return read[domain.FinancialPeriod](r.db.QueryRowContext(ctx, query, tenantID, id))
}
func (r PeriodRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FinancialPeriod], error) {
	return searchCatalog[domain.FinancialPeriod](ctx, r.db, "financial_periods", tenantID, q)
}
func (r PeriodRepository) Update(ctx context.Context, item domain.FinancialPeriod) (domain.FinancialPeriod, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	query := `UPDATE financial_periods SET status = $3, version = version + 1, updated_at = $4 WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL`
	res, err := exec.ExecContext(ctx, query, item.TenantID, item.ID, string(item.Status), time.Now().UTC())
	if err := changed(res, err); err != nil {
		return domain.FinancialPeriod{}, err
	}
	return item, nil
}
func (r PeriodRepository) FindClosedForDate(ctx context.Context, tenantID, dateStr string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM financial_periods WHERE tenant_id = $1 AND status = 'closed' AND $2::date BETWEEN start_date AND end_date AND deleted_at IS NULL)`
	var closed bool
	err := r.db.QueryRowContext(ctx, query, tenantID, dateStr).Scan(&closed)
	return closed, mapError(err)
}

// BudgetRepository implements ports.BudgetRepository
type BudgetRepository struct{ db *sql.DB }

func (r BudgetRepository) Create(ctx context.Context, item domain.Budget) (domain.Budget, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	query := `INSERT INTO financial_budgets (id, tenant_id, period_id, cost_center_id, category_id, name, status, planned_amount, actual_amount, notes, version, created_at, updated_at)
	VALUES ($1, $2, nullif($3,'')::uuid, nullif($4,'')::uuid, nullif($5,'')::uuid, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err := exec.ExecContext(ctx, query, item.ID, item.TenantID, item.PeriodID, item.CostCenterID, item.CategoryID, item.Name, item.Status, item.PlannedAmount, item.ActualAmount, item.Notes, item.Version, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return domain.Budget{}, mapError(err)
	}
	return item, nil
}
func (r BudgetRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Budget], error) {
	return searchCatalog[domain.Budget](ctx, r.db, "financial_budgets", tenantID, q)
}

// AdjustmentRepository implements ports.AdjustmentRepository
type AdjustmentRepository struct{ db *sql.DB }

func (r AdjustmentRepository) Create(ctx context.Context, item domain.FinancialAdjustment) (domain.FinancialAdjustment, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	query := `INSERT INTO financial_adjustments (id, tenant_id, transaction_id, adjustment_type, reason, adjusted_amount, notes, version, created_at, updated_at)
	VALUES ($1, $2, $3::uuid, $4, $5, $6, $7, 1, $8, $8)`
	_, err := exec.ExecContext(ctx, query, item.ID, item.TenantID, item.TransactionID, item.AdjustmentType, item.Reason, item.AdjustedAmount, item.Notes, item.CreatedAt)
	if err != nil {
		return domain.FinancialAdjustment{}, mapError(err)
	}
	return item, nil
}
func (r AdjustmentRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FinancialAdjustment], error) {
	return searchCatalog[domain.FinancialAdjustment](ctx, r.db, "financial_adjustments", tenantID, q)
}

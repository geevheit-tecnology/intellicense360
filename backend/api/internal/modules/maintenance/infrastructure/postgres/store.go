package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/ports"
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
			return application.ErrValidation
		}
		if pg.Code == "23503" || pg.Code == "23502" || pg.Code == "23514" || pg.Code == "22P02" {
			return application.ErrValidation
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

func (s *Store) WorkOrders() ports.WorkOrderRepository {
	return WorkOrderRepository{db: s.db}
}
func (s *Store) PreventivePlans() ports.PreventivePlanRepository {
	return PreventivePlanRepository{db: s.db}
}
func (s *Store) ServiceTypes() ports.ServiceTypeRepository {
	return ServiceTypeRepository{db: s.db}
}
func (s *Store) Categories() ports.CatalogRepository {
	return CatalogRepository{db: s.db, table: "maintenance_categories"}
}
func (s *Store) Types() ports.CatalogRepository {
	return CatalogRepository{db: s.db, table: "maintenance_types"}
}
func (s *Store) Priorities() ports.CatalogRepository {
	return CatalogRepository{db: s.db, table: "maintenance_priorities"}
}
func (s *Store) Reasons() ports.CatalogRepository {
	return CatalogRepository{db: s.db, table: "maintenance_reasons"}
}
func (s *Store) Workshops() ports.WorkshopRepository {
	return WorkshopRepository{db: s.db}
}
func (s *Store) Technicians() ports.TechnicianRepository {
	return TechnicianRepository{db: s.db}
}
func (s *Store) Labor() ports.LaborRepository {
	return LaborRepository{db: s.db}
}
func (s *Store) Downtime() ports.DowntimeRepository {
	return DowntimeRepository{db: s.db}
}
func (s *Store) History() ports.HistoryRepository {
	return HistoryRepository{db: s.db}
}

// ==================== WorkOrderRepository ====================

type WorkOrderRepository struct{ db *sql.DB }

func (r WorkOrderRepository) Create(ctx context.Context, item domain.WorkOrder) (domain.WorkOrder, error) {
	meta, err := json.Marshal(item.Metadata)
	if err != nil {
		meta = []byte("{}")
	}
	query := `INSERT INTO maintenance_work_orders AS t (
		id, tenant_id, code, asset_id, vehicle_id, title, description,
		kind, status, priority, category_id, workshop_id, technician_id,
		reason_id, service_type_id, opened_at, scheduled_at, started_at,
		completed_at, cancelled_at, estimated_hours, actual_hours,
		invoice_number, invoice_key, total_cost, payment_status,
		is_third_party, third_party_identifier,
		created_by, updated_by, metadata, created_at, updated_at, deleted_at, version
	) VALUES (
		$1, $2::uuid, $3, NULLIF($4, '')::uuid, NULLIF($5, '')::uuid, $6, $7,
		$8, $9, $10, NULLIF($11, '')::uuid, NULLIF($12, '')::uuid, NULLIF($13, '')::uuid,
		NULLIF($14, '')::uuid, NULLIF($15, '')::uuid, $16, $17, $18,
		$19, $20, $21, $22,
		$23, $24, $25, $26,
		$27, $28,
		NULLIF($29, '')::uuid, NULLIF($30, '')::uuid, $31::jsonb, $32, $33, $34, $35
	) RETURNING row_to_json(t)`

	if item.PaymentStatus == "" {
		item.PaymentStatus = "pending"
	}

	return read[domain.WorkOrder](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.Code, item.AssetID, item.VehicleID, item.Title, item.Description,
		item.Kind, item.Status, item.Priority, item.CategoryID, item.WorkshopID, item.TechnicianID,
		item.ReasonID, item.ServiceTypeID, item.OpenedAt, item.ScheduledAt, item.StartedAt,
		item.CompletedAt, item.CancelledAt, item.EstimatedHours, item.ActualHours,
		item.InvoiceNumber, item.InvoiceKey, item.TotalCost, item.PaymentStatus,
		item.IsThirdParty, item.ThirdPartyIdentifier,
		item.CreatedBy, item.UpdatedBy, meta, item.CreatedAt, item.UpdatedAt, item.DeletedAt, item.Version,
	))
}

func (r WorkOrderRepository) FindByID(ctx context.Context, tenantID string, id string) (domain.WorkOrder, error) {
	query := `SELECT row_to_json(t) FROM maintenance_work_orders AS t WHERE tenant_id = $1::uuid AND id = $2::uuid AND deleted_at IS NULL`
	return read[domain.WorkOrder](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, id))
}

func (r WorkOrderRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.WorkOrder], error) {
	where := []string{"tenant_id = $1::uuid", "deleted_at IS NULL"}
	args := []any{tenantID}

	if q.Search != "" {
		args = append(args, "%"+strings.ToLower(q.Search)+"%")
		where = append(where, fmt.Sprintf("(lower(title) LIKE $%d OR lower(code) LIKE $%d OR lower(invoice_number) LIKE $%d)", len(args), len(args), len(args)))
	}
	if v, ok := q.Filters["vehicle_id"]; ok && v != "" {
		args = append(args, v)
		where = append(where, fmt.Sprintf("vehicle_id = $%d::uuid", len(args)))
	}
	if v, ok := q.Filters["status"]; ok && v != "" {
		args = append(args, v)
		where = append(where, fmt.Sprintf("status = $%d", len(args)))
	}
	if v, ok := q.Filters["kind"]; ok && v != "" {
		args = append(args, v)
		where = append(where, fmt.Sprintf("kind = $%d", len(args)))
	}
	if v, ok := q.Filters["priority"]; ok && v != "" {
		args = append(args, v)
		where = append(where, fmt.Sprintf("priority = $%d", len(args)))
	}
	if v, ok := q.Filters["payment_status"]; ok && v != "" {
		args = append(args, v)
		where = append(where, fmt.Sprintf("payment_status = $%d", len(args)))
	}
	if v, ok := q.Filters["is_third_party"]; ok && v != "" {
		args = append(args, v == "true")
		where = append(where, fmt.Sprintf("is_third_party = $%d", len(args)))
	}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM maintenance_work_orders WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.WorkOrder]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	orderCol := "created_at"
	orderDir := "DESC"
	if q.SortBy == "code" {
		orderCol = "code"
	} else if q.SortBy == "opened_at" {
		orderCol = "opened_at"
	}
	if strings.ToLower(q.SortOrder) == "asc" {
		orderDir = "ASC"
	}

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM maintenance_work_orders
			WHERE %s
			ORDER BY %s %s
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, orderCol, orderDir, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.WorkOrder]{}, mapError(err)
	}

	var items []domain.WorkOrder
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.WorkOrder]{}, err
	}

	return ports.Page[domain.WorkOrder]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

func (r WorkOrderRepository) Update(ctx context.Context, item domain.WorkOrder) (domain.WorkOrder, error) {
	meta, err := json.Marshal(item.Metadata)
	if err != nil {
		meta = []byte("{}")
	}
	query := `UPDATE maintenance_work_orders AS t SET
		code = $3,
		asset_id = NULLIF($4, '')::uuid,
		vehicle_id = NULLIF($5, '')::uuid,
		title = $6,
		description = $7,
		kind = $8,
		status = $9,
		priority = $10,
		category_id = NULLIF($11, '')::uuid,
		workshop_id = NULLIF($12, '')::uuid,
		technician_id = NULLIF($13, '')::uuid,
		reason_id = NULLIF($14, '')::uuid,
		service_type_id = NULLIF($15, '')::uuid,
		scheduled_at = $16,
		started_at = $17,
		completed_at = $18,
		cancelled_at = $19,
		estimated_hours = $20,
		actual_hours = $21,
		invoice_number = $22,
		invoice_key = $23,
		total_cost = $24,
		payment_status = $25,
		is_third_party = $26,
		third_party_identifier = $27,
		updated_by = NULLIF($28, '')::uuid,
		metadata = $29::jsonb,
		updated_at = $30,
		version = t.version + 1
	WHERE t.tenant_id = $1::uuid AND t.id = $2::uuid AND t.deleted_at IS NULL AND t.version = $31
	RETURNING row_to_json(t)`

	if item.PaymentStatus == "" {
		item.PaymentStatus = "pending"
	}

	return read[domain.WorkOrder](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.TenantID, item.ID, item.Code, item.AssetID, item.VehicleID, item.Title, item.Description,
		item.Kind, item.Status, item.Priority, item.CategoryID, item.WorkshopID, item.TechnicianID,
		item.ReasonID, item.ServiceTypeID, item.ScheduledAt, item.StartedAt,
		item.CompletedAt, item.CancelledAt, item.EstimatedHours, item.ActualHours,
		item.InvoiceNumber, item.InvoiceKey, item.TotalCost, item.PaymentStatus,
		item.IsThirdParty, item.ThirdPartyIdentifier,
		item.UpdatedBy, meta, item.UpdatedAt, item.Version,
	))
}

func (r WorkOrderRepository) Delete(ctx context.Context, tenantID string, id string) error {
	query := `UPDATE maintenance_work_orders SET deleted_at = now() WHERE tenant_id = $1::uuid AND id = $2::uuid AND deleted_at IS NULL`
	res, err := platform.ExecutorFor(ctx, r.db).ExecContext(ctx, query, tenantID, id)
	return changed(res, err)
}

func (r WorkOrderRepository) Exists(ctx context.Context, tenantID string, id string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM maintenance_work_orders WHERE tenant_id = $1::uuid AND id = $2::uuid AND deleted_at IS NULL)`
	var exists bool
	err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, id).Scan(&exists)
	return exists, mapError(err)
}

func (r WorkOrderRepository) ExistsCode(ctx context.Context, tenantID string, code string, excludeID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM maintenance_work_orders WHERE tenant_id = $1::uuid AND code = $2 AND deleted_at IS NULL AND ($3 = '' OR id <> $3::uuid))`
	var exists bool
	err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, code, excludeID).Scan(&exists)
	return exists, mapError(err)
}

// ==================== PreventivePlanRepository ====================

type PreventivePlanRepository struct{ db *sql.DB }

func (r PreventivePlanRepository) Create(ctx context.Context, item domain.PreventivePlan) (domain.PreventivePlan, error) {
	query := `INSERT INTO maintenance_preventive_plans AS t (
		id, tenant_id, name, description, asset_id, vehicle_id, service_type_id,
		frequency, interval_value, next_due_at, next_due_value, active,
		created_at, updated_at, deleted_at, version
	) VALUES (
		$1, $2::uuid, $3, $4, NULLIF($5, '')::uuid, NULLIF($6, '')::uuid, NULLIF($7, '')::uuid,
		$8, $9, $10, $11, $12,
		$13, $14, $15, $16
	) RETURNING row_to_json(t)`
	return read[domain.PreventivePlan](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.Name, item.Description, item.AssetID, item.VehicleID, item.ServiceTypeID,
		item.Frequency, item.IntervalValue, item.NextDueAt, item.NextDueValue, item.Active,
		item.CreatedAt, item.UpdatedAt, item.DeletedAt, item.Version,
	))
}

func (r PreventivePlanRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.PreventivePlan], error) {
	where := []string{"tenant_id = $1::uuid", "deleted_at IS NULL"}
	args := []any{tenantID}

	if q.Search != "" {
		args = append(args, "%"+strings.ToLower(q.Search)+"%")
		where = append(where, fmt.Sprintf("(lower(name) LIKE $%d OR lower(description) LIKE $%d)", len(args), len(args)))
	}
	if v, ok := q.Filters["vehicle_id"]; ok && v != "" {
		args = append(args, v)
		where = append(where, fmt.Sprintf("vehicle_id = $%d::uuid", len(args)))
	}
	if v, ok := q.Filters["active"]; ok && v != "" {
		args = append(args, v == "true")
		where = append(where, fmt.Sprintf("active = $%d", len(args)))
	}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM maintenance_preventive_plans WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.PreventivePlan]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM maintenance_preventive_plans
			WHERE %s
			ORDER BY created_at DESC
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.PreventivePlan]{}, mapError(err)
	}

	var items []domain.PreventivePlan
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.PreventivePlan]{}, err
	}

	return ports.Page[domain.PreventivePlan]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

func (r PreventivePlanRepository) Update(ctx context.Context, item domain.PreventivePlan) (domain.PreventivePlan, error) {
	query := `UPDATE maintenance_preventive_plans AS t SET
		name = $3,
		description = $4,
		asset_id = NULLIF($5, '')::uuid,
		vehicle_id = NULLIF($6, '')::uuid,
		service_type_id = NULLIF($7, '')::uuid,
		frequency = $8,
		interval_value = $9,
		next_due_at = $10,
		next_due_value = $11,
		active = $12,
		updated_at = $13,
		version = t.version + 1
	WHERE t.tenant_id = $1::uuid AND t.id = $2::uuid AND t.deleted_at IS NULL AND t.version = $14
	RETURNING row_to_json(t)`
	return read[domain.PreventivePlan](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.TenantID, item.ID, item.Name, item.Description, item.AssetID, item.VehicleID, item.ServiceTypeID,
		item.Frequency, item.IntervalValue, item.NextDueAt, item.NextDueValue, item.Active,
		item.UpdatedAt, item.Version,
	))
}

func (r PreventivePlanRepository) Delete(ctx context.Context, tenantID string, id string) error {
	query := `UPDATE maintenance_preventive_plans SET deleted_at = now() WHERE tenant_id = $1::uuid AND id = $2::uuid AND deleted_at IS NULL`
	res, err := platform.ExecutorFor(ctx, r.db).ExecContext(ctx, query, tenantID, id)
	return changed(res, err)
}

// ==================== ServiceTypeRepository ====================

type ServiceTypeRepository struct{ db *sql.DB }

func (r ServiceTypeRepository) Create(ctx context.Context, item domain.ServiceType) (domain.ServiceType, error) {
	query := `INSERT INTO maintenance_service_types AS t (
		id, tenant_id, name, code, description, created_at, updated_at, deleted_at, version
	) VALUES (
		$1, $2::uuid, $3, $4, $5, $6, $7, $8, $9
	) RETURNING row_to_json(t)`
	return read[domain.ServiceType](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.Name, item.Code, item.Description,
		item.CreatedAt, item.UpdatedAt, item.DeletedAt, item.Version,
	))
}

func (r ServiceTypeRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.ServiceType], error) {
	where := []string{"tenant_id = $1::uuid", "deleted_at IS NULL"}
	args := []any{tenantID}

	if q.Search != "" {
		args = append(args, "%"+strings.ToLower(q.Search)+"%")
		where = append(where, fmt.Sprintf("(lower(name) LIKE $%d OR lower(code) LIKE $%d)", len(args), len(args)))
	}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM maintenance_service_types WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.ServiceType]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM maintenance_service_types
			WHERE %s
			ORDER BY name ASC
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.ServiceType]{}, mapError(err)
	}

	var items []domain.ServiceType
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.ServiceType]{}, err
	}

	return ports.Page[domain.ServiceType]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

func (r ServiceTypeRepository) Update(ctx context.Context, item domain.ServiceType) (domain.ServiceType, error) {
	query := `UPDATE maintenance_service_types AS t SET
		name = $3,
		code = $4,
		description = $5,
		updated_at = $6,
		version = t.version + 1
	WHERE t.tenant_id = $1::uuid AND t.id = $2::uuid AND t.deleted_at IS NULL AND t.version = $7
	RETURNING row_to_json(t)`
	return read[domain.ServiceType](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.TenantID, item.ID, item.Name, item.Code, item.Description, item.UpdatedAt, item.Version,
	))
}

func (r ServiceTypeRepository) Delete(ctx context.Context, tenantID string, id string) error {
	query := `UPDATE maintenance_service_types SET deleted_at = now() WHERE tenant_id = $1::uuid AND id = $2::uuid AND deleted_at IS NULL`
	res, err := platform.ExecutorFor(ctx, r.db).ExecContext(ctx, query, tenantID, id)
	return changed(res, err)
}

// ==================== CatalogRepository ====================

type CatalogRepository struct {
	db    *sql.DB
	table string
}

func (r CatalogRepository) Create(ctx context.Context, item domain.MaintenanceCatalog) (domain.MaintenanceCatalog, error) {
	query := fmt.Sprintf(`INSERT INTO %s AS t (
		id, tenant_id, name, code, description, created_at, updated_at, deleted_at, version
	) VALUES (
		$1, $2::uuid, $3, $4, $5, $6, $7, $8, $9
	) RETURNING row_to_json(t)`, r.table)
	return read[domain.MaintenanceCatalog](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.Name, item.Code, item.Description,
		item.CreatedAt, item.UpdatedAt, item.DeletedAt, item.Version,
	))
}

func (r CatalogRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.MaintenanceCatalog], error) {
	where := []string{"tenant_id = $1::uuid", "deleted_at IS NULL"}
	args := []any{tenantID}

	if q.Search != "" {
		args = append(args, "%"+strings.ToLower(q.Search)+"%")
		where = append(where, fmt.Sprintf("(lower(name) LIKE $%d OR lower(code) LIKE $%d)", len(args), len(args)))
	}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM %s WHERE %s", r.table, whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.MaintenanceCatalog]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM %s
			WHERE %s
			ORDER BY name ASC
			LIMIT $%d OFFSET $%d
		) AS t`, r.table, whereClause, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.MaintenanceCatalog]{}, mapError(err)
	}

	var items []domain.MaintenanceCatalog
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.MaintenanceCatalog]{}, err
	}

	return ports.Page[domain.MaintenanceCatalog]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

func (r CatalogRepository) Update(ctx context.Context, item domain.MaintenanceCatalog) (domain.MaintenanceCatalog, error) {
	query := fmt.Sprintf(`UPDATE %s AS t SET
		name = $3,
		code = $4,
		description = $5,
		updated_at = $6,
		version = t.version + 1
	WHERE t.tenant_id = $1::uuid AND t.id = $2::uuid AND t.deleted_at IS NULL AND t.version = $7
	RETURNING row_to_json(t)`, r.table)
	return read[domain.MaintenanceCatalog](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.TenantID, item.ID, item.Name, item.Code, item.Description, item.UpdatedAt, item.Version,
	))
}

func (r CatalogRepository) Delete(ctx context.Context, tenantID string, id string) error {
	query := fmt.Sprintf(`UPDATE %s SET deleted_at = now() WHERE tenant_id = $1::uuid AND id = $2::uuid AND deleted_at IS NULL`, r.table)
	res, err := platform.ExecutorFor(ctx, r.db).ExecContext(ctx, query, tenantID, id)
	return changed(res, err)
}

// ==================== WorkshopRepository ====================

type WorkshopRepository struct{ db *sql.DB }

func (r WorkshopRepository) Create(ctx context.Context, item domain.Workshop) (domain.Workshop, error) {
	query := `INSERT INTO maintenance_workshops AS t (
		id, tenant_id, name, document, phone, email, created_at, updated_at, deleted_at, version
	) VALUES (
		$1, $2::uuid, $3, $4, $5, $6, $7, $8, $9, $10
	) RETURNING row_to_json(t)`
	return read[domain.Workshop](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.Name, item.Document, item.Phone, item.Email,
		item.CreatedAt, item.UpdatedAt, item.DeletedAt, item.Version,
	))
}

func (r WorkshopRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Workshop], error) {
	where := []string{"tenant_id = $1::uuid", "deleted_at IS NULL"}
	args := []any{tenantID}

	if q.Search != "" {
		args = append(args, "%"+strings.ToLower(q.Search)+"%")
		where = append(where, fmt.Sprintf("(lower(name) LIKE $%d OR lower(document) LIKE $%d)", len(args), len(args)))
	}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM maintenance_workshops WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.Workshop]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM maintenance_workshops
			WHERE %s
			ORDER BY name ASC
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.Workshop]{}, mapError(err)
	}

	var items []domain.Workshop
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.Workshop]{}, err
	}

	return ports.Page[domain.Workshop]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

func (r WorkshopRepository) Update(ctx context.Context, item domain.Workshop) (domain.Workshop, error) {
	query := `UPDATE maintenance_workshops AS t SET
		name = $3,
		document = $4,
		phone = $5,
		email = $6,
		updated_at = $7,
		version = t.version + 1
	WHERE t.tenant_id = $1::uuid AND t.id = $2::uuid AND t.deleted_at IS NULL AND t.version = $8
	RETURNING row_to_json(t)`
	return read[domain.Workshop](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.TenantID, item.ID, item.Name, item.Document, item.Phone, item.Email,
		item.UpdatedAt, item.Version,
	))
}

func (r WorkshopRepository) Delete(ctx context.Context, tenantID string, id string) error {
	query := `UPDATE maintenance_workshops SET deleted_at = now() WHERE tenant_id = $1::uuid AND id = $2::uuid AND deleted_at IS NULL`
	res, err := platform.ExecutorFor(ctx, r.db).ExecContext(ctx, query, tenantID, id)
	return changed(res, err)
}

// ==================== TechnicianRepository ====================

type TechnicianRepository struct{ db *sql.DB }

func (r TechnicianRepository) Create(ctx context.Context, item domain.Technician) (domain.Technician, error) {
	query := `INSERT INTO maintenance_technicians AS t (
		id, tenant_id, name, document, phone, email, active, created_at, updated_at, deleted_at, version
	) VALUES (
		$1, $2::uuid, $3, $4, $5, $6, $7, $8, $9, $10, $11
	) RETURNING row_to_json(t)`
	return read[domain.Technician](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.Name, item.Document, item.Phone, item.Email, item.Active,
		item.CreatedAt, item.UpdatedAt, item.DeletedAt, item.Version,
	))
}

func (r TechnicianRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Technician], error) {
	where := []string{"tenant_id = $1::uuid", "deleted_at IS NULL"}
	args := []any{tenantID}

	if q.Search != "" {
		args = append(args, "%"+strings.ToLower(q.Search)+"%")
		where = append(where, fmt.Sprintf("(lower(name) LIKE $%d OR lower(document) LIKE $%d)", len(args), len(args)))
	}
	if v, ok := q.Filters["active"]; ok && v != "" {
		args = append(args, v == "true")
		where = append(where, fmt.Sprintf("active = $%d", len(args)))
	}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM maintenance_technicians WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.Technician]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM maintenance_technicians
			WHERE %s
			ORDER BY name ASC
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.Technician]{}, mapError(err)
	}

	var items []domain.Technician
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.Technician]{}, err
	}

	return ports.Page[domain.Technician]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

func (r TechnicianRepository) Update(ctx context.Context, item domain.Technician) (domain.Technician, error) {
	query := `UPDATE maintenance_technicians AS t SET
		name = $3,
		document = $4,
		phone = $5,
		email = $6,
		active = $7,
		updated_at = $8,
		version = t.version + 1
	WHERE t.tenant_id = $1::uuid AND t.id = $2::uuid AND t.deleted_at IS NULL AND t.version = $9
	RETURNING row_to_json(t)`
	return read[domain.Technician](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.TenantID, item.ID, item.Name, item.Document, item.Phone, item.Email, item.Active,
		item.UpdatedAt, item.Version,
	))
}

func (r TechnicianRepository) Delete(ctx context.Context, tenantID string, id string) error {
	query := `UPDATE maintenance_technicians SET deleted_at = now() WHERE tenant_id = $1::uuid AND id = $2::uuid AND deleted_at IS NULL`
	res, err := platform.ExecutorFor(ctx, r.db).ExecContext(ctx, query, tenantID, id)
	return changed(res, err)
}

// ==================== LaborRepository ====================

type LaborRepository struct{ db *sql.DB }

func (r LaborRepository) Create(ctx context.Context, item domain.LaborEntry) (domain.LaborEntry, error) {
	query := `INSERT INTO maintenance_labor_entries AS t (
		id, tenant_id, work_order_id, technician, hours, cost,
		started_at, finished_at, created_at, updated_at, deleted_at
	) VALUES (
		$1, $2::uuid, $3::uuid, $4, $5, $6,
		$7, $8, $9, $10, $11
	) RETURNING row_to_json(t)`
	return read[domain.LaborEntry](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.WorkOrderID, item.Technician, item.Hours, item.Cost,
		item.StartedAt, item.FinishedAt, item.CreatedAt, item.UpdatedAt, item.DeletedAt,
	))
}

func (r LaborRepository) List(ctx context.Context, tenantID string, workOrderID string, q ports.Query) (ports.Page[domain.LaborEntry], error) {
	where := []string{"tenant_id = $1::uuid", "work_order_id = $2::uuid", "deleted_at IS NULL"}
	args := []any{tenantID, workOrderID}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM maintenance_labor_entries WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.LaborEntry]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM maintenance_labor_entries
			WHERE %s
			ORDER BY created_at ASC
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.LaborEntry]{}, mapError(err)
	}

	var items []domain.LaborEntry
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.LaborEntry]{}, err
	}

	return ports.Page[domain.LaborEntry]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

func (r LaborRepository) Delete(ctx context.Context, tenantID string, workOrderID string, id string) error {
	query := `UPDATE maintenance_labor_entries SET deleted_at = now() WHERE tenant_id = $1::uuid AND work_order_id = $2::uuid AND id = $3::uuid AND deleted_at IS NULL`
	res, err := platform.ExecutorFor(ctx, r.db).ExecContext(ctx, query, tenantID, workOrderID, id)
	return changed(res, err)
}

// ==================== DowntimeRepository ====================

type DowntimeRepository struct{ db *sql.DB }

func (r DowntimeRepository) Create(ctx context.Context, item domain.Downtime) (domain.Downtime, error) {
	query := `INSERT INTO maintenance_downtime AS t (
		id, tenant_id, work_order_id, asset_id, vehicle_id, started_at, ended_at, reason,
		created_at, updated_at, deleted_at
	) VALUES (
		$1, $2::uuid, $3::uuid, NULLIF($4, '')::uuid, NULLIF($5, '')::uuid, $6, $7, $8,
		$9, $10, $11
	) RETURNING row_to_json(t)`
	return read[domain.Downtime](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.WorkOrderID, item.AssetID, item.VehicleID, item.StartedAt, item.EndedAt, item.Reason,
		item.CreatedAt, item.UpdatedAt, item.DeletedAt,
	))
}

func (r DowntimeRepository) List(ctx context.Context, tenantID string, workOrderID string, q ports.Query) (ports.Page[domain.Downtime], error) {
	where := []string{"tenant_id = $1::uuid", "work_order_id = $2::uuid", "deleted_at IS NULL"}
	args := []any{tenantID, workOrderID}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM maintenance_downtime WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.Downtime]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM maintenance_downtime
			WHERE %s
			ORDER BY started_at DESC
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.Downtime]{}, mapError(err)
	}

	var items []domain.Downtime
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.Downtime]{}, err
	}

	return ports.Page[domain.Downtime]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

func (r DowntimeRepository) End(ctx context.Context, tenantID string, workOrderID string, id string) (domain.Downtime, error) {
	now := time.Now().UTC()
	query := `UPDATE maintenance_downtime AS t SET
		ended_at = $4,
		updated_at = $4
	WHERE t.tenant_id = $1::uuid AND t.work_order_id = $2::uuid AND t.id = $3::uuid AND t.deleted_at IS NULL AND t.ended_at IS NULL
	RETURNING row_to_json(t)`
	return read[domain.Downtime](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, workOrderID, id, now))
}

// ==================== HistoryRepository ====================

type HistoryRepository struct{ db *sql.DB }

func (r HistoryRepository) Create(ctx context.Context, item domain.MaintenanceHistory) (domain.MaintenanceHistory, error) {
	query := `INSERT INTO maintenance_history AS t (
		id, tenant_id, work_order_id, event, actor_id, notes, created_at
	) VALUES (
		$1, $2::uuid, $3::uuid, $4, NULLIF($5, '')::uuid, $6, $7
	) RETURNING row_to_json(t)`
	return read[domain.MaintenanceHistory](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.WorkOrderID, item.Event, item.ActorID, item.Notes, item.CreatedAt,
	))
}

func (r HistoryRepository) List(ctx context.Context, tenantID string, workOrderID string, q ports.Query) (ports.Page[domain.MaintenanceHistory], error) {
	where := []string{"tenant_id = $1::uuid", "work_order_id = $2::uuid"}
	args := []any{tenantID, workOrderID}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM maintenance_history WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.MaintenanceHistory]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM maintenance_history
			WHERE %s
			ORDER BY created_at DESC
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.MaintenanceHistory]{}, mapError(err)
	}

	var items []domain.MaintenanceHistory
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.MaintenanceHistory]{}, err
	}

	return ports.Page[domain.MaintenanceHistory]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}


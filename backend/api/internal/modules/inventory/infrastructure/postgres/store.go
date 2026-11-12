package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/ports"
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

func (s *Store) Parts() ports.PartRepository {
	return PartRepository{db: s.db}
}
func (s *Store) Categories() ports.CatalogRepository {
	return CatalogRepository{db: s.db, table: "inventory_categories"}
}
func (s *Store) Brands() ports.CatalogRepository {
	return CatalogRepository{db: s.db, table: "inventory_brands"}
}
func (s *Store) Models() ports.CatalogRepository {
	return CatalogRepository{db: s.db, table: "inventory_models"}
}
func (s *Store) Units() ports.CatalogRepository {
	return CatalogRepository{db: s.db, table: "inventory_units"}
}
func (s *Store) Warehouses() ports.WarehouseRepository {
	return WarehouseRepository{db: s.db}
}
func (s *Store) Locations() ports.LocationRepository {
	return LocationRepository{db: s.db}
}

// PartRepository implements ports.PartRepository
type PartRepository struct{ db *sql.DB }

func (r PartRepository) Create(ctx context.Context, item domain.Part) (domain.Part, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	metaBytes, _ := json.Marshal(item.Metadata)
	query := `INSERT INTO inventory_parts (
		id, tenant_id, sku, internal_code, name, description,
		category_id, brand_id, model_id, unit_id, status,
		unit_cost, min_stock, current_stock, location, metadata,
		version, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6,
		nullif($7,'')::uuid, nullif($8,'')::uuid, nullif($9,'')::uuid, $10, $11,
		$12, $13, $14, $15, $16,
		$17, $18, $19
	)`
	_, err := exec.ExecContext(ctx, query,
		item.ID, item.TenantID, item.SKU, item.InternalCode, item.Name, item.Description,
		item.CategoryID, item.BrandID, item.ModelID, item.UnitID, string(item.Status),
		item.UnitCost, item.MinStock, item.CurrentStock, item.Location, metaBytes,
		item.Version, item.CreatedAt, item.UpdatedAt,
	)
	if err != nil {
		return domain.Part{}, mapError(err)
	}
	return item, nil
}

func (r PartRepository) FindByID(ctx context.Context, tenantID, id string) (domain.Part, error) {
	query := `SELECT row_to_json(p) FROM inventory_parts p WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL`
	return read[domain.Part](r.db.QueryRowContext(ctx, query, tenantID, id))
}

func (r PartRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Part], error) {
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
		where += fmt.Sprintf(" AND (name ILIKE '%%' || $%d || '%%' OR sku ILIKE '%%' || $%d || '%%' OR internal_code ILIKE '%%' || $%d || '%%')", len(args), len(args), len(args))
	}
	for k, v := range q.Filters {
		if v == "" {
			continue
		}
		if k == "category_id" {
			args = append(args, v)
			where += fmt.Sprintf(" AND category_id = $%d::uuid", len(args))
		} else if k == "status" {
			args = append(args, v)
			where += fmt.Sprintf(" AND status = $%d", len(args))
		}
	}

	var count int
	if err := r.db.QueryRowContext(ctx, "SELECT count(*) FROM inventory_parts WHERE "+where, args...).Scan(&count); err != nil {
		return ports.Page[domain.Part]{}, mapError(err)
	}

	limit := q.PageSize
	offset := (q.Page - 1) * q.PageSize
	order := "name ASC, created_at DESC"
	query := fmt.Sprintf("SELECT coalesce(json_agg(row_to_json(p)), '[]'::json) FROM (SELECT * FROM inventory_parts WHERE %s ORDER BY %s LIMIT %d OFFSET %d) p", where, order, limit, offset)

	var body []byte
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&body); err != nil {
		return ports.Page[domain.Part]{}, mapError(err)
	}
	items := []domain.Part{}
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.Part]{}, err
	}
	pages := (count + q.PageSize - 1) / q.PageSize
	return ports.Page[domain.Part]{Items: items, Page: q.Page, PageSize: q.PageSize, TotalItems: count, TotalPages: pages}, nil
}

func (r PartRepository) Update(ctx context.Context, item domain.Part) (domain.Part, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	metaBytes, _ := json.Marshal(item.Metadata)
	query := `UPDATE inventory_parts SET
		sku = $3, internal_code = $4, name = $5, description = $6,
		category_id = nullif($7,'')::uuid, brand_id = nullif($8,'')::uuid,
		model_id = nullif($9,'')::uuid, unit_id = $10, status = $11,
		unit_cost = $12, min_stock = $13, current_stock = $14, location = $15,
		metadata = $16, version = version + 1, updated_at = $17
	WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL`
	res, err := exec.ExecContext(ctx, query,
		item.TenantID, item.ID, item.SKU, item.InternalCode, item.Name, item.Description,
		item.CategoryID, item.BrandID, item.ModelID, item.UnitID, string(item.Status),
		item.UnitCost, item.MinStock, item.CurrentStock, item.Location, metaBytes, time.Now().UTC(),
	)
	if err := changed(res, err); err != nil {
		return domain.Part{}, err
	}
	return item, nil
}

func (r PartRepository) Delete(ctx context.Context, tenantID, id string) error {
	exec := platform.ExecutorFor(ctx, r.db)
	res, err := exec.ExecContext(ctx, `UPDATE inventory_parts SET deleted_at = now(), version = version + 1 WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL`, tenantID, id)
	return changed(res, err)
}

func (r PartRepository) ExistsSKU(ctx context.Context, tenantID, sku, excludeID string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM inventory_parts WHERE tenant_id = $1 AND lower(sku) = lower($2) AND id != $3 AND deleted_at IS NULL)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, tenantID, sku, excludeID).Scan(&exists)
	return exists, mapError(err)
}

func (r PartRepository) ExistsInternalCode(ctx context.Context, tenantID, code, excludeID string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM inventory_parts WHERE tenant_id = $1 AND lower(internal_code) = lower($2) AND id != $3 AND deleted_at IS NULL)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, tenantID, code, excludeID).Scan(&exists)
	return exists, mapError(err)
}

// CatalogRepository implements ports.CatalogRepository for categories, units, brands, models
type CatalogRepository struct {
	db    *sql.DB
	table string
}

func (r CatalogRepository) Create(ctx context.Context, item domain.Catalog) (domain.Catalog, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	query := fmt.Sprintf(`INSERT INTO %s (id, tenant_id, name, code, description, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, r.table)
	_, err := exec.ExecContext(ctx, query, item.ID, item.TenantID, item.Name, item.Code, item.Description, item.Version, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return domain.Catalog{}, mapError(err)
	}
	return item, nil
}

func (r CatalogRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Catalog], error) {
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
		where += fmt.Sprintf(" AND (name ILIKE '%%' || $%d || '%%' OR code ILIKE '%%' || $%d || '%%')", len(args), len(args))
	}
	var count int
	if err := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT count(*) FROM %s WHERE %s", r.table, where), args...).Scan(&count); err != nil {
		return ports.Page[domain.Catalog]{}, mapError(err)
	}
	limit := q.PageSize
	offset := (q.Page - 1) * q.PageSize
	query := fmt.Sprintf("SELECT coalesce(json_agg(row_to_json(c)), '[]'::json) FROM (SELECT * FROM %s WHERE %s ORDER BY name ASC LIMIT %d OFFSET %d) c", r.table, where, limit, offset)
	var body []byte
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&body); err != nil {
		return ports.Page[domain.Catalog]{}, mapError(err)
	}
	items := []domain.Catalog{}
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.Catalog]{}, err
	}
	pages := (count + q.PageSize - 1) / q.PageSize
	return ports.Page[domain.Catalog]{Items: items, Page: q.Page, PageSize: q.PageSize, TotalItems: count, TotalPages: pages}, nil
}

func (r CatalogRepository) Update(ctx context.Context, item domain.Catalog) (domain.Catalog, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	query := fmt.Sprintf(`UPDATE %s SET name = $3, code = $4, description = $5, version = version + 1, updated_at = $6 WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL`, r.table)
	res, err := exec.ExecContext(ctx, query, item.TenantID, item.ID, item.Name, item.Code, item.Description, time.Now().UTC())
	if err := changed(res, err); err != nil {
		return domain.Catalog{}, err
	}
	return item, nil
}

func (r CatalogRepository) Delete(ctx context.Context, tenantID, id string) error {
	exec := platform.ExecutorFor(ctx, r.db)
	query := fmt.Sprintf(`UPDATE %s SET deleted_at = now(), version = version + 1 WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL`, r.table)
	res, err := exec.ExecContext(ctx, query, tenantID, id)
	return changed(res, err)
}

func (r CatalogRepository) ExistsCode(ctx context.Context, tenantID, code, excludeID string) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM %s WHERE tenant_id = $1 AND lower(code) = lower($2) AND id != $3 AND deleted_at IS NULL)`, r.table)
	var exists bool
	err := r.db.QueryRowContext(ctx, query, tenantID, code, excludeID).Scan(&exists)
	return exists, mapError(err)
}

// WarehouseRepository implements ports.WarehouseRepository
type WarehouseRepository struct{ db *sql.DB }

func (r WarehouseRepository) Create(ctx context.Context, item domain.Warehouse) (domain.Warehouse, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	query := `INSERT INTO inventory_warehouses (id, tenant_id, name, code, address, active, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := exec.ExecContext(ctx, query, item.ID, item.TenantID, item.Name, item.Code, item.Address, item.Active, item.Version, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return domain.Warehouse{}, mapError(err)
	}
	return item, nil
}

func (r WarehouseRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Warehouse], error) {
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
		where += fmt.Sprintf(" AND (name ILIKE '%%' || $%d || '%%' OR code ILIKE '%%' || $%d || '%%')", len(args), len(args))
	}
	var count int
	if err := r.db.QueryRowContext(ctx, "SELECT count(*) FROM inventory_warehouses WHERE "+where, args...).Scan(&count); err != nil {
		return ports.Page[domain.Warehouse]{}, mapError(err)
	}
	limit := q.PageSize
	offset := (q.Page - 1) * q.PageSize
	query := fmt.Sprintf("SELECT coalesce(json_agg(row_to_json(w)), '[]'::json) FROM (SELECT * FROM inventory_warehouses WHERE %s ORDER BY name ASC LIMIT %d OFFSET %d) w", where, limit, offset)
	var body []byte
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&body); err != nil {
		return ports.Page[domain.Warehouse]{}, mapError(err)
	}
	items := []domain.Warehouse{}
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.Warehouse]{}, err
	}
	pages := (count + q.PageSize - 1) / q.PageSize
	return ports.Page[domain.Warehouse]{Items: items, Page: q.Page, PageSize: q.PageSize, TotalItems: count, TotalPages: pages}, nil
}

func (r WarehouseRepository) Update(ctx context.Context, item domain.Warehouse) (domain.Warehouse, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	query := `UPDATE inventory_warehouses SET name = $3, code = $4, address = $5, active = $6, version = version + 1, updated_at = $7 WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL`
	res, err := exec.ExecContext(ctx, query, item.TenantID, item.ID, item.Name, item.Code, item.Address, item.Active, time.Now().UTC())
	if err := changed(res, err); err != nil {
		return domain.Warehouse{}, err
	}
	return item, nil
}

func (r WarehouseRepository) Delete(ctx context.Context, tenantID, id string) error {
	exec := platform.ExecutorFor(ctx, r.db)
	res, err := exec.ExecContext(ctx, `UPDATE inventory_warehouses SET deleted_at = now(), version = version + 1 WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL`, tenantID, id)
	return changed(res, err)
}

// LocationRepository implements ports.LocationRepository
type LocationRepository struct{ db *sql.DB }

func (r LocationRepository) Create(ctx context.Context, item domain.WarehouseLocation) (domain.WarehouseLocation, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	query := `INSERT INTO inventory_locations (id, tenant_id, warehouse_id, name, code, version, created_at, updated_at) VALUES ($1, $2, $3::uuid, $4, $5, $6, $7, $8)`
	_, err := exec.ExecContext(ctx, query, item.ID, item.TenantID, item.WarehouseID, item.Name, item.Code, item.Version, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return domain.WarehouseLocation{}, mapError(err)
	}
	return item, nil
}

func (r LocationRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.WarehouseLocation], error) {
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
		where += fmt.Sprintf(" AND (name ILIKE '%%' || $%d || '%%' OR code ILIKE '%%' || $%d || '%%')", len(args), len(args))
	}
	var count int
	if err := r.db.QueryRowContext(ctx, "SELECT count(*) FROM inventory_locations WHERE "+where, args...).Scan(&count); err != nil {
		return ports.Page[domain.WarehouseLocation]{}, mapError(err)
	}
	limit := q.PageSize
	offset := (q.Page - 1) * q.PageSize
	query := fmt.Sprintf("SELECT coalesce(json_agg(row_to_json(l)), '[]'::json) FROM (SELECT * FROM inventory_locations WHERE %s ORDER BY name ASC LIMIT %d OFFSET %d) l", where, limit, offset)
	var body []byte
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&body); err != nil {
		return ports.Page[domain.WarehouseLocation]{}, mapError(err)
	}
	items := []domain.WarehouseLocation{}
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.WarehouseLocation]{}, err
	}
	pages := (count + q.PageSize - 1) / q.PageSize
	return ports.Page[domain.WarehouseLocation]{Items: items, Page: q.Page, PageSize: q.PageSize, TotalItems: count, TotalPages: pages}, nil
}

func (r LocationRepository) Update(ctx context.Context, item domain.WarehouseLocation) (domain.WarehouseLocation, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	query := `UPDATE inventory_locations SET warehouse_id = $3::uuid, name = $4, code = $5, version = version + 1, updated_at = $6 WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL`
	res, err := exec.ExecContext(ctx, query, item.TenantID, item.ID, item.WarehouseID, item.Name, item.Code, time.Now().UTC())
	if err := changed(res, err); err != nil {
		return domain.WarehouseLocation{}, err
	}
	return item, nil
}

func (r LocationRepository) Delete(ctx context.Context, tenantID, id string) error {
	exec := platform.ExecutorFor(ctx, r.db)
	res, err := exec.ExecContext(ctx, `UPDATE inventory_locations SET deleted_at = now(), version = version + 1 WHERE tenant_id = $1 AND id = $2 AND deleted_at IS NULL`, tenantID, id)
	return changed(res, err)
}

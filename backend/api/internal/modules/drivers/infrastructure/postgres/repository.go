package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/drivers/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/drivers/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/drivers/ports"
	platform "github.com/geevheit/intelligence360/backend/api/internal/platform/postgres"
	"github.com/lib/pq"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return Repository{db: db} }

var _ ports.DriverRepository = Repository{}

func mapError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return application.ErrNotFound
	}
	var pg *pq.Error
	if errors.As(err, &pg) {
		switch pg.Code {
		case "23505":
			return application.ErrConflict
		case "23503":
			return application.ErrReferenceNotFound
		case "23514", "23502", "22P02", "22007", "22008":
			return application.ErrValidation
		}
	}
	return err
}
func read(row *sql.Row) (domain.Driver, error) {
	var item domain.Driver
	var body []byte
	if err := row.Scan(&body); err != nil {
		return item, mapError(err)
	}
	err := json.Unmarshal(body, &item)
	return item, err
}
func (r Repository) FindByID(ctx context.Context, tenant, id string) (domain.Driver, error) {
	return read(platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `SELECT row_to_json(t) FROM drivers t WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenant, id))
}
func (r Repository) Exists(ctx context.Context, tenant, id string) (bool, error) {
	var exists bool
	err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM drivers WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL)`, tenant, id).Scan(&exists)
	return exists, mapError(err)
}
func (r Repository) Create(ctx context.Context, item domain.Driver) (domain.Driver, error) {
	return read(platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO drivers AS t(id,tenant_id,code,name,tax_id,phone,email,vehicle_id,status,notes,created_at,updated_at,version) VALUES($1,$2,$3,$4,NULLIF($5,''),$6,$7,NULLIF($8,'')::uuid,$9,$10,$11,$12,$13) RETURNING row_to_json(t)`, item.ID, item.TenantID, item.Code, item.Name, item.TaxID, item.Phone, item.Email, item.VehicleID, item.Status, item.Notes, item.CreatedAt, item.UpdatedAt, item.Version))
}
func (r Repository) Update(ctx context.Context, item domain.Driver) (domain.Driver, error) {
	saved, err := read(platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `UPDATE drivers AS t SET code=$3,name=$4,tax_id=NULLIF($5,''),phone=$6,email=$7,vehicle_id=NULLIF($8,'')::uuid,status=$9,notes=$10,updated_at=$11,version=version+1 WHERE tenant_id=$2::uuid AND id=$1::uuid AND version=$12 AND deleted_at IS NULL RETURNING row_to_json(t)`, item.ID, item.TenantID, item.Code, item.Name, item.TaxID, item.Phone, item.Email, item.VehicleID, item.Status, item.Notes, item.UpdatedAt, item.Version))
	if errors.Is(err, application.ErrNotFound) {
		if _, findErr := r.FindByID(ctx, item.TenantID, item.ID); findErr == nil {
			return saved, application.ErrConflict
		}
	}
	return saved, err
}
func (r Repository) Delete(ctx context.Context, tenant, id string, version int64) error {
	result, err := platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `UPDATE drivers SET deleted_at=now(),updated_at=now(),version=version+1 WHERE tenant_id=$1::uuid AND id=$2::uuid AND version=$3 AND deleted_at IS NULL`, tenant, id, version)
	if err != nil {
		return mapError(err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		if _, err := r.FindByID(ctx, tenant, id); err != nil {
			return err
		}
		return application.ErrConflict
	}
	return nil
}
func (r Repository) Search(ctx context.Context, tenant string, q ports.Query) (ports.Page, error) {
	result := ports.Page{Items: []domain.Driver{}, Page: q.Page, PageSize: q.PageSize}
	args := []any{tenant}
	where := "tenant_id=$1::uuid AND deleted_at IS NULL"
	add := func(column, value string, asUUID bool) {
		if value != "" {
			args = append(args, value)
			cast := ""
			if asUUID {
				cast = "::uuid"
			}
			where += fmt.Sprintf(" AND %s=$%d%s", column, len(args), cast)
		}
	}
	add("status", q.Status, false)
	add("vehicle_id", q.VehicleID, true)
	if q.Search != "" {
		args = append(args, q.Search)
		where += fmt.Sprintf(" AND (name ILIKE '%%' || $%d || '%%' OR code ILIKE '%%' || $%d || '%%')", len(args), len(args))
	}
	order := "created_at,id"

	args = append(args, q.PageSize, (int64(q.Page)-1)*int64(q.PageSize))
	query := fmt.Sprintf(`WITH filtered AS (SELECT * FROM drivers WHERE %s), page AS (SELECT * FROM filtered ORDER BY %s LIMIT $%d OFFSET $%d) SELECT (SELECT count(*) FROM filtered),coalesce((SELECT json_agg(p) FROM page p),'[]'::json)`, where, order, len(args)-1, len(args))
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, args...).Scan(&result.TotalItems, &body); err != nil {
		return result, mapError(err)
	}
	if err := json.Unmarshal(body, &result.Items); err != nil {
		return result, err
	}
	if result.TotalItems > 0 {
		result.TotalPages = (result.TotalItems + q.PageSize - 1) / q.PageSize
	}
	return result, nil
}

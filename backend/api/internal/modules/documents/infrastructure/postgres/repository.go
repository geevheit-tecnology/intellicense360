package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/documents/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/documents/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/documents/ports"
	platform "github.com/geevheit/intelligence360/backend/api/internal/platform/postgres"
	"github.com/lib/pq"
)

type Repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return Repository{db: db} }

var _ ports.DocumentRepository = Repository{}

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
func read(row *sql.Row) (domain.Document, error) {
	var item domain.Document
	var body []byte
	if err := row.Scan(&body); err != nil {
		return item, mapError(err)
	}
	err := json.Unmarshal(body, &item)
	return item, err
}
func (r Repository) FindByID(ctx context.Context, tenant, id string) (domain.Document, error) {
	return read(platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `SELECT row_to_json(t) FROM documents t WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenant, id))
}
func (r Repository) Exists(ctx context.Context, tenant, id string) (bool, error) {
	var exists bool
	err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM documents WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL)`, tenant, id).Scan(&exists)
	return exists, mapError(err)
}
func (r Repository) Create(ctx context.Context, item domain.Document) (domain.Document, error) {
	return read(platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO documents AS t(id,tenant_id,driver_id,vehicle_id,kind,title,number,issuer,issued_on,expires_on,status,notes,created_at,updated_at,version) VALUES($1,$2,NULLIF($3,'')::uuid,NULLIF($4,'')::uuid,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15) RETURNING row_to_json(t)`, item.ID, item.TenantID, item.DriverID, item.VehicleID, item.Kind, item.Title, item.Number, item.Issuer, item.IssuedOn, item.ExpiresOn, item.Status, item.Notes, item.CreatedAt, item.UpdatedAt, item.Version))
}
func (r Repository) Update(ctx context.Context, item domain.Document) (domain.Document, error) {
	saved, err := read(platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `UPDATE documents AS t SET driver_id=NULLIF($3,'')::uuid,vehicle_id=NULLIF($4,'')::uuid,kind=$5,title=$6,number=$7,issuer=$8,issued_on=$9,expires_on=$10,status=$11,notes=$12,updated_at=$13,version=version+1 WHERE tenant_id=$2::uuid AND id=$1::uuid AND version=$14 AND deleted_at IS NULL RETURNING row_to_json(t)`, item.ID, item.TenantID, item.DriverID, item.VehicleID, item.Kind, item.Title, item.Number, item.Issuer, item.IssuedOn, item.ExpiresOn, item.Status, item.Notes, item.UpdatedAt, item.Version))
	if errors.Is(err, application.ErrNotFound) {
		if _, findErr := r.FindByID(ctx, item.TenantID, item.ID); findErr == nil {
			return saved, application.ErrConflict
		}
	}
	return saved, err
}
func (r Repository) Delete(ctx context.Context, tenant, id string, version int64) error {
	result, err := platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `UPDATE documents SET deleted_at=now(),updated_at=now(),version=version+1 WHERE tenant_id=$1::uuid AND id=$2::uuid AND version=$3 AND deleted_at IS NULL`, tenant, id, version)
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
	result := ports.Page{Items: []domain.Document{}, Page: q.Page, PageSize: q.PageSize}
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
	add("kind", q.Kind, false)
	add("status", q.Status, false)
	add("driver_id", q.DriverID, true)
	add("vehicle_id", q.VehicleID, true)
	if q.Search != "" {
		args = append(args, q.Search)
		where += fmt.Sprintf(" AND (title ILIKE '%%' || $%d || '%%' OR number ILIKE '%%' || $%d || '%%')", len(args), len(args))
	}
	if q.ExpiryStatus != "" {
		args = append(args, q.AsOf.Format("2006-01-02"))
		dateArg := len(args)
		args = append(args, q.ExpiryStatus)
		where += fmt.Sprintf(` AND (CASE WHEN status='archived' THEN 'archived' WHEN expires_on IS NULL THEN 'no_expiry' WHEN expires_on<$%d::date THEN 'expired' WHEN expires_on<=$%d::date+30 THEN 'due_soon' ELSE 'valid' END)=$%d`, dateArg, dateArg, len(args))
	}
	order := "expires_on ASC NULLS LAST,id"

	args = append(args, q.PageSize, (int64(q.Page)-1)*int64(q.PageSize))
	query := fmt.Sprintf(`WITH filtered AS (SELECT * FROM documents WHERE %s), page AS (SELECT * FROM filtered ORDER BY %s LIMIT $%d OFFSET $%d) SELECT (SELECT count(*) FROM filtered),coalesce((SELECT json_agg(p) FROM page p),'[]'::json)`, where, order, len(args)-1, len(args))
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

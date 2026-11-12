package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/tires/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/tires/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/tires/ports"
	platform "github.com/geevheit/intelligence360/backend/api/internal/platform/postgres"
	"github.com/lib/pq"
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
			switch pg.Constraint {
			case "ux_tires_serial_number_active":
				return application.ErrSerialNumberTaken
			case "ux_tires_fire_number_active":
				return application.ErrFireNumberTaken
			}
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

func (s *Store) Tires() ports.TireRepository {
	return TireRepository{db: s.db}
}
func (s *Store) Inspections() ports.TireInspectionRepository {
	return InspectionRepository{db: s.db}
}
func (s *Store) Movements() ports.TireMovementRepository {
	return MovementRepository{db: s.db}
}
func (s *Store) History() ports.TireHistoryRepository {
	return HistoryRepository{db: s.db}
}

// TireRepository
type TireRepository struct{ db *sql.DB }

func (r TireRepository) Create(ctx context.Context, item domain.Tire) (domain.Tire, error) {
	query := `INSERT INTO tires AS t (
		id, tenant_id, serial_number, fire_number, brand, model, size, construction, tire_type,
		position_type, manufacturing_date, purchase_date, purchase_value, supplier, dot,
		current_tread_mm, original_tread_mm, minimum_tread_mm, status, vehicle_id, position,
		current_km, total_km, recap_count, notes, created_by, updated_by, created_at, updated_at,
		dimension, load_index, speed_rating, supplier_reference, warranty, condition, version
	) VALUES (
		$1, $2::uuid, $3, $4, $5, $6, $7, $8, $9,
		$10, $11, $12, $13, $14, $15,
		$16, $17, $18, $19, NULLIF($20, '')::uuid, $21,
		$22, $23, $24, $25, NULLIF($26, '')::uuid, NULLIF($27, '')::uuid, $28, $29,
		$30, $31, $32, $33, $34, $35, $36
	) RETURNING row_to_json(t)`
	return read[domain.Tire](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.SerialNumber, item.FireNumber, item.Brand, item.Model, item.Size, item.Construction, item.TireType,
		item.PositionType, item.ManufacturingDate, item.PurchaseDate, item.PurchaseValue, item.Supplier, item.DOT,
		item.CurrentTreadMM, item.OriginalTreadMM, item.MinimumTreadMM, item.Status, item.VehicleID, item.Position,
		item.CurrentKM, item.TotalKM, item.RecapCount, item.Notes, item.CreatedBy, item.UpdatedBy, item.CreatedAt, item.UpdatedAt,
		item.Dimension, item.LoadIndex, item.SpeedRating, item.SupplierReference, item.Warranty, item.Condition, item.Version,
	))
}

func (r TireRepository) FindByID(ctx context.Context, tenantID, id string) (domain.Tire, error) {
	query := `SELECT row_to_json(t) FROM tires t WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`
	return read[domain.Tire](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, id))
}

func (r TireRepository) Update(ctx context.Context, item domain.Tire) (domain.Tire, error) {
	query := `UPDATE tires AS t SET
		serial_number=$3, fire_number=$4, brand=$5, model=$6, size=$7, construction=$8, tire_type=$9,
		position_type=$10, manufacturing_date=$11, purchase_date=$12, purchase_value=$13, supplier=$14, dot=$15,
		current_tread_mm=$16, original_tread_mm=$17, minimum_tread_mm=$18, status=$19, vehicle_id=NULLIF($20, '')::uuid, position=$21,
		current_km=$22, total_km=$23, recap_count=$24, notes=$25, updated_by=NULLIF($26, '')::uuid, updated_at=$27,
		dimension=$28, load_index=$29, speed_rating=$30, supplier_reference=$31, warranty=$32, condition=$33, version=t.version+1
		WHERE tenant_id=$2::uuid AND id=$1::uuid AND deleted_at IS NULL AND version=$34-1
		RETURNING row_to_json(t)`
	saved, err := read[domain.Tire](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.SerialNumber, item.FireNumber, item.Brand, item.Model, item.Size, item.Construction, item.TireType,
		item.PositionType, item.ManufacturingDate, item.PurchaseDate, item.PurchaseValue, item.Supplier, item.DOT,
		item.CurrentTreadMM, item.OriginalTreadMM, item.MinimumTreadMM, item.Status, item.VehicleID, item.Position,
		item.CurrentKM, item.TotalKM, item.RecapCount, item.Notes, item.UpdatedBy, item.UpdatedAt,
		item.Dimension, item.LoadIndex, item.SpeedRating, item.SupplierReference, item.Warranty, item.Condition,
		item.Version,
	))
	if errors.Is(err, application.ErrNotFound) {
		if _, findErr := r.FindByID(ctx, item.TenantID, item.ID); findErr == nil {
			return saved, application.ErrValidation
		}
	}
	return saved, err
}

func (r TireRepository) Delete(ctx context.Context, tenantID, id string) error {
	query := `UPDATE tires SET deleted_at=now(), updated_at=now(), version=version+1 WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, query, tenantID, id))
}

func (r TireRepository) Exists(ctx context.Context, tenantID, id string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM tires WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL)`
	err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, id).Scan(&exists)
	return exists, mapError(err)
}

func (r TireRepository) ExistsSerialNumber(ctx context.Context, tenantID, serialNumber, exceptTireID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM tires WHERE tenant_id=$1::uuid AND upper(serial_number)=upper($2) AND deleted_at IS NULL AND id <> NULLIF($3, '')::uuid)`
	err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, serialNumber, exceptTireID).Scan(&exists)
	return exists, mapError(err)
}

func (r TireRepository) ExistsFireNumber(ctx context.Context, tenantID, fireNumber, exceptTireID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM tires WHERE tenant_id=$1::uuid AND upper(fire_number)=upper($2) AND deleted_at IS NULL AND id <> NULLIF($3, '')::uuid)`
	err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, fireNumber, exceptTireID).Scan(&exists)
	return exists, mapError(err)
}

func (r TireRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Tire], error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 20
	}
	result := ports.Page[domain.Tire]{Items: []domain.Tire{}, Page: q.Page, PageSize: q.PageSize}
	args := []any{tenantID}
	where := "tenant_id=$1 AND deleted_at IS NULL"
	if q.Search != "" {
		args = append(args, q.Search)
		where += fmt.Sprintf(" AND (serial_number || ' ' || fire_number || ' ' || brand || ' ' || model) ILIKE '%%' || $%d || '%%'", len(args))
	}
	for key, value := range q.Filters {
		if value == "" {
			continue
		}
		if key == "status" || key == "vehicle_id" || key == "brand" {
			args = append(args, value)
			if key == "vehicle_id" {
				where += fmt.Sprintf(" AND %s=$%d::uuid", key, len(args))
			} else {
				where += fmt.Sprintf(" AND %s=$%d", key, len(args))
			}
		}
	}

	err := (platform.Transactions{DB: r.db}).Run(ctx, func(ctx context.Context) error {
		exec := platform.ExecutorFor(ctx, r.db)
		if err := exec.QueryRowContext(ctx, "SELECT count(*) FROM tires WHERE "+where, args...).Scan(&result.TotalItems); err != nil {
			return mapError(err)
		}
		sortBy := "created_at"
		switch q.SortBy {
		case "created_at", "updated_at", "id", "serial_number", "fire_number":
			sortBy = q.SortBy
		}
		order := "DESC"
		if strings.EqualFold(q.SortOrder, "asc") {
			order = "ASC"
		}
		args = append(args, q.PageSize, (int64(q.Page)-1)*int64(q.PageSize))
		rows, err := exec.QueryContext(ctx, fmt.Sprintf("SELECT row_to_json(t) FROM tires t WHERE %s ORDER BY %s %s,id LIMIT $%d OFFSET $%d", where, sortBy, order, len(args)-1, len(args)), args...)
		if err != nil {
			return mapError(err)
		}
		defer rows.Close()
		for rows.Next() {
			var body []byte
			var item domain.Tire
			if err := rows.Scan(&body); err != nil {
				return err
			}
			if err := json.Unmarshal(body, &item); err != nil {
				return err
			}
			result.Items = append(result.Items, item)
		}
		return rows.Err()
	})
	if result.TotalItems > 0 {
		result.TotalPages = (result.TotalItems + q.PageSize - 1) / q.PageSize
	}
	return result, err
}

// InspectionRepository
type InspectionRepository struct{ db *sql.DB }

func (r InspectionRepository) Create(ctx context.Context, item domain.TireInspection) (domain.TireInspection, error) {
	query := `INSERT INTO tire_inspections AS t (
		id, tenant_id, tire_id, inspection_date, tread_mm, pressure, temperature, condition, observations, inspector, created_at, updated_at
	) VALUES (
		$1, $2::uuid, $3::uuid, $4, $5, $6, $7, $8, $9, $10, $11, $12
	) RETURNING row_to_json(t)`
	return read[domain.TireInspection](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.TireID, item.InspectionDate, item.TreadMM, item.Pressure, item.Temperature,
		item.Condition, item.Observations, item.Inspector, item.CreatedAt, item.UpdatedAt,
	))
}

func (r InspectionRepository) FindByID(ctx context.Context, tenantID, tireID, inspectionID string) (domain.TireInspection, error) {
	query := `SELECT row_to_json(t) FROM tire_inspections t WHERE tenant_id=$1::uuid AND tire_id=$2::uuid AND id=$3::uuid AND deleted_at IS NULL`
	return read[domain.TireInspection](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, tireID, inspectionID))
}

func (r InspectionRepository) List(ctx context.Context, tenantID, tireID string, q ports.Query) (ports.Page[domain.TireInspection], error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 20
	}
	result := ports.Page[domain.TireInspection]{Items: []domain.TireInspection{}, Page: q.Page, PageSize: q.PageSize}
	where := "tenant_id=$1::uuid AND tire_id=$2::uuid AND deleted_at IS NULL"
	args := []any{tenantID, tireID}

	err := (platform.Transactions{DB: r.db}).Run(ctx, func(ctx context.Context) error {
		exec := platform.ExecutorFor(ctx, r.db)
		if err := exec.QueryRowContext(ctx, "SELECT count(*) FROM tire_inspections WHERE "+where, args...).Scan(&result.TotalItems); err != nil {
			return mapError(err)
		}
		offset := (int64(q.Page) - 1) * int64(q.PageSize)
		rows, err := exec.QueryContext(ctx, fmt.Sprintf("SELECT row_to_json(t) FROM tire_inspections t WHERE %s ORDER BY inspection_date DESC, created_at DESC LIMIT $%d OFFSET $%d", where, len(args)+1, len(args)+2), tenantID, tireID, q.PageSize, offset)
		if err != nil {
			return mapError(err)
		}
		defer rows.Close()
		for rows.Next() {
			var body []byte
			var item domain.TireInspection
			if err := rows.Scan(&body); err != nil {
				return err
			}
			if err := json.Unmarshal(body, &item); err != nil {
				return err
			}
			result.Items = append(result.Items, item)
		}
		return rows.Err()
	})
	if result.TotalItems > 0 {
		result.TotalPages = (result.TotalItems + q.PageSize - 1) / q.PageSize
	}
	return result, err
}

func (r InspectionRepository) Update(ctx context.Context, item domain.TireInspection) (domain.TireInspection, error) {
	query := `UPDATE tire_inspections AS t SET
		inspection_date=$4, tread_mm=$5, pressure=$6, temperature=$7, condition=$8, observations=$9, inspector=$10, updated_at=$11
		WHERE tenant_id=$2::uuid AND tire_id=$3::uuid AND id=$1::uuid AND deleted_at IS NULL
		RETURNING row_to_json(t)`
	return read[domain.TireInspection](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.TireID, item.InspectionDate, item.TreadMM, item.Pressure, item.Temperature,
		item.Condition, item.Observations, item.Inspector, item.UpdatedAt,
	))
}

func (r InspectionRepository) Delete(ctx context.Context, tenantID, tireID, inspectionID string) error {
	query := `UPDATE tire_inspections SET deleted_at=now(), updated_at=now() WHERE tenant_id=$1::uuid AND tire_id=$2::uuid AND id=$3::uuid AND deleted_at IS NULL`
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, query, tenantID, tireID, inspectionID))
}

// MovementRepository
type MovementRepository struct{ db *sql.DB }

func (r MovementRepository) Create(ctx context.Context, item domain.TireMovement) (domain.TireMovement, error) {
	query := `INSERT INTO tire_movements AS t (
		id, tenant_id, tire_id, movement_type, vehicle_id, position, km, reason, performed_by, movement_date, created_at, updated_at
	) VALUES (
		$1, $2::uuid, $3::uuid, $4, NULLIF($5, '')::uuid, $6, $7, $8, NULLIF($9, '')::uuid, $10, $11, $12
	) RETURNING row_to_json(t)`
	return read[domain.TireMovement](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.TireID, item.MovementType, item.VehicleID, item.Position, item.KM,
		item.Reason, item.PerformedBy, item.MovementDate, item.CreatedAt, item.UpdatedAt,
	))
}

func (r MovementRepository) List(ctx context.Context, tenantID, tireID string, q ports.Query) (ports.Page[domain.TireMovement], error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 20
	}
	result := ports.Page[domain.TireMovement]{Items: []domain.TireMovement{}, Page: q.Page, PageSize: q.PageSize}
	where := "tenant_id=$1::uuid AND tire_id=$2::uuid AND deleted_at IS NULL"

	err := (platform.Transactions{DB: r.db}).Run(ctx, func(ctx context.Context) error {
		exec := platform.ExecutorFor(ctx, r.db)
		if err := exec.QueryRowContext(ctx, "SELECT count(*) FROM tire_movements WHERE "+where, tenantID, tireID).Scan(&result.TotalItems); err != nil {
			return mapError(err)
		}
		offset := (int64(q.Page) - 1) * int64(q.PageSize)
		rows, err := exec.QueryContext(ctx, fmt.Sprintf("SELECT row_to_json(t) FROM tire_movements t WHERE %s ORDER BY movement_date DESC, created_at DESC LIMIT $3 OFFSET $4", where), tenantID, tireID, q.PageSize, offset)
		if err != nil {
			return mapError(err)
		}
		defer rows.Close()
		for rows.Next() {
			var body []byte
			var item domain.TireMovement
			if err := rows.Scan(&body); err != nil {
				return err
			}
			if err := json.Unmarshal(body, &item); err != nil {
				return err
			}
			result.Items = append(result.Items, item)
		}
		return rows.Err()
	})
	if result.TotalItems > 0 {
		result.TotalPages = (result.TotalItems + q.PageSize - 1) / q.PageSize
	}
	return result, err
}

// HistoryRepository
type HistoryRepository struct{ db *sql.DB }

func (r HistoryRepository) Create(ctx context.Context, item domain.TireHistory) (domain.TireHistory, error) {
	query := `INSERT INTO tire_history AS t (
		id, tenant_id, tire_id, event, from_status, to_status, actor_id, notes, created_at
	) VALUES (
		$1, $2::uuid, $3::uuid, $4, $5, $6, NULLIF($7, '')::uuid, $8, $9
	) RETURNING row_to_json(t)`
	return read[domain.TireHistory](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.TireID, item.Event, item.FromStatus, item.ToStatus, item.ActorID, item.Notes, item.CreatedAt,
	))
}

func (r HistoryRepository) List(ctx context.Context, tenantID, tireID string, q ports.Query) (ports.Page[domain.TireHistory], error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 20
	}
	result := ports.Page[domain.TireHistory]{Items: []domain.TireHistory{}, Page: q.Page, PageSize: q.PageSize}
	where := "tenant_id=$1::uuid AND tire_id=$2::uuid"

	err := (platform.Transactions{DB: r.db}).Run(ctx, func(ctx context.Context) error {
		exec := platform.ExecutorFor(ctx, r.db)
		if err := exec.QueryRowContext(ctx, "SELECT count(*) FROM tire_history WHERE "+where, tenantID, tireID).Scan(&result.TotalItems); err != nil {
			return mapError(err)
		}
		offset := (int64(q.Page) - 1) * int64(q.PageSize)
		rows, err := exec.QueryContext(ctx, fmt.Sprintf("SELECT row_to_json(t) FROM tire_history t WHERE %s ORDER BY created_at DESC LIMIT $3 OFFSET $4", where), tenantID, tireID, q.PageSize, offset)
		if err != nil {
			return mapError(err)
		}
		defer rows.Close()
		for rows.Next() {
			var body []byte
			var item domain.TireHistory
			if err := rows.Scan(&body); err != nil {
				return err
			}
			if err := json.Unmarshal(body, &item); err != nil {
				return err
			}
			result.Items = append(result.Items, item)
		}
		return rows.Err()
	})
	if result.TotalItems > 0 {
		result.TotalPages = (result.TotalItems + q.PageSize - 1) / q.PageSize
	}
	return result, err
}


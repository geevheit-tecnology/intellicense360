package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/ports"
	platform "github.com/geevheit/intelligence360/backend/api/internal/platform/postgres"
	"github.com/lib/pq"
	"strings"
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
			case "ux_fleet_vehicles_plate_active":
				return application.ErrLicensePlateTaken
			case "ux_fleet_vehicles_chassis_active":
				return application.ErrChassisTaken
			case "ux_fleet_vehicles_renavam_active":
				return application.ErrRenavamTaken
			}
			return application.ErrConflict
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

// The table/search/sort fragments below are selected exclusively by this package.
func search[T any](ctx context.Context, db *sql.DB, table, searchColumn string, tenant string, q ports.Query) (ports.Page[T], error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 20
	}
	result := ports.Page[T]{Items: []T{}, Page: q.Page, PageSize: q.PageSize}
	args := []any{tenant}
	where := "tenant_id=$1 AND deleted_at IS NULL"
	if q.Search != "" {
		args = append(args, q.Search)
		where += fmt.Sprintf(" AND (%s) ILIKE '%%' || $%d || '%%'", searchColumn, len(args))
	}
	for _, key := range []string{"status", "category_id"} {
		value := q.Filters[key]
		if value == "" {
			continue
		}
		if (key == "status" && (table == "fleet_vehicles" || table == "fleet_assets")) || (key == "category_id" && (table == "fleet_vehicles" || table == "fleet_vehicle_models" || table == "fleet_vehicle_types")) {
			args = append(args, value)
			if key == "category_id" {
				where += fmt.Sprintf(" AND %s=$%d::uuid", key, len(args))
			} else {
				where += fmt.Sprintf(" AND %s=$%d", key, len(args))
			}
		}
	}
	// Keep count and page in one transaction.
	err := (platform.Transactions{DB: db}).Run(ctx, func(ctx context.Context) error {
		exec := platform.ExecutorFor(ctx, db)
		if err := exec.QueryRowContext(ctx, "SELECT count(*) FROM "+table+" WHERE "+where, args...).Scan(&result.TotalItems); err != nil {
			return mapError(err)
		}
		sortBy := "updated_at"
		switch q.SortBy {
		case "created_at", "updated_at", "id":
			sortBy = q.SortBy
		case "name":
			if table != "fleet_vehicles" {
				sortBy = "name"
			}
		case "license_plate":
			if table == "fleet_vehicles" {
				sortBy = "license_plate"
			}
		}
		order := "ASC"
		if strings.EqualFold(q.SortOrder, "desc") {
			order = "DESC"
		}
		args = append(args, q.PageSize, (int64(q.Page)-1)*int64(q.PageSize))
		rows, err := exec.QueryContext(ctx, fmt.Sprintf("SELECT row_to_json(t) FROM %s t WHERE %s ORDER BY %s %s,id LIMIT $%d OFFSET $%d", table, where, sortBy, order, len(args)-1, len(args)), args...)
		if err != nil {
			return mapError(err)
		}
		defer rows.Close()
		for rows.Next() {
			var body []byte
			var item T
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

func (s *Store) Vehicles() ports.VehicleRepository { return VehicleRepository{db: s.db} }

type VehicleRepository struct{ db *sql.DB }

var _ ports.VehicleRepository = VehicleRepository{}

func (r VehicleRepository) Create(ctx context.Context, item domain.Vehicle) (domain.Vehicle, error) {
	metadata, err := json.Marshal(item.Metadata)
	if err != nil {
		return item, err
	}
	if item.Metadata == nil {
		metadata = []byte("{}")
	}
	return read[domain.Vehicle](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO fleet_vehicles AS t (id,tenant_id,audit_id,category_id,type_id,brand_id,model_id,asset_id,license_plate,renavam,chassis,engine,color,fuel_type,transmission,axle_configuration,emission_standard,ownership_type,status,year_manufacture,year_model,metadata,created_at,updated_at,version) VALUES ($1,$2,NULLIF($3,'')::uuid,NULLIF($4,'')::uuid,NULLIF($5,'')::uuid,NULLIF($6,'')::uuid,NULLIF($7,'')::uuid,NULLIF($8,'')::uuid,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22::jsonb,$23,$24,$25) RETURNING row_to_json(t)`, item.ID, item.TenantID, item.AuditID, item.CategoryID, item.TypeID, item.BrandID, item.ModelID, item.AssetID, item.LicensePlate, item.Renavam, item.Chassis, item.Engine, item.Color, item.FuelType, item.Transmission, item.AxleConfiguration, item.EmissionStandard, item.OwnershipType, item.Status, item.YearManufacture, item.YearModel, string(metadata), item.CreatedAt, item.UpdatedAt, item.Version))
}
func (r VehicleRepository) FindByID(ctx context.Context, tenant, id string) (domain.Vehicle, error) {
	return read[domain.Vehicle](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `SELECT row_to_json(t) FROM fleet_vehicles t WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenant, id))
}
func (r VehicleRepository) Update(ctx context.Context, item domain.Vehicle) (domain.Vehicle, error) {
	metadata, err := json.Marshal(item.Metadata)
	if err != nil {
		return item, err
	}
	if item.Metadata == nil {
		metadata = []byte("{}")
	}
	saved, err := read[domain.Vehicle](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `UPDATE fleet_vehicles AS t SET audit_id=NULLIF($3,'')::uuid,category_id=NULLIF($4,'')::uuid,type_id=NULLIF($5,'')::uuid,brand_id=NULLIF($6,'')::uuid,model_id=NULLIF($7,'')::uuid,asset_id=NULLIF($8,'')::uuid,license_plate=$9,renavam=$10,chassis=$11,engine=$12,color=$13,fuel_type=$14,transmission=$15,axle_configuration=$16,emission_standard=$17,ownership_type=$18,status=$19,year_manufacture=$20,year_model=$21,metadata=$22::jsonb,updated_at=$23,version=t.version+1 WHERE tenant_id=$2::uuid AND id=$1::uuid AND deleted_at IS NULL AND version=$24-1 RETURNING row_to_json(t)`, item.ID, item.TenantID, item.AuditID, item.CategoryID, item.TypeID, item.BrandID, item.ModelID, item.AssetID, item.LicensePlate, item.Renavam, item.Chassis, item.Engine, item.Color, item.FuelType, item.Transmission, item.AxleConfiguration, item.EmissionStandard, item.OwnershipType, item.Status, item.YearManufacture, item.YearModel, string(metadata), item.UpdatedAt, item.Version))
	if errors.Is(err, application.ErrNotFound) {
		if _, findErr := r.FindByID(ctx, item.TenantID, item.ID); findErr == nil {
			return saved, application.ErrConflict
		}
	}
	return saved, err
}
func (r VehicleRepository) Delete(ctx context.Context, tenant, id string) error {
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `UPDATE fleet_vehicles SET deleted_at=now(),updated_at=now(),version=version+1 WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenant, id))
}
func (r VehicleRepository) Search(ctx context.Context, tenant string, q ports.Query) (ports.Page[domain.Vehicle], error) {
	return search[domain.Vehicle](ctx, r.db, "fleet_vehicles", "license_plate || ' ' || chassis || ' ' || renavam", tenant, q)
}

func (s *Store) Brands() ports.VehicleBrandRepository { return VehicleBrandRepository{db: s.db} }

type VehicleBrandRepository struct{ db *sql.DB }

var _ ports.VehicleBrandRepository = VehicleBrandRepository{}

func (r VehicleBrandRepository) Create(ctx context.Context, item domain.VehicleBrand) (domain.VehicleBrand, error) {
	return read[domain.VehicleBrand](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO fleet_vehicle_brands AS t (id,tenant_id,name,created_at,updated_at,version) VALUES ($1,$2,$3,$4,$5,$6) RETURNING row_to_json(t)`, item.ID, item.TenantID, item.Name, item.CreatedAt, item.UpdatedAt, item.Version))
}
func (r VehicleBrandRepository) FindByID(ctx context.Context, tenant, id string) (domain.VehicleBrand, error) {
	return read[domain.VehicleBrand](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `SELECT row_to_json(t) FROM fleet_vehicle_brands t WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenant, id))
}
func (r VehicleBrandRepository) Update(ctx context.Context, item domain.VehicleBrand) (domain.VehicleBrand, error) {
	saved, err := read[domain.VehicleBrand](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `UPDATE fleet_vehicle_brands AS t SET name=$3,updated_at=$4,version=t.version+1 WHERE tenant_id=$2::uuid AND id=$1::uuid AND deleted_at IS NULL AND version=$5-1 RETURNING row_to_json(t)`, item.ID, item.TenantID, item.Name, item.UpdatedAt, item.Version))
	if errors.Is(err, application.ErrNotFound) {
		if _, findErr := r.FindByID(ctx, item.TenantID, item.ID); findErr == nil {
			return saved, application.ErrConflict
		}
	}
	return saved, err
}
func (r VehicleBrandRepository) Delete(ctx context.Context, tenant, id string) error {
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `UPDATE fleet_vehicle_brands SET deleted_at=now(),updated_at=now(),version=version+1 WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenant, id))
}
func (r VehicleBrandRepository) Search(ctx context.Context, tenant string, q ports.Query) (ports.Page[domain.VehicleBrand], error) {
	return search[domain.VehicleBrand](ctx, r.db, "fleet_vehicle_brands", "name", tenant, q)
}

func (s *Store) Models() ports.VehicleModelRepository { return VehicleModelRepository{db: s.db} }

type VehicleModelRepository struct{ db *sql.DB }

var _ ports.VehicleModelRepository = VehicleModelRepository{}

func (r VehicleModelRepository) Create(ctx context.Context, item domain.VehicleModel) (domain.VehicleModel, error) {
	return read[domain.VehicleModel](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO fleet_vehicle_models AS t (id,tenant_id,brand_id,name,category_id,type_id,created_at,updated_at,version) VALUES ($1,$2,NULLIF($3,'')::uuid,$4,NULLIF($5,'')::uuid,NULLIF($6,'')::uuid,$7,$8,$9) RETURNING row_to_json(t)`, item.ID, item.TenantID, item.BrandID, item.Name, item.CategoryID, item.TypeID, item.CreatedAt, item.UpdatedAt, item.Version))
}
func (r VehicleModelRepository) FindByID(ctx context.Context, tenant, id string) (domain.VehicleModel, error) {
	return read[domain.VehicleModel](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `SELECT row_to_json(t) FROM fleet_vehicle_models t WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenant, id))
}
func (r VehicleModelRepository) Update(ctx context.Context, item domain.VehicleModel) (domain.VehicleModel, error) {
	saved, err := read[domain.VehicleModel](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `UPDATE fleet_vehicle_models AS t SET brand_id=NULLIF($3,'')::uuid,name=$4,category_id=NULLIF($5,'')::uuid,type_id=NULLIF($6,'')::uuid,updated_at=$7,version=t.version+1 WHERE tenant_id=$2::uuid AND id=$1::uuid AND deleted_at IS NULL AND version=$8-1 RETURNING row_to_json(t)`, item.ID, item.TenantID, item.BrandID, item.Name, item.CategoryID, item.TypeID, item.UpdatedAt, item.Version))
	if errors.Is(err, application.ErrNotFound) {
		if _, findErr := r.FindByID(ctx, item.TenantID, item.ID); findErr == nil {
			return saved, application.ErrConflict
		}
	}
	return saved, err
}
func (r VehicleModelRepository) Delete(ctx context.Context, tenant, id string) error {
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `UPDATE fleet_vehicle_models SET deleted_at=now(),updated_at=now(),version=version+1 WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenant, id))
}
func (r VehicleModelRepository) Search(ctx context.Context, tenant string, q ports.Query) (ports.Page[domain.VehicleModel], error) {
	return search[domain.VehicleModel](ctx, r.db, "fleet_vehicle_models", "name", tenant, q)
}

func (s *Store) Categories() ports.VehicleCategoryRepository {
	return VehicleCategoryRepository{db: s.db}
}

type VehicleCategoryRepository struct{ db *sql.DB }

var _ ports.VehicleCategoryRepository = VehicleCategoryRepository{}

func (r VehicleCategoryRepository) Create(ctx context.Context, item domain.VehicleCategoryEntity) (domain.VehicleCategoryEntity, error) {
	return read[domain.VehicleCategoryEntity](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO fleet_vehicle_categories AS t (id,tenant_id,name,code,created_at,updated_at,version) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING row_to_json(t)`, item.ID, item.TenantID, item.Name, item.Code, item.CreatedAt, item.UpdatedAt, item.Version))
}
func (r VehicleCategoryRepository) FindByID(ctx context.Context, tenant, id string) (domain.VehicleCategoryEntity, error) {
	return read[domain.VehicleCategoryEntity](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `SELECT row_to_json(t) FROM fleet_vehicle_categories t WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenant, id))
}
func (r VehicleCategoryRepository) Update(ctx context.Context, item domain.VehicleCategoryEntity) (domain.VehicleCategoryEntity, error) {
	saved, err := read[domain.VehicleCategoryEntity](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `UPDATE fleet_vehicle_categories AS t SET name=$3,code=$4,updated_at=$5,version=t.version+1 WHERE tenant_id=$2::uuid AND id=$1::uuid AND deleted_at IS NULL AND version=$6-1 RETURNING row_to_json(t)`, item.ID, item.TenantID, item.Name, item.Code, item.UpdatedAt, item.Version))
	if errors.Is(err, application.ErrNotFound) {
		if _, findErr := r.FindByID(ctx, item.TenantID, item.ID); findErr == nil {
			return saved, application.ErrConflict
		}
	}
	return saved, err
}
func (r VehicleCategoryRepository) Delete(ctx context.Context, tenant, id string) error {
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `UPDATE fleet_vehicle_categories SET deleted_at=now(),updated_at=now(),version=version+1 WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenant, id))
}
func (r VehicleCategoryRepository) Search(ctx context.Context, tenant string, q ports.Query) (ports.Page[domain.VehicleCategoryEntity], error) {
	return search[domain.VehicleCategoryEntity](ctx, r.db, "fleet_vehicle_categories", "name", tenant, q)
}

func (s *Store) Types() ports.VehicleTypeRepository { return VehicleTypeRepository{db: s.db} }

type VehicleTypeRepository struct{ db *sql.DB }

var _ ports.VehicleTypeRepository = VehicleTypeRepository{}

func (r VehicleTypeRepository) Create(ctx context.Context, item domain.VehicleTypeEntity) (domain.VehicleTypeEntity, error) {
	return read[domain.VehicleTypeEntity](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO fleet_vehicle_types AS t (id,tenant_id,category_id,name,created_at,updated_at,version) VALUES ($1,$2,NULLIF($3,'')::uuid,$4,$5,$6,$7) RETURNING row_to_json(t)`, item.ID, item.TenantID, item.CategoryID, item.Name, item.CreatedAt, item.UpdatedAt, item.Version))
}
func (r VehicleTypeRepository) FindByID(ctx context.Context, tenant, id string) (domain.VehicleTypeEntity, error) {
	return read[domain.VehicleTypeEntity](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `SELECT row_to_json(t) FROM fleet_vehicle_types t WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenant, id))
}
func (r VehicleTypeRepository) Update(ctx context.Context, item domain.VehicleTypeEntity) (domain.VehicleTypeEntity, error) {
	saved, err := read[domain.VehicleTypeEntity](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `UPDATE fleet_vehicle_types AS t SET category_id=NULLIF($3,'')::uuid,name=$4,updated_at=$5,version=t.version+1 WHERE tenant_id=$2::uuid AND id=$1::uuid AND deleted_at IS NULL AND version=$6-1 RETURNING row_to_json(t)`, item.ID, item.TenantID, item.CategoryID, item.Name, item.UpdatedAt, item.Version))
	if errors.Is(err, application.ErrNotFound) {
		if _, findErr := r.FindByID(ctx, item.TenantID, item.ID); findErr == nil {
			return saved, application.ErrConflict
		}
	}
	return saved, err
}
func (r VehicleTypeRepository) Delete(ctx context.Context, tenant, id string) error {
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `UPDATE fleet_vehicle_types SET deleted_at=now(),updated_at=now(),version=version+1 WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenant, id))
}
func (r VehicleTypeRepository) Search(ctx context.Context, tenant string, q ports.Query) (ports.Page[domain.VehicleTypeEntity], error) {
	return search[domain.VehicleTypeEntity](ctx, r.db, "fleet_vehicle_types", "name", tenant, q)
}

func (s *Store) Assets() ports.AssetRepository { return AssetRepository{db: s.db} }

type AssetRepository struct{ db *sql.DB }

var _ ports.AssetRepository = AssetRepository{}

func (r AssetRepository) Create(ctx context.Context, item domain.Asset) (domain.Asset, error) {
	return read[domain.Asset](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO fleet_assets AS t (id,tenant_id,code,name,description,status,created_at,updated_at,version) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING row_to_json(t)`, item.ID, item.TenantID, item.Code, item.Name, item.Description, item.Status, item.CreatedAt, item.UpdatedAt, item.Version))
}
func (r AssetRepository) FindByID(ctx context.Context, tenant, id string) (domain.Asset, error) {
	return read[domain.Asset](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `SELECT row_to_json(t) FROM fleet_assets t WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenant, id))
}
func (r AssetRepository) Update(ctx context.Context, item domain.Asset) (domain.Asset, error) {
	saved, err := read[domain.Asset](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `UPDATE fleet_assets AS t SET code=$3,name=$4,description=$5,status=$6,updated_at=$7,version=t.version+1 WHERE tenant_id=$2::uuid AND id=$1::uuid AND deleted_at IS NULL AND version=$8-1 RETURNING row_to_json(t)`, item.ID, item.TenantID, item.Code, item.Name, item.Description, item.Status, item.UpdatedAt, item.Version))
	if errors.Is(err, application.ErrNotFound) {
		if _, findErr := r.FindByID(ctx, item.TenantID, item.ID); findErr == nil {
			return saved, application.ErrConflict
		}
	}
	return saved, err
}
func (r AssetRepository) Delete(ctx context.Context, tenant, id string) error {
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `UPDATE fleet_assets SET deleted_at=now(),updated_at=now(),version=version+1 WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenant, id))
}
func (r AssetRepository) Search(ctx context.Context, tenant string, q ports.Query) (ports.Page[domain.Asset], error) {
	return search[domain.Asset](ctx, r.db, "fleet_assets", "name || ' ' || code", tenant, q)
}

func (r VehicleRepository) Exists(ctx context.Context, tenant, id string) (bool, error) {
	var exists bool
	err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM fleet_vehicles WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL)`, tenant, id).Scan(&exists)
	return exists, err
}
func (r VehicleRepository) existsField(ctx context.Context, tenant, column, value, except string) (bool, error) {
	var exists bool
	err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM fleet_vehicles WHERE tenant_id=$1::uuid AND `+column+`=$2 AND id IS DISTINCT FROM NULLIF($3,'')::uuid AND deleted_at IS NULL)`, tenant, value, except).Scan(&exists)
	return exists, err
}
func (r VehicleRepository) ExistsLicensePlate(ctx context.Context, tenant string, value domain.LicensePlate, except string) (bool, error) {
	return r.existsField(ctx, tenant, "upper(replace(license_plate, '-', ''))", value.Normalized(), except)
}
func (r VehicleRepository) ExistsChassis(ctx context.Context, tenant string, value domain.Chassis, except string) (bool, error) {
	return r.existsField(ctx, tenant, "upper(chassis)", strings.ToUpper(strings.TrimSpace(string(value))), except)
}
func (r VehicleRepository) ExistsRenavam(ctx context.Context, tenant string, value domain.Renavam, except string) (bool, error) {
	return r.existsField(ctx, tenant, "renavam", strings.TrimSpace(string(value)), except)
}

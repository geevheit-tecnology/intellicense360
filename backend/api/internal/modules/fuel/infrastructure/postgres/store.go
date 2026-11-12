package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/ports"
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
			return application.ErrConflict
		}
		if pg.Code == "23503" || pg.Code == "23502" || pg.Code == "23514" || pg.Code == "22P02" {
			return application.ErrInvalidFuelData
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
	if q.Search != "" && searchColumn != "" {
		args = append(args, q.Search)
		where += fmt.Sprintf(" AND (%s) ILIKE '%%' || $%d || '%%'", searchColumn, len(args))
	}
	for key, value := range q.Filters {
		if value == "" {
			continue
		}
		if key == "status" || key == "fuel_kind" || key == "payment_status" || key == "station_id" {
			args = append(args, value)
			where += fmt.Sprintf(" AND %s=$%d", key, len(args))
		}
	}

	err := (platform.Transactions{DB: db}).Run(ctx, func(ctx context.Context) error {
		exec := platform.ExecutorFor(ctx, db)
		if err := exec.QueryRowContext(ctx, "SELECT count(*) FROM "+table+" WHERE "+where, args...).Scan(&result.TotalItems); err != nil {
			return mapError(err)
		}
		sortBy := "created_at"
		switch q.SortBy {
		case "created_at", "updated_at", "id", "transaction_date", "receipt_date":
			sortBy = q.SortBy
		}
		order := "DESC"
		if strings.EqualFold(q.SortOrder, "asc") {
			order = "ASC"
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

func (s *Store) Transactions() ports.FuelTransactionRepository {
	return TransactionRepository{db: s.db}
}
func (s *Store) Types() ports.FuelTypeRepository {
	return TypeRepository{db: s.db}
}
func (s *Store) Stations() ports.FuelStationRepository {
	return StationRepository{db: s.db}
}
func (s *Store) Tanks() ports.FuelTankRepository {
	return TankRepository{db: s.db}
}
func (s *Store) Nozzles() ports.FuelNozzleRepository {
	return NozzleRepository{db: s.db}
}
func (s *Store) Readings() ports.FuelReadingRepository {
	return ReadingRepository{db: s.db}
}
func (s *Store) Prices() ports.FuelPriceRepository {
	return PriceRepository{db: s.db}
}
func (s *Store) Receipts() ports.FuelReceiptRepository {
	return ReceiptRepository{db: s.db}
}
func (s *Store) Adjustments() ports.FuelAdjustmentRepository {
	return AdjustmentRepository{db: s.db}
}

// TransactionRepository
type TransactionRepository struct{ db *sql.DB }

func (r TransactionRepository) Create(ctx context.Context, item domain.FuelTransaction) (domain.FuelTransaction, error) {
	query := `INSERT INTO fuel_transactions AS t (
		id, tenant_id, transaction_date, fuel_type_id, fuel_kind, quantity, unit_price, total_amount,
		odometer_reading, engine_hour_reading, station_id, tank_id, nozzle_id, receipt_id, receipt_number,
		driver_reference, vehicle_reference, asset_reference, payment_method, notes, status,
		cancellation_reason, completed_at, canceled_at, created_at, updated_at, version,
		requisition_number, is_third_party, third_party_identifier, invoice_number, invoice_key, payment_status
	) VALUES (
		$1, $2::uuid, $3, NULLIF($4, '')::uuid, $5, $6, $7, $8,
		$9, $10, NULLIF($11, '')::uuid, NULLIF($12, '')::uuid, NULLIF($13, '')::uuid, NULLIF($14, '')::uuid, $15,
		$16, $17, $18, $19, $20, $21,
		$22, $23, $24, $25, $26, $27,
		$28, $29, $30, $31, $32, $33
	) RETURNING row_to_json(t)`
	return read[domain.FuelTransaction](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.TransactionDate, item.FuelTypeID, item.FuelKind, item.Quantity, item.UnitPrice, item.TotalAmount,
		item.OdometerReading, item.EngineHourReading, item.StationID, item.TankID, item.NozzleID, item.ReceiptID, item.ReceiptNumber,
		item.DriverReference, item.VehicleReference, item.AssetReference, item.PaymentMethod, item.Notes, item.Status,
		item.CancellationReason, item.CompletedAt, item.CanceledAt, item.CreatedAt, item.UpdatedAt, item.Version,
		item.RequisitionNumber, item.IsThirdParty, item.ThirdPartyIdentifier, item.InvoiceNumber, item.InvoiceKey, item.PaymentStatus,
	))
}

func (r TransactionRepository) FindByID(ctx context.Context, tenantID, id string) (domain.FuelTransaction, error) {
	query := `SELECT row_to_json(t) FROM fuel_transactions t WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`
	return read[domain.FuelTransaction](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, id))
}

func (r TransactionRepository) Update(ctx context.Context, item domain.FuelTransaction) (domain.FuelTransaction, error) {
	query := `UPDATE fuel_transactions AS t SET
		transaction_date=$3, fuel_type_id=NULLIF($4, '')::uuid, fuel_kind=$5, quantity=$6, unit_price=$7, total_amount=$8,
		odometer_reading=$9, engine_hour_reading=$10, station_id=NULLIF($11, '')::uuid, tank_id=NULLIF($12, '')::uuid,
		nozzle_id=NULLIF($13, '')::uuid, receipt_id=NULLIF($14, '')::uuid, receipt_number=$15, driver_reference=$16,
		vehicle_reference=$17, asset_reference=$18, payment_method=$19, notes=$20, status=$21,
		cancellation_reason=$22, completed_at=$23, canceled_at=$24, updated_at=$25, version=t.version+1,
		requisition_number=$26, is_third_party=$27, third_party_identifier=$28, invoice_number=$29, invoice_key=$30, payment_status=$31
		WHERE tenant_id=$2::uuid AND id=$1::uuid AND deleted_at IS NULL AND version=$32-1
		RETURNING row_to_json(t)`
	saved, err := read[domain.FuelTransaction](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.TransactionDate, item.FuelTypeID, item.FuelKind, item.Quantity, item.UnitPrice, item.TotalAmount,
		item.OdometerReading, item.EngineHourReading, item.StationID, item.TankID, item.NozzleID, item.ReceiptID, item.ReceiptNumber,
		item.DriverReference, item.VehicleReference, item.AssetReference, item.PaymentMethod, item.Notes, item.Status,
		item.CancellationReason, item.CompletedAt, item.CanceledAt, item.UpdatedAt,
		item.RequisitionNumber, item.IsThirdParty, item.ThirdPartyIdentifier, item.InvoiceNumber, item.InvoiceKey, item.PaymentStatus,
		item.Version,
	))
	if errors.Is(err, application.ErrNotFound) {
		if _, findErr := r.FindByID(ctx, item.TenantID, item.ID); findErr == nil {
			return saved, application.ErrConflict
		}
	}
	return saved, err
}

func (r TransactionRepository) Delete(ctx context.Context, tenantID, id string) error {
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `UPDATE fuel_transactions SET deleted_at=now(), updated_at=now(), version=version+1 WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenantID, id))
}

func (r TransactionRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelTransaction], error) {
	return search[domain.FuelTransaction](ctx, r.db, "fuel_transactions", "receipt_number || ' ' || vehicle_reference || ' ' || requisition_number || ' ' || invoice_number", tenantID, q)
}

// TypeRepository
type TypeRepository struct{ db *sql.DB }

func (r TypeRepository) Create(ctx context.Context, item domain.FuelType) (domain.FuelType, error) {
	return read[domain.FuelType](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO fuel_types AS t (id, tenant_id, name, kind, code, description, active, created_at, updated_at, version) VALUES ($1, $2::uuid, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING row_to_json(t)`, item.ID, item.TenantID, item.Name, item.Kind, item.Code, item.Description, item.Active, item.CreatedAt, item.UpdatedAt, item.Version))
}
func (r TypeRepository) FindByID(ctx context.Context, tenantID, id string) (domain.FuelType, error) {
	return read[domain.FuelType](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `SELECT row_to_json(t) FROM fuel_types t WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenantID, id))
}
func (r TypeRepository) Update(ctx context.Context, item domain.FuelType) (domain.FuelType, error) {
	return read[domain.FuelType](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `UPDATE fuel_types AS t SET name=$3, kind=$4, code=$5, description=$6, active=$7, updated_at=$8, version=t.version+1 WHERE tenant_id=$2::uuid AND id=$1::uuid AND deleted_at IS NULL AND version=$9-1 RETURNING row_to_json(t)`, item.ID, item.TenantID, item.Name, item.Kind, item.Code, item.Description, item.Active, item.UpdatedAt, item.Version))
}
func (r TypeRepository) Delete(ctx context.Context, tenantID, id string) error {
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `UPDATE fuel_types SET deleted_at=now(), updated_at=now(), version=version+1 WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenantID, id))
}
func (r TypeRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelType], error) {
	return search[domain.FuelType](ctx, r.db, "fuel_types", "name || ' ' || code", tenantID, q)
}

// StationRepository
type StationRepository struct{ db *sql.DB }

func (r StationRepository) Create(ctx context.Context, item domain.FuelStation) (domain.FuelStation, error) {
	return read[domain.FuelStation](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO fuel_stations AS t (id, tenant_id, name, legal_name, cnpj, address, city, state, country, latitude, longitude, active, notes, created_at, updated_at, version) VALUES ($1, $2::uuid, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING row_to_json(t)`, item.ID, item.TenantID, item.Name, item.LegalName, item.CNPJ, item.Address, item.City, item.State, item.Country, item.Latitude, item.Longitude, item.Active, item.Notes, item.CreatedAt, item.UpdatedAt, item.Version))
}
func (r StationRepository) FindByID(ctx context.Context, tenantID, id string) (domain.FuelStation, error) {
	return read[domain.FuelStation](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `SELECT row_to_json(t) FROM fuel_stations t WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenantID, id))
}
func (r StationRepository) Update(ctx context.Context, item domain.FuelStation) (domain.FuelStation, error) {
	return read[domain.FuelStation](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `UPDATE fuel_stations AS t SET name=$3, legal_name=$4, cnpj=$5, address=$6, city=$7, state=$8, country=$9, latitude=$10, longitude=$11, active=$12, notes=$13, updated_at=$14, version=t.version+1 WHERE tenant_id=$2::uuid AND id=$1::uuid AND deleted_at IS NULL AND version=$15-1 RETURNING row_to_json(t)`, item.ID, item.TenantID, item.Name, item.LegalName, item.CNPJ, item.Address, item.City, item.State, item.Country, item.Latitude, item.Longitude, item.Active, item.Notes, item.UpdatedAt, item.Version))
}
func (r StationRepository) Delete(ctx context.Context, tenantID, id string) error {
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `UPDATE fuel_stations SET deleted_at=now(), updated_at=now(), version=version+1 WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenantID, id))
}
func (r StationRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelStation], error) {
	return search[domain.FuelStation](ctx, r.db, "fuel_stations", "name || ' ' || legal_name || ' ' || coalesce(cnpj, '')", tenantID, q)
}
func (r StationRepository) ExistsCNPJ(ctx context.Context, tenantID, cnpj, exceptID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM fuel_stations WHERE tenant_id=$1::uuid AND cnpj=$2 AND deleted_at IS NULL AND id <> NULLIF($3, '')::uuid)`
	err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, cnpj, exceptID).Scan(&exists)
	return exists, mapError(err)
}

// TankRepository
type TankRepository struct{ db *sql.DB }

func (r TankRepository) Create(ctx context.Context, item domain.FuelTank) (domain.FuelTank, error) {
	return read[domain.FuelTank](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO fuel_tanks AS t (id, tenant_id, code, name, capacity, current_reading, fuel_type_id, fuel_kind, location_reference, status, notes, created_at, updated_at, version) VALUES ($1, $2::uuid, $3, $4, $5, $6, NULLIF($7, '')::uuid, $8, $9, $10, $11, $12, $13, $14) RETURNING row_to_json(t)`, item.ID, item.TenantID, item.Code, item.Name, item.Capacity, item.CurrentReading, item.FuelTypeID, item.FuelKind, item.LocationRef, item.Status, item.Notes, item.CreatedAt, item.UpdatedAt, item.Version))
}
func (r TankRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelTank], error) {
	return search[domain.FuelTank](ctx, r.db, "fuel_tanks", "code || ' ' || coalesce(name, '')", tenantID, q)
}
func (r TankRepository) Update(ctx context.Context, item domain.FuelTank) (domain.FuelTank, error) {
	return read[domain.FuelTank](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `UPDATE fuel_tanks AS t SET code=$3, name=$4, capacity=$5, current_reading=$6, fuel_type_id=NULLIF($7, '')::uuid, fuel_kind=$8, location_reference=$9, status=$10, notes=$11, updated_at=$12, version=t.version+1 WHERE tenant_id=$2::uuid AND id=$1::uuid AND deleted_at IS NULL AND version=$13-1 RETURNING row_to_json(t)`, item.ID, item.TenantID, item.Code, item.Name, item.Capacity, item.CurrentReading, item.FuelTypeID, item.FuelKind, item.LocationRef, item.Status, item.Notes, item.UpdatedAt, item.Version))
}
func (r TankRepository) Delete(ctx context.Context, tenantID, id string) error {
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `UPDATE fuel_tanks SET deleted_at=now(), updated_at=now(), version=version+1 WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenantID, id))
}

// NozzleRepository
type NozzleRepository struct{ db *sql.DB }

func (r NozzleRepository) Create(ctx context.Context, item domain.FuelNozzle) (domain.FuelNozzle, error) {
	return read[domain.FuelNozzle](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO fuel_nozzles AS t (id, tenant_id, code, fuel_type_id, fuel_kind, tank_id, status, meter_reading, notes, created_at, updated_at, version) VALUES ($1, $2::uuid, $3, NULLIF($4, '')::uuid, $5, NULLIF($6, '')::uuid, $7, $8, $9, $10, $11, $12) RETURNING row_to_json(t)`, item.ID, item.TenantID, item.Code, item.FuelTypeID, item.FuelKind, item.TankID, item.Status, item.MeterReading, item.Notes, item.CreatedAt, item.UpdatedAt, item.Version))
}
func (r NozzleRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelNozzle], error) {
	return search[domain.FuelNozzle](ctx, r.db, "fuel_nozzles", "code", tenantID, q)
}
func (r NozzleRepository) Update(ctx context.Context, item domain.FuelNozzle) (domain.FuelNozzle, error) {
	return read[domain.FuelNozzle](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `UPDATE fuel_nozzles AS t SET code=$3, fuel_type_id=NULLIF($4, '')::uuid, fuel_kind=$5, tank_id=NULLIF($6, '')::uuid, status=$7, meter_reading=$8, notes=$9, updated_at=$10, version=t.version+1 WHERE tenant_id=$2::uuid AND id=$1::uuid AND deleted_at IS NULL AND version=$11-1 RETURNING row_to_json(t)`, item.ID, item.TenantID, item.Code, item.FuelTypeID, item.FuelKind, item.TankID, item.Status, item.MeterReading, item.Notes, item.UpdatedAt, item.Version))
}
func (r NozzleRepository) Delete(ctx context.Context, tenantID, id string) error {
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `UPDATE fuel_nozzles SET deleted_at=now(), updated_at=now(), version=version+1 WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenantID, id))
}

// ReadingRepository
type ReadingRepository struct{ db *sql.DB }

func (r ReadingRepository) Create(ctx context.Context, item domain.FuelReading) (domain.FuelReading, error) {
	return read[domain.FuelReading](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO fuel_readings AS t (id, tenant_id, reading_type, reference_id, value, reading_date, source, notes, created_at) VALUES ($1, $2::uuid, $3, $4, $5, $6, $7, $8, $9) RETURNING row_to_json(t)`, item.ID, item.TenantID, item.ReadingType, item.ReferenceID, item.Value, item.ReadingDate, item.Source, item.Notes, item.CreatedAt))
}
func (r ReadingRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelReading], error) {
	return search[domain.FuelReading](ctx, r.db, "fuel_readings", "reference_id || ' ' || source", tenantID, q)
}

// PriceRepository
type PriceRepository struct{ db *sql.DB }

func (r PriceRepository) Create(ctx context.Context, item domain.FuelPrice) (domain.FuelPrice, error) {
	return read[domain.FuelPrice](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO fuel_prices AS t (id, tenant_id, fuel_type_id, fuel_kind, unit_price, effective_date, station_id, source, notes, created_at, updated_at, version) VALUES ($1, $2::uuid, NULLIF($3, '')::uuid, $4, $5, $6, NULLIF($7, '')::uuid, $8, $9, $10, $11, $12) RETURNING row_to_json(t)`, item.ID, item.TenantID, item.FuelTypeID, item.FuelKind, item.UnitPrice, item.EffectiveDate, item.StationID, item.Source, item.Notes, item.CreatedAt, item.UpdatedAt, item.Version))
}
func (r PriceRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelPrice], error) {
	return search[domain.FuelPrice](ctx, r.db, "fuel_prices", "source", tenantID, q)
}
func (r PriceRepository) Update(ctx context.Context, item domain.FuelPrice) (domain.FuelPrice, error) {
	return read[domain.FuelPrice](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `UPDATE fuel_prices AS t SET fuel_type_id=NULLIF($3, '')::uuid, fuel_kind=$4, unit_price=$5, effective_date=$6, station_id=NULLIF($7, '')::uuid, source=$8, notes=$9, updated_at=$10, version=t.version+1 WHERE tenant_id=$2::uuid AND id=$1::uuid AND deleted_at IS NULL AND version=$11-1 RETURNING row_to_json(t)`, item.ID, item.TenantID, item.FuelTypeID, item.FuelKind, item.UnitPrice, item.EffectiveDate, item.StationID, item.Source, item.Notes, item.UpdatedAt, item.Version))
}
func (r PriceRepository) Delete(ctx context.Context, tenantID, id string) error {
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `UPDATE fuel_prices SET deleted_at=now(), updated_at=now(), version=version+1 WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenantID, id))
}

// ReceiptRepository
type ReceiptRepository struct{ db *sql.DB }

func (r ReceiptRepository) Create(ctx context.Context, item domain.FuelReceipt) (domain.FuelReceipt, error) {
	return read[domain.FuelReceipt](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO fuel_receipts AS t (id, tenant_id, receipt_number, receipt_date, amount, attachment_reference, notes, created_at, updated_at, version) VALUES ($1, $2::uuid, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING row_to_json(t)`, item.ID, item.TenantID, item.ReceiptNumber, item.ReceiptDate, item.Amount, item.AttachmentReference, item.Notes, item.CreatedAt, item.UpdatedAt, item.Version))
}
func (r ReceiptRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelReceipt], error) {
	return search[domain.FuelReceipt](ctx, r.db, "fuel_receipts", "receipt_number", tenantID, q)
}
func (r ReceiptRepository) Update(ctx context.Context, item domain.FuelReceipt) (domain.FuelReceipt, error) {
	return read[domain.FuelReceipt](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `UPDATE fuel_receipts AS t SET receipt_number=$3, receipt_date=$4, amount=$5, attachment_reference=$6, notes=$7, updated_at=$8, version=t.version+1 WHERE tenant_id=$2::uuid AND id=$1::uuid AND deleted_at IS NULL AND version=$9-1 RETURNING row_to_json(t)`, item.ID, item.TenantID, item.ReceiptNumber, item.ReceiptDate, item.Amount, item.AttachmentReference, item.Notes, item.UpdatedAt, item.Version))
}
func (r ReceiptRepository) Delete(ctx context.Context, tenantID, id string) error {
	return changed(platform.ExecutorFor(ctx, r.db).ExecContext(ctx, `UPDATE fuel_receipts SET deleted_at=now(), updated_at=now(), version=version+1 WHERE tenant_id=$1::uuid AND id=$2::uuid AND deleted_at IS NULL`, tenantID, id))
}

// AdjustmentRepository
type AdjustmentRepository struct{ db *sql.DB }

func (r AdjustmentRepository) Create(ctx context.Context, item domain.FuelAdjustment) (domain.FuelAdjustment, error) {
	return read[domain.FuelAdjustment](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, `INSERT INTO fuel_adjustments AS t (id, tenant_id, transaction_id, adjustment_type, reason, original_reference, adjusted_quantity, adjusted_unit_price, adjusted_total_amount, notes, created_by, created_at) VALUES ($1, $2::uuid, $3::uuid, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING row_to_json(t)`, item.ID, item.TenantID, item.TransactionID, item.AdjustmentType, item.Reason, item.OriginalReference, item.AdjustedQuantity, item.AdjustedUnitPrice, item.AdjustedTotalAmount, item.Notes, item.CreatedBy, item.CreatedAt))
}
func (r AdjustmentRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.FuelAdjustment], error) {
	return search[domain.FuelAdjustment](ctx, r.db, "fuel_adjustments", "reason || ' ' || adjustment_type", tenantID, q)
}


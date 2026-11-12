package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/ports"
	platform "github.com/geevheit/intelligence360/backend/api/internal/platform/postgres"
)

func (s *Store) Odometer() ports.OdometerRepository { return OdometerRepository{db: s.db} }

type OdometerRepository struct{ db *sql.DB }

var _ ports.OdometerRepository = OdometerRepository{}

func (r OdometerRepository) Record(ctx context.Context, item domain.OdometerReading) (domain.OdometerReading, error) {
	driverID := ""
	if item.DriverID != nil {
		driverID = *item.DriverID
	}
	query := `INSERT INTO fleet_odometer_readings AS t 
		(id, tenant_id, vehicle_id, reading, read_at, source, driver_id, actor_id, notes, is_correction, correction_reason, created_at) 
		VALUES ($1, $2::uuid, $3::uuid, $4, $5, $6, NULLIF($7, '')::uuid, $8::uuid, $9, $10, $11, $12) 
		RETURNING row_to_json(t)`
	return read[domain.OdometerReading](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID,
		item.TenantID,
		item.VehicleID,
		item.Reading,
		item.ReadAt,
		item.Source,
		driverID,
		item.ActorID,
		item.Notes,
		item.IsCorrection,
		item.CorrectionReason,
		item.CreatedAt,
	))
}

func (r OdometerRepository) GetLatest(ctx context.Context, tenantID, vehicleID string) (*domain.OdometerReading, error) {
	query := `SELECT row_to_json(t) FROM fleet_odometer_readings t 
		WHERE tenant_id=$1::uuid AND vehicle_id=$2::uuid 
		ORDER BY read_at DESC, created_at DESC LIMIT 1`
	var body []byte
	err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, vehicleID).Scan(&body)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, mapError(err)
	}
	var item domain.OdometerReading
	if err := json.Unmarshal(body, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

func (r OdometerRepository) FindByID(ctx context.Context, tenantID, id string) (domain.OdometerReading, error) {
	query := `SELECT row_to_json(t) FROM fleet_odometer_readings t WHERE tenant_id=$1::uuid AND id=$2::uuid`
	return read[domain.OdometerReading](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, id))
}

func (r OdometerRepository) ListByVehicle(ctx context.Context, tenantID, vehicleID string, q ports.Query) (ports.Page[domain.OdometerReading], error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 20
	}
	result := ports.Page[domain.OdometerReading]{
		Items:    []domain.OdometerReading{},
		Page:     q.Page,
		PageSize: q.PageSize,
	}

	err := (platform.Transactions{DB: r.db}).Run(ctx, func(ctx context.Context) error {
		exec := platform.ExecutorFor(ctx, r.db)
		countQuery := `SELECT count(*) FROM fleet_odometer_readings WHERE tenant_id=$1::uuid AND vehicle_id=$2::uuid`
		if err := exec.QueryRowContext(ctx, countQuery, tenantID, vehicleID).Scan(&result.TotalItems); err != nil {
			return mapError(err)
		}

		offset := (int64(q.Page) - 1) * int64(q.PageSize)
		listQuery := `SELECT row_to_json(t) FROM fleet_odometer_readings t 
			WHERE tenant_id=$1::uuid AND vehicle_id=$2::uuid 
			ORDER BY read_at DESC, created_at DESC 
			LIMIT $3 OFFSET $4`

		rows, err := exec.QueryContext(ctx, listQuery, tenantID, vehicleID, q.PageSize, offset)
		if err != nil {
			return mapError(err)
		}
		defer rows.Close()

		for rows.Next() {
			var body []byte
			var item domain.OdometerReading
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

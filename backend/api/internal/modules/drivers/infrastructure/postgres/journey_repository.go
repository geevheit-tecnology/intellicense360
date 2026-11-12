package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/drivers/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/drivers/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/drivers/ports"
	platform "github.com/geevheit/intelligence360/backend/api/internal/platform/postgres"
)

type JourneyRepository struct {
	db *sql.DB
}

func NewJourneyRepository(db *sql.DB) JourneyRepository {
	return JourneyRepository{db: db}
}

var _ ports.JourneyRepository = JourneyRepository{}

func readJourney(row *sql.Row) (domain.DriverJourneyRecord, error) {
	var item domain.DriverJourneyRecord
	var body []byte
	if err := row.Scan(&body); err != nil {
		return item, mapError(err)
	}
	err := json.Unmarshal(body, &item)
	return item, err
}

func (r JourneyRepository) Create(ctx context.Context, item domain.DriverJourneyRecord) (domain.DriverJourneyRecord, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	query := `INSERT INTO driver_journey_records AS t (
		id, tenant_id, driver_id, vehicle_id, trip_id, event_type,
		start_at, end_at, duration_minutes, location, odometer_km,
		is_violation, violation_reason, daily_allowance_amount, notes,
		created_at, updated_at, version
	) VALUES (
		$1, $2::uuid, $3::uuid, nullif($4,'')::uuid, nullif($5,'')::uuid, $6,
		$7, $8, $9, $10, $11,
		$12, $13, $14, $15,
		$16, $17, $18
	) RETURNING row_to_json(t)`

	var vehicleIDStr, tripIDStr string
	if item.VehicleID != nil {
		vehicleIDStr = *item.VehicleID
	}
	if item.TripID != nil {
		tripIDStr = *item.TripID
	}

	row := exec.QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.DriverID, vehicleIDStr, tripIDStr, string(item.EventType),
		item.StartAt, item.EndAt, item.DurationMinutes, item.Location, item.OdometerKM,
		item.IsViolation, item.ViolationReason, item.DailyAllowanceAmount, item.Notes,
		item.CreatedAt, item.UpdatedAt, item.Version,
	)
	return readJourney(row)
}

func (r JourneyRepository) FindByID(ctx context.Context, tenantID, id string) (domain.DriverJourneyRecord, error) {
	query := `SELECT row_to_json(t) FROM driver_journey_records t WHERE tenant_id = $1::uuid AND id = $2::uuid`
	return readJourney(r.db.QueryRowContext(ctx, query, tenantID, id))
}

func (r JourneyRepository) ListByDriver(ctx context.Context, tenantID, driverID string, limit int) ([]domain.DriverJourneyRecord, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	query := `SELECT coalesce(json_agg(row_to_json(p)), '[]'::json) FROM (
		SELECT * FROM driver_journey_records
		WHERE tenant_id = $1::uuid AND driver_id = $2::uuid
		ORDER BY start_at DESC
		LIMIT $3
	) p`

	var body []byte
	if err := r.db.QueryRowContext(ctx, query, tenantID, driverID, limit).Scan(&body); err != nil {
		return nil, mapError(err)
	}
	var items []domain.DriverJourneyRecord
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (r JourneyRepository) GetLastActiveRecord(ctx context.Context, tenantID, driverID string) (*domain.DriverJourneyRecord, error) {
	query := `SELECT row_to_json(t) FROM driver_journey_records t
		WHERE tenant_id = $1::uuid AND driver_id = $2::uuid AND end_at IS NULL
		ORDER BY start_at DESC LIMIT 1`

	rec, err := readJourney(r.db.QueryRowContext(ctx, query, tenantID, driverID))
	if err != nil {
		if errors.Is(err, application.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rec, nil
}

func (r JourneyRepository) CloseActiveRecord(ctx context.Context, tenantID, id string, endAt time.Time, durationMinutes int, isViolation bool, violationReason string) error {
	exec := platform.ExecutorFor(ctx, r.db)
	query := `UPDATE driver_journey_records SET
		end_at = $3,
		duration_minutes = $4,
		is_violation = $5,
		violation_reason = $6,
		updated_at = $7,
		version = version + 1
	WHERE tenant_id = $1::uuid AND id = $2::uuid AND end_at IS NULL`

	res, err := exec.ExecContext(ctx, query, tenantID, id, endAt, durationMinutes, isViolation, violationReason, time.Now().UTC())
	if err != nil {
		return mapError(err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return application.ErrNotFound
	}
	return nil
}

func (r JourneyRepository) GetSummary(ctx context.Context, tenantID, driverID string) (domain.JourneySummary, error) {
	query := `SELECT
		coalesce(sum(duration_minutes) FILTER (WHERE event_type = 'driving'), 0),
		coalesce(sum(duration_minutes) FILTER (WHERE event_type IN ('break_rest', 'meal', 'daily_rest')), 0),
		coalesce(sum(duration_minutes) FILTER (WHERE event_type = 'waiting'), 0),
		coalesce(count(*) FILTER (WHERE is_violation = true), 0),
		coalesce(sum(daily_allowance_amount), 0)
	FROM driver_journey_records
	WHERE tenant_id = $1::uuid AND driver_id = $2::uuid`

	var s domain.JourneySummary
	err := r.db.QueryRowContext(ctx, query, tenantID, driverID).Scan(
		&s.TotalDrivingMinutes,
		&s.TotalRestMinutes,
		&s.TotalWaitingMinutes,
		&s.ViolationsCount,
		&s.TotalAllowanceBRL,
	)
	if err != nil {
		return s, mapError(err)
	}

	// Check if there is an active event currently running
	var activeType string
	err = r.db.QueryRowContext(ctx, `SELECT event_type FROM driver_journey_records WHERE tenant_id = $1::uuid AND driver_id = $2::uuid AND end_at IS NULL ORDER BY start_at DESC LIMIT 1`, tenantID, driverID).Scan(&activeType)
	if err == nil {
		s.ActiveEvent = &activeType
	}

	return s, nil
}


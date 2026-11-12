package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/trips/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/trips/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/trips/ports"
	platform "github.com/geevheit/intelligence360/backend/api/internal/platform/postgres"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func mapError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return application.ErrNotFound
	}
	var pg *pq.Error
	if errors.As(err, &pg) {
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

func (s *Store) Trips() ports.TripRepository {
	return TripRepository{db: s.db}
}

type TripRepository struct {
	db *sql.DB
}

func (r TripRepository) Create(ctx context.Context, item domain.Trip) (domain.Trip, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	query := `INSERT INTO operational_trips (
		id, tenant_id, trip_number, status, vehicle_id, trailer_id, driver_id,
		origin_city, origin_state, destination_city, destination_state, distance_km,
		departure_at, expected_arrival_at, delivered_at, cargo_description,
		cargo_weight_kg, cargo_value, customer_name, mdfe_key, cte_key, ciot_number,
		freight_amount, advance_amount, toll_amount, balance_amount,
		start_odometer_km, end_odometer_km, validation_notes, notes,
		created_at, updated_at, version
	) VALUES (
		$1, $2, $3, $4, $5::uuid, nullif($6,'')::uuid, $7::uuid,
		$8, $9, $10, $11, $12,
		$13, $14, $15, $16,
		$17, $18, $19, $20, $21, $22,
		$23, $24, $25, $26,
		$27, $28, $29, $30,
		$31, $32, $33
	) RETURNING row_to_json(operational_trips.*)`

	var trailerIDStr string
	if item.TrailerID != nil {
		trailerIDStr = *item.TrailerID
	}

	row := exec.QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.TripNumber, string(item.Status), item.VehicleID, trailerIDStr, item.DriverID,
		item.OriginCity, item.OriginState, item.DestinationCity, item.DestinationState, item.DistanceKM,
		item.DepartureAt, item.ExpectedArrivalAt, item.DeliveredAt, item.CargoDescription,
		item.CargoWeightKG, item.CargoValue, item.CustomerName, item.MDFEKey, item.CTEKey, item.CIOTNumber,
		item.FreightAmount, item.AdvanceAmount, item.TollAmount, item.BalanceAmount,
		item.StartOdometerKM, item.EndOdometerKM, item.ValidationNotes, item.Notes,
		item.CreatedAt, item.UpdatedAt, item.Version,
	)

	return read[domain.Trip](row)
}

func (r TripRepository) FindByID(ctx context.Context, tenantID string, id string) (domain.Trip, error) {
	query := `SELECT row_to_json(t) FROM operational_trips t WHERE tenant_id = $1::uuid AND id = $2::uuid AND deleted_at IS NULL`
	return read[domain.Trip](r.db.QueryRowContext(ctx, query, tenantID, id))
}

func (r TripRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Trip], error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 20
	}
	args := []any{tenantID}
	where := "tenant_id = $1::uuid AND deleted_at IS NULL"
	if q.Search != "" {
		args = append(args, q.Search)
		where += fmt.Sprintf(" AND (trip_number ILIKE '%%' || $%d || '%%' OR customer_name ILIKE '%%' || $%d || '%%' OR origin_city ILIKE '%%' || $%d || '%%' OR destination_city ILIKE '%%' || $%d || '%%')", len(args), len(args), len(args), len(args))
	}
	if status, ok := q.Filters["status"]; ok && status != "" {
		args = append(args, status)
		where += fmt.Sprintf(" AND status = $%d", len(args))
	}

	var count int
	if err := r.db.QueryRowContext(ctx, "SELECT count(*) FROM operational_trips WHERE "+where, args...).Scan(&count); err != nil {
		return ports.Page[domain.Trip]{}, mapError(err)
	}

	limit := q.PageSize
	offset := (q.Page - 1) * q.PageSize
	order := "created_at DESC"
	query := fmt.Sprintf("SELECT coalesce(json_agg(row_to_json(p)), '[]'::json) FROM (SELECT * FROM operational_trips WHERE %s ORDER BY %s LIMIT %d OFFSET %d) p", where, order, limit, offset)

	var body []byte
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&body); err != nil {
		return ports.Page[domain.Trip]{}, mapError(err)
	}
	items := []domain.Trip{}
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.Trip]{}, err
	}
	pages := (count + q.PageSize - 1) / q.PageSize
	return ports.Page[domain.Trip]{Items: items, Page: q.Page, PageSize: q.PageSize, TotalItems: count, TotalPages: pages}, nil
}

func (r TripRepository) Update(ctx context.Context, item domain.Trip) (domain.Trip, error) {
	exec := platform.ExecutorFor(ctx, r.db)
	query := `UPDATE operational_trips SET
		status = $3,
		trailer_id = nullif($4,'')::uuid,
		departure_at = $5,
		expected_arrival_at = $6,
		delivered_at = $7,
		cargo_description = $8,
		cargo_weight_kg = $9,
		cargo_value = $10,
		customer_name = $11,
		mdfe_key = $12,
		cte_key = $13,
		ciot_number = $14,
		freight_amount = $15,
		advance_amount = $16,
		toll_amount = $17,
		balance_amount = $18,
		start_odometer_km = $19,
		end_odometer_km = $20,
		validation_notes = $21,
		notes = $22,
		distance_km = $23,
		updated_at = $24,
		version = $25
	WHERE tenant_id = $1::uuid AND id = $2::uuid AND deleted_at IS NULL
	RETURNING row_to_json(operational_trips.*)`

	var trailerIDStr string
	if item.TrailerID != nil {
		trailerIDStr = *item.TrailerID
	}

	row := exec.QueryRowContext(ctx, query,
		item.TenantID, item.ID, string(item.Status), trailerIDStr,
		item.DepartureAt, item.ExpectedArrivalAt, item.DeliveredAt,
		item.CargoDescription, item.CargoWeightKG, item.CargoValue,
		item.CustomerName, item.MDFEKey, item.CTEKey, item.CIOTNumber,
		item.FreightAmount, item.AdvanceAmount, item.TollAmount, item.BalanceAmount,
		item.StartOdometerKM, item.EndOdometerKM, item.ValidationNotes, item.Notes,
		item.DistanceKM, item.UpdatedAt, item.Version,
	)

	return read[domain.Trip](row)
}

func (r TripRepository) Delete(ctx context.Context, tenantID string, id string) error {
	exec := platform.ExecutorFor(ctx, r.db)
	res, err := exec.ExecContext(ctx, `UPDATE operational_trips SET deleted_at = $3, version = version + 1 WHERE tenant_id = $1::uuid AND id = $2::uuid AND deleted_at IS NULL`, tenantID, id, time.Now().UTC())
	if err != nil {
		return mapError(err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return application.ErrNotFound
	}
	return nil
}

func (r TripRepository) GetSummary(ctx context.Context, tenantID string) (domain.TripSummary, error) {
	query := `SELECT
		count(*),
		coalesce(count(*) FILTER (WHERE status = 'dispatched'), 0),
		coalesce(count(*) FILTER (WHERE status = 'delivered'), 0),
		coalesce(sum(freight_amount), 0),
		coalesce(sum(advance_amount), 0),
		coalesce(sum(distance_km), 0)
	FROM operational_trips
	WHERE tenant_id = $1::uuid AND deleted_at IS NULL`

	var s domain.TripSummary
	err := r.db.QueryRowContext(ctx, query, tenantID).Scan(
		&s.TotalTrips,
		&s.DispatchedTrips,
		&s.DeliveredTrips,
		&s.TotalFreightBRL,
		&s.TotalAdvancesBRL,
		&s.TotalDistanceKM,
	)
	if err != nil {
		return s, mapError(err)
	}
	return s, nil
}


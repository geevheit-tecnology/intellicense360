package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/ports"
	platform "github.com/geevheit/intelligence360/backend/api/internal/platform/postgres"
)

type Store struct{ db *sql.DB }

func NewStore(db *sql.DB) *Store { return &Store{db: db} }

func mapError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return application.ErrCommandItemNotFound
	}
	var pg *pq.Error
	if errors.As(err, &pg) {
		if pg.Code == "23505" {
			return application.ErrDuplicateCommandItem
		}
		if pg.Code == "23503" || pg.Code == "23502" || pg.Code == "23514" || pg.Code == "22P02" {
			return application.ErrInvalidData
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

func (s *Store) Items() ports.CommandItemRepository {
	return ItemRepository{db: s.db}
}
func (s *Store) Events() ports.CommandEventRepository {
	return EventRepository{db: s.db}
}
func (s *Store) Actions() ports.CommandActionRepository {
	return ActionRepository{db: s.db}
}
func (s *Store) Snapshots() ports.OperationalSnapshotRepository {
	return SnapshotRepository{db: s.db}
}
func (s *Store) Idempotency() ports.IdempotencyRepository {
	return IdempotencyRepository{db: s.db}
}

// ==================== ItemRepository ====================

type ItemRepository struct{ db *sql.DB }

func (r ItemRepository) Create(ctx context.Context, item domain.CommandItem) (domain.CommandItem, error) {
	meta, err := json.Marshal(item.Metadata)
	if err != nil {
		meta = []byte("{}")
	}
	query := `INSERT INTO mission_control_command_items AS t (
		id, tenant_id, type, category, severity, priority, status,
		title, description, source, source_type, source_id,
		confidence, impact_score, risk_score, urgency_score,
		due_at, sla_status, sla_hours, detected_at, acknowledged_at, resolved_at,
		assigned_to, fingerprint, metadata, created_at, updated_at, deleted_at, version
	) VALUES (
		$1, $2::uuid, $3, $4, $5, $6, $7,
		$8, $9, $10, $11, $12,
		$13, $14, $15, $16,
		$17, $18, $19, $20, $21, $22,
		NULLIF($23, '')::uuid, $24, $25::jsonb, $26, $27, NULL, $28
	) RETURNING row_to_json(t)`

	if item.SLAStatus == "" {
		item.SLAStatus = domain.SLANotApplicable
	}

	return read[domain.CommandItem](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.Type, item.Category, item.Severity, item.Priority, item.Status,
		item.Title, item.Description, item.Source, item.SourceType, item.SourceID,
		item.Confidence, item.ImpactScore, item.RiskScore, item.UrgencyScore,
		item.DueAt, item.SLAStatus, item.SLAHours, item.DetectedAt, item.AcknowledgedAt, item.ResolvedAt,
		item.AssignedTo, item.Fingerprint, meta, item.CreatedAt, item.UpdatedAt, item.Version,
	))
}

func (r ItemRepository) GetByID(ctx context.Context, tenantID, id string) (domain.CommandItem, error) {
	query := `SELECT row_to_json(t) FROM mission_control_command_items AS t WHERE tenant_id = $1::uuid AND id = $2::uuid AND deleted_at IS NULL`
	return read[domain.CommandItem](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, id))
}

func (r ItemRepository) Update(ctx context.Context, item domain.CommandItem) (domain.CommandItem, error) {
	meta, err := json.Marshal(item.Metadata)
	if err != nil {
		meta = []byte("{}")
	}
	query := `UPDATE mission_control_command_items AS t SET
		type = $3,
		category = $4,
		severity = $5,
		priority = $6,
		status = $7,
		title = $8,
		description = $9,
		source = $10,
		source_type = $11,
		source_id = $12,
		confidence = $13,
		impact_score = $14,
		risk_score = $15,
		urgency_score = $16,
		due_at = $17,
		sla_status = $18,
		sla_hours = $19,
		acknowledged_at = $20,
		resolved_at = $21,
		assigned_to = NULLIF($22, '')::uuid,
		metadata = $23::jsonb,
		updated_at = $24,
		version = t.version + 1
	WHERE t.tenant_id = $1::uuid AND t.id = $2::uuid AND t.deleted_at IS NULL AND t.version = $25
	RETURNING row_to_json(t)`
	return read[domain.CommandItem](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.TenantID, item.ID, item.Type, item.Category, item.Severity, item.Priority, item.Status,
		item.Title, item.Description, item.Source, item.SourceType, item.SourceID,
		item.Confidence, item.ImpactScore, item.RiskScore, item.UrgencyScore,
		item.DueAt, item.SLAStatus, item.SLAHours, item.AcknowledgedAt, item.ResolvedAt,
		item.AssignedTo, meta, item.UpdatedAt, item.Version,
	))
}

func (r ItemRepository) List(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.CommandItem], error) {
	where := []string{"tenant_id = $1::uuid", "deleted_at IS NULL"}
	args := []any{tenantID}

	if q.Search != "" {
		args = append(args, "%"+strings.ToLower(q.Search)+"%")
		where = append(where, fmt.Sprintf("(lower(title) LIKE $%d OR lower(description) LIKE $%d)", len(args), len(args)))
	}
	if v, ok := q.Filters["status"]; ok && v != "" {
		args = append(args, v)
		where = append(where, fmt.Sprintf("status = $%d", len(args)))
	}
	if v, ok := q.Filters["severity"]; ok && v != "" {
		args = append(args, v)
		where = append(where, fmt.Sprintf("severity = $%d", len(args)))
	}
	if v, ok := q.Filters["category"]; ok && v != "" {
		args = append(args, v)
		where = append(where, fmt.Sprintf("category = $%d", len(args)))
	}
	if v, ok := q.Filters["type"]; ok && v != "" {
		args = append(args, v)
		where = append(where, fmt.Sprintf("type = $%d", len(args)))
	}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM mission_control_command_items WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.CommandItem]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	perPage := q.PerPage
	if perPage < 1 {
		perPage = 20
	}
	totalPages := (totalItems + perPage - 1) / perPage
	offset := (page - 1) * perPage

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM mission_control_command_items
			WHERE %s
			ORDER BY detected_at DESC
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, len(args)+1, len(args)+2)

	args = append(args, perPage, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.CommandItem]{}, mapError(err)
	}

	var items []domain.CommandItem
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.CommandItem]{}, err
	}

	return ports.Page[domain.CommandItem]{
		Data:       items,
		Page:       page,
		PerPage:    perPage,
		Total:      totalItems,
		TotalPages: totalPages,
	}, nil
}

func (r ItemRepository) FindByFingerprint(ctx context.Context, tenantID, fingerprint string) (domain.CommandItem, error) {
	query := `SELECT row_to_json(t) FROM mission_control_command_items AS t
		WHERE tenant_id = $1::uuid AND fingerprint = $2 AND status NOT IN ('resolved', 'dismissed') AND deleted_at IS NULL`
	return read[domain.CommandItem](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, fingerprint))
}

func (r ItemRepository) AllOpen(ctx context.Context, tenantID string) ([]domain.CommandItem, error) {
	query := `SELECT coalesce(json_agg(row_to_json(t)), '[]'::json) FROM (
		SELECT * FROM mission_control_command_items
		WHERE tenant_id = $1::uuid AND status NOT IN ('resolved', 'dismissed') AND deleted_at IS NULL
		ORDER BY detected_at DESC
	) AS t`
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID).Scan(&body); err != nil {
		return nil, mapError(err)
	}
	var items []domain.CommandItem
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, err
	}
	return items, nil
}

// ==================== EventRepository ====================

type EventRepository struct{ db *sql.DB }

func (r EventRepository) Create(ctx context.Context, event domain.CommandEvent) (domain.CommandEvent, error) {
	payload, _ := json.Marshal(event.Payload)
	query := `INSERT INTO mission_control_command_events AS t (
		id, tenant_id, command_item_id, event_type, previous_status, new_status, actor_id, payload, occurred_at
	) VALUES (
		$1, $2::uuid, $3::uuid, $4, $5, $6, NULLIF($7, '')::uuid, $8::jsonb, $9
	) RETURNING row_to_json(t)`
	return read[domain.CommandEvent](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		event.ID, event.TenantID, event.CommandItemID, event.EventType, event.PreviousStatus, event.NewStatus,
		event.ActorID, payload, event.OccurredAt,
	))
}

func (r EventRepository) ListByItem(ctx context.Context, tenantID, itemID string, q ports.Query) (ports.Page[domain.CommandEvent], error) {
	where := "tenant_id = $1::uuid AND command_item_id = $2::uuid"
	args := []any{tenantID, itemID}

	countQuery := fmt.Sprintf("SELECT count(*) FROM mission_control_command_events WHERE %s", where)
	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.CommandEvent]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	perPage := q.PerPage
	if perPage < 1 {
		perPage = 50
	}
	totalPages := (totalItems + perPage - 1) / perPage
	offset := (page - 1) * perPage

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM mission_control_command_events
			WHERE %s
			ORDER BY occurred_at DESC
			LIMIT $%d OFFSET $%d
		) AS t`, where, len(args)+1, len(args)+2)

	args = append(args, perPage, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.CommandEvent]{}, mapError(err)
	}

	var items []domain.CommandEvent
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.CommandEvent]{}, err
	}

	return ports.Page[domain.CommandEvent]{
		Data:       items,
		Page:       page,
		PerPage:    perPage,
		Total:      totalItems,
		TotalPages: totalPages,
	}, nil
}

// ==================== ActionRepository ====================

type ActionRepository struct{ db *sql.DB }

func (r ActionRepository) Create(ctx context.Context, action domain.CommandAction) (domain.CommandAction, error) {
	payload, _ := json.Marshal(action.Payload)
	query := `INSERT INTO mission_control_command_actions AS t (
		id, tenant_id, command_item_id, type, label, status, priority, payload, created_at, updated_at
	) VALUES (
		$1, $2::uuid, $3::uuid, $4, $5, $6, $7, $8::jsonb, $9, $10
	) RETURNING row_to_json(t)`
	return read[domain.CommandAction](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		action.ID, action.TenantID, action.CommandItemID, action.Type, action.Label,
		action.Status, action.Priority, payload, action.CreatedAt, action.UpdatedAt,
	))
}

func (r ActionRepository) ListByItem(ctx context.Context, tenantID, itemID string, q ports.Query) (ports.Page[domain.CommandAction], error) {
	where := "tenant_id = $1::uuid AND command_item_id = $2::uuid"
	args := []any{tenantID, itemID}

	countQuery := fmt.Sprintf("SELECT count(*) FROM mission_control_command_actions WHERE %s", where)
	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.CommandAction]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	perPage := q.PerPage
	if perPage < 1 {
		perPage = 50
	}
	totalPages := (totalItems + perPage - 1) / perPage
	offset := (page - 1) * perPage

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM mission_control_command_actions
			WHERE %s
			ORDER BY created_at ASC
			LIMIT $%d OFFSET $%d
		) AS t`, where, len(args)+1, len(args)+2)

	args = append(args, perPage, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.CommandAction]{}, mapError(err)
	}

	var items []domain.CommandAction
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.CommandAction]{}, err
	}

	return ports.Page[domain.CommandAction]{
		Data:       items,
		Page:       page,
		PerPage:    perPage,
		Total:      totalItems,
		TotalPages: totalPages,
	}, nil
}

// ==================== SnapshotRepository ====================

type SnapshotRepository struct{ db *sql.DB }

func (r SnapshotRepository) Create(ctx context.Context, snapshot domain.OperationalSnapshot) (domain.OperationalSnapshot, error) {
	query := `INSERT INTO mission_control_snapshots AS t (
		id, tenant_id, snapshot_at, open_items, critical_items, high_priority_items,
		active_risks, active_alerts, open_incidents, opportunities, breached_slas,
		average_resolution_time, operational_score, risk_score, health_score
	) VALUES (
		$1, $2::uuid, $3, $4, $5, $6,
		$7, $8, $9, $10, $11,
		$12, $13, $14, $15
	) RETURNING row_to_json(t)`
	return read[domain.OperationalSnapshot](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		snapshot.ID, snapshot.TenantID, snapshot.SnapshotAt, snapshot.OpenItems, snapshot.CriticalItems, snapshot.HighPriorityItems,
		snapshot.ActiveRisks, snapshot.ActiveAlerts, snapshot.OpenIncidents, snapshot.Opportunities, snapshot.BreachedSLAs,
		snapshot.AverageResolutionTime, snapshot.OperationalScore, snapshot.RiskScore, snapshot.HealthScore,
	))
}

func (r SnapshotRepository) Latest(ctx context.Context, tenantID string) (domain.OperationalSnapshot, error) {
	query := `SELECT row_to_json(t) FROM mission_control_snapshots AS t
		WHERE tenant_id = $1::uuid
		ORDER BY snapshot_at DESC
		LIMIT 1`
	return read[domain.OperationalSnapshot](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID))
}

// ==================== IdempotencyRepository ====================

type IdempotencyRepository struct{ db *sql.DB }

func (r IdempotencyRepository) Exists(ctx context.Context, tenantID, key string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM mission_control_idempotency WHERE tenant_id = $1::uuid AND idempotency_key = $2)`
	var exists bool
	err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, key).Scan(&exists)
	return exists, mapError(err)
}

func (r IdempotencyRepository) Save(ctx context.Context, tenantID, key, resourceID string) error {
	query := `INSERT INTO mission_control_idempotency (tenant_id, idempotency_key, resource_id) VALUES ($1::uuid, $2, $3::uuid) ON CONFLICT DO NOTHING`
	_, err := platform.ExecutorFor(ctx, r.db).ExecContext(ctx, query, tenantID, key, resourceID)
	return mapError(err)
}

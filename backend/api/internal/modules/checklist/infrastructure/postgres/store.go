package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/ports"
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

func (s *Store) Checklists() ports.ChecklistRepository {
	return ChecklistRepository{db: s.db}
}
func (s *Store) Items() ports.ChecklistItemRepository {
	return ItemRepository{db: s.db}
}
func (s *Store) Answers() ports.ChecklistAnswerRepository {
	return AnswerRepository{db: s.db}
}
func (s *Store) Templates() ports.ChecklistTemplateRepository {
	return TemplateRepository{db: s.db}
}
func (s *Store) Versions() ports.ChecklistTemplateVersionRepository {
	return TemplateVersionRepository{db: s.db}
}
func (s *Store) Types() ports.ChecklistTypeRepository {
	return TypeRepository{db: s.db}
}
func (s *Store) Sections() ports.ChecklistSectionRepository {
	return SectionRepository{db: s.db}
}
func (s *Store) EngineItems() ports.ChecklistEngineItemRepository {
	return EngineItemRepository{db: s.db}
}
func (s *Store) Executions() ports.ChecklistExecutionRepository {
	return ExecutionRepository{db: s.db}
}
func (s *Store) Responses() ports.ChecklistResponseRepository {
	return ResponseRepository{db: s.db}
}
func (s *Store) Evidence() ports.ChecklistEvidenceRepository {
	return EvidenceRepository{db: s.db}
}
func (s *Store) NonConformities() ports.ChecklistNonConformityRepository {
	return NonConformityRepository{db: s.db}
}
func (s *Store) Signatures() ports.ChecklistSignatureRepository {
	return SignatureRepository{db: s.db}
}
func (s *Store) History() ports.ChecklistHistoryRepository {
	return HistoryRepository{db: s.db}
}

// ==================== ChecklistRepository ====================

type ChecklistRepository struct{ db *sql.DB }

func (r ChecklistRepository) Create(ctx context.Context, item domain.Checklist) (domain.Checklist, error) {
	query := `INSERT INTO checklists AS t (
		id, tenant_id, vehicle_id, title, description, type, status,
		started_at, finished_at, driver_name, driver_document,
		created_by, updated_by, deleted_at, created_at, updated_at
	) VALUES (
		$1, $2::uuid, $3::uuid, $4, $5, $6, $7,
		$8, $9, $10, $11,
		NULLIF($12, '')::uuid, NULLIF($13, '')::uuid, $14, $15, $16
	) RETURNING row_to_json(t)`
	return read[domain.Checklist](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.VehicleID, item.Title, item.Description, item.Type, item.Status,
		item.StartedAt, item.FinishedAt, item.DriverName, item.DriverDocument,
		item.CreatedBy, item.UpdatedBy, item.DeletedAt, item.CreatedAt, item.UpdatedAt,
	))
}

func (r ChecklistRepository) FindByID(ctx context.Context, tenantID string, id string) (domain.Checklist, error) {
	query := `SELECT row_to_json(t) FROM checklists AS t WHERE tenant_id = $1::uuid AND id = $2::uuid AND deleted_at IS NULL`
	return read[domain.Checklist](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, id))
}

func (r ChecklistRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.Checklist], error) {
	where := []string{"tenant_id = $1::uuid", "deleted_at IS NULL"}
	args := []any{tenantID}

	if q.Search != "" {
		args = append(args, "%"+strings.ToLower(q.Search)+"%")
		where = append(where, fmt.Sprintf("(lower(title) LIKE $%d OR lower(driver_name) LIKE $%d)", len(args), len(args)))
	}
	if v, ok := q.Filters["vehicle_id"]; ok && v != "" {
		args = append(args, v)
		where = append(where, fmt.Sprintf("vehicle_id = $%d::uuid", len(args)))
	}
	if v, ok := q.Filters["status"]; ok && v != "" {
		args = append(args, v)
		where = append(where, fmt.Sprintf("status = $%d", len(args)))
	}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM checklists WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.Checklist]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM checklists
			WHERE %s
			ORDER BY created_at DESC
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.Checklist]{}, mapError(err)
	}

	var items []domain.Checklist
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.Checklist]{}, err
	}

	return ports.Page[domain.Checklist]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

func (r ChecklistRepository) Update(ctx context.Context, item domain.Checklist) (domain.Checklist, error) {
	query := `UPDATE checklists AS t SET
		vehicle_id = $3::uuid,
		title = $4,
		description = $5,
		type = $6,
		status = $7,
		started_at = $8,
		finished_at = $9,
		driver_name = $10,
		driver_document = $11,
		updated_by = NULLIF($12, '')::uuid,
		updated_at = $13
	WHERE t.tenant_id = $1::uuid AND t.id = $2::uuid AND t.deleted_at IS NULL
	RETURNING row_to_json(t)`
	return read[domain.Checklist](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.TenantID, item.ID, item.VehicleID, item.Title, item.Description, item.Type, item.Status,
		item.StartedAt, item.FinishedAt, item.DriverName, item.DriverDocument,
		item.UpdatedBy, item.UpdatedAt,
	))
}

func (r ChecklistRepository) Delete(ctx context.Context, tenantID string, id string) error {
	query := `UPDATE checklists SET deleted_at = now() WHERE tenant_id = $1::uuid AND id = $2::uuid AND deleted_at IS NULL`
	res, err := platform.ExecutorFor(ctx, r.db).ExecContext(ctx, query, tenantID, id)
	return changed(res, err)
}

func (r ChecklistRepository) Exists(ctx context.Context, tenantID string, id string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM checklists WHERE tenant_id = $1::uuid AND id = $2::uuid AND deleted_at IS NULL)`
	var exists bool
	err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, id).Scan(&exists)
	return exists, mapError(err)
}

// ==================== ItemRepository ====================

type ItemRepository struct{ db *sql.DB }

func (r ItemRepository) Create(ctx context.Context, item domain.ChecklistItem) (domain.ChecklistItem, error) {
	query := `INSERT INTO checklist_items AS t (
		id, tenant_id, checklist_id, title, description, category,
		required, order_index, answer_type, expected_value, deleted_at, created_at, updated_at
	) VALUES (
		$1, $2::uuid, $3::uuid, $4, $5, $6,
		$7, $8, $9, $10, $11, $12, $13
	) RETURNING row_to_json(t)`
	return read[domain.ChecklistItem](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.ChecklistID, item.Title, item.Description, item.Category,
		item.Required, item.OrderIndex, item.AnswerType, item.ExpectedValue, item.DeletedAt, item.CreatedAt, item.UpdatedAt,
	))
}

func (r ItemRepository) FindByID(ctx context.Context, tenantID string, checklistID string, itemID string) (domain.ChecklistItem, error) {
	query := `SELECT row_to_json(t) FROM checklist_items AS t WHERE tenant_id = $1::uuid AND checklist_id = $2::uuid AND id = $3::uuid AND deleted_at IS NULL`
	return read[domain.ChecklistItem](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, checklistID, itemID))
}

func (r ItemRepository) List(ctx context.Context, tenantID string, checklistID string, q ports.Query) (ports.Page[domain.ChecklistItem], error) {
	where := []string{"tenant_id = $1::uuid", "checklist_id = $2::uuid", "deleted_at IS NULL"}
	args := []any{tenantID, checklistID}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM checklist_items WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.ChecklistItem]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM checklist_items
			WHERE %s
			ORDER BY order_index ASC
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.ChecklistItem]{}, mapError(err)
	}

	var items []domain.ChecklistItem
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.ChecklistItem]{}, err
	}

	return ports.Page[domain.ChecklistItem]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

func (r ItemRepository) Update(ctx context.Context, item domain.ChecklistItem) (domain.ChecklistItem, error) {
	query := `UPDATE checklist_items AS t SET
		title = $4,
		description = $5,
		category = $6,
		required = $7,
		order_index = $8,
		answer_type = $9,
		expected_value = $10,
		updated_at = $11
	WHERE t.tenant_id = $1::uuid AND t.checklist_id = $2::uuid AND t.id = $3::uuid AND t.deleted_at IS NULL
	RETURNING row_to_json(t)`
	return read[domain.ChecklistItem](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.TenantID, item.ChecklistID, item.ID, item.Title, item.Description, item.Category,
		item.Required, item.OrderIndex, item.AnswerType, item.ExpectedValue, item.UpdatedAt,
	))
}

func (r ItemRepository) Delete(ctx context.Context, tenantID string, checklistID string, itemID string) error {
	query := `UPDATE checklist_items SET deleted_at = now() WHERE tenant_id = $1::uuid AND checklist_id = $2::uuid AND id = $3::uuid AND deleted_at IS NULL`
	res, err := platform.ExecutorFor(ctx, r.db).ExecContext(ctx, query, tenantID, checklistID, itemID)
	return changed(res, err)
}

// ==================== AnswerRepository ====================

type AnswerRepository struct{ db *sql.DB }

func (r AnswerRepository) Create(ctx context.Context, item domain.ChecklistAnswer) (domain.ChecklistAnswer, error) {
	query := `INSERT INTO checklist_answers AS t (
		id, tenant_id, checklist_id, checklist_item_id, answer, notes, photo_url,
		answered_by, answered_at, deleted_at, created_at, updated_at
	) VALUES (
		$1, $2::uuid, $3::uuid, $4::uuid, $5, $6, $7,
		NULLIF($8, '')::uuid, $9, NULL, $10, $11
	) RETURNING row_to_json(t)`
	return read[domain.ChecklistAnswer](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.ChecklistID, item.ChecklistItemID, item.Answer, item.Notes, item.PhotoURL,
		item.AnsweredBy, item.AnsweredAt, item.CreatedAt, item.UpdatedAt,
	))
}

func (r AnswerRepository) List(ctx context.Context, tenantID string, checklistID string, q ports.Query) (ports.Page[domain.ChecklistAnswer], error) {
	where := []string{"tenant_id = $1::uuid", "checklist_id = $2::uuid", "deleted_at IS NULL"}
	args := []any{tenantID, checklistID}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM checklist_answers WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.ChecklistAnswer]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM checklist_answers
			WHERE %s
			ORDER BY answered_at ASC
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.ChecklistAnswer]{}, mapError(err)
	}

	var items []domain.ChecklistAnswer
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.ChecklistAnswer]{}, err
	}

	return ports.Page[domain.ChecklistAnswer]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

// ==================== TemplateRepository ====================

type TemplateRepository struct{ db *sql.DB }

func (r TemplateRepository) Create(ctx context.Context, item domain.ChecklistTemplate) (domain.ChecklistTemplate, error) {
	query := `INSERT INTO checklist_templates AS t (
		id, tenant_id, name, description, type_id, type, status, active,
		current_version_number, audit_id, version, created_at, updated_at, deleted_at
	) VALUES (
		$1, $2::uuid, $3, $4, NULLIF($5, '')::uuid, $6, $7, $8,
		$9, NULLIF($10, '')::uuid, $11, $12, $13, $14
	) RETURNING row_to_json(t)`
	return read[domain.ChecklistTemplate](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.Name, item.Description, item.TypeID, item.Type, item.Status, item.Active,
		item.CurrentVersionNumber, item.CreatedBy, item.Version, item.CreatedAt, item.UpdatedAt, item.DeletedAt,
	))
}

func (r TemplateRepository) FindByID(ctx context.Context, tenantID string, id string) (domain.ChecklistTemplate, error) {
	query := `SELECT row_to_json(t) FROM checklist_templates AS t WHERE tenant_id = $1::uuid AND id = $2::uuid AND deleted_at IS NULL`
	return read[domain.ChecklistTemplate](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, id))
}

func (r TemplateRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.ChecklistTemplate], error) {
	where := []string{"tenant_id = $1::uuid", "deleted_at IS NULL"}
	args := []any{tenantID}

	if q.Search != "" {
		args = append(args, "%"+strings.ToLower(q.Search)+"%")
		where = append(where, fmt.Sprintf("(lower(name) LIKE $%d OR lower(description) LIKE $%d)", len(args), len(args)))
	}
	if v, ok := q.Filters["status"]; ok && v != "" {
		args = append(args, v)
		where = append(where, fmt.Sprintf("status = $%d", len(args)))
	}
	if v, ok := q.Filters["active"]; ok && v != "" {
		args = append(args, v == "true")
		where = append(where, fmt.Sprintf("active = $%d", len(args)))
	}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM checklist_templates WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.ChecklistTemplate]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM checklist_templates
			WHERE %s
			ORDER BY created_at DESC
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.ChecklistTemplate]{}, mapError(err)
	}

	var items []domain.ChecklistTemplate
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.ChecklistTemplate]{}, err
	}

	return ports.Page[domain.ChecklistTemplate]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

func (r TemplateRepository) Update(ctx context.Context, item domain.ChecklistTemplate) (domain.ChecklistTemplate, error) {
	query := `UPDATE checklist_templates AS t SET
		name = $3,
		description = $4,
		type_id = NULLIF($5, '')::uuid,
		type = $6,
		status = $7,
		active = $8,
		current_version_number = $9,
		updated_at = $10,
		version = t.version + 1
	WHERE t.tenant_id = $1::uuid AND t.id = $2::uuid AND t.deleted_at IS NULL AND t.version = $11
	RETURNING row_to_json(t)`
	return read[domain.ChecklistTemplate](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.TenantID, item.ID, item.Name, item.Description, item.TypeID, item.Type, item.Status, item.Active,
		item.CurrentVersionNumber, item.UpdatedAt, item.Version,
	))
}

func (r TemplateRepository) Delete(ctx context.Context, tenantID string, id string) error {
	query := `UPDATE checklist_templates SET deleted_at = now() WHERE tenant_id = $1::uuid AND id = $2::uuid AND deleted_at IS NULL`
	res, err := platform.ExecutorFor(ctx, r.db).ExecContext(ctx, query, tenantID, id)
	return changed(res, err)
}

// ==================== TemplateVersionRepository ====================

type TemplateVersionRepository struct{ db *sql.DB }

func (r TemplateVersionRepository) Create(ctx context.Context, item domain.ChecklistTemplateVersion) (domain.ChecklistTemplateVersion, error) {
	scoring, _ := json.Marshal(item.ScoringConfig)
	severity, _ := json.Marshal(item.SeverityConfig)
	query := `INSERT INTO checklist_template_versions AS t (
		id, tenant_id, template_id, version_number, status, instructions,
		scoring_config, severity_config, evidence_required, signature_required,
		published_at, audit_id, version, created_at, updated_at, deleted_at
	) VALUES (
		$1, $2::uuid, $3::uuid, $4, $5, $6,
		$7::jsonb, $8::jsonb, $9, $10,
		$11, NULLIF($12, '')::uuid, $13, $14, $15, $16
	) RETURNING row_to_json(t)`
	return read[domain.ChecklistTemplateVersion](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.TemplateID, item.VersionNumber, item.Status, item.Instructions,
		scoring, severity, item.EvidenceRequired, item.SignatureRequired,
		item.PublishedAt, "", item.Version, item.CreatedAt, item.UpdatedAt, item.DeletedAt,
	))
}

func (r TemplateVersionRepository) FindByID(ctx context.Context, tenantID string, id string) (domain.ChecklistTemplateVersion, error) {
	query := `SELECT row_to_json(t) FROM checklist_template_versions AS t WHERE tenant_id = $1::uuid AND id = $2::uuid AND deleted_at IS NULL`
	return read[domain.ChecklistTemplateVersion](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, id))
}

func (r TemplateVersionRepository) ListByTemplate(ctx context.Context, tenantID string, templateID string, q ports.Query) (ports.Page[domain.ChecklistTemplateVersion], error) {
	where := []string{"tenant_id = $1::uuid", "template_id = $2::uuid", "deleted_at IS NULL"}
	args := []any{tenantID, templateID}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM checklist_template_versions WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.ChecklistTemplateVersion]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM checklist_template_versions
			WHERE %s
			ORDER BY version_number DESC
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.ChecklistTemplateVersion]{}, mapError(err)
	}

	var items []domain.ChecklistTemplateVersion
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.ChecklistTemplateVersion]{}, err
	}

	return ports.Page[domain.ChecklistTemplateVersion]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

func (r TemplateVersionRepository) Update(ctx context.Context, item domain.ChecklistTemplateVersion) (domain.ChecklistTemplateVersion, error) {
	scoring, _ := json.Marshal(item.ScoringConfig)
	severity, _ := json.Marshal(item.SeverityConfig)
	query := `UPDATE checklist_template_versions AS t SET
		status = $3,
		instructions = $4,
		scoring_config = $5::jsonb,
		severity_config = $6::jsonb,
		evidence_required = $7,
		signature_required = $8,
		published_at = $9,
		updated_at = $10,
		version = t.version + 1
	WHERE t.tenant_id = $1::uuid AND t.id = $2::uuid AND t.deleted_at IS NULL AND t.version = $11
	RETURNING row_to_json(t)`
	return read[domain.ChecklistTemplateVersion](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.TenantID, item.ID, item.Status, item.Instructions,
		scoring, severity, item.EvidenceRequired, item.SignatureRequired,
		item.PublishedAt, item.UpdatedAt, item.Version,
	))
}

// ==================== TypeRepository ====================

type TypeRepository struct{ db *sql.DB }

func (r TypeRepository) Create(ctx context.Context, item domain.ChecklistType) (domain.ChecklistType, error) {
	query := `INSERT INTO checklist_types AS t (
		id, tenant_id, name, code, description, active, audit_id, version, created_at, updated_at, deleted_at
	) VALUES (
		$1, $2::uuid, $3, $4, $5, $6, NULLIF($7, '')::uuid, $8, $9, $10, $11
	) RETURNING row_to_json(t)`
	return read[domain.ChecklistType](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.Name, item.Code, item.Description, item.Active,
		"", item.Version, item.CreatedAt, item.UpdatedAt, item.DeletedAt,
	))
}

func (r TypeRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.ChecklistType], error) {
	where := []string{"tenant_id = $1::uuid", "deleted_at IS NULL"}
	args := []any{tenantID}

	if q.Search != "" {
		args = append(args, "%"+strings.ToLower(q.Search)+"%")
		where = append(where, fmt.Sprintf("(lower(name) LIKE $%d OR lower(code) LIKE $%d)", len(args), len(args)))
	}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM checklist_types WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.ChecklistType]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM checklist_types
			WHERE %s
			ORDER BY name ASC
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.ChecklistType]{}, mapError(err)
	}

	var items []domain.ChecklistType
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.ChecklistType]{}, err
	}

	return ports.Page[domain.ChecklistType]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

// ==================== SectionRepository ====================

type SectionRepository struct{ db *sql.DB }

func (r SectionRepository) Create(ctx context.Context, item domain.ChecklistSection) (domain.ChecklistSection, error) {
	query := `INSERT INTO checklist_sections AS t (
		id, tenant_id, template_version_id, name, description, order_index,
		active, audit_id, version, created_at, updated_at, deleted_at
	) VALUES (
		$1, $2::uuid, $3::uuid, $4, $5, $6,
		$7, NULLIF($8, '')::uuid, $9, $10, $11, $12
	) RETURNING row_to_json(t)`
	return read[domain.ChecklistSection](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.TemplateVersionID, item.Name, item.Description, item.OrderIndex,
		item.Active, "", item.Version, item.CreatedAt, item.UpdatedAt, item.DeletedAt,
	))
}

func (r SectionRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.ChecklistSection], error) {
	where := []string{"tenant_id = $1::uuid", "deleted_at IS NULL"}
	args := []any{tenantID}

	if v, ok := q.Filters["template_version_id"]; ok && v != "" {
		args = append(args, v)
		where = append(where, fmt.Sprintf("template_version_id = $%d::uuid", len(args)))
	}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM checklist_sections WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.ChecklistSection]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM checklist_sections
			WHERE %s
			ORDER BY order_index ASC
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.ChecklistSection]{}, mapError(err)
	}

	var items []domain.ChecklistSection
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.ChecklistSection]{}, err
	}

	return ports.Page[domain.ChecklistSection]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

// ==================== EngineItemRepository ====================

type EngineItemRepository struct{ db *sql.DB }

func (r EngineItemRepository) Create(ctx context.Context, item domain.ChecklistEngineItem) (domain.ChecklistEngineItem, error) {
	query := `INSERT INTO checklist_engine_items AS t (
		id, tenant_id, template_version_id, section_id, question, description,
		item_type, required, order_index, help_text, severity, evidence_required,
		active, audit_id, version, created_at, updated_at, deleted_at
	) VALUES (
		$1, $2::uuid, $3::uuid, NULLIF($4, '')::uuid, $5, $6,
		$7, $8, $9, $10, $11, $12,
		$13, NULLIF($14, '')::uuid, $15, $16, $17, $18
	) RETURNING row_to_json(t)`
	return read[domain.ChecklistEngineItem](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.TemplateVersionID, item.SectionID, item.Question, item.Description,
		item.ItemType, item.Required, item.OrderIndex, item.HelpText, item.Severity, item.EvidenceRequired,
		item.Active, "", item.Version, item.CreatedAt, item.UpdatedAt, item.DeletedAt,
	))
}

func (r EngineItemRepository) FindByID(ctx context.Context, tenantID string, id string) (domain.ChecklistEngineItem, error) {
	query := `SELECT row_to_json(t) FROM checklist_engine_items AS t WHERE tenant_id = $1::uuid AND id = $2::uuid AND deleted_at IS NULL`
	return read[domain.ChecklistEngineItem](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, id))
}

func (r EngineItemRepository) ListByVersion(ctx context.Context, tenantID string, versionID string, q ports.Query) (ports.Page[domain.ChecklistEngineItem], error) {
	where := []string{"tenant_id = $1::uuid", "template_version_id = $2::uuid", "deleted_at IS NULL"}
	args := []any{tenantID, versionID}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM checklist_engine_items WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.ChecklistEngineItem]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM checklist_engine_items
			WHERE %s
			ORDER BY order_index ASC
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.ChecklistEngineItem]{}, mapError(err)
	}

	var items []domain.ChecklistEngineItem
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.ChecklistEngineItem]{}, err
	}

	return ports.Page[domain.ChecklistEngineItem]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

// ==================== ExecutionRepository ====================

type ExecutionRepository struct{ db *sql.DB }

func (r ExecutionRepository) Create(ctx context.Context, item domain.ChecklistExecution) (domain.ChecklistExecution, error) {
	query := `INSERT INTO checklist_executions AS t (
		id, tenant_id, template_version_id, vehicle_id, status, performed_by,
		location_reference, notes, score, final_result, started_at, ended_at,
		audit_id, version, created_at, updated_at, deleted_at
	) VALUES (
		$1, $2::uuid, $3::uuid, NULLIF($4, '')::uuid, $5, $6,
		$7, $8, $9, $10, $11, $12,
		NULLIF($13, '')::uuid, $14, $15, $16, $17
	) RETURNING row_to_json(t)`
	return read[domain.ChecklistExecution](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.TemplateVersionID, item.VehicleID, item.Status, item.PerformedBy,
		item.LocationReference, item.Notes, item.Score, item.FinalResult, item.StartedAt, item.EndedAt,
		item.PerformedBy, item.Version, item.CreatedAt, item.UpdatedAt, item.DeletedAt,
	))
}

func (r ExecutionRepository) FindByID(ctx context.Context, tenantID string, id string) (domain.ChecklistExecution, error) {
	query := `SELECT row_to_json(t) FROM checklist_executions AS t WHERE tenant_id = $1::uuid AND id = $2::uuid AND deleted_at IS NULL`
	return read[domain.ChecklistExecution](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, id))
}

func (r ExecutionRepository) Search(ctx context.Context, tenantID string, q ports.Query) (ports.Page[domain.ChecklistExecution], error) {
	where := []string{"tenant_id = $1::uuid", "deleted_at IS NULL"}
	args := []any{tenantID}

	if v, ok := q.Filters["template_version_id"]; ok && v != "" {
		args = append(args, v)
		where = append(where, fmt.Sprintf("template_version_id = $%d::uuid", len(args)))
	}
	if v, ok := q.Filters["vehicle_id"]; ok && v != "" {
		args = append(args, v)
		where = append(where, fmt.Sprintf("vehicle_id = $%d::uuid", len(args)))
	}
	if v, ok := q.Filters["status"]; ok && v != "" {
		args = append(args, v)
		where = append(where, fmt.Sprintf("status = $%d", len(args)))
	}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM checklist_executions WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.ChecklistExecution]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM checklist_executions
			WHERE %s
			ORDER BY created_at DESC
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.ChecklistExecution]{}, mapError(err)
	}

	var items []domain.ChecklistExecution
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.ChecklistExecution]{}, err
	}

	return ports.Page[domain.ChecklistExecution]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

func (r ExecutionRepository) Update(ctx context.Context, item domain.ChecklistExecution) (domain.ChecklistExecution, error) {
	query := `UPDATE checklist_executions AS t SET
		status = $3,
		notes = $4,
		score = $5,
		final_result = $6,
		started_at = $7,
		ended_at = $8,
		updated_at = $9,
		version = t.version + 1
	WHERE t.tenant_id = $1::uuid AND t.id = $2::uuid AND t.deleted_at IS NULL AND t.version = $10
	RETURNING row_to_json(t)`
	return read[domain.ChecklistExecution](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.TenantID, item.ID, item.Status, item.Notes, item.Score, item.FinalResult,
		item.StartedAt, item.EndedAt, item.UpdatedAt, item.Version,
	))
}

// ==================== ResponseRepository ====================

type ResponseRepository struct{ db *sql.DB }

func (r ResponseRepository) Create(ctx context.Context, item domain.ChecklistResponse) (domain.ChecklistResponse, error) {
	query := `INSERT INTO checklist_responses AS t (
		id, tenant_id, execution_id, item_id, value, result, notes,
		responder, severity, responded_at, audit_id, version, created_at, updated_at, deleted_at
	) VALUES (
		$1, $2::uuid, $3::uuid, $4::uuid, $5, $6, $7,
		$8, $9, $10, NULLIF($11, '')::uuid, 1, $12, $12, NULL
	) RETURNING row_to_json(t)`
	return read[domain.ChecklistResponse](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.ExecutionID, item.ItemID, item.Value, item.Result, item.Notes,
		item.Responder, item.Severity, item.RespondedAt, item.Responder, item.CreatedAt,
	))
}

func (r ResponseRepository) ListByExecution(ctx context.Context, tenantID string, executionID string, q ports.Query) (ports.Page[domain.ChecklistResponse], error) {
	where := []string{"tenant_id = $1::uuid", "execution_id = $2::uuid", "deleted_at IS NULL"}
	args := []any{tenantID, executionID}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM checklist_responses WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.ChecklistResponse]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 100
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM checklist_responses
			WHERE %s
			ORDER BY responded_at ASC
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.ChecklistResponse]{}, mapError(err)
	}

	var items []domain.ChecklistResponse
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.ChecklistResponse]{}, err
	}

	return ports.Page[domain.ChecklistResponse]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

func (r ResponseRepository) ExistsForItem(ctx context.Context, tenantID string, executionID string, itemID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM checklist_responses WHERE tenant_id = $1::uuid AND execution_id = $2::uuid AND item_id = $3::uuid AND deleted_at IS NULL)`
	var exists bool
	err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, executionID, itemID).Scan(&exists)
	return exists, mapError(err)
}

// ==================== EvidenceRepository ====================

type EvidenceRepository struct{ db *sql.DB }

func (r EvidenceRepository) Create(ctx context.Context, item domain.ChecklistEvidence) (domain.ChecklistEvidence, error) {
	query := `INSERT INTO checklist_evidence AS t (
		id, tenant_id, execution_id, response_id, evidence_type, reference,
		file_name, mime_type, notes, audit_id, version, created_at, updated_at, deleted_at
	) VALUES (
		$1, $2::uuid, $3::uuid, $4, $5, $6,
		$7, $8, $9, NULL, 1, $10, $10, NULL
	) RETURNING row_to_json(t)`
	return read[domain.ChecklistEvidence](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.ExecutionID, item.ResponseID, item.EvidenceType, item.Reference,
		item.FileName, item.MimeType, item.Notes, item.CreatedAt,
	))
}

func (r EvidenceRepository) ListByExecution(ctx context.Context, tenantID string, executionID string, q ports.Query) (ports.Page[domain.ChecklistEvidence], error) {
	where := []string{"tenant_id = $1::uuid", "execution_id = $2::uuid", "deleted_at IS NULL"}
	args := []any{tenantID, executionID}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM checklist_evidence WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.ChecklistEvidence]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM checklist_evidence
			WHERE %s
			ORDER BY created_at ASC
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.ChecklistEvidence]{}, mapError(err)
	}

	var items []domain.ChecklistEvidence
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.ChecklistEvidence]{}, err
	}

	return ports.Page[domain.ChecklistEvidence]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

func (r EvidenceRepository) ExistsForResponse(ctx context.Context, tenantID string, executionID string, responseID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM checklist_evidence WHERE tenant_id = $1::uuid AND execution_id = $2::uuid AND response_id = $3 AND deleted_at IS NULL)`
	var exists bool
	err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, executionID, responseID).Scan(&exists)
	return exists, mapError(err)
}

// ==================== NonConformityRepository ====================

type NonConformityRepository struct{ db *sql.DB }

func (r NonConformityRepository) Create(ctx context.Context, item domain.ChecklistNonConformity) (domain.ChecklistNonConformity, error) {
	query := `INSERT INTO checklist_non_conformities AS t (
		id, tenant_id, execution_id, response_id, title, description,
		severity, status, recommendation, resolved_at, audit_id, version,
		created_at, updated_at, deleted_at
	) VALUES (
		$1, $2::uuid, $3::uuid, $4, $5, $6,
		$7, $8, $9, $10, NULL, 1,
		$11, $11, NULL
	) RETURNING row_to_json(t)`
	return read[domain.ChecklistNonConformity](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.ExecutionID, item.ResponseID, item.Title, item.Description,
		item.Severity, item.Status, item.Recommendation, item.ResolvedAt, item.CreatedAt,
	))
}

func (r NonConformityRepository) ListByExecution(ctx context.Context, tenantID string, executionID string, q ports.Query) (ports.Page[domain.ChecklistNonConformity], error) {
	where := []string{"tenant_id = $1::uuid", "execution_id = $2::uuid", "deleted_at IS NULL"}
	args := []any{tenantID, executionID}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM checklist_non_conformities WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.ChecklistNonConformity]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM checklist_non_conformities
			WHERE %s
			ORDER BY created_at ASC
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.ChecklistNonConformity]{}, mapError(err)
	}

	var items []domain.ChecklistNonConformity
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.ChecklistNonConformity]{}, err
	}

	return ports.Page[domain.ChecklistNonConformity]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

// ==================== SignatureRepository ====================

type SignatureRepository struct{ db *sql.DB }

func (r SignatureRepository) Create(ctx context.Context, item domain.ChecklistSignature) (domain.ChecklistSignature, error) {
	query := `INSERT INTO checklist_signatures AS t (
		id, tenant_id, execution_id, signer, signature_reference, signature_type,
		signed_at, audit_id, version, created_at, updated_at, deleted_at
	) VALUES (
		$1, $2::uuid, $3::uuid, $4, $5, $6,
		$7, NULL, 1, now(), now(), NULL
	) RETURNING row_to_json(t)`
	return read[domain.ChecklistSignature](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.ExecutionID, item.Signer, item.SignatureReference, item.SignatureType,
		item.SignedAt,
	))
}

func (r SignatureRepository) ExistsForExecution(ctx context.Context, tenantID string, executionID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM checklist_signatures WHERE tenant_id = $1::uuid AND execution_id = $2::uuid AND deleted_at IS NULL)`
	var exists bool
	err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query, tenantID, executionID).Scan(&exists)
	return exists, mapError(err)
}

// ==================== HistoryRepository ====================

type HistoryRepository struct{ db *sql.DB }

func (r HistoryRepository) Create(ctx context.Context, item domain.ChecklistHistory) (domain.ChecklistHistory, error) {
	query := `INSERT INTO checklist_history AS t (
		id, tenant_id, execution_id, event, actor_id, notes,
		audit_id, version, created_at, updated_at, deleted_at
	) VALUES (
		$1, $2::uuid, $3::uuid, $4, $5, $6,
		NULL, 1, $7, $7, NULL
	) RETURNING row_to_json(t)`
	return read[domain.ChecklistHistory](platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, query,
		item.ID, item.TenantID, item.ExecutionID, item.Event, item.ActorID, item.Notes,
		item.CreatedAt,
	))
}

func (r HistoryRepository) ListByExecution(ctx context.Context, tenantID string, executionID string, q ports.Query) (ports.Page[domain.ChecklistHistory], error) {
	where := []string{"tenant_id = $1::uuid", "execution_id = $2::uuid", "deleted_at IS NULL"}
	args := []any{tenantID, executionID}

	whereClause := strings.Join(where, " AND ")
	countQuery := fmt.Sprintf("SELECT count(*) FROM checklist_history WHERE %s", whereClause)

	var totalItems int
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return ports.Page[domain.ChecklistHistory]{}, mapError(err)
	}

	page := q.Page
	if page < 1 {
		page = 1
	}
	pageSize := q.PageSize
	if pageSize < 1 {
		pageSize = 50
	}
	totalPages := (totalItems + pageSize - 1) / pageSize
	offset := (page - 1) * pageSize

	dataQuery := fmt.Sprintf(`SELECT coalesce(json_agg(row_to_json(t)), '[]'::json)
		FROM (
			SELECT * FROM checklist_history
			WHERE %s
			ORDER BY created_at ASC
			LIMIT $%d OFFSET $%d
		) AS t`, whereClause, len(args)+1, len(args)+2)

	args = append(args, pageSize, offset)
	var body []byte
	if err := platform.ExecutorFor(ctx, r.db).QueryRowContext(ctx, dataQuery, args...).Scan(&body); err != nil {
		return ports.Page[domain.ChecklistHistory]{}, mapError(err)
	}

	var items []domain.ChecklistHistory
	if err := json.Unmarshal(body, &items); err != nil {
		return ports.Page[domain.ChecklistHistory]{}, err
	}

	return ports.Page[domain.ChecklistHistory]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}, nil
}

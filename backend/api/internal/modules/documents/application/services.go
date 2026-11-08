package application

import (
	"context"
	"errors"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/documents/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/documents/ports"
	"github.com/geevheit/intelligence360/backend/api/internal/shared/contextkeys"
	"github.com/google/uuid"
	"strings"
	"time"
)

var (
	ErrNotFound          = errors.New("document not found")
	ErrValidation        = errors.New("invalid document fields or dates")
	ErrConflict          = errors.New("document number or version conflict")
	ErrReferenceNotFound = errors.New("document owner not found")
)

type Service struct {
	repo     ports.DocumentRepository
	subjects ports.SubjectAccess
	tx       ports.Transactions
	audit    ports.Audit
	now      func() time.Time
}

func NewService(repo ports.DocumentRepository, subjects ports.SubjectAccess, tx ports.Transactions, audit ports.Audit) Service {
	return Service{repo: repo, subjects: subjects, tx: tx, audit: audit, now: time.Now}
}
func (s Service) WithClock(now func() time.Time) Service { s.now = now; return s }
func validID(value string) bool                          { _, err := uuid.Parse(value); return err == nil }
func normalize(d *domain.Document) error {
	d.Title = strings.TrimSpace(d.Title)
	d.Number = strings.TrimSpace(d.Number)
	d.Issuer = strings.TrimSpace(d.Issuer)
	d.Kind = strings.ToLower(strings.TrimSpace(d.Kind))
	if !validID(d.TenantID) || d.Title == "" || len(d.Title) > 200 || d.Number == "" || len(d.Number) > 100 || len(d.Issuer) > 200 || len(d.Notes) > 4000 {
		return ErrValidation
	}
	if (d.DriverID == "") == (d.VehicleID == "") {
		return ErrValidation
	}
	if (d.DriverID != "" && !validID(d.DriverID)) || (d.VehicleID != "" && !validID(d.VehicleID)) {
		return ErrValidation
	}
	switch d.Kind {
	case "cnh":
		if d.DriverID == "" {
			return ErrValidation
		}
	case "crlv":
		if d.VehicleID == "" {
			return ErrValidation
		}
	case "insurance", "other":
	default:
		return ErrValidation
	}
	if d.Status == "" {
		d.Status = "active"
	}
	if d.Status != "active" && d.Status != "archived" {
		return ErrValidation
	}
	for _, value := range []*string{d.IssuedOn, d.ExpiresOn} {
		if value != nil {
			date, err := time.Parse(time.DateOnly, *value)
			if err != nil || date.Year() < 1900 || date.Year() > 9999 {
				return ErrValidation
			}
		}
	}
	if d.IssuedOn != nil && d.ExpiresOn != nil && *d.ExpiresOn < *d.IssuedOn {
		return ErrValidation
	}
	return nil
}
func (s Service) subject(ctx context.Context, d domain.Document) error {
	if d.DriverID != "" {
		return s.subjects.ValidateDriverAccess(ctx, d.TenantID, d.DriverID)
	}
	return s.subjects.ValidateVehicleAccess(ctx, d.TenantID, d.VehicleID)
}
func (s Service) record(ctx context.Context, d domain.Document, action string, previous *domain.Document) error {
	actor, _ := contextkeys.ActorIDFromContext(ctx)
	date := func(value *string) string {
		if value == nil {
			return ""
		}
		return *value
	}
	metadata := map[string]string{"driver_id": d.DriverID, "vehicle_id": d.VehicleID, "kind": d.Kind, "status": d.Status, "issued_on": date(d.IssuedOn), "expires_on": date(d.ExpiresOn)}
	if previous != nil {
		metadata["previous_status"] = previous.Status
		metadata["previous_issued_on"] = date(previous.IssuedOn)
		metadata["previous_expires_on"] = date(previous.ExpiresOn)
	}
	return s.audit.RecordDocumentEvent(ctx, ports.AuditEvent{TenantID: d.TenantID, ActorID: actor, Action: action, ResourceID: d.ID, Metadata: metadata})
}

func (s Service) Create(ctx context.Context, d domain.Document) (domain.Document, error) {
	if err := normalize(&d); err != nil {
		return domain.Document{}, err
	}
	now := s.now().UTC()
	d.ID = uuid.NewString()
	d.CreatedAt = now
	d.UpdatedAt = now
	d.DeletedAt = nil
	d.Version = 1
	var saved domain.Document
	err := s.tx.Run(ctx, func(ctx context.Context) error {
		if err := s.subject(ctx, d); err != nil {
			return err
		}
		var err error
		saved, err = s.repo.Create(ctx, d)
		if err != nil {
			return err
		}
		return s.record(ctx, saved, "document.created", nil)
	})
	saved.ExpiryStatus = saved.ExpiryAt(now)
	return saved, err
}
func (s Service) FindByID(ctx context.Context, tenant, id string) (domain.Document, error) {
	if !validID(tenant) || !validID(id) {
		return domain.Document{}, ErrNotFound
	}
	d, err := s.repo.FindByID(ctx, tenant, id)
	if err == nil {
		d.ExpiryStatus = d.ExpiryAt(s.now())
	}
	return d, err
}
func (s Service) ValidateDocumentAccess(ctx context.Context, tenant, id string) error {
	_, err := s.FindByID(ctx, tenant, id)
	return err
}
func (s Service) Search(ctx context.Context, tenant string, q ports.Query) (ports.Page, error) {
	if !validID(tenant) || len(q.Search) > 200 || (q.DriverID != "" && !validID(q.DriverID)) || (q.VehicleID != "" && !validID(q.VehicleID)) || (q.DriverID != "" && q.VehicleID != "") {
		return ports.Page{}, ErrValidation
	}
	if q.Kind != "" && q.Kind != "cnh" && q.Kind != "crlv" && q.Kind != "insurance" && q.Kind != "other" {
		return ports.Page{}, ErrValidation
	}
	if q.Status != "" && q.Status != "active" && q.Status != "archived" {
		return ports.Page{}, ErrValidation
	}
	if q.ExpiryStatus != "" && q.ExpiryStatus != "expired" && q.ExpiryStatus != "due_soon" && q.ExpiryStatus != "valid" && q.ExpiryStatus != "no_expiry" && q.ExpiryStatus != "archived" {
		return ports.Page{}, ErrValidation
	}
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Page > 1000000 {
		return ports.Page{}, ErrValidation
	}
	if q.PageSize < 1 || q.PageSize > 100 {
		q.PageSize = 20
	}
	q.AsOf = s.now().UTC()
	page, err := s.repo.Search(ctx, tenant, q)
	if err != nil {
		return page, err
	}
	for i := range page.Items {
		page.Items[i].ExpiryStatus = page.Items[i].ExpiryAt(q.AsOf)
	}
	return page, nil
}
func (s Service) Update(ctx context.Context, d domain.Document) (domain.Document, error) {
	if !validID(d.ID) {
		return domain.Document{}, ErrNotFound
	}
	if d.Version < 1 {
		return domain.Document{}, ErrValidation
	}
	if err := normalize(&d); err != nil {
		return domain.Document{}, err
	}
	now := s.now().UTC()
	var saved domain.Document
	err := s.tx.Run(ctx, func(ctx context.Context) error {
		current, err := s.repo.FindByID(ctx, d.TenantID, d.ID)
		if err != nil {
			return err
		}
		if current.Version != d.Version {
			return ErrConflict
		}
		// Keep historical ownership stable; corrections create a new document record.
		if current.DriverID != d.DriverID || current.VehicleID != d.VehicleID {
			return ErrValidation
		}
		d.CreatedAt = current.CreatedAt
		d.UpdatedAt = now
		d.DeletedAt = nil
		saved, err = s.repo.Update(ctx, d)
		if err != nil {
			return err
		}
		return s.record(ctx, saved, "document.updated", &current)
	})
	saved.ExpiryStatus = saved.ExpiryAt(now)
	return saved, err
}
func (s Service) Delete(ctx context.Context, tenant, id string, version int64) error {
	if !validID(tenant) || !validID(id) {
		return ErrNotFound
	}
	if version < 1 {
		return ErrValidation
	}
	return s.tx.Run(ctx, func(ctx context.Context) error {
		d, err := s.repo.FindByID(ctx, tenant, id)
		if err != nil {
			return err
		}
		if err := s.repo.Delete(ctx, tenant, id, version); err != nil {
			return err
		}
		return s.record(ctx, d, "document.deleted", nil)
	})
}

package application

import (
	"context"
	"errors"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/drivers/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/drivers/ports"
	"github.com/geevheit/intelligence360/backend/api/internal/shared/contextkeys"
	"github.com/google/uuid"
	"net/mail"
	"strings"
	"time"
	"unicode"
)

var (
	ErrNotFound          = errors.New("driver not found")
	ErrValidation        = errors.New("invalid driver fields")
	ErrConflict          = errors.New("driver code, tax ID or version conflict")
	ErrReferenceNotFound = errors.New("vehicle not found")
)

type Service struct {
	repo     ports.DriverRepository
	vehicles ports.VehicleAccess
	tx       ports.Transactions
	audit    ports.Audit
}

func NewService(repo ports.DriverRepository, vehicles ports.VehicleAccess, tx ports.Transactions, audit ports.Audit) Service {
	return Service{repo, vehicles, tx, audit}
}
func validID(value string) bool { _, err := uuid.Parse(value); return err == nil }
func normalize(d *domain.Driver) error {
	d.Code = strings.ToUpper(strings.TrimSpace(d.Code))
	d.Name = strings.TrimSpace(d.Name)
	d.Email = strings.ToLower(strings.TrimSpace(d.Email))
	d.TaxID = strings.Map(func(r rune) rune {
		if r == '.' || r == '-' || unicode.IsSpace(r) {
			return -1
		}
		return r
	}, d.TaxID)
	if !validID(d.TenantID) || d.Code == "" || len(d.Code) > 64 || d.Name == "" || len(d.Name) > 200 || len(d.Notes) > 4000 || len(d.Phone) > 40 {
		return ErrValidation
	}
	if d.Status == "" {
		d.Status = "active"
	}
	if d.Status != "active" && d.Status != "inactive" {
		return ErrValidation
	}
	if d.Email != "" {
		address, err := mail.ParseAddress(d.Email)
		if err != nil || address.Address != d.Email || len(d.Email) > 254 {
			return ErrValidation
		}
	}
	if d.TaxID != "" {
		if len(d.TaxID) != 11 {
			return ErrValidation
		}
		for _, r := range d.TaxID {
			if r < '0' || r > '9' {
				return ErrValidation
			}
		}
	}
	if d.VehicleID != "" && !validID(d.VehicleID) {
		return ErrValidation
	}
	return nil
}
func (s Service) validateVehicle(ctx context.Context, d domain.Driver) error {
	if d.VehicleID == "" {
		return nil
	}
	return s.vehicles.ValidateVehicleAccess(ctx, d.TenantID, d.VehicleID)
}
func (s Service) record(ctx context.Context, d domain.Driver, action, previousVehicle string) error {
	actor, _ := contextkeys.ActorIDFromContext(ctx)
	return s.audit.RecordDriverEvent(ctx, ports.AuditEvent{TenantID: d.TenantID, ActorID: actor, Action: action, ResourceID: d.ID, Metadata: map[string]string{"previous_vehicle_id": previousVehicle, "vehicle_id": d.VehicleID}})
}
func (s Service) Create(ctx context.Context, d domain.Driver) (domain.Driver, error) {
	if err := normalize(&d); err != nil {
		return domain.Driver{}, err
	}
	now := time.Now().UTC()
	d.ID = uuid.NewString()
	d.CreatedAt = now
	d.UpdatedAt = now
	d.DeletedAt = nil
	d.Version = 1
	var saved domain.Driver
	err := s.tx.Run(ctx, func(ctx context.Context) error {
		if err := s.validateVehicle(ctx, d); err != nil {
			return err
		}
		var err error
		saved, err = s.repo.Create(ctx, d)
		if err != nil {
			return err
		}
		return s.record(ctx, saved, "driver.created", "")
	})
	return saved, err
}
func (s Service) FindByID(ctx context.Context, tenant, id string) (domain.Driver, error) {
	if !validID(tenant) || !validID(id) {
		return domain.Driver{}, ErrNotFound
	}
	return s.repo.FindByID(ctx, tenant, id)
}
func (s Service) ValidateDriverAccess(ctx context.Context, tenant, id string) error {
	_, err := s.FindByID(ctx, tenant, id)
	return err
}
func (s Service) Search(ctx context.Context, tenant string, q ports.Query) (ports.Page, error) {
	if !validID(tenant) || (q.Status != "" && q.Status != "active" && q.Status != "inactive") || (q.VehicleID != "" && !validID(q.VehicleID)) || len(q.Search) > 200 {
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
	return s.repo.Search(ctx, tenant, q)
}
func (s Service) Update(ctx context.Context, d domain.Driver) (domain.Driver, error) {
	if !validID(d.ID) {
		return domain.Driver{}, ErrNotFound
	}
	if d.Version < 1 {
		return domain.Driver{}, ErrValidation
	}
	if err := normalize(&d); err != nil {
		return domain.Driver{}, err
	}
	var saved domain.Driver
	err := s.tx.Run(ctx, func(ctx context.Context) error {
		current, err := s.repo.FindByID(ctx, d.TenantID, d.ID)
		if err != nil {
			return err
		}
		if current.Version != d.Version {
			return ErrConflict
		}
		if err := s.validateVehicle(ctx, d); err != nil {
			return err
		}
		d.CreatedAt = current.CreatedAt
		d.DeletedAt = nil
		d.UpdatedAt = time.Now().UTC()
		saved, err = s.repo.Update(ctx, d)
		if err != nil {
			return err
		}
		return s.record(ctx, saved, "driver.updated", current.VehicleID)
	})
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
		current, err := s.repo.FindByID(ctx, tenant, id)
		if err != nil {
			return err
		}
		if err := s.repo.Delete(ctx, tenant, id, version); err != nil {
			return err
		}
		return s.record(ctx, current, "driver.deleted", current.VehicleID)
	})
}

package di

import (
	"context"
	"errors"
	documentsapp "github.com/geevheit/intelligence360/backend/api/internal/modules/documents/application"
	documentsports "github.com/geevheit/intelligence360/backend/api/internal/modules/documents/ports"
	driversapp "github.com/geevheit/intelligence360/backend/api/internal/modules/drivers/application"
	driversports "github.com/geevheit/intelligence360/backend/api/internal/modules/drivers/ports"
	fleetapp "github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/application"
	fleetports "github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/ports"
	identitydomain "github.com/geevheit/intelligence360/backend/api/internal/modules/identity/domain"
	identityports "github.com/geevheit/intelligence360/backend/api/internal/modules/identity/ports"
)

// Composition translates errors at bounded-context contracts without coupling
// Drivers or Documents to the implementation details of another module.
func driverVehicleAccess(fleet fleetports.VehicleService) driversports.VehicleAccess {
	return driversports.VehicleAccessFunc(func(ctx context.Context, tenant, id string) error {
		err := fleet.ValidateVehicleAccess(ctx, tenant, id)
		if errors.Is(err, fleetapp.ErrNotFound) {
			return driversapp.ErrReferenceNotFound
		}
		return err
	})
}

type documentSubjects struct {
	drivers driversports.DriversService
	fleet   fleetports.VehicleService
}

func (s documentSubjects) ValidateDriverAccess(ctx context.Context, tenant, id string) error {
	err := s.drivers.ValidateDriverAccess(ctx, tenant, id)
	if errors.Is(err, driversapp.ErrNotFound) {
		return documentsapp.ErrReferenceNotFound
	}
	return err
}
func (s documentSubjects) ValidateVehicleAccess(ctx context.Context, tenant, id string) error {
	err := s.fleet.ValidateVehicleAccess(ctx, tenant, id)
	if errors.Is(err, fleetapp.ErrNotFound) {
		return documentsapp.ErrReferenceNotFound
	}
	return err
}

type driverAudit struct{ service identityports.AuditService }

func (a driverAudit) RecordDriverEvent(ctx context.Context, event driversports.AuditEvent) error {
	return a.service.Record(ctx, identitydomain.AuditLog{TenantID: event.TenantID, ActorID: event.ActorID, Action: event.Action, ResourceType: "driver", ResourceID: event.ResourceID, Metadata: event.Metadata})
}

type documentAudit struct{ service identityports.AuditService }

func (a documentAudit) RecordDocumentEvent(ctx context.Context, event documentsports.AuditEvent) error {
	return a.service.Record(ctx, identitydomain.AuditLog{TenantID: event.TenantID, ActorID: event.ActorID, Action: event.Action, ResourceType: "document", ResourceID: event.ResourceID, Metadata: event.Metadata})
}

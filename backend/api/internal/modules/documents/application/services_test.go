package application

import (
	"github.com/geevheit/intelligence360/backend/api/internal/modules/documents/domain"
	"github.com/google/uuid"
	"testing"
)

func TestDocumentOwnershipAndDates(t *testing.T) {
	base := domain.Document{TenantID: uuid.NewString(), DriverID: uuid.NewString(), Kind: "cnh", Title: "CNH de teste", Number: "TESTE"}
	tests := []struct {
		name string
		edit func(*domain.Document)
	}{
		{"owner required", func(d *domain.Document) { d.DriverID = "" }},
		{"only one owner", func(d *domain.Document) { d.VehicleID = uuid.NewString() }},
		{"CNH requires driver", func(d *domain.Document) { d.VehicleID = d.DriverID; d.DriverID = "" }},
		{"CRLV requires vehicle", func(d *domain.Document) { d.Kind = "crlv" }},
		{"unsupported kind", func(d *domain.Document) { d.Kind = "arbitrary" }},
		{"invalid owner ID", func(d *domain.Document) { d.DriverID = "bad" }},
		{"invalid calendar date", func(d *domain.Document) { value := "2026-02-30"; d.ExpiresOn = &value }},
		{"empty date is not null", func(d *domain.Document) { value := ""; d.ExpiresOn = &value }},
		{"expiry before issue", func(d *domain.Document) {
			issue, expiry := "2026-11-07", "2026-11-06"
			d.IssuedOn = &issue
			d.ExpiresOn = &expiry
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := base
			tt.edit(&d)
			if err := normalize(&d); err != ErrValidation {
				t.Fatalf("expected validation error, got %v", err)
			}
		})
	}
	if err := normalize(&base); err != nil {
		t.Fatal(err)
	}
	vehicle := base
	vehicle.DriverID = ""
	vehicle.VehicleID = uuid.NewString()
	vehicle.Kind = "crlv"
	if err := normalize(&vehicle); err != nil {
		t.Fatal(err)
	}
}

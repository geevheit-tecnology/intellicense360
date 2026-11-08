package domain

import (
	"testing"
	"time"
)

func TestExpiryCalendarBoundaries(t *testing.T) {
	now := time.Date(2026, 11, 7, 23, 59, 59, 0, time.UTC)
	tests := []struct {
		name   string
		date   string
		status string
		want   string
	}{
		{"no date", "", "active", "no_expiry"},
		{"previous day", "2026-11-06", "active", "expired"},
		{"last second of expiry day", "2026-11-07", "active", "due_soon"},
		{"30-day boundary", "2026-12-07", "active", "due_soon"},
		{"31 days", "2026-12-08", "active", "valid"},
		{"archived old date", "2020-01-01", "archived", "archived"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Document{Status: tt.status}
			if tt.date != "" {
				d.ExpiresOn = &tt.date
			}
			if got := d.ExpiryAt(now); got != tt.want {
				t.Fatalf("got %s want %s", got, tt.want)
			}
		})
	}
	date := "2026-11-07"
	d := Document{Status: "active", ExpiresOn: &date}
	if got := d.ExpiryAt(now.Add(time.Second)); got != "expired" {
		t.Fatalf("midnight rollover: %s", got)
	}
}

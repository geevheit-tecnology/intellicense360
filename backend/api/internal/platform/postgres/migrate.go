package postgres

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Migrate applies only the schema needed by the current durable runtime.
// Other modules and development seeds are deliberately excluded. Historical SQL
// remains unchanged; extending persistence requires extending this reviewed list.
// The batch and all checksums commit together.
func Migrate(ctx context.Context, db *sql.DB, directory string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(directory, "*.sql"))
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no migrations in %s", directory)
	}
	found := map[string]bool{}
	for _, path := range files {
		found[filepath.Base(path)] = true
	}
	for name := range runtimeMigrations {
		if !found[name] {
			return nil, fmt.Errorf("required migration missing: %s", name)
		}
	}
	sort.Strings(files)
	conn, err := db.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	// Transaction-scoped lock avoids leaking a session advisory lock on cancellation.
	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock(3600020)`); err != nil {
		return nil, err
	}
	if _, err = tx.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (name text PRIMARY KEY, checksum text NOT NULL, applied_at timestamptz NOT NULL DEFAULT now())`); err != nil {
		return nil, err
	}
	applied := []string{}
	for _, path := range files {
		name := filepath.Base(path)
		if !runtimeMigrations[name] {
			continue
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		checksum := fmt.Sprintf("%x", sha256.Sum256(body))
		var existing string
		err = tx.QueryRowContext(ctx, `SELECT checksum FROM schema_migrations WHERE name=$1`, name).Scan(&existing)
		if err == nil {
			if checksum != existing {
				return nil, fmt.Errorf("migration checksum changed: %s", name)
			}
			continue
		}
		if err != sql.ErrNoRows {
			return nil, err
		}
		if _, err := tx.ExecContext(ctx, string(body)); err != nil {
			return nil, fmt.Errorf("migration %s: %w", name, err)
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations(name,checksum) VALUES ($1,$2)`, name, checksum); err != nil {
			return nil, err
		}
		applied = append(applied, name)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return applied, nil
}

var runtimeMigrations = map[string]bool{
	"000001_init.sql":                                    true,
	"000002_identity_platform.sql":                       true,
	"000003_fleet_core.sql":                              true,
	"000005_checklist_core.sql":                          true,
	"000006_tires_core.sql":                              true,
	"000008_maintenance_core.sql":                       true,
	"000009_maintenance_prompt_alignment.sql":           true,
	"000010_inventory_core.sql":                         true,
	"000012_tire_lifecycle_core.sql":                     true,
	"000013_fuel_core.sql":                               true,
	"000014_checklist_engine_core.sql":                  true,
	"000015_financial_cost_core.sql":                    true,
	"000017_event_bus_outbox.sql":                        true,
	"000018_intelligence_core.sql":                      true,
	"000019_mission_control_core.sql":                   true,
	"000020_identity_fleet_tenant_constraints.sql":       true,
	"000021_drivers_documents_core.sql":                  true,
	"000022_odometer_km_history.sql":                     true,
	"000023_fuel_operational_financial_alignment.sql":   true,
	"000024_tires_lifecycle_alignment.sql":              true,
	"000025_maintenance_operational_alignment.sql":       true,
	"000026_checklist_operational_alignment.sql":         true,
	"000027_mission_control_alignment.sql":               true,
	"000028_financial_cashflow_reconciliation.sql":       true,
	"000029_inventory_operational_alignment.sql":        true,
	"000030_trips_dispatch_core.sql":                    true,
	"000031_driver_journey_compliance.sql":              true,
}

// CheckMigrations is read-only; startup never performs DDL or seeds users.
func CheckMigrations(ctx context.Context, db *sql.DB) error {
	// Single round trip: remote/serverless databases (e.g. Neon) add latency per query.
	rows, err := db.QueryContext(ctx, `SELECT name FROM schema_migrations`)
	if err != nil {
		return fmt.Errorf("schema not ready; run cmd/migrate explicitly: %w", err)
	}
	defer rows.Close()
	applied := make(map[string]bool, len(runtimeMigrations))
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return fmt.Errorf("schema not ready; run cmd/migrate explicitly: %w", err)
		}
		applied[name] = true
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("schema not ready; run cmd/migrate explicitly: %w", err)
	}
	for name := range runtimeMigrations {
		if !applied[name] {
			return fmt.Errorf("migration pending: %s; run cmd/migrate explicitly", name)
		}
	}
	return nil
}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/domain"
	platform "github.com/geevheit/intelligence360/backend/api/internal/platform/postgres"
)

type SignalScanner struct {
	db *sql.DB
}

func NewSignalScanner(db *sql.DB) *SignalScanner {
	return &SignalScanner{db: db}
}

// ScanOperationalSignals scans database tables for proactive alert conditions.
func (s *SignalScanner) ScanOperationalSignals(ctx context.Context, tenantID string) ([]domain.CommandItem, error) {
	now := time.Now().UTC()
	var signals []domain.CommandItem

	exec := platform.ExecutorFor(ctx, s.db)

	// 1. Expired or upcoming document expirations (CNH, CRLV, etc.)
	docQuery := `
		SELECT id, category, subject_type, subject_id, expires_at
		FROM drivers_documents
		WHERE tenant_id = $1::uuid
		  AND deleted_at IS NULL
		  AND expires_at IS NOT NULL
		  AND expires_at <= $2
	`
	thirtyDaysAhead := now.AddDate(0, 0, 30)
	docRows, err := exec.QueryContext(ctx, docQuery, tenantID, thirtyDaysAhead)
	if err == nil {
		defer docRows.Close()
		for docRows.Next() {
			var docID, category, subjectType, subjectID string
			var expiresAt time.Time
			if err := docRows.Scan(&docID, &category, &subjectType, &subjectID, &expiresAt); err == nil {
				isExpired := expiresAt.Before(now)
				severity := domain.SeverityHigh
				priority := domain.PriorityHigh
				itemType := domain.TypeWarning
				title := fmt.Sprintf("Documento %s (%s) a vencer em %s", category, subjectType, expiresAt.Format("02/01/2006"))
				if isExpired {
					severity = domain.SeverityCritical
					priority = domain.PriorityUrgent
					itemType = domain.TypeAlert
					title = fmt.Sprintf("Documento %s (%s) VENCIDO em %s", category, subjectType, expiresAt.Format("02/01/2006"))
				}

				signals = append(signals, domain.CommandItem{
					TenantID:     tenantID,
					Type:         itemType,
					Category:     domain.CategoryDocument,
					Severity:     severity,
					Priority:     priority,
					Title:        title,
					Description:  fmt.Sprintf("Vencimento do documento %s para %s %s", category, subjectType, subjectID),
					Source:       "rules_engine",
					SourceType:   "drivers_documents",
					SourceID:     docID,
					Confidence:   1.0,
					ImpactScore:  0.8,
					RiskScore:    0.9,
					UrgencyScore: 0.9,
					DetectedAt:   now,
					Fingerprint:  domain.Fingerprint(tenantID, itemType, domain.CategoryDocument, "drivers_documents", docID, title),
				})
			}
		}
	}

	// 2. Tire tread wear safety limit reached
	tireQuery := `
		SELECT id, fire_number, current_tread_mm, minimum_tread_mm
		FROM tires
		WHERE tenant_id = $1::uuid
		  AND deleted_at IS NULL
		  AND current_tread_mm <= minimum_tread_mm
	`
	tireRows, err := exec.QueryContext(ctx, tireQuery, tenantID)
	if err == nil {
		defer tireRows.Close()
		for tireRows.Next() {
			var tireID, fireNumber string
			var currentTread, minTread float64
			if err := tireRows.Scan(&tireID, &fireNumber, &currentTread, &minTread); err == nil {
				title := fmt.Sprintf("Pneu %s atingiu limite de segurança: %.1f mm (mínimo: %.1f mm)", fireNumber, currentTread, minTread)
				signals = append(signals, domain.CommandItem{
					TenantID:     tenantID,
					Type:         domain.TypeAlert,
					Category:     domain.CategoryTire,
					Severity:     domain.SeverityCritical,
					Priority:     domain.PriorityUrgent,
					Title:        title,
					Description:  "Necessário rodízio imediato, recapagem ou troca para evitar estouro/multa rodoviária.",
					Source:       "rules_engine",
					SourceType:   "tires",
					SourceID:     tireID,
					Confidence:   1.0,
					ImpactScore:  0.85,
					RiskScore:    0.9,
					UrgencyScore: 0.95,
					DetectedAt:   now,
					Fingerprint:  domain.Fingerprint(tenantID, domain.TypeAlert, domain.CategoryTire, "tires", tireID, title),
				})
			}
		}
	}

	// 3. Preventive maintenance due by KM
	planQuery := `
		SELECT p.id, p.name, p.vehicle_id, p.next_due_value,
		       COALESCE(o.current_km, 0) AS current_km
		FROM maintenance_preventive_plans p
		LEFT JOIN LATERAL (
			SELECT current_km FROM fleet_odometer_readings
			WHERE tenant_id = p.tenant_id AND vehicle_id = p.vehicle_id
			ORDER BY reading_at DESC LIMIT 1
		) o ON true
		WHERE p.tenant_id = $1::uuid
		  AND p.deleted_at IS NULL
		  AND p.active = true
		  AND p.frequency = 'km'
		  AND p.vehicle_id IS NOT NULL
		  AND p.next_due_value > 0
		  AND COALESCE(o.current_km, 0) >= (p.next_due_value - 1000)
	`
	planRows, err := exec.QueryContext(ctx, planQuery, tenantID)
	if err == nil {
		defer planRows.Close()
		for planRows.Next() {
			var planID, planName, vehicleID string
			var nextDueKM, currentKM int64
			if err := planRows.Scan(&planID, &planName, &vehicleID, &nextDueKM, &currentKM); err == nil {
				isOverdue := currentKM >= nextDueKM
				title := fmt.Sprintf("Plano preventivo %q próximo (%d km / meta %d km)", planName, currentKM, nextDueKM)
				severity := domain.SeverityMedium
				priority := domain.PriorityNormal
				itemType := domain.TypeWarning
				if isOverdue {
					title = fmt.Sprintf("Plano preventivo %q VENCIDO (%d km / meta %d km)", planName, currentKM, nextDueKM)
					severity = domain.SeverityHigh
					priority = domain.PriorityHigh
					itemType = domain.TypeAlert
				}

				signals = append(signals, domain.CommandItem{
					TenantID:     tenantID,
					Type:         itemType,
					Category:     domain.CategoryMaintenance,
					Severity:     severity,
					Priority:     priority,
					Title:        title,
					Description:  fmt.Sprintf("Veículo %s atingiu marco de revisão preventiva.", vehicleID),
					Source:       "rules_engine",
					SourceType:   "maintenance_preventive_plans",
					SourceID:     planID,
					Confidence:   0.95,
					ImpactScore:  0.75,
					RiskScore:    0.8,
					UrgencyScore: 0.8,
					DetectedAt:   now,
					Fingerprint:  domain.Fingerprint(tenantID, itemType, domain.CategoryMaintenance, "maintenance_preventive_plans", planID, title),
				})
			}
		}
	}

	// 4. Critical checklist non-conformities unresolved
	ncQuery := `
		SELECT id, execution_id, title, severity, status
		FROM checklist_non_conformities
		WHERE tenant_id = $1::uuid
		  AND deleted_at IS NULL
		  AND status != 'resolved'
		  AND severity IN ('critical', 'high')
	`
	ncRows, err := exec.QueryContext(ctx, ncQuery, tenantID)
	if err == nil {
		defer ncRows.Close()
		for ncRows.Next() {
			var ncID, execID, ncTitle, ncSeverity, ncStatus string
			if err := ncRows.Scan(&ncID, &execID, &ncTitle, &ncSeverity, &ncStatus); err == nil {
				title := fmt.Sprintf("Não-conformidade crítica de checklist: %s", ncTitle)
				signals = append(signals, domain.CommandItem{
					TenantID:     tenantID,
					Type:         domain.TypeAlert,
					Category:     domain.CategorySafety,
					Severity:     domain.SeverityCritical,
					Priority:     domain.PriorityUrgent,
					Title:        title,
					Description:  fmt.Sprintf("Identificada durante execução de checklist %s. Status: %s", execID, ncStatus),
					Source:       "rules_engine",
					SourceType:   "checklist_non_conformities",
					SourceID:     ncID,
					Confidence:   1.0,
					ImpactScore:  0.9,
					RiskScore:    0.95,
					UrgencyScore: 0.95,
					DetectedAt:   now,
					Fingerprint:  domain.Fingerprint(tenantID, domain.TypeAlert, domain.CategorySafety, "checklist_non_conformities", ncID, title),
				})
			}
		}
	}

	// 5. Driver Journey Violations & Continuous Driving Exceeded (Lei 13.103/2015)
	journeyQuery := `
		SELECT r.id, r.driver_id, COALESCE(d.name, 'Motorista'), r.event_type, r.start_at,
		       COALESCE(r.violation_reason, ''), r.is_violation,
		       EXTRACT(EPOCH FROM ($2 - r.start_at)) / 60.0 AS elapsed_minutes
		FROM driver_journey_records r
		LEFT JOIN drivers d ON d.id = r.driver_id
		WHERE r.tenant_id = $1::uuid
		  AND (
		      r.is_violation = true OR
		      (r.end_at IS NULL AND r.event_type = 'driving' AND EXTRACT(EPOCH FROM ($2 - r.start_at)) > 19800)
		  )
		ORDER BY r.start_at DESC
		LIMIT 15
	`
	journeyRows, err := exec.QueryContext(ctx, journeyQuery, tenantID, now)
	if err == nil {
		defer journeyRows.Close()
		for journeyRows.Next() {
			var recID, driverID, driverName, eventType string
			var startAt time.Time
			var violationReason string
			var isViolation bool
			var elapsedMin float64
			if err := journeyRows.Scan(&recID, &driverID, &driverName, &eventType, &startAt, &violationReason, &isViolation, &elapsedMin); err == nil {
				title := fmt.Sprintf("Infração de Jornada: %s em direção contínua sem pausa (Lei 13.103/2015)", driverName)
				desc := violationReason
				if desc == "" {
					desc = fmt.Sprintf("Motorista atingiu %.0f minutos de direção contínua sem registro de pausa obrigatória de 30 min.", elapsedMin)
				}
				signals = append(signals, domain.CommandItem{
					TenantID:     tenantID,
					Type:         domain.TypeAlert,
					Category:     domain.CategoryCompliance,
					Severity:     domain.SeverityCritical,
					Priority:     domain.PriorityUrgent,
					Title:        title,
					Description:  desc,
					Source:       "rules_engine",
					SourceType:   "driver_journey_records",
					SourceID:     recID,
					Confidence:   1.0,
					ImpactScore:  0.9,
					RiskScore:    0.95,
					UrgencyScore: 0.95,
					DetectedAt:   now,
					Fingerprint:  domain.Fingerprint(tenantID, domain.TypeAlert, domain.CategoryCompliance, "driver_journey_records", recID, title),
				})
			}
		}
	}

	// 6. Inventory Items Below Minimum Safety Stock (Almoxarifado)
	invQuery := `
		SELECT id, sku, name, current_stock, min_stock, location
		FROM inventory_parts
		WHERE tenant_id = $1::uuid
		  AND deleted_at IS NULL
		  AND min_stock > 0
		  AND current_stock <= min_stock
		ORDER BY (min_stock - current_stock) DESC
		LIMIT 15
	`
	invRows, err := exec.QueryContext(ctx, invQuery, tenantID)
	if err == nil {
		defer invRows.Close()
		for invRows.Next() {
			var partID, sku, partName, location string
			var currentStock, minStock float64
			if err := invRows.Scan(&partID, &sku, &partName, &currentStock, &minStock, &location); err == nil {
				title := fmt.Sprintf("Estoque crítico de %s (SKU: %s): saldo %.1f (mínimo: %.1f)", partName, sku, currentStock, minStock)
				signals = append(signals, domain.CommandItem{
					TenantID:     tenantID,
					Type:         domain.TypeWarning,
					Category:     domain.CategoryInventory,
					Severity:     domain.SeverityHigh,
					Priority:     domain.PriorityHigh,
					Title:        title,
					Description:  fmt.Sprintf("Item de almoxarifado em nível crítico de reposição. Localização: %s", location),
					Source:       "rules_engine",
					SourceType:   "inventory_parts",
					SourceID:     partID,
					Confidence:   1.0,
					ImpactScore:  0.75,
					RiskScore:    0.7,
					UrgencyScore: 0.8,
					DetectedAt:   now,
					Fingerprint:  domain.Fingerprint(tenantID, domain.TypeWarning, domain.CategoryInventory, "inventory_parts", partID, title),
				})
			}
		}
	}

	// 7. Trips Overdue in Transit (Despacho / Prazo de Chegada)
	tripQuery := `
		SELECT id, trip_number, driver_name, vehicle_plate, origin_city, destination_city, expected_arrival_at
		FROM operational_trips
		WHERE tenant_id = $1::uuid
		  AND deleted_at IS NULL
		  AND status = 'dispatched'
		  AND expected_arrival_at IS NOT NULL
		  AND expected_arrival_at < $2
		ORDER BY expected_arrival_at ASC
		LIMIT 10
	`
	tripRows, err := exec.QueryContext(ctx, tripQuery, tenantID, now)
	if err == nil {
		defer tripRows.Close()
		for tripRows.Next() {
			var tripID, tripNum, driverName, plate, origin, dest string
			var expectedArrival time.Time
			if err := tripRows.Scan(&tripID, &tripNum, &driverName, &plate, &origin, &dest, &expectedArrival); err == nil {
				title := fmt.Sprintf("Viagem #%s atrasada (%s ➔ %s) - Veículo %s", tripNum, origin, dest, plate)
				signals = append(signals, domain.CommandItem{
					TenantID:     tenantID,
					Type:         domain.TypeWarning,
					Category:     domain.CategoryOperational,
					Severity:     domain.SeverityHigh,
					Priority:     domain.PriorityHigh,
					Title:        title,
					Description:  fmt.Sprintf("Previsão de chegada expirou em %s. Condutor: %s", expectedArrival.Format("02/01/2006 15:04"), driverName),
					Source:       "rules_engine",
					SourceType:   "operational_trips",
					SourceID:     tripID,
					Confidence:   0.95,
					ImpactScore:  0.8,
					RiskScore:    0.8,
					UrgencyScore: 0.85,
					DetectedAt:   now,
					Fingerprint:  domain.Fingerprint(tenantID, domain.TypeWarning, domain.CategoryOperational, "operational_trips", tripID, title),
				})
			}
		}
	}

	return signals, nil
}


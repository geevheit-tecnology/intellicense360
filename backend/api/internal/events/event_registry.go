package events

import (
	"strconv"
	"strings"
)

const (
	AssetCreatedV1   = "asset.created.v1"
	AssetUpdatedV1   = "asset.updated.v1"
	DriverCreatedV1  = "driver.created.v1"
	VehicleCreatedV1 = "vehicle.created.v1"
	VehicleUpdatedV1 = "vehicle.updated.v1"

	MaintenanceOrderCreatedV1   = "maintenance.order.created.v1"
	MaintenanceOrderCompletedV1 = "maintenance.order.completed.v1"
	MaintenanceOrderCanceledV1  = "maintenance.order.canceled.v1"

	InventoryPartCreatedV1 = "inventory.part.created.v1"
	InventoryPartUpdatedV1 = "inventory.part.updated.v1"

	SupplierCreatedV1 = "supplier.created.v1"
	SupplierUpdatedV1 = "supplier.updated.v1"

	TireCreatedV1   = "tire.created.v1"
	TireInstalledV1 = "tire.installed.v1"
	TireRemovedV1   = "tire.removed.v1"
	TireInspectedV1 = "tire.inspected.v1"
	TireRetreadedV1 = "tire.retreaded.v1"
	TireDamagedV1   = "tire.damaged.v1"
	TireDisposedV1  = "tire.disposed.v1"

	FuelTransactionCreatedV1   = "fuel.transaction.created.v1"
	FuelTransactionCompletedV1 = "fuel.transaction.completed.v1"
	FuelTransactionCanceledV1  = "fuel.transaction.canceled.v1"
	FuelTransactionAdjustedV1  = "fuel.transaction.adjusted.v1"

	ChecklistTemplatePublishedV1    = "checklist.template.published.v1"
	ChecklistExecutionStartedV1     = "checklist.execution.started.v1"
	ChecklistExecutionCompletedV1   = "checklist.execution.completed.v1"
	ChecklistNonConformityCreatedV1 = "checklist.non_conformity.created.v1"

	FinancialExpenseCreatedV1    = "financial.expense.created.v1"
	FinancialExpenseApprovedV1   = "financial.expense.approved.v1"
	FinancialExpensePaidV1       = "financial.expense.paid.v1"
	FinancialRevenueCreatedV1    = "financial.revenue.created.v1"
	FinancialRevenueReceivedV1   = "financial.revenue.received.v1"
	FinancialAdjustmentCreatedV1 = "financial.adjustment.created.v1"

	CIOTCreatedV1     = "ciot.created.v1"
	CIOTActivatedV1   = "ciot.activated.v1"
	CIOTSuspendedV1   = "ciot.suspended.v1"
	CIOTReactivatedV1 = "ciot.reactivated.v1"
	CIOTClosedV1      = "ciot.closed.v1"
	CIOTCanceledV1    = "ciot.canceled.v1"
)

func VersionFromType(eventType string) (int, error) {
	parts := strings.Split(eventType, ".")
	if len(parts) < 3 || !strings.HasPrefix(parts[len(parts)-1], "v") {
		return 0, ErrInvalidEventType
	}
	version, err := strconv.Atoi(strings.TrimPrefix(parts[len(parts)-1], "v"))
	if err != nil || version <= 0 {
		return 0, ErrInvalidEventType
	}
	return version, nil
}

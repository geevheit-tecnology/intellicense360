package di

import (
	"github.com/geevheit/intelligence360/backend/api/internal/core/config"
	corehttp "github.com/geevheit/intelligence360/backend/api/internal/core/http"
	"github.com/geevheit/intelligence360/backend/api/internal/core/logger"
	assetsapp "github.com/geevheit/intelligence360/backend/api/internal/modules/assets/application"
	assetsinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/assets/infrastructure"
	assetshttp "github.com/geevheit/intelligence360/backend/api/internal/modules/assets/transport/http"
	checklistapp "github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/application"
	checklistinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/infrastructure"
	checklisthttp "github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/transport/http"
	ciotapp "github.com/geevheit/intelligence360/backend/api/internal/modules/ciot/application"
	ciotinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/ciot/infrastructure"
	ciothttp "github.com/geevheit/intelligence360/backend/api/internal/modules/ciot/transport/http"
	financialapp "github.com/geevheit/intelligence360/backend/api/internal/modules/financial/application"
	financialinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/financial/infrastructure"
	financialhttp "github.com/geevheit/intelligence360/backend/api/internal/modules/financial/transport/http"
	fleetapp "github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/application"
	fleetinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/infrastructure"
	fleethttp "github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/transport/http"
	fuelapp "github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/application"
	fuelinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/infrastructure"
	fuelhttp "github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/transport/http"
	identityapp "github.com/geevheit/intelligence360/backend/api/internal/modules/identity/application"
	identityinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/identity/infrastructure"
	identityhttp "github.com/geevheit/intelligence360/backend/api/internal/modules/identity/transport/http"
	intelligenceapp "github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/application"
	intelligencehttp "github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/transport/http"
	inventoryapp "github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/application"
	inventoryinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/infrastructure"
	inventoryhttp "github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/transport/http"
	maintenanceapp "github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/application"
	maintenanceinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/infrastructure"
	maintenancehttp "github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/transport/http"
	suppliersapp "github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/application"
	suppliersinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/infrastructure"
	suppliershttp "github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/transport/http"
	tiresapp "github.com/geevheit/intelligence360/backend/api/internal/modules/tires/application"
	tiresinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/tires/infrastructure"
	tireshttp "github.com/geevheit/intelligence360/backend/api/internal/modules/tires/transport/http"
	"github.com/gin-gonic/gin"
)

type Container struct {
	config config.Config
	logger logger.Logger
}

func NewContainer(cfg config.Config) *Container {
	return &Container{
		config: cfg,
		logger: logger.New(cfg.Env),
	}
}

func (c *Container) HTTPRouter() *gin.Engine {
	identityStore := identityinfra.NewMemoryStore()
	tokenService := identityapp.NewTokenService(c.config.AuthIssuer, c.config.AuthSecret)
	authService := identityapp.NewAuthService(identityStore.Users(), identityStore.Roles(), identityStore.Sessions(), identityStore.RefreshTokens(), tokenService)
	userService := identityapp.NewUserService(identityStore.Users())
	roleService := identityapp.NewRoleService(identityStore.Roles())
	permissionService := identityapp.NewPermissionService(identityStore.Permissions())
	tenantService := identityapp.NewTenantService(identityStore.Tenants())
	auditService := identityapp.NewAuditService(identityStore.Audit())
	identityHandler := identityhttp.NewHandler(authService, userService, roleService, permissionService, tenantService, auditService)

	assetsStore := assetsinfra.NewMemoryStore()
	assetsHandler := assetshttp.NewHandler(
		assetsapp.NewAssetService(assetsStore.Assets(), assetsStore.Audit()),
		assetsapp.NewCategoryService(assetsStore.Categories()),
		assetsapp.NewTypeService(assetsStore.Types()),
		assetsapp.NewManufacturerService(assetsStore.Manufacturers()),
		assetsapp.NewModelService(assetsStore.Models()),
		assetsapp.NewEquipmentService(assetsStore.Equipment()),
	)

	checklistStore := checklistinfra.NewMemoryStore()
	checklistHandler := checklisthttp.NewHandler(
		checklistapp.NewChecklistService(checklistStore.Checklists()),
		checklistapp.NewChecklistItemService(checklistStore.Checklists(), checklistStore.Items()),
		checklistapp.NewChecklistAnswerService(checklistStore.Items(), checklistStore.Answers()),
		checklistapp.NewChecklistTemplateService(checklistStore.Templates()),
		checklistapp.NewChecklistTemplateVersionService(checklistStore.Templates(), checklistStore.Versions()),
		checklistapp.NewChecklistTypeService(checklistStore.Types()),
		checklistapp.NewChecklistSectionService(checklistStore.Sections()),
		checklistapp.NewChecklistEngineItemService(checklistStore.Versions(), checklistStore.EngineItems()),
		checklistapp.NewChecklistExecutionService(checklistStore.Versions(), checklistStore.EngineItems(), checklistStore.Executions(), checklistStore.Responses(), checklistStore.Evidence(), checklistStore.Signatures(), checklistStore.History()),
		checklistapp.NewChecklistEngineResponseService(checklistStore.Executions(), checklistStore.EngineItems(), checklistStore.Responses(), checklistStore.History()),
		checklistapp.NewChecklistEvidenceService(checklistStore.Evidence(), checklistStore.History()),
		checklistapp.NewChecklistNonConformityService(checklistStore.NonConformities(), checklistStore.History()),
		checklistapp.NewChecklistHistoryService(checklistStore.History()),
	)

	fleetStore := fleetinfra.NewMemoryStore()
	fleetHandler := fleethttp.NewHandler(
		fleetapp.NewVehicleService(fleetStore.Vehicles()),
		fleetapp.NewVehicleBrandService(fleetStore.Brands()),
		fleetapp.NewVehicleModelService(fleetStore.Models()),
		fleetapp.NewVehicleCategoryService(fleetStore.Categories()),
		fleetapp.NewVehicleTypeService(fleetStore.Types()),
		fleetapp.NewAssetService(fleetStore.Assets()),
	)

	fuelStore := fuelinfra.NewMemoryStore()
	fuelHandler := fuelhttp.NewHandler(
		fuelapp.NewFuelTransactionService(fuelStore.Transactions(), fuelStore.Adjustments()),
		fuelapp.NewFuelTypeService(fuelStore.Types()),
		fuelapp.NewFuelStationService(fuelStore.Stations()),
		fuelapp.NewFuelTankService(fuelStore.Tanks()),
		fuelapp.NewFuelNozzleService(fuelStore.Nozzles()),
		fuelapp.NewFuelReadingService(fuelStore.Readings()),
		fuelapp.NewFuelPriceService(fuelStore.Prices()),
		fuelapp.NewFuelReceiptService(fuelStore.Receipts()),
		fuelapp.NewFuelAdjustmentService(fuelStore.Adjustments()),
	)

	financialStore := financialinfra.NewMemoryStore()
	financialHandler := financialhttp.NewHandler(
		financialapp.NewTransactionService(financialStore.Transactions(), financialStore.Periods(), financialStore.Adjustments()),
		financialapp.NewCatalogService(financialStore.Categories(), financialapp.InitCategory),
		financialapp.NewCatalogService(financialStore.Types(), financialapp.InitType),
		financialapp.NewCatalogService(financialStore.Centers(), financialapp.InitCenter),
		financialapp.NewCatalogService(financialStore.Accounts(), financialapp.InitAccount),
		financialapp.NewCatalogService(financialStore.Methods(), financialapp.InitPaymentMethod),
		financialapp.NewPeriodService(financialStore.Periods()),
		financialapp.NewBudgetService(financialStore.Budgets()),
		financialapp.NewAdjustmentService(financialStore.Adjustments()),
	)

	ciotStore := ciotinfra.NewMemoryStore()
	ciotHandler := ciothttp.NewHandler(
		ciotapp.NewCIOTService(ciotStore.CIOTs(), ciotStore.History(), ciotStore.Errors()),
		ciotapp.NewCatalogService(ciotStore.Contracts(), ciotapp.InitContract),
		ciotapp.NewCatalogService(ciotStore.Carriers(), ciotapp.InitCarrier),
		ciotapp.NewCatalogService(ciotStore.Transporters(), ciotapp.InitTransporter),
		ciotapp.NewCatalogService(ciotStore.Operations(), ciotapp.InitOperation),
		ciotapp.NewCatalogService(ciotStore.Vehicles(), ciotapp.InitVehicleReference),
		ciotapp.NewCatalogService(ciotStore.Drivers(), ciotapp.InitDriverReference),
		ciotapp.NewCatalogService(ciotStore.Amounts(), ciotapp.InitAmount),
		ciotapp.NewStatusHistoryService(ciotStore.History()),
		ciotapp.NewPaymentService(ciotStore.Payments(), ciotStore.History()),
		ciotapp.NewProviderAttemptService(ciotStore.Attempts(), ciotStore.History()),
		ciotapp.NewExternalReferenceService(ciotStore.ExternalReferences()),
		ciotapp.NewDocumentService(ciotStore.Documents()),
		ciotapp.NewErrorService(ciotStore.Errors()),
	)

	maintenanceStore := maintenanceinfra.NewMemoryStore()
	maintenanceHandler := maintenancehttp.NewHandler(
		maintenanceapp.NewWorkOrderService(maintenanceStore.WorkOrders(), maintenanceStore.History()),
		maintenanceapp.NewPreventivePlanService(maintenanceStore.PreventivePlans()),
		maintenanceapp.NewServiceTypeService(maintenanceStore.ServiceTypes()),
		maintenanceapp.NewCatalogService(maintenanceStore.Categories()),
		maintenanceapp.NewCatalogService(maintenanceStore.Types()),
		maintenanceapp.NewCatalogService(maintenanceStore.Priorities()),
		maintenanceapp.NewCatalogService(maintenanceStore.Reasons()),
		maintenanceapp.NewWorkshopService(maintenanceStore.Workshops()),
		maintenanceapp.NewTechnicianService(maintenanceStore.Technicians()),
		maintenanceapp.NewLaborService(maintenanceStore.Labor(), maintenanceStore.WorkOrders()),
		maintenanceapp.NewDowntimeService(maintenanceStore.Downtime(), maintenanceStore.WorkOrders()),
		maintenanceapp.NewHistoryService(maintenanceStore.History()),
	)

	inventoryStore := inventoryinfra.NewMemoryStore()
	inventoryHandler := inventoryhttp.NewHandler(
		inventoryapp.NewPartService(inventoryStore.Parts(), inventoryStore.Units()),
		inventoryapp.NewCatalogService(inventoryStore.Categories()),
		inventoryapp.NewCatalogService(inventoryStore.Units()),
		inventoryapp.NewWarehouseService(inventoryStore.Warehouses()),
		inventoryapp.NewLocationService(inventoryStore.Locations()),
	)

	suppliersStore := suppliersinfra.NewMemoryStore()
	suppliersHandler := suppliershttp.NewHandler(
		suppliersapp.NewSupplierService(suppliersStore.Suppliers(), suppliersStore.Categories(), suppliersStore.Audit()),
		suppliersapp.NewCategoryService(suppliersStore.Categories()),
		suppliersapp.NewTypeService(suppliersStore.Types()),
		suppliersapp.NewContactService(suppliersStore.Contacts(), suppliersStore.Suppliers()),
		suppliersapp.NewAddressService(suppliersStore.Addresses(), suppliersStore.Suppliers()),
		suppliersapp.NewDocumentService(suppliersStore.Documents(), suppliersStore.Suppliers(), suppliersStore.Audit()),
		suppliersapp.NewContractService(suppliersStore.Contracts(), suppliersStore.Suppliers(), suppliersStore.Audit()),
		suppliersapp.NewRatingService(suppliersStore.Ratings(), suppliersStore.Suppliers()),
	)

	tiresStore := tiresinfra.NewMemoryStore()
	tiresHandler := tireshttp.NewHandler(
		tiresapp.NewTireLifecycleService(tiresStore.Tires(), tiresStore.Movements(), tiresStore.History()),
		tiresapp.NewTireInspectionService(tiresStore.Tires(), tiresStore.Inspections()),
		tiresapp.NewTireMovementService(tiresStore.Movements()),
		tiresapp.NewTireHistoryService(tiresStore.History()),
	)

	recommendationService := intelligenceapp.NewRecommendationService()
	intelligenceHandler := intelligencehttp.NewHandler(recommendationService)

	return corehttp.NewRouter(corehttp.RouterDependencies{
		Config:              c.config,
		Logger:              c.logger,
		AssetsHandler:       assetsHandler,
		IdentityHandler:     identityHandler,
		AuthService:         authService,
		ChecklistHandler:    checklistHandler,
		FleetHandler:        fleetHandler,
		FuelHandler:         fuelHandler,
		FinancialHandler:    financialHandler,
		CIOTHandler:         ciotHandler,
		InventoryHandler:    inventoryHandler,
		MaintenanceHandler:  maintenanceHandler,
		SuppliersHandler:    suppliersHandler,
		TiresHandler:        tiresHandler,
		IntelligenceHandler: intelligenceHandler,
	})
}

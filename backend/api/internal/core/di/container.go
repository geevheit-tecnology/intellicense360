package di

import (
	"database/sql"
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
	documentsapp "github.com/geevheit/intelligence360/backend/api/internal/modules/documents/application"
	documentspg "github.com/geevheit/intelligence360/backend/api/internal/modules/documents/infrastructure/postgres"
	documentshttp "github.com/geevheit/intelligence360/backend/api/internal/modules/documents/transport/http"
	driversapp "github.com/geevheit/intelligence360/backend/api/internal/modules/drivers/application"
	driverspg "github.com/geevheit/intelligence360/backend/api/internal/modules/drivers/infrastructure/postgres"
	drivershttp "github.com/geevheit/intelligence360/backend/api/internal/modules/drivers/transport/http"
	financialapp "github.com/geevheit/intelligence360/backend/api/internal/modules/financial/application"
	financialinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/financial/infrastructure"
	financialhttp "github.com/geevheit/intelligence360/backend/api/internal/modules/financial/transport/http"
	fleetapp "github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/application"
	fleetinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/infrastructure"
	fleetpg "github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/infrastructure/postgres"
	fleetports "github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/ports"
	fleethttp "github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/transport/http"
	fuelapp "github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/application"
	fuelinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/infrastructure"
	fuelhttp "github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/transport/http"
	identityapp "github.com/geevheit/intelligence360/backend/api/internal/modules/identity/application"
	identityinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/identity/infrastructure"
	identitypg "github.com/geevheit/intelligence360/backend/api/internal/modules/identity/infrastructure/postgres"
	identityports "github.com/geevheit/intelligence360/backend/api/internal/modules/identity/ports"
	identityhttp "github.com/geevheit/intelligence360/backend/api/internal/modules/identity/transport/http"
	intelligenceapp "github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/application"
	intelligenceinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/infrastructure"
	intelligencehttp "github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/transport/http"
	inventoryapp "github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/application"
	inventoryinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/infrastructure"
	inventoryhttp "github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/transport/http"
	maintenanceapp "github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/application"
	maintenanceinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/infrastructure"
	maintenancehttp "github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/transport/http"
	missionapp "github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/application"
	missioninfra "github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/infrastructure"
	missionhttp "github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/transport/http"
	suppliersapp "github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/application"
	suppliersinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/infrastructure"
	suppliershttp "github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/transport/http"
	tiresapp "github.com/geevheit/intelligence360/backend/api/internal/modules/tires/application"
	tiresinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/tires/infrastructure"
	tireshttp "github.com/geevheit/intelligence360/backend/api/internal/modules/tires/transport/http"
	platform "github.com/geevheit/intelligence360/backend/api/internal/platform/postgres"
	"github.com/gin-gonic/gin"
)

type Container struct {
	db     *sql.DB
	config config.Config
	logger logger.Logger
}

func NewContainer(cfg config.Config) *Container {
	return &Container{
		config: cfg,
		logger: logger.New(cfg.Env),
	}
}

func (c *Container) WithDatabase(db *sql.DB) *Container { c.db = db; return c }

func (c *Container) HTTPRouter() *gin.Engine {
	var users identityports.UserRepository
	var roles identityports.RoleRepository
	var permissions identityports.PermissionRepository
	var tenants identityports.TenantRepository
	var sessions identityports.SessionRepository
	var refreshTokens identityports.RefreshTokenRepository
	var audit identityports.AuditRepository
	if c.db != nil {
		store := identitypg.NewStore(c.db)
		users, roles, permissions, tenants, sessions, refreshTokens, audit = store.Users(), store.Roles(), store.Permissions(), store.Tenants(), store.Sessions(), store.RefreshTokens(), store.Audit()
	} else {
		if c.config.StorageDriver != "memory" {
			panic("postgres storage requires a database")
		}
		store := identityinfra.NewMemoryStore()
		users, roles, permissions, tenants, sessions, refreshTokens, audit = store.Users(), store.Roles(), store.Permissions(), store.Tenants(), store.Sessions(), store.RefreshTokens(), store.Audit()
	}
	tokenService := identityapp.NewTokenService(c.config.AuthIssuer, c.config.AuthSecret)
	authService := identityapp.NewAuthService(users, roles, sessions, refreshTokens, tokenService)
	if c.db != nil {
		authService = authService.WithTransactions(platform.Transactions{DB: c.db})
	}
	userService := identityapp.NewUserService(users)
	roleService := identityapp.NewRoleService(roles)
	permissionService := identityapp.NewPermissionService(permissions)
	tenantService := identityapp.NewTenantService(tenants)
	auditService := identityapp.NewAuditService(audit)
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

	var vehicles fleetports.VehicleRepository
	var brands fleetports.VehicleBrandRepository
	var models fleetports.VehicleModelRepository
	var categories fleetports.VehicleCategoryRepository
	var types fleetports.VehicleTypeRepository
	var fleetAssets fleetports.AssetRepository
	if c.db != nil {
		store := fleetpg.NewStore(c.db)
		vehicles, brands, models, categories, types, fleetAssets = store.Vehicles(), store.Brands(), store.Models(), store.Categories(), store.Types(), store.Assets()
	} else {
		store := fleetinfra.NewMemoryStore()
		vehicles, brands, models, categories, types, fleetAssets = store.Vehicles(), store.Brands(), store.Models(), store.Categories(), store.Types(), store.Assets()
	}
	fleetService := fleetapp.NewVehicleService(vehicles)
	fleetHandler := fleethttp.NewHandler(
		fleetService,
		fleetapp.NewVehicleBrandService(brands),
		fleetapp.NewVehicleModelService(models),
		fleetapp.NewVehicleCategoryService(categories),
		fleetapp.NewVehicleTypeService(types),
		fleetapp.NewAssetService(fleetAssets),
	)

	var driversHandler *drivershttp.Handler
	var documentsHandler *documentshttp.Handler
	if c.db != nil {
		tx := platform.Transactions{DB: c.db}
		driversService := driversapp.NewService(driverspg.NewRepository(c.db), driverVehicleAccess(fleetService), tx, driverAudit{service: auditService})
		driversHandler = drivershttp.NewHandler(driversService)
		documentsService := documentsapp.NewService(documentspg.NewRepository(c.db), documentSubjects{drivers: driversService, fleet: fleetService}, tx, documentAudit{service: auditService})
		documentsHandler = documentshttp.NewHandler(documentsService)
	}

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

	intelligenceStore := intelligenceinfra.NewMemoryStore()
	intelligenceCore := intelligenceapp.NewIntelligenceService(
		intelligenceStore.Anomalies(),
		intelligenceStore.Risks(),
		intelligenceStore.Opportunities(),
		intelligenceStore.Recommendations(),
		intelligenceStore.Insights(),
	)
	intelligenceHandler := intelligencehttp.NewHandler(
		intelligenceapp.NewMetricService(intelligenceStore.Metrics()),
		intelligenceCore,
	)

	missionStore := missioninfra.NewMemoryStore()
	missionHandler := missionhttp.NewHandler(
		missionapp.NewCommandService(missionStore.Items(), missionStore.Events(), missionStore.Actions(), missionStore.Snapshots(), missionStore.Idempotency()),
	)

	return corehttp.NewRouter(corehttp.RouterDependencies{
		Config:                c.config,
		DriversHandler:        driversHandler,
		DocumentsHandler:      documentsHandler,
		Logger:                c.logger,
		AssetsHandler:         assetsHandler,
		IdentityHandler:       identityHandler,
		AuthService:           authService,
		ChecklistHandler:      checklistHandler,
		FleetHandler:          fleetHandler,
		FuelHandler:           fuelHandler,
		FinancialHandler:      financialHandler,
		CIOTHandler:           ciotHandler,
		InventoryHandler:      inventoryHandler,
		MaintenanceHandler:    maintenanceHandler,
		SuppliersHandler:      suppliersHandler,
		TiresHandler:          tiresHandler,
		IntelligenceHandler:   intelligenceHandler,
		MissionControlHandler: missionHandler,
	})
}

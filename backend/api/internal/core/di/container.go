package di

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/geevheit/intelligence360/backend/api/internal/core/config"
	corehttp "github.com/geevheit/intelligence360/backend/api/internal/core/http"
	"github.com/geevheit/intelligence360/backend/api/internal/core/logger"
	assetsapp "github.com/geevheit/intelligence360/backend/api/internal/modules/assets/application"
	assetsinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/assets/infrastructure"
	assetshttp "github.com/geevheit/intelligence360/backend/api/internal/modules/assets/transport/http"
	checklistapp "github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/application"
	checklistinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/infrastructure"
	checklistpg "github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/infrastructure/postgres"
	checklistports "github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/ports"
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
	financialdomain "github.com/geevheit/intelligence360/backend/api/internal/modules/financial/domain"
	financialinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/financial/infrastructure"
	financialpg "github.com/geevheit/intelligence360/backend/api/internal/modules/financial/infrastructure/postgres"
	financialports "github.com/geevheit/intelligence360/backend/api/internal/modules/financial/ports"
	financialhttp "github.com/geevheit/intelligence360/backend/api/internal/modules/financial/transport/http"
	fleetapp "github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/application"
	fleetinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/infrastructure"
	fleetpg "github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/infrastructure/postgres"
	fleetports "github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/ports"
	fleethttp "github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/transport/http"
	fuelapp "github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/application"
	fuelinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/infrastructure"
	fuelpg "github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/infrastructure/postgres"
	fuelports "github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/ports"
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
	inventorypg "github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/infrastructure/postgres"
	inventoryports "github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/ports"
	inventoryhttp "github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/transport/http"
	maintenanceapp "github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/application"
	maintenanceinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/infrastructure"
	maintenancepg "github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/infrastructure/postgres"
	maintenanceports "github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/ports"
	maintenancehttp "github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/transport/http"
	missionapp "github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/application"
	missioninfra "github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/infrastructure"
	missionpg "github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/infrastructure/postgres"
	missionports "github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/ports"
	missionhttp "github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/transport/http"
	suppliersapp "github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/application"
	suppliersinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/infrastructure"
	suppliershttp "github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/transport/http"
	tiresapp "github.com/geevheit/intelligence360/backend/api/internal/modules/tires/application"
	tiresinfra "github.com/geevheit/intelligence360/backend/api/internal/modules/tires/infrastructure"
	tirespg "github.com/geevheit/intelligence360/backend/api/internal/modules/tires/infrastructure/postgres"
	tiresports "github.com/geevheit/intelligence360/backend/api/internal/modules/tires/ports"
	tireshttp "github.com/geevheit/intelligence360/backend/api/internal/modules/tires/transport/http"
	tripsapp "github.com/geevheit/intelligence360/backend/api/internal/modules/trips/application"
	tripspg "github.com/geevheit/intelligence360/backend/api/internal/modules/trips/infrastructure/postgres"
	tripshttp "github.com/geevheit/intelligence360/backend/api/internal/modules/trips/transport/http"
	platform "github.com/geevheit/intelligence360/backend/api/internal/platform/postgres"
)

type Container struct {
	db             *sql.DB
	config         config.Config
	logger         logger.Logger
	commandService missionapp.CommandService
}

func NewContainer(cfg config.Config) *Container {
	return &Container{
		config: cfg,
		logger: logger.New(cfg.Env),
	}
}

func (c *Container) WithDatabase(db *sql.DB) *Container { c.db = db; return c }

// StartAutonomousScanner launches the 24/7 background AI signal evaluation engine.
func (c *Container) StartAutonomousScanner(ctx context.Context, interval time.Duration) {
	if c.db == nil {
		return
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				rows, err := c.db.QueryContext(ctx, "SELECT id FROM tenants WHERE deleted_at IS NULL")
				if err != nil {
					continue
				}
				var tenantIDs []string
				for rows.Next() {
					var id string
					if err := rows.Scan(&id); err == nil {
						tenantIDs = append(tenantIDs, id)
					}
				}
				rows.Close()

				for _, tid := range tenantIDs {
					count, err := c.commandService.EvaluateSignals(ctx, tid)
					if err == nil && count > 0 {
						c.logger.Info(ctx, "motor autônomo de inteligência detectou novos alertas operacionais",
							slog.String("tenant_id", tid),
							slog.Int("alert_count", count),
						)
					}
				}
			}
		}
	}()
}

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

	var clChecklists checklistports.ChecklistRepository
	var clItems checklistports.ChecklistItemRepository
	var clAnswers checklistports.ChecklistAnswerRepository
	var clTemplates checklistports.ChecklistTemplateRepository
	var clVersions checklistports.ChecklistTemplateVersionRepository
	var clTypes checklistports.ChecklistTypeRepository
	var clSections checklistports.ChecklistSectionRepository
	var clEngineItems checklistports.ChecklistEngineItemRepository
	var clExecutions checklistports.ChecklistExecutionRepository
	var clResponses checklistports.ChecklistResponseRepository
	var clEvidence checklistports.ChecklistEvidenceRepository
	var clNonConformities checklistports.ChecklistNonConformityRepository
	var clSignatures checklistports.ChecklistSignatureRepository
	var clHistory checklistports.ChecklistHistoryRepository

	if c.db != nil {
		store := checklistpg.NewStore(c.db)
		clChecklists = store.Checklists()
		clItems = store.Items()
		clAnswers = store.Answers()
		clTemplates = store.Templates()
		clVersions = store.Versions()
		clTypes = store.Types()
		clSections = store.Sections()
		clEngineItems = store.EngineItems()
		clExecutions = store.Executions()
		clResponses = store.Responses()
		clEvidence = store.Evidence()
		clNonConformities = store.NonConformities()
		clSignatures = store.Signatures()
		clHistory = store.History()
	} else {
		store := checklistinfra.NewMemoryStore()
		clChecklists = store.Checklists()
		clItems = store.Items()
		clAnswers = store.Answers()
		clTemplates = store.Templates()
		clVersions = store.Versions()
		clTypes = store.Types()
		clSections = store.Sections()
		clEngineItems = store.EngineItems()
		clExecutions = store.Executions()
		clResponses = store.Responses()
		clEvidence = store.Evidence()
		clNonConformities = store.NonConformities()
		clSignatures = store.Signatures()
		clHistory = store.History()
	}

	checklistHandler := checklisthttp.NewHandler(
		checklistapp.NewChecklistService(clChecklists),
		checklistapp.NewChecklistItemService(clChecklists, clItems),
		checklistapp.NewChecklistAnswerService(clItems, clAnswers),
		checklistapp.NewChecklistTemplateService(clTemplates),
		checklistapp.NewChecklistTemplateVersionService(clTemplates, clVersions),
		checklistapp.NewChecklistTypeService(clTypes),
		checklistapp.NewChecklistSectionService(clSections),
		checklistapp.NewChecklistEngineItemService(clVersions, clEngineItems),
		checklistapp.NewChecklistExecutionService(clVersions, clEngineItems, clExecutions, clResponses, clEvidence, clSignatures, clHistory),
		checklistapp.NewChecklistEngineResponseService(clExecutions, clEngineItems, clResponses, clHistory),
		checklistapp.NewChecklistEvidenceService(clEvidence, clHistory),
		checklistapp.NewChecklistNonConformityService(clNonConformities, clHistory),
		checklistapp.NewChecklistHistoryService(clHistory),
	)

	var vehicles fleetports.VehicleRepository
	var brands fleetports.VehicleBrandRepository
	var models fleetports.VehicleModelRepository
	var categories fleetports.VehicleCategoryRepository
	var types fleetports.VehicleTypeRepository
	var fleetAssets fleetports.AssetRepository
	var odometerRepo fleetports.OdometerRepository
	if c.db != nil {
		store := fleetpg.NewStore(c.db)
		vehicles, brands, models, categories, types, fleetAssets, odometerRepo = store.Vehicles(), store.Brands(), store.Models(), store.Categories(), store.Types(), store.Assets(), store.Odometer()
	} else {
		store := fleetinfra.NewMemoryStore()
		vehicles, brands, models, categories, types, fleetAssets, odometerRepo = store.Vehicles(), store.Brands(), store.Models(), store.Categories(), store.Types(), store.Assets(), store.Odometer()
	}
	fleetService := fleetapp.NewVehicleService(vehicles)
	odometerService := fleetapp.NewOdometerService(odometerRepo, vehicles)
	fleetHandler := fleethttp.NewHandler(
		fleetService,
		fleetapp.NewVehicleBrandService(brands),
		fleetapp.NewVehicleModelService(models),
		fleetapp.NewVehicleCategoryService(categories),
		fleetapp.NewVehicleTypeService(types),
		fleetapp.NewAssetService(fleetAssets),
		odometerService,
	)

	var driversHandler *drivershttp.Handler
	var documentsHandler *documentshttp.Handler
	if c.db != nil {
		tx := platform.Transactions{DB: c.db}
		driversService := driversapp.NewService(driverspg.NewRepository(c.db), driverVehicleAccess(fleetService), tx, driverAudit{service: auditService})
		journeyRepo := driverspg.NewJourneyRepository(c.db)
		journeyService := driversapp.NewJourneyService(journeyRepo)
		driversHandler = drivershttp.NewHandler(driversService, journeyService)
		documentsService := documentsapp.NewService(documentspg.NewRepository(c.db), documentSubjects{drivers: driversService, fleet: fleetService}, tx, documentAudit{service: auditService})
		documentsHandler = documentshttp.NewHandler(documentsService)
	}

	var fuelTransactions fuelports.FuelTransactionRepository
	var fuelTypes fuelports.FuelTypeRepository
	var fuelStations fuelports.FuelStationRepository
	var fuelTanks fuelports.FuelTankRepository
	var fuelNozzles fuelports.FuelNozzleRepository
	var fuelReadings fuelports.FuelReadingRepository
	var fuelPrices fuelports.FuelPriceRepository
	var fuelReceipts fuelports.FuelReceiptRepository
	var fuelAdjustments fuelports.FuelAdjustmentRepository

	if c.db != nil {
		store := fuelpg.NewStore(c.db)
		fuelTransactions = store.Transactions()
		fuelTypes = store.Types()
		fuelStations = store.Stations()
		fuelTanks = store.Tanks()
		fuelNozzles = store.Nozzles()
		fuelReadings = store.Readings()
		fuelPrices = store.Prices()
		fuelReceipts = store.Receipts()
		fuelAdjustments = store.Adjustments()
	} else {
		store := fuelinfra.NewMemoryStore()
		fuelTransactions = store.Transactions()
		fuelTypes = store.Types()
		fuelStations = store.Stations()
		fuelTanks = store.Tanks()
		fuelNozzles = store.Nozzles()
		fuelReadings = store.Readings()
		fuelPrices = store.Prices()
		fuelReceipts = store.Receipts()
		fuelAdjustments = store.Adjustments()
	}

	fuelHandler := fuelhttp.NewHandler(
		fuelapp.NewFuelTransactionService(fuelTransactions, fuelAdjustments),
		fuelapp.NewFuelTypeService(fuelTypes),
		fuelapp.NewFuelStationService(fuelStations),
		fuelapp.NewFuelTankService(fuelTanks),
		fuelapp.NewFuelNozzleService(fuelNozzles),
		fuelapp.NewFuelReadingService(fuelReadings),
		fuelapp.NewFuelPriceService(fuelPrices),
		fuelapp.NewFuelReceiptService(fuelReceipts),
		fuelapp.NewFuelAdjustmentService(fuelAdjustments),
	)

	var finTransactions financialports.TransactionRepository
	var finCategories financialports.CatalogRepository[financialdomain.CostCategory]
	var finTypes financialports.CatalogRepository[financialdomain.CostType]
	var finCenters financialports.CatalogRepository[financialdomain.CostCenter]
	var finAccounts financialports.CatalogRepository[financialdomain.Account]
	var finMethods financialports.CatalogRepository[financialdomain.PaymentMethod]
	var finPeriods financialports.PeriodRepository
	var finBudgets financialports.BudgetRepository
	var finAdjustments financialports.AdjustmentRepository
	var finClosings financialports.CashClosingRepository

	if c.config.StorageDriver == "postgres" && c.db != nil {
		pgStore := financialpg.NewStore(c.db)
		finTransactions = pgStore.Transactions()
		finCategories = pgStore.Categories()
		finTypes = pgStore.Types()
		finCenters = pgStore.Centers()
		finAccounts = pgStore.Accounts()
		finMethods = pgStore.Methods()
		finPeriods = pgStore.Periods()
		finBudgets = pgStore.Budgets()
		finAdjustments = pgStore.Adjustments()
		finClosings = pgStore.Closings()
	} else {
		memStore := financialinfra.NewMemoryStore()
		finTransactions = memStore.Transactions()
		finCategories = memStore.Categories()
		finTypes = memStore.Types()
		finCenters = memStore.Centers()
		finAccounts = memStore.Accounts()
		finMethods = memStore.Methods()
		finPeriods = memStore.Periods()
		finBudgets = memStore.Budgets()
		finAdjustments = memStore.Adjustments()
		finClosings = memStore.Closings()
	}

	financialHandler := financialhttp.NewHandler(
		financialapp.NewTransactionService(finTransactions, finPeriods, finAdjustments),
		financialapp.NewCatalogService(finCategories, financialapp.InitCategory),
		financialapp.NewCatalogService(finTypes, financialapp.InitType),
		financialapp.NewCatalogService(finCenters, financialapp.InitCenter),
		financialapp.NewCatalogService(finAccounts, financialapp.InitAccount),
		financialapp.NewCatalogService(finMethods, financialapp.InitPaymentMethod),
		financialapp.NewPeriodService(finPeriods),
		financialapp.NewBudgetService(finBudgets),
		financialapp.NewAdjustmentService(finAdjustments),
		financialapp.NewCashFlowService(finTransactions, finClosings),
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

	var workOrders maintenanceports.WorkOrderRepository
	var plans maintenanceports.PreventivePlanRepository
	var serviceTypes maintenanceports.ServiceTypeRepository
	var maintCategories maintenanceports.CatalogRepository
	var maintTypes maintenanceports.CatalogRepository
	var maintPriorities maintenanceports.CatalogRepository
	var maintReasons maintenanceports.CatalogRepository
	var workshops maintenanceports.WorkshopRepository
	var technicians maintenanceports.TechnicianRepository
	var labor maintenanceports.LaborRepository
	var downtime maintenanceports.DowntimeRepository
	var history maintenanceports.HistoryRepository

	if c.db != nil {
		store := maintenancepg.NewStore(c.db)
		workOrders = store.WorkOrders()
		plans = store.PreventivePlans()
		serviceTypes = store.ServiceTypes()
		maintCategories = store.Categories()
		maintTypes = store.Types()
		maintPriorities = store.Priorities()
		maintReasons = store.Reasons()
		workshops = store.Workshops()
		technicians = store.Technicians()
		labor = store.Labor()
		downtime = store.Downtime()
		history = store.History()
	} else {
		store := maintenanceinfra.NewMemoryStore()
		workOrders = store.WorkOrders()
		plans = store.PreventivePlans()
		serviceTypes = store.ServiceTypes()
		maintCategories = store.Categories()
		maintTypes = store.Types()
		maintPriorities = store.Priorities()
		maintReasons = store.Reasons()
		workshops = store.Workshops()
		technicians = store.Technicians()
		labor = store.Labor()
		downtime = store.Downtime()
		history = store.History()
	}

	maintenanceHandler := maintenancehttp.NewHandler(
		maintenanceapp.NewWorkOrderService(workOrders, history),
		maintenanceapp.NewPreventivePlanService(plans),
		maintenanceapp.NewServiceTypeService(serviceTypes),
		maintenanceapp.NewCatalogService(maintCategories),
		maintenanceapp.NewCatalogService(maintTypes),
		maintenanceapp.NewCatalogService(maintPriorities),
		maintenanceapp.NewCatalogService(maintReasons),
		maintenanceapp.NewWorkshopService(workshops),
		maintenanceapp.NewTechnicianService(technicians),
		maintenanceapp.NewLaborService(labor, workOrders),
		maintenanceapp.NewDowntimeService(downtime, workOrders),
		maintenanceapp.NewHistoryService(history),
	)

	var invParts inventoryports.PartRepository
	var invCategories inventoryports.CatalogRepository
	var invUnits inventoryports.CatalogRepository
	var invWarehouses inventoryports.WarehouseRepository
	var invLocations inventoryports.LocationRepository

	if c.config.StorageDriver == "postgres" && c.db != nil {
		pgStore := inventorypg.NewStore(c.db)
		invParts = pgStore.Parts()
		invCategories = pgStore.Categories()
		invUnits = pgStore.Units()
		invWarehouses = pgStore.Warehouses()
		invLocations = pgStore.Locations()
	} else {
		memStore := inventoryinfra.NewMemoryStore()
		invParts = memStore.Parts()
		invCategories = memStore.Categories()
		invUnits = memStore.Units()
		invWarehouses = memStore.Warehouses()
		invLocations = memStore.Locations()
	}

	inventoryHandler := inventoryhttp.NewHandler(
		inventoryapp.NewPartService(invParts, invUnits),
		inventoryapp.NewCatalogService(invCategories),
		inventoryapp.NewCatalogService(invUnits),
		inventoryapp.NewWarehouseService(invWarehouses),
		inventoryapp.NewLocationService(invLocations),
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

	var tiresRepo tiresports.TireRepository
	var tireInspections tiresports.TireInspectionRepository
	var tireMovements tiresports.TireMovementRepository
	var tireHistory tiresports.TireHistoryRepository

	if c.db != nil {
		store := tirespg.NewStore(c.db)
		tiresRepo = store.Tires()
		tireInspections = store.Inspections()
		tireMovements = store.Movements()
		tireHistory = store.History()
	} else {
		store := tiresinfra.NewMemoryStore()
		tiresRepo = store.Tires()
		tireInspections = store.Inspections()
		tireMovements = store.Movements()
		tireHistory = store.History()
	}

	tiresHandler := tireshttp.NewHandler(
		tiresapp.NewTireLifecycleService(tiresRepo, tireMovements, tireHistory),
		tiresapp.NewTireInspectionService(tiresRepo, tireInspections),
		tiresapp.NewTireMovementService(tireMovements),
		tiresapp.NewTireHistoryService(tireHistory),
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

	var mcItems missionports.CommandItemRepository
	var mcEvents missionports.CommandEventRepository
	var mcActions missionports.CommandActionRepository
	var mcSnapshots missionports.OperationalSnapshotRepository
	var mcIdempotency missionports.IdempotencyRepository
	var mcScanner missionports.SignalScanner

	if c.db != nil {
		store := missionpg.NewStore(c.db)
		mcItems = store.Items()
		mcEvents = store.Events()
		mcActions = store.Actions()
		mcSnapshots = store.Snapshots()
		mcIdempotency = store.Idempotency()
		mcScanner = missionpg.NewSignalScanner(c.db)
	} else {
		store := missioninfra.NewMemoryStore()
		mcItems = store.Items()
		mcEvents = store.Events()
		mcActions = store.Actions()
		mcSnapshots = store.Snapshots()
		mcIdempotency = store.Idempotency()
	}

	commandService := missionapp.NewCommandService(mcItems, mcEvents, mcActions, mcSnapshots, mcIdempotency)
	if mcScanner != nil {
		commandService = commandService.WithSignalScanner(mcScanner)
	}
	c.commandService = commandService
	missionHandler := missionhttp.NewHandler(commandService)

	var tripsHandler *tripshttp.Handler
	if c.db != nil {
		tripsStore := tripspg.NewStore(c.db)
		tripsSvc := tripsapp.NewTripService(tripsStore.Trips())
		h := tripshttp.NewHandler(tripsSvc)
		tripsHandler = &h
	}

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
		TripsHandler:          tripsHandler,
	})
}

package http

import (
	stdhttp "net/http"
	"time"

	"github.com/geevheit/intelligence360/backend/api/internal/core/config"
	"github.com/geevheit/intelligence360/backend/api/internal/core/logger"
	"github.com/geevheit/intelligence360/backend/api/internal/core/middleware"
	assetshttp "github.com/geevheit/intelligence360/backend/api/internal/modules/assets/transport/http"
	checklisthttp "github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/transport/http"
	ciothttp "github.com/geevheit/intelligence360/backend/api/internal/modules/ciot/transport/http"
	financialhttp "github.com/geevheit/intelligence360/backend/api/internal/modules/financial/transport/http"
	fleethttp "github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/transport/http"
	fuelhttp "github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/transport/http"
	identityports "github.com/geevheit/intelligence360/backend/api/internal/modules/identity/ports"
	identityhttp "github.com/geevheit/intelligence360/backend/api/internal/modules/identity/transport/http"
	intelligencehttp "github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/transport/http"
	inventoryhttp "github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/transport/http"
	maintenancehttp "github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/transport/http"
	suppliershttp "github.com/geevheit/intelligence360/backend/api/internal/modules/suppliers/transport/http"
	tireshttp "github.com/geevheit/intelligence360/backend/api/internal/modules/tires/transport/http"
	"github.com/gin-gonic/gin"
)

type RouterDependencies struct {
	Config              config.Config
	Logger              logger.Logger
	AssetsHandler       assetshttp.Handler
	IdentityHandler     identityhttp.Handler
	AuthService         identityports.AuthenticationService
	ChecklistHandler    checklisthttp.Handler
	FleetHandler        fleethttp.Handler
	FuelHandler         fuelhttp.Handler
	FinancialHandler    financialhttp.Handler
	CIOTHandler         ciothttp.Handler
	InventoryHandler    inventoryhttp.Handler
	MaintenanceHandler  maintenancehttp.Handler
	SuppliersHandler    suppliershttp.Handler
	TiresHandler        tireshttp.Handler
	IntelligenceHandler intelligencehttp.Handler
}

func NewRouter(deps RouterDependencies) *gin.Engine {
	if deps.Config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	_ = router.SetTrustedProxies(nil)

	router.Use(
		gin.Recovery(),
		middleware.RequestLogger(deps.Logger),
		middleware.ErrorHandler(deps.Logger),
		middleware.Audit(deps.Logger),
	)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(stdhttp.StatusOK, gin.H{
			"status":    "ok",
			"service":   deps.Config.AppName,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	v1 := router.Group("/api/v1")
	deps.IdentityHandler.RegisterPublicRoutes(v1)

	protected := v1.Group("")
	protected.Use(middleware.Authentication(deps.AuthService), middleware.Tenant())
	{
		deps.AssetsHandler.RegisterRoutes(protected)
		deps.IdentityHandler.RegisterProtectedRoutes(protected)
		deps.ChecklistHandler.RegisterRoutes(protected)
		deps.FleetHandler.RegisterRoutes(protected)
		deps.FuelHandler.RegisterRoutes(protected)
		deps.FinancialHandler.RegisterRoutes(protected)
		deps.CIOTHandler.RegisterRoutes(protected)
		deps.InventoryHandler.RegisterRoutes(protected)
		deps.MaintenanceHandler.RegisterRoutes(protected)
		deps.SuppliersHandler.RegisterRoutes(protected)
		deps.TiresHandler.RegisterRoutes(protected)
		deps.IntelligenceHandler.RegisterRoutes(protected)
	}

	return router
}

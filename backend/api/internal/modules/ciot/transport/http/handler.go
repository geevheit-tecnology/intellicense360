package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/geevheit/intelligence360/backend/api/internal/core/middleware"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/ciot/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/ciot/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/ciot/ports"
	"github.com/geevheit/intelligence360/backend/api/internal/shared/contextkeys"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	ciots        ports.CIOTService
	contracts    ports.CatalogService[domain.CIOTContract]
	carriers     ports.CatalogService[domain.CIOTCarrier]
	transporters ports.CatalogService[domain.CIOTTransporter]
	operations   ports.CatalogService[domain.CIOTOperation]
	vehicles     ports.CatalogService[domain.CIOTVehicleReference]
	drivers      ports.CatalogService[domain.CIOTDriverReference]
	amounts      ports.CatalogService[domain.CIOTAmount]
	history      ports.StatusHistoryService
	payments     ports.PaymentService
	attempts     ports.ProviderAttemptService
	references   ports.ExternalReferenceService
	documents    ports.DocumentService
	errors       ports.ErrorService
}

func NewHandler(ciots ports.CIOTService, contracts ports.CatalogService[domain.CIOTContract], carriers ports.CatalogService[domain.CIOTCarrier], transporters ports.CatalogService[domain.CIOTTransporter], operations ports.CatalogService[domain.CIOTOperation], vehicles ports.CatalogService[domain.CIOTVehicleReference], drivers ports.CatalogService[domain.CIOTDriverReference], amounts ports.CatalogService[domain.CIOTAmount], history ports.StatusHistoryService, payments ports.PaymentService, attempts ports.ProviderAttemptService, references ports.ExternalReferenceService, documents ports.DocumentService, errors ports.ErrorService) Handler {
	return Handler{ciots: ciots, contracts: contracts, carriers: carriers, transporters: transporters, operations: operations, vehicles: vehicles, drivers: drivers, amounts: amounts, history: history, payments: payments, attempts: attempts, references: references, documents: documents, errors: errors}
}

func (h Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/ciot")
	group.Use(middleware.RequirePermission("ciot.ciot.manage"))
	group.GET("", h.searchCIOT)
	group.POST("", h.createCIOT)
	group.GET("/types", h.types)
	group.GET("/contracts", h.searchContracts)
	group.POST("/contracts", h.createContract)
	group.GET("/carriers", h.searchCarriers)
	group.POST("/carriers", h.createCarrier)
	group.GET("/transporters", h.searchTransporters)
	group.POST("/transporters", h.createTransporter)
	group.GET("/operations", h.searchOperations)
	group.POST("/operations", h.createOperation)
	group.GET("/vehicle-references", h.searchVehicles)
	group.POST("/vehicle-references", h.createVehicle)
	group.GET("/driver-references", h.searchDrivers)
	group.POST("/driver-references", h.createDriver)
	group.GET("/amounts", h.searchAmounts)
	group.POST("/amounts", h.createAmount)
	group.GET("/:id", h.getCIOT)
	group.PATCH("/:id", h.updateCIOT)
	group.DELETE("/:id", h.deleteCIOT)
	group.GET("/:id/history", h.historyByCIOT)
	group.POST("/:id/submit", h.submit)
	group.POST("/:id/generated", h.generated)
	group.POST("/:id/activate", h.activate)
	group.POST("/:id/suspend", h.suspend)
	group.POST("/:id/reactivate", h.reactivate)
	group.POST("/:id/close", h.close)
	group.POST("/:id/cancel", h.cancel)
	group.POST("/:id/error", h.recordError)
	group.POST("/:id/retry", h.retry)
	group.GET("/:id/payments", h.paymentsByCIOT)
	group.POST("/:id/payments", h.createPayment)
	group.GET("/:id/provider-attempts", h.attemptsByCIOT)
	group.POST("/:id/provider-attempts", h.createAttempt)
	group.GET("/:id/external-reference", h.externalReference)
	group.POST("/:id/external-reference", h.upsertExternalReference)
	group.GET("/:id/documents", h.documentsByCIOT)
	group.POST("/:id/documents", h.createDocument)
	group.GET("/:id/errors", h.errorsByCIOT)
}

func (h Handler) createCIOT(c *gin.Context) {
	var req domain.CIOT
	if !bind(c, &req) {
		return
	}
	req.TenantID, req.CreatedBy, req.UpdatedBy = tenantID(c), actorID(c), actorID(c)
	saved, err := h.ciots.Create(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchCIOT(c *gin.Context) {
	page, err := h.ciots.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) getCIOT(c *gin.Context) {
	item, err := h.ciots.FindByID(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) updateCIOT(c *gin.Context) {
	var req domain.CIOT
	if !bind(c, &req) {
		return
	}
	req.ID, req.TenantID, req.UpdatedBy = c.Param("id"), tenantID(c), actorID(c)
	saved, err := h.ciots.Update(c.Request.Context(), req)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteCIOT(c *gin.Context) {
	respond(c, http.StatusNoContent, gin.H{}, h.ciots.Delete(c.Request.Context(), tenantID(c), c.Param("id")))
}
func (h Handler) submit(c *gin.Context) {
	saved, err := h.ciots.Submit(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) generated(c *gin.Context) {
	var req struct {
		ExternalProtocol string `json:"external_protocol"`
	}
	_ = c.ShouldBindJSON(&req)
	saved, err := h.ciots.MarkGenerated(c.Request.Context(), tenantID(c), c.Param("id"), req.ExternalProtocol, actorID(c))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) activate(c *gin.Context) {
	saved, err := h.ciots.Activate(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) suspend(c *gin.Context) {
	saved, err := h.ciots.Suspend(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) reactivate(c *gin.Context) {
	saved, err := h.ciots.Reactivate(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) close(c *gin.Context) {
	saved, err := h.ciots.Close(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) cancel(c *gin.Context) {
	saved, err := h.ciots.Cancel(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) recordError(c *gin.Context) {
	var req struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if !bind(c, &req) {
		return
	}
	saved, err := h.ciots.RecordError(c.Request.Context(), tenantID(c), c.Param("id"), req.Code, req.Message, actorID(c))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) retry(c *gin.Context) {
	saved, err := h.ciots.RetryFromError(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) types(c *gin.Context) {
	respond(c, http.StatusOK, []gin.H{{"code": domain.TypeTACAgregado, "label": "TAC Agregado"}, {"code": domain.TypeTACIndependente, "label": "TAC Independente"}, {"code": domain.TypeOther, "label": "Other"}}, nil)
}

func (h Handler) createContract(c *gin.Context)     { createCatalog(c, h.contracts) }
func (h Handler) searchContracts(c *gin.Context)    { searchCatalog(c, h.contracts) }
func (h Handler) createCarrier(c *gin.Context)      { createCatalog(c, h.carriers) }
func (h Handler) searchCarriers(c *gin.Context)     { searchCatalog(c, h.carriers) }
func (h Handler) createTransporter(c *gin.Context)  { createCatalog(c, h.transporters) }
func (h Handler) searchTransporters(c *gin.Context) { searchCatalog(c, h.transporters) }
func (h Handler) createOperation(c *gin.Context)    { createCatalog(c, h.operations) }
func (h Handler) searchOperations(c *gin.Context)   { searchCatalog(c, h.operations) }
func (h Handler) createVehicle(c *gin.Context)      { createCatalog(c, h.vehicles) }
func (h Handler) searchVehicles(c *gin.Context)     { searchCatalog(c, h.vehicles) }
func (h Handler) createDriver(c *gin.Context)       { createCatalog(c, h.drivers) }
func (h Handler) searchDrivers(c *gin.Context)      { searchCatalog(c, h.drivers) }
func (h Handler) createAmount(c *gin.Context)       { createCatalog(c, h.amounts) }
func (h Handler) searchAmounts(c *gin.Context)      { searchCatalog(c, h.amounts) }

func createCatalog[T any](c *gin.Context, service ports.CatalogService[T]) {
	var req T
	if !bind(c, &req) {
		return
	}
	setTenant(&req, tenantID(c))
	saved, err := service.Create(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func searchCatalog[T any](c *gin.Context, service ports.CatalogService[T]) {
	page, err := service.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}

func (h Handler) historyByCIOT(c *gin.Context) {
	page, err := h.history.ListByCIOT(c.Request.Context(), tenantID(c), c.Param("id"), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) createPayment(c *gin.Context) {
	var req domain.CIOTPayment
	if !bind(c, &req) {
		return
	}
	req.TenantID, req.CIOTID = tenantID(c), c.Param("id")
	saved, err := h.payments.Create(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) paymentsByCIOT(c *gin.Context) {
	page, err := h.payments.ListByCIOT(c.Request.Context(), tenantID(c), c.Param("id"), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) createAttempt(c *gin.Context) {
	var req domain.CIOTProviderAttempt
	if !bind(c, &req) {
		return
	}
	req.TenantID, req.CIOTID = tenantID(c), c.Param("id")
	saved, err := h.attempts.Create(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) attemptsByCIOT(c *gin.Context) {
	page, err := h.attempts.ListByCIOT(c.Request.Context(), tenantID(c), c.Param("id"), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) upsertExternalReference(c *gin.Context) {
	var req domain.CIOTExternalReference
	if !bind(c, &req) {
		return
	}
	req.TenantID, req.CIOTID = tenantID(c), c.Param("id")
	saved, err := h.references.Upsert(c.Request.Context(), req)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) externalReference(c *gin.Context) {
	item, err := h.references.FindByCIOT(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) createDocument(c *gin.Context) {
	var req domain.CIOTDocument
	if !bind(c, &req) {
		return
	}
	req.TenantID, req.CIOTID = tenantID(c), c.Param("id")
	saved, err := h.documents.Create(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) documentsByCIOT(c *gin.Context) {
	page, err := h.documents.ListByCIOT(c.Request.Context(), tenantID(c), c.Param("id"), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) errorsByCIOT(c *gin.Context) {
	page, err := h.errors.ListByCIOT(c.Request.Context(), tenantID(c), c.Param("id"), query(c))
	respond(c, http.StatusOK, page, err)
}

func bind(c *gin.Context, target any) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request", "message": err.Error()})
		return false
	}
	return true
}
func respond(c *gin.Context, status int, payload any, err error) {
	if err != nil {
		code := http.StatusInternalServerError
		name := "internal_error"
		switch {
		case errors.Is(err, application.ErrNotFound):
			code, name = http.StatusNotFound, "not_found"
		case errors.Is(err, application.ErrInvalidData):
			code, name = http.StatusBadRequest, "invalid_data"
		case errors.Is(err, application.ErrInvalidTransition):
			code, name = http.StatusBadRequest, "invalid_transition"
		case errors.Is(err, application.ErrDuplicateRequest):
			code, name = http.StatusConflict, "duplicate_request"
		case errors.Is(err, application.ErrFinalizedImmutable):
			code, name = http.StatusConflict, "finalized_immutable"
		}
		c.JSON(code, gin.H{"error": name, "message": err.Error()})
		return
	}
	if status == http.StatusNoContent {
		c.Status(status)
		return
	}
	c.JSON(status, payload)
}
func tenantID(c *gin.Context) string {
	if value, ok := c.Get(string(contextkeys.TenantID)); ok {
		if text, ok := value.(string); ok {
			return text
		}
	}
	return ""
}
func actorID(c *gin.Context) string {
	if value, ok := c.Get(string(contextkeys.ActorID)); ok {
		if text, ok := value.(string); ok {
			return text
		}
	}
	return ""
}
func query(c *gin.Context) ports.Query {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "25"))
	return ports.Query{Search: c.Query("search"), Page: page, PerPage: perPage, Sort: c.Query("sort"), Filters: map[string]string{"status": c.Query("status"), "type": c.Query("type")}}
}

func setTenant(target any, tenantID string) {
	switch item := target.(type) {
	case *domain.CIOTContract:
		item.TenantID = tenantID
	case *domain.CIOTCarrier:
		item.TenantID = tenantID
	case *domain.CIOTTransporter:
		item.TenantID = tenantID
	case *domain.CIOTOperation:
		item.TenantID = tenantID
	case *domain.CIOTVehicleReference:
		item.TenantID = tenantID
	case *domain.CIOTDriverReference:
		item.TenantID = tenantID
	case *domain.CIOTAmount:
		item.TenantID = tenantID
	}
}

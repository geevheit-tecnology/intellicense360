package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/geevheit/intelligence360/backend/api/internal/core/middleware"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/ports"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/transport/dto"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fuel/transport/mapper"
	"github.com/geevheit/intelligence360/backend/api/internal/shared/contextkeys"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	transactions ports.FuelTransactionService
	types        ports.FuelTypeService
	stations     ports.FuelStationService
	tanks        ports.FuelTankService
	nozzles      ports.FuelNozzleService
	readings     ports.FuelReadingService
	prices       ports.FuelPriceService
	receipts     ports.FuelReceiptService
	adjustments  ports.FuelAdjustmentService
}

func NewHandler(transactions ports.FuelTransactionService, types ports.FuelTypeService, stations ports.FuelStationService, tanks ports.FuelTankService, nozzles ports.FuelNozzleService, readings ports.FuelReadingService, prices ports.FuelPriceService, receipts ports.FuelReceiptService, adjustments ports.FuelAdjustmentService) Handler {
	return Handler{transactions: transactions, types: types, stations: stations, tanks: tanks, nozzles: nozzles, readings: readings, prices: prices, receipts: receipts, adjustments: adjustments}
}

func (h Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/fuel")
	group.Use(middleware.RequirePermission("fuel.fuel.manage"))
	group.GET("/transactions", h.searchTransactions)
	group.POST("/transactions", h.createTransaction)
	group.GET("/transactions/:id", h.getTransaction)
	group.PUT("/transactions/:id", h.updateTransaction)
	group.DELETE("/transactions/:id", h.deleteTransaction)
	group.POST("/transactions/:id/complete", h.completeTransaction)
	group.POST("/transactions/:id/cancel", h.cancelTransaction)
	group.POST("/transactions/:id/adjust", h.adjustTransactionByID)
	group.GET("/types", h.searchTypes)
	group.POST("/types", h.createType)
	group.GET("/types/:id", h.getType)
	group.PUT("/types/:id", h.updateType)
	group.DELETE("/types/:id", h.deleteType)
	group.GET("/stations", h.searchStations)
	group.POST("/stations", h.createStation)
	group.GET("/stations/:id", h.getStation)
	group.PUT("/stations/:id", h.updateStation)
	group.DELETE("/stations/:id", h.deleteStation)
	group.GET("/tanks", h.searchTanks)
	group.POST("/tanks", h.createTank)
	group.PUT("/tanks/:id", h.updateTank)
	group.DELETE("/tanks/:id", h.deleteTank)
	group.GET("/nozzles", h.searchNozzles)
	group.POST("/nozzles", h.createNozzle)
	group.PUT("/nozzles/:id", h.updateNozzle)
	group.DELETE("/nozzles/:id", h.deleteNozzle)
	group.GET("/readings", h.searchReadings)
	group.POST("/readings", h.recordReading)
	group.GET("/prices", h.searchPrices)
	group.POST("/prices", h.recordPrice)
	group.PUT("/prices/:id", h.updatePrice)
	group.DELETE("/prices/:id", h.deletePrice)
	group.GET("/receipts", h.searchReceipts)
	group.POST("/receipts", h.createReceipt)
	group.PUT("/receipts/:id", h.updateReceipt)
	group.DELETE("/receipts/:id", h.deleteReceipt)
	group.GET("/adjustments", h.searchAdjustments)
	group.POST("/adjustments", h.adjustTransaction)
}

func (h Handler) createTransaction(c *gin.Context) {
	var req dto.FuelTransactionRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.TransactionFromRequest(req)
	item.TenantID = tenantID(c)
	item.CreatedBy = actorID(c)
	item.UpdatedBy = actorID(c)
	saved, err := h.transactions.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) updateTransaction(c *gin.Context) {
	var req dto.FuelTransactionRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.TransactionFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	item.UpdatedBy = actorID(c)
	saved, err := h.transactions.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) searchTransactions(c *gin.Context) {
	page, err := h.transactions.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) getTransaction(c *gin.Context) {
	item, err := h.transactions.FindByID(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) deleteTransaction(c *gin.Context) {
	if err := h.transactions.Delete(c.Request.Context(), tenantID(c), c.Param("id")); err != nil {
		respondErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h Handler) completeTransaction(c *gin.Context) {
	saved, err := h.transactions.Complete(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) cancelTransaction(c *gin.Context) {
	var req dto.FuelCancelRequest
	if !bind(c, &req) {
		return
	}
	saved, err := h.transactions.Cancel(c.Request.Context(), tenantID(c), c.Param("id"), req.Reason, actorID(c))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) adjustTransactionByID(c *gin.Context) {
	var req dto.FuelAdjustmentRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.AdjustmentFromRequest(req)
	item.TransactionID = c.Param("id")
	item.TenantID = tenantID(c)
	item.CreatedBy = actorID(c)
	saved, err := h.transactions.Adjust(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) adjustTransaction(c *gin.Context) {
	var req dto.FuelAdjustmentRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.AdjustmentFromRequest(req)
	item.TenantID = tenantID(c)
	item.CreatedBy = actorID(c)
	saved, err := h.transactions.Adjust(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}

func (h Handler) createType(c *gin.Context) {
	var req dto.FuelTypeRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.TypeFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.types.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) updateType(c *gin.Context) {
	var req dto.FuelTypeRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.TypeFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.types.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) searchTypes(c *gin.Context) {
	page, err := h.types.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) getType(c *gin.Context) {
	item, err := h.types.FindByID(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) deleteType(c *gin.Context) {
	if err := h.types.Delete(c.Request.Context(), tenantID(c), c.Param("id")); err != nil {
		respondErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h Handler) createStation(c *gin.Context) {
	var req dto.FuelStationRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.StationFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.stations.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) updateStation(c *gin.Context) {
	var req dto.FuelStationRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.StationFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.stations.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) searchStations(c *gin.Context) {
	page, err := h.stations.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) getStation(c *gin.Context) {
	item, err := h.stations.FindByID(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) deleteStation(c *gin.Context) {
	if err := h.stations.Delete(c.Request.Context(), tenantID(c), c.Param("id")); err != nil {
		respondErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h Handler) createTank(c *gin.Context) {
	var req dto.FuelTankRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.TankFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.tanks.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) updateTank(c *gin.Context) {
	var req dto.FuelTankRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.TankFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.tanks.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) searchTanks(c *gin.Context) {
	page, err := h.tanks.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) deleteTank(c *gin.Context) {
	if err := h.tanks.Delete(c.Request.Context(), tenantID(c), c.Param("id")); err != nil {
		respondErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h Handler) createNozzle(c *gin.Context) {
	var req dto.FuelNozzleRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.NozzleFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.nozzles.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) updateNozzle(c *gin.Context) {
	var req dto.FuelNozzleRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.NozzleFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.nozzles.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) searchNozzles(c *gin.Context) {
	page, err := h.nozzles.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) deleteNozzle(c *gin.Context) {
	if err := h.nozzles.Delete(c.Request.Context(), tenantID(c), c.Param("id")); err != nil {
		respondErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h Handler) recordReading(c *gin.Context) {
	var req dto.FuelReadingRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.ReadingFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.readings.Record(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchReadings(c *gin.Context) {
	page, err := h.readings.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) recordPrice(c *gin.Context) {
	var req dto.FuelPriceRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.PriceFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.prices.Record(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) updatePrice(c *gin.Context) {
	var req dto.FuelPriceRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.PriceFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.prices.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) searchPrices(c *gin.Context) {
	page, err := h.prices.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) deletePrice(c *gin.Context) {
	if err := h.prices.Delete(c.Request.Context(), tenantID(c), c.Param("id")); err != nil {
		respondErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h Handler) createReceipt(c *gin.Context) {
	var req dto.FuelReceiptRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.ReceiptFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.receipts.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) updateReceipt(c *gin.Context) {
	var req dto.FuelReceiptRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.ReceiptFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.receipts.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) searchReceipts(c *gin.Context) {
	page, err := h.receipts.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) deleteReceipt(c *gin.Context) {
	if err := h.receipts.Delete(c.Request.Context(), tenantID(c), c.Param("id")); err != nil {
		respondErr(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h Handler) searchAdjustments(c *gin.Context) {
	page, err := h.adjustments.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}

func query(c *gin.Context) ports.Query {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	return ports.Query{Search: c.Query("search"), Page: page, PageSize: pageSize, SortBy: c.Query("sort_by"), SortOrder: c.Query("sort_order"), Filters: map[string]string{"status": c.Query("status"), "vehicle_reference": c.Query("vehicle_reference")}}
}
func bind(c *gin.Context, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json body"})
		return false
	}
	return true
}
func respond(c *gin.Context, status int, payload any, err error) {
	if err != nil {
		respondErr(c, err)
		return
	}
	c.JSON(status, payload)
}
func respondErr(c *gin.Context, err error) {
	code := http.StatusBadRequest
	if errors.Is(err, application.ErrNotFound) {
		code = http.StatusNotFound
	}
	c.JSON(code, gin.H{"error": err.Error()})
}
func tenantID(c *gin.Context) string {
	value, _ := c.Get(string(contextkeys.TenantID))
	if tenant, ok := value.(string); ok {
		return tenant
	}
	return ""
}
func actorID(c *gin.Context) string {
	value, _ := c.Get(string(contextkeys.ActorID))
	if actor, ok := value.(string); ok {
		return actor
	}
	return ""
}

package http

import (
	"net/http"
	"strconv"

	"github.com/geevheit/intelligence360/backend/api/internal/core/middleware"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/tires/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/tires/ports"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/tires/transport/dto"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/tires/transport/mapper"
	"github.com/geevheit/intelligence360/backend/api/internal/shared/contextkeys"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	tires       ports.TireService
	inspections ports.TireInspectionService
	movements   ports.TireMovementService
	history     ports.TireHistoryService
}

func NewHandler(tires ports.TireService, inspections ports.TireInspectionService, movements ports.TireMovementService, history ports.TireHistoryService) Handler {
	return Handler{tires: tires, inspections: inspections, movements: movements, history: history}
}

func (h Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/tires")
	group.Use(middleware.RequirePermission("tires.tires.manage"))

	group.GET("", h.searchTires)
	group.POST("", h.createTire)
	group.GET("/brands", h.emptyCollection)
	group.GET("/models", h.emptyCollection)
	group.GET("/specifications", h.emptyCollection)
	group.GET("/positions", h.emptyCollection)
	group.GET("/inspections", h.searchInspections)
	group.POST("/inspections", h.registerInspectionTopLevel)
	group.GET("/measurements", h.emptyCollection)
	group.GET("/installations", h.emptyCollection)
	group.GET("/removals", h.emptyCollection)
	group.GET("/movements", h.searchMovements)
	group.GET("/retreads", h.emptyCollection)
	group.GET("/history", h.searchHistory)
	group.GET("/attachments", h.emptyCollection)
	group.GET("/:id", h.getTire)
	group.PUT("/:id", h.updateTire)
	group.DELETE("/:id", h.deleteTire)
	group.GET("/:id/inspections", h.listInspections)
	group.POST("/:id/inspections", h.registerInspection)
	group.PUT("/:id/inspections/:inspection_id", h.updateInspection)
	group.DELETE("/:id/inspections/:inspection_id", h.deleteInspection)
	group.GET("/:id/movements", h.listMovements)
	group.POST("/:id/movements", h.registerMovement)
	group.GET("/:id/history", h.listHistory)
	group.POST("/:id/receive", h.receive)
	group.POST("/:id/install", h.install)
	group.POST("/:id/remove", h.remove)
	group.POST("/:id/rotate", h.rotate)
	group.POST("/:id/recap", h.recap)
	group.POST("/:id/return", h.returnFromRecap)
	group.POST("/:id/dispose", h.dispose)
}

func (h Handler) createTire(c *gin.Context) {
	var req dto.TireRequest
	if !bind(c, &req) {
		return
	}
	tire := mapper.TireFromRequest(req)
	tire.TenantID = tenantID(c)
	tire.CreatedBy = actorID(c)
	tire.UpdatedBy = actorID(c)
	saved, err := h.tires.Create(c.Request.Context(), tire)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) updateTire(c *gin.Context) {
	var req dto.TireRequest
	if !bind(c, &req) {
		return
	}
	tire := mapper.TireFromRequest(req)
	tire.ID = c.Param("id")
	tire.TenantID = tenantID(c)
	tire.UpdatedBy = actorID(c)
	saved, err := h.tires.Update(c.Request.Context(), tire)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteTire(c *gin.Context) {
	err := h.tires.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h Handler) getTire(c *gin.Context) {
	tire, err := h.tires.FindByID(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, tire, err)
}
func (h Handler) searchTires(c *gin.Context) {
	page, err := h.tires.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}

func (h Handler) registerInspection(c *gin.Context) {
	var req dto.TireInspectionRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.InspectionFromRequest(req)
	item.TenantID = tenantID(c)
	item.TireID = c.Param("id")
	saved, err := h.inspections.Register(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) registerInspectionTopLevel(c *gin.Context) {
	var req dto.TireInspectionRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.InspectionFromRequest(req)
	item.TenantID = tenantID(c)
	item.TireID = req.TireID
	saved, err := h.inspections.Register(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) updateInspection(c *gin.Context) {
	var req dto.TireInspectionRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.InspectionFromRequest(req)
	item.ID = c.Param("inspection_id")
	item.TenantID = tenantID(c)
	item.TireID = c.Param("id")
	saved, err := h.inspections.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteInspection(c *gin.Context) {
	err := h.inspections.Delete(c.Request.Context(), tenantID(c), c.Param("id"), c.Param("inspection_id"))
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h Handler) listInspections(c *gin.Context) {
	page, err := h.inspections.List(c.Request.Context(), tenantID(c), c.Param("id"), query(c))
	respond(c, http.StatusOK, page, err)
}

func (h Handler) registerMovement(c *gin.Context) {
	var req dto.TireMovementRequest
	if !bind(c, &req) {
		return
	}
	movement := mapper.MovementFromRequest(req)
	movement.TenantID = tenantID(c)
	movement.TireID = c.Param("id")
	movement.PerformedBy = actorID(c)
	saved, err := h.movements.Register(c.Request.Context(), movement)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) listMovements(c *gin.Context) {
	page, err := h.movements.List(c.Request.Context(), tenantID(c), c.Param("id"), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) searchMovements(c *gin.Context) {
	page, err := h.movements.List(c.Request.Context(), tenantID(c), c.Query("tire_id"), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) listHistory(c *gin.Context) {
	page, err := h.history.List(c.Request.Context(), tenantID(c), c.Param("id"), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) searchHistory(c *gin.Context) {
	page, err := h.history.List(c.Request.Context(), tenantID(c), c.Query("tire_id"), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) searchInspections(c *gin.Context) {
	page, err := h.inspections.List(c.Request.Context(), tenantID(c), c.Query("tire_id"), query(c))
	respond(c, http.StatusOK, page, err)
}

func (h Handler) receive(c *gin.Context) {
	var req dto.TireOperationRequest
	_ = c.ShouldBindJSON(&req)
	saved, err := h.tires.Receive(c.Request.Context(), tenantID(c), c.Param("id"), req.Reason, actorID(c))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) install(c *gin.Context) {
	var req dto.TireOperationRequest
	if !bind(c, &req) {
		return
	}
	saved, err := h.tires.Install(c.Request.Context(), tenantID(c), c.Param("id"), req.VehicleID, req.Position, req.KM, actorID(c))
	respond(c, http.StatusOK, saved, err)
}

func (h Handler) emptyCollection(c *gin.Context) {
	respond(c, http.StatusOK, ports.Page[any]{Items: []any{}, Page: 1, PageSize: 20, TotalItems: 0, TotalPages: 0}, nil)
}
func (h Handler) remove(c *gin.Context) {
	var req dto.TireOperationRequest
	if !bind(c, &req) {
		return
	}
	saved, err := h.tires.Remove(c.Request.Context(), tenantID(c), c.Param("id"), req.KM, req.Reason, actorID(c))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) rotate(c *gin.Context) {
	var req dto.TireOperationRequest
	if !bind(c, &req) {
		return
	}
	saved, err := h.tires.Rotate(c.Request.Context(), tenantID(c), c.Param("id"), req.Position, req.KM, actorID(c))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) recap(c *gin.Context) {
	var req dto.TireOperationRequest
	if !bind(c, &req) {
		return
	}
	saved, err := h.tires.SendToRecap(c.Request.Context(), tenantID(c), c.Param("id"), req.Reason, actorID(c))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) returnFromRecap(c *gin.Context) {
	var req dto.TireOperationRequest
	_ = c.ShouldBindJSON(&req)
	saved, err := h.tires.ReturnFromRecap(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) dispose(c *gin.Context) {
	var req dto.TireOperationRequest
	if !bind(c, &req) {
		return
	}
	saved, err := h.tires.Dispose(c.Request.Context(), tenantID(c), c.Param("id"), req.Reason, actorID(c))
	respond(c, http.StatusOK, saved, err)
}

func query(c *gin.Context) ports.Query {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	return ports.Query{Search: c.Query("search"), Page: page, PageSize: pageSize, SortBy: c.Query("sort_by"), SortOrder: c.Query("sort_order"), Filters: map[string]string{"vehicle_id": c.Query("vehicle_id"), "status": c.Query("status")}}
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, payload)
}
func tenantID(c *gin.Context) string {
	value, _ := c.Get(string(contextkeys.TenantID))
	tenantID, _ := value.(string)
	return tenantID
}
func actorID(c *gin.Context) string {
	value, _ := c.Get(string(contextkeys.ActorID))
	actorID, _ := value.(string)
	return actorID
}

var _ = domain.Tire{}

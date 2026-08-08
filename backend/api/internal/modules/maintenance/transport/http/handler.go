package http

import (
	"net/http"
	"strconv"

	"github.com/geevheit/intelligence360/backend/api/internal/core/middleware"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/ports"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/transport/dto"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/maintenance/transport/mapper"
	"github.com/geevheit/intelligence360/backend/api/internal/shared/contextkeys"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	workOrders   ports.WorkOrderService
	plans        ports.PreventivePlanService
	serviceTypes ports.ServiceTypeService
	categories   ports.CatalogService
	types        ports.CatalogService
	priorities   ports.CatalogService
	reasons      ports.CatalogService
	workshops    ports.WorkshopService
	technicians  ports.TechnicianService
	labor        ports.LaborService
	downtime     ports.DowntimeService
	history      ports.HistoryService
}

func NewHandler(workOrders ports.WorkOrderService, plans ports.PreventivePlanService, serviceTypes ports.ServiceTypeService, categories ports.CatalogService, types ports.CatalogService, priorities ports.CatalogService, reasons ports.CatalogService, workshops ports.WorkshopService, technicians ports.TechnicianService, labor ports.LaborService, downtime ports.DowntimeService, history ports.HistoryService) Handler {
	return Handler{workOrders: workOrders, plans: plans, serviceTypes: serviceTypes, categories: categories, types: types, priorities: priorities, reasons: reasons, workshops: workshops, technicians: technicians, labor: labor, downtime: downtime, history: history}
}

func (h Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/maintenance")
	group.Use(middleware.RequirePermission("maintenance.maintenance.manage"))
	group.GET("/work-orders", h.searchWorkOrders)
	group.POST("/work-orders", h.createWorkOrder)
	group.GET("/work-orders/:id", h.getWorkOrder)
	group.PUT("/work-orders/:id", h.updateWorkOrder)
	group.DELETE("/work-orders/:id", h.deleteWorkOrder)
	group.POST("/work-orders/:id/start", h.startWorkOrder)
	group.POST("/work-orders/:id/complete", h.completeWorkOrder)
	group.POST("/work-orders/:id/cancel", h.cancelWorkOrder)
	group.GET("/work-orders/:id/labor", h.listLabor)
	group.POST("/work-orders/:id/labor", h.addLabor)
	group.DELETE("/work-orders/:id/labor/:labor_id", h.deleteLabor)
	group.GET("/work-orders/:id/downtime", h.listDowntime)
	group.POST("/work-orders/:id/downtime", h.startDowntime)
	group.POST("/work-orders/:id/downtime/:downtime_id/end", h.endDowntime)
	group.GET("/work-orders/:id/history", h.listHistory)
	group.GET("/preventive-plans", h.searchPlans)
	group.POST("/preventive-plans", h.createPlan)
	group.PUT("/preventive-plans/:id", h.updatePlan)
	group.DELETE("/preventive-plans/:id", h.deletePlan)
	group.GET("/service-types", h.searchServiceTypes)
	group.POST("/service-types", h.createServiceType)
	group.PUT("/service-types/:id", h.updateServiceType)
	group.DELETE("/service-types/:id", h.deleteServiceType)
	group.GET("/categories", h.searchCategories)
	group.POST("/categories", h.createCategory)
	group.PUT("/categories/:id", h.updateCategory)
	group.DELETE("/categories/:id", h.deleteCategory)
	group.GET("/types", h.searchTypes)
	group.POST("/types", h.createType)
	group.PUT("/types/:id", h.updateType)
	group.DELETE("/types/:id", h.deleteType)
	group.GET("/priorities", h.searchPriorities)
	group.POST("/priorities", h.createPriority)
	group.PUT("/priorities/:id", h.updatePriority)
	group.DELETE("/priorities/:id", h.deletePriority)
	group.GET("/reasons", h.searchReasons)
	group.POST("/reasons", h.createReason)
	group.PUT("/reasons/:id", h.updateReason)
	group.DELETE("/reasons/:id", h.deleteReason)
	group.GET("/workshops", h.searchWorkshops)
	group.POST("/workshops", h.createWorkshop)
	group.PUT("/workshops/:id", h.updateWorkshop)
	group.DELETE("/workshops/:id", h.deleteWorkshop)
	group.GET("/technicians", h.searchTechnicians)
	group.POST("/technicians", h.createTechnician)
	group.PUT("/technicians/:id", h.updateTechnician)
	group.DELETE("/technicians/:id", h.deleteTechnician)
}

func (h Handler) createWorkOrder(c *gin.Context) {
	var req dto.WorkOrderRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.WorkOrderFromRequest(req)
	item.TenantID = tenantID(c)
	item.CreatedBy = actorID(c)
	item.UpdatedBy = actorID(c)
	saved, err := h.workOrders.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchWorkOrders(c *gin.Context) {
	page, err := h.workOrders.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) getWorkOrder(c *gin.Context) {
	item, err := h.workOrders.FindByID(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) updateWorkOrder(c *gin.Context) {
	var req dto.WorkOrderRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.WorkOrderFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	item.UpdatedBy = actorID(c)
	saved, err := h.workOrders.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteWorkOrder(c *gin.Context) {
	err := h.workOrders.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h Handler) startWorkOrder(c *gin.Context) {
	item, err := h.workOrders.Start(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) completeWorkOrder(c *gin.Context) {
	item, err := h.workOrders.Complete(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) cancelWorkOrder(c *gin.Context) {
	item, err := h.workOrders.Cancel(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, item, err)
}

func (h Handler) createPlan(c *gin.Context) {
	var req dto.PreventivePlanRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.PreventivePlanFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.plans.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchPlans(c *gin.Context) {
	page, err := h.plans.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updatePlan(c *gin.Context) {
	var req dto.PreventivePlanRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.PreventivePlanFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.plans.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deletePlan(c *gin.Context) {
	_ = h.plans.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createServiceType(c *gin.Context) {
	var req dto.ServiceTypeRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.ServiceTypeFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.serviceTypes.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchServiceTypes(c *gin.Context) {
	page, err := h.serviceTypes.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateServiceType(c *gin.Context) {
	var req dto.ServiceTypeRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.ServiceTypeFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.serviceTypes.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteServiceType(c *gin.Context) {
	_ = h.serviceTypes.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createCategory(c *gin.Context)   { h.createCatalog(c, h.categories) }
func (h Handler) searchCategories(c *gin.Context) { h.searchCatalog(c, h.categories) }
func (h Handler) updateCategory(c *gin.Context)   { h.updateCatalog(c, h.categories) }
func (h Handler) deleteCategory(c *gin.Context) {
	_ = h.categories.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}
func (h Handler) createType(c *gin.Context)  { h.createCatalog(c, h.types) }
func (h Handler) searchTypes(c *gin.Context) { h.searchCatalog(c, h.types) }
func (h Handler) updateType(c *gin.Context)  { h.updateCatalog(c, h.types) }
func (h Handler) deleteType(c *gin.Context) {
	_ = h.types.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}
func (h Handler) createPriority(c *gin.Context)   { h.createCatalog(c, h.priorities) }
func (h Handler) searchPriorities(c *gin.Context) { h.searchCatalog(c, h.priorities) }
func (h Handler) updatePriority(c *gin.Context)   { h.updateCatalog(c, h.priorities) }
func (h Handler) deletePriority(c *gin.Context) {
	_ = h.priorities.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}
func (h Handler) createReason(c *gin.Context)  { h.createCatalog(c, h.reasons) }
func (h Handler) searchReasons(c *gin.Context) { h.searchCatalog(c, h.reasons) }
func (h Handler) updateReason(c *gin.Context)  { h.updateCatalog(c, h.reasons) }
func (h Handler) deleteReason(c *gin.Context) {
	_ = h.reasons.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createCatalog(c *gin.Context, service ports.CatalogService) {
	var req dto.CatalogRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.CatalogFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := service.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}

func (h Handler) searchCatalog(c *gin.Context, service ports.CatalogService) {
	page, err := service.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}

func (h Handler) updateCatalog(c *gin.Context, service ports.CatalogService) {
	var req dto.CatalogRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.CatalogFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := service.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}

func (h Handler) createWorkshop(c *gin.Context) {
	var req dto.WorkshopRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.WorkshopFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.workshops.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchWorkshops(c *gin.Context) {
	page, err := h.workshops.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateWorkshop(c *gin.Context) {
	var req dto.WorkshopRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.WorkshopFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.workshops.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteWorkshop(c *gin.Context) {
	_ = h.workshops.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createTechnician(c *gin.Context) {
	var req dto.TechnicianRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.TechnicianFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.technicians.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchTechnicians(c *gin.Context) {
	page, err := h.technicians.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateTechnician(c *gin.Context) {
	var req dto.TechnicianRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.TechnicianFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.technicians.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteTechnician(c *gin.Context) {
	_ = h.technicians.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) addLabor(c *gin.Context) {
	var req dto.LaborRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.LaborFromRequest(req)
	item.TenantID = tenantID(c)
	item.WorkOrderID = c.Param("id")
	saved, err := h.labor.Add(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) listLabor(c *gin.Context) {
	page, err := h.labor.List(c.Request.Context(), tenantID(c), c.Param("id"), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) deleteLabor(c *gin.Context) {
	_ = h.labor.Delete(c.Request.Context(), tenantID(c), c.Param("id"), c.Param("labor_id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) startDowntime(c *gin.Context) {
	var req dto.DowntimeRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.DowntimeFromRequest(req)
	item.TenantID = tenantID(c)
	item.WorkOrderID = c.Param("id")
	saved, err := h.downtime.Start(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) endDowntime(c *gin.Context) {
	item, err := h.downtime.End(c.Request.Context(), tenantID(c), c.Param("id"), c.Param("downtime_id"))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) listDowntime(c *gin.Context) {
	page, err := h.downtime.List(c.Request.Context(), tenantID(c), c.Param("id"), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) listHistory(c *gin.Context) {
	page, err := h.history.List(c.Request.Context(), tenantID(c), c.Param("id"), query(c))
	respond(c, http.StatusOK, page, err)
}

func query(c *gin.Context) ports.Query {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	return ports.Query{Search: c.Query("search"), Page: page, PageSize: pageSize, SortBy: c.Query("sort_by"), SortOrder: c.Query("sort_order"), Filters: map[string]string{"status": c.Query("status"), "kind": c.Query("kind"), "vehicle_id": c.Query("vehicle_id"), "asset_id": c.Query("asset_id")}}
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

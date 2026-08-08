package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/geevheit/intelligence360/backend/api/internal/core/middleware"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/ports"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/transport/dto"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/inventory/transport/mapper"
	"github.com/geevheit/intelligence360/backend/api/internal/shared/contextkeys"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	parts      ports.PartService
	categories ports.CatalogService
	units      ports.CatalogService
	warehouses ports.WarehouseService
	locations  ports.LocationService
}

func NewHandler(parts ports.PartService, categories ports.CatalogService, units ports.CatalogService, warehouses ports.WarehouseService, locations ports.LocationService) Handler {
	return Handler{parts: parts, categories: categories, units: units, warehouses: warehouses, locations: locations}
}

func (h Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/inventory")
	group.Use(middleware.RequirePermission("inventory.inventory.manage"))
	group.GET("/parts", h.searchParts)
	group.POST("/parts", h.createPart)
	group.GET("/parts/:id", h.getPart)
	group.PUT("/parts/:id", h.updatePart)
	group.DELETE("/parts/:id", h.deletePart)
	group.GET("/categories", h.searchCategories)
	group.POST("/categories", h.createCategory)
	group.PUT("/categories/:id", h.updateCategory)
	group.DELETE("/categories/:id", h.deleteCategory)
	group.GET("/units", h.searchUnits)
	group.POST("/units", h.createUnit)
	group.PUT("/units/:id", h.updateUnit)
	group.DELETE("/units/:id", h.deleteUnit)
	group.GET("/warehouses", h.searchWarehouses)
	group.POST("/warehouses", h.createWarehouse)
	group.PUT("/warehouses/:id", h.updateWarehouse)
	group.DELETE("/warehouses/:id", h.deleteWarehouse)
	group.GET("/locations", h.searchLocations)
	group.POST("/locations", h.createLocation)
	group.PUT("/locations/:id", h.updateLocation)
	group.DELETE("/locations/:id", h.deleteLocation)
}

func (h Handler) createPart(c *gin.Context) {
	var req dto.PartRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.PartFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.parts.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchParts(c *gin.Context) {
	page, err := h.parts.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) getPart(c *gin.Context) {
	item, err := h.parts.FindByID(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) updatePart(c *gin.Context) {
	var req dto.PartRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.PartFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.parts.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deletePart(c *gin.Context) {
	err := h.parts.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h Handler) createCategory(c *gin.Context)   { h.createCatalog(c, h.categories) }
func (h Handler) searchCategories(c *gin.Context) { h.searchCatalog(c, h.categories) }
func (h Handler) updateCategory(c *gin.Context)   { h.updateCatalog(c, h.categories) }
func (h Handler) deleteCategory(c *gin.Context)   { h.deleteCatalog(c, h.categories) }
func (h Handler) createUnit(c *gin.Context)       { h.createCatalog(c, h.units) }
func (h Handler) searchUnits(c *gin.Context)      { h.searchCatalog(c, h.units) }
func (h Handler) updateUnit(c *gin.Context)       { h.updateCatalog(c, h.units) }
func (h Handler) deleteUnit(c *gin.Context)       { h.deleteCatalog(c, h.units) }

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
func (h Handler) deleteCatalog(c *gin.Context, service ports.CatalogService) {
	_ = service.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createWarehouse(c *gin.Context) {
	var req dto.WarehouseRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.WarehouseFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.warehouses.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchWarehouses(c *gin.Context) {
	page, err := h.warehouses.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateWarehouse(c *gin.Context) {
	var req dto.WarehouseRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.WarehouseFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.warehouses.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteWarehouse(c *gin.Context) {
	_ = h.warehouses.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createLocation(c *gin.Context) {
	var req dto.LocationRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.LocationFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.locations.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchLocations(c *gin.Context) {
	page, err := h.locations.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateLocation(c *gin.Context) {
	var req dto.LocationRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.LocationFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.locations.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteLocation(c *gin.Context) {
	_ = h.locations.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func query(c *gin.Context) ports.Query {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	return ports.Query{Search: c.Query("search"), Page: page, PageSize: pageSize, SortBy: c.Query("sort_by"), SortOrder: c.Query("sort_order"), Filters: map[string]string{"status": c.Query("status"), "category_id": c.Query("category_id"), "warehouse_id": c.Query("warehouse_id")}}
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
		code := http.StatusBadRequest
		if errors.Is(err, application.ErrNotFound) {
			code = http.StatusNotFound
		}
		c.JSON(code, gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, payload)
}
func tenantID(c *gin.Context) string {
	value, _ := c.Get(string(contextkeys.TenantID))
	if tenant, ok := value.(string); ok {
		return tenant
	}
	return ""
}

var _ = domain.Part{}

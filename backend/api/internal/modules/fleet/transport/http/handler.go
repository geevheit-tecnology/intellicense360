package http

import (
	"net/http"
	"strconv"

	"github.com/geevheit/intelligence360/backend/api/internal/core/middleware"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/fleet/ports"
	"github.com/geevheit/intelligence360/backend/api/internal/shared/contextkeys"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	vehicles   ports.VehicleService
	brands     ports.VehicleBrandService
	models     ports.VehicleModelService
	categories ports.VehicleCategoryService
	types      ports.VehicleTypeService
	assets     ports.AssetService
}

func NewHandler(vehicles ports.VehicleService, brands ports.VehicleBrandService, models ports.VehicleModelService, categories ports.VehicleCategoryService, types ports.VehicleTypeService, assets ports.AssetService) Handler {
	return Handler{vehicles: vehicles, brands: brands, models: models, categories: categories, types: types, assets: assets}
}

func (h Handler) RegisterRoutes(router *gin.RouterGroup) {
	fleet := router.Group("/fleet")
	fleet.Use(middleware.RequirePermission("fleet.vehicles.manage"))

	fleet.POST("/vehicles", h.createVehicle)
	fleet.GET("/vehicles", h.searchVehicles)
	fleet.GET("/vehicles/:id", h.getVehicle)
	fleet.PUT("/vehicles/:id", h.updateVehicle)
	fleet.DELETE("/vehicles/:id", h.deleteVehicle)

	fleet.POST("/brands", h.createBrand)
	fleet.GET("/brands", h.searchBrands)
	fleet.PUT("/brands/:id", h.updateBrand)
	fleet.DELETE("/brands/:id", h.deleteBrand)

	fleet.POST("/models", h.createModel)
	fleet.GET("/models", h.searchModels)
	fleet.PUT("/models/:id", h.updateModel)
	fleet.DELETE("/models/:id", h.deleteModel)

	fleet.POST("/categories", h.createCategory)
	fleet.GET("/categories", h.searchCategories)
	fleet.PUT("/categories/:id", h.updateCategory)
	fleet.DELETE("/categories/:id", h.deleteCategory)

	fleet.POST("/types", h.createType)
	fleet.GET("/types", h.searchTypes)
	fleet.PUT("/types/:id", h.updateType)
	fleet.DELETE("/types/:id", h.deleteType)

	fleet.POST("/assets", h.createAsset)
	fleet.GET("/assets", h.searchAssets)
	fleet.PUT("/assets/:id", h.updateAsset)
	fleet.DELETE("/assets/:id", h.deleteAsset)
}

func (h Handler) createVehicle(c *gin.Context) {
	var item domain.Vehicle
	if !bind(c, &item) {
		return
	}
	item.TenantID = tenantID(c)
	saved, err := h.vehicles.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchVehicles(c *gin.Context) {
	page, err := h.vehicles.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) getVehicle(c *gin.Context) {
	item, err := h.vehicles.FindByID(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) updateVehicle(c *gin.Context) {
	var item domain.Vehicle
	if !bind(c, &item) {
		return
	}
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.vehicles.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteVehicle(c *gin.Context) {
	err := h.vehicles.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	if err != nil {
		respond(c, http.StatusNotFound, nil, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h Handler) createBrand(c *gin.Context) {
	var item domain.VehicleBrand
	if !bind(c, &item) {
		return
	}
	item.TenantID = tenantID(c)
	saved, err := h.brands.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchBrands(c *gin.Context) {
	page, err := h.brands.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateBrand(c *gin.Context) {
	var item domain.VehicleBrand
	if !bind(c, &item) {
		return
	}
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.brands.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteBrand(c *gin.Context) {
	_ = h.brands.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createModel(c *gin.Context) {
	var item domain.VehicleModel
	if !bind(c, &item) {
		return
	}
	item.TenantID = tenantID(c)
	saved, err := h.models.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchModels(c *gin.Context) {
	page, err := h.models.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateModel(c *gin.Context) {
	var item domain.VehicleModel
	if !bind(c, &item) {
		return
	}
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.models.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteModel(c *gin.Context) {
	_ = h.models.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createCategory(c *gin.Context) {
	var item domain.VehicleCategoryEntity
	if !bind(c, &item) {
		return
	}
	item.TenantID = tenantID(c)
	saved, err := h.categories.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchCategories(c *gin.Context) {
	page, err := h.categories.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateCategory(c *gin.Context) {
	var item domain.VehicleCategoryEntity
	if !bind(c, &item) {
		return
	}
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.categories.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteCategory(c *gin.Context) {
	_ = h.categories.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createType(c *gin.Context) {
	var item domain.VehicleTypeEntity
	if !bind(c, &item) {
		return
	}
	item.TenantID = tenantID(c)
	saved, err := h.types.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchTypes(c *gin.Context) {
	page, err := h.types.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateType(c *gin.Context) {
	var item domain.VehicleTypeEntity
	if !bind(c, &item) {
		return
	}
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.types.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteType(c *gin.Context) {
	_ = h.types.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createAsset(c *gin.Context) {
	var item domain.Asset
	if !bind(c, &item) {
		return
	}
	item.TenantID = tenantID(c)
	saved, err := h.assets.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchAssets(c *gin.Context) {
	page, err := h.assets.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateAsset(c *gin.Context) {
	var item domain.Asset
	if !bind(c, &item) {
		return
	}
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.assets.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteAsset(c *gin.Context) {
	_ = h.assets.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func query(c *gin.Context) ports.Query {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	return ports.Query{Search: c.Query("search"), Page: page, PageSize: pageSize, SortBy: c.Query("sort_by"), SortOrder: c.Query("sort_order"), Filters: map[string]string{"status": c.Query("status"), "category_id": c.Query("category_id")}}
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

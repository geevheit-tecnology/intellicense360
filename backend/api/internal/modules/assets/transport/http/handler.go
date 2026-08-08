package http

import (
	"net/http"
	"strconv"

	"github.com/geevheit/intelligence360/backend/api/internal/core/middleware"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/assets/ports"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/assets/transport/dto"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/assets/transport/mapper"
	"github.com/geevheit/intelligence360/backend/api/internal/shared/contextkeys"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	assets        ports.AssetService
	categories    ports.CategoryService
	types         ports.TypeService
	manufacturers ports.ManufacturerService
	models        ports.ModelService
	equipment     ports.EquipmentService
}

func NewHandler(assets ports.AssetService, categories ports.CategoryService, types ports.TypeService, manufacturers ports.ManufacturerService, models ports.ModelService, equipment ports.EquipmentService) Handler {
	return Handler{assets: assets, categories: categories, types: types, manufacturers: manufacturers, models: models, equipment: equipment}
}

func (h Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/assets")
	group.Use(middleware.RequirePermission("assets.assets.manage"))
	group.GET("", h.searchAssets)
	group.POST("", h.createAsset)
	group.GET("/categories", h.searchCategories)
	group.POST("/categories", h.createCategory)
	group.PUT("/categories/:id", h.updateCategory)
	group.DELETE("/categories/:id", h.deleteCategory)
	group.GET("/types", h.searchTypes)
	group.POST("/types", h.createType)
	group.PUT("/types/:id", h.updateType)
	group.DELETE("/types/:id", h.deleteType)
	group.GET("/manufacturers", h.searchManufacturers)
	group.POST("/manufacturers", h.createManufacturer)
	group.PUT("/manufacturers/:id", h.updateManufacturer)
	group.DELETE("/manufacturers/:id", h.deleteManufacturer)
	group.GET("/models", h.searchModels)
	group.POST("/models", h.createModel)
	group.PUT("/models/:id", h.updateModel)
	group.DELETE("/models/:id", h.deleteModel)
	group.GET("/equipment", h.searchEquipment)
	group.POST("/equipment", h.createEquipment)
	group.PUT("/equipment/:id", h.updateEquipment)
	group.DELETE("/equipment/:id", h.deleteEquipment)
	group.GET("/:id", h.getAsset)
	group.PUT("/:id", h.updateAsset)
	group.DELETE("/:id", h.deleteAsset)
}

func (h Handler) createAsset(c *gin.Context) {
	var req dto.AssetRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.AssetFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.assets.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchAssets(c *gin.Context) {
	page, err := h.assets.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) getAsset(c *gin.Context) {
	item, err := h.assets.FindByID(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) updateAsset(c *gin.Context) {
	var req dto.AssetRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.AssetFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.assets.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteAsset(c *gin.Context) {
	err := h.assets.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h Handler) createCategory(c *gin.Context) {
	var req dto.CategoryRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.CategoryFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.categories.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchCategories(c *gin.Context) {
	page, err := h.categories.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateCategory(c *gin.Context) {
	var req dto.CategoryRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.CategoryFromRequest(req)
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
	var req dto.TypeRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.TypeFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.types.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchTypes(c *gin.Context) {
	page, err := h.types.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateType(c *gin.Context) {
	var req dto.TypeRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.TypeFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.types.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteType(c *gin.Context) {
	_ = h.types.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createManufacturer(c *gin.Context) {
	var req dto.ManufacturerRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.ManufacturerFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.manufacturers.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchManufacturers(c *gin.Context) {
	page, err := h.manufacturers.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateManufacturer(c *gin.Context) {
	var req dto.ManufacturerRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.ManufacturerFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.manufacturers.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteManufacturer(c *gin.Context) {
	_ = h.manufacturers.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createModel(c *gin.Context) {
	var req dto.ModelRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.ModelFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.models.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchModels(c *gin.Context) {
	page, err := h.models.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateModel(c *gin.Context) {
	var req dto.ModelRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.ModelFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.models.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteModel(c *gin.Context) {
	_ = h.models.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func (h Handler) createEquipment(c *gin.Context) {
	var req dto.EquipmentRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.EquipmentFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.equipment.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchEquipment(c *gin.Context) {
	page, err := h.equipment.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) updateEquipment(c *gin.Context) {
	var req dto.EquipmentRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.EquipmentFromRequest(req)
	item.ID = c.Param("id")
	item.TenantID = tenantID(c)
	saved, err := h.equipment.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) deleteEquipment(c *gin.Context) {
	_ = h.equipment.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	c.Status(http.StatusNoContent)
}

func query(c *gin.Context) ports.Query {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	return ports.Query{Search: c.Query("search"), Page: page, PageSize: pageSize, SortBy: c.Query("sort_by"), SortOrder: c.Query("sort_order"), Filters: map[string]string{"status": c.Query("status"), "category_id": c.Query("category_id"), "type_id": c.Query("type_id")}}
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

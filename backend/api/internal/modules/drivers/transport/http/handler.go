package http

import (
	"errors"
	"github.com/geevheit/intelligence360/backend/api/internal/core/middleware"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/drivers/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/drivers/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/drivers/ports"
	"github.com/geevheit/intelligence360/backend/api/internal/shared/contextkeys"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handler struct{ service ports.ManagementService }

func NewHandler(service ports.ManagementService) *Handler { return &Handler{service: service} }
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/drivers")
	group.GET("", middleware.RequirePermission("drivers.read"), h.search)
	group.GET("/:id", middleware.RequirePermission("drivers.read"), h.get)
	group.POST("", middleware.RequirePermission("drivers.manage"), h.create)
	group.PUT("/:id", middleware.RequirePermission("drivers.manage"), h.update)
	group.DELETE("/:id", middleware.RequirePermission("drivers.manage"), h.delete)
}
func tenant(c *gin.Context) string {
	id, _ := contextkeys.TenantIDFromContext(c.Request.Context())
	return id
}
func bind(c *gin.Context, item any) bool {
	if err := c.ShouldBindJSON(item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return false
	}
	return true
}
func respond(c *gin.Context, status int, item any, err error) {
	if err != nil {
		message := "internal server error"
		status = http.StatusInternalServerError
		switch {
		case errors.Is(err, application.ErrNotFound), errors.Is(err, application.ErrReferenceNotFound):
			status = http.StatusNotFound
			message = err.Error()
		case errors.Is(err, application.ErrConflict):
			status = http.StatusConflict
			message = err.Error()
		case errors.Is(err, application.ErrValidation):
			status = http.StatusBadRequest
			message = err.Error()
		}
		c.JSON(status, gin.H{"error": message})
		return
	}
	if status == http.StatusNoContent {
		c.Status(status)
		return
	}
	c.JSON(status, item)
}
func (h *Handler) create(c *gin.Context) {
	var item domain.Driver
	if !bind(c, &item) {
		return
	}
	item.TenantID = tenant(c)
	saved, err := h.service.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}
func (h *Handler) get(c *gin.Context) {
	item, err := h.service.FindByID(c.Request.Context(), tenant(c), c.Param("id"))
	respond(c, http.StatusOK, item, err)
}
func (h *Handler) update(c *gin.Context) {
	var item domain.Driver
	if !bind(c, &item) {
		return
	}
	item.ID = c.Param("id")
	item.TenantID = tenant(c)
	saved, err := h.service.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}
func (h *Handler) delete(c *gin.Context) {
	version, err := strconv.ParseInt(c.Query("version"), 10, 64)
	if err != nil || version < 1 {
		respond(c, 0, nil, application.ErrValidation)
		return
	}
	err = h.service.Delete(c.Request.Context(), tenant(c), c.Param("id"), version)
	respond(c, http.StatusNoContent, nil, err)
}
func (h *Handler) search(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 || page > 1000000 {
		respond(c, 0, nil, application.ErrValidation)
		return
	}
	size, err := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if err != nil || size < 1 || size > 100 {
		respond(c, 0, nil, application.ErrValidation)
		return
	}
	q := ports.Query{Search: c.Query("search"), Status: c.Query("status"), VehicleID: c.Query("vehicle_id"), Page: page, PageSize: size}
	result, err := h.service.Search(c.Request.Context(), tenant(c), q)
	respond(c, http.StatusOK, result, err)
}

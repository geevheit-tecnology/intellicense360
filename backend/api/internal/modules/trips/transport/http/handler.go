package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/geevheit/intelligence360/backend/api/internal/core/middleware"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/trips/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/trips/ports"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/trips/transport/dto"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/trips/transport/mapper"
	"github.com/geevheit/intelligence360/backend/api/internal/shared/contextkeys"
)

type Handler struct {
	trips ports.TripService
}

func NewHandler(trips ports.TripService) Handler {
	return Handler{trips: trips}
}

func (h Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/trips")
	group.Use(middleware.RequirePermission("trips.read"))
	group.GET("", h.searchTrips)
	group.GET("/summary", h.getSummary)
	group.GET("/:id", h.getTrip)

	manageGroup := group.Group("")
	manageGroup.Use(middleware.RequirePermission("trips.manage"))
	manageGroup.POST("", h.createTrip)
	manageGroup.POST("/:id/dispatch", h.dispatchTrip)
	manageGroup.POST("/:id/deliver", h.deliverTrip)
	manageGroup.POST("/:id/cancel", h.cancelTrip)
}

func (h Handler) createTrip(c *gin.Context) {
	var req dto.TripRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.TripFromRequest(req)
	item.TenantID = tenantID(c)
	saved, err := h.trips.Create(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}

func (h Handler) searchTrips(c *gin.Context) {
	page, err := h.trips.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}

func (h Handler) getTrip(c *gin.Context) {
	item, err := h.trips.FindByID(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, item, err)
}

func (h Handler) dispatchTrip(c *gin.Context) {
	var req dto.DispatchRequest
	_ = c.ShouldBindJSON(&req)
	saved, err := h.trips.DispatchTrip(c.Request.Context(), tenantID(c), c.Param("id"), req.StartOdometerKM)
	respond(c, http.StatusOK, saved, err)
}

func (h Handler) deliverTrip(c *gin.Context) {
	var req dto.DeliverRequest
	if !bind(c, &req) {
		return
	}
	saved, err := h.trips.DeliverTrip(c.Request.Context(), tenantID(c), c.Param("id"), req.EndOdometerKM)
	respond(c, http.StatusOK, saved, err)
}

func (h Handler) cancelTrip(c *gin.Context) {
	var req dto.CancelRequest
	_ = c.ShouldBindJSON(&req)
	saved, err := h.trips.CancelTrip(c.Request.Context(), tenantID(c), c.Param("id"), req.Reason)
	respond(c, http.StatusOK, saved, err)
}

func (h Handler) getSummary(c *gin.Context) {
	summary, err := h.trips.GetSummary(c.Request.Context(), tenantID(c))
	respond(c, http.StatusOK, summary, err)
}

func query(c *gin.Context) ports.Query {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	return ports.Query{
		Search:    c.Query("search"),
		Page:      page,
		PageSize:  size,
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
		Filters: map[string]string{
			"status": c.Query("status"),
		},
	}
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
	v, _ := c.Get(string(contextkeys.TenantID))
	s, _ := v.(string)
	return s
}


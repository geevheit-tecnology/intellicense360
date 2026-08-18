package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/geevheit/intelligence360/backend/api/internal/core/middleware"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/missioncontrol/ports"
	"github.com/geevheit/intelligence360/backend/api/internal/shared/contextkeys"
	"github.com/gin-gonic/gin"
)

type Handler struct{ commands ports.CommandService }

func NewHandler(commands ports.CommandService) Handler { return Handler{commands: commands} }

func (h Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/mission-control")
	group.Use(middleware.RequirePermission("mission_control.read"))
	group.GET("/summary", h.summary)
	group.GET("/items", h.listItems)
	group.GET("/items/:id", h.getItem)
	group.POST("/items", middleware.RequirePermission("mission_control.create"), h.createItem)
	group.PATCH("/items/:id", middleware.RequirePermission("mission_control.update"), h.updateItem)
	group.POST("/items/:id/acknowledge", middleware.RequirePermission("mission_control.acknowledge"), h.acknowledge)
	group.POST("/items/:id/start", middleware.RequirePermission("mission_control.update"), h.start)
	group.POST("/items/:id/resolve", middleware.RequirePermission("mission_control.resolve"), h.resolve)
	group.POST("/items/:id/dismiss", middleware.RequirePermission("mission_control.dismiss"), h.dismiss)
	group.GET("/items/:id/actions", h.listActions)
	group.POST("/items/:id/actions", middleware.RequirePermission("mission_control.update"), h.createAction)
	group.GET("/items/:id/history", h.history)
	group.GET("/snapshot", h.latestSnapshot)
	group.POST("/snapshot/rebuild", middleware.RequirePermission("mission_control.snapshot"), h.rebuildSnapshot)
	group.GET("/recommendations", h.listRecommendations)
	group.POST("/recommendations/evaluate", middleware.RequirePermission("mission_control.update"), h.evaluateRecommendations)
}

func (h Handler) summary(c *gin.Context) {
	summary, err := h.commands.Summary(c.Request.Context(), tenantID(c))
	respond(c, http.StatusOK, summary, err)
}
func (h Handler) createItem(c *gin.Context) {
	var req domain.CommandItem
	if !bind(c, &req) {
		return
	}
	req.TenantID = tenantID(c)
	saved, err := h.commands.Create(c.Request.Context(), req, actorID(c))
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) updateItem(c *gin.Context) {
	var req domain.CommandItem
	if !bind(c, &req) {
		return
	}
	req.ID, req.TenantID = c.Param("id"), tenantID(c)
	saved, err := h.commands.Update(c.Request.Context(), req, actorID(c))
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) getItem(c *gin.Context) {
	item, err := h.commands.Get(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) listItems(c *gin.Context) {
	page, err := h.commands.List(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) acknowledge(c *gin.Context) {
	item, err := h.commands.Acknowledge(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) start(c *gin.Context) {
	item, err := h.commands.Start(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) resolve(c *gin.Context) {
	item, err := h.commands.Resolve(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) dismiss(c *gin.Context) {
	item, err := h.commands.Dismiss(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) createAction(c *gin.Context) {
	var req domain.CommandAction
	if !bind(c, &req) {
		return
	}
	req.TenantID, req.CommandItemID = tenantID(c), c.Param("id")
	saved, err := h.commands.CreateAction(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) listActions(c *gin.Context) {
	page, err := h.commands.ListActions(c.Request.Context(), tenantID(c), c.Param("id"), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) history(c *gin.Context) {
	page, err := h.commands.History(c.Request.Context(), tenantID(c), c.Param("id"), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) latestSnapshot(c *gin.Context) {
	snapshot, err := h.commands.LatestSnapshot(c.Request.Context(), tenantID(c))
	if errors.Is(err, application.ErrCommandItemNotFound) {
		snapshot, err = h.commands.RebuildSnapshot(c.Request.Context(), tenantID(c))
	}
	respond(c, http.StatusOK, snapshot, err)
}
func (h Handler) rebuildSnapshot(c *gin.Context) {
	snapshot, err := h.commands.RebuildSnapshot(c.Request.Context(), tenantID(c))
	respond(c, http.StatusCreated, snapshot, err)
}
func (h Handler) evaluateRecommendations(c *gin.Context) {
	actions, err := h.commands.EvaluateRecommendations(c.Request.Context(), tenantID(c))
	respond(c, http.StatusCreated, actions, err)
}
func (h Handler) listRecommendations(c *gin.Context) {
	q := query(c)
	q.Filters["type"] = string(domain.TypeRecommendation)
	page, err := h.commands.List(c.Request.Context(), tenantID(c), q)
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
		code, name := http.StatusInternalServerError, "internal_error"
		switch {
		case errors.Is(err, application.ErrCommandItemNotFound):
			code, name = http.StatusNotFound, "not_found"
		case errors.Is(err, application.ErrDuplicateCommandItem):
			code, name = http.StatusConflict, "duplicate_command_item"
		case errors.Is(err, application.ErrInvalidStatusTransition):
			code, name = http.StatusUnprocessableEntity, "invalid_status_transition"
		case errors.Is(err, application.ErrInvalidCommandItemType), errors.Is(err, application.ErrInvalidSeverity), errors.Is(err, application.ErrInvalidPriority), errors.Is(err, application.ErrInvalidRiskScore), errors.Is(err, application.ErrInvalidImpactScore), errors.Is(err, application.ErrInvalidConfidence), errors.Is(err, application.ErrInvalidData):
			code, name = http.StatusBadRequest, "invalid_data"
		}
		c.JSON(code, gin.H{"error": name, "message": err.Error()})
		return
	}
	c.JSON(status, payload)
}
func tenantID(c *gin.Context) string {
	v, _ := c.Get(string(contextkeys.TenantID))
	s, _ := v.(string)
	return s
}
func actorID(c *gin.Context) string {
	v, _ := c.Get(string(contextkeys.ActorID))
	s, _ := v.(string)
	return s
}
func query(c *gin.Context) ports.Query {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "25"))
	return ports.Query{Search: c.Query("search"), Page: page, PerPage: perPage, Sort: c.Query("sort"), Filters: map[string]string{
		"type": c.Query("type"), "category": c.Query("category"), "severity": c.Query("severity"), "priority": c.Query("priority"), "status": c.Query("status"), "source_type": c.Query("source_type"), "assigned_to": c.Query("assigned_to"), "sla_status": c.Query("sla_status"),
	}}
}

package http

import (
	"net/http"
	"strconv"

	"github.com/geevheit/intelligence360/backend/api/internal/core/middleware"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/ports"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/transport/dto"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/transport/mapper"
	"github.com/geevheit/intelligence360/backend/api/internal/shared/contextkeys"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	checklists ports.ChecklistService
	items      ports.ChecklistItemService
	answers    ports.ChecklistAnswerService
}

func NewHandler(checklists ports.ChecklistService, items ports.ChecklistItemService, answers ports.ChecklistAnswerService) Handler {
	return Handler{checklists: checklists, items: items, answers: answers}
}

func (h Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/checklists")
	group.Use(middleware.RequirePermission("checklist.checklists.manage"))

	group.GET("", h.searchChecklists)
	group.POST("", h.createChecklist)
	group.GET("/:id", h.getChecklist)
	group.PUT("/:id", h.updateChecklist)
	group.DELETE("/:id", h.deleteChecklist)
	group.POST("/:id/start", h.startChecklist)
	group.POST("/:id/finish", h.finishChecklist)
	group.POST("/:id/cancel", h.cancelChecklist)

	group.GET("/:id/items", h.listItems)
	group.POST("/:id/items", h.addItem)
	group.PUT("/:id/items/:item_id", h.updateItem)
	group.DELETE("/:id/items/:item_id", h.deleteItem)

	group.POST("/:id/answers", h.answerItem)
	group.GET("/:id/answers", h.listAnswers)
}

func (h Handler) createChecklist(c *gin.Context) {
	var req dto.ChecklistRequest
	if !bind(c, &req) {
		return
	}
	checklist := mapper.ChecklistFromRequest(req)
	checklist.TenantID = tenantID(c)
	checklist.CreatedBy = actorID(c)
	checklist.UpdatedBy = actorID(c)
	saved, err := h.checklists.Create(c.Request.Context(), checklist)
	respond(c, http.StatusCreated, saved, err)
}

func (h Handler) updateChecklist(c *gin.Context) {
	var req dto.ChecklistRequest
	if !bind(c, &req) {
		return
	}
	checklist := mapper.ChecklistFromRequest(req)
	checklist.ID = c.Param("id")
	checklist.TenantID = tenantID(c)
	checklist.UpdatedBy = actorID(c)
	saved, err := h.checklists.Update(c.Request.Context(), checklist)
	respond(c, http.StatusOK, saved, err)
}

func (h Handler) deleteChecklist(c *gin.Context) {
	err := h.checklists.Delete(c.Request.Context(), tenantID(c), c.Param("id"))
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h Handler) getChecklist(c *gin.Context) {
	checklist, err := h.checklists.FindByID(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, checklist, err)
}

func (h Handler) searchChecklists(c *gin.Context) {
	page, err := h.checklists.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}

func (h Handler) startChecklist(c *gin.Context) {
	checklist, err := h.checklists.Start(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, checklist, err)
}

func (h Handler) finishChecklist(c *gin.Context) {
	checklist, err := h.checklists.Finish(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, checklist, err)
}

func (h Handler) cancelChecklist(c *gin.Context) {
	checklist, err := h.checklists.Cancel(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, checklist, err)
}

func (h Handler) addItem(c *gin.Context) {
	var req dto.ChecklistItemRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.ItemFromRequest(req)
	item.TenantID = tenantID(c)
	item.ChecklistID = c.Param("id")
	saved, err := h.items.Add(c.Request.Context(), item)
	respond(c, http.StatusCreated, saved, err)
}

func (h Handler) updateItem(c *gin.Context) {
	var req dto.ChecklistItemRequest
	if !bind(c, &req) {
		return
	}
	item := mapper.ItemFromRequest(req)
	item.ID = c.Param("item_id")
	item.TenantID = tenantID(c)
	item.ChecklistID = c.Param("id")
	saved, err := h.items.Update(c.Request.Context(), item)
	respond(c, http.StatusOK, saved, err)
}

func (h Handler) deleteItem(c *gin.Context) {
	err := h.items.Delete(c.Request.Context(), tenantID(c), c.Param("id"), c.Param("item_id"))
	if err != nil {
		respond(c, http.StatusBadRequest, nil, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h Handler) listItems(c *gin.Context) {
	page, err := h.items.List(c.Request.Context(), tenantID(c), c.Param("id"), query(c))
	respond(c, http.StatusOK, page, err)
}

func (h Handler) answerItem(c *gin.Context) {
	var req dto.ChecklistAnswerRequest
	if !bind(c, &req) {
		return
	}
	answer := mapper.AnswerFromRequest(req)
	answer.TenantID = tenantID(c)
	answer.ChecklistID = c.Param("id")
	answer.AnsweredBy = actorID(c)
	saved, err := h.answers.AnswerItem(c.Request.Context(), answer)
	respond(c, http.StatusCreated, saved, err)
}

func (h Handler) listAnswers(c *gin.Context) {
	page, err := h.answers.List(c.Request.Context(), tenantID(c), c.Param("id"), query(c))
	respond(c, http.StatusOK, page, err)
}

func query(c *gin.Context) ports.Query {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	return ports.Query{
		Search:    c.Query("search"),
		Page:      page,
		PageSize:  pageSize,
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
		Filters: map[string]string{
			"vehicle_id":        c.Query("vehicle_id"),
			"status":            c.Query("status"),
			"type":              c.Query("type"),
			"category":          c.Query("category"),
			"checklist_item_id": c.Query("checklist_item_id"),
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

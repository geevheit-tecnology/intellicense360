package http

import (
	"net/http"
	"strconv"

	"github.com/geevheit/intelligence360/backend/api/internal/core/middleware"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/ports"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/transport/dto"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/checklist/transport/mapper"
	"github.com/geevheit/intelligence360/backend/api/internal/shared/contextkeys"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	checklists      ports.ChecklistService
	items           ports.ChecklistItemService
	answers         ports.ChecklistAnswerService
	templates       ports.ChecklistTemplateService
	versions        ports.ChecklistTemplateVersionService
	types           ports.ChecklistTypeService
	sections        ports.ChecklistSectionService
	engineItems     ports.ChecklistEngineItemService
	executions      ports.ChecklistExecutionService
	responses       ports.ChecklistEngineResponseService
	evidence        ports.ChecklistEvidenceService
	nonConformities ports.ChecklistNonConformityService
	history         ports.ChecklistHistoryService
}

func NewHandler(checklists ports.ChecklistService, items ports.ChecklistItemService, answers ports.ChecklistAnswerService, templates ports.ChecklistTemplateService, versions ports.ChecklistTemplateVersionService, types ports.ChecklistTypeService, sections ports.ChecklistSectionService, engineItems ports.ChecklistEngineItemService, executions ports.ChecklistExecutionService, responses ports.ChecklistEngineResponseService, evidence ports.ChecklistEvidenceService, nonConformities ports.ChecklistNonConformityService, history ports.ChecklistHistoryService) Handler {
	return Handler{checklists: checklists, items: items, answers: answers, templates: templates, versions: versions, types: types, sections: sections, engineItems: engineItems, executions: executions, responses: responses, evidence: evidence, nonConformities: nonConformities, history: history}
}

func (h Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/checklists")
	group.Use(middleware.RequirePermission("checklist.checklists.manage"))
	group.Use(middleware.RequirePermission("checklist.checklist.manage"))

	group.GET("/templates", h.searchTemplates)
	group.POST("/templates", h.createTemplate)
	group.GET("/templates/:id", h.getTemplate)
	group.PUT("/templates/:id", h.updateTemplate)
	group.POST("/templates/:id/archive", h.archiveTemplate)
	group.GET("/templates/:id/versions", h.listTemplateVersions)
	group.POST("/templates/:id/versions", h.createTemplateVersion)
	group.POST("/templates/versions/:version_id/publish", h.publishTemplateVersion)
	group.GET("/types", h.searchEngineTypes)
	group.POST("/types", h.createEngineType)
	group.GET("/sections", h.searchSections)
	group.POST("/sections", h.createSection)
	group.GET("/items", h.searchEngineItems)
	group.POST("/items", h.createEngineItem)
	group.GET("/executions", h.searchExecutions)
	group.POST("/executions", h.startExecution)
	group.GET("/executions/:id", h.getExecution)
	group.POST("/executions/:id/complete", h.completeExecution)
	group.POST("/executions/:id/cancel", h.cancelExecution)
	group.POST("/executions/:id/invalidate", h.invalidateExecution)
	group.GET("/executions/:id/responses", h.listEngineResponses)
	group.POST("/executions/:id/responses", h.recordEngineResponse)
	group.GET("/executions/:id/evidence", h.listEvidence)
	group.POST("/executions/:id/evidence", h.addEvidence)
	group.GET("/executions/:id/non-conformities", h.listNonConformities)
	group.POST("/executions/:id/non-conformities", h.createNonConformity)
	group.GET("/executions/:id/history", h.listHistory)

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

func (h Handler) createTemplate(c *gin.Context) {
	var req domain.ChecklistTemplate
	if !bind(c, &req) {
		return
	}
	req.TenantID, req.CreatedBy, req.UpdatedBy = tenantID(c), actorID(c), actorID(c)
	saved, err := h.templates.Create(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) updateTemplate(c *gin.Context) {
	var req domain.ChecklistTemplate
	if !bind(c, &req) {
		return
	}
	req.ID, req.TenantID, req.UpdatedBy = c.Param("id"), tenantID(c), actorID(c)
	saved, err := h.templates.Update(c.Request.Context(), req)
	respond(c, http.StatusOK, saved, err)
}
func (h Handler) searchTemplates(c *gin.Context) {
	page, err := h.templates.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) getTemplate(c *gin.Context) {
	item, err := h.templates.FindByID(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) archiveTemplate(c *gin.Context) {
	item, err := h.templates.Archive(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) createTemplateVersion(c *gin.Context) {
	var req domain.ChecklistTemplateVersion
	if !bind(c, &req) {
		return
	}
	req.TenantID, req.TemplateID = tenantID(c), c.Param("id")
	saved, err := h.versions.Create(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) publishTemplateVersion(c *gin.Context) {
	item, err := h.versions.Publish(c.Request.Context(), tenantID(c), c.Param("version_id"), actorID(c))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) listTemplateVersions(c *gin.Context) {
	page, err := h.versions.ListByTemplate(c.Request.Context(), tenantID(c), c.Param("id"), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) createEngineType(c *gin.Context) {
	var req domain.ChecklistType
	if !bind(c, &req) {
		return
	}
	req.TenantID = tenantID(c)
	saved, err := h.types.Create(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchEngineTypes(c *gin.Context) {
	page, err := h.types.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) createSection(c *gin.Context) {
	var req domain.ChecklistSection
	if !bind(c, &req) {
		return
	}
	req.TenantID = tenantID(c)
	saved, err := h.sections.Create(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchSections(c *gin.Context) {
	page, err := h.sections.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) createEngineItem(c *gin.Context) {
	var req domain.ChecklistEngineItem
	if !bind(c, &req) {
		return
	}
	req.TenantID = tenantID(c)
	saved, err := h.engineItems.Create(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchEngineItems(c *gin.Context) {
	page, err := h.engineItems.ListByVersion(c.Request.Context(), tenantID(c), c.Query("template_version_id"), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) startExecution(c *gin.Context) {
	var req domain.ChecklistExecution
	if !bind(c, &req) {
		return
	}
	req.TenantID = tenantID(c)
	if req.PerformedBy == "" {
		req.PerformedBy = actorID(c)
	}
	saved, err := h.executions.Start(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchExecutions(c *gin.Context) {
	page, err := h.executions.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) getExecution(c *gin.Context) {
	item, err := h.executions.FindByID(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) completeExecution(c *gin.Context) {
	item, err := h.executions.Complete(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) cancelExecution(c *gin.Context) {
	item, err := h.executions.Cancel(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) invalidateExecution(c *gin.Context) {
	item, err := h.executions.Invalidate(c.Request.Context(), tenantID(c), c.Param("id"), actorID(c))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) recordEngineResponse(c *gin.Context) {
	var req domain.ChecklistResponse
	if !bind(c, &req) {
		return
	}
	req.TenantID, req.ExecutionID = tenantID(c), c.Param("id")
	if req.Responder == "" {
		req.Responder = actorID(c)
	}
	saved, err := h.responses.Record(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) listEngineResponses(c *gin.Context) {
	page, err := h.responses.ListByExecution(c.Request.Context(), tenantID(c), c.Param("id"), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) addEvidence(c *gin.Context) {
	var req domain.ChecklistEvidence
	if !bind(c, &req) {
		return
	}
	req.TenantID, req.ExecutionID = tenantID(c), c.Param("id")
	saved, err := h.evidence.Add(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) listEvidence(c *gin.Context) {
	page, err := h.evidence.ListByExecution(c.Request.Context(), tenantID(c), c.Param("id"), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) createNonConformity(c *gin.Context) {
	var req domain.ChecklistNonConformity
	if !bind(c, &req) {
		return
	}
	req.TenantID, req.ExecutionID = tenantID(c), c.Param("id")
	saved, err := h.nonConformities.Create(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) listNonConformities(c *gin.Context) {
	page, err := h.nonConformities.ListByExecution(c.Request.Context(), tenantID(c), c.Param("id"), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) listHistory(c *gin.Context) {
	page, err := h.history.ListByExecution(c.Request.Context(), tenantID(c), c.Param("id"), query(c))
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
			"vehicle_id":          c.Query("vehicle_id"),
			"status":              c.Query("status"),
			"type":                c.Query("type"),
			"category":            c.Query("category"),
			"checklist_item_id":   c.Query("checklist_item_id"),
			"template_version_id": c.Query("template_version_id"),
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

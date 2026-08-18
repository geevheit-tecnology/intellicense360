package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/geevheit/intelligence360/backend/api/internal/core/middleware"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/application"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/domain"
	"github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/ports"
	"github.com/geevheit/intelligence360/backend/api/internal/shared/contextkeys"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	metrics      ports.MetricService
	intelligence ports.IntelligenceService
}

func NewHandler(metrics ports.MetricService, intelligence ports.IntelligenceService) Handler {
	return Handler{metrics: metrics, intelligence: intelligence}
}

func (h Handler) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/intelligence")
	group.Use(middleware.RequirePermission("intelligence.intelligence.manage"))
	group.GET("/health", h.health)
	group.GET("/metrics", h.searchMetrics)
	group.POST("/metrics", h.createMetric)
	group.GET("/metrics/:id", h.getMetric)
	group.GET("/anomalies", h.searchAnomalies)
	group.POST("/anomalies", h.createAnomaly)
	group.GET("/risks", h.searchRisks)
	group.POST("/risks", h.createRisk)
	group.GET("/opportunities", h.searchOpportunities)
	group.POST("/opportunities", h.createOpportunity)
	group.GET("/recommendations", h.searchRecommendations)
	group.POST("/recommendations", h.createRecommendation)
	group.GET("/insights", h.searchInsights)
	group.POST("/insights", h.createInsight)
	group.GET("/insights/:id", h.getInsight)
	group.POST("/insights/:id/acknowledge", h.acknowledgeInsight)
	group.POST("/insights/:id/resolve", h.resolveInsight)
	group.POST("/insights/:id/dismiss", h.dismissInsight)
}

func (h Handler) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "engine": "deterministic", "external_ai": false})
}
func (h Handler) createMetric(c *gin.Context) {
	var req domain.IntelligenceMetric
	if !bind(c, &req) {
		return
	}
	req.TenantID = tenantID(c)
	saved, err := h.metrics.Create(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchMetrics(c *gin.Context) {
	page, err := h.metrics.Search(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) getMetric(c *gin.Context) {
	item, err := h.metrics.FindByID(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) createAnomaly(c *gin.Context) {
	var req domain.Anomaly
	if !bind(c, &req) {
		return
	}
	req.TenantID = tenantID(c)
	saved, err := h.intelligence.CreateAnomaly(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchAnomalies(c *gin.Context) {
	page, err := h.intelligence.SearchAnomalies(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) createRisk(c *gin.Context) {
	var req domain.Risk
	if !bind(c, &req) {
		return
	}
	req.TenantID = tenantID(c)
	saved, err := h.intelligence.CreateRisk(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchRisks(c *gin.Context) {
	page, err := h.intelligence.SearchRisks(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) createOpportunity(c *gin.Context) {
	var req domain.Opportunity
	if !bind(c, &req) {
		return
	}
	req.TenantID = tenantID(c)
	saved, err := h.intelligence.CreateOpportunity(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchOpportunities(c *gin.Context) {
	page, err := h.intelligence.SearchOpportunities(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) createRecommendation(c *gin.Context) {
	var req domain.Recommendation
	if !bind(c, &req) {
		return
	}
	req.TenantID = tenantID(c)
	saved, err := h.intelligence.CreateRecommendation(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchRecommendations(c *gin.Context) {
	page, err := h.intelligence.SearchRecommendations(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) createInsight(c *gin.Context) {
	var req domain.Insight
	if !bind(c, &req) {
		return
	}
	req.TenantID = tenantID(c)
	saved, err := h.intelligence.CreateInsight(c.Request.Context(), req)
	respond(c, http.StatusCreated, saved, err)
}
func (h Handler) searchInsights(c *gin.Context) {
	page, err := h.intelligence.SearchInsights(c.Request.Context(), tenantID(c), query(c))
	respond(c, http.StatusOK, page, err)
}
func (h Handler) getInsight(c *gin.Context) {
	item, err := h.intelligence.FindInsight(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) acknowledgeInsight(c *gin.Context) {
	item, err := h.intelligence.AcknowledgeInsight(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) resolveInsight(c *gin.Context) {
	item, err := h.intelligence.ResolveInsight(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, item, err)
}
func (h Handler) dismissInsight(c *gin.Context) {
	item, err := h.intelligence.DismissInsight(c.Request.Context(), tenantID(c), c.Param("id"))
	respond(c, http.StatusOK, item, err)
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
		case errors.Is(err, application.ErrNotFound):
			code, name = http.StatusNotFound, "not_found"
		case errors.Is(err, application.ErrInvalidData), errors.Is(err, application.ErrInsufficientData):
			code, name = http.StatusBadRequest, "invalid_data"
		case errors.Is(err, application.ErrDuplicateInsight):
			code, name = http.StatusConflict, "duplicate_insight"
		}
		c.JSON(code, gin.H{"error": name, "message": err.Error()})
		return
	}
	c.JSON(status, payload)
}
func tenantID(c *gin.Context) string {
	value, _ := c.Get(string(contextkeys.TenantID))
	tenantID, _ := value.(string)
	return tenantID
}
func query(c *gin.Context) ports.Query {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "25"))
	return ports.Query{Search: c.Query("search"), Page: page, PerPage: perPage, Sort: c.Query("sort"), Filters: map[string]string{"status": c.Query("status"), "severity": c.Query("severity"), "category": c.Query("category")}}
}

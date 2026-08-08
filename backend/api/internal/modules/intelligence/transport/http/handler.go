package http

import (
	stdhttp "net/http"

	"github.com/geevheit/intelligence360/backend/api/internal/modules/intelligence/ports"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	recommendationService ports.RecommendationService
}

func NewHandler(recommendationService ports.RecommendationService) Handler {
	return Handler{recommendationService: recommendationService}
}

func (h Handler) RegisterRoutes(router gin.IRoutes) {
	router.GET("/mission-control/summary", h.currentOperationalSummary)
}

func (h Handler) currentOperationalSummary(c *gin.Context) {
	summary := h.recommendationService.CurrentOperationalSummary(c.Request.Context())
	c.JSON(stdhttp.StatusOK, summary)
}
